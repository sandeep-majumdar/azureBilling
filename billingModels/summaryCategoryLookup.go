package billingModels

import (
	"encoding/csv"
	"fmt"
	"io"
	"os"
	"strconv"
	"strings"

	"github.com/adeturner/azureBilling/config"
	"github.com/adeturner/observability"
)

func (scl *summaryCategoryLookup) print(cnt int) {

	i := 0

	for k, v := range scl.Items {
		observability.Logger("Info", fmt.Sprintf("%s -> %v", k, v))
		i++
		if i > cnt {
			break
		}
	}
}

func (scl *summaryCategoryLookup) printRecord(ReportingCategory, ReportingSubCategory, UnitOfMeasure string) {

	observability.Logger("Info", fmt.Sprintf("|%s|%s|%s|", ReportingCategory, ReportingSubCategory, UnitOfMeasure))

	for k, v := range scl.Items {

		if v.ReportingCategory == "Data PaaS" {
			observability.Logger("Info", fmt.Sprintf("M0|%s|%s|%s|", v.ReportingCategory, v.ReportingSubCategory, v.UnitOfMeasure))
		}

		if v.ReportingCategory == ReportingCategory {
			observability.Logger("Info", fmt.Sprintf("M1 %s -> %v", k, v))
		}

		if v.ReportingCategory == ReportingCategory && v.ReportingSubCategory == ReportingSubCategory {
			observability.Logger("Info", fmt.Sprintf("M2 %s -> %v", k, v))
		}

		if v.ReportingCategory == ReportingCategory && v.ReportingSubCategory == ReportingSubCategory && v.UnitOfMeasure == UnitOfMeasure {
			observability.Logger("Info", fmt.Sprintf("M3 %s -> %v", k, v))
		}
	}
}

func (scl *summaryCategoryLookup) printKey(str string) {

	i := 0
	b := false

	for k, v := range scl.Items {
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
	observability.Logger("Info", fmt.Sprintf("summaryCategoryLookup has %d records\n", len(scl.Items)))
}

func (scl *summaryCategoryLookup) init() {
	scl.Items = make(map[string]SummaryCategoryLookupItem)
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
				i := SummaryCategoryLookupItem{}
				i.setValues(record)

				key = scl.getKey(i.ReportingCategory, i.ReportingSubCategory, i.UnitOfMeasure)

				scl.Items[key] = i
			}
		}
	}

	observability.LogMemory("Info")
	scl.printCount()
	// scl.print(300)

	return err

}

// effectivedate = US format: mm/dd/yyyy
func (scl *summaryCategoryLookup) GetDivisor(quantityDivisor, effectiveDate string) float64 {

	var d float64

	n, err := strconv.Atoi(quantityDivisor)
	if err == nil {

		d = float64(n)

	} else {

		days, err := strconv.Atoi(config.ConfigMap.NumDaysInMonth)
		if err != nil {
			observability.Logger("Error", fmt.Sprintf("Failed to parse days, err := strconv.Atoi(.NumDaysInMonth) from %s", config.ConfigMap.NumDaysInMonth))
		}

		// NumDaysInMonthTimes24HrsDiv10hrs

		switch quantityDivisor {
		case "NumDaysInMonthTimes24HrsDiv1024":
			d = float64(days) * 24 * 1024
		case "NumDaysInMonthTimes24Hrs":
			d = float64(days) * 24
		case "NumDaysInMonthTimes24HrsDiv10":
			d = float64(days) * 24.0 / 10.0
		case "NumDaysInMonthTimes24HoursDiv100":
			d = float64(days) * 24.0 / 100.0
		case "NumDaysInMonthTimes24HrsDiv10kmins":
			d = float64(days) * 24.0 * 60.0 / 10000.0
		case "NumDaysInMonthTimes24HrsDiv1000mins":
			d = float64(days) * 24.0 * 60.0 / 1000.0
		case "NumDaysInMonth":
			d = float64(days)
		case "Times10":
			d = 0.1
		case "ManagedDisksOnly":
			d = 1.0
		case "DayTo24Hours": // = 1/24
			d = 0.04
		default: // n/a,
			d = 1.0
			observability.Logger("Error", fmt.Sprintf("Unexpected default, quantityDivisor=%s", quantityDivisor))
		}

	}

	return d
}

func (scl *summaryCategoryLookup) getKey(ReportingCategory, ReportingSubCategory, UnitOfMeasure string) string {
	return strings.ToLower(fmt.Sprintf(":%s:%s:%s:", ReportingCategory, ReportingSubCategory, UnitOfMeasure))
}

func (scl *summaryCategoryLookup) Get(ReportingCategory, ReportingSubCategory, UnitOfMeasure, MeterCategory, MeterSubCategory string) (SummaryCategoryLookupItem, bool) {

	// scl.printCount()
	// scl.printRecord(ReportingCategory, ReportingSubCategory, UnitOfMeasure)

	key1 := scl.getKey(ReportingCategory, ReportingSubCategory, UnitOfMeasure)
	scli, ok := scl.Items[key1]

	if !ok {

		var ok2 bool

		key2 := scl.getKey(ReportingCategory, ReportingSubCategory, "Other")
		scli, ok2 = scl.Items[key2]

		if !ok2 {
			observability.Logger("Error", fmt.Sprintf("Unable to find summaryCategoryLookupItem for key1=%s, key2=%s", key1, key2))
			// fmt.Println(fmt.Sprintf("%s,%s,%s,%s,%s", MeterCategory, MeterSubCategory, ReportingCategory, ReportingSubCategory, UnitOfMeasure))
		}
	}

	return scli, ok
}
