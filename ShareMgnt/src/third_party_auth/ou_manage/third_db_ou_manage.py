#!/usr/bin/python3
# -*- coding:utf-8 -*-
from src.common.db.oracle_manage import OracleManage
from src.common.db.mssql_manage import MsSQLManage
from src.common.db.mysql_manage import MySQLManage
from ShareMgnt.ttypes import *
from src.third_party_auth.ou_manage.base_ou_manage import *


MYSQL = 1
ORACLE = 2
MSSQL = 3

manage_class_dict = {MYSQL: MySQLManage,
                     ORACLE: OracleManage,
                     MSSQL: MsSQLManage}

SQL_TEMPLATE = """
SELECT * FROM %s
"""


class ThirdDBOuManage(BaseOuManage):
    """
    第三方数据库组织用户管理类
    """
    def __init__(self, b_eacplog=False):
        """
        初始化函数
        """
        super(ThirdDBOuManage, self).__init__()
        self.ou_user_tree = {}
        self.ou_depart_tree = {}
        self.group_dict = {}
        self.db_manage = None
        self.db_info = None
        self.table_info = None
        self.sync_info = None
        self.root_id = None

    def init_db_manage(self):
        """
        获取第三方数据管理器
        """
        manage_class = manage_class_dict[self.db_info.dbType]
        if manage_class:
            self.db_manage = manage_class(self.db_info.ip,
                                          int(self.db_info.port),
                                          self.db_info.admin,
                                          self.db_info.password,
                                          self.db_info.database)
            self.db_operate = self.db_manage.get_db_operator()

    def init_server_info(self, server_info):
        """
        """
        self.db_info = server_info['dbInfo']
        self.table_info = server_info['tableInfo']
        self.sync_info = server_info['syncInfo']
        self.init_db_manage()

        self.ou_user_tree = {}
        self.ou_depart_tree = {}

        self.parse_departs()
        self.parse_depart_relation()
        self.parse_users()
        self.parse_user_depart_relation()

    def decode(self, content):
        """
        """
        if self.db_info.charset and content:
            content = content.decode(self.db_info.charset)

        return content

    def execute_sql(self, sql):
        """
        执行sql
        """
        if self.db_info.dbType == ORACLE:
            cursor = self.db_manage.get_cursor()
            cursor.execute(sql)
            desc = [d[0] for d in cursor.description]
            result = [dict(list(zip(desc, line))) for line in cursor]
            cursor.close()
            return result
        else:
            return self.db_operate.all(sql)

    def query_result(self, table_info):
        """
        """
        sql = SQL_TEMPLATE % table_info.tableName
        if table_info.filter:
            sql = sql + " where " + table_info.filter

        return self.execute_sql(sql)

    def parse_result(self, res, fields_expr):
        """
        """
        result = ""

        # 使用或“||”来划分子语句
        sub_expr_list = fields_expr.split("||")

        # 按顺序解析每一个字段，只到数据库中有一个字段的值不为空为止
        for sub_expr in sub_expr_list:
            result = ""
            b_field_exist = False

            field_list = sub_expr.split("+")
            for field in field_list:
                if not field:
                    continue

                field = field.strip()
                if field.find("@") == 0 and field[1:]:
                    result += field[1:]
                else:
                    # 解析要替换处理的字符
                    replace_dict = {}
                    start = field.find("(")
                    end = field.find(")")

                    if start > 0 and end > 0 and start < end:
                        replace_expr = field[start + 1: end].split(",")
                        field = field[0: start]
                        for r_expr in replace_expr:
                            if r_expr.find("->") != -1:
                                c1, c2 = r_expr.split("->")
                                if c1:
                                    replace_dict[c1] = c2

                    # 解析要取字段范围
                    start = field.find("[")
                    end = field.find("]")
                    field_start = field_end = 0
                    if start > 0 and end > 0 and start < end:
                        field_start, field_end = field[start + 1: end].split(":")
                        field = field[0: start]

                    if field in res and res[field]:
                        value = str(res[field])

                        # 截取值
                        if field_start != field_end:
                            value = value[int(field_start): int(field_end)]

                        # 替换字符
                        for c1, c2 in list(replace_dict.items()):
                            value = value.replace(c1, c2)

                        b_field_exist = True
                        result += self.decode(value)

            # 如果表达式中一个字段也没有值，则进行下一个或语句处理
            if result and b_field_exist:
                break

        return result

    def add_sub_group(self, table_info, parent_ou_info):
        """
        增加子分组信息
        """
        if not table_info.customSubGroupNames:
            return

        for sub_group_name in table_info.customSubGroupNames:
            # 统计自定义组的字典，结构为tableId->groupName->[A,B,C]
            if sub_group_name not in self.group_dict[table_info.tableId]:
                self.group_dict[table_info.tableId][sub_group_name] = set()

            group_info = OuInfo()
            group_info.third_id = parent_ou_info.third_id + sub_group_name
            group_info.ou_name = sub_group_name

            if group_info.third_id not in self.ou_depart_tree:
                self.ou_depart_tree[group_info.third_id] = group_info

            parent_ou_info.sub_third_ou_ids.append(group_info.third_id)
            self.group_dict[table_info.tableId][sub_group_name].add(group_info.third_id)

    def parse_departs(self):
        """
        解析部门信息
        """
        # 新建根组织
        self.root_id = self.sync_info.thirdRootId if self.sync_info.thirdRootId else "-1"
        root_info = OuInfo()
        root_info.third_id = self.root_id
        root_info.ou_name = self.sync_info.thirdRootName
        self.ou_depart_tree[root_info.third_id] = root_info

        # 解析部门信息
        if self.table_info.thirdDepartTableInfos:
            for table_info in self.table_info.thirdDepartTableInfos:
                # 统计自定义组的字典，结构为tableId->groupName->[A,B,C]
                self.group_dict[table_info.tableId] = {}

                # 如果表名为ROOT_SUB_GROUP,则是在根组织下添加分组
                if table_info.tableName == "ROOT_SUB_GROUP":
                    # 添加子分组
                    self.add_sub_group(table_info, root_info)
                else:
                    results = self.query_result(table_info)
                    for res in results:
                        # 添加部门信息
                        ou_info = OuInfo()
                        ou_info.third_id = self.parse_result(res, table_info.departmentIdField)
                        ou_info.ou_name = self.parse_result(res, table_info.departmentNameField)

                        if table_info.departmentPriorityField:
                            ou_info.priority = res[table_info.departmentPriorityField]
                            if ou_info.priority:
                                ou_info.priority = int(ou_info.priority)
                        self.ou_depart_tree[ou_info.third_id] = ou_info

                        # 添加子分组
                        self.add_sub_group(table_info, ou_info)

    def parse_depart_relation(self):
        """
        解析部门关系
        """
        if self.table_info.thirdDepartRelationTableInfos:
            for table_info in self.table_info.thirdDepartRelationTableInfos:
                # 获取应该添加到的组id
                parent_group_ids = []
                p_table_id = table_info.parentCustomGroupTableId
                p_group_names = table_info.parentCustomGroupName

                if p_table_id and p_group_names and p_table_id in self.group_dict:
                    for group_name in p_group_names:
                        if group_name in self.group_dict[p_table_id]:
                            parent_group_ids += list(self.group_dict[p_table_id][group_name])

                results = self.query_result(table_info)
                for res in results:
                    depart_id = self.parse_result(res, table_info.departmentIdField)

                    # 优先使用数据库里的父部门字段
                    parent_id = None
                    if table_info.parentDepartmentIdField:
                        parent_id = self.parse_result(res, table_info.parentDepartmentIdField)
                    elif parent_group_ids:
                        # 部门只支持添加在一个群组下
                        parent_id = parent_group_ids[0]

                    if not depart_id or depart_id not in self.ou_depart_tree:
                        continue

                    if parent_id in self.ou_depart_tree:
                        parent_ou_info = self.ou_depart_tree[parent_id]
                        if depart_id not in parent_ou_info.sub_third_ou_ids:
                            parent_ou_info.sub_third_ou_ids.append(depart_id)
                    else:
                        # 没有父部门，则默认添加到根组织下
                        root_ou_info = self.ou_depart_tree[self.root_id]
                        if depart_id not in root_ou_info.sub_third_ou_ids:
                            root_ou_info.sub_third_ou_ids.append(depart_id)

    def parse_users(self):
        """
        解析用户信息
        """
        if self.table_info.thirdUserTableInfos:
            for table_info in self.table_info.thirdUserTableInfos:
                results = self.query_result(table_info)
                for res in results:
                    user_info = UserInfo()
                    user_info.third_id = self.parse_result(res, table_info.userIdField)
                    user_info.login_name = self.parse_result(res, table_info.userLoginNameField)

                    if not user_info.login_name or not user_info.third_id:
                        continue

                    if table_info.userDisplayNameField:
                        user_info.display_name = self.parse_result(res, table_info.userDisplayNameField)
                    else:
                        user_info.display_name = user_info.login_name

                    if table_info.userEmailField:
                        user_info.email = self.parse_result(res, table_info.userEmailField)

                    if table_info.userStatusField:
                        user_info.status = res[table_info.userStatusField]

                    if table_info.userPriorityField:
                        user_info.priority = res[table_info.userPriorityField]
                        if user_info.priority:
                            user_info.priority = int(user_info.priority)

                    if table_info.userPasswordField:
                        user_info.password = res[table_info.userPasswordField]

                    if self.sync_info.spaceSize:
                        user_info.space_size = self.sync_info.spaceSize

                    user_info.type = self.sync_info.userType

                    self.ou_user_tree[user_info.third_id] = user_info

    def parse_user_depart_relation(self):
        """
        解析用户部门关系
        """
        def __add_user_ou_relation(user_id, parent_id):
            if not user_id or not parent_id:
                return

            if parent_id in self.ou_depart_tree:
                parent_ou_info = self.ou_depart_tree[parent_id]
                if user_id not in parent_ou_info.sub_third_user_ids:
                    parent_ou_info.sub_third_user_ids.append(user_id)

        if self.table_info.thirdUserRelationTableInfos:
            # 遍历每一个用户部门关系表
            for table_info in self.table_info.thirdUserRelationTableInfos:
                # 获取应该添加到的群组id
                parent_group_ids = []
                p_table_id = table_info.parentCustomGroupTableId
                p_group_names = table_info.parentCustomGroupName

                if p_table_id and p_group_names and p_table_id in self.group_dict:
                    for group_name in p_group_names:
                        if group_name in self.group_dict[p_table_id]:
                            parent_group_ids += list(self.group_dict[p_table_id][group_name])

                # 获取用户id
                results = self.query_result(table_info)
                for res in results:
                    user_id = self.parse_result(res, table_info.userIdField)

                    if not user_id or user_id not in self.ou_user_tree:
                        continue

                    # 添加用户和数据库中的部门关系
                    if table_info.parentDepartmentIdField:
                        ou_parent_id = self.parse_result(res, table_info.parentDepartmentIdField)
                        __add_user_ou_relation(user_id, ou_parent_id)

                    # 添加用户和自定义组关系
                    for parent_id in parent_group_ids:
                        __add_user_ou_relation(user_id, parent_id)

