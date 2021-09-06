package rightsizing

import (
	"encoding/csv"
	"fmt"
	"io"
	"os"

	"github.com/adeturner/azureBilling/observability"
	"github.com/adeturner/azureBilling/utils"
)

type vmDetails struct {
	vmMap map[vmResourceId]*vmDetail
}

func NewVmDetails() (vmd *vmDetails, err error) {
	vmd = &vmDetails{}
	vmd.vmMap = make(map[vmResourceId]*vmDetail)
	return vmd, err
}

func (vmd *vmDetails) addValue(i *vmDetail) (err error) {

	k := vmd.getvmResourceId(i.ResourceId)
	vmd.vmMap[k] = i
	return err
}

func (vmd *vmDetails) Add(resourceId, portfolio, platform, productCode, environmentType, resourceLocation, meterName string) (err error) {

	v, ok := vmd.Get(resourceId)
	if !ok {
		v = &vmDetail{}
	}
	v.ResourceId = resourceId
	v.Portfolio = portfolio
	v.Platform = platform
	v.ProductCode = productCode
	v.EnvironmentType = environmentType
	v.ResourceLocation = resourceLocation
	v.MeterName = meterName
	k := vmd.getvmResourceId(resourceId)
	vmd.vmMap[k] = v
	return err
}

func (vmd *vmDetails) getvmResourceId(resourceId string) vmResourceId {
	return vmResourceId(resourceId)
}

func (vmd *vmDetails) Get(resourceId string) (v *vmDetail, ok bool) {
	k := vmd.getvmResourceId(resourceId)
	v, ok = vmd.vmMap[k]
	return v, ok
}

func (vmd *vmDetails) ReadFile(filename string) error {

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
				i := &vmDetail{}
				if cnt == 2 {
					//observability.Info(fmt.Sprintf("%s %s %s", record[1], record[2], record[3]))
				}
				i.ResourceId = record[1]
				i.Portfolio = record[2]
				i.Platform = record[3]
				i.ProductCode = record[4]
				i.EnvironmentType = record[5]
				i.ResourceLocation = record[6]
				i.MeterName = record[7]
				vmd.addValue(i)
			}
		}
	}
	observability.Info(fmt.Sprintf("Loaded %d vmDetails from file %s", len(vmd.vmMap), filename))
	observability.LogMemory("Info")
	return err
}

func (vmd *vmDetails) WriteFile(filename string) error {
	var fs utils.FileSystem = utils.LocalFS{}
	file, err := fs.Create(filename)
	if err != nil {
		observability.Logger("Error", fmt.Sprintf("Failed to open file: %v", err))
	}
	defer file.Close()
	if err == nil {
		err = vmd.WriteCSVHeader(file)
	}
	if err == nil {
		err = vmd.WriteCSVOutput(file)
	}
	return err
}

func (vmd *vmDetails) FileExists(filename string) bool {
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

func (vmd *vmDetails) getCSVHeader() []byte {
	return []byte(fmt.Sprintf("\"%s\",\"%s\",\"%s\",\"%s\",\"%s\",\"%s\",\"%s\",\"%s\"\n",
		"MapKey", "ResourceId", "Portfolio", "Platform", "ProductCode", "EnvironmentType", "ResourceLocation", "MeterName"))
}

func (vmd *vmDetails) WriteCSVHeader(w io.Writer) (err error) {
	_, err = w.Write(vmd.getCSVHeader())
	return err
}

func (vmd *vmDetails) WriteCSVOutput(w io.Writer) (err error) {
	var b []byte
	for k, v := range vmd.vmMap {
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
func (vmd *vmDetails)  queryvmResourceId(val vmResourceId) (subscriptionid, resourcename, datestr string, err error) {
	s := strings.Split(string(val), "/")
	subscriptionid = s[0]
	resourcename = s[1]
	datestr = s[2]
	return subscriptionid, resourcename, datestr, err
}
*/
