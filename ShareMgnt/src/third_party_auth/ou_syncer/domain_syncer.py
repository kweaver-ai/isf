#!/usr/bin/python3
# -*- coding:utf-8 -*-

import traceback
from collections import deque
from src.common.sharemgnt_logger import ShareMgnt_Log
from src.modules.department_manage import DepartmentManage
from src.modules.domain_manage import DomainManage
from src.modules.ldap_manage import (LdapManage, name2dn)
from src.third_party_auth.ou_syncer.base_syncer import BaseSyncer


class DomainSyncer(BaseSyncer):
    """
    """
    def __init__(self, app_id, src_ou_manage):
        """
        """
        super(DomainSyncer, self).__init__(app_id, src_ou_manage)
        self.app_id = app_id
        self.domain_manage = DomainManage()
        self.depart_manage = DepartmentManage()
        self.server_info = self.get_server_info()
        self.app_name = self.server_info.name
        self.sync_disable_user = True
        self.forced_sync = self.server_info.config.forcedSync

    @classmethod
    def get_register_app_id(self):
        """
        返回域控id作为第三方appid
        """
        app_ids = []
        domains = DomainManage().get_all_domains()
        for domain in domains:
            if domain.status and domain.syncStatus == 0:
                app_ids.append(str(int(domain.id)))
        return app_ids

    def get_server_info(self):
        """
        获取服务器信息
        """
        return self.domain_manage.get_available_domain_by_id(self.app_id)

    @classmethod
    def check_server_status(self, app_id):
        """
        检查服务器状态
        """
        server_info = DomainManage().get_domain_by_id(app_id)
        self.ldap_manage = LdapManage(server_info.ipAddress,
                                      server_info.adminName,
                                      server_info.password,
                                      server_info.port,
                                      use_ssl=server_info.useSSL)

    def get_sync_interval(self):
        """
        获取同步间隔, 默认返回300秒
        """
        server_config = self.domain_manage.get_domain_sync_config(self.app_id)
        if server_config:
            if server_config.syncInterval:
                return server_config.syncInterval * 60
        return 300

    def get_sync_status(self):
        """
        获取同步状态
        """
        sync_status = self.domain_manage.get_domain_sync_status(self.server_info.id)
        if sync_status == 0:
            return True
        else:
            return False

    def get_dest_depart(self):
        """
        获取同步到目的部门
        """
        dest_depart_id = None
        base_dn = name2dn(self.server_info.name)

        server_config = self.domain_manage.get_domain_sync_config(self.app_id)
        if server_config:
            dest_depart_id = server_config.destDepartId

        b_exist = True
        if not dest_depart_id:
            b_exist = False
        else:
            b_exist = self.depart_manage.check_depart_exists(dest_depart_id, True, False)

        if not b_exist:
            dest_depart_id = self.dest_ou_manage.get_depart_id(base_dn)
            if not dest_depart_id:
                dest_depart_id = self.depart_manage.add_third_depart_to_db(base_dn,
                                                                           self.server_info.name)
        return dest_depart_id

    def get_src_ous(self):
        """
        获取同步的域组织
        """
        src_ou_ids = []
        base_dn = name2dn(self.server_info.name)

        # 没有指定同步的域组织,则使用域根组织
        server_config = self.domain_manage.get_domain_sync_config(self.app_id)
        if server_config:
            src_ou_dns = server_config.ouPath
            if src_ou_dns:
                for ou_dn in src_ou_dns:
                    third_id = self.src_ou_manage.get_third_id_by_dn(ou_dn)
                    if third_id:
                        src_ou_ids.append(third_id)

        if not src_ou_ids:
            src_ou_ids = [base_dn]

        return src_ou_ids

    def set_third_root_ou_name(self):
        """
        """
        pass

    def sync_ou(self, parent_depart_id, src_ou_ids):
        """
        parent_depart_id: anyshare中的父部门id
        src_ou_ids: ad中的id
        同步所选的组织结构（包括上层组织）
        """
        try:
            # 生成组织用户关系
            ou_dict, user_dict = \
                self.src_ou_manage.generate_ou_user_relation(parent_depart_id,
                                                             src_ou_ids,
                                                             self.app_id)

            src_user_infos = user_dict["-1"] if "-1" in user_dict else []
            src_ou_ids = ou_dict["-1"] if "-1" in ou_dict else []

            # 目的部门不能在第三方同步对象内
            dept_info = self.dest_ou_manage.get_depart_info_by_id(parent_depart_id)
            if dept_info and dept_info.third_id in ou_dict:
                raise Exception("dept third_id can't been in src_ou_ids")

            # 同步根组织的子部门及用户
            src_ou_infos = self.src_ou_manage.get_ou_infos_by_ou_id(src_ou_ids)
            self.syn_sub_ous(parent_depart_id, src_ou_infos)
            self.sync_sub_users(parent_depart_id, src_user_infos)

            ou_queue = deque(src_ou_ids)
            while len(ou_queue):
                ou_id = ou_queue.popleft()
                if ou_id in ou_dict:
                    sub_ou_ids = ou_dict[ou_id]
                    src_user_infos = user_dict[ou_id] if ou_id in user_dict else []
                    src_ou_infos = self.src_ou_manage.get_ou_infos_by_ou_id(sub_ou_ids)
                    # 目的部门不存在，跳过
                    dest_depart_info = self.dest_ou_manage.get_ou(ou_id)
                    if not dest_depart_info:
                        continue
                    self.syn_sub_ous(dest_depart_info.depart_id, src_ou_infos)
                    self.sync_sub_users(dest_depart_info.depart_id, src_user_infos)
                    ou_queue.extend(sub_ou_ids)

            # 处理需移除的用户
            ShareMgnt_Log("处理需移除的用户开始")
            self.process_removed_user()
            ShareMgnt_Log("处理需移除的用户结束")

            # 删除部门
            ShareMgnt_Log("删除部门开始")
            self.process_deleted_ou()
            ShareMgnt_Log("删除部门结束")
        except Exception:
            ShareMgnt_Log(traceback.format_exc())
