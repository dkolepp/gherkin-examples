"""behave hooks for the checkout team's own feature files.

Each step library has an install(context) that puts its driver on the context.
Order matters: a library may need what an earlier one installed.
"""

import shop_bdd_common.steps_behave
import shop_cart.bdd.steps_behave
import shop_payments.bdd.steps_behave
import shop_shipping.bdd.steps_behave
import shop_checkout.bdd.steps_behave

LIBRARIES = [
    shop_bdd_common.steps_behave,
    shop_cart.bdd.steps_behave,
    shop_payments.bdd.steps_behave,
    shop_shipping.bdd.steps_behave,
    shop_checkout.bdd.steps_behave,
]


def before_scenario(context, scenario):
    for library in LIBRARIES:
        library.install(context)
