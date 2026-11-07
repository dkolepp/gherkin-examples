package bddcommon

import (
	"fmt"
	"io"
	"os"
)

// DeprecationOutput is where deprecation warnings go. Tests may swap it.
var DeprecationOutput io.Writer = os.Stderr

// Deprecated retires a step sentence without breaking the feature files that
// use it. Changing a sentence is a breaking change for everyone who uses it,
// so for one minor release keep the old sentence working and warn:
//
//	sc.Then(oldPattern, func(ctx context.Context, qty int, country string) error {
//		bddcommon.Deprecated(`old words`, `new words`, "1.1.0")
//		return newStep(ctx, qty, country)
//	})
//
// Then remove it in the next major release.
func Deprecated(old, replacement, since string) {
	fmt.Fprintf(DeprecationOutput, "StepDeprecationWarning: %s\n", DeprecationMessage(old, replacement, since))
}

// DeprecationMessage is the text Deprecated prints.
func DeprecationMessage(old, replacement, since string) string {
	return fmt.Sprintf("Step %q is deprecated since %s. Use %q instead.", old, since, replacement)
}
