package azureBilling

import (
	"fmt"
	"io"

	"github.com/adeturner/observability"
)

func (aggp *aggregatePlatform) print(cnt int) {

	i := 0

	for k, v := range aggp.items {
		observability.Logger("Info", fmt.Sprintf("%s -> %v\n", k, v))
		i++
		if i > cnt {
			break
		}
	}
}

func (aggp *aggregatePlatform) printCount() {
	observability.Logger("Info", fmt.Sprintf("managedDiskLookup has %d records\n", len(aggp.items)))
}

func (aggp *aggregatePlatform) init() {
	aggp.items = make(map[string]*aggregatePlatformItem)
}

func (aggp *aggregatePlatform) add(reportingCategory, reportingSubCategory, platform string, unitOfMeasure string, quantity float64, costInBillingCurrency float64) {

	key := fmt.Sprintf("%s:%s:%s:%s", reportingCategory, reportingSubCategory, platform, unitOfMeasure)

	// initializes two variables - api will receive either the value of "key" from the map
	// or a "zero value" (in this case the empty string)
	// ok will receive a bool that will be set to true if "key" was actually present in the map
	// evaluates ok, which will be true if "key" was in the map
	if _, ok := aggp.items[key]; !ok {

		// if not found initialise
		api := aggregatePlatformItem{}
		api.platform = platform
		api.reportingCategory = reportingCategory
		api.reportingSubCategory = reportingSubCategory
		api.unitOfMeasure = unitOfMeasure
		api.CostInBillingCurrency = 0
		api.Quantity = 0
		aggp.items[key] = &api
	}

	aggp.items[key].Quantity += quantity
	aggp.items[key].CostInBillingCurrency += costInBillingCurrency
}

/*
   #######################################################
   Below here is about producing output only
   #######################################################
*/

func (aggp *aggregatePlatform) check(e error) {
	if e != nil {
		observability.Logger("Error", fmt.Sprintf("%v", e))
		panic(e)
	}
}

func (aggp *aggregatePlatform) getCSVHeader() []byte {
	return []byte(fmt.Sprintf("\"%s\",\"%s\",\"%s\",\"%s\",\"%s\",\"%s\"\n",
		"reportingCategory", "reportingSubCategory", "platform", "UnitOfMeasure", "Quantity", "CostInBillingCurrency"))

}

func (aggp *aggregatePlatform) WriteCSVHeader(w io.Writer) {

	_, err := w.Write(aggp.getCSVHeader())
	aggp.check(err)
}

func (aggp *aggregatePlatform) WriteCSVOutput(w io.Writer) {

	var csvRow string

	for _, v := range aggp.items {

		csvRow = fmt.Sprintf("\"%s\",\"%s\",\"%s\",\"%s\",\"%f\",\"%f\"\n",
			v.reportingCategory, v.reportingSubCategory, v.platform, v.unitOfMeasure, v.Quantity, v.CostInBillingCurrency)

		_, err := w.Write([]byte(csvRow))
		aggp.check(err)
	}

}

func (aggp *aggregatePlatform) WriteFile(filename string) {

	observability.Logger("Info", fmt.Sprintf("Writing to %s", filename))

	var fs fileSystem = localFS{}

	file, err := fs.Create(filename)
	if err != nil {
		observability.Logger("Error", fmt.Sprintf("Failed to open file: %v", err))
	}
	defer file.Close()

	aggp.WriteCSVHeader(file)
	aggp.WriteCSVOutput(file)

}
