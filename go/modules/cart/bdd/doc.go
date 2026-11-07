// Package bdd is the cart team's published BDD kit: driver + step library +
// bundled data.
//
//	driver.go   ShopDriver: business actions on the cart (shopping, discounts, receipt)
//	steps.go    InitializeScenario: the godog steps and the Before hook
//	data/       bundled test data: standard_shop.json
//
// Contract
//
//	Register   cartbdd.InitializeScenario, after bddcommon.InitializeScenario
//	Provides   Install(ctx) puts a fresh *ShopDriver in the scenario's context
//	           (a Before hook); ShopFrom(ctx) reads it back
//	Needs      the Outcome from bddcommon
//
// Import it as cartbdd "example.com/shop/cart/bdd": every team's kit is a
// package named bdd.
package bdd
