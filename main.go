package main

import (
	"log"
	"os"

	"github.com/joho/godotenv"
	"github.com/smatand/vinted_go/bot"
)

func main() {
	err := godotenv.Load()
	if err != nil {
		log.Fatalf("Error loading the .env file: %s", err)
	}

	discordToken := os.Getenv("DISCORD_TOKEN")

	discordBot.Run(discordToken)
}