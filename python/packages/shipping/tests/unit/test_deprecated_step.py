"""The 1.0.0 sentence must keep working in 1.1.x, with a warning."""

import pytest

from shop_bdd_common.deprecation import StepDeprecationWarning
from shop_shipping.bdd.driver import ShippingDriver
from shop_shipping.bdd.steps_pytest import shipped_old


def test_old_sentence_still_works_but_warns():
    shipping = ShippingDriver()
    shipping.set_rate("Canada", "12.00")
    shipping.service.ship("ORD-0001", "Canada", 3)
    with pytest.warns(StepDeprecationWarning, match="deprecated since 1.1.0"):
        shipped_old(shipping, 3, "Canada")
