# gherkin-examples

Learn **Gherkin** and **BDD** (Behaviour-Driven Development) with a tiny shop
split the way a real organisation might split it: four teams (cart, payments,
shipping, checkout) each publish a versioned *BDD kit* (a driver and a step
library), and a light journey project composes them into an end-to-end
checkout scenario without writing any step code of its own.

Three lessons run through it:

1. **Specs are separate from test code.** The same Gherkin runs under several
   runners and, soon, several languages.
2. **Teams own their parts.** Each team's rules, driver and steps are its own.
3. **Boundaries need contracts.** Step sentences, driver state and versions are
   public API between teams.

## Layout

```
specs/                  language-neutral Gherkin: the shared contract
├── cart/               cart.feature  discounts.feature  receipt.feature
├── payments/           payments.feature
├── shipping/           shipping.feature
├── checkout/           checkout.feature
└── journeys/           checkout.feature (the end-to-end journey)
python/                 the Python kits (behave and pytest-bdd): see python/README.md
├── packages/<team>/    one versioned package per team, plus its runner wiring
├── journey/            the journey, run from the built wheels
└── Makefile
docs/                   design documents and the Gherkin support matrix
Makefile                delegates to each language's Makefile
```

A Go version (`go/`, using godog) is coming. It will run the same `specs/`.

## Why the specs live on their own

A `.feature` file says *what* the shop does, in plain language. It doesn't
care which language or runner checks it. So the specs are the contract, and
each language's kits implement it:

- `specs/` holds only Gherkin: no step code, no runner config. Each team owns
  its folder there, and its kit in every language (see
  [`docs/CODEOWNERS.example`](docs/CODEOWNERS.example)).
- Each language keeps its own wiring next to its code and points its runners
  at `specs/`. For Python, see "How the runners find the shared specs" in
  [`python/README.md`](python/README.md).
- A changed sentence in `specs/` is a change to every language's kit. When
  each language's run is green, every kit agrees with the contract.

## Quick start

Needs `make`, and Python 3.10+ for the Python side.

```bash
make test             # every language's tests (today: Python)
make journey          # every language's journey (today: Python)
make python-test      # python/: each team's behave and pytest runs, unit tests
make python-journey   # python/: build the wheels, run the journey from them
make python-check     # python/: journey dry run: undefined or duplicate steps
make help             # all targets
```

The lessons, the Gherkin cheat sheet, the contracts between teams and the
exercises are in [`python/README.md`](python/README.md).

## Further reading

| Document | What it covers |
|---|---|
| [Gherkin support matrix](docs/gherkin-support-matrix.html) | Which Gherkin syntax behave, godog and pytest-bdd support, tested against each tool, plus Gherkin's version history. Download and open in a browser; [live version](https://claude.ai/artifact/MT6CLm637feEN39eXGMrjU) |
| [Journey Testing Platform](docs/journey-testing-platform.md) | Journey testing as a platform problem: personas, responsibilities, contracts at each boundary, friction and how to measure it |
| [Jira Xray: Test Management Model](docs/jira-xray-test-management.md) | What Xray is for, its entity model, Gherkin support, limits, and how it compares with the platform model |

The Markdown files are exported copies; each links to its live version.
