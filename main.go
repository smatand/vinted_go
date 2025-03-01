package main

import (
	"fmt"
	"log"
	"os"

	"github.com/joho/godotenv"
	discordBot "github.com/smatand/vinted_go/bot"
	"github.com/smatand/vinted_go/db"
)

func main() {
	err := godotenv.Load()
	if err != nil {
		log.Fatalf("Error loading the .env file: %s", err)
	}

	// token
	token := os.Getenv("DISCORD_TOKEN")

	//	const exampleUrl = "https://www.vinted.sk/catalog?search_text=&catalog[]=2050&size_ids[]=206&size_ids[]=208&brand_ids[]=53&search_id=20979747721&order=newest_first"
	//
	//	vinted := vinted.Vinted{}
	//	vinted.ParseParams(exampleUrl)

	toWrite := db.WatcherURL{
		URL:             "https://www.vinted.sk/catalog?search_text=&catalog[]=2050&size_ids[]=206&size_ids[]=208&brand_ids[]=53&search_id=20979747721&order=newest_first",
		Seller_currency: []string{"EUR", "USD"},
	}
	if db.WriteWatcherURL("watchers.json", toWrite) != nil {
		fmt.Println("error writing")
	}

	toRead, err := db.ReadWatcherURLs("watchers.json")
	if err != nil {
		fmt.Println(err)
	}
	fmt.Println(toRead)

	toWrite2 := []db.ItemID{
		{Id: "2813",},
		{Id: "2814",},
	}
	if db.WriteItemIDs("items.json", toWrite2) != nil {
		fmt.Println("error writing")
	}

	toRead2, err := db.ReadItemIDs("items.json")
	if err != nil {
		fmt.Println(err)
	}
	fmt.Println(toRead2)


	// discordBot
	discordBot.Run(token)
}
