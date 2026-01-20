#!/usr/bin/python3
# -*- coding:utf-8 -*-
"""更新用户自动禁用线程"""
import threading
import time
import datetime

from src.common.sharemgnt_logger import ShareMgnt_Log
from src.modules.config_manage import ConfigManage
from src.modules.user_manage import UserManage
from src.common.db.connector import DBConnector
from src.common.business_date import BusinessDate
from ShareMgnt.constants import (NCT_USER_ADMIN,
                                 NCT_USER_AUDIT,
                                 NCT_USER_SYSTEM,
                                 NCT_USER_SECURIT)

WAIT_TIME = 86400   # 24  * 3600
USER_AUTO_DISABLED = 0x00000001     # 用户长时间不登录禁用


class UserAutoDisableThread(threading.Thread):

    def __init__(self):
        """
        初始化
        """
        super(UserAutoDisableThread, self).__init__()
        self.config_manage = ConfigManage()
        self.user_manage = UserManage()
        self.db_operator = DBConnector()

    def get_need_disable_userids(self, days):
        """
        获取需要禁用的用户id
        """
        now = BusinessDate.now()

        delta = datetime.timedelta(days)

        allow_max_last_request_time = (now - delta).strftime('%Y-%m-%d %H:%M:%S')
        query_sql = """
        SELECT f_user_id
        FROM t_user
        WHERE f_last_request_time < %s
              AND f_auto_disable_status & {4} = 0
              AND `f_user_id` != '{0}'
              AND `f_user_id` != '{1}'
              AND `f_user_id` != '{2}'
              AND `f_user_id` != '{3}'
        """.format(NCT_USER_ADMIN,
                   NCT_USER_AUDIT, NCT_USER_SYSTEM, NCT_USER_SECURIT,
                   USER_AUTO_DISABLED)
        results = self.db_operator.r_db.all(query_sql, allow_max_last_request_time)

        user_ids = []
        for res in results:
            user_ids.append(res["f_user_id"])

        return user_ids

    def run(self):
        """
        执行
        """
        ShareMgnt_Log("**************** user auto disable thread start *****************")

        while True:
            try:

                # 获取自动禁用配置
                user_auto_disable_config = self.config_manage.get_auto_disable_config()
                if user_auto_disable_config.isEnabled:
                    # 如果开启自动禁用， 则自动禁用
                    user_ids = self.get_need_disable_userids(user_auto_disable_config.days)
                    for user_id in user_ids:
                        self.user_manage.set_user_auto_disable_status(user_id, USER_AUTO_DISABLED)

            except Exception as e:
                ShareMgnt_Log("user auto disable thread run error: %s", str(e))
            time.sleep(WAIT_TIME)

        ShareMgnt_Log("**************** user auto disable thread end *****************")
