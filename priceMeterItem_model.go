package azureBilling

import "time"

type priceMeterItem struct {
	EffectiveStartDate time.Time
	TierMinimumUnits   float64
	ReservationTerm    string
	RetailPrice        float64
	UnitPrice          float64
	UnitOfMeasure      string
	SkuName            string
	ArmSkuName         string
}
