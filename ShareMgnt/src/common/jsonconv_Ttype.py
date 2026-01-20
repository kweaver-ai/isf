#!/usr/bin/python3
# -*- coding:utf-8 -*-
"""一堆thrift类型转换器"""
import json
from ShareMgnt.ttypes import (ncTSmtpSrvConf,
                              ncTAlarmConfig)


class SmtpConfEnc(json.JSONEncoder):
    def default(self, obj):
        if isinstance(obj, ncTSmtpSrvConf):
            return {"server": obj.server,
                    "safeMode": obj.safeMode,
                    "port": obj.port,
                    "email": obj.email,
                    "password": obj.password,
                    "openRelay": obj.openRelay}
        return json.JSONEncoder.default(self, obj)


class SmtpConfDec(json.JSONDecoder):
    def decode(self, obj):
        dct = json.JSONDecoder.decode(self, obj)
        return ncTSmtpSrvConf(dct["server"],
                              dct["safeMode"],
                              dct["port"],
                              dct["email"],
                              dct["password"],
                              dct["openRelay"])

# ----------------------------------------------------------------------------


class AlarmSenderEnum:
    Email = 1


class AlarmConfEnc(json.JSONEncoder):
    def default(self, obj):
        if isinstance(obj, ncTAlarmConfig):
            dct = {}
            # 添加等级配置
            level = {}
            if (obj.infoConfig & AlarmSenderEnum.Email) != 0:
                level["info"] = ["email"]
            if (obj.warnConfig & AlarmSenderEnum.Email) != 0:
                level["warn"] = ["email"]
            if level:
                dct["level"] = level

            # 添加邮箱列表
            if obj.emailToList:
                dct["email"] = obj.emailToList

            return dct
        return json.JSONEncoder.default(self, obj)


class AlarmConfDec(json.JSONDecoder):
    def decode(self, obj):
        dct = json.JSONDecoder.decode(self, obj)
        ret = ncTAlarmConfig(0, 0, [])
        # 等级配置
        if "level" in dct:
            for level, levelConf in dct["level"].items():
                tmpConf = 0
                for item in levelConf:
                    if item == "email":
                        tmpConf |= AlarmSenderEnum.Email
                if level == "info":
                    ret.infoConfig = tmpConf
                elif level == "warn":
                    ret.warnConfig = tmpConf

        if "email" in dct:
            ret.emailToList = dct["email"]
        return ret
