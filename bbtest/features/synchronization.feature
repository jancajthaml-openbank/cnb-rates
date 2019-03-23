Feature: Verify download

  Scenario: eventually downloads historic rates from CNB cloud
    Given current time is "Mon Jan 4 14:29:59 1993"
    And   cnb-rates is running with mocked CNB Gateway

    Then  all fx-main CNB data are downloaded
    And   all fx-main CNB data are processed
