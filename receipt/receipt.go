package receipt

import (
	"math"
	"strconv"
	"strings"
	"unicode"
)

const multipleConstant = 0.25

type Receipt struct {
	Retailer     string       `json:"retailer"`
	PurchaseDate purchaseDate `json:"purchaseDate"`
	PurchaseTime purchaseTime `json:"purchaseTime"`
	Items        []Item       `json:"items"`
	Total        string       `json:"total"`
}

func (r Receipt) GetScore() int {
	score := 0

	score += r.scoreRetailer()
	score += r.scoreTotal()
	score += r.scoreItems()

	for _, i := range r.Items {
		score += i.scoreDescription()
	}

	score += r.PurchaseDate.scoreDay()
	score += r.PurchaseTime.scoreTime()

	return score
}

// scoreRetailer counts the number of alphanumeric characters
// in the string
func (r Receipt) scoreRetailer() int {
	retailer := strings.TrimSpace(r.Retailer)

	score := 0
	for _, v := range retailer {
		if unicode.IsNumber(v) || unicode.IsLetter(v) {
			score++
		}
	}

	return score
}

func (r Receipt) scoreTotal() int {
	parts := strings.Split(r.Total, ".")

	score := 0
	if parts[1] == "00" {
		score += 50
	}

	total, _ := strconv.ParseFloat(r.Total, 32)
	if math.Mod(total, multipleConstant) == 0 {
		score += 25
	}

	return score
}

func (r Receipt) scoreItems() int {
	return (len(r.Items) / 2) * 5
}
