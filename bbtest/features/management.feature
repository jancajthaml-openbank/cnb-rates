Feature: Properly behaving unit

  Scenario: control
    Given systemctl contains following
    """
      cnb-rates.service
      cnb-rates-import.service
      cnb-rates-rest.service
    """

    When stop unit "cnb-rates.service"
    Then unit "cnb-rates-import.service" is not running
    Then unit "cnb-rates-rest.service" is not running

    When start unit "cnb-rates.service"
    Then unit "cnb-rates-import.service" is running
    Then unit "cnb-rates-rest.service" is running

    When restart unit "cnb-rates.service"
    Then unit "cnb-rates-import.service" is running
    Then unit "cnb-rates-rest.service" is running
