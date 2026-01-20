# -*- coding:utf-8 -*-

""" Http协议客户端 """

import urllib
from httplib2 import Http
from components.json_wrapper import json_encode
from . import get_host

class HttpClient(object):
    def __init__(self, host=get_host(), port=''):
        self.host = host
        self.port = port
        self.client = Http(disable_ssl_certificate_validation=True)

    def _url(self, version='v1', interface='', callback=''):
        """ 构建请求URL
        arguments:
            version -- 接口版本
            interface -- 接口名
            callback -- 方法名

        return:
            {str} 拼接后的url
        """
        base = 'https://{host}:{port}/{version}/{interface}'.format( host=self.host,
                                                                     port=self.port,
                                                                     version=version,
                                                                     interface=interface )
        query = urllib.urlencode({'method': callback})

        return '?'.join((base, query))

    def request(self, method='POST', interface='', callback='', data=''):
        """ 发送请求 

        arguments: 
            {str} method -- 请求方法
            {str} interface -- 接口名
            {str} callback -- 方法名
            {json|dict} data -- body参数，可以传入dict对象，最终会转换成json字符串

        return:
            {tuple} 返回包含响应头和响应体的元组
        """

        url = self._url(interface=interface, callback=callback)
        if isinstance(data, dict):
            data = json_encode(data)

        headers, body = self.client.request(url, method, data)
        if int(headers['status']) >= 400:
            raise Exception(body)
        return (headers, body)
