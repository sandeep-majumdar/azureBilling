package cli

import (
	"fmt"

	"github.com/adeturner/azureBilling/azure"
	"github.com/adeturner/azureBilling/billing"
	"github.com/adeturner/azureBilling/billingModels"
	"github.com/adeturner/azureBilling/config"
	"github.com/adeturner/azureBilling/observability"
	"github.com/adeturner/azureBilling/rightsizing"
	"github.com/spf13/cobra"
)

var azureCmd = &cobra.Command{
	Use:   "azure",
	Short: "Azure commands",
	Long: `
	Invoke azure commands i.e. the azure cli
	`,
	Run: azureVerb,
}

var azureCliCmd = &cobra.Command{
	Use:   "cli",
	Short: "Install the environment",
	Long:  `Install the environment according to the config file`,
	Run:   azcli,
}

var azureBillingCmd = &cobra.Command{
	Use:   "billing",
	Short: "Process the bill file",
	Long:  `Process the bill file, could take a while`,
	Run:   azbilling,
}

var azureRightsizingCmd = &cobra.Command{
	Use:   "rightsizing",
	Short: "Recommend rightsizing opportunity from the bill file",
	Long:  `Recommend rightsizing opportunity from the bill file`,
	Run:   rightsizingVerb,
}

var azureVmCmd = &cobra.Command{
	Use:   "vm",
	Short: "Recommend rightsizing opportunity for VMs",
	Long:  `Recommend rightsizing opportunity for VMs, could take a while`,
	Run:   azrightsizing,
}

func init() {

	// ##########################
	// go run cmd/main.go azure cli -l login.txt -f file.txt -s subscriptionID
	// go run cmd/main.go azure cli -f examples/azureMonitor.txt
	// ##########################
	azureCliCmd.Flags().StringVarP(&login, "login", "l", "", "The filename containing the azure cli login commands")
	azureCliCmd.MarkFlagRequired("login")

	azureCliCmd.Flags().StringVarP(&file, "file", "f", "", "The filename containing the azure cli commands to run")
	azureCliCmd.MarkFlagRequired("file")

	azureCliCmd.Flags().StringVarP(&subscription, "subscription", "s", "", "The subscription to connect to")
	azureCliCmd.MarkFlagRequired("subscription")

	// ##########################
	// go run cmd/main.go azure billing -c ./config.json
	// ##########################
	azureBillingCmd.Flags().StringVarP(&configFile, "config", "c", "", "The filename containing the configuration")
	azureBillingCmd.MarkFlagRequired("config")

	// ##########################
	// go run cmd/main.go azure rightsizing -c ./config.json -l login.txt
	// ##########################
	azureVmCmd.Flags().StringVarP(&configFile, "config", "c", "", "The filename containing the configuration")
	azureVmCmd.MarkFlagRequired("config")

	// ####################################################
	azureRightsizingCmd.AddCommand(azureVmCmd)
	azureCmd.AddCommand(azureCliCmd)
	azureCmd.AddCommand(azureBillingCmd)
	azureCmd.AddCommand(azureRightsizingCmd)
	azureCmd.Flags().BoolVarP(&verbose, "verbose", "v", false, "Print more information about request")
}

func azcli(ccmd *cobra.Command, args []string) {

	observability.SetCausationId("1")
	observability.GenCorrId()

	// load config
	config.ConfigMap.LoadConfiguration("config.json")
	c := config.ConfigMap
	observability.Info(fmt.Sprintf("Loaded config: %v", c))

	a := azure.NewAzureCli()
	a.SetSubscription(subscription)
	err := a.ExecFileLineByLine(login)
	if err == nil {
		err = a.ExecFileLineByLine(file)
	}
	if err != nil {
		observability.Error(err.Error())
	}
}

func azbilling(ccmd *cobra.Command, args []string) {
	// set logging values
	observability.SetCausationId("1")
	observability.GenCorrId()

	// load config
	config.ConfigMap.LoadConfiguration(configFile)
	c := config.ConfigMap
	billingCSVFile := c.WorkingDirectory + c.BillingCSVFile
	azurePricesCSVFile := c.WorkingDirectory + c.OutputAzurePricesCSVFile
	billingCSVMaxDate := c.BillingCSVMaxDate
	lookupDirectory := c.LookupDirectory

	// if AzurePricesCSVFile doesnt exist, create a new one
	observability.Logger("Info", azurePricesCSVFile)
	ap := billingModels.AzurePrices{}
	ap.SetFile(azurePricesCSVFile)
	if !ap.FileExists() {
		ap.GeneratePrices(azurePricesCSVFile)
	}
	ap.ReadAzurePrices(billingCSVMaxDate)

	// Lookups are expected to exist and must be manually maintained
	billingModels.VmSizeLookup.Read(lookupDirectory + "vmSizes.csv")
	billingModels.ManagedDiskLookup.Read(lookupDirectory + "managedDisks.csv")
	billingModels.PlatformMapLookup.Read(lookupDirectory + "platformMap.csv")
	billingModels.ReportingCategoryLookup.Read(lookupDirectory + "reportingCategories.csv")
	billingModels.SummaryCategoryLookup.Read(lookupDirectory + "summaryCategories.csv")

	// 6514840 records in test file in 5 mins
	bcsv := billing.BillingCSV{}
	bcsv.SetFile(billingCSVFile)
	bcsv.ProcessFile()
}

func azrightsizing(ccmd *cobra.Command, args []string) {

	// set logging values
	observability.SetCausationId("1")
	observability.GenCorrId()

	// load config
	config.ConfigMap.LoadConfiguration(configFile)
	c := config.ConfigMap
	billingCSVFile := c.WorkingDirectory + c.BillingCSVFile
	lookupDirectory := c.LookupDirectory

	billingModels.PlatformMapLookup.Read(lookupDirectory + "platformMap.csv")

	rsz := rightsizing.VmRightsizing{}
	rsz.SetFile(billingCSVFile)
	err := rsz.ProcessBillFile()

	if err == nil {
		err = rsz.ProcessMetrics()
	}
}

func azureVerb(ccmd *cobra.Command, args []string) {
	switch {
	default:
		fmt.Printf("azure verb must pass a sub-command\n")
		ccmd.Help()
	}
}

func rightsizingVerb(ccmd *cobra.Command, args []string) {
	switch {
	default:
		fmt.Printf("rightsizing verb must pass a sub-command\n")
		ccmd.Help()
	}
}
