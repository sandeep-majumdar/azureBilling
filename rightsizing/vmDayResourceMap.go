package rightsizing

import (
	"encoding/csv"
	"fmt"
	"io"
	"os"
	"strconv"

	"github.com/adeturner/azureBilling/observability"
	"github.com/adeturner/azureBilling/utils"
)

type vmDayResourceMap struct {
	vmMap map[vmResourceId]*vmDayValueMap
}

func NewVmDayValueResourceMap() (vdrm *vmDayResourceMap, err error) {
	vdrm = &vmDayResourceMap{}
	vdrm.vmMap = make(map[vmResourceId]*vmDayValueMap)
	return vdrm, err
}

func (vdrm *vmDayResourceMap) incrementValue(dm *vmDayValueMap, datestr string, quantity, costInBillingCurrency float64) (err error) {
	n := dm.dayMap[vmDatestring(datestr)]
	n.CostInBillingCurrency = n.CostInBillingCurrency + costInBillingCurrency
	n.Quantity = n.Quantity + quantity
	return err
}

func (vdrm *vmDayResourceMap) Add(resourceId, datestr, effectivePrice, unitPrice string, quantity, costInBillingCurrency float64, unitOfMeasure string) (err error) {

	// does the vmDayResourceMapExist for this resourceid?
	dayValueMap, ok := vdrm.Get(resourceId)

	// if not, create it
	if !ok {
		//observability.Debug("creating new value map: " + resourceId)
		dayValueMap = NewVmDayValueMap()
		vdrm.vmMap[vmResourceId(resourceId)] = dayValueMap
	}

	// does the vmDayValue exist for this day map?
	vdv, ok2 := dayValueMap.Get(datestr)

	if !ok2 {
		//observability.Debug("creating new day value for datestr=" + datestr + " resource=" + resourceId)
		// if not, create it
		vdv = &vmDayValue{}
		vdv.ResourceId = resourceId
		vdv.Datestr = datestr
		vdv.Quantity = quantity
		vdv.CostInBillingCurrency = costInBillingCurrency
		vdv.UnitOfMeasure = unitOfMeasure

		vdv.EffectivePrice, err = strconv.ParseFloat(effectivePrice, 64)
		if err == nil {
			vdv.UnitPrice, err = strconv.ParseFloat(unitPrice, 64)
		}

		if err == nil {
			dayValueMap.dayMap[vmDatestring(datestr)] = vdv
		}

	} else {
		// if yes, update it
		//observability.Debug("Updating existing day value" + resourceId)
		err = vdrm.incrementValue(dayValueMap, datestr, quantity, costInBillingCurrency)
	}

	return err
}

func (vdrm *vmDayResourceMap) getResourceId(resourceId string) vmResourceId {
	return vmResourceId(resourceId)
}

func (vdrm *vmDayResourceMap) Get(resourceId string) (v *vmDayValueMap, ok bool) {
	v, ok = vdrm.vmMap[vmResourceId(resourceId)]
	return v, ok
}

func (vdrm *vmDayResourceMap) ReadFile(filename string) error {

	var resourceId, datestr, effectivePrice, unitPrice, unitOfMeasure string
	var quantity, costInBillingCurrency float64

	f, err := os.Open(filename)
	if err != nil {
		observability.Logger("Error", fmt.Sprintf("Unable to read input file=%s err=%s", filename, err))
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
				observability.Logger("Error", fmt.Sprintf("Unable to parse file as CSV; file=%s err=%s", filename, err))
				break
			}
			cnt++
			// skip the first row (header)
			if cnt > 1 {
				// first column is mapkey - ignore
				resourceId = record[1]
				datestr = record[2]
				unitOfMeasure = record[3]
				quantity, err = strconv.ParseFloat(record[4], 64)
				effectivePrice = record[5]
				costInBillingCurrency, err = strconv.ParseFloat(record[6], 64)
				unitPrice = record[7]

				vdrm.Add(resourceId, datestr, effectivePrice, unitPrice, quantity, costInBillingCurrency, unitOfMeasure)
			}
		}
	}
	observability.Info(fmt.Sprintf("Loaded %d vmDayValues from file %s", len(vdrm.vmMap), filename))
	observability.LogMemory("Info")
	return err

}

func (vdrm *vmDayResourceMap) WriteFile(filename string) error {
	var fs utils.FileSystem = utils.LocalFS{}
	file, err := fs.Create(filename)
	if err != nil {
		observability.Logger("Error", fmt.Sprintf("Failed to open file: %v", err))
	}
	defer file.Close()
	if err == nil {
		err = vdrm.WriteCSVHeader(file)
	}
	if err == nil {
		err = vdrm.WriteCSVOutput(file)
	}
	return err
}

func (vdrm *vmDayResourceMap) FileExists(filename string) bool {
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

func (vdrm *vmDayResourceMap) getCSVHeader() []byte {
	return []byte(fmt.Sprintf("\"%s\",\"%s\",\"%s\",\"%s\",\"%s\",\"%s\",\"%s\",\"%s\"\n",
		"MapKey", "ResourceId", "Date", "UnitOfMeasure", "Quantity", "EffectivePrice", "CostInBillingCurrency", "UnitPrice"))
}

func (vdrm *vmDayResourceMap) WriteCSVHeader(w io.Writer) (err error) {
	_, err = w.Write(vdrm.getCSVHeader())
	return err
}

func (vdrm *vmDayResourceMap) WriteCSVOutput(w io.Writer) (err error) {
	var b []byte
	for k, v := range vdrm.vmMap {
		for _, v2 := range v.dayMap {
			b = []byte("\"" + k + "\"" + ",")
			b = append(b, v2.getCSVRow()...)
			b = append(b)
			_, err = w.Write(b)
			if err != nil {
				break
			}
		}
	}
	return err
}
