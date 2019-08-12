Feature: Properly behaving units

  Scenario: lifecycle
    Then  systemctl contains following active units
      | name                   | type    |
      | cnb-rates              | path    |
      | cnb-rates              | service |
      | cnb-rates-rest         | service |
      | cnb-rates-import       | service |
      | cnb-rates-batch        | service |
    And unit "cnb-rates-rest.service" is running

    When stop unit "cnb-rates-rest.service"
    Then unit "cnb-rates-rest.service" is not running

    When start unit "cnb-rates-rest.service"
    Then unit "cnb-rates-rest.service" is running

    When restart unit "cnb-rates-rest.service"
    Then unit "cnb-rates-rest.service" is running
