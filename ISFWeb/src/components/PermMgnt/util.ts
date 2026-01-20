import { SystemRoleType, UserRole } from "@/core/role/role";
import { InnerRoleIdEnum } from "./types";
import { PermissionType, PolicyConfigListType } from "@/core/apis/console/authorization/type";
import intl from "react-intl-universal";
import dayjs from "dayjs";

export enum AccessorTypeEnum {
    User = 'user',
    Department = 'department',
    Group = 'group',
    App = 'app',
}

/**
 * 永久有效时间
 * @description 用于表示永久有效时间的字符串，格式为 ISO 8601 标准。
 * @example '1970-01-01T08:00:00+08:00'
 */
export const foreverExpire = '1970-01-01T08:00:00+08:00';

/**
 * 格式化权限
 * @param cur 当前权限
 * @param permConfig 权限配置
 * @returns 权限字符串`
 */
export const formatPerm = (cur: {operation: PermissionType}, permConfig: PolicyConfigListType[]) => {
    const {
        operation: { allow, deny },
    } = cur;
    let permString: string = '';
    if (allow.length) {
        permString = allow.map(item => item.name).join('/');
    }
  
    if (deny.length) {
        const denyString = deny.map(item => item.name).join('/');
        permString = `${permString}${deny.length >= permConfig.length || !permString ? intl.get('access.refused') : `(${intl.get('deny')} ${denyString})`}`;
    }
    return permString;
};

/**
 * 格式化过期时间
 * @param expires 过期时间
 * @returns 格式化后的过期时间
 */
export const formatExpires = (expires: string) => {
    return expires === foreverExpire
        ? intl.get('forever.expire')
        : dayjs(expires).format('YYYY/MM/DD HH:mm');
};

/**
 * 格式化当前项
 * @param item 当前项
 * @param needName 是否需要名称
 * @returns 格式化后的当前项
 */
export const fromatItem = (item: any, needName = false) => {
    switch (item.type) {
        case AccessorTypeEnum.User:
            return needName ? { id: item.id, name: (item.user && item.user.displayName) || item.name || '', type: item.type }: { id: item.id, type: item.type };
        case AccessorTypeEnum.Department:
            return needName ?{ id: item.id, name: item.name, type: item.type }: { id: item.id, type: item.type };
        case AccessorTypeEnum.Group:
            return needName ?{ id: item.id, name: item.name, type: item.type }: { id: item.id, type: item.type };
        case AccessorTypeEnum.App:
            return needName ?{ id: item.id, name: item.name, type: item.type }: { id: item.id, type: item.type };
        default:
            return item;
    }
};

export const getRoleType = (roles): UserRole => {
    const roleIds = roles.map((role) => role.id)

    switch (true) {
        case roleIds.includes(SystemRoleType.Supper):
            return UserRole.Super
        case roleIds.includes(SystemRoleType.Admin):
            return UserRole.Admin
        case roleIds.includes(SystemRoleType.Securit):
            return UserRole.Security
        case roleIds.includes(SystemRoleType.Audit):
            return UserRole.Audit
        case roleIds.includes(SystemRoleType.OrgManager):
            return UserRole.OrgManager
        case roleIds.includes(SystemRoleType.OrgAudit):
            return UserRole.OrgAudit
        default:
            return UserRole.Super
    }
}

/**
 * 格式化请求
 * @param instance_url 请求url
 * @returns 格式化后的请求信息
 */
export const formatRequest = (instance_url) => {
    const [method, url] = instance_url.split(' ');
    const [path, queryString] = url?.split('?') || ['', ''];
    //通过url的query中是否有keyword参数，来判断是否支持搜索，使用时需要过滤掉keyword参数
    let filteredQueryString = '';
    if (queryString) {
        const params = queryString.split('&');
        const filteredParams = params.filter(param => {
            const [key] = param.split('=');
            return key !== 'keyword';
        });
        
        if (filteredParams.length > 0) {
            filteredQueryString = '?' + filteredParams.join('&');
        }
    }
    
    const processedUrl = path + filteredQueryString;
    const urlWithoutApi = processedUrl.replace('/api', '');
    const urlParams = urlWithoutApi.split('/').filter(part => part);

    return {method: method.toLowerCase(), urlParams}
}