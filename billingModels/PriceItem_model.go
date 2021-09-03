package billingModels

// "MeterId","MeterName","ProductName","SkuName","ArmSkuName","ServiceFamily",
// "ServiceName","Location","UnitOfMeasure","ItemType","ReservationTerm",
// "EffectiveStartDate","TierMinimumUnits","UnitPrice","RetailPrice"

type PriceItem struct {
	MeterId            string
	MeterName          string
	ProductName        string
	SkuName            string
	ArmSkuName         string
	ServiceFamily      string
	ServiceName        string
	Location           string
	UnitOfMeasure      string
	ItemType           string
	ReservationTerm    string
	EffectiveStartDate string
	TierMinimumUnits   float64
	UnitPrice          float64
	RetailPrice        float64
}
