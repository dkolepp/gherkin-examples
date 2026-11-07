// Package bdd is the payments team's published BDD kit: driver + step library.
//
//	driver.go   PaymentsDriver: issue a card, charge it, check charges
//	steps.go    InitializeScenario: the godog steps and the Before hook
//
// Contract
//
//	Register   paymentsbdd.InitializeScenario, after bddcommon.InitializeScenario
//	Provides   Install(ctx) puts a fresh *PaymentsDriver in the scenario's
//	           context (a Before hook); PaymentsFrom(ctx) reads it back
//	Needs      the Outcome from bddcommon
//
// Import it as paymentsbdd "example.com/shop/payments/bdd".
package bdd
