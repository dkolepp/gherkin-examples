// The team kits this journey uses, pinned to exact versions.
// Like an Ansible requirements.yml: change a version to upgrade a team's kit.
// `make install` takes the example.com/shop modules from dist/ only (see Makefile).
module example.com/shop/journey

go 1.24

require (
	example.com/shop/bddcommon v1.1.0
	example.com/shop/cart v1.0.0
	example.com/shop/checkout v1.0.0
	example.com/shop/payments v1.0.0
	example.com/shop/shipping v1.1.0
)

require (
	github.com/cucumber/gherkin/go/v42 v42.0.0 // indirect
	github.com/cucumber/godog v0.16.0 // indirect
	github.com/cucumber/messages/go/v34 v34.2.0 // indirect
	github.com/google/uuid v1.6.0 // indirect
	github.com/hashicorp/go-immutable-radix v1.3.1 // indirect
	github.com/hashicorp/go-memdb v1.3.5 // indirect
	github.com/hashicorp/golang-lru v0.5.4 // indirect
	github.com/spf13/pflag v1.0.10 // indirect
)
