package azureBilling

// stringKey: reportingCategory + "/" + reportingSubCategory + "/" + SubscriptionId + "/" + ResourceGroup + "/" + MeterId

type aggregateResourceGroup struct {
	items map[string]*aggregateResourceGroupItem
}

type aggregateResourceGroupItem struct {
	// key fields
	reportingCategory    string
	reportingSubCategory string
	SubscriptionId       string
	ResourceGroup        string
	MeterId              string
	// core values
	summaryCategory       string
	quantityDivisor       string
	portfolio             string
	Platform              string
	product               string
	EnvironmentType       string
	summaryQuantity       float64
	CostInBillingCurrency float64
	// detail for deep dive
	ResourceLocation string
	ProductName      string
	MeterCategory    string
	MeterSubCategory string
	MeterName        string
	MeterRegion      string
	UnitOfMeasure    string
	EffectivePrice   string
	CostCenter       string
	ConsumedService  string
	ReservationId    string
	Term             string
	Quantity         float64
	UnitPrice        string
}
