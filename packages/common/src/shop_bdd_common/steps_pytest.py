"""Shared steps and fixtures for pytest-bdd.

Use with a star import in conftest.py:  from shop_bdd_common.steps_pytest import *

Contract
    fixture `outcome`  an Outcome, fresh for each scenario
"""

import pytest
from pytest_bdd import then

from shop_bdd_common.outcome import Outcome
from shop_bdd_common.pytest_support import p


@pytest.fixture
def outcome():
    return Outcome()


@then(p('I should be told "{message}"'))
def told(outcome, message):
    """Check the message the shopper saw when the last 'try to' step was refused."""
    assert outcome.last_error is not None, "expected an error, but the action succeeded"
    assert str(outcome.last_error) == message
