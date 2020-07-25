#!/usr/bin/env python3
# -*- coding: utf-8 -*-

import time
from helpers.unit import UnitHelper
from helpers.shell import execute
from mocks.cnb.server import CNBMock


def after_feature(context, feature):
  context.unit.collect_logs()


def before_all(context):
  context.unit = UnitHelper(context)
  context.cnb = CNBMock(context)
  context.cnb.start()
  context.unit.configure()
  context.unit.download()
  execute(['timedatectl', 'set-ntp', '0'])
  execute(['timedatectl', 'set-local-rtc', '0'])


def after_all(context):
  context.unit.teardown()
  context.cnb.stop()
  execute(['timedatectl', 'set-ntp', '1'])
  execute(['timedatectl', 'set-local-rtc', '1'])
  time.sleep(2)


