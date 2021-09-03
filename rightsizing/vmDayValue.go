package rightsizing

import "fmt"

type vmDayValue struct {
	ResourceId            string
	Datestr               string
	UnitOfMeasure         string
	Quantity              float64
	EffectivePrice        float64
	CostInBillingCurrency float64
	UnitPrice             float64
	// TODO add timeseries here?
}

func (vs *vmDayValue) getCSVRow() []byte {
	return []byte(fmt.Sprintf("\"%s\",\"%s\",\"%s\",\"%f\",\"%f\",\"%f\",\"%f\"\n",
		vs.ResourceId, vs.Datestr, vs.UnitOfMeasure, vs.Quantity, vs.EffectivePrice, vs.CostInBillingCurrency, vs.UnitPrice))

}
