@uninstall
Feature: Unnstall package

  Scenario: uninstall
    Given package "cnb-rates" is uninstalled
    Then  systemctl does not contains following
    """
      cnb-rates.service
      cnb-rates.path
      cnb-rates-rest.service
    """
