#!/usr/bin/python3
# -*- coding:utf-8 -*-
"""
文档抓取管理类
"""
import re
from src.common.db.connector import DBConnector
from src.common.db.db_manager import get_db_name
from src.modules.config_manage import ConfigManage
from src.modules.user_manage import UserManage
from ShareMgnt.ttypes import (ncTFileCrawlConfig, ncTShareMgntError)
from src.common.lib import (raise_exception, escape_key, check_start_limit)


class FileCrawlManage(DBConnector):
    def __init__(self):
        """
        """
        self.config_manage = ConfigManage()
        self.user_manage = UserManage()

    def __check_crawl_file_type(self, fileCrawlType):
        """
        检查文件抓取类型
        """
        # 空格间隔的后缀类型，如：".txt .doc .exe"
        pattern = r'^\.([a-zA-Z0-9]+)(\s\.[a-zA-Z0-9]+)*$'
        if not fileCrawlType or not re.match(pattern, fileCrawlType):
            raise_exception(exp_msg=_("INVALID_PARAM") % fileCrawlType,
                            exp_num=ncTShareMgntError.NCT_INVALID_PARAMTER)

    def __check_custom_doc_exist(self, docId):
        """
        检查文档库是否存在
        """
        sql = f"""
        SELECT `f_doc_id` FROM `{get_db_name('anyshare')}`.`t_acs_doc`
        WHERE `f_doc_id` = %s AND `f_doc_type` = 3
        """
        result = self.r_db.one(sql, docId)
        if not result:
            raise_exception(exp_msg=_("IDS_CUSTOM_DOC_NOT_EXIST"),
                            exp_num=ncTShareMgntError.NCT_CUSTOM_DOC_NOT_EXIST)

    def __check_file_crawl_config_exist(self, userId):
        """
        检查文件抓取策略是否存在
        """
        sql = """
        SELECT `f_strategy_id` FROM `t_file_crawl_strategy`
        WHERE `f_user_id` = %s
        """
        result = self.r_db.one(sql, userId)
        if result:
            raise_exception(exp_msg=_("IDS_FILE_CRAWL_STRATEGY_EXIST"),
                            exp_num=ncTShareMgntError.NCT_FILE_CRAWL_STRATEGY_EXIST)

    def __check_another_file_crawl_config_exist(self, userId, strategyId):
        """
        检查相同用户文件抓取策略是否存在
        """
        sql = """
        SELECT `f_strategy_id` FROM `t_file_crawl_strategy`
        WHERE `f_user_id` = %s AND `f_strategy_id` != %s
        """
        result = self.r_db.one(sql, userId, strategyId)
        if result:
            raise_exception(exp_msg=_("IDS_FILE_CRAWL_STRATEGY_EXIST"),
                            exp_num=ncTShareMgntError.NCT_FILE_CRAWL_STRATEGY_EXIST)

    def set_file_crawl_status(self, status):
        """
        设置文档抓取总开关
        """
        self.config_manage.set_config('file_crawl_status', int(status))

    def get_file_crawl_status(self):
        """
        获取文档抓取总开关
        """
        return bool(int(self.config_manage.get_config("file_crawl_status")))

    def add_file_crawl_config(self, fileCrawlConfig):
        """
        新建文档抓取配置
        """
        self.__check_crawl_file_type(fileCrawlConfig.fileCrawlType)
        self.user_manage.check_user_exists(fileCrawlConfig.userId)
        self.__check_custom_doc_exist(fileCrawlConfig.docId)
        self.__check_file_crawl_config_exist(fileCrawlConfig.userId)

        insert_sql = """
        INSERT INTO `t_file_crawl_strategy`
        (`f_user_id`, `f_doc_id`, `f_file_crawl_type`)
        VALUES (%s, %s, %s)
        """
        self.w_db.query(insert_sql,
                        fileCrawlConfig.userId,
                        fileCrawlConfig.docId,
                        self.w_db.escape(fileCrawlConfig.fileCrawlType))

        sql = """
        SELECT `f_strategy_id` FROM `t_file_crawl_strategy`
        WHERE `f_user_id` = %s
        """
        result = self.r_db.one(sql, fileCrawlConfig.userId)
        strategyId = result['f_strategy_id']

        return strategyId

    def set_file_crawl_config(self, fileCrawlConfig):
        """
        设置文档抓取配置
        """
        self.__check_crawl_file_type(fileCrawlConfig.fileCrawlType)
        self.user_manage.check_user_exists(fileCrawlConfig.userId)
        self.__check_custom_doc_exist(fileCrawlConfig.docId)
        self.__check_another_file_crawl_config_exist(fileCrawlConfig.userId, fileCrawlConfig.strategyId)

        update_sql = """
        UPDATE `t_file_crawl_strategy`
        SET `f_user_id` = %s, `f_doc_id`= %s, `f_file_crawl_type` = %s
        WHERE `f_strategy_id` = %s
        """
        self.w_db.query(update_sql,
                        fileCrawlConfig.userId,
                        fileCrawlConfig.docId,
                        self.w_db.escape(fileCrawlConfig.fileCrawlType),
                        fileCrawlConfig.strategyId)

    def delete_file_crawl_config(self, strategyId):
        """
        删除文档抓取配置
        """
        delete_sql = """
        DELETE FROM `t_file_crawl_strategy`
        WHERE `f_strategy_id` = %s
        """
        self.w_db.query(delete_sql, strategyId)

    def get_file_crawl_config_count(self):
        """
        获取文档抓取配置数
        """
        sql = """
        SELECT COUNT(*) AS cnt FROM t_file_crawl_strategy
        """
        result = self.r_db.one(sql)
        return result['cnt']

    def get_file_crawl_config(self, start, limit):
        """
        分页获取文档抓取配置
        """
        limit_statement = check_start_limit(start, limit)
        sql = f"""
        SELECT `f_file_crawl_type`, `f_strategy_id`,
        `t_file_crawl_strategy`.`f_user_id`,
        `t_file_crawl_strategy`.`f_doc_id`,
        `t_user`.`f_login_name` AS `f_login_name`,
        `t_user`.`f_display_name` AS `f_display_name`,
        `t_acs_doc`.`f_name` AS `f_doc_name`
        FROM `t_file_crawl_strategy` JOIN `t_user`
        ON `t_file_crawl_strategy`.`f_user_id` = `t_user`.`f_user_id`
        LEFT JOIN `{get_db_name('anyshare')}`.`t_acs_doc`
        ON `t_file_crawl_strategy`.`f_doc_id` = `t_acs_doc`.`f_doc_id`
        AND `t_acs_doc`.`f_status` = %s
        ORDER BY `f_strategy_id` DESC
        {limit_statement}
        """
        results = self.r_db.all(sql, 1)

        fileCrawlConfigList = []
        for result in results:
            fileCrawlConfig = ncTFileCrawlConfig()
            fileCrawlConfig.strategyId = result['f_strategy_id']
            fileCrawlConfig.userId = result['f_user_id']
            fileCrawlConfig.loginName = result['f_login_name']
            fileCrawlConfig.displayName = result['f_display_name']
            fileCrawlConfig.docId = result['f_doc_id']
            fileCrawlConfig.docName = result['f_doc_name'] if result['f_doc_name'] else ""
            fileCrawlConfig.fileCrawlType = result['f_file_crawl_type']
            fileCrawlConfigList.append(fileCrawlConfig)
        return fileCrawlConfigList

    def get_search_file_crawl_config_count(self, searchKey):
        """
        根据关键字分页搜索文档抓取配置数
        """
        search_sql = f"""
        SELECT COUNT(*) AS cnt
        FROM `t_file_crawl_strategy` JOIN `t_user`
        ON `t_file_crawl_strategy`.`f_user_id` = `t_user`.`f_user_id`
        LEFT JOIN `{get_db_name('anyshare')}`.`t_acs_doc`
        ON `t_file_crawl_strategy`.`f_doc_id` = `t_acs_doc`.`f_doc_id`
        AND `t_acs_doc`.`f_status` = 1
        WHERE `t_user`.`f_login_name` LIKE %s
        OR `t_user`.`f_display_name` LIKE %s
        """
        esckey = "%%%s%%" % escape_key(searchKey)
        result = self.r_db.one(search_sql, esckey, esckey)
        return result['cnt']

    def search_file_crawl_config(self, searchKey, start, limit):
        """
        根据关键字搜索文档抓取配置
        """
        limit_statement = check_start_limit(start, limit)

        search_sql = f"""
        SELECT `f_file_crawl_type`, `f_strategy_id`,
        `t_file_crawl_strategy`.`f_user_id`,
        `t_file_crawl_strategy`.`f_doc_id`,
        `t_user`.`f_login_name` AS `f_login_name`,
        `t_user`.`f_display_name` AS `f_display_name`,
        `t_acs_doc`.`f_name` AS `f_doc_name`
        FROM `t_file_crawl_strategy` JOIN `t_user`
        ON `t_file_crawl_strategy`.`f_user_id` = `t_user`.`f_user_id`
        LEFT JOIN `{get_db_name('anyshare')}`.`t_acs_doc`
        ON `t_file_crawl_strategy`.`f_doc_id` = `t_acs_doc`.`f_doc_id`
        AND `t_acs_doc`.`f_status` = 1
        WHERE `t_user`.`f_login_name` LIKE %s
        OR `t_user`.`f_display_name` LIKE %s
        ORDER BY `f_strategy_id` DESC
        {limit_statement}
        """

        esckey = "%%%s%%" % escape_key(searchKey)
        results = self.r_db.all(search_sql, esckey, esckey)

        fileCrawlConfigList = []
        for result in results:
            fileCrawlConfig = ncTFileCrawlConfig()
            fileCrawlConfig.strategyId = result['f_strategy_id']
            fileCrawlConfig.userId = result['f_user_id']
            fileCrawlConfig.loginName = result['f_login_name']
            fileCrawlConfig.displayName = result['f_display_name']
            fileCrawlConfig.docId = result['f_doc_id']
            fileCrawlConfig.docName = result['f_doc_name'] if result['f_doc_name'] else ""
            fileCrawlConfig.fileCrawlType = result['f_file_crawl_type']
            fileCrawlConfigList.append(fileCrawlConfig)
        return fileCrawlConfigList

    def get_file_crawl_config_by_userid(self, userId):
        """
        根据用户ID获取文件抓取配置
        """
        sql = f"""
        SELECT `f_strategy_id`, `t_file_crawl_strategy`.`f_doc_id`,
               `f_file_crawl_type`,`t_acs_doc`.`f_name`AS `f_doc_name`
        FROM `t_file_crawl_strategy` LEFT JOIN `{get_db_name('anyshare')}`.`t_acs_doc`
        ON `t_file_crawl_strategy`.`f_doc_id` = `t_acs_doc`.`f_doc_id`
        AND `t_acs_doc`.`f_status` = 1
        WHERE `f_user_id` = %s
        """
        result = self.r_db.one(sql, userId)
        fileCrawlConfig = ncTFileCrawlConfig()
        fileCrawlConfig.strategyId = -1
        if result:
            fileCrawlConfig.strategyId = result['f_strategy_id']
            fileCrawlConfig.docId = result['f_doc_id']
            fileCrawlConfig.fileCrawlType = result['f_file_crawl_type']
            fileCrawlConfig.docName = result['f_doc_name'] if result['f_doc_name'] else ""
        return fileCrawlConfig

    def set_file_crawl_show_status(self, status):
        """
        设置控制台是否显示抓取策略开关
        """
        self.config_manage.set_config('file_crawl_show_status', int(status))

    def get_file_crawl_show_status(self):
        """
        获取控制台是否显示抓取策略开关
        """
        return bool(int(self.config_manage.get_config("file_crawl_show_status")))
