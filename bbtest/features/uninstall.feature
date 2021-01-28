Feature: Uninstall package

  Scenario: uninstall
    Given package cnb-rates is uninstalled
    Then  systemctl does not contain following active units
      | name              | type    |
      | cnb-rates         | service |
      | cnb-rates-import  | service |
      | cnb-rates-rest    | service |
      | cnb-rates-watcher | path    |
      | cnb-rates-watcher | service |
