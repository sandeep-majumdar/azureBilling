package rightsizing

type MeasureType int

const (
	MEASURETYPE_PRECENTAGE_CPU MeasureType = iota
	MEASURETYPE_AVAILABLE_MEMORY
)

func (m MeasureType) String() string {
	return [...]string{
		"Percentage CPU",
		"Available Memory Bytes",
	}[m]
}

var MeasureTypeList = []MeasureType{
	MEASURETYPE_PRECENTAGE_CPU,
	MEASURETYPE_AVAILABLE_MEMORY,
}
