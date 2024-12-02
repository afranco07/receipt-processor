package receipt

import (
	"errors"
	"fmt"
	"math"
	"strconv"
	"strings"
	"unicode"

	"github.com/go-playground/validator/v10"
)

const multipleConstant = 0.25

type Receipt struct {
	Retailer     string       `json:"retailer" validate:"required"`
	PurchaseDate purchaseDate `json:"purchaseDate" validate:"required"`
	PurchaseTime purchaseTime `json:"purchaseTime" validate:"required"`
	Items        []Item       `json:"items" validate:"gt=0,dive"`
	Total        string       `json:"total" validate:"required,numeric"`
}

// GetScore gets the total number of points that is
// awarded to receipt
func (r Receipt) GetScore() (int, error) {
	score := 0

	score += r.scoreRetailer()

	total, err := r.scoreTotal()
	if err != nil {
		return 0, nil
	}
	score += total

	score += r.scoreItems()

	for _, i := range r.Items {
		score += i.scoreDescription()
	}

	score += r.PurchaseDate.scoreDay()
	score += r.PurchaseTime.scoreTime()

	return score, nil
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

// scoreTotal checks if the receipt total is a multiple of
// 0.25 and has no cents
func (r Receipt) scoreTotal() (int, error) {
	score := 0

	total, err := strconv.ParseFloat(r.Total, 32)
	if err != nil {
		return 0, err
	}

	if math.Mod(total, multipleConstant) == 0 {
		score += 25
	}

	if math.Mod(total*100, 100) == 0 {
		score += 50
	}

	return score, nil
}

func (r Receipt) scoreItems() int {
	return (len(r.Items) / 2) * 5
}

func (r Receipt) ValidateReceipt(v *validator.Validate) (validator.ValidationErrors, error) {
	err := v.Struct(r)
	if err == nil {
		return nil, nil
	}

	var validationErrors validator.ValidationErrors
	e := errors.As(err, &validationErrors)
	if !e {
		return nil, nil
	}

	var fields string
	for _, v := range validationErrors {
		fields += v.Field() + ", "
	}

	return validationErrors, errors.New(fmt.Sprintf("validation errors for the following fields: %s", fields))
}
