package rightsizing

import (
	"fmt"
	"io"
	"os"

	"github.com/adeturner/azureBilling/config"
	"github.com/adeturner/azureBilling/observability"
	"github.com/adeturner/azureBilling/utils"
)

type vmAggregateMap struct {
	vmMap map[string]*vmAggregateOutput
}

func NewVmAggregateMap(details *vmDetails, dayResourceMap *vmDayResourceMap, monitorMetrics *vmMonitorMetrics) *vmAggregateMap {

	var key string

	agg := &vmAggregateMap{}
	agg.vmMap = make(map[string]*vmAggregateOutput)
	for resourceId, detail := range details.vmMap {

		//observability.Info(fmt.Sprintf("dayMap len=%d", len(dayResourceMap.vmMap[resourceId].dayMap)))

		for datestr, dayValue := range dayResourceMap.vmMap[resourceId].dayMap {

			key = string(resourceId) + "/" + datestr.CanonicalFormat()

			agg.vmMap[key] = NewVmAggregateOutput(datestr, detail, dayValue)

			if monitorMetrics.VmMap[vmResourceId(resourceId)] != nil {

				o := *monitorMetrics.VmMap[vmResourceId(resourceId)]

				if o.Observations != nil {
					//observability.Info(fmt.Sprintf("non-nil observemap"))
					agg.vmMap[key].setMetrics(datestr, &o)
				}
			}
		}
	}
	observability.Info(fmt.Sprintf("%d aggregate records loaded\n", len(agg.vmMap)))
	observability.LogMemory("Info")

	return agg
}

func (vag *vmAggregateMap) WriteFile(filename string) error {
	observability.Info("Writing aggregates to file " + filename)
	var fs utils.FileSystem = utils.LocalFS{}
	file, err := fs.Create(filename)
	if err != nil {
		observability.Logger("Error", fmt.Sprintf("Failed to open file: %v", err))
	}
	defer file.Close()
	if err == nil {
		err = vag.WriteCSVHeader(file)
	}
	if err == nil {
		err = vag.WriteCSVOutput(file)
	}
	return err
}

func (vag *vmAggregateMap) FileExists(filename string) bool {
	retval := true
	f, err := os.Open(filename)
	if err != nil {
		//observability.Logger("Error", fmt.Sprintf("Unable to read input file=%s err=%s", filename, err))
		retval = false
	} else {
		//observability.Logger("Info", fmt.Sprintf("Successfully found file=%s", filename))
	}
	defer f.Close()
	return retval
}

func (vag *vmAggregateMap) getCSVHeader() (retval []byte) {

	header := fmt.Sprintf(
		"\"ResourceId\"," +
			"\"Portfolio\"," +
			"\"Platform\"," +
			"\"ProductCode\"," +
			"\"EnvironmentType\"," +
			"\"ResourceLocation\"," +
			"\"MeterName\"," +
			"\"dayOfWeek\"," +
			"\"Datestr\"," +
			"\"UnitOfMeasure\"," +
			"\"Quantity\"," +
			"\"EffectivePrice\"," +
			"\"CostInBillingCurrency\"," +
			"\"UnitPrice\"," +
			"\"ErrorString\"," +
			"\"MeasureType\"," +
			"\"ObservationType\",")

	var hr string
	var detail string
	for i := 0; i < 24; i++ {
		hr = utils.HourString(i)
		detail = detail + fmt.Sprintf("\"%s\",", hr)
	}

	return []byte(header + detail + "\n")

}

func (vag *vmAggregateMap) WriteCSVHeader(w io.Writer) (err error) {
	observability.Info("Writing CSV header")
	_, err = w.Write(vag.getCSVHeader())
	return err
}

func (vag *vmAggregateMap) WriteCSVOutput(w io.Writer) (err error) {
	observability.Info(fmt.Sprintf("Writing CSV output from len(map)=%d records", len(vag.vmMap)))
	observability.LogMemory("Info")

	for _, v := range vag.vmMap {
		for _, mt := range MeasureTypeList {
			for _, ot := range observationTypeList {
				if ot == AVG || ot == MIN || ot == MAX {
					_, err = w.Write(v.getCSVRow(mt, ot))
					if err != nil {
						break
					}
				}
			}
		}
	}

	observability.LogMemory("Info")
	return err
}

func (vag *vmAggregateMap) getOutputFile() string {
	return config.ConfigMap.WorkingDirectory + "/" + config.ConfigMap.OutputVmAggregateCSVFile
}
