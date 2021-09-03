package billingModels

import "time"

// map["MeterId"]value)
type priceMeter struct {
	PeriodEndDate time.Time
	Items         map[string]PriceMeterItem
}
