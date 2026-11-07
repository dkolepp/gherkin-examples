"""Checkout steps (owned by the checkout team). behave version.

Only checkout actions live here. The journey reuses the cart, payments and
shipping teams' steps for everything else, instead of writing its own
versions of them.
"""

from behave import given, register_type, then, when

from bdd_support.types import TYPES

register_type(**TYPES)


@given('the shopper lives in "{country}"')
def step_lives_in(context, country):
    """Set the country the order will be shipped to."""
    context.checkout.live_in(country)


@when('I check out paying with card "{card}"')
def step_check_out(context, card):
    """Place the order for the current cart. Fails the scenario if refused."""
    context.checkout.check_out(card)


@when('I try to check out paying with card "{card}"')
def step_try_check_out(context, card):
    """Try to place the order, which may be refused by payments or shipping.

    Follow with 'Then I should be told "..."' to check the reason.
    """
    context.outcome.attempt(context.checkout.check_out, card)


@then("the order total should be {amount:Money}")
def step_order_total(context, amount):
    """Check what the order cost in total: goods after discount, plus shipping."""
    total = context.checkout.order.total
    assert total == amount, f"order total was {total}"
