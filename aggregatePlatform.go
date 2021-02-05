package azureBilling

import (
	"fmt"

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
	aggp.items = make(map[string]aggregatePlatformItem)
}

func (aggp *aggregatePlatform) add(reportingCategory, reportingSubCategory, platform string, quantity float64, costInBillingCurrency float64) {

	key := fmt.Sprintf("%s/%s/%s", reportingCategory, reportingSubCategory, platform)

	var api aggregatePlatformItem

	// initializes two variables - api will receive either the value of "key" from the map
	// or a "zero value" (in this case the empty string)
	// ok will receive a bool that will be set to true if "key" was actually present in the map
	// evaluates ok, which will be true if "key" was in the map
	if api, ok := aggp.items[key]; !ok {

		// if not found initialise
		api = aggregatePlatformItem{}
		api.platform = platform
		api.reportingCategory = reportingCategory
		api.reportingSubCategory = reportingSubCategory
		api.CostInBillingCurrency = 0
		api.Quantity = 0
		aggp.items[key] = api
	}

	api.Quantity += quantity
	api.CostInBillingCurrency += costInBillingCurrency

}
