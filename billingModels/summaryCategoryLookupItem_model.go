package billingModels

type SummaryCategoryLookupItem struct {
	// key
	ReportingCategory    string
	ReportingSubCategory string
	UnitOfMeasure        string
	// values
	Summary         string
	QuantityDivisor string
}
