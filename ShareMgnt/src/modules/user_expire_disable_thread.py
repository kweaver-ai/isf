#!/usr/bin/python3
# -*- coding:utf-8 -*-
"""用户账号过期自动禁用线程"""
from src.common import global_info
import time
import datetime
import threading

from src.common.sharemgnt_logger import ShareMgnt_Log
from src.modules.user_manage import UserManage
from src.common.db.connector import DBConnector
from src.common.eacp_log import eacp_log
from src.common.business_date import BusinessDate
from ShareMgnt.constants import (NCT_USER_ADMIN,
                                 NCT_USER_AUDIT,
                                 NCT_USER_SYSTEM,
                                 NCT_USER_SECURIT)
from eisoo.tclients import TClient

WAIT_TIME_HOUR = 23
WAIT_TIME_MINUTE = 59
WAIT_TIME_SECOND = 59
USER_EXPIRE_DISABLED = 0x00000002   # 用户账号过期禁用
USER_DISABLED = 0x00000040


class UserExpireDisableThread(threading.Thread):

    def __init__(self):
        """
        初始化
        """
        super(UserExpireDisableThread, self).__init__()
        self.user_manage = UserManage()
        self.db_operator = DBConnector()

    def get_need_disable_userinfos(self):
        """
        获取需要禁用的用户信息
        """
        query_sql = """
        SELECT `f_user_id`, `f_display_name`, `f_login_name`
        FROM `t_user`
        WHERE `f_expire_time` <= {0}
              AND `f_expire_time` != %s
              AND f_auto_disable_status & {1} = 0
              AND `f_user_id` != '{2}'
              AND `f_user_id` != '{3}'
              AND `f_user_id` != '{4}'
              AND `f_user_id` != '{5}'
        """.format(int(BusinessDate.time()), USER_EXPIRE_DISABLED, NCT_USER_ADMIN,
                   NCT_USER_AUDIT, NCT_USER_SYSTEM, NCT_USER_SECURIT)
        results = self.db_operator.r_db.all(query_sql, -1)

        user_infos = []
        for res in results:
            user_infos.append([res["f_user_id"], res["f_display_name"], res["f_login_name"]])

        return user_infos

    def get_wait_time(self):
        """
        获取距离下一次线程运行需要等待的时间
        """
        cur_time = BusinessDate.now()
        cur_hour = cur_time.hour
        cur_minute = cur_time.minute
        cur_second = cur_time.second
        wait_time = (WAIT_TIME_HOUR - cur_hour) * 60 * 60 + \
                    (WAIT_TIME_MINUTE - cur_minute) * 60 +   \
                    (WAIT_TIME_SECOND - cur_second)

        # 为防止等待时间为0时多次运行线程, 用+1来处理边界情况
        if wait_time == 0:
            wait_time += 1

        return wait_time

    def run(self):
        """
        执行
        """
        ShareMgnt_Log("**************** user expire disable thread start *****************")

        while True:
            try:
                time.sleep(self.get_wait_time())

                # 获取待禁用的用户信息
                user_infos = self.get_need_disable_userinfos()
                for user_info in user_infos:
                    # 更新禁用标志
                    self.user_manage.set_user_auto_disable_status(user_info[0], USER_EXPIRE_DISABLED)

                    # 记录禁用日志
                    msg = _("IDS_DISABLE_EXPIRED_USER") % (user_info[1], user_info[2])
                    ex_msg = _("IDS_DISABLE_EXPIRED_USER_EXMSG")
                    eacp_log(_("IDS_SYSTEM"),
                            global_info.LOG_TYPE_MANAGE,
                            global_info.USER_TYPE_INTER,
                            global_info.LOG_LEVEL_WARN,
                            global_info.LOG_OP_TYPE_SET,
                             msg,
                             ex_msg,
                             raise_ex=True)

            except Exception as e:
                ShareMgnt_Log("user expire disable thread run error: %s", str(e))

        ShareMgnt_Log("**************** user expire disable thread end *****************")
