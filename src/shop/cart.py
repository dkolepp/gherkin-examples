from decimal import Decimal

from shop.catalog import Catalog, Product
from shop.discounts import Discount
from shop.money import to_money


class OutOfStock(ValueError):
    def __init__(self, product: Product):
        super().__init__(f"Only {product.stock} {product.name} left in stock")


class Cart:
    def __init__(self, catalog: Catalog):
        self._catalog = catalog
        self._quantities: dict[str, int] = {}
        self.discount: Discount | None = None

    # --- contents -------------------------------------------------------
    def add(self, name: str, quantity: int = 1) -> None:
        product = self._catalog.get(name)
        wanted = self.quantity_of(name) + quantity
        if wanted > product.stock:
            raise OutOfStock(product)
        self._quantities[name] = wanted

    def remove(self, name: str, quantity: int = 1) -> None:
        left = self.quantity_of(name) - quantity
        if left > 0:
            self._quantities[name] = left
        else:
            self._quantities.pop(name, None)

    def quantity_of(self, name: str) -> int:
        return self._quantities.get(name, 0)

    @property
    def is_empty(self) -> bool:
        return not self._quantities

    def lines(self):
        """Yield (product, quantity, line_total) in the order items were added."""
        for name, qty in self._quantities.items():
            product = self._catalog.get(name)
            yield product, qty, to_money(product.price * qty)

    # --- money ----------------------------------------------------------
    def apply_discount(self, discount: Discount) -> None:
        discount.check_eligible(self.subtotal)
        self.discount = discount

    @property
    def subtotal(self) -> Decimal:
        return to_money(sum((line for _, _, line in self.lines()), Decimal("0")))

    @property
    def discount_amount(self) -> Decimal:
        if self.discount is None:
            return to_money(0)
        return self.discount.amount_for(self.subtotal)

    @property
    def total(self) -> Decimal:
        return self.subtotal - self.discount_amount
