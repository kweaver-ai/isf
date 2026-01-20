#-*- coding:utf-8 -*-

''' 获取服务端IP '''
import os
from django.conf import settings

def get_host():
    """
    获取服务端IP
    """

    if settings.DEBUG and hasattr(settings, 'DEBUG_HOST'):
        return settings.DEBUG_HOST
    else:
        return os.getenv('HOST_IP')

