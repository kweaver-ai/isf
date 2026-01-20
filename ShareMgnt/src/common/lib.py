#!/usr/bin/python3
# -*- coding:utf-8 -*-
"""
sharemgnt 公共函数
"""
import collections
import configparser
import hashlib
import hmac
import base64
import os
import re
import socket
import subprocess
import sys
import time
import traceback
import uuid
from contextlib import contextmanager
from hashlib import md5
from socket import inet_aton
from struct import unpack
from datetime import datetime
import random
import string
import MySQLdb

import ldap
import base64
from eisoo import nodeconf
from eisoo.tclients import TClient
from EThriftException.ttypes import ncTException, ncTExpType
from ShareMgnt.ttypes import (ncTAlarmConfig, ncTShareMgntError,
                              ncTSmtpSrvConf, ncTUsrmAddUserInfo,
                              ncTUsrmDomainInfo, ncTUsrmGetUserInfo,
                              ncTUsrmImportContent, ncTUsrmImportOption)
from src.common.sharemgnt_logger import ShareMgnt_Log
from src.common.business_date import BusinessDate
from thrift import Thrift
from thrift.protocol import TBinaryProtocol
from thrift.Thrift import TException
from thrift.transport import TSocket, TTransport
from src.common.db.connector import DBConnector

ERROR_INFO_BLACK_TUPLE = (MySQLdb.Error,)


def check_args(func):
    """
    check the parameters in sharemgnt_mgnt_server.py
    """
    def _deco(*args, **kwargs):
        """
        check args
        """
        for i in range(len(args)):
            if i == 0:
                continue
            if args[i] is None:
                raise_my_exception(ncTExpType.NCT_INFO,
                                   traceback.extract_stack()[-2][1],
                                   traceback.extract_stack()[-2][0],
                                   _("INVALID_PARAM") % (args[i]))
            if isinstance(args[i], ncTUsrmAddUserInfo):
                if args[i].user.loginName is None or \
                   args[i].password is None or \
                   args[i].user.departmentIds is None:
                    raise_my_exception(ncTExpType.NCT_INFO,
                                       get_line_number(),
                                       __file__,
                                       _("INVALID_PARAM") % (args[i]))

                if len(args[i].user.loginName) == 0 or \
                   len(args[i].password) == 0 or \
                   len(args[i].user.departmentIds) == 0 or \
                   len(args[i].user.departmentIds) > 1:
                    raise_my_exception(ncTExpType.NCT_INFO,
                                       get_line_number(),
                                       __file__,
                                       _("INVALID_PARAM") % (args[i]))

            if isinstance(args[i], ncTUsrmGetUserInfo):
                if args[i].id is None or args[i].user.loginName is None:
                    raise_my_exception(ncTExpType.NCT_INFO,
                                       get_line_number(),
                                       __file__,
                                       _("INVALID_PARAM") % (args[i]))

                if len(args[i].id) == 0 or len(args[i].user.loginName) == 0:
                    raise_my_exception(ncTExpType.NCT_INFO,
                                       get_line_number(),
                                       __file__,
                                       _("INVALID_PARAM") % (args[i]))

            if isinstance(args[i], ncTUsrmDomainInfo):
                if args[i].id is None or \
                   args[i].type is None or \
                   args[i].parentId is None or \
                   args[i].name is None or \
                   args[i].ipAddress is None or \
                   args[i].adminName is None or \
                   args[i].password is None or \
                   args[i].status is None:
                    raise_my_exception(ncTExpType.NCT_INFO,
                                       get_line_number(),
                                       __file__,
                                       _("INVALID_PARAM") % (args[i]))

            if isinstance(args[i], ncTUsrmImportContent):
                if args[i].domain.name is None or \
                   args[i].domain.ipAddress is None or \
                   args[i].domain.adminName is None or \
                   args[i].domain.password is None:
                    raise_my_exception(ncTExpType.NCT_INFO,
                                       get_line_number(),
                                       __file__,
                                       _("INVALID_PARAM") % (args[i]))

                if args[i].domainName and \
                   (args[i].users or args[i].ous):
                    raise_my_exception(ncTExpType.NCT_INFO,
                                       get_line_number(),
                                       __file__,
                                       _("INVALID_PARAM") % (args[i]))

            if isinstance(args[i], ncTUsrmImportOption):
                if args[i].userEmail is None or \
                   args[i].userDisplayName is None or \
                   args[i].userCover is None or \
                   args[i].departmentId is None:
                    raise_my_exception(ncTExpType.NCT_INFO,
                                       get_line_number(),
                                       __file__,
                                       _("INVALID_PARAM") % (args[i]))


            if isinstance(args[i], ncTSmtpSrvConf):
                if args[i].server is None or \
                   args[i].safeMode is None or \
                   args[i].port is None or \
                   args[i].email is None:
                    raise_my_exception(ncTExpType.NCT_INFO,
                                       get_line_number(),
                                       __file__,
                                       _("INVALID_PARAM") % (args[i]))


            if isinstance(args[i], ncTAlarmConfig):
                if args[i].infoConfig is None or \
                   args[i].infoConfig < 0 or \
                   args[i].warnConfig is None or \
                   args[i].warnConfig < 0 or \
                   args[i].emailToList is None:
                    raise_my_exception(ncTExpType.NCT_INFO,
                                       get_line_number(),
                                       __file__,
                                       _("INVALID_PARAM") % (args[i]))
                for e in args[i].emailToList:
                    if not check_email(e):
                        raise_my_exception(ncTExpType.NCT_INFO,
                                           get_line_number(),
                                           __file__,
                                           _("INVALID_PARAM") % (args[i]))

        ret = func(*args, **kwargs)
        return ret
    return _deco


def check_str_available(value):
    """
    检查字符串是否包含不合法的字符
    """
    dirty_stuff = ["\"", "\\", "/", "*", "'", "=", "-",
                   "#", ";", "<", ">", "+", "%"]
    for stuff in dirty_stuff:
        if stuff in value:
            return False
    return True


def check_is_uuid(value):
    """
    检查是否为合法UUID，用于优化ID存在判断
    UUID规则，数字为多少个16进制字符，0-9 a-f：
    8-4-4-4-12
    """
    if not isinstance(value, str):
        return False

    uuid_sliced = value.split("-")
    if len(uuid_sliced) != 5:
        return False

    for i, hexlen in enumerate([8, 4, 4, 4, 12]):
        if len(uuid_sliced[i]) != hexlen:
            return False
        try:
            int(uuid_sliced[i], 16)
        except ValueError:
            return False

    return True


def check_name(name):
    """
    检测名称是否合法
    检测规则：不允许使用特殊字符 / : * ? " < > | ，且长度范围为 1~128 个字符
    """
    if name is None:
        return False
    if not isinstance(name, str):
        name = name.decode('utf8')

    name_len = len(name)
    if name_len < 1 or name_len > 128:
        return False

    # 不合法字符列表
    # char_list = ['/', ':', '?', '\"', '<', '>', '|',
    #             '：', '？', '’', '“', '”', '《', '》']
    char_list = ['\\', '/', ':', '*', '?', '\"', '<', '>', '|']
    for char in char_list:
        if char in name:
            return False
    return True

def check_name2(name):
    """
    检测名称是否合法
    检测规则：不允许使用特殊字符 / * ? " < > | ，且长度范围为 1~128 个字符
    中广核去除:字符检查
    """
    if name is None:
        return False
    if not isinstance(name, str):
        name = name.decode('utf8')

    name_len = len(name)
    if name_len < 1 or name_len > 128:
        return False

    # 不合法字符列表
    char_list = ['\\', '/', '*', '?', '\"', '<', '>', '|']
    for char in char_list:
        if char in name:
            return False
    return True

def check_email(email):
    """
    检测邮箱是否合法
    """
    if len(email) < 5 or len(email) > 100:
        return False
    email_regex = r"^[a-zA-Z0-9_\.\-]+@[a-zA-Z0-9\-_]+(\.[a-zA-Z0-9\-_]+)+$"
    reobj = re.compile(email_regex)
    return reobj.search(email) is not None


def check_server(server):
    """
    检测邮件服务器是否合法
    """
    if len(server) < 3 or len(server) > 100:
        return False
    server_regex = r"^[a-zA-Z0-9@\-_\.]+$"
    reobj = re.compile(server_regex)
    return reobj.search(server) is not None

def check_smtp_params(conf):
    """
    检测邮件服务器参数是否合法
    """
    # 验证配置是否为空
    if conf is None:
        raise_exception(exp_msg=_("stmp not set"),
                        exp_num=ncTShareMgntError.NCT_SMTP_SERVER_NOT_SET)

    # 验证port参数格式正确性
    if not (conf.port > 0
            and conf.port < 65536):
        raise_exception(exp_msg=_("port illegal"),
                        exp_num=ncTShareMgntError.
                        NCT_INVALID_PORT)

    # 验证safeMode参数格式正确性
    if not (conf.safeMode > -1
            and conf.safeMode < 3):
        raise_exception(exp_msg=_("safeMode illegal"),
                        exp_num=ncTShareMgntError.
                        NCT_INVALID_SAFEMODE)

    # 验证server邮件服务器格式正确性
    if not check_server(conf.server):
        raise_exception(exp_msg=_("server illegal"),
                        exp_num=ncTShareMgntError.
                        NCT_INVALID_SERVER)

    # 验证email是电子邮件格式正确性
    if not check_email(conf.email):
        raise_exception(exp_msg=_("email illegal"),
                        exp_num=ncTShareMgntError.
                        NCT_INVALID_EMAIL)

def check_start_limit(start, limit):
    """
    检查分页参数，并返回LIMIT子句
    """
    if start < 0:
        raise_exception(exp_msg=_("IDS_START_LESS_THAN_ZERO"),
                        exp_num=ncTShareMgntError.NCT_START_LESS_THAN_ZERO)

    if limit == -1:
        # 适配ocean base,
        # ocean base获取limit offset 是通过limit+offset获取所有数据
        # 当limit为9223372036854775807时，limit+offfset超过bigint范围报错
        # 后续ocean base 会适配9223372036854775807
        return 'LIMIT {0}, {1}'.format(start, 10000000000)
    elif limit < 0:
        raise_exception(exp_msg=_("IDS_LIMIT_LESS_THAN_MINUS_ONE"),
                        exp_num=ncTShareMgntError.NCT_LIMIT_LESS_THAN_MINUS_ONE)
    else:
        return 'LIMIT {0}, {1}'.format(start, limit)


def escape_key(key):
    """
    去除关键字
    """
    esckey = ""
    for c in key:
        if c == "\n":
            esckey += "\\"
            esckey += "n"
        elif c == "\r":
            esckey += "\\"
            esckey += "r"
        elif c == "\'" or c == "\"" or c == "\\" or c == "%" or c == "_":
            esckey += "\\"
            esckey += c
        else:
            esckey += c
    return esckey


def escape_format_percent(key):
    """
    转换格式化字符中的%
    """
    esckey = []
    for c in key:
        if c == "%":
            esckey.append("%")
            esckey.append(c)
        else:
            esckey.append(c)
    return ''.join(esckey)


def generate_group_str(ids):
    """
    生成组字符串
    """
    str = ""
    for i, group_id in enumerate(ids):
        str += "'"
        str += escape_key(group_id)
        str += "'"
        if i != (len(ids) - 1):
            str += ","
    return str


def has_chinese(value):
    """
    检测值中是否包含中文
    """
    # 如果不是unicode，先转换为unicode
    if not isinstance(value, str):
        try:
            value = value.decode("utf8")
        # UTF8转码失败，尝试gbk
        except UnicodeEncodeError:
            try:
                value = value.decode("gbk")
            # 同样失败，忽略转换
            except UnicodeEncodeError:
                pass

    # 中文正则
    zh_pattern = re.compile('[\u4e00-\u9fa5]+')
    if zh_pattern.search(value):
        return True
    else:
        return False

def check_XSS_safe(value):
    """
    检测值中如果包含 < 和 >，认为有XSS风险
    """
    if '<' in value and '>' in value:
        return False
    return True


def warp_exception(func):
    """
    捕获所有非ncTException类异常
    并转换为ncTException类异常抛出
    """
    def warp(*args, **kwargs):
        """
        warp
        """
        try:
            return func(*args, **kwargs)

        except Exception as ex:
            ShareMgnt_Log(traceback.format_exc())

            if not isinstance(ex, ncTException):
                frame = traceback.extract_stack()
                # 处理黑名单内的错误，防止泄露内部信息
                if isinstance(ex, ERROR_INFO_BLACK_TUPLE):
                    msg = _("IDS_INNER_ERROR")
                else:
                    msg = str(ex)

                raise_exception(msg,
                                frame[-3][0],
                                frame[-3][1],
                                exp_num=99)
            raise ex
    return warp


def raise_exception(exp_msg, file_name="", line_no=0,
                    exp_type=ncTExpType.NCT_WARN, exp_num=0, exp_detail=""):
    """
    丢出异常
    - exp_msg    异常内容
    - file_name 异常文件，默认为执行本函数的上一个堆栈所在文件
    - line_no   异常行号，默认为执行本函数的上一个堆栈所在行
    - exp_type  异常类型，默认警告
    - exp_num     异常编号
    """
    frame = traceback.extract_stack()
    if not file_name:
        file_name = frame[-2][0]
        line_no = frame[-2][1]

    exp = ncTException()
    object.__setattr__(exp, 'expType', exp_type)
    object.__setattr__(exp, 'fileName', file_name)
    object.__setattr__(exp, 'codeLine', line_no)
    object.__setattr__(exp, 'errID', exp_num)
    object.__setattr__(exp, 'expMsg', exp_msg.replace("<", "").replace(">", ""))
    object.__setattr__(exp, 'errProvider', "ShareMgnt")
    object.__setattr__(exp, 'time', BusinessDate.now().strftime("%a %b %d %H:%M:%S %Y"))
    object.__setattr__(exp, 'errDetail', exp_detail)
    raise exp


def raise_my_exception(ex_type, code_line, file_name, message):
    """
    丢出异常
    - extype            类型，ncTExpType.NCT_FATAL, ncTExpType.NCT_CRITICAL,
                            ncTExpType.NCT_WARN = 2, ncTExpType.NCT_INFO
    - codeline          代码行
    - filename          文件
    - message           消息
    """
    ex = ncTException()
    object.__setattr__(ex, 'expType', ex_type)
    object.__setattr__(ex, 'fileName', file_name)
    object.__setattr__(ex, 'codeLine', code_line)
    object.__setattr__(ex, 'errID', 0)
    object.__setattr__(ex, 'expMsg', message)
    object.__setattr__(ex, 'errProvider', "ShareMgnt")
    object.__setattr__(ex, 'time', BusinessDate.now().strftime("%a %b %d %H:%M:%S %Y"))
    raise ex


def get_line_number():
    """
    获取当前文件的行号
    """
    return traceback.extract_stack()[-2][1]


def get_exec_info():
    """
    exc_traceback.tb_lineno
    """
    _, _, exc_traceback = sys.exc_info()
    return exc_traceback


def check_is_active():
    """
    check whether local uuid equal to the active uuid in sharemgnt_db.
    if not, raise. Else, pass
    """
    return


def replace_char(name):
    """
    替换单引号
    """
    if name:
        name = name.replace("'", "''")
    return name


def encrypt_pwd(value):
    """
    对密码明文进行MD5加密，空字符串不加密
    """
    if not value:
        return value
    else:
        return md5(value.encode("utf-8")).hexdigest().lower()


def ntlm_md4(value):
    """
    明文先转成utf-16，再进行MD4 hash
    """
    return hashlib.new("md4", value.encode("utf-16le")).hexdigest()


def sha2_encrypt(data):
    """
    明文sha256进行加密
    """
    return hashlib.sha256(data.encode('utf-8')).hexdigest()


def get_text(elem, key, default=None):
    """
    获取节点文本
    """
    node = elem.find(key)
    if node is None:
        return default

    text = node.text
    if text is None:
        return default
    else:
        return text.strip().encode("utf8")


def get_local_uuid():
    """
    获取本节点的uuid
    """
    return nodeconf.NodeConfig.get_node_uuid()

def check_service_node():
    """
    检查服务节点
    """
    try:
        config = configparser.ConfigParser()
        config.read('/sysvol/conf/service_conf/app_default.conf')
        value = config.getboolean("ShareMgnt", "is_single")
    except Exception:
        raise_exception('read sharemgnt single mark failed')

    if value:
        ShareMgnt_Log("service node running...")
        return True
    else:
        ShareMgnt_Log("non-service node running...")
        return False


def get_app_config(section, option):
    """
    获取app下的配置
    """
    value = ""
    try:
        config = configparser.ConfigParser()
        config.read('/sysvol/conf/service_conf/app_default.conf')
        value = config.get(section, option)
    except Exception:
        raise_exception(
            'read /sysvol/conf/service_conf/app_default.conf failed')

    return value

def get_server_port():
    """
    获取服务端口
    """
    value = ""
    try:
        config = configparser.ConfigParser()
        config.read('/sysvol/conf/service_conf/app_default.conf')
        value = config.getint("ShareMgnt", "port")
    except Exception:
        raise_exception('read sharemgnt server port failed')

    return value

def check_is_md5(url_str):
    """
    检查是否是md5
    """
    if (len(re.findall(r"^([a-fA-F\d]{32,32})$", url_str)) == 0):
        return False
    else:
        return True


def check_is_valid_password(pwd):
    """
    密码可以包含ASCII字符的可见范围符号任意组合 ，且长度[6-100]，需求713330
    """
    if pwd is None or (len(re.findall(r"^([\x20-\x7E]{6,100})$", pwd)) == 0):
        return False
    else:
        return True


def check_is_strong_password(pwd):
    """
    检测是否是强密码:
        最少由管理员控制, 最多100个字符，需同时包含大小写字母及数字， 特殊字符为半角字符可见范围内所有字符，需求713330.
    """
    exec_sql = """
    SELECT f_value FROM t_sharemgnt_config
    WHERE f_key = 'strong_pwd_length'
    """
    results = db_operator(exec_sql)
    min_length = int(results[0]['f_value'])

    if len(pwd) < min_length or len(pwd) > 100:
        return False

    lowcase_res = re.search(r'[a-z]+', pwd)
    capital_res = re.search(r'[A-Z]+', pwd)
    digital_res = re.search(r'[0-9]+', pwd)
    special_res = re.search(r'[\x20-\x2F\x3A-\x40\x5B-\x60\x7B-\x7E]+', pwd)

    if lowcase_res and capital_res and digital_res and special_res:
        if check_is_valid_password(pwd):
            return True

    return False


def exec_command(command, shell=True):
    if isinstance(command, str):
        command = [command]
    proc = subprocess.Popen(command, shell=shell, close_fds=True,
                            stdout=subprocess.PIPE, stderr=subprocess.PIPE)
    (outdata, errdata) = proc.communicate()
    if proc.returncode == 0:
        return (proc.returncode, outdata)
    else:
        return (proc.returncode, errdata)


def is_port_listened(ip, port):
    s = socket.socket(socket.AF_INET, socket.SOCK_STREAM)
    try:
        s.connect((ip, int(port)))
        s.shutdown(2)
        return True
    except:
        return False

def generate_sign(key, data):
    """
    生成签名
    """
    return base64.b64encode(hmac.new(key, data, hashlib.sha1).digest()).rstrip()


def generate_search_order_sql(key):
    """
    生成搜索排序语句
    key: 列表，按关键字排序
    eg:
        key: ['aaa','bbb']
        return: 'case when aaa = %s then 0 when bbb = %s then 1
                 when aaa like %s then 2 when bbb like %s then 3
                 else 4 end '
    """
    if not isinstance(key, list) or (len(key) == 0):
        return

    exact_order_str = """
    when {field} = %s then {value}
    """

    like_order_str = """
     when {field} like %s then {value}
    """

    format_str = ''
    temp_str = ''
    # 获取精确匹配排序语句
    for i in range(len(key)):
        format_str = exact_order_str.format(field=key[i], value=i)
        temp_str += format_str
    exact_order_str = temp_str

    # 获取模糊匹配排序语句
    temp_str = ''
    for i in range(len(key), 2 * len(key)):
        format_str = like_order_str.format(field=key[i - len(key)], value=i)
        temp_str += format_str
    like_order_str = temp_str

    order_str = 'case ' + exact_order_str + like_order_str

    order_str = order_str + "else {0} end ".format(2 * len(key))
    return order_str


def check_net_ip(ip_str):
    """
    """
    pattern = r'^(25[0-5]|2[0-4][0-9]|[01]?[0-9][0-9]?)\.(25[0-5]|2[0-4][0-9]|[01]?[0-9][0-9]?)\.(25[0-5]|2[0-4][0-9]|[01]?[0-9][0-9]?)\.(25[0-5]|2[0-4][0-9]|[01]?[0-9][0-9]?)$'
    return True if re.match(pattern, ip_str) else False


def check_net_mask(mask):
    """
    Return the validity of the mask

     >>>eisoopylib.isValidMask("255.255.255.0")
    True
    >>>eisoopylib.isValidMask("192.168.0")
    False
    >>>eisoopylib.isValidMask("test")
    False
    >>>eisoopylib.isValidMask("0.0.0.0")
    False
    >>>eisoopylib.isValidMask("255.255.255.255")
    True

    etc.
    """
    try:
        if check_net_ip(mask):
            mask_num, = unpack("!I", inet_aton(mask))
            if mask_num == 0:
                return False

            # get inverted
            mask_num = ~mask_num + 1
            binstr = bin(mask_num)[3:]
            # convert to positive integer
            binstr = '0b%s' % ''.join('1' if b == '0' else '0' for b in binstr)
            mask_num = int(binstr, 2) + 1
            # check 2^n
            return (mask_num & (mask_num - 1) == 0)
        return False
    except Exception:
        return False


def remove_duplicate_item_from_list(items, key=None):
    """
    从列表中移除重复项
    """
    seen = set()
    for item in items:
        val = item if key is None else key(item)
        if val not in seen:
            yield item
            seen.add(val)


def check_filename(file_name):
    """
    检查文件名是否合法
    """
    legal = True
    # 文件名长度不能超过255
    if len(file_name) == 0 or len(file_name) > 255:
        legal = False

    # 避免使用加号、减号或者"."作为普通文件的第一个字符
    black_list = ['+', '-', '.']
    if file_name[0] in black_list:
        legal = False

    # 文件名避免使用下列特殊字符,包括制表符和退格符
    black_list = ['/', '\t', '\b', '@', '#', '$', '%', '^', '&', '*', '(', ')', '[', ']']
    intersection = set(black_list) & set(file_name)
    if len(intersection) != 0:
        legal = False

    if not legal:
        raise_exception(exp_msg=_("IDS_INVALID_FILENAME"),
                        exp_num=ncTShareMgntError.NCT_INVALID_FILENAME)


def get_machine_code():
    """
    获取本节点的机器码
    """
    node = uuid.getnode()
    return uuid.UUID(int=node).hex[-12:]


def get_os_version():
    """
    获取本机的linux发行版本
    """
    return "CentOS release 6.5 (Final)"


def remove_prefix_u(data):
    """
    去除掉data中的字符串中的 u
    """
    if isinstance(data, str):
        return data.encode('utf-8')
    elif isinstance(data, collections.Mapping):
        return dict(list(map(remove_prefix_u, iter(data.items()))))
    elif isinstance(data, collections.Iterable):
        return type(data)(list(map(remove_prefix_u, data)))
    else:
        return data


def check_tel_number(tel_number):
    """
    检查手机号合法性
    """
    pattern = re.compile(r'^\d{1,20}$')
    return True if pattern.match(tel_number) else False


def generate_random_key(num):
    """
    生成随机秘钥， num为随机数个数，只包含字母和数字
    """
    return ''.join(random.sample(string.ascii_letters + string.digits, num))

def db_operator(exec_sql):
    """
    数据库操作
    """
    results = DBConnector().r_db.all(exec_sql)
    return results

def strip_whitespace(*args):
    """
    批量去除参数两边的空格
    """
    return (each.strip() for each in args if each)

def check_url(url):
    """
    检测 url 合法性
    """
    if not url.strip():
        raise_exception(exp_msg=_("IDS_URL_EMPTY"),
                        exp_num=ncTShareMgntError.NCT_URL_EMPTY)
    regex = re.compile(
        r'^(?:http)s?://'                                                                     # http:// 或 https://
        r'(?:(?:[A-Z0-9](?:[A-Z0-9-]{0,61}[A-Z0-9])?\.)+(?:[A-Z]{2,6}\.?|[A-Z0-9-]{2,}\.?)|'  # 域名
        r'localhost|'                                                                         # localhost
        r'\d{1,3}\.\d{1,3}\.\d{1,3}\.\d{1,3}|'                                                # IPv4
        r'\[(([0-9a-fA-F]{1,4}:){7,7}[0-9a-fA-F]{1,4}|([0-9a-fA-F]{1,4}:){1,7}:|([0-9a-fA-F]{1,4}:){1,6}:[0-9a-fA-F]{1,4}|([0-9a-fA-F]{1,4}:){1,5}(:[0-9a-fA-F]{1,4}){1,2}|([0-9a-fA-F]{1,4}:){1,4}(:[0-9a-fA-F]{1,4}){1,3}|([0-9a-fA-F]{1,4}:){1,3}(:[0-9a-fA-F]{1,4}){1,4}|([0-9a-fA-F]{1,4}:){1,2}(:[0-9a-fA-F]{1,4}){1,5}|[0-9a-fA-F]{1,4}:((:[0-9a-fA-F]{1,4}){1,6})|:((:[0-9a-fA-F]{1,4}){1,7}|:)|fe80:(:[0-9a-fA-F]{0,4}){0,4}%[0-9a-zA-Z]{1,}|::(ffff(:0{1,4}){0,1}:){0,1}((25[0-5]|(2[0-4]|1{0,1}[0-9]){0,1}[0-9])\.){3,3}(25[0-5]|(2[0-4]|1{0,1}[0-9]){0,1}[0-9])|([0-9a-fA-F]{1,4}:){1,4}:((25[0-5]|(2[0-4]|1{0,1}[0-9]){0,1}[0-9])\.){3,3}(25[0-5]|(2[0-4]|1{0,1}[0-9]){0,1}[0-9]))\])'  # IPv6
        r'(?::\d+)?'                                                                          # 可选择的端口
        r'(?:/?|[/?]\S+)$', re.IGNORECASE)
    if not re.match(regex, url):
        raise_exception(exp_msg=_("IDS_INVALID_URL"),
                        exp_num=ncTShareMgntError.NCT_INVALID_URL)

def is_valid_string(value):
    """
    检查是否为合法的字符串
    不能包含 \ / : * ? " < > | 特殊字符
    长度最大为128字节
    Return:
        bool
    """
    if not isinstance(value, str):
        value = value.strip().decode('utf-8')

    return re.match(r'^[^\\\/\:\*\?\"\<\>\|]{1,128}$', value) is not None

def is_valid_string2(value):
    """
    检查是否为合法的字符串
    不能包含 \ / * ? " < > | 特殊字符
    长度最大为128字节
    中广核去除:字符检查
    Return:
        bool
    """
    if not isinstance(value, str):
        value = value.strip().decode('utf-8')

    return re.match(r'^[^\\\/\*\?\"\<\>\|]{1,128}$', value) is not None

def is_code_string(value):
    """
    检查编码是否为合法字段
    只支持大小写英文，数字，下划线，横线
    最大长度255
    """
    if not isinstance(value, str):
        value = value.strip().decode('utf-8')

    return re.match(r'^[a-zA-Z0-9_-]{1,255}$', value) is not None

def merge_dicts(dict1, dict2):
    """
    合并两个多级字典
    """
    res = {**dict1, **dict2}  # 创建一个新字典,合并两个字典的键值对

    # 对于存在于两个字典的键,如果值都是字典,则递归合并
    for key, value1 in dict1.items():
        if key in dict2:
            value2 = dict2[key]
            if isinstance(value1, dict) and isinstance(value2, dict):
                res[key] = merge_dicts(value1, value2)

    return res
