# example.com/shop/bddcommon 1.1.0

Shared building blocks for every team's BDD kit. Owned by all teams: any
change needs every team's review. It is test code: application packages never
import it.

## Contract

| | godog |
|---|---|
| Require | `example.com/shop/bddcommon v1.1.0` in your `go.mod` |
| Import | `"example.com/shop/bddcommon"` |
| Register | `bddcommon.InitializeScenario`, before every other kit |
| Provides | A Before hook, `bddcommon.Install(ctx)`, puts a fresh `*Outcome` (the error from the last "try to" step) in the scenario's `context.Context`. Read it with `bddcommon.OutcomeFrom(ctx)` or `bddcommon.LookupOutcome(ctx)` |
| Provides | Step `Then I should be told "{message}"` |
| Needs | nothing |
| Python equivalent | `install(context)` sets `context.outcome` |

Also:

| Name | What it is |
|---|---|
| `Money`, `ParseMoney`, `MustMoney`, `ToMoney` | Money in whole cents: `$12.50` in steps, `"12.50"` in tables and JSON, rounded half up |
| `MoneyPattern`, `QtyPattern` | Regular-expression pieces for step patterns, the Go form of `{amount:Money}` and `{qty:Qty}` |
| `Rows` | A data table as one `map[string]string` per row, keyed by the header |
| `LoadJSON` | Bundled test data from a kit's embedded `data/`; a file of that name in the working folder wins |
| `Deprecated` | Warns (`StepDeprecationWarning:` on stderr) that a step sentence is deprecated, naming its replacement |
| `RunSpecs` | Runs feature files kept in another folder (the shared `specs/`) with godog, strict, each scenario as a Go subtest. Takes the kits' `InitializeScenario` functions in order. The Go counterpart of the Python `SpecsRunner` |

`RunSpecs` reads godog's own flags (`-godog.tags`, `-godog.format`,
`-godog.definitions`, `-godog.paths`, ...). It turns on `Strict` by default:
plain godog passes a suite with undefined or pending steps.

## Versioning

New exported name: minor version. Changed behaviour or removed name: major
version, with a new module path (`example.com/shop/bddcommon/v2`). The version
is in [`VERSION`](VERSION). See [CHANGELOG.md](CHANGELOG.md).
