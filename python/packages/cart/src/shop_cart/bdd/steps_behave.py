"""Cart steps (owned by the cart team). behave version.

Each function matches one Gherkin sentence, pulls the values out of it, and
hands them to the ShopDriver (context.shop). No business logic lives here.

Compare this file with steps_pytest.py: the same sentences, glued to the
same driver, using pytest-bdd instead.

Contract
    install(context)  sets context.shop (a ShopDriver)
    requires          context.outcome, from shop_bdd_common
"""

from behave import given, step, then, when

import shop_bdd_common.steps_behave  # noqa: F401  registers Money/Qty types
from shop_cart.bdd.driver import ShopDriver


def install(context) -> None:
    context.shop = ShopDriver()


# --- Given: put the world in a known state ---------------------------------

@given("the shop sells these products:")
def step_shop_sells(context):
    """Stock the shop with the products in the table.

    Columns: name, price, stock. A stock of 0 means sold out.
    """
    # context.table is the data table under the step; rows act like dicts.
    for row in context.table:
        context.shop.stock_product(row["name"], row["price"], row["stock"])


@given("these discount codes exist:")
def step_discount_codes(context):
    """Create the discount codes in the table.

    Columns: code, percent off, minimum spend (use 0.00 for no minimum).
    """
    for row in context.table:
        context.shop.add_discount_code(
            row["code"], row["percent off"], row["minimum spend"]
        )


@given('the shop is set up from "{file_name}"')
def step_shop_from_file(context, file_name):
    """Stock the shop from a JSON file of products and discount codes.

    Bundled: "standard_shop.json". Or give the path of your own file.
    """
    context.shop.load_shop(file_name)


@given("I have an empty cart")
def step_empty_cart(context):
    """Start a new, empty cart."""
    context.shop.new_cart()


@given('I have a cart containing {qty:Qty} "{product}"')
def step_cart_containing(context, qty, product):
    """Start a new cart that already holds some of one product."""
    context.shop.new_cart()
    context.shop.add_to_cart(product, qty)


# --- When: the action under test --------------------------------------------
# @step matches Given, When, Then, And and But, so one sentence can be used
# for setup ("And I add ...") as well as for the action itself.

@step('I add {qty:Qty} "{product}" to the cart')
def step_add(context, qty, product):
    """Put products in the cart. Fails the scenario if that isn't possible.

    Works as Given, When or And.
    """
    context.shop.add_to_cart(product, qty)


@when('I try to add {qty:Qty} "{product}" to the cart')
def step_try_add(context, qty, product):
    """Try to put products in the cart, which may be refused.

    Follow with 'Then I should be told "..."' to check the reason.
    """
    context.outcome.attempt(context.shop.add_to_cart, product, qty)


@when('I remove {qty:Qty} "{product}" from the cart')
def step_remove(context, qty, product):
    """Take products out of the cart. Removing more than is there empties that line."""
    context.shop.remove_from_cart(product, qty)


# Two sentences, one function: stack the decorators.
@when('I apply the discount code "{code}"')
@given('I have applied the discount code "{code}"')
def step_apply_code(context, code):
    """Apply a discount code. Fails the scenario if the code is refused."""
    context.shop.apply_code(code)


@when('I try to apply the discount code "{code}"')
def step_try_apply_code(context, code):
    """Try a discount code, which may be refused.

    Follow with 'Then I should be told "..."' to check the reason.
    """
    context.outcome.attempt(context.shop.apply_code, code)


@when("I print the receipt")
def step_print_receipt(context):
    """Print the receipt for the current cart."""
    context.shop.print_receipt()


# --- Then: check the outcome ------------------------------------------------

@then("the cart should be empty")
def step_cart_empty(context):
    """Check that the cart holds nothing."""
    assert context.shop.cart.is_empty, "expected the cart to be empty"


@then("the cart total should be {amount:Money}")
def step_cart_total(context, amount):
    """Check the amount to pay, after any discount. Write money as $12.50."""
    assert context.shop.cart.total == amount, f"total was {context.shop.cart.total}"


@then("the discount should be {amount:Money}")
def step_discount(context, amount):
    """Check how much the discount takes off. $0.00 means no discount applies."""
    actual = context.shop.cart.discount_amount
    assert actual == amount, f"discount was {actual}"


@then("the cart should contain:")
def step_cart_contains_table(context):
    """Check how many of each product are in the cart.

    Columns: product, quantity.
    """
    for row in context.table:
        actual = context.shop.cart.quantity_of(row["product"])
        assert actual == int(row["quantity"]), f"{row['product']}: {actual}"


@then('the cart should contain {qty:Qty} "{product}"')
def step_cart_contains(context, qty, product):
    """Check how many of one product are in the cart."""
    actual = context.shop.cart.quantity_of(product)
    assert actual == qty, f"cart has {actual} {product}"



@then("the receipt should read:")
def step_receipt(context):
    """Check the printed receipt matches the text block below the step exactly."""
    # context.text holds the doc string (the text between the """ marks).
    assert context.shop.receipt == context.text, "\n" + context.shop.receipt
