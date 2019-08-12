import os
from helpers.unit import UnitHelper
from mocks.cnb.server import CNBMock


def after_feature(context, feature):
  context.unit.cleanup()


def before_all(context):
  context.unit = UnitHelper(context)
  context.cnb = CNBMock(context)

  os.system('mkdir -p /tmp/reports /tmp/reports/blackbox-tests /tmp/reports/blackbox-tests/logs /tmp/reports/blackbox-tests/metrics')
  os.system('rm -rf /tmp/reports/blackbox-tests/logs/*.log /tmp/reports/blackbox-tests/metrics/*.json')

  context.cnb.start()
  context.unit.download()
  context.unit.configure()


def after_all(context):
  context.cnb.stop()
  context.unit.teardown()
