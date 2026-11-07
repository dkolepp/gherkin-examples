# Owner: checkout team (the journey owner).
# Its steps come from four step modules: cart, payments, shipping and checkout
# (see features/steps/). The checkout team owns only the checkout steps.
@journey @checkout
Feature: Checkout journey
  A shopper fills a cart, pays by card and gets the order shipped.

  This journey crosses three domains owned by different teams: cart,
  payments and shipping. It only checks that they work together.
  Each domain's detailed rules live in that team's own feature file.

  Background:
    Given the shop is set up from "standard_shop.json"
    And shipping is set up from "standard_shipping.json"
    And the shopper has a card "4242" with a limit of $500.00
    And the shopper lives in "Canada"

  @smoke
  Scenario: The order is paid for and shipped
    Given I have a cart containing 2 "Coffee Beans"
    And I add 1 "Mug" to the cart
    When I check out paying with card "4242"
    Then the order total should be $45.00
    And card "4242" should be charged $45.00
    And a shipment of 3 items should be sent to "Canada"
    And the cart should be empty

  # A rule that belongs to no single team: shipping is free from $50.00,
  # but is that before or after the cart's discount? Only a journey sees it.
  Rule: Free shipping is judged on the price after any discount

    Scenario: A $50.00 order ships free
      Given I have a cart containing 4 "Coffee Beans"
      When I check out paying with card "4242"
      Then the order total should be $50.00

    Scenario: A discount that drops the order below $50.00 brings shipping back
      Given I have a cart containing 4 "Coffee Beans"
      And I have applied the discount code "SAVE10"
      When I check out paying with card "4242"
      Then the order total should be $57.00

  Rule: If any part fails, nothing half-happens

    Scenario: A declined card means nothing is shipped
      Given the shopper has a card "0002" with a limit of $10.00
      And I have a cart containing 2 "Coffee Beans"
      When I try to check out paying with card "0002"
      Then I should be told "Card ending 0002 was declined"
      And no shipment should be sent
      And the cart should contain 2 "Coffee Beans"

    Scenario: An address we can't ship to means the card is not charged
      Given the shopper lives in "Antarctica"
      And I have a cart containing 1 "Mug"
      When I try to check out paying with card "4242"
      Then I should be told "We don't ship to Antarctica"
      And card "4242" should not be charged
