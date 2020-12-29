package bot

import (
	"log"

	"github.com/kevinburke/twilio-go"
	"github.com/kimrgrey/go-telegram"
	"github.com/spf13/viper"
)

// Clients struct combines all available clients
type Clients struct {
	Telegram TelegramConfig `mapstructure:"telegram"`
	Twilio   TwilioConfig   `mapstructure:"twilio"`
}

// TelegramConfig models Telegram configuration
type TelegramConfig struct {
	BotKey string `mapstructure:"bot_key"`
}

// TwilioConfig models Twilio configuration
type TwilioConfig struct {
	AccountSid string `mapstructure:"account_sid"`
	AuthToken  string `mapstructure:"auth_token"`
	Number     string `mapstructure:"number"`
}

// Twilio contains a Twilio client as well as the Twilio number
type Twilio struct {
	Client *twilio.Client
	Number string
}

// LoadClients loads registered clients/channels in the chn.yml file
func LoadClients(path *string) map[string]interface{} {
	config := viper.New()
	config.SetConfigName("chn")
	config.AddConfigPath(*path)

	clients := make(map[string]interface{})

	if err := config.ReadInConfig(); err != nil {
		log.Println(err)
		return clients
	}

	var end Clients
	if err := config.Unmarshal(&end); err != nil {
		log.Println(err)
		return clients
	}

	// TELEGRAM
	if end.Telegram != (TelegramConfig{}) {
		telegramClient := telegram.NewClient(end.Telegram.BotKey)
		clients["telegram"] = telegramClient
		log.Printf("Added Telegram client: %v\n", telegramClient.GetMe())
	}

	// TWILIO
	if end.Twilio != (TwilioConfig{}) {
		twilioClient := twilio.NewClient(end.Twilio.AccountSid, end.Twilio.AuthToken, nil)
		clients["twilio"] = Twilio{twilioClient, end.Twilio.Number}
		log.Printf("Added Twilio client: %v\n", twilioClient.AccountSid)
	}

	return clients
}
