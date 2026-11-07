"""Shared steps for behave. Every team's step library depends on this one.

Contract
    install(context)  sets context.outcome (an Outcome). Call it first.
"""

from behave import register_type, then

from shop_bdd_common.outcome import Outcome
from shop_bdd_common.types import TYPES

register_type(**TYPES)


def install(context) -> None:
    context.outcome = Outcome()


@then('I should be told "{message}"')
def step_told(context, message):
    """Check the message the shopper saw when the last 'try to' step was refused."""
    error = context.outcome.last_error
    assert error is not None, "expected an error, but the action succeeded"
    assert str(error) == message, f"got: {error}"
