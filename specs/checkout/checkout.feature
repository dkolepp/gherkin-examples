# Owner: checkout team. Checkout's own rules. The end-to-end journey lives in
# the separate journey/ project, which installs this package as a dependency.
@checkout
Feature: Placing an order

  Background:
    Given the shop is set up from "standard_shop.json"
    And shipping is set up from "standard_shipping.json"
    And the shopper has a card "4242" with a limit of $500.00
    And the shopper lives in "Canada"

  Scenario: An empty cart can't be checked out
    Given I have an empty cart
    When I try to check out paying with card "4242"
    Then I should be told "Your cart is empty"
    And card "4242" should not be charged

  Scenario: Orders are numbered in sequence
    Given I have a cart containing 1 "Mug"
    And I check out paying with card "4242"
    And I add 1 "Mug" to the cart
    When I check out paying with card "4242"
    Then the order number should be "ORD-0002"
