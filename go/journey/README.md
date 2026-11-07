# Checkout journey (Go)

An end-to-end journey across the cart, payments, shipping and checkout teams.
Like an Ansible playbook, it only composes other teams' published work:

```
go/journey/
├── go.mod            which team modules, at which versions (no replace directives)
├── journey_test.go   register the teams' kits, in order, and run the spec (about 20 lines)
├── Makefile          install, test, catalog, check
└── README.md

specs/journeys/
└── checkout.feature  the journey itself, shared by every language
```

There is no step code and no driver code here. Every sentence in
`specs/journeys/checkout.feature` comes from a team module. The journey does
not even import godog: `bddcommon.RunSpecs` runs it.

## Run it

Build the team modules first (in real life they would already be on your
private module proxy):

```bash
make -C go build           # from the repo root: modules into go/dist
cd go/journey
make test                  # fetch the pinned modules, run the journey
make catalog               # every step you can use
make check                 # duplicate sentences, then a strict run
```

Pass godog flags with `GODOG`: `make test GODOG='-godog.tags=@smoke'`, or
`GODOG=-godog.format=progress`. godog's tag syntax is the old one (`~@a` for
"not", `&&`, `,` for "or"): see "Working on one team's module" in
[`go/README.md`](../README.md).

## Upgrading a team's kit

Change its version in `go.mod`, run `make check`, then `make test`. Read the
module's CHANGELOG first: a deprecated sentence still works but prints a
`StepDeprecationWarning` naming the new wording. A new major version is a new
module path (`example.com/shop/shipping/v2`), so it is a new `require` line
and a new import path in `journey_test.go`.

## Why installation has two steps

`make install` mirrors the Python journey's two-step install:

1. **Team modules from `dist/` only.** `GOPROXY=file://.../dist`, nothing else
   in the list, so a module missing from `dist/` is an error, never a
   fall-back to a public proxy. `GONOSUMDB=example.com/shop` keeps their names
   away from the public checksum database. `GONOPROXY=none` (matches no
   module) stops a `GOPRIVATE` set elsewhere from sending them to version
   control instead.
2. **Third-party modules (godog and its dependencies) from proxy.golang.org.**
   Here `GOPRIVATE=example.com/shop` keeps any team module off the public
   proxy, and `GOVCS=example.com/shop:off` forbids fetching one from version
   control, so a team module that step 1 did not fetch fails the install.

The run itself uses `GOPROXY=off`: everything is in the journey's own module
cache by then. `GOWORK=off` makes sure `go/go.work` (which joins the team
modules from source, for development) plays no part.

`go.sum` is written on install, not committed, because `make build` rebuilds
the same versions from source and their hashes change whenever the source
does (see the exercise in [`go/README.md`](../README.md)). With real,
immutable releases, commit `go.sum`: it then pins the exact content of every
module, and a swapped artifact fails with a checksum mismatch.

## What `make check` can and can't catch

godog has no dry run. `make check` does what it can:

| Problem | `make check` | behave `--dry-run` |
|---|---|---|
| Undefined step | Yes, by running the scenarios strictly | Yes, without running |
| Pending step (`godog.ErrPending`) | Yes, strict run | n/a (behave has no pending) |
| Two definitions match a step in use | Yes: strict run reports "ambiguous" | Yes, at load time |
| The exact same sentence registered twice, used or not | Yes: compares the `-godog.definitions` list | Yes, at load time |
| Two different patterns that could overlap, in no scenario yet | No | No |
| A step that fails | Yes, but only because it runs everything | No |

Without strict mode godog passes a suite with undefined or pending steps, and
for a step that matches two definitions it silently uses the one registered
first. `bddcommon.RunSpecs` turns strict mode on for every run, so `make
test` already fails on all three. The difference is the duplicate check,
which catches a clash before any feature file uses the sentence.
