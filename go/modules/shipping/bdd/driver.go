package bdd

// Shipping driver (owned by the shipping team). Part of the public contract.

import (
	"embed"
	"encoding/json"

	"example.com/shop/bddcommon"
	"example.com/shop/shipping"
)

//go:embed data/*.json
var bundled embed.FS

// ShippingDriver drives the shipping application.
type ShippingDriver struct {
	Service   *shipping.Service
	LastQuote *bddcommon.Money
}

// NewShippingDriver returns a driver for a service with no rates.
func NewShippingDriver() *ShippingDriver {
	return &ShippingDriver{Service: shipping.NewService()}
}

// SetRate sets the flat rate for a country.
func (d *ShippingDriver) SetRate(country string, rate bddcommon.Money) {
	d.Service.SetRate(country, int64(rate))
}

// FreeShippingFrom makes orders of at least this value ship free.
func (d *ShippingDriver) FreeShippingFrom(amount bddcommon.Money) {
	d.Service.SetFreeFrom(int64(amount))
}

type shippingData struct {
	Rates []struct {
		Country string      `json:"country"`
		Rate    json.Number `json:"rate"`
	} `json:"rates"`
	FreeShippingFrom json.Number `json:"free_shipping_from"`
}

// Load loads rates and the free-shipping threshold from a JSON file.
//
// Bundled: "standard_shipping.json". A file of that name in the current
// folder is used instead, if there is one.
func (d *ShippingDriver) Load(fileName string) error {
	var data shippingData
	if err := bddcommon.LoadJSON(bundled, fileName, &data); err != nil {
		return err
	}
	for _, r := range data.Rates {
		rate, err := bddcommon.ToMoney(r.Rate.String())
		if err != nil {
			return err
		}
		d.SetRate(r.Country, rate)
	}
	if data.FreeShippingFrom != "" {
		free, err := bddcommon.ToMoney(data.FreeShippingFrom.String())
		if err != nil {
			return err
		}
		d.FreeShippingFrom(free)
	}
	return nil
}

// Quote asks what shipping costs, and keeps the answer in LastQuote.
func (d *ShippingDriver) Quote(country string, orderValue bddcommon.Money) error {
	cost, err := d.Service.Quote(country, int64(orderValue))
	if err != nil {
		return err
	}
	quote := bddcommon.Money(cost)
	d.LastQuote = &quote
	return nil
}

// Shipments lists the shipments sent so far.
func (d *ShippingDriver) Shipments() []shipping.Shipment {
	return append([]shipping.Shipment(nil), d.Service.Shipments...)
}
