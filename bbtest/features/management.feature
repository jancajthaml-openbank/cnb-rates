Feature: Properly behaving unit

  Scenario: control
    Given systemctl contains following
    """
      cnb-rates.service
    """

    When stop unit "cnb-rates.service"
    Then unit "cnb-rates.service" is not running

    When start unit "cnb-rates.service"
    Then unit "cnb-rates.service" is running

    When restart unit "cnb-rates.service"
    Then unit "cnb-rates.service" is running
