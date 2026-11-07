# This journey writes no step code. Importing the teams' published step
# libraries registers all their steps with behave.
import shop_bdd_common.steps_behave  # noqa: F401
import shop_cart.bdd.steps_behave  # noqa: F401
import shop_checkout.bdd.steps_behave  # noqa: F401
import shop_payments.bdd.steps_behave  # noqa: F401
import shop_shipping.bdd.steps_behave  # noqa: F401
