# shop-payments 1.0.0

Owned by the payments team. Covers card payments with a spending limit per card.

This package is two things:

- `shop_payments`: the application code.
- `shop_payments.bdd`: the team's BDD kit, a driver and step library other teams
  can use in their own feature files.

## Contract

| | behave | pytest-bdd |
|---|---|---|
| Install | `pip install "shop-payments[behave]"` | `pip install "shop-payments[pytest]"` |
| Register steps | `import shop_payments.bdd.steps_behave` | `from shop_payments.bdd.steps_pytest import *` in `conftest.py` |
| Provides | `install(context)` sets `context.payments` (`PaymentsDriver`) | fixture `payments` (`PaymentsDriver`) |
| Needs | `context.outcome` | fixtures `outcome` |

## Steps

`*` means the step works as Given, When or Then.

| | Sentence | What it does |
|---|---|---|
| Given | `the shopper has a card "{card}" with a limit of {limit:Money}` | Give the shopper a card that can be charged up to the limit in total. |
| When | `I charge {amount:Money} to card "{card}"` | Charge a card directly. Fails the scenario if the payment is declined. |
| When | `I try to charge {amount:Money} to card "{card}"` | Try to charge a card, which may be declined. |
| Then | `card "{card}" should be charged {amount:Money}` | Check the total of all successful payments on a card. |
| Then | `card "{card}" should not be charged` | Check that no payment went through on a card. |

## Versioning

The step sentences are this package's public API.

- New step: minor version.
- Changed wording or meaning of a step: major version. Deprecate the old sentence
  for at least one minor version first (see `shop_bdd_common.deprecation`).
- Driver methods follow the same rules.

See [CHANGELOG.md](CHANGELOG.md).
