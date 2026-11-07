# example.com/shop/payments 1.0.0

Owned by the payments team. Covers card payments with a spending limit per card.

This module is two things:

- `example.com/shop/payments`: the application code. It knows nothing about godog.
- `example.com/shop/payments/bdd`: the team's BDD kit, a driver and step library other teams
  can use in their own feature files.

The team's own feature files are the shared, language-neutral specs in
`specs/payments/` at the repo root. `features/` here holds only the godog wiring:
a test that registers the kit, exactly as consumers do, and points godog at
the specs.

## Contract

| | godog |
|---|---|
| Require | `example.com/shop/payments v1.0.0` in your `go.mod` |
| Import | `paymentsbdd "example.com/shop/payments/bdd"` (every kit's package is named `bdd`) |
| Register | `paymentsbdd.InitializeScenario`, after `bddcommon.InitializeScenario` |
| Provides | A Before hook, `paymentsbdd.Install(ctx)`, puts a fresh `*PaymentsDriver` in the scenario's `context.Context`. Read it with `paymentsbdd.PaymentsFrom(ctx)` (panics with a clear message if missing) or `paymentsbdd.LookupPayments(ctx)` |
| Needs | the `*bddcommon.Outcome` |
| Python equivalent | `install(context)` sets `context.payments` |

## Steps

`*` means the step works as Given, When or Then (`sc.Step`); the others are
registered with `sc.Given`, `sc.When` or `sc.Then` and match only that keyword
(an `And` takes the keyword before it).

| | Sentence | What it does |
|---|---|---|
| Given | `the shopper has a card "{card}" with a limit of {limit:Money}` | Give the shopper a card that can be charged up to the limit in total. |
| When | `I charge {amount:Money} to card "{card}"` | Charge a card directly. Fails the scenario if the payment is declined. |
| When | `I try to charge {amount:Money} to card "{card}"` | Try to charge a card, which may be declined. |
| Then | `card "{card}" should be charged {amount:Money}` | Check the total of all successful payments on a card. |
| Then | `card "{card}" should not be charged` | Check that no payment went through on a card. |

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
module path in Go (`example.com/shop/payments/v2`): see "Versions" in
[`go/README.md`](../../README.md).

See [CHANGELOG.md](CHANGELOG.md).
