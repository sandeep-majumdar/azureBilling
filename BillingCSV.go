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

	AggregateTotal.init()
	AggregatePlatform.init()
	AggregateResourceGroup.init()

	var uom string
	var cat, subcat string
	var plat string

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
					plat = plmi.platform
				} else {
					plat = "Other"
				}

				pmi, ok3 := MeterLookup.get(l.MeterId)
				if ok3 {
					uom = pmi.UnitOfMeasure
				} else {
					uom = l.UnitOfMeasure + "?"
				}

				quantity := l.Quantity

				// adjust quantity for managed disks
				// note its not perfect, because selecting performance option for a
				// small disk will allocate a larger disk without the volume
				if l.MeterCategory[len(l.MeterCategory)-5:] == " Disks" {
					mdli, ok4 := ManagedDiskLookup.get(l.MeterName)
					if ok4 {
						quantity = float64(mdli.SizeGB) * l.Quantity
					}
				}

				AggregateTotal.add(cat, subcat, uom, quantity, l.CostInBillingCurrency)
				AggregatePlatform.add(cat, subcat, plat, uom, quantity, l.CostInBillingCurrency)
				AggregateResourceGroup.add(cat, subcat, plat, uom, quantity, l)

				if mod(cnt, 100000) == 0 {
					observability.Logger("Info", fmt.Sprintf("Processed %d rows of billing CSV", cnt))
					observability.LogMemory("Info")
				}
			}
		}

		AggregateTotal.print(1000)

		outputAggregateTotalCsvFile := ConfigMap.WorkingDirectory + ConfigMap.OutputAggregateTotalCsvFile
		outputAggregatePlatformCsvFile := ConfigMap.WorkingDirectory + ConfigMap.OutputAggregatePlatformCsvFile
		outputAggregateRGCsvFile := ConfigMap.WorkingDirectory + ConfigMap.OutputAggregateRGCsvFile

		AggregateTotal.WriteFile(outputAggregateTotalCsvFile)
		AggregatePlatform.WriteFile(outputAggregatePlatformCsvFile)
		AggregateResourceGroup.WriteFile(outputAggregateRGCsvFile)

		t1.EndAndPrint(true)

		observability.Logger("Info", fmt.Sprintf("Complete. Processed %d rows of billing CSV", cnt))
		observability.LogMemory("Info")

	}

	return err

}
