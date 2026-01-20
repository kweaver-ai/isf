import { ErrorCode } from '../apis/openapiconsole/errorcode';
import { ErrorCode as Code } from '@/core/thrift/sharemgnt/errcode';
import { DocLibType } from '../doclibs/doclibs'
import __ from './locale';

// 异常码对应中文描述，用于在国际化资源中查找对应信息
// 国际化资源需要在locale中定义
const ErrorKeyHash = {
    [ErrorCode.ResourceInaccessible]: '该资源已不存在。',
    [ErrorCode.ResourceConflict]: '该资源已存在。',
    [ErrorCode.DomainLinkFalied]: '添加失败，指定的文档域无法连接。',
    [ErrorCode.PolicyBoundDomain]: '无法删除，此策略配置与文档域有绑定关系，请先解除文档域绑定，再执行此操作。',
    [ErrorCode.DocNotExist]: '该文档库已不存在。',
    [ErrorCode.DocNameSameAsUserName]: '库名称不能和用户显示名重名。',
    [ErrorCode.DocNameConflict]: (docType: DocLibType) => {
        switch (docType) {
            case DocLibType.UserDocLib:
                return '已存在同名的个人文档库。'

            case DocLibType.DepDocLib:
                return '已存在同名的部门文档库。'

            case DocLibType.CustomDocLib:
                return '已存在同名的自定义文档库。'

            case DocLibType.KnowledgeDocLib:
                return '已存在同名的知识库。'
        }
    },
    [ErrorCode.QuotaGreaterThanAvailable]: '当前文档管理剩余可分配空间为${availableQuota}。',
    [ErrorCode.QuotaLessThanUsed]: '配额空间不能小于当前已使用空间。',
    [ErrorCode.StorageNotExist]: '所指定的存储位置已不可用，请更换。',
    [ErrorCode.StorageUnAvailable]: '所指定的存储位置已不可用，请更换。',
    [Code.OSSNotExist]: '所指定的存储位置已不可用，请更换。',
    [Code.OSSDisabled]: '所指定的存储位置已不可用，请更换。',
    [Code.OSSInvalid]: '所指定的存储位置已不可用，请更换。',
    [Code.OSSUnabled]: '所指定的存储位置已不可用，请更换。',
    [ErrorCode.UserGroupNotExist]: '用户组“${groupName}”已不存在。',
    [ErrorCode.UserNotExist]: '您的选择包含不存在的用户',
    [ErrorCode.DepartmentsNotExist]: '您的选择包含不存在的部门',
    [ErrorCode.UserGruopsNotExist]: '您的选择包含不存在的用户组',
    [ErrorCode.InternalError]: '内部错误',
    [ErrorCode.InvalidRequest]: '参数不合法',
}

/**
 * 通过errcode查找资源文件的Key
 * @param errcode 异常码
 * @param optype 操作类型
 * @return 返回异常信息的Key
 */
function findLocale(errcode, ...params): string {
    const match = ErrorKeyHash[errcode];

    if (!match) {
        return '内部错误';
    } else if (typeof match === 'string') {
        return match;
    } else if (typeof match === 'function') {
        return match(...params);
    }
}

/**
 * 获取异常提示模版
 * @param errcode 异常码
 * @param optype 操作类型
 * @return 返回异常提示模版函数
 */
export function getErrorTemplate(errcode: number, ...params): (Object) => string {
    return function (args) {
        return __(findLocale(errcode, ...params), args);
    }
}

/**
 * 根据异常码获取异常提示
 * @param errcode 异常码
 * @param args 模版填充信息
 * @return 返回异常提示信息
 */
export function getErrorMessage(errcode: number, ...params): string {
    return __(findLocale(errcode, ...params)) || '';
}