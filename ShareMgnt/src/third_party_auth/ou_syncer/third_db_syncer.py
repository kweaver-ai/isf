#!/usr/bin/python3
# -*- coding:utf-8 -*-

from src.modules.department_manage import DepartmentManage
from src.modules.third_db_manage import ThirdDBManage
from src.third_party_auth.ou_syncer.base_syncer import BaseSyncer


class ThirdDbSyncer(BaseSyncer):
    """
    """
    def __init__(self, app_id, src_ou_manage):
        """
        """
        super(ThirdDbSyncer, self).__init__(app_id, src_ou_manage)
        self.app_id = app_id
        self.third_db_manage = ThirdDBManage()
        self.depart_manage = DepartmentManage()
        self.server_info = self.get_server_info()
        self.app_name = self.server_info['dbInfo'].name

    def get_server_info(self):
        """
        获取服务器信息
        """
        third_db_info = self.third_db_manage.get_third_db_info(self.app_id)
        third_table_info = self.third_db_manage.get_third_db_table_infos(self.app_id)
        third_sync_info = self.third_db_manage.get_third_db_sync_config(self.app_id)

        server_info = {}
        server_info["destDeptId"] = third_sync_info.parentDepartId
        server_info["syncInterval"] = third_sync_info.syncInterval
        server_info["spaceSize"] = third_sync_info.spaceSize
        server_info["rootName"] = third_sync_info.thirdRootName
        server_info["tableInfo"] = third_table_info
        server_info["dbInfo"] = third_db_info
        server_info['syncInfo'] = third_sync_info
        return server_info

    def get_sync_status(self):
        """
        获取同步状态
        """
        status = self.third_db_manage.get_status(self.app_id)
        return status
