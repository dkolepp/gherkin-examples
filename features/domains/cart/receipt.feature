@receipt
Feature: Printed receipt

  # Test data can live outside the spec too: this loads test_data/standard_shop.json.
  # Use this for big or shared data sets; keep small, meaningful data inline.
  Background:
    Given the shop is set up from "standard_shop.json"

  Scenario: Receipt lists each line, the discount and the total
    Given I have a cart containing 2 "Coffee Beans"
    And I add 1 "Mug" to the cart
    And I have applied the discount code "SAVE10"
    When I print the receipt
    # Text between """ marks is a "doc string": a multi-line block of text.
    Then the receipt should read:
      """
      Coffee Beans    2 x  12.50 =   25.00
      Mug             1 x   8.00 =    8.00
      Subtotal                        33.00
      Discount (SAVE10)               -3.30
      TOTAL                           29.70
      """
