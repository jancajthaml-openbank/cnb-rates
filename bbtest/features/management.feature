Feature: Properly behaving unit

  Scenario: control
    Given systemctl contains following
    """
      cnb-rates.service
      cnb-rates-rest.service
    """

    When stop unit "cnb-rates-rest.service"
    Then unit "cnb-rates-rest.service" is not running

    When start unit "cnb-rates-rest.service"
    Then unit "cnb-rates-rest.service" is running

    When restart unit "cnb-rates-rest.service"
    Then unit "cnb-rates-rest.service" is running
