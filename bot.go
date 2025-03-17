package main

import (
	"log"
	"os"
	"time"

	"github.com/Pato05/unobot/bot"
	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api/v5"
)

func main() {
	log.Default().SetFlags(log.Ldate | log.Lshortfile)
	token := os.Getenv("BOT_TOKEN")
	if token == "" {
		log.Panic("the BOT_TOKEN environment variable is empty.")
	}

	tbot, err := tgbotapi.NewBotAPI(token)

	if err != nil {
		log.Panic("Bot authorization failed: ", err)
	}

	log.Printf("Logged in as %s", tbot.Self.UserName)

	handler := bot.NewBotHandler(tbot, true)

	config := tgbotapi.NewUpdate(-1)

	for {
		updates, err := tbot.GetUpdates(config)
		// stop here if the error happened because of another instance running
		if err != nil {
			log.Println("GetUpdates call failed:")
			log.Println(err)
			log.Println("\nWaiting 5 seconds before next call...")
			time.Sleep(time.Second * 5)
		}

		for _, update := range updates {
			// the goroutine is spawned inside after a copy is made
			handler.ProcessUpdate(update)
		}

		if len(updates) > 0 {
			config.Offset = updates[len(updates)-1].UpdateID + 1
		}
	}
}
