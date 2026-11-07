# One repo, one set of Gherkin specs (specs/), one implementation per language.
# Each language has its own Makefile; this one only delegates to them.
#
#   make test             run every language's tests (today: Python only)
#   make journey          run every language's journey (today: Python only)
#   make clean            clean every language's build output
#
#   make python-test      python/: each team's behave and pytest runs, unit tests
#   make python-journey   python/: build the wheels, run the journey from them
#   make python-check     python/: journey dry run (undefined or duplicate steps)
#   make python-clean     python/: remove venvs, wheels and build folders
#
# Go targets (go-test, go-journey, ...) will be added with go/.

LANGUAGES ?= python

.PHONY: help test journey clean python-test python-journey python-check python-clean

help:
	@sed -n 's/^#   //p' Makefile

test: $(addsuffix -test,$(LANGUAGES))
journey: $(addsuffix -journey,$(LANGUAGES))
clean: $(addsuffix -clean,$(LANGUAGES))

# --- Python -----------------------------------------------------------------
python-test:
	$(MAKE) -C python test

python-journey:
	$(MAKE) -C python journey

python-check:
	$(MAKE) -C python build
	$(MAKE) -C python/journey check

python-clean:
	$(MAKE) -C python clean

# --- Go (coming) ------------------------------------------------------------
# go-test:
# 	$(MAKE) -C go test
