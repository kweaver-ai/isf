#!/usr/bin/python3
# -*- coding: utf-8 -*-

import time
from src.common.db.connector import DBConnector
from src.common.lib import raise_exception
from src.common.lib import (raise_exception, escape_key, check_start_limit)
from ShareMgnt.ttypes import (ncTShareMgntError, ncTWatermarkDocInfo)

class AntivirusManage(DBConnector):
    """
    """
    def __init__(self):
        pass

    def __check_user_is_already_added(self, userId):
        """
        """
        sql = """
        SELECT COUNT(*) AS cnt FROM `t_antivirus_admin`
        WHERE `f_user_id` = %s
        """
        return True if self.r_db.one(sql, userId)["cnt"] == 1 else False

    def get_all_antivirus_admin(self):
        sql = "SELECT f_user_id FROM t_antivirus_admin"
        admins = []
        results = self.r_db.all(sql)
        for res in results:
            admins.append(res['f_user_id'])
        return admins

    def add_antivirus_admin(self, loginName):
        sql = "SELECT `f_user_id` FROM `t_user` WHERE `f_login_name` = %s"
        adminInfo = self.r_db.one(sql, loginName)
        if not adminInfo:
            raise_exception(exp_msg=_("user not exists"),
                            exp_num=ncTShareMgntError.NCT_USER_NOT_EXIST)
        if not self.__check_user_is_already_added(adminInfo['f_user_id']):
            self.w_db.insert("t_antivirus_admin", {"f_user_id": adminInfo['f_user_id']})
