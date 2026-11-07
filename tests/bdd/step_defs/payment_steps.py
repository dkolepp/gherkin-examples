"""Payment steps (owned by the payments team). pytest-bdd version."""

from pytest_bdd import given, then, when

from .helpers import p


@given(p('the shopper has a card "{card}" with a limit of {limit:Money}'))
def card(payments, card, limit):
    """Give the shopper a card that can be charged up to the limit in total."""
    payments.issue_card(card, limit)


@when(p('I charge {amount:Money} to card "{card}"'))
def charge(payments, amount, card):
    """Charge a card directly. Fails the scenario if the payment is declined."""
    payments.charge(card, amount)


@when(p('I try to charge {amount:Money} to card "{card}"'))
def try_charge(payments, outcome, amount, card):
    """Try to charge a card, which may be declined.

    Follow with 'Then I should be told "..."' to check the reason.
    """
    outcome.attempt(payments.charge, card, amount)


@then(p('card "{card}" should be charged {amount:Money}'))
def charged(payments, card, amount):
    """Check the total of all successful payments on a card."""
    assert payments.amount_charged(card) == amount


@then(p('card "{card}" should not be charged'))
def not_charged(payments, card):
    """Check that no payment went through on a card."""
    assert payments.amount_charged(card) == 0
