#!/usr/bin/python3
# -*- coding:utf-8 -*-
"""
屏蔽组织架构显示
"""
from src.common.db.connector import DBConnector
from src.modules.department_manage import DepartmentManage
from src.common.lib import (raise_exception, escape_key, generate_group_str)
from ShareMgnt.ttypes import (ncTShareMgntError, ncTUsrmDepartmentInfo)
from src.modules.user_manage import UserManage
from src.modules.config_manage import ConfigManage


class HideOuManage(DBConnector):
    def __init__(self):
        """
        """
        self.dept_manage = DepartmentManage()
        self.user_manage = UserManage()
        self.config_manage = ConfigManage()

    def add(self, departmentIds):
        """
        添加部门
        """
        if not departmentIds:
            raise_exception(exp_msg=_("depart or organ not exists"),
                            exp_num=ncTShareMgntError.NCT_ORG_OR_DEPART_NOT_EXIST)

        # 检查部门列表是否存在,并且去重
        depart_ids = []
        for department_id in departmentIds:
            if department_id in depart_ids:
                continue
            self.dept_manage.check_depart_exists(department_id, True)
            depart_ids.append(department_id)

        if len(depart_ids) > 0:
            # 同时插入多行记录 INSERT INTO tbl_name (a,b,c) VALUES(1,2,3),(4,5,6),(7,8,9);

            groupStr = generate_group_str(depart_ids)
            # 获取存在数据库里面的值
            check_sql = """
            select f_department_id from t_hide_ou where f_department_id in ({0})
            """.format(groupStr)
            results = self.r_db.all(check_sql)
            exist_depart_ids = []
            for result in results:
                exist_depart_ids.append(result['f_department_id'])

            # 获取需要插入的值
            not_exist_depart_ids = []
            for department_id in depart_ids:
                if department_id not in exist_depart_ids:
                    tmp = "('%s')" % self.w_db.escape(department_id)
                    not_exist_depart_ids.append(tmp)

            if len(not_exist_depart_ids) == 0:
                return

            insert_sql = """
            INSERT INTO t_hide_ou (f_department_id)
            VALUES {0}
            """.format(','.join(not_exist_depart_ids))
            self.w_db.query(insert_sql)

    def get(self):
        """
        获取屏蔽组织架构的部门信息
        """
        query_sql = """
        SELECT d.f_department_id, d.f_name
        FROM t_hide_ou AS h
        INNER JOIN t_department AS d
        ON h.f_department_id = d.f_department_id
        ORDER BY upper(d.f_name)
        """
        results = self.r_db.all(query_sql)

        depart_infos = []
        for res in results:
            depart_info = ncTUsrmDepartmentInfo()
            depart_info.departmentId = res['f_department_id']
            depart_info.departmentName = res['f_name']
            depart_infos.append(depart_info)
        return depart_infos

    def delete(self, departmentId):
        """
        根据id删除屏蔽组织架构中的部门
        """
        delete_sql = """
        DELETE FROM `t_hide_ou`
        WHERE `f_department_id` = %s
        """
        affect_row = self.w_db.query(delete_sql, departmentId)
        if not affect_row:
            raise_exception(exp_msg=_("depart or organ not exists"),
                            exp_num=ncTShareMgntError.NCT_ORG_OR_DEPART_NOT_EXIST)

    def search(self, search_key):
        """
        根据部门名搜索屏蔽组织架构的部门信息
        """
        esckey = "%%%s%%" % escape_key(search_key)
        search_sql = """
        SELECT `d`.`f_department_id`, `d`.`f_name`
        FROM `t_hide_ou` AS `h`
        INNER JOIN `t_department` AS `d`
        ON `h`.`f_department_id` = `d`.`f_department_id`
        WHERE `d`.`f_name` LIKE %s
        ORDER BY upper(`d`.`f_name`)
        """
        results = self.r_db.all(search_sql, esckey)

        depart_infos = []
        for res in results:
            depart_info = ncTUsrmDepartmentInfo()
            depart_info.departmentId = res['f_department_id']
            depart_info.departmentName = res['f_name']
            depart_infos.append(depart_info)
        return depart_infos

    def check(self, user_id):
        """
        检查是否需要屏蔽组织架构
        """

        if int(self.config_manage.get_config('hide_ou_info')) == 0:
            return False

        # 获取用户的所有父部门
        parent_dept_ids = self.user_manage.get_all_path_dept_id(user_id)
        if parent_dept_ids:
            groupStr = generate_group_str(parent_dept_ids)
            query_sql = """
            SELECT COUNT(*) AS cnt
            FROM t_hide_ou
            WHERE f_department_id in ({0})
            """.format(groupStr)

            result = self.r_db.one(query_sql)
            if result['cnt'] == 0:
                return False

        return True
