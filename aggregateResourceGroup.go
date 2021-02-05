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

func (aggrg *aggregateResourceGroup) add(reportingCategory, reportingSubCategory, portfolio, platform string, unitOfMeasure string, quantity float64, l billingLine) {

	key := fmt.Sprintf("%s:%s:%s:%s:%s:%s", reportingCategory, reportingSubCategory, l.SubscriptionId, l.ResourceGroup, l.MeterId, unitOfMeasure)

	// initializes two variables - api will receive either the value of "key" from the map
	// or a "zero value" (in this case the empty string)
	// ok will receive a bool that will be set to true if "key" was actually present in the map
	// evaluates ok, which will be true if "key" was in the map
	if _, ok := aggrg.items[key]; !ok {

		// if not found initialise
		argi := aggregateResourceGroupItem{}
		argi.Platform = platform
		argi.portfolio = portfolio
		argi.reportingCategory = reportingCategory
		argi.reportingSubCategory = reportingSubCategory
		argi.UnitOfMeasure = unitOfMeasure
		argi.CostInBillingCurrency = 0
		argi.Quantity = 0
		argi.SubscriptionId = l.SubscriptionId
		argi.ResourceGroup = l.ResourceGroup
		argi.MeterId = l.MeterId
		argi.ResourceLocation = l.ResourceLocation
		argi.ProductName = l.ProductName
		argi.MeterCategory = l.MeterCategory
		argi.MeterSubCategory = l.MeterSubCategory
		argi.MeterName = l.MeterName
		argi.MeterRegion = l.MeterRegion
		argi.EffectivePrice = l.EffectivePrice
		argi.CostCenter = l.CostCenter
		argi.ConsumedService = l.ConsumedService
		argi.ResourceId = l.ResourceId
		argi.ReservationId = l.ReservationId
		argi.Term = l.Term
		argi.UnitPrice = l.UnitPrice

		aggrg.items[key] = &argi
	}

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
		"\"%s\",\"%s\",\"%s\",\"%s\",\"%s\",\"%s\",\"%s\",\"%s\",\"%s\",\"%s\",\"%s\",\"%s\",\"%s\",\"%s\",\"%s\",\"%s\",\"%s\",\"%s\",\"%s\",\"%s\",\"%s\",\"%s\",\"%s\"\n",
		"reportingCategory", "reportingSubCategory", "Portfolio", "Platform", "SubscriptionId", "ResourceGroup", "MeterId", "UnitOfMeasure",
		"ProductName", "ResourceLocation", "MeterCategory", "MeterSubCategory", "MeterName", "MeterRegion",
		"EffectivePrice", "CostCenter", "ConsumedService", "ResourceId", "ReservationId", "Term", "Quantity", "UnitPrice", "CostInBillingCurrency"))
}

func (aggrg *aggregateResourceGroup) WriteCSVHeader(w io.Writer) {

	_, err := w.Write(aggrg.getCSVHeader())
	aggrg.check(err)
}

func (aggrg *aggregateResourceGroup) WriteCSVOutput(w io.Writer) {

	var csvRow string

	for _, v := range aggrg.items {

		csvRow = fmt.Sprintf(
			"\"%s\",\"%s\",\"%s\",\"%s\",\"%s\",\"%s\",\"%s\",\"%s\",\"%s\",\"%s\",\"%s\",\"%s\",\"%s\",\"%s\",\"%s\",\"%s\",\"%s\",\"%s\",\"%s\",\"%s\",\"%f\",\"%s\",\"%f\"\n",
			v.reportingCategory, v.reportingSubCategory, v.portfolio, v.Platform, v.SubscriptionId, v.ResourceGroup, v.MeterId, v.UnitOfMeasure,
			v.ProductName, v.ResourceLocation, v.MeterCategory, v.MeterSubCategory, v.MeterName, v.MeterRegion,
			v.EffectivePrice, v.CostCenter, v.ConsumedService, v.ResourceId, v.ReservationId, v.Term, v.Quantity, v.UnitPrice, v.CostInBillingCurrency)

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
