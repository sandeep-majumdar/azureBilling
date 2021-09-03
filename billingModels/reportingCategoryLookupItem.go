package billingModels

func (rcli *ReportingCategoryLookupItem) setValues(i []string) {
	rcli.MeterCategory = i[0]
	rcli.MeterSubCategory = i[1]
	rcli.ReportingCategory = i[2]
	rcli.ReportingSubCategory = i[3]
}
