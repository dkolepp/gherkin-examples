"""Turn every scenario in features/cart.feature into a pytest test.

That one line is the whole test file: the scenarios are in the .feature file
and the steps are in conftest.py.
"""

from pytest_bdd import scenarios

scenarios("cart.feature")
