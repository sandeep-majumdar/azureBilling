package azureBilling

type platformMapLookupItem struct {
	// key
	subscriptionId string
	rgName         string
	//
	platform        string
	productCode     string
	environmentType string
}
