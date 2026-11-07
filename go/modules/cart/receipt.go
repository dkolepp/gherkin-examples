package cart

import (
	"fmt"
	"strings"
)

// RenderReceipt prints the cart as a fixed-width receipt.
func RenderReceipt(c *Cart) string {
	var rows []string
	for _, l := range c.Lines() {
		rows = append(rows, fmt.Sprintf("%-14s%3d x %6s = %7s",
			l.Product.Name, l.Quantity, Decimal(l.Product.Price), Decimal(l.Total)))
	}
	rows = append(rows, fmt.Sprintf("%-30s%7s", "Subtotal", Decimal(c.Subtotal())))
	if amount := c.DiscountAmount(); amount != 0 {
		label := fmt.Sprintf("Discount (%s)", c.Discount.Code)
		rows = append(rows, fmt.Sprintf("%-30s%7s", label, Decimal(-amount)))
	}
	rows = append(rows, fmt.Sprintf("%-30s%7s", "TOTAL", Decimal(c.Total())))
	return strings.Join(rows, "\n")
}
