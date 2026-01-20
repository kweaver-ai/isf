import configparser
from interface.sharemgnt import ShareMgntClient
from components.requests_handler import requests
import json
from components.json_wrapper import json_encode, json_decode
from ShareMgnt.constants import NCT_SYSTEM_ROLE_SUPPER, NCT_SYSTEM_ROLE_ADMIN, NCT_SYSTEM_ROLE_SECURIT, NCT_SYSTEM_ROLE_AUDIT, NCT_SYSTEM_ROLE_ORG_MANAGER, NCT_SYSTEM_ROLE_ORG_AUDIT

def get_info_dict(name=''):
    try:
        _info_dict = get_client_info()
        if name:
            return _info_dict[name]
        else:
            return _info_dict
    except KeyError:
        return {}

def get_client_info():
    """ 获取client_id等配置信息 """
    config = configparser.ConfigParser()
    config.read('/config/service_access.conf')
    hydra_public_host = config.get('hydra', 'publicHost')
    hydra_public_port = config.get('hydra', 'publicPort')
    hydra_admin_host = config.get('hydra', 'administrativeHost')
    hydra_admin_port = config.get('hydra', 'administrativePort')
    authentication_publicHost = config.get('authentication', 'publicHost')
    authentication_publicPort = config.get('authentication', 'publicPort')
    authentication_privateHost = config.get('authentication', 'privateHost')
    authentication_privatePort = config.get('authentication', 'privatePort')
    user_mgmt_publicHost = config.get('user-management', 'publicHost')
    user_mgmt_publicPort =  config.get('user-management', 'publicPort')
    user_mgmt_privateHost = config.get('user-management', 'privateHost')
    user_mgmt_privatePort =  config.get('user-management', 'privatePort')
    return {
        'hydra': {
            'public_host': hydra_public_host,
            'public_port': hydra_public_port,
            'admin_host': hydra_admin_host,
            'admin_port': hydra_admin_port
        },
        'authentication': {
            'publicHost': authentication_publicHost,
            'publicPort': authentication_publicPort,
            'privateHost': authentication_privateHost,
            'privatePort': authentication_privatePort,
        },
        'user-mgmt': {
            'publicHost': user_mgmt_publicHost,
            'publicPort': user_mgmt_publicPort,
            'privateHost': user_mgmt_privateHost,
            'privatePort': user_mgmt_privatePort
        }
    }

def get_user_mgnt_info_by_userid(userid, token):
    user_mgmt = get_info_dict()['user-mgmt']
    fields_str = ','.join(['custom_attr'])
    url='http://{host}:{port}/api/user-management/v1/users/{user_ids}/{fields}'.format(host=user_mgmt['privateHost'], port=user_mgmt['privatePort'], user_ids=userid, fields=fields_str)
    payload = 'token={token}'.format(token=token)
    headers = {
        'content-type': 'application/x-www-form-urlencoded',
        'cache-control': 'no-cache',
    }
    usermgent_response = requests.request('GET', url, data=payload, headers=headers)
    usermgent_info = json.loads(usermgent_response.text)
    return usermgent_info

def get_userid_by_token(token):
    hydra = get_info_dict()['hydra']
    """ 通过token获取userid """
    url = 'http://{host}:{port}/admin/oauth2/introspect'.format(host=hydra['admin_host'], port=hydra['admin_port'])
    payload = 'token={token}'.format(token=token)
    headers = {
        'content-type': 'application/x-www-form-urlencoded',
        'cache-control': 'no-cache',
    }
    response = requests.request('POST', url, data=payload, headers=headers)
    return response

def get_user_info_by_userid(userid, token):
    """ 通过userid获取userinfo """
    sharemgnt_inst = ShareMgntClient()
    sharemgnt_response = sharemgnt_inst.call_interface('Usrm_GetUserInfo', userid)
    sharemgnt_info = json_decode(json_encode(sharemgnt_response))

    user_mgnt_info = get_user_mgnt_info_by_userid(userid, token)
    sharemgnt_info['user']['custom_attr'] = user_mgnt_info[0].get('custom_attr')
    return sharemgnt_info

def get_user_info(token):
    get_userid_res = get_userid_by_token(token)
    if get_userid_res.status_code < 400:
        userid = json.loads(get_userid_res.text).get('sub')
        user = get_user_info_by_userid(userid, token)
        return user

def get_authorization_header(request):
    auth = request.META.get('HTTP_AUTHORIZATION', b'').split()
    csrftoken = request.META.get('HTTP_X_CSRFTOKEN', b'')
    querytoken = request.GET.get('token', '')

    if not ((auth and auth[0].lower() == 'bearer') or csrftoken or querytoken):
        return None
    try:
        return csrftoken or querytoken or auth[1]
    except:
        return None

def is_console_role(user):
    CONSOLE_ROLES = [
        NCT_SYSTEM_ROLE_SUPPER, # 超级管理员
        NCT_SYSTEM_ROLE_ADMIN, # 系统管理员
        NCT_SYSTEM_ROLE_SECURIT, # 安全管理员
        NCT_SYSTEM_ROLE_AUDIT, # 审计管理员
        NCT_SYSTEM_ROLE_ORG_MANAGER, # 组织管理员
        NCT_SYSTEM_ROLE_ORG_AUDIT  # 组织审计员
    ]

    roles = user['user']['roles']

    result = [(role['id'] in CONSOLE_ROLES) for role in roles]
    return any(result)

# 特殊的API限制
def specialAPICheck(request, userInfo):
    whiteList = [
        NCT_SYSTEM_ROLE_SUPPER, # 超级管理员
        NCT_SYSTEM_ROLE_ADMIN, # 系统管理员
        NCT_SYSTEM_ROLE_SECURIT, # 安全管理员
        NCT_SYSTEM_ROLE_AUDIT # 审计管理员
    ]
    roles = userInfo['user']['roles']
    
    # 组织管理、组织审计员不允许调该接口
    if request.path.endswith('Usrm_GetAllUsers'):
         return any([(role['id'] in whiteList) for role in roles])
    # 组织审计员只允许查看自己的userInfo
    elif request.path.endswith('Usrm_GetUserInfo'):
        whiteList.append(NCT_SYSTEM_ROLE_ORG_MANAGER)
        if not any([(role['id'] in whiteList) for role in roles]):
            data = json.loads(request.body)
            return data[0] == userInfo['id']
    return True

def verify(request):
    try:
        access_token = get_authorization_header(request)
        user_info = get_user_info(access_token)
        if is_console_role(user_info):
            return specialAPICheck(request, user_info)
        return False
    except:
        return False
