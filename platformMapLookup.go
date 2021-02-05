package azureBilling

import (
	"encoding/csv"
	"fmt"
	"io"
	"os"

	"github.com/adeturner/observability"
)

func (pml *platformMapLookup) print(cnt int) {

	i := 0

	for k, v := range pml.items {
		observability.Logger("Info", fmt.Sprintf("%s -> %v\n", k, v))
		i++
		if i > cnt {
			break
		}
	}
}

func (pml *platformMapLookup) printCount() {
	observability.Logger("Info", fmt.Sprintf("platformMapLookup has %d records\n", len(pml.items)))
}

func (pml *platformMapLookup) init() {
	pml.items = make(map[string]platformMapLookupItem)
}

func (pml *platformMapLookup) Read(fileLocation string) error {

	pml.init()

	f, err := os.Open(fileLocation)
	if err != nil {
		observability.Logger("Error", fmt.Sprintf("Unable to read input file=%s err=%s", fileLocation, err))
	}
	defer f.Close()

	cnt := 0

	if err == nil {

		r := csv.NewReader(f)
		for {

			record, err := r.Read()
			if err == io.EOF {
				break
			} else if err != nil {
				observability.Logger("Error", fmt.Sprintf("Unable to parse file as CSV; file=%s err=%s", fileLocation, err))
				break
			}

			cnt++

			// skip the first row (header)
			if cnt > 1 {
				i := platformMapLookupItem{}
				i.setValues(record)

				pml.items[i.subscriptionId+"/"+i.rgName] = i
			}
		}
	}

	observability.LogMemory("Info")
	pml.printCount()
	// pml.print(10)

	return err

}

func (pml *platformMapLookup) get(subscriptionId, rgName string) (platformMapLookupItem, bool) {

	key := fmt.Sprintf("%s/%s", subscriptionId, rgName)

	pmli, ok := pml.items[key]
	if !ok {

		// some rgNames like databricks are dynamic, so just try to match the subscriptionid
		l := len(subscriptionId)

		ok2 := false
		for k, v := range pml.items {
			if k[:l] == subscriptionId {
				pmli, ok2 = v, true
			}
		}

		if !ok2 {
			// observability.Logger("Error", fmt.Sprintf("Unable to find platformMapLookupItem for key=%s", key))
		}

	}

	return pmli, ok
}
