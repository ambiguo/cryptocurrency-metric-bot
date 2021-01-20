package main

import (
  "log"
  "fmt"
  "math"
  "strconv"
	"github.com/go-telegram-bot-api/telegram-bot-api"
)

func main() {

  telegramApiKey := "<TELEGRAM API KEY>"
  nomicsApiKey := "<NOMICS API KEY>"
  btcWatcherChannelChatId := -0
  btcPriceChangeAlertChatId := -0
  minimumMovement := 5.0 // to the alert

  finished := make(chan bool)

  bot, err := tgbotapi.NewBotAPI(telegramApiKey)

  if err != nil {
    log.Panic(err)
  }
 
  bot.Debug = false

  log.Printf("Authorized on account %s", bot.Self.UserName)

  nomics := GetNomics("BTC,ETH", nomicsApiKey, "USD", 1800)

  list, err := nomics.GetCurrencyUpdateChannel()

  if err != nil {
    panic(err)
  }

  go func(updateList MessageChannel, conversion string) {

    for updates := range updateList {
      for _, update := range updates {

        actualPrice := cs2f(update.Price)

        text := fmt.Sprintf("#%s \xF0\x9F\x92\xB0 \nPrice: USD$ %.2f \n", update.ID, actualPrice)

        btcWatcherChannel := tgbotapi.NewMessage(btcWatcherChannelChatId, text)

        _, err = bot.Send(btcWatcherChannel)

        if err != nil {
            finished <- true
        }

        priceChange := cs2f(update.OneHour.PriceChange) 
        lastPrice := 0.0

        if priceChange < 0 {
          lastPrice = actualPrice + math.Abs(priceChange)
        } else { 
          lastPrice = actualPrice - priceChange
        }

        change := (((lastPrice-actualPrice)/actualPrice) * 100) * -1

        if change > minimumMovement || change < (minimumMovement*-1) {

          text = fmt.Sprintf("#%s changed %.2f%% in 1 hour from %.2f to %.2f \xE2\x9A\xA0 ", update.ID, change, lastPrice, actualPrice)

          btcPriceChangeAlert := tgbotapi.NewMessage(btcPriceChangeAlertChatId, text)

          _, err = bot.Send(btcPriceChangeAlert)

          if err != nil {
              finished <- true
          }
        }
      }
	  }
  }(list, nomics.Conversion)

  select {
      case <-finished:
        fmt.Println("Ending with error %v", err)
  }

}


func cs2f(number string) float64 {
   newnumber, _ := strconv.ParseFloat(number, 64)
   return math.Round(newnumber*100)/100
}
