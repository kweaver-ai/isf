# coding: utf-8
"""
第三方认证系统管理模块
"""
import os
import sys
import json
from importlib import reload
from src.third_party_auth.third_config_manage import ThirdConfigManage
from src.common.lib import raise_exception
from src.common.sharemgnt_logger import ShareMgnt_Log
from src.common.db.connector import DBConnector
from src.common.redis_connector import OPRedis
from src.common.plugin_lock import PluginVersion
from ShareMgnt.ttypes import (ncTShareMgntError,
                              ncTPluginType,
                              ncTMFAType,
                              ncTReturnInfo,
                              ncTVcodeCreateInfo,
                              ncTUsrmGetUserInfo,
                              ncTUsrmUserInfo)
from src.modules.vcode_manage import VcodeManage
from src.modules.user_manage import UserManage
from src.third_party_auth.auth_manage import auth_interface

from EThriftException.ttypes import ncTException

AUTH_CLASS = None
LOG_FLAG = 0


class ThirdAuthManage(DBConnector):
    """
    第三方认证管理模块
    """

    def __init__(self):
        """
        支持的第三方认证列表
        """
        self.third_cfg_manage = ThirdConfigManage()
        self.vcode_manage = VcodeManage()
        self.user_manage = UserManage()
        self.opredis = OPRedis()
        self.plugin_version = PluginVersion()
        self.config = None

        self.modify_time = 0
        self.import_error_str = "from auth_interface import *"

    def __check_file_modified(self, file_path):
        """
        检查插件是否被修改过
        """
        try:
            modified = os.stat(file_path).st_mtime
        except Exception:
            return True

        if not self.modify_time:
            self.modify_time = modified
            return True

        if modified == self.modify_time:
            return False
        else:
            self.modify_time = modified
            return True

    def __check_need_reload(self, file_path):
        """
        检查是否需要重新加载插件
        """
        if not AUTH_CLASS or self.__check_file_modified(file_path):
            return True
        else:
            return False

    def __log_once(self, log_type):
        """
        根据log_type判断是否记录日志
        通过此方法记录的日志只记录一次
        """
        global LOG_FLAG
        offset = 1 << log_type
        if LOG_FLAG & offset:
            return False
        else:
            LOG_FLAG = LOG_FLAG | offset
            return True

    def load_auth_module(self):
        """
        加载认证类（暂时只有一个）
        """
        third_infos = self.third_cfg_manage.get_third_party_config(ncTPluginType.AUTHENTICATION)
        if third_infos and third_infos[0].enabled and third_infos[0].config:
            try:
                # 同时使用config和internalConfig
                self.config = json.loads(third_infos[0].config)
                self.config.update(json.loads(third_infos[0].internalConfig))

                try:
                    import_str = "from third_party_auth.auth_manage import *"
                    exec(import_str)
                except ImportError:
                    # 源码调试时，需要从src目录导入
                    import_src_str = "from src.third_party_auth.auth_manage import *"
                    exec(import_src_str)

                # 导入第三方插件认证模块
                import_path, module = self.third_cfg_manage.get_format_plugin_path(third_infos[0], "auth_module.py")

                # 检查是否需要重新加载插件
                # if not self.__check_need_reload(module):
                #     return

                # 添加插件路径到环境变量
                plugin_path = "/sysvol/plugin/%s" % import_path
                if plugin_path not in sys.path:
                    sys.path.append(plugin_path)

               # 获取数据库中对象存储的版本
                auth_store_version = third_infos[0].plugin.objectId

                # 1 本地版本不存在, 2 本地版本与对象存储中版本不一致, 3 本地插件不存在, 则从对象存储中下载插件.
                if not PluginVersion.AUTH_LOCAL_VERSION \
                    or (auth_store_version != PluginVersion.AUTH_LOCAL_VERSION) \
                    or (not os.path.exists(plugin_path)):
                    try:
                        # 复制插件到对应路径，会覆盖原先内容
                        third_infos[0].plugin.data = self.third_cfg_manage.download_third_party_plugin(
                            third_infos[0].thirdPartyId)
                        if third_infos[0].plugin.data:
                            from src.third_party_auth.third_party_manage import ThirdPartyManage
                            ThirdPartyManage().add_local_third_party_plugin(third_infos[0].plugin)
                    except ncTException as ex:
                        ShareMgnt_Log("The auth plugin is download failed")
                        if ex.errID == 16778576:
                            raise_exception(exp_msg=ex.expMsg, exp_num=ex.errID)

                    if os.path.exists(module):
                        # 卸载已加载的第三方认证插件
                        self.third_cfg_manage.unload_plugin("auth_module")

                if os.path.exists(module):
                    # redis中锁的key
                    POD_ID = os.getenv("POD_IP", "127.0.0.1")
                    key = POD_ID + " " + str(third_infos[0].plugin.type)

                    # 获取锁
                    # import pdb; pdb.set_trace()
                    if not self.opredis.get_lock(key, 30):
                        authModule = __import__('auth_module')
                        reload(authModule)
                        import_str = "from auth_module import *"
                        exec(import_str)

                global AUTH_CLASS
                AUTH_CLASS = locals()[self.config['authModule']]

            except Exception as ex:
                raise_exception(exp_msg=_("IDS_LOAD_THIRD_PARTY_AUTH_MODULE_FAILED") % str(ex),
                                exp_num=ncTShareMgntError.NCT_FAILED_THIRD_CONFIG)
        else:
            raise_exception(exp_msg=_("third party auth no open"),
                            exp_num=ncTShareMgntError.NCT_THIRD_PARTY_AUTH_NOT_OPEN)

    def login(self, user_name, password, user):
        """
        第三方登录
        """
        self.load_auth_module()
        try:
            return AUTH_CLASS(self.config).login(user_name, password, user)
        except ValueError as ex:
            if ex.args[0] == ncTShareMgntError.NCT_INVALID_ACCOUNT_OR_PASSWORD:
                raise_exception(exp_msg=_("invalid account or password"),
                                exp_num=ncTShareMgntError.NCT_INVALID_ACCOUNT_OR_PASSWORD)
            elif ex.args[0] == ncTShareMgntError.NCT_CANNOT_CONNECT_THID_PARTY_SERVER:
                raise_exception(exp_msg=_("could not connect to third party auth server"),
                                exp_num=ncTShareMgntError.NCT_CANNOT_CONNECT_THID_PARTY_SERVER)
            else:
                raise ex

    def validate(self, validate_info):
        """
        第三方验证
        """
        if validate_info:
            validate_info = json.loads(validate_info)
            params = validate_info["params"]
            if "deviceinfo" in validate_info:
                params['deviceinfo'] = validate_info['deviceinfo']
        self.load_auth_module()
        try:
            return AUTH_CLASS(self.config).validate(params)
        except ValueError as ex:
            if ex.args[0] == ncTShareMgntError.NCT_INVALID_ACCOUNT_OR_PASSWORD:
                raise_exception(exp_msg=_("invalid account or password"),
                                exp_num=ncTShareMgntError.NCT_INVALID_ACCOUNT_OR_PASSWORD)
            elif ex.args[0] == ncTShareMgntError.NCT_CANNOT_CONNECT_THID_PARTY_SERVER:
                raise_exception(exp_msg=_("could not connect to third party auth server"),
                                exp_num=ncTShareMgntError.NCT_CANNOT_CONNECT_THID_PARTY_SERVER)
            else:
                raise ex

    def login_console_new(self, validate_info):
        """
        标准的第三方单点登录
        """
        try:
            if validate_info:
                validate_info_json = json.loads(validate_info)
                params = validate_info_json["params"]
                params["__eisoo_login_type__"] = "webconsole"
                validate_info = json.dumps(validate_info_json)
            return self.validate(validate_info)
        except Exception:
            raise

    def login_console(self, validate_info):
        """
        第三方单点登录
        """
        if validate_info:
            params = json.loads(validate_info)["params"]
        self.load_auth_module()
        return AUTH_CLASS(self.config).login_console(params)

    def validate_security_device(self, validate_info):
        """
        二次安全设备认证
        """
        if validate_info:
            params = json.loads(validate_info)["params"]
        self.load_auth_module()
        return AUTH_CLASS(self.config).validate_security_device(params)

    def get_third_auth_type_status(self, auth_type):
        """
        获取第三方认证插件是否开启
        """
        third_configs = self.third_cfg_manage.get_third_party_config(ncTPluginType.AUTHENTICATION)
        if not third_configs:
            return False

        for third_config in third_configs:
            if not third_config.enabled:
                return False

        # todo:后续会支持多认证插件，与现在的检测功能合并
        return True

    def convert_tel_num(self, user_info, telnum = None):
        """
        手机号脱敏处理
        """
        convert_way = lambda tel:tel[:3] + len(tel[3:-3]) * '*' + tel[-3:]
        if telnum and len(telnum) > 10:
            return convert_way(telnum)
        else:
            return convert_way(user_info.telNumber) if user_info.telNumber else '*' * 11

    def __check_params_from_plugin(self, params, check_type, error_string):
        """
        判断第三方插件传入的参数是否符合auth_interface中的定义
        params：传入的参数
        check_type：检查的类型
        error_string：提示的错误参数类型
        """
        for param in params:
            if not isinstance(param, check_type):
                raise TypeError("Illegal %s from third plugin:%s" % (error_string, param))

    def __create_auth_vcode(self, user_id, vcode_type):
        """
        生成MFA验证码
        """
        try:
            # 插件中生成验证码
            exec(self.import_error_str)
            user_info = self.user_manage.get_user_by_id(user_id).user
            auth_module = AUTH_CLASS(self.config)
            vcode_info = auth_module.create_auth_vcode_info(user_info)

            if vcode_info:
                # 参数检查
                self.__check_params_from_plugin([vcode_info], auth_interface.VcodeInfo, 'vcode_info')
                self.__check_params_from_plugin([vcode_info.vcode, vcode_info.uuid], str, 'vcode_info')
                self.__check_params_from_plugin([vcode_info.is_duplicate_sended], bool, 'vcode_info')

                return_vcode_info = ncTVcodeCreateInfo()
                return_vcode_info.vcode = vcode_info.vcode
                return_vcode_info.uuid = vcode_info.uuid
                return_vcode_info.isDuplicateSended = vcode_info.is_duplicate_sended
                return return_vcode_info

        except NotImplementedError:
            # 插件中没有生成验证码功能时，使用AS的验证码生成逻辑
            if self.__log_once(1):
                ShareMgnt_Log("There is no create_auth_vcode_info function in third plugin, implement in AnyShare")
            return self.vcode_manage.create_vcode_info(user_id, vcode_type)

    def __create_return_info(self, user_info, old_telnum, is_duplicateSended):
        """
        生成返回信息
        """
        final_return_info = ncTReturnInfo()
        try:
            exec(self.import_error_str)
            auth_module = AUTH_CLASS(self.config)
            return_info = auth_module.get_return_info(user_info, old_telnum, is_duplicateSended)

            if return_info:
                # 参数检查
                self.__check_params_from_plugin([return_info], auth_interface.ReturnInfo, 'return_info')
                self.__check_params_from_plugin([return_info.send_interval], int, 'return_info')
                self.__check_params_from_plugin([return_info.is_duplicate_sended], bool, 'return_info')
                self.__check_params_from_plugin([return_info.tel_number], str, 'return_info')
                if return_info.send_interval < 0:
                    return_info.send_interval = 60

                final_return_info.sendInterval = return_info.send_interval
                final_return_info.isDuplicateSended = return_info.is_duplicate_sended
                final_return_info.telNumber = self.convert_tel_num(user_info, return_info.tel_number)
                return final_return_info

        except NotImplementedError:
            if self.__log_once(2):
                ShareMgnt_Log("There is no get_return_info function in third plugin, implement in AnyShare")

            if not user_info.telNumber:
                raise auth_interface.UserHasNotBountPhone
            final_return_info.telNumber = self.convert_tel_num(user_info)
            if old_telnum and old_telnum != final_return_info.telNumber:
                raise auth_interface.PhoneNumberHasBeenChanged
            final_return_info.sendInterval = self.config.get("sendInterval", 60)
            if not isinstance(final_return_info.sendInterval, int) or final_return_info.sendInterval < 0:
                final_return_info.sendInterval = 60
            final_return_info.isDuplicateSended = is_duplicateSended

            return final_return_info

    def send_auth_vcode(self, user_id, vcode_type, old_telnum):
        """
        发送双因子认证短信验证码
        """
        try:
            exec(self.import_error_str)
            user_info = self.user_manage.get_user_by_id(user_id).user
            self.load_auth_module()
            auth_module = AUTH_CLASS(self.config)

            # 生成验证码
            vcode_info = self.__create_auth_vcode(user_id, vcode_type)

            # 生成返回信息
            return_info = self.__create_return_info(user_info, old_telnum, vcode_info.isDuplicateSended)

            # 验证码在发送间隔内重复发送，不再发送，直接返回
            if vcode_info.isDuplicateSended:
                return return_info

            # 发送验证码
            auth_module.send_auth_vcode(vcode_info, user_info)
            return return_info

        except NotImplementedError:
            # 发送失败，删除AS中的验证码
            self.vcode_manage.delete_vcode_info(user_id)
            raise_exception(exp_msg=(_("IDS_MFA_CONFIG_ERROR")),
                            exp_num=ncTShareMgntError.NCT_MFA_CONFIG_ERROR)

        except auth_interface.MFASMSSeverError as ex:
            # 发送失败，删除AS中的验证码
            self.vcode_manage.delete_vcode_info(user_id)
            raise_exception(exp_msg=_("could not connect to third party auth server"),
                            exp_num=ncTShareMgntError.NCT_MFA_SMS_SERVER_ERROR)

        except auth_interface.UserHasNotBountPhone as ex:
            self.vcode_manage.delete_vcode_info(user_id)
            raise_exception(exp_msg=_("you didn't bind the phone"),
                            exp_num=ncTShareMgntError.NCT_USER_HAS_NOT_BOUND_PHONE)

        except auth_interface.PhoneNumberHasBeenChanged as ex:
            self.vcode_manage.delete_vcode_info(user_id)
            raise_exception(exp_msg=_("your phone number has been changed"),
                            exp_num=ncTShareMgntError.NCT_PHONE_NUMBER_HAS_BEEN_CHANGED)

        except auth_interface.SendVerifyCodeFailed as ex:
            self.vcode_manage.delete_vcode_info(user_id)
            raise_exception(exp_msg=_("IDS_SEND_VERIFY_CODE_FAIL"),
                            exp_num=ncTShareMgntError.NCT_SEND_VERIFY_CODE_FAIL)

        except auth_interface.ThirdPluginInterError as ex:
            self.vcode_manage.delete_vcode_info(user_id)
            ShareMgnt_Log("Error in third plugin: %s" % ex.value)
            raise_exception(exp_msg=(_("IDS_THIRD_PLUGIN_INTER_ERROR") % ex.value),
                            exp_num=ncTShareMgntError.NCT_THIRD_PLUGIN_INTER_ERROR)

    def send_sms_vcode(self, telNumber, content):
        """
        发送短信验证码
        这里捕捉的异常，是参考上面 发送双因子认证短信验证码 接口定义的，但是省略了 NCT_USER_HAS_NOT_BOUND_PHONE、NCT_PHONE_NUMBER_HAS_BEEN_CHANGED 这两种异常情况。
        对于访问匿名外链，需要发送手机验证码的情况，这里的用户是任意的，只能获取到手机号信息，用户的其它信息无法获取。
        """
        try:
            exec(self.import_error_str)
            self.load_auth_module()
            auth_module = AUTH_CLASS(self.config)

            user_info = ncTUsrmUserInfo()
            user_info.telNumber = telNumber

            vcode_info = ncTVcodeCreateInfo()
            vcode_info.vcode = content

            # 发送验证码
            auth_module.send_auth_vcode(vcode_info, user_info)
        except NotImplementedError:
            raise_exception(exp_msg=(_("IDS_MFA_CONFIG_ERROR")),
                            exp_num=ncTShareMgntError.NCT_MFA_CONFIG_ERROR)

        except auth_interface.MFASMSSeverError as ex:
            raise_exception(exp_msg=_("could not connect to third party auth server"),
                            exp_num=ncTShareMgntError.NCT_MFA_SMS_SERVER_ERROR)

        except auth_interface.SendVerifyCodeFailed as ex:
            raise_exception(exp_msg=_("IDS_SEND_VERIFY_CODE_FAIL"),
                            exp_num=ncTShareMgntError.NCT_SEND_VERIFY_CODE_FAIL)

        except auth_interface.ThirdPluginInterError as ex:
            ShareMgnt_Log("Error in third plugin: %s" % ex.value)
            raise_exception(exp_msg=(_("IDS_THIRD_PLUGIN_INTER_ERROR") % ex.value),
                            exp_num=ncTShareMgntError.NCT_THIRD_PLUGIN_INTER_ERROR)

    def sms_validate(self, user_id, auth_vcode, vcode_type, delete_after_check):
        """
        双因子认证（短信验证码）
        """
        try:
            exec(self.import_error_str)
            self.load_auth_module()
            user_info = self.user_manage.get_user_by_id(user_id).user

            try:
                # 去除验证码两边的空格
                auth_vcode = auth_vcode.strip()
                AUTH_CLASS(self.config).sms_validate(auth_vcode, user_info)
            except NotImplementedError as ex:
                if self.__log_once(3):
                    ShareMgnt_Log("There is no sms_validate module in third plugin, validate in AnyShare")

                self.vcode_manage.verify_vcode_info(user_id, auth_vcode, vcode_type, delete_after_check)
                self.vcode_manage.delete_vcode_info(user_id)

        except auth_interface.MFASMSSeverError as ex:
            raise_exception(exp_msg=_("could not connect to third party auth server"),
                            exp_num=ncTShareMgntError.NCT_MFA_SMS_SERVER_ERROR)

        except auth_interface.CheckVcodeIsTimeout as ex:
            raise_exception(exp_msg=_("IDS_VCODE_TIMEOUT"),
                            exp_num=ncTShareMgntError.NCT_CHECK_VCODE_IS_TIMEOUT)

        except auth_interface.CheckVcodeWrong as ex:
            raise_exception(exp_msg=_("IDS_VCODE_ERROR"),
                            exp_num=ncTShareMgntError.NCT_CHECK_VCODE_IS_WRONG)

        except auth_interface.ThirdPluginInterError as ex:
            ShareMgnt_Log("Error in third plugin: %s" % ex.value)
            raise_exception(exp_msg=(_("IDS_THIRD_PLUGIN_INTER_ERROR") % ex.value),
                            exp_num=ncTShareMgntError.NCT_THIRD_PLUGIN_INTER_ERROR)

    def OTP_validate(self, validate_info, user_id):
        """
        双因子认证（动态密码）
        """
        exec(self.import_error_str)
        self.load_auth_module()
        try:
            # 去除空格
            validate_info = validate_info.strip()
            user_info = self.user_manage.get_user_by_id(user_id).user
            AUTH_CLASS(self.config).OTP_validate(validate_info, user_info)

        except NotImplementedError:
            raise_exception(exp_msg=(_("IDS_MFA_CONFIG_ERROR")),
                            exp_num=ncTShareMgntError.NCT_MFA_CONFIG_ERROR)

        except auth_interface.MFAOTPServerError as ex:
            raise_exception(exp_msg=_("could not connect to third party auth server"),
                            exp_num=ncTShareMgntError.NCT_MFA_OTP_SERVER_ERROR)

        except auth_interface.OTPWrong as ex:
            raise_exception(exp_msg=_("IDS_OTP_WRONG"),
                            exp_num=ncTShareMgntError.NCT_OTP_WRONG)

        except auth_interface.OTPTimeout as ex:
            raise_exception(exp_msg=_("IDS_OTP_WRONG"),
                            exp_num=ncTShareMgntError.NCT_OTP_TIMEOUT)

        except auth_interface.OTPTooManyWrongTime as ex:
            raise_exception(exp_msg=_("IDS_OTP_WRONG"),
                            exp_num=ncTShareMgntError.NCT_OTP_TOO_MANY_WRONG_TIME)

        except auth_interface.ThirdPluginInterError as ex:
            ShareMgnt_Log("Error in third plugin: %s" % ex.value)
            raise_exception(exp_msg=(_("IDS_THIRD_PLUGIN_INTER_ERROR") % ex.value),
                            exp_num=ncTShareMgntError.NCT_THIRD_PLUGIN_INTER_ERROR)