# One repo, one set of Gherkin specs (specs/), one implementation per language.
# Each language has its own Makefile; this one only delegates to them.
#
#   make test             run every language's tests (Python and Go)
#   make journey          run every language's journey (Python and Go)
#   make clean            clean every language's build output
#   make LANGUAGES=go ... limit to some languages
#
#   make python-test      python/: each team's behave and pytest runs, unit tests
#   make python-journey   python/: build the wheels, run the journey from them
#   make python-check     python/: journey dry run (undefined or duplicate steps)
#   make python-clean     python/: remove venvs, wheels and build folders
#
#   make go-test          go/: gofmt, go vet, each module's unit tests and godog runs
#   make go-journey       go/: publish the modules into go/dist, run the journey from them
#   make go-check         go/: journey check (duplicate sentences, then a strict run)
#   make go-clean         go/: remove go/dist and the journey's caches

LANGUAGES ?= python go

.PHONY: help test journey clean python-test python-journey python-check python-clean \
	go-test go-journey go-check go-clean

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

# --- Go ---------------------------------------------------------------------
go-test:
	$(MAKE) -C go test

go-journey:
	$(MAKE) -C go journey

go-check:
	$(MAKE) -C go build
	$(MAKE) -C go/journey check

go-clean:
	$(MAKE) -C go clean
