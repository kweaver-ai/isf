#!/usr/bin/python3
# -*- coding:utf-8 -*-
"""This is group manage class"""
import uuid
from eisoo.tclients import TClient
from src.common.db.connector import DBConnector, ConnectorManager
from src.common.http import pub_nsq_msg
from src.common.lib import (raise_exception,
                            check_start_limit,
                            check_name,
                            generate_search_order_sql,
                            generate_group_str,
                            escape_key,
                            check_is_uuid)
from src.modules.user_manage import UserManage
from ShareMgnt.ttypes import (ncTPersonGroup,
                              ncTSearchPersonGroup,
                              ncTShareMgntError)

TOPIC_ORG_NAME_MODIFY = "core.org.name.modify"


class GroupManage(DBConnector):
    """
    group manage
    """

    def __init__(self):
        """
        init
        """
        self.user_manage = UserManage()
        self.tmp_group = _("IDS_TMP_PERSON_GROUP")

    def check_self_group(self, user_id, group_id):
        """
        检测联系人组是否属于当前登录用户
        """
        count = 0
        if check_is_uuid(group_id):
            sql = """
            SELECT COUNT(*) AS cnt
            FROM `t_person_group`
            WHERE `f_user_id` = %s
                  AND `f_group_id` = %s
            """
            count = self.r_db.one(sql, user_id, group_id)['cnt']
        if not count:
            raise_exception(exp_msg=_("group not exists"),
                            exp_num=ncTShareMgntError. NCT_GROUP_NOT_EXIST)

    def check_group_name_available(self, group_name):
        """
        检测联系人组名是否合法
        """
        # 检测合法性
        if not group_name or not check_name(group_name):
            raise_exception(exp_msg=_("group name illegal"),
                            exp_num=ncTShareMgntError.NCT_INVALID_GROUP_NAME)

        if group_name == self.tmp_group:
            raise_exception(exp_msg=_("group is exists"),
                            exp_num=ncTShareMgntError.NCT_GROUP_HAS_EXIST)

    def check_group_exists(self, user_id, group_name, group_id=''):
        """
        检测联系人组是否存在
        """
        where = ''
        if group_id:
            where = "AND `f_group_id` <> '{0}'".format(
                self.w_db.escape(group_id))
        sql = """
        SELECT COUNT(*) AS cnt FROM `t_person_group`
        WHERE `f_user_id` = %s AND `f_group_name` = %s {0}
        """.format(where)
        result = self.r_db.one(sql, user_id, group_name)['cnt']
        if result:
            raise_exception(exp_msg=_("group is exists"),
                            exp_num=ncTShareMgntError.NCT_GROUP_HAS_EXIST)

    def check_tmp_group(self, group_id):
        """
        检测这个组是否为临时组
        Args:
            group_id: string ID
        Raise:
            如果是临时组，丢出异常
        """
        if not check_is_uuid(group_id):
            raise_exception(exp_msg=_("group not exists"),
                            exp_num=ncTShareMgntError.NCT_GROUP_NOT_EXIST)

        sql = """
        SELECT `f_group_name` FROM `t_person_group`
        WHERE `f_group_id` = %s
        """
        result = self.r_db.one(sql, group_id)
        if result and result['f_group_name'] == self.tmp_group:
            raise_exception(exp_msg=_("cannot control tmp group"),
                            exp_num=ncTShareMgntError.NCT_CANNOT_OPERATE_TMP_GROUP)

        name = ""
        if result:
            name = result['f_group_name']
        return name

    def create_tmp_group(self, uid):
        """
        创建临时联系人组
        Args:
            uid: string, 用户ID
        """
        group_id = str(uuid.uuid1())
        self.w_db.insert("t_person_group", [group_id, uid,
                                            self.tmp_group, 0])
        return group_id

    def create_person_group(self, cur_user_id, group_name):
        """
        创建联系人组
        """
        group_name = group_name.strip()

        # 检查联系人组名
        self.check_group_name_available(group_name)

        # 检查用户是否存在
        self.user_manage.check_user_exists(cur_user_id)

        # 检查用户是否存在同名联系人组
        self.check_group_exists(cur_user_id, group_name)

        # 新增组
        group_id = str(uuid.uuid1())
        self.w_db.insert("t_person_group", [group_id, cur_user_id,
                                            group_name, 0])
        return group_id

    def edit_person_group(self, cur_user_id, group_id, new_name):
        """
        编辑联系人组
        """
        # 检查是否是临时联系人组
        old_name = self.check_tmp_group(group_id)

        # 检查联系组名是否正常
        self.check_group_name_available(new_name)

        # 检查用户是否存在
        self.user_manage.check_user_exists(cur_user_id)

        # 联系人组是否是当前用户联系人组
        self.check_self_group(cur_user_id, group_id)

        # 检查新的联系人组名是否冲突
        self.check_group_exists(cur_user_id, new_name, group_id)

        sql = """
        UPDATE `t_person_group`
        SET `f_group_name` = %s
        WHERE `f_group_id` = %s
        """
        self.w_db.query(sql, new_name, group_id)

        if old_name != new_name:
            # 发送部门显示名更新nsq消息
            pub_nsq_msg(TOPIC_ORG_NAME_MODIFY, {
                        "id": group_id, "new_name": new_name, "type": "contactor"})

        return

    def get_person_groups(self, cur_user_id):
        """
        获取联系人组
        """
        # 检查用户是否存在
        self.user_manage.check_user_exists(cur_user_id)

        # 获取当前用户的联系人
        sql = """
        SELECT f_group_id, f_group_name, f_person_count
        FROM t_person_group
        WHERE f_user_id = '{0}'
        ORDER BY FIELD(f_group_name, '{1}') desc,
        upper(f_group_name)
        """.format(cur_user_id, self.tmp_group)
        result = self.r_db.all(sql)
        groups = []
        for row in result:
            group = ncTPersonGroup(row['f_group_id'],
                                   row['f_group_name'],
                                   row['f_person_count'])
            groups.append(group)
        return groups

    def __add_person(self, user_id, group_id):
        """
        添加用户到数据库
        """
        data = {
            "f_group_id": group_id,
            "f_user_id": user_id
        }
        self.w_db.insert("t_contact_person", data)

        # 更新person_count
        sql = """
        UPDATE `t_person_group` SET `f_person_count` = `f_person_count` + 1
        WHERE `f_group_id` = %s
        """
        self.w_db.query(sql, group_id)

    def check_person_in_group(self, user_id, group_id):
        """
        检查联系人是否已添加
        """
        sql = """
        SELECT `f_id` FROM `t_contact_person`
            WHERE `f_user_id` = %s AND `f_group_id` = %s
        """
        result = self.r_db.one(sql, user_id, group_id)
        return True if result else False

    def add_person_by_id(self, cur_user_id, users_id, group_id):
        """
        根据ID添加联系人
        """
        # 检查当前用户是否除存在
        self.user_manage.check_user_exists(cur_user_id)

        # 检查联系人组和用户关系
        self.check_self_group(cur_user_id, group_id)

        # 去掉用户ID两侧空格
        users_id = [uid.strip() for uid in users_id]

        # 去重
        users_id = list(set(users_id))

        for uid in users_id:
            # 自己不能添加自己
            if cur_user_id == uid:
                continue
            # 用户不存在，忽略
            elif not self.user_manage.check_user_exists(uid, False):
                continue
            # 用户已经在组，忽略
            elif self.check_person_in_group(uid, group_id):
                continue
            # 验证通过，进行添加
            else:
                self.__add_person(uid, group_id)
        return

    def del_person(self, cur_user_id, users_id, group_id):
        """
        删除联系人
        """
        # 检查当前用户是否除存在
        self.user_manage.check_user_exists(cur_user_id)

        # 检查联系人组和用户关系
        self.check_self_group(cur_user_id, group_id)

        # 去重
        users_id = list(set(users_id))

        # 验证用户是否在组中
        for uid in users_id:
            self.user_manage.check_user_exists(uid)
            if not self.check_person_in_group(uid, group_id):
                raise_exception(exp_msg=_("person not exists"),
                                exp_num=ncTShareMgntError.NCT_CONTACT_NOT_EXIST)

        sql = """
        DELETE FROM `t_contact_person`
        WHERE `f_group_id` = %s AND `f_user_id` = %s
        """
        for uid in users_id:
            self.w_db.query(sql, group_id, uid)

        # 更新person_count
        sql = """
        UPDATE `t_person_group` AS `group` SET `f_person_count` = (
            SELECT COUNT(*) FROM `t_contact_person`
            WHERE `f_group_id` = `group`.`f_group_id`
        )
        WHERE `f_group_id` = %s
        """
        self.w_db.query(sql, group_id)

    def get_person_from_group(self, cur_user_id, group_id, start, limit):
        """
        从联系人组获取联系人
        """
        self.user_manage.check_user_exists(cur_user_id)
        self.check_self_group(cur_user_id, group_id)

        limit_statement = check_start_limit(start, limit)

        sql = """
        SELECT `t_user`.`f_user_id`, `t_user`.`f_login_name`,
            `t_user`.`f_display_name`, `t_user`.`f_mail_address`, `t_user`.`f_idcard_number`,
            `t_user`.`f_status`, `t_user`.`f_remark`, `t_user`.`f_priority`, `t_user`.`f_csf_level`,
            `t_user`.`f_pwd_control`, `t_user`.`f_oss_id`, `t_user`.`f_auto_disable_status`
        FROM `t_user`
        JOIN `t_contact_person`
            ON `t_contact_person`.`f_user_id` = `t_user`.`f_user_id`
        JOIN `t_person_group`
            ON `t_person_group`.`f_group_id` = `t_contact_person`.`f_group_id`
        WHERE `t_contact_person`.`f_group_id` = %s
        order by t_user.f_priority, upper(`t_user`.`f_display_name`)
        {0}
        """.format(limit_statement)
        result = self.r_db.all(sql, group_id)

        users = []
        for row in result:
            users.append(self.user_manage.fetch_user(row, get_quota=False))

        return users

    def search_person_from_group(self, userId, search_key):
        """
        从联系人组信息中根据用户显示名搜索用户
        """
        self.user_manage.check_user_exists(userId)
        # 搜索用户所创建的所有联系人组
        sql = """
        select f_group_id, f_person_count
        from t_person_group
        where f_user_id = %s;
        """
        results = self.r_db.all(sql, userId)
        group_ids = []
        if results:
            group_ids = [result["f_group_id"]
                         for result in results if result["f_person_count"] > 0]
        if len(group_ids) == 0:
            return []

        groupStr = generate_group_str(group_ids)
        order_by_str = generate_search_order_sql(['f_display_name'])
        esckey = "%%%s%%" % escape_key(search_key)
        sql = """
        select u.f_user_id,u.f_login_name,u.f_display_name, r.f_group_id, g.f_group_name
        from t_user as u
        inner join t_contact_person as r
        on u.f_user_id = r.f_user_id
        inner join t_person_group as g
        on g.f_group_id = r.f_group_id
            and r.f_group_id in ({0})
        WHERE u.f_display_name LIKE %s
        order by {1}, upper(f_display_name)
        """.format(groupStr, order_by_str)
        results = self.r_db.all(sql, esckey, escape_key(
            search_key), escape_key(search_key))
        search_person_infos = []
        for result in results:
            search_person_info = ncTSearchPersonGroup()
            search_person_info.userId = result["f_user_id"]
            search_person_info.loginName = result["f_login_name"]
            search_person_info.displayName = result["f_display_name"]
            search_person_info.groupId = result["f_group_id"]
            search_person_info.groupName = result["f_group_name"]
            search_person_infos.append(search_person_info)
        return search_person_infos

    def set_person_group(self, user_id, group_name, new_contact_ids):
        """
        user_id中，如果不存在group_name的联系人组，则创建，并将contact_ids添加到组中
        如果存在，则将联系人组更新为contact_ids（删除掉不在contact_ids中的，并创建新的）
        """

        # 执行情况
        ret_code = 0
        ret_group_id = ""

        # 使用事务保证原子性
        conn = ConnectorManager.get_db_conn()
        cursor = conn.cursor()
        cursor.execute("begin")

        try:
            # 检查user_id对应的group_name是否存在
            strsql = """
            select f_group_id from t_person_group
            where f_user_id = %s and f_group_name = %s
            """
            cursor.execute(strsql, (user_id, group_name))
            results = cursor.fetchone()
            if results is None:
                ret_group_id = str(uuid.uuid1())
                strsql = """
                insert into t_person_group(f_group_id,f_user_id,f_group_name,f_person_count)
                values(%s,%s,%s,%s)
                """
                cursor.execute(strsql, (ret_group_id, user_id, group_name, 0))

                for contact_id in new_contact_ids:
                    strsql = """
                    insert into t_contact_person(f_group_id,f_user_id)
                    values(%s, %s)
                    """
                    cursor.execute(strsql, (ret_group_id, contact_id))

                count = len(new_contact_ids)
                strsql = """
                update t_person_group set f_person_count = %s
                where f_group_id = %s
                """
                cursor.execute(strsql, (count, ret_group_id))

                # 返回1表示新建成功
                ret_code = 1
            else:
                ret_group_id = results["f_group_id"]
                old_contact_ids = []
                strsql = """
                select f_user_id from t_contact_person where f_group_id = %s
                """
                cursor.execute(strsql, (ret_group_id,))
                contact_results = cursor.fetchall()
                for tmp_result in contact_results:
                    old_contact_ids.append(tmp_result["f_user_id"])

                delete_ids = []
                add_ids = []
                old_contact_ids.sort()
                new_contact_ids.sort()
                if old_contact_ids != new_contact_ids:
                    # 需要删除掉的
                    for contact_id in old_contact_ids:
                        if contact_id not in new_contact_ids:
                            delete_ids.append(contact_id)
                    for contact_id in new_contact_ids:
                        if contact_id not in old_contact_ids:
                            add_ids.append(contact_id)

                    for contact_id in add_ids:
                        strsql = """
                        insert into t_contact_person(f_group_id,f_user_id)
                        values(%s,%s)
                        """
                        cursor.execute(strsql, (ret_group_id, contact_id))

                    for contact_id in delete_ids:
                        strsql = """
                        delete from t_contact_person
                        where f_group_id = %s and f_user_id = %s
                        """
                        cursor.execute(strsql, (ret_group_id, contact_id))

                    count = len(new_contact_ids)

                    strsql = """
                    update t_person_group set f_person_count = %s
                    where f_group_id = %s
                    """
                    cursor.execute(strsql, (count, ret_group_id))

                    # 返回2表示修改
                    ret_code = 2
                else:
                    # 返回0表示不需要修改
                    ret_code = 0

            cursor.execute("commit")
        except Exception as e:
            cursor.execute("rollback")
            raise_exception("set_person_group error: %s", str(e))

        return ret_code, ret_group_id
