package azureBilling

func (rcli *reportingCategoryLookupItem) setValues(i []string) {
	rcli.meterCategory = i[0]
	rcli.meterSubCategory = i[1]
	rcli.reportingCategory = i[2]
	rcli.reportingSubCategory = i[3]
}
