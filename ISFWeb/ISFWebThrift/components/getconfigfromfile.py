# coding: utf-8
import configparser
from components.requests_handler import requests
import json

def format_host(host):
    ''' 处理ipv6的url '''

    if ':' in host:
        # ipv6

        if '[' in host:
            # 已经加过[]了，无需处理，直接返回
            return host

        # 加上[]
        return '[{host}]'.format(host=host)

    # ipv4，直接返回
    return host

def getconfigfromfile():
    """ 获取client_id等配置信息 """

    config = configparser.ConfigParser()
    config.read('/config/service_access.conf')
    hydra_public_host = config.get('hydra','publicHost')
    hydra_public_port = config.get('hydra','publicPort')
    hydra_admin_host = config.get('hydra','administrativeHost')
    hydra_admin_port = config.get('hydra','administrativePort')
    sharemgnt_host = config.get('sharemgnt', 'host')
    sharemgnt_port = config.get('sharemgnt', 'port')
    sharemgnt_single_host = config.get('sharemgnt-single', 'host')
    sharemgnt_single_port = config.get('sharemgnt-single', 'port')
    eacp_publicHttpHost = config.get('eacp', 'publicHttpHost')
    eacp_publicHttpPort = config.get('eacp', 'publicHttpPort')
    eacp_privateHttpHost = config.get('eacp', 'privateHttpHost')
    eacp_privateHttpPort = config.get('eacp', 'privateHttpPort')
    eacp_thriftHost = config.get('eacp', 'thriftHost')
    eacp_thriftPort = config.get('eacp', 'thriftPort')
    ossgateway_privateHttpHost = config.get('ossgatewaymanager', 'privateHost')
    ossgateway_privateHttpPort = config.get('ossgatewaymanager', 'privatePort')
    auditlog_privateHttpHost = config.get('audit-log', 'privateHost')
    auditlog_privateHttpPort = config.get('audit-log', 'privatePort')

    return {
        'ossgateway': {
            'privateHttpHost': ossgateway_privateHttpHost,
            'privateHttpPort': ossgateway_privateHttpPort,
        },
        'hydra': {
            'public_host': hydra_public_host,
            'public_port': hydra_public_port,
            'admin_host': hydra_admin_host,
            'admin_port': hydra_admin_port
        },
        'sharemgnt': {
            'host': sharemgnt_host,
            'port': sharemgnt_port,
        },
        'sharemgnt-single': {
            'host': sharemgnt_single_host,
            'port': sharemgnt_single_port,
        },
        'eacp': {
            'publicHttpHost': eacp_publicHttpHost,
            'publicHttpPort': eacp_publicHttpPort,
            'privateHttpHost': eacp_privateHttpHost,
            'privateHttpPort': eacp_privateHttpPort,
            'thriftHost': eacp_thriftHost,
            'thriftPort': eacp_thriftPort,
        },
        'audit-log': {
            'privateHttpHost': auditlog_privateHttpHost,
            'privateHttpPort': auditlog_privateHttpPort,
        },
    }
    
def getHost():
    info = getconfigfromfile()
    requrl = 'http://{host}:9703/api/deploy-manager/v1/access-addr/app'.format(host=info['deploy']['host'])

    res = requests.request('GET', url=requrl)

    if int(res.status_code) >= 400:
        raise Exception(res.text)

    data = json.loads(res.text)
    data['host'] = format_host(data['host'])
    # 如果path是/，则prefix为空，否则prefix为path
    if data['path'] == '/':
        data['prefix'] = ''
    else:
        data['prefix'] = data['path']
    data['host_port_prefix'] = '{host}:{port}{prefix}'.format(host=data['host'],port=data['port'],prefix=data['prefix'])

    return data