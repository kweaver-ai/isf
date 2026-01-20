#!/usr/bin/python3
# -*- coding:utf-8 -*-
"""This is perm share manage class"""
import uuid
import json
from src.common.db.connector import DBConnector
from src.common.lib import (escape_key,
                            check_start_limit,
                            generate_search_order_sql,
                            raise_exception)
from src.modules.config_manage import ConfigManage
from src.modules.user_manage import UserManage
from src.modules.department_manage import DepartmentManage
from src.common.db.connector import ConnectorManager
from src.common.sharemgnt_logger import ShareMgnt_Log
from ShareMgnt.ttypes import (ncTShareMgntError,
                              ncTPermShareInfo,
                              ncTShareObjInfo)
from ShareMgnt.constants import (NCT_ALL_USER_GROUP,
                                 NCT_DIRECT_DEPARTMENT,
                                 NCT_DIRECT_ORGANIZATION)


SHARER_TYPE = 1
SCOPE_TYPE = 2
USER_OBJ = 1
DEPART_OBJ = 2


class PermShareManage(DBConnector):

    """
    Perm share manage
    """

    def __init__(self):
        """
        init
        """
        self.config_manage = ConfigManage()
        self.user_manage = UserManage()
        self.dept_manage = DepartmentManage()

    def check_strategy_info(self, strategy_info, b_strategy_id=False):
        """
        检查权限共享范围信息
        b_strategy_id 是否检查strategy id
        """
        if not strategy_info.sharerUsers and not strategy_info.sharerDepts:
            raise_exception(exp_msg=_("sharer is empty"),
                            exp_num=ncTShareMgntError.NCT_SHARER_IS_EMPTY)

        if not strategy_info.scopeUsers and not strategy_info.scopeDepts:
            raise_exception(exp_msg=_("share scope is empty"),
                            exp_num=ncTShareMgntError.NCT_SHARE_SCOPE_IS_EMPTY)

        notfound_params = []
        # 检查用户是否存在
        user_id_name_map = {}
        for each in (strategy_info.sharerUsers or []) + (strategy_info.scopeUsers or []):
            user_id_name_map[each.id or ""] = each.name or ""
        if len(user_id_name_map):
            for userid, name in list(user_id_name_map.items()):
                if not self.user_manage.check_user_exists(userid, False):
                    notfound_params.append({"id":userid, "name":name, "type":"user"})

        # 检查部门是否存在
        depart_id_name_map = {}
        for each in (strategy_info.sharerDepts or []) + (strategy_info.scopeDepts or []):
            depart_id_name_map[each.id or ""] = each.name or ""
        if len(depart_id_name_map):
            for departid, name in list(depart_id_name_map.items()):
                if not self.dept_manage.check_depart_exists(departid, True, False):
                    notfound_params.append({"id":departid, "name":name, "type":"department"})

        if len(notfound_params) > 0:
            raise_exception(exp_msg=_("IDS_OBJ_NOT_EXIST"),
                            exp_num=ncTShareMgntError.NCT_OBJ_NOT_EXIST,
                            exp_detail=json.dumps({"notfound_params": notfound_params}, ensure_ascii=False))

        if strategy_info.scopeUsers:
            for obj_info in strategy_info.scopeUsers:
                if not hasattr(obj_info, "parentId") or \
                   not obj_info.parentId or \
                   not self.dept_manage.check_depart_exists(obj_info.parentId, True, False):
                    raise_exception(exp_msg=_("depart not exists"),
                                    exp_num=ncTShareMgntError.NCT_DEPARTMENT_NOT_EXIST)

        if b_strategy_id:
            query_sql = """
            SELECT * FROM `t_perm_share_strategy`
            WHERE `f_strategy_id` = %s
            """
            result = self.r_db.one(query_sql, strategy_info.strategyId)
            if not result:
                raise_exception(exp_msg=_("perm share item not exists"),
                                exp_num=ncTShareMgntError.NCT_SHARE_STRATEGY_NOT_EXISTS)

    def get_user_obj_info(self, user_id, parent_id=None):
        """
        根据用户id获取用户对象信息
        """
        obj_info = ncTShareObjInfo()

        if user_id == NCT_ALL_USER_GROUP:
            obj_info.id = NCT_ALL_USER_GROUP
            obj_info.name = _("any user")
            return obj_info
        else:
            try:
                user_info = self.user_manage.get_user_by_id(user_id)
                obj_info.id = user_id
                obj_info.name = user_info.user.displayName

                # 获取部门信息
                if parent_id:
                    depart_info = self.dept_manage.get_department_info(parent_id, True)
                    obj_info.parentId = depart_info.departmentId
                    obj_info.parentName = depart_info.departmentName
                return obj_info
            except Exception as ex:
                ShareMgnt_Log(str(ex))

    def __get_dept_obj_info(self, dept_id):
        """
        根据部门id获取部门对象信息
        """
        obj_info = ncTShareObjInfo()
        if dept_id == NCT_DIRECT_DEPARTMENT:
            obj_info.id = NCT_DIRECT_DEPARTMENT
            obj_info.name = _("dept user belongs to")
            return obj_info
        elif dept_id == NCT_DIRECT_ORGANIZATION:
            obj_info.id = NCT_DIRECT_ORGANIZATION
            obj_info.name = _("organ user belongs to")
            return obj_info
        else:
            try:
                dept_info = self.dept_manage.get_department_info(dept_id, True)
                obj_info.id = dept_id
                obj_info.name = dept_info.departmentName
                return obj_info
            except Exception as ex:
                ShareMgnt_Log(str(ex))

    def check_user_direct_dept_scope(self):
        """
        检查是否开启所有用户的权限范围为直属部门
        """
        sql = """
        SELECT `f_status` FROM `t_perm_share_strategy`
        WHERE `f_strategy_id` = %s
        """

        res = self.r_db.one(sql, '-1')
        return True if res['f_status'] == 1 else False

    def check_user_direct_org_scope(self):
        """
        检查是否开启所有用户的权限范围为直属组织
        """
        sql = """
        SELECT `f_status` FROM `t_perm_share_strategy`
        WHERE `f_strategy_id` = %s
        """

        res = self.r_db.one(sql, '-2')
        return True if res['f_status'] == 1 else False

    def get_perm_scope_by_sharer_ids(self, user_id):
        """
        根据共享者id获取权限范围
        """
        scope_dict = {}

        sharer_ids = self.user_manage.get_all_path_dept_id(user_id)
        sharer_ids.add(user_id)

        # 根据共享者id获取共享策略id
        group_str = ["\'" + self.w_db.escape(s_id) + "\'" for s_id in sharer_ids]
        group_str = ",".join(group_str)

        sql = """
        SELECT distinct `f_strategy_id` FROM `t_perm_share_strategy`
        WHERE `f_obj_id` IN ({0}) AND `f_sharer_or_scope` = %s;
        """.format(group_str)

        results = self.r_db.all(sql, 1)

        if results:
            strategy_ids = [res['f_strategy_id'] for res in results]

            # 根据共享策略id获取相应的共享范围
            group_str = ["\'" + s_id + "\'" for s_id in strategy_ids]
            group_str = ",".join(group_str)
            sql = """
            SELECT `f_obj_id`,`f_parent_id`,`f_obj_type` FROM `t_perm_share_strategy`
            WHERE `f_strategy_id` IN ({0}) AND `f_sharer_or_scope` = %s;
            """.format(group_str)
            results = self.r_db.all(sql, 2)

            scope_dict = {}
            for res in results:
                obj_id = res['f_obj_id']
                obj_type = res['f_obj_type']
                parent_id = res['f_parent_id']

                # 检查范围对象用户和父部门的关系是否存在
                if (obj_type == USER_OBJ):
                    if not self.dept_manage.check_user_in_depart(obj_id, parent_id, False):
                        continue

                scope_dict[obj_id] = obj_type

        # 检查是否开启所有用户的权限范围为直属部门
        b_direct_ok = self.check_user_direct_dept_scope()

        # 如果所有用户的直属部门权限范围开启，则添加直属部门到权限范围
        if b_direct_ok:
            # 获取用户的所有直属父部门id
            direct_dept_ids = self.user_manage.get_belong_depart_id(user_id)
            for dept_id in direct_dept_ids:
                scope_dict[dept_id] = DEPART_OBJ

        # 检查是否开启所有用户的权限范围为直属组织
        b_direct_ok = self.check_user_direct_org_scope()

        # 如果所有用户的直属组织权限范围开启，则添加直属组织到权限范围
        if b_direct_ok:
            # 获取用户的所有直属组织id
            direct_ou_ids = self.dept_manage.get_ou_by_user_id(user_id)
            if direct_ou_ids:
                for ou_id in direct_ou_ids:
                    scope_dict[ou_id] = DEPART_OBJ

        return scope_dict

    def check_user_in_perm_scope(self, user_id, check_user_id):
        """
        检查某个用户是否在权限范围内
        """
        # 获取共享范围
        if not self.get_system_perm_share_status():
            return True

        perm_scope = self.get_perm_scope_by_sharer_ids(user_id)

        path_dept_ids = self.user_manage.get_all_path_dept_id(check_user_id)

        if check_user_id in perm_scope:
            return True
        else:
            for dept_id in path_dept_ids:
                if dept_id in perm_scope:
                    return True

        return False

    def set_system_perm_share_status(self, status):
        """
        设置系统权限共享范围状态：
        参数：
           status：
                True： 开启
                Flase：关闭
        """
        status = 1 if status else 0
        self.config_manage.set_config("perm_share_status", status)

    def get_system_perm_share_status(self):
        """
        获取系统权限共享范围状态：
        参数：
           status：
                True： 开启
                Flase：关闭
        """
        result = self.config_manage.get_config('perm_share_status')
        status = True if int(result) == 1 else False
        return status

    def add_perm_share_info(self, strategy_info):
        """
        增加一条权限共享策略信息
        """
        # 检查权限共享范围信息
        self.check_strategy_info(strategy_info)

        # 生成策略id，共享者组id，共享范围组id
        stg_id = str(uuid.uuid1())

        # 使用事务插入数据
        conn = ConnectorManager.get_db_conn()
        cursor = conn.cursor()

        insert_sql = """
        INSERT INTO `t_perm_share_strategy`
        (`f_strategy_id`, `f_obj_id`, `f_parent_id`, `f_obj_type`,
         `f_sharer_or_scope`, `f_status`)
        VALUES(%s, %s, %s, %s, %s, %s)
        """

        try:
            # 插入共享者组信息
            for sharer in strategy_info.sharerUsers:
                cursor.execute(insert_sql, (stg_id, sharer.id, "", USER_OBJ, SHARER_TYPE, 1))

            for sharer in strategy_info.sharerDepts:
                cursor.execute(insert_sql, (stg_id, sharer.id, "", DEPART_OBJ, SHARER_TYPE, 1))

            # 插入共享范围组信息
            for scope in strategy_info.scopeUsers:
                cursor.execute(insert_sql, (stg_id, scope.id, scope.parentId,
                                            USER_OBJ, SCOPE_TYPE, 1))

            for scope in strategy_info.scopeDepts:
                cursor.execute(insert_sql, (stg_id, scope.id, "", DEPART_OBJ, SCOPE_TYPE, 1))
            conn.commit()
        except Exception as e:
            conn.rollback()
            raise_exception(exp_msg=str(e),
                            exp_num=ncTShareMgntError.NCT_DB_OPERATE_FAILED)
        finally:
            cursor.close()
            conn.close()

        return stg_id

    def edit_perm_share_info(self, strategy_info):
        """
        编辑一条权限共享范围信息
        """
        # 检查策略信息
        self.check_strategy_info(strategy_info, True)

        # 统计编辑的共享者和共享范围的id信息
        shr_user_ids = [user.id for user in strategy_info.sharerUsers]
        shr_dept_ids = [dept.id for dept in strategy_info.sharerDepts]
        scp_user_ids = [(user.id, user.parentId) for user in strategy_info.scopeUsers]
        scp_dept_ids = [dept.id for dept in strategy_info.scopeDepts]

        # 获取当前共享者和共享范围的id
        cur_shr_user_ids = []
        cur_shr_dept_ids = []
        cur_scp_user_ids = []
        cur_scp_dept_ids = []
        query_sql = """
        SELECT * FROM `t_perm_share_strategy`
        WHERE `f_strategy_id` = %s
        """
        results = self.r_db.all(query_sql, strategy_info.strategyId)
        for res in results:
            obj_id = res['f_obj_id']
            parent_id = res['f_parent_id']
            obj_type = res['f_obj_type']
            sharer_or_scope = res['f_sharer_or_scope']

            # 1:共享者，2：共享范围
            if sharer_or_scope == SHARER_TYPE:
                if obj_type == USER_OBJ:
                    cur_shr_user_ids.append(obj_id)
                else:
                    cur_shr_dept_ids.append(obj_id)
            else:
                if obj_type == USER_OBJ:
                    cur_scp_user_ids.append((obj_id, parent_id))
                else:
                    cur_scp_dept_ids.append(obj_id)

        # 统计增加、删除信息
        shr_add_user_ids = list(set(shr_user_ids) - set(cur_shr_user_ids))
        shr_del_user_ids = list(set(cur_shr_user_ids) - set(shr_user_ids))
        shr_add_dept_ids = list(set(shr_dept_ids) - set(cur_shr_dept_ids))
        shr_del_dept_ids = list(set(cur_shr_dept_ids) - set(shr_dept_ids))
        scp_add_user_ids = list(set(scp_user_ids) - set(cur_scp_user_ids))
        scp_del_user_ids = list(set(cur_scp_user_ids) - set(scp_user_ids))
        scp_add_dept_ids = list(set(scp_dept_ids) - set(cur_scp_dept_ids))
        scp_del_dept_ids = list(set(cur_scp_dept_ids) - set(scp_dept_ids))

        # 增加用户共享者
        for user_id in shr_add_user_ids:
            self.w_db.insert(
                't_perm_share_strategy',
                {
                    'f_strategy_id': strategy_info.strategyId,
                    'f_obj_id': user_id,
                    'f_obj_type': USER_OBJ,
                    'f_sharer_or_scope': SHARER_TYPE,
                    'f_status': 1
                }
            )

        # 增加部门共享者
        for dept_id in shr_add_dept_ids:
            self.w_db.insert(
                't_perm_share_strategy',
                {
                    'f_strategy_id': strategy_info.strategyId,
                    'f_obj_id': dept_id,
                    'f_obj_type': DEPART_OBJ,
                    'f_sharer_or_scope': SHARER_TYPE,
                    'f_status': 1
                }
            )

        # 增加范围用户
        for user_id, parent_id in scp_add_user_ids:
            self.w_db.insert(
                't_perm_share_strategy',
                {
                    'f_strategy_id': strategy_info.strategyId,
                    'f_obj_id': user_id,
                    'f_obj_type': USER_OBJ,
                    'f_parent_id': parent_id,
                    'f_sharer_or_scope': SCOPE_TYPE,
                    'f_status': 1
                }
            )

        # 增加范围部门
        for dept_id in scp_add_dept_ids:
            self.w_db.insert(
                't_perm_share_strategy',
                {
                    'f_strategy_id': strategy_info.strategyId,
                    'f_obj_id': dept_id,
                    'f_obj_type': DEPART_OBJ,
                    'f_sharer_or_scope': SCOPE_TYPE,
                    'f_status': 1
                }
            )

        del_sql = """
        DELETE FROM `t_perm_share_strategy`
        WHERE `f_strategy_id` = %s AND `f_obj_id` = %s
        AND `f_obj_type` = %s AND `f_sharer_or_scope` = %s
        """
        del_scope_user_sql = """
        DELETE FROM `t_perm_share_strategy`
        WHERE `f_strategy_id` = %s AND `f_obj_id` = %s AND `f_parent_id`= %s
        AND `f_obj_type` = %s AND `f_sharer_or_scope` = %s
        """

        # 删除用户共享者
        for user_id in shr_del_user_ids:
            self.w_db.query(del_sql, strategy_info.strategyId, user_id, USER_OBJ, SHARER_TYPE)

        # 删除部门共享者
        for dept_id in shr_del_dept_ids:
            self.w_db.query(del_sql, strategy_info.strategyId, dept_id, DEPART_OBJ, SHARER_TYPE)

        # 删除范围用户
        for user_id, parent_id in scp_del_user_ids:
            self.w_db.query(del_scope_user_sql, strategy_info.strategyId,
                            user_id, parent_id, USER_OBJ, SCOPE_TYPE)

        # 删除范围部门
        for dept_id in scp_del_dept_ids:
            self.w_db.query(del_sql, strategy_info.strategyId, dept_id, DEPART_OBJ, SCOPE_TYPE)

    def delete_perm_share_info(self, strategy_id):
        """
        删除一条权限共享范围信息
        """
        if str(strategy_id) == '-1':
            return

        del_sql = """
        DELETE FROM `t_perm_share_strategy`
        WHERE `f_strategy_id` = %s
        """
        self.r_db.query(del_sql, strategy_id)

    def set_perm_share_status(self, strategy_id, status):
        """
        设置权限共享范围状态：
        参数：
           status：
                True： 开启
                Flase：关闭
        """
        status = 1 if status else 0
        update_sql = """
        UPDATE `t_perm_share_strategy`
        SET `f_status` = %s
        WHERE `f_strategy_id` = %s
        """
        self.w_db.query(update_sql, status, strategy_id)

    def get_perm_share_info_cnt(self):
        """
        获取权限共享范围信息总条数
        """
        query_sql = """
        SELECT COUNT(DISTINCT f_strategy_id) AS cnt
        FROM t_perm_share_strategy
        """
        result = self.r_db.one(query_sql)
        return result['cnt']

    def get_perm_share_info_by_id(self, strategy_id):
        """
        根据权限共享id获取共享范围限制信息
        """
        query_sql = """
        SELECT * FROM `t_perm_share_strategy`
        WHERE `f_strategy_id` = %s
        """
        results = self.r_db.all(query_sql, strategy_id)

        if results:
            share_info = ncTPermShareInfo(strategy_id, [], [], [], [])
            share_info.strategyId = strategy_id
            for res in results:
                share_info.status = True if res['f_status'] == 1 else False
                obj_id = res['f_obj_id']
                obj_type = res['f_obj_type']
                parent_id = res['f_parent_id']
                sharer_or_scope = res['f_sharer_or_scope']

                if sharer_or_scope == SHARER_TYPE:
                    if obj_type == USER_OBJ:
                        obj_info = self.get_user_obj_info(obj_id)
                        if obj_info:
                            share_info.sharerUsers.append(obj_info)
                    else:
                        obj_info = self.__get_dept_obj_info(obj_id)
                        if obj_info:
                            share_info.sharerDepts.append(obj_info)
                else:
                    if obj_type == USER_OBJ:
                        obj_info = self.get_user_obj_info(obj_id, parent_id)
                        if obj_info:
                            share_info.scopeUsers.append(obj_info)
                    else:
                        obj_info = self.__get_dept_obj_info(obj_id)
                        if obj_info:
                            share_info.scopeDepts.append(obj_info)
            return share_info
        else:
            return None

    def get_perm_share_info_by_page(self, start, limit):
        """
        分页获取权限共享策略信息
        """
        # 检查参数
        limit_statement = check_start_limit(start, limit)

        query_sql = """
        SELECT DISTINCT f_strategy_id FROM t_perm_share_strategy
        ORDER BY f_strategy_id
        {0}
        """.format(limit_statement)
        results = self.r_db.all(query_sql)

        strategy_infos = []
        for res in results:
            share_info = self.get_perm_share_info_by_id(res['f_strategy_id'])
            strategy_infos.append(share_info)

        return strategy_infos

    def search_perm_share_info(self, start, limit, search_key):
        """
        搜索权限共享策略信息
        """
        # 检查参数
        limit_statement = check_start_limit(start, limit)

        # 增加搜索排序子句
        order_by_str = generate_search_order_sql(['MIN(`t_user`.`f_display_name`)',
                                                  'MIN(`t_department`.`f_name`)'])

        search_sql = """
        SELECT `t_perm_share_strategy`.`f_strategy_id`,
        MIN(`t_user`.`f_display_name`),
        MIN(`t_department`.`f_name`)
        FROM `t_perm_share_strategy`
        LEFT JOIN `t_user`
        ON `t_perm_share_strategy`.`f_obj_id` = `t_user`.`f_user_id`
        LEFT JOIN `t_department`
        ON (`t_perm_share_strategy`.`f_obj_id` = `t_department`.`f_department_id`
        OR `t_perm_share_strategy`.`f_parent_id` = `t_department`.`f_department_id`)
        WHERE (`t_user`.`f_display_name` LIKE %s OR `t_department`.`f_name` LIKE %s)
        GROUP BY `t_perm_share_strategy`.`f_strategy_id`
        ORDER BY {0}
        {1}
        """.format(order_by_str, limit_statement)
        esckey = "%%%s%%" % escape_key(search_key)
        results = self.r_db.all(search_sql, esckey, esckey,
                                escape_key(search_key), escape_key(search_key), esckey, esckey)

        strategy_infos = []
        for res in results:
            share_info = self.get_perm_share_info_by_id(res['f_strategy_id'])
            strategy_infos.append(share_info)

        return strategy_infos

    def get_defaul_strategy_superim_status(self):
        """
        是否叠加默认策略
        """
        return True if int(self.config_manage.get_config("default_strategy_superim_status")) else False

    def set_defaul_strategy_superim_status(self, status):
        """
        设置是否叠加默认策略
        """
        status = 1 if status else 0
        self.config_manage.set_config("default_strategy_superim_status", status)
