package cart

import "fmt"

// UnknownDiscountCodeError says a code doesn't exist.
type UnknownDiscountCodeError struct{ Code string }

func (e *UnknownDiscountCodeError) Error() string {
	return e.Code + " is not a valid discount code"
}

// DiscountNotApplicableError says a code can't be used on this cart.
type DiscountNotApplicableError struct{ Message string }

func (e *DiscountNotApplicableError) Error() string { return e.Message }

// Discount takes a percentage off the subtotal, above a minimum spend.
type Discount struct {
	Code         string
	PercentOff   int64 // whole percent
	MinimumSpend int64 // cents
}

// IsEligible reports whether the discount applies to a subtotal.
func (d Discount) IsEligible(subtotal int64) bool { return subtotal >= d.MinimumSpend }

// CheckEligible returns an error if the discount doesn't apply to a subtotal.
func (d Discount) CheckEligible(subtotal int64) error {
	if !d.IsEligible(subtotal) {
		return &DiscountNotApplicableError{
			fmt.Sprintf("%s needs a minimum spend of %s", d.Code, FormatMoney(d.MinimumSpend)),
		}
	}
	return nil
}

// AmountFor is how much the discount takes off a subtotal, in cents.
func (d Discount) AmountFor(subtotal int64) int64 {
	if !d.IsEligible(subtotal) {
		return 0
	}
	return (subtotal*d.PercentOff + 50) / 100 // round half up
}

// DiscountBook is all the discount codes the shop accepts.
type DiscountBook struct {
	codes map[string]Discount
}

// NewDiscountBook returns a book holding discounts.
func NewDiscountBook(discounts ...Discount) *DiscountBook {
	b := &DiscountBook{codes: map[string]Discount{}}
	for _, d := range discounts {
		b.Add(d)
	}
	return b
}

// Add adds or replaces a discount code.
func (b *DiscountBook) Add(d Discount) { b.codes[d.Code] = d }

// Lookup finds a discount by its code.
func (b *DiscountBook) Lookup(code string) (Discount, error) {
	d, ok := b.codes[code]
	if !ok {
		return Discount{}, &UnknownDiscountCodeError{code}
	}
	return d, nil
}
