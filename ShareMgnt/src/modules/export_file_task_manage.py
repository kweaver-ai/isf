#!/usr/bin/python3
# -*- coding:utf-8 -*-

import time
import os
import shutil
import threading
import uuid

from src.common import global_info
from src.common.sharemgnt_logger import ShareMgnt_Log
from src.common.lib import raise_exception
from src.common.business_date import BusinessDate
from ShareMgnt.ttypes import ncTShareMgntError


class TaskFinishedStatus:
    """
    任务完成状态
    """
    TASK_ERROR = -1
    TASK_IN_PROCESS = 0
    TASK_FINISHED = 1

class BaseTaskInfo(object):
    """
    @ todo: 生成任务结构体
    @ create_time: 任务创建时间
    @ file_path: 生成文件路径
    @ finished_status: 任务处理状态
    @ name: 文件名
    """
    def __init__(self, create_time=BusinessDate.time(), file_path="",
                    finished_status=TaskFinishedStatus.TASK_IN_PROCESS, name="",):
        self.create_time = create_time
        self.file_path = file_path
        self.finished_status = finished_status
        self.name = name
        self.lock = threading.Lock()

    def get_finished_status(self):
        with self.lock:
            return self.finished_status

    def set_finished_status(self, status):
        with self.lock:
            self.finished_status = status

    @classmethod
    def add_gen_file_task(cls, taskInfo):
        """
        添加任务
        返回任务 id
        """
        taskId = str(uuid.uuid1())
        with global_info.EXPORT_FILE_THREADLOCK:
            global_info.TASK_DICT[taskId] = taskInfo
        return taskId

    @classmethod
    def del_gen_file_task(cls, taskId, errMsg, errId):
        """
        删除任务
        """
        with global_info.EXPORT_FILE_THREADLOCK:
            # 任务不存在
            if taskId not in global_info.TASK_DICT:
                raise_exception(exp_msg=errMsg, exp_num=errId)
            else:
                del(global_info.TASK_DICT[taskId])

    @classmethod
    def get_gen_task_info(self, taskId, errMsg, errId):
        """
        获取任务信息
        """
        with global_info.EXPORT_FILE_THREADLOCK:
            # 任务不存在
            if taskId not in global_info.TASK_DICT:
                raise_exception(exp_msg=errMsg, exp_num=errId)
            else:
                return global_info.TASK_DICT[taskId]

    @classmethod
    def check_task_exist(self):
        """
        检查任务是否已存在，现在只支持单任务（导入导出的进度使用的全局变量，多任务时会显示出错，后续会优化）
        """
        with global_info.EXPORT_FILE_THREADLOCK:
            # 有任务正在进行，抛错
            for v in list(global_info.TASK_DICT.values()):
                if v.finished_status == TaskFinishedStatus.TASK_IN_PROCESS:
                    raise_exception(exp_msg=_("IDS_BATCH_USERS_EXPORTING"),
                                    exp_num=ncTShareMgntError.NCT_BATCH_USERS_EXPORTING)

class GenFileThread(threading.Thread):
    """
    生成文件线程
    """
    def __init__(self, taskId, taskInfo, gen_file_handler, delete_file_path):
        super(GenFileThread, self).__init__()
        self.taskId = taskId
        self.taskInfo = taskInfo
        self.gen_file_handler = gen_file_handler
        self.delete_file_path = delete_file_path

    def run(self):
        """
        执行
        """
        ShareMgnt_Log("**************** generate file thread start *************** *")

        try:
            self.taskInfo.file_path = self.gen_file_handler(self.taskId, self.taskInfo)
        except Exception as ex:
            # 任务处理异常，更新任务状态
            fileDir = os.path.join(self.delete_file_path, self.taskId)
            shutil.rmtree(fileDir)
            self.taskInfo.set_finished_status(TaskFinishedStatus.TASK_ERROR)
            ShareMgnt_Log("generate file thread run error: %s", str(ex))

        ShareMgnt_Log("**************** generate file thread end *******************")


class DeleteFileThread(threading.Thread):
    """
    清理文件线程
    """
    def __init__(self, interval = 3600, exist_time = 3600):
        """
        @ interval：线程的执行间隔
        @ file_path：需要清理的文件路径
        @ exist_time：需要清理文件、任务的存活时间
        @ thead_name：线程名
        """
        super(DeleteFileThread, self).__init__()
        self.interval = interval
        self.exist_time = exist_time

    def delete_overtime_task(self):
        """
        删除创建时间超过 self.exist_time 的任务
        """
        # 清理指定路径的文件
        for file_path in global_info.DELETE_FILE_PATHS:
            if os.path.exists(file_path):
                shutil.rmtree(file_path)

        with global_info.EXPORT_FILE_THREADLOCK:
            items = list(global_info.TASK_DICT.items())

        # 清理不在任务列表的文件
        for (taskId, taskInfo) in items:
            if taskInfo.create_time < (BusinessDate.time() - self.exist_time):
                with global_info.EXPORT_FILE_THREADLOCK:
                    if taskId in global_info.TASK_DICT:
                        # 删除文件夹
                        if os.path.exists(taskInfo.file_path):
                            shutil.rmtree(taskInfo.file_path)
                        # 删除任务
                        del(global_info.TASK_DICT[taskId])
                        ShareMgnt_Log("delete task: %s success.", taskId)

    def run(self):
        """
        执行
        """
        ShareMgnt_Log("**************** auto clear file path thread start *****************")

        while True:
            try:
                self.delete_overtime_task()
            except Exception as e:
                ShareMgnt_Log(" auto clear file path thread run error: %s", str(e))
            time.sleep(self.interval)

        ShareMgnt_Log("**************** auto clear file path thread end *****************")
