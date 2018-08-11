package cnb

/*
import (
	"fmt"
	"time"

	decimal "gopkg.in/inf.v0"
)

var oneUnit *decimal.Dec

func init() {
	one, _ := new(decimal.Dec).SetString("1")
	oneUnit = one
}

type ExchangeRates struct {
	Id    int
	Date  time.Time
	Rates []ExchangeRate
}

type ExchangeRate struct {
	From Money
	To   Money
}

func (rate ExchangeRate) String() string {
	return fmt.Sprintf("ExchangeRate{ %s %s = %s %s }", rate.From.Amount.String(), rate.From.Currency, rate.To.Amount.String(), rate.To.Currency)
}

type Money struct {
	Amount   *decimal.Dec
	Currency string
}

func (rate *ExchangeRate) Normalize() {
	if rate.From.Amount.Cmp(oneUnit) != 1 {
		return
	}
	rate.To.Amount = new(decimal.Dec).QuoRound(rate.To.Amount, rate.From.Amount, 35, decimal.RoundCeil)
	rate.From.Amount = new(decimal.Dec).Set(oneUnit)
}
*/