package discordBot

import (
	"strings"
	"log"
	"os"
	"os/signal"

	"github.com/bwmarrin/discordgo"
)

func newMessage(discordSession *discordgo.Session, message *discordgo.MessageCreate) {
	if message.Author.ID == discordSession.State.User.ID {
		return
	}

	switch {
	case strings.Contains(message.Content, "ahoj"):
		discordSession.ChannelMessageSend(message.ChannelID, "caw")
	default:
		discordSession.ChannelMessageSend(message.ChannelID, message.Content)
	}
}

func Run(botToken string) {
	var err error
	bot, err := discordgo.New("Bot " + botToken)
	if err != nil {
		log.Fatalf("Wrong parameters for bot: %v", err)
	}

	bot.AddHandler(newMessage)

	bot.Open()
	defer bot.Close()

	c := make(chan os.Signal, 1)
	signal.Notify(c, os.Interrupt)
	<-c
}