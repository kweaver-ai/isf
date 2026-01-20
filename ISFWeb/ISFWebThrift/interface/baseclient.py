# coding:utf-8
'''
Created on 2013-7-10

@author: Chensi Yuan(yuan.chensi@eisoo)
'''
import time
from thrift.transport import TSocket
from thrift.transport import TTransport
from thrift.protocol import TBinaryProtocol
from . import get_host

class BaseClient(object):
    host = ''
    port = 0
    transport = None
    client = None

    def __init__(self):
        """
        Construct
        """
        
        if not self.host:
            self.host = get_host()
        self.transport = TSocket.TSocket(self.host, self.port)
        self.transport = TTransport.TBufferedTransport(self.transport)
        protocol = TBinaryProtocol.TBinaryProtocol(self.transport)
        self.client = self.client.Client(protocol)

        i = 0
        while i < 4:
            try:
               self.transport.open()
               break
            except Exception as e:
                if i > 0:
                    print("retry connect to:", self.host, time.strftime("%y-%m-%d %H:%M:%S"), "retry time:", i)
                time.sleep(0.2)
                i += 1

    def __del__ (self, exception_type=None, exception_val=None, trace=None):
        """
        Del class instance
        """
        self.close()

    def close(self):
        """
        关闭连接
        """
        try:
            self.transport.close()
        except Exception:
            pass

    def call_interface(self, interface, *args):
        """
        调用指定接口，并获得返回值
        """
        method = getattr(self.client, interface)
        if not method:
            raise Exception("Method %s not implemented" % interface)
        result = method(*args)
        return result
