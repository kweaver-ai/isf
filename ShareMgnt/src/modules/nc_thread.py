#!/usr/bin/python3
# -*- coding:utf-8 -*-
"""通知中心运行线程"""
import threading
from src.common.sharemgnt_logger import ShareMgnt_Log
from src.common.nc_senders import email_send
from src.modules.smtp_manage import (JsonConfManage, MailRecipient)
from src.common.jsonconv_Ttype import (SmtpConfEnc, SmtpConfDec)
from ShareMgnt.ttypes import (ncTSmtpSrvConf,
                              ncTAlarmConfig)
import queue


class NotifiCenterThread(threading.Thread):

    def __init__(self):
        """
        初始化
        """
        super(NotifiCenterThread, self).__init__()
        self.smtpConf = None
        self.alertConf = ncTAlarmConfig(0, 0, [])
        self.jobs = queue.Queue()

        self.__notifies = {
            "logsave_timeout": ["warn", _("log validity alarm"), _("log validity alarm content")],
            "logsave_full": ["warn", _("log space alarm"), _("log space alarm content")]
        }

    _instance_lock = threading.Lock()

    @staticmethod
    def instance():
        """
        这个类单例化
        """
        if not hasattr(NotifiCenterThread, "_instance"):
            with NotifiCenterThread._instance_lock:
                if not hasattr(NotifiCenterThread, "_instance"):
                    # 一定要2次检验
                    NotifiCenterThread._instance = NotifiCenterThread()
        return NotifiCenterThread._instance

    def close(self):
        """
        关闭
        """
        self.jobs.put("close")

    def add(self, task):
        """
        添加任务
        """
        self.jobs.put(task)

    def run(self):
        """
        执行
        """
        ShareMgnt_Log("**************** 通知中心线程开始 *****************")

        # 先加载配置
        try:
            smtp_manage = JsonConfManage("smtp_config", SmtpConfEnc, SmtpConfDec)
            self.smtpConf = smtp_manage.get_config()

        except Exception as e:
            ShareMgnt_Log("通知中心载入异常，采用默认配置：%s", str(e))

        while True:
            try:
                task = self.jobs.get()
                if isinstance(task, str):
                    if task == "close":
                        break
                elif isinstance(task, ncTSmtpSrvConf):
                    self.smtpConf = task

                elif isinstance(task, tuple):
                    taskid = task[0]
                    ShareMgnt_Log("通知中心接收到：%s", taskid)
                    # taskp = task[1]
                    if taskid in self.__notifies:
                        if self.smtpConf is not None:
                            # 获取邮箱收件人列表
                            email_list = MailRecipient("smtp_Recipient_config").get_config()
                            if email_list:
                                email_send(self.smtpConf,
                                           email_list,
                                           self.__notifies[taskid][1],
                                           self.__notifies[taskid][2])

            except Exception as e:
                ShareMgnt_Log("通知中心运行异常：%s", str(e))

        ShareMgnt_Log("**************** 通知中心线程结束 *****************")
