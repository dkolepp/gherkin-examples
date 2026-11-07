"""Checkout driver (owned by the checkout team, who own the journey).

It builds the real Checkout service from the other teams' drivers: the cart
team's cart, the payments team's gateway and the shipping team's service.
Journey steps call this driver for checkout actions, and the other teams'
drivers for everything else.
"""

from bdd_support.cart_driver import ShopDriver
from bdd_support.payments_driver import PaymentsDriver
from bdd_support.shipping_driver import ShippingDriver
from checkout import Checkout, Order


class CheckoutDriver:
    def __init__(self, shop: ShopDriver, payments: PaymentsDriver, shipping: ShippingDriver):
        self.shop = shop
        self.checkout = Checkout(payments.gateway, shipping.service)
        self.country: str | None = None
        self.order: Order | None = None

    def live_in(self, country: str) -> None:
        self.country = country

    def check_out(self, card: str) -> None:
        self.order = self.checkout.place_order(self.shop.cart, card, self.country)
