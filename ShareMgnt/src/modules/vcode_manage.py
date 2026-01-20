#!/usr/bin/python3
# -*- coding:utf-8 -*-
"""
登录验证码管理类
"""
import MySQLdb
import uuid
import json
import string
import datetime
import random
import io
import os
import sys
import base64
from src.common.db.connector import DBConnector
from src.common.lib import (raise_exception, escape_key, check_name)
from src.common.business_date import BusinessDate
from src.modules.vcode_image_manage import VcodeImageManage
from src.third_party_auth.third_config_manage import ThirdConfigManage
from EThriftException.ttypes import ncTException
from ShareMgnt.ttypes import (ncTShareMgntError,
                              ncTVcodeConfig,
                              ncTVcodeCreateInfo,
                              ncTVcodeType,
                              ncTPluginType)
from eisoo.tclients import TClient
from src.common.sharemgnt_logger import ShareMgnt_Log


class VcodeManage(DBConnector):
    def __init__(self):
        font_type_path = os.path.realpath(os.path.join(os.path.dirname(sys.argv[0]), "../conf/simsun.ttc"))
        self.vcode_image_manage = VcodeImageManage(font_type=font_type_path, draw_lines=False, draw_points=False)
        self.third_config_manage = ThirdConfigManage()

    def set_vcode_config(self, vcode_config):
        """
        设置登录验证码配置信息
        """
        update_config_sql = """
        UPDATE `t_sharemgnt_config`
        SET `f_value` = %s
        WHERE `f_key` = %s
        """

        if vcode_config.passwdErrCnt is not None:
            if vcode_config.passwdErrCnt < 0 or vcode_config.passwdErrCnt > 99:
                raise_exception(exp_msg=_("IDS_INVALID_PASSWORD_ERR_CNT"),
                                exp_num=ncTShareMgntError.NCT_INVALID_PASSWORD_ERR_CNT)

        if vcode_config.isEnable is not None:
            data = {}
            data['isEnable'] = vcode_config.isEnable
            if vcode_config.isEnable:
                data['passwdErrCnt'] = vcode_config.passwdErrCnt
            else:
                data['passwdErrCnt'] = self.get_vcode_config().passwdErrCnt
            self.w_db.query(update_config_sql, json.dumps(data), 'vcode_login_config')

    def get_vcode_config(self):
        """
        获取用户密码配置信息
        """
        select_config_sql = """
        SELECT `f_value` FROM `t_sharemgnt_config`
        WHERE `f_key` = %s
        """
        vcode_config = self.r_db.one(select_config_sql, "vcode_login_config")   # 返回结果是一个字典
        config_dict = json.loads(vcode_config["f_value"])   # 将需要的数据拿出来
        config = ncTVcodeConfig()
        config.isEnable = config_dict['isEnable']
        config.passwdErrCnt = config_dict['passwdErrCnt']
        return config

    def get_vcode_with_b64(self):
        """
        生成验证码字符串以及验证码图片经过 base64 编码后字符串
        """
        img, vcode = self.vcode_image_manage.create_check_code()

        stm = io.BytesIO()
        img.save(stm, "jpeg")
        data = stm.getvalue()
        b64 = base64.b64encode(data)
        return vcode, bytes.decode(b64)

    def create_vcode_info(self, uuidIn, vcodeType=ncTVcodeType.IMAGE_VCODE):
        """
        生成登录验证码/忘记密码生成验证码
        vcodeType为3时，uuidIn等于用户id
        """
        # 限制验证码发送间隔
        vcode_info = ncTVcodeCreateInfo()
        if vcodeType == ncTVcodeType.DAUL_AUTH_VCODE:
            if self._check_vcode_by_third_auth(uuidIn):
                vcode_info.isDuplicateSended = True
                return vcode_info
            else:
                vcode_info.isDuplicateSended = False
        if uuidIn:
            self.delete_vcode_info(uuidIn)
        # 传入用户id时，uuid为用户id，且删除表里uuid为用户id的旧验证码
        if vcodeType == ncTVcodeType.IMAGE_VCODE:
            # 用户主动点击图片刷新验证码时，需要删除上一条获取到的验证码信息
            vcode, b64 = self.get_vcode_with_b64()
            vcode_info.vcode = b64
        else:
            # 忘记密码生成6位随机数
            vcode_info.vcode = str(random.randint(100000, 999999))
            vcode = vcode_info.vcode
        # 如果ncTVcodeType为DAUL_AUTH_VCODE， uuid等于传入的uuid
        vcode_info.uuid = uuidIn if vcodeType == ncTVcodeType.DAUL_AUTH_VCODE else str(uuid.uuid1())
        self.save_vcode_info(vcode_info.uuid, vcode, vcodeType)
        return vcode_info

    def verify_vcode_info(self, uuidIn, vcodeIn, vcodeType=ncTVcodeType.IMAGE_VCODE, delete_after_check=True):
        """
        校验验证码
        """
        if not uuidIn:
            raise_exception(exp_msg=_("IDS_VCODE_NULL"),
                            exp_num=ncTShareMgntError.NCT_CHECK_VCODE_IS_NULL)

        now = BusinessDate.now()

        delta = datetime.timedelta(minutes=5)   # days; minutes; seconds. 默认为 days

        allow_max_last_useful_time = (now - delta).strftime('%Y-%m-%d %H:%M:%S')
        select_sql = """
        SELECT f_vcode, f_createtime, f_vcode_error_cnt
        FROM t_vcode
        WHERE f_uuid = %s
        AND f_vcode_type = %s
        """
        results = self.r_db.one(select_sql, uuidIn, vcodeType)
        if results is None:
            raise_exception(exp_msg=_("IDS_VCODE_ERROR"),
                            exp_num=ncTShareMgntError.NCT_CHECK_VCODE_IS_WRONG)
        if delete_after_check:
            self.delete_vcode_info(uuidIn)

        vcode_str = results["f_vcode"]
        vcode_createtime = results["f_createtime"].strftime('%Y-%m-%d %H:%M:%S')
        vocde_vcode_error_cnt = results["f_vcode_error_cnt"] + 1

        # 验证码输入已达到限定次数
        if vocde_vcode_error_cnt > 3:
            raise_exception(exp_msg=_("IDS_VCODE_MORE_THAN_THE_LIMIT"),
                            exp_num=ncTShareMgntError.NCT_CHECK_VCODE_MORE_THAN_THE_LIMIT)
        # 验证码为空
        elif not vcodeIn:
            self.update_vcode_err_cnt(vocde_vcode_error_cnt, uuidIn)
            raise_exception(exp_msg=_("IDS_VCODE_NULL"),
                            exp_num=ncTShareMgntError.NCT_CHECK_VCODE_IS_NULL)

        # 验证码出错
        elif vcode_str.upper() != vcodeIn.upper():
            self.update_vcode_err_cnt(vocde_vcode_error_cnt, uuidIn)
            raise_exception(exp_msg=_("IDS_VCODE_ERROR"),
                            exp_num=ncTShareMgntError.NCT_CHECK_VCODE_IS_WRONG)

        # 验证码超时
        elif vcode_createtime < allow_max_last_useful_time:
            raise_exception(exp_msg=_("IDS_VCODE_TIMEOUT"),
                            exp_num=ncTShareMgntError.NCT_CHECK_VCODE_IS_TIMEOUT)


    def save_vcode_info(self, uuidIn, vcodeIn, vcodeType=ncTVcodeType.IMAGE_VCODE):
        now = BusinessDate.now().strftime("%Y-%m-%d %H:%M:%S")
        sql = """
        INSERT INTO t_vcode (f_uuid, f_vcode, f_vcode_type, f_createtime) VALUES (%s, %s, %s, %s)
        """
        self.w_db.query(sql, uuidIn, vcodeIn, vcodeType, now)

    def delete_vcode_info(self, uuidIn):
        # 检验是否存在验证码
        checkSql = """
        SELECT * FROM t_vcode WHERE f_uuid='%s'
        """ % self.w_db.escape(uuidIn)
        result = self.w_db.one(checkSql)

        if result:
            sql = """
            DELETE FROM t_vcode WHERE f_uuid='%s'
            """ % self.w_db.escape(uuidIn)
            self.w_db.query(sql)

    def is_user_need_check_vcode(self, db_user):
        """
        检查该用户是否需要校验
        """
        # 开启并且次数超过限制，则需要检查验证码
        vcode_config = self.get_vcode_config()
        if vcode_config.isEnable:
            # 1. 如果账户不存在，则不需要检查验证码，由前端来控制是否需要校验验证码
            if not db_user:
                return False

            # 2. 用户存在，且密码错误次数超过配置次数，需要校验验证码
            if db_user['f_pwd_error_cnt'] >= vcode_config.passwdErrCnt:
                return True

        return False

    def is_user_need_display_vcode(self, db_user):
        """
        检查输入的账号下一次是否需要显示验证码
        """
        # 1. 如果验证码未开启，不需要显示验证码
        vcode_config = self.get_vcode_config()
        if vcode_config.isEnable:
            # 2. 如果用户不存在，需要显示验证码
            if not db_user:
                return True

            # 3. 如果用户密码超过次数限制，需要显示验证码
            if db_user['f_pwd_error_cnt'] >= vcode_config.passwdErrCnt:
                return True

        return False

    def update_vcode_err_cnt(self, vcodeErrorCnt, uuidIn):
        """
        获取验证码输错更新次数
        """
        update_sql = """
        UPDATE t_vcode
        SET f_vcode_error_cnt = %s
        WHERE f_uuid = %s
        """
        self.w_db.query(update_sql, vcodeErrorCnt, uuidIn)

    def _check_vcode_by_third_auth(self, uuidIn):
        """
        第三方配置控制验证码的发送间隔
        """
        checkSql = """
        SELECT `f_createtime` FROM t_vcode WHERE `f_uuid` = %s
        """
        result = self.w_db.one(checkSql, uuidIn)
        if result:
            try:
                vcode_createtime = result["f_createtime"]
                now = BusinessDate.now().strftime('%Y-%m-%d %H:%M:%S')
                third_infos = self.third_config_manage.get_third_party_config(ncTPluginType.AUTHENTICATION)

                # 获取内外配置
                config = {}
                if third_infos and third_infos[0].enabled and third_infos[0].config:
                    try:
                        # 同时使用config和internalConfig
                        config = json.loads(third_infos[0].config)
                        config.update(json.loads(third_infos[0].internalConfig))
                    except Exception as ex:
                        ShareMgnt_Log("获取服务器配置信息异常: ex=%s", str(ex))

                interval = config.get("sendInterval", 60)
                if not isinstance(interval, int) or interval < 0:
                    interval = 60
                delta = datetime.timedelta(seconds=interval)
                expire_time = (vcode_createtime + delta).strftime('%Y-%m-%d %H:%M:%S')

                return True if expire_time > now else False
            except Exception as e:
                raise e
