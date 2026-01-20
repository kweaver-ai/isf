#!/usr/bin/python3
# -*- coding:utf-8 -*-
"""
用户组织管理模块，包括用户和组织的数据结构已经用户组织管理基类
"""
from src.common.db.connector import DBConnector
from ShareMgnt.ttypes import *
from src.common.global_info import DEFAULT_DEPART_PRIORITY


class OuInfo(object):
    """
    组织部门信息
    ou_name: 组织部门名
    depart_id: 部门id
    sub_third_ou_ids: 子部门的第三方id
    sub_third_user_ids: 子用户的第三方id
    """

    def __init__(self, depart_id=None, ou_name=None, third_id=None, parent_id=None, depart_path=None,
                 remark=None, status=None, code=None, position=None, manager_id=None):
        self.ou_name = ou_name
        self.depart_id = depart_id
        self.third_id = third_id
        self.sub_third_ou_ids = []
        self.sub_third_user_ids = []
        self.dn = ''
        self.oss_id = ''
        self.priority = DEFAULT_DEPART_PRIORITY
        self.is_enterprise = False
        self.type = ncTUsrmDepartType.NCT_DEPART_TYPE_THIRD
        self.parent_id = parent_id
        self.depart_path = depart_path
        self.remark = remark
        self.status = status
        self.code = code
        self.position = position
        self.manager_id = manager_id

    def __repr__(self):
        L = ['%s=%r' % (key, value)
             for key, value in self.__dict__.items()]
        return '%s(%s)' % (self.__class__.__name__, ', '.join(L))


class UserInfo(object):
    """
    用户信息
    user_id: 用户id
    login_name: 登录名
    display_name: 显示名
    idcard_number: 身份证号
    tel_number: 电话号
    email: 邮箱
    password: 密码
    status: 状态
    third_id: 第三方id
    third_ou_id: 所属的第三方组织id
    type: 用户类型，默认为第三方用户
    doc_status: 是否创建个人文档
    """
    def __init__(self, user_id=None, login_name=None, display_name=None, idcard_number=None,
                 tel_number=None,email=None, password='', status=None, third_id=None,
                 third_ou_id=None, space_size=None, user_type=None,
                 doc_status=None, third_ou_ids=[], third_attr=None, csf_level=None,
                 position=None, code=None, manager_id=None, csf_level2=None):
        self.user_id = user_id
        self.login_name = login_name
        self.display_name = display_name
        self.email = email
        self.idcard_number = idcard_number
        self.tel_number = tel_number
        self.password = password
        self.status = status
        self.third_id = third_id
        self.third_ou_id = third_ou_id
        self.third_ou_ids = third_ou_ids if third_ou_ids else []
        self.dn = ''
        self.server_type = 0
        self.space_size = space_size
        self.oss_id = ''
        self.priority = 999
        if user_type is None or not isinstance(user_type, int):
            self.type = ncTUsrmUserType.NCT_USER_TYPE_THIRD
        elif user_type not in ncTUsrmUserType._VALUES_TO_NAMES:
            self.type = ncTUsrmUserType.NCT_USER_TYPE_THIRD
        else:
            self.type = user_type
        self.doc_status = doc_status
        self.third_attr = ''
        self.csf_level = csf_level
        self.position = position
        self.csf_level2 = csf_level2
        self.code = code
        self.manager_id = manager_id

    def __repr__(self):
        L = ['%s=%r' % (key, value)
             for key, value in self.__dict__.items()]
        return '%s(%s)' % (self.__class__.__name__, ', '.join(L))


def raise_exception(func_name):
        ex_str = "function 'BaseOuManage.%s' is not implemented" % func_name
        raise AttributeError(ex_str)


class BaseOuManage(DBConnector):
    """
    组织用户管理基类
    继承此类的继承类如果是作为源组织管理类，则需要实现下面接口：
        get_root_id
        get_ou
        get_sub_ous
        get_sub_users

    如果继承类是作为目的组织管理类，则需要实现下面接口：
        add_ou
        update_ou
        delete_ou
        add_user
        update_user
        delete_user
        update_manager
    """
    def __init__(self, b_eacplog=False):
        """
        初始化函数
        """
        self.root_id = "-1"
        self.ou_user_tree = {}
        self.ou_depart_tree = {}
        self.ou_user_manager_tree = {}
        self.ou_depart_manager_tree = {}

    def init_server_info(self, server_info):
        """
        初始化服务器信息
        """
        raise_exception("init_server_info")

    def add_ou(self, parent_id, ou_info):
        """
        增加组织部门
        """
        raise_exception("add_ou")

    def update_ou(self, parent_id, ou_info):
        """
        更新组织部门
        """
        raise_exception("update_ou")

    def delete_ou(self, user_info):
        """
        删除组织部门
        """
        raise_exception("delete_ou")

    def add_user(self, parent_id, user_info):
        """
        添加用户
        """
        raise_exception("add_user")

    def update_user(self, user_info):
        """
        更新用户
        """
        raise_exception("update_user")

    def disable_user(self):
        """
        禁用用户
        """
        raise_exception("disable_user")

    def enable_user(self):
        """
        启用用户
        """
        raise_exception("enable_user")

    def delete_user(self, third_user_id):
        """
        删除用户
        """
        raise_exception("delete_user")

    def update_manager(self, user_manager_infos, depart_manager_infos):
        """
        更新用户上级和部门负责人
        """
        raise_exception("update_manager")

    def get_ous_num_info(self):
        """
        获取部门和用户总数信息
        """
        ou_total = len(self.ou_depart_tree)
        user_total = 0
        for _, ou_info in list(self.ou_depart_tree.items()):
            user_total += len(ou_info.sub_third_user_ids)
        return ou_total, user_total

    def get_root_id(self):
        """
        获取根组织的第三方id
        """
        return self.root_id

    def get_ou(self, third_ou_id):
        """
        获取组织部门信息
        """
        return self.ou_depart_tree[third_ou_id]
    
    def get_ou_manager(self):
        """
        获取部门负责人
        """
        return self.ou_depart_manager_tree
    
    def get_user_manager(self):
        """
        获取用户上级
        """
        return self.ou_user_manager_tree

    def get_sub_ous(self, third_ou_id):
        """
        获取子组织或部门
        """
        org_info = self.ou_depart_tree[third_ou_id]
        sub_ou_ids = org_info.sub_third_ou_ids

        sub_ous = []
        for third_id in sub_ou_ids:
            sub_ou = self.ou_depart_tree[third_id]
            sub_ou.ou_name = sub_ou.ou_name
            sub_ou.third_id = sub_ou.third_id
            sub_ous.append(sub_ou)

        return sub_ous

    def get_sub_users(self, third_ou_id):
        """
        获取子用户
        """
        org_info = self.ou_depart_tree[third_ou_id]
        sub_user_ids = org_info.sub_third_user_ids

        sub_users = []
        for third_id in sub_user_ids:
            sub_user = self.ou_user_tree[third_id]
            sub_user.login_name = sub_user.login_name
            sub_user.third_id = sub_user.third_id
            sub_user.display_name = \
                sub_user.display_name if sub_user.display_name else None
            sub_users.append(sub_user)

        return sub_users

    def get_sub_ou_ids(self, third_ou_id):
        """
        获取子组织的id
        """
        return self.ou_depart_tree[third_ou_id].sub_third_ou_ids

    def get_sync_doc_user_ids(self):
        """
        """
        return set(self.ou_user_tree)

    def get_user_disable_status(self, third_ou_id):
        # 将用户从部门中移到未分配时，获取用户的禁用状态
        # 如果返回True，表示禁用
        # 返回False，表示不禁用
        return True

    def get_undistributed_users(self):
        # 获取未分配的用户
        return []

    def get_delete_disable_undist_users_flag(self):
        # 默认不删除禁用的未分配用户
        return False

    def get_root_node(self, *args):
        # 获取第三方根组织节点
        raise NotImplementedError("Method not implement")

    def expand_node(self, *args):
        # 展开第三方节点
        raise NotImplementedError("Method not implement")

    def get_ous(self, *args):
        # 获取选中的组织信息
        raise NotImplementedError("Method not implement")

    def call_after_sync(self):
        """
        组织结构同步完成后，调用
        """
        pass
