package rightsizing

import (
	"encoding/csv"
	"encoding/json"
	"fmt"
	"io"
	"os"
	"strings"
	"sync"

	"github.com/adeturner/azureBilling/azure"
	"github.com/adeturner/azureBilling/billing"
	"github.com/adeturner/azureBilling/billingModels"
	"github.com/adeturner/azureBilling/config"
	"github.com/adeturner/azureBilling/observability"
	"github.com/adeturner/azureBilling/utils"
)

type VmRightsizing struct {
	fileLocation string
	vmd          *vmDetails
	vdv          *vmDayValues
	vmm          *vmMonitorMetrics
}

func (rsz *VmRightsizing) SetFile(filePath string) {
	rsz.fileLocation = filePath
}

func (rsz *VmRightsizing) ProcessBillFile() (err error) {

	var f *os.File
	var plat, portfolio, product, envType string
	loadFromBillfile := true

	if err == nil {
		rsz.vmd, err = NewVmDetails()
	}

	if err == nil {
		rsz.vdv, err = NewVmDayValues()
	}

	if rsz.vmd.FileExists(rsz.getOutputVmDetailsCSVFile()) && rsz.vdv.FileExists(rsz.getOutputVmDayValuesCSVFile()) {
		observability.Info("Loading from output files")
		loadFromBillfile = false
		err = rsz.vmd.ReadFile(rsz.getOutputVmDetailsCSVFile())
		if err == nil {
			rsz.vdv.ReadFile(rsz.getOutputVmDayValuesCSVFile())
		}
	}

	if err == nil && loadFromBillfile {
		f, err = os.Open(rsz.fileLocation)
		if err != nil {
			observability.Logger("Error", fmt.Sprintf("Unable to read input file=%s err=%s", rsz.fileLocation, err))
		}
		defer f.Close()
	}

	cnt := 0

	if err == nil && loadFromBillfile {

		r := csv.NewReader(f)
		t1 := observability.Timer{}
		t1.Start(true, "BillingCSV")
		for {

			record, err := r.Read()
			if err == io.EOF {
				break
			} else if err != nil {
				observability.Logger("Error", fmt.Sprintf("Unable to parse file as CSV; file=%s err=%s", rsz.fileLocation, err))
				break
			}

			cnt++

			// skip the header row
			if cnt > 1 {

				l := billing.BillingLine{}
				l.SetValues(record)

				if l.MeterCategory == "Virtual Machines" {

					plmi, ok2 := billingModels.PlatformMapLookup.Get(l.SubscriptionId, l.ResourceGroup)

					if ok2 {
						portfolio = plmi.Portfolio
						plat = plmi.Platform
						product = plmi.ProductCode
						envType = plmi.EnvironmentType
					} else {
						portfolio = "Other"
						plat = "Other"
						envType = "Other"
						product = "Other"
					}

					if strings.Contains(l.Tags, "\"Vendor\": \"Databricks\"") || l.ConsumedService != "Microsoft.Compute" {
						// Ignore unless its Microsoft.Compute and not Databricks
					} else {
						err = rsz.vmd.Add(l.ResourceId, portfolio, plat, product, envType, l.ResourceLocation, l.MeterName)
						if err == nil {
							err = rsz.vdv.Add(l.ResourceId, l.Date, l.EffectivePrice, l.UnitPrice, l.Quantity, l.CostInBillingCurrency, l.UnitOfMeasure)
						}
					}
				}

				if utils.Mod(cnt, 100000) == 0 {
					//observability.Info(fmt.Sprintf("%s %s %s %f %v", l.Date, l.SubscriptionId, l.ResourceName, l.Quantity, rsz.vdv))
					observability.Info(fmt.Sprintf("Processed %d rows of billing CSV", cnt))
					observability.Info(fmt.Sprintf("%d detail records; %d measure records", len(rsz.vmd.vmMap), len(rsz.vdv.vmMap)))
					observability.LogMemory("Info")
				}
			}
		}

		t1.EndAndPrint(true)
		observability.Logger("Info", fmt.Sprintf("Complete. Processed %d rows of billing CSV", cnt))
		observability.LogMemory("Info")

		if err == nil {
			filename := rsz.getOutputVmDetailsCSVFile()
			observability.Info(fmt.Sprintf("Writing vm details file to %s", filename))
			err = rsz.vmd.WriteFile(filename)
		}

		if err == nil {
			filename := rsz.getOutputVmDayValuesCSVFile()
			observability.Info(fmt.Sprintf("Writing vm day values file to %s", filename))
			err = rsz.vdv.WriteFile(filename)
		}
	}
	return err
}

func (rsz *VmRightsizing) getOutputVmDetailsCSVFile() string {
	return config.ConfigMap.WorkingDirectory + "/" + config.ConfigMap.OutputVmDetailsCSVFile
}

func (rsz *VmRightsizing) getOutputVmDayValuesCSVFile() string {
	return config.ConfigMap.WorkingDirectory + "/" + config.ConfigMap.OutputVmDayValuesCSVFile
}

func (rsz *VmRightsizing) getOutputVmMonitorMetricsFile() string {
	return config.ConfigMap.WorkingDirectory + "/" + config.ConfigMap.OutputVmMonitorMetricsFile
}

func (rsz *VmRightsizing) ProcessMetrics() (err error) {
	observability.Info("Processing Metrics")

	if err == nil {
		rsz.vmm, err = NewVmMonitorMetrics()
	}

	loadFromFile := false

	if rsz.vmm.FileExists(rsz.getOutputVmMonitorMetricsFile()) {
		observability.Info("Loading from output file")
		err = rsz.vmm.ReadFile(rsz.getOutputVmMonitorMetricsFile())
		loadFromFile = true
	}

	var azClis []*azure.AzureCli

	if err == nil && !loadFromFile {
		i := 0
		for i < config.ConfigMap.RightsizingMaxThreads {
			azClis = append(azClis, azure.NewAzureCli())
			observability.Info(fmt.Sprintf("Created azure cli %d : len(azClis)=%d", i, len(azClis)))
			err = azClis[i].Login()
			i++
		}
		observability.Info(fmt.Sprintf("Azure cli creation complete"))
	}

	if err == nil && !loadFromFile {

		cnt := 0
		var wg sync.WaitGroup
		i := 0
		for k, v := range rsz.vmd.vmMap {
			wg.Add(1)

			go func(cliIndex int, key vmdConcatKey, val *vmDetail) {
				rsz.ProcessMetricsForVM(azClis[cliIndex], key, val)
				wg.Done()
			}(i, k, v)

			i++
			cnt++
			if i >= config.ConfigMap.RightsizingMaxThreads {
				wg.Wait()
				i = 0
			}

			if utils.Mod(cnt, 100) == 0 {
				observability.Info(fmt.Sprintf("Processed %d virtual machines", cnt))
				observability.Info(fmt.Sprintf("Have %d vm metric records", rsz.vmm.GetMapLen()))
				observability.LogMemory("Info")
			}
		}

		if err == nil {
			filename := rsz.getOutputVmMonitorMetricsFile()
			observability.Info(fmt.Sprintf("Processed %d virtual machines", cnt))
			observability.Info(fmt.Sprintf("Have %d vm metric records", rsz.vmm.GetMapLen()))
			observability.Info(fmt.Sprintf("Writing vm day values file to %s", filename))
			err = rsz.vmm.WriteFile(filename)
		}
	}

	if err != nil {
		observability.Error(fmt.Sprintf(err.Error()))
	}

	return err
}

func (rsz *VmRightsizing) ProcessMetricsForVM(azCli *azure.AzureCli, k vmdConcatKey, v *vmDetail) (err error) {

	var cmd, out, errStr string
	var mt *azMonitorMetricsType

	//if rsz.vmm.Exists(v.ResourceId) {
	//	observability.Info("Skipping resource")
	//} else {

	azCli.SetSubscription(v.getSubscription())

	// intervals... PT24H PT1H PT15M PT5M PT1M PT60S
	cmd = "az monitor metrics list --resource " +
		v.ResourceId + " --metric \"Percentage CPU\" \"Available Memory Bytes\" " +
		" --start-time 2021-08-01T00:00:00Z --end-time 2021-08-31T00:00:00Z" +
		" --interval PT1H --aggregation Average Count Maximum Minimum"

	out, err = azCli.ExecCmd(cmd)

	if err != nil {
		// Encode known errors, and update the resourceid for tracking
		errStr = rsz.getErrorString(out)
		if errStr == "" {
			// output unknown errors
			observability.Error(fmt.Sprintf("%s, %s %s", v.ResourceId, err.Error(), out))
			errStr = "Unknown"
		}
		rsz.vmm.add(v.ResourceId, errStr, mt)
	} else {
		mt, err = rsz.ProcessMetricOutput(k, v, out)
		if err == nil {
			err = rsz.vmm.add(v.ResourceId, "", mt)
		} else {
			observability.Error(err.Error())
		}

	}

	//}
	return err
}

func (rsz *VmRightsizing) ProcessMetricOutput(k vmdConcatKey, v *vmDetail, out string) (mt *azMonitorMetricsType, err error) {

	mt = &azMonitorMetricsType{}
	err = json.Unmarshal([]byte(out), &mt)
	if err == nil {
		//observability.Info(fmt.Sprintf("%v", mt))
	}
	return mt, err
}

func (rsz *VmRightsizing) getErrorString(out string) (errStr string) {

	if strings.Contains(out, "InvalidAuthenticationTokenTenant") {
		errStr = "InvalidAuthenticationTokenTenant"
	} else if strings.Contains(out, "ResourceNotFound") {
		errStr = "ResourceNotFound"
	} else if strings.Contains(out, "ResourceGroupNotFound") {
		errStr = "ResourceGroupNotFound"
	} else if strings.Contains(out, "AuthorizationFailed") {
		errStr = "AuthorizationFailed"
	} else if strings.Contains(out, "Logging error") {
		/*
			exit status 1 --- Logging error ---
			Traceback (most recent call last):
			  File "/opt/az/lib/python3.6/site-packages/knack/cli.py", line 231, in invoke
			  File "/opt/az/lib/python3.6/site-packages/azure/cli/core/commands/__init__.py", line 582, in execute
			  File "/opt/az/lib/python3.6/site-packages/knack/parser.py",
		*/
		errStr = "Logging error"
	} else if strings.Contains(out, " Auto upgrade failed") {
		/*
			exit status 120 WARNING: Auto upgrade failed. name 'exit_code' is not defined
			Error in atexit._run_exitfuncs:
			Error in atexit._run_exitfuncs:
			Error in atexit._run_exitfuncs:
			Error in atexit._run_exitfuncs:
		*/
		errStr = " Auto upgrade failed"
	}
	return errStr
}
