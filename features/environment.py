"""behave hooks: code that runs around features and scenarios.

Think of these like pytest fixtures. Here each scenario gets a fresh
ShopDriver, so scenarios never leak state into each other.
"""

from bdd_support.driver import ShopDriver


def before_scenario(context, scenario):
    context.shop = ShopDriver()
