#!/usr/bin/python3
# -*- coding:utf-8 -*-
import os
import sys
import json
import threading
import time
from src.common import global_info
from src.common.sharemgnt_logger import ShareMgnt_Log
from src.common.lib import raise_exception
from src.common.db.connector import DBConnector
from src.common.redis_connector import OPRedis
from src.common.plugin_lock import PluginVersion
from src.modules.domain_manage import DomainManage
from src.modules.third_db_manage import ThirdDBManage
from src.third_party_auth.third_config_manage import ThirdConfigManage
from src.third_party_auth.ou_manage.as_sql import *
from src.third_party_auth.ou_syncer.base_syncer import BaseSyncer

from EThriftException.ttypes import ncTException
from ShareMgnt.ttypes import ncTShareMgntError, ncTSyncType

SYNC_RETRY_LOCK = threading.Lock()
SYNC_RETRY_WAIT_TIME = 1.5

class ThirdSyncManage(DBConnector):
    """
    第三方同步管理类
    """
    thread_dict = {}

    def __init__(self):
        """
        初始化
        """
        self.third_config_manager = ThirdConfigManage()
        self.opredis = OPRedis()
        self.plugin_version = PluginVersion()

    def __get_sync_module(self, app_id):
        """
        根据app_id获取同步组件
        """
        sync_class = None
        ou_class = None
        sync_type = ""

        try:
            # 二进制服务运行，从third_party_auth导入
            import_ou_str = "from third_party_auth.ou_manage import *"
            import_syncer_str = "from third_party_auth.ou_syncer import *"
            exec(import_ou_str)
            exec(import_syncer_str)
        except Exception as ex:
            # 源码调试时，需要从src目录导入
            import traceback
            ShareMgnt_Log("ImportError_traceback:%s", traceback.format_exc())
            ShareMgnt_Log("ImportError:%s", str(ex))
            ShareMgnt_Log("pythonpath:%s", str(sys.path))
            import_src_ou_str = "from src.third_party_auth.ou_manage import *"
            import_src_syncer_str = "from src.third_party_auth.ou_syncer import *"
            exec(import_src_ou_str)
            exec(import_src_syncer_str)

        try:
            domain_info = DomainManage().get_domain_by_id(app_id)
            if domain_info and domain_info.status:
                sync_class = locals()["DomainSyncer"]
                ou_class = locals()["DomainOuManage"]
                sync_type = ncTSyncType.DOMAIN_SYNC
                return sync_class, ou_class, sync_type
        except ncTException as ex:
            if ex.errID != ncTShareMgntError.NCT_DOMAIN_NOT_EXIST:
                raise ex

        third_db_info = ThirdDBManage().get_third_db_info(app_id)
        if third_db_info and third_db_info.status:
            sync_class = locals()["ThirdDbSyncer"]
            ou_class = locals()["ThirdDBOuManage"]
            sync_type = ncTSyncType.THIRD_DB_SYNC
            return sync_class, ou_class, sync_type

        third_info = self.third_config_manager.get_third_party_config_by_appid(app_id)
        if third_info and third_info.enabled and third_info.config:

            # 导入第三方插件认证模块
            import_path, module = self.third_config_manager.get_format_plugin_path(third_info, "ou_module.py")

            # 添加插件路径到环境变量
            plugin_path = "/sysvol/plugin/%s" % import_path
            if plugin_path not in sys.path:
                sys.path.append(plugin_path)    

            # 获取数据库中对象存储的版本
            auth_store_version = third_info.plugin.objectId

            # 1 本地版本不存在, 2 本地版本与对象存储中版本不一致, 3 本地插件不存在, 则从对象存储中下载插件.
            if not PluginVersion.AUTH_LOCAL_VERSION \
                or (auth_store_version != PluginVersion.AUTH_LOCAL_VERSION) \
                or (not os.path.exists(plugin_path)):
                try:
                    # 复制插件到对应路径，会覆盖原先内容
                    third_info.plugin.data = self.third_config_manager.download_third_party_plugin(third_info.thirdPartyId)
                    if third_info.plugin.data:
                        from src.third_party_auth.third_party_manage import ThirdPartyManage
                        ThirdPartyManage().add_local_third_party_plugin(third_info.plugin)
                except ncTException as ex:
                    ShareMgnt_Log("The sync plugin is download failed")                          

            # 导入第三方插件组织结构解析模块
            if os.path.exists(module):
                # 卸载已加载的第三方同步插件
                self.third_config_manager.unload_plugin("ou_module")

                # redis中锁的key
                POD_ID = os.getenv("POD_IP", "127.0.0.1")
                key = POD_ID + " " + str(third_info.plugin.type)

                # 获取锁
                if not self.opredis.get_lock(key, 30):
                    import_str = "from ou_module import *"
                    exec(import_str)        

            config = json.loads(third_info.config)
            config.update(json.loads(third_info.internalConfig))
            if "syncModule" in config:
                sync_class = locals()[config['syncModule']]

            if "ouModule" in config:
                ou_class = locals()[config['ouModule']]

            sync_type = ncTSyncType.THIRD_SYNC

        return sync_class, ou_class, sync_type

    def register_thread(self, app_id, sync_thread):
        """
        注册线程
        """
        self.thread_dict[str(app_id)] = sync_thread

    def remove_thread(self, app_id):
        """
        移除线程
        """
        if app_id in self.thread_dict:
            del self.thread_dict[app_id]

    def start_sync(self, app_id, auto_sync=True, update_config=False):
        """
        开启第三方同步
        Args:
            app_id: 第三方：app id
                    域同步：域id
                    数据库同步：第三方数据库id  
            auto_sync:
                    True: 定期自动同步
                    False: 仅同步一次
        """    
        # 第三方同步
        try:
            # 获取第三方插件的一些配置
            third_info = self.third_config_manager.get_third_party_config_by_appid(app_id)

            # 获取数据库中对象存储的版本
            auth_store_version = third_info.plugin.objectId

            # 本地版本存在, 本地版本与对象存储中版本一致且appid的自动同步已开启 则不重新加载插件.
            if PluginVersion.AUTH_LOCAL_VERSION \
                and (auth_store_version == PluginVersion.AUTH_LOCAL_VERSION) \
                and (app_id in self.thread_dict) \
                and (not update_config):
                self.thread_dict[app_id].sync_immediately()
                return

        except ncTException as ex:
            if (ex.errID == ncTShareMgntError.NCT_INVALID_APPID_OR_APPKEY) and (app_id in self.thread_dict):
                self.thread_dict[app_id].sync_immediately()
                return
        
        # 获取同步组件
        sync_class = None
        ou_class = None
        try:
            sync_class, ou_class, sync_type = self.__get_sync_module(app_id)
        except Exception as ex:
            ShareMgnt_Log("appid: %s, get sync module failed: %s", app_id, str(ex))

        # 检测同步类和ou类是否设置
        if not sync_class:
            ShareMgnt_Log("appid: %s, syncModule not set", app_id)
            return
        if not ou_class:
            ShareMgnt_Log("appid: %s, ouModule not set", app_id)
            return
        
        # 同步类型是第三方同步且已经存在同步的线程
        retry = False
        if sync_type == ncTSyncType.THIRD_SYNC and global_info.THIRD_SYNC_THREAD:
            # 关闭同步的线程
            global_info.THIRD_SYNC_THREAD.close()
            if global_info.THIRD_SYNC_THREAD.is_syncing():
                with SYNC_RETRY_LOCK:
                    global_info.RETRY_SYNC_THREAD["THIRD"] = app_id
                retry = True
            else:
                # 移除线程的app_id
                self.remove_thread(app_id)

        # 同步类型是域同步且该域已经存在同步的线程
        if sync_type == ncTSyncType.DOMAIN_SYNC and global_info.DOMAIN_SYNC_THREAD.get(app_id):
            global_info.DOMAIN_SYNC_THREAD[app_id].close()
            if global_info.DOMAIN_SYNC_THREAD[app_id].is_syncing():
                with SYNC_RETRY_LOCK:
                    global_info.RETRY_SYNC_THREAD["DOMAIN"].add(global_info.DOMAIN_SYNC_THREAD[app_id])
                retry = True
            else:
                self.remove_thread(app_id)

        # 同步类型是第三方数据库同步且该数据库已经存在同步的线程
        if sync_type == ncTSyncType.THIRD_DB_SYNC and global_info.THIRD_DB_SYNC_THREAD.get(app_id):
            global_info.THIRD_DB_SYNC_THREAD[app_id].close()
            if global_info.THIRD_DB_SYNC_THREAD[app_id].is_syncing():
                with SYNC_RETRY_LOCK:
                    global_info.RETRY_SYNC_THREAD["THIRD_DB"].add(global_info.THIRD_DB_SYNC_THREAD[app_id])
                retry = True
            else:
                self.remove_thread(app_id)

        # 保证同步只会同时出现一次
        if not retry:
            sync_thread = SyncThread(app_id, sync_class, ou_class, self.remove_thread, auto_sync)
            sync_thread.daemon = True
            sync_thread.start()

            # 同步类型是第三方同步
            if sync_type == ncTSyncType.THIRD_SYNC:
                # 将同步线程对象存入全局变量, 为了后面需要关闭同步线程的时候可以定位到
                global_info.THIRD_SYNC_THREAD = sync_thread

            # 同步类型是域同步
            if sync_type == ncTSyncType.DOMAIN_SYNC:
                global_info.DOMAIN_SYNC_THREAD[app_id] = sync_thread
            
            # 同步类型是第三方数据库
            if sync_type == ncTSyncType.THIRD_DB_SYNC:
                global_info.THIRD_DB_SYNC_THREAD[app_id] = sync_thread    

            # 自动同步，注册当前appid对应的线程
            self.register_thread(app_id, sync_thread)
            ShareMgnt_Log("自动同步线程 appid: %s 已开启并注册到管理器", app_id)  

class SyncThread(threading.Thread):
    """
    第三方同步线程类
    """
    def __init__(self, app_id, sync_class, ou_class, call_back, auto_sync=True):
        """
        初始化
        """
        super(SyncThread, self).__init__()
        self.syncer = sync_class(app_id, ou_class)
        self.app_id = app_id
        self.auto_sync = auto_sync
        self.terminate = False
        self.evt = threading.Event()
        self.call_back = call_back

    def close(self):
        """
        关闭
        """
        self.terminate = True

    def is_syncing(self):
        """
        是否正在同步
        """
        if self.evt.isSet():
            return True
        else:
            return False

    def sync_immediately(self):
        """
        立刻同步
        """
        if not self.evt.isSet():
            self.evt.set()

    def run(self):
        """
        """
        if not self.syncer:
            ShareMgnt_Log(">>>>>>>>>>>>>>>>>>>>>>> appid: %s, 加载同步器失败 <<<<<<<<<<<<<<<<<<<<<<<<<<<", self.app_id)
            return

        ShareMgnt_Log(">>>>>>>>>>>>>>>>>>>>>>> appid: %s, 同步线程开启 <<<<<<<<<<<<<<<<<<<<<<<<<<", self.app_id)

        cnt = 1
        while self.auto_sync:
            interval_time = 1800
            try:
                interval_time = self.syncer.get_sync_interval()
                sync_status = self.syncer.get_sync_status()
                if not sync_status:
                    break

                if self.terminate:
                    break

                ShareMgnt_Log(">>>>>>>>>>>>>>>>>>>>>>>>> appid: %s, 第%s次同步开始 <<<<<<<<<<<<<<<<<<<<<<<<<<<", self.app_id, str(cnt))
                self.evt.set()
                self.syncer.sync()

                ShareMgnt_Log(">>>>>>>>>>>>>>>>>>>>>>>>> appid: %s, 第%s次同步结束 <<<<<<<<<<<<<<<<<<<<<<<<<<<", self.app_id, str(cnt))

                cnt += 1

            except Exception as ex:
                ShareMgnt_Log("run， 异常：%s" % (str(ex)))

            if self.terminate:
                break

            self.evt.clear()
            self.evt.wait(interval_time)
            self.evt.set()

        else:
            ShareMgnt_Log(">>>>>>>>>>>>>>>>>>>>>>>>> appid: %s, 第%s次同步开始 <<<<<<<<<<<<<<<<<<<<<<<<<<<", self.app_id, str(cnt))
            self.evt.set()
            self.syncer.sync()
            ShareMgnt_Log(">>>>>>>>>>>>>>>>>>>>>>>>> appid: %s, 第%s次同步结束 <<<<<<<<<<<<<<<<<<<<<<<<<<<", self.app_id, str(cnt))
        
        self.evt.clear()
        # 回调，从third_sync_manage中移除当前线程
        self.call_back(self.app_id)
        ShareMgnt_Log(">>>>>>>>>>>>>>>>>>>>>>>>>> appid:%s, 同步线程结束 <<<<<<<<<<<<<<<<<<<<<<<<<<<", self.app_id)

class SyncRetryThread(threading.Thread):

    def __init__(self):
        """
        初始化
        """
        super(SyncRetryThread, self).__init__()

    def sync_retry_task(self):
        """
        同步过程中，客户触发同步待上次同步完成后，重新触发一次同步
        """
        if global_info.RETRY_SYNC_THREAD["THIRD"] or global_info.RETRY_SYNC_THREAD["THIRD_DB"] or global_info.RETRY_SYNC_THREAD["DOMAIN"]:
            try:
                third_sync_manage = ThirdSyncManage()
                if global_info.RETRY_SYNC_THREAD["THIRD"] and not global_info.THIRD_SYNC_THREAD.is_syncing():
                    third_sync_manage.start_sync(global_info.RETRY_SYNC_THREAD["THIRD"])
                    with SYNC_RETRY_LOCK:
                        global_info.RETRY_SYNC_THREAD["THIRD"] = ""
                if global_info.RETRY_SYNC_THREAD["THIRD_DB"]:
                    for id in global_info.RETRY_SYNC_THREAD["THIRD_DB"]:
                        if not global_info.THIRD_DB_SYNC_THREAD[id].is_syncing():
                            third_sync_manage.start_sync(id)
                            with SYNC_RETRY_LOCK:
                                global_info.RETRY_SYNC_THREAD["THIRD_DB"].remove(id)
                if global_info.RETRY_SYNC_THREAD["DOMAIN"]:
                    for id in global_info.RETRY_SYNC_THREAD["DOMAIN"]:
                        if not global_info.DOMAIN_SYNC_THREAD[id].is_syncing():
                            third_sync_manage.start_sync(id)
                            with SYNC_RETRY_LOCK:
                                global_info.RETRY_SYNC_THREAD["DOMAIN"].remove(id)
            except Exception as e:
                ShareMgnt_Log("sync retry task recording thread run error: %s", str(e))
            
    def run(self):
        """
        执行
        """
        ShareMgnt_Log("**************** sync retry task recording thread start *****************")

        while True:
            self.sync_retry_task()
            time.sleep(SYNC_RETRY_WAIT_TIME)
