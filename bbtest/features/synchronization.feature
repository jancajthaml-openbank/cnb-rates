Feature: Verify synchronization

  Scenario: eventually synchronizes everything from CNB cloud

    Given current time is "27.1.1992"
    When cnb-rates is running with mocked CNB Gateway
    Then all CNB data are eventually synchronized
