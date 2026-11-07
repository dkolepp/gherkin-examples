"""Checkout steps (owned by the checkout team). pytest-bdd version.

Only checkout actions live here. The journey reuses the cart, payments and
shipping teams' steps for everything else.

Contract
    fixture `checkout`  a CheckoutDriver, fresh for each scenario
    requires            fixtures `shop`, `payments`, `shipping`, `outcome`
"""

import pytest
from pytest_bdd import given, step, then, when

from shop_bdd_common.pytest_support import p
from shop_checkout.bdd.driver import CheckoutDriver


@pytest.fixture
def checkout(shop, payments, shipping):
    return CheckoutDriver(shop, payments, shipping)


@given(p('the shopper lives in "{country}"'))
def lives_in(checkout, country):
    """Set the country the order will be shipped to."""
    checkout.live_in(country)


@step(p('I check out paying with card "{card}"'))
def check_out(checkout, card):
    """Place the order for the current cart. Fails the scenario if refused.

    Works as Given, When or And.
    """
    checkout.check_out(card)


@when(p('I try to check out paying with card "{card}"'))
def try_check_out(checkout, outcome, card):
    """Try to place the order, which may be refused by payments or shipping.

    Follow with 'Then I should be told "..."' to check the reason.
    """
    outcome.attempt(checkout.check_out, card)


@then(p("the order total should be {amount:Money}"))
def order_total(checkout, amount):
    """Check what the order cost in total: goods after discount, plus shipping."""
    assert checkout.order.total == amount


@then(p('the order number should be "{order_id}"'))
def order_number(checkout, order_id):
    """Check the number given to the last order placed."""
    assert checkout.order.id == order_id
