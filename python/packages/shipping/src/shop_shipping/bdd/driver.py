"""Shipping driver (owned by the shipping team). Part of the public contract."""

from decimal import Decimal

from shop_bdd_common.data import load_json
from shop_bdd_common.money import to_money
from shop_shipping import Shipment, ShippingService


class ShippingDriver:
    def __init__(self):
        self.service = ShippingService()
        self.last_quote: Decimal | None = None

    def set_rate(self, country: str, rate) -> None:
        self.service.set_rate(country, to_money(rate))

    def free_shipping_from(self, amount) -> None:
        self.service.free_from = to_money(amount)

    def load(self, file_name: str) -> None:
        """Load rates and the free-shipping threshold from a JSON file.

        Bundled: "standard_shipping.json". A file of that name in the current
        folder is used instead, if there is one.
        """
        data = load_json("shop_shipping.bdd", file_name)
        for r in data["rates"]:
            self.set_rate(r["country"], r["rate"])
        if "free_shipping_from" in data:
            self.free_shipping_from(data["free_shipping_from"])

    def quote(self, country: str, order_value) -> None:
        self.last_quote = self.service.quote(country, to_money(order_value))

    def shipments(self) -> list[Shipment]:
        return list(self.service.shipments)
