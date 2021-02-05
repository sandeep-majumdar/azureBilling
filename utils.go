package azureBilling

import (
	"fmt"
	"os"
	"time"

	"github.com/adeturner/observability"
)

// expects DD/MM/YYYY
func dateStrToTime(dateStr string) (time.Time, error) {

	layout := "2006-01-02T15:04:05.000Z"

	str := fmt.Sprintf("%s-%s-%sT00:00:00.000Z", dateStr[6:10], dateStr[3:5], dateStr[0:2])

	t, err := time.Parse(layout, str)

	if err != nil {
		observability.Logger("Error", fmt.Sprintf("Failed to convert date from %s", str))
	} else {
		// observability.Logger("Info", fmt.Sprintf("found date %v from %s", t, str))
	}

	return t, err
}

func mod(a, b int) int {
	m := a % b
	if a < 0 && b < 0 {
		m -= b
	}
	if a < 0 && b > 0 {
		m += b
	}
	return m
}

func FileExists(fileLocation string) bool {

	retval := true

	f, err := os.Open(fileLocation)

	if err != nil {
		observability.Logger("Error", fmt.Sprintf("Unable to read input file=%s err=%s", fileLocation, err))
		retval = false
	} else {
		observability.Logger("Info", fmt.Sprintf("Successfully found file=%s", fileLocation))
	}

	defer f.Close()

	return retval
}
