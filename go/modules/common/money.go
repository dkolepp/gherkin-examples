package bddcommon

import (
	"fmt"
	"math/big"
	"regexp"
	"strconv"
)

// Money is an amount in whole cents. Steps write it as $12.50.
type Money int64

// String formats m as $12.50 (or -$3.30).
func (m Money) String() string {
	sign := ""
	if m < 0 {
		sign, m = "-", -m
	}
	return fmt.Sprintf("%s$%d.%02d", sign, m/100, m%100)
}

// Step pattern pieces. Use them inside a step's regular expression:
//
//	"^the cart total should be " + bddcommon.MoneyPattern + "$"
//
// Each one is a single capture group, so the step function receives the
// matched text as one string argument: parse it with MustMoney or let godog
// convert a Qty to int.
const (
	MoneyPattern = `(\$\d+(?:\.\d{2})?)` // $12.50 or $12
	QtyPattern   = `(\d+)`               // 3
)

var moneyText = regexp.MustCompile(`^\$\d+(?:\.\d{2})?$`)

// ParseMoney turns step text such as "$12.50" into Money.
func ParseMoney(text string) (Money, error) {
	if !moneyText.MatchString(text) {
		return 0, fmt.Errorf("%q is not money: write it as $12.50", text)
	}
	return ToMoney(text[1:])
}

// MustMoney is ParseMoney for text a step pattern has already matched with
// MoneyPattern. It panics on anything else, which godog reports as a failed step.
func MustMoney(text string) Money {
	m, err := ParseMoney(text)
	if err != nil {
		panic(err)
	}
	return m
}

// ToMoney turns a number written as text ("12.50", "8", "1.005") into Money,
// rounding half up to whole cents. Use it for table cells and JSON data.
func ToMoney(text string) (Money, error) {
	r, ok := new(big.Rat).SetString(text)
	if !ok {
		return 0, fmt.Errorf("%q is not a number", text)
	}
	cents := new(big.Rat).Mul(r, big.NewRat(100, 1))
	// Round half up (away from zero for negatives, like Python's ROUND_HALF_UP).
	half := big.NewRat(1, 2)
	if cents.Sign() < 0 {
		half.Neg(half)
	}
	cents.Add(cents, half)
	q := new(big.Int).Quo(cents.Num(), cents.Denom()) // truncates toward zero
	if !q.IsInt64() {
		return 0, fmt.Errorf("%q is too large", text)
	}
	return Money(q.Int64()), nil
}

// ParseQty turns "3" into 3. Steps rarely need it: godog converts a QtyPattern
// group to int when the step function takes an int.
func ParseQty(text string) (int, error) {
	return strconv.Atoi(text)
}
