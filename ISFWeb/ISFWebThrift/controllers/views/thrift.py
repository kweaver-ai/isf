# -*- coding: utf-8 -*-

from components.json_wrapper import json_response
from components.thrift_agent import ThriftAgent
from components.getconfigfromfile import getconfigfromfile

@json_response
def controller(request, module, method):
    request_method = ['GET', 'POST','PUT', 'DELETE', 'HEAD']
    if request.method in request_method:
        client = ThriftAgent(module)
        
        return client.call(method, request.body)
    else:
        raise Exception('HTTP request method error')

def getAddr(module):
    info = getconfigfromfile()
    hosts = {
       'ShareMgnt': info['sharemgnt']['host'],
       'ShareMgntSingle': info['sharemgnt-single']['host'],
       'EACP': info['eacp']['thriftHost'],
    }

    return hosts.get(module)
