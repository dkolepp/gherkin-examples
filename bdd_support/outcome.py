"""Shared by all teams: remembers what went wrong in a "try to" step.

    "When I add ..."        -> must succeed; an error fails the scenario.
    "When I try to add ..." -> calls outcome.attempt(); a later
                               'Then I should be told "..."' checks the error.

Keeping this in one shared place means every team's "try to" steps work with
the same "I should be told" step.
"""


class Outcome:
    def __init__(self):
        self.last_error: Exception | None = None

    def attempt(self, action, *args) -> None:
        self.last_error = None
        try:
            action(*args)
        except Exception as error:  # domain errors are what we want to see
            self.last_error = error
