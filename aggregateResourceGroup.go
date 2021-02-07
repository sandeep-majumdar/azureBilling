package azureBilling

import (
	"fmt"
	"io"

	"github.com/adeturner/observability"
)

func (aggrg *aggregateResourceGroup) print(cnt int) {

	i := 0

	for k, v := range aggrg.items {
		observability.Logger("Info", fmt.Sprintf("%s -> %v\n", k, v))
		i++
		if i > cnt {
			break
		}
	}
}

func (aggrg *aggregateResourceGroup) printCount() {
	observability.Logger("Info", fmt.Sprintf("managedDiskLookup has %d records\n", len(aggrg.items)))
}

func (aggrg *aggregateResourceGroup) init() {
	aggrg.items = make(map[string]*aggregateResourceGroupItem)
}

func (aggrg *aggregateResourceGroup) add(reportingCategory, reportingSubCategory, portfolio, platform, unitOfMeasure, summaryCategory, quantityDivisor string, summaryQuantity, quantity float64, l billingLine) {

	key := fmt.Sprintf("%s:%s:%s:%s:%s:%s", reportingCategory, reportingSubCategory, l.SubscriptionId, l.ResourceGroup, l.MeterId, unitOfMeasure)

	// initializes two variables - api will receive either the value of "key" from the map
	// or a "zero value" (in this case the empty string)
	// ok will receive a bool that will be set to true if "key" was actually present in the map
	// evaluates ok, which will be true if "key" was in the map
	if _, ok := aggrg.items[key]; !ok {

		// if not found initialise
		argi := aggregateResourceGroupItem{}
		// key fields
		argi.reportingCategory = reportingCategory
		argi.reportingSubCategory = reportingSubCategory
		argi.SubscriptionId = l.SubscriptionId
		argi.ResourceGroup = l.ResourceGroup
		argi.MeterId = l.MeterId
		argi.UnitOfMeasure = unitOfMeasure
		// core values
		argi.summaryCategory = summaryCategory
		argi.quantityDivisor = quantityDivisor
		argi.portfolio = portfolio
		argi.Platform = platform
		argi.summaryQuantity = 0
		argi.Quantity = 0
		argi.CostInBillingCurrency = 0
		// detail for deep dive
		argi.ResourceLocation = l.ResourceLocation
		argi.ProductName = l.ProductName
		argi.MeterCategory = l.MeterCategory
		argi.MeterSubCategory = l.MeterSubCategory
		argi.MeterName = l.MeterName
		argi.MeterRegion = l.MeterRegion
		argi.EffectivePrice = l.EffectivePrice
		argi.CostCenter = l.CostCenter
		argi.ConsumedService = l.ConsumedService
		argi.ReservationId = l.ReservationId
		argi.Term = l.Term
		argi.UnitPrice = l.UnitPrice

		aggrg.items[key] = &argi
	}

	aggrg.items[key].summaryQuantity += summaryQuantity
	aggrg.items[key].Quantity += quantity
	aggrg.items[key].CostInBillingCurrency += l.CostInBillingCurrency
}

/*
   #######################################################
   Below here is about producing output only
   #######################################################
*/

func (aggrg *aggregateResourceGroup) check(e error) {
	if e != nil {
		observability.Logger("Error", fmt.Sprintf("%v", e))
		panic(e)
	}
}

func (aggrg *aggregateResourceGroup) getCSVHeader() []byte {
	return []byte(fmt.Sprintf(
		"\"%s\",\"%s\",\"%s\",\"%s\",\"%s\",\"%s\",\"%s\",\"%s\",\"%s\",\"%s\",\"%s\",\"%s\",\"%s\",\"%s\",\"%s\",\"%s\",\"%s\",\"%s\",\"%s\",\"%s\",\"%s\",\"%s\",\"%s\",\"%s\",\"%s\"\n",
		"reportingCategory", "reportingSubCategory", "SubscriptionId", "ResourceGroup", "MeterId", "UnitOfMeasure",
		"SummaryCategory", "quantityDivisor", "Portfolio", "Platform",
		"summaryQuantity", "Quantity", "CostInBillingCurrency", "UnitPrice",
		"ResourceLocation", "ProductName", "MeterCategory", "MeterSubCategory", "MeterName", "MeterRegion",
		"EffectivePrice", "CostCenter", "ConsumedService", "ReservationId", "Term"))
}

func (aggrg *aggregateResourceGroup) WriteCSVHeader(w io.Writer) {

	_, err := w.Write(aggrg.getCSVHeader())
	aggrg.check(err)
}

func (aggrg *aggregateResourceGroup) WriteCSVOutput(w io.Writer) {

	var csvRow string

	for _, v := range aggrg.items {

		csvRow = fmt.Sprintf(
			"\"%s\",\"%s\",\"%s\",\"%s\",\"%s\",\"%s\",\"%s\",\"%s\",\"%s\",\"%s\",\"%f\",\"%f\",\"%f\",\"%s\",\"%s\",\"%s\",\"%s\",\"%s\",\"%s\",\"%s\",\"%s\",\"%s\",\"%s\",\"%s\",\"%s\"\n",
			v.reportingCategory, v.reportingSubCategory, v.SubscriptionId, v.ResourceGroup, v.MeterId, v.UnitOfMeasure,
			v.summaryCategory, v.quantityDivisor, v.portfolio, v.Platform,
			v.summaryQuantity, v.Quantity, v.CostInBillingCurrency, v.UnitPrice,
			v.ResourceLocation, v.ProductName, v.MeterCategory, v.MeterSubCategory, v.MeterName, v.MeterRegion,
			v.EffectivePrice, v.CostCenter, v.ConsumedService, v.ReservationId, v.Term)

		_, err := w.Write([]byte(csvRow))
		aggrg.check(err)
	}

}

func (aggrg *aggregateResourceGroup) WriteFile(filename string) {

	observability.Logger("Info", fmt.Sprintf("Writing to %s", filename))

	var fs fileSystem = localFS{}

	file, err := fs.Create(filename)
	if err != nil {
		observability.Logger("Error", fmt.Sprintf("Failed to open file: %v", err))
	}
	defer file.Close()

	aggrg.WriteCSVHeader(file)
	aggrg.WriteCSVOutput(file)

}
