#!/usr/bin/python3
# -*- coding:utf-8 -*-
"""
ShareMgnt日志记录，主要用来记录调试信息
"""

import logging
import sys
import os
from logging.handlers import RotatingFileHandler

# logger名前缀
LOGNAME_PREFIX = "AnyShare_"

# 获取当前文件名
if hasattr(sys, 'frozen'):  # support for pyinstaller
    _srcfile = "src.common.sharemgnt_logger"
elif __file__[-4:].lower() in ['.pyc', '.pyo']:
    _srcfile = __file__[:-4] + '.py'
else:
    _srcfile = __file__
_srcfile = os.path.normcase(_srcfile)


def ShareMgnt_Log(msg, *args, **kwargs):
    """
    默认日志记录函数
    """
    ShareMgnt_Log2("sharemgnt.log", msg, *args, **kwargs)


def ShareMgnt_Log2(filename, msg, *args, **kwargs):
    """
    指定日志文件名记录函数
    """
    # 去掉扩展名，点号会干扰logging内部的分级处理
    logger = logging.getLogger("{0}{1}".format(LOGNAME_PREFIX, os.path.splitext(filename)[0]))
    logger.debug(msg, *args, **kwargs)


class ShareMgntLogger(logging.Logger):
    def __init__(self, name):
        logging.Logger.__init__(self, name)

        # 检查是不是自己的log
        if name.startswith(LOGNAME_PREFIX):
            self._isMine = True

            self.setLevel(logging.DEBUG)

            # 原来的配置handler逻辑
            formatter = logging.Formatter("%(asctime)s, Thread-%(thread)d:"
                                          "%(filename)s:%(lineno)s: %(message)s")

            stdout_handler = logging.StreamHandler(sys.stdout)
            stdout_handler.setFormatter(formatter)
            self.addHandler(stdout_handler)

logging.setLoggerClass(ShareMgntLogger)
