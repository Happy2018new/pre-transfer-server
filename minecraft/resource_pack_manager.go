package minecraft

import (
	"encoding/json"
	"fmt"
	"os"
	"sync"
)

var mutex *sync.RWMutex
var record map[string]map[string]bool

func init() {
	mutex = new(sync.RWMutex)
	record = make(map[string]map[string]bool)

	fileBytes, err := os.ReadFile("record.json")
	if err != nil {
		return
	}

	err = json.Unmarshal(fileBytes, &record)
	if err != nil {
		return
	}
}

// GetPackDownloadStates ..
func GetPackDownloadStates(playerIdentity string, packIdentity string) bool {
	mutex.RLock()
	defer mutex.RUnlock()
	mapping, ok := record[playerIdentity]
	if !ok {
		return false
	}
	return mapping[packIdentity]
}

// SetPackDownloadStates ..
func SetPackDownloadStates(playerIdentity string, packIdentity string, states bool) error {
	mutex.Lock()
	defer mutex.Unlock()

	mapping, ok := record[playerIdentity]
	if !ok {
		mapping = make(map[string]bool)
	}
	mapping[packIdentity] = states
	record[playerIdentity] = mapping

	fileBytes, err := json.MarshalIndent(record, "", "\t")
	if err != nil {
		return fmt.Errorf("SetPackDownloadStates: %v", err)
	}

	err = os.WriteFile("record.json", fileBytes, 0600)
	if err != nil {
		return fmt.Errorf("SetPackDownloadStates: %v", err)
	}

	return nil
}
