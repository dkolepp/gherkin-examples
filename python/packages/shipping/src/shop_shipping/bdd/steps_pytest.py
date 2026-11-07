"""Shipping steps (owned by the shipping team). pytest-bdd version.

Contract
    fixture `shipping`  a ShippingDriver, fresh for each scenario
    requires            fixture `outcome`, from shop_bdd_common
"""

import pytest
from pytest_bdd import given, then, when

from shop_bdd_common.deprecation import deprecated
from shop_bdd_common.pytest_support import as_dicts, p
from shop_shipping.bdd.driver import ShippingDriver


@pytest.fixture
def shipping():
    return ShippingDriver()


@given("shipping costs:")
def shipping_costs(shipping, datatable):
    """Set a flat shipping rate per country.

    Columns: country, rate. Countries not listed can't be shipped to.
    """
    for row in as_dicts(datatable):
        shipping.set_rate(row["country"], row["rate"])


@given(p("orders of {amount:Money} or more ship free"))
def free_from(shipping, amount):
    """Make shipping free for orders at or above this value."""
    shipping.free_shipping_from(amount)


@given(p('shipping is set up from "{file_name}"'))
def shipping_from_file(shipping, file_name):
    """Load shipping rates and the free-shipping threshold from a JSON file.

    Bundled: "standard_shipping.json". Or give the path of your own file.
    """
    shipping.load(file_name)


@when(p('I get a shipping quote to "{country}" for {amount:Money}'))
def quote(shipping, country, amount):
    """Ask what shipping costs for an order of this value."""
    shipping.quote(country, amount)


@when(p('I try to get a shipping quote to "{country}" for {amount:Money}'))
def try_quote(shipping, outcome, country, amount):
    """Ask for a shipping quote, which may be refused.

    Follow with 'Then I should be told "..."' to check the reason.
    """
    outcome.attempt(shipping.quote, country, amount)


@then(p("shipping should cost {amount:Money}"))
def shipping_cost(shipping, amount):
    """Check the last shipping quote."""
    assert shipping.last_quote == amount


@then(p('a shipment of {qty:Qty} items should be sent to "{country}"'))
def shipment_sent(shipping, qty, country):
    """Check that exactly one shipment went out, with this many items."""
    shipments = shipping.shipments()
    assert len(shipments) == 1
    assert (shipments[0].items, shipments[0].country) == (qty, country)


@then("no shipment should be sent")
def no_shipment(shipping):
    """Check that nothing was shipped."""
    assert shipping.shipments() == []


# Version 1.0.0 used this sentence. Kept working until 2.0.0.
@then(p('{qty:Qty} items should be shipped to "{country}"'))
@deprecated(old='{qty:Qty} items should be shipped to "{country}"', new='a shipment of {qty:Qty} items should be sent to "{country}"', since="1.1.0")
def shipped_old(shipping, qty, country):
    """Deprecated since 1.1.0: use 'a shipment of N items should be sent to "..."'."""
    shipment_sent(shipping, qty, country)
