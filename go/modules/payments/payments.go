// Package payments is the payments domain (owned by the payments team).
//
// A fake card gateway: each card has a spending limit, and a payment is
// declined if it would take the card over that limit. Amounts are int64
// counts of cents.
package payments

import "fmt"

// DeclinedError says a payment did not go through.
type DeclinedError struct{ Message string }

func (e *DeclinedError) Error() string { return e.Message }

// Payment is a successful charge.
type Payment struct {
	ID     string
	Card   string
	Amount int64 // cents
}

// Gateway charges cards.
type Gateway struct {
	limits   map[string]int64
	Payments []Payment
}

// NewGateway returns a gateway with no cards.
func NewGateway() *Gateway { return &Gateway{limits: map[string]int64{}} }

// AddCard registers a card that can be charged up to limit in total.
func (g *Gateway) AddCard(number string, limit int64) { g.limits[number] = limit }

// ChargedTo is the total of all successful payments on a card.
func (g *Gateway) ChargedTo(number string) int64 {
	var sum int64
	for _, p := range g.Payments {
		if p.Card == number {
			sum += p.Amount
		}
	}
	return sum
}

// Charge takes amount from a card, unless that would go over its limit.
func (g *Gateway) Charge(number string, amount int64) (Payment, error) {
	ending := number
	if len(ending) > 4 {
		ending = ending[len(ending)-4:]
	}
	limit, ok := g.limits[number]
	if !ok {
		return Payment{}, &DeclinedError{fmt.Sprintf("Card ending %s is not recognised", ending)}
	}
	if g.ChargedTo(number)+amount > limit {
		return Payment{}, &DeclinedError{fmt.Sprintf("Card ending %s was declined", ending)}
	}
	p := Payment{fmt.Sprintf("PAY-%04d", len(g.Payments)+1), number, amount}
	g.Payments = append(g.Payments, p)
	return p, nil
}
