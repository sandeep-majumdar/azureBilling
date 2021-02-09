package azureBilling

import (
	"fmt"
	"strings"

	"github.com/adeturner/observability"
)

func (pm *priceMeter) print(cnt int) {

	i := 0

	for k, v := range pm.items {
		observability.Logger("Info", fmt.Sprintf("%s -> %v\n", k, v))
		i++
		if i > cnt {
			break
		}
	}
}

func (pm *priceMeter) printCount() {
	observability.Logger("Info", fmt.Sprintf("MeterLookup has %d records\n", len(pm.items)))
}

func (pm *priceMeter) init(dateStr string) {

	var err error

	MeterLookup.periodEndDate, err = dateStrToTime(dateStr)
	if err != nil {
		observability.Logger("Fatal", fmt.Sprintf("Invalid date passed to function %s", dateStr))
	}

	MeterLookup.items = make(map[string]priceMeterItem)

}

func (pm *priceMeter) dateBefore(dateStr string) bool {

	var retval bool

	t1, err := dateStrToTime(dateStr)

	if err != nil {
		// do nothing
	} else {
		retval = pm.periodEndDate.Before(t1)
	}

	return retval
}

func (pm *priceMeter) add(i PriceItem) {

	itemDate, err := dateStrToTime(i.EffectiveStartDate)

	var key string

	// added b = true to bypass date checks temporarily (maybe forever...)
	b := true

	if err != nil {
		// do nothing
	} else {
		// make sure the item effiective date is before the billing period period
		if b || itemDate.Before(pm.periodEndDate) {

			v := priceMeterItem{}
			v.setValues(i)
			key = pm.getKey(i.MeterId)

			// if the meterid date is not set, set it to the item effective date
			if b || pm.items[key].EffectiveStartDate.IsZero() ||
				pm.items[key].EffectiveStartDate.After(itemDate) {

				pm.items[key] = v
			}
		}
	}

}

func (pm *priceMeter) getKey(meterId string) string {
	return strings.ToLower(fmt.Sprintf(":%s:", meterId))
}

func (pm *priceMeter) get(meterId string) (priceMeterItem, bool) {

	key := pm.getKey(meterId)

	rcli, ok := pm.items[key]
	if !ok {
		//observability.Logger("Error", fmt.Sprintf("Unable to find priceMeterItem for key=%s", key))
	}

	return rcli, ok
}
