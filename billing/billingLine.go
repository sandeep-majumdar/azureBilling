package billing

import (
	"fmt"
	"strconv"

	"github.com/adeturner/observability"
)

func (bl *BillingLine) print() {
	observability.Logger("Info", fmt.Sprintf("%v", bl))
}

func (bl *BillingLine) SetValues(records []string) {
	var err error

	bl.InvoiceSectionName = records[0]
	bl.AccountName = records[1]
	bl.AccountOwnerId = records[2]
	bl.SubscriptionId = records[3]
	bl.SubscriptionName = records[4]
	bl.ResourceGroup = records[5]
	bl.ResourceLocation = records[6]
	bl.Date = records[7]
	bl.ProductName = records[8]
	bl.MeterCategory = records[9]

	bl.MeterSubCategory = records[10]
	bl.MeterId = records[11]
	bl.MeterName = records[12]
	bl.MeterRegion = records[13]
	bl.UnitOfMeasure = records[14]

	// bl.Quantity = records[15]
	bl.Quantity, err = strconv.ParseFloat(records[15], 64)
	if err != nil {
		observability.Logger("Error", fmt.Sprintf("Failed to parse bl.Quantity, err = strconv.ParseFloat(records[15], 64) from %s", records[15]))
	}

	bl.EffectivePrice = records[16]

	//bl.CostInBillingCurrency = records[17]
	bl.CostInBillingCurrency, err = strconv.ParseFloat(records[17], 64)
	if err != nil {
		observability.Logger("Error", fmt.Sprintf("Failed to parse bl.CostInBillingCurrency, err = strconv.ParseFloat(records[17], 64) from %s", records[17]))
	}

	bl.CostCenter = records[18]
	bl.ConsumedService = records[19]

	bl.ResourceId = records[20]
	bl.Tags = records[21]
	bl.OfferId = records[22]
	bl.AdditionalInfo = records[23]
	bl.ServiceInfo1 = records[24]
	bl.ServiceInfo2 = records[25]
	bl.ResourceName = records[26]
	bl.ReservationId = records[27]
	bl.ReservationName = records[28]
	bl.UnitPrice = records[29]
	bl.ProductOrderId = records[30]

	bl.ProductOrderName = records[31]
	bl.Term = records[32]
	bl.PublisherType = records[33]
	bl.PublisherName = records[34]
	bl.ChargeType = records[35]
	bl.Frequency = records[36]
	bl.PricingModel = records[37]
	bl.AvailabilityZone = records[38]
	bl.BillingAccountId = records[39]
	bl.BillingAccountName = records[40]

	bl.BillingCurrencyCode = records[41]
	bl.BillingPeriodStartDate = records[42]
	bl.BillingPeriodEndDate = records[43]
	bl.BillingProfileId = records[44]
	bl.BillingProfileName = records[45]
	bl.InvoiceSectionId = records[46]
	bl.IsAzureCreditEligible = records[47]
	bl.PartNumber = records[48]
	bl.PayGPrice = records[49]
	bl.PlanName = records[50]

	bl.ServiceFamily = records[51]
}
