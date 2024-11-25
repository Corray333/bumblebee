package telegram

import (
	"log"
	"os"

	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api/v5"
)

type TelegramClient struct {
	Bot *tgbotapi.BotAPI
}

func New() *TelegramClient {
	token := os.Getenv("BOT_TOKEN")
	bot, err := tgbotapi.NewBotAPI(token)
	if err != nil {
		log.Fatal("failed to create bot: ", err)
	}

	bot.Debug = true

	return &TelegramClient{
		Bot: bot,
	}
}
