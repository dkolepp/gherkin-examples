package bddcommon

import (
	"flag"
	"testing"

	"github.com/cucumber/godog"
)

// options are the godog options for RunSpecs. Command-line flags override
// them: go test -godog.tags=@smoke, -godog.format=pretty, -godog.definitions.
// (godog's tag syntax is the old one: "@a && ~@b", and "," for OR.)
var options = godog.Options{
	Format: "progress",
	Strict: true, // fail on undefined, pending and ambiguous steps
}

func init() {
	godog.BindFlags("godog.", flag.CommandLine, &options)
}

// RunSpecs runs the feature files under paths with godog, each scenario as a
// Go subtest of t, with each kit's InitializeScenario in the order given.
//
// The feature files are the shared, language-neutral specs at the repo root,
// so paths usually point outside the module ("../../../specs/cart").
// -godog.paths=a,b on the command line replaces them.
//
// With -godog.definitions it lists every step and runs nothing.
func RunSpecs(t *testing.T, paths []string, kits ...func(*godog.ScenarioContext)) {
	t.Helper()
	opts := options
	opts.TestingT = t
	if len(opts.Paths) == 0 {
		opts.Paths = paths
	}
	suite := godog.TestSuite{
		Name: t.Name(),
		ScenarioInitializer: func(sc *godog.ScenarioContext) {
			for _, kit := range kits {
				kit(sc)
			}
		},
		Options: &opts,
	}
	status := suite.Run()
	if opts.ShowStepDefinitions {
		return // godog listed the steps and ran nothing
	}
	if status != 0 {
		t.Fatalf("godog finished with status %d", status)
	}
}
