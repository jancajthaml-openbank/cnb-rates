from behave import *
from helpers.shell import execute
from helpers.eventually import eventually
import datetime
import time
import os


@then('all fx-main CNB data are downloaded')
def step_impl(context):
  from_datetime = datetime.datetime.strptime('1991-01-01', '%Y-%m-%d').replace(tzinfo=datetime.timezone.utc)
  to_datetime = context.timeshift

  days = []
  while from_datetime.date() <= to_datetime.date():
    if from_datetime.isoweekday() not in (6, 7):
      days.append(from_datetime)
    from_datetime += datetime.timedelta(days=1)

  @eventually(10)
  def impl():
    for day in days:
      day_file = '/data/rates/cnb/raw/fx-main/{}'.format(day.strftime('%d.%m.%Y'))
      assert os.path.isfile(day_file), 'file {} not found'.format(day_file)
  impl()


@then('all fx-main CNB data are processed')
def step_impl(context):
  from_datetime = datetime.datetime.strptime('1991-01-01', '%Y-%m-%d').replace(tzinfo=datetime.timezone.utc)
  to_datetime = context.timeshift

  days = []
  while from_datetime.date() <= to_datetime.date():
    if from_datetime.isoweekday() not in (6, 7):
      days.append(from_datetime)
    from_datetime += datetime.timedelta(days=1)

  @eventually(10)
  def impl():
    for day in days:
      day_file = '/data/rates/cnb/processed/fx-main/d/{}'.format(day.strftime('%d.%m.%Y'))
      assert os.path.isfile(day_file), 'file {} not found'.format(day_file)
  impl()
