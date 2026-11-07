"""pytest-bdd: make the published step libraries (and their fixtures) available."""

from shop_bdd_common.steps_pytest import *  # noqa: F401,F403
from shop_cart.bdd.steps_pytest import *  # noqa: F401,F403
from shop_payments.bdd.steps_pytest import *  # noqa: F401,F403
from shop_shipping.bdd.steps_pytest import *  # noqa: F401,F403
from shop_checkout.bdd.steps_pytest import *  # noqa: F401,F403
