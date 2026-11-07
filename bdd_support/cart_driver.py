"""Cart driver (owned by the cart team).

The ShopDriver: business-level actions, with no Gherkin and no runner code.

Step definitions stay one or two lines long because the real work lives here.
If the shop later gets a web UI or REST API, you can swap this driver for one
that clicks buttons or sends HTTP requests, and the .feature files and step
definitions do not change at all.
"""

import json
from decimal import Decimal
from pathlib import Path

from shop import (
    Cart,
    Catalog,
    Discount,
    DiscountBook,
    Product,
    render_receipt,
)
from shop.money import to_money

TEST_DATA = Path(__file__).resolve().parent.parent / "test_data"


class ShopDriver:
    def __init__(self):
        self.catalog = Catalog()
        self.discounts = DiscountBook()
        self.cart: Cart | None = None
        self.receipt: str | None = None

    # --- setting up the world -------------------------------------------
    def stock_product(self, name: str, price, stock) -> None:
        self.catalog.add(Product(name, to_money(price), int(stock)))

    def add_discount_code(self, code: str, percent_off, minimum_spend="0") -> None:
        self.discounts.add(
            Discount(code, Decimal(str(percent_off)), to_money(minimum_spend))
        )

    def load_shop(self, file_name: str) -> None:
        """Load products and discount codes from a JSON file in test_data/."""
        data = json.loads((TEST_DATA / file_name).read_text())
        for p in data["products"]:
            self.stock_product(p["name"], p["price"], p["stock"])
        for d in data.get("discount_codes", []):
            self.add_discount_code(d["code"], d["percent_off"], d.get("minimum_spend", "0"))

    # --- shopper actions ------------------------------------------------
    def new_cart(self) -> Cart:
        self.cart = Cart(self.catalog)
        return self.cart

    def add_to_cart(self, name: str, quantity: int) -> None:
        self.cart.add(name, quantity)

    def remove_from_cart(self, name: str, quantity: int) -> None:
        self.cart.remove(name, quantity)

    def apply_code(self, code: str) -> None:
        self.cart.apply_discount(self.discounts.lookup(code))

    def print_receipt(self) -> str:
        self.receipt = render_receipt(self.cart)
        return self.receipt
