package main

import (
	"math"
	"math/rand"
)

func mapValue(
	value uint8,
	minIn uint8,
	maxIn uint8,
	minOut uint8,
	maxOut uint8,
) uint8 {
	finalValue := value

	if value > maxIn {
		finalValue = maxIn
	} else if value < minIn {
		finalValue = minIn
	}

	a := float64(maxIn-finalValue) / float64(maxIn-minIn)

	b := a * float64(maxOut-minOut)

	c := math.Round(b)

	if c < 0 {
		c = 0
	} else if c > 255 {
		c = 255
	}

	return maxOut - uint8(c)
}

func getRandomToken(slice string) string {
	randomIndex := rand.Intn(len(slice))
	return string(slice[randomIndex])
}
