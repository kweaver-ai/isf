#!/usr/bin/python3
# -*- coding:utf-8 -*-

import time
import json
import copy
import traceback
from collections import deque
from eisoo.tclients import TClient
from src.common import global_info
from src.common.sharemgnt_logger import ShareMgnt_Log
from src.common.business_date import BusinessDate
from src.common.redis_connector import OPRedis
from src.common.eacp_log import eacp_log
from src.third_party_auth.third_config_manage import ThirdConfigManage
from src.third_party_auth.ou_manage.as_ou_manage import ASOuManage
from ShareMgnt.constants import NCT_USER_ADMIN, NCT_UNDISTRIBUTE_USER_GROUP, ncTUsrmUserStatus
from ShareMgnt.ttypes import ncTPluginType
import uuid


class OuProgressInfo(object):
    """
    组织同步进度信息统计类
    """
    def __init__(self):
        """
        初始化函数
        """
        self.reset()

    def reset(self):
        """
        """
        self.total_num = 0    # 要同步的部门总数
        self.synced_num = 0   # 已同步的部门数
        self.added_num = 0    # 已增加的部门数
        self.deleted_num = 0  # 已删除的部门数
        self.updated_num = 0  # 已更新的部门数
        self.moved_num = 0    # 已移动的部门数
        self.failed_num = 0   # 同步失败的部门数


class UserProgressInfo(object):
    """
    用户同步进度信息统计类
    """
    def __init__(self):
        """
        """
        self.reset()

    def reset(self):
        """
        """
        self.total_num = 0    # 要同步的用户总数
        self.synced_num = 0   # 已同步的用户数
        self.added_num = 0    # 已增加的用户数
        self.moved_num = 0    # 移动的用户数
        self.deleted_num = 0  # 已删除的用户数
        self.updated_num = 0  # 已更新的用户数
        self.failed_num = 0   # 已失败的用户数


class BaseSyncer(object):
    """
    第三方同步基本件
    """
    def __init__(self, app_id=None, src_ou_manage=None, b_eacplog=True):
        """
        初始化函数
        """
        self.app_id = app_id
        self.app_name = app_id
        self.ou_pregress_info = OuProgressInfo()
        self.user_progress_info = UserProgressInfo()
        self.src_ou_manage = src_ou_manage(b_eacplog)
        self.dest_ou_manage = ASOuManage(b_eacplog)
        self.b_eacplog = b_eacplog
        self.new_user_ids = set()
        self.sync_disable_user = True
        self.forced_sync = True  # 是否强制同步用户状态
        self.opredis = OPRedis()
        self.pending_deleted_ou_ids = []
        self.pending_removed_user_ids = []

    def init_progress_info(self):
        """
        初始换同步进度信息
        """
        self.ou_pregress_info.reset()
        self.user_progress_info.reset()
        ou_total, user_total = self.src_ou_manage.get_ous_num_info()
        self.ou_pregress_info.total_num = ou_total
        self.user_progress_info.total_num = user_total
        self.pending_deleted_ou_ids = []
        self.pending_removed_user_ids = []
        self.moved_ou_ids = []

    def sync_eacp_log(self, op_type, msg, ex_msg=None):
        """
        记录组织操作日志
        参数：
            loglevel - {str} 日志级别
            NCT_LL_ALL = 0,    // 所有日志级别
            NCT_LL_INFO = 1,   // 信息
            NCT_LL_WARN = 2,   // 警告

            optype - {str} 操作类型
            NCT_SOT_ALL = 0,        // 所有操作
            NCT_SOT_CREATE = 1,     // 新建操作
            NCT_SOT_MODIFY = 2,     // 修改操作
            NCT_SOT_DELETE = 3,     // 删除操作
            NCT_OOT_ADD_USER_TO_DEP = 4,    // 添加用户到部门
            NCT_OOT_MOVE_USER_AND_DEP = 5,  // 迁移用户和部门
            NCT_OOT_DISABLE = 6,            // 禁用
            NCT_OOT_ENABLE = 7,             // 启用
        """
        eacp_log(_("IDS_SYNCER"),
                        global_info.LOG_TYPE_MANAGE,
                        global_info.USER_TYPE_INTER,
                        global_info.LOG_LEVEL_INFO,
                        op_type,
                         msg,
                         ex_msg,
                         raise_ex=True)

    @classmethod
    def get_register_app_id(self):
        """
        获取需要注册的app id(子类必须实现)
        """
        raise

    def get_server_info(self):
        """
        获取第三方服务器配置信息
        """
        third_info_manage = ThirdConfigManage()
        third_infos = third_info_manage.get_third_party_config(ncTPluginType.AUTHENTICATION)
        if third_infos and third_infos[0].enabled and third_infos[0].config:
            try:
                # 同时使用config和internalConfig
                config = json.loads(third_infos[0].config)
                config.update(json.loads(third_infos[0].internalConfig))
                return config
            except Exception as ex:
                ShareMgnt_Log("获取服务器配置信息异常: ex=%s", str(ex))

    @classmethod
    def check_server_status(self, app_id):
        """
        检查服务器状态
        """
        return True

    def get_sync_interval(self):
        """
        获取同步间隔
        """
        config = self.get_server_info()
        if config and "syncInterval" in config:
            return int(config["syncInterval"])
        else:
            return 1800

    def get_sync_status(self):
        """
        获取同步状态:
            每次同步开始，线程检测同步状态，如果同步状态为False，则结束同步.
        Args:
            返回值：
                True： 同步
                Flase：不同步
        """
        third_info_manage = ThirdConfigManage()
        third_info = third_info_manage.get_third_party_auth_by_appid(self.app_id)
        return third_info.enabled

    def get_dest_depart(self):
        """
        获取同步到的目的组织或部门
        """
        dest_depart_id = None
        server_info = self.get_server_info()

        # 这里因为第三方组织结构同步的配置是字典,而域控同步的配置是ncTUsrmDomainInfo, 故需要区分处理
        if isinstance(server_info, dict) and "destDeptId" in server_info:
            dept_id = server_info["destDeptId"]
            dept_info = self.dest_ou_manage.get_depart_info_by_id(dept_id)
            if dept_info:
                dest_depart_id = dept_info.depart_id

        # 检查根组织名
        root_third_id = self.src_ou_manage.get_root_id()
        root_info = self.src_ou_manage.get_ou(root_third_id)

        # 这里因为第三方组织结构同步的配置是字典,而域控同步的配置是ncTUsrmDomainInfo, 故需要区分处理
        if isinstance(server_info, dict) and "rootOrgName" in server_info and server_info["rootOrgName"]:
            root_info.ou_name = server_info["rootOrgName"].encode("utf-8")
        if not isinstance(root_info.ou_name, str):
            root_info.ou_name = root_info.ou_name.decode('utf8')

        # 不存在，则使用根组织来创建目的部门
        if not dest_depart_id:
            dept_info = self.dest_ou_manage.get_ou(str(root_third_id))

            # 默认根部门已存在，判断是否需要更新部门名
            if dept_info:
                if dept_info.ou_name != root_info.ou_name:
                    self.dest_ou_manage.update_ou(dept_info.depart_id,
                                                  root_info,
                                                  self.ou_pregress_info)
                dest_depart_id = dept_info.depart_id
            else:
                # 先设置站点信息
                dest_depart_id = self.dest_ou_manage.add_ou(-1, root_info, self.ou_pregress_info)

        return dest_depart_id

    def get_src_ous(self):
        """
        获取同步的域组织:
        """
        return [self.src_ou_manage.get_root_id()]

    def syn_sub_ous(self, depart_id, ou_infos):
        """
        同步子部门
        """
        if not depart_id:
            return

        dest_sub_ous = self.dest_ou_manage.get_sub_ous_by_depart_id(depart_id)

        dest_sub_ous_dict = {}
        for dest_ou in dest_sub_ous:
            dest_sub_ous_dict[dest_ou.third_id] = dest_ou

        src_sub_ids = [ou.third_id for ou in ou_infos]
        dest_sub_ids = [ou.third_id for ou in dest_sub_ous]

        # 删除组织部门
        for ou_id in dest_sub_ids:
            if ou_id not in src_sub_ids:
                # 将部门标记为删除，等部门移动完毕后，再进行删除
                self.pending_deleted_ou_ids.append(ou_id)

        for ou in ou_infos:
            if ou.third_id in dest_sub_ids:
                # 更新组织部门
                try:
                    dest_ou = dest_sub_ous_dict[ou.third_id]
                    if not self.dest_ou_manage.compare_ou_info(dest_ou, ou):
                        self.dest_ou_manage.update_ou(depart_id,
                                                      ou,
                                                      self.ou_pregress_info)
                except Exception as ex:
                    ShareMgnt_Log("修改部门异常: third_ou_id=%s, ex=%s", ou.third_id, str(ex))
            else:
                # 增加组织部门
                try:
                    b_exist = self.dest_ou_manage.check_depart_exists(ou.third_id)
                    if not b_exist:
                        self.dest_ou_manage.add_ou(depart_id,
                                                   ou,
                                                   self.ou_pregress_info)
                    else:
                        # 记录移动的部门，避免移动失败时，将父部门删除掉时，删除了自身
                        self.moved_ou_ids.append(ou.third_id)

                        # 移动部门，move_ou内部会处理重名问题
                        self.dest_ou_manage.move_ou(depart_id,
                                                    ou,
                                                    self.ou_pregress_info)

                        # 如果部门成功被移动，从待删除部门中删除
                        if ou.third_id in self.pending_deleted_ou_ids:
                            self.pending_deleted_ou_ids.remove(ou.third_id)

                except Exception as ex:
                    ShareMgnt_Log("增加部门异常: third_ou_id=%s, ex=%s", ou.third_id, str(ex))

    def sync_sub_users(self, depart_id, user_infos):
        """
        同步子用户
        """
        if not depart_id:
            return

        dest_sub_users = self.dest_ou_manage.get_sub_users_by_depart_id(depart_id)

        dest_sub_users_dict = {}
        for dest_user in dest_sub_users:
            dest_sub_users_dict[dest_user.third_id] = dest_user

        src_user_ids = [user.third_id for user in user_infos]
        dest_user_ids = [user.third_id for user in dest_sub_users]

        for user_id in dest_user_ids:
            if user_id not in src_user_ids:
                # 移除用户
                # 记录需移除的用户，最后处理
                try:
                    disable_flag = self.src_ou_manage.get_user_disable_status(user_id)
                    self.pending_removed_user_ids.append({
                        "depart_id": depart_id,
                        "user_id": user_id,
                        "disable_flag": disable_flag
                    })

                except Exception as ex:
                    ShareMgnt_Log("记录需从部门移除用户异常: user_id=%s, ex=%s", user_id, str(ex))

        for user in user_infos:
            old_status = user.status
            if user.third_id in dest_user_ids:
                # 更新用户
                try:
                    dest_user = dest_sub_users_dict[user.third_id]
                    if not self.dest_ou_manage.compare_user_info(dest_user, user):
                        self.dest_ou_manage.update_user(dest_user, user, self.user_progress_info)
                    else:
                        self.user_progress_info.synced_num += 1
                except Exception as ex:
                    ShareMgnt_Log("更新用户异常: name=%s, ex=%s", user.login_name, str(ex))
            else:
                b_exist = self.dest_ou_manage.chec_user_exists(user.third_id)
                # 增加用户
                if not b_exist:
                    try:
                        # 是否同步禁用用户
                        if self.sync_disable_user is False and old_status is False:
                            self.user_progress_info.synced_num += 1
                            continue

                        self.dest_ou_manage.add_user_to_ou(depart_id,
                                                           copy.deepcopy(user),
                                                           self.user_progress_info)

                        # 记录新增用户的third_id
                        self.new_user_ids.add(user.third_id)
                    except Exception as ex:
                        ShareMgnt_Log("添加用户异常: name=%s, ex=%s", user.login_name, str(ex))

                # 添加已存在的用户到该部门
                else:
                    try:
                        # 这里要先判断下用户是否需要更新
                        dest_user = self.dest_ou_manage.get_user(user.third_id)
                        if dest_user and (not self.dest_ou_manage.compare_user_info(dest_user, user)):
                            self.dest_ou_manage.update_user(dest_user, user, self.user_progress_info)

                        self.dest_ou_manage.move_user_to_ou(depart_id,
                                                            user,
                                                            self.user_progress_info)
                    except Exception as ex:
                        ShareMgnt_Log("移动用户异常: name=%s, ex=%s", user.login_name, str(ex))

            # 设置了禁用字段时会根据字段启用/禁用用户，否则不做处理
            if self.forced_sync and old_status is not None:

                # 如果用户已在目的部门,则不需要查询数据库
                user_info = None
                if user.third_id in dest_sub_users_dict:
                    user_info = dest_sub_users_dict[user.third_id]

                # 如果用户不在目的部门,则需要查询数据库
                if not user_info:
                    user_info = self.dest_ou_manage.get_user(user.third_id)

                # 不处理不同步的用户
                if user_info:
                    if old_status:
                        try:
                            self.dest_ou_manage.enable_user(user_info)
                        except Exception as ex:
                            ShareMgnt_Log("启用用户失败: name=%s, ex=%s", user.login_name, str(ex))
                    else:
                        self.dest_ou_manage.disable_user(user_info)

    def sync_ou(self, parent_depart_id, src_ou_ids):
        """
        同步组织
        """
        try:
            src_ou_infos = []
            src_user_infos = []

            # 处理域根组织和目的部门相同的情况
            if len(src_ou_ids) == 1:
                root_id = src_ou_ids[0]
                ou_info = self.src_ou_manage.get_ou(root_id)
                dept_info = self.dest_ou_manage.get_depart_info_by_id(parent_depart_id)
                if dept_info and ou_info and str(dept_info.third_id) == str(ou_info.third_id):
                    # 如果是域根组织，则获取根组织下的所有子组织
                    src_user_infos = self.src_ou_manage.get_sub_users(root_id)
                    src_ou_ids = self.src_ou_manage.get_sub_ou_ids(root_id)

            for ou_id in src_ou_ids:
                ou_info = self.src_ou_manage.get_ou(ou_id)
                ou_info.ou_name = ou_info.ou_name
                ou_info.third_id = ou_info.third_id
                src_ou_infos.append(ou_info)

            # 同步根组织下的部门
            self.syn_sub_ous(parent_depart_id, src_ou_infos)

            # 同步根组织下的用户
            self.sync_sub_users(parent_depart_id, src_user_infos)

            # 按广度优先获取所有需要同步的第三方ouid
            ou_queue = deque(src_ou_ids)
            try:
                while True:
                    ou_id = ou_queue.popleft()
                    sub_ou_ids = self.src_ou_manage.get_sub_ou_ids(ou_id)
                    src_ou_ids.extend(sub_ou_ids)
                    ou_queue.extend(sub_ou_ids)
            # 所有子部门获取完毕，则会丢异常
            except IndexError:
                pass

            for ou_id in src_ou_ids:
                dest_depart_info = self.dest_ou_manage.get_ou(ou_id)
                sub_users = self.src_ou_manage.get_sub_users(ou_id)
                sub_ous = self.src_ou_manage.get_sub_ous(ou_id)

                if not dest_depart_info:
                    continue

                # 同步一级子用户
                self.sync_sub_users(dest_depart_info.depart_id, sub_users)

                # 同步一级子部门
                self.syn_sub_ous(dest_depart_info.depart_id, sub_ous)

            # 处理需移除的用户
            ShareMgnt_Log("处理需移除的用户开始")
            self.process_removed_user()
            ShareMgnt_Log("处理需移除的用户结束")

            # 删除部门
            ShareMgnt_Log("删除部门开始")
            self.process_deleted_ou()
            ShareMgnt_Log("删除部门结束")

            # 同步未分配用户组
            ShareMgnt_Log("同步未分配用户组开始")
            self.sync_undistributed_user()
            ShareMgnt_Log("同步未分配用户组结束")
        except Exception:
            ShareMgnt_Log(traceback.format_exc())

    def sync_undistributed_user(self):
        """
        同步第三方未分配用户组
        """
        src_undist_users = self.src_ou_manage.get_undistributed_users()
        src_undist_users_dict = {}
        for src_user in src_undist_users:
            src_undist_users_dict[src_user.third_id] = src_user

        dst_undist_users = self.dest_ou_manage.get_undistributed_users()
        dst_undist_users_dict = {}
        for dst_user in dst_undist_users:
            dst_undist_users_dict[dst_user.third_id] = dst_user

        # 在源端，不在目的端，需要未分配用户组中新建用户
        for src_user in src_undist_users:
            if src_user.third_id not in dst_undist_users_dict:
                self.dest_ou_manage.add_user_to_ou(NCT_UNDISTRIBUTE_USER_GROUP, src_user, self.user_progress_info)

        # 在目的端，不在源端，需要在未分配用户中禁用该用户
        for dst_user in dst_undist_users:
            if dst_user.third_id not in src_undist_users_dict:
                self.dest_ou_manage.disable_user(dst_user)

        # 在源端，也在目的端，需要在未分配中启用该用户
        for dst_user in dst_undist_users:
            if dst_user.third_id in src_undist_users_dict:
                # 因为用户名称可能变化，所以要先更新用户
                self.dest_ou_manage.update_user(dst_user, src_undist_users_dict[dst_user.third_id], self.user_progress_info)
                self.dest_ou_manage.enable_user(dst_user)

    def delete_disabled_undist_users(self):
        """
        删除掉禁用的未分配用户的个人文档和用户
        """
        if self.src_ou_manage.get_delete_disable_undist_users_flag() is False:
            ShareMgnt_Log("delete_disable_undist_users_flag is False.")
        else:
            ShareMgnt_Log("delete_disable_undist_users_flag is True.")

            # 获取所有未分配用户
            dst_undist_users = self.dest_ou_manage.get_undistributed_users()
            for user in dst_undist_users:
                if user.status == ncTUsrmUserStatus.NCT_STATUS_DISABLE:
                    # 删除用户
                    self.dest_ou_manage.delete_user(user.third_id, self.user_progress_info)
            return

    def process_removed_user(self):
        """
        处理需移除的用户
        """
        for user in self.pending_removed_user_ids:
            try:
                self.dest_ou_manage.remove_user_from_ou(user["depart_id"], user["user_id"], user["disable_flag"], self.user_progress_info)
            except Exception as ex:
                ShareMgnt_Log("处理需移除的用户异常: ex=%s, user_id=%s, depart_id=%s", str(ex), user["user_id"], user["depart_id"])

    def process_deleted_ou(self):
        """
        sync_ou会标记要删除的部门，这里执行真正的删除
        如果要删除的部门包含移动的部门，则不能被删除
        """
        # 先将moved_ou_ids（第三方id）转成 anyshare的 departid（uuid格式）
        all_moved_depart_ids = []
        for moved_id in self.moved_ou_ids:
            depart_id = self.dest_ou_manage.get_depart_id(moved_id)
            if depart_id:
                all_moved_depart_ids.append(depart_id)

        try:
            for third_ou_id in self.pending_deleted_ou_ids:

                # 如果删除的部门包含已经移动过的部门，则不进行处理
                if self.dest_ou_manage.is_contain_moved_ou(third_ou_id, all_moved_depart_ids):
                    ShareMgnt_Log("删除部门异常: third_ou_id=%s, contained moved ou", third_ou_id)
                    continue

                self.dest_ou_manage.delete_ou(third_ou_id,
                                              self.src_ou_manage.get_user_disable_status,
                                              self.ou_pregress_info,
                                              self.user_progress_info)
        except Exception as ex:
            ShareMgnt_Log("删除部门异常: third_ou_id=%s, ex=%s", third_ou_id, str(ex))

    def sync_update_ou_manager(self):
        """
        更新部门的负责人
        """
        # 获取客户所有的部门负责人信息（depart_third_id + user_manager_third_id)
        src_ou_manager_tree = self.src_ou_manage.get_ou_manager()

        # 获取客户所有的用户负责人信息（user_third_id + user_manager_third_id)
        src_user_manager_tree = self.src_ou_manage.get_user_manager()

        self.dest_ou_manage.update_manager(src_user_manager_tree, src_ou_manager_tree)

    def sync(self):
        """
        同步开始函数
        """
        try:
            # 初始化服务器信息
            server_info = self.get_server_info()
            self.src_ou_manage.init_server_info(server_info)

            # 初始化导入的一些全局配置
            self.dest_ou_manage.init_sync_config(server_info)

            # 初始化进度信息
            self.init_progress_info()

            # 获取目的部门和源组织
            dest_depart_id = self.get_dest_depart()
            src_ou_ids = self.get_src_ous()

            # 同步组织
            self.sync_ou(dest_depart_id, src_ou_ids)

            # 删除掉禁用的未分配用户的个人文档和用户
            self.delete_disabled_undist_users()

            ShareMgnt_Log("组织管理员配额更新开始")
            # 更新组织管理员配额空间
            self.dest_ou_manage.sync_update_responsible_person_space()
            ShareMgnt_Log("组织管理员配额更新结束")

            # 更新部门的负责人以及用户上级
            ShareMgnt_Log("更新部门的负责人以及用户上级开始")
            self.sync_update_ou_manager()
            ShareMgnt_Log("更新部门的负责人以及用户上级结束")

            # 加载各同步器自己定制功能
            self.src_ou_manage.call_after_sync()

            # 记录审计日志
            dest_depart_name = dest_depart_id
            try:
                dept_info = self.dest_ou_manage.get_depart_info_by_id(dest_depart_id)
                dest_depart_name = dept_info.ou_name
            except:
                pass

            self.sync_eacp_log(global_info.LOG_OP_TYPE_IMPORT,
                          _('IDS_SYNC_SUCCESS') % (self.app_name, dest_depart_name))

            ShareMgnt_Log("部门统计信息----总数：%d, 增加: %d, 删除: %d, 更新：%d, 失败: %d",
                          self.ou_pregress_info.total_num,
                          self.ou_pregress_info.added_num,
                          self.ou_pregress_info.deleted_num,
                          self.ou_pregress_info.updated_num,
                          self.ou_pregress_info.failed_num)

            ShareMgnt_Log("用户统计信息----总数：%d, 增加: %d, 移动：%d, 删除: %d, 更新：%d, 失败: %d",
                          self.user_progress_info.total_num,
                          self.user_progress_info.added_num,
                          self.user_progress_info.moved_num,
                          self.user_progress_info.deleted_num,
                          self.user_progress_info.updated_num,
                          self.user_progress_info.failed_num)

        except Exception:
            ShareMgnt_Log(traceback.format_exc())
