package db

import (
	"encoding/json"
	"fmt"
	"log"
	"os"
)

// JSON structure containing the URL of the watcher and the list of the seller_currency.
type WatcherURL struct {
	URL             string   `json:"url"`
	SellerCurrency []string `json:"seller_currency"`
}

// JSON structure containing the id of the item.
type ItemID struct {
	Id int `json:"id"`
}

// Loads teh content of the file filePath, appends the new items to the unmarshaled content and updates the file filePath.
// Returns error if reading, marshalling or writing fails.
func AppendWatcher(filePath string, watcher WatcherURL) error {
	if filePath == "" {
		filePath = "watchers.json"
	}

	// load the content of json file
	watchers, err := ReadWatchers(filePath)
	if err != nil {
		return fmt.Errorf("error reading watcherURL: %v", err)
	}

	// append the new watcher
	watchers = append(watchers, watcher)

	updatedContent, err := json.MarshalIndent(watchers, "", "  ")
	if err != nil {
		return fmt.Errorf("error marshalling watchers: %v", err)
	}

	if err := os.WriteFile(filePath, updatedContent, 0644); err != nil {
		return fmt.Errorf("error writing file while updating the json content: %v", err)
	}

	return nil
}

// Function changes the content of data parameter. Returns nil if file is empty/not found.
func readBytes(filePath string, data *[]byte) error {
	bytes, err := os.ReadFile(filePath)
	if os.IsNotExist(err) {
		return nil
	}

	if err != nil {
		return fmt.Errorf("error reading file: %v", err)
	}

	if len(bytes) == 0 {
		return nil
	}

	*data = bytes
	return nil
}

// Reads the content of the given file filePath and returns the slice of WatcherURLs.
// Returns nil if file is empty/not found.
func ReadWatchers(filePath string) ([]WatcherURL, error) {
	var watchers []WatcherURL

	var bytes []byte
	if err := readBytes(filePath, &bytes); err != nil {
		return nil, fmt.Errorf("error reading %v: %v", filePath, err)
	}

	if bytes == nil {
		return watchers, nil
	}

	if err := json.Unmarshal(bytes, &watchers); err != nil {
		return nil, fmt.Errorf("error unmarshalling: %v", err)
	}

	return watchers, nil
}

// Loads the content of the fiile filePath, appends the new items to the unmarshaled content and updates the file.
// Returns error if reading, marshalling or writing fails.
// Default filePath is "items.json".
func AppendItemIDs(filePath string, items []ItemID) error {
	if filePath == "" {
		filePath = "items.json"
	}

	itemsToWrite, err := ReadItemIDs(filePath)
	if err != nil {
		return fmt.Errorf("error reading itemIDs: %v", err)
	}

	itemsToWrite = append(itemsToWrite, items...)

	updatedContent, err := json.Marshal(itemsToWrite)
	if err != nil {
		return fmt.Errorf("error marshalling items: %v", err)
	}

	if err := os.WriteFile(filePath, updatedContent, 0644); err != nil {
		return fmt.Errorf("error writing file while updating the json content: %v", err)
	}

	return nil
}

func ItemExists(item ItemID) bool {
	ids, err := ReadItemIDs("items.json")
	if err != nil {
		log.Printf("error reading itemIDs: %v", err)
		return false
	}

	for _, id := range ids {
		if id.Id == item.Id {
			return true
		}
	}

	return false
}

// Reads the content of the given file filePath and returns the slice of ItemID.
// Returns nil if file is empty/not found.
func ReadItemIDs(filePath string) ([]ItemID, error) {
	var items []ItemID

	var bytes []byte
	if err := readBytes(filePath, &bytes); err != nil {
		return nil, fmt.Errorf("error reading %v: %v", filePath, err)
	}

	if bytes == nil {
		return nil, nil
	}

	if err := json.Unmarshal(bytes, &items); err != nil {
		return nil, fmt.Errorf("error unmarshalling: %v", err)
	}

	return items, nil
}
