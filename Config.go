package azureBilling

import (
	"encoding/json"
	"fmt"
	"os"

	"github.com/adeturner/observability"
)

type Config struct {
	WorkingDirectory               string `json:"workingDirectory"`
	BillingCSVFile                 string `json:"billingCSVFile"`
	OutputAzurePricesCSVFile       string `json:"outputAzurePricesCSVFile"`
	BillingCSVMaxDate              string `json:"billingCSVMaxDate"`
	LookupDirectory                string `json:"lookupDirectory"`
	OutputAggregateTotalCsvFile    string `json:"outputAggregateTotalCsvFile"`
	OutputAggregatePlatformCsvFile string `json:"outputAggregatePlatformCsvFile"`
	OutputAggregateRGCsvFile       string `json:"outputAggregateRGCsvFile"`
}

func (cfg *Config) LoadConfiguration(file string) {

	configFile, err := os.Open(file)
	defer configFile.Close()
	if err != nil {
		fmt.Println(err.Error())
	}
	jsonParser := json.NewDecoder(configFile)
	jsonParser.Decode(cfg)

	observability.Logger("Info", fmt.Sprintf("%v", cfg))

}
