"""pytest-bdd step definitions and fixtures.

These match the SAME sentences as features/steps/shop_steps.py (behave), and
call the SAME ShopDriver. Only the glue differs:

    behave                         pytest-bdd
    ------                         ----------
    context.shop                   a `shop` fixture
    context.table                  a `datatable` argument (list of rows)
    context.text                   a `docstring` argument
    "{qty:Qty}" + register_type    parsers.parse("{qty:Qty}", extra_types=...)

Steps live in conftest.py so every test_*.py file in this folder can use them.
"""

import pytest
from pytest_bdd import given, parsers, step, then, when

from bdd_support.driver import ShopDriver
from bdd_support.types import TYPES


def p(pattern: str):
    """Shorthand: a parse-style step pattern that knows our Money/Qty types."""
    return parsers.parse(pattern, extra_types=TYPES)


def as_dicts(datatable):
    """pytest-bdd gives a table as a list of lists; the first row is the header."""
    header, *rows = datatable
    return [dict(zip(header, row)) for row in rows]


@pytest.fixture
def shop():
    # A new driver for every scenario, just like before_scenario in behave.
    return ShopDriver()


# --- Given ------------------------------------------------------------------

@given("the shop sells these products:")
def shop_sells(shop, datatable):
    for row in as_dicts(datatable):
        shop.stock_product(row["name"], row["price"], row["stock"])


@given("these discount codes exist:")
def discount_codes(shop, datatable):
    for row in as_dicts(datatable):
        shop.add_discount_code(row["code"], row["percent off"], row["minimum spend"])


@given(p('the shop is set up from "{file_name}"'))
def shop_from_file(shop, file_name):
    shop.load_shop(file_name)


@given("I have an empty cart")
def empty_cart(shop):
    shop.new_cart()


@given(p('I have a cart containing {qty:Qty} "{product}"'))
def cart_containing(shop, qty, product):
    shop.new_cart()
    shop.add_to_cart(product, qty)


# --- When -------------------------------------------------------------------

@step(p('I add {qty:Qty} "{product}" to the cart'))
def add(shop, qty, product):
    shop.add_to_cart(product, qty)


@when(p('I try to add {qty:Qty} "{product}" to the cart'))
def try_add(shop, qty, product):
    shop.attempt(shop.add_to_cart, product, qty)


@when(p('I remove {qty:Qty} "{product}" from the cart'))
def remove(shop, qty, product):
    shop.remove_from_cart(product, qty)


@when(p('I apply the discount code "{code}"'))
@given(p('I have applied the discount code "{code}"'))
def apply_code(shop, code):
    shop.apply_code(code)


@when(p('I try to apply the discount code "{code}"'))
def try_apply_code(shop, code):
    shop.attempt(shop.apply_code, code)


@when("I print the receipt")
def print_receipt(shop):
    shop.print_receipt()


# --- Then -------------------------------------------------------------------

@then("the cart should be empty")
def cart_empty(shop):
    assert shop.cart.is_empty


@then(p("the cart total should be {amount:Money}"))
def cart_total(shop, amount):
    assert shop.cart.total == amount


@then(p("the discount should be {amount:Money}"))
def discount(shop, amount):
    assert shop.cart.discount_amount == amount


@then("the cart should contain:")
def cart_contains_table(shop, datatable):
    for row in as_dicts(datatable):
        assert shop.cart.quantity_of(row["product"]) == int(row["quantity"])


@then(p('the cart should contain {qty:Qty} "{product}"'))
def cart_contains(shop, qty, product):
    assert shop.cart.quantity_of(product) == qty


@then(p('I should be told "{message}"'))
def told(shop, message):
    assert shop.last_error is not None, "expected an error, but the action succeeded"
    assert str(shop.last_error) == message


@then("the receipt should read:")
def receipt(shop, docstring):
    assert shop.receipt == docstring
