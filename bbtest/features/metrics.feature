@metrics
Feature: Metrics test

  Scenario: cnb-rates-rest metrics have expected keys
    Given metrics file /tmp/reports/metrics.json should have following keys:
    """
      gatewayLatency
      importLatency
    """
    And metrics file /tmp/reports/metrics.json has permissions -rw-r--r--

  Scenario: cnb-rates-batch metrics have expected keys
    Given metrics file /tmp/reports/metrics.batch.json should have following keys:
    """
      daysProcessed
      monthsProcessed
    """
    And metrics file /tmp/reports/metrics.batch.json has permissions -rw-r--r--

  Scenario: cnb-rates-import metrics have expected keys
    Given metrics file /tmp/reports/metrics.import.json should have following keys:
    """
      daysImported
      gatewayLatency
      importLatency
      monthsImported
    """
    And metrics file /tmp/reports/metrics.import.json has permissions -rw-r--r--
