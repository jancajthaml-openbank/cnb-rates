#!/usr/bin/env python3
# -*- coding: utf-8 -*-

import threading
from http.server import HTTPServer
import ssl
import os
import time
import tempfile
from .handler import RequestHandler
from .logic import BussinessLogic


class CNBMock(threading.Thread):

  def __init__(self, context):
    threading.Thread.__init__(self)
    self.context = context
    self.port = 4000
    self.__keyfile = tempfile.NamedTemporaryFile()
    self.__certfile = tempfile.NamedTemporaryFile()

    # https://stackoverflow.com/questions/10175812/how-to-create-a-self-signed-certificate-with-openssl/27931596#27931596
    # https://stackoverflow.com/questions/21297139/how-do-you-sign-a-certificate-signing-request-with-your-certification-authority/21340898#21340898

    os.system("sed -i '/^.*v3_ca.*/a subjectAltName = IP:127.0.0.1' /etc/ssl/openssl.cnf")
    os.system('faketime "1991-01-01 00:00:00" openssl req -x509 -nodes -newkey rsa:2048 -keyout "{}" -out "{}" -days 36500 -subj "/CN=127.0.0.1" > /dev/null 2>&1'.format(self.__keyfile.name, self.__certfile.name))
    os.system('cp {} /usr/local/share/ca-certificates/ > /dev/null 2>&1'.format(self.__certfile.name))
    os.system('update-ca-certificates > /dev/null 2>&1')

  def start(self):
    self.httpd = HTTPServer(('127.0.0.1', self.port), RequestHandler)
    self.httpd.socket = ssl.wrap_socket(self.httpd.socket, certfile=self.__certfile.name, keyfile=self.__keyfile.name, server_side=True)
    self.httpd.logic = BussinessLogic()
    threading.Thread.start(self)

  def run(self):
    self.httpd.serve_forever()

  def stop(self):
    if self.httpd:
      self.httpd.shutdown()
    try:
      self.join()
    except:
      pass
    self.__keyfile.close()
    self.__certfile.close()
    del self.httpd
