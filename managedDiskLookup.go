package azureBilling

import (
	"encoding/csv"
	"fmt"
	"io"
	"os"
	"strings"

	"github.com/adeturner/observability"
)

func (mdl *managedDiskLookup) print(cnt int) {

	i := 0

	for k, v := range mdl.items {
		observability.Logger("Info", fmt.Sprintf("%s -> %v\n", k, v))
		i++
		if i > cnt {
			break
		}
	}
}

func (mdl *managedDiskLookup) printCount() {
	observability.Logger("Info", fmt.Sprintf("managedDiskLookup has %d records\n", len(mdl.items)))
}

func (mdl *managedDiskLookup) init() {
	mdl.items = make(map[string]managedDiskLookupItem)
}

func (mdl *managedDiskLookup) Read(fileLocation string) error {

	mdl.init()

	var key string

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
				i := managedDiskLookupItem{}
				i.setValues(record)

				key = mdl.getKey(i.MeterName)
				mdl.items[strings.ToLower(key)] = i
			}
		}
	}

	observability.LogMemory("Info")
	mdl.printCount()
	// mdl.print(10)

	return err

}

func (mdl *managedDiskLookup) getKey(meterName string) string {
	return strings.ToLower(fmt.Sprintf(":%s:", meterName))
}

func (mdl *managedDiskLookup) get(meterName string) (managedDiskLookupItem, bool) {

	key := mdl.getKey(meterName)

	mdli, ok := mdl.items[key]
	if !ok {
		observability.Logger("Error", fmt.Sprintf("Unable to find managedDiskLookupItem for key=%s", key))
	}

	return mdli, ok
}
