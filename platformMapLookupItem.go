package azureBilling

func (pmli *platformMapLookupItem) setValues(i []string) {

	pmli.portfolio = i[0]
	pmli.platform = i[1]
	pmli.productCode = i[2]
	pmli.environmentType = i[3]
	pmli.subscriptionId = i[4]
	pmli.rgId = i[5]
	pmli.rgName = i[6]

}
