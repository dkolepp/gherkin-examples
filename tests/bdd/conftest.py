"""pytest-bdd fixtures, and the import that makes every team's steps available.

The steps live in step_defs/, one module per owning team, matching
features/steps/ for behave. Only the glue differs between the two runners:

    behave                         pytest-bdd
    ------                         ----------
    context.shop                   a `shop` fixture
    context.table                  a `datatable` argument (list of rows)
    context.text                   a `docstring` argument
    "{qty:Qty}" + register_type    parsers.parse("{qty:Qty}", extra_types=...)
"""

import pytest

from bdd_support.cart_driver import ShopDriver
from bdd_support.checkout_driver import CheckoutDriver
from bdd_support.outcome import Outcome
from bdd_support.payments_driver import PaymentsDriver
from bdd_support.shipping_driver import ShippingDriver

# pytest-bdd registers each step as a fixture in the module that defines it.
# Star-importing the step modules here makes them visible to every test file.
from step_defs.cart_steps import *  # noqa: E402,F401,F403
from step_defs.checkout_steps import *  # noqa: E402,F401,F403
from step_defs.common_steps import *  # noqa: E402,F401,F403
from step_defs.payment_steps import *  # noqa: E402,F401,F403
from step_defs.shipping_steps import *  # noqa: E402,F401,F403


# A fresh set of drivers for every scenario, like before_scenario in behave.

@pytest.fixture
def outcome():
    return Outcome()


@pytest.fixture
def shop():
    return ShopDriver()


@pytest.fixture
def payments():
    return PaymentsDriver()


@pytest.fixture
def shipping():
    return ShippingDriver()


@pytest.fixture
def checkout(shop, payments, shipping):
    return CheckoutDriver(shop, payments, shipping)
