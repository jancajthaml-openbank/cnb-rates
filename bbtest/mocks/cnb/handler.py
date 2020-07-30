#!/usr/bin/env python3
# -*- coding: utf-8 -*-

from http.server import BaseHTTPRequestHandler
import json
import datetime


class RequestHandler(BaseHTTPRequestHandler):

  def log_message(self, format, *args):
    pass

  def do_GET(self):
    if self.path.startswith('/en/financial_markets/foreign_exchange_market/exchange_rate_fixing/daily.txt'):

      options = dict()
      for pair in self.path.split('/')[-1].split('?')[-1].split('&'):
        (k, v) = pair.split('=')
        options[k] = v

      if 'date' in options:
        try:
          options['date'] = datetime.datetime.strptime(options['date'], '%d+%m+%Y').replace(tzinfo=datetime.timezone.utc)
        except:
          return self.__respond(400)
      else:
        options['date'] = datetime.datetime.utcnow().replace(hour=0, minute=0, second=0, microsecond=0, tzinfo=datetime.timezone.utc)

      response = self.server.logic.get_daily_main_fx(options['date'])

      return self.__respond(200, response)

    if self.path.startswith('/en/financial_markets/foreign_exchange_market/other_currencies_fx_rates/fx_rates.txt'):
      for _ in range(20):
        print('CNB get other fx {}'.format(self.path))
      return self.__respond(200)

    return self.__respond(404)

  def __respond(self, status, body=None):
    self.send_response(status)
    self.end_headers()
    if body:
      self.wfile.write(body.encode('utf-8'))
