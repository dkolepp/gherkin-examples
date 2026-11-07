# shop-cart 1.0.0

Owned by the cart team. Covers the cart: products, discount codes and receipts.

This package is two things:

- `shop_cart`: the application code.
- `shop_cart.bdd`: the team's BDD kit, a driver and step library other teams
  can use in their own feature files.

The team's own feature files are the shared, language-neutral specs in
`specs/cart/` at the repo root. `features/` here holds only the behave wiring,
and `tests/bdd/` the pytest-bdd wiring; `pyproject.toml` points both at the specs.

## Contract

| | behave | pytest-bdd |
|---|---|---|
| Install | `pip install "shop-cart[behave]"` | `pip install "shop-cart[pytest]"` |
| Register steps | `import shop_cart.bdd.steps_behave` | `from shop_cart.bdd.steps_pytest import *` in `conftest.py` |
| Provides | `install(context)` sets `context.shop` (`ShopDriver`) | fixture `shop` (`ShopDriver`) |
| Needs | `context.outcome` | fixtures `outcome` |

Bundled test data: `standard_shop.json`. A file with that name in your working folder takes precedence.

## Steps

`*` means the step works as Given, When or Then.

| | Sentence | What it does |
|---|---|---|
| Given | `the shop sells these products:` | Stock the shop with the products in the table. |
| Given | `these discount codes exist:` | Create the discount codes in the table. |
| Given | `the shop is set up from "{file_name}"` | Stock the shop from a JSON file of products and discount codes. |
| Given | `I have an empty cart` | Start a new, empty cart. |
| Given | `I have a cart containing {qty:Qty} "{product}"` | Start a new cart that already holds some of one product. |
| * | `I add {qty:Qty} "{product}" to the cart` | Put products in the cart. Fails the scenario if that isn't possible. |
| When | `I try to add {qty:Qty} "{product}" to the cart` | Try to put products in the cart, which may be refused. |
| When | `I remove {qty:Qty} "{product}" from the cart` | Take products out of the cart. Removing more than is there empties that line. |
| When | `I apply the discount code "{code}"` | Apply a discount code. Fails the scenario if the code is refused. |
| Given | `I have applied the discount code "{code}"` | Apply a discount code. Fails the scenario if the code is refused. |
| When | `I try to apply the discount code "{code}"` | Try a discount code, which may be refused. |
| When | `I print the receipt` | Print the receipt for the current cart. |
| Then | `the cart should be empty` | Check that the cart holds nothing. |
| Then | `the cart total should be {amount:Money}` | Check the amount to pay, after any discount. Write money as $12.50. |
| Then | `the discount should be {amount:Money}` | Check how much the discount takes off. $0.00 means no discount applies. |
| Then | `the cart should contain:` | Check how many of each product are in the cart. |
| Then | `the cart should contain {qty:Qty} "{product}"` | Check how many of one product are in the cart. |
| Then | `the receipt should read:` | Check the printed receipt matches the text block below the step exactly. |

## Versioning

The step sentences are this package's public API.

- New step: minor version.
- Changed wording or meaning of a step: major version. Deprecate the old sentence
  for at least one minor version first (see `shop_bdd_common.deprecation`).
- Driver methods follow the same rules.

See [CHANGELOG.md](CHANGELOG.md).
