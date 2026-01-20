# coding:utf-8

from django.http import HttpResponse
from components.json_wrapper import json_encode
import json
from EACP.ttypes import *
from EInfoworksLogger.ttypes import *
from EThriftException.ttypes import *
from ShareMgnt.ttypes import *
from interface.eacp import EACPClient
from interface.sharemgnt import ShareMgntClient
from interface.sharemgntsingle import ShareMgntSingleClient

'''
Thrift协议代理
'''
class ThriftAgent(object):

    def __init__(self, protocol):
        '''
        初始化协议代理

        Arguments:
            protocol {string} - 协议名
        '''

        self.protocol = protocol


    def call(self, method_name, json_params=''):
        '''
        调用协议方法

        Arguments:
            method_name {string} - 方法名
            json_params {string} - JSON参数
        '''

        clients = {
            'EACP': EACPClient,
            'ShareMgnt': ShareMgntClient,
            'ShareMgntSingle': ShareMgntSingleClient,
        }

        client = clients.get(self.protocol, '')

        if client:
            client_inst = client()
            return client_inst.call_interface(method_name, *self._convert_json_to_list(json_params))
        else:
            resp = HttpResponse(json_encode(None), content_type='application/json')
            resp.status_code = 501
            return resp


    @staticmethod
    def _convert_dict_to_struct(arg):
        '''
        转换dict为struct

        Arguments:
            arg {dict} - 要转换的dict

        Returns:
            struct 返回dict对应的结构体
        '''

        struct_name, struct_body = list(arg.items())[0]

        struct = globals().get(struct_name)

        # for key, value in struct_body.items():
        #     struct_body[key] = value.encode('utf-8') if isinstance(value, string_types) else value

        return struct(**struct_body)


    def _convert_json_to_list(self, json_params=''):
        '''
        转换参数

        Arguments:
            json_params {string} - JSON参数

        Returns:
            list 返回参数列表
        '''

        args = json.loads(json_params) if json_params else []

        return self.build_arguments(args)

    def build_arguments(self, args):
        ''' 
        构建参数列表
        '''

        ret = []

        for arg in args:
            if isinstance(arg, dict):
                if set(globals().keys()).intersection(arg.keys()):
                    ret.append(self._build_struct(arg))
                else:
                    ret.append(arg)
            elif isinstance(arg, list):
                ret.append(self.build_arguments(arg))
            else:
                # 不明白之前为啥要 encode('utf-8')，所以注释留着
                # ret.append(arg.encode('utf-8') if isinstance(arg, string_types) else arg)
                ret.append(arg)

        return ret


    def _build_struct(self, arg):
        '''
        构造结构体

        Arguments:
            arg {dict} - 要转换的dict

        Returns:
            {struct} 返回结构体对象 

        '''

        struct_body = list(arg.values())[0]

        for key, value in struct_body.items():
            if isinstance(value, dict):
                struct_body[key] = self._build_struct(value)
            elif isinstance(value, list):
                collection = []
                
                for v in value:
                    if isinstance(v, dict) or isinstance(v, list):
                        collection.append(self._build_struct(v))
                    else:
                        # collection.append(v.encode('utf-8') if isinstance(v, string_types) else v)
                        collection.append(v)

                struct_body[key] = collection

        return self._convert_dict_to_struct(arg)