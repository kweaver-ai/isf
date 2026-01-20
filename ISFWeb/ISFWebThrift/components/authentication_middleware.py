# coding: utf-8

from django.http import HttpResponseForbidden, HttpResponseRedirect
from django.utils.deprecation import MiddlewareMixin
from components.common import verify
from components.getconfigfromfile import getHost

class AuthenMiddleware(MiddlewareMixin):
    """ 用户认证中间件 """

    # 以下请求不需要检查权限
    POST_WHITELIST = (
        '/api/ShareMgnt/OEM_GetConfigBySection',
        '/api/ShareMgnt/Usrm_CreateVcodeInfo',
        '/api/ShareMgnt/Usrm_GetVcodeConfig',
        '/api/ShareMgnt/Usrm_GetPasswordConfig',
    )

    GET_WHITELIST = (
        '/',
        '/api/deploy-manager/access-addr/app',
    )

    def process_request(self, request):
        if request.method == 'GET':
            if not (request.path in self.GET_WHITELIST):
                if not verify(request):
                    return HttpResponseForbidden()
        elif request.method == 'POST':
            if not request.path in self.POST_WHITELIST:
                if not verify(request):
                    return HttpResponseForbidden()