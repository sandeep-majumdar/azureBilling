package azureBilling

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

	for k, v := range rcl.items {
		observability.Logger("Info", fmt.Sprintf("%s -> %v\n", k, v))
		i++
		if i > cnt {
			break
		}
	}
}

func (rcl *reportingCategoryLookup) printCount() {
	observability.Logger("Info", fmt.Sprintf("reportingCategoryLookup has %d records\n", len(rcl.items)))
}

func (rcl *reportingCategoryLookup) init() {
	rcl.items = make(map[string]reportingCategoryLookupItem)
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
				i := reportingCategoryLookupItem{}
				i.setValues(record)

				key = rcl.getKey(i.meterCategory, i.meterSubCategory)
				rcl.items[key] = i
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

func (rcl *reportingCategoryLookup) get(meterCategory, meterSubCategory string) (reportingCategoryLookupItem, bool) {

	key := rcl.getKey(meterCategory, meterSubCategory)

	rcli, ok := rcl.items[key]

	if !ok {

		// try again with wildcard
		key := rcl.getKey(meterCategory, "*")

		rcli, ok = rcl.items[key]

		if !ok {
			observability.Logger("Error", fmt.Sprintf("Unable to find reportingCategoryLookupItem for key=%s", key))
		}

	}

	return rcli, ok
}
