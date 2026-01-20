# coding: utf-8
"""
第三方认证配置管理
"""
import json
import uuid
import os
import sys
import requests
import os
from datetime import datetime
from eisoo.tclients import TClient
from src.common.sharemgnt_logger import ShareMgnt_Log
from src.common.business_date import BusinessDate
from src.common.db.connector import DBConnector
from src.common.lib import raise_exception, strip_whitespace
from src.driven.service_access.ossgateway_config import OssgatewayDriven
from ShareMgnt.ttypes import (ncTThirdPartyAuthConf,
                              ncTShareMgntError,
                              ncTThirdPartyPluginInfo,
                              ncTPluginType,
                              ncTThirdPartyConfig)
from EThriftException.ttypes import ncTException


class ThirdConfigManage(DBConnector):
    """
    第三方配置管理模块
    """

    def __init__(self):
        """
        初始化函数
        """
        self.pluginCid = "a1c47132fbb04d40911ebe2eda1a624f"

        # 添加插件环境变量
        self.add_third_party_plugin_path()
        self.ossgateway_driven = OssgatewayDriven()

    def add_third_party_plugin_path(self):
        """
        添加 auth_manage 和 ou_manage 路径给第三方插件调用
        """
        paths = []
        paths.append(os.path.realpath(os.path.join(
            os.path.dirname(sys.argv[0]), "third_party_auth/auth_manage")))
        paths.append(os.path.realpath(os.path.join(
            os.path.dirname(sys.argv[0]), "third_party_auth/ou_manage")))

        # 添加第三方插件路径
        paths.append("/sysvol/plugin")

        for path in paths:
            if path not in sys.path:
                sys.path.append(path)

    def get_third_party_info_auth(self):
        """
        获取第三方认证配置信息，暂时保留以兼容，后续需要删除
        """
        sql = """
        SELECT `f_id`,`f_app_id`,`f_app_name`,`f_config`, `f_enable`
        FROM `t_third_party_auth`
        WHERE `f_plugin_type` = %s
        """

        result = self.r_db.one(sql, ncTPluginType.AUTHENTICATION)

        third_info = ncTThirdPartyAuthConf()
        if result:
            third_info.indexId = result['f_id']
            third_info.thirdPartyId = result['f_app_id']
            third_info.thirdPartyName = result['f_app_name']
            third_info.config = result['f_config']
            third_info.enabled = True if result['f_enable'] else False

            try:
                json.loads(third_info.config)
            except ValueError:
                raise_exception(exp_msg=_("IDS_THIRD_CONFIG_FAILED"),
                                exp_num=ncTShareMgntError.NCT_FAILED_THIRD_CONFIG)
        else:
            third_info.enabled = False

        return third_info

    def get_third_party_auth_by_appid(self, appid):
        """
        根据appid获取第三方认证信息，暂时保留以兼容，后续需要删除
        """
        sql = """
        SELECT `f_app_id`,`f_enable`,`f_config`
        FROM `t_third_party_auth`
        WHERE `f_app_id` = %s AND `f_plugin_type` = %s
        """
        result = self.r_db.one(sql, appid, ncTPluginType.AUTHENTICATION)
        third_info = ncTThirdPartyAuthConf()
        if result:
            third_info.thirdPartyId = result["f_app_id"]
            third_info.enabled = result["f_enable"]

            try:
                json_config = json.loads(result['f_config'])
            except ValueError:
                raise_exception(exp_msg=_("IDS_THIRD_CONFIG_FAILED"),
                                exp_num=ncTShareMgntError.NCT_FAILED_THIRD_CONFIG)

            third_info.config = json.dumps(json_config)

        return third_info

    def get_third_party_plugin_status(self, appId):
        """
        根据appid获取第三方插件状态
        """
        sql = """
        SELECT `f_object_id` FROM `t_third_party_auth`
        WHERE `f_app_id` = %s
        """
        result = self.r_db.one(sql, appId)
        if result:
            return True
        else:
            return False

    def get_third_party_config_by_index_id(self, indexId):
        """
        根据indexId获取第三方配置
        """
        sql = """
        SELECT `f_app_id`,`f_app_name`,`f_config`, `f_internal_config`, `f_enable`,
               `f_plugin_name`, `f_plugin_type`, `f_object_id`, `f_id`
        FROM `t_third_party_auth`
        WHERE `f_id` = %s
        """
        result = self.r_db.one(sql, indexId)
        if not result:
            raise_exception(exp_msg=_("IDS_THIRD_APPID_SECRET_FAILED"),
                            exp_num=ncTShareMgntError.NCT_INVALID_APPID_OR_APPKEY)

        third_info = ncTThirdPartyConfig()
        third_info.thirdPartyId = result['f_app_id']
        third_info.thirdPartyName = result['f_app_name']
        third_info.config = result['f_config']
        third_info.internalConfig = result['f_internal_config']
        third_info.enabled = True if result['f_enable'] else False
        third_info.indexId = result['f_id']

        plugin_info = ncTThirdPartyPluginInfo()
        plugin_info.thirdPartyId = result['f_app_id']
        plugin_info.filename = result['f_plugin_name']
        plugin_info.type = result['f_plugin_type']
        plugin_info.objectId = result['f_object_id']
        plugin_info.indexId = result['f_id']
        third_info.plugin = plugin_info

        try:
            json.loads(third_info.config)
            json.loads(third_info.internalConfig)
        except ValueError:
            raise_exception(exp_msg=_("IDS_THIRD_CONFIG_FAILED"),
                            exp_num=ncTShareMgntError.NCT_FAILED_THIRD_CONFIG)

        return third_info

    def get_third_party_config_by_appid(self, appId):
        """
        根据appid获取第三方配置
        """
        sql = """
        SELECT `f_app_id`,`f_app_name`,`f_config`, `f_internal_config`, `f_enable`,
               `f_plugin_name`, `f_plugin_type`, `f_object_id`, `f_id`
        FROM `t_third_party_auth`
        WHERE `f_app_id` = %s
        """
        result = self.r_db.one(sql, appId)
        if not result:
            raise_exception(exp_msg=_("IDS_THIRD_APPID_SECRET_FAILED"),
                            exp_num=ncTShareMgntError.NCT_INVALID_APPID_OR_APPKEY)

        third_info = ncTThirdPartyConfig()
        third_info.thirdPartyId = result['f_app_id']
        third_info.thirdPartyName = result['f_app_name']
        third_info.config = result['f_config']
        third_info.internalConfig = result['f_internal_config']
        third_info.enabled = True if result['f_enable'] else False
        third_info.indexId = result['f_id']

        plugin_info = ncTThirdPartyPluginInfo()
        plugin_info.thirdPartyId = result['f_app_id']
        plugin_info.filename = result['f_plugin_name']
        plugin_info.type = result['f_plugin_type']
        plugin_info.objectId = result['f_object_id']
        plugin_info.indexId = result['f_id']
        third_info.plugin = plugin_info

        try:
            json.loads(third_info.config)
            json.loads(third_info.internalConfig)
        except ValueError:
            raise_exception(exp_msg=_("IDS_THIRD_CONFIG_FAILED"),
                            exp_num=ncTShareMgntError.NCT_FAILED_THIRD_CONFIG)

        return third_info

    def get_third_party_config(self, pluginType):
        """
        获取第三方配置信息
        """
        # -1 为获取所有配置
        if (pluginType != -1):
            self.__check_third_party_plugin_type(pluginType, b_raise=True)

        sql = """
        SELECT f_app_id,f_app_name,f_config, f_internal_config, f_enable,
               f_plugin_name, f_plugin_type, f_object_id, f_id
        FROM t_third_party_auth
        WHERE f_plugin_type = %s
        """ % pluginType

        if pluginType == -1:
            sql = """
            SELECT f_app_id,f_app_name,f_config, f_internal_config, f_enable,
                   f_plugin_name, f_plugin_type, f_object_id, f_id
            FROM t_third_party_auth
            """

        results = self.r_db.all(sql)
        third_infos = []
        for result in results:
            third_info = ncTThirdPartyConfig()
            third_info.thirdPartyId = result['f_app_id']
            third_info.thirdPartyName = result['f_app_name']
            third_info.config = result['f_config']
            third_info.internalConfig = result['f_internal_config']
            third_info.enabled = True if result['f_enable'] else False
            third_info.indexId = result['f_id']

            plugin_info = ncTThirdPartyPluginInfo()
            plugin_info.thirdPartyId = result['f_app_id']
            plugin_info.filename = result['f_plugin_name']
            plugin_info.type = result['f_plugin_type']
            plugin_info.objectId = result['f_object_id']
            plugin_info.indexId = result['f_id']
            third_info.plugin = plugin_info

            try:
                json.loads(third_info.config)
                json.loads(third_info.internalConfig)
            except ValueError:
                raise_exception(exp_msg=_("IDS_THIRD_CONFIG_FAILED"),
                                exp_num=ncTShareMgntError.NCT_FAILED_THIRD_CONFIG)
            third_infos.append(third_info)

        return third_infos

    def add_third_party_config(self, config):
        """
        新增第三方配置
        """
        # 参数检查
        config.thirdPartyId, config.thirdPartyName = strip_whitespace(
            config.thirdPartyId, config.thirdPartyName)
        self.__check_third_party_config(config)
        self.__check_appid_exist(
            config.indexId, config.thirdPartyId, b_raise=True)

        # 暂时限制只能上传一个认证插件
        sql = """
        SELECT `f_app_id` from `t_third_party_auth` WHERE `f_plugin_type` = %s
        """
        result = self.r_db.one(sql, ncTPluginType.AUTHENTICATION)
        if config.plugin.type == ncTPluginType.AUTHENTICATION and result:
            raise_exception(exp_msg=_("IDS_THIRD_CONFIG_FAILED"),
                            exp_num=ncTShareMgntError.NCT_FAILED_THIRD_CONFIG)

        sql = """
        INSERT INTO `t_third_party_auth` (`f_app_id`, `f_app_name`, `f_enable`, `f_config`, `f_internal_config`, `f_plugin_type`, `f_plugin_name`, `f_object_id`, `f_oss_id`)
        VALUES (%s, %s, %s, %s, %s, %s, '', '', '')
        """
        self.w_db.query(sql, config.thirdPartyId, config.thirdPartyName,
                        config.enabled, config.config, config.internalConfig, config.plugin.type)

        sql = """
        SELECT `f_id` from `t_third_party_auth` WHERE `f_app_id` = %s
        """
        result = self.r_db.one(sql, config.thirdPartyId)
        indexId = result['f_id']
        return indexId

    def set_third_party_config(self, config):
        """
        设置第三方配置
        """
        config.thirdPartyId, config.thirdPartyName = strip_whitespace(
            config.thirdPartyId, config.thirdPartyName)
        self.__check_third_exists_by_indexid(config.indexId, b_raise=True)
        self.__check_third_party_config(config)
        self.__check_appid_exist(
            config.indexId, config.thirdPartyId, b_raise=True)

        sql = """
        UPDATE `t_third_party_auth`
        SET `f_app_name` = %s, `f_enable` = %s, `f_config` = %s, `f_internal_config` = %s, `f_plugin_type` = %s, `f_app_id` = %s
        WHERE `f_id` = %s
        """
        self.w_db.query(sql, config.thirdPartyName, config.enabled, config.config,
                        config.internalConfig, config.plugin.type, config.thirdPartyId, config.indexId)

    def delete_third_party_config(self, indexId):
        """
        删除第三方配置
        """
        self.__check_third_exists_by_indexid(indexId, b_raise=True)

        # 删除第三方插件
        self.delete_third_party_plugin(indexId)

        # 删除配置
        sql = """
        DELETE FROM `t_third_party_auth`
        WHERE `f_id` = %s
        """
        self.w_db.query(sql, indexId)

    def set_third_party_plugin(self, plugin):
        """
        插件持久化
        """
        # 参数检查
        self.__check_third_party_plugin_config(plugin)

        # 上传到存储
        ossId, objectId = self.upload_third_party_plugin(plugin)

        # 删除存储原有插件
        self.delete_third_party_plugin(ossId, objectId)

    def delete_third_party_plugin(self, ossId, objectId):
        """
        删除旧插件
        """
        # 没有旧插件, 直接返回
        if not ossId:
            return

        # 删除存储数据
        code, deleteinfo = self.ossgateway_driven.get_delete_info(ossId, objectId)
        if code == 403:
            if deleteinfo.get("code") == 403031020:
                ShareMgnt_Log("delete plugin error, %s", deleteinfo.get("message"))
                return
            ShareMgnt_Log(f'get delete url failed: {code},{deleteinfo.get("message")},{deleteinfo.get("cause")}')
            raise_exception(exp_msg=_("GET_URL_FAILED"),
                            exp_num=ncTShareMgntError.NCT_GET_URL_FAILED)
        elif code != 200:
            ShareMgnt_Log(f'get delete url failed: {code},{deleteinfo.get("message")},{deleteinfo.get("cause")}')
            raise_exception(exp_msg=_("GET_URL_FAILED"),
                                exp_num=ncTShareMgntError.NCT_GET_URL_FAILED)
        if not deleteinfo.get("url"):
            raise_exception(exp_msg=_("GET_URL_FAILED"),
                                exp_num=ncTShareMgntError.NCT_GET_URL_FAILED)

        rsp = requests.request(deleteinfo.get("method"),
                            deleteinfo.get("url"),
                            headers=deleteinfo.get("headers"),
                            verify=False)
        if ("[2" not in str(rsp)) and ("[404]" not in str(rsp)):
            raise_exception(exp_msg=rsp.text,
                                exp_num=ncTShareMgntError.NCT_DEL_PKG_INFO_FAILED)
        else:
            ShareMgnt_Log("Delete third party plugin success. objectId: %s" % objectId)


    def upload_third_party_plugin(self, plugin):
        """
        上传插件到存储
        """
        sql = """
        SELECT `f_oss_id`, `f_object_id`
        FROM `t_third_party_auth`
        WHERE `f_id` = %s
        """
        result = self.r_db.one(sql, plugin.indexId)
        oldobjectId = result['f_object_id']
        oldossId = result['f_oss_id']

        ossIds = []
        uploadOSSId = ""
        ossIds = self.ossgateway_driven.get_local_storages()
        if not ossIds:
            raise_exception(exp_msg="no available oss",
                                    exp_num=ncTShareMgntError.NCT_NO_AVAILABLE_OSS)
        else:
            for ossinfo in ossIds:
                if ossinfo.get("default"):
                    uploadOSSId = ossinfo.get("id")
            if not uploadOSSId:
                uploadOSSId = ossIds[0].get("id")
        objectId = str(uuid.uuid1()).replace('-', '')
        code, uploadinfo = self.ossgateway_driven.get_upload_info(uploadOSSId, objectId)
        if code != 200:
            raise_exception(exp_msg=f'get upload url failed: {code},{uploadinfo.get("message")},{uploadinfo.get("cause")}',
                            exp_num=ncTShareMgntError.NCT_UPLOAD_OSS_FAILED)
        rsp = requests.request(uploadinfo.get("method"),
                               uploadinfo.get("url"),
                               headers=uploadinfo.get("headers"),
                               data=plugin.data,
                               verify=False)
        if "[2" not in str(rsp):
            raise_exception(exp_msg=rsp.text,
                            exp_num=ncTShareMgntError.NCT_UPLOAD_OSS_FAILED)

        # 插件信息持久化到数据库
        sql = """
        UPDATE `t_third_party_auth`
        SET `f_plugin_type` = %s, `f_plugin_name` = %s, `f_object_id` = %s, `f_oss_id` = %s
        WHERE `f_app_id` = %s
        """
        self.w_db.query(sql, plugin.type, plugin.filename,
                        objectId, uploadOSSId, plugin.thirdPartyId)

        return oldossId, oldobjectId

    def download_third_party_plugin(self, thirdPartyId):
        """
        下载插件
        """
        sql = """
        SELECT `f_object_id`, `f_plugin_name`, `f_oss_id`
        FROM `t_third_party_auth`
        WHERE `f_app_id` = %s
        """
        result = self.r_db.one(sql, thirdPartyId)
        objectId = result['f_object_id']
        ossId = result['f_oss_id']

        code, downloadinfo = self.ossgateway_driven.get_download_info(ossId, objectId)
        if code == 403:
            if downloadinfo.get("code") == 403031020:
                # 如果对象存储被禁用，则不删除原对象存储中的旧插件
                ShareMgnt_Log("please upload plugin again, %s", downloadinfo.get("message"))
            else:
                ShareMgnt_Log(f'get download url failed: {code},{downloadinfo.get("message")},{downloadinfo.get("cause")}')
            raise_exception(exp_msg=_("GET_URL_FAILED"),
                        exp_num=ncTShareMgntError.NCT_GET_URL_FAILED)
        elif code // 100 == 5:
            raise_exception(exp_msg="ossgateway exception",
                        exp_num=ncTShareMgntError.NCT_NO_AVAILABLE_OSS)
        elif code != 200:
            ShareMgnt_Log(f'get download url failed: {code},{downloadinfo.get("message")},{downloadinfo.get("cause")}')
            raise_exception(exp_msg=_("GET_URL_FAILED"),
                        exp_num=ncTShareMgntError.NCT_GET_URL_FAILED)
        if not downloadinfo.get("url"):
            ShareMgnt_Log("get url failed ")
            raise_exception(exp_msg=_("GET_URL_FAILED"),
                        exp_num=ncTShareMgntError.NCT_GET_URL_FAILED)
        ShareMgnt_Log("get plugin url address: %s", downloadinfo.get("url"))
        r = requests.get(downloadinfo.get("url"), verify=False)
        ShareMgnt_Log("get plugin data code: %s", r.status_code)
        if r.status_code // 100 == 5:
            raise_exception(exp_msg="storage exception",
                        exp_num=ncTShareMgntError.NCT_NO_AVAILABLE_OSS)
        return r.content

    def get_plugin_type_str(self, pluginType):
        """
        根据插件类型获取目录名用于拼接插件存放路径
        """
        pluginTypeStr = ""
        if pluginType == ncTPluginType.AUTHENTICATION:
            pluginTypeStr = "auth"
        else:
            raise_exception(exp_msg=_("IDS_THIRD_CONFIG_FAILED"),
                            exp_num=ncTShareMgntError.NCT_FAILED_THIRD_CONFIG)
        return pluginTypeStr

    def get_format_plugin_path(self, config, fileName):
        """
        获取插件路径
        """
        THIRD_PARTY_PLUGIN_PATH_ROOT = "/sysvol/plugin"
        pluginTypeStr = self.get_plugin_type_str(config.plugin.type)
        # auth-1
        import_path = "%s_%s" % (pluginTypeStr, str(config.plugin.indexId))
        # /sysvol/plugin/auth-1/auth_module.py
        module_path = """%s/%s_%s/%s""" % (THIRD_PARTY_PLUGIN_PATH_ROOT,
                                           pluginTypeStr, str(config.plugin.indexId), fileName)
        return import_path, module_path

    def __check_appid_exist(self, indexId, appId, b_raise=False):
        """
        检查appid冲突
        """
        # 检查indexId
        if indexId:
            sql = """
            SELECT `f_id`
            FROM `t_third_party_auth`
            WHERE `f_id` = %s
            """
            r = self.r_db.one(sql, indexId)
            if not r:
                raise_exception(exp_msg=_("IDS_THIRD_CONFIG_FAILED"),
                                exp_num=ncTShareMgntError.NCT_FAILED_THIRD_CONFIG)

        if not indexId:
            indexId = -1

        sql = """
        SELECT `f_config`
        FROM `t_third_party_auth`
        WHERE `f_id` != %s AND `f_app_id` = %s
        """
        db_config = self.r_db.one(sql, indexId, appId)
        if b_raise and db_config:
            raise_exception(exp_msg=_("APPID_HAS_EXIST"),
                            exp_num=ncTShareMgntError.APPID_HAS_EXIST)
        else:
            return True if db_config else False

    def __check_third_party_plugin_config(self, pluginConfig):
        """
        检查插件参数
        """
        correct = True
        if not pluginConfig.thirdPartyId or not pluginConfig.filename or not pluginConfig.data:
            correct = False

        if pluginConfig.type != ncTPluginType.AUTHENTICATION:
            correct = False

        if not correct:
            raise_exception(exp_msg=_("IDS_THIRD_CONFIG_FAILED"),
                            exp_num=ncTShareMgntError.NCT_FAILED_THIRD_CONFIG)

    def __check_third_party_config(self, config):
        """
        第三方参数检查
        """
        correct = True
        if not config.thirdPartyId or not config.thirdPartyName or not config.plugin or not config.config or not config.internalConfig:
            raise_exception(exp_msg=_("IDS_THIRD_CONFIG_FAILED"),
                            exp_num=ncTShareMgntError.NCT_FAILED_THIRD_CONFIG)

        correct = self.__check_third_party_plugin_type(config.plugin.type)

        if correct:
            try:
                json.loads(config.config)
                json.loads(config.internalConfig)
            except ValueError:
                correct = False

        if not correct:
            raise_exception(exp_msg=_("IDS_THIRD_CONFIG_FAILED"),
                            exp_num=ncTShareMgntError.NCT_FAILED_THIRD_CONFIG)

    def __check_third_exists_by_indexid(self, index_id, b_raise=False):
        """
        根据indexId检查第三方存不存在
        """
        sql = """
        SELECT `f_config`
        FROM `t_third_party_auth`
        WHERE `f_id` = %s
        """
        db_config = self.r_db.one(sql, index_id)
        if b_raise and not db_config:
            raise_exception(exp_msg=_("APPID_NOT_EXIST"),
                            exp_num=ncTShareMgntError.APPID_NOT_EXIST)
        else:
            return True if db_config else False

    def __check_third_party_plugin_type(self, pluginType, b_raise=False):
        """
        检查第三方插件类型
        """
        if pluginType != ncTPluginType.AUTHENTICATION:
            if b_raise:
                raise_exception(exp_msg=_("IDS_THIRD_CONFIG_FAILED"),
                                exp_num=ncTShareMgntError.NCT_FAILED_THIRD_CONFIG)
            else:
                return False
        return True

    def unload_plugin(self, module_name):
        """
        卸载已加载的插件模块
        """
        if module_name in sys.modules:
            del sys.modules[module_name]
