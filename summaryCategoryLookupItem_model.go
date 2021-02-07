package azureBilling

type summaryCategoryLookupItem struct {
	// key
	reportingCategory    string
	reportingSubCategory string
	UnitOfMeasure        string
	// values
	Summary         string
	QuantityDivisor string
}
