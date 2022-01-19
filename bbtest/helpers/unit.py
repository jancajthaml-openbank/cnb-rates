#!/usr/bin/env python3
# -*- coding: utf-8 -*-

import os
from openbank_testkit import Shell, Platform, Package


class UnitHelper(object):

  @staticmethod
  def default_config():
    return {
      "STORAGE": "{}/reports/blackbox-tests/data".format(os.getcwd()),
      "LOG_LEVEL": "DEBUG",
      "CNB_GATEWAY": "https://127.0.0.1:4000",
      "HTTP_PORT": "443",
      "SERVER_KEY": "/etc/cnb-rates/secrets/domain.local.key",
      "SERVER_CERT": "/etc/cnb-rates/secrets/domain.local.crt",
      "METRICS_CONTINUOUS": True,
      "METRICS_REFRESHRATE": "1s",
      "METRICS_OUTPUT": "{}/reports/blackbox-tests/metrics".format(os.getcwd())
    }

  def __init__(self, context):
    self.store = dict()
    self.units = list()
    self.context = context

  def download(self):
    version = os.environ.get('VERSION', '')
    meta = os.environ.get('META', '')

    if version.startswith('v'):
      version = version[1:]

    assert version, 'VERSION not provided'
    assert meta, 'META not provided'

    package = Package('cnb-rates')

    cwd = os.path.realpath('{}/../..'.format(os.path.dirname(__file__)))

    assert package.download(version, meta, '{}/packaging/bin'.format(cwd)), 'unable to download package cnb-rates'

    self.binary = '{}/packaging/bin/cnb-rates_{}_{}.deb'.format(cwd, version, Platform.arch)

  def configure(self, params = None):
    options = dict()
    options.update(UnitHelper.default_config())
    if params:
      options.update(params)

    os.makedirs('/etc/cnb-rates/conf.d', exist_ok=True)
    with open('/etc/cnb-rates/conf.d/init.conf', 'w') as fd:
      fd.write(str(os.linesep).join("CNB_RATES_{!s}={!s}".format(k, v) for (k, v) in options.items()))

  def collect_logs(self):
    (code, result, error) = Shell.run(['journalctl', '-o', 'cat', '--no-pager'])
    if code == 'OK':
      with open('reports/blackbox-tests/logs/journal.log', 'w') as fd:
        fd.write(result)

    for unit in set(self.__get_systemd_units() + self.units):
      (code, result, error) = Shell.run(['journalctl', '-o', 'cat', '-u', unit, '--no-pager'])
      if code != 'OK' or not result:
        continue
      with open('reports/blackbox-tests/logs/{}.log'.format(unit), 'w') as fd:
        fd.write(result)

  def teardown(self):
    self.collect_logs()
    for unit in self.__get_systemd_units():
      Shell.run(['systemctl', 'stop', unit])
    self.collect_logs()

  def __get_systemd_units(self):
    (code, result, error) = Shell.run(['systemctl', 'list-units', '--all', '--no-legend'])
    result = [item.replace('*', '').strip().split(' ')[0].strip() for item in result.split(os.linesep)]
    result = [item for item in result if "cnb-rates" in item]
    return result
