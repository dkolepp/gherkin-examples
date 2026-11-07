# Python: develop, test and build every team package (run from python/,
# or from the repo root with `make -C python <target>`).
# The feature files are the shared specs in ../specs/; each package's
# pyproject.toml points behave and pytest-bdd at its team's folder there.
# The journey has its own Makefile in journey/.
#
#   make dev        install all packages in editable mode (for team development)
#   make test       run each team's own feature files and tests
#   make build      build each package into a wheel in python/dist/ (our "Galaxy")
#   make journey    build, then run the journey against the built wheels
#
# Limit to some packages with:  make test PACKAGES="cart payments"

PACKAGES ?= common cart payments shipping checkout
PYTHON   ?= python3
VENV     := .venv
BIN      := $(abspath $(VENV))/bin
DIST     := $(abspath dist)

.PHONY: help dev test build journey clean

help:
	@sed -n 's/^#   //p' Makefile

$(BIN)/python:
	$(PYTHON) -m venv $(VENV)
	$(BIN)/pip install -q --upgrade pip

dev: $(BIN)/python
	$(BIN)/pip install -q $(foreach p,$(PACKAGES),-e "packages/$(p)[dev]")

test: dev
	@set -e; for p in $(PACKAGES); do \
	  echo "== $$p"; \
	  if [ -d packages/$$p/features ]; then (cd packages/$$p && $(BIN)/behave -f progress --no-summary); fi; \
	  (cd packages/$$p && $(BIN)/pytest -q); \
	done

build: $(BIN)/python
	rm -rf $(DIST)
	@# Clear old build folders first: setuptools can reuse stale files from them.
	@set -e; for p in $(PACKAGES); do \
	  rm -rf packages/$$p/build; \
	  $(BIN)/pip wheel -q --no-deps -w $(DIST) packages/$$p; \
	  rm -rf packages/$$p/build; \
	done
	@ls -1 $(DIST)

journey: build
	$(MAKE) -C journey test DIST=$(DIST)

clean:
	rm -rf $(VENV) $(DIST) packages/*/build packages/*/src/*.egg-info
	$(MAKE) -C journey clean
