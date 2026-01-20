#!/usr/bin/python3
# -*- coding:utf-8 -*-
"""
短信服务管理模块
"""

import hashlib
import json
import random
import time
import base64
import hmac
from datetime import datetime
from src.common.db.connector import DBConnector
from src.common.http import send_request
from src.common.lib import (raise_exception, encrypt_pwd)
from src.common.lib import (raise_exception,
                            encrypt_pwd,
                            check_email,
                            check_tel_number)
from src.common.sharemgnt_logger import ShareMgnt_Log
from src.modules.config_manage import ConfigManage
from src.modules.login_manage import LoginManage
from src.modules.user_manage import UserManage
from ShareMgnt.ttypes import ncTShareMgntError
from src.common.encrypt.simple import eisoo_rsa_decrypt
from src.common.business_date import BusinessDate
from src.common.http import pub_nsq_msg

from tencentcloud.common import credential
from tencentcloud.common.profile.client_profile import ClientProfile
from tencentcloud.common.profile.http_profile import HttpProfile
from tencentcloud.common.exception.tencent_cloud_sdk_exception import TencentCloudSDKException
from tencentcloud.sms.v20190711 import sms_client, models

TOPIC_USER_MODIFIED = "user_management.user.modified"


class SmsManage(DBConnector):
    """
    """
    def __init__(self):
        """
        """
        self.config_manage = ConfigManage()
        self.login_manage = LoginManage()
        self.user_manage = UserManage()

    def get_sms_config(self):
        """
        获取短信服务器配置
        """
        sql = """
        SELECT `f_value` FROM `t_sharemgnt_config`
        WHERE `f_key` = %s
        """
        result = self.r_db.one(sql, 'sms_config')
        if result["f_value"]:
            return result["f_value"]
        else:
            # 返回默认配置
            return '{"server_id": "", "server_name": "", "app_id": "", \
                     "secret_id":"", "secret_key":"", "international":0, \
                     "template_id": "", "expire_time": "30"}'

    def set_sms_config(self, config):
        """
        设置短信服务器配置
        """
        # 检查短信服务器配置
        self.check_sms_config(config)

        self.config_manage.replace_config('sms_config', config)

    def check_sms_config(self, config):
        """
        检查短信服务器配置
        """
        try:
            sms_config = json.loads(config)
        except:
            raise_exception(exp_msg=_("IDS_INVALID_SMS_CONFIG"),
                            exp_num=ncTShareMgntError.NCT_INVALID_SMS_CONFIG)

        if "server_id" in sms_config and sms_config["server_id"] == "tencent_cloud":
            self.check_tencent_sms_config(sms_config)
        else:
            raise_exception(exp_msg=_("IDS_NOT_SUPPORT_SMS_SERVER"),
                            exp_num=ncTShareMgntError.NCT_NOT_SUPPORT_SMS_SERVER)

    def check_tencent_sms_config(self, sms_config):
        """
        """
        if not ("app_id" in sms_config and sms_config["app_id"] != "" \
            and "secret_id" in sms_config and sms_config["secret_id"] != "" \
            and "secret_key" in sms_config and sms_config["secret_key"] != "" \
            and "template_id" in sms_config and sms_config["template_id"] != ""):
            raise_exception(exp_msg=_("IDS_INVALID_SMS_CONFIG"),
                            exp_num=ncTShareMgntError.NCT_INVALID_SMS_CONFIG)
        try:
            secretId = bytes.decode(eisoo_rsa_decrypt(sms_config["secret_id"]))
            secretKey = bytes.decode(eisoo_rsa_decrypt(sms_config["secret_key"]))
            cred = credential.Credential(secretId, secretKey)
            httpProfile = HttpProfile()
            httpProfile.endpoint = "sms.tencentcloudapi.com"
            clientProfile = ClientProfile()
            clientProfile.httpProfile = httpProfile
            client = sms_client.SmsClient(cred, "", clientProfile)

            # appid检查（使用套餐包信息统计接口）
            req = models.SmsPackagesStatisticsRequest()
            req.SmsSdkAppid = sms_config["app_id"]
            req.Limit = 10
            req.Offset = 0
            resp = client.SmsPackagesStatistics(req)

            # 模板检查
            req = models.DescribeSmsTemplateListRequest()
            templateIdSet = []
            templateIdSet.append(sms_config["template_id"])
            req.TemplateIdSet = templateIdSet
            if "international" not in sms_config:
                sms_config["international"] = 0
            req.International = sms_config["international"];

            resp = client.DescribeSmsTemplateList(req)
            if len(resp.DescribeTemplateStatusSet) <= 0:
                #模板不存在
                raise_exception("sms template not exists")
        except Exception as e:
            ShareMgnt_Log("connect sms server error: %s", str(e))
            raise_exception(exp_msg=_("IDS_INVALID_SMS_CONFIG"),
                            exp_num=ncTShareMgntError.NCT_CONNECT_SMS_SERVER_ERROR)

    def send_vcode(self, account, password, telNumber):
        """
        发送短信验证码
        """
        if not self.config_manage.get_sms_activate_status():
            raise_exception(exp_msg=_("IDS_SMS_ACTIVATE_DISABLED"),
                            exp_num=ncTShareMgntError.NCT_SMS_ACTIVATE_DISABLED)

        # 根据输入账号名获取用户信息
        db_user = self.login_manage.match_account(account)

        if not db_user:
            raise_exception(exp_msg=_("invalid account or password"),
                            exp_num=ncTShareMgntError.NCT_INVALID_ACCOUNT_OR_PASSWORD)

        if db_user["f_is_activate"] == 1:
            raise_exception(exp_msg=_("IDS_USER_IS_ACTIVATE"),
                            exp_num=ncTShareMgntError.NCT_USER_IS_ACTIVATE)

        # 检查用户密码
        self.login_manage.switch_login(db_user, password)

        # 检查手机号
        self.user_manage.check_user_tel_number(account, telNumber)

        try:
            config = self.get_sms_config()
            json_config = json.loads(config)
        except:
            raise_exception(exp_msg=_("IDS_INVALID_SMS_CONFIG"),
                            exp_num=ncTShareMgntError.NCT_INVALID_SMS_CONFIG)

        if "server_id" in json_config and json_config["server_id"] == "tencent_cloud":
            self.sendsms_by_tencent(json_config, telNumber)
        else:
            raise_exception(exp_msg=_("IDS_NOT_SUPPORT_SMS_SERVER"),
                            exp_num=ncTShareMgntError.NCT_NOT_SUPPORT_SMS_SERVER)

    def sendsms_by_tencent(self, smsConfig, telNumber):
        """
        """
        try:
            secretId = bytes.decode(eisoo_rsa_decrypt(smsConfig["secret_id"]))
            secretKey = bytes.decode(eisoo_rsa_decrypt(smsConfig["secret_key"]))
            cred = credential.Credential(secretId, secretKey)
            httpProfile = HttpProfile()
            httpProfile.endpoint = "sms.tencentcloudapi.com"

            clientProfile = ClientProfile()
            clientProfile.httpProfile = httpProfile
            client = sms_client.SmsClient(cred, "", clientProfile)
            verify_code = self.get_random()

            req = models.SendSmsRequest()
            req.SmsSdkAppid = smsConfig["app_id"]
            req.PhoneNumberSet = ["+86" + telNumber]
            req.TemplateID = smsConfig["template_id"]
            templateParamSet = []
            templateParamSet.append(str(verify_code))
            templateParamSet.append(str(smsConfig["expire_time"]))

            req.TemplateParamSet = templateParamSet
            resp = client.SendSms(req)
            if resp.SendStatusSet:
                sendStatusSet = resp.SendStatusSet
                if sendStatusSet[0].Code != "Ok":
                    ShareMgnt_Log("connect sms server error: %s", sendStatusSet[0].Message)
                    raise_exception(exp_msg=_(sendStatusSet[0].Message),
                                    exp_num=ncTShareMgntError.NCT_CONNECT_SMS_PARAM_ERROR)
            else:
                ShareMgnt_Log("send verify code fail: %s", resp.to_json_string())
                raise_exception(exp_msg=_("IDS_SEND_VERIFY_CODE_FAIL"),
                                exp_num=ncTShareMgntError.NCT_SEND_VERIFY_CODE_FAIL)
            # 保存验证码
            check_sql = """
            select f_verify_code from t_sms_code where f_tel_number = %s
            """
            result = self.r_db.one(check_sql, telNumber)
            now = BusinessDate.now().strftime("%Y-%m-%d %H:%M:%S")
            if result:
                update_sql = """
                update t_sms_code set f_verify_code = %s, f_create_time= %s where f_tel_number = %s
                """
                self.w_db.query(update_sql, verify_code, now, telNumber)
            else:
                insert_sql = """
                insert into t_sms_code(f_tel_number, f_verify_code, f_create_time) values (%s, %s, %s)
                """
                self.w_db.query(insert_sql, telNumber, verify_code, now)
        except Exception as e:
            ShareMgnt_Log("send verify code fail: %s", str(e))
            raise_exception(exp_msg=_("IDS_SEND_VERIFY_CODE_FAIL"),
                            exp_num=ncTShareMgntError.NCT_SEND_VERIFY_CODE_FAIL)

    def activate(self, account, password, telNumber, mailAddress, verifyCode):
        """
        激活账号
        """
        if not self.config_manage.get_sms_activate_status():
            raise_exception(exp_msg=_("IDS_SMS_ACTIVATE_DISABLED"),
                            exp_num=ncTShareMgntError.NCT_SMS_ACTIVATE_DISABLED)

        # 根据输入账号名获取用户信息
        db_user = self.login_manage.match_account(account)

        # 用户不存在
        if not db_user:
            raise_exception(exp_msg=_("invalid account or password"),
                            exp_num=ncTShareMgntError.NCT_INVALID_ACCOUNT_OR_PASSWORD)

        # 是否已激活
        if db_user["f_is_activate"] == 1:
            raise_exception(exp_msg=_("IDS_USER_IS_ACTIVATE"),
                            exp_num=ncTShareMgntError.NCT_USER_IS_ACTIVATE)

        # 检查用户密码
        self.login_manage.switch_login(db_user, password)

        # 检查手机号
        self.user_manage.check_user_tel_number(account, telNumber)

        # 检查邮箱
        self.user_manage.check_user_email(db_user["f_user_id"], mailAddress)

        # 检查验证码是否正确
        self.check_verify_code(telNumber, verifyCode)

        # 标记用户为已激活, 并解除禁用
        update_sql = """
        UPDATE `t_user`
        SET `f_status` = 0, `f_is_activate` = 1,
        `f_tel_number` = %s, `f_mail_address` = %s
        WHERE `f_user_id` = %s
        """
        self.w_db.query(update_sql, telNumber, mailAddress, db_user["f_user_id"])

        user_modify_info = {}
        if db_user['f_tel_number'] != telNumber:
            user_modify_info["new_telephone"] = telNumber
        if db_user['f_mail_address'] != mailAddress:
            user_modify_info["new_email"] = mailAddress
        if len(user_modify_info) > 0:
            user_modify_info["user_id"] = db_user["f_user_id"]
            pub_nsq_msg(TOPIC_USER_MODIFIED, user_modify_info)

        return db_user["f_user_id"]

    def check_verify_code(self, telNumber, verifyCode):
        """
        """
        if not verifyCode:
            raise_exception(exp_msg=_("IDS_SMS_VERIFY_CODE_ERROR"),
                            exp_num=ncTShareMgntError.NCT_SMS_VERIFY_CODE_ERROR)

        try:
            config = self.get_sms_config()
            json_config = json.loads(config)
        except:
            raise_exception(exp_msg=_("IDS_INVALID_SMS_CONFIG"),
                            exp_num=ncTShareMgntError.NCT_INVALID_SMS_CONFIG)

        expire_time = int(json_config["expire_time"]) * 60

        sql = """
        SELECT f_create_time
        FROM `t_sms_code`
        WHERE `f_tel_number` = %s AND `f_verify_code` = %s
        """
        result = self.r_db.one(sql, telNumber, verifyCode)
        if result:
            if (int(result["f_create_time"].timestamp()) + expire_time) < int(BusinessDate.time()):
                raise_exception(exp_msg=_("IDS_SMS_VERIFY_CODE_TIMEOUT"),
                                exp_num=ncTShareMgntError.NCT_SMS_VERIFY_CODE_TIMEOUT)
        else:
            raise_exception(exp_msg=_("IDS_SMS_VERIFY_CODE_ERROR"),
                            exp_num=ncTShareMgntError.NCT_SMS_VERIFY_CODE_ERROR)

    def get_random(self):
        """
        生成六位随机数
        """
        return random.randint(100000, 999999)
