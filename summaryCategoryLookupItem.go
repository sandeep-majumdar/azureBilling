package azureBilling

func (scli *summaryCategoryLookupItem) setValues(i []string) {

	scli.reportingCategory = i[0]
	scli.reportingSubCategory = i[1]
	scli.UnitOfMeasure = i[2]
	// values
	scli.Summary = i[3]
	scli.QuantityDivisor = i[4]

}
