Feature: Metrics test

  Scenario: metrics have expected keys
    Given cnb-rates is configured with
      | property            | value |
      | METRICS_REFRESHRATE |    1s |

    Then metrics file reports/blackbox-tests/metrics/metrics.json should have following keys:
      | key             |
      | gatewayLatency  |
      | importLatency   |
    And metrics file reports/blackbox-tests/metrics/metrics.json has permissions -rw-r--r--
