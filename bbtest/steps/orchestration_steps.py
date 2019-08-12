from behave import *
from helpers.shell import execute
import os
from helpers.eventually import eventually
import datetime


@given('current time is "{value}"')
def timeshift(context, value):
  import time

  (code, result, error) = execute([
    'timedatectl', 'set-ntp', '0'
  ])
  assert code == 0, "{}{}".format(result, error)

  (code, result, error) = execute([
    'timedatectl', 'set-local-rtc', '0'
  ])
  assert code == 0, "{}{}".format(result, error)

  context.timeshift = datetime.datetime.strptime(value, '%Y-%m-%dT%H:%M:%S%z').astimezone(datetime.timezone.utc)

  for _ in range(4):
    (code, result, error) = execute([
      'timedatectl', 'set-time', context.timeshift.strftime('%Y-%m-%d %H:%M:%S')
    ])
    assert code == 0, "{}{}".format(result, error)

    context.timeshift += datetime.timedelta(seconds=1)
    time.sleep(1)

    (code, result, error) = execute([
      "systemctl", "show", "-p", "SubState", 'cnb-rates-import.service'
    ])
    if 'SubState=running' in result:
      break


@given('package {package} is {operation}')
def step_impl(context, package, operation):
  if operation == 'installed':
    (code, result, error) = execute([
      "apt-get", "-y", "install", "-f", "/tmp/packages/{}.deb".format(package)
    ])
    assert code == 0
    assert os.path.isfile('/etc/init/cnb-rates.conf') is True

  elif operation == 'uninstalled':
    (code, result, error) = execute([
      "apt-get", "-y", "remove", package
    ])
    assert code == 0
    assert os.path.isfile('/etc/init/cnb-rates.conf') is False

  else:
    assert False


@given('systemctl contains following active units')
@then('systemctl contains following active units')
def step_impl(context):
  (code, result, error) = execute([
    "systemctl", "list-units", "--no-legend"
  ])
  assert code == 0

  items = []
  for row in context.table:
    items.append(row['name'] + '.' + row['type'])

  result = [item.split(' ')[0].strip() for item in result.split('\n')]
  result = [item for item in result if item in items]

  assert len(result) > 0


@given('systemctl does not contain following active units')
@then('systemctl does not contain following active units')
def step_impl(context):
  (code, result, error) = execute([
    "systemctl", "list-units", "--no-legend"
  ])
  assert code == 0

  items = []
  for row in context.table:
    items.append(row['name'] + '.' + row['type'])

  result = [item.split(' ')[0].strip() for item in result.split('\n')]
  result = [item for item in result if item in items]

  assert len(result) == 0


@given('unit "{unit}" is running')
@then('unit "{unit}" is running')
def unit_running(context, unit):
  @eventually(2)
  def impl():
    (code, result, error) = execute([
      "systemctl", "show", "-p", "SubState", unit
    ])

    assert code == 0
    assert 'SubState=running' in result
  impl()


@given('unit "{unit}" is not running')
@then('unit "{unit}" is not running')
def unit_not_running(context, unit):
  (code, result, error) = execute([
    "systemctl", "show", "-p", "SubState", unit
  ])

  assert code == 0
  assert 'SubState=dead' in result


@given('{operation} unit "{unit}"')
@when('{operation} unit "{unit}"')
def operation_unit(context, operation, unit):
  (code, result, error) = execute([
    "systemctl", operation, unit
  ])
  assert code == 0

  if operation == 'restart':
    unit_running(context, unit)


@given('cnb-rates is configured with')
def unit_is_configured(context):
  params = dict()
  for row in context.table:
    params[row['property']] = row['value']
  context.unit.configure(params)

  operation_unit(context, 'restart', 'cnb-rates-rest.service')
