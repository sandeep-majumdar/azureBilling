package azureBilling

type aggregateTotal struct {
	// stringKey: reportingCategory +"/"+  reportingSubCategory + "/" + UnitOfMeasure
	items map[string]*aggregateTotalItem
}

type aggregateTotalItem struct {
	reportingCategory     string
	reportingSubCategory  string
	UnitOfMeasure         string
	Quantity              float64
	CostInBillingCurrency float64
}
