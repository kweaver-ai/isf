#!/usr/bin/python3
# -*- coding:utf-8 -*-
"""
限速管理类
"""
import uuid
import time
import json
import threading
from collections import defaultdict
from src.common.db.connector import DBConnector
from src.modules.user_manage import UserManage
from src.modules.department_manage import DepartmentManage
from src.common.lib import (raise_exception, escape_key,
                            check_start_limit, generate_search_order_sql)
from EThriftException.ttypes import ncTException
from ShareMgnt.ttypes import (ncTShareMgntError, ncTLimitRateInfo, ncTLimitRateObject,
                              ncTLimitRateType, ncTLimitRateConfig, ncTLimitRateObjInfo)
from src.common.db.connector import ConnectorManager
from src.modules.handle_task_thread import (CallableTask, HandleTaskThread)
from src.modules.config_manage import ConfigManage
from src.common.sharemgnt_logger import ShareMgnt_Log

USER_OBJ = 1
DEPART_OBJ = 2

# 设定最大速度初始值，（实际速度配置不超过5位数）

MAX_UPLOAD_RATE = 100000
MAX_DOWNLOAD_RATE = 100000

threadLock = threading.Lock()
checkThread = None
WAIT_TIME = 10


class UserInfo:
    """
    @ todo: 限速对象为用户的结构体信息
    @ param uploadRate: 规则表中的上传限速值
    @ param downloadRate: 规则表中的下载限速值
    """
    def __init__(self):
        self.uploadRate = -1
        self.downloadRate = -1


class DeptInfo:
    """
    @ todo: 限速对象为部门的结构体信息
    @ param parentId: 父部门结构体
    @ param subDeptIds: 子部门结构体列表
    @ param uploadUser: 正在上传的用户人数
    @ param uploadRate: 规则表中的上传限速值
    @ param uploadRealRate: 实际该部门中用户被分配的上传速度
    @ param downloadUser: 正在下载的用户人数
    @ param downloadRate: 规则表中的下载限速值
    @ param downloadRealRate: 实际该部门中用户被分配的下载速度
    """
    def __init__(self):
        self.parentId = "-1"
        self.subDeptIds = []
        self.uploadUser = 0
        self.uploadRate = -1
        self.uploadRealRate = -1
        self.downloadUser = 0
        self.downloadRate = -1
        self.downloadRealRate = -1


class LimitRateManage(DBConnector):
    def __init__(self):
        """
        """
        self.user_manage = UserManage()
        self.dept_manage = DepartmentManage()
        self.handle_task_thread = HandleTaskThread()
        self.config_manage = ConfigManage()

    def __check_obj_limit_rate_exists(self, objInfo, objType, limitType, limitId):
        """
        检查该对象是否已存在同类型的限速规则
        """
        sql = """
        SELECT COUNT(*) AS cnt FROM `t_limit_rate` WHERE `f_obj_id` = '%s
        AND `f_limit_type` = %s
        AND `f_id` != %s
        """
        count = self.r_db.one(sql, objInfo.objectId, limitType, limitId)['cnt']

        if count != 0:
            if limitType == ncTLimitRateType.LIMIT_USER:
                if objType == USER_OBJ:
                    raise_exception(exp_msg=_("IDS_LIMIT_OBJ_EXISTS") % objInfo.objectName,
                                    exp_num=ncTShareMgntError.NCT_LIMIT_USER_EXIST)
                elif objType == DEPART_OBJ:
                    raise_exception(exp_msg=_("IDS_LIMIT_OBJ_EXISTS") % objInfo.objectName,
                                    exp_num=ncTShareMgntError.NCT_LIMIT_DEPART_EXIST)
            else:
                if objType == USER_OBJ:
                    raise_exception(exp_msg=_("IDS_LIMIT_USER_EXISTS"),
                                    exp_num=ncTShareMgntError.NCT_LIMIT_USER_EXIST)
                elif objType == DEPART_OBJ:
                    raise_exception(exp_msg=_("IDS_LIMIT_DEPT_EXISTS"),
                                    exp_num=ncTShareMgntError.NCT_LIMIT_DEPART_EXIST)

    def __check_user_list(self, userList, limitType, limitId):
        """
        检查用户列表中用户是否存在
        Args:
            userList: list 用户列表
            limitType: 0 - 用户级别的限速; 1 - 用户组级别的限速
            limitId: string, ""为新增规则; 非空为编辑规则
        """
        if not userList:
            return
        for user_obj in userList:
            self.user_manage.check_user_exists(user_obj.objectId, True)
            self.__check_obj_limit_rate_exists(user_obj, USER_OBJ, limitType, limitId)

    def __check_dept_list(self, deptList, limitType, limitId):
        """
        检查部门列表中部门是否存在
        """
        if not deptList:
            return

        for dept_obj in deptList:
            self.dept_manage.check_depart_exists(dept_obj.objectId, True)
            self.__check_obj_limit_rate_exists(dept_obj, DEPART_OBJ, limitType, limitId)

    def __check_limit_rate_value(self, uploadRate, downloadRate):
        """
        检查速度值是否有效
        """
        if uploadRate == downloadRate == -1:
            raise_exception(exp_msg=_("IDS_AT_LEAST_SET_ONE_SPEED"),
                            exp_num=ncTShareMgntError.NCT_AT_LEAST_SET_ONE_SPEED)

        if uploadRate is None or \
            (uploadRate != -1 and (uploadRate < 200 or uploadRate > 99999)) or \
            downloadRate is None or \
            (downloadRate != -1 and (downloadRate <= 0 or downloadRate > 99999)):
                raise_exception(exp_msg=_("IDS_INVALID_LIMIT_RATE_VALUES"),
                                exp_num=ncTShareMgntError.NCT_INVALID_LIMIT_RATE_VALUES)

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

    def add(self, limitRateInfo):
        """
        添加一条限速配置
        """
        # 检查限速类型是否合法
        if limitRateInfo.limitType is None or limitRateInfo.limitType < ncTLimitRateType.LIMIT_USER or  \
            limitRateInfo.limitType > ncTLimitRateType.LIMIT_USER_GROUP:
                raise_exception(exp_msg=_("IDS_INVALID_LIMIT_RATE_TYPE"),
                                exp_num=ncTShareMgntError.NCT_INVALID_LIMIT_RATE_TYPE)

        # 检查限速对象是否配置
        if not limitRateInfo.userInfos and not limitRateInfo.depInfos:
            raise_exception(exp_msg=_("IDS_LIMIT_RATE_OBJECT_NOT_SET"),
                            exp_num=ncTShareMgntError.NCT_LIMIT_RATE_OBJECT_NOT_SET)

        # 检查用户列表中用户是否存在
        if limitRateInfo.userInfos:
            limitRateInfo.userInfos = self.__remove_duplicate_obj(limitRateInfo.userInfos)
            self.__check_user_list(limitRateInfo.userInfos, limitRateInfo.limitType, "")

        # 检查部门列表中部门是否存在
        if limitRateInfo.depInfos:
            limitRateInfo.depInfos = self.__remove_duplicate_obj(limitRateInfo.depInfos)
            self.__check_dept_list(limitRateInfo.depInfos, limitRateInfo.limitType, "")

        # 检查用户组类型的限速对象个数是否合法
        if limitRateInfo.limitType == ncTLimitRateType.LIMIT_USER_GROUP and \
            (len(limitRateInfo.userInfos) + len(limitRateInfo.depInfos)) != 1:
                raise_exception(exp_msg=_("IDS_ONLY_ONE_LIMIT_RATE_OBJECT"),
                            exp_num=ncTShareMgntError.NCT_ONLY_ONE_LIMIT_RATE_OBJECT)

        # 检查限速值
        self.__check_limit_rate_value(limitRateInfo.uploadRate, limitRateInfo.downloadRate)

        # 生成一条唯一id
        limitRateInfo.id = str(uuid.uuid1())

        # 保存数据到数据库
        self.add_limit_rate_info_to_db(limitRateInfo)

        # 用户级别限速模式下增加任务更新nginx限速值
        if limitRateInfo.limitType == ncTLimitRateType.LIMIT_USER:
            self.add_update_user_limit_rate_task()

        return limitRateInfo.id

    def add_limit_rate_info_to_db(self, limitRateInfo):
        """
        保存数据到数据库
        """
        # 使用事务插入数据
        conn = ConnectorManager.get_db_conn()
        cursor = conn.cursor()

        insert_sql = """
        INSERT INTO `t_limit_rate`
        (`f_id`, `f_obj_id`, `f_obj_type`, `f_limit_type`, `f_upload_rate`, `f_download_rate`)
        VALUES(%s, %s, %s, %s, %s, %s)
        """

        try:
            # 插入用户列表配置
            if limitRateInfo.userInfos:
                for userInfo in limitRateInfo.userInfos:
                    cursor.execute(insert_sql, (limitRateInfo.id, userInfo.objectId, USER_OBJ, limitRateInfo.limitType,
                                   limitRateInfo.uploadRate, limitRateInfo.downloadRate))

            # 插入部门列表配置
            if limitRateInfo.depInfos:
                for depInfo in limitRateInfo.depInfos:
                    cursor.execute(insert_sql, (limitRateInfo.id, depInfo.objectId, DEPART_OBJ, limitRateInfo.limitType,
                                   limitRateInfo.uploadRate, limitRateInfo.downloadRate))

            conn.commit()

        except Exception as e:
            conn.rollback()
            raise_exception(exp_msg=str(e),
                            exp_num=ncTShareMgntError.NCT_DB_OPERATE_FAILED)
        finally:
            cursor.close()
            conn.close()

    def edit(self, limitRateInfo):
        """
        编辑一条限速配置
        """
        # 检查限速类型是否合法
        if limitRateInfo.limitType is None or limitRateInfo.limitType < ncTLimitRateType.LIMIT_USER or  \
            limitRateInfo.limitType > ncTLimitRateType.LIMIT_USER_GROUP:
                raise_exception(exp_msg=_("IDS_INVALID_LIMIT_RATE_TYPE"),
                                exp_num=ncTShareMgntError.NCT_INVALID_LIMIT_RATE_TYPE)

        # 检查限速对象是否配置
        if not limitRateInfo.userInfos and not limitRateInfo.depInfos:
            raise_exception(exp_msg=_("IDS_LIMIT_RATE_OBJECT_NOT_SET"),
                            exp_num=ncTShareMgntError.NCT_LIMIT_RATE_OBJECT_NOT_SET)

        # 检查该条限速规则是否存在
        if not self.__check_limit_rate_exist(limitRateInfo.id):
            raise_exception(exp_msg=_("IDS_LIMIT_RATE_NOT_EXIST"),
                            exp_num=ncTShareMgntError.NCT_LIMIT_RATE_NOT_EXIST)

        # 检查用户列表中用户是否存在
        if limitRateInfo.userInfos:
            limitRateInfo.userInfos = self.__remove_duplicate_obj(limitRateInfo.userInfos)
            self.__check_user_list(limitRateInfo.userInfos, limitRateInfo.limitType, limitRateInfo.id)

        # 检查部门列表中部门是否存在
        if limitRateInfo.depInfos:
            limitRateInfo.depInfos = self.__remove_duplicate_obj(limitRateInfo.depInfos)
            self.__check_dept_list(limitRateInfo.depInfos, limitRateInfo.limitType, limitRateInfo.id)

        # 检查用户组类型的限速对象个数是否合法
        if limitRateInfo.limitType == ncTLimitRateType.LIMIT_USER_GROUP and \
            (len(limitRateInfo.userInfos) + len(limitRateInfo.depInfos)) != 1:
                raise_exception(exp_msg=_("IDS_ONLY_ONE_LIMIT_RATE_OBJECT"),
                            exp_num=ncTShareMgntError.NCT_ONLY_ONE_LIMIT_RATE_OBJECT)

        # 检查限速值
        self.__check_limit_rate_value(limitRateInfo.uploadRate, limitRateInfo.downloadRate)

        # 更新数据库
        self.update_limit_rate_info_to_db(limitRateInfo)

        # 用户级别限速模式下增加任务更新nginx限速值
        if limitRateInfo.limitType == ncTLimitRateType.LIMIT_USER:
            self.add_update_user_limit_rate_task()

    def get(self, start, limit, limitType):
        """
        获取已配置的限速设置
        """
        limit_statement = check_start_limit(start, limit)

        query_sql = """
        SELECT DISTINCT `f_id`
        FROM `t_limit_rate`
        WHERE `f_limit_type` = %s
        ORDER BY `f_id`
        {0}
        """.format(limit_statement)
        results = self.r_db.all(query_sql, limitType)

        limit_rate_infos = []
        for res in results:
            limit_rate_info = self.get_limit_rate_info_by_id(res['f_id'])
            if limit_rate_info:
                limit_rate_info.limitType = limitType
                limit_rate_infos.append(limit_rate_info)

        return limit_rate_infos

    def get_cnt(self, limitType):
        """
        获取已配置的限速条数
        """
        sql = """
        SELECT COUNT(DISTINCT `f_id`) AS cnt
        FROM `t_limit_rate`
        WHERE `f_limit_type` = %s
        """
        count = self.r_db.one(sql, limitType)
        return count["cnt"]

    def get_limit_rate_info_by_id(self, limit_rate_id):
        """
        根据id获取限速信息
        """
        query_sql = """
        SELECT *
        FROM `t_limit_rate`
        WHERE `f_id` = %s
        """
        results = self.r_db.all(query_sql, limit_rate_id)

        need_delete_obj_id = []

        if results:
            limit_info = ncTLimitRateInfo()
            limit_info.id = limit_rate_id
            limit_info.uploadRate = results[0]['f_upload_rate']
            limit_info.downloadRate = results[0]['f_download_rate']
            limit_info.userInfos = []
            limit_info.depInfos = []

            for res in results:
                obj_id = res['f_obj_id']
                obj_type = res['f_obj_type']

                if obj_type == USER_OBJ:
                    try:
                        obj_info = self.user_manage.get_user_by_id(obj_id)

                        user_info = ncTLimitRateObject()
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

                        dept_info = ncTLimitRateObject()
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
        删除一条限速配置信息
        """
        for obj_id in deleteObjIdList:
            delete_sql = """
            DELETE FROM `t_limit_rate`
            WHERE `f_obj_id` = %s
            """
            self.w_db.query(delete_sql, obj_id)

    def update_limit_rate_info_to_db(self, limitRateInfo):
        """
        更新数据库中数据
        """
        # 获取原配置中的限速值
        select_sql = """
        SELECT `f_upload_rate`, `f_download_rate`
        FROM `t_limit_rate`
        WHERE `f_id` = %s
        """
        result = self.r_db.one(select_sql, limitRateInfo.id)
        upload_rate = result['f_upload_rate']
        download_rate = result['f_download_rate']

        # 获取已经配置的用户id与部门id
        user_ids = set()
        dept_ids = set()
        select_sql = """
        SELECT `f_obj_id`, `f_obj_type`
        FROM `t_limit_rate`
        WHERE `f_id` = %s
        AND `f_limit_type` = %s
        """
        result = self.r_db.all(select_sql, limitRateInfo.id, limitRateInfo.limitType)
        for res in result:
            if res['f_obj_type'] == USER_OBJ:
                user_ids.add(res['f_obj_id'])
            if res['f_obj_type'] == DEPART_OBJ:
                dept_ids.add(res['f_obj_id'])

        # 获取新配置的用户对象id
        new_user_ids = set()
        if limitRateInfo.userInfos:
            new_user_ids = set([userObj.objectId for userObj in limitRateInfo.userInfos])

        # 获取新配置的部门对象id
        new_dept_ids = set()
        if limitRateInfo.depInfos:
            new_dept_ids = set([deptObj.objectId for deptObj in limitRateInfo.depInfos])

        # 需要删除的对象id
        need_delete_obj_id_list = (user_ids - new_user_ids) | (dept_ids - new_dept_ids)

        # 需要更新的对象id
        need_update_obj_id_list = (user_ids & new_user_ids) | (dept_ids & new_dept_ids)

        # 需要增加的用户对象id
        need_add_user_ids = new_user_ids - user_ids

        # 需要增加的部门对象id
        need_add_dept_ids = new_dept_ids - dept_ids

        # 使用事务更新数据
        conn = ConnectorManager.get_db_conn()
        cursor = conn.cursor()

        delete_sql = """
        DELETE
        FROM `t_limit_rate`
        WHERE `f_id` = '{0}' AND `f_obj_id` = '{1}'
        """

        update_sql = """
        UPDATE `t_limit_rate`
        SET f_upload_rate = {0}, f_download_rate = {1}
        WHERE `f_id` = '{2}' AND `f_obj_id` = '{3}'
        """

        insert_sql = """
        INSERT INTO `t_limit_rate`
        (`f_id`, `f_obj_id`, `f_obj_type`, `f_limit_type`, `f_upload_rate`, `f_download_rate`)
        VALUES('{0}', '{1}', {2}, {3}, {4}, {5})
        """
        try:
            # 删除已存在的配置项
            for obj_id in need_delete_obj_id_list:
                cursor.execute(delete_sql.format(limitRateInfo.id, obj_id))

            # 更新原有配置项的限速值
            if (upload_rate != limitRateInfo.uploadRate) or (download_rate != limitRateInfo.downloadRate):
                for obj_id in need_update_obj_id_list:
                    cursor.execute(update_sql.format(limitRateInfo.uploadRate, limitRateInfo.downloadRate,
                                limitRateInfo.id, obj_id))

            # 插入用户列表配置
            for user_id in need_add_user_ids:
                cursor.execute(insert_sql.format(limitRateInfo.id, user_id, USER_OBJ, limitRateInfo.limitType,
                               limitRateInfo.uploadRate, limitRateInfo.downloadRate))

            # 插入部门列表配置
            for dept_id in need_add_dept_ids:
                cursor.execute(insert_sql.format(limitRateInfo.id, dept_id, DEPART_OBJ, limitRateInfo.limitType,
                               limitRateInfo.uploadRate, limitRateInfo.downloadRate))

            conn.commit()

        except Exception as e:
            conn.rollback()
            raise_exception(exp_msg=str(e),
                            exp_num=ncTShareMgntError.NCT_DB_OPERATE_FAILED)
        finally:
            cursor.close()
            conn.close()

    def __check_limit_rate_exist(self, rateId):
        """
        检查限速是否配置
        """
        sql = """
        SELECT count(*) as cnt
        FROM `t_limit_rate`
        WHERE `f_id` = %s
        """
        result = self.r_db.one(sql, rateId)
        return True if result['cnt'] != 0 else False

    def delete(self, deleteId, limitType):
        """
        删除一条限速配置信息
        """
        delete_sql = """
        DELETE FROM `t_limit_rate`
        WHERE `f_id` = %s
        AND `f_limit_type`= %s
        """
        affect_row = self.w_db.query(delete_sql, self.w_db.escape(deleteId), limitType)
        if not affect_row:
            raise_exception(exp_msg=_("IDS_LIMIT_RATE_NOT_EXIST"),
                            exp_num=ncTShareMgntError.NCT_LIMIT_RATE_NOT_EXIST)

        # 增加任务更新nginx限速值
        if limitType == ncTLimitRateType.LIMIT_USER:
            self.add_update_user_limit_rate_task()

    def search(self, searchKey, start, limit, limitType):
        """
        搜索限速配置信息
        """
        # 检查参数
        limit_statement = check_start_limit(start, limit)

        # 增加搜索排序子句
        order_by_str = generate_search_order_sql(['MIN(`t_user`.`f_display_name`)',
                                                  'MIN(`t_department`.`f_name`)'])

        search_sql = """
        SELECT `t_limit_rate`.`f_id`,
        MIN(`t_user`.`f_display_name`),
        MIN(`t_department`.`f_name`)
        FROM `t_limit_rate`
        LEFT JOIN `t_user`
        ON `t_limit_rate`.`f_obj_id` = `t_user`.`f_user_id`
        LEFT JOIN `t_department`
        ON `t_limit_rate`.`f_obj_id` = `t_department`.`f_department_id`
        WHERE ((`t_user`.`f_display_name` LIKE %s OR `t_department`.`f_name` LIKE %s)
        AND `t_limit_rate`.`f_limit_type` = {2})
        GROUP BY `t_limit_rate`.`f_id`
        ORDER BY {0}
        {1}
        """.format(order_by_str, limit_statement, limitType)
        esckey = "%%%s%%" % escape_key(searchKey)
        results = self.r_db.all(search_sql, esckey, esckey,
                                escape_key(searchKey), escape_key(searchKey), esckey, esckey)

        limit_rate_infos = []
        for res in results:
            limit_rate_info = self.get_limit_rate_info_by_id(res['f_id'])
            if limit_rate_info:
                limit_rate_info.limitType = limitType
                limit_rate_infos.append(limit_rate_info)

        return limit_rate_infos

    def search_cnt(self, searchKey, limitType):
        """
        搜索限速配置信息条数
        """
        search_sql = """
        SELECT COUNT(DISTINCT `t_limit_rate`.`f_id`) as cnt
        FROM `t_limit_rate`
        LEFT JOIN `t_user`
        ON `t_limit_rate`.`f_obj_id` = `t_user`.`f_user_id`
        LEFT JOIN `t_department`
        ON `t_limit_rate`.`f_obj_id` = `t_department`.`f_department_id`
        WHERE ((`t_user`.`f_display_name` LIKE %s OR `t_department`.`f_name` LIKE %s)
        AND `t_limit_rate`.`f_limit_type` = %s)
        """
        esckey = "%%{0}%%".format(escape_key(searchKey))
        result = self.r_db.one(search_sql, esckey, esckey, limitType)
        return result['cnt']

    def update_nginx_limit_rate(self, args):
        """
        更新nginx表中限速值配置
        """
        select_sql = """
        SELECT `f_obj_id`, `f_upload_rate`, `f_download_rate`
        FROM `t_limit_rate`
        WHERE `f_limit_type` = %s
        """
        results = self.r_db.all(select_sql, 0)

        # 获取用户级别限速的所有规则信息
        limit_obj_map = {}
        for res in results:
            if res["f_obj_id"] in limit_obj_map:
                if 0 < res["f_upload_rate"] < limit_obj_map[res["f_obj_id"]].uploadRate:
                    limit_obj_map[res["f_obj_od"]].uploadRate = res["f_upload_rate"]
                if 0 < res["f_download_rate"] < limit_obj_map[res["f_obj_id"].downloadRate]:
                    limit_obj_map[res["f_obj_od"]].downloadRate = res["f_download_rate"]
            else:
                limit_obj_map[res["f_obj_id"]] = obj = UserInfo()
                obj.uploadRate = res["f_upload_rate"]
                obj.downloadRate = res["f_download_rate"]

        # 从nginx限速表中获取需要更新的用户id
        select_sql = """
        SELECT f_userid, f_upload_rate, f_download_rate
        FROM t_nginx_user_rate
        """
        results = self.r_db.all(select_sql)

        # 部门映射关系
        dept_relation_map = {}
        for res in results:
            upload_rate = MAX_UPLOAD_RATE
            download_rate = MAX_DOWNLOAD_RATE

            if res["f_userid"] in limit_obj_map:
                if 0 < limit_obj_map[res["f_userid"]].uploadRate < upload_rate:
                    upload_rate = limit_obj_map[res["f_userid"]].uploadRate
                if 0 < limit_obj_map[res["f_userid"]].downloadRate < download_rate:
                    download_rate = limit_obj_map[res["f_userid"]].downloadRate
            else:
                # 获取用户的直属部门id
                direct_dept_ids = self.user_manage.get_belong_depart_id(res["f_userid"])

                # 获取用户上一层有限速的部门id
                limit_part_dept_id_set = set()
                for direct_id in direct_dept_ids:
                    part_id = direct_id
                    limit_part_id = "-1"
                    need_add_relation = []

                    while True:
                        if not part_id:
                            break
                        elif part_id in dept_relation_map:
                            limit_part_id = dept_relation_map[part_id]
                            break
                        elif part_id in limit_obj_map:
                            limit_part_id = part_id
                            break
                        else:
                            need_add_relation.append(part_id)

                            part_id = self.dept_manage.get_parent_id(part_id)

                    for dept_id in need_add_relation:
                        dept_relation_map[dept_id] = limit_part_id

                    if limit_part_id != "-1":
                        limit_part_dept_id_set.add(limit_part_id)

                for dept_id in limit_part_dept_id_set:
                    if 0 < limit_obj_map[dept_id].uploadRate < upload_rate:
                        upload_rate = limit_obj_map[dept_id].uploadRate
                    if 0 < limit_obj_map[dept_id].downloadRate < download_rate:
                        download_rate = limit_obj_map[dept_id].downloadRate

            if upload_rate == MAX_UPLOAD_RATE or upload_rate < 0:
                upload_rate = 0

            if download_rate == MAX_DOWNLOAD_RATE or download_rate < 0:
                download_rate = 0

            if (res["f_upload_rate"] != upload_rate) or (res["f_download_rate"] != download_rate):
                self.update_nginx_limit_rate_db(res['f_userid'], upload_rate, download_rate)

    def update_nginx_limit_rate_db(self, userid, upload_rate, download_rate):
        """
        更新nginx表中限速值配置
        """
        # 从nginx限速表中获取需要更新的用户id
        update_sql = """
        UPDATE t_nginx_user_rate
        SET f_upload_rate = %s, f_download_rate = %s
        WHERE f_userid = %s
        """
        self.w_db.all(update_sql, upload_rate, download_rate, userid)

    def add_update_user_limit_rate_task(self):
        """
        添加更新用户限速值任务
        """
        task = CallableTask()
        task.module_name = "limit_rate_manage"
        task.function_name = "update_nginx_limit_rate"
        self.handle_task_thread.add(task)

    def get_limit_rate_config(self):
        """
        获取限速配置信息
        """
        config_dict = json.loads(self.config_manage.get_config("limit_rate_config"))
        config = ncTLimitRateConfig()
        config.isEnabled = config_dict["isEnabled"]
        config.limitType = config_dict["limitType"]
        return config

    def start_update_user_limit_rate_thread(self):
        """
        启动更新用户限速值线程
        """
        global checkThread
        with threadLock:
            if checkThread is None:
                checkThread = UpdateUserLimitRateThread()
                checkThread.daemon = True
                checkThread.start()

    def update_parent_department_ids(self):
        """
        更新部门和用户有规则的父部门id
        """
        # step 1: 重构拷贝规则表, 获取部门有规则的父部门id
        # 清空拷贝规则表
        delete_sql = """
        TRUNCATE TABLE t_copy_limit_rate
        """
        self.w_db.query(delete_sql)

        # 拷贝当前规则表中所有用户组级别的限速规则
        select_sql = """
        SELECT `f_obj_id`, `f_obj_type`, `f_upload_rate`, `f_download_rate`
        FROM `t_limit_rate`
        WHERE `f_limit_type` = %s
        """

        insert_sql = """
        INSERT INTO t_copy_limit_rate
        (f_obj_id, f_parent_id, f_obj_type, f_upload_rate, f_download_rate)
        VALUES {0}
        """
        insert_clause = ", ('{0}', '{1}', {2}, {3}, {4})"

        results = self.r_db.all(select_sql, 1)

        # 获取用户组级别限速的所有部门id
        limit_dept_id_set = set()
        for res in results:
            limit_dept_id_set.add(res["f_obj_id"])

        # 部门映射关系 {"dept_id": "part_id"}
        # part_id: dept_id 有规则的父部门id
        dept_relation_map = {}

        format_str = ""
        for res in results:
            if res["f_obj_type"] == USER_OBJ:
                # 添加限速对象为用户的记录
                format_str += insert_clause.format(res["f_obj_id"], "", res["f_obj_type"],
                                            res["f_upload_rate"], res["f_download_rate"])
                continue

            part_id = res["f_obj_id"]
            limit_part_id = "-1"
            need_add_relation = []

            while True:
                part_id = self.dept_manage.get_parent_id(part_id)

                if not part_id:
                    break
                elif part_id in dept_relation_map:
                    limit_part_id = dept_relation_map[part_id]
                    break
                elif part_id in limit_dept_id_set:
                    limit_part_id = part_id
                    break
                else:
                    need_add_relation.append(part_id)

            # 添加新的映射关系
            for dept_id in need_add_relation:
                dept_relation_map[dept_id] = limit_part_id

            # 添加限速对象为部门的记录
            format_str += insert_clause.format(res["f_obj_id"], limit_part_id, res["f_obj_type"],
                                        res["f_upload_rate"], res["f_download_rate"])

        if format_str:
            self.w_db.query(insert_sql.format(format_str[1:]))

        # step 2: 更新用户连接表中的父部门id
        select_sql = """
        SELECT f_userid, f_parent_deptids
        FROM t_nginx_user_rate
        """

        update_sql = """
        UPDATE t_nginx_user_rate
        SET f_parent_deptids = '{0}'
        WHERE f_userid = '{1}'
        """

        results = self.r_db.all(select_sql)

        for res in results:
            # 获取用户的直属部门
            direct_dept_ids = self.user_manage.get_belong_depart_id(res["f_userid"])

            if (not direct_dept_ids) and (res["f_parent_deptids"] != "-1"):
                self.w_db.query(update_sql.format("-1", res["f_userid"]))
                continue

            # 最终保存的有规则的父部门id
            parent_deptids_list = set()
            for direct_id in direct_dept_ids:
                part_id = direct_id
                limit_part_id = "-1"
                need_add_relation = []

                while True:
                    if not part_id:
                        break
                    elif part_id in dept_relation_map:
                        limit_part_id = dept_relation_map[part_id]
                        break
                    elif part_id in limit_dept_id_set:
                        limit_part_id = part_id
                        break
                    else:
                        need_add_relation.append(part_id)
                        part_id = self.dept_manage.get_parent_id(part_id)

                # 添加新的映射关系
                for dept_id in need_add_relation:
                    dept_relation_map[dept_id] = limit_part_id

                if limit_part_id != "-1":
                    parent_deptids_list.add(self.w_db.escape(limit_part_id))

            # 构造最终的父部门id
            partids_text = ",".join(sorted(parent_deptids_list)) if len(parent_deptids_list) > 0 else "-1"

            # 更新用户有规则的父部门id
            if res["f_parent_deptids"] != partids_text:
                self.w_db.query(update_sql.format(partids_text, res["f_userid"]))

    def update_user_limit_rate(self):
        """
        计算各部门实际能被分配到的速度
        并更新用户的最终连接速度
        """
        # step 1: 构造简易组织架构并初始化
        dept_map = {}
        user_map = {}
        top_dept_list = []

        select_sql = """
        SELECT * FROM t_copy_limit_rate
        """
        results = self.r_db.all(select_sql)

        for res in results:
            if res["f_obj_type"] == USER_OBJ:
                user_map[res["f_obj_id"]] = obj = UserInfo()
                obj.uploadRate = res["f_upload_rate"]
                obj.downloadRate = res["f_download_rate"]
            else:
                # 添加顶层部门id, 更新有规则的子部门id
                if res["f_parent_id"] == "-1":
                    top_dept_list.append(res["f_obj_id"])
                elif res["f_parent_id"] in dept_map:
                    dept_map[res["f_parent_id"]].subDeptIds.append(res["f_obj_id"])
                else:
                    dept_map[res["f_parent_id"]] = obj = DeptInfo()
                    obj.subDeptIds = [res["f_obj_id"]]

                # 添加部门映射关系
                if res["f_obj_id"] in dept_map:
                    obj = dept_map[res["f_obj_id"]]
                else:
                    dept_map[res["f_obj_id"]] = obj = DeptInfo()

                obj.uploadRate = res["f_upload_rate"]
                obj.downloadRate = res["f_download_rate"]
                obj.parentId = res["f_parent_id"]

        # step 2: 计算当前连接人数
        select_sql = """
        SELECT `f_userid`, `f_parent_deptids`, `f_upload_req_cnt`,
        `f_download_req_cnt`, `f_upload_rate`, `f_download_rate`
        FROM `t_nginx_user_rate`
        WHERE `f_upload_req_cnt` > %s OR `f_download_req_cnt` > %s
        """
        results = self.r_db.all(select_sql, 0, 0)

        for res in results:
            if not res["f_parent_deptids"] or res["f_parent_deptids"] == "-1":
                continue

            parent_deptids_list = res["f_parent_deptids"].split(",")

            # 需要增加人数的所有父部门id
            need_add_user_count = set()
            for dept_id in parent_deptids_list:
                part_id = dept_id
                while part_id != "-1":
                    if part_id in dept_map:
                        need_add_user_count.add(part_id)
                        part_id = dept_map[part_id].parentId
                    else:
                        part_id = "-1"

            for dept_id in need_add_user_count:
                if res["f_upload_req_cnt"] > 0:
                    dept_map[dept_id].uploadUser += 1

                if res["f_download_req_cnt"] > 0:
                    dept_map[dept_id].downloadUser += 1

        # step 3: 从顶层开始计算部门速度
        calc_list = top_dept_list
        while calc_list:
            new_calc_list = []
            for dept_id in calc_list:
                obj = dept_map[dept_id]
                new_calc_list += obj.subDeptIds

                # 计算部门实际能分配的速度值
                if obj.uploadUser > 0 and obj.uploadRate != -1:
                    rate = obj.uploadRate / obj.uploadUser
                    obj.uploadRealRate = rate if rate > 0 else rate + 1
                if obj.downloadUser > 0 and obj.downloadRate != -1:
                    rate = obj.downloadRate / obj.downloadUser
                    obj.downloadRealRate = rate if rate > 0 else rate + 1

                if obj.parentId != "-1":
                    part_upload_rate = dept_map[obj.parentId].uploadRealRate
                    part_download_rate = dept_map[obj.parentId].downloadRealRate

                    if part_upload_rate != -1 and (obj.uploadRealRate == -1 or  \
                        part_upload_rate < obj.uploadRealRate):
                            obj.uploadRealRate = part_upload_rate
                    if part_download_rate != -1 and (obj.downloadRealRate == -1 or  \
                        part_download_rate < obj.downloadRealRate):
                            obj.downloadRealRate = part_download_rate

            calc_list = new_calc_list

        # step 4: 获取每个用户实际能分配的速度
        upload_rate_map = defaultdict(list)
        download_rate_map = defaultdict(list)

        for res in results:
            final_upload_rate = 0
            final_download_rate = 0

            if res["f_userid"] in user_map:
                if user_map[res["f_userid"]].uploadRate != -1:
                    final_upload_rate = user_map[res["f_userid"]].uploadRate
                if user_map[res["f_userid"]].downloadRate != -1:
                    final_download_rate = user_map[res["f_userid"]].downloadRate

            parent_deptids_list = res["f_parent_deptids"].split(",")
            for dept_id in parent_deptids_list:
                if dept_id in dept_map:
                    if dept_map[dept_id].uploadRealRate != -1 and (final_upload_rate == 0 or \
                        dept_map[dept_id].uploadRealRate < final_upload_rate):
                            final_upload_rate = dept_map[dept_id].uploadRealRate

                    if dept_map[dept_id].downloadRealRate != -1 and (final_download_rate == 0 or \
                        dept_map[dept_id].downloadRealRate < final_download_rate):
                            final_download_rate = dept_map[dept_id].downloadRealRate

            # 上传/下载速度不清零
            if res["f_upload_req_cnt"] == 0:
                final_upload_rate = res["f_upload_rate"]
            if res["f_download_req_cnt"] == 0:
                final_download_rate = res["f_download_rate"]

            userid = "'{0}'".format(self.w_db.escape(res["f_userid"]))

            upload_rate_map[final_upload_rate].append(userid)
            download_rate_map[final_download_rate].append(userid)

        # step 5: 更新用户速度
        update_upload_rate_sql = """
        UPDATE t_nginx_user_rate
        SET f_upload_rate = {0}
        WHERE f_userid in ({1})
        """

        update_download_rate_sql = """
        UPDATE t_nginx_user_rate
        SET f_download_rate = {0}
        WHERE f_userid in ({1})
        """

        # 更新上传速度
        for rate in upload_rate_map:
            self.w_db.query(update_upload_rate_sql.format(rate, ",".join(upload_rate_map[rate])))

        # 更新下载速度
        for rate in download_rate_map:
            self.w_db.query(update_download_rate_sql.format(rate, ",".join(download_rate_map[rate])))

    def get_exist_object_info(self, userInfos, depInfos, limitType, limitId):
        """
        获取已存在其他限速规则的对象信息
        """
        obj_map = {}
        where_clause = []
        for user_info in userInfos:
            obj_map[user_info.objectId] = user_info
            where_clause.append("'{0}'".format(self.w_db.escape(user_info.objectId)))

        for dept_info in depInfos:
            obj_map[dept_info.objectId] = dept_info
            where_clause.append("'{0}'".format(self.w_db.escape(dept_info.objectId)))

        select_sql = """
        SELECT `f_obj_id`, `f_obj_type`
        FROM `t_limit_rate`
        WHERE `f_obj_id` in ({0})
        AND `f_limit_type` = %s
        AND `f_id` != %s
        """
        where_clause = ",".join(where_clause) if len(where_clause) > 0 else "''"
        results = self.r_db.all(select_sql.format(where_clause), limitType, self.w_db.escape(limitId))

        exist_obj_info = ncTLimitRateObjInfo()
        exist_obj_info.userInfos = []
        exist_obj_info.depInfos = []
        for res in results:
            if res["f_obj_type"] == USER_OBJ:
                exist_obj_info.userInfos.append(obj_map[res["f_obj_id"]])
            elif res["f_obj_type"] == DEPART_OBJ:
                exist_obj_info.depInfos.append(obj_map[res["f_obj_id"]])

        return exist_obj_info


class UpdateUserLimitRateThread(threading.Thread):
    """
    更新用户限速值线程
    """
    def __init__(self):
        """
        初始化
        """
        super(UpdateUserLimitRateThread, self).__init__()
        self.terminate = False
        self.limit_rate_manage = LimitRateManage()

    def close(self):
        """
        关闭
        """
        self.terminate = True

    def run(self):
        """
        执行
        """
        ShareMgnt_Log("**************** update user limit rate thread start *****************")

        cnt = 0
        while True:
            if self.terminate:
                break

            try:
                if cnt % 18 == 0:
                    cnt = 0
                    self.limit_rate_manage.update_parent_department_ids()

                cnt += 1
                self.limit_rate_manage.update_user_limit_rate()
            except Exception as ex:
                ShareMgnt_Log("update user limit rate thread run error: %s", str(ex))

            time.sleep(WAIT_TIME)

        ShareMgnt_Log("**************** update user limit rate thread end *****************")
