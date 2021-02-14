package main

import (
	"fmt"
	"github.com/go-telegram-bot-api/telegram-bot-api"
)

func notifyPrice(message Message, btcWatcherChannelChatId int64, bot *tgbotapi.BotAPI, finished chan bool) {

	actualPrice := cs2f(message.Price)

	text := fmt.Sprintf("#%s \xF0\x9F\x92\xB0 \nPrice: USD$ %.2f \n", message.ID, actualPrice)

	btcWatcherChannel := tgbotapi.NewMessage(btcWatcherChannelChatId, text)

	_, err := bot.Send(btcWatcherChannel)

	if err != nil {
		finished <- true
	}

}

func alertVariation(message Message, minimumMovementHour float64, minimumMovementDayli float64, btcPriceChangeAlertChatId int64, bot *tgbotapi.BotAPI, finished chan bool, lastMessage Message) {
	
	lastAlertExceed := false

	exceed, change, lastPrice := calculateIfExceedLimitAndReturnDifferenceAndLastPrice(
										 cs2f(message.Price),
										 cs2f(message.OneHour.PriceChange), 
										 minimumMovementHour)

	if exceed {

		text := fmt.Sprintf("#%s changed %.2f%% in 1 hour from %.2f to %.2f \xE2\x9A\xA0 ", message.ID,
																							 change,
																							  lastPrice,
																							   cs2f(message.Price))

		btcPriceChangeAlert := tgbotapi.NewMessage(btcPriceChangeAlertChatId, text)

		_, err := bot.Send(btcPriceChangeAlert)

		if err != nil {
			finished <- true
		}
	}

	if lastMessage.ID != "" { //empty
		lastAlertExceed, _, _ = calculateIfExceedLimitAndReturnDifferenceAndLastPrice(
											 cs2f(lastMessage.Price),
											 cs2f(lastMessage.OneD.PriceChange), 
											 minimumMovementHour)
	}
	
	exceed, change, lastPrice = calculateIfExceedLimitAndReturnDifferenceAndLastPrice(
										 cs2f(message.Price),
										 cs2f(message.OneD.PriceChange), 
										 minimumMovementDayli)

	if !lastAlertExceed && exceed {

		text := fmt.Sprintf("#%s changed %.2f%% in 24 hours from %.2f to %.2f \xE2\x9A\xA0 ", message.ID,
																							 change,
																							  lastPrice,
																							   cs2f(message.Price))

		btcPriceChangeAlert := tgbotapi.NewMessage(btcPriceChangeAlertChatId, text)

		_, err := bot.Send(btcPriceChangeAlert)

		if err != nil {
			finished <- true
		}
	}

}
