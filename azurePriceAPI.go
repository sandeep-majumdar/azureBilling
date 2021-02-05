package azureBilling

import "fmt"

func (api *AzurePriceAPI) getCSVRow() []byte {
	// 2018-05-01T00:00:00Z
	effdate := fmt.Sprintf("%s-%s-%s", api.EffectiveStartDate[8:10], api.EffectiveStartDate[5:7], api.EffectiveStartDate[0:4])
	return []byte(fmt.Sprintf("\"%s\",\"%s\",\"%s\",\"%s\",\"%s\",\"%s\",\"%s\",\"%s\",\"%s\",\"%s\",\"%s\",\"%s\",%f,%f,%f\n", api.MeterId, api.MeterName, api.ProductName, api.SkuName, api.ArmSkuName, api.ServiceFamily, api.ServiceName, api.Location, api.UnitOfMeasure, api.ItemType, api.ReservationTerm, effdate, api.TierMinimumUnits, api.UnitPrice, api.RetailPrice))
}
