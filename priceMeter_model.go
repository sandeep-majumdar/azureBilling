package azureBilling

import "time"

// map["MeterId"]value)
type priceMeter struct {
	periodEndDate time.Time
	items         map[string]priceMeterItem
}
