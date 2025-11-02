package main

import (
	"log"
	"os"

	"github.com/joho/godotenv"
 	"github.com/smatand/vinted_go/bot"
	"github.com/smatand/vinted_go/db"
)

func main() {
	err := godotenv.Load()
	if err != nil {
		log.Fatalf("Error loading the .env file: %s", err)
	}

	mongoURI := os.Getenv("MONGODB_URI")
	if mongoURI != "" {
		err = db.InitDB(mongoURI)
		if err != nil {
			log.Fatalf("Error initializing DB: %s", err)
		}
	}

	token := os.Getenv("DISCORD_TOKEN")
	guildID := os.Getenv("GUILD_ID")

	bot.Run(token, guildID)
}
