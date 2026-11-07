// Package cart is a tiny shop domain used as the "system under test":
// products, a cart, discount codes and a printed receipt.
//
// Nothing in this package knows about Gherkin or godog. That is the point:
// the application code stays clean, and tests talk to it from outside
// (see the bdd package for the cart team's BDD kit).
//
// Amounts of money are int64 counts of cents.
package cart
