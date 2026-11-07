"""Cart steps (owned by the cart team). pytest-bdd version.

Same sentences as steps_behave.py, same driver.
Use with a star import in conftest.py:  from shop_cart.bdd.steps_pytest import *

Contract
    fixture `shop`  a ShopDriver, fresh for each scenario
    requires        fixture `outcome`, from shop_bdd_common
"""

import pytest
from pytest_bdd import given, step, then, when

from shop_bdd_common.pytest_support import as_dicts, p
from shop_cart.bdd.driver import ShopDriver


@pytest.fixture
def shop():
    return ShopDriver()


# --- Given ------------------------------------------------------------------

@given("the shop sells these products:")
def shop_sells(shop, datatable):
    """Stock the shop with the products in the table.

    Columns: name, price, stock. A stock of 0 means sold out.
    """
    for row in as_dicts(datatable):
        shop.stock_product(row["name"], row["price"], row["stock"])


@given("these discount codes exist:")
def discount_codes(shop, datatable):
    """Create the discount codes in the table.

    Columns: code, percent off, minimum spend (use 0.00 for no minimum).
    """
    for row in as_dicts(datatable):
        shop.add_discount_code(row["code"], row["percent off"], row["minimum spend"])


@given(p('the shop is set up from "{file_name}"'))
def shop_from_file(shop, file_name):
    """Stock the shop from a JSON file of products and discount codes.

    Bundled: "standard_shop.json". Or give the path of your own file.
    """
    shop.load_shop(file_name)


@given("I have an empty cart")
def empty_cart(shop):
    """Start a new, empty cart."""
    shop.new_cart()


@given(p('I have a cart containing {qty:Qty} "{product}"'))
def cart_containing(shop, qty, product):
    """Start a new cart that already holds some of one product."""
    shop.new_cart()
    shop.add_to_cart(product, qty)


# --- When -------------------------------------------------------------------

@step(p('I add {qty:Qty} "{product}" to the cart'))
def add(shop, qty, product):
    """Put products in the cart. Fails the scenario if that isn't possible.

    Works as Given, When or And.
    """
    shop.add_to_cart(product, qty)


@when(p('I try to add {qty:Qty} "{product}" to the cart'))
def try_add(shop, outcome, qty, product):
    """Try to put products in the cart, which may be refused.

    Follow with 'Then I should be told "..."' to check the reason.
    """
    outcome.attempt(shop.add_to_cart, product, qty)


@when(p('I remove {qty:Qty} "{product}" from the cart'))
def remove(shop, qty, product):
    """Take products out of the cart. Removing more than is there empties that line."""
    shop.remove_from_cart(product, qty)


@when(p('I apply the discount code "{code}"'))
@given(p('I have applied the discount code "{code}"'))
def apply_code(shop, code):
    """Apply a discount code. Fails the scenario if the code is refused."""
    shop.apply_code(code)


@when(p('I try to apply the discount code "{code}"'))
def try_apply_code(shop, outcome, code):
    """Try a discount code, which may be refused.

    Follow with 'Then I should be told "..."' to check the reason.
    """
    outcome.attempt(shop.apply_code, code)


@when("I print the receipt")
def print_receipt(shop):
    """Print the receipt for the current cart."""
    shop.print_receipt()


# --- Then -------------------------------------------------------------------

@then("the cart should be empty")
def cart_empty(shop):
    """Check that the cart holds nothing."""
    assert shop.cart.is_empty


@then(p("the cart total should be {amount:Money}"))
def cart_total(shop, amount):
    """Check the amount to pay, after any discount. Write money as $12.50."""
    assert shop.cart.total == amount


@then(p("the discount should be {amount:Money}"))
def discount(shop, amount):
    """Check how much the discount takes off. $0.00 means no discount applies."""
    assert shop.cart.discount_amount == amount


@then("the cart should contain:")
def cart_contains_table(shop, datatable):
    """Check how many of each product are in the cart.

    Columns: product, quantity.
    """
    for row in as_dicts(datatable):
        assert shop.cart.quantity_of(row["product"]) == int(row["quantity"])


@then(p('the cart should contain {qty:Qty} "{product}"'))
def cart_contains(shop, qty, product):
    """Check how many of one product are in the cart."""
    assert shop.cart.quantity_of(product) == qty



@then("the receipt should read:")
def receipt(shop, docstring):
    """Check the printed receipt matches the text block below the step exactly."""
    assert shop.receipt == docstring
