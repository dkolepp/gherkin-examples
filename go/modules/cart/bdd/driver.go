package bdd

// Cart driver (owned by the cart team). Part of the public contract.
//
// The ShopDriver: business-level actions, with no Gherkin and no runner code.
//
// Step definitions stay one or two lines long because the real work lives
// here. If the shop later gets a web UI or REST API, you can swap this driver
// for one that clicks buttons or sends HTTP requests, and the .feature files
// and step definitions do not change at all.

import (
	"embed"
	"encoding/json"

	"example.com/shop/bddcommon"
	"example.com/shop/cart"
)

//go:embed data/*.json
var bundled embed.FS

// ShopDriver drives the cart application.
type ShopDriver struct {
	Catalog   *cart.Catalog
	Discounts *cart.DiscountBook
	Cart      *cart.Cart
	Receipt   string
}

// NewShopDriver returns a driver for an empty shop.
func NewShopDriver() *ShopDriver {
	return &ShopDriver{Catalog: cart.NewCatalog(), Discounts: cart.NewDiscountBook()}
}

// --- setting up the world --------------------------------------------------

// StockProduct adds a product. price is text such as "12.50".
func (d *ShopDriver) StockProduct(name, price string, stock int) error {
	cents, err := bddcommon.ToMoney(price)
	if err != nil {
		return err
	}
	d.Catalog.Add(cart.Product{Name: name, Price: int64(cents), Stock: stock})
	return nil
}

// AddDiscountCode adds a discount code. minimumSpend is text such as "50.00".
func (d *ShopDriver) AddDiscountCode(code string, percentOff int, minimumSpend string) error {
	minimum, err := bddcommon.ToMoney(minimumSpend)
	if err != nil {
		return err
	}
	d.Discounts.Add(cart.Discount{Code: code, PercentOff: int64(percentOff), MinimumSpend: int64(minimum)})
	return nil
}

type shopData struct {
	Products []struct {
		Name  string      `json:"name"`
		Price json.Number `json:"price"`
		Stock int         `json:"stock"`
	} `json:"products"`
	DiscountCodes []struct {
		Code         string      `json:"code"`
		PercentOff   int         `json:"percent_off"`
		MinimumSpend json.Number `json:"minimum_spend"`
	} `json:"discount_codes"`
}

// LoadShop loads products and discount codes from a JSON file.
//
// Bundled: "standard_shop.json". A file of that name in the current folder
// is used instead, if there is one.
func (d *ShopDriver) LoadShop(fileName string) error {
	var data shopData
	if err := bddcommon.LoadJSON(bundled, fileName, &data); err != nil {
		return err
	}
	for _, p := range data.Products {
		if err := d.StockProduct(p.Name, p.Price.String(), p.Stock); err != nil {
			return err
		}
	}
	for _, c := range data.DiscountCodes {
		minimum := c.MinimumSpend.String()
		if minimum == "" {
			minimum = "0"
		}
		if err := d.AddDiscountCode(c.Code, c.PercentOff, minimum); err != nil {
			return err
		}
	}
	return nil
}

// --- shopper actions -------------------------------------------------------

// NewCart starts a new, empty cart.
func (d *ShopDriver) NewCart() *cart.Cart {
	d.Cart = cart.New(d.Catalog)
	return d.Cart
}

// AddToCart puts products in the cart.
func (d *ShopDriver) AddToCart(name string, quantity int) error {
	return d.Cart.Add(name, quantity)
}

// RemoveFromCart takes products out of the cart.
func (d *ShopDriver) RemoveFromCart(name string, quantity int) {
	d.Cart.Remove(name, quantity)
}

// ApplyCode applies a discount code to the cart.
func (d *ShopDriver) ApplyCode(code string) error {
	discount, err := d.Discounts.Lookup(code)
	if err != nil {
		return err
	}
	return d.Cart.ApplyDiscount(discount)
}

// PrintReceipt prints the receipt for the current cart.
func (d *ShopDriver) PrintReceipt() string {
	d.Receipt = cart.RenderReceipt(d.Cart)
	return d.Receipt
}
