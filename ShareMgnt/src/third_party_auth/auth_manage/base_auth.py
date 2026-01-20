# coding=utf8
"""第三方授权接口"""

class BaseAuth(object):
    """第三方授权接口"""
    def __init__(self, config):
        pass

    def login(self, user_name, pwd, db_user):
        """账号密码登录接口"""
        raise NotImplementedError("Method not implement")

    def validate(self, params):
        """单点登录验证接口"""
        raise NotImplementedError("Method not implement")
