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
	//
	portfolio             string
	Platform              string
	ResourceLocation      string
	ProductName           string
	MeterCategory         string
	MeterSubCategory      string
	MeterName             string
	MeterRegion           string
	UnitOfMeasure         string
	EffectivePrice        string
	CostInBillingCurrency float64
	CostCenter            string
	ConsumedService       string
	ResourceId            string
	ReservationId         string
	Term                  string
	Quantity              float64
	UnitPrice             string
}
