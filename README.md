# gherkin-examples

Learn **Gherkin** and **BDD** (Behaviour-Driven Development) in Python, using
both [behave](https://behave.readthedocs.io/) and
[pytest-bdd](https://pytest-bdd.readthedocs.io/).

The example is a tiny shop: products, a cart, discount codes and a printed
receipt. The same `.feature` files are run by **both** tools, which shows the
main lesson: *the test cases (specs) are separate from the test code.*

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
  tests/bdd/conftest.py (pytest-bdd)
        │
        ▼
bdd_support/driver.py     3. HOW    Business actions in Python ("add to cart",
                                    "apply code"). Shared by both runners.
        │
        ▼
src/shop/                 4. APP    The real application. Knows nothing about tests.

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
| `features/cart.feature` | Background, data tables, tags, Scenario Outline with two Examples tables |
| `features/discounts.feature` | `Rule:` blocks, Outline, expecting errors |
| `features/receipt.feature` | Doc strings (`"""`), loading data from `test_data/` |
| `features/environment.py` | behave hooks (fresh driver per scenario) |
| `features/steps/shop_steps.py` | behave step definitions |
| `tests/bdd/conftest.py` | pytest-bdd step definitions + `shop` fixture |
| `tests/bdd/test_*.py` | Bind feature files to pytest (`scenarios(...)` or `@scenario`) |
| `tests/unit/` | Plain unit tests, for contrast |
| `bdd_support/driver.py` | `ShopDriver`, shared by both runners |
| `bdd_support/types.py` | Custom `{amount:Money}` / `{qty:Qty}` step types |
| `src/shop/` | The application |

Tip: open `features/steps/shop_steps.py` and `tests/bdd/conftest.py` side by
side. Same sentences, same driver calls, different glue.

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

## Running a subset

| Goal | behave | pytest-bdd |
|------|--------|-----------|
| Only `@smoke` scenarios | `behave --tags=smoke` | `pytest -m smoke` |
| Skip a tag | `behave --tags="not receipt"` | `pytest -m "not receipt"` |
| One feature file | `behave features/cart.feature` | `pytest tests/bdd/test_cart.py` |
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

1. Add a `Scenario` to `cart.feature`: removing a "Mug" from an empty cart
   leaves the cart empty. You can write it using only existing steps, so it
   should pass straight away in both tools.
2. Write a new step, e.g. `Then the cart should have {n:Qty} different products`.
   Implement it for behave *and* pytest-bdd (one line each, plus a driver method
   if needed).
3. Add a "buy 3, pay for 2" offer. Write the `.feature` file **first**, watch it
   fail, then implement it in `src/shop/`. That's the BDD loop.
4. Write a step sentence with no definition and run both tools. behave prints a
   ready-to-paste code snippet for it; pytest-bdd fails with
   `StepDefinitionNotFoundError` naming the missing sentence.
