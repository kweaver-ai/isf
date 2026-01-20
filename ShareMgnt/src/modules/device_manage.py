#!/usr/bin/python3
# -*- coding:utf-8 -*-
"""
设备绑定管理类
"""
from src.common.db.connector import DBConnector
from src.common.db.db_manager import get_db_name
from src.common.lib import (raise_exception, escape_key,
                            check_start_limit)
from ShareMgnt.ttypes import (ncTShareMgntError, ncTDeviceBindUserInfo, ncTBindStatusSearchScope)
from ShareMgnt.constants import (NCT_USER_ADMIN,
                                 NCT_USER_AUDIT,
                                 NCT_USER_SYSTEM,
                                 NCT_USER_SECURIT,
                                 NCT_ALL_USER_GROUP)


class DeviceManage(DBConnector):
    def __init__(self):
        """
        """

    def search_users_bind_status(self, scope, key, start, limit, cnt_only=False):
        """
        搜索用户设备绑定状态信息
        """
        if scope not in list(ncTBindStatusSearchScope._NAMES_TO_VALUES.values()):
            raise_exception(exp_msg=_("search scope illegal"),
                            exp_num=ncTShareMgntError.NCT_INVALID_SEARCH_SCOPE)

        limit_statement = check_start_limit(start, limit)
        esckey = "%%%s%%" % escape_key(key)

        # 根据scope增加过滤子句
        filter_str = ""
        needed_bind_flag = -1 # 所有状态
        if scope == ncTBindStatusSearchScope.NCT_BIND:
            filter_str = " having max(d.f_bind_flag) = 1"
            needed_bind_flag = 1
        elif scope == ncTBindStatusSearchScope.NCT_UNBIND:
            filter_str = " having max(d.f_bind_flag) is null or max(d.f_bind_flag) = 0"
            needed_bind_flag = 0

        query_sql = f"""
        SELECT u.f_user_id, u.f_login_name, u.f_display_name, max(d.f_bind_flag) as bind_flag
        FROM t_user as u
        LEFT JOIN {get_db_name('anyshare')}.t_device as d
            ON u.f_user_id = d.f_user_id
        WHERE u.f_display_name like %s
            and u.f_user_id not in ('{NCT_USER_ADMIN}', '{NCT_USER_AUDIT}', '{NCT_USER_SYSTEM}', '{NCT_USER_SECURIT}')
        group by u.f_user_id
            {filter_str}
        order by upper(u.f_display_name)
        {limit_statement}
        """
        results = self.r_db.all(query_sql, esckey)

        # 获取“所有用户”配置
        query_all_user_group_sql = f"""
        SELECT f_user_id, max(f_bind_flag) as bind_flag
        FROM {get_db_name('anyshare')}.t_device
        WHERE f_user_id = '{NCT_ALL_USER_GROUP}'
        """
        all_user_group_results = self.r_db.all(query_all_user_group_sql)

        all_user_group_bind_flag = 0
        if len(all_user_group_results) > 0:
            res = all_user_group_results[0]
            all_user_group_bind_flag = res['bind_flag'] if res['bind_flag'] else 0

        if cnt_only:
            count_num = len(results)
            # 需要统计“所有用户”配置项
            if len(key) == 0 and (needed_bind_flag == -1 or all_user_group_bind_flag == needed_bind_flag):
                count_num += 1
            return count_num

        user_infos = []
        # 需要添加“所有用户”配置项
        if start == 0 and len(key) == 0 and (needed_bind_flag == -1 or all_user_group_bind_flag == needed_bind_flag):
            info = ncTDeviceBindUserInfo()
            info.id = NCT_ALL_USER_GROUP
            info.loginName = _("all user")
            info.displayName = _("all user")
            info.bindStatus = all_user_group_bind_flag
            user_infos.append(info)

        for res in results:
            info = ncTDeviceBindUserInfo()
            info.id = res['f_user_id']
            info.loginName = res['f_login_name']
            info.displayName = res['f_display_name']
            info.bindStatus = res['bind_flag'] if res['bind_flag'] else 0
            user_infos.append(info)

        return user_infos
