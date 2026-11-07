# Go: develop, test and build every team module (run from go/,
# or from the repo root with `make -C go <target>`).
# The feature files are the shared specs in ../specs/; each module's
# features/ test points godog at its team's folder there.
# The journey has its own Makefile in journey/.
#
#   make dev        download the dependencies (go.work joins the modules from source)
#   make test       gofmt and go vet, then each module's unit tests and own feature files
#   make build      publish each module at its version into go/dist/ (our module proxy)
#   make journey    build, then run the journey against the published modules
#   make clean      remove dist/ and the journey's module cache
#
# Limit to some modules with:  make test MODULES="cart payments"

MODULES ?= common cart payments shipping checkout
GO      ?= go

# Use the Go on this machine; never download another toolchain.
export GOTOOLCHAIN := local

.PHONY: help dev test build journey clean

help:
	@sed -n 's/^#   //p' Makefile

dev:
	@set -e; for m in $(MODULES); do (cd modules/$$m && $(GO) list -deps -test ./... >/dev/null); done
	cd tools && $(GO) list -deps ./... >/dev/null

# In features/, `go test` with no package argument streams godog's output.
test:
	@unformatted="$$(gofmt -l modules tools journey/*.go)"; \
	  if [ -n "$$unformatted" ]; then echo "gofmt needed:"; echo "$$unformatted"; exit 1; fi
	@set -e; for m in $(MODULES); do \
	  echo "== $$m"; \
	  (cd modules/$$m && $(GO) vet ./... \
	    && $(GO) test -count=1 $$($(GO) list ./... | grep -v '/features$$') \
	    && if [ -d features ]; then cd features && $(GO) test -count=1; fi); \
	done
	@echo "== tools"; cd tools && $(GO) vet ./...

build:
	rm -rf ./dist
	$(GO) run ./tools/publish -out dist $(addprefix modules/,$(MODULES))

journey: build
	$(MAKE) -C journey test DIST=$(CURDIR)/dist

clean:
	rm -rf ./dist
	$(MAKE) -C journey clean
