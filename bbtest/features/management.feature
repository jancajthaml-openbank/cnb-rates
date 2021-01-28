Feature: System control

  Scenario: check units presence
    Then  systemctl contains following active units
      | name                   | type    |
      | cnb-rates              | service |
      | cnb-rates-watcher      | path    |
      | cnb-rates-watcher      | service |
      | cnb-rates-rest         | service |
      | cnb-rates-import       | service |
      | cnb-rates-batch        | service |
    And unit "cnb-rates-rest.service" is running

  Scenario: stop
    When stop unit "cnb-rates.service"
    Then unit "cnb-rates-rest.service" is not running

  Scenario: start
    When start unit "cnb-rates.service"
    Then unit "cnb-rates-rest.service" is running

  Scenario: restart
    When restart unit "cnb-rates.service"
    Then unit "cnb-rates-rest.service" is running
