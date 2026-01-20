import session from '@/util/session'
import __ from './locale';

export enum SystemAccountId {
    /**
     * 超级管理员和系统管理员
     */
    Admin = '266c6a42-6131-4d62-8f39-853e7093701c',

    /**
     * 安全管理员
     */
    Securit = '4bb41612-a040-11e6-887d-005056920bea',

    /**
     * 审计管理员
     */
    Audit = '94752844-BDD0-4B9E-8927-1CA8D427E699',
}

export enum SystemRoleType {
    /*
     * 超级管理员
     */
    Supper = '7dcfcc9c-ad02-11e8-aa06-000c29358ad6',

    /**
     * 系统管理员
     */
    Admin = 'd2bd2082-ad03-11e8-aa06-000c29358ad6',

    /**
     * 安全管理员
     */
    Securit = 'd8998f72-ad03-11e8-aa06-000c29358ad6',

    /**
     * 审计管理员
     */
    Audit = 'def246f2-ad03-11e8-aa06-000c29358ad6',

    /**
     * 组织管理员
     */
    OrgManager = 'e63e1c88-ad03-11e8-aa06-000c29358ad6',

    /**
     * 组织审计员
     */
    OrgAudit = 'f06ac18e-ad03-11e8-aa06-000c29358ad6',

    /**
     * 共享审核员
     */
    SharedApprove = 'f58622b2-ad03-11e8-aa06-000c29358ad6',

    /**
     * 文档审核员
     */
    DocApprove = 'fb648fac-ad03-11e8-aa06-000c29358ad6',

    /**
     * 定密审核员
     */
    CsfApprove = '01a78ac2-ad04-11e8-aa06-000c29358ad6',

}

/**
 * 用户角色
 */
export enum UserRole {
    /**
     * 超级管理员
     */
    Super = 'super_admin',

    /**
     * 系统管理员
     */
    Admin = 'sys_admin',

    /**
     * 安全管理员
     */
    Security = 'sec_admin',

    /**
     * 审计管理员
     */
    Audit = 'audit_admin',

    /**
     * 组织管理员
     */
    OrgManager = 'org_manager',

    /**
     * 组织审计员
     */
    OrgAudit = 'org_audit',

    /**
     * 普通用户
     */
    NormalUser = 'normal_user',
}

/**
 * 系统角色id与用户角色的映射
 */
export const SysUserRoles = {
    [SystemRoleType.Supper]: UserRole.Super,
    [SystemRoleType.Admin]: UserRole.Admin,
    [SystemRoleType.Securit]: UserRole.Security,
    [SystemRoleType.Audit]: UserRole.Audit,
    [SystemRoleType.OrgManager]: UserRole.OrgManager,
    [SystemRoleType.OrgAudit]: UserRole.OrgAudit,
}

/**
 * 获取角色名
 */
export const getRoleName = (role: any): string => {
    switch (role.id) {
        case SystemRoleType.Supper:
            return __('超级管理员')
        case SystemRoleType.Admin:
            return __('系统管理员')
        case SystemRoleType.Securit:
            return __('安全管理员')
        case SystemRoleType.Audit:
            return __('审计管理员')
        case SystemRoleType.OrgManager:
            return __('组织管理员');
        case SystemRoleType.OrgAudit:
            return __('组织审计员');
        default:
            return role.name || '---'
    }
}

/**
 * 根据用户角色获取角色名
 */
export const getNameByRole = (role: UserRole): string => {
    switch (role) {
        case UserRole.Super:
            return __('超级管理员')
        case UserRole.Admin:
            return __('系统管理员')
        case UserRole.Security:
            return __('安全管理员')
        case UserRole.Audit:
            return __('审计管理员')
        case UserRole.OrgManager:
            return __('组织管理员');
        case UserRole.OrgAudit:
            return __('组织审计员');
        case UserRole.NormalUser:
            return __('普通用户');
        default:
            return '---'
    }
}

/**
 * 获取角色类型
 */
export const getRoleType = (): UserRole => {
    const { user: { roles } } = session.get('isf.userInfo')

    const roleIds = roles.map((role) => role.id)

    switch (true) {
        case roleIds.includes(SystemRoleType.Supper):
            return UserRole.Super
        case roleIds.includes(SystemRoleType.Admin):
            return UserRole.Admin
        case roleIds.includes(SystemRoleType.Audit):
            return UserRole.Audit
        case roleIds.includes(SystemRoleType.Securit):
            return UserRole.Security
        case roleIds.includes(SystemRoleType.OrgManager):
            return UserRole.OrgManager
        case roleIds.includes(SystemRoleType.OrgAudit):
            return UserRole.OrgAudit
        default:
            return UserRole.NormalUser
    }
}

/**
 * 获取角色职能
 */
export const getRoleFuntional = (role: any): string => {

    switch (role.id) {
        case SystemRoleType.Supper:
            return __('组织管理、文档管理、运营监控、系统维护和安全管控，审计所有用户（包括自己）的行为、管理所有用户角色')
        case SystemRoleType.Admin:
            return __('组织管理、文档管理、运营监控和系统维护，管理组织管理员、各个审核员以及自定义角色')
        case SystemRoleType.Securit:
            return __('安全管控，审计审计管理员、组织管理员、组织审计员及各个审核员和普通用户的行为以及管理组织审计员角色')
        case SystemRoleType.Audit:
            return __('审计系统管理员和安全管理员的行为')
        case SystemRoleType.OrgManager:
            return __('管理用户组织和文档组织，管理自定义角色，在管辖范围内管理组织管理员和各个审核员角色');
        case SystemRoleType.OrgAudit:
            return __('审计管辖范围内用户的行为，在管辖范围内管理组织审计员角色');
        default:
            return role.description || '---'
    }
}

/**
 * 审核员审核对象类型
 */
export enum AuditObjectType {
    /**
     * 用户
     */
    AuditObjectUser = 1,

    /**
     * 部门
     */
    AuditObjectDept = 2,

    /**
     * 文档库
     */
    AuditObjectCustom = 3,
}

export interface ManagerType {
    is_knowledge_manager: boolean;
    is_other_manager?: boolean;
    managerType?: UserRole;
}
/**
 * 用户角色
 */
export function getRole() {
    const user_info = session.get('isf.userInfo')
    const roles = user_info.user.roles
    const is_knowledge_manager = user_info.user.custom_attr.is_knowledge
    const managerType = getRoleType()
    if(!roles.length && is_knowledge_manager) {
        return { is_knowledge_manager: true, is_other_manager: false, managerType }
    } else {
        return { is_knowledge_manager, is_other_manager: true, managerType }
    }
}

/**
 * 获取所有角色类型
 */
export const getRoleTypes = (): UserRole[] => {
    const { user: { roles } } = session.get('isf.userInfo')

    const roleIds = roles.map((role) => role.id)

    const userRoles = roleIds.map((role) => SysUserRoles[role])

    return userRoles.length ? userRoles : [UserRole.NormalUser]
}