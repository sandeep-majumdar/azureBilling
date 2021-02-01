# azureBilling

Generate a CSV from the public Azure Prices API: https://prices.azure.com/api/retail/prices

## Run the code

```bash
go run cmd/main.go
```

## Azure Documentation

https://docs.microsoft.com/en-us/rest/api/cost-management/retail-prices/azure-retail-prices

```json
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
```
