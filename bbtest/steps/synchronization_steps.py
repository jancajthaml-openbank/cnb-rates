#!/usr/bin/env python3
# -*- coding: utf-8 -*-

from behave import *
from helpers.shell import execute
from helpers.eventually import eventually
import datetime
import time
import os


@then('fx-main CNB data are downloaded until "{endtime}"')
def step_impl(context, endtime):
  from_datetime = datetime.datetime.strptime('1991-01-01', '%Y-%m-%d').replace(tzinfo=datetime.timezone.utc)
  to_datetime = datetime.datetime.strptime(endtime, '%d.%m.%Y').replace(tzinfo=datetime.timezone.utc)

  days = []
  while from_datetime.date() <= to_datetime.date():
    if from_datetime.isoweekday() not in (6, 7):
      days.append(from_datetime)
    from_datetime += datetime.timedelta(days=1)

  @eventually(30)
  def wait_for_files():
    for day in days:
      day_file = '/tmp/reports/blackbox-tests/data/rates/cnb/raw/fx-main/{}'.format(day.strftime('%d.%m.%Y'))
      assert os.path.isfile(day_file), 'file {} not found'.format(day_file)
  wait_for_files()


@then('fx-main CNB data are processed until "{endtime}"')
def step_impl(context, endtime):
  from_datetime = datetime.datetime.strptime('1991-01-01', '%Y-%m-%d').replace(tzinfo=datetime.timezone.utc)
  to_datetime = datetime.datetime.strptime(endtime, '%d.%m.%Y').replace(tzinfo=datetime.timezone.utc)

  days = []
  while from_datetime.date() <= to_datetime.date():
    if from_datetime.isoweekday() not in (6, 7):
      days.append(from_datetime)
    from_datetime += datetime.timedelta(days=1)

  @eventually(30)
  def wait_for_files():
    for day in days:
      day_file = '/tmp/reports/blackbox-tests/data/rates/cnb/processed/fx-main/d/{}'.format(day.strftime('%d.%m.%Y'))
      assert os.path.isfile(day_file), 'file {} not found'.format(day_file)
  wait_for_files()
