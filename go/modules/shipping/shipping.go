// Package shipping is the shipping domain (owned by the shipping team).
//
// Flat rates per country, free shipping above a threshold, and a record of
// shipments sent. Amounts are int64 counts of cents.
package shipping

import "fmt"

// UnsupportedDestinationError says the shop doesn't ship to a country.
type UnsupportedDestinationError struct{ Country string }

func (e *UnsupportedDestinationError) Error() string { return "We don't ship to " + e.Country }

// Shipment is an order sent to a country.
type Shipment struct {
	ID      string
	OrderID string
	Country string
	Items   int
}

// Service quotes shipping costs and sends shipments.
type Service struct {
	rates     map[string]int64
	FreeFrom  *int64 // cents; nil means shipping is never free
	Shipments []Shipment
}

// NewService returns a service with no rates.
func NewService() *Service { return &Service{rates: map[string]int64{}} }

// SetRate sets the flat rate for a country, in cents.
func (s *Service) SetRate(country string, rate int64) { s.rates[country] = rate }

// SetFreeFrom makes orders of at least this value (in cents) ship free.
func (s *Service) SetFreeFrom(cents int64) { s.FreeFrom = &cents }

// Quote is what shipping costs to a country for an order of this value.
func (s *Service) Quote(country string, orderValue int64) (int64, error) {
	rate, ok := s.rates[country]
	if !ok {
		return 0, &UnsupportedDestinationError{country}
	}
	if s.FreeFrom != nil && orderValue >= *s.FreeFrom {
		return 0, nil
	}
	return rate, nil
}

// Ship sends an order's items to a country.
func (s *Service) Ship(orderID, country string, items int) Shipment {
	sh := Shipment{fmt.Sprintf("SHP-%04d", len(s.Shipments)+1), orderID, country, items}
	s.Shipments = append(s.Shipments, sh)
	return sh
}
