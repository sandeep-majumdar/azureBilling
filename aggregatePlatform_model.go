package azureBilling

type aggregatePlatform struct {
	// stringKey: reportingCategory +  reportingSubCategory + platform
	items map[string]aggregatePlatformItem
}

type aggregatePlatformItem struct {
	reportingCategory     string
	reportingSubCategory  string
	platform              string
	Quantity              float64
	CostInBillingCurrency float64
}
