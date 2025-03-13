package main

import (
	"log"
	"os"

	"github.com/joho/godotenv"
	discordBot "github.com/smatand/vinted_go/bot"
)

func main() {
	err := godotenv.Load()
	if err != nil {
		log.Fatalf("Error loading the .env file: %s", err)
	}

	token := os.Getenv("DISCORD_TOKEN")
	guildID := os.Getenv("GUILD_ID")

	discordBot.Run(token, guildID)
}
