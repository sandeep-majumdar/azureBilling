package azureBilling

import (
	"fmt"

	"github.com/adeturner/observability"
)

func (aggrg *aggregateResourceGroup) print(cnt int) {

	i := 0

	for k, v := range aggrg.items {
		observability.Logger("Info", fmt.Sprintf("%s -> %v\n", k, v))
		i++
		if i > cnt {
			break
		}
	}
}

func (aggrg *aggregateResourceGroup) printCount() {
	observability.Logger("Info", fmt.Sprintf("managedDiskLookup has %d records\n", len(aggrg.items)))
}

func (aggrg *aggregateResourceGroup) init() {
	aggrg.items = make(map[string]aggregateResourceGroupItem)
}
