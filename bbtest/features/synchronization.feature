Feature: Verify download

  Scenario: eventually downloads historic rates from CNB cloud

    Given cnb-rates is running with mocked CNB Gateway
    And current time is "Mon Jan 4 14:29:59 1993"

    Then all CNB data are eventually synchronized
