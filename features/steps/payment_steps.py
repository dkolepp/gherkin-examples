"""Payment steps (owned by the payments team). behave version."""

from behave import given, register_type, then, when

from bdd_support.types import TYPES

register_type(**TYPES)


@given('the shopper has a card "{card}" with a limit of {limit:Money}')
def step_card(context, card, limit):
    """Give the shopper a card that can be charged up to the limit in total."""
    context.payments.issue_card(card, limit)


@when('I charge {amount:Money} to card "{card}"')
def step_charge(context, amount, card):
    """Charge a card directly. Fails the scenario if the payment is declined."""
    context.payments.charge(card, amount)


@when('I try to charge {amount:Money} to card "{card}"')
def step_try_charge(context, amount, card):
    """Try to charge a card, which may be declined.

    Follow with 'Then I should be told "..."' to check the reason.
    """
    context.outcome.attempt(context.payments.charge, card, amount)


@then('card "{card}" should be charged {amount:Money}')
def step_charged(context, card, amount):
    """Check the total of all successful payments on a card."""
    actual = context.payments.amount_charged(card)
    assert actual == amount, f"card was charged {actual}"


@then('card "{card}" should not be charged')
def step_not_charged(context, card):
    """Check that no payment went through on a card."""
    actual = context.payments.amount_charged(card)
    assert actual == 0, f"card was charged {actual}"
