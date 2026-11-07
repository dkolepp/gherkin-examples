"""behave step definitions.

Each function matches one Gherkin sentence, pulls the values out of it, and
hands them to the ShopDriver (context.shop). No business logic lives here.

Compare this file with tests/bdd/steps.py: the same sentences, glued to the
same driver, using pytest-bdd instead.
"""

from behave import given, register_type, step, then, when

from bdd_support.types import TYPES

register_type(**TYPES)  # enables {amount:Money} and {qty:Qty} in patterns


# --- Given: put the world in a known state ---------------------------------

@given("the shop sells these products:")
def step_shop_sells(context):
    # context.table is the data table under the step; rows act like dicts.
    for row in context.table:
        context.shop.stock_product(row["name"], row["price"], row["stock"])


@given("these discount codes exist:")
def step_discount_codes(context):
    for row in context.table:
        context.shop.add_discount_code(
            row["code"], row["percent off"], row["minimum spend"]
        )


@given('the shop is set up from "{file_name}"')
def step_shop_from_file(context, file_name):
    context.shop.load_shop(file_name)


@given("I have an empty cart")
def step_empty_cart(context):
    context.shop.new_cart()


@given('I have a cart containing {qty:Qty} "{product}"')
def step_cart_containing(context, qty, product):
    context.shop.new_cart()
    context.shop.add_to_cart(product, qty)


# --- When: the action under test --------------------------------------------
# @step matches Given, When, Then, And and But, so one sentence can be used
# for setup ("And I add ...") as well as for the action itself.

@step('I add {qty:Qty} "{product}" to the cart')
def step_add(context, qty, product):
    context.shop.add_to_cart(product, qty)


@when('I try to add {qty:Qty} "{product}" to the cart')
def step_try_add(context, qty, product):
    context.shop.attempt(context.shop.add_to_cart, product, qty)


@when('I remove {qty:Qty} "{product}" from the cart')
def step_remove(context, qty, product):
    context.shop.remove_from_cart(product, qty)


# Two sentences, one function: stack the decorators.
@when('I apply the discount code "{code}"')
@given('I have applied the discount code "{code}"')
def step_apply_code(context, code):
    context.shop.apply_code(code)


@when('I try to apply the discount code "{code}"')
def step_try_apply_code(context, code):
    context.shop.attempt(context.shop.apply_code, code)


@when("I print the receipt")
def step_print_receipt(context):
    context.shop.print_receipt()


# --- Then: check the outcome ------------------------------------------------

@then("the cart should be empty")
def step_cart_empty(context):
    assert context.shop.cart.is_empty, "expected the cart to be empty"


@then("the cart total should be {amount:Money}")
def step_cart_total(context, amount):
    assert context.shop.cart.total == amount, f"total was {context.shop.cart.total}"


@then("the discount should be {amount:Money}")
def step_discount(context, amount):
    actual = context.shop.cart.discount_amount
    assert actual == amount, f"discount was {actual}"


@then("the cart should contain:")
def step_cart_contains_table(context):
    for row in context.table:
        actual = context.shop.cart.quantity_of(row["product"])
        assert actual == int(row["quantity"]), f"{row['product']}: {actual}"


@then('the cart should contain {qty:Qty} "{product}"')
def step_cart_contains(context, qty, product):
    actual = context.shop.cart.quantity_of(product)
    assert actual == qty, f"cart has {actual} {product}"


@then('I should be told "{message}"')
def step_told(context, message):
    error = context.shop.last_error
    assert error is not None, "expected an error, but the action succeeded"
    assert str(error) == message, f"got: {error}"


@then("the receipt should read:")
def step_receipt(context):
    # context.text holds the doc string (the text between the """ marks).
    assert context.shop.receipt == context.text, "\n" + context.shop.receipt
