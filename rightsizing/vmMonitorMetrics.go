package rightsizing

import (
	"encoding/json"
	"fmt"
	"os"
	"sync"

	"github.com/adeturner/azureBilling/observability"
	"github.com/adeturner/azureBilling/utils"
)

type vmmConcatKey string

type vmMonitorMetrics struct {
	VmMap map[vmmConcatKey]*vmMonitorMetric `json:"vmMap"`
	mutex sync.RWMutex
}

func NewVmMonitorMetrics() (vmm *vmMonitorMetrics, err error) {
	vmm = &vmMonitorMetrics{}
	vmm.VmMap = make(map[vmmConcatKey]*vmMonitorMetric)
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

func (vmm *vmMonitorMetrics) getConcatKey(resourceId string) vmmConcatKey {
	return vmmConcatKey(resourceId)
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
		decoder := json.NewDecoder(f)
		err = decoder.Decode(vmm)
	}
	if err == nil {
		observability.Info(fmt.Sprintf("Loaded %d metric records from file %s", vmm.GetMapLen(), filename))
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
