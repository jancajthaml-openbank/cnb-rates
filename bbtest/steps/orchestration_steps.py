#!/usr/bin/env python3
# -*- coding: utf-8 -*-

from behave import *
import datetime
import os
from helpers.eventually import eventually
from openbank_testkit import Shell, Request


@given('current time is "{value}"')
def timeshift(context, value):
  context.new_epoch = datetime.datetime.strptime(value, '%Y-%m-%dT%H:%M:%S%z').astimezone(datetime.timezone.utc)

  @eventually(30)
  def wait_for_import_to_start():
    context.timeshift.set_date_time(context.new_epoch)
    context.new_epoch += datetime.timedelta(seconds=1)
    (code, result, error) = Shell.run(["systemctl", "restart", "cnb-rates-import.timer"])
    assert code == 'OK', str(result) + ' ' + str(error)
    (code, result, error) = Shell.run(["systemctl", "show", "-p", "SubState", 'cnb-rates-import.service'])
    assert code == 'OK', str(result) + ' ' + str(error)
    assert 'SubState=running' in result, str(result) + ' ' + str(error)

  @eventually(30)
  def wait_for_import_to_stop():
    (code, result, error) = Shell.run(["systemctl", "stop", "cnb-rates-import.service"])
    assert code == 'OK', str(result) + ' ' + str(error)
    (code, result, error) = Shell.run(["systemctl", "show", "-p", "SubState", 'cnb-rates-import.service'])
    assert code == 'OK', str(result) + ' ' + str(error)
    assert 'SubState=dead' in result, str(result) + ' ' + str(error)

  wait_for_import_to_start()
  del context.new_epoch
  wait_for_import_to_stop()


@given('package {package} is {operation}')
def step_impl(context, package, operation):
  if operation == 'installed':
    (code, result, error) = Shell.run(["apt-get", "install", "-f", "-qq", "-o=Dpkg::Use-Pty=0", "-o=Dpkg::Options::=--force-confold", context.unit.binary])
    assert code == 'OK', "unable to install with code {} and {} {}".format(code, result, error)
    assert os.path.isfile('/etc/cnb-rates/conf.d/init.conf') is True
  elif operation == 'uninstalled':
    (code, result, error) = Shell.run(["apt-get", "-y", "remove", package])
    assert code == 'OK', "unable to uninstall with code {} and {} {}".format(code, result, error)
    (code, result, error) = Shell.run(["apt-get", "-y", "purge", package])
    assert code == 'OK', "unable to purge with code {} and {} {}".format(code, result, error)
    assert os.path.isfile('/etc/cnb-rates/conf.d/init.conf') is False
  else:
    assert False, 'unknown operation {}'.format(operation)


@given('systemctl contains following active units')
@then('systemctl contains following active units')
def step_impl(context):
  (code, result, error) = Shell.run(["systemctl", "list-units", "--all", "--no-legend", "--state=active"])
  assert code == 'OK', str(result) + ' ' + str(error)
  items = []
  for row in context.table:
    items.append(row['name'] + '.' + row['type'])
  result = [item.replace('*', '').strip().split(' ')[0].strip() for item in result.split(os.linesep)]
  result = [item for item in result if item in items]
  assert len(result) > 0, 'units not found'


@given('systemctl does not contain following active units')
@then('systemctl does not contain following active units')
def step_impl(context):
  (code, result, error) = Shell.run(["systemctl", "list-units", "--all", "--no-legend", "--state=active"])
  assert code == 'OK', str(result) + ' ' + str(error)
  items = []
  for row in context.table:
    items.append(row['name'] + '.' + row['type'])
  result = [item.replace('*', '').strip().split(' ')[0].strip() for item in result.split(os.linesep)]
  result = [item for item in result if item in items]
  assert len(result) == 0, 'units {} found'.format(result)


@given('unit "{unit}" is running')
@then('unit "{unit}" is running')
def unit_running(context, unit):
  @eventually(10)
  def wait_for_unit_state_change():
    (code, result, error) = Shell.run(["systemctl", "show", "-p", "SubState", unit])
    assert code == 'OK', str(result) + ' ' + str(error)
    assert 'SubState=running' in result, result

  wait_for_unit_state_change()

  if 'cnb-rates-rest' in unit:
    request = Request(method='GET', url="https://127.0.0.1/health")

    @eventually(5)
    def wait_for_healthy():
      response = request.do()
      assert response.status == 200, str(response.status)

    wait_for_healthy()


@given('unit "{unit}" is not running')
@then('unit "{unit}" is not running')
def unit_not_running(context, unit):
  @eventually(10)
  def wait_for_unit_state_change():
    (code, result, error) = Shell.run(["systemctl", "show", "-p", "SubState", unit])
    assert code == 'OK', str(result) + ' ' + str(error)
    assert 'SubState=dead' in result, result

  wait_for_unit_state_change()


@given('{operation} unit "{unit}"')
@when('{operation} unit "{unit}"')
def operation_unit(context, operation, unit):
  (code, result, error) = Shell.run(["systemctl", operation, unit])
  assert code == 'OK', str(result) + ' ' + str(error)


@given('{unit} is configured with')
def unit_is_configured(context, unit):
  params = dict()
  for row in context.table:
    params[row['property']] = row['value']
  context.unit.configure(params)
