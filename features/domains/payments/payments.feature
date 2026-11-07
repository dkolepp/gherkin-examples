# Owner: payments team. The detailed payment rules live here, not in journeys.
@payments
Feature: Card payments
  Each card has a spending limit. A payment that would go over the
  limit is declined.

  Background:
    Given the shopper has a card "4242" with a limit of $100.00

  Scenario: A payment within the limit is accepted
    When I charge $60.00 to card "4242"
    Then card "4242" should be charged $60.00

  Scenario: A payment over the limit is declined
    When I try to charge $150.00 to card "4242"
    Then I should be told "Card ending 4242 was declined"
    And card "4242" should not be charged

  Scenario: The limit covers all payments together
    When I charge $60.00 to card "4242"
    And I try to charge $50.00 to card "4242"
    Then I should be told "Card ending 4242 was declined"
    And card "4242" should be charged $60.00

  Scenario: An unknown card is refused
    When I try to charge $5.00 to card "9999"
    Then I should be told "Card ending 9999 is not recognised"
