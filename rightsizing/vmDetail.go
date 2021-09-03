package rightsizing

import (
	"fmt"
	"strings"
)

type vmDetail struct {
	ResourceId       string
	Portfolio        string
	Platform         string
	ProductCode      string
	EnvironmentType  string
	ResourceLocation string
	MeterName        string
}

func (vd *vmDetail) getCSVRow() []byte {
	return []byte(fmt.Sprintf("\"%s\",\"%s\",\"%s\",\"%s\",\"%s\",\"%s\",\"%s\"\n",
		vd.ResourceId, vd.Portfolio, vd.Platform, vd.ProductCode, vd.EnvironmentType, vd.ResourceLocation, vd.MeterName))
}

func (vd *vmDetail) getSubscription() string {
	s := strings.Split(vd.ResourceId, "/")
	return s[1]
}
