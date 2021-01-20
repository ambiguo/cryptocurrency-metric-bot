package main

import (
  "net/http"
  "encoding/json"
  "log"
  "time"
)

func GetNomics(currency string,  apiKey string, conversion string, interval int32) *NomicsAPI {
  return &NomicsAPI{
      Currency: currency, 
      ApiKey: apiKey,
      Conversion: conversion, 
      Interval: interval, 
      Buffer: 4096, 
      HttpClient: &http.Client{},
      baseUrl: "https://api.nomics.com/v1",
  }
}

func (api *NomicsAPI) GetCurrencyUpdateChannel() (MessageChannel, error) {
	ch := make(chan []Message, api.Buffer)

	go func() {
		for {
			select {
			case <-api.shutdownConsumer:
				close(ch)
				return
			default:
			}

			messages, err := api.GetCurrencyUpdate()

			if err != nil {

				log.Println(err)

				time.Sleep(time.Second * 3)

				continue
			}

			ch <- messages

      time.Sleep(time.Duration(api.Interval) * time.Second)
		}
	}()

	return ch, nil
}


func (api *NomicsAPI) GetCurrencyUpdate() ([]Message, error) {

  var message []Message 

	err := api.makeRequest("/currencies/ticker?key="+api.ApiKey+ "&ids="+api.Currency+"&interval=1h&convert="+api.Conversion,
  &message)

	return message, err
}

func (api *NomicsAPI) makeRequest(path string, r interface{}) (error) {

  resp, err := api.HttpClient.Get(api.baseUrl + path)

  defer resp.Body.Close()

  if err != nil {
		return err
	}

  decoder := json.NewDecoder(resp.Body)

  return decoder.Decode(&r)
} 

