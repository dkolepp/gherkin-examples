from dataclasses import dataclass
from decimal import Decimal

from shop_cart.money import format_money, to_money


class UnknownDiscountCode(LookupError):
    def __init__(self, code: str):
        super().__init__(f"{code} is not a valid discount code")


class DiscountNotApplicable(ValueError):
    pass


@dataclass(frozen=True)
class Discount:
    code: str
    percent_off: Decimal
    minimum_spend: Decimal = Decimal("0")

    def is_eligible(self, subtotal: Decimal) -> bool:
        return subtotal >= self.minimum_spend

    def check_eligible(self, subtotal: Decimal) -> None:
        if not self.is_eligible(subtotal):
            raise DiscountNotApplicable(
                f"{self.code} needs a minimum spend of {format_money(self.minimum_spend)}"
            )

    def amount_for(self, subtotal: Decimal) -> Decimal:
        if not self.is_eligible(subtotal):
            return to_money(0)
        return to_money(subtotal * self.percent_off / 100)


class DiscountBook:
    """All discount codes the shop accepts."""

    def __init__(self, discounts=()):
        self._codes = {d.code: d for d in discounts}

    def add(self, discount: Discount) -> None:
        self._codes[discount.code] = discount

    def lookup(self, code: str) -> Discount:
        try:
            return self._codes[code]
        except KeyError:
            raise UnknownDiscountCode(code) from None
