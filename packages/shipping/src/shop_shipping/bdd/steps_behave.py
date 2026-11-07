"""Shipping steps (owned by the shipping team). behave version.

Contract
    install(context)  sets context.shipping (a ShippingDriver)
    requires          context.outcome, from shop_bdd_common
"""

from behave import given, then, when

import shop_bdd_common.steps_behave  # noqa: F401  registers Money/Qty types
from shop_bdd_common.deprecation import deprecated
from shop_shipping.bdd.driver import ShippingDriver


def install(context) -> None:
    context.shipping = ShippingDriver()


@given("shipping costs:")
def step_shipping_costs(context):
    """Set a flat shipping rate per country.

    Columns: country, rate. Countries not listed can't be shipped to.
    """
    for row in context.table:
        context.shipping.set_rate(row["country"], row["rate"])


@given("orders of {amount:Money} or more ship free")
def step_free_from(context, amount):
    """Make shipping free for orders at or above this value."""
    context.shipping.free_shipping_from(amount)


@given('shipping is set up from "{file_name}"')
def step_shipping_from_file(context, file_name):
    """Load shipping rates and the free-shipping threshold from a JSON file.

    Bundled: "standard_shipping.json". Or give the path of your own file.
    """
    context.shipping.load(file_name)


@when('I get a shipping quote to "{country}" for {amount:Money}')
def step_quote(context, country, amount):
    """Ask what shipping costs for an order of this value."""
    context.shipping.quote(country, amount)


@when('I try to get a shipping quote to "{country}" for {amount:Money}')
def step_try_quote(context, country, amount):
    """Ask for a shipping quote, which may be refused.

    Follow with 'Then I should be told "..."' to check the reason.
    """
    context.outcome.attempt(context.shipping.quote, country, amount)


@then("shipping should cost {amount:Money}")
def step_shipping_cost(context, amount):
    """Check the last shipping quote."""
    assert context.shipping.last_quote == amount, f"quote was {context.shipping.last_quote}"


@then('a shipment of {qty:Qty} items should be sent to "{country}"')
def step_shipment_sent(context, qty, country):
    """Check that exactly one shipment went out, with this many items."""
    shipments = context.shipping.shipments()
    assert len(shipments) == 1, f"{len(shipments)} shipments were sent"
    assert (shipments[0].items, shipments[0].country) == (qty, country), shipments[0]


@then("no shipment should be sent")
def step_no_shipment(context):
    """Check that nothing was shipped."""
    assert context.shipping.shipments() == [], context.shipping.shipments()


# Version 1.0.0 used this sentence. It was too loose: it didn't say "exactly
# one shipment". Kept working until 2.0.0 so existing feature files don't break.
@then('{qty:Qty} items should be shipped to "{country}"')
@deprecated(old='{qty:Qty} items should be shipped to "{country}"', new='a shipment of {qty:Qty} items should be sent to "{country}"', since="1.1.0")
def step_shipped_old(context, qty, country):
    """Deprecated since 1.1.0: use 'a shipment of N items should be sent to "..."'."""
    step_shipment_sent(context, qty, country)
