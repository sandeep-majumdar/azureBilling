package observability

import (
	"fmt"
	"sync"
	"time"
)

// MetricTypeID gives the type of the metric
type MetricTypeID uint8

/*
const (
	duration MetricTypeID = iota
	integer
	float
)
*/

// Metric type to contain observations
type Metric struct {
	//ID    MetricTypeID
	value float64
}

// Metrics type to contain observations
type Metrics struct {
	mutex sync.Mutex
	m     map[string]Metric
}

// Init the map
func (ms *Metrics) Init() {
	ms.m = make(map[string]Metric)
}

func (ms *Metrics) setKeyValue(key string, m Metric) {
	ms.mutex.Lock()
	ms.m[key] = m
	ms.mutex.Unlock()
}

// SetDuration blah
func (ms *Metrics) SetDuration(key string, d time.Duration) {
	m := Metric{value: d.Seconds()}
	ms.setKeyValue(key, m)
}

// SetInteger blah
func (ms *Metrics) SetInteger(key string, i int) {
	m := Metric{value: float64(i)}
	ms.setKeyValue(key, m)
}

// SetFloat blah
func (ms *Metrics) SetFloat(key string, f float64) {
	m := Metric{value: f}
	ms.setKeyValue(key, m)
}

// Dump values
func (ms *Metrics) Dump() {
	for k, v := range ms.m {
		Info(fmt.Sprintf("metric %s %f", k, v.value))
	}
}
