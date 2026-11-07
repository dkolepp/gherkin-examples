package cart

import "fmt"

// OutOfStockError says the shop can't supply that many of a product.
type OutOfStockError struct{ Product Product }

func (e *OutOfStockError) Error() string {
	return fmt.Sprintf("Only %d %s left in stock", e.Product.Stock, e.Product.Name)
}

// Line is one product in the cart.
type Line struct {
	Product  Product
	Quantity int
	Total    int64 // cents
}

// Cart holds the products a shopper wants to buy.
type Cart struct {
	catalog    *Catalog
	order      []string // product names, in the order they were added
	quantities map[string]int
	Discount   *Discount
}

// New returns an empty cart for products from catalog.
func New(catalog *Catalog) *Cart {
	return &Cart{catalog: catalog, quantities: map[string]int{}}
}

// --- contents ------------------------------------------------------------

// Add puts quantity of a product in the cart, if there is enough stock.
func (c *Cart) Add(name string, quantity int) error {
	p, err := c.catalog.Get(name)
	if err != nil {
		return err
	}
	wanted := c.QuantityOf(name) + quantity
	if wanted > p.Stock {
		return &OutOfStockError{p}
	}
	if _, ok := c.quantities[name]; !ok {
		c.order = append(c.order, name)
	}
	c.quantities[name] = wanted
	return nil
}

// Remove takes quantity of a product out. Removing more than is there
// empties that line.
func (c *Cart) Remove(name string, quantity int) {
	left := c.QuantityOf(name) - quantity
	if left > 0 {
		c.quantities[name] = left
		return
	}
	if _, ok := c.quantities[name]; !ok {
		return
	}
	delete(c.quantities, name)
	for i, n := range c.order {
		if n == name {
			c.order = append(c.order[:i], c.order[i+1:]...)
			break
		}
	}
}

// Clear empties the cart and drops any discount.
func (c *Cart) Clear() {
	c.order = nil
	c.quantities = map[string]int{}
	c.Discount = nil
}

// QuantityOf is how many of a product are in the cart.
func (c *Cart) QuantityOf(name string) int { return c.quantities[name] }

// IsEmpty reports whether the cart holds nothing.
func (c *Cart) IsEmpty() bool { return len(c.quantities) == 0 }

// Lines lists the cart's contents in the order products were added.
func (c *Cart) Lines() []Line {
	lines := make([]Line, 0, len(c.order))
	for _, name := range c.order {
		p, _ := c.catalog.Get(name) // only products from the catalog get in
		q := c.quantities[name]
		lines = append(lines, Line{p, q, p.Price * int64(q)})
	}
	return lines
}

// --- money ---------------------------------------------------------------

// ApplyDiscount applies a discount, if the subtotal is eligible.
func (c *Cart) ApplyDiscount(d Discount) error {
	if err := d.CheckEligible(c.Subtotal()); err != nil {
		return err
	}
	c.Discount = &d
	return nil
}

// Subtotal is the cost of the goods before any discount, in cents.
func (c *Cart) Subtotal() int64 {
	var sum int64
	for _, l := range c.Lines() {
		sum += l.Total
	}
	return sum
}

// DiscountAmount is how much the discount takes off, in cents. It is 0 when
// there is no discount or the subtotal has dropped below its minimum spend.
func (c *Cart) DiscountAmount() int64 {
	if c.Discount == nil {
		return 0
	}
	return c.Discount.AmountFor(c.Subtotal())
}

// Total is the amount to pay, after any discount, in cents.
func (c *Cart) Total() int64 { return c.Subtotal() - c.DiscountAmount() }
