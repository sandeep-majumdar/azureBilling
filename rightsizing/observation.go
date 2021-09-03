package rightsizing

type observationType int

const (
	AVG observationType = iota
	MAX
	MIN
	COUNT
	TOTAL
)

type Observation struct {
	TimeStamp   string    `json:"timeStamp"`
	ValuesArray []float64 // indexed by observationType
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
