"""behave hooks: code that runs around features and scenarios.

Think of these like pytest fixtures. Each scenario gets fresh drivers, so
scenarios never leak state into each other. Every team's driver is created
here; a scenario simply ignores the ones it doesn't use.
"""

from bdd_support.cart_driver import ShopDriver
from bdd_support.checkout_driver import CheckoutDriver
from bdd_support.outcome import Outcome
from bdd_support.payments_driver import PaymentsDriver
from bdd_support.shipping_driver import ShippingDriver


def before_scenario(context, scenario):
    context.outcome = Outcome()
    context.shop = ShopDriver()
    context.payments = PaymentsDriver()
    context.shipping = ShippingDriver()
    context.checkout = CheckoutDriver(context.shop, context.payments, context.shipping)
