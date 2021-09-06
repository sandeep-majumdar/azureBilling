package rightsizing

import (
	"fmt"

	"github.com/adeturner/azureBilling/observability"
)

type vmMonitorMetric struct {
	ResourceId   string          `json:"resourceId"`
	Observations *ObservationMap `json:"observations"`
	ErrorString  string          `json:"errorStr"`
}

func (vmmt *vmMonitorMetric) setErrorString(errStr string) {
	vmmt.ErrorString = errStr
}

func NewVmMonitorMetric(resourceid string, errStr string, om *ObservationMap) *vmMonitorMetric {
	m := &vmMonitorMetric{}
	m.ResourceId = resourceid
	m.ErrorString = errStr
	m.Observations = om
	return m
}

func (vmmt *vmMonitorMetric) print() {
	observability.Info(fmt.Sprintf("%v", vmmt))
	if vmmt.Observations != nil {
		vmmt.Observations.print()
	}
}
