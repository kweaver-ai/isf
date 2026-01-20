#!/usr/bin/python3
# -*- coding:utf-8 -*-
"""
简易的存取json config的manage类
"""
from src.common.db.connector import DBConnector
import json
from src.common.lib import (check_email,
                            check_smtp_params,
                            raise_exception)
from ShareMgnt.ttypes import ncTShareMgntError
from src.modules.config_manage import ConfigManage


class JsonConfManage(DBConnector):
    def __init__(self, key, encCls, decCls):
        self.key = key
        self.encCls = encCls
        self.decCls = decCls
        self.config_manage = ConfigManage()

    def get_config(self):
        """
        获得配置
        PS:这个接口会返回None，在外层再处理
        """
        sql = """
        SELECT `f_value`
        FROM `t_sharemgnt_config`
        WHERE `f_key` = %s
        """
        result = self.r_db.one(sql, self.key)
        if result:
            return json.loads(result["f_value"], cls=self.decCls)
        return None

    def set_config(self, conf):
        """
        设置配置
        PS:外层在通知线程更新配置
        """
        if not conf.openRelay:
            if not conf.password:
                raise_exception(exp_msg=_("IDS_SMTP_PASSWORD_NOT_SET"),
                                exp_num=ncTShareMgntError.NCT_SMTP_PASSWORD_NOT_SET)

        # 检查SMTP配置格式
        check_smtp_params(conf)

        oldConf = self.get_config()
        jsonConf = json.dumps(conf, cls=self.encCls)
        if oldConf:
            sql = """
            UPDATE `t_sharemgnt_config`
            SET `f_value` = %s
            WHERE `f_key` = %s
            """
            self.w_db.query(sql, jsonConf, self.key)
        else:
            self.w_db.insert("t_sharemgnt_config", [self.key, jsonConf])


class MailRecipient(DBConnector):
    def __init__(self, key):
        self.key = key
        self.config_manage = ConfigManage()

    def get_config(self):
        """
        获得邮箱收件人列表
        """
        sql = """
        SELECT `f_value`
        FROM `t_sharemgnt_config`
        WHERE `f_key` = %s
        """
        result = self.r_db.one(sql, self.key)
        if result:
            return json.loads(result["f_value"])
        return []

    def set_config(self, tomail):
        """
        配置邮箱收件人列表
        """
        # 验证toList是电子邮件格式正确性
        for e in tomail:
            if not check_email(e):
                raise_exception(exp_msg=_("email illegal"),
                                exp_num=ncTShareMgntError.NCT_INVALID_EMAIL)
        jsonConf = json.dumps(tomail)

        self.config_manage.replace_config(self.key, jsonConf)
