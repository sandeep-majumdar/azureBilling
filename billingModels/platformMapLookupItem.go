package billingModels

func (pmli *PlatformMapLookupItem) SetValues(i []string) {

	pmli.Portfolio = i[0]
	pmli.Platform = i[1]
	pmli.ProductCode = i[2]
	pmli.EnvironmentType = i[3]
	pmli.SubscriptionId = i[4]
	pmli.RgId = i[5]
	pmli.RgName = i[6]

}
