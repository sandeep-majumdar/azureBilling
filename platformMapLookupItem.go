package azureBilling

func (pmli *platformMapLookupItem) setValues(i []string) {

	pmli.platform = i[0]
	pmli.productCode = i[1]
	pmli.environmentType = i[2]
	pmli.subscriptionId = i[3]
	pmli.rgName = i[4]

}
