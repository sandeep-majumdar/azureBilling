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

	//AggregateTotal.init()
	//AggregatePlatform.init()
	AggregateResourceGroup.init()

	var uom string
	var cat, subcat string
	var plat, portfolio string
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
				} else {
					portfolio = "Other"
					plat = "Other"
				}

				pmi, ok3 := MeterLookup.get(l.MeterId)
				if ok3 {
					uom = pmi.UnitOfMeasure
				} else {
					uom = l.UnitOfMeasure
				}

				scli, ok2 := SummaryCategoryLookup.get(cat, subcat, uom)
				if ok2 {
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

				AggregateResourceGroup.add(cat, subcat, portfolio, plat, uom, summaryCategory, quantityDivisor, summaryQuantity, quantity, l)

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
