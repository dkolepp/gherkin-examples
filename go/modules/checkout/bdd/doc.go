// Package bdd is the checkout team's published BDD kit: driver + step library.
//
//	driver.go   CheckoutDriver: place orders, built from the other teams' drivers
//	steps.go    InitializeScenario: the godog steps and the Before hook
//
// Only checkout actions live here. The journey reuses the cart, payments and
// shipping teams' steps for everything else, instead of writing its own
// versions of them.
//
// Contract
//
//	Register   checkoutbdd.InitializeScenario LAST: after bddcommon, cartbdd,
//	           paymentsbdd and shippingbdd. godog runs Before hooks in the
//	           order they were registered.
//	Provides   Install(ctx) puts a *CheckoutDriver in the scenario's context
//	           (a Before hook); CheckoutFrom(ctx) reads it back
//	Needs      the ShopDriver, PaymentsDriver and ShippingDriver from the
//	           other kits, and the Outcome from bddcommon
//
// Import it as checkoutbdd "example.com/shop/checkout/bdd".
package bdd
