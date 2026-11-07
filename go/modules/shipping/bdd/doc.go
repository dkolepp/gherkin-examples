// Package bdd is the shipping team's published BDD kit: driver + step
// library + bundled data.
//
//	driver.go   ShippingDriver: rates, quotes and shipments
//	steps.go    InitializeScenario: the godog steps and the Before hook
//	data/       bundled test data: standard_shipping.json
//
// Contract
//
//	Register   shippingbdd.InitializeScenario, after bddcommon.InitializeScenario
//	Provides   Install(ctx) puts a fresh *ShippingDriver in the scenario's
//	           context (a Before hook); ShippingFrom(ctx) reads it back
//	Needs      the Outcome from bddcommon
//
// Import it as shippingbdd "example.com/shop/shipping/bdd".
package bdd
