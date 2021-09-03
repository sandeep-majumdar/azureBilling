package billingModels

import "github.com/adeturner/azureBilling/utils"

func (pmv *PriceMeterItem) setValues(i PriceItem) {
	pmv.EffectiveStartDate, _ = utils.DateStrToTime(i.EffectiveStartDate)
	pmv.TierMinimumUnits = i.TierMinimumUnits
	pmv.ReservationTerm = i.ReservationTerm
	pmv.RetailPrice = i.RetailPrice
	pmv.UnitPrice = i.UnitPrice
	pmv.UnitOfMeasure = i.UnitOfMeasure
	pmv.SkuName = i.SkuName
	pmv.ArmSkuName = i.ArmSkuName
}
