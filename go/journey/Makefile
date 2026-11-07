# The checkout journey, run against the team modules published in DIST.
#
#   make install   fetch the pinned team modules from DIST, the rest from the Go proxy
#   make test      run the journey (specs/journeys/ at the repo root)
#   make catalog   list every step the journey can use
#   make check     fail on duplicate sentences, then a strict run: undefined, pending, ambiguous
#   make clean     remove the journey's caches and go.sum
#
# Team modules come from DIST (default ../dist, built by `make build` in go/).
# In real life, point it at your private module proxy instead.

DIST ?= ../dist
GO   ?= go

# Which team modules to fetch, read from the require block in go.mod.
TEAM_MODULES := $(shell sed -n 's/^[[:space:]]*\(example\.com\/shop\/[^ ]*\) \(v[^ ]*\).*/\1@\2/p' go.mod)

# Like the Python journey's own venv: a module cache and build cache of the
# journey's own, so installing never mixes with development.
# The go command assumes a module version never changes, and keeps facts
# about it in the build cache. `make build` republishes the same versions
# from changed source (try the exercise in go/README.md), so install
# empties the build cache too.
export GOMODCACHE := $(CURDIR)/.gomodcache
export GOCACHE    := $(CURDIR)/.gocache
export GOFLAGS    := -modcacherw
export GOWORK     := off
export GOTOOLCHAIN := local

# Step 1: team modules from our own proxy only. GONOPROXY=none matches no
# module, so a GOPRIVATE set elsewhere cannot send them anywhere else.
# GONOSUMDB keeps their names away from the public checksum database.
FROM_DIST   := GOPROXY=file://$(abspath $(DIST)) GONOPROXY=none GONOSUMDB=example.com/shop
# Step 2: third-party modules (godog and its dependencies) from the public
# proxy. A team module that step 1 did not fetch must not be fetched now:
# GOPRIVATE keeps it off the public proxy and GOVCS=off blocks a direct fetch.
FROM_PUBLIC := GOPROXY=https://proxy.golang.org GOPRIVATE=example.com/shop GOVCS=example.com/shop:off
# Running: everything is in the cache by now, so no network at all.
OFFLINE     := GOPROXY=off

.PHONY: help install test catalog check clean

help:
	@sed -n 's/^#   //p' Makefile

install:
	@test -d $(DIST) || { echo "No $(DIST): run 'make build' in go/ first."; exit 1; }
	@# 0. Forget team modules fetched before (like pip --force-reinstall).
	rm -rf .gomodcache/example.com .gomodcache/cache/download/example.com .gocache
	rm -f go.sum
	@# 1. Team modules only from our own proxy, never from a public one,
	@#    so nobody can slip in a look-alike module with the same path.
	$(FROM_DIST) $(GO) mod download $(TEAM_MODULES)
	@# 2. Their third-party dependencies (godog and friends) from the Go proxy.
	@#    ("all" also records their checksums in go.sum.)
	$(FROM_PUBLIC) $(GO) mod download all

# `go test` with no package argument streams godog's output as it runs;
# `go test .` would show it only on failure.
# Extra godog flags: make test GODOG='-godog.tags=@smoke'
test: install
	$(OFFLINE) $(GO) test -count=1 -godog.format=pretty $(GODOG)

# godog's step list: one regular expression per line, then "# file:line -> func".
# It is printed in colour whatever -godog.no-colors says, so strip the codes.
STEP_LIST := $(OFFLINE) $(GO) test -count=1 -godog.definitions | sed -e 's/\x1b\[[0-9;]*m//g' -e '/^ok\|^PASS/d'

catalog: install
	@$(STEP_LIST)

check: install
	$(OFFLINE) $(GO) vet .
	@# Static: the same sentence registered twice, used or not.
	@dups="$$($(STEP_LIST) | sed -n 's/^\(\^.*\$$\) *# .*/\1/p' | sort | uniq -d)"; \
	  if [ -n "$$dups" ]; then echo "Duplicate step sentences:"; echo "$$dups"; exit 1; fi; \
	  echo "No duplicate step sentences."
	@# Strict run: undefined, pending and ambiguous steps fail. It runs the
	@# scenarios: godog has no dry run.
	$(OFFLINE) $(GO) test -count=1 -godog.strict -godog.format=progress

clean:
	@# The module cache is read-only by default; -modcacherw (GOFLAGS) and
	@# go clean -modcache both handle that.
	-$(GO) clean -modcache
	rm -rf .gomodcache .gocache
	rm -f go.sum
