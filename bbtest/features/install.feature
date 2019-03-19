@install
Feature: Install package

  Scenario: install
    Given package "cnb-rates.deb" is installed
    Then  systemctl contains following
    """
      cnb-rates.service
      cnb-rates.path
      cnb-rates-rest.service
    """
