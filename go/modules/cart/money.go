package cart

import (
	"fmt"
	"math/big"
)

// ToCents turns a number written as text ("12.50", "2", "1.005") into cents,
// rounding half up.
func ToCents(text string) (int64, error) {
	r, ok := new(big.Rat).SetString(text)
	if !ok {
		return 0, fmt.Errorf("%q is not an amount of money", text)
	}
	return roundHalfUp(r.Mul(r, big.NewRat(100, 1))), nil
}

// FormatMoney formats cents as $12.50.
func FormatMoney(cents int64) string {
	return "$" + Decimal(cents)
}

// Decimal formats cents as 12.50 (or -3.30), without a currency sign.
func Decimal(cents int64) string {
	sign := ""
	if cents < 0 {
		sign, cents = "-", -cents
	}
	return fmt.Sprintf("%s%d.%02d", sign, cents/100, cents%100)
}

// roundHalfUp rounds r to the nearest whole number, halves away from zero.
func roundHalfUp(r *big.Rat) int64 {
	half := big.NewRat(1, 2)
	if r.Sign() < 0 {
		half.Neg(half)
	}
	r = new(big.Rat).Add(r, half)
	return new(big.Int).Quo(r.Num(), r.Denom()).Int64()
}
