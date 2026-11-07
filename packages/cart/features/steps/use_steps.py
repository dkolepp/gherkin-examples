# The team tests its own published step library exactly the way consumers
# use it: by importing it. Importing registers the steps with behave.
import shop_bdd_common.steps_behave  # noqa: F401
import shop_cart.bdd.steps_behave  # noqa: F401
