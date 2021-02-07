package azureBilling

import (
	"encoding/csv"
	"fmt"
	"io"
	"os"
	"strconv"
	"strings"

	"github.com/adeturner/observability"
)

func (scl *summaryCategoryLookup) print(cnt int) {

	i := 0

	for k, v := range scl.items {
		observability.Logger("Info", fmt.Sprintf("%s -> %v", k, v))
		i++
		if i > cnt {
			break
		}
	}
}

func (scl *summaryCategoryLookup) printKey(str string) {

	i := 0
	b := false

	for k, v := range scl.items {
		if k == str {
			observability.Logger("Info", fmt.Sprintf("PrintKey %s -> %v", k, v))
			b = true
			break
		}

		i++
	}

	if !b {
		observability.Logger("Info", fmt.Sprintf("PrintKey %s -> NOT FOUND", str))
	}
}

func (scl *summaryCategoryLookup) printCount() {
	observability.Logger("Info", fmt.Sprintf("summaryCategoryLookup has %d records\n", len(scl.items)))
}

func (scl *summaryCategoryLookup) init() {
	scl.items = make(map[string]summaryCategoryLookupItem)
}

func (scl *summaryCategoryLookup) Read(fileLocation string) error {

	scl.init()

	f, err := os.Open(fileLocation)
	if err != nil {
		observability.Logger("Error", fmt.Sprintf("Unable to read input file=%s err=%s", fileLocation, err))
	}
	defer f.Close()

	cnt := 0
	var key string

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
				i := summaryCategoryLookupItem{}
				i.setValues(record)

				key = scl.getKey(i.reportingCategory, i.reportingSubCategory, i.UnitOfMeasure)

				scl.items[key] = i
			}
		}
	}

	observability.LogMemory("Info")
	scl.printCount()
	// scl.print(100)

	return err

}

// effectivedate = US format: mm/dd/yyyy
func (scl *summaryCategoryLookup) getDivisor(quantityDivisor, effectiveDate string) float64 {

	var d float64

	n, err := strconv.Atoi(quantityDivisor)
	if err == nil {

		d = float64(n)

	} else {

		days, err := strconv.Atoi(ConfigMap.NumDaysInMonth)
		if err != nil {
			observability.Logger("Error", fmt.Sprintf("Failed to parse days, err := strconv.Atoi(ConfigMap.NumDaysInMonth) from %s", ConfigMap.NumDaysInMonth))
		}

		switch quantityDivisor {
		case "NumDaysInMonthTimes24Hrs":
			d = float64(days) * 24
		case "NumDaysInMonth":
			d = float64(days)
		case "ManagedDisksOnly":
			d = 1.0
		default: // n/a,
			d = 1.0
			observability.Logger("Error", fmt.Sprintf("Unexpected default, quantityDivisor=%s", quantityDivisor))
		}

	}

	return d
}

func (scl *summaryCategoryLookup) getKey(reportingCategory, reportingSubCategory, UnitOfMeasure string) string {
	return strings.ToLower(fmt.Sprintf(":%s:%s:%s:", reportingCategory, reportingSubCategory, UnitOfMeasure))
}

func (scl *summaryCategoryLookup) get(reportingCategory, reportingSubCategory, UnitOfMeasure string) (summaryCategoryLookupItem, bool) {

	key1 := scl.getKey(reportingCategory, reportingSubCategory, UnitOfMeasure)
	scli, ok := scl.items[key1]

	if !ok {

		var ok2 bool

		key2 := scl.getKey(reportingCategory, reportingSubCategory, "Other")
		scli, ok2 = scl.items[key2]

		if !ok2 {
			observability.Logger("Error", fmt.Sprintf("Unable to find summaryCategoryLookupItem for key1=%s, key2=%s", key1, key2))
		}
	}

	return scli, ok
}
