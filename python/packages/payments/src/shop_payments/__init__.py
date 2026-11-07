"""Payments domain (owned by the payments team).

A fake card gateway: each card has a spending limit, and a payment is
declined if it would take the card over that limit.
"""

from dataclasses import dataclass
from decimal import Decimal


class PaymentDeclined(ValueError):
    pass


@dataclass(frozen=True)
class Payment:
    id: str
    card: str
    amount: Decimal


class PaymentGateway:
    def __init__(self):
        self._limits: dict[str, Decimal] = {}
        self.payments: list[Payment] = []

    def add_card(self, number: str, limit: Decimal) -> None:
        self._limits[number] = limit

    def charged_to(self, number: str) -> Decimal:
        return sum((p.amount for p in self.payments if p.card == number), Decimal("0.00"))

    def charge(self, number: str, amount: Decimal) -> Payment:
        ending = number[-4:]
        if number not in self._limits:
            raise PaymentDeclined(f"Card ending {ending} is not recognised")
        if self.charged_to(number) + amount > self._limits[number]:
            raise PaymentDeclined(f"Card ending {ending} was declined")
        payment = Payment(f"PAY-{len(self.payments) + 1:04d}", number, amount)
        self.payments.append(payment)
        return payment
