from decimal import Decimal

import pytest

from shop_bdd_common.deprecation import StepDeprecationWarning, deprecated
from shop_bdd_common.types import parse_money


def test_money_type_parses_dollars():
    assert parse_money("$12.5") == Decimal("12.50")


def test_deprecated_step_warns_and_still_runs():
    @deprecated(old="old words", new="new words", since="1.1.0")
    def step(x):
        return x * 2

    with pytest.warns(StepDeprecationWarning, match='Use "new words"'):
        assert step(2) == 4
