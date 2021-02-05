package azureBilling

type aggregatePlatform struct {
	// stringKey: reportingCategory +  reportingSubCategory + platform + unitofmeasure
	items map[string]*aggregatePlatformItem
}

type aggregatePlatformItem struct {
	reportingCategory     string
	reportingSubCategory  string
	portfolio             string
	platform              string
	unitOfMeasure         string
	Quantity              float64
	CostInBillingCurrency float64
}
