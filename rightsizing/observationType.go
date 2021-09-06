package rightsizing

type observationType int

const (
	AVG observationType = iota
	MAX
	MIN
	COUNT
	TOTAL
)

func (ot observationType) String() string {
	return [...]string{
		"Avg",
		"Max",
		"Min",
		"Count",
		"Total",
	}[ot]
}

var observationTypeList = []observationType{
	AVG,
	MAX,
	MIN,
	COUNT,
	TOTAL,
}
