"""The other way to bind scenarios: one at a time with @scenario.

Use this when you want a named test function, for example to add extra
pytest markers or to run some code after the scenario's steps.
"""

from pytest_bdd import scenario


@scenario("domains/cart/receipt.feature", "Receipt lists each line, the discount and the total")
def test_receipt_shows_lines_discount_and_total():
    pass  # the steps do the work; anything here runs after them
