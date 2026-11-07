"""Steps shared by every team. behave version.

Changes here affect everyone, so they need review from all teams.
"""

from behave import then


@then('I should be told "{message}"')
def step_told(context, message):
    """Check the message the shopper saw when the last 'try to' step was refused."""
    error = context.outcome.last_error
    assert error is not None, "expected an error, but the action succeeded"
    assert str(error) == message, f"got: {error}"
