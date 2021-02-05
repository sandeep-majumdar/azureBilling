package azureBilling

import (
	"fmt"
	"strconv"

	"github.com/adeturner/observability"
)

func (i *PriceItem) print() {
	observability.Logger("Info", fmt.Sprintf("%v", i))
}

func (i *PriceItem) setValues(records []string) {

	var err error

	i.MeterId = records[0]
	i.MeterName = records[1]
	i.ProductName = records[2]
	i.SkuName = records[3]
	i.ArmSkuName = records[4]
	i.ServiceFamily = records[5]
	i.ServiceName = records[6]
	i.Location = records[7]
	i.UnitOfMeasure = records[8]
	i.ItemType = records[9]
	i.ReservationTerm = records[10]
	i.EffectiveStartDate = records[11]
	i.TierMinimumUnits, err = strconv.ParseFloat(records[12], 64)
	if err != nil {
		observability.Logger("Error", fmt.Sprintf("Failed to parse i.TierMinimumUnits, err = strconv.ParseFloat(records[12], 64) from %s", records[12]))
	}
	i.UnitPrice, err = strconv.ParseFloat(records[13], 64)
	if err != nil {
		observability.Logger("Error", fmt.Sprintf("Failed to parse i.UnitPrice, err = strconv.ParseFloat(records[13], 64) from %s", records[13]))
	}
	i.RetailPrice, err = strconv.ParseFloat(records[14], 64)
	if err != nil {
		observability.Logger("Error", fmt.Sprintf("Failed to parse i.RetailPrice, err = strconv.ParseFloat(records[14], 64) from %s", records[14]))
	}

}
