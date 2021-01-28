#!/usr/bin/env python3
# -*- coding: utf-8 -*-

import time
from helpers.unit import UnitHelper
from helpers.timeshift import TimeshiftHelper
from helpers.shell import execute
from mocks.cnb.server import CNBMock
from helpers.logger import logger


def before_feature(context, feature):
  context.log.info('')
  context.log.info('  (FEATURE) {}'.format(feature.name))


def before_scenario(context, scenario):
  context.log.info('')
  context.log.info('  (SCENARIO) {}'.format(scenario.name))
  context.log.info('')


def after_feature(context, feature):
  context.unit.collect_logs()


def before_all(context):
  context.log = logger()
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
