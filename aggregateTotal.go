package azureBilling

import (
	"fmt"

	"github.com/adeturner/observability"
)

func (aggtot *aggregateTotal) print(cnt int) {

	i := 0

	for k, v := range aggtot.items {
		observability.Logger("Info", fmt.Sprintf("%s -> %v\n", k, v))
		i++
		if i > cnt {
			break
		}
	}
}

func (aggtot *aggregateTotal) printCount() {
	observability.Logger("Info", fmt.Sprintf("managedDiskLookup has %d records\n", len(aggtot.items)))
}

func (aggtot *aggregateTotal) init() {
	aggtot.items = make(map[string]*aggregateTotalItem)
}

func (aggtot *aggregateTotal) add(reportingCategory, reportingSubCategory, unitOfMeasure string, quantity float64, costInBillingCurrency float64) {

	// stringKey: reportingCategory +"/"+  reportingSubCategory
	key := fmt.Sprintf("%s:%s:%s", reportingCategory, reportingSubCategory, unitOfMeasure)

	// var aggtoti aggregateTotalItem

	// initializes two variables - api will receive either the value of "key" from the map
	// or a "zero value" (in this case the empty string)
	// ok will receive a bool that will be set to true if "key" was actually present in the map
	// evaluates ok, which will be true if "key" was in the map
	if _, ok := aggtot.items[key]; !ok {

		// if not found initialise
		aggtoti := aggregateTotalItem{}
		aggtoti.reportingCategory = reportingCategory
		aggtoti.reportingSubCategory = reportingSubCategory
		aggtoti.UnitOfMeasure = unitOfMeasure
		aggtoti.CostInBillingCurrency = 0
		aggtoti.Quantity = 0
		aggtot.items[key] = &aggtoti

	} else {
		// observability.Logger("Info", fmt.Sprintf("oldQ=%f quantity=%f", aggtoti.Quantity, quantity))
	}

	aggtot.items[key].Quantity += quantity
	aggtot.items[key].CostInBillingCurrency += costInBillingCurrency

	// observability.Logger("Info", fmt.Sprintf("newQ=%f", aggtot.items[key].Quantity))

}
