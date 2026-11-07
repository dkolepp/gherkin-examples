# Owner: shipping team.
@shipping
Feature: Shipping costs
  Each country has a flat shipping rate. Large orders ship free.

  Background:
    Given shipping costs:
      | country       | rate  |
      | United States | 5.00  |
      | Canada        | 12.00 |
    And orders of $50.00 or more ship free

  Scenario Outline: Cost depends on the destination and order value
    When I get a shipping quote to "<country>" for <order>
    Then shipping should cost <cost>

    Examples:
      | country       | order  | cost   |
      | United States | $20.00 | $5.00  |
      | Canada        | $20.00 | $12.00 |
      | Canada        | $49.99 | $12.00 |
      | Canada        | $50.00 | $0.00  |

  Scenario: Some places are not served
    When I try to get a shipping quote to "Antarctica" for $20.00
    Then I should be told "We don't ship to Antarctica"
