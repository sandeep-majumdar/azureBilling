package rightsizing

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"os"
	"sync"

	"github.com/adeturner/azureBilling/observability"
	"github.com/adeturner/azureBilling/utils"
)

type vmMonitorMetrics struct {
	VmMap map[vmResourceId]*vmMonitorMetric `json:"vmMap"`
	mutex sync.RWMutex
}

func NewVmMonitorMetrics() (vmm *vmMonitorMetrics, err error) {
	vmm = &vmMonitorMetrics{}
	vmm.VmMap = make(map[vmResourceId]*vmMonitorMetric)
	return vmm, err
}

func (vmm *vmMonitorMetrics) add(resourceId, errStr string, mt *azMonitorMetricsType) (err error) {
	m := NewVmMonitorMetric(resourceId, errStr, NewObservationsFromAzMonitor(mt))
	k := vmm.getConcatKey(resourceId)
	vmm.mutex.Lock()
	vmm.VmMap[k] = m
	vmm.mutex.Unlock()
	return err
}

func (vmm *vmMonitorMetrics) getConcatKey(resourceId string) vmResourceId {
	return vmResourceId(resourceId)
}

func (vmm *vmMonitorMetrics) WriteFile(filename string) error {
	var fs utils.FileSystem = utils.LocalFS{}
	file, err := fs.Create(filename)
	if err != nil {
		observability.Logger("Error", fmt.Sprintf("Failed to open file: %v", err))
	}
	defer file.Close()
	if err == nil {

		encoder := json.NewEncoder(file)
		encoder.Encode(vmm)
	}
	return err
}

func (vmm *vmMonitorMetrics) ReadFile(filename string) error {
	f, err := os.Open(filename)
	defer f.Close()

	if err == nil {
		file, _ := ioutil.ReadFile(filename)
		_ = json.Unmarshal([]byte(file), &vmm)
	}
	if err == nil {
		observability.Info(fmt.Sprintf("Loaded %d metric records from file %s", vmm.GetMapLen(), filename))
		observability.LogMemory("Info")
	}
	return err
}

func (vmm *vmMonitorMetrics) FileExists(filename string) bool {
	retval := true
	f, err := os.Open(filename)
	defer f.Close()
	if err != nil {
		retval = false
	}
	return retval
}

func (vmm *vmMonitorMetrics) GetMapLen() int {
	vmm.mutex.RLock()
	defer vmm.mutex.RUnlock()
	return len(vmm.VmMap)
}

func (vmm *vmMonitorMetrics) print(r vmResourceId) {
	observability.Info(fmt.Sprintf("%v", vmm.VmMap[r]))
	if vmm.VmMap[r].Observations != nil {
		vmm.VmMap[r].print()
	}
}
