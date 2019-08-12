Feature: Uninstall package

  Scenario: uninstall
    Given package cnb-rates is uninstalled
    Then  systemctl does not contain following active units
      | name             | type    |
      | cnb-rates-import | timer   |
      | cnb-rates-rest   | service |
      | cnb-rates        | service |
      | cnb-rates        | path    |
