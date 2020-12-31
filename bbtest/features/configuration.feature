Feature: Service can be configured

  Scenario: configure log level to ERROR
    Given cnb-rates is configured with
      | property  | value |
      | LOG_LEVEL | ERROR |

    When restart unit "cnb-rates-rest.service"

    Then journalctl of "cnb-rates-rest.service" contains following
    """
      Log level set to ERROR
    """

  Scenario: configure log level to INFO
    Given cnb-rates is configured with
      | property  | value |
      | LOG_LEVEL | INFO  |

    When restart unit "cnb-rates-rest.service"

    Then journalctl of "cnb-rates-rest.service" contains following
    """
      Log level set to INFO
    """

  Scenario: configure log level to DEBUG
    Given cnb-rates is configured with
      | property  | value |
      | LOG_LEVEL | DEBUG |

    When restart unit "cnb-rates-rest.service"

    Then journalctl of "cnb-rates-rest.service" contains following
    """
      Log level set to DEBUG
    """

  Scenario: configure log level to INVALID
    Given cnb-rates is configured with
      | property  | value   |
      | LOG_LEVEL | INVALID |

    When restart unit "cnb-rates-rest.service"

    Then journalctl of "cnb-rates-rest.service" contains following
    """
      Log level set to INFO
    """
