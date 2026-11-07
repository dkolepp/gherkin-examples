# shop-shipping 1.1.0

Owned by the shipping team. Covers shipping rates per country, free shipping and shipments.

This package is two things:

- `shop_shipping`: the application code.
- `shop_shipping.bdd`: the team's BDD kit, a driver and step library other teams
  can use in their own feature files.

## Contract

| | behave | pytest-bdd |
|---|---|---|
| Install | `pip install "shop-shipping[behave]"` | `pip install "shop-shipping[pytest]"` |
| Register steps | `import shop_shipping.bdd.steps_behave` | `from shop_shipping.bdd.steps_pytest import *` in `conftest.py` |
| Provides | `install(context)` sets `context.shipping` (`ShippingDriver`) | fixture `shipping` (`ShippingDriver`) |
| Needs | `context.outcome` | fixtures `outcome` |

Bundled test data: `standard_shipping.json`. A file with that name in your working folder takes precedence.

## Steps

`*` means the step works as Given, When or Then.

| | Sentence | What it does |
|---|---|---|
| Given | `shipping costs:` | Set a flat shipping rate per country. |
| Given | `orders of {amount:Money} or more ship free` | Make shipping free for orders at or above this value. |
| Given | `shipping is set up from "{file_name}"` | Load shipping rates and the free-shipping threshold from a JSON file. |
| When | `I get a shipping quote to "{country}" for {amount:Money}` | Ask what shipping costs for an order of this value. |
| When | `I try to get a shipping quote to "{country}" for {amount:Money}` | Ask for a shipping quote, which may be refused. |
| Then | `shipping should cost {amount:Money}` | Check the last shipping quote. |
| Then | `a shipment of {qty:Qty} items should be sent to "{country}"` | Check that exactly one shipment went out, with this many items. |
| Then | `no shipment should be sent` | Check that nothing was shipped. |
| Then | `{qty:Qty} items should be shipped to "{country}"` | Deprecated since 1.1.0: use 'a shipment of N items should be sent to "..."'. |

## Versioning

The step sentences are this package's public API.

- New step: minor version.
- Changed wording or meaning of a step: major version. Deprecate the old sentence
  for at least one minor version first (see `shop_bdd_common.deprecation`).
- Driver methods follow the same rules.

See [CHANGELOG.md](CHANGELOG.md).
