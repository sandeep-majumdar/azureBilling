package azureBilling

func (pmv *priceMeterItem) setValues(i PriceItem) {
	pmv.EffectiveStartDate, _ = dateStrToTime(i.EffectiveStartDate)
	pmv.TierMinimumUnits = i.TierMinimumUnits
	pmv.ReservationTerm = i.ReservationTerm
	pmv.RetailPrice = i.RetailPrice
	pmv.UnitPrice = i.UnitPrice
	pmv.UnitOfMeasure = i.UnitOfMeasure
}
