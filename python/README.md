# gherkin-examples: Python

Learn **Gherkin** and **BDD** (Behaviour-Driven Development) in Python, using
[behave](https://behave.readthedocs.io/) and
[pytest-bdd](https://pytest-bdd.readthedocs.io/).

The example is a tiny shop, split the way a real organisation might split it:

- **Four teams** (cart, payments, shipping, checkout) each publish a versioned
  package: their application code plus a *BDD kit* (a driver and a step library).
- **A light journey project** installs those packages as built wheels and runs
  an end-to-end checkout scenario. It writes no step code of its own.

The feature files are not in this folder. They are the shared,
language-neutral specs in [`../specs/`](../specs), which every language's kits
implement (see the [root README](../README.md)). This folder holds only the
Python side: the packages, their step code and the runner wiring.

All paths and commands below are relative to this `python/` folder. From the
repo root, use `make python-test`, `make python-journey`, or
`make -C python <target>`.

Three lessons run through it:

1. **Specs are separate from test code.** The same Gherkin runs under behave
   and pytest-bdd (and, with go/, under godog).
2. **Teams own their parts.** Each team's rules, driver and steps live in its
   own package.
3. **Boundaries need contracts.** Step sentences, driver state and versions are
   public API between teams.

## Quick start

Needs Python 3.10+ and `make`.

```bash
make dev        # create .venv and install every team package (editable)
make test       # run each team's own feature files and tests
make build      # build each package into a wheel in dist/
make journey    # build, then run the journey against the built wheels
```

`make help` lists the targets. `make test PACKAGES="cart payments"` limits
them to some packages.

## The big picture

```
../specs/                     the shared Gherkin specs, one folder per team
├── cart/ payments/ shipping/ checkout/
└── journeys/                 the end-to-end journey
packages/                     one folder per team; each would be its own repo
├── common/     shop-bdd-common 1.1.0   shared steps and helpers (all teams)
├── cart/       shop-cart       1.0.0   cart, discounts, receipt
├── payments/   shop-payments   1.0.0   card payments
├── shipping/   shop-shipping   1.1.0   shipping costs and shipments
└── checkout/   shop-checkout   1.0.0   placing an order (uses the other three)
dist/                         built wheels: our stand-in for a package index
journey/                      runs the journey (light: no step code)
Makefile                      team tasks: dev, test, build, journey
```

It maps closely onto the Ansible ecosystem:

| Ansible | Here |
|---|---|
| Playbook, written by an integrator | `../specs/journeys/checkout.feature`, by the journey owner |
| Role or collection, by domain experts | A team package: driver + step library + bundled data |
| Module (does the actual work) | Driver method (`charge`, `ship`) |
| Role's `argument_specs` | Step sentences and their typed placeholders (`{amount:Money}`) |
| Role defaults | Bundled test data (`standard_shop.json`) |
| Molecule tests for a role | Each team's own feature files in `../specs/<team>/` |
| Galaxy | `dist/` (in real life, a private package index) |
| `requirements.yml` | `journey/requirements.txt` with pinned versions |
| `ansible-doc` | `make -C journey catalog` |
| Controller: RBAC, job templates | CI with protected environments (not shown) |

## Inside each team package

```
packages/payments/
├── pyproject.toml            name, version, and the [behave] / [pytest] extras
├── README.md                 the contract: what it provides and needs, every step
├── CHANGELOG.md
├── src/shop_payments/
│   ├── __init__.py           APP   the application code
│   └── bdd/
│       ├── driver.py         HOW   business actions in Python
│       ├── steps_behave.py   GLUE  step definitions for behave
│       └── steps_pytest.py   GLUE  the same steps for pytest-bdd
├── features/                 behave wiring for ../specs/payments/ (no .feature files)
│   ├── environment.py        gives each scenario fresh drivers
│   └── steps/use_steps.py    imports the step library, just like consumers do
└── tests/bdd/                pytest-bdd wiring for the same feature files

../specs/payments/            WHAT  the team's own rules, in Gherkin
└── payments.feature
```

### How the runners find the shared specs

Both runners are configured in the package's `pyproject.toml`:

```toml
[tool.pytest.ini_options]
bdd_features_base_dir = "../../../specs/payments/"   # scenarios("payments.feature") looks here

[tool.behave]
runner = "shop_bdd_common.behave_runner:SpecsRunner"
paths = ["features", "../../../specs/payments"]      # wiring first, then the specs
```

pytest-bdd only needs the base folder. behave needs help: it takes the first
path as its base directory and insists that it holds both `steps/` (or
`environment.py`) and at least one `.feature` file. Plain behave therefore
fails both ways round: with the specs first it finds no `steps/`, with the
wiring first it finds no feature files. `SpecsRunner`, a small `Runner`
subclass in shop-bdd-common, uses the first path only to find the wiring
and runs the remaining paths.

The layers, from most readable to most technical:

```
*.feature        WHAT   Specs in plain language. Product owners can read them.
steps_*.py       GLUE   Match a sentence, pull out values, call the driver.
driver.py        HOW    Business actions ("charge a card"). Shared by both runners.
src/shop_*/      APP    The real application. Knows nothing about tests.
```

Each team tests its step library by importing it, exactly as other teams will.
If a team's own specs pass, the published steps work.

## Contracts between teams

### 1. Step sentences are public API

Every package README has a table of its step sentences. Treat them like
function signatures:

- **New step**: minor version.
- **Changed wording or meaning**: major version. Deprecate first.

Shipping 1.1.0 shows the deprecation path. It added
`a shipment of {qty} items should be sent to "{country}"` and kept the old
1.0.0 sentence `{qty} items should be shipped to "{country}"` working, with a
`StepDeprecationWarning` naming the replacement. The old sentence goes in
2.0.0. See `packages/shipping/CHANGELOG.md`.

### 2. Shared state is declared

Steps from different teams share one world, so each library says what it puts
there and what it needs:

| Library | Provides (behave / pytest-bdd) | Needs |
|---|---|---|
| shop-bdd-common | `context.outcome` / `outcome` | – |
| shop-cart | `context.shop` / `shop` | outcome |
| shop-payments | `context.payments` / `payments` | outcome |
| shop-shipping | `context.shipping` / `shipping` | outcome |
| shop-checkout | `context.checkout` / `checkout` | shop, payments, shipping, outcome |

For behave, each library has `install(context)`; the journey calls them in
order. For pytest-bdd, they are fixtures.

### 3. Gherkin has no namespaces

The sentence text *is* the identifier, shared by every team. Two teams
defining `I have an empty cart` would clash. behave refuses to load duplicate
sentences (`AmbiguousStep`), and `make -C journey check` runs that check
without executing anything. It also fails on sentences nobody defined.
pytest-bdd silently lets the closer definition win, so run the behave check
even if you use pytest.

### 4. The runner is part of the contract

Install a kit with the extra for your runner: `shop-cart[behave]` or
`shop-cart[pytest]`. The journey uses behave only, so it never installs
pytest-bdd.

### 5. Install team packages only from your own index

`journey/Makefile` installs the team packages with `--no-index` from `dist/`,
then their third-party dependencies from PyPI. Otherwise a public package with
the same name could win (dependency confusion).

## Why the journey earns its place

Shipping is free from $50.00. Is that before or after the cart's discount? The
cart and shipping teams' tests pass either way, and so do the checkout team's
own. Only the journey scenario *"A discount that drops the order below $50.00
brings shipping back"* catches it.

Try it: in `packages/checkout/src/shop_checkout/__init__.py`, change
`cart.total` to `cart.subtotal`. Then:

```bash
make test PACKAGES=checkout   # checkout's own tests: still pass
make journey                  # rebuilds the wheels: exactly that scenario fails
```

So keep journeys few: one happy path, rules no single team owns, and failures
that must not half-happen (a declined card must not ship anything).

## Finding steps

| Goal | Command |
|------|---------|
| Every step the journey can use, with descriptions | `make -C journey catalog` |
| One team's steps | its `README.md`, or `../../.venv/bin/behave --steps-catalog` in its folder |
| Where each step is used, and unused ones | `behave --dry-run -f steps.usage` in a package |
| pytest-bdd steps (sentences only) | `pytest --fixtures tests/bdd \| grep pytestbdd_stepdef_` in a package |
| godog steps (Go) | `godog --definitions` |

Each step's docstring is written for scenario writers, and the catalog shows
it. Day to day, editor plugins (VS Code's Cucumber extension, PyCharm
Professional) autocomplete step sentences and jump to their code.

## Working on one team's package

```bash
source .venv/bin/activate
cd packages/cart
behave                       # the team's feature files, with behave
pytest                       # the same feature files with pytest-bdd, plus unit tests
```

| Goal | behave | pytest-bdd |
|------|--------|-----------|
| Only `@smoke` scenarios | `behave --tags=smoke` | `pytest -m smoke` |
| Skip a tag | `behave --tags="not receipt"` | `pytest -m "not receipt"` |
| One feature file | `behave ../../../specs/cart/cart.feature` | `pytest tests/bdd/test_cart.py` |
| By scenario name | `behave -n "stock"` | `pytest -k stock` |
| Compact output | `behave -f progress` | `pytest -q` |

(pytest needs each tag registered under `markers` in the package's `pyproject.toml`.)

## Gherkin cheat sheet

```gherkin
@tag                                # tags: filter what runs
Feature: Short name
  Free text describing the feature (often "As a … I want … So that …").

  Background:                       # runs before every scenario in the file
    Given some shared setup

  Rule: One business rule           # optional grouping (Gherkin 6+)

    Scenario: One concrete example
      Given a starting state        # Given = arrange
      When something happens        # When  = act
      Then this is the result       # Then  = assert
      And another result            # And / But continue the previous keyword

    Scenario Outline: Same steps, many data rows
      When I enter <input>
      Then I see <output>

      Examples:
        | input | output |
        | 1     | one    |
        | 2     | two    |

  Scenario: Step arguments
    Given these rows:               # data table → context.table / datatable
      | name | price |
      | Mug  | 8.00  |
    Then the text is:               # doc string → context.text / docstring
      """
      multi-line text
      """
```

(The `#` notes above are for reading only. In a real feature file, comments
must be on their own line: text after a step, `# like this`, becomes part of
the step.)

## Writing good scenarios

- **Describe behaviour, not clicks.** `When I apply the discount code "SAVE10"`,
  not `When I type SAVE10 into field #code and press Enter`.
- **One behaviour per scenario.** If you need "And then … And then …", split it.
- **Use real examples.** `$12.50 → $11.25` is clearer than "the price is reduced".
- **Keep Given/When/Then in order.** Arrange, act, assert.
- **Make failure explicit.** `I add` must succeed, while `I try to add`
  records the error so a `Then I should be told "…"` step can check it.
- **Keep the technical edge cases in unit tests** (see `packages/cart/tests/unit/`).

## Exercises

1. In `../specs/cart/cart.feature`, add a scenario: removing a "Mug"
   from an empty cart leaves the cart empty. It needs only existing steps.
2. Add a step to the cart kit, e.g. `Then the cart should have {n:Qty} different
   products`, in both `steps_behave.py` and `steps_pytest.py`. Bump the cart
   version to 1.1.0, add a CHANGELOG entry, `make build`, update the pin in
   `journey/requirements.txt`, and use the new step in the journey
   (`../specs/journeys/checkout.feature`). A new sentence in a shared spec
   needs a step in every language's kit before their runs pass.
3. Retire a sentence: rename a payments step and deprecate the old wording the
   way `packages/shipping` does. Check the warning appears when the old
   sentence runs.
4. Add a "buy 3, pay for 2" offer to the cart. Write the scenario **first**,
   watch it fail, then implement it. That's the BDD loop.
5. Break the contract on purpose: define `I have an empty cart` again in
   `journey/features/steps/team_steps.py` and run `make -C journey check`.

## Further reading

| Document | What it covers |
|---|---|
| [Gherkin support matrix](../docs/gherkin-support-matrix.html) | Which Gherkin syntax behave, godog and pytest-bdd support, tested against each tool, plus Gherkin's version history. Download and open in a browser; [live version](https://claude.ai/artifact/MT6CLm637feEN39eXGMrjU) |
| [Journey Testing Platform](../docs/journey-testing-platform.md) | Journey testing as a platform problem: personas, responsibilities, contracts at each boundary, friction and how to measure it |
| [Jira Xray: Test Management Model](../docs/jira-xray-test-management.md) | What Xray is for, its entity model, Gherkin support, limits, and how it compares with the platform model |

The Markdown files are exported copies; each links to its live version.

## Task runner

Plain `make` is used because it's already installed almost everywhere. If you
prefer [mise](https://mise.jdx.dev/) or [invoke](https://www.pyinvoke.org/),
their tasks can simply call these same `make` targets; mise can also pin the
Python version.
