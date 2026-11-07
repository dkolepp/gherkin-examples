"""Give every scenario fresh drivers from each team's step library.

The journey itself is specs/journeys/checkout.feature at the repo root;
behave.ini lists this folder first so behave finds this file and steps/.

Order matters: checkout is built from the cart, payments and shipping drivers.
"""

import shop_bdd_common.steps_behave as common
import shop_cart.bdd.steps_behave as cart
import shop_checkout.bdd.steps_behave as checkout
import shop_payments.bdd.steps_behave as payments
import shop_shipping.bdd.steps_behave as shipping

LIBRARIES = [common, cart, payments, shipping, checkout]


def before_scenario(context, scenario):
    for library in LIBRARIES:
        library.install(context)
