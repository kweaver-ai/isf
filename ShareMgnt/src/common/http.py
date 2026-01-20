# coding: utf-8
"""
HTTP请求库
"""
import urllib.request, urllib.parse, urllib.error
import json
import configparser
from src.common.lib import (raise_exception, check_url)
from src.common.sharemgnt_logger import ShareMgnt_Log
import ssl
import requests
import requests.packages.urllib3
from mq_sdk.proton_mq import Connector
import asyncio


def send_request(url, data=None, method="GET", headers=None):
    """
    发送请求
    Args:
        url: string 请求URL
        data: dict 请求数据
        method: string 请求类型，例如 POST | GET
        headers: dict http头信息
    Return:
        response: string 请求结果
    """
    ssl.match_hostname = lambda cert, hostname: True
    if method == "GET":
        if data:
            url = "%s?%s" % (url, urllib.parse.urlencode(data))
        request = urllib.request.Request(url)

    if method == "POST":
        if isinstance(data, dict):
            data = json.dumps(data, ensure_ascii=False)
        request = urllib.request.Request(url, data.encode('utf-8'))

    if method == "PUT":
        if isinstance(data, dict):
            data = json.dumps(data, ensure_ascii=False)
        request = urllib.request.Request(url, data)
        request.get_method = lambda: 'PUT'

    if not request:
        raise_exception("Invalide method type")

    # 添加必要的请求头
    if headers:
        for k, v in list(headers.items()):
            request.add_header(k, v)

    try:
        response = urllib.request.urlopen(request, timeout=60)
    except urllib.error.HTTPError as e:
        ShareMgnt_Log("test error code: %s, content: %s", e.code, e.reason)
        return e.code, "{\"errorMsg\":\"" + e.reason + "\"}"

    content = response.read()
    code = response.getcode()
    response.close()
    if code != 200 and code != 204:
        ShareMgnt_Log("send_request error, url: %s, code: %s, content: %s", url, code, content)

    return code, content

def test_connection(url, method="POST"):
    """
    测试 url 是否可用
    """
    requests.packages.urllib3.disable_warnings()

    # url 合法性检测
    check_url(url)

    # url 可访问性检测
    try:
        requests.request(method, url, verify=False, timeout=3)
        return True
    except Exception as e:
        ShareMgnt_Log(str(e))
        return False


cnt = Connector.get_connector_from_file("/sysvol/conf/service_conf/mq_config.yaml")
def pub_nsq_msg(topic, content):
    """
    发布nsq消息
    """

    try:
        message = json.dumps(content, ensure_ascii=False)
        asyncio.run(cnt.pub(topic, message))
    except Exception as err:
        raise Exception(("Pub mq message failed, topic = {0}, err = {1}").format(topic, err))
