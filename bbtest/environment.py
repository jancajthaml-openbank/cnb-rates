#!/usr/bin/env python3
# -*- coding: utf-8 -*-

import time
from helpers.unit import UnitHelper
from helpers.timeshift import TimeshiftHelper
from helpers.shell import execute
from mocks.cnb.server import CNBMock


def after_feature(context, feature):
  context.unit.collect_logs()


def before_all(context):
  context.timeshift = TimeshiftHelper(context)
  context.unit = UnitHelper(context)
  context.cnb = CNBMock(context)
  context.cnb.start()
  context.unit.configure()
  context.unit.download()
  context.timeshift.bootstrap()


def after_all(context):
  context.unit.teardown()
  context.cnb.stop()
  context.timeshift.teardown()
