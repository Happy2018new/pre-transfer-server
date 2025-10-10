package minecraft

import (
	"encoding/json"
	"fmt"
	"os"
	"sync"
)

var mu *sync.RWMutex
var record map[string]bool

func init() {
	mu = new(sync.RWMutex)
	record = make(map[string]bool)

	fileBytes, err := os.ReadFile("record.json")
	if err != nil {
		return
	}

	err = json.Unmarshal(fileBytes, &record)
	if err != nil {
		return
	}
}

// markResourcePackDownloaded ..
func markResourcePackDownloaded(playerIdentity string) error {
	mu.Lock()
	defer mu.Unlock()

	record[playerIdentity] = true

	fileBytes, err := json.Marshal(record)
	if err != nil {
		return fmt.Errorf("markResourcePackDownloaded: %v", err)
	}

	err = os.WriteFile("record.json", fileBytes, 0600)
	if err != nil {
		return fmt.Errorf("markResourcePackDownloaded: %v", err)
	}

	return nil
}

// checkResourcePackDownloaded ..
func checkResourcePackDownloaded(playerIdentity string) bool {
	mu.RLock()
	defer mu.RUnlock()
	return record[playerIdentity]
}
