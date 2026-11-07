# example.com/shop/shipping 1.1.0

Owned by the shipping team. Covers shipping rates per country, free shipping and shipments.

This module is two things:

- `example.com/shop/shipping`: the application code. It knows nothing about godog.
- `example.com/shop/shipping/bdd`: the team's BDD kit, a driver and step library other teams
  can use in their own feature files.

The team's own feature files are the shared, language-neutral specs in
`specs/shipping/` at the repo root. `features/` here holds only the godog wiring:
a test that registers the kit, exactly as consumers do, and points godog at
the specs.

## Contract

| | godog |
|---|---|
| Require | `example.com/shop/shipping v1.1.0` in your `go.mod` |
| Import | `shippingbdd "example.com/shop/shipping/bdd"` (every kit's package is named `bdd`) |
| Register | `shippingbdd.InitializeScenario`, after `bddcommon.InitializeScenario` |
| Provides | A Before hook, `shippingbdd.Install(ctx)`, puts a fresh `*ShippingDriver` in the scenario's `context.Context`. Read it with `shippingbdd.ShippingFrom(ctx)` (panics with a clear message if missing) or `shippingbdd.LookupShipping(ctx)` |
| Needs | the `*bddcommon.Outcome` |
| Python equivalent | `install(context)` sets `context.shipping` |

Bundled test data: `standard_shipping.json` (embedded with `go:embed`). A file with that name in your working folder takes precedence.

## Steps

`*` means the step works as Given, When or Then (`sc.Step`); the others are
registered with `sc.Given`, `sc.When` or `sc.Then` and match only that keyword
(an `And` takes the keyword before it).

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
| Then | `{qty:Qty} items should be shipped to "{country}"` | Deprecated since 1.1.0: use 'a shipment of N items should be sent to "..."'. Prints a `StepDeprecationWarning` to stderr. |

In Go the placeholders are regular expressions: `{qty:Qty}` is
`bddcommon.QtyPattern` (`(\d+)`), `{amount:Money}` is `bddcommon.MoneyPattern`
(`$12.50`), and `"{name}"` is `"([^"]*)"`. `make -C ../../journey catalog`
prints them.

## Versioning

The step sentences are this module's public API.

- New step: minor version.
- Changed wording or meaning of a step: major version. Deprecate the old sentence
  for at least one minor version first (see `bddcommon.Deprecated`).
- Driver methods and exported Go names follow the same rules.

The version is in [`VERSION`](VERSION); `make build` in `go/` publishes it and
refuses a version with no CHANGELOG entry. A major version also changes the
module path in Go (`example.com/shop/shipping/v2`): see "Versions" in
[`go/README.md`](../../README.md).

See [CHANGELOG.md](CHANGELOG.md).
