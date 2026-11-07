"""Checkout (owned by the checkout team).

Connects the other domains through their public APIs only:
cart total -> shipping quote -> card payment -> shipment.
"""

from dataclasses import dataclass
from decimal import Decimal

from shop_cart import Cart
from shop_payments import Payment, PaymentGateway
from shop_shipping import Shipment, ShippingService


class CheckoutError(ValueError):
    pass


@dataclass(frozen=True)
class Order:
    id: str
    goods: Decimal
    shipping: Decimal
    total: Decimal
    payment: Payment
    shipment: Shipment


class Checkout:
    def __init__(self, payments: PaymentGateway, shipping: ShippingService):
        self.payments = payments
        self.shipping = shipping
        self.orders: list[Order] = []

    def place_order(self, cart: Cart, card: str, country: str) -> Order:
        if cart.is_empty:
            raise CheckoutError("Your cart is empty")
        goods = cart.total  # after any discount
        shipping_cost = self.shipping.quote(country, goods)  # may refuse
        total = goods + shipping_cost
        payment = self.payments.charge(card, total)  # may decline: nothing ships
        order_id = f"ORD-{len(self.orders) + 1:04d}"
        items = sum(qty for _, qty, _ in cart.lines())
        shipment = self.shipping.ship(order_id, country, items)
        order = Order(order_id, goods, shipping_cost, total, payment, shipment)
        self.orders.append(order)
        cart.clear()
        return order
