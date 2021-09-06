package rightsizing

import (
	"fmt"

	"github.com/adeturner/azureBilling/observability"
)

type Observation struct {
	TimeStamp   string    `json:"timeStamp"`
	ValuesArray []float64 `json:"ValuesArray"` // indexed by observationType
}

func (o *Observation) print() {
	observability.Info(fmt.Sprintf("%v", o))
}

func (o *Observation) Hour() string {
	// 2021-08-01T00:00:00+00:00
	// 0123456789012
	//observability.Info(fmt.Sprintf("%s %s", o.TimeStamp, o.TimeStamp[11:13]))
	return o.TimeStamp[11:13]
}

func (o *Observation) avg() float64 {
	return o.ValuesArray[AVG]
}
func (o *Observation) max() float64 {
	return o.ValuesArray[MAX]
}
func (o *Observation) min() float64 {
	return o.ValuesArray[MIN]
}
func (o *Observation) count() float64 {
	return o.ValuesArray[AVG]
}
func (o *Observation) total() float64 {
	return o.ValuesArray[TOTAL]
}
