package azureBilling

import (
	"encoding/json"
	"fmt"
	"io"

	"github.com/adeturner/observability"
)

func check(e error) {
	if e != nil {
		observability.Logger("Error", fmt.Sprintf("%v", e))
		panic(e)
	}
}

func WritePriceHeader(w io.Writer, prices azurePrices) {

	_, err := w.Write(prices.getCSVHeader())
	check(err)

}

func WritePriceOutput(w io.Writer, prices azurePrices) {
	// Appends a line of text to the file text.log. It creates the file if it doesn’t already exist.
	// f, err := os.OpenFile("fname", os.O_APPEND|os.O_CREATE|os.O_WRONLY, 0644)

	for i := 0; i < prices.getItemCount(); i++ {
		_, err := w.Write(prices.Items[i].getCSVRow())
		check(err)
	}
}

func GeneratePrices(filename string) {

	var r RestClient
	var prices azurePrices
	var fs fileSystem = localFS{}
	var lastNextPageLink string

	file, err := fs.Create(filename)
	if err != nil {
		observability.Logger("Error", fmt.Sprintf("Failed to open file: %v", err))
	}
	defer file.Close()

	WritePriceHeader(file, prices)

	if err == nil {
		r.Init()

		url := "https://prices.azure.com/api/retail/prices"

		for true {

			respData, err := r.GET(url)

			if err == nil {
				err = json.Unmarshal(respData, &prices)

				check(err)

				url = prices.NextPageLink

				observability.Logger("Debug", fmt.Sprintf("Successfully unmarshalled prices: %d %s", prices.Count, prices.NextPageLink))

				if lastNextPageLink == url {
					break
				} else {
					lastNextPageLink = url
				}

			}

			if err == nil {
				WritePriceOutput(file, prices)
			}

			if err != nil || url == "" {
				break
			}

		}
	}
}
