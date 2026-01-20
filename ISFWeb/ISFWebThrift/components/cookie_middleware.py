# coding:utf-8

from django.conf import settings
from components.cookie_util import convert_cookies_to_dict
from django.utils.deprecation import MiddlewareMixin

'''
cookie中包含中文，request.COOKIES.get获取不到sessionid
通过request.META.get('HTTP_COOKIE')获取cookie再重新赋值给request.COOKIES
'''
class CookieMiddleware(MiddlewareMixin):
    def process_request(self, request):
        ck = request.META.get('HTTP_COOKIE', None)
        if ck:
            ck_dict = convert_cookies_to_dict(ck)
            if settings.SESSION_COOKIE_NAME in ck_dict:
                request.COOKIES[settings.SESSION_COOKIE_NAME] = ck_dict[settings.SESSION_COOKIE_NAME]