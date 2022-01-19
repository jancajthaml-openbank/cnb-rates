#!/usr/bin/env python3
# -*- coding: utf-8 -*-

from openbank_testkit import Shell
import datetime


class TimeshiftHelper(object):

  def __init__(self, context):
    self.__original_epoch = datetime.datetime.now()

  def set_date_time(self, new_epoch):
    (code, result, error) = Shell.run(['timedatectl', 'set-time', new_epoch.strftime('%Y-%m-%d %H:%M:%S')])
    assert code == 'OK', "timedatectl set-time failed with: {} {}".format(result, error)

  def bootstrap(self):
    Shell.run(['timedatectl', 'set-ntp', '0'])
    Shell.run(['timedatectl', 'set-local-rtc', '0'])

  def teardown(self):
    Shell.run(['timedatectl', 'set-local-rtc', '1'])
    Shell.run(['timedatectl', 'set-ntp', '1'])
    Shell.run(['systemctl', 'restart', 'systemd-timedated'])
    Shell.run(['systemctl', 'restart', 'systemd-timesyncd'])
    Shell.run(['timedatectl', 'set-time', self.__original_epoch.strftime('%Y-%m-%d %H:%M:%S')])
