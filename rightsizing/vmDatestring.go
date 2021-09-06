package rightsizing

import (
	"time"
)

type vmDatestring string

/*
	MS provides in American date format
	08/14/2021
	0123456789
*/

func (vds vmDatestring) Day() string {
	s := string(vds)
	return s[3:5]
}

func (vds vmDatestring) Month() string {
	s := string(vds)
	return s[0:2]
}

func (vds vmDatestring) Year() string {
	s := string(vds)
	return s[6:10]
}

func (vds vmDatestring) DatestrToDayOfWeek() (string, error) {
	var t time.Time
	layout := time.RFC3339 // "2006-01-02T15:04:05Z07:00"
	str := vds.Year() + "-" + vds.Month() + "-" + vds.Day() + "T00:00:00Z"
	//
	t, err := time.Parse(layout, str)
	// observability.Info(fmt.Sprintf("%s => %s ==== %v", t.Weekday(), str, t))
	return t.Weekday().String(), err
}

func (vds vmDatestring) Match(ts string) (ok bool) {

	// 08/13/2021 2021-08-23T01:00:00+00:00
	// 0123456789
	ok = vds.CanonicalFormat() == ts[0:10]
	if ok {
		//observability.Info(fmt.Sprintf("%s => %s", vds.CanonicalFormat(), ts[0:10]))
	}
	return ok
}

func (vds vmDatestring) CanonicalFormat() string {
	return vds.Year() + "-" + vds.Month() + "-" + vds.Day()
}
