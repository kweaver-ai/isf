#!/usr/bin/python3
# -*- coding:utf-8 -*-
"""回收站管理"""

from src.common.db.connector import DBConnector
from ShareMgnt.ttypes import (ncTRecycleInfo)

class RecycleManage(DBConnector):

    def set_info(self, info, cid):
        """
        设置回收站配置
        """
        check_sql = """
        select f_gns from t_recycle where f_cid = %s
        """

        insert_sql = """
        insert into t_recycle (f_cid, f_gns, f_setter, f_retention_days) values (%s, %s, %s, %s)
        """

        update_sql = """
        update t_recycle set f_gns = %s, f_setter = %s, f_retention_days = %s where f_cid = %s
        """
        result = self.r_db.one(check_sql, cid)
        if result:
            self.w_db.query(update_sql, info.gns, info.setter, info.retentionDays, cid)
        else:
            self.w_db.query(insert_sql, cid, info.gns, info.setter, info.retentionDays)

    def get_info(self, cid):
        """
        获取回收站配置
        """
        sql = """
        SELECT `f_gns`, `f_setter`, `f_retention_days`
        FROM `t_recycle` WHERE `f_cid` = %s
        """
        info = self.r_db.one(sql, cid)
        ret = ncTRecycleInfo()
        if not info:
            ret.gns = ""
            ret.setter = ""
            ret.retentionDays = -1
        else:
            ret.gns = info["f_gns"]
            ret.setter = info["f_setter"]
            ret.retentionDays = int(info["f_retention_days"])
        return ret

    def del_info(self, cid):
        """
        删除回收站配置
        """
        sql = """
        DELETE FROM `t_recycle` WHERE `f_cid` = %s
        """
        self.w_db.query(sql, cid)

    def get_all_info(self):
        """
        获取所有回收站配置
        """
        sql = """
        SELECT f_gns, f_setter, f_retention_days
        FROM t_recycle
        """
        results = self.r_db.all(sql)

        infos = []
        for res in results:
            info = ncTRecycleInfo()
            info.gns = res["f_gns"]
            info.setter = res["f_setter"]
            info.retentionDays = int(res["f_retention_days"])
            infos.append(info)

        return infos
