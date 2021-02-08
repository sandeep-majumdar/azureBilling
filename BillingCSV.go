package azureBilling

import (
	"encoding/csv"
	"fmt"
	"io"
	"os"

	"github.com/adeturner/observability"
)

func (bcsv *BillingCSV) SetFile(filePath string) {
	bcsv.fileLocation = filePath
}

func (bcsv *BillingCSV) ProcessFile() error {

	f, err := os.Open(bcsv.fileLocation)
	if err != nil {
		observability.Logger("Error", fmt.Sprintf("Unable to read input file=%s err=%s", bcsv.fileLocation, err))
	}
	defer f.Close()

	cnt := 0

	AggregateResourceGroup.init()

	var uom string
	var cat, subcat string
	var plat, portfolio, product, envType string
	var summaryCategory, quantityDivisor string
	var divisor float64

	if err == nil {

		r := csv.NewReader(f)

		t1 := observability.Timer{}
		t1.Start(true, "BillingCSV")
		for {

			record, err := r.Read()
			if err == io.EOF {
				break
			} else if err != nil {
				observability.Logger("Error", fmt.Sprintf("Unable to parse file as CSV; file=%s err=%s", bcsv.fileLocation, err))
				break
			}

			cnt++

			// skip the header row
			if cnt > 1 {

				l := billingLine{}
				l.setValues(record)

				rcli, ok1 := ReportingCategoryLookup.get(l.MeterCategory)
				if ok1 {
					cat = rcli.reportingCategory
					subcat = rcli.reportingSubCategory
				} else {
					cat = l.MeterCategory
					subcat = ""
				}

				plmi, ok2 := PlatformMapLookup.get(l.SubscriptionId, l.ResourceGroup)
				if ok2 {
					portfolio = plmi.portfolio
					plat = plmi.platform
					product = plmi.productCode
					envType = plmi.environmentType
				} else {
					portfolio = "Other"
					plat = "Other"
					envType = "Other"
					product = "Other"
				}

				pmi, ok3 := MeterLookup.get(l.MeterId)
				if ok3 {
					uom = pmi.UnitOfMeasure
				} else {
					uom = l.UnitOfMeasure
				}

				scli, ok4 := SummaryCategoryLookup.get(cat, subcat, uom)
				if ok4 {
					summaryCategory = scli.Summary
					quantityDivisor = scli.QuantityDivisor
					divisor = SummaryCategoryLookup.getDivisor(quantityDivisor, l.BillingPeriodEndDate)
				} else {
					summaryCategory = "Other"
					summaryCategory = "Other"
					divisor = 1.0
				}

				quantity := l.Quantity
				summaryQuantity := l.Quantity / divisor

				// adjust quantity for managed disks
				// note its not perfect, because selecting performance option for a
				// small disk will allocate a larger disk without the volume
				if l.MeterCategory[len(l.MeterCategory)-5:] == " Disks" {
					mdli, ok4 := ManagedDiskLookup.get(l.MeterName)
					if ok4 {
						summaryQuantity = float64(mdli.SizeGB) * l.Quantity / divisor
					}
				}

				// For databases, we will ignore DTUs etc and just count the line items
				if cat == "Data PaaS" && subcat == "Database" && summaryCategory == "ResourceUnits" {
					summaryQuantity = 1.0 / divisor
				}

				AggregateResourceGroup.add(cat, subcat, portfolio, plat, product, envType, uom, summaryCategory, quantityDivisor, summaryQuantity, quantity, l)

				/*
					When we see a vm being used, calculate Cores and Memory from the Count
					Compute IaaS,Virtual Machines,1 Hour,Count,NumDaysInMonthTimes24Hrs
					Compute IaaS,Virtual Machines,n/a,CPU,n/a
					Compute IaaS,Virtual Machines,n/a,Memory (GB),n/a
				*/

				if cat == "Compute IaaS" && subcat == "Virtual Machines" && summaryCategory == "ResourceUnits" {

					if ok3 {
						// observability.Logger("Info", fmt.Sprintf("sku=%s armsku=%s", pmi.SkuName, pmi.ArmSkuName))
						// sku=DS12-2 v2 armsku=Standard_DS12-2_v2
						vmli, ok5 := VmSizeLookup.get(pmi.ArmSkuName)

						if ok5 {

							// observability.Logger("Info", fmt.Sprintf("Matched ArmSkuName=%s", pmi.ArmSkuName))
							cores := l.Quantity * float64(vmli.Cores) / divisor
							memgb := l.Quantity * float64(vmli.MemGB) / divisor

							// set quantity = 0 because these dont exist in the source csv
							AggregateResourceGroup.add(cat, subcat, portfolio, plat, product, envType, "CPU", "CPU", quantityDivisor, cores, 0, l)
							AggregateResourceGroup.add(cat, subcat, portfolio, plat, product, envType, "MemGB", "MemGB", quantityDivisor, memgb, 0, l)

						} else {
							// logging happens in VmSizeLookup
						}

					} else {
						observability.Logger("Info", fmt.Sprintf("Failed meterlookup, cannot retrieve sku for vm"))
					}

				}

				if mod(cnt, 100000) == 0 {
					observability.Logger("Info", fmt.Sprintf("Processed %d rows of billing CSV", cnt))
					observability.LogMemory("Info")
				}
			}
		}

		outputAggregateRGCsvFile := ConfigMap.WorkingDirectory + ConfigMap.OutputAggregateRGCsvFile

		AggregateResourceGroup.WriteFile(outputAggregateRGCsvFile)

		t1.EndAndPrint(true)

		observability.Logger("Info", fmt.Sprintf("Complete. Processed %d rows of billing CSV", cnt))
		observability.LogMemory("Info")

	}

	return err

}
