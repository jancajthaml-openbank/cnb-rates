Feature: Verify service

  Scenario: properly installed debian package

    Given cnb-rates is running
    Then systemctl contains following
    """
      cnb-rates.service
    """

  Scenario: configure log level

    Given cnb-rates is running with following configuration
    """
      CNB_RATES_LOG_LEVEL=DEBUG
      CNB_RATES_SYNC_RATE=1h
    """
    Then journalctl of "cnb-rates.service" contains following
    """
      Log level set to DEBUG
    """

    Given cnb-rates is running with following configuration
    """
      CNB_RATES_LOG_LEVEL=ERROR
      CNB_RATES_SYNC_RATE=1h
    """
    Then journalctl of "cnb-rates.service" contains following
    """
      Log level set to ERROR
    """

    Given cnb-rates is running with following configuration
    """
      CNB_RATES_LOG_LEVEL=INFO
      CNB_RATES_SYNC_RATE=1h
    """
    Then journalctl of "cnb-rates.service" contains following
    """
      Log level set to INFO
    """
