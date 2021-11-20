package main

import (
	"math"
	"strconv"
)

func calculateIfExceedLimitAndReturnDifferenceAndLastPrice(actualPrice, priceChange, minimumMovement float64) (bool, float64, float64) {

	lastPrice := 0.0

	if priceChange < 0 {
		lastPrice = actualPrice + math.Abs(priceChange)
	} else {
		lastPrice = actualPrice - priceChange
	}

	change := (((lastPrice - actualPrice) / actualPrice) * 100) * -1
	exceed := (change > minimumMovement || change < (minimumMovement*-1))

	return exceed, change, lastPrice
}

func cs2f(number string) float64 {
	newnumber, _ := strconv.ParseFloat(number, 64)
	return math.Round(newnumber*100) / 100
}
