# shop-bdd-common 1.1.0

Shared building blocks for every team's BDD step library. Owned by all teams:
any change needs every team's review.

| Provides | behave | pytest-bdd |
|----------|--------|-----------|
| Error from the last "try to" step | `context.outcome` (set by `install(context)`) | `outcome` fixture |
| Step `Then I should be told "{message}"` | `steps_behave` | `steps_pytest` |
| Placeholder types `{amount:Money}`, `{qty:Qty}` | registered on import | `pytest_support.p(...)` |

Also: `data.load_json()` for bundled test data, `deprecation.deprecated()` for
retiring step sentences, and `behave_runner.SpecsRunner` for running feature
files kept in another folder than their `environment.py` and `steps/`
(set `runner = "shop_bdd_common.behave_runner:SpecsRunner"` in your behave config).

Install with the extra for your runner: `shop-bdd-common[behave]` or
`shop-bdd-common[pytest]`.
