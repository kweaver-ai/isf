#!/usr/bin/python3
# -*- coding:utf-8 -*-
"""
内外链模板配置管理
"""
import time
import uuid
import json
from src.common.db.connector import DBConnector
from src.common.lib import (raise_exception,
                            escape_key)
from src.common.business_date import BusinessDate
from src.modules.user_manage import UserManage
from src.modules.config_manage import ConfigManage
from src.modules.department_manage import DepartmentManage
from ShareMgnt.ttypes import (ncTShareMgntError,
                              ncTAuditObjectType,
                              ncTLinkTemplateInfo,
                              ncTTemplateType,
                              ncTLinkShareInfo)
from ShareMgnt.constants import NCT_ALL_USER_GROUP


NO_DEL_PERM = 0x1F
INTERNAL_PERM_MAX = 0x3F # 同ncIACSPermManager.idl文件中的ncAtomPermValue.ACS_CP_MAX
EXTERNAL_PERM_MAX = 0x1F # 同ncIACSPermManager.idl文件中的ncAtomPermValue.ACS_CP_ANONYMOUS_MAX


class LinkTemplateManage(DBConnector):
    """
    内外链模板配置管理类
    """
    def __init__(self):
        self.depart_manage = DepartmentManage()
        self.user_manage = UserManage()
        self.config_manage = ConfigManage()

    def __check_template_type(self, templateType):
        """
        检查模板类型
        """
        if templateType != ncTTemplateType.INTERNAL_LINK and templateType != ncTTemplateType.EXTERNAL_LINK:
            raise_exception(exp_msg=_("IDS_INVALID_LINK_TEMPLATE_TYPE"),
                            exp_num=ncTShareMgntError.NCT_INVALID_LINK_TEMPLATE_TYPE)

    def __check_sharer_type(self, sharerType):
        """
        检查共享者类型是否合法
        """
        if sharerType != ncTAuditObjectType.NCT_AUDIT_OBJECT_USER and sharerType != ncTAuditObjectType.NCT_AUDIT_OBJECT_DEPT:
            raise_exception(exp_msg=_("IDS_INVALID_SHARER_TYPE"),
                            exp_num=ncTShareMgntError.NCT_INVALID_SHAERER_TYPE)

    def __check_sharer_exists(self, sharerInfo):
        """
        检查共享者是否存在
        """
        if sharerInfo.sharerType == ncTAuditObjectType.NCT_AUDIT_OBJECT_USER:
            if not self.user_manage.check_user_exists(sharerInfo.sharerId, raise_ex=False):
                return False

        elif sharerInfo.sharerType == ncTAuditObjectType.NCT_AUDIT_OBJECT_DEPT:
            if not self.depart_manage.check_depart_exists(sharerInfo.sharerId, include_organ=True, raise_ex=False):
                return False

        return True

    def __check_default_link_template(self, templateId):
        """
        检查是否为默认模板
        """
        internal_dafault_template_id = self.get_link_template_by_shareId(ncTTemplateType.INTERNAL_LINK, NCT_ALL_USER_GROUP).templateId
        external_dafault_template_id = self.get_link_template_by_shareId(ncTTemplateType.EXTERNAL_LINK, NCT_ALL_USER_GROUP).templateId
        if templateId == internal_dafault_template_id or templateId == external_dafault_template_id:
            return True

    def __check_secret_internal_link_template(self, jsonConfig):
        """
        涉密模式下检查内链模板参数是否合法
        """
        # 不能设置所有者，且默认权限不能为永久有效
        if self.config_manage.get_secret_mode_status():
            if jsonConfig["limitExpireDays"] is False and jsonConfig["allowExpireDays"] == -1:
                raise_exception(exp_msg=_("IDS_EXCEED_MAX_INTERNAL_LINK_EXPIRE_DAY"),
                                exp_num=ncTShareMgntError.NCT_EXCEED_MAX_INTERNAL_LINK_EXPIRE_DAY)
            if jsonConfig["allowOwner"] is True or jsonConfig["defaultOwner"] is True:
                raise_exception(exp_msg=_("IDS_CANNOT_SET_OWNER"),
                                exp_num=ncTShareMgntError.NCT_CANNOT_SET_OWNER)

    def get_link_template_by_shareId(self, templateType, shareId):
        """
        根据共享者ID获取模板
        """
        sql = """
        SELECT `f_template_id`, `f_config`, `f_create_time`
        FROM `t_link_template`
        WHERE `f_sharer_id` = %s and `f_template_type` = %s
        """
        result = self.r_db.one(sql, shareId, templateType)

        if result:
            template = ncTLinkTemplateInfo()
            template.templateId = result["f_template_id"]
            template.config = result["f_config"]
            template.createTime = result["f_create_time"]
            return template

    def get_one_link_template_by_templateId(self, templateId, raise_ex=False):
        """
        根据模板ID获取模板公共参数
        """
        sql = """
        SELECT `f_template_type`, `f_create_time`, `f_config`
        FROM `t_link_template`
        WHERE `f_template_id` = %s
        """

        result = self.r_db.one(sql, templateId)
        if not result:
            if raise_ex:
                raise_exception(exp_msg=_("IDS_TEMPLATE_NOT_EXIST"),
                                exp_num=ncTShareMgntError.NCT_TEMPLATE_NOT_EXIST)
        else:
            return result

    def convert_link_template_info(self, db_link_template_infos):
        """
        转换数据库中获取的模板信息
        """
        link_template_infos = []

        for db_link_template_info in db_link_template_infos:
            link_template_info = ncTLinkTemplateInfo()
            link_template_info.templateType = db_link_template_info["f_template_type"]
            link_template_info.templateId = db_link_template_info["f_template_id"]
            link_template_info.config = db_link_template_info["f_config"]
            link_template_info.sharerInfos = self.get_sharer_infos_by_templateId(link_template_info.templateId)
            link_template_infos.append(link_template_info)

        return link_template_infos

    def get_sharer_infos_by_templateId(self, template_id):
        """
        根据模板ID获取共享者信息
        """
        sharerInfos = []

        sql = """
            SELECT DISTINCT `f_sharer_id`, `f_sharer_type`
            FROM `t_link_template`
            WHERE `f_template_id` = %s
        """

        results = self.r_db.all(sql, template_id)

        for res in results:
            sharerInfo = ncTLinkShareInfo()
            sharerInfo.sharerId = res["f_sharer_id"]
            sharerInfo.sharerType = int(res["f_sharer_type"])

            if sharerInfo.sharerId == NCT_ALL_USER_GROUP:
                sharerInfo.sharerName = _("all user")
            else:
                # 共享者为用户
                if sharerInfo.sharerType == ncTAuditObjectType.NCT_AUDIT_OBJECT_USER:
                    # 过滤无效用户
                    try:
                        user_info = self.user_manage.get_user_by_id(sharerInfo.sharerId)
                        sharerInfo.sharerName = user_info.user.displayName
                    except Exception:
                        pass

                # 共享者为部门或组织
                elif sharerInfo.sharerType == ncTAuditObjectType.NCT_AUDIT_OBJECT_DEPT:
                    # 过滤无效部门
                    try:
                        depart_info = self.depart_manage.get_department_info(sharerInfo.sharerId, b_include_org=True)
                        sharerInfo.sharerName = depart_info.departmentName
                    except Exception:
                        pass

            if sharerInfo.sharerName:
                sharerInfos.append(sharerInfo)

        return sharerInfos

    def get_calculated_link_template_by_userId(self, templateType, userId):
        """
        获取计算过的内外链模板信息，规则：个人>子部门>父部们，相同级别部门按最新创建的部门为准
        """
        # 检查用户是否存在
        self.user_manage.check_user_exists(userId)

        # 快速判断，如果只有一条默认策略则直接返回
        sql = """
        SELECT COUNT(*) AS cnt FROM `t_link_template`
        WHERE  `f_template_type` = %s
        """
        cnt = self.r_db.one(sql, templateType)['cnt']
        if 1 == cnt:
            # 获取默认策略
            defaultConfig = self.get_link_template_by_shareId(templateType, NCT_ALL_USER_GROUP)
            return defaultConfig

        # 获取用户模板，如果存在直接返回
        template = self.get_link_template_by_shareId(templateType, userId)

        # 存在用户级别的模板则直接返回
        if template:
            return template

        # 该用户相关的部门树
        depart_tree = self.depart_manage.get_depart_tree_of_user(userId)

        # 遍历部门树，获取部门树所有部门的模板
        all_template_dict = {}
        for depart_id in depart_tree:
            template = self.get_link_template_by_shareId(templateType, depart_id)
            if not template:
                continue

            # 所有部门的模板先置为有效
            template.valid = True
            all_template_dict[depart_id] = template

        # 每个叶子节点向上遍历
        for depart_id in depart_tree:
            if depart_tree[depart_id].subDepartIds:
                continue

            # 向上遍历
            all_path_depart_ids = [depart_id]
            parent_id = depart_tree[depart_id].parentDepartId
            while parent_id:
                all_path_depart_ids.append(parent_id)
                parent_id = depart_tree[parent_id].parentDepartId

            # 遍历路径，如果路径上有设置的模版，则剩余的路径上存在的模板必须置为失效
            for i in range(0, len(all_path_depart_ids)):
                # 不存在模板，则继续向上遍历
                if all_path_depart_ids[i] not in all_template_dict:
                    continue

                # 存在模板，则需要将剩余路径上存在的模板设置为失效
                for j in range(i + 1, len(all_path_depart_ids)):
                    remain_depart_id = all_path_depart_ids[j]
                    if remain_depart_id in all_template_dict:
                        all_template_dict[remain_depart_id].valid = False

        # 针对所有部门模板进行计算，取有效的，时间最近的
        result = None
        for template_id in all_template_dict:
            if not all_template_dict[template_id].valid:
                continue

            # 第一次设置result
            if not result:
                result = all_template_dict[template_id]
                continue

            # 比result新，则更新为该模板
            if all_template_dict[template_id].createTime > result.createTime:
                result = all_template_dict[template_id]

        # 删除掉动态添加的valid属性
        if result:
            del result.valid

        # 如果没有配置模板则返回默认模板
        if not result:
            result = self.get_link_template_by_shareId(templateType, NCT_ALL_USER_GROUP)

        return result

    def check_link_perm(self, setPerm, allowPerm):
        """
        检查设定的权限是否在可设定的范围内
        """
        # 设定的访问权限不能超过可设定的权限
        if (setPerm ^ allowPerm) & setPerm:
            raise_exception(exp_msg=_("IDS_EXCEED_MAX_LINK_PERM"),
                            exp_num=ncTShareMgntError.NCT_EXCEED_MAX_LINK_PERM)

    def check_external_link_perm(self, linkInfo):
        """
        检查外链配置是否符合模板
        """
        # 根据用户ID获取有效的外链共享模板
        template = self.get_calculated_link_template_by_userId(ncTTemplateType.EXTERNAL_LINK, linkInfo.userId)

        # 获取外链模板配置
        if template:
            config = json.loads(template.config)

            # 限制外链的有效期，有效期不能超过模板设定的最大天数
            if config["limitExpireDays"] and (linkInfo.allowExpireDays > config["allowExpireDays"] or linkInfo.allowExpireDays == -1):
                raise_exception(exp_msg=_("IDS_EXCEED_MAX_OUT_LINK_DATE"),
                                exp_num=ncTShareMgntError.NCT_EXCEED_MAX_EXTERNAL_LINK_EXPIRE_DAY)

            # 限制外链的访问次数，设置的最大访问次数不能超过模板设定的最大次数
            if config["limitAccessTimes"] and (linkInfo.accessLimit > config["allowAccessTimes"] or linkInfo.accessLimit == -1):
                raise_exception(exp_msg=_("IDS_EXCEED_MAX_OUT_LINK_ACCESS_TIME"),
                                exp_num=ncTShareMgntError.NCT_EXCEED_MAX_EXTERNAL_LINK_EXPIRE_TIME)

            # 强制使用密码，未设置访问密码
            if config["accessPassword"] and not linkInfo.password:
                raise_exception(exp_msg=_("IDS_NEED_SET_ACCESS_PASSWORD"),
                                exp_num=ncTShareMgntError.NCT_NEED_SET_ACCESS_PASSWORD)

            # 设定的访问权限不能超过可设定的权限
            self.check_link_perm(linkInfo.permValue, config["allowPerm"])

    def __check_template_config(self, templateInfo):
        """
        检查模板参数
        """
        config = json.loads(templateInfo.config)
        if ncTTemplateType.INTERNAL_LINK == templateInfo.templateType:
            # 允许的权限值是否合法
            self.check_link_perm(config["allowPerm"], INTERNAL_PERM_MAX)

            # 涉密模式下检查内链模板参数是否合法
            self.__check_secret_internal_link_template(config)
        elif ncTTemplateType.EXTERNAL_LINK == templateInfo.templateType:
            # 允许的权限值是否合法
            self.check_link_perm(config["allowPerm"], EXTERNAL_PERM_MAX)

    def add_link_template(self, templateInfo):
        """
        添加模板策略
        """
        # 检查模板类型是否合法
        self.__check_template_type(templateInfo.templateType)

        # 检查模板参数是否合法
        self.__check_template_config(templateInfo)

        # 检查共享者冲突
        conflict_sharers = []
        filter_sharers_names = []
        tmp_share_infos = templateInfo.sharerInfos[:]

        for sharerInfo in tmp_share_infos:

            # 检查共享者类型是否合法
            self.__check_sharer_type(sharerInfo.sharerType)

            # 直接过滤不存在的共享者，处理添加模板时用户、部门恰好被删除的情况
            if not self.__check_sharer_exists(sharerInfo):

                # 记录过滤的共享者
                filter_sharers_names.append(sharerInfo.sharerName)

                templateInfo.sharerInfos.remove(sharerInfo)
                continue

            # 检查共享者是否已经配置一条模板
            if self.get_link_template_by_shareId(templateInfo.templateType, sharerInfo.sharerId):
                conflict_sharers.append(sharerInfo.sharerName)

        # 存在冲突共享者时，直接返回冲突共享者名称列表
        if conflict_sharers:
            return conflict_sharers

        # 共享者为空
        if not templateInfo.sharerInfos:
            filter_sharers_names_str = ",".join(filter_sharers_names)
            raise_exception(exp_msg=_("IDS_SHARER_NOT_EXISTS") % filter_sharers_names_str,
                            exp_num=ncTShareMgntError.NCT_SHARER_IS_EMPTY)

        template_id = str(uuid.uuid1())
        insert_sql = """
        INSERT INTO `t_link_template`
        (`f_template_id`, `f_template_type`, `f_sharer_id`, `f_sharer_type`,`f_create_time`, `f_config`)
        VALUES(%s, %s, %s, %s, %s, %s)
        """

        for sharerInfo in templateInfo.sharerInfos:
            self.w_db.query(insert_sql, template_id, templateInfo.templateType,
                            sharerInfo.sharerId, sharerInfo.sharerType,
                            int(BusinessDate.time() * 1000000), templateInfo.config)

        # 添加成功返回空列表
        return []

    def get_sharer_ids_by_template_id(self, template_id):
        """
        根据模板ID获取共享者ID
        """
        sql = """
            SELECT DISTINCT `f_sharer_id`
            FROM `t_link_template`
            WHERE `f_template_id` = %s
        """

        results = self.r_db.all(sql, template_id)

        sharer_ids = []
        for res in results:
            sharer_id = res["f_sharer_id"]
            sharer_ids.append(sharer_id)

        return sharer_ids

    def edit_link_template(self, templateInfo):
        """
        设置内外链模板信息
        """
        # 检查模板类型是否合法
        self.__check_template_type(templateInfo.templateType)

        # 检查模板参数是否合法
        self.__check_template_config(templateInfo)

        # 获取模板创建时间
        result = self.get_one_link_template_by_templateId(templateInfo.templateId, raise_ex=True)
        createTime = int(result['f_create_time'])

        # 不允许编辑默认模板的共享者
        if self.__check_default_link_template(templateInfo.templateId):
            if len(templateInfo.sharerInfos) != 1 or templateInfo.sharerInfos[0].sharerId != NCT_ALL_USER_GROUP:
                raise_exception(exp_msg=_("IDS_CANNOT_EDIT_DEFAULT_LINK_TEMPLATE_SHARER"),
                                exp_num=ncTShareMgntError.NCT_CANNOT_EDIT_DEFAULT_LINK_TEMPLATE_SHARER)

        # 检查共享者冲突
        conflict_sharers = []
        src_sharer_ids = []
        filter_sharers_names = []
        tmp_share_infos = templateInfo.sharerInfos[:]

        for sharerInfo in tmp_share_infos:
            src_sharer_ids.append(sharerInfo.sharerId)
            # 检查共享者类型是否合法
            self.__check_sharer_type(sharerInfo.sharerType)

            # 直接过滤不存在的共享者，处理添加模板时用户、部门恰好被删除的情况
            if not self.__check_sharer_exists(sharerInfo):

                # 记录过滤的共享者
                filter_sharers_names.append(sharerInfo.sharerName)

                templateInfo.sharerInfos.remove(sharerInfo)
                continue

            template = self.get_link_template_by_shareId(templateInfo.templateType, sharerInfo.sharerId)

            if template and template.templateId != templateInfo.templateId:
                conflict_sharers.append(sharerInfo.sharerName)

        # 存在已经配置其他模板的共享者时直接返回
        if conflict_sharers:
            return conflict_sharers

        # 共享者为空
        if not templateInfo.sharerInfos:
            filter_sharers_names_str = ",".join(filter_sharers_names)
            raise_exception(exp_msg=_("IDS_SHARER_NOT_EXISTS") % filter_sharers_names_str,
                            exp_num=ncTShareMgntError.NCT_SHARER_IS_EMPTY)

        # 删除共享者
        delete_sharer_ids = []
        dest_sharer_ids = self.get_sharer_ids_by_template_id(templateInfo.templateId)
        for sharer_id in dest_sharer_ids:
            if sharer_id not in src_sharer_ids:
                delete_sharer_ids.append(self.w_db.escape(sharer_id))

        if delete_sharer_ids:
            delete_sql = """
            DELETE FROM t_link_template
            WHERE f_template_id = %s AND f_sharer_id IN {0}
            """.format("('" + "','".join(delete_sharer_ids) + "')")
            self.w_db.query(delete_sql, templateInfo.templateId)

        for sharerInfo in templateInfo.sharerInfos:
            if sharerInfo.sharerId not in dest_sharer_ids:
                # 共享者不存在
                insert_sql = """
                INSERT INTO `t_link_template`
                (`f_template_id`, `f_template_type`, `f_sharer_id`, `f_sharer_type`, `f_create_time`, `f_config`)
                VALUES (%s, %s, %s, %s, %s, %s)
                """
                self.w_db.query(insert_sql,
                                templateInfo.templateId,
                                templateInfo.templateType,
                                sharerInfo.sharerId,
                                sharerInfo.sharerType,
                                createTime,
                                templateInfo.config)
            else:
                # 共享者存在，更新配置
                update_sql = """
                UPDATE `t_link_template`
                SET `f_config` = %s
                WHERE `f_template_id` = %s AND `f_sharer_id` = %s
                """
                self.w_db.query(update_sql, templateInfo.config, templateInfo.templateId, sharerInfo.sharerId)

        return []

    def delete_link_template_by_templateId(self, template_id):
        """
        删除模板
        """
        if self.__check_default_link_template(template_id):
            raise_exception(exp_msg=_("IDS_CANNOT_DELETE_DEFAULT_LINK_TEMPLATE"),
                            exp_num=ncTShareMgntError.NCT_CANNOT_DELETE_DEFAULT_LINK_TEMPLATE)
        else:
            delete_sql = """
            DELETE FROM `t_link_template`
            WHERE `f_template_id` = %s
            """
            affect_row = self.r_db.query(delete_sql, template_id)
            if not affect_row:
                # 模板不存在
                raise_exception(exp_msg=_("IDS_TEMPLATE_NOT_EXIST"),
                                exp_num=ncTShareMgntError.NCT_TEMPLATE_NOT_EXIST)

    def get_link_template(self, templateType):
        """
        获取内外链模板信息
        """
        self.__check_template_type(templateType)

        sql = """
            SELECT `t_link_template`.`f_template_id` AS `f_template_id`,
            MIN(`t_link_template`.`f_config`) AS `f_config`,
            MIN(`t_link_template`.`f_template_type`) AS `f_template_type`,
            MIN(`t_link_template`.`f_sharer_id`),
            MIN(`t_link_template`.`f_create_time`) AS `f_create_time`
            FROM `t_link_template`
            WHERE `t_link_template`.`f_template_type` = %s
            GROUP BY `f_template_id`
            ORDER BY CASE MIN(`t_link_template`.`f_sharer_id`) WHEN %s THEN 0 ELSE 1 END,
            `f_create_time` DESC
        """
        results = self.r_db.all(sql, templateType, NCT_ALL_USER_GROUP)

        return self.convert_link_template_info(results)

    def search_link_template(self, templateType, searchKey):
        """
        搜索模板
        """
        self.__check_template_type(templateType)

        search_sql = """
            SELECT `t_link_template`.`f_template_id` AS `f_template_id`,
            MIN(`t_link_template`.`f_config`) AS `f_config`,
            MIN(`t_link_template`.`f_template_type`) AS `f_template_type`,
            MIN(`t_link_template`.`f_sharer_id`),
            MIN(`t_link_template`.`f_create_time`) AS `f_create_time`
            FROM `t_link_template`
            LEFT JOIN `t_user`
            ON `t_link_template`.`f_sharer_id` = `t_user`.`f_user_id`
            LEFT JOIN `t_department`
            ON `t_link_template`.`f_sharer_id` = `t_department`.`f_department_id`
            WHERE (`t_link_template`.`f_template_type` = %s AND
            (`t_user`.`f_display_name` LIKE %s OR `t_department`.`f_name` LIKE %s))
            GROUP BY `f_template_id`
            ORDER BY CASE MIN(`t_link_template`.`f_sharer_id`) WHEN %s THEN 0 ELSE 1 END,
            `f_create_time` DESC
        """
        esckey = "%%%s%%" % escape_key(searchKey)
        results = self.r_db.all(search_sql, templateType, esckey, esckey, NCT_ALL_USER_GROUP)

        return self.convert_link_template_info(results)

    def remove_del_perm_from_internal_link_template(self):
        """
        去除内链模板中的删除权限
        """
        search_sql = """
            SELECT DISTINCT `f_template_id`, `f_config` FROM `t_link_template`
            WHERE `f_template_type` = %s
        """

        update_sql = """
        UPDATE `t_link_template`
        SET `f_config` = %s
        WHERE `f_template_id` = %s
        """

        results = self.r_db.all(search_sql, '0')
        for res in results:
            old_config = json.loads(res["f_config"])

            # 可设定权限中去掉删除权限
            old_config["allowPerm"] = old_config["allowPerm"] & NO_DEL_PERM

            # 默认权限中去掉删除权限
            old_config["defaultPerm"] = old_config["defaultPerm"] & NO_DEL_PERM

            now_config = json.dumps(old_config)

            self.w_db.query(update_sql, now_config, res["f_template_id"])
