package billingModels

import "time"

type PriceMeterItem struct {
	EffectiveStartDate time.Time
	TierMinimumUnits   float64
	ReservationTerm    string
	RetailPrice        float64
	UnitPrice          float64
	UnitOfMeasure      string
	SkuName            string
	ArmSkuName         string
}
