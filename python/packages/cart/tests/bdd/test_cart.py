"""Turn every scenario in specs/cart/cart.feature into a pytest test.

The path is relative to bdd_features_base_dir in pyproject.toml, which points
at the shared specs folder at the repo root.

That one line is the whole test file: the scenarios are in the .feature file
and the steps are in conftest.py.
"""

from pytest_bdd import scenarios

scenarios("cart.feature")
