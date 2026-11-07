# The checkout journey. It writes no step code: it installs the teams'
# published packages (built wheels) and runs its feature files with them.
#
#   make install   create .venv and install the pinned team packages
#   make test      run the journey
#   make catalog   list every step the installed packages offer
#   make check     dry run: fail on undefined or duplicate steps
#   make clean     remove .venv
#
# Packages come from DIST (default ../dist, built by `make build` at the
# repo root). In real life, point it at your private package index instead.

DIST   ?= ../dist
PYTHON ?= python3
VENV   := .venv
BIN    := $(VENV)/bin

.PHONY: help install test catalog check clean

help:
	@sed -n 's/^#   //p' Makefile

$(BIN)/python:
	$(PYTHON) -m venv $(VENV)
	$(BIN)/pip install -q --upgrade pip

install: $(BIN)/python
	@# 1. Team packages only from our own wheels, never from public PyPI,
	@#    so nobody can slip in a look-alike package with the same name.
	$(BIN)/pip install -q --force-reinstall --no-deps --no-index --find-links $(DIST) -r requirements.txt
	@# 2. Their third-party dependencies (behave, parse) from PyPI.
	$(BIN)/pip install -q -r requirements.txt

test: install
	$(BIN)/behave

catalog: install
	$(BIN)/behave --steps-catalog

check: install
	$(BIN)/behave --dry-run --format progress

clean:
	rm -rf $(VENV)
