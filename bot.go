package main

import (
	"log"
	"os"

	"github.com/Pato05/unobot/bot"
	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api/v5"
)

func main() {
	token := os.Getenv("BOT_TOKEN")
	if token == "" {
		log.Panic("the BOT_TOKEN environment variable is empty.")
	}

	tbot, err := tgbotapi.NewBotAPI(token)
	tbot.Debug = true

	if err != nil {
		log.Panic("Bot authorization failed: ", err)
	}

	log.Printf("Logged in as %s", tbot.Self.UserName)
	log.Default().SetFlags(log.Ldate | log.Llongfile)

	handler := bot.NewBotHandler(tbot, true)

	updates := tbot.GetUpdatesChan(tgbotapi.NewUpdate(-1))

	for update := range updates {
		//gen.HandleGreyCardsGen(update, tbot)
		go handler.ProcessUpdate(update)
	}
}
