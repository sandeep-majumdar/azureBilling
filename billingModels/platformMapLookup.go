package billingModels

import (
	"encoding/csv"
	"fmt"
	"io"
	"os"
	"strings"

	"github.com/adeturner/observability"
)

func (pml *platformMapLookup) print(cnt int) {

	i := 0

	for k, v := range pml.Items {
		observability.Logger("Info", fmt.Sprintf("%s -> %v\n", k, v))
		i++
		if i > cnt {
			break
		}
	}
}

func (pml *platformMapLookup) printCount() {
	observability.Logger("Info", fmt.Sprintf("PlatformMapLookup has %d records\n", len(pml.Items)))
}

func (pml *platformMapLookup) init() {
	pml.Items = make(map[string]PlatformMapLookupItem)
}

func (pml *platformMapLookup) Read(fileLocation string) error {

	pml.init()

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
				i := PlatformMapLookupItem{}
				i.SetValues(record)

				key = pml.getKey(i.SubscriptionId, i.RgName)
				pml.Items[key] = i
				key = pml.getKey(i.SubscriptionId, "")
				pml.Items[key] = i

			}
		}
	}

	observability.LogMemory("Info")
	pml.printCount()
	// pml.print(10)

	return err

}

func (pml *platformMapLookup) getKey(subscriptionId, rgName string) string {
	return strings.ToLower(fmt.Sprintf(":%s:%s:", subscriptionId, rgName))
}

func (pml *platformMapLookup) Get(subscriptionId, rgName string) (PlatformMapLookupItem, bool) {

	key := pml.getKey(subscriptionId, rgName)

	//if subscriptionId == "ad88d8c8-5739-4619-b8dd-4cab5fd3c075" && rgName == "DATABRICKS-RG-ADBAZEWTDATALAKEPLATFORM-6PRSSGZRQWVLC" {
	//	observability.Logger("Error", fmt.Sprintf("Debugging key=%s", key))
	//}

	pmli, ok := pml.Items[key]
	if !ok {

		// some rgNames like databricks are dynamic, so just try to match the subscriptionid
		l := len(subscriptionId)

		//if subscriptionId == "ad88d8c8-5739-4619-b8dd-4cab5fd3c075" && rgName == "DATABRICKS-RG-ADBAZEWTDATALAKEPLATFORM-6PRSSGZRQWVLC" {
		//	observability.Logger("Debug", fmt.Sprintf("Debugging key=%s failed once, subscriptionLen=%d", key, l))
		//}

		ok2 := false
		for k, v := range pml.Items {

			if k[1:l+1] == subscriptionId {
				pmli, ok2 = v, true

				//if subscriptionId == "ad88d8c8-5739-4619-b8dd-4cab5fd3c075" && rgName == "DATABRICKS-RG-ADBAZEWTDATALAKEPLATFORM-6PRSSGZRQWVLC" {
				//	observability.Logger("Debug", fmt.Sprintf("Debugging pmli=%s", v))
				//}

				break
			} else {
				//if subscriptionId == "ad88d8c8-5739-4619-b8dd-4cab5fd3c075" && rgName == "DATABRICKS-RG-ADBAZEWTDATALAKEPLATFORM-6PRSSGZRQWVLC" {
				//	if strings.Contains(k, "ad88d8c8-5739-4619-b8dd-4cab5fd3c075") {
				//		observability.Logger("Debug", fmt.Sprintf("Debugging key=%s failed once, k=%s, k[1:l+1]=%s", key, k, k[1:l+1]))
				//	}
				//}
			}
		}

		ok = ok2

		/*
			if !ok2 {
				if subscriptionId == "ad88d8c8-5739-4619-b8dd-4cab5fd3c075" && rgName == "DATABRICKS-RG-ADBAZEWTDATALAKEPLATFORM-6PRSSGZRQWVLC" {
					observability.Logger("Debug", fmt.Sprintf("Debugging key=%s failed totally", key))
				}
				//observability.Logger("Error", fmt.Sprintf("Unable to find PlatformMapLookupItem for key=%s", key))
			}
		*/

	}

	return pmli, ok
}
