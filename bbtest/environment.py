#!/usr/bin/env python
# -*- coding: utf-8 -*-

import os
from helpers.unit import UnitHelper
from mocks.cnb.server import CNBMock


def after_feature(context, feature):
  context.unit.cleanup()


def before_all(context):
  context.unit = UnitHelper(context)
  context.cnb = CNBMock(context)
  context.cnb.start()
  context.unit.download()
  context.unit.configure()


def after_all(context):
  context.cnb.stop()
  context.unit.teardown()
  if os.path.isdir('/data'):
    os.system('cp -r /data/* /tmp/reports/blackbox-tests/data/')
