"""Steps shared by every team. pytest-bdd version.

Changes here affect everyone, so they need review from all teams.
"""

from pytest_bdd import then

from .helpers import p


@then(p('I should be told "{message}"'))
def told(outcome, message):
    """Check the message the shopper saw when the last 'try to' step was refused."""
    assert outcome.last_error is not None, "expected an error, but the action succeeded"
    assert str(outcome.last_error) == message
