package azureBilling

import (
	"fmt"
)

type azurePrices struct {
	CustomerBillingCurrencyEntityId string
	CustomerEntityType              string
	Items                           []AzurePriceItem
	NextPageLink                    string
	Count                           int
}

type AzurePriceItem struct {
	CurrencyCode         string
	TierMinimumUnits     float64
	ReservationTerm      string
	RetailPrice          float64
	UnitPrice            float64
	ArmRegionName        string
	Location             string
	EffectiveStartDate   string
	MeterId              string
	MeterName            string
	ProductId            string
	SkuId                string
	ProductName          string
	SkuName              string
	ServiceName          string
	ServiceId            string
	ServiceFamily        string
	UnitOfMeasure        string
	ItemType             string `json:"type"`
	IsPrimaryMeterRegion bool
	ArmSkuName           string
}

func (ap *azurePrices) getItemCount() int {
	return ap.Count
}

func (ap *azurePrices) getCSVHeader() []byte {
	return []byte(fmt.Sprintf("\"%s\",\"%s\",\"%s\",\"%s\",\"%s\",\"%s\",\"%s\",\"%s\",\"%s\",\"%s\",\"%s\",\"%s\",\"%s\",\"%s\",\"%s\"\n",
		"MeterId", "MeterName", "ProductName", "SkuName", "ArmSkuName", "ServiceFamily", "ServiceName", "Location", "UnitOfMeasure", "ItemType", "ReservationTerm", "EffectiveStartDate", "TierMinimumUnits", "UnitPrice", "RetailPrice"))
}

func (i *AzurePriceItem) getCSVRow() []byte {
	// 2018-05-01T00:00:00Z
	effdate := fmt.Sprintf("%s-%s-%s", i.EffectiveStartDate[8:10], i.EffectiveStartDate[5:7], i.EffectiveStartDate[0:4])
	return []byte(fmt.Sprintf("\"%s\",\"%s\",\"%s\",\"%s\",\"%s\",\"%s\",\"%s\",\"%s\",\"%s\",\"%s\",\"%s\",\"%s\",%f,%f,%f\n", i.MeterId, i.MeterName, i.ProductName, i.SkuName, i.ArmSkuName, i.ServiceFamily, i.ServiceName, i.Location, i.UnitOfMeasure, i.ItemType, i.ReservationTerm, effdate, i.TierMinimumUnits, i.UnitPrice, i.RetailPrice))
}
