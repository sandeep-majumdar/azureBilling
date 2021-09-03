package billingModels

func (scli *SummaryCategoryLookupItem) setValues(i []string) {

	scli.ReportingCategory = i[0]
	scli.ReportingSubCategory = i[1]
	scli.UnitOfMeasure = i[2]
	// values
	scli.Summary = i[3]
	scli.QuantityDivisor = i[4]

}
