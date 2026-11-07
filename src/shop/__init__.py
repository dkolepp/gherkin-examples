"""A tiny shop domain used as the "system under test".

Nothing in this package knows about Gherkin, behave or pytest. That is the
point: the application code stays clean, and tests talk to it from outside.
"""

from shop.catalog import Catalog, Product, UnknownProduct
from shop.cart import Cart, OutOfStock
from shop.discounts import (
    Discount,
    DiscountBook,
    DiscountNotApplicable,
    UnknownDiscountCode,
)
from shop.receipt import render_receipt

__all__ = [
    "Cart",
    "Catalog",
    "Discount",
    "DiscountBook",
    "DiscountNotApplicable",
    "OutOfStock",
    "Product",
    "UnknownDiscountCode",
    "UnknownProduct",
    "render_receipt",
]
