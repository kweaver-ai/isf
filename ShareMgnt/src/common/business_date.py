#!/usr/bin/python3
# -*- coding:utf-8 -*-

from datetime import datetime
from datetime import date
import time
import threading
import configparser

class BusinessDate:
    _instance_lock = threading.Lock()

    @classmethod
    def get_time_offset(cls):
        '''
        从配置文件得到时间偏移量，businessdate._timeoffset单位为秒
        '''
        if not hasattr(BusinessDate, "_timeoffset"):
            with BusinessDate._instance_lock:
                if not hasattr(BusinessDate, "_timeoffset"):
                    try:
                        cf = configparser.ConfigParser()
                        cf.read("/sysvol/conf/service_conf/app_default.conf")
                        BusinessDate._timeoffset = cf.getint("Global","business_time_offset")
                    except Exception:
                        BusinessDate._timeoffset = 0
        return BusinessDate._timeoffset

    @classmethod
    def time(cls) :
        '''
        对time.time()函数的封装
        '''
        return time.time() + cls.get_time_offset()

    @classmethod
    def now(cls,tz=None):
        '''
        对datetime.datetime.now()函数的封装
        '''
        return datetime.fromtimestamp(cls.time(), tz)

    @classmethod
    def today(cls):
        '''
        对datetime.date.today()函数的封装
        '''
        return date.fromtimestamp(cls.time())
