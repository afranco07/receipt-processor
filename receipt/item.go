package receipt

import (
	"math"
	"strconv"
	"strings"
)

type Item struct {
	ShortDescription string
	Price            string
}

func (i Item) scoreDescription() int {
	desc := strings.TrimSpace(i.ShortDescription)

	// length of description is not a multiple of 3
	if len(desc)%3 != 0 {
		return 0
	}

	price, _ := strconv.ParseFloat(i.Price, 32)

	price *= 0.2

	price = math.Ceil(price)

	return int(price)
}
