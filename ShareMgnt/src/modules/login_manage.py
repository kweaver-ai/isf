# -*- coding:utf-8 -*-
"""
用户认证模块
"""
import hashlib
import hmac
import json
import codecs
import binascii
from Crypto.Cipher import DES
from eisoo.tclients import TClient
from src.common import global_info
from src.common.lib import (raise_exception,
                            encrypt_pwd,
                            sha2_encrypt,
                            check_is_strong_password)
from src.common.db.connector import DBConnector
from src.modules.ldap_manage import LdapManage
from src.modules.domain_manage import DomainManage
from src.third_party_auth.third_auth_manage import ThirdAuthManage
from src.third_party_auth.third_config_manage import ThirdConfigManage
from src.modules.user_manage import UserManage
from src.modules.login_access_control_manage import LoginAccessControlManage
from src.modules.config_manage import ConfigManage
from src.modules.openapi import OpenApi
from src.modules.vcode_manage import VcodeManage
from src.common.encrypt.simple import (des_encrypt_with_padzero,eisoo_rsa_decrypt,eisoo_rsa2048_decrypt)

from ShareMgnt.constants import (NCT_USER_ADMIN,
                                 NCT_USER_SYSTEM,
                                 NCT_USER_AUDIT,
                                 NCT_USER_SECURIT,
                                 NCT_SYSTEM_ROLE_SUPPER,
                                 NCT_SYSTEM_ROLE_ADMIN,
                                 NCT_SYSTEM_ROLE_SECURIT,
                                 NCT_SYSTEM_ROLE_AUDIT,
                                 NCT_SYSTEM_ROLE_ORG_MANAGER,
                                 NCT_SYSTEM_ROLE_ORG_AUDIT)
from ShareMgnt.ttypes import (ncTUsrmAuthenType,
                              ncTUsrmUserType,
                              ncTUsrmUserStatus,
                              ncTShareMgntError,
                              ncTNTLMResponse,
                              ncTVcodeType)
from EACP.ttypes import ncTCheckTokenInfo
from EThriftException.ttypes import ncTException
from src.modules.role_manage import RoleManage


class LoginManage(DBConnector):
    """
    LoginManage module
    """
    def __init__(self):
        self.third_auth_manage = ThirdAuthManage()
        self.third_config_manage = ThirdConfigManage()
        self.user_manage = UserManage()
        self.domain_manage = DomainManage()
        self.login_access_control_manage = LoginAccessControlManage()
        self.config_manage = ConfigManage()
        self.openapi = OpenApi()
        self.vcode_manage = VcodeManage()
        self.role_manage = RoleManage()

    def login(self, user_name, password, authen_type, option=None, b2048=False):
        """
        登陆集成接口
        """
        try:
            # 根据输入账号名获取用户信息
            db_user = self.match_account(user_name, authen_type)

            # 校验登录验证码
            uuid = vcode = OTP = ""
            isModify = False
            if option:
                uuid = option.uuid if option.uuid else uuid
                vcode = option.vcode if option.vcode else vcode
                isModify = True if option.isModify else isModify
                OTP = option.OTP if option.OTP else ""

            # 校验短信验证
            if option and option.vcodeType and option.vcodeType == ncTVcodeType.DAUL_AUTH_VCODE:
                if db_user:
                    is_delete_after_check = False
                    uuid = db_user["f_user_id"]
                    self.third_auth_manage.sms_validate(uuid, vcode, option.vcodeType, is_delete_after_check)
                else:
                    raise_exception(exp_msg=_("IDS_VCODE_ERROR"),
                            exp_num=ncTShareMgntError.NCT_CHECK_VCODE_IS_WRONG)
                    
            # 校验OTP
            if OTP:
                if db_user:
                    self.third_auth_manage.OTP_validate(OTP, db_user["f_user_id"])
                else:
                    raise_exception(exp_msg=_("IDS_OTP_WRONG"),
                            exp_num=ncTShareMgntError.NCT_OTP_WRONG)
                    
            # 校验图片验证码
            if (uuid or self.vcode_manage.is_user_need_check_vcode(db_user)) and (not isModify):
                self.vcode_manage.verify_vcode_info(uuid, vcode)

            # 用户账号不存在
            if not db_user:
                raise_exception(exp_msg=_("invalid account or password"),
                                exp_num=ncTShareMgntError.NCT_INVALID_ACCOUNT_OR_PASSWORD)

            # 获取admin account map
            admin_account_map = self.user_manage.get_all_admin_account()

            # 判断用户是登录webconsole还是webclient
            if authen_type == ncTUsrmAuthenType.NCT_AUTHEN_TYPE_MANAGER:
                return self.login_console(db_user, password, option, admin_account_map, b2048)

            return self.login_client(db_user, password, option)
        except ncTException as ex:
            if ex.errID == ncTShareMgntError.NCT_INVALID_ACCOUNT_OR_PASSWORD and    \
                self.vcode_manage.get_vcode_config().isEnable and not self.user_manage.get_password_config().lockStatus:
                    if db_user:
                        self.user_manage.modify_pwd_lock_info(db_user['f_user_id'], False)

            if db_user:
                # 用户名或密码错误时，会更新用户密码错误次数，db_user 中存放的是未更改的数据，需要重新读取
                db_user = self.match_account(user_name)
            detail = {}
            if ex.errDetail:
                detail = json.loads(ex.errDetail)
            detail["isShowStatus"] = self.vcode_manage.is_user_need_display_vcode(db_user)
            raise_exception(exp_msg=ex.expMsg, exp_num=ex.errID, exp_detail=json.dumps(detail))

    def login_with_console_log(self, user_name, password, authen_type, option=None , os_type=''):
        """
        用户登录WEB，如果登录控制台失败，记录日志
        """
        try:
            return self.login(user_name, password, authen_type, option, True)
        except ncTException as ex:
            bHasRaise = False
            if authen_type == ncTUsrmAuthenType.NCT_AUTHEN_TYPE_MANAGER:
                db_user = self.match_account(user_name)
                user_id = ''
                if db_user:
                    user_id = db_user["f_user_id"]

                if ex.errID == ncTShareMgntError.NCT_INVALID_ACCOUNT_OR_PASSWORD or \
                    ex.errID  == ncTShareMgntError.NCT_PWD_THIRD_FAILED:
                    detail = {}
                    if ex.errDetail:
                        detail = json.loads(ex.errDetail)
                    detail["id"] = user_id
                    bHasRaise = True
                    raise_exception(exp_msg=ex.expMsg, exp_num=ex.errID, exp_detail=json.dumps(detail))

            if bHasRaise == False:
                raise ex

    def get_db_user_by_account(self, account):
        """
        根据用户账号获取用户信息
        """
        db_user = self.match_account(account)
        # 用户账号不存在
        if not db_user:
            raise_exception(exp_msg=_("invalid account or password"),
                            exp_num=ncTShareMgntError.NCT_INVALID_ACCOUNT_OR_PASSWORD)

        return db_user

    def login_console(self, db_user, password, option=None, admin_account_map=None, b2048=False):
        """
        登陆控制台
        """
        # 控制台密码进行rsa解密
        if not option or not option.isPlainPwd:
            try:
                if b2048:
                    password = bytes.decode(eisoo_rsa2048_decrypt(password))
                else:
                    password = bytes.decode(eisoo_rsa_decrypt(password))
            except Exception:
                raise_exception(exp_msg=_("invalid paramter"),
                                exp_num=ncTShareMgntError.NCT_INVALID_PARAMTER)

        user_id = self.switch_login(db_user, password)

        self.check_login_console(db_user, option, admin_account_map)

        # 记录用户登录控制台请求时间
        self.user_manage.update_user_last_request_time(user_id)

        return user_id

    def check_login_console(self, db_user, option=None, admin_account_map=None):
        """
        检查登录控制台账号
        """
        if not admin_account_map:
            admin_account_map = self.user_manage.get_all_admin_account()

        # 系统未初始化时，非admin管理员账号不能登录 (涉密模式)
        user_name = db_user['f_login_name'].lower()
        if (user_name != admin_account_map[NCT_USER_ADMIN] and
                self.config_manage.get_secret_mode_status() and
                not self.config_manage.get_system_init_status()):
            raise_exception(exp_msg=_("IDS_ACCOUNT_CANNOT_LOGIN_IN_SECRET_NODE"),
                            exp_num=ncTShareMgntError.NCT_ACCOUNT_CANNOT_LOGIN_IN_SECRET_NODE)

        # 管理员system被禁用
        if db_user['f_login_name'].lower() == admin_account_map[NCT_USER_SYSTEM]:
            raise_exception(exp_msg=_("invalid account or password"),
                            exp_num=ncTShareMgntError.NCT_INVALID_ACCOUNT_OR_PASSWORD)

        # 非三权分立情况下不允许audit和security账户登录
        if user_name in [admin_account_map[NCT_USER_AUDIT],
                         admin_account_map[NCT_USER_SECURIT]]:
            if not self.user_manage.get_trisystem_status():
                raise_exception(exp_msg=_("invalid account or password"),
                                exp_num=ncTShareMgntError.NCT_INVALID_ACCOUNT_OR_PASSWORD)

        #  检查账号是否已被禁用，系统内置管理员账号不检查
        if user_name not in list(admin_account_map.values()):
            if self.check_user_disable(db_user):
                raise_exception(exp_msg=_("user is disabled"),
                                exp_num=ncTShareMgntError.NCT_USER_DISABLED)

    def check_user_disable(self, db_user):
        """
        检查用户是否禁用
        """
        if (db_user['f_status'] == ncTUsrmUserStatus.NCT_STATUS_ENABLE and
                db_user['f_auto_disable_status'] == 0):
            return False
        return True

    def check_user_activate(self, db_user):
        """
        检查用户是否已激活
        """
        return True if (db_user['f_is_activate']) == 1 else False

    def login_client(self, db_user, password, option=None):
        """
        登录客户端
        """
        user_id = self.switch_login(db_user, password)

        self.check_login_client(db_user, option)

        # 记录用户登录客户端 请求时间
        self.user_manage.update_user_last_request_time(user_id)

        # 记录用户登录客户端 客户端请求时间
        self.user_manage.update_user_last_client_request_time(user_id)
        return user_id

    def check_login_client(self, db_user, option=None):
        """
        检查登录客户端账号
        """

        # 检查账号是否已被禁用
        if self.check_user_disable(db_user):
            # 短信激活开启时，检查用户是否已激活
            if self.config_manage.get_sms_activate_status() and not self.check_user_activate(db_user):
                raise_exception(exp_msg=_("IDS_USER_NOT_ACTIVATE"),
                                exp_num=ncTShareMgntError.NCT_USER_NOT_ACTIVATE)
            raise_exception(exp_msg=_("user is disabled"),
                            exp_num=ncTShareMgntError.NCT_USER_DISABLED)

    def match_account(self, user_name, authen_type = None):
        """
        根据输入的帐号匹配出可能的帐号信息
        """
        # 优先使用login_name进行登录，然后是身份证号登录
        sql = """
        SELECT * FROM `t_user`
        WHERE `f_login_name` = %s
        """
        db_user = self.r_db.one(sql, user_name)

        if (not db_user and authen_type and authen_type != ncTUsrmAuthenType.NCT_AUTHEN_TYPE_MANAGER and
            self.config_manage.get_custom_config_of_bool("id_card_login_status")):
            user_idcard = bytes.decode(des_encrypt_with_padzero(global_info.des_key,
                        user_name,
                        global_info.des_key))
            sql = """
            SELECT * FROM `t_user`
            WHERE `f_idcard_number` = %s
            """
            db_user = self.r_db.one(sql, user_idcard)

        if not db_user:
            # 判断域是否开启
            domain_infos = self.domain_manage.get_all_domains()
            for domain_info in domain_infos:
                if domain_info.status:
                    # 搜索域帐号
                    sql = """
                    SELECT * FROM `t_user`
                    WHERE `f_login_name` like %s and f_auth_type = %s
                    order by f_login_name asc
                    """
                    db_user = self.r_db.one(sql, user_name + "@%%",
                                            ncTUsrmUserType.NCT_USER_TYPE_DOMAIN)

                    # 为了处理：域用户如果登录名为zhangying@qq.com@test2.develop.cn，
                    # 而zhangying, zhangying@qq.com zhangying@qq.com@test2.develop.cn
                    # 都可以登录的问题
                    if (db_user and
                            user_name.upper() != db_user['f_login_name'].rsplit('@', 1)[0].upper()):
                        db_user = None

                    break

            # 检查广西检察院认证是否开启
            third_auth_info = self.third_config_manage.get_third_party_info_auth()
            if third_auth_info.enabled and third_auth_info.thirdPartyId == 'softanywhere':
                # 搜索登录用户的显示名
                sql = """
                SELECT * FROM `t_user`
                WHERE `f_display_name` = %s and `f_auth_type` = %s
                """
                db_user = self.r_db.one(sql, user_name,
                                        ncTUsrmAuthenType.NCT_AUTHEN_TYPE_THIRD)

        return db_user

    def switch_login(self, db_user, password):
        """
        根据用户授权类型选择登陆方式
        """
        if db_user['f_auth_type'] == ncTUsrmUserType.NCT_USER_TYPE_LOCAL:
            # 本地用户登录
            return self.__check_user_login_password(db_user, password,
                                                    self.__check_local_password_func)
        elif db_user['f_auth_type'] == ncTUsrmUserType.NCT_USER_TYPE_DOMAIN:
            # 域控用户通过域控验证
            return self.__check_user_login_password(db_user, password,
                                                    self.__check_domain_password_func)
        elif db_user['f_auth_type'] == ncTUsrmUserType.NCT_USER_TYPE_THIRD:
            # 第三方认证
            return self.__check_user_login_password(db_user, password,
                                                    self.__check_third_password_func)
        else:
            raise_exception(exp_msg=_("forbidden login"),
                            exp_num=ncTShareMgntError.NCT_FORBIDDEN_LOGIN)

    def __check_local_password_func(self, db_user, password):
        """
        本地本地用户密码
        """
        if db_user['f_password'].lower() == "":
            return db_user['f_sha2_password'] == sha2_encrypt(password)
        else:
            return db_user['f_password'].lower() == encrypt_pwd(password)

    def __check_third_password_func(self, db_user, password):
        """
        本地第三方用户密码
        """
        try:
            self.third_auth_manage.login(db_user['f_login_name'], password, db_user)
            return True
        except ncTException as ex:
            if ex.errID == ncTShareMgntError.NCT_INVALID_ACCOUNT_OR_PASSWORD:
                return False
            else:
                raise ex

    def check_third_password_func(self, loginName, password):
        """
        本地第三方用户密码
        """
        try:
            db_user = self.match_account(loginName, ncTUsrmAuthenType.NCT_AUTHEN_TYPE_NORMAL)
            self.third_auth_manage.login(db_user['f_login_name'], password, db_user)
            return True
        except ncTException as ex:
            if ex.errID == ncTShareMgntError.NCT_INVALID_ACCOUNT_OR_PASSWORD:
                return False
            else:
                raise ex

    def __check_domain_password_func(self, db_user, password):
        """
        域控用户登录
        """
        return self.domain_manage.user_login(db_user['f_login_name'], db_user['f_ldap_server_type'], db_user['f_domain_path'], password)

    def check_domain_password_func(self, loginName, ldapServerType, domainPath, password):
        """
        域控用户登录
        """
        return self.domain_manage.user_login(loginName, ldapServerType, domainPath, password)

    def __check_user_login_password(self, db_user, password, func):
        """
        检查用户登录密码
        """
        user_id = db_user['f_user_id']
        pwd_config = self.user_manage.get_password_config()

        # 密码错误
        b_handle_pwd_err = False
        if not func(db_user, password):
            if (pwd_config.lockStatus and
                    (db_user['f_auth_type'] == ncTUsrmUserType.NCT_USER_TYPE_LOCAL or
                     self.config_manage.get_third_pwd_lock())):
                    # 更新用户登录锁定信息
                    self.user_manage.modify_pwd_lock_info(user_id, False)

            raise_exception(exp_msg=_("invalid account or password"),
                            exp_num=ncTShareMgntError.NCT_INVALID_ACCOUNT_OR_PASSWORD)
        else:
            b_handle_pwd_err = True

        # 检查用户是否被锁定
        if pwd_config.lockStatus:
            # 判断域认证和第三方认证是否开启限制
            if (db_user['f_auth_type'] == ncTUsrmUserType.NCT_USER_TYPE_LOCAL or
                    self.config_manage.get_third_pwd_lock()):
                if self.user_manage.check_user_locked(user_id):
                    self.user_manage.handle_pwd_lock_except(user_id)
        
        # 如果正常登录，则重置锁定
        if b_handle_pwd_err:
            self.user_manage.modify_pwd_lock_info(user_id, True)

        # 如果是系统账户，需要检查用户密码是否是初始密码
        if db_user['f_user_id'] in [NCT_USER_ADMIN, NCT_USER_AUDIT, NCT_USER_SYSTEM, NCT_USER_SECURIT]:
            if db_user['f_password'] != "" and encrypt_pwd(password) == encrypt_pwd(self.user_manage.admin_default_passwd):
                raise_exception(exp_msg=_("password is initial"),
                                exp_num=ncTShareMgntError.
                                NCT_PASSWORD_IS_INITIAL)
            elif db_user['f_sha2_password'] != "" and sha2_encrypt(password) == sha2_encrypt(self.user_manage.admin_default_passwd):
                raise_exception(exp_msg=_("password is initial"),
                                exp_num=ncTShareMgntError.
                                NCT_PASSWORD_IS_INITIAL)

        # 如果是本地用户类型，需要检查用户密码是否是初始密码
        if db_user['f_auth_type'] == ncTUsrmUserType.NCT_USER_TYPE_LOCAL:
            # 旧版本默认密码为“123456”
            if db_user['f_password'] != "" and encrypt_pwd(password) == self.user_manage.user_default_password.md5_pwd:
                raise_exception(exp_msg=_("password is initial"),
                                exp_num=ncTShareMgntError.
                                NCT_PASSWORD_IS_INITIAL)
            elif db_user['f_sha2_password'] != "" and sha2_encrypt(password) == self.user_manage.user_default_password.sha2_pwd:
                raise_exception(exp_msg=_("password is initial"),
                                exp_num=ncTShareMgntError.
                                NCT_PASSWORD_IS_INITIAL)

            # 检查强密码
            pwd_config = self.user_manage.get_password_config()
            if pwd_config.strongStatus:
                b_valid = check_is_strong_password(password)
                if not b_valid:
                    raise_exception(exp_msg=_("password not safe"),
                                    exp_num=ncTShareMgntError.NCT_PASSWORD_NOT_SAFE)

            # 检查用户密码是否过期
            if pwd_config.expireTime != -1:
                b_expire = self.user_manage.check_password_expire(user_id)
                if b_expire:
                    if db_user['f_pwd_control']:
                        raise_exception(exp_msg=_("IDS_CONTROLED_PASSWORD_EXPIRED"),
                                        exp_num=ncTShareMgntError.NCT_CONTROLED_PASSWORD_EXPIRE)
                    raise_exception(exp_msg=_("password expire"),
                                    exp_num=ncTShareMgntError.NCT_PASSWORD_EXPIRE)

        return user_id

    def login_by_ntlmv1(self, account, challenge, password):
        """
        根据ntlmv1进行登录
        """
        # 根据输入账号名匹配可能的账号信息
        db_user = self.get_db_user_by_account(account)

        # 获取admin account map
        admin_account_map = self.user_manage.get_all_admin_account()

        # 如果是系统管理员，不允许登陆
        if db_user["f_login_name"].lower() in list(admin_account_map.values()):
            raise_exception(exp_msg=_("invalid account or password"),
                            exp_num=ncTShareMgntError.NCT_INVALID_ACCOUNT_OR_PASSWORD)

        # 检查用户登录客户端合法性
        self.check_login_client(db_user)

        pwd_config = self.user_manage.get_password_config()

        if pwd_config.lockStatus:
            # 检查用户是否被锁定
            if self.user_manage.check_user_locked(db_user['f_user_id']):
                self.user_manage.handle_pwd_lock_except(db_user['f_user_id'])

        # 密码错误
        if db_user['f_ntlm_password'] is None or len(db_user['f_ntlm_password']) == 0:
            raise_exception(exp_msg=_("password not safe"),
                            exp_num=ncTShareMgntError.NCT_PASSWORD_NOT_SAFE)

        ok, sesskey = self.ntlmv1(db_user['f_ntlm_password'], challenge, password)
        if not ok:
            if pwd_config.lockStatus:
                # 更新用户登录锁定信息
                self.user_manage.modify_pwd_lock_info(db_user['f_user_id'], False)

                # 密码错误异常处理
                self.user_manage.handle_pwd_err_except(db_user['f_user_id'])
            else:
                raise_exception(exp_msg=_("invalid account or password"),
                                exp_num=ncTShareMgntError.NCT_INVALID_ACCOUNT_OR_PASSWORD)
        # 密码通过验证，则清空失败记录次数
        else:
            self.user_manage.modify_pwd_lock_info(db_user['f_user_id'], True)

        # 检查用户密码是否过期
        if pwd_config.expireTime != -1:
            b_expire = self.user_manage.check_password_expire(db_user['f_user_id'])
            if b_expire:
                if db_user['f_pwd_control']:
                    raise_exception(exp_msg=_("IDS_CONTROLED_PASSWORD_EXPIRED"),
                                    exp_num=ncTShareMgntError.NCT_CONTROLED_PASSWORD_EXPIRE)
                raise_exception(exp_msg=_("password expire"),
                                exp_num=ncTShareMgntError.NCT_PASSWORD_EXPIRE)

        ret = ncTNTLMResponse()
        ret.userId = db_user['f_user_id']
        ret.sessKey = sesskey
        return ret

    def des56(self, C, p7):
        key = bytearray(8)
        key[0] = p7[0] >> 1
        key[1] = ((p7[0] & 0x01) << 6) | (p7[1] >> 2)
        key[2] = ((p7[1] & 0x03) << 5) | (p7[2] >> 3)
        key[3] = ((p7[2] & 0x07) << 4) | (p7[3] >> 4)
        key[4] = ((p7[3] & 0x0F) << 3) | (p7[4] >> 5)
        key[5] = ((p7[4] & 0x1F) << 2) | (p7[5] >> 6)
        key[6] = ((p7[5] & 0x3F) << 1) | (p7[6] >> 7)
        key[7] = p7[6] & 0x7F
        for i in range(8):
            key[i] = (key[i] << 1)

        des = DES.new(key,DES.MODE_ECB)
        return binascii.b2a_hex(des.encrypt(C))

    def ntlmv1(self, pmd4, C, password):
        pmd4 = binascii.a2b_hex(pmd4 + "0" * 10)
        C = binascii.a2b_hex(C)
        return (self.des56(C, pmd4[:7]) + self.des56(C, pmd4[7:14]) + self.des56(C, pmd4[14:]) == password.encode('utf-8'),
                hashlib.new("md4", pmd4[:16]).hexdigest())

    def login_by_ntlmv2(self, account, domain, challenge, password):
        """
        根据ntlmv2进行登录
        """
        # 根据输入账号名匹配可能的账号信息
        db_user = self.get_db_user_by_account(account)

        # 获取admin account map
        admin_account_map = self.user_manage.get_all_admin_account()

        # 如果是系统管理员，不允许登陆
        if db_user["f_login_name"].lower() in list(admin_account_map.values()):
            raise_exception(exp_msg=_("invalid account or password"),
                            exp_num=ncTShareMgntError.NCT_INVALID_ACCOUNT_OR_PASSWORD)

        # 检查用户登录客户端合法性
        self.check_login_client(db_user)

        pwd_config = self.user_manage.get_password_config()
        if pwd_config.lockStatus:
            # 检查用户是否被锁定
            if self.user_manage.check_user_locked(db_user['f_user_id']):
                self.user_manage.handle_pwd_lock_except(db_user['f_user_id'])

        # 密码错误
        if db_user['f_ntlm_password'] is None or len(db_user['f_ntlm_password']) == 0:
            raise_exception(exp_msg=_("password not safe"),
                            exp_num=ncTShareMgntError.NCT_PASSWORD_NOT_SAFE)

        ok, sesskey = self.ntlmv2(account, domain, challenge, password, db_user['f_ntlm_password'])
        if not ok:
            if pwd_config.lockStatus:
                # 更新用户登录锁定信息
                self.user_manage.modify_pwd_lock_info(db_user['f_user_id'], False)

                # 密码错误异常处理
                self.user_manage.handle_pwd_err_except(db_user['f_user_id'])
            else:
                raise_exception(exp_msg=_("invalid account or password"),
                                exp_num=ncTShareMgntError.NCT_INVALID_ACCOUNT_OR_PASSWORD)
        # 密码通过验证，则清空失败记录次数
        else:
            self.user_manage.modify_pwd_lock_info(db_user['f_user_id'], True)

        # 检查用户密码是否过期
        if pwd_config.expireTime != -1:
            b_expire = self.user_manage.check_password_expire(db_user['f_user_id'])
            if b_expire:
                if db_user['f_pwd_control']:
                    raise_exception(exp_msg=_("IDS_CONTROLED_PASSWORD_EXPIRED"),
                                    exp_num=ncTShareMgntError.NCT_CONTROLED_PASSWORD_EXPIRE)
                raise_exception(exp_msg=_("password expire"),
                                exp_num=ncTShareMgntError.NCT_PASSWORD_EXPIRE)

        ret = ncTNTLMResponse()
        ret.userId = db_user['f_user_id']
        ret.sessKey = sesskey
        return ret

    def ntv2_owf_gen(self, pmd4, account, domain):
        md5 = hmac.new(pmd4, digestmod='md5')
        md5.update(account.upper().encode("utf-16le"))
        md5.update(domain.encode("utf-16le"))
        return md5.digest()

    def ntlmv2(self, account, domain, challenge, password, ntlm_password):
        decode_hex = codecs.getdecoder("hex_codec")
        kr = self.ntv2_owf_gen(decode_hex(ntlm_password)[0], account, domain)
        tmp = decode_hex(password)[0]
        cli_chal = tmp[16:]

        md5 = hmac.new(kr, digestmod='md5')
        md5.update(decode_hex(challenge)[0])
        md5.update(cli_chal)

        key = hmac.new(kr, digestmod='md5')
        key.update(tmp[:16])

        return md5.digest() == tmp[:16], key.hexdigest()

    def validate(self, token):
        """
        根据Token信息向第三方认证平台进行认证
        """
        return self.third_auth_manage.validate(token)

    def login_console_by_third_party(self, params):
        """
        单点登录控制台
        """
        return self.third_auth_manage.login_console(params)

    def login_console_by_third_party_new(self, params):
        """
        标准的第三方单点登录控制台
        """
        user_name = self.third_auth_manage.login_console_new(params)
        db_user = self.match_account(user_name)
        # 用户账号不存在
        if not db_user:
            raise_exception(exp_msg=_("user not exists"),
                            exp_num=ncTShareMgntError.NCT_USER_NOT_EXIST)
        self.check_login_console(db_user)

         # 记录用户第三方登录请求时间
        self.user_manage.update_user_last_request_time(db_user["f_user_id"])
        return user_name

    def get_login_client_info(self, userid):
        """
        获取用户登录webclient的信息
        """
        # 1. 根据用户id获取用户account
        user_info = self.user_manage.get_t_user_info_by_id(userid)
        if not user_info:
            raise_exception(_("userid not exist"))
        account = user_info['f_login_name']

        # 2. 获取用户登录web客户端appid: LoginWebClient
        appid, appkey = self.openapi.get_web_client_auth_info()

        # 3. 生成签名
        md5 = hashlib.md5()
        md5.update(appid.encode('utf-8'))
        md5.update(appkey.encode('utf-8'))
        md5.update(account.encode('utf-8'))
        sign = md5.hexdigest().lower()

        data = {}
        data["account"] = account
        data["appid"] = appid
        data["key"] = sign

        return json.dumps(data)
