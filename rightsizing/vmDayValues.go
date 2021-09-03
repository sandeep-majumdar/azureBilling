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

type vdvConcatKey string

type vmDayValues struct {
	vmMap map[vdvConcatKey]*vmDayValue
}

func NewVmDayValues() (vdv *vmDayValues, err error) {
	vdv = &vmDayValues{}
	vdv.vmMap = make(map[vdvConcatKey]*vmDayValue)
	return vdv, err
}

func (vdv *vmDayValues) addValue(i vmDayValue) (err error) {
	k := vdv.getVmmConcatKey(i.ResourceId, i.Datestr)
	vdv.vmMap[k] = &i
	return err
}

func (vdv *vmDayValues) Add(resourceId, datestr, effectivePrice, unitPrice string, quantity, costInBillingCurrency float64, unitOfMeasure string) (err error) {

	k := vdv.getVmmConcatKey(resourceId, datestr)

	var m *vmDayValue
	n, ok := vdv.Get(resourceId, datestr)

	if !ok {
		m = &vmDayValue{}
		m.ResourceId = resourceId
		m.Datestr = datestr
		m.Quantity = quantity
		m.CostInBillingCurrency = costInBillingCurrency
		m.UnitOfMeasure = unitOfMeasure

		m.EffectivePrice, err = strconv.ParseFloat(effectivePrice, 64)
		if err == nil {
			m.UnitPrice, err = strconv.ParseFloat(unitPrice, 64)
		}

		if err == nil {
			vdv.vmMap[k] = m
		}

	} else {
		n.CostInBillingCurrency = n.CostInBillingCurrency + costInBillingCurrency
		n.Quantity = n.Quantity + quantity
		//observability.Error(fmt.Sprintf("Unexpected found existing vmDayValue with datestr=%s, quantity=%f cost=%f resourceId=%s", datestr, n.Quantity, n.CostInBillingCurrency, resourceId))
	}

	return err
}

func (vdv *vmDayValues) getVmmConcatKey(resourceId, datestr string) vdvConcatKey {
	return vdvConcatKey(resourceId + "/" + datestr)
}

func (vdv *vmDayValues) Get(resourceId, datestr string) (v *vmDayValue, ok bool) {
	k := vdv.getVmmConcatKey(resourceId, datestr)
	v, ok = vdv.vmMap[k]
	return v, ok
}

func (vdv *vmDayValues) ReadFile(filename string) error {

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
				i := vmDayValue{}
				i.ResourceId = record[0]
				i.Datestr = record[1]
				i.UnitOfMeasure = record[2]
				i.Quantity, err = strconv.ParseFloat(record[3], 64)
				i.EffectivePrice, err = strconv.ParseFloat(record[4], 64)
				i.CostInBillingCurrency, err = strconv.ParseFloat(record[5], 64)
				i.UnitPrice, err = strconv.ParseFloat(record[6], 64)

				// meterLookup is a global variable declared in AzurePriceMeter.go
				vdv.addValue(i)
			}
		}
	}
	observability.Info(fmt.Sprintf("Loaded %d vmDayValues from file %s", len(vdv.vmMap), filename))
	observability.LogMemory("Info")
	return err

}

func (vdv *vmDayValues) WriteFile(filename string) error {
	var fs utils.FileSystem = utils.LocalFS{}
	file, err := fs.Create(filename)
	if err != nil {
		observability.Logger("Error", fmt.Sprintf("Failed to open file: %v", err))
	}
	defer file.Close()
	if err == nil {
		err = vdv.WriteCSVHeader(file)
	}
	if err == nil {
		err = vdv.WriteCSVOutput(file)
	}
	return err
}

func (vdv *vmDayValues) FileExists(filename string) bool {
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

func (vdv *vmDayValues) getCSVHeader() []byte {
	return []byte(fmt.Sprintf("\"%s\",\"%s\",\"%s\",\"%s\",\"%s\",\"%s\",\"%s\",\"%s\"\n",
		"MapKey", "ResourceId", "Date", "UnitOfMeasure", "Quantity", "EffectivePrice", "CostInBillingCurrency", "UnitPrice"))
}

func (vdv *vmDayValues) WriteCSVHeader(w io.Writer) (err error) {
	_, err = w.Write(vdv.getCSVHeader())
	return err
}

func (vdv *vmDayValues) WriteCSVOutput(w io.Writer) (err error) {
	var b []byte
	for k, v := range vdv.vmMap {
		b = []byte("\"" + k + "\"" + ",")
		b = append(b, v.getCSVRow()...)
		b = append(b)
		_, err = w.Write(b)
		if err != nil {
			break
		}
	}
	return err
}

/*
func (vdv *vmDayValues) queryVmmConcatKey(val vdvConcatKey) (subscriptionid, resourcename, datestr string, err error) {
	s := strings.Split(string(val), "/")
	subscriptionid = s[0]
	resourcename = s[1]
	datestr = s[2]
	return subscriptionid, resourcename, datestr, err
}
*/
