package azureBilling

type AzurePrices struct {
	fileLocation                    string
	CustomerBillingCurrencyEntityId string
	CustomerEntityType              string
	Items                           []AzurePriceAPI
	NextPageLink                    string
	Count                           int
}
