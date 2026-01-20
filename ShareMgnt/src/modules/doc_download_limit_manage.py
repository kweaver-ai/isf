#!/usr/bin/python3
# -*- coding:utf-8 -*-
"""
用户文档下载限制管理类
"""
import sys
import uuid
import time
import datetime
from eisoo.tclients import TClient
from src.modules.user_manage import UserManage
from src.modules.department_manage import DepartmentManage
from src.modules.oem_manage import OEMManage
from src.modules.config_manage import ConfigManage
from src.modules.handle_task_thread import (CallableTask, HandleTaskThread)
from EThriftException.ttypes import ncTException
from ShareMgnt.ttypes import (ncTShareMgntError, ncTDocDownloadLimitInfo, ncTDocDownloadLimitObject)
from ShareMgnt.constants import (NCT_SYSTEM_ROLE_AUDIT)
from EVFS.ttypes import (ncTUserDownloadLimitInfo)
from src.common.db.connector import DBConnector
from src.common.db.connector import ConnectorManager
from src.common.lib import (raise_exception, escape_key,
                            check_start_limit, generate_search_order_sql)
from src.common.nc_senders import (email_send_html_content)
from src.common.global_info import (IS_SINGLE, NC_EVFS_NAME_IOC_DATAEXCHANGE_ID)
from src.common.business_date import BusinessDate
from src.modules.role_manage import RoleManage


USER_OBJ = 1
DEPART_OBJ = 2

MAX_DOC_DOWNLOAD_LIMIT_VALUE = sys.maxsize


class DocDownloadLimitManage(DBConnector):
    def __init__(self):
        """
        """
        self.user_manage = UserManage()
        self.dept_manage = DepartmentManage()
        self.oem_manage = OEMManage()
        self.handle_task_thread = HandleTaskThread()
        self.config_manage = ConfigManage()
        self.role_manage = RoleManage()

    def __check_user_list(self, userList):
        """
        检查用户列表中用户是否存在
        """
        if not userList:
            return
        for user_obj in userList:
            self.user_manage.check_user_exists(user_obj.objectId, True)

    def __check_dept_list(self, deptList):
        """
        检查部门列表中用户是否存在
        """
        if not deptList:
            return

        for dept_obj in deptList:
            self.dept_manage.check_depart_exists(dept_obj.objectId, True)

    def __check_limit_value(self, limitValue):
        """
        检查限制值是否有效
        """
        if limitValue != -1 and limitValue <= 0:
            raise_exception(exp_msg=_("IDS_INVALID_DOC_DOWNLOAD_LIMIT_VALUE"),
                            exp_num=ncTShareMgntError.NCT_INVALID_DOC_DOWNLOAD_LIMIT_VALUE)

    def __remove_duplicate_obj(self, objList):
        """
        去重重复项
        """
        obj_id_set = set()
        obj_list = []
        for obj in objList:
            if obj.objectId not in obj_id_set:
                obj_list.append(obj)
                obj_id_set.add(obj.objectId)

        return obj_list

    def add(self, limitInfo):
        """
        添加一条文档下载量限制
        """
        global IS_SINGLE
        if not IS_SINGLE:
            with TClient('ShareMgntSingle') as client:
                return client.Usrm_AddDocDownloadLimitInfo(limitInfo)

        # 检查限制对象是否配置
        if not limitInfo.userInfos and not limitInfo.depInfos:
            raise_exception(exp_msg=_("IDS_DOC_DOWNLOAD_LIMIT_OBJECT_NOT_SET"),
                            exp_num=ncTShareMgntError.NCT_DOC_DOWNLOAD_LIMIT_OBJECT_NOT_SET)

        # 检查用户列表中用户是否存在
        if limitInfo.userInfos:
            limitInfo.userInfos = self.__remove_duplicate_obj(limitInfo.userInfos)
            self.__check_user_list(limitInfo.userInfos)

        # 检查部门列表中部门是否存在
        if limitInfo.depInfos:
            limitInfo.depInfos = self.__remove_duplicate_obj(limitInfo.depInfos)
            self.__check_dept_list(limitInfo.depInfos)

        # 检查限制值
        self.__check_limit_value(limitInfo.limitValue)

        # 生成一条唯一id
        limitInfo.id = str(uuid.uuid1())

        # 保存数据到数据库
        self.add_doc_download_limit_info_to_db(limitInfo)

        # 增加任务更新 EFAST 中用户的下载限制值配置
        self.add_update_user_doc_download_limit_task()

        return limitInfo.id

    def add_doc_download_limit_info_to_db(self, limitInfo):
        """
        保存数据到数据库
        """
        # 使用事务插入数据
        conn = ConnectorManager.get_db_conn()
        cursor = conn.cursor()

        insert_sql = """
        INSERT INTO `t_doc_download_limit`
        (`f_id`, `f_obj_id`, `f_obj_type`, `f_download_limit_value`, `f_time`)
        VALUES(%s, %s, %s, %s, %s)
        """

        try:
            # 插入用户列表配置
            if limitInfo.userInfos:
                for userInfo in limitInfo.userInfos:
                    cursor.execute(insert_sql, (limitInfo.id, userInfo.objectId,
                                   USER_OBJ, limitInfo.limitValue, int(BusinessDate.time()*1000000)))

            # 插入部门列表配置
            if limitInfo.depInfos:
                for depInfo in limitInfo.depInfos:
                    cursor.execute(insert_sql, (limitInfo.id, depInfo.objectId,
                                   DEPART_OBJ, limitInfo.limitValue, int(BusinessDate.time()*1000000)))

            conn.commit()

        except Exception as e:
            conn.rollback()
            raise_exception(exp_msg=str(e),
                            exp_num=ncTShareMgntError.NCT_DB_OPERATE_FAILED)
        finally:
            cursor.close()
            conn.close()

    def get(self, start, limit):
        """
        获取已配置的文档下载量限制
        """
        limit_statement = check_start_limit(start, limit)

        query_sql = """
        SELECT `f_id`,
        MIN(`f_time`) AS `min_create_time`
        FROM `t_doc_download_limit`
        WHERE 1=%s
        GROUP BY `f_id`
        ORDER BY `min_create_time` DESC
        {0}
        """.format(limit_statement)
        results = self.r_db.all(query_sql, 1)

        limit_infos = []
        for res in results:
            limit_info = self.get_doc_download_limit_info_by_id(res['f_id'])
            if limit_info:
                limit_infos.append(limit_info)

        return limit_infos

    def get_cnt(self):
        """
        获取已配置的文档下载量限制条数
        """
        sql = """
        SELECT COUNT(DISTINCT f_id) as cnt
        FROM t_doc_download_limit
        """
        count = self.r_db.one(sql)
        return count["cnt"]

    def get_doc_download_limit_info_by_id(self, limit_id):
        """
        根据id获取文档下载量限制信息
        """
        query_sql = """
        SELECT *
        FROM `t_doc_download_limit`
        WHERE `f_id` = %s
        """
        results = self.r_db.all(query_sql, limit_id)

        need_delete_obj_id = []

        if results:
            limit_info = ncTDocDownloadLimitInfo()
            limit_info.id = limit_id
            limit_info.limitValue = results[0]['f_download_limit_value']
            limit_info.userInfos = []
            limit_info.depInfos = []

            for res in results:
                obj_id = res['f_obj_id']
                obj_type = res['f_obj_type']

                if obj_type == USER_OBJ:
                    try:
                        obj_info = self.user_manage.get_user_by_id(obj_id)

                        user_info = ncTDocDownloadLimitObject()
                        user_info.objectId = obj_id
                        user_info.objectName = obj_info.user.displayName
                        limit_info.userInfos.append(user_info)

                    except ncTException as ex:
                        if ex.errID == ncTShareMgntError.NCT_USER_NOT_EXIST:
                            # 用户已被删除，需移除对应配置
                            need_delete_obj_id.append(obj_id)
                        else:
                            raise ex

                if obj_type == DEPART_OBJ:
                    try:
                        obj_info = self.dept_manage.get_department_info(obj_id, True)

                        dept_info = ncTDocDownloadLimitObject()
                        dept_info.objectId = obj_id
                        dept_info.objectName = obj_info.departmentName

                        limit_info.depInfos.append(dept_info)

                    except ncTException as ex:
                        if ex.errID == ncTShareMgntError.NCT_ORG_OR_DEPART_NOT_EXIST:
                            # 部门已被删除，需移除对应配置
                            need_delete_obj_id.append(obj_id)
                        else:
                            raise ex

            if need_delete_obj_id:
                self.delete_obj_from_db(need_delete_obj_id)

            if limit_info.userInfos or limit_info.depInfos:
                return limit_info

    def delete_obj_from_db(self, deleteObjIdList):
        """
        删除一条文档下载量限制信息
        """
        for obj_id in deleteObjIdList:
            delete_sql = """
            DELETE FROM `t_doc_download_limit`
            WHERE `f_obj_id` = %s
            """
            self.w_db.query(delete_sql, obj_id)

    def edit_object(self, editId, userList, deptList):
        """
        编辑文档下载量限制对象
        """
        global IS_SINGLE
        if not IS_SINGLE:
            with TClient('ShareMgntSingle') as client:
                return client.Usrm_EditDocDownloadLimitObject(editId, userList, deptList)
        # 检查参数合法性
        if not userList and not deptList:
            raise_exception(exp_msg=_("IDS_DOC_DOWNLOAD_LIMIT_OBJECT_NOT_SET"),
                            exp_num=ncTShareMgntError.NCT_DOC_DOWNLOAD_LIMIT_OBJECT_NOT_SET)

        if not self.__check_doc_download_limit_exist(editId):
            raise_exception(exp_msg=_("IDS_DOC_DOWNLOAD_LIMIT_NOT_EXIST"),
                            exp_num=ncTShareMgntError.NCT_DOC_DOWNLOAD_LIMIT_NOT_EXIST)

        if userList:
            userList = self.__remove_duplicate_obj(userList)
            self.__check_user_list(userList)
        if deptList:
            deptList = self.__remove_duplicate_obj(deptList)
            self.__check_dept_list(deptList)

        self.update_doc_download_limit_object_to_db(editId, userList, deptList)

        # 增加任务更新 EFAST 中用户的下载限制值配置
        self.add_update_user_doc_download_limit_task()

    def update_doc_download_limit_object_to_db(self, updateId, userList, deptList):
        """
        更新数据库中数据
        """
        select_sql = """
        SELECT f_download_limit_value, f_time
        FROM `t_doc_download_limit`
        WHERE `f_id` = %s
        """
        result = self.r_db.one(select_sql, updateId)
        limit_value = result['f_download_limit_value']
        create_time = result['f_time']

        # 获取已经配置的用户id与部门id
        user_ids = []
        dept_ids = []
        select_sql = """
        SELECT f_obj_id, f_obj_type
        FROM `t_doc_download_limit`
        WHERE `f_id` = %s
        """
        result = self.r_db.all(select_sql, updateId)
        for res in result:
            if res['f_obj_type'] == USER_OBJ:
                user_ids.append(res['f_obj_id'])
            if res['f_obj_type'] == DEPART_OBJ:
                dept_ids.append(res['f_obj_id'])

        # 获取新配置的用户对象id
        new_user_ids = []
        if userList:
            new_user_ids = [userObj.objectId for userObj in userList]

        # 获取新配置的部门对象id
        new_dept_ids = []
        if deptList:
            new_dept_ids = [deptObj.objectId for deptObj in deptList]

        # 需要删除的对象id
        need_delete_obj_id_list = []
        for user_id in user_ids:
            if user_id not in new_user_ids:
                need_delete_obj_id_list.append(user_id)
        for dept_id in dept_ids:
            if dept_id not in new_dept_ids:
                need_delete_obj_id_list.append(dept_id)

        # 需要增加的用户对象id
        need_add_user_ids = []
        for user_id in new_user_ids:
            if user_id not in user_ids:
                need_add_user_ids.append(user_id)

        # 需要增加的部门对象id
        need_add_dept_ids = []
        for dept_id in new_dept_ids:
            if dept_id not in dept_ids:
                need_add_dept_ids.append(dept_id)

        # 使用事务更新数据
        conn = ConnectorManager.get_db_conn()
        cursor = conn.cursor()

        delete_sql = """
        DELETE
        FROM `t_doc_download_limit`
        WHERE `f_id` = %s and `f_obj_id` = %s
        """

        insert_sql = """
        INSERT INTO `t_doc_download_limit`
        (`f_id`, `f_obj_id`, `f_obj_type`, `f_download_limit_value`, `f_time`)
        VALUES(%s, %s, %s, %s, %s)
        """
        try:
            # 删除已存在的配置项
            for obj_id in need_delete_obj_id_list:
                cursor.execute(delete_sql, (updateId, obj_id))

            # 插入用户列表配置
            for user_id in need_add_user_ids:
                cursor.execute(insert_sql, (updateId, user_id, USER_OBJ, limit_value, create_time))

            # 插入部门列表配置
            for dept_id in need_add_dept_ids:
                cursor.execute(insert_sql, (updateId, dept_id, DEPART_OBJ,
                               limit_value, create_time))

            conn.commit()

        except Exception as e:
            conn.rollback()
            raise_exception(exp_msg=str(e),
                            exp_num=ncTShareMgntError.NCT_DB_OPERATE_FAILED)
        finally:
            cursor.close()
            conn.close()

    def edit_value(self, editId, limitValue):
        """
        编辑文档下载量限制值
        """
        global IS_SINGLE
        if not IS_SINGLE:
            with TClient('ShareMgntSingle') as client:
                return client.Usrm_EditDocDownloadLimitValue(editId, limitValue)
        # 检查参数合法性
        if not self.__check_doc_download_limit_exist(editId):
            raise_exception(exp_msg=_("IDS_DOC_DOWNLOAD_LIMIT_NOT_EXIST"),
                            exp_num=ncTShareMgntError.NCT_DOC_DOWNLOAD_LIMIT_NOT_EXIST)

        # 检查限速值
        self.__check_limit_value(limitValue)

        update_sql = """
        UPDATE `t_doc_download_limit`
        SET f_download_limit_value = %s
        WHERE `f_id` = %s
        """
        self.w_db.query(update_sql, limitValue, editId)

        # 增加任务更新 EFAST 中用户的下载限制值配置
        self.add_update_user_doc_download_limit_task()

    def __check_doc_download_limit_exist(self, recordId):
        """
        检查文档下载量限制是否配置
        """
        sql = """
        SELECT count(*) as cnt
        FROM `t_doc_download_limit`
        WHERE `f_id` = %s
        """
        result = self.r_db.one(sql, recordId)
        return True if result['cnt'] != 0 else False

    def delete(self, deleteId):
        """
        删除一条文档下载量限制
        """
        global IS_SINGLE
        if not IS_SINGLE:
            with TClient('ShareMgntSingle') as client:
                return client.Usrm_DeleteDocDownloadLimitInfo(deleteId)
        delete_sql = """
        DELETE FROM `t_doc_download_limit`
        WHERE `f_id` = %s
        """
        affect_row = self.w_db.query(delete_sql, deleteId)
        if not affect_row:
            raise_exception(exp_msg=_("IDS_DOC_DOWNLOAD_LIMIT_NOT_EXIST"),
                            exp_num=ncTShareMgntError.NCT_DOC_DOWNLOAD_LIMIT_NOT_EXIST)

        # 增加任务更新 EFAST 中用户的下载限制值配置
        self.add_update_user_doc_download_limit_task()

    def search(self, search_key, start, limit):
        """
        搜索文档下载量限制信息
        """
        # 检查参数
        limit_statement = check_start_limit(start, limit)

        # 增加搜索排序子句
        order_by_str = generate_search_order_sql(['MIN(`t_user`.`f_display_name`)',
                                                  'MIN(`t_department`.`f_name`)'])

        search_sql = """
        SELECT `t_doc_download_limit`.`f_id`,
        MIN(`t_user`.`f_display_name`),
        MIN(`t_department`.`f_name`)
        FROM `t_doc_download_limit`
        LEFT JOIN `t_user`
        ON `t_doc_download_limit`.`f_obj_id` = `t_user`.`f_user_id`
        LEFT JOIN `t_department`
        ON `t_doc_download_limit`.`f_obj_id` = `t_department`.`f_department_id`
        WHERE (`t_user`.`f_display_name` LIKE %s OR `t_department`.`f_name` LIKE %s)
        GROUP BY `t_doc_download_limit`.`f_id`
        ORDER BY {0}
        {1}
        """.format(order_by_str, limit_statement)
        esckey = "%%%s%%" % escape_key(search_key)
        results = self.r_db.all(search_sql, esckey, esckey,
                                escape_key(search_key), escape_key(search_key), esckey, esckey)

        limit_infos = []
        for res in results:
            limit_info = self.get_doc_download_limit_info_by_id(res['f_id'])
            if limit_info:
                limit_infos.append(limit_info)

        return limit_infos

    def search_cnt(self, search_key):
        """
        搜索文档下载量限制条数
        """
        search_sql = """
        SELECT COUNT(DISTINCT `t_doc_download_limit`.`f_id`) as cnt
        FROM `t_doc_download_limit`
        LEFT JOIN `t_user`
        ON `t_doc_download_limit`.`f_obj_id` = `t_user`.`f_user_id`
        LEFT JOIN `t_department`
        ON `t_doc_download_limit`.`f_obj_id` = `t_department`.`f_department_id`
        WHERE (`t_user`.`f_display_name` LIKE %s OR `t_department`.`f_name` LIKE %s)
        """
        esckey = "%%%s%%" % escape_key(search_key)
        result = self.r_db.one(search_sql, esckey, esckey)
        return result['cnt']

    def get_limit_value_by_userid(self, userid):
        """
        根据用户userid获取对应的配置
        """
        global NC_EVFS_NAME_IOC_DATAEXCHANGE_ID
        if userid == NC_EVFS_NAME_IOC_DATAEXCHANGE_ID:
            return -1
        self.user_manage.check_user_exists(userid, True)

        limit_value = MAX_DOC_DOWNLOAD_LIMIT_VALUE

        select_sql = """
        SELECT f_download_limit_value
        FROM `t_doc_download_limit`
        WHERE f_obj_id = %s
        """
        results = self.r_db.all(select_sql, userid)

        if len(results) > 0:
            for res in results:
                if 0 < res['f_download_limit_value'] < limit_value:
                    limit_value = res['f_download_limit_value']
        else:
            # 叶子节点集合
            leaf_nodes_ids = []

            # 各部门的配置
            departs_limit_info = {}

            # 该用户相关的部门树
            depart_tree = self.dept_manage.get_depart_tree_of_user(userid)

            for depart_id in depart_tree:
                if not depart_tree[depart_id].subDepartIds:
                    leaf_nodes_ids.append(depart_id)
                depart_limit_info = {}
                depart_limit_info["valid"] = True   # 各部门配置默认有效
                depart_limit_info["limit_value"] = MAX_DOC_DOWNLOAD_LIMIT_VALUE
                departs_limit_info[depart_id] = depart_limit_info

            # 每个叶子节点向上遍历
            for leaf_nodes_id in leaf_nodes_ids:
                node_id = leaf_nodes_id
                while True:
                    # 本部门配置被标记为无效，则其各父部门配置也被标记为无效，无需向上遍历
                    if not departs_limit_info[node_id]["valid"]:
                        break

                    results = self.r_db.all(select_sql, node_id)
                    if len(results) > 0:
                        tmp_limit_value = MAX_DOC_DOWNLOAD_LIMIT_VALUE
                        for res in results:
                            if 0 < res['f_download_limit_value'] < tmp_limit_value:
                                tmp_limit_value = res['f_download_limit_value']
                        departs_limit_info[node_id]["limit_value"] = tmp_limit_value
                        # 标记所有父部门配置无效
                        cur_node_id = node_id
                        while True:
                            if not depart_tree[cur_node_id].parentDepartId:
                                break
                            else:
                                cur_node_id = depart_tree[cur_node_id].parentDepartId
                                departs_limit_info[cur_node_id]["valid"] = False
                        break
                    else:
                        # 标记本部门限制配置无效
                        departs_limit_info[node_id]["valid"] = False

                    # 父部门不存在，结束遍历
                    if not depart_tree[node_id].parentDepartId:
                        break
                    node_id = depart_tree[node_id].parentDepartId

            # 从有效的配置中，取交集
            for (depart_id, limit_info) in list(departs_limit_info.items()):
                if limit_info["valid"] and (0 < limit_info["limit_value"] < limit_value):
                    limit_value = limit_info["limit_value"]

        if limit_value == MAX_DOC_DOWNLOAD_LIMIT_VALUE or limit_value < 0:
            limit_value = -1

        return limit_value

    def update_efast_download_limit(self, args):
        """
        更新 EFAST 中的下载量限制
        """
        with TClient('EVFS') as Client:
            # 获取所有用户下载量限制信息
            results = Client.GetUserDownloadLimitInfos()

        # 获取用户下载量变化的配置
        changed_limit_infos = []
        need_to_delete_limit_infos = []
        for res in results:
            try:
                limit_value = self.get_limit_value_by_userid(res.userId)
                if limit_value != res.limitValue:
                    limit_info = ncTUserDownloadLimitInfo()
                    limit_info.userId = res.userId
                    limit_info.limitValue = limit_value
                    changed_limit_infos.append(limit_info)
            except ncTException as ex:
                if ex.errID == ncTShareMgntError.NCT_USER_NOT_EXIST:
                    # 用户已被删除，需移除对应配置
                    need_to_delete_limit_infos.append(res.userId)
                else:
                    raise ex

        # 更新用户下载量限制信息
        with TClient('EVFS') as Client:
            Client.AddUserDownloadLimitInfos(changed_limit_infos)
            Client.DeleteUserDownloadLimitInfos(need_to_delete_limit_infos)

    def add_update_user_doc_download_limit_task(self):
        """
        添加更新文件下载量任务
        """
        task = CallableTask()
        task.module_name = "doc_download_limit_manage"
        task.function_name = "update_efast_download_limit"
        self.handle_task_thread.add(task)

        # 三权分立模式下，如果audit开启了接收通知，则当securit修改用户阈值时，audit收到邮件通知
        if (self.user_manage.get_trisystem_status() and
                self.config_manage.get_ddl_email_notify_mode_status()):
            # 获取audit角色下所有用户邮箱
            toEmailList = self.role_manage.get_role_mails(NCT_SYSTEM_ROLE_AUDIT)
            if 0 != len(toEmailList):
                product_name = self.oem_manage.get_config_by_option('shareweb_en-us', 'product')
                subject = _("IDS_EVFS_DOC_DOWNLOAD_LIMIT_CONFIG_EMAIL_SUBJECT") % (product_name)
                time_string = BusinessDate.now().strftime('%Y-%m-%d %H:%M:%S')
                content = _("IDS_EVFS_DOC_DOWNLOAD_LIMIT_CONFIG_EMAIL_CONTENT") % (time_string,
                                                                                   product_name)
                email_send_html_content(toEmailList, subject, content)
