package discordBot

import (
	"fmt"
	"log"
	"os"
	"os/signal"
	"time"

	"github.com/bwmarrin/discordgo"
)

type MessageConfig struct {
	Title        string
	Description  string
	URL          string
	ThumbnailURL string
	Color        int
	Fields       []EmbedField
}

type EmbedField struct {
	Name   string
	Value  string
	Inline bool
}

func SendMessageEmbed(s *discordgo.Session, channelID string, config MessageConfig) error {
	if s == nil {
		return fmt.Errorf("discord session cannot be nil")
	}

	emb := &discordgo.MessageEmbed{
		Title:       config.Title,
		Description: config.Description,
		URL:         config.URL,
		Color:       config.Color,
		Timestamp:   time.Now().Format(time.RFC3339),
		Footer: &discordgo.MessageEmbedFooter{
			Text: "vinted",
		},
	}

	if config.ThumbnailURL != "" {
		emb.Thumbnail = &discordgo.MessageEmbedThumbnail{
			URL: config.ThumbnailURL,
		}
	}

	if len(config.Fields) > 0 {
		emb.Fields = make([]*discordgo.MessageEmbedField, len(config.Fields))
		for i, field := range config.Fields {
			emb.Fields[i] = &discordgo.MessageEmbedField{
				Name:   field.Name,
				Value:  field.Value,
				Inline: field.Inline,
			}
		}
	}

	_, err := s.ChannelMessageSendComplex(channelID, &discordgo.MessageSend{
		Embed: emb,
	})

	return err
}

func newMessage(s *discordgo.Session, m *discordgo.MessageCreate) {
	if m.Author.ID == s.State.User.ID {
		return
	}

	config := MessageConfig{
		Title:        "item details",
		Description:  "product description",
		URL:          "https://example.com/ex",
		ThumbnailURL: "https://icons.veryicon.com/png/o/miscellaneous/construction-of-fengying-website-in-xian/no-copy.png",
		Color:        0x00ff00,
		Fields: []EmbedField{
			{
				Name:   "price",
				Value:  "$42.99",
				Inline: true,
			},
		},
	}

	log.Printf("sending message to channel %v for %v", m.ChannelID, m.Author.Username)
	if err := SendMessageEmbed(s, m.ChannelID, config); err != nil {
		log.Printf("error sending message: %v", err)
		return
	}

}

func Run(botToken string) {
	var err error
	bot, err := discordgo.New("Bot " + botToken)
	if err != nil {
		log.Fatalf("wrong parameters for bot: %v", err)
	}

	bot.AddHandler(newMessage)

	bot.Open()
	defer bot.Close()

	c := make(chan os.Signal, 1)
	signal.Notify(c, os.Interrupt)
	<-c
}
