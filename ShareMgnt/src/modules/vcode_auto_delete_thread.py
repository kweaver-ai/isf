#!/usr/bin/python3
# -*- coding:utf-8 -*-
"""更新用户自动禁用线程"""
import threading
import time
import datetime

from src.common.sharemgnt_logger import ShareMgnt_Log
from src.modules.vcode_manage import VcodeManage
from src.common.db.connector import DBConnector
from src.common.business_date import BusinessDate

WAIT_TIME = 86400   # 24  * 3600


class VcodeAutoDeleteThread(threading.Thread):

    def __init__(self):
        """
        初始化
        """
        super(VcodeAutoDeleteThread, self).__init__()
        self.vcode_manage = VcodeManage()
        self.db_operator = DBConnector()

    def delete_all_invalid_vcodes(self, limit_time):
        """
        删除所有失效的验证码
        """
        now = BusinessDate.now()

        delta = datetime.timedelta(limit_time) # days; minutes; seconds. 默认为 days

        allow_max_last_useful_time = (now - delta).strftime('%Y-%m-%d %H:%M:%S')
        query_sql = """
        DELETE FROM t_vcode WHERE f_createtime < '{0}'
        """.format(allow_max_last_useful_time)
        self.db_operator.w_db.query(query_sql)

    def run(self):
        """
        执行
        """
        ShareMgnt_Log("**************** vcode auto delete thread start *****************")

        while True:
            try:
                self.delete_all_invalid_vcodes(1)

            except Exception as e:
                print(("vcode auto delete thread run error: %s", str(e)))
            time.sleep(WAIT_TIME)

        ShareMgnt_Log("**************** vcode auto delete thread end *****************")
