"""Custom placeholder types for step patterns.

Both behave and pytest-bdd use the `parse` library, so one set of converters
works for both:

    "the cart total should be {amount:Money}"   ->  amount == Decimal("33.00")
"""

import parse

from shop_bdd_common.money import to_money


@parse.with_pattern(r"\$\d+(?:\.\d{2})?")
def parse_money(text: str):
    return to_money(text.lstrip("$"))


@parse.with_pattern(r"\d+")
def parse_quantity(text: str) -> int:
    return int(text)


TYPES = {"Money": parse_money, "Qty": parse_quantity}
