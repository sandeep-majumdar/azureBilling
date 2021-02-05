package azureBilling

type AzurePriceAPI struct {
	CurrencyCode         string
	TierMinimumUnits     float64
	ReservationTerm      string
	RetailPrice          float64
	UnitPrice            float64
	ArmRegionName        string
	Location             string
	EffectiveStartDate   string
	MeterId              string
	MeterName            string
	ProductId            string
	SkuId                string
	ProductName          string
	SkuName              string
	ServiceName          string
	ServiceId            string
	ServiceFamily        string
	UnitOfMeasure        string
	ItemType             string `json:"type"`
	IsPrimaryMeterRegion bool
	ArmSkuName           string
}
