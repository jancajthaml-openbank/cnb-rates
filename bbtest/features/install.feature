Feature: Install package

  Scenario: install
    Given package cnb-rates is installed
    Then  systemctl contains following active units
      | name              | type    |
      | cnb-rates         | service |
      | cnb-rates-import  | service |
      | cnb-rates-rest    | service |
      | cnb-rates-watcher | path    |
      | cnb-rates-watcher | service |