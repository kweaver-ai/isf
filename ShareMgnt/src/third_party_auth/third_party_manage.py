# coding: utf-8
"""
第三方系统管理模块
"""
import os
import sys
import time
import tarfile
import shutil
import traceback
import subprocess
from thrift.transport import TTransport
from eisoo.tclients import TClient
from thrift.Thrift import TException
from src.common.lib import (raise_exception, check_filename, strip_whitespace)
from src.common.sharemgnt_logger import ShareMgnt_Log
from src.common.db.connector import DBConnector
from src.common.plugin_lock import PluginVersion
from src.common.redis_connector import OPRedis
from src.modules.domain_manage import DomainManage
from src.modules.third_db_manage import ThirdDBManage
from src.third_party_auth.third_config_manage import ThirdConfigManage
from src.third_party_auth.third_auth_manage import ThirdAuthManage
from src.third_party_auth.third_sync_manage import ThirdSyncManage
from ShareMgnt.ttypes import (ncTShareMgntError, ncTPluginType)
from EThriftException.ttypes import ncTException

THIRD_PARTY_PLUGIN_PATH = "/sysvol/plugin/"


class ThirdPartyManage(DBConnector):
    """
    第三方管理模块
    """

    def __init__(self):
        """
        """
        self.third_config_manage = ThirdConfigManage()
        self.third_auth_manage = ThirdAuthManage()
        self.third_sync_mange = ThirdSyncManage()
        self.domain_manage = DomainManage()
        self.third_db_manage = ThirdDBManage()
        self.plugin_version = PluginVersion()
        self.opredis = OPRedis()

        # 添加插件环境变量
        self.add_third_party_plugin_path()

    def add_third_party(self, config):
        """
        开启第三方
        """
        indexId = self.third_config_manage.add_third_party_config(config)
        # 同步插件则开始同步
        if config.plugin.type == ncTPluginType.AUTHENTICATION:
            self.third_sync_mange.start_sync(config.thirdPartyId)

        return indexId

    def set_third_party(self, config):
        """
        设置第三方
        """
        self.third_config_manage.set_third_party_config(config)
        # 同步插件则开始同步
        if config.plugin.type == ncTPluginType.AUTHENTICATION:
            self.third_sync_mange.start_sync(config.thirdPartyId, update_config=True)

    def start_third_sync(self, app_id, auto_sync=True):
        """
        根据appid启动第三方同步
        """
        self.third_sync_mange.start_sync(app_id, auto_sync)

    def start_all_third_sync(self):
        """
        启动第三方同步
        """
        try:
            third_info = self.third_config_manage.get_third_party_info_auth()
            if third_info.thirdPartyId:
                self.third_sync_mange.start_sync(third_info.thirdPartyId)
        except Exception:
            ShareMgnt_Log(traceback.format_exc())

    def start_all_domain_sync(self):
        """
        启动所有域自动同步
        """
        try:
            # 获取所有的域id
            all_domains = self.domain_manage.get_all_domains()

            for domain in all_domains:
                if domain.status and domain.syncStatus == 0:
                    self.third_sync_mange.start_sync(domain.id, True)
        except Exception:
            ShareMgnt_Log(traceback.format_exc())

    def start_all_third_db_sync(self):
        """
        启动所有第三方数据库同步
        """
        try:
            third_db_ids = self.third_db_manage.get_enable_third_db_ids()
            for db_id in third_db_ids:
                self.third_sync_mange.start_sync(db_id, True)
        except Exception:
            ShareMgnt_Log(traceback.format_exc())

    def add_global_third_party_plugin(self, plugin_info):
        """
        向所有节点添加第三方认证插件
        """
        # 去除参数的空格
        plugin_info.thirdPartyId, plugin_info.filename = strip_whitespace(plugin_info.thirdPartyId, plugin_info.filename)

        # 文件名合法性检测
        check_filename(plugin_info.filename)

        # 第三方插件持久化
        self.third_config_manage.set_third_party_plugin(plugin_info)

        # 同步插件则开始同步
        if plugin_info.type == ncTPluginType.AUTHENTICATION:
            self.third_sync_mange.start_sync(plugin_info.thirdPartyId, update_config=True)

    def __extract_plugin(self, filepath):
        '''
        安全解压tar.gz文件，防止被解压的文件名中包含../这类相对路径定位，导致文件被覆盖
        '''
        resolve = lambda x: os.path.realpath(os.path.abspath(x))

        def bad_path(base, path):
            return not resolve(os.path.join(base, path)).startswith(base)

        def bad_link(base, info):
            tip = resolve(os.path.join(base, os.path.dirname(info.name)))
            return bad_path(tip, info.linkname)

        def safemembers(members):
            '''
            排除非法的路径
            '''
            base = THIRD_PARTY_PLUGIN_PATH
            for info in members:
                if bad_path(base, info.name):
                    raise Exception("illegal path")
                elif info.issym() and bad_link(base, info):
                    raise Exception("illegal path")
                elif info.islnk() and bad_link(base, info):
                    raise Exception("illegal path")
                else:
                    yield info

        try:
            tar = tarfile.open(filepath, "r:gz")
            tar.extractall(path=os.path.dirname(resolve(filepath)), members=safemembers(tar))
        except Exception:
            raise_exception(exp_msg=_("IDS_INVALID_THIRD_PARTY_PLUGIN"),
                            exp_num=ncTShareMgntError.NCT_INVALID_THIRD_PARTY_PLUGIN)
        finally:
            # 关闭文件
            try:
                tar.close()
            except UnboundLocalError:
                pass

            # 删除插件包
            try:
                os.remove(filepath)
            except OSError:
                pass

    def add_local_third_party_plugin(self, plugin_info):
        """
        向单个节点添加第三方认证插件
        """       
        # 去除参数的空格
        plugin_info.thirdPartyId, plugin_info.filename = strip_whitespace(plugin_info.thirdPartyId, plugin_info.filename)

        config = self.third_config_manage.get_third_party_config_by_appid(plugin_info.thirdPartyId)

        # 检查插件类型合法
        if config.plugin.type != plugin_info.type:
            raise_exception(exp_msg=_("IDS_INVALID_THIRD_PARTY_PLUGIN"),
                            exp_num=ncTShareMgntError.NCT_INVALID_THIRD_PARTY_PLUGIN)

        plugin_type_str = self.third_config_manage.get_plugin_type_str(config.plugin.type)
     
        # redis中锁的key
        POD_ID = os.getenv("POD_IP", "127.0.0.1")
        key = POD_ID + " " + str(plugin_info.type)

        # 上锁
        if self.opredis.set_lock(key, 0, 30):
            # 第三方插件路径，类似: /sysvol/plugin/auth-1/auth_module.py
            file_path = """%s/%s_%s/%s""" % (THIRD_PARTY_PLUGIN_PATH, plugin_type_str, str(config.plugin.indexId), plugin_info.filename)
            dir_path = """%s/%s_%s""" % (THIRD_PARTY_PLUGIN_PATH, plugin_type_str, str(config.plugin.indexId))

            # 清空插件目录
            if os.path.exists(dir_path):
                shutil.rmtree(dir_path)

            # 检查插件安装路径，不存在则创建
            if not os.path.exists(dir_path):
                os.makedirs(dir_path)

            # 复制插件到对应路径，会覆盖原先内容
            with open(file_path, 'wb+') as file:
                file.write(plugin_info.data)

            # 创建 __init__.py 文件用于导入模块
            init_file_name = "%s/%s_%s/__init__.py" % (THIRD_PARTY_PLUGIN_PATH, plugin_type_str, str(config.plugin.indexId))
            fd = open(init_file_name, 'a')
            fd.close()

            self.__extract_plugin(file_path)

            # 将数据库的对象存储版本覆盖到本地版本
            if plugin_info.type == 0:
                self.plugin_version.update_auth_local_version(plugin_info.objectId)
            elif plugin_info.type == 1:
                self.plugin_version.update_msg_local_version(plugin_info.objectId)

            # 解锁
            self.opredis.del_redis(key)

    def add_third_party_plugin_path(self):
        """
        添加 auth_manage 和 ou_manage 路径给第三方插件调用
        """
        paths = []
        paths.append(os.path.realpath(os.path.join(os.path.dirname(sys.argv[0]), "third_party_auth/auth_manage")))
        paths.append(os.path.realpath(os.path.join(os.path.dirname(sys.argv[0]), "third_party_auth/ou_manage")))

        # 添加第三方插件路径
        paths.append("/sysvol/plugin")

        for path in paths:
            if path not in sys.path:
                sys.path.append(path)

    def sync_third_party_plugin(self):
        """
        新增节点时同步第三方插件
        """
        # 读取配置，检查所有插件是否已经加载
        configs = self.third_config_manage.get_third_party_config(-1)
        for config in configs:
            # 检查数据库中是否有第三方插件
            if self.third_config_manage.get_third_party_plugin_status(config.thirdPartyId):
                # 复制插件到对应路径，会覆盖原先内容
                while True:
                    try:
                        config.plugin.data = self.third_config_manage.download_third_party_plugin(config.thirdPartyId)
                        # 退出重试
                        break
                    except ncTException as ex:
                        if ex.errID != 21107:
                            ShareMgnt_Log("run, 异常：%s" % (str(ex)))
                            # 重试间隔2s
                            time.sleep(2)
                            continue
                        else:
                            break
            
                self.add_local_third_party_plugin(config.plugin)
        
        ShareMgnt_Log("downloading plugin completed")

        try:
            ShareMgnt_Log("Starting sync third party plugin...")
            # 开启第三方自动同步
            self.start_all_third_sync()
        except Exception as ex:
            ShareMgnt_Log("sync third party plugin failed:%s", str(ex))

    def start_third_plugin_service(self):
        """
        启动第三方插件特定服务脚本
        """
        THIRD_PARTY_PLUGIN_SERVICE_PATH = "/sysvol/plugin/third_plugin_service.py"

        # 检测脚本是否存在
        if not os.path.exists(THIRD_PARTY_PLUGIN_SERVICE_PATH):
            ShareMgnt_Log(THIRD_PARTY_PLUGIN_SERVICE_PATH + " not exists")
        else:
            try:
                subprocess.Popen(["python", THIRD_PARTY_PLUGIN_SERVICE_PATH], close_fds=True)
            except Exception:
                ShareMgnt_Log(traceback.format_exc())
