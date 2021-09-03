package rightsizing

type vmMonitorMetric struct {
	ResourceId  string          `json:"resourceId"`
	ObserveMap  *ObservationMap `json:"observeMap"`
	ErrorString string          `json:"errorStr"`
}

func (vmmt *vmMonitorMetric) setErrorString(errStr string) {
	vmmt.ErrorString = errStr
}

func NewVmMonitorMetric(resourceid string, errStr string, om *ObservationMap) *vmMonitorMetric {
	m := &vmMonitorMetric{}
	m.ResourceId = resourceid
	m.ErrorString = errStr
	m.ObserveMap = om
	return m
}
