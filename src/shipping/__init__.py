"""Shipping domain (owned by the shipping team).

Flat rates per country, free shipping above a threshold, and a record of
shipments sent.
"""

from dataclasses import dataclass
from decimal import Decimal


class UnsupportedDestination(ValueError):
    def __init__(self, country: str):
        super().__init__(f"We don't ship to {country}")


@dataclass(frozen=True)
class Shipment:
    id: str
    order_id: str
    country: str
    items: int


class ShippingService:
    def __init__(self):
        self._rates: dict[str, Decimal] = {}
        self.free_from: Decimal | None = None
        self.shipments: list[Shipment] = []

    def set_rate(self, country: str, rate: Decimal) -> None:
        self._rates[country] = rate

    def quote(self, country: str, order_value: Decimal) -> Decimal:
        if country not in self._rates:
            raise UnsupportedDestination(country)
        if self.free_from is not None and order_value >= self.free_from:
            return Decimal("0.00")
        return self._rates[country]

    def ship(self, order_id: str, country: str, items: int) -> Shipment:
        shipment = Shipment(f"SHP-{len(self.shipments) + 1:04d}", order_id, country, items)
        self.shipments.append(shipment)
        return shipment
