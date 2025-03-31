package agent

import (
	"log"
	"math/rand"
	"strings"
	"time"

	"github.com/smatand/vinted_go/db"
	vintedApi "github.com/smatand/vinted_go/vinted_api"
)

const (
	watchersFilePath = "watchers.json"
	itemsFilePath    = "items.json"
	maxRandWait      = 120
)

func itemContainsCurrency(item vintedApi.VintedItemResp, currencies []string) bool {
	itemCurrency := item.Conversion.SellerCurrency
	// If the item's currency is empty, probably it is from same country as user
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

func Run(newItemsChan chan<- []vintedApi.VintedItemResp) {
	for {
		watcher, err := db.ReadWatchers(watchersFilePath)
		if err != nil {
			log.Fatalf("error while getting urls to watch: %v", err)
		}

		if watcher == nil {
			log.Println("no watcher to watch")
		}

		// Parse user given url and then fethc item from the parsed API url
		for _, url := range watcher {
			items, err := vintedApi.GetVintedItems(url.URL)
			if err != nil {
				log.Printf("error while getting items: %v", err)

				break
			}

			var itemIDs []db.ItemID
			var uniqueItems []vintedApi.VintedItemResp
			for _, item := range items.Items {
				itemID := db.ItemID{Id: item.ID}

				// The item was already added to the db, skip
				if db.ItemExists(itemID) {
					continue
				}

				itemIDs = append(itemIDs, itemID)

				// The item is sold by other seller's nationality than the user wants, skip
				// But keep the record of it so it won't have to be processed later again
				if !itemContainsCurrency(item, url.SellerCurrency) {
					continue
				}

				uniqueItems = append(uniqueItems, item)
			}
			db.AppendItemIDs(itemsFilePath, itemIDs)

			// Pass the details of items to discordBot
			newItemsChan <- uniqueItems

			// To prevent API overload
			randomWait := time.Duration(rand.Intn(10)+1) * time.Second

			time.Sleep(randomWait)
		}

		// And again, another wait
		randomInterval := time.Duration(rand.Intn(maxRandWait)+maxRandWait) * time.Second
		time.Sleep(randomInterval)
	}
}
