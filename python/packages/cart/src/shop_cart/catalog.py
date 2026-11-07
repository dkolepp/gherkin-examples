from dataclasses import dataclass
from decimal import Decimal


class UnknownProduct(LookupError):
    def __init__(self, name: str):
        super().__init__(f"We don't sell {name}")


@dataclass(frozen=True)
class Product:
    name: str
    price: Decimal
    stock: int


class Catalog:
    """The products the shop sells, looked up by name."""

    def __init__(self, products=()):
        self._products: dict[str, Product] = {}
        for product in products:
            self.add(product)

    def add(self, product: Product) -> None:
        self._products[product.name] = product

    def get(self, name: str) -> Product:
        try:
            return self._products[name]
        except KeyError:
            raise UnknownProduct(name) from None
