# gherkin-examples: Go

Learn **Gherkin** and **BDD** (Behaviour-Driven Development) in Go, using
[godog](https://github.com/cucumber/godog) v0.16.0.

The example is the same tiny shop as the Python version, split the same way:

- **Four teams** (cart, payments, shipping, checkout) each publish a versioned
  Go module: their application code plus a *BDD kit* (a driver and a step library).
- **A light journey project** takes those modules as built artifacts from a
  local module proxy and runs an end-to-end checkout scenario. It writes no
  step code of its own.

The feature files are not in this folder. They are the shared,
language-neutral specs in [`../specs/`](../specs), the same files the Python
kits run, unchanged (see the [root README](../README.md)). This folder holds
only the Go side: the modules, their step code and the runner wiring.

All paths and commands below are relative to this `go/` folder. From the repo
root, use `make go-test`, `make go-journey`, or `make -C go <target>`.

Three lessons run through it:

1. **Specs are separate from test code.** The same Gherkin runs under behave,
   pytest-bdd and godog.
2. **Teams own their parts.** Each team's rules, driver and steps live in its
   own module.
3. **Boundaries need contracts.** Step sentences, driver state and versions are
   public API between teams.

## Quick start

Needs Go 1.24+ and `make`. godog and its dependencies come from
proxy.golang.org on first use.

```bash
make dev        # download the dependencies; go.work joins the modules from source
make test       # gofmt, go vet, each module's unit tests and own feature files
make build      # publish each module at its version into dist/
make journey    # build, then run the journey against the published modules
```

`make help` lists the targets. `make test MODULES="cart payments"` limits
them to some modules. The Makefiles set `GOTOOLCHAIN=local`, so Go never
downloads another toolchain.

## The big picture

```
../specs/                     the shared Gherkin specs, one folder per team
├── cart/ payments/ shipping/ checkout/
└── journeys/                 the end-to-end journey
modules/                      one module per team; each would be its own repo
├── common/     example.com/shop/bddcommon 1.1.0   shared steps and helpers (all teams)
├── cart/       example.com/shop/cart      1.0.0   cart, discounts, receipt
├── payments/   example.com/shop/payments  1.0.0   card payments
├── shipping/   example.com/shop/shipping  1.1.0   shipping costs and shipments
└── checkout/   example.com/shop/checkout  1.0.0   placing an order (uses the other three)
tools/publish/                publishes a module into dist/ (uses golang.org/x/mod/zip)
dist/                         built modules: a file-based module proxy (gitignored)
journey/                      runs the journey (light: no step code)
go.work                       joins the modules from source, for development only
Makefile                      team tasks: dev, test, build, journey
```

The Ansible mapping from the Python README holds here too:

| Ansible | Here |
|---|---|
| Playbook, written by an integrator | `../specs/journeys/checkout.feature`, by the journey owner |
| Role or collection, by domain experts | A team module: driver + step library + embedded data |
| Module (does the actual work) | Driver method (`Charge`, `CheckOut`) |
| Role's `argument_specs` | Step sentences and their patterns (`bddcommon.MoneyPattern`) |
| Role defaults | Bundled test data (`standard_shop.json`, via `go:embed`) |
| Molecule tests for a role | Each team's own feature files in `../specs/<team>/` |
| Galaxy | `dist/` (in real life, a private module proxy such as Athens or Artifactory) |
| `requirements.yml` | `journey/go.mod` with exact versions |
| `ansible-doc` | `make -C journey catalog` |

## Inside each team module

```
modules/payments/
├── go.mod                    module example.com/shop/payments; requires bddcommon and godog
├── VERSION                   v1.0.0: what `make build` publishes
├── README.md                 the contract: what it provides and needs, every step
├── CHANGELOG.md
├── payments.go               APP   package payments: the application code
├── bdd/                            package bdd: the published kit
│   ├── doc.go                      the contract, in Go doc form
│   ├── driver.go             HOW   business actions in Go
│   └── steps.go              GLUE  InitializeScenario: Before hook + godog steps
└── features/                 godog wiring for ../specs/payments/ (no .feature files)
    └── features_test.go      registers the kit, just like consumers do

../specs/payments/            WHAT  the team's own rules, in Gherkin
└── payments.feature
```

### Naming

- **One module per team**, path `example.com/shop/<team>`, like one Python
  package per team. The shared one is `example.com/shop/bddcommon` (folder
  `modules/common/`), because "common" says nothing at an import site.
- **The app is the module's root package** (`example.com/shop/cart`, package
  `cart`), so other apps import it by its plain name.
- **Every kit is a package named `bdd`** (`example.com/shop/cart/bdd`), the
  Go form of `shop_cart.bdd`. Import it with an alias named after the team:
  `cartbdd "example.com/shop/cart/bdd"`. Go style usually avoids packages that
  need an alias, but here the identical names are the point: every kit has the
  same shape (`InitializeScenario`, `Install`, `<Driver>From`), and the alias
  says whose it is.
- **`features/` holds only a `_test.go` file**, the Go form of the Python
  `features/` wiring folder. A folder of test files only is a valid Go package
  for `go test`, and nothing can import it.

### How the runner finds the shared specs

Each `features/features_test.go` calls `bddcommon.RunSpecs` with the path of its
team's specs and the kits to register, in order:

```go
func TestPaymentsSpecs(t *testing.T) {
	bddcommon.RunSpecs(t, []string{"../../../../specs/payments"},
		bddcommon.InitializeScenario,
		paymentsbdd.InitializeScenario,
	)
}
```

godog needs no help with specs in another folder (unlike behave, which needed
`SpecsRunner`): `Options.Paths` takes any path. `RunSpecs` adds what every team
would otherwise repeat: strict mode on, godog's command-line flags bound, each
scenario a Go subtest, and no failure when `-godog.definitions` lists the
steps without running anything.

The layers, from most readable to most technical:

```
*.feature        WHAT   Specs in plain language. Product owners can read them.
bdd/steps.go     GLUE   Match a sentence, pull out values, call the driver.
bdd/driver.go    HOW    Business actions ("charge a card").
*.go (root)      APP    The real application. Knows nothing about tests.
```

Each team tests its kit by registering it, exactly as other teams will. If a
team's own specs pass, the published steps work.

## Contracts between teams

### 1. Step sentences are public API

Every module README has a table of its step sentences, the same sentences as
the Python kits. Treat them like function signatures:

- **New step**: minor version.
- **Changed wording or meaning**: major version. Deprecate first.

Shipping 1.1.0 shows the deprecation path. It added
`a shipment of {qty} items should be sent to "{country}"` and kept the old
1.0.0 sentence `{qty} items should be shipped to "{country}"` working: it
calls `bddcommon.Deprecated`, which prints a `StepDeprecationWarning` naming the
replacement to stderr, then runs the new step. The old sentence goes in
2.0.0. See `modules/shipping/CHANGELOG.md` and its unit test
`modules/shipping/bdd/deprecated_test.go`.

godog has no placeholder types like `{amount:Money}`; steps are regular
expressions. `bddcommon.MoneyPattern` (`$12.50`) and `bddcommon.QtyPattern`
keep the patterns the same everywhere, and godog converts a matched number to
`int` when the step function takes one.

### 2. Shared state is declared

Steps from different teams share one scenario. In godog that is the
scenario's `context.Context`: every step function takes `ctx` as its first
argument. Each kit's `InitializeScenario` registers a **Before hook** that
puts its driver in the context (`Install(ctx)`), and an **accessor** reads it
back:

| Kit | Provides (Before hook → accessor) | Needs |
|---|---|---|
| bddcommon | `*Outcome` → `bddcommon.OutcomeFrom(ctx)` | – |
| cartbdd | `*ShopDriver` → `cartbdd.ShopFrom(ctx)` | outcome |
| paymentsbdd | `*PaymentsDriver` → `paymentsbdd.PaymentsFrom(ctx)` | outcome |
| shippingbdd | `*ShippingDriver` → `shippingbdd.ShippingFrom(ctx)` | outcome |
| checkoutbdd | `*CheckoutDriver` → `checkoutbdd.CheckoutFrom(ctx)` | shop, payments, shipping, outcome |

This is Python's `install(context)` contract. The context keys are
unexported types, so no kit can overwrite another's driver by accident; the
accessors are the only way in. Each accessor panics with a message naming the
kit to register (godog reports the panic as a failed step), and a
`Lookup...(ctx)` variant returns `(driver, ok)` instead.

**Hook order matters.** godog runs Before hooks in the order they were
registered, so register `bddcommon` first and `checkoutbdd` last. Get it
wrong and checkout's hook fails the scenario with *"the checkout kit needs
paymentsbdd: register paymentsbdd.InitializeScenario before
checkoutbdd.InitializeScenario"* (see `modules/checkout/bdd/install_test.go`).

### 3. Gherkin has no namespaces

The sentence text *is* the identifier, shared by every team. Two teams
defining `I have an empty cart` would clash. What godog v0.16.0 does about it
(checked in its source, `suite.go`, and by running it):

| | godog, default | godog, `Strict` | behave |
|---|---|---|---|
| Undefined step | Reported, but the suite **passes** | Suite fails | Fails |
| Pending step (`godog.ErrPending`) | Reported, but the suite **passes** | Suite fails | (no pending status) |
| A step matched by two definitions | The one registered **first wins, silently** | "ambiguous step definition", suite fails | `AmbiguousStep` at load time |
| Same sentence registered twice, never used | Nothing | Nothing | `AmbiguousStep` at load time |
| Dry run (check without running) | **None** | None | `--dry-run` |

So `bddcommon.RunSpecs` turns on `Strict` for every run. godog only checks
ambiguity for steps a scenario actually uses, and only among definitions for
the same keyword: `sc.Given` steps never clash with `sc.When` steps, and
`sc.Step` matches every keyword. `make -C journey check` adds the one static
check godog makes possible: it lists the step definitions
(`-godog.definitions`) and fails on any sentence registered twice, used or
not. Then it does a strict run, which does execute the scenarios. See
[`journey/README.md`](journey/README.md) for the full comparison with behave's
dry run.

### 4. The runner is part of the contract

Python kits have extras (`shop-cart[behave]`). Go modules have no extras: each
team module requires godog because its `bdd` package uses it. That costs an
app-only user little: Go compiles only the packages you import, so importing
`example.com/shop/cart` compiles nothing from godog. The rule that apps never
import test code is easy to check:

```bash
cd modules/checkout && go list -deps . | grep -E 'godog|bdd'   # prints nothing
```

To keep godog out of the app's module graph as well, split the kit into its
own module (`example.com/shop/cart/bdd` with its own `go.mod`). This example
keeps one module per team, like one package per team in Python.

### 5. Take team modules only from your own proxy

`make build` publishes each module into `dist/` in the layout the go command
reads from any module proxy (`go help goproxy`):

```
dist/example.com/shop/cart/@v/list           v1.0.0
dist/example.com/shop/cart/@v/v1.0.0.info    {"Version":"v1.0.0","Time":"..."}
dist/example.com/shop/cart/@v/v1.0.0.mod     the module's go.mod
dist/example.com/shop/cart/@v/v1.0.0.zip     made by golang.org/x/mod/zip, the code the go command uses
```

The journey reads it as `GOPROXY=file://.../go/dist`, with `GOWORK=off` and no
`replace` directives, so it uses exactly what was published. Its two-step
install (in `journey/Makefile`) is the Go form of the Python `--no-index`
install:

| Step | Environment | What it guards |
|---|---|---|
| 1. Team modules | `GOPROXY=file://.../dist` and nothing else, `GONOSUMDB=example.com/shop`, `GONOPROXY=none` | A module missing from `dist/` is an error, never a fall-back to a public proxy. Their names never reach sum.golang.org |
| 2. godog and friends | `GOPROXY=https://proxy.golang.org`, `GOPRIVATE=example.com/shop`, `GOVCS=example.com/shop:off` | A team module that step 1 did not fetch can't come from the public proxy (`GOPRIVATE`) or from version control (`GOVCS=off`): the install fails |
| Running | `GOPROXY=off` | Nothing is fetched at all |

Why not one `GOPROXY=file://.../dist,https://proxy.golang.org`? Because the go
command falls through to the next proxy when one answers "not found". A team
module missing from `dist/` would then be looked up on the public proxy, and
whoever controls a public module at that path would win. That is dependency
confusion.

Go is less exposed than pip to begin with: a module path starts with a domain
name, and the public proxy only serves what that domain's owner publishes. So
use a domain you own for your modules. `example.com` is safe for this example
for another reason: it is reserved for documentation (RFC 2606), so nobody can
publish a module under it (try `GOPROXY=https://proxy.golang.org go mod download
example.com/shop/cart@v1.0.0`: 404). For your real domain, set
`GOPRIVATE=yourcompany.com` everywhere, so names never leak to the public
proxy and checksum database, and serve the modules from your own proxy.

To see that the journey really resolves from `dist/`, point it at an empty one:

```bash
mkdir -p /tmp/empty && make -C journey test DIST=/tmp/empty
# go: example.com/shop/checkout@v1.0.0: reading file:///tmp/empty/example.com/shop/checkout/@v/v1.0.0.info: no such file or directory
```

## Versions

- Each module's version is in its `VERSION` file. `tools/publish` refuses to
  publish a version that has no `## x.y.z` heading in the CHANGELOG.
- Modules require each other at exact versions in `go.mod` (checkout requires
  cart v1.0.0, payments v1.0.0, shipping v1.1.0). Go's minimal version
  selection then uses the highest version anybody requires: there are no
  ranges like Python's `>=1.1,<2`.
- **A major version is a new module path.** Go's semantic import versioning
  makes shipping 2.0.0 `example.com/shop/shipping/v2`. Consumers move to it on
  purpose, by changing an import path, so 1.x and 2.x can even live side by
  side.
- `go.work` joins the modules from source for development, and needs a
  `replace` line for each version a module requires (the go command reads that
  version's `go.mod` even when the module is in the workspace). Bump those
  lines with the versions.

## Why the journey earns its place

Shipping is free from $50.00. Is that before or after the cart's discount? The
cart and shipping teams' tests pass either way, and so do the checkout team's
own. Only the journey scenario *"A discount that drops the order below $50.00
brings shipping back"* catches it.

Try it: in `modules/checkout/checkout.go`, change `basket.Total()` to
`basket.Subtotal()`. Then:

```bash
make test MODULES=checkout   # checkout's own tests: still pass
make journey                 # republishes the modules: exactly that scenario fails
```

```
--- FAIL: TestJourney/A_discount_that_drops_the_order_below_$50.00_brings_shipping_back
    order total was $50.00
5 scenarios (4 passed, 1 failed)
```

`make build` republishes the same version numbers from the changed source,
which a real proxy would refuse: published versions are immutable. That is why
the journey's install empties its own module and build caches each time (the go
command trusts a cached module version never to change) and does not commit
`go.sum`.

So keep journeys few: one happy path, rules no single team owns, and failures
that must not half-happen (a declined card must not ship anything).

## Finding steps

| Goal | Command |
|------|---------|
| Every step the journey can use | `make -C journey catalog` |
| One team's steps | its `README.md`, or `go test -godog.definitions` in its `features/` folder |
| What a step does | the comment above its registration in `bdd/steps.go` |

godog's list shows each step's regular expression and the function behind it
(`steps.go:62 -> example.com/shop/shipping/bdd.stepShipmentSent`). It has no
descriptions and doesn't say Given, When or Then, which is why each module
README has a full table. Editor plugins (VS Code's Cucumber extension, GoLand)
autocomplete sentences and jump to their code.

## Working on one team's module

```bash
cd modules/cart/features
go test                      # the team's feature files, streamed as godog runs them
cd .. && go test ./...       # plus the unit tests (go test ./... shows output only on failure)
```

Pass godog's flags after the package, with the `godog.` prefix:

| Goal | godog |
|------|-------|
| Only `@smoke` scenarios | `go test -godog.tags=@smoke` |
| Skip a tag | `go test -godog.tags='~@receipt'` |
| Both tags | `go test -godog.tags='@smoke && @cart'` |
| Either tag | `go test -godog.tags='@receipt,@discounts'` |
| Smoke but not discounts | `go test -godog.tags='@smoke && ~@discounts'` |
| One feature file | `go test -godog.paths=../../../../specs/cart/cart.feature` |
| By scenario name | `go test -run 'TestCartSpecs/stock'` (see below) |
| Full output | `go test -godog.format=pretty` (the default is `progress`) |

**godog's tag syntax is the old Cucumber one**: `~@tag` for "not", `&&` for
"and", `,` for "or". It does not understand the current Cucumber and behave
syntax (`not @receipt`, `@a and @b`, `@a or @b`): `-godog.tags='not @receipt'`
looks for a tag called "not receipt" and runs **no scenarios**, without an
error. Python's `--tags="not receipt"` is `-godog.tags='~@receipt'` here.

Each scenario is a Go subtest, so `-run` can pick scenarios by name (spaces
become `_`), but godog's summary then counts the scenarios it skipped as
"undefined". Prefer `-godog.tags` or `-godog.paths`.

## Gherkin in godog

The cheat sheet and the advice on writing good scenarios are in
[`../python/README.md`](../python/README.md); they hold for godog too. What is
different in the step code:

| Gherkin | Step function argument |
|---|---|
| `{qty}` matched by `(\d+)` | `int` (godog converts) |
| `"{name}"` matched by `"([^"]*)"` | `string` |
| `$12.50` matched by `bddcommon.MoneyPattern` | `string`; `bddcommon.MustMoney(s)` makes it `Money` |
| Data table | `*godog.Table`; `bddcommon.Rows(t)` gives one map per row |
| Doc string | `*godog.DocString`; its `Content` is the text between the `"""` marks, indentation removed |
| Shared state | `ctx context.Context`, first argument |

## Exercises

1. In `../specs/cart/cart.feature`, add a scenario: removing a "Mug" from an
   empty cart leaves the cart empty. It needs only existing steps, so both the
   Go and Python runs pass.
2. Add a step to the cart kit, e.g. `Then the cart should have {n} different
   products`. Bump `modules/cart/VERSION` to v1.1.0, add a CHANGELOG entry,
   update the `replace` line in `go.work` and the requires in checkout's and
   the journey's `go.mod`, `make build`, and use the new step in
   `../specs/journeys/checkout.feature`. The Python kit needs the same step
   before its run passes.
3. Retire a sentence: rename a payments step and deprecate the old wording the
   way `modules/shipping` does. Check the warning appears when the old
   sentence runs.
4. Add a "buy 3, pay for 2" offer to the cart. Write the scenario **first**,
   watch it fail, then implement it. That's the BDD loop.
5. Break the contract on purpose: register `^I have an empty cart$` a second
   time in the journey (add `func(sc *godog.ScenarioContext) { sc.Given(...) }`
   to the list in `journey_test.go`). `make -C journey test` still passes,
   because no journey scenario uses that sentence; `make -C journey check`
   fails. Then try a sentence the journey does use.

## Task runner

Plain `make`, as for Python. Every target is ordinary `go` commands; the
Makefiles only set the environment and loop over the modules.
