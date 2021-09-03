package billingModels

import (
	"fmt"
	"strings"

	"github.com/adeturner/azureBilling/utils"
	"github.com/adeturner/observability"
)

func (pm *priceMeter) print(cnt int) {

	i := 0

	for k, v := range pm.Items {
		observability.Logger("Info", fmt.Sprintf("%s -> %v\n", k, v))
		i++
		if i > cnt {
			break
		}
	}
}

func (pm *priceMeter) printCount() {
	observability.Logger("Info", fmt.Sprintf("MeterLookup has %d records\n", len(pm.Items)))
}

func (pm *priceMeter) Init(dateStr string) {

	var err error

	MeterLookup.PeriodEndDate, err = utils.DateStrToTime(dateStr)
	if err != nil {
		observability.Logger("Fatal", fmt.Sprintf("Invalid date passed to function %s", dateStr))
	}

	MeterLookup.Items = make(map[string]PriceMeterItem)

}

func (pm *priceMeter) dateBefore(dateStr string) bool {

	var retval bool

	t1, err := utils.DateStrToTime(dateStr)

	if err != nil {
		// do nothing
	} else {
		retval = pm.PeriodEndDate.Before(t1)
	}

	return retval
}

func (pm *priceMeter) add(i PriceItem) {

	itemDate, err := utils.DateStrToTime(i.EffectiveStartDate)

	var key string

	// added b = true to bypass date checks temporarily (maybe forever...)
	b := true

	if err != nil {
		// do nothing
	} else {
		// make sure the item effiective date is before the billing period period
		if b || itemDate.Before(pm.PeriodEndDate) {

			v := PriceMeterItem{}
			v.setValues(i)
			key = pm.getKey(i.MeterId)

			// if the meterid date is not set, set it to the item effective date
			if b || pm.Items[key].EffectiveStartDate.IsZero() ||
				pm.Items[key].EffectiveStartDate.After(itemDate) {

				pm.Items[key] = v
			}
		}
	}

}

func (pm *priceMeter) getKey(meterId string) string {
	return strings.ToLower(fmt.Sprintf(":%s:", meterId))
}

func (pm *priceMeter) Get(meterId string) (PriceMeterItem, bool) {

	key := pm.getKey(meterId)

	rcli, ok := pm.Items[key]
	if !ok {
		//observability.Logger("Error", fmt.Sprintf("Unable to find PriceMeterItem for key=%s", key))
	}

	return rcli, ok
}
