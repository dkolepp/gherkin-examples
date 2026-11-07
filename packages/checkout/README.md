# shop-checkout 1.0.0

Owned by the checkout team. Covers placing an order: cart total, shipping quote, card payment, shipment.

This package is two things:

- `shop_checkout`: the application code.
- `shop_checkout.bdd`: the team's BDD kit, a driver and step library other teams
  can use in their own feature files.

## Contract

| | behave | pytest-bdd |
|---|---|---|
| Install | `pip install "shop-checkout[behave]"` | `pip install "shop-checkout[pytest]"` |
| Register steps | `import shop_checkout.bdd.steps_behave` | `from shop_checkout.bdd.steps_pytest import *` in `conftest.py` |
| Provides | `install(context)` sets `context.checkout` (`CheckoutDriver`) | fixture `checkout` (`CheckoutDriver`) |
| Needs | `context.shop`, `context.payments`, `context.shipping`, `context.outcome`: install those libraries first | fixtures `shop`, `payments`, `shipping`, `outcome` |

## Steps

`*` means the step works as Given, When or Then.

| | Sentence | What it does |
|---|---|---|
| Given | `the shopper lives in "{country}"` | Set the country the order will be shipped to. |
| * | `I check out paying with card "{card}"` | Place the order for the current cart. Fails the scenario if refused. |
| When | `I try to check out paying with card "{card}"` | Try to place the order, which may be refused by payments or shipping. |
| Then | `the order total should be {amount:Money}` | Check what the order cost in total: goods after discount, plus shipping. |
| Then | `the order number should be "{order_id}"` | Check the number given to the last order placed. |

## Versioning

The step sentences are this package's public API.

- New step: minor version.
- Changed wording or meaning of a step: major version. Deprecate the old sentence
  for at least one minor version first (see `shop_bdd_common.deprecation`).
- Driver methods follow the same rules.

See [CHANGELOG.md](CHANGELOG.md).
