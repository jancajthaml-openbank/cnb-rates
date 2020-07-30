#!/usr/bin/env python3
# -*- coding: utf-8 -*-


from helpers.shell import execute
import datetime


class TimeshiftHelper(object):

  def __init__(self, context):
    self.__original_epoch = datetime.datetime.now()

  def set_date_time(self, new_epoch):
    (code, result, error) = execute(['timedatectl', 'set-time', new_epoch.strftime('%Y-%m-%d %H:%M:%S')])
    assert code == 0, "timedatectl set-time failed with: {} {}".format(result, error)

  def bootstrap(self):
    execute(['timedatectl', 'set-ntp', '0'])
    execute(['timedatectl', 'set-local-rtc', '0'])

  def teardown(self):
    execute(['timedatectl', 'set-local-rtc', '1'])
    execute(['timedatectl', 'set-ntp', '1'])
    execute(['systemctl', 'restart', 'systemd-timedated'])
    execute(['systemctl', 'restart', 'systemd-timesyncd'])
    execute(['timedatectl', 'set-time', self.__original_epoch.strftime('%Y-%m-%d %H:%M:%S')])
