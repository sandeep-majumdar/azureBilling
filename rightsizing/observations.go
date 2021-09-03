package rightsizing

type MetricName string

type ObservationMap struct {
	ObserveMap map[MetricName]*Observations `json:"ObserveMap"`
}

type Observations struct {
	Name             MetricName     `json:"name"`
	ObservationArray []*Observation `json:"observation"`
}

func NewObservationMap() *ObservationMap {
	om := &ObservationMap{}
	om.ObserveMap = make(map[MetricName]*Observations)
	return om
}

func (om *ObservationMap) NewObservations(metricName string) *Observations {
	obs := &Observations{}
	obs.Name = MetricName(metricName)
	om.ObserveMap[obs.Name] = obs
	return obs
}

func (os *Observations) add(o *Observation) {
	os.ObservationArray = append(os.ObservationArray, o)
}

func NewObservationsFromAzMonitor(azmmt *azMonitorMetricsType) (om *ObservationMap) {

	om = NewObservationMap()
	if azmmt != nil {
		for _, v := range azmmt.Value {
			obs := om.NewObservations(v.Name.Value)
			for _, ts := range v.Timeseries {
				for _, metricObservation := range ts.Data {
					obs.add(NewObservation(metricObservation.TimeStamp, metricObservation.getValuesArray()))
				}
			}
		}
	} else {
		//observability.Info("Nil metrics are apparently a thing?! TODO")
	}

	return om
}

func NewObservation(ts string, a []float64) (o *Observation) {
	return &Observation{
		TimeStamp:   ts,
		ValuesArray: a,
	}
}

/*
[timestamp, [avg, min, max, ]]

func (o *Observations) WorkloadAnalyse() {

	for i, v := range o.ObservationArray {

	}

}

func (o *Observations) CostAnalyse() {

}

*/
