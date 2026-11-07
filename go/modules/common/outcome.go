package bddcommon

// Outcome remembers what went wrong in a "try to" step. It is shared by all
// teams, so every team's "try to" steps work with the same
// `Then I should be told "..."` step:
//
//	"When I add ..."        must succeed; an error fails the scenario.
//	"When I try to add ..." calls Outcome.Attempt; a later
//	                        `Then I should be told "..."` checks the error.
type Outcome struct {
	LastError error
}

// Attempt runs action and keeps its error (or nil) for a later check.
// Domain errors are what we want to see, so Attempt itself never fails.
func (o *Outcome) Attempt(action func() error) {
	o.LastError = action()
}
