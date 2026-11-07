# gherkin-examples

Learn **Gherkin** and **BDD** (Behaviour-Driven Development) in Python, using
both [behave](https://behave.readthedocs.io/) and
[pytest-bdd](https://pytest-bdd.readthedocs.io/).

The example is a tiny shop: a cart with discount codes and a printed receipt,
card payments, shipping, and a checkout journey that ties them together. The
same `.feature` files are run by **both** tools, which shows the main lesson:
*the test cases (specs) are separate from the test code.*

The second lesson is about ownership: each part of the shop belongs to a
different (imaginary) team, and a cross-team journey reuses their work instead
of duplicating it. See [A journey across teams](#a-journey-across-teams).

## Quick start

```bash
python -m venv .venv
source .venv/bin/activate
pip install -r requirements.txt

behave        # run the specs with behave
pytest        # run the same specs with pytest-bdd (plus unit tests)
```

## How the concerns are separated

```
features/*.feature        1. WHAT   Specs in plain language. No code.
        │                           Product owners and testers can read and write these.
        ▼
step definitions          2. GLUE   Match a sentence, pull out values, call the driver.
  features/steps/ (behave)          One or two lines each.
  tests/bdd/step_defs/ (pytest-bdd)
        │
        ▼
bdd_support/*_driver.py   3. HOW    Business actions in Python ("add to cart",
                                    "apply code"). Shared by both runners.
        │
        ▼
src/                      4. APP    The real application. Knows nothing about tests.

test_data/*.json          Test data that is too big or too shared to live in a spec.
```

Why bother?

- **Specs change for business reasons, code changes for technical reasons.**
  Keeping them apart means one kind of change doesn't break the other.
- **Steps are reusable.** `Given I have an empty cart` is written once and used
  in many scenarios.
- **The driver hides the "how".** If the shop got a web UI tomorrow, you'd write
  a new driver that clicks buttons. The `.feature` files and steps stay the same.
- **The runner is replaceable.** This repo proves it: behave and pytest-bdd read
  identical specs.

## Repository tour

| Path | What it is |
|------|-----------|
| `features/domains/cart/cart.feature` | Background, data tables, tags, Scenario Outline with two Examples tables |
| `features/domains/cart/discounts.feature` | `Rule:` blocks, Outline, expecting errors |
| `features/domains/cart/receipt.feature` | Doc strings (`"""`), loading data from `test_data/` |
| `features/domains/payments/`, `…/shipping/` | The payments and shipping teams' own rules |
| `features/journeys/checkout.feature` | A journey across cart, payments and shipping |
| `features/environment.py` | behave hooks (fresh drivers per scenario) |
| `features/steps/*_steps.py` | behave step definitions, one module per team |
| `tests/bdd/step_defs/*_steps.py` | pytest-bdd step definitions, one module per team |
| `tests/bdd/conftest.py` | pytest-bdd fixtures; imports all step modules |
| `tests/bdd/test_*.py` | Bind feature files to pytest (`scenarios(...)` or `@scenario`) |
| `tests/unit/` | Plain unit tests, for contrast |
| `bdd_support/*_driver.py` | One driver per team, shared by both runners |
| `bdd_support/outcome.py` | Remembers errors from "try to" steps, for all teams |
| `bdd_support/types.py` | Custom `{amount:Money}` / `{qty:Qty}` step types |
| `src/shop/`, `src/payments/`, `src/shipping/`, `src/checkout/` | The application, one package per domain |
| `docs/CODEOWNERS.example` | How the team ownership would be enforced on GitHub |

Tip: open `features/steps/cart_steps.py` and `tests/bdd/step_defs/cart_steps.py`
side by side. Same sentences, same driver calls, different glue.

## A journey across teams

The checkout journey crosses three domains, each owned by a different team:

| Team | Owns |
|------|------|
| Cart | `src/shop/`, `features/domains/cart/`, `cart_driver.py`, `cart_steps.py` |
| Payments | `src/payments/`, `features/domains/payments/`, `payments_driver.py`, `payment_steps.py` |
| Shipping | `src/shipping/`, `features/domains/shipping/`, `shipping_driver.py`, `shipping_steps.py` |
| Checkout (journey owner) | `src/checkout/`, `features/journeys/`, `checkout_driver.py`, `checkout_steps.py` |
| Everyone | `common_steps.py`, `outcome.py`, `types.py`: every team should review changes |

How it's put together:

- **Each team's detailed rules live in its own feature file.** Card limits are
  tested in `payments.feature`, not in the journey.
- **The journey reuses other teams' steps.** `checkout.feature` uses the cart
  team's `I have a cart containing…`, the payments team's
  `card "4242" should be charged…` and the shipping team's
  `a shipment of 3 items should be sent…`. The checkout team only wrote the
  four checkout steps.
- **Teams talk through drivers, not internals.** `CheckoutDriver` is built from
  the other teams' drivers, and `src/checkout/` only calls the other packages'
  public classes. A team can change its internals without breaking the journey.
- **Journeys stay few.** They check that the pieces connect. Here that means
  one happy path, a rule no single team owns, and failures that must not
  half-happen (a declined card must not ship anything).

**Why the journey earns its place:** shipping is free from $50.00, but is that
before or after the cart's discount? The cart and shipping teams' tests each
pass whichever way checkout does it. Only the journey scenario
*"A discount that drops the order below $50.00 brings shipping back"* catches
checkout using the wrong number. Try it: in `src/checkout/__init__.py` change
`cart.total` to `cart.subtotal` and run either tool. Exactly that one scenario
fails.

Run by owner with tags: `behave --tags=journey`, `pytest -m payments`.

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

## Finding steps

Steps are reusable, so before writing a new sentence, check whether one already
exists. Each step function has a plain-language docstring, which acts as the
catalog's description.

| Goal | Command |
|------|---------|
| List every step with its description | `behave --steps-catalog` |
| See where each step is used, and spot unused ones | `behave --dry-run -f steps.usage` |
| List pytest-bdd steps (sentences only) | `pytest --fixtures tests/bdd \| grep pytestbdd_stepdef_` |
| List godog steps (Go) | `godog --definitions` |

Example `behave --steps-catalog` output:

```
Given I have a cart containing {qty:Qty} "{product}"
    Start a new cart that already holds some of one product.

When I try to add {qty:Qty} "{product}" to the cart
    Try to put products in the cart, which may be refused.

    Follow with 'Then I should be told "..."' to check the reason.
```

pytest-bdd has no catalog command. It registers each step as a pytest fixture
named `pytestbdd_stepdef_<keyword>_<sentence>`, so the fixture list doubles as a
step list, without the docstrings. Because both runners share the same
sentences, the behave catalog describes the pytest-bdd steps too.

Day to day, an editor plugin is the quickest way to find steps: VS Code's
Cucumber extension and PyCharm Professional both autocomplete step sentences
and jump to their code.

Keeping the catalog useful:

- Write each docstring for someone writing scenarios, not for programmers.
  Mention table columns and which step to pair it with.
- Keep the docstrings in `features/steps/` and `tests/bdd/step_defs/` in sync.
- Group steps by owning team or business area (`cart_steps.py`,
  `payment_steps.py`), as this repo does.
- Before adding a step, search the catalog for one that already says it.

## Running a subset

| Goal | behave | pytest-bdd |
|------|--------|-----------|
| Only `@smoke` scenarios | `behave --tags=smoke` | `pytest -m smoke` |
| Skip a tag | `behave --tags="not receipt"` | `pytest -m "not receipt"` |
| One feature file | `behave features/domains/cart/cart.feature` | `pytest tests/bdd/test_cart.py` |
| By scenario name | `behave -n "stock"` | `pytest -k stock` |
| Compact output | `behave -f progress` | `pytest -q` |

(pytest needs each tag registered under `markers` in `pyproject.toml`.)

## Writing good scenarios

- **Describe behaviour, not clicks.** `When I apply the discount code "SAVE10"`,
  not `When I type SAVE10 into field #code and press Enter`.
- **One behaviour per scenario.** If you need "And then … And then …", split it.
- **Use real examples.** `$12.50 → $11.25` is clearer than "the price is reduced".
- **Keep Given/When/Then in order.** Arrange, act, assert.
- **Make failure explicit.** Here, `I add` must succeed, while `I try to add`
  records the error so a `Then I should be told "…"` step can check it.
- **Keep the technical edge cases in unit tests** (see `tests/unit/`).

## Exercises

1. Add a `Scenario` to `features/domains/cart/cart.feature`: removing a "Mug" from an empty cart
   leaves the cart empty. You can write it using only existing steps, so it
   should pass straight away in both tools.
2. Write a new step, e.g. `Then the cart should have {n:Qty} different products`.
   Implement it for behave *and* pytest-bdd (one line each, plus a driver method
   if needed).
3. Add a "buy 3, pay for 2" offer. Write the `.feature` file **first**, watch it
   fail, then implement it in `src/shop/`. That's the BDD loop.
4. Add a journey scenario: checking out with an empty cart (start with
   `Given I have an empty cart`). The checkout code already refuses it with
   "Your cart is empty", so you only need existing steps.
5. Write a step sentence with no definition and run both tools. behave prints a
   ready-to-paste code snippet for it; pytest-bdd fails with
   `StepDefinitionNotFoundError` naming the missing sentence.
