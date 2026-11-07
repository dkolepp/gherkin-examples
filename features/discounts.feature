@discounts
Feature: Discount codes
  Shoppers can enter a discount code to get a percentage off their order.
  Some codes only work above a minimum spend.

  Background:
    Given the shop sells these products:
      | name         | price | stock |
      | Coffee Beans | 12.50 | 10    |
      | Mug          |  8.00 | 5     |
    And these discount codes exist:
      | code   | percent off | minimum spend |
      | SAVE10 | 10          | 0.00          |
      | BIG20  | 20          | 50.00         |

  # A Rule groups scenarios that illustrate one business rule.
  Rule: A valid code takes a percentage off the subtotal

    @smoke
    Scenario Outline: Applying a code
      Given I have a cart containing <beans> "Coffee Beans"
      When I apply the discount code "<code>"
      Then the discount should be <discount>
      And the cart total should be <total>

      Examples:
        | beans | code   | discount | total  |
        | 1     | SAVE10 | $1.25    | $11.25 |
        | 4     | SAVE10 | $5.00    | $45.00 |
        | 4     | BIG20  | $10.00   | $40.00 |

  Rule: Some codes need a minimum spend

    Scenario: Code is rejected below the minimum spend
      Given I have a cart containing 1 "Coffee Beans"
      When I try to apply the discount code "BIG20"
      Then I should be told "BIG20 needs a minimum spend of $50.00"
      And the cart total should be $12.50

    Scenario: Discount stops applying if the cart drops below the minimum
      Given I have a cart containing 4 "Coffee Beans"
      And I have applied the discount code "BIG20"
      When I remove 1 "Coffee Beans" from the cart
      Then the discount should be $0.00
      And the cart total should be $37.50

  Rule: Unknown codes are rejected

    Scenario: A typo in the code
      Given I have a cart containing 1 "Mug"
      When I try to apply the discount code "SAVE100"
      Then I should be told "SAVE100 is not a valid discount code"
      And the cart total should be $8.00
