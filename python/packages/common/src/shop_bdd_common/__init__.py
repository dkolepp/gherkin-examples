"""Shared BDD building blocks used by every team's step library.

    money.py          to_money(): turn "12.50" into Decimal("12.50")
    types.py          {amount:Money} and {qty:Qty} step placeholder types
    outcome.py        remembers the error from a "try to" step
    data.py           loads JSON test data (bundled defaults or your own file)
    deprecation.py    marks a step sentence as deprecated
    behave_runner.py  runs specs that live apart from their behave wiring
    steps_behave.py   the shared 'I should be told "..."' step (behave)
    steps_pytest.py   the same step (pytest-bdd)
"""

__version__ = "1.1.0"
