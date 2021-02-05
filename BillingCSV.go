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
				plmi, ok2 := PlatformMapLookup.get(l.SubscriptionId, l.ResourceGroup)
				pmi, ok3 := MeterLookup.get(l.MeterId)

				quantity := l.Quantity

				// adjust quantity for managed disks
				// note its not perfect, because selecting performance option for a
				// small disk will allocate a larger disk without the volume
				if l.MeterCategory[len(l.MeterCategory)-5:] == " Disks" {
					mdli, ok4 := ManagedDiskLookup.get(l.MeterName)
					if ok4 {
						quantity = float64(mdli.SizeGB)
					}
				}

				if ok1 && ok2 && ok3 {

					AggregateTotal.add(rcli.reportingCategory, rcli.reportingSubCategory, pmi.UnitOfMeasure, quantity, l.CostInBillingCurrency)
					AggregatePlatform.add(l.SubscriptionId, l.ResourceGroup, plmi.platform, quantity, l.CostInBillingCurrency)

				}

				if mod(cnt, 100000) == 0 {
					observability.Logger("Info", fmt.Sprintf("Processed %d rows of billing CSV", cnt))
				}
			}
		}

		AggregateTotal.print(1000)

		t1.EndAndPrint(true)

		observability.Logger("Info", fmt.Sprintf("Complete. Processed %d rows of billing CSV", cnt))
	}

	return err

}
