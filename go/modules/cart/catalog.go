package cart

// Product is something the shop sells.
type Product struct {
	Name  string
	Price int64 // cents
	Stock int
}

// UnknownProductError says the shop doesn't sell a product.
type UnknownProductError struct{ Name string }

func (e *UnknownProductError) Error() string { return "We don't sell " + e.Name }

// Catalog is the products the shop sells, looked up by name.
type Catalog struct {
	products map[string]Product
}

// NewCatalog returns a catalog holding products.
func NewCatalog(products ...Product) *Catalog {
	c := &Catalog{products: map[string]Product{}}
	for _, p := range products {
		c.Add(p)
	}
	return c
}

// Add adds or replaces a product.
func (c *Catalog) Add(p Product) { c.products[p.Name] = p }

// Get looks a product up by name.
func (c *Catalog) Get(name string) (Product, error) {
	p, ok := c.products[name]
	if !ok {
		return Product{}, &UnknownProductError{name}
	}
	return p, nil
}
