# Lines starting with # are comments.
# Tags (@cart, @smoke) let you pick which scenarios to run.
@cart
Feature: Shopping cart
  As a shopper
  I want to put products in a cart
  So that I can buy several things at once

  # Background runs before EVERY scenario in this file.
  # The table is a "data table": the step definition receives it as rows.
  Background:
    Given the shop sells these products:
      | name         | price | stock |
      | Coffee Beans | 12.50 | 10    |
      | Mug          |  8.00 | 3     |
      | Tea Bags     |  4.25 | 0     |

  @smoke
  Scenario: A new cart is empty
    Given I have an empty cart
    Then the cart should be empty
    And the cart total should be $0.00

  @smoke
  Scenario: Adding products updates the total
    Given I have an empty cart
    When I add 2 "Coffee Beans" to the cart
    And I add 1 "Mug" to the cart
    Then the cart should contain:
      | product      | quantity |
      | Coffee Beans | 2        |
      | Mug          | 1        |
    And the cart total should be $33.00

  Scenario: Removing a product
    Given I have a cart containing 2 "Coffee Beans"
    When I remove 1 "Coffee Beans" from the cart
    Then the cart should contain 1 "Coffee Beans"
    And the cart total should be $12.50

  # "try to" means we expect it might fail; a Then step checks the message.
  Scenario: Cannot add more than is in stock
    Given I have an empty cart
    When I try to add 4 "Mug" to the cart
    Then I should be told "Only 3 Mug left in stock"
    And the cart should be empty

  # A Scenario Outline runs once per row of each Examples table.
  # <placeholders> are replaced with the values from the row.
  Scenario Outline: Stock limits are enforced
    Given I have an empty cart
    When I try to add <requested> "<product>" to the cart
    Then the cart should contain <added> "<product>"

    Examples: Within stock
      | product      | requested | added |
      | Coffee Beans | 10        | 10    |
      | Mug          | 3         | 3     |

    Examples: Beyond stock
      | product      | requested | added |
      | Mug          | 4         | 0     |
      | Tea Bags     | 1         | 0     |
