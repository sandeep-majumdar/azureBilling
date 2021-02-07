package azureBilling

func (rcli *reportingCategoryLookupItem) setValues(i []string) {
	rcli.meterCategory = i[0]
	rcli.reportingCategory = i[1]
	rcli.reportingSubCategory = i[2]
}
