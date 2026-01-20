#!/usr/bin/python3
# -*- coding:utf-8 -*-
"""This is spcae manage class"""
import uuid
from src.common import global_info
from src.common.db.connector import DBConnector
from src.common.lib import raise_exception
from ShareMgnt.ttypes import (ncTShareMgntError,
                              ncTThirdDBInfo,
                              ncTThirdTableInfo,
                              ncTThirdDepartTableInfo,
                              ncTThirdDepartRelationTableInfo,
                              ncTThirdUserTableInfo,
                              ncTThirdUserDepartRelationTableInfo,
                              ncTThirdDbSyncConfig)


class ThirdDBManage(DBConnector):
    """
    Third manage
    """
    def __init__(self):
        super(ThirdDBManage, self).__init__()

    def check_third_db_exists(self, third_db_id, raise_ex=True):
        """
        检查第三方数据库id是否存在
        """
        sql = """
        SELECT count(*) AS cnt FROM `t_third_party_db`
        WHERE `f_third_db_id` = %s
        """
        result = self.r_db.one(sql, third_db_id)

        b_exist = True if result['cnt'] else False

        if raise_ex and not b_exist:
            raise_exception(exp_msg=_("IDS_THIRD_DB_NOT_EXISTS"),
                            exp_num=ncTShareMgntError.NCT_THIRD_DB_NOT_EXISTS)

        return b_exist

    def check_third_table_exists(self, third_table_id, raise_ex=True):
        """
        检查第三方数据库id是否存在
        """
        table_names = ['t_third_depart_table',
                       't_third_depart_relation_table',
                       't_third_user_table',
                       't_third_user_relation_table']

        b_exist = False
        for table in table_names:
            sql = """
            SELECT count(*) AS cnt FROM {0}
            WHERE `f_table_id` = %s
            """.format(table)
            result = self.r_db.one(sql, third_table_id)
            b_exist = True if result['cnt'] else False
            if b_exist:
                break

        if raise_ex and not b_exist:
            raise_exception(exp_msg=_("IDS_THIRD_TABLE_NOT_EXISTS"),
                            exp_num=ncTShareMgntError.NCT_THIRD_TABLE_NOT_EXISTS)

        return b_exist

    def add_third_db_info(self, third_db_info):
        """
        增加第三方数据库信息
        """
        third_db_id = str(uuid.uuid1())
        data = {
            "f_third_db_id": third_db_id,
            "f_name": third_db_info.name,
            "f_ip": third_db_info.ip,
            "f_port": third_db_info.port,
            "f_admin": third_db_info.admin,
            "f_password": third_db_info.password,
            "f_database": third_db_info.database,
            "f_charset": third_db_info.charset,
            "f_db_type": third_db_info.dbType,
            "f_status": third_db_info.status,
        }

        self.w_db.insert('t_third_party_db', data)
        return third_db_id

    def get_third_db_info(self, third_db_id):
        """
        获取第三方数据库信息
        """
        self.check_third_db_exists(third_db_id, False)

        sql = """
        SELECT `f_third_db_id`, `f_name`, `f_ip`, `f_port`, `f_admin`,
        `f_password`, `f_database`, `f_db_type`, `f_charset`, `f_status`
        FROM `t_third_party_db` WHERE `f_third_db_id` = %s
        """
        result = self.r_db.one(sql, third_db_id)

        third_db_info = ncTThirdDBInfo()
        if result:
            third_db_info.id = str(result['f_third_db_id'])
            third_db_info.name = result['f_name']
            third_db_info.ip = result['f_ip']
            third_db_info.port = str(result['f_port'])
            third_db_info.admin = result['f_admin']
            third_db_info.password = result['f_password']
            third_db_info.database = result['f_database']
            third_db_info.dbType = int(result['f_db_type'])
            third_db_info.charset = result['f_charset']
            third_db_info.status = bool(result['f_status'])
        return third_db_info

    def edit_third_db_info(self, third_db_info):
        """
        编辑第三方数据信息
        """
        self.check_third_db_exists(third_db_info.id)

        sql = """
        UPDATE `t_third_party_db` SET `f_name` = %s, `f_ip` = %s,
        `f_port` = %s, `f_admin` = %s, `f_password` = %s, `f_database`= %s,
        `f_charset` = %s, `f_db_type` = %s, `f_status` = %s
        WHERE `f_third_db_id` = %s
        """
        self.w_db.query(sql,
                        third_db_info.name,
                        third_db_info.ip,
                        int(third_db_info.port),
                        third_db_info.admin,
                        third_db_info.password,
                        third_db_info.database,
                        third_db_info.charset,
                        int(third_db_info.dbType),
                        bool(third_db_info.status),
                        third_db_info.id)

    def delete_third_db_info(self, third_db_id):
        """
        删除第三方数据库信息
        """
        self.check_third_db_exists(third_db_id)

        sql = """
        DELETE FROM `t_third_party_db` WHERE `f_third_db_id` = %s
        """
        self.w_db.query(sql, third_db_id)

        # 关闭指定第三方数据库同步线程
        if global_info.THIRD_DB_SYNC_THREAD.get(third_db_id):
            global_info.THIRD_DB_SYNC_THREAD[third_db_id].close()

    def get_third_db_table_infos(self, third_db_id):
        """
        获取第三方数据库的所有表信息
        """
        self.check_third_db_exists(third_db_id)

        third_table_infos = ncTThirdTableInfo()

        depart_table_infos = self.get_third_depart_table_infos(third_db_id)
        depart_rela_table_infos = self.get_third_depart_relation_table_infos(third_db_id)
        user_table_infos = self.get_third_user_table_infos(third_db_id)
        user_rela_table_infos = self.get_third_user_relation_table_infos(third_db_id)

        third_table_infos.thirdDepartTableInfos = depart_table_infos
        third_table_infos.thirdDepartRelationTableInfos = depart_rela_table_infos
        third_table_infos.thirdUserTableInfos = user_table_infos
        third_table_infos.thirdUserRelationTableInfos = user_rela_table_infos
        return third_table_infos

    def add_third_depart_table_info(self, table_info):
        """
        增加第三方部门表信息
        """
        self.check_third_db_exists(table_info.thirdDbId)

        table_id = str(uuid.uuid1())
        sub_group_name = ""
        if table_info.customSubGroupNames:
            sub_group_name = ",".join([name for name in table_info.customSubGroupNames])
        data = {
            "f_table_id": table_id,
            "f_third_db_id": table_info.thirdDbId,
            "f_table_name": table_info.tableName,
            "f_department_id": table_info.departmentIdField,
            "f_department_name": table_info.departmentNameField,
            "f_deparment_priority": table_info.departmentPriorityField,
            "f_filter": table_info.filter,
            "f_sub_group": sub_group_name
        }
        self.w_db.insert('t_third_depart_table', data)
        return table_id

    def get_third_depart_table_infos(self, third_db_id):
        """
        """
        self.check_third_db_exists(third_db_id)

        sql = """
        SELECT * FROM `t_third_depart_table`
        WHERE `f_third_db_id` = %s
        """

        depart_table_infos = []
        results = self.r_db.all(sql, third_db_id)
        for res in results:
            table_info = ncTThirdDepartTableInfo()
            table_info.thirdDbId = res['f_third_db_id']
            table_info.tableId = res['f_table_id']
            table_info.tableName = res['f_table_name']
            table_info.departmentIdField = res['f_department_id']
            table_info.departmentNameField = res['f_department_name']
            table_info.departmentPriorityField = res['f_deparment_priority']
            table_info.filter = res['f_filter']
            if res['f_sub_group']:
                table_info.customSubGroupNames = res['f_sub_group'].split(",")

            depart_table_infos.append(table_info)
        return depart_table_infos

    def edit_third_depart_table_info(self, table_info):
        """
        编辑第三方部门表信息
        """
        self.check_third_db_exists(table_info.thirdDbId)
        self.check_third_table_exists(table_info.tableId)

        sub_group_name = ""
        if table_info.customSubGroupNames:
            sub_group_name = ",".join([name for name in table_info.customSubGroupNames])

        sql = """
        UPDATE `t_third_depart_table` SET `f_table_name` = %s, `f_department_id` = %s,
        `f_department_name` = %s, `f_deparment_priority` = %s, `f_filter` = %s,
        `f_sub_group` = %s WHERE `f_table_id` = %s
        """

        self.w_db.query(sql,
                        table_info.tableName,
                        table_info.departmentIdField,
                        table_info.departmentNameField,
                        table_info.departmentPriorityField,
                        table_info.filter,
                        sub_group_name,
                        table_info.tableId)

    def add_third_depart_relation_table_info(self, table_info):
        """
        增加第三方部门关系表信息
        """
        self.check_third_db_exists(table_info.thirdDbId)

        sub_group_name = ""
        if table_info.parentCustomGroupName:
            sub_group_name = ",".join([name for name in table_info.parentCustomGroupName])

        table_id = str(uuid.uuid1())
        data = {
            "f_table_id": table_id,
            "f_third_db_id": table_info.thirdDbId,
            "f_table_name": table_info.tableName,
            "f_department_id": table_info.departmentIdField,
            "f_parent_department_id": table_info.parentDepartmentIdField,
            "f_parent_group_table_id": table_info.parentCustomGroupTableId,
            "f_parent_group_name": sub_group_name,
            "f_filter": table_info.filter,
        }

        self.w_db.insert('t_third_depart_relation_table', data)
        return table_id

    def edit_third_depart_relation_table_info(self, table_info):
        """
        编辑第三方部门关系表信息
        """
        self.check_third_db_exists(table_info.thirdDbId)
        self.check_third_table_exists(table_info.tableId)

        sub_group_name = ""
        if table_info.parentCustomGroupName:
            sub_group_name = ",".join([name for name in table_info.parentCustomGroupName])

        sql = """
        UPDATE `t_third_depart_relation_table` SET `f_table_name` = %s, `f_department_id` = %s,
        `f_parent_department_id` = %s, `f_parent_group_table_id` = %s, `f_filter` = %s,
        `f_parent_group_name` = %s WHERE `f_table_id` = %s
        """

        self.w_db.query(sql,
                        table_info.tableName,
                        table_info.departmentIdField,
                        table_info.parentDepartmentIdField,
                        table_info.parentCustomGroupTableId,
                        table_info.filter,
                        sub_group_name,
                        table_info.tableId)

    def get_third_depart_relation_table_infos(self, third_db_id):
        """
        获取第三方部门关系表信息
        """
        self.check_third_db_exists(third_db_id)

        sql = """
        SELECT * FROM `t_third_depart_relation_table`
        WHERE `f_third_db_id` = %s
        """

        depart_relation_table_infos = []
        results = self.r_db.all(sql, third_db_id)
        for res in results:
            table_info = ncTThirdDepartRelationTableInfo()
            table_info.thirdDbId = res['f_third_db_id']
            table_info.tableId = res['f_table_id']
            table_info.tableName = res['f_table_name']
            table_info.departmentIdField = res['f_department_id']
            table_info.parentDepartmentIdField = res['f_parent_department_id']
            table_info.parentCustomGroupTableId = res['f_parent_group_table_id']

            table_info.parentCustomGroupName = []
            if res['f_parent_group_name']:
                table_info.parentCustomGroupName = res['f_parent_group_name'].split(",")

            table_info.filter = res['f_filter']
            depart_relation_table_infos.append(table_info)
        return depart_relation_table_infos

    def add_third_user_table_info(self, table_info):
        """
        增加第三方用户信息表
        """
        self.check_third_db_exists(table_info.thirdDbId)

        table_id = str(uuid.uuid1())
        data = {
            "f_table_id": table_id,
            "f_third_db_id": table_info.thirdDbId,
            "f_table_name": table_info.tableName,
            "f_user_id": table_info.userIdField,
            "f_user_login_name": table_info.userLoginNameField,
            "f_user_display_name": table_info.userDisplayNameField,
            "f_user_email": table_info.userEmailField,
            "f_user_password": table_info.userPasswordField,
            "f_user_status": table_info.userStatusField,
            "f_user_priority": table_info.userPriorityField,
            "f_filter": table_info.filter,
        }
        self.w_db.insert('t_third_user_table', data)
        return table_id

    def edit_third_user_table_info(self, table_info):
        """
        增加第三方用户信息表
        """
        self.check_third_db_exists(table_info.thirdDbId)
        self.check_third_table_exists(table_info.tableId)

        sql = """
        UPDATE `t_third_user_table` SET `f_table_name` = %s, `f_user_id` = %s,
        `f_user_login_name` = %s, `f_user_display_name` = %s, `f_user_email` = %s,
        `f_user_password` = %s, `f_user_status` = %s, `f_user_priority` = %s,
        `f_filter` = %s WHERE `f_table_id` = %s
        """

        self.w_db.query(sql,
                        table_info.tableName,
                        table_info.userIdField,
                        table_info.userLoginNameField,
                        table_info.userDisplayNameField,
                        table_info.userEmailField,
                        table_info.userPasswordField,
                        table_info.userStatusField,
                        table_info.userPriorityField,
                        table_info.filter,
                        table_info.tableId)

    def get_third_user_table_infos(self, third_db_id):
        """
        根据第三方数据库id获取第三方用户表信息
        """
        self.check_third_db_exists(third_db_id)

        sql = """
        SELECT * FROM `t_third_user_table`
        WHERE `f_third_db_id` = %s
        """

        third_user_table_infos = []
        results = self.r_db.all(sql, third_db_id)
        for res in results:
            table_info = ncTThirdUserTableInfo()
            table_info.thirdDbId = res['f_third_db_id']
            table_info.tableId = res['f_table_id']
            table_info.tableName = res['f_table_name']
            table_info.userIdField = res['f_user_id']
            table_info.userLoginNameField = res['f_user_login_name']
            table_info.userDisplayNameField = res['f_user_display_name']
            table_info.userEmailField = res['f_user_email']
            table_info.userPasswordField = res['f_user_password']
            table_info.userStatusField = res['f_user_status']
            table_info.userPriorityField = res['f_user_priority']
            table_info.filter = res['f_filter']
            third_user_table_infos.append(table_info)
        return third_user_table_infos

    def add_third_user_relation_table_info(self, table_info):
        """
        增加第三方用户关系表
        """
        self.check_third_db_exists(table_info.thirdDbId)

        parent_group_name = ""
        if table_info.parentCustomGroupName:
            parent_group_name = ",".join([name for name in table_info.parentCustomGroupName])

        table_id = str(uuid.uuid1())
        data = {
            "f_table_id": table_id,
            "f_third_db_id": table_info.thirdDbId,
            "f_table_name": table_info.tableName,
            "f_user_id": table_info.userIdField,
            "f_parent_department_id": table_info.parentDepartmentIdField,
            "f_parent_group_table_id": table_info.parentCustomGroupTableId,
            "f_parent_group_name": parent_group_name,
            "f_filter": table_info.filter,
        }
        self.w_db.insert('t_third_user_relation_table', data)
        return table_id

    def edit_third_user_relation_table_info(self, table_info):
        """
        编辑第三方用户关系表
        """
        self.check_third_db_exists(table_info.thirdDbId)
        self.check_third_table_exists(table_info.tableId)

        parent_group_name = ""
        if table_info.parentCustomGroupName:
            parent_group_name = ",".join([name for name in table_info.parentCustomGroupName])

        sql = """
        UPDATE `t_third_user_relation_table` SET `f_table_name` = %s, `f_user_id` = %s,
        `f_parent_department_id` = %s, `f_parent_group_table_id` = %s,
        `f_parent_group_name` = %s, `f_filter` = %s WHERE `f_table_id` = %s
        """

        self.w_db.query(sql,
                        table_info.tableName,
                        table_info.userIdField,
                        table_info.parentDepartmentIdField,
                        table_info.parentCustomGroupTableId,
                        parent_group_name,
                        table_info.filter,
                        table_info.tableId)

    def get_third_user_relation_table_infos(self, third_db_id):
        """
        """
        self.check_third_db_exists(third_db_id)

        sql = """
        SELECT * FROM `t_third_user_relation_table`
        WHERE `f_third_db_id` = %s
        """

        third_user_relation_table_infos = []
        results = self.r_db.all(sql, third_db_id)
        for res in results:
            table_info = ncTThirdUserDepartRelationTableInfo()
            table_info.thirdDbId = res['f_third_db_id']
            table_info.tableId = res['f_table_id']
            table_info.tableName = res['f_table_name']
            table_info.userIdField = res['f_user_id']
            table_info.parentDepartmentIdField = res['f_parent_department_id']
            table_info.parentCustomGroupTableId = res['f_parent_group_table_id']

            table_info.parentCustomGroupName = []
            if res['f_parent_group_name']:
                table_info.parentCustomGroupName = res['f_parent_group_name'].split(",")

            table_info.filter = res['f_filter']
            third_user_relation_table_infos.append(table_info)

        return third_user_relation_table_infos

    def delete_third_table(self, third_table_id):
        """
        删除第三方表信息
        """
        self.check_third_table_exists(third_table_id)

        table_names = ['t_third_depart_table',
                       't_third_depart_relation_table',
                       't_third_user_table',
                       't_third_user_relation_table']

        for table in table_names:
            sql = """
            DELETE FROM {0} WHERE `f_table_id` = %s
            """.format(table)
            self.w_db.query(sql, third_table_id)

    def get_third_db_sync_config(self, third_db_id):
        """
        获取第三方数据库的同步配置
        """
        self.check_third_db_exists(third_db_id)

        sql = """
        SELECT * FROM `t_third_party_db`
        WHERE `f_third_db_id` = %s
        """

        result = self.r_db.one(sql, third_db_id)
        sync_config = ncTThirdDbSyncConfig()
        if result:
            sync_config.parentDepartId = result['f_parent_department_id']
            sync_config.thirdRootName = result['f_third_root_name']
            sync_config.thirdRootId = result['f_third_root_id']
            sync_config.syncInterval = int(result['f_sync_interval'])
            sync_config.spaceSize = int(result['f_space_size'])
            sync_config.userType = int(result['f_user_type'])
        return sync_config

    def add_third_db_sync_config(self, third_db_id, sync_config):
        """
        设置第三方数据库的同步配置
        """
        self.check_third_db_exists(third_db_id)

        sql = """
        UPDATE `t_third_party_db` SET `f_parent_department_id` = %s,
        `f_third_root_name` = %s, `f_third_root_id` = %s,
        `f_sync_interval` = %s, `f_space_size` = %s, `f_user_type` = %s
        WHERE `f_third_db_id` = %s
        """
        self.w_db.query(sql,
                        sync_config.parentDepartId,
                        sync_config.thirdRootName,
                        sync_config.thirdRootId,
                        sync_config.syncInterval,
                        sync_config.spaceSize,
                        sync_config.userType,
                        third_db_id)

    def get_status(self, third_db_id):
        """
        获取第三方数据库同步的启用状态
        """
        sql = """
        SELECT `f_status` FROM `t_third_party_db`
        WHERE `f_third_db_id` = %s
        """
        result = self.r_db.one(sql, third_db_id)
        return bool(result['f_status'])

    def set_status(self, third_db_id, status):
        """
        设置第三方数据库同步的状态
        """
        sql = """
        UPDATE `t_third_party_db` SET `f_status` = %s
        WHERE `f_third_db_id` = %s
        """
        status = 1 if status else 0
        self.w_db.query(sql, status, third_db_id)

    def get_enable_third_db_ids(self):
        """
        获取状态开启的第三方数据库id
        """
        sql = """
        SELECT `f_third_db_id` FROM `t_third_party_db`
        WHERE `f_status` = %s
        """
        results = self.r_db.all(sql, 1)

        third_db_ids = []
        if results:
            third_db_ids = [res['f_third_db_id'] for res in results]

        return third_db_ids
