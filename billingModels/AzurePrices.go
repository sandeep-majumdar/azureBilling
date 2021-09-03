package billingModels

import (
	"encoding/csv"
	"encoding/json"
	"fmt"
	"io"
	"os"

	"github.com/adeturner/azureBilling/utils"
	"github.com/adeturner/observability"
)

func (ap *AzurePrices) ReadAzurePrices(dateStr string) error {

	// initialise the meterid lookup
	MeterLookup.Init(dateStr)

	f, err := os.Open(ap.fileLocation)
	if err != nil {
		observability.Logger("Error", fmt.Sprintf("Unable to read input file=%s err=%s", ap.fileLocation, err))
	}
	defer f.Close()

	cnt := 0

	if err == nil {

		r := csv.NewReader(f)
		for {

			record, err := r.Read()
			if err == io.EOF {
				break
			} else if err != nil {
				observability.Logger("Error", fmt.Sprintf("Unable to parse file as CSV; file=%s err=%s", ap.fileLocation, err))
				break
			}

			cnt++

			// skip the first row (header)
			if cnt > 1 {
				i := PriceItem{}
				i.setValues(record)

				// meterLookup is a global variable declared in AzurePriceMeter.go
				MeterLookup.add(i)
			}

		}
	}

	observability.LogMemory("Info")
	//meterLookup.print(5)
	MeterLookup.printCount()

	// meterLookup.print(50000)

	return err

}

func (ap *AzurePrices) FileExists() bool {

	retval := true

	f, err := os.Open(ap.fileLocation)

	if err != nil {
		observability.Logger("Error", fmt.Sprintf("Unable to read input file=%s err=%s", ap.fileLocation, err))
		retval = false
	} else {
		observability.Logger("Info", fmt.Sprintf("Successfully found file=%s", ap.fileLocation))
	}

	defer f.Close()

	return retval
}

func (ap *AzurePrices) SetFile(filePath string) {
	ap.fileLocation = filePath
}

func (ap *AzurePrices) getItemCount() int {
	return ap.Count
}

func (ap *AzurePrices) getCSVHeader() []byte {
	return []byte(fmt.Sprintf("\"%s\",\"%s\",\"%s\",\"%s\",\"%s\",\"%s\",\"%s\",\"%s\",\"%s\",\"%s\",\"%s\",\"%s\",\"%s\",\"%s\",\"%s\"\n",
		"MeterId", "MeterName", "ProductName", "SkuName", "ArmSkuName", "ServiceFamily", "ServiceName", "Location", "UnitOfMeasure", "ItemType", "ReservationTerm", "EffectiveStartDate", "TierMinimumUnits", "UnitPrice", "RetailPrice"))
}

func (ap *AzurePrices) check(e error) {
	if e != nil {
		observability.Logger("Error", fmt.Sprintf("%v", e))
		panic(e)
	}
}

func (ap *AzurePrices) WritePriceHeader(w io.Writer, prices AzurePrices) {

	_, err := w.Write(prices.getCSVHeader())
	ap.check(err)

}

func (ap *AzurePrices) WritePriceOutput(w io.Writer, prices AzurePrices) {
	// Appends a line of text to the file text.log. It creates the file if it doesn’t already exist.
	// f, err := os.OpenFile("fname", os.O_APPEND|os.O_CREATE|os.O_WRONLY, 0644)

	for i := 0; i < prices.getItemCount(); i++ {
		_, err := w.Write(prices.Items[i].getCSVRow())
		ap.check(err)
	}
}

func (ap *AzurePrices) GeneratePrices(filename string) {

	var r RestClient
	var prices AzurePrices
	var fs utils.FileSystem = utils.LocalFS{}
	var lastNextPageLink string

	file, err := fs.Create(filename)
	if err != nil {
		observability.Logger("Error", fmt.Sprintf("Failed to open file: %v", err))
	}
	defer file.Close()

	ap.WritePriceHeader(file, prices)

	if err == nil {
		r.Init()
		url := "https://prices.azure.com/api/retail/prices"
		for true {
			respData, err := r.GET(url)
			if err == nil {
				err = json.Unmarshal(respData, &prices)
				ap.check(err)
				url = prices.NextPageLink
				observability.Logger("Debug", fmt.Sprintf("Successfully unmarshalled prices: %d %s", prices.Count, prices.NextPageLink))
				if lastNextPageLink == url {
					break
				} else {
					lastNextPageLink = url
				}
			}
			if err == nil {
				ap.WritePriceOutput(file, prices)
			}
			if err != nil || url == "" {
				break
			}
		}
	}
}
