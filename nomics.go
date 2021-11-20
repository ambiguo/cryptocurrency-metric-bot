package main

import (
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"time"
)

func GetNomics(currency string, apiKey string, conversion string, interval int32) *NomicsAPI {
	return &NomicsAPI{
		Currency:   currency,
		ApiKey:     apiKey,
		Conversion: conversion,
		Interval:   interval,
		Buffer:     4096,
		HttpClient: &http.Client{},
		baseUrl:    "https://api.nomics.com/v1",
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

				time.Sleep(time.Second * time.Duration(api.Interval))

				continue
			}

			ch <- messages

			time.Sleep(time.Second * time.Duration(api.Interval))
		}
	}()

	return ch, nil
}

func (api *NomicsAPI) GetCurrencyUpdate() ([]Message, error) {

	var message []Message

	request := fmt.Sprintf("/currencies/ticker?key=%s&ids=%s&interval=1h,1d&convert=%s", api.ApiKey, api.Currency, api.Conversion)

	err := api.makeRequest(request, &message)

	return message, err
}

func (api *NomicsAPI) makeRequest(path string, r interface{}) error {

	resp, err := api.HttpClient.Get(api.baseUrl + path)

	if err != nil {
		return err
	}

	defer resp.Body.Close()

	decoder := json.NewDecoder(resp.Body)

	return decoder.Decode(&r)
}
