Feature: CNB Rates import

  Scenario: eventually downloads historic rates from cloud
    Given current time is "1991-03-09T14:29:57+0058"

    Then  all fx-main CNB data are downloaded
    And   all fx-main CNB data are processed
