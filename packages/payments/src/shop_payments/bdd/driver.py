"""Payments driver (owned by the payments team). Part of the public contract.

Other teams' steps use these methods instead of touching PaymentGateway
directly, so the payments team can change how payments work inside.
"""

from decimal import Decimal

from shop_bdd_common.money import to_money
from shop_payments import PaymentGateway


class PaymentsDriver:
    def __init__(self):
        self.gateway = PaymentGateway()

    def issue_card(self, number: str, limit) -> None:
        self.gateway.add_card(number, to_money(limit))

    def charge(self, number: str, amount) -> None:
        self.gateway.charge(number, to_money(amount))

    def amount_charged(self, number: str) -> Decimal:
        return self.gateway.charged_to(number)
