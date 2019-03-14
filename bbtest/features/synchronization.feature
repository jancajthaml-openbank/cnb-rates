Feature: Verify download

  Scenario: eventually downloads historic rates from CNB cloud
    Given current time is "Mon Jan 4 14:29:59 1993"
    And cnb-rates is running with mocked CNB Gateway
    Then all CNB data are eventually synchronized
