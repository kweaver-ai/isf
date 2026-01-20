#!/usr/bin/python3
# -*- coding:utf-8 -*-


class AuthInterface(object):
    def __init__(self, config):
        self.config = config

    def login(self, user_name, passwd, db_user):
        """账号密码登录接口"""
        raise NotImplementedError("Method not implement")

    def validate(self, params):
        """单点登录验证接口"""
        raise NotImplementedError("Method not implement")

    def get_return_info(self, user_info, old_telnum, is_duplicate_sended):
        """"返回信息（用户手机号、发送间隔、是否已经重复发送了验证码）生成接口"""
        raise NotImplementedError("Method not implement")

    def create_auth_vcode_info(self, user_info):
        """生成验证码接口"""
        raise NotImplementedError("Method not implement")

    def send_auth_vcode(self, vcode_info, user_info):
        """发送验证码接口"""
        raise NotImplementedError("Method not implement")

    def sms_validate(self, auth_vcode, user_info):
        """短信验证码接口"""
        raise NotImplementedError("Method not implement")

    def OTP_validate(self, validate_info, user_info):
        """动态密码验证接口"""
        raise NotImplementedError("Method not implement")

class ReturnInfo(object):
    """
    获取短信验证码时返回给前端的信息
    """
    def __init__(self, tel_number = "", send_interval = 60, is_duplicate_sended = False):
        self.tel_number = tel_number
        self.send_interval = send_interval
        self.is_duplicate_sended = is_duplicate_sended

class VcodeInfo(object):
    """
    验证码信息
    """
    def __init__(self, vcode = "", uuid = "", is_duplicate_sended = False):
        self.vcode = vcode
        self.uuid = uuid
        self.is_duplicate_sended = is_duplicate_sended

class ThirdPluginError(Exception):
    """
    第三方认证插错误
    """
    def __init__(self, value = ""):
        self.value = value
    def __str__(self):
        return repr(self.value)

class ThirdPluginInterError(ThirdPluginError):
    """
    第三方认证插件内部其他错误
    """
    def __init__(self, value = ""):
        super(ThirdPluginInterError, self).__init__(value)

class MFASMSSeverError(ThirdPluginError):
    """
    短信服务器出错
    """
    def __init__(self):
        super(MFASMSSeverError, self).__init__()

class UserHasNotBountPhone(ThirdPluginError):
    """
    用户没有绑定手机
    """
    def __init__(self):
        super(UserHasNotBountPhone, self).__init__()

class PhoneNumberHasBeenChanged(ThirdPluginError):
    """
    用户手机号被修改
    """
    def __init__(self):
        super(PhoneNumberHasBeenChanged, self).__init__()

class SendVerifyCodeFailed(ThirdPluginError):
    """
    验证码发送失败
    """
    def __init__(self):
        super(SendVerifyCodeFailed, self).__init__()

class CheckVcodeIsTimeout(ThirdPluginError):
    """
    验证码过期
    """
    def __init__(self):
        super(CheckVcodeIsTimeout, self).__init__()

class CheckVcodeWrong(ThirdPluginError):
    """
    验证码错误
    """
    def __init__(self):
        super(CheckVcodeWrong, self).__init__()

class MFAOTPServerError(ThirdPluginError):
    """
    动态密码服务器异常
    """
    def __init__(self):
        super(MFAOTPServerError, self).__init__()

class OTPWrong(ThirdPluginError):
    """
    动态密码错误
    """
    def __init__(self):
        super(OTPWrong, self).__init__()

class OTPTimeout(ThirdPluginError):
    """
    动态密码过期
    """
    def __init__(self):
        super(OTPTimeout, self).__init__()

class OTPTooManyWrongTime(ThirdPluginError):
    """
    动态密码错误次数过多
    """
    def __init__(self):
        super(OTPTooManyWrongTime, self).__init__()