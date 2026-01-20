#!/usr/bin/python3
# -*- coding:utf-8 -*-
"""This is sharemgnt config manage class"""

from src.common.db.connector import DBConnector
from src.common.lib import (raise_exception,
                            check_is_valid_password,
                            check_url)
from ShareMgnt.ttypes import (ncTShareMgntError,
                              ncTThirdCSFSysConfig,
                              ncTTemplateType,
                              ncTDocType,
                              ncTUserAutoDisableConfig,
                              ncTSearchRange,
                              ncTSearchResults,
                              ncTSearchUserConfig)
from eisoo.tclients import TClient
import json


MIN_USER_SPACE_BYTES = 0


class ConfigManage(DBConnector):
    """
    config manage
    """
    def __init__(self):
        """
        init
        """
        self.custom_config_list = ['exit_pwd',                      # 是否启用退出口令功能
                                   'enable_exit_pwd',               # 退出口令
                                   'id_card_login_status',          # 是否启用身份证号登录
                                   'enable_recycle_delay_delete',   # 是否开启系统回收站
                                   'recycle_delete_delay_time',     # 系统回收站延迟删除时间
                                   'enable_antivirus',              # 杀毒服务开关
                                   'only_share_to_user',            # 是否只能共享给用户
                                   'antivirus_config',              # 杀毒服务器配置
                                   'enable_pwd_control',            # 屏蔽【管控密码】界面中的 "不允许用户自主修改密码" 和 "用户密码：输入框 随机密码"
                                   'enable_set_delete_perm',        # 内链共享时不允许设置 "删除" 权限
                                   'enable_set_folder_security_level',   # 不允许设置文件夹密级
                                   'vcode_server_status',           # 发送验证码服务器类型（邮件/短信）开启状态
                                   'enable_outlink_watermark',      # 允许用户配置外链水印
                                   'dualfactor_auth_server_status', # 双因子登录方式的开关
                                   'enable_update_virus_db',        # 病毒库自动更新开关
                                   'recycle_delete_delay_time_unit',# 系统回收站延迟删除时间计量单位
                                   'catelogue_template_count',      # 编目模板限制条数
                                   'catelogue_count',               # 编目限制条数
                                   'update_virus_db_method',        # 病毒库更新方法
                                   'enable_get_subobj_csf_level',    # 是否获取当前目录子对象的密级
                                   'tag_max_num'                    # 标签限制条数
                                  ]

    def set_config(self, key, value):
        """
        """
        # t_sharemgnt_config 表中 f_value 的类型是 varchar
        # update t_sharemgnt_config set f_value = ture 时, 人大金仓数据库里f_value存储的是'true',其他数据库里f_value存储的是'1'
        # 为了兼容人大金仓数据库, 将bool转成0/1
        if type(value) is bool:
            if value:
                value = 1 
            else:
                value = 0
        update_config_sql = """
        update t_sharemgnt_config
        set f_value = %s
        where f_key = %s
        """
        self.w_db.query(update_config_sql, value, key)

    def replace_config(self, key, value):
        """
        """
        check_config_sql = """
        select f_value from t_sharemgnt_config where f_key = %s
        """
        result = self.r_db.one(check_config_sql, key)

        update_config_sql = ""
        if result:
            update_config_sql = """
            update t_sharemgnt_config set f_value = %s where f_key = %s
            """
        else:
            update_config_sql = """
            insert into t_sharemgnt_config (f_value, f_key)
            values(%s, %s)
            """
        self.w_db.query(update_config_sql, value, key)

    def get_config(self, key):
        """
        """
        select_config_sql = """
        SELECT `f_value` FROM `t_sharemgnt_config`
        WHERE `f_key` = %s
        """
        result = self.r_db.one(select_config_sql, key)
        return result['f_value']

    def set_default_space_size(self, space_size):
        """
        设置默认的用户配额空间
        """
        space_size = int(space_size)
        if space_size <= MIN_USER_SPACE_BYTES:
            raise_exception(exp_msg=_("IDS_USER_QUOTA_NUM_WRONG"),
                            exp_num=ncTShareMgntError.NCT_INVALID_SAPCE_SIZE)

        self.set_config('default_space_size', space_size)

    def get_default_space_size(self):
        """
        获取默认的用户配额空间
        """
        space_size = self.get_config('default_space_size')
        return int(space_size)

    def set_user_doc_status(self, status):
        """
        设置是否默认开启个人文档
        """
        status = 0 if status else -1
        self.set_config('enable_user_doc', status)

    def get_user_doc_status(self):
        """
        获取个人文档默认开启状态
        """
        result = self.get_config('enable_user_doc')
        status = True if int(result) == 0 else False
        return status

    def init_csf_levels(self, csf_levels):
        """
        初始化密级枚举
        """
        if self.get_csf_levels():
            raise_exception(exp_msg=_("IDS_CSFLEVELS_HAS_BEEN_INITIALIZED"),
                            exp_num=ncTShareMgntError.NCT_CSF_LEVEL_ENUM_HAS_BEEN_INITIALIZED)

        if csf_levels and len(csf_levels) <= 11:
            csf_level_dict = dict()
            for csf_level in csf_levels:
                csf_level_dict[csf_level.name] = csf_level.value
            self.replace_config("csf_level_enum", json.dumps(csf_level_dict, ensure_ascii=False))
            # 密级枚举不能小于5
            if min(csf_level_dict.values()) < 5:
                raise_exception(exp_msg=(_("INVALID_PARAM") % csf_levels),
                                exp_num=ncTShareMgntError.NCT_INVALID_PARAMTER)
        else:
            raise_exception(exp_msg=(_("INVALID_PARAM") % csf_levels),
                            exp_num=ncTShareMgntError.NCT_INVALID_PARAMTER)

    def get_csf_levels(self):
        """
        获取密级枚举
        """
        db_config = self.get_config("csf_level_enum")
        if db_config:
            csf_levels = json.loads(db_config)
            return csf_levels
        else:
            return dict()
        
    def get_csf_levels2(self):
        """
        获取密级2枚举
        """
        db_config = self.get_config("csf_level2_enum")
        if db_config:
            csf_levels = json.loads(db_config)
            return csf_levels
        else:
            return dict()

    def get_max_csf_level(self):
        """
        获取最大密级值
        """
        csf_levels = self.get_csf_levels()
        if not csf_levels:
            raise_exception(exp_msg=_("IDS_CSFLEVELS_HAS_NOT_BEEN_INITIALIZED"),
                            exp_num=ncTShareMgntError.NCT_CSF_LEVEL_ENUM_HAS_NOT_BEEN_INITIALIZED)
        return max(csf_levels.values())

    def get_min_csf_level2(self):
        """
        获取最小密级2值
        """
        csf_levels = self.get_csf_levels2()
        if not csf_levels:
            csf_level_dict = dict()
            min_csf_level = 51
            csf_level_dict["公开"] = min_csf_level
            self.replace_config("csf_level2_enum", json.dumps(csf_level_dict, ensure_ascii=False))
            return min_csf_level
        return min(csf_levels.values())
    
    def get_min_csf_level(self):
        """
        获取最小密级值
        """
        csf_levels = self.get_csf_levels()
        if not csf_levels:
            # 问题原因：
            # 在单独部署认证包时，用户密级枚举尚未做初始化，接着部署其它安装包时(例如：AD的主模块包)
            # 安装包内有服务初始化时会调用ShareMgnt的接口去创建用户（AD gmanager-kg-user-rbac 服务），此时密级枚举未初始化，创建用户的接口报错
            # 导致服务panic 从而安装包安装失败
            # 解决办法：
            # 当密级枚举不存在时，初始化一个密级枚举 {"非密": 5}，保证其它服务在调用接口时不报错
            #
            # 特殊场景：
            # 当同时部署AD后，接着部署AS，在AS的管理控制台对系统配置做初始化时，会由于密级枚举已经初始化，而导致无法再次设置用户密级枚举
            # 初始化系统配置界面无法关闭（刷新界面后可解决），此时系统中用户的密级枚举始终是非密
            # 解决办法：
            # 提供数据库脚本，单独设置用户的密级枚举
            csf_level_dict = dict()
            min_csf_level = 5
            csf_level_dict["非密"] = min_csf_level
            self.replace_config("csf_level_enum", json.dumps(csf_level_dict, ensure_ascii=False))
            return min_csf_level
        return min(csf_levels.values())

    def get_clear_cache_interval(self):
        """
        获取清除缓存的时间间隔
        """
        return int(self.get_config("clear_cache_interval"))

    def set_clear_cache_interval(self, interval):
        """
        设置清除缓存的时间间隔
        """
        self.set_config('clear_cache_interval', interval)

    def get_clear_cache_size(self):
        """
        获取清除缓存的配额空间大小
        """
        return int(self.get_config("clear_cache_size"))

    def set_clear_cache_size(self, size):
        """
        设置清除缓存的配额空间大小
        """
        self.set_config('clear_cache_size', size)

    def set_login_strategy_status(self, status):
        """
        获取清除缓存的配额空间大小
        """
        value = 1 if status else 0
        self.set_config("login_strategy_status", value),

    def get_login_strategy_status(self):
        """
        设置清除缓存的配额空间大小
        """
        return int(self.get_config("login_strategy_status")) == 1

    def set_force_clear_cache_status(self, status):
        """
        设置客户端是否强制清除缓存
        """
        self.set_config("force_clear_client_cache", int(status))

    def get_force_clear_cache_status(self):
        """
        获取客户端强制清除缓存状态
        """
        return bool(int(self.get_config("force_clear_client_cache")))

    def set_hide_cache_setting_status(self, status):
        """
        设置是否隐藏客户端缓存设置
        """
        self.set_config("hide_client_cache_setting", int(status))

    def get_hide_cache_setting_status(self):
        """
        获取是否隐藏客户端缓存设置的状态
        """
        return bool(int(self.get_config("hide_client_cache_setting")))

    def set_multi_tenant_status(self, status):
        """
        设置系统多租户状态
        """
        self.set_config("multi_tenant", int(status))

    def get_multi_tenant_status(self):
        """
        获取系统多租户状态
        """
        return bool(int(self.get_config("multi_tenant")))

    def get_secret_mode_status(self):
        """
        获取涉密模式总开关状态
        """
        try:
            return bool(int(self.get_config("enable_secret_mode")))
        except TypeError:
            return False

    def get_system_init_status(self):
        """
        获取系统初始化状态
        """
        try:
            return bool(int(self.get_config("system_init_status")))
        except TypeError:
            return False

    def init_system(self):
        """
        设置控制台https协议状态
        """
        self.set_config("system_init_status", '1')

    def get_uninstall_pwd_status(self):
        """
        获取pc客户端卸载口令状态
        """
        return bool(int(self.get_config("enable_uninstall_pwd")))

    def _check_uninstall_pwd_enabled(self):
        """
        检查pc端卸载口令是否开启，未开启抛出错误
        """

        if not self.get_uninstall_pwd_status():
            raise_exception(exp_msg=_("IDS_UNINSTALL_PWD_NOT_ENABLED"),
                            exp_num=ncTShareMgntError.NCT_UNINSTALL_PWD_NOT_ENABLED)

    def set_uninstall_pwd(self, pwd):
        """
        设置pc端卸载口令
        """
        # 检查pc端卸载口令是否开启，未开启抛出错误
        self._check_uninstall_pwd_enabled()

        # 检查密码有效性
        if not check_is_valid_password(pwd):
            raise_exception(exp_msg=_("invalid password"),
                            exp_num=ncTShareMgntError.NCT_INVALID_PASSWORD)

        self.set_config("uninstall_pwd", pwd)

    def get_uninstall_pwd(self):
        """
        设置pc端卸载口令
        """
        # 检查pc端卸载口令是否开启，未开启抛出错误
        self._check_uninstall_pwd_enabled()

        return self.get_config('uninstall_pwd')

    def check_uninstall_pwd(self, pwd):
        """
        检查pc端卸载口令是否正确
        """
        if pwd != self.get_uninstall_pwd():
            raise_exception(exp_msg=_("IDS_UNINSTALL_PWD_INCORRECT"),
                            exp_num=ncTShareMgntError.NCT_UNINSTALL_PWD_INCORRECT)

    def get_third_csfsys_config(self):
        """
        获取设置第三方标密系统配置
        """
        db_config = self.get_config("third_csfsys_config")
        config_dict = json.loads(db_config)
        config = ncTThirdCSFSysConfig()
        config.isEnabled = config_dict["isEnabled"]
        config.id = config_dict["id"]
        config.only_upload_classified = config_dict["only_upload_classified"]
        config.only_share_classified = config_dict["only_share_classified"]
        config.auto_match_doc_classfication = config_dict["auto_match_doc_classfication"]
        return config

    def set_net_docs_limit_status(self, status):
        """
        设置网段文档库绑定开关状态
        """

        self.set_config('enable_net_docs_limit', status)

    def get_net_docs_limit_status(self):
        """
        获取网段文档库绑定开关状态
        """
        return bool(int(self.get_config('enable_net_docs_limit')))

    def set_ddl_email_notify_mode_status(self, status):
        """
        设置下载量限制配置的邮件通知状态
        """
        self.set_config('enable_ddl_email_notify', status)

    def get_ddl_email_notify_mode_status(self):
        """
        获取下载量限制配置的邮件通知状态
        """
        return bool(int(self.get_config('enable_ddl_email_notify')))

    def set_third_pwd_lock(self, status):
        """
        设置是否启用域认证或第三方认证密码锁策略
        """
        self.set_config('enable_third_pwd_lock', status)

    def get_third_pwd_lock(self):
        """
        获取域认证或第三方认证是否启用密码锁策略状态
        """
        return bool(int(self.get_config('enable_third_pwd_lock')))

    def get_share_doc_status(self, docType, linkType):
        """
        获取共享文档配置
        """
        self.check_doc_type(docType)
        self.check_link_type(linkType)

        if docType == ncTDocType.NCT_USER_DOC:
            # 个人文档
            if linkType:
                # 外链
                return bool(int(self.get_config('enable_user_doc_out_link')))
            else:
                return bool(int(self.get_config('enable_user_doc_inner_link')))

    def check_doc_type(self, docType):
        """
        """
        if docType < ncTDocType.NCT_USER_DOC or docType > ncTDocType.NCT_ARCHIVE_DOC:
            raise_exception(exp_msg=_("IDS_INVALID_DOC_TYPE"),
                            exp_num=ncTShareMgntError.NCT_INVALID_DOC_TYPE)

    def check_link_type(self, linkType):
        """
        """
        if linkType < ncTTemplateType.INTERNAL_LINK or linkType > ncTTemplateType.EXTERNAL_LINK:
            raise_exception(exp_msg=_("IDS_INVALID_LINK_TEMPLATE_TYPE"),
                            exp_num=ncTShareMgntError.NCT_INVALID_LINK_TEMPLATE_TYPE)

    def set_share_doc_status(self, docType, linkType, status):
        """
        设置共享文档开关状态
        """
        self.check_doc_type(docType)
        self.check_link_type(linkType)

        status = 1 if int(status) else 0

        if docType == ncTDocType.NCT_USER_DOC:
            # 个人文档
            if linkType:
                # 外链
                self.set_config('enable_user_doc_out_link', status)
            else:
                self.set_config('enable_user_doc_inner_link', status)

    def get_retain_file_status(self):
        """
        获取文件留底开关状态
        """
        delay_delete_status = self.get_custom_config_of_bool("enable_recycle_delay_delete")
        delay_delete_time = self.get_custom_config_of_int64("recycle_delete_delay_time")
        return delay_delete_status and delay_delete_time == -1

    def set_auto_disable_config(self, config):
        """
        设置用户自动禁用配置
        """
        param = {}
        param["isEnabled"] = config.isEnabled
        param["days"] = config.days
        self.set_config("auto_disable_config", json.dumps(param, ensure_ascii=False))

    def get_auto_disable_config(self):
        """
        获取用户自动禁用配置
        """
        old_config = json.loads(self.get_config("auto_disable_config"))
        config = ncTUserAutoDisableConfig()
        config.isEnabled = old_config["isEnabled"]
        config.days = old_config["days"]
        return config

    def set_retain_out_link_status(self, status):
        """
        设置文件留底开关状态
        """
        self.set_config("retain_out_link_status", int(status))

    def get_retain_out_link_status(self):
        """
        获取文件留底开关状态
        """
        return bool(int(self.get_config("retain_out_link_status")))

    def set_freeze_status(self, status):
        """
        开启关闭冻结功能，True:开启，False:关闭
        """
        self.set_config("enable_freeze", int(status))

    def get_freeze_status(self):
        """
        开启关闭冻结功能，True:开启，False:关闭
        """
        return bool(int(self.get_config("enable_freeze")))

    def get_real_name_auth_status(self):
        """
        获取实名认证开关状态
        """
        return bool(int(self.get_config("enable_real_name_auth")))

    def set_real_name_auth_status(self, status):
        """
        设置实名认证开关状态
        """
        self.set_config("enable_real_name_auth", int(status))

    def set_search_user_config(self, config):
        """
        设置用户共享时搜索配置
        """
        # 检查参数是否为空
        if config is None or config.exactSearch is None or      \
            config.searchRange is None or config.searchResults is None:
                raise_exception(exp_msg=_("parameter is none"),
                            exp_num=ncTShareMgntError.NCT_PARAMETER_IS_NULL)

        # 检查参数合法性
        if config.searchRange < ncTSearchRange.NCT_LOGIN_NAME or    \
            config.searchRange > ncTSearchRange.NCT_LOGIN_AND_DISPLAY or    \
            config.searchResults < ncTSearchResults.NCT_DISPLAY_NAME or     \
            config.searchResults > ncTSearchResults.NCT_LOGIN_AND_DISPLAY:
                raise_exception(exp_msg=_("IDS_INVALID_SEARCH_CONFIG_PARAM"),
                            exp_num=ncTShareMgntError.NCT_INVALID_SEARCH_CONFIG_PARAM)

        param = {}
        param["exactSearch"] = config.exactSearch
        param["searchRange"] = config.searchRange
        param["searchResults"] = config.searchResults
        self.set_config("search_user_config", json.dumps(param, ensure_ascii=False))

    def get_search_user_config(self):
        """
        获取用户共享时搜索配置
        """
        old_config = json.loads(self.get_config("search_user_config"))
        config = ncTSearchUserConfig()
        config.exactSearch = old_config["exactSearch"]
        config.searchRange = old_config["searchRange"]
        config.searchResults = old_config["searchResults"]
        return config

    def get_sms_activate_status(self):
        """
        短信激活是否开启
        """
        sql = """
        SELECT `f_value` FROM `t_sharemgnt_config`
        WHERE `f_key` = %s
        """
        result = self.r_db.one(sql, 'sms_activate')
        return True if result and int(result["f_value"]) == 1 else False

    def __check_custom_config(self, key, value, checkValue=False):
        """
        检查参数合法
        """
        # 检查 key 合法性
        if key not in self.custom_config_list:
            raise_exception(exp_msg=(_("INVALID_PARAM") % key),
                            exp_num=ncTShareMgntError.NCT_INVALID_PARAMTER)

        # 检查 value 合法性
        if checkValue:
            #退出密码有效性
            if "exit_pwd" == key:
                if not check_is_valid_password(value):
                    raise_exception(exp_msg=_("invalid password"),
                                    exp_num=ncTShareMgntError.NCT_INVALID_PASSWORD)
            #系统回收站延迟删除时间有效性
            if "recycle_delete_delay_time" == key:
                if value < 0 and value != -1:
                    raise_exception(exp_msg=_(_("INVALID_PARAM") % value),
                                    exp_num=ncTShareMgntError.NCT_INVALID_PARAMTER)
            # 开启屏蔽 "删除" 选项时, 去除内链模板中的删除权限
            if "enable_set_delete_perm" == key and False == value:
                from src.modules.link_template_manage import LinkTemplateManage
                LinkTemplateManage().remove_del_perm_from_internal_link_template()

            #双因子认证同时只能开启一种
            if "dualfactor_auth_server_status" == key:
                if value.count("true") > 1:
                    raise_exception(exp_msg=_(_("INVALID_PARAM") % value),
                                    exp_num=ncTShareMgntError.NCT_INVALID_PARAMTER)

            # 系统回收站清理时间单位检测
            if "recycle_delete_delay_time_unit" == key:
                if value not in ["day", "week", "month", "year"]:
                    raise_exception(exp_msg=_(_("INVALID_PARAM") % value),
                                    exp_num=ncTShareMgntError.NCT_INVALID_PARAMTER)

    def set_custom_config_of_string(self, key, value):
        """
        设置自定义配置，String
        """
        self.__check_custom_config(key, value, checkValue=True)
        self.set_config(key, value)

    def set_custom_config_of_int64(self, key, value):
        """
        设置自定义配置，Int64
        """
        self.__check_custom_config(key, value, checkValue=True)
        self.set_config(key, value)

    def set_custom_config_of_bool(self, key, value):
        """
        设置自定义配置，Bool
        """
        self.__check_custom_config(key, value, checkValue=True)
        self.set_config(key, value)

    def get_custom_config_of_string(self, key):
        """
        获取自定义配置，String
        """
        self.__check_custom_config(key, None, checkValue=False)
        return self.get_config(key)

    def get_custom_config_of_int64(self, key):
        """
        获取自定义配置，Int64
        """
        self.__check_custom_config(key, None, checkValue=False)
        return int(self.get_config(key))

    def get_custom_config_of_bool(self, key):
        """
        获取自定义配置，Bool
        """
        self.__check_custom_config(key, None, checkValue=False)
        return bool(int(self.get_config(key)))
