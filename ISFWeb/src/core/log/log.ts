import { noop, merge } from 'lodash';
import { currify } from '@/util/currify';
import session from '@/util/session';
import { privateAPI } from '../privateAPI/privateAPI';

const Config = {
    CSRFToken: noop,
}

/**
 * 设置Log
 */
export function setup(options) {
    merge(Config, options);
}

/**
 * 日志分类
 */
export enum LogCategory {
    Active = 1, // 活跃日志
    History = 2, // 历史日志
}

/**
 * 日志类型
 */
export enum LogType {
    NCT_LT_LOGIN = 10, // 登录日志
    NCT_LT_MANAGEMENT = 11, // 管理操作日志
    NCT_LT_OPEARTION = 12, // 操作日志
}

/**
 * 日至级别
 */
export enum Level {
    ALL, // 所有
    INFO, // 信息
    WARN,  // 警告
}

// 登录操作
export enum LoginOps {
    ALL = 0, // 所有操作
    LOGIN = 1, // 登录操作
    LOGOUT = 2, // 退出操作
    AUTHENICATION = 3, // 认证操作
    OTHER = 127, // 其它操作
}

// 管理操作
export enum ManagementOps {
    ALL = 0, // 所有操作
    CREATE = 1, // 新建操作
    ADD = 2, // 添加操作
    SET = 3, // 设置操作
    DELETE = 4, // 删除操作
    COPY = 5, // 复制
    MOVE = 6, // 移动
    REMOVE = 7, // 移除
    IMPORT = 8, // 导入操作
    EXPORT = 9, // 导出操作
    AUDIT_MGM = 10, // 审核操作
    QUARANTINE = 11, // 隔离
    UPLOAD = 12, // 上传
    PREVIEW = 13, // 预览
    DOWNLOAD = 14, // 下载
    RESTORE = 15, // 还原
    QUARANTINE_APPEAL = 16, // 隔离区申诉
    RESTART = 17, // 重启
    SEND_EMAIL = 18,  // 发送邮件
    RECOVER = 19,      // 恢复
    EDIT = 20,  // 编辑
    OTHER = 127, // 其他操作
}

// 文档操作
export enum OperationOps {
    ALL = 0, // 所有操作
    PREVIEW = 1, // 预览作用
    UPLOAD = 2, // 上传
    DOWNLOAD = 3, // 下载
    EDIT = 4, // 修改
    RENAME = 5, // 重名命
    DELETE = 6, // 删除操作
    COPY = 7, // 复制
    MOVE = 8, // 移动
    RESTORE_FROM_RECYCLE = 9, // 从回收站还原
    DELETE_FROM_RECYCLE = 10, // 彻底删除，
    PERM_MGM = 11, // 权限共享
    LINK_MGM = 12, // 外链共享
    FINDER_MGM = 13, // 发现共享
    BACKUP_BEGIN = 14, // 备份恢复
    LOCK_MGM = 16, // 文件锁
    ENTRY_DOC_MGM = 17, // 共享管理
    DEVICE_MGM = 18, // 登陆设备管理
    SET_CSF = 19,
    SYSREC_DELETE = 20, // 从系统回收站删除
    SYSREC_RESTORE = 21, // 从系统回收站还原
    CREATE_FOLDER = 22, // 新建文件夹
    SUBMIT_DOC_RELAY = 23, // 提交文档流转
    AUDIT_MGM = 24, // 审核管理
    DOC_RELAY = 25, // 文档流转
    NCT_DOT_FILECOLLECTOR = 26, // 文档收集
    NCT_DOT_CACHE = 27, // 缓存
    NCT_DOT_AUTOMATION = 28, // 自动化
    NCT_DOT_EXPORT = 29, // 导出
    DOC_DOMAIN_SYNC = 30, // 文档域同步
    NCT_DOT_SYNC = 31, // 同步流程管理
    NCT_DOT_FEEDBACK = 32, // 反馈
    NCT_DOT_KNOWLEDGE = 33, // 知识操作
    NCT_DOT_PRODUCT_PUBLISH = 34, // 数据产品发布
    NCT_DOT_PRODUCT_ARCHIVER = 35, // 数据产品归档
    RESTORE_REV = 50, // 还原版本
    OTHER = 127, // 其它
}

/**
 * 记录审计日志
 * @param logtype 日志类型
 * @param optype 操作类型
 * @param msg 日志信息
 * @param exmsg 附加信息
 * @param loglevel 日志级别
 * @param userid 用户id
 */
async function log(logtype: string, optype: ManagementOps, msg: string = '', exmsg: string = '', loglevel: Level = Level.INFO, userid: string = session.get('isf.userid')): Promise<any> {
    const params = { op_type: optype, msg, ex_msg: exmsg, level: loglevel, user_id: userid, user_type: 'authenticated_user' }

    try {
        return await privateAPI('audit-log', logtype, 'post', '', params)
    } catch {}
}

// 记录管理日志
export const manageLog = currify(log, 'log/management');
