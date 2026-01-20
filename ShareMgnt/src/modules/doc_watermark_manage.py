#!/usr/bin/python3
# -*- coding: utf-8 -*-

import time
from src.common.db.connector import DBConnector
from src.common.db.db_manager import get_db_name
from src.common.lib import raise_exception
from src.common.lib import (raise_exception, escape_key, check_start_limit)
from src.common.business_date import BusinessDate
from ShareMgnt.ttypes import (ncTShareMgntError, ncTWatermarkDocInfo, ncTDocType)
import json


class DocWatermarkManage(DBConnector):
    """
    """
    def __init__(self):
        """
        """

    def get_doc_watermark_config(self):
        """
        获取文件水印策略配置
        """
        select_sql = """
        SELECT f_config FROM t_watermark_config limit 1
        """
        result = self.r_db.one(select_sql)

        if result is not None and result["f_config"]:
            if not isinstance(result["f_config"], bytes):
                """
                适配GoldenDB
                """
                result["f_config"] = result["f_config"].encode("utf8")
            return bytes.decode(result["f_config"])
        else:
            raise_exception(exp_msg=_("IDS_INVALID_DOC_WATERMARK_CONFIG"),
                            exp_num=ncTShareMgntError.NCT_INVALID_DOC_WATERMARK_CONFIG)

    def set_doc_watermark_config(self, config):
        """
        设置文件水印策略配置
        """
        config = config.strip()

        if not config:
            raise_exception(exp_msg=_("IDS_INVALID_DOC_WATERMARK_CONFIG"),
                            exp_num=ncTShareMgntError.NCT_INVALID_DOC_WATERMARK_CONFIG)

        update_sql = """
        UPDATE `t_watermark_config`
        SET `f_config` = %s
        """
        self.w_db.query(update_sql, config)

    def add_watermark_doc(self, addId, watermarkType):
        """
        添加开启水印的文档库, 水印类型: 0为无水印(即不对所有文档库开启水印) ，1为预览，2为下载，3为预览与下载
        """
        self.__check_watermark_type(watermarkType)

        if not self.__check_obj_is_exist(addId):
            raise_exception(exp_msg=_("IDS_OBJ_NOT_EXIST"),
                            exp_num=ncTShareMgntError.NCT_OBJ_NOT_EXIST)

        if not self.__check_obj_is_already_added(addId):
            self.w_db.insert("t_watermark_doc", {
                "f_obj_id": addId,
                "f_watermark_type": watermarkType,
                "f_time": int(BusinessDate.time() * 1000000)
                })
        else:
            raise_exception(exp_msg=_("IDS_OBJ_EXIST"),
                            exp_num=ncTShareMgntError.NCT_OBJ_EXIST)

    def update_watermark_doc(self, objId, watermarkType):
        """
        更新开启水印的文档库, 水印类型: 0为无水印(即不对所有文档库开启水印) ，1为预览，2为下载，3为预览与下载
        """
        self.__check_watermark_type(watermarkType)

        if self.__set_doc_lib_watermark_type(objId, watermarkType) == True:
            return

        if not self.__check_obj_is_exist(objId):
            raise_exception(exp_msg=_("IDS_OBJ_NOT_EXIST"),
                            exp_num=ncTShareMgntError.NCT_OBJ_NOT_EXIST)
        if not self.__check_obj_is_already_added(objId):
            raise_exception(exp_msg=_("IDS_OBJ_NOT_EXIST"),
                            exp_num=ncTShareMgntError.NCT_OBJ_NOT_EXIST)
        update_sql = """
        UPDATE `t_watermark_doc`
        SET `f_watermark_type` = %s
        WHERE `f_obj_id` = %s
        """
        self.w_db.query(update_sql, watermarkType, objId)

    def get_watermark_docs(self):
        """
        获取所有开启水印的文档库
        """
        query_sql = f"""
        SELECT `t_watermark_doc`.`f_obj_id` as `f_obj_id`,
        `t_watermark_doc`.`f_watermark_type` as `f_watermark_type`,
        `acs_doc`.`f_name` as `obj_name`,
        `acs_doc`.`f_doc_type` as `doc_type`
        FROM `t_watermark_doc`
        INNER JOIN `{get_db_name('anyshare')}`.`t_acs_doc` as `acs_doc`
        ON `t_watermark_doc`.`f_obj_id` = `acs_doc`.`f_doc_id`
            AND `acs_doc`.`f_status` = %s
        """
        results = self.r_db.all(query_sql, 1)

        retInfos = []
        objIds = []
        for res in results:
            info = ncTWatermarkDocInfo()
            info.objId = res['f_obj_id']
            info.objName = res['obj_name']
            info.objType = res['doc_type']
            info.watermarkType = res['f_watermark_type']
            retInfos.append(info)
            objIds.append(self.w_db.escape(info.objId))

        if objIds:
            delete_sql = """
            DELETE FROM t_watermark_doc
            WHERE f_obj_id NOT IN {0}
            """.format("('" + "','".join(objIds) + "')")
            self.w_db.query(delete_sql)

        select_sql = """
        SELECT f_for_user_doc, f_for_custom_doc, f_for_archive_doc FROM t_watermark_config limit 1
        """
        result = self.r_db.one(select_sql)
        if result is not None:
            info1 = ncTWatermarkDocInfo()
            info1.objId = "1"
            info1.objType = 1
            info1.objName = "userdoc"
            info1.watermarkType = result['f_for_user_doc']
            retInfos.append(info1)
            info3 = ncTWatermarkDocInfo()
            info3.objId = "3"
            info3.objType = 3
            info3.objName = "customdoc"
            info3.watermarkType = result['f_for_custom_doc']
            retInfos.append(info3)
            # info5 = ncTWatermarkDocInfo()
            # info5.objId = "5"
            # info5.objType = 5
            # info5.objName = "archivedoc"
            # info5.watermarkType = result['f_for_archive_doc']
            # retInfos.append(info5)
        return retInfos

    def delete_watermark_doc(self, deleteId):
        """
        删除开启水印的文档库
        """
        delete_sql = """
        DELETE FROM `t_watermark_doc`
        WHERE `f_obj_id` = %s
        """
        affect_row = self.w_db.query(delete_sql, deleteId)
        if not affect_row:
            raise_exception(exp_msg=_("IDS_OBJ_NOT_EXIST"),
                            exp_num=ncTShareMgntError.NCT_OBJ_NOT_EXIST)

    def get_watermark_doc_cnt(self):
        """
        获取开启水印的文档库总数
        """
        query_sql = f"""
        SELECT COUNT(`t_watermark_doc`.`f_obj_id`) as cnt
        FROM `t_watermark_doc`
        INNER JOIN `{get_db_name('anyshare')}`.`t_acs_doc` as `acs_doc`
        ON `t_watermark_doc`.`f_obj_id` = `acs_doc`.`f_doc_id`
            AND `acs_doc`.`f_status` = %s
        """
        result = self.r_db.one(query_sql, 1)

        # 默认加上两种不同类型文档库的水印配置数（个人文档库，自定义文档库）
        return result['cnt'] + 2

    def get_watermark_doc_by_page(self, start, limit):
        """
        分页获取开启水印的文档库信息
        """
        limit_statement = check_start_limit(start, limit)

        query_sql = f"""
        SELECT `t_watermark_doc`.`f_obj_id` as `f_obj_id`,
        `t_watermark_doc`.`f_watermark_type` as `f_watermark_type`,
        `acs_doc`.`f_name` as `obj_name`,
        `acs_doc`.`f_doc_type` as `doc_type`
        FROM `t_watermark_doc`
        INNER JOIN `{get_db_name('anyshare')}`.`t_acs_doc` as `acs_doc`
        ON `t_watermark_doc`.`f_obj_id` = `acs_doc`.`f_doc_id`
            AND `acs_doc`.`f_status` = %s
        ORDER BY `t_watermark_doc`.`f_time` desc {limit_statement}
        """
        results = self.r_db.all(query_sql, 1)

        retInfos = []
        if start == 0:
            select_sql = """
            SELECT f_for_user_doc, f_for_custom_doc, f_for_archive_doc FROM t_watermark_config limit 1
            """
            result = self.r_db.one(select_sql)
            if result is not None:
                info1 = ncTWatermarkDocInfo()
                info1.objId = "1"
                info1.objType = 1
                info1.objName = "userdoc"
                info1.watermarkType = result['f_for_user_doc']
                retInfos.append(info1)
                info3 = ncTWatermarkDocInfo()
                info3.objId = "3"
                info3.objType = 3
                info3.objName = "customdoc"
                info3.watermarkType = result['f_for_custom_doc']
                retInfos.append(info3)
                # info5 = ncTWatermarkDocInfo()
                # info5.objId = "5"
                # info5.objType = 5
                # info5.objName = "archivedoc"
                # info5.watermarkType = result['f_for_archive_doc']
                # retInfos.append(info5)
        for res in results:
            info = ncTWatermarkDocInfo()
            info.objId = res['f_obj_id']
            info.objName = res['obj_name']
            info.objType = res['doc_type']
            info.watermarkType = res['f_watermark_type']
            retInfos.append(info)

        return retInfos

    def search_watermark_doc_cnt(self, searchKey):
        """
        搜索开启水印的文档库总数
        """
        # 按文档库名称搜索
        search_sql = f"""
        SELECT COUNT(`t_watermark_doc`.`f_obj_id`) as cnt
        FROM `t_watermark_doc`
        INNER JOIN `{get_db_name('anyshare')}`.`t_acs_doc` as `acs_doc`
        ON `t_watermark_doc`.`f_obj_id` = `acs_doc`.`f_doc_id`
            AND `acs_doc`.`f_status` = 1 AND `acs_doc`.`f_name` LIKE %s
        """
        esckey = "%%%s%%" % escape_key(searchKey)
        result = self.r_db.one(search_sql, esckey)
        count_doc = result['cnt']

        # 按个人文档显示名搜索
        search_user_doc_sql = f"""
        SELECT COUNT(`t_watermark_doc`.`f_obj_id`) as cnt
        FROM `t_watermark_doc`
        INNER JOIN `{get_db_name('anyshare')}`.`t_acs_doc` as `acs_doc`
        ON `t_watermark_doc`.`f_obj_id` = `acs_doc`.`f_doc_id` AND `acs_doc`.`f_doc_type` = 1
        INNER JOIN `t_user` as `user`
        ON `acs_doc`.`f_creater_id` = `user`.`f_user_id` AND `user`.`f_display_name` LIKE %s
            AND `acs_doc`.`f_status` = 1
        """
        result = self.r_db.one(search_user_doc_sql, esckey)
        count_user_doc = result['cnt']

        return count_doc + count_user_doc

    def search_watermark_doc_by_page(self, searchKey, start, limit):
        """
        搜索开启水印的文档库信息
        """
        check_start_limit(start, limit)

        # 按文档库名称搜索
        search_sql = f"""
        SELECT `t_watermark_doc`.`f_obj_id` as `f_obj_id`,
        `t_watermark_doc`.`f_watermark_type` as `f_watermark_type`,
        `acs_doc`.`f_name` as `obj_name`,
        `acs_doc`.`f_doc_type` as `doc_type`,
        `t_watermark_doc`.`f_time` as `time`
        FROM `t_watermark_doc`
        INNER JOIN `{get_db_name('anyshare')}`.`t_acs_doc` as `acs_doc`
        ON `t_watermark_doc`.`f_obj_id` = `acs_doc`.`f_doc_id`
            AND `acs_doc`.`f_status` = 1 AND `acs_doc`.`f_name` LIKE %s
        """
        esckey = "%%%s%%" % escape_key(searchKey)
        results = self.r_db.all(search_sql, esckey)
        retInfos = []
        for res in results:
            info = ncTWatermarkDocInfo()
            info.objId = res['f_obj_id']
            info.objName = res['obj_name']
            info.objType = res['doc_type']
            info.time = res['time']
            info.watermarkType = res['f_watermark_type']
            retInfos.append(info)

        # 按个人文档显示名搜索
        search_user_doc_sql = f"""
        SELECT `t_watermark_doc`.`f_obj_id` as `f_obj_id`,
        `t_watermark_doc`.`f_watermark_type` as `f_watermark_type`,
        `acs_doc`.`f_doc_type` as `doc_type`,
        `user`.`f_display_name` as `display_name`,
        `t_watermark_doc`.`f_time` as `time`
        FROM `t_watermark_doc`
        INNER JOIN `{get_db_name('anyshare')}`.`t_acs_doc` as `acs_doc`
        ON `t_watermark_doc`.`f_obj_id` = `acs_doc`.`f_doc_id` AND `acs_doc`.`f_doc_type` = 1
        INNER JOIN `t_user` as `user`
        ON `acs_doc`.`f_creater_id` = `user`.`f_user_id` AND `user`.`f_display_name` LIKE %s
            AND `acs_doc`.`f_status` = 1
        """
        results = self.r_db.all(search_user_doc_sql, esckey)
        for res in results:
            info = ncTWatermarkDocInfo()
            info.objId = res['f_obj_id']
            info.objName = res['display_name']
            info.objType = res['doc_type']
            info.time = res['time']
            info.watermarkType = res['f_watermark_type']
            retInfos.append(info)

        # 时间戳大的放前面
        retInfos.sort(key=lambda x: x.time, reverse=True)
        if limit == -1:
            return retInfos
        else:
            return retInfos[0:limit]

    def set_watermark_type_for_libs(self, docType, watermarkType):
        """
        设置指定类型文档库水印类型，0为无水印(即不对所有文档库开启水印) ，1为预览，2为下载，3为预览与下载
        """
        docTypeStr = self.__check_doc_type(docType)
        self.__check_watermark_type(watermarkType)

        update_sql = """
        UPDATE `t_watermark_config`
        SET `%s` = %s
        """ % (docTypeStr, watermarkType)
        self.w_db.query(update_sql)

    def __check_obj_is_exist(self, obj_id):
        """
        """
        is_exist = True
        try:
            sql = f"""
            SELECT COUNT(*) AS cnt FROM `{get_db_name('anyshare')}`.`t_acs_doc`
            WHERE `f_doc_id` = %s
            """
            count = self.r_db.one(sql, obj_id)['cnt']
            is_exist = True if count else False
        except:
            is_exist = False
        return is_exist

    def __check_obj_is_already_added(self, objId):
        """
        """
        sql = """
        SELECT COUNT(*) AS cnt FROM `t_watermark_doc`
        WHERE `f_obj_id` = %s
        """
        return True if self.r_db.one(sql, objId)["cnt"] == 1 else False

    def __check_doc_type(self, docType):
        if docType < ncTDocType.NCT_USER_DOC or docType > ncTDocType.NCT_ARCHIVE_DOC or docType == ncTDocType.NCT_SHARE_DOC:
            raise_exception(exp_msg=_("IDS_INVALID_DOC_TYPE"),
                            exp_num=ncTShareMgntError.NCT_INVALID_DOC_TYPE)
        if docType == ncTDocType.NCT_USER_DOC:
            return "f_for_user_doc"
        if docType == ncTDocType.NCT_CUSTOM_DOC:
            return "f_for_custom_doc"
        if docType == ncTDocType.NCT_ARCHIVE_DOC:
            return "f_for_archive_doc"

    def __check_watermark_type(self, watermarkType):
        if watermarkType < 0 or watermarkType > 3:
            raise_exception(exp_msg=_("IDS_INVALID_WATERMARK_TYPE"),
                            exp_num=ncTShareMgntError.NCT_INVALID_WATERMARK_TYPE)

    def __set_doc_lib_watermark_type(self, objId, watermarkType):
        if "1" == objId:
            self.set_watermark_type_for_libs (ncTDocType.NCT_USER_DOC, watermarkType)
            return True
        elif "3" == objId:
            self.set_watermark_type_for_libs (ncTDocType.NCT_CUSTOM_DOC, watermarkType)
            return True
        elif "5" == objId:
            self.set_watermark_type_for_libs (ncTDocType.NCT_ARCHIVE_DOC, watermarkType)
            return True
        return False
