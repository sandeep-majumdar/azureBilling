# azureBilling

Generate a CSV from the public Azure Prices API: https://prices.azure.com/api/retail/prices

## Run the code

```bash
go run cmd/main.go
```

## Azure Pricing API

See https://docs.microsoft.com/en-us/rest/api/cost-management/retail-prices/azure-retail-prices

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

Running the code calls the API repeatedly, and writes a file in the current directory called azureBilling.csv

```txt
"MeterId","MeterName","ProductName","SkuName","ArmSkuName","ServiceFamily","ServiceName","Location","UnitOfMeasure","ItemType","ReservationTerm","EffectiveStartDate","TierMinimumUnits","UnitPrice","RetailPrice"
"0001d427-82df-4d83-8ab2-b60768527e08","E10 Disks","Standard SSD Managed Disks","E10 LRS","","Storage","Storage","UK South","1/Month","Consumption","","01-11-2018",0.000000,10.560000,10.560000
"0001e46a-9285-5fa8-b48a-240e307a24f7","A3 Spot","Virtual Machines A Series Windows","A3 Spot","Standard_A3","Compute","Virtual Machines","UK North","1 Hour","DevTestConsumption","","16-10-2019",0.000000,0.062988,0.062988
"0001e46a-9285-5fa8-b48a-240e307a24f7","A3 Spot","Virtual Machines A Series Windows","A3 Spot","Standard_A3","Compute","Virtual Machines","UK North","1 Hour","Consumption","","16-10-2019",0.000000,0.190000,0.190000
```
