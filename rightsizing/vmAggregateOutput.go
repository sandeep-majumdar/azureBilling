package rightsizing

import (
	"fmt"

	"github.com/adeturner/azureBilling/observability"
	"github.com/adeturner/azureBilling/utils"
)

type vmAggregateOutput struct {
	ResourceId       string
	Portfolio        string
	Platform         string
	ProductCode      string
	EnvironmentType  string
	ResourceLocation string
	MeterName        string
	//
	dayOfWeek             string
	Datestr               string
	UnitOfMeasure         string
	Quantity              float64
	EffectivePrice        float64
	CostInBillingCurrency float64
	UnitPrice             float64
	//
	ErrorString string
	AvailMemAvg map[vmAggregateHour]float64
	AvailMemMin map[vmAggregateHour]float64
	AvailMemMax map[vmAggregateHour]float64
	PctCPUAvg   map[vmAggregateHour]float64
	PctCPUMin   map[vmAggregateHour]float64
	PctCPUMax   map[vmAggregateHour]float64
}

func NewVmAggregateOutput(datestr vmDatestring, detail *vmDetail, dayValue *vmDayValue) *vmAggregateOutput {

	output := &vmAggregateOutput{}
	output.ResourceId = detail.ResourceId
	output.Portfolio = detail.Portfolio
	output.Platform = detail.Platform
	output.ProductCode = detail.ProductCode
	output.EnvironmentType = detail.EnvironmentType
	output.ResourceLocation = detail.ResourceLocation
	output.MeterName = detail.MeterName
	output.dayOfWeek, _ = vmDatestring(dayValue.Datestr).DatestrToDayOfWeek()
	output.Datestr = vmDatestring(dayValue.Datestr).CanonicalFormat()
	output.UnitOfMeasure = dayValue.UnitOfMeasure
	output.Quantity = dayValue.Quantity
	output.EffectivePrice = dayValue.EffectivePrice
	output.CostInBillingCurrency = dayValue.CostInBillingCurrency
	output.UnitPrice = dayValue.UnitPrice
	output.AvailMemAvg = make(map[vmAggregateHour]float64)
	output.AvailMemMin = make(map[vmAggregateHour]float64)
	output.AvailMemMax = make(map[vmAggregateHour]float64)
	output.PctCPUAvg = make(map[vmAggregateHour]float64)
	output.PctCPUMin = make(map[vmAggregateHour]float64)
	output.PctCPUMax = make(map[vmAggregateHour]float64)

	return output
}

func (vao *vmAggregateOutput) setMetrics(datestr vmDatestring, vmm *vmMonitorMetric) error {

	vao.ErrorString = vmm.ErrorString

	for k, v := range vmm.Observations.ObserveMap {

		if k == MetricName(MEASURETYPE_AVAILABLE_MEMORY.String()) {

			for _, j := range v.ObservationArray {
				if datestr.Match(j.TimeStamp) {
					key := vmAggregateHour(j.Hour())
					vao.AvailMemAvg[key] = j.avg()
					vao.AvailMemMin[key] = j.min()
					vao.AvailMemMax[key] = j.max()
					//observability.Info(fmt.Sprintf("%v", vao.AvailMemAvg[key]))
				}
			}

		} else if k == MetricName(MEASURETYPE_PRECENTAGE_CPU.String()) {
			for _, j := range v.ObservationArray {

				if datestr.Match(j.TimeStamp) {
					key := vmAggregateHour(j.Hour())
					vao.PctCPUAvg[key] = j.avg()
					vao.PctCPUMin[key] = j.min()
					vao.PctCPUMax[key] = j.max()
					//observability.Info(fmt.Sprintf("%s, %v", string(key), vao.PctCPUAvg[key]))
				}
			}

		} else {
			observability.Info("Ignoring metric name: " + string(k))
		}
	}

	return nil
}

func (vao *vmAggregateOutput) getCSVRow(mt MeasureType, ot observationType) []byte {

	header := fmt.Sprintf(
		"\"%s\",\"%s\",\"%s\",\"%s\",\"%s\",\"%s\",\"%s\",\"%s\",\"%s\",\"%s\",\"%f\",\"%f\",\"%f\",\"%f\",\"%s\",\"%s\",\"%s\",",
		vao.ResourceId,
		vao.Portfolio,
		vao.Platform,
		vao.ProductCode,
		vao.EnvironmentType,
		vao.ResourceLocation,
		vao.MeterName,
		vao.dayOfWeek,
		vao.Datestr,
		vao.UnitOfMeasure,
		vao.Quantity,
		vao.EffectivePrice,
		vao.CostInBillingCurrency,
		vao.UnitPrice,
		vao.ErrorString,
		mt.String(),
		ot.String())

	// 24 hours
	var detail string
	var hr vmAggregateHour
	for i := 0; i < 24; i++ {
		//observability.Info(fmt.Sprintf("%s %v %f", string(hr), vao.PctCPUMax, vao.PctCPUMax[hr]))
		hr = vmAggregateHour(utils.HourString(i))
		detail = detail + fmt.Sprintf("%f,", vao.getValForMetricType(mt, ot, hr))
	}

	return []byte(header + detail + "\n")
}

func (vao *vmAggregateOutput) getValForMetricType(mt MeasureType, ot observationType, hr vmAggregateHour) (retval float64) {

	if mt == MEASURETYPE_AVAILABLE_MEMORY {
		switch ot {
		case AVG:
			retval = vao.AvailMemAvg[hr] / 1024.0 / 1024.0 / 1024.0
		case MIN:
			retval = vao.AvailMemMin[hr] / 1024.0 / 1024.0 / 1024.0
		case MAX:
			retval = vao.AvailMemMax[hr] / 1024.0 / 1024.0 / 1024.0
		default:
			observability.Error("Invalid observationType")
		}

	} else if mt == MEASURETYPE_PRECENTAGE_CPU {
		switch ot {
		case AVG:
			retval = vao.PctCPUAvg[hr]
		case MIN:
			retval = vao.PctCPUMin[hr]
		case MAX:
			retval = vao.PctCPUMax[hr]
		default:
			observability.Error("Invalid observationType")
		}

	}
	return retval
}
