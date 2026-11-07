# Checkout journey

An end-to-end journey across the cart, payments, shipping and checkout teams.
Like an Ansible playbook, it only composes other teams' published work:

```
journey/
├── requirements.txt          which team packages, at which versions
├── Makefile                  install, test, catalog, check
└── features/
    ├── checkout.feature      the journey itself
    ├── environment.py        give each scenario fresh drivers (15 lines)
    └── steps/team_steps.py   import the teams' step libraries (5 lines)
```

There is no step code and no driver code here. Every sentence in
`checkout.feature` comes from a team package.

## Run it

From the repo root, build the team packages first (in real life they would
already be on your private package index):

```bash
make build                 # at the repo root: wheels into ../dist
cd journey
make test                  # install the pinned wheels, run the journey
make catalog               # every step you can use, with descriptions
make check                 # dry run: fails on undefined or duplicate steps
```

## Upgrading a team's library

Change its pin in `requirements.txt`, run `make check`, then `make test`.
Read the package's CHANGELOG first: a deprecated sentence still works but
prints a `StepDeprecationWarning` telling you the new wording.

## Why installation has two steps

`make install` first installs the team packages with `--no-index`, from your
own wheels only. Then it installs their third-party dependencies (behave,
parse) from PyPI. If both came from one command, a public package with the
same name and a higher version number could be installed instead of yours.
This is called dependency confusion.
