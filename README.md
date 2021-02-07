# AzureBilling

This product offers a solution to enable executive reporting for Azure Billing in terms they should understand

The target scope of the output is to obtain Price (P) and Quantity (Q) by org unit/functional area/product.

Many other views are possible!

The output is in the form of a resource group level aggregate CSV that can be turned into Pivot tables as required

## Why is this program needed?

Unfortunately MS really dont make it easy with regards their billing csv format or to enable enterprise reporting

Our monthly billing CSV is 6GB+ so excel based reporting isnt possible.

UnitOfMeasure in the billing CSV is particularly useless, so we have to go to the effort of downloading from the Azure Prices API, described at the end of the readme.

Quantity in the billing CSV is some devious value that at the day level does not reflect the actual real world quantity; it does make sense when you aggregate across the month

When put together the failings of UnitOfMeasure and Quantity are large enough to justify this program

## How To

### 1. Customise the lookup CSVs

#### managedDisks.csv

The lookup for managed disk names providing the SizeGB. Note There is a challenge with "performance-enabled" managed disks, which will report the size of the underlying disk rather than the available volume

#### platformMap.csv

Provides a method to map your business to subscriptions/resource groups

Copy *platformMap_example.csv* to create your own as ours has been .gitignored

There are three levels available

- TopPlatform, e.g org, portfolio
- Platform, e.g. functional area, platform
- ProductCode, e.g. product, application

#### reportingCategories.csv

Maps azure products to reporting categories

#### summaryCategories.csv

Summarises reporting categories for executive reporting

#### vmSizes.csv

Lookup for VM names to CPU and MemGB


### 2. Create the config file

Copy config.example.json and edit the values as required

- workingDirectory: directory where the billing file is found and where output gets created
- azurePricesCSVFile: name of the billing file
- NumDaysInMonth: how many days in the month
- billingCSVMaxDate: max date of the billing file (used in Azure Price API generation)
- outputAzurePricesCSVFile - name of the generated azure prices lookup
- outputAggregateRGCSVFile - name of the main output file, the one you want
- lookup directory, where the supplied/customised lookups are located

Example:

```json
{
"workingDirectory": "/mnt/c/Users/nnn/Downloads/",
"billingCSVFile": "ActualCostMtD_174c3ad5-7777-4e8d-824d-52f47a1f84dd.csv",
"azurePricesCSVFile": "azurePrices.csv",
"billingCSVMaxDate": "31/01/2021",
"NumDaysInMonth": "31",
"lookupDirectory": "./lookups/",
"outputAzurePricesCSVFile": "azurePrices.csv",
"outputAggregateRGCSVFile": "aggregateRG.csv"
}
```

### 3. Clone and run the code

With core i7 and SSD you can expect ~3.5mins per 1 million records

```bash
go run cmd/main.go
```

### 4. Create reports from the output

Create pivot tables from the data in whatever fashion you prefer. Below are a couple of sample executive views:

### Price

Filter: Portfolio, Platform, ProductName
Rows: ReportingCategory, ReportingSubCategory
Columns: SummaryCategory
Values: Sum of summaryQuantity

### Quantity

Filter: Portfolio, Platform, ProductName
Rows: ReportingCategory, ReportingSubCategory
Columns: SummaryCategory
Values: Sum of CostInBillingCurrency

## Sample execution output

```txt
I 2021-02-07 09:35:37.0074 [1 => 20e252d9-5b25-495f-8598-78c53d277b4e]  Config.go:32    &{/mnt/c/Users/adria/Downloads/ ActualCostMtD_174c3ad5-77af-4e8d-824d-52f47a1f84dd.csv azurePrices.csv 31/01/2021 31 ./lookups/ aggregateRG.csv execSummary.csv}
I 2021-02-07 09:35:37.0074 [1 => 20e252d9-5b25-495f-8598-78c53d277b4e]  main.go:23      /mnt/c/Users/adria/Downloads/azurePrices.csv
I 2021-02-07 09:35:37.0097 [1 => 20e252d9-5b25-495f-8598-78c53d277b4e]  AzurePrices.go:73       Successfully found file=/mnt/c/Users/adria/Downloads/azurePrices.csv
I 2021-02-07 09:35:39.6759 [1 => 20e252d9-5b25-495f-8598-78c53d277b4e]  AzurePrices.go:53       MEMORY Alloc = 84 MiB   TotalAlloc = 215 MiB    Sys = 138 MiB   NumGC = 16
I 2021-02-07 09:35:39.6760 [1 => 20e252d9-5b25-495f-8598-78c53d277b4e]  priceMeter.go:23        MeterLookup has 109476 records

I 2021-02-07 09:35:39.6766 [1 => 20e252d9-5b25-495f-8598-78c53d277b4e]  vmLookup.go:72  MEMORY Alloc = 84 MiB   TotalAlloc = 215 MiB    Sys = 138 MiB   NumGC = 16
I 2021-02-07 09:35:39.6767 [1 => 20e252d9-5b25-495f-8598-78c53d277b4e]  vmLookup.go:27  vmSizeLookup has 310 records

I 2021-02-07 09:35:39.6768 [1 => 20e252d9-5b25-495f-8598-78c53d277b4e]  managedDiskLookup.go:71 MEMORY Alloc = 84 MiB   TotalAlloc = 215 MiB    Sys = 138 MiB   NumGC = 16
I 2021-02-07 09:35:39.6768 [1 => 20e252d9-5b25-495f-8598-78c53d277b4e]  managedDiskLookup.go:27 managedDiskLookup has 37 records

I 2021-02-07 09:35:39.6833 [1 => 20e252d9-5b25-495f-8598-78c53d277b4e]  platformMapLookup.go:70 MEMORY Alloc = 88 MiB   TotalAlloc = 220 MiB    Sys = 138 MiB   NumGC = 16
I 2021-02-07 09:35:39.6834 [1 => 20e252d9-5b25-495f-8598-78c53d277b4e]  platformMapLookup.go:26 platformMapLookup has 6533 records

I 2021-02-07 09:35:39.6839 [1 => 20e252d9-5b25-495f-8598-78c53d277b4e]  reportingCategoryLoo:70 MEMORY Alloc = 89 MiB   TotalAlloc = 220 MiB    Sys = 138 MiB   NumGC = 16
I 2021-02-07 09:35:39.6840 [1 => 20e252d9-5b25-495f-8598-78c53d277b4e]  reportingCategoryLoo:26 reportingCategoryLookup has 80 records

I 2021-02-07 09:35:39.6850 [1 => 20e252d9-5b25-495f-8598-78c53d277b4e]  summaryCategoryLooku:75 MEMORY Alloc = 89 MiB   TotalAlloc = 220 MiB    Sys = 138 MiB   NumGC = 16
I 2021-02-07 09:35:39.6850 [1 => 20e252d9-5b25-495f-8598-78c53d277b4e]  summaryCategoryLooku:28 summaryCategoryLookup has 89 records

I 2021-02-07 09:36:01.7223 [1 => 20e252d9-5b25-495f-8598-78c53d277b4e]  BillingCSV.go:112       Processed 100000 rows of billing CSV
I 2021-02-07 09:36:01.7225 [1 => 20e252d9-5b25-495f-8598-78c53d277b4e]  BillingCSV.go:113       MEMORY Alloc = 119 MiB  TotalAlloc = 472 MiB    Sys = 273 MiB   NumGC = 20
I 2021-02-07 09:36:20.2349 [1 => 20e252d9-5b25-495f-8598-78c53d277b4e]  BillingCSV.go:112       Processed 200000 rows of billing CSV
I 2021-02-07 09:36:20.2351 [1 => 20e252d9-5b25-495f-8598-78c53d277b4e]  BillingCSV.go:113       MEMORY Alloc = 173 MiB  TotalAlloc = 718 MiB    Sys = 274 MiB   NumGC = 22
I 2021-02-07 09:36:40.3830 [1 => 20e252d9-5b25-495f-8598-78c53d277b4e]  BillingCSV.go:112       Processed 300000 rows of billing CSV
I 2021-02-07 09:36:40.3832 [1 => 20e252d9-5b25-495f-8598-78c53d277b4e]  BillingCSV.go:113       MEMORY Alloc = 155 MiB  TotalAlloc = 957 MiB    Sys = 341 MiB   NumGC = 24
I 2021-02-07 09:36:59.6619 [1 => 20e252d9-5b25-495f-8598-78c53d277b4e]  BillingCSV.go:112       Processed 400000 rows of billing CSV
...
...
...
I 2021-02-07 09:56:54.4787 [1 => 20e252d9-5b25-495f-8598-78c53d277b4e]  BillingCSV.go:112       Processed 6400000 rows of billing CSV
I 2021-02-07 09:56:54.4789 [1 => 20e252d9-5b25-495f-8598-78c53d277b4e]  BillingCSV.go:113       MEMORY Alloc = 293 MiB  TotalAlloc = 15542 MiB  Sys = 409 MiB   NumGC = 127
I 2021-02-07 09:57:12.3306 [1 => 20e252d9-5b25-495f-8598-78c53d277b4e]  BillingCSV.go:112       Processed 6500000 rows of billing CSV
I 2021-02-07 09:57:12.3307 [1 => 20e252d9-5b25-495f-8598-78c53d277b4e]  BillingCSV.go:113       MEMORY Alloc = 231 MiB  TotalAlloc = 15771 MiB  Sys = 409 MiB   NumGC = 129
I 2021-02-07 09:57:15.3958 [1 => 20e252d9-5b25-495f-8598-78c53d277b4e]  aggregateResourceGro:131        Writing to /mnt/c/Users/adria/Downloads/aggregateRG.csv
I 2021-02-07 09:57:23.8507 [1 => 20e252d9-5b25-495f-8598-78c53d277b4e]  BillingCSV.go:122       BillingCSV completed in 1304162.61 ms
I 2021-02-07 09:57:23.8508 [1 => 20e252d9-5b25-495f-8598-78c53d277b4e]  BillingCSV.go:124       Complete. Processed 6514840 rows of billing CSV
I 2021-02-07 09:57:23.8509 [1 => 20e252d9-5b25-495f-8598-78c53d277b4e]  BillingCSV.go:125       MEMORY Alloc = 190 MiB  TotalAlloc = 15875 MiB  Sys = 409 MiB   NumGC = 130

```

## Azure Pricing API

The code makes use of the public Azure Pricing API to get a useful Unit Of Measure

The first time you invoke azureBilling product it will scrape the API and write a file called azureBilling.csv (see config.json)

For the docs, see https://docs.microsoft.com/en-us/rest/api/cost-management/retail-prices/azure-retail-prices

### Example api return value

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

### Sample of csv generated

```txt
"MeterId","MeterName","ProductName","SkuName","ArmSkuName","ServiceFamily","ServiceName","Location","UnitOfMeasure","ItemType","ReservationTerm","EffectiveStartDate","TierMinimumUnits","UnitPrice","RetailPrice"
"0001d427-82df-4d83-8ab2-b60768527e08","E10 Disks","Standard SSD Managed Disks","E10 LRS","","Storage","Storage","UK South","1/Month","Consumption","","01-11-2018",0.000000,10.560000,10.560000
"0001e46a-9285-5fa8-b48a-240e307a24f7","A3 Spot","Virtual Machines A Series Windows","A3 Spot","Standard_A3","Compute","Virtual Machines","UK North","1 Hour","DevTestConsumption","","16-10-2019",0.000000,0.062988,0.062988
"0001e46a-9285-5fa8-b48a-240e307a24f7","A3 Spot","Virtual Machines A Series Windows","A3 Spot","Standard_A3","Compute","Virtual Machines","UK North","1 Hour","Consumption","","16-10-2019",0.000000,0.190000,0.190000
...
```
