package discordBot

import (
	"fmt"
	"log"
	"os"
	"os/signal"
	"strings"

	"github.com/bwmarrin/discordgo"
)

var (
	commands = []*discordgo.ApplicationCommand{
		{
			Name:        "watch",
			Description: "Insert Vinted's URL to watch items. You may choose which currencies to include.",
			Type:        discordgo.ChatApplicationCommand,
			Options: []*discordgo.ApplicationCommandOption{
				{
					Name:         "url",
					Description:  "Insert url of vinted page, e. g. https://www.vinted.cz/catalog/1206-outerwear",
					Type:         discordgo.ApplicationCommandOptionString,
					Required:     true,
					Autocomplete: false,
				},
				{
					Name:        "currency_eur",
					Description: "Include items in EUR",
					Type:        discordgo.ApplicationCommandOptionBoolean,
					Required:    false,
				},
				{
					Name:        "currency_czk",
					Description: "Include items in CZK",
					Type:        discordgo.ApplicationCommandOptionBoolean,
					Required:    false,
				},
				{
					Name:        "currency_pln",
					Description: "Include items in PLN",
					Type:        discordgo.ApplicationCommandOptionBoolean,
					Required:    false,
				},
			},
		},
	}

	commandHandlers = map[string]func(s *discordgo.Session, i *discordgo.InteractionCreate){
		"watch": handleWatcher,
	}
)

func handleWatcher(s *discordgo.Session, i *discordgo.InteractionCreate) {
	switch i.Type {
	case discordgo.InteractionApplicationCommand:
		data := i.ApplicationCommandData()

		var url string
		var selectedCurrencies []string
		for _, opt := range data.Options {
			switch opt.Name {
			case "url":
				if opt.Value != "" {
					url = opt.StringValue()
				}

			case "currency_eur":
				if opt.BoolValue() {
					selectedCurrencies = append(selectedCurrencies, "EUR")
				}
			case "currency_czk":
				if opt.BoolValue() {
					selectedCurrencies = append(selectedCurrencies, "CZK")
				}
			case "currency_pln":
				if opt.BoolValue() {
					selectedCurrencies = append(selectedCurrencies, "PLN")
				}
			}
		}

		if len(selectedCurrencies) == 0 {
			selectedCurrencies = []string{"EUR", "CZK", "PLN"}
		}

		currenciesStr := strings.Join(selectedCurrencies, ", ")

		err := s.InteractionRespond(i.Interaction, &discordgo.InteractionResponse{
			Type: discordgo.InteractionResponseChannelMessageWithSource,
			Data: &discordgo.InteractionResponseData{
				Content: fmt.Sprintf(
					"you entered %s and currencies %v to watch", url, currenciesStr,
				),
			},
		})
		if err != nil {
			log.Printf("error responding to interaction: %v", err)
		}

		log.Printf("user entered URL %s with currencies %v", url, currenciesStr)
	}
}

func Run(botToken string) error {
	GuildID := ""

	if botToken != "" && !strings.HasPrefix(botToken, "Bot ") {
		botToken = "Bot " + botToken
	} else {
		log.Fatalf("invalid bot token: %s", botToken)
	}

	bot, err := discordgo.New(botToken)
	if err != nil {
		return err
	}

	bot.AddHandler(func(s *discordgo.Session, r *discordgo.Ready) { log.Println("Bot is up!") })
	bot.AddHandler(func(s *discordgo.Session, i *discordgo.InteractionCreate) {
		if h, ok := commandHandlers[i.ApplicationCommandData().Name]; ok {
			h(s, i)
		}
	})

	if err := bot.Open(); err != nil {
		return fmt.Errorf("error opening connection: %v", err)
	}
	defer bot.Close()

	createdCommands, err := bot.ApplicationCommandBulkOverwrite(bot.State.User.ID, GuildID, commands)

	if err != nil {
		log.Fatalf("cannot register commands: %v", err)
	}

	c := make(chan os.Signal, 1)
	signal.Notify(c, os.Interrupt)
	<-c

	for _, cmd := range createdCommands {
		err := bot.ApplicationCommandDelete(bot.State.User.ID, GuildID, cmd.ID)
		if err != nil {
			log.Printf("cannot delete %q command: %v", cmd.Name, err)
		}
	}

	return nil
}
