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

/*

https://docs.microsoft.com/en-us/rest/api/cost-management/retail-prices/azure-retail-prices

{
	"BillingCurrency": "USD",
	"CustomerEntityId": "Default",
	"CustomerEntityType": "Retail",
	"Items": [{
		"currencyCode": "USD",
		"tierMinimumUnits": 0.0,
		"retailPrice": 0.06,
		"unitPrice": 0.06,
		"armRegionName": "Global",
		"location": "Global",
		"effectiveStartDate": "2017-09-19T00:00:00Z",
		"meterId": "ed8a651a-e0a3-4de6-a8ae-3b4ce8cb72cf",
		"meterName": "LRS Data Stored",
		"productId": "DZH318Z0BP0B",
		"skuId": "DZH318Z0BP0B/004M",
		"productName": "Files",
		"skuName": "Standard LRS",
		"serviceName": "Storage",
		"serviceId": "DZH317F1HKN0",
		"serviceFamily": "Storage",
		"unitOfMeasure": "1 GB/Month",
		"type": "Consumption",
		"isPrimaryMeterRegion": true,
		"armSkuName": ""
	}],
	"NextPageLink": null,
	"Count": 1
}
*/
