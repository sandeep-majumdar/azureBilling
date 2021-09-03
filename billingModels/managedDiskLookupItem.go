package billingModels

import (
	"fmt"
	"strconv"

	"github.com/adeturner/observability"
)

func (mdli *managedDiskLookupItem) setValues(i []string) {

	var err error

	mdli.MeterName = i[0]

	mdli.SizeGB, err = strconv.Atoi(i[1])
	if err != nil {
		observability.Logger("Error", fmt.Sprintf("Failed to parse SizeGB.SizeGB, err = strconv.Atoi(i[1]) from %s", i[1]))
	}

}
