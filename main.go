package main

import (
	"fmt"
	"log"
	"os"
	"strconv"

	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api"
	"github.com/joho/godotenv"
)

func main() {

	err := godotenv.Load("app.env")

	if err != nil {
		log.Fatal("Error loading app.env file")
	}

	telegramApiKey := os.Getenv("TELEGRAM_API_KEY")
	nomicsApiKey := os.Getenv("NOMICS_API_KEY")

	criptoWatcherChannelChatId, _ := strconv.ParseInt(os.Getenv("CRIPTO_WATCHER_CHANNEL_CHAT_ID"), 10, 64)
	criptoPriceChangeAlertChatId, _ := strconv.ParseInt(os.Getenv("CRIPTO_PRICE_CHANGE_ALERT_CHAT_ID"), 10, 64)
	minimumMovementHour := 5.0   // to the alert
	minimumMovementDayli := 10.0 // to the alert

	finished := make(chan bool)

	lastMessages := make(map[string]Message)

	bot, err := tgbotapi.NewBotAPI(telegramApiKey)

	if err != nil {
		log.Panic(err)
	}

	log.Printf("Authorized on account %s", bot.Self.UserName)

	nomics := GetNomics("BTC,ETH,LINK", nomicsApiKey, "USD", 1800)

	list, err := nomics.GetCurrencyUpdateChannel()

	if err != nil {
		panic(err)
	}

	go func(updateList MessageChannel, conversion string) {
		for updates := range updateList {
			for _, update := range updates {
				notifyPrice(update, criptoWatcherChannelChatId, bot, finished)
				alertVariation(update, minimumMovementHour, minimumMovementDayli, criptoPriceChangeAlertChatId, bot, finished, lastMessages[update.ID])
				lastMessages[update.ID] = update
			}
		}
	}(list, nomics.Conversion)

	failed := <-finished
	if failed {
		fmt.Printf("Ending with error %v", err)
	}

}
