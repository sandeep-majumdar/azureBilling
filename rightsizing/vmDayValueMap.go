package rightsizing

type vmDayValueMap struct {
	dayMap map[vmDatestring]*vmDayValue
}

func NewVmDayValueMap() *vmDayValueMap {
	vdvm := &vmDayValueMap{}
	vdvm.dayMap = make(map[vmDatestring]*vmDayValue)
	return vdvm
}

func (vdvm *vmDayValueMap) Get(datestr string) (v *vmDayValue, ok bool) {
	v, ok = vdvm.dayMap[vmDatestring(datestr)]
	return v, ok
}
