"""Checkout steps (owned by the checkout team). pytest-bdd version.

Only checkout actions live here. The journey reuses the cart, payments and
shipping teams' steps for everything else.
"""

from pytest_bdd import given, then, when

from .helpers import p


@given(p('the shopper lives in "{country}"'))
def lives_in(checkout, country):
    """Set the country the order will be shipped to."""
    checkout.live_in(country)


@when(p('I check out paying with card "{card}"'))
def check_out(checkout, card):
    """Place the order for the current cart. Fails the scenario if refused."""
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
