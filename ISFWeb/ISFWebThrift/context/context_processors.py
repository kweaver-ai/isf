# coding=utf-8
from django.conf import settings


def view(request):

    if request.path == '/login/' or request.path == '/oauth/':
        return {}
    else:
       return {}

def user(request):
    userinfo = request.session.get('userinfo', {})
    return userinfo

def scripts(request):
    SCRIPTS_ROOT = settings.DEBUG and 'src' or 'build'
    return {'SCRIPTS_ROOT':SCRIPTS_ROOT}