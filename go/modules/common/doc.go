// Package bddcommon holds the building blocks every team's BDD kit shares.
//
//	outcome.go      Outcome: remembers the error from a "try to" step
//	money.go        Money ($12.50 in steps, "12.50" in tables) and Qty parsing
//	tables.go       data tables as rows of column name -> cell
//	data.go         loads JSON test data (bundled defaults or your own file)
//	deprecation.go  marks a step sentence as deprecated
//	steps.go        Install and the shared `I should be told "..."` step
//	runner.go       RunSpecs: runs feature files kept in specs/ with godog
//
// It is test code: application packages never import it.
package bddcommon
