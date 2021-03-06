#!/usr/bin/env python3
# -*- coding: utf-8 -*-

from behave import *
from helpers.shell import execute
import os
from helpers.eventually import eventually
import datetime


@given('current time is "{value}"')
def timeshift(context, value):
  context.new_epoch = datetime.datetime.strptime(value, '%Y-%m-%dT%H:%M:%S%z').astimezone(datetime.timezone.utc)

  @eventually(30)
  def wait_for_import_to_start():
    context.timeshift.set_date_time(context.new_epoch)
    context.new_epoch += datetime.timedelta(seconds=1)
    (code, result, error) = execute(["systemctl", "restart", "cnb-rates-import.timer"])
    assert code == 0, str(result) + ' ' + str(error)
    (code, result, error) = execute(["systemctl", "show", "-p", "SubState", 'cnb-rates-import.service'])
    assert code == 0, str(result) + ' ' + str(error)
    assert 'SubState=running' in result, str(result) + ' ' + str(error)

  @eventually(30)
  def wait_for_import_to_stop():
    (code, result, error) = execute(["systemctl", "stop", "cnb-rates-import.service"])
    assert code == 0, str(result) + ' ' + str(error)
    (code, result, error) = execute(["systemctl", "show", "-p", "SubState", 'cnb-rates-import.service'])
    assert code == 0, str(result) + ' ' + str(error)
    assert 'SubState=dead' in result, str(result) + ' ' + str(error)

  wait_for_import_to_start()
  del context.new_epoch
  wait_for_import_to_stop()


@given('package {package} is {operation}')
def step_impl(context, package, operation):
  if operation == 'installed':
    (code, result, error) = execute(["apt-get", "install", "-f", "-qq", "-o=Dpkg::Use-Pty=0", "-o=Dpkg::Options::=--force-confold", "/tmp/packages/{}.deb".format(package)])
    assert code == 0, "unable to install with code {} and {} {}".format(code, result, error)
    assert os.path.isfile('/etc/cnb-rates/conf.d/init.conf') is True
  elif operation == 'uninstalled':
    (code, result, error) = execute(["apt-get", "-y", "remove", package])
    assert code == 0, "unable to uninstall with code {} and {} {}".format(code, result, error)
    (code, result, error) = execute(["apt-get", "-y", "purge", package])
    assert code == 0, "unable to purge with code {} and {} {}".format(code, result, error)
    assert os.path.isfile('/etc/cnb-rates/conf.d/init.conf') is False
  else:
    assert False, 'unknown operation {}'.format(operation)


@given('systemctl contains following active units')
@then('systemctl contains following active units')
def step_impl(context):
  (code, result, error) = execute(["systemctl", "list-units", "--no-legend", "--state=active"])
  assert code == 0, str(result) + ' ' + str(error)
  items = []
  for row in context.table:
    items.append(row['name'] + '.' + row['type'])
  result = [item.split(' ')[0].strip() for item in result.split(os.linesep)]
  result = [item for item in result if item in items]
  assert len(result) > 0, 'units not found'


@given('systemctl does not contain following active units')
@then('systemctl does not contain following active units')
def step_impl(context):
  (code, result, error) = execute(["systemctl", "list-units", "--no-legend", "--state=active"])
  assert code == 0, str(result) + ' ' + str(error)
  items = []
  for row in context.table:
    items.append(row['name'] + '.' + row['type'])
  result = [item.split(' ')[0].strip() for item in result.split(os.linesep)]
  result = [item for item in result if item in items]
  assert len(result) == 0, 'units {} found'.format(result)


@given('unit "{unit}" is running')
@then('unit "{unit}" is running')
def unit_running(context, unit):
  @eventually(10)
  def wait_for_unit_state_change():
    (code, result, error) = execute(["systemctl", "show", "-p", "SubState", unit])
    assert code == 0, str(result) + ' ' + str(error)
    assert 'SubState=running' in result, '{} {}'.format(unit, result)

  wait_for_unit_state_change()


@given('unit "{unit}" is not running')
@then('unit "{unit}" is not running')
def unit_not_running(context, unit):
  @eventually(10)
  def wait_for_unit_state_change():
    (code, result, error) = execute(["systemctl", "show", "-p", "SubState", unit])
    assert code == 0, str(result) + ' ' + str(error)
    assert 'SubState=running' not in result, '{} {}'.format(unit, result)

  wait_for_unit_state_change()


@given('{operation} unit "{unit}"')
@when('{operation} unit "{unit}"')
def operation_unit(context, operation, unit):
  (code, result, error) = execute(["systemctl", operation, unit])
  assert code == 0, str(result) + ' ' + str(error)


@given('{unit} is configured with')
def unit_is_configured(context, unit):
  params = dict()
  for row in context.table:
    params[row['property']] = row['value']
  context.unit.configure(params)
