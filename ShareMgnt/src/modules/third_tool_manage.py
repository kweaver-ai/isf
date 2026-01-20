#!/usr/bin/python3
# -*- coding:utf-8 -*-
"""
第三方工具管理类
"""
import re
import urllib.request, urllib.error, urllib.parse
from ShareMgnt.ttypes import (ncTThirdPartyToolConfig,
                              ncTShareMgntError,
                              ncTThirdToolAuthInfo)
from src.common.db.connector import DBConnector
from src.common.sharemgnt_logger import ShareMgnt_Log
from src.common.lib import (raise_exception, check_url)
from src.modules.config_manage import ConfigManage
from src.common.encrypt import simple
from src.common.business_date import BusinessDate
import json
import time
import hashlib
import random
import ssl


class ThirdPartyToolManage(DBConnector):
    """
    第三方工具管理类
    """
    def __init__(self):
        """
        初始化配置信息
        """
        # 第三方工具列表
        self.thirdPartyToolIdList = set(["CAD", "OFFICE", "ANYROBOT", "WOPI", "DOCCONVERT","SURSEN", "TENCENTVTS"])
        self.cadToolNameList = set(["hc", "mx"])
        self.config_manage = ConfigManage()

    def get_third_party_tool_config(self, thirdPartyToolId):
        """
        获取第三方工具配置信息
        """
        # 对 thirdPartyToolId 进行合法性检测
        if not thirdPartyToolId:
            raise_exception(exp_msg=_("IDS_TOOL_ID_EMPTY"),
                            exp_num=ncTShareMgntError.NCT_TOOL_ID_IS_EMPTY)
        if thirdPartyToolId not in self.thirdPartyToolIdList:
                raise_exception(exp_msg=_("IDS_TOOL_ID_NOT_SUPPORT"),
                                exp_num=ncTShareMgntError.NCT_TOOL_ID_NOT_SUPPORT)

        sql = """
        SELECT *
        FROM `t_third_party_tool_config`
        WHERE `f_tool_id` = %s
        """
        result = self.r_db.one(sql, thirdPartyToolId)
        third_config = ncTThirdPartyToolConfig()
        authInfo = ncTThirdToolAuthInfo()
        if result:
            third_config.thirdPartyToolId = result['f_tool_id']
            third_config.enabled = result['f_enabled']
            third_config.url = result['f_url']
            third_config.thirdPartyToolName = result['f_tool_name']
            authInfo.appid = result['f_app_id']
            authInfo.appkey = ""
            if result['f_app_key'] is not None:
                authInfo.appkey = bytes.decode(simple.des_decrypt_simple(result['f_app_key']))
            third_config.authInfo = authInfo

        return third_config

    def set_third_party_tool_config(self, thirdPartyToolConfig):
        """
        设置第三方工具配置信息
        """

        # 对 thirdPartyToolId 进行合法性检测
        if not thirdPartyToolConfig.thirdPartyToolId:
            raise_exception(exp_msg=_("IDS_TOOL_ID_EMPTY"),
                            exp_num=ncTShareMgntError.NCT_TOOL_ID_IS_EMPTY)
        if thirdPartyToolConfig.thirdPartyToolId not in self.thirdPartyToolIdList:
            raise_exception(exp_msg=_("IDS_TOOL_ID_NOT_SUPPORT"),
                            exp_num=ncTShareMgntError.NCT_TOOL_ID_NOT_SUPPORT)
        if thirdPartyToolConfig.thirdPartyToolName is None or \
            thirdPartyToolConfig.thirdPartyToolId != "CAD":
            thirdPartyToolConfig.thirdPartyToolName = ""
        if thirdPartyToolConfig.thirdPartyToolId == "CAD" and \
            thirdPartyToolConfig.thirdPartyToolName not in self.cadToolNameList:
            raise_exception(exp_msg=_("IDS_TOOL_NAME_NOT_SUPPORT"),
                            exp_num=ncTShareMgntError.NCT_TOOL_NAME_NOT_SUPPORT)

        authInfo = ncTThirdToolAuthInfo()
        if thirdPartyToolConfig.authInfo is not None:
            if len(thirdPartyToolConfig.authInfo.appid) > 50 or len(thirdPartyToolConfig.authInfo.appid) == 0:
                raise_exception(exp_msg=_("APPID_LENGTH_MORE_THAN_50_OR_EMPTY"),
                                exp_num=ncTShareMgntError.APPID_LENGTH_MORE_THAN_50)
            else:
                authInfo.appid = thirdPartyToolConfig.authInfo.appid
            if len(thirdPartyToolConfig.authInfo.appkey) > 50 or len(thirdPartyToolConfig.authInfo.appkey) == 0:
                raise_exception(exp_msg=_("APPKEY_LENGTH_MORE_THAN_50_OR_EMPTY"),
                                exp_num=ncTShareMgntError.APPKEY_LENGTH_MORE_THAN_50_OR_EMPTY)
            else:
                authInfo.appkey = simple.des_encrypt_simple(thirdPartyToolConfig.authInfo.appkey)

        # 当第三方工具设置为关闭时不覆盖url，不覆盖工具名
        check_sql = """
        select f_url from t_third_party_tool_config where f_tool_id = %s
        """
        result = self.r_db.one(check_sql, thirdPartyToolConfig.thirdPartyToolId)

        if thirdPartyToolConfig.enabled:
            # url 合法性检测
            check_url(thirdPartyToolConfig.url)

            if result:
                sql = """
                update t_third_party_tool_config set f_enabled = %s, f_url = %s, f_tool_name = %s, f_app_id = %s, f_app_key = %s
                where f_tool_id = %s
                """
                self.w_db.query(sql, thirdPartyToolConfig.enabled,
                            thirdPartyToolConfig.url,
                            thirdPartyToolConfig.thirdPartyToolName,
                            authInfo.appid,
                            authInfo.appkey,
                            thirdPartyToolConfig.thirdPartyToolId)
            else:
                sql = """
                INSERT INTO `t_third_party_tool_config`
                (`f_tool_id`, `f_enabled`, `f_url`, `f_tool_name`, `f_app_id`, `f_app_key`)
                VALUES (%s, %s, %s, %s, %s, %s)
                """

                self.w_db.query(sql, thirdPartyToolConfig.thirdPartyToolId,
                            thirdPartyToolConfig.enabled,
                            thirdPartyToolConfig.url,
                            thirdPartyToolConfig.thirdPartyToolName,
                            authInfo.appid,
                            authInfo.appkey)

        else:
            if result:
                sql = """
                update t_third_party_tool_config set f_enabled = %s where f_tool_id = %s
                """
                self.w_db.query(sql, thirdPartyToolConfig.enabled, thirdPartyToolConfig.thirdPartyToolId)
            else:
                sql = """
                INSERT INTO `t_third_party_tool_config`
                (`f_tool_id`, `f_enabled`, `f_app_id`, `f_app_key`, `f_tool_name`)
                VALUES (%s, %s, %s, %s, %s)
                """
                self.w_db.query(sql, thirdPartyToolConfig.thirdPartyToolId,
                            thirdPartyToolConfig.enabled,
                            authInfo.appid,
                            authInfo.appkey,
                            thirdPartyToolConfig.thirdPartyToolName)

        ShareMgnt_Log('set third party tool confid success, %s', thirdPartyToolConfig.thirdPartyToolId)

    def test_third_party_tool_config(self, url):
        """
        测试第三方工具配置信息
        """
        # url 合法性检测
        check_url(url)

        # url 可访问性检测
        try:
            context = ssl._create_unverified_context()
            result = urllib.request.urlopen(url, context=context, timeout=3)
        except Exception as e:
            ShareMgnt_Log(str(e))
            return False
        if result.getcode() == 200:
            return True

    def get_anyrobot_url(self, host, account):
        """
        获取AnyRobot跳转URL
        """
        check_url(host)
        db_config = self.config_manage.get_config("anyrobot_config")
        config = json.loads(db_config)
        appId = config["appId"]
        appSecret = config["appSecret"]
        timestamp = '%i' % int(BusinessDate.time())
        m = hashlib.md5()
        m.update(appId + appSecret + account + timestamp)
        sign = m.hexdigest()
        reqArgs = {}
        reqArgs["appId"] = appId
        reqArgs["account"] = account
        reqArgs["timestamp"] = timestamp
        reqArgs["sign"] = sign
        url = host + config["uri"] + '?' + urllib.parse.urlencode(reqArgs)
        return url