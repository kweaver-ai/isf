#!/usr/bin/python3
# -*- coding:utf-8 -*-
"""This is link share class"""
from src.common.db.connector import DBConnector
from src.common.lib import (escape_key,
                            check_start_limit,
                            generate_search_order_sql)
from src.modules.config_manage import ConfigManage
from src.modules.user_manage import UserManage
from src.modules.department_manage import DepartmentManage
from ShareMgnt.ttypes import (ncTLinkShareInfo)


class LinkShareManage(DBConnector):

    """
    link share manage
    """

    def __init__(self):
        """
        init
        """
        self.config_manage = ConfigManage()
        self.user_manage = UserManage()
        self.dept_manage = DepartmentManage()

    def __check_share_info(self, share_info):
        """
        检查共享信息
        """
        # 检查部门或用户是否存在
        if share_info.sharerType == 1:
            self.user_manage.check_user_exists(share_info.sharerId)
        else:
            self.dept_manage.check_depart_exists(share_info.sharerId, True)

    def convert_share_info(self, db_infos):
        """
        转换数据库信息
        """
        share_infos = []
        if db_infos:
            for res in db_infos:
                info = ncTLinkShareInfo()
                info.sharerType = res['type']
                info.sharerId = res['sharer_id']
                info.sharerName = res['user_name'] if info.sharerType == 1 else res['dept_name']
                share_infos.append(info)
        return share_infos

    def set_link_share_status(self, status):
        """
        设置系统外链共享状态：
        参数：
           status：
                True： 开启
                Flase：关闭
        """
        status = 1 if status else 0
        self.config_manage.set_config("link_share_status", status)

    def get_link_share_status(self):
        """
        设置系统外链共享状态：
        参数：
           status：
                True： 开启
                Flase：关闭
        """
        result = self.config_manage.get_config('link_share_status')
        status = True if int(result) == 1 else False
        return status

    def add_link_share_info(self, share_info):
        """
        增加一条外链共享策略信息
        """
        self.__check_share_info(share_info)

        self.w_db.insert(
            't_link_share_strategy',
            {
                "f_sharer_id": share_info.sharerId,
                "f_sharer_type": share_info.sharerType
            }
        )
        return share_info.sharerId

    def delete_link_share_info(self, sharer_id):
        """
        删除一条外链共享策略信息
        """
        delete_sql = """
        DELETE FROM `t_link_share_strategy`
        WHERE `f_sharer_id` = %s
        """
        self.r_db.query(delete_sql, sharer_id)

    def get_link_share_info_by_page(self, start, limit):
        """
        分页获取外链共享策略信息
        """
        # 检查参数
        limit_statement = check_start_limit(start, limit)

        query_sql = """
        SELECT `strategy`.`f_sharer_id` AS `sharer_id`,
        `strategy`.`f_sharer_type` AS `type`,
        `t_user`.`f_display_name` AS `user_name`,
        `t_department`.`f_name` AS `dept_name`
        FROM  `t_link_share_strategy` AS `strategy`
        LEFT JOIN `t_user`
        ON `strategy`.`f_sharer_id` = `t_user`.`f_user_id`
        LEFT JOIN `t_department`
        ON `strategy`.`f_sharer_id` = `t_department`.`f_department_id`
        WHERE 1=%s
        ORDER BY upper(`type`),
        upper(`dept_name`),
        upper(`user_name`)
        {0}
        """.format(limit_statement)
        results = self.r_db.all(query_sql, 1)
        return self.convert_share_info(results)

    def get_link_share_info_cnt(self):
        """
        获取外链共享策略信息总条数
        """
        query_sql = """
        SELECT count(*) as cnt
        FROM t_link_share_strategy
        """
        result = self.r_db.one(query_sql)
        return result['cnt']

    def search_link_share_info(self, start, limit, search_key):
        """
        搜索外链贡献策略信息
        """
        limit_statement = check_start_limit(start, limit)

        # 增加搜索排序子句
        order_by_str = generate_search_order_sql(
            ['t_user.f_display_name', 't_department.f_name'])

        query_sql = """
        SELECT `strategy`.`f_sharer_id` AS `sharer_id`,
        `strategy`.`f_sharer_type` AS `type`,
        `t_user`.`f_display_name` AS `user_name`,
        `t_department`.`f_name` AS `dept_name`
        FROM  `t_link_share_strategy` AS `strategy`
        LEFT JOIN `t_user`
        ON `strategy`.`f_sharer_id` = `t_user`.`f_user_id`
        LEFT JOIN `t_department`
        ON `strategy`.`f_sharer_id` = `t_department`.`f_department_id`
        WHERE (`t_user`.`f_display_name` LIKE %s OR  `t_department`.`f_name` LIKE %s)
        ORDER BY upper(`type`),
        {0},
        upper(`dept_name`),
        upper(`user_name`)
        {1}
        """.format(order_by_str, limit_statement)
        esckey = "%%%s%%" % escape_key(search_key)

        results = self.r_db.all(query_sql, esckey, esckey,
                                escape_key(search_key), escape_key(search_key), esckey, esckey)
        return self.convert_share_info(results)
