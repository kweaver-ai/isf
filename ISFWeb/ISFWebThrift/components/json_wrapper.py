# coding:utf-8
'''
Created on 2013-7-11
重写python内置json模块
@author: Chensi Yuan(yuan.chensi@eisoo)
'''
from json import JSONEncoder, JSONDecoder, dumps
from django.http import HttpResponse

class Encoder(JSONEncoder):
    '''
    重载JSONEncoder，让其可以解析class
    '''
    def default (self, obj):
        return obj.__dict__

def json_encode(obj):
    '''
    针对json库不方便的问题，进行的重写
    '''
    return Encoder(ensure_ascii=False).encode(obj)

def json_decode(json):
    ''' 
    封装JSON的decode方法 

    @param json String json格式数据
    @return Dict 返回字典

    '''

    return JSONDecoder().decode(json)

def json_response(func):
    '''
    自动封装返回值为json格式
    '''
    def warp(*args, **kwargs):
        '''
        封装函数
        '''
        result = func (*args, **kwargs)
        if isinstance(result, HttpResponse):
            result["Content-Type"] = 'application/json'
            return result
        return HttpResponse(json_encode(result),
                            content_type='application/json')

    return warp

def json_dumps(d):
    '''
    将Dict对象转为JSON String
    '''

    return dumps(d, ensure_ascii=True)