"""Plain pytest unit tests, for contrast with the BDD scenarios.

BDD scenarios describe behaviour a shopper or product owner cares about.
Unit tests cover small technical details (rounding, edge cases) that would
only clutter a .feature file. You usually want both.
"""

from decimal import Decimal

import pytest

from shop import Cart, Catalog, Product, UnknownProduct
from shop.money import format_money, to_money


@pytest.mark.parametrize(
    "raw, expected",
    [("1.005", "1.01"), ("2", "2.00"), (Decimal("0.004"), "0.00")],
)
def test_to_money_rounds_half_up_to_cents(raw, expected):
    assert to_money(raw) == Decimal(expected)


def test_format_money():
    assert format_money(Decimal("3.5")) == "$3.50"


def test_removing_more_than_in_cart_just_empties_it():
    cart = Cart(Catalog([Product("Mug", Decimal("8.00"), 5)]))
    cart.add("Mug", 2)
    cart.remove("Mug", 10)
    assert cart.is_empty


def test_unknown_product_raises():
    with pytest.raises(UnknownProduct, match="We don't sell Teapot"):
        Cart(Catalog()).add("Teapot")
