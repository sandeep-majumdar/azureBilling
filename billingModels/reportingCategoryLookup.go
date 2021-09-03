package billingModels

import (
	"encoding/csv"
	"fmt"
	"io"
	"os"
	"strings"

	"github.com/adeturner/observability"
)

func (rcl *reportingCategoryLookup) print(cnt int) {

	i := 0

	for k, v := range rcl.Items {
		observability.Logger("Info", fmt.Sprintf("%s -> %v\n", k, v))
		i++
		if i > cnt {
			break
		}
	}
}

func (rcl *reportingCategoryLookup) printCount() {
	observability.Logger("Info", fmt.Sprintf("ReportingCategoryLookup has %d records\n", len(rcl.Items)))
}

func (rcl *reportingCategoryLookup) init() {
	rcl.Items = make(map[string]ReportingCategoryLookupItem)
}

func (rcl *reportingCategoryLookup) Read(fileLocation string) error {

	rcl.init()

	var key string

	f, err := os.Open(fileLocation)
	if err != nil {
		observability.Logger("Error", fmt.Sprintf("Unable to read input file=%s err=%s", fileLocation, err))
	}
	defer f.Close()

	cnt := 0

	if err == nil {

		r := csv.NewReader(f)
		for {

			record, err := r.Read()
			if err == io.EOF {
				break
			} else if err != nil {
				observability.Logger("Error", fmt.Sprintf("Unable to parse file as CSV; file=%s err=%s", fileLocation, err))
				break
			}

			cnt++

			// skip the first row (header)
			if cnt > 1 {
				i := ReportingCategoryLookupItem{}
				i.setValues(record)

				key = rcl.getKey(i.MeterCategory, i.MeterSubCategory)
				rcl.Items[key] = i
			}
		}
	}

	observability.LogMemory("Info")
	rcl.printCount()
	// rcl.print(10)

	return err

}

func (rcl *reportingCategoryLookup) getKey(meterCategory, meterSubCategory string) string {
	return strings.ToLower(fmt.Sprintf(":%s:%s:", meterCategory, meterSubCategory))
}

func (rcl *reportingCategoryLookup) Get(meterCategory, meterSubCategory string) (ReportingCategoryLookupItem, bool) {

	key := rcl.getKey(meterCategory, meterSubCategory)

	rcli, ok := rcl.Items[key]

	if !ok {

		// try again with wildcard
		key := rcl.getKey(meterCategory, "*")

		rcli, ok = rcl.Items[key]

		if !ok {
			observability.Logger("Error", fmt.Sprintf("Unable to find ReportingCategoryLookupItem for key=%s", key))
		}

	}

	return rcli, ok
}
