package agent

import (
	"log"
	"math/rand"
	"strings"
	"time"

	"github.com/smatand/vinted_go/db"
	vintedApi "github.com/smatand/vinted_go/vintedApi"
)

const (
	DefaultWatchersFilePath = "watchers.json"
	DefaultItemsFilePath    = "items.json"
	DefaultMaxRandWait      = 120
	DefaultMaxRandWaitBetweenURLs = 10
)

type VintedAgent struct {
	WatchersFilePath string
	ItemsFilePath string
	MaxWaitLoops int
	MaxWaitURLs int
}

func NewVintedAgent(watchersFilePath, itemsFilePath string) * VintedAgent {
	if watchersFilePath == "" {
		watchersFilePath = DefaultWatchersFilePath	
	}

	if itemsFilePath == "" {
		itemsFilePath = DefaultItemsFilePath
	}
	return &VintedAgent{
		WatchersFilePath: watchersFilePath,
		ItemsFilePath: itemsFilePath,
		MaxWaitLoops: DefaultMaxRandWait,
		MaxWaitURLs: DefaultMaxRandWaitBetweenURLs,
	}
}

func (ag * VintedAgent) Start(newItemsChan chan<- []vintedApi.VintedItemResp) {
	for {
		ag.checkWatchers(newItemsChan)

		ag.randomSleep(ag.MaxWaitLoops)
	}
}

func (ag * VintedAgent) checkWatchers(newItemsChan chan<- []vintedApi.VintedItemResp) {
	watcherURLs, err := db.ReadWatchers(ag.WatchersFilePath)	
	if err != nil {
		log.Fatalf("error while reading watcher urls: %v", err)
	}

	if len(watcherURLs) == 0 {
		log.Println("no watcher to watch")
	}

	// Parse user given url and then fetch the item from the parsed API url.
	for _, watcherURL := range watcherURLs {
		ag.processWatcher(watcherURL, newItemsChan)

		// We don't want to fetch all watcher URLs at once, rather make pauses between each url fetch.
		ag.randomSleep(ag.MaxWaitURLs)
	}
}

func (ag * VintedAgent) processWatcher(watcher db.WatcherURL, newItemsChan chan<- []vintedApi.VintedItemResp) {
	items, err := vintedApi.GetVintedItems(watcher.URL)
	if err != nil {
		log.Printf("error while getting items: %v", err)

		return
	}

	uniqueItems := ag.filterItems(items, watcher)

	// Pass the details of items to discordBot.
	if len(uniqueItems) > 0 {
		newItemsChan <- uniqueItems
	}
}

func (ag * VintedAgent) filterItems(items *vintedApi.VintedItemsResp, watcher db.WatcherURL) []vintedApi.VintedItemResp {
	var itemIDs []db.ItemID
	var uniqueItems []vintedApi.VintedItemResp

	for _, item := range items.Items {
		itemID := db.ItemID{Id: item.ID}

		// Skip items already in the db.
		if db.ItemExists(itemID) {
			continue
		}

		// Track the ID regardless of the currency.
		itemIDs = append(itemIDs, itemID)

		// The item is sold by other seller's nationality than the user wants, skip.
		// But keep the record of it so it won't have to be processed later again.
		if !itemContainsCurrency(item, watcher.SellerCurrency) {
			continue
		}

		uniqueItems = append(uniqueItems, item)
	}

	// Only update if new items have been found.
	if len(itemIDs) > 0 {
		db.AppendItemIDs(ag.ItemsFilePath, itemIDs)
	}

	return uniqueItems
}

func (ag * VintedAgent) randomSleep(max int) {
	randomWait := time.Duration(rand.Intn(max) + max)
	time.Sleep(randomWait)
}


func itemContainsCurrency(item vintedApi.VintedItemResp, currencies []string) bool {
	itemCurrency := item.Conversion.SellerCurrency
	// If the item's currency is empty, probably it is from same country as user.
	if itemCurrency == "" {
		return true
	}

	for _, currency := range currencies {
		if strings.Contains(itemCurrency, currency) {
			return true
		}
	}

	return false
}