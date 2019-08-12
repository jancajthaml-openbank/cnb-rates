Feature: Install package

  Scenario: install
    Given package cnb-rates is installed
    Then  systemctl contains following active units
      | name             | type    |
      | cnb-rates-import | timer   |
      | cnb-rates-rest   | service |
      | cnb-rates        | service |
      | cnb-rates        | path    |
