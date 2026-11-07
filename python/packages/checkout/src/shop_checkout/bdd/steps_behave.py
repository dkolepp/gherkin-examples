"""Checkout steps (owned by the checkout team). behave version.

Only checkout actions live here. The journey reuses the cart, payments and
shipping teams' steps for everything else, instead of writing its own
versions of them.

Contract
    install(context)  sets context.checkout (a CheckoutDriver)
    requires          context.shop, context.payments, context.shipping and
                      context.outcome: install those libraries first
"""

from behave import given, step, then, when

import shop_bdd_common.steps_behave  # noqa: F401  registers Money/Qty types
from shop_checkout.bdd.driver import CheckoutDriver


def install(context) -> None:
    context.checkout = CheckoutDriver(context.shop, context.payments, context.shipping)


@given('the shopper lives in "{country}"')
def step_lives_in(context, country):
    """Set the country the order will be shipped to."""
    context.checkout.live_in(country)


@step('I check out paying with card "{card}"')
def step_check_out(context, card):
    """Place the order for the current cart. Fails the scenario if refused.

    Works as Given, When or And.
    """
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


@then('the order number should be "{order_id}"')
def step_order_number(context, order_id):
    """Check the number given to the last order placed."""
    actual = context.checkout.order.id
    assert actual == order_id, f"order number was {actual}"
