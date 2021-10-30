#!/usr/bin/env python3
# -*- coding: utf-8 -*-

import docker
import platform
import tarfile
import tempfile
import errno
import os
import subprocess
from helpers.shell import execute


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

  def get_arch(self):
    return {
      'x86_64': 'amd64',
      'armv7l': 'armhf',
      'armv8': 'arm64'
    }.get(platform.uname().machine, 'amd64')

  def __init__(self, context):
    self.arch = self.get_arch()

    self.store = dict()
    self.image_version = None
    self.debian_version = None
    self.units = list()
    self.docker = docker.from_env()
    self.context = context

  def download(self):
    failure = None
    os.makedirs('/tmp/packages', exist_ok=True)

    self.image_version = os.environ.get('IMAGE_VERSION', '')
    self.debian_version = os.environ.get('UNIT_VERSION', '')

    if self.debian_version.startswith('v'):
      self.debian_version = self.debian_version[1:]

    assert self.image_version, 'IMAGE_VERSION not provided'
    assert self.debian_version, 'UNIT_VERSION not provided'

    image = 'openbank/cnb-rates:{}'.format(self.image_version)
    package = '/opt/artifacts/cnb-rates_{}_{}.deb'.format(self.debian_version, self.arch)
    target = '/tmp/packages/cnb-rates.deb'

    temp = tempfile.NamedTemporaryFile(delete=True)
    try:
      with open(temp.name, 'w') as fd:
        fd.write(str(os.linesep).join([
          'FROM alpine',
          'COPY --from={} {} {}'.format(image, package, target)
        ]))

      image, stream = self.docker.images.build(fileobj=temp, rm=True, pull=False, tag='bbtest_artifacts-scratch')
      for chunk in stream:
        if not 'stream' in chunk:
          continue
        for line in chunk['stream'].splitlines():
          l = line.strip(os.linesep)
          if not len(l):
            continue
          print(l)

      scratch = self.docker.containers.run('bbtest_artifacts-scratch', ['/bin/true'], detach=True)

      tar_name = tempfile.NamedTemporaryFile(delete=True)
      with open(tar_name.name, 'wb') as fd:
        bits, stat = scratch.get_archive(target)
        for chunk in bits:
          fd.write(chunk)

      archive = tarfile.TarFile(tar_name.name)
      archive.extract(os.path.basename(target), os.path.dirname(target))

      (code, result, error) = execute(['dpkg', '-c', target])
      if code != 'OK':
        raise RuntimeError('code: {}, stdout: [{}], stderr: [{}]'.format(code, result, error))
      else:
        with open('reports/blackbox-tests/meta/debian.cnb-rates.txt', 'w') as fd:
          fd.write(result)

        result = [item for item in result.split(os.linesep)]
        result = [item.rsplit('/', 1)[-1].strip() for item in result if "/lib/systemd/system/cnb-rates" in item]

        self.units = result

      scratch.remove()
    except Exception as ex:
      failure = ex
    finally:
      temp.close()
      try:
        self.docker.images.remove('bbtest_artifacts-scratch', force=True)
      except:
        pass

    if failure:
      raise failure

  def configure(self, params = None):
    options = dict()
    options.update(UnitHelper.default_config())
    if params:
      options.update(params)

    os.makedirs('/etc/cnb-rates/conf.d', exist_ok=True)
    with open('/etc/cnb-rates/conf.d/init.conf', 'w') as fd:
      fd.write(str(os.linesep).join("CNB_RATES_{!s}={!s}".format(k, v) for (k, v) in options.items()))

  def collect_logs(self):
    (code, result, error) = execute(['journalctl', '-o', 'cat', '--no-pager'])
    if code == 'OK':
      with open('reports/blackbox-tests/logs/journal.log', 'w') as fd:
        fd.write(result)

    for unit in set(self.__get_systemd_units() + self.units):
      (code, result, error) = execute(['journalctl', '-o', 'cat', '-u', unit, '--no-pager'])
      if code != 'OK' or not result:
        continue
      with open('reports/blackbox-tests/logs/{}.log'.format(unit), 'w') as fd:
        fd.write(result)

  def teardown(self):
    self.collect_logs()
    for unit in self.__get_systemd_units():
      execute(['systemctl', 'stop', unit])
    self.collect_logs()

  def __get_systemd_units(self):
    (code, result, error) = execute(['systemctl', 'list-units', '--all', '--no-legend'])
    result = [item.replace('*', '').strip().split(' ')[0].strip() for item in result.split(os.linesep)]
    result = [item for item in result if "cnb-rates" in item]
    return result
