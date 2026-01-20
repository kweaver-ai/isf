# --*-- coding: utf-8 --*--

""" privateAPI代理转发 """ 
import time
import json
import uuid
from components.requests_handler import requests
from django.http import HttpResponse,HttpResponseForbidden
from components.getconfigfromfile import getconfigfromfile
from components.ipaddress_util import get_ipaddress
from components.common import verify

def controller(request, service, module, method = None):
    # 根据request.path拼接url
    pathArr = request.get_full_path().split('/')
    length = len(pathArr)
    newPathArr = pathArr[3:length]
    pathUrl = '/'.join(newPathArr)

    info = getconfigfromfile()

    svcInfo = {
        'ossgateway': {
            'host': info['ossgateway']['privateHttpHost'],
            'port': info['ossgateway']['privateHttpPort']
        },
        'audit-log': {
            'host': info['audit-log']['privateHttpHost'],
            'port': info['audit-log']['privateHttpPort']
        }
    }
        
    requrl = 'http://{host}:{port}/api/{service}/v1/{pathUrl}'.format(host=svcInfo[service]['host'],
                                                                     port=svcInfo[service]['port'],
                                                                     service=service,
                                                                     pathUrl=pathUrl)
    
    payload = request.body
    
    # 下载对象存储文件，设置有效期1个小时
    if service == 'ossgateway' and 'type=query_string' in pathUrl:
        expires_time = int(time.time()) + 3600
        requrl += f'&Expires={expires_time}'
    # 记录管理日志，设置ip、date和out_biz_id字段
    elif service == 'audit-log':
        payload = json.loads(request.body.decode('utf-8'))
        payload['ip'] = get_ipaddress(request)
        payload['date'] = int(time.time()) * 1000 * 1000
        payload['out_biz_id'] = str(uuid.uuid4()).replace('-', '')
        payload = json.dumps(payload)

    res = requests.request(method=request.method, url=requrl, data=payload)

    if int(res.status_code) >= 400:
        if not verify(request):
            return HttpResponseForbidden()
        else:
            raise Exception(res.text)

    return HttpResponse(res.text, content_type='application/json')

