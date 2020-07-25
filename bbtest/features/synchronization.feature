Feature: CNB Rates import

  Scenario: eventually downloads historic rates from cloud
    Given current time is "1991-03-09T14:29:57+0058"

    Then  fx-main CNB data are downloaded until "31.01.1991"
    And   fx-main CNB data are processed until "31.01.1991"
