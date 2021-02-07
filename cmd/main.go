package main

import (
	"github.com/adeturner/azureBilling"
	"github.com/adeturner/observability"
)

func main() {

	// set logging values
	observability.SetCausationId("1")
	observability.GenCorrId()

	// load config
	azureBilling.ConfigMap.LoadConfiguration("config.json")
	c := azureBilling.ConfigMap
	billingCSVFile := c.WorkingDirectory + c.BillingCSVFile
	azurePricesCSVFile := c.WorkingDirectory + c.OutputAzurePricesCSVFile
	billingCSVMaxDate := c.BillingCSVMaxDate
	lookupDirectory := c.LookupDirectory

	// if AzurePricesCSVFile doesnt exist, create a new one
	observability.Logger("Info", azurePricesCSVFile)
	ap := azureBilling.AzurePrices{}
	ap.SetFile(azurePricesCSVFile)
	if !ap.FileExists() {
		ap.GeneratePrices(azurePricesCSVFile)
	}
	ap.ReadAzurePrices(billingCSVMaxDate)

	// Lookups are expected to exist and must be manually maintained
	azureBilling.VmSizeLookup.Read(lookupDirectory + "vmSizes.csv")
	azureBilling.ManagedDiskLookup.Read(lookupDirectory + "managedDisks.csv")
	azureBilling.PlatformMapLookup.Read(lookupDirectory + "platformMap.csv")
	azureBilling.ReportingCategoryLookup.Read(lookupDirectory + "reportingCategories.csv")
	azureBilling.SummaryCategoryLookup.Read(lookupDirectory + "summaryCategories.csv")

	// 6514840 records in test file in 5 mins
	bcsv := azureBilling.BillingCSV{}
	bcsv.SetFile(billingCSVFile)
	bcsv.ProcessFile()

}
