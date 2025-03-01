package discordBot

import (
	"fmt"
	"os"
	"os/signal"
	"strings"
	"sync"

	//"time"

	"github.com/bwmarrin/discordgo"
	// "github.com/smatand/vinted_go/vinted"
	// vintedApi "github.com/smatand/vinted_go/vinted_api"
)

type WatcherURL struct {
	URL        string
	Currencies []string
}

type Bot struct {
	ses       *discordgo.Session
	channelID string
	watchers  []WatcherURL
	mu        sync.RWMutex // Protect concurrent access to watchers
}

func newBot(token string) (*Bot, error) {
	ses, err := discordgo.New("Bot " + token)
	if err != nil {
		return nil, fmt.Errorf("wrong parameters for bot: %v", err)
	}

	return &Bot{
		ses:       ses,
		channelID: "1200879123220942889",
		watchers:  make([]WatcherURL, 0),
	}, nil
}

func (b *Bot) handleMessage(s *discordgo.Session, m *discordgo.MessageCreate) {
	// Ignore messages from the bot itself
	if m.Author.ID == s.State.User.ID {
		return
	}

	// Check if message starts with !addWatcher
	if !strings.HasPrefix(m.Content, "!addWatcher") {
		return
	}

	// Parse the command
	watcher, err := parseWatcherCommand(m.Content)
	if err != nil {
		s.ChannelMessageSend(m.ChannelID, fmt.Sprintf("error: %v", err))
		return
	}

	// Add to watchers list
	b.mu.Lock()
	b.watchers = append(b.watchers, watcher)
	b.mu.Unlock()

	s.ChannelMessageSend(m.ChannelID, fmt.Sprintf("added watcher for URL: %s with currencies: %v", watcher.URL, watcher.Currencies))
}

func parseWatcherCommand(content string) (WatcherURL, error) {
	parts := strings.Fields(content)
	if len(parts) < 2 {
		return WatcherURL{}, fmt.Errorf("invalid command format")
	}

	var url string
	var currencies []string
	defaultCurrencies := []string{"EUR", "CZK", "PLN"}

	// Parse parameters
	for _, part := range parts[1:] {
		if strings.HasPrefix(part, "url:") {
			url = strings.TrimPrefix(part, "url:")
		} else if strings.HasPrefix(part, "seller_currency:") {
			currency := strings.TrimPrefix(part, "seller_currency:")
			if !isValidCurrency(currency) {
				return WatcherURL{}, fmt.Errorf("invalid currency: %s", currency)
			}
			currencies = append(currencies, currency)
		}
	}

	if url == "" {
		return WatcherURL{}, fmt.Errorf("url is mandatory")
	}

	// If no currencies specified, use all defaults
	if len(currencies) == 0 {
		currencies = defaultCurrencies
	}

	fmt.Println("URL:", url)
	fmt.Println("currencies:", currencies)

	return WatcherURL{
		URL:        url,
		Currencies: currencies,
	}, nil
}

func isValidCurrency(currency string) bool {
	validCurrencies := map[string]bool{
		"EUR": true,
		"CZK": true,
		"PLN": true,
	}
	return validCurrencies[currency]
}

//func (b *Bot) startWatcherLoop(watchings vinted.Vinted) {
//	ticker := time.NewTicker(1 * time.Minute)
//	go func() {
//		for range ticker.C {
//			b.checkWatchers(watchings)
//		}
//	}()
//}
//
//func (b *Bot) checkWatchers(watchings vinted.Vinted) {
//	b.mu.RLock()
//	defer b.mu.RUnlock()
//
//	for _, watch := range b.watchers {
//		var vintedParams vinted.Vinted
//		vintedParams.ParseParams(watch.URL)
//		items, err := vintedApi.GetVintedItems(vintedParams)
//		if err != nil {
//			fmt.Printf("Error checking URL %s: %v\n", watch.URL, err)
//			continue
//		}
//
//		// Handle the items (you might want to customize this part)
//		if len(items) > 0 {
//			message := fmt.Sprintf("New items found for %s:\n", watch.URL)
//			b.ses.ChannelMessageSend(b.channelID, message)
//			// Add logic to handle/display items
//		}
//	}
//}

func Run(botToken string) error {
	bot, err := newBot(botToken)
	if err != nil {
		return err
	}

	// Add message handler
	bot.ses.AddHandler(bot.handleMessage)

	// Open connection
	if err := bot.ses.Open(); err != nil {
		return fmt.Errorf("error opening connection: %v", err)
	}
	defer bot.ses.Close()

	// Start the watcher loop
	//	bot.startWatcherLoop(watchings)

	// Wait for interrupt signal
	c := make(chan os.Signal, 1)
	signal.Notify(c, os.Interrupt)
	<-c

	return nil
}
