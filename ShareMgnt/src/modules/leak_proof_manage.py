#!/usr/bin/python3
# -*- coding:utf-8 -*-
"""
防泄密管理实现
"""
from src.common.db.connector import DBConnector
from src.common.lib import (raise_exception, escape_key,
                            generate_search_order_sql,
                            check_start_limit)
from src.modules.config_manage import ConfigManage
from src.modules.user_manage import UserManage
from src.modules.department_manage import DepartmentManage
from ShareMgnt.ttypes import ncTShareMgntError, ncTLeakProofStrategyInfo


class LeakProofManage(DBConnector):
    """
    LeakProofManage
    """
    def __init__(self):
        """
        init
        """
        self.config_manage = ConfigManage()
        self.user_manage = UserManage()
        self.depart_manage = DepartmentManage()

    def set_leak_proof_status(self, status):
        """
        设置系统防泄密状态
        @param status: bool, true为开启, false为关闭
        """
        if status:
            self.config_manage.set_config("leak_proof_status", "1")
        else:
            self.config_manage.set_config("leak_proof_status", "0")

    def get_leak_proof_status(self):
        """
        获取系统防泄密状态
        @ret bool: true为开启, false为关闭
        """
        if self.config_manage.get_config("leak_proof_status") == "1":
            return True
        else:
            return False

    def add_strategy(self, param):
        """
        添加防泄密策略
        @param param: ncTAddLeakProofStrategyParam
        @ret int: 返回成功后的策略id
        """
        # 检查param.permValue是否合法
        if param.permValue < 1 or param.permValue > 3:
            raise_exception(exp_msg=_("IDS_INVALID_PERM_VALUE"),
                            exp_num=ncTShareMgntError.NCT_INVALID_PERM_VALUE)

        # 检查param.accessorType是否合法
        if param.accessorType != 1 and param.accessorType != 2:
            raise_exception(exp_msg=_("IDS_INVALID_ACCESSOR_TYPE"),
                            exp_num=ncTShareMgntError.NCT_INVALID_ACCESSOR_TYPE)

        # 检查accessorId是否存在
        if param.accessorType == 1:
            if not self.user_manage.is_user_id_exists(param.accessorId):
                raise_exception(exp_msg=_("IDS_ACCESSOR_ID_NOT_EXISTS"),
                                exp_num=ncTShareMgntError.NCT_ACCESSOR_ID_NOT_EXISTS)
        else:
            if not self.depart_manage.is_depart_id_exists(param.accessorId):
                raise_exception(exp_msg=_("IDS_ACCESSOR_ID_NOT_EXISTS"),
                                exp_num=ncTShareMgntError.NCT_ACCESSOR_ID_NOT_EXISTS)

        sql = """
        select f_strategy_id,f_perm_value from t_leak_proof_strategy where f_accessor_id = %s
        """
        result = self.w_db.one(sql, param.accessorId)

        if result:
            # 如果accessorId已经存在，抛出异常
            raise_exception(exp_msg=_("IDS_ACCESSOR_ID_ALREADY_EXISTS"),
                            exp_num=ncTShareMgntError.NCT_ACCESSOR_ID_ALREADY_EXISTS)
        else:
            # 如果不存在，则进行添加
            sql = """
            insert into t_leak_proof_strategy(f_accessor_id, f_accessor_type, f_perm_value)
            values(%s, %s, %s)
            """
            self.w_db.query(sql, param.accessorId, param.accessorType, param.permValue)

            sql = """
            select f_strategy_id from t_leak_proof_strategy where f_accessor_id = %s
            """

            last_id = self.w_db.one(sql, param.accessorId)['f_strategy_id']

            return last_id

    def edit_strategy(self, param):
        """
        编辑防泄密策略
        @param param: ncTEditLeakProofStrategyParam
        @ret: 无返回值
        """
        # 检查startegyid是否存在
        sql = """
        select f_strategy_id,f_perm_value from t_leak_proof_strategy where f_strategy_id = %s
        """
        result = self.w_db.one(sql, param.strategyId)

        # 策略id不存在，抛出异常
        if not result:
            raise_exception(exp_msg=_("IDS_STRATEGY_ID_NOT_EXISTS"),
                            exp_num=ncTShareMgntError.NCT_STRATEGY_ID_NOT_EXISTS)

        # 如果未设置permValue，返回
        if not param.permValue:
            return

        # 检查param.permValue是否合法
        if param.permValue < 1 or param.permValue > 3:
            raise_exception(exp_msg=_("IDS_INVALID_PERM_VALUE"),
                            exp_num=ncTShareMgntError.NCT_INVALID_PERM_VALUE)

        # 权限值不同，编辑权限
        if result["f_perm_value"] != param.permValue:
            sql = """
            update t_leak_proof_strategy set f_perm_value = %s
            where f_strategy_id = %s
            """
            self.w_db.query(sql, param.permValue, param.strategyId)

    def delete_strategy(self, strategyId):
        """
        删除防泄密策略
        @param param: strategyId, 策略id
        @ret: 无返回值
        """
        sql = """
        delete from t_leak_proof_strategy
        where f_strategy_id = %s
        """
        self.w_db.query(sql, strategyId)

    def get_strategy_count(self):
        """
        获取防泄密记录总数，用来分页
        @ret int: 总记录条数
        """
        sql = """
        select count(f_strategy_id) as cnt from t_leak_proof_strategy
        """
        result = self.w_db.one(sql)
        return result["cnt"]

    def get_strategy_infos_by_page(self, start, limit):
        """
        分页获取防泄密策略信息
        @param start: int，开始位置 >= 0
        @param limit: -1, 取的条数 >= -1
        @ret int: 总记录条数
        """
        limit_statement = check_start_limit(start, limit)

        # 分页获取策略信息
        sql = """
        SELECT `t_leak_proof_strategy`.`f_strategy_id` AS `strategy_id`,
        `t_leak_proof_strategy`.`f_accessor_id` AS `accessor_id`,
        `t_leak_proof_strategy`.`f_accessor_type` AS `accessor_type`,
        `t_leak_proof_strategy`.`f_perm_value` AS `perm_value`,
        `t_user`.`f_display_name` AS `user_name`,
        `t_department`.`f_name` AS `dept_name`
        FROM  `t_leak_proof_strategy`
        LEFT JOIN `t_user`
        ON `t_leak_proof_strategy`.`f_accessor_id` = `t_user`.`f_user_id`
        LEFT JOIN `t_department`
        ON `t_leak_proof_strategy`.`f_accessor_id` = `t_department`.`f_department_id`
        WHERE 1=%s
        ORDER BY accessor_type,
        upper(`dept_name`),
        upper(`user_name`)
        {0}
        """.format(limit_statement)
        results = self.w_db.all(sql, 1)
        return self._convert_strategy_info(results)

    def search_strategy_count(self, key):
        """
        根据key搜索防泄密策略信息，返回搜索到的条数
        """
        sql = """
        SELECT count(`t_leak_proof_strategy`.`f_strategy_id`) as cnt
        FROM  `t_leak_proof_strategy`
        LEFT JOIN `t_user`
        ON `t_leak_proof_strategy`.`f_accessor_id` = `t_user`.`f_user_id`
        LEFT JOIN `t_department`
        ON `t_leak_proof_strategy`.`f_accessor_id` = `t_department`.`f_department_id`
        WHERE (`t_user`.`f_display_name` LIKE %s OR  `t_department`.`f_name` LIKE %s)
        """
        esckey = "%%%s%%" % escape_key(key)
        result = self.r_db.one(sql, esckey, esckey)
        return result["cnt"]

    def search_strategy_infos_by_page(self, key, start, limit):
        """
        根据key搜索防泄密策略信息，返回搜索到的条数
        """
        limit_statement = check_start_limit(start, limit)

        # 增加搜索排序子句
        order_by_str = generate_search_order_sql(['t_user.f_display_name', 't_department.f_name'])

        sql = """
        SELECT `t_leak_proof_strategy`.`f_strategy_id` AS `strategy_id`,
        `t_leak_proof_strategy`.`f_accessor_id` AS `accessor_id`,
        `t_leak_proof_strategy`.`f_accessor_type` AS `accessor_type`,
        `t_leak_proof_strategy`.`f_perm_value` AS `perm_value`,
        `t_user`.`f_display_name` AS `user_name`,
        `t_department`.`f_name` AS `dept_name`
        FROM  `t_leak_proof_strategy`
        LEFT JOIN `t_user`
        ON `t_leak_proof_strategy`.`f_accessor_id` = `t_user`.`f_user_id`
        LEFT JOIN `t_department`
        ON `t_leak_proof_strategy`.`f_accessor_id` = `t_department`.`f_department_id`
        WHERE (`t_user`.`f_display_name` LIKE %s OR  `t_department`.`f_name` LIKE %s)
        ORDER BY accessor_type,
        {0},
        upper(`dept_name`),
        upper(`user_name`)
        {1}
        """.format(order_by_str, limit_statement)

        esckey = "%%%s%%" % escape_key(key)
        results = self.r_db.all(sql, esckey, esckey,
                                escape_key(key), escape_key(key), esckey, esckey)
        return self._convert_strategy_info(results)

    def _convert_strategy_info(self, db_infos):
        """
        转化db results to ncTLeakProofStrategyInfo
        """
        strategy_infos = []
        if db_infos:
            for res in db_infos:
                info = ncTLeakProofStrategyInfo()
                info.strategyId = res['strategy_id']
                info.accessorId = res['accessor_id']
                info.accessorType = int(res['accessor_type'])
                info.permValue = int(res['perm_value'])
                info.accessorName = res['user_name'] if info.accessorType == 1 else res['dept_name']
                strategy_infos.append(info)
        return strategy_infos
