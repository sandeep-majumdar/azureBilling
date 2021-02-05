package azureBilling

import (
	"fmt"
	"strconv"

	"github.com/adeturner/observability"
)

func (vmli *vmLookupItem) setValues(i []string) {

	var err error

	vmli.VM = i[0]

	vmli.Cores, err = strconv.Atoi(i[1])
	if err != nil {
		observability.Logger("Error", fmt.Sprintf("Failed to parse vmli.Cores, err = strconv.Atoi(i[1]) from %s", i[1]))
	}

	vmli.MemGB, err = strconv.ParseFloat(i[2], 64)
	if err != nil {
		observability.Logger("Error", fmt.Sprintf("Failed to parse vmli.MemGB, err = strconv.ParseFloat(i[2], 64) from %s", i[2]))
	}

}
