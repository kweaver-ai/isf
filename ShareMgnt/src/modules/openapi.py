#!/usr/bin/python3
# -*- coding:utf-8 -*-
"""openapi相关"""

import uuid
from src.common.db.connector import DBConnector
from src.common.lib import (raise_exception, generate_sign)
from ShareMgnt.ttypes import ncTShareMgntError


class OpenApi(DBConnector):

    """
    openapi相关
    """

    def __init__(self):
        """
        init
        """

    def __get_key_by_id(self, appid):
        """
        通过appid获取appkey
        """
        sql = """
        SELECT `f_app_key`, `f_enabled` FROM `t_third_auth_info`
        WHERE `f_app_id` = %s
        """
        result = self.r_db.one(sql, appid)
        if not result:
            raise_exception(exp_msg=_("APPID_NOT_EXIST"),
                            exp_num=ncTShareMgntError.APPID_NOT_EXIST)

        if not result['f_enabled']:
            raise_exception(exp_msg=_("APPID_DISABLED"),
                            exp_num=ncTShareMgntError.APPID_DISABLED)

        return result['f_app_key']

    def __check_appid_exist(self, appid):
        """
        检查第三方用户id是否存在
        """
        sql = """
        SELECT `f_app_key` FROM `t_third_auth_info`
        WHERE `f_app_id` = %s
        """
        result = self.r_db.one(sql, appid)
        return True if result else False

    def check_third_app_sign(self, appid, method, body, sign):
        """
        检查第三方签名认证
        """

        # 1. 通过appid获取appkey
        appkey = self.__get_key_by_id(appid)

        # 2. 进行签名计算
        gen_sign = self.gen_third_app_sign(appid, appkey, method, body)

        # 3. 签名比较
        if gen_sign != sign:
            raise_exception(exp_msg=_("SIGN_AUTH_FAILED"),
                            exp_num=ncTShareMgntError.SIGN_AUTH_FAILED)

    def gen_third_app_sign(self, appid, appkey, method, body):
        """
        生成签名
        """
        if not appid or not appkey or not method:
            return ''

        # 拼接需要签名的数据
        body = body if body is not None else ''
        data = appid + method + body

        return bytes.decode(generate_sign(appkey.encode('utf-8'), data.encode('utf-8')))

    def add_third_app(self, appid):
        """
        添加第三方应用
        """
        if self.__check_appid_exist(appid):
            raise_exception(exp_msg=_("APPID_HAS_EXIST"),
                            exp_num=ncTShareMgntError.APPID_HAS_EXIST)

        if len(appid) > 50 or len(appid) == 0:
            raise_exception(exp_msg=_("APPID_LENGTH_MORE_THAN_50_OR_EMPTY"),
                            exp_num=ncTShareMgntError.APPID_LENGTH_MORE_THAN_50)

        # 生成随机的appkey
        appkey = str(uuid.uuid1())

        self.w_db.insert('t_third_auth_info', {'f_app_id': appid,
                                               'f_app_key': appkey})
        return appkey

    def set_third_app_status(self, appid, status):
        """
        设置第三方应用状态
        """
        if not self.__check_appid_exist(appid):
            raise_exception(exp_msg=_("APPID_NOT_EXIST"),
                            exp_num=ncTShareMgntError.APPID_NOT_EXIST)

        update_sql = """
        UPDATE `t_third_auth_info` SET `f_enabled` = %s
        WHERE `f_app_id` = %s
        """
        self.w_db.query(update_sql, status, appid)

    def get_web_client_auth_info(self):
        """
        获取登录web client的认证信息
        """
        # 系统内置登录web client appid
        appid = "anyshare"

        sql = """
        SELECT `f_app_key`
        FROM `t_third_auth_info`
        WHERE `f_app_id` = %s
        """
        result = self.r_db.one(sql, appid)
        if not result:
            appkey = str(uuid.uuid1())
            self.w_db.insert('t_third_auth_info', {'f_app_id': appid,
                                                   'f_app_key': appkey})
        else:
            appkey = result['f_app_key']

        return (appid, appkey)