import { difference } from 'lodash'
import { getConfidentialConfig } from '@/core/apis/eachttp/config/config'
import { SystemRoleType } from '@/core/role/role'
import { transformAdmins } from '@/core/confidential/confidential'
import session from '@/util/session'
import * as depMgmImg from './assets/images/depMgm.png'
import * as userMgmImg from './assets/images/userMgm.png'
import * as exportMgmImg from './assets/images/exportMgm.png'
import * as FreezeUserImg from './assets/images/freezeUser.png'
import __ from './locale'

export enum TabEnum {
    User = "user",
    Department = "department",
    UserGroup = "user-group",
    AppAccount = "app-account"
}

export enum Range {
    DEPARTMENT_DEEP, // 部门及其子部门

    DEPARTMENT, // 当前部门

    USERS, // 当前选中的用户
}

/**
 * 特殊分组
 */
export enum SpecialDep {
    /**
     * 未分配组
     */
    Unassigned = '-1',

    /**
     * 所有用户
     */
    AllUsers = '-2',
}

/**
 * 用户状态
 */
export enum EnableStatus {
    /**
     * 已启用
     */
    Enabled,

    /**
     * 已禁用
     */
    Disabled,
}

export enum Action {
    /**
     * 无操作
     */
    None,

    /**
     * 分割线
     */
    Line,

    /**
     * 新建组织
     */
    CreateOrg,

    /**
     * 编辑组织
     */
    EditOrg,

    /**
     * 删除组织
     */
    DelOrg,

    /**
     * 新建部门
     */
    CreateDep,

    /**
     * 编辑部门
     */
    EditDep,

    /**
     * 删除部门
     */
    DelDep,

    /**
     * 移动部门
     */
    MoveDep,

    /**
     * 添加用户至部门
     */
    AddUsersToDep,

    /**
     * 新建用户
     */
    CreateUser,

    /**
     * 编辑用户
     */
    EditUser,

    /**
     * 删除用户
     */
    DelUser,

    /**
     * 设置用户有效期
     */
    SetExpiration,

    /**
     * 移动用户
     */
    MoveUser,

    /**
     * 移除用户
     */
    RemoveUser,

    /**
     * 启用用户
     */
    EnableUser,

    /**
     * 禁用用户
     */
    DisableUser,

    /**
     * 设置角色
     */
    SetRole,

    /**
     * 产品授权概览
     */
    ProductLicense,

    /**
     * 交接工作
     */
    WorkHandover,

    /**
     * 管控密码
     */
    ManagePwd,

    /**
     * 导入导出用户组织
     */
    ImportOrg,

    /**
     * 导入域用户组织
     */
    ImportDomain,

    /**
     * 导入第三方用户组织
     */
    ImportThirdOrg,

    /**
     * 冻结用户
     */
    FreezeUser,

    /**
     * 解冻用户
     */
    UnfreezeUser,
}

/**
 * 菜单
 */
export interface Menu {
    /**
     * 操作类型
     */
    type: Action;

    /**
     * 操作名称
     */
    text?: string;

    /**
     * 是否禁用
     */
    disabled?: boolean;
}

/**
 * 菜单分组名称
 */
export enum MenuGroupName {
    /**
     * 部门管理
     */
    DepMgnt = 'depMgnt',

    /**
     * 用户管理
     */
    UserMgnt = 'userMgnt',

    /**
     * 导入导出用户组织
     */
    ImportOrg = 'importOrg',

    /**
     * 用户冻结管理
     */
    FreezeUser = 'freezeUser',
}

/**
 * 菜单栏
 */
export interface MenuGroup {
    /**
     * 菜单分组名称
     */
    name: MenuGroupName;

    /**
     * 显示名称
     */
    text: string;

    /**
     * 图标
     */
    fallback: string;

    /**
     * 具体菜单操作
     */
    actions: ReadonlyArray<Menu>;
}

/**
 * 部门管理操作
 */
export const DepActions = {
    [Action.Line]: {
        type: Action.Line,
    },
    [Action.CreateOrg]: {
        type: Action.CreateOrg,
        text: __('新建组织'),
        disabled: false,
    },
    [Action.EditOrg]: {
        type: Action.EditOrg,
        text: __('编辑组织'),
        disabled: false,
    },
    [Action.DelOrg]: {
        type: Action.DelOrg,
        text: __('删除组织'),
        disabled: false,
    },
    [Action.CreateDep]: {
        type: Action.CreateDep,
        text: __('新建部门'),
        disabled: false,
    },
    [Action.EditDep]: {
        type: Action.EditDep,
        text: __('编辑部门'),
        disabled: false,
    },
    [Action.DelDep]: {
        type: Action.DelDep,
        text: __('删除部门'),
        disabled: false,
    },
    [Action.MoveDep]: {
        type: Action.MoveDep,
        text: __('移动部门至'),
        disabled: false,
    },
    [Action.AddUsersToDep]: {
        type: Action.AddUsersToDep,
        text: __('添加用户至部门'),
        disabled: false,
    },
}

/**
 * 用户管理操作
 */
export const UserActions = {
    [Action.Line]: {
        type: Action.Line,
    },
    [Action.CreateUser]: {
        type: Action.CreateUser,
        text: __('新建用户'),
        disabled: false,
    },
    [Action.EditUser]: {
        type: Action.EditUser,
        text: __('编辑用户'),
        disabled: false,
    },
    [Action.DelUser]: {
        type: Action.DelUser,
        text: __('删除用户'),
        disabled: false,
    },
    [Action.SetExpiration]: {
        type: Action.SetExpiration,
        text: __('设置用户有效期'),
        disabled: false,
    },
    [Action.MoveUser]: {
        type: Action.MoveUser,
        text: __('移动用户至'),
        disabled: false,
    },
    [Action.RemoveUser]: {
        type: Action.RemoveUser,
        text: __('移除用户'),
        disabled: false,
    },
    [Action.EnableUser]: {
        type: Action.EnableUser,
        text: __('启用用户'),
        disabled: false,
    },
    [Action.DisableUser]: {
        type: Action.DisableUser,
        text: __('禁用用户'),
        disabled: false,
    },
    [Action.SetRole]: {
        type: Action.SetRole,
        text: __('设置系统角色'),
        disabled: false,
    },
    [Action.ProductLicense]: {
        type: Action.ProductLicense,
        text: __('产品授权'),
        disabled: false,
    },
    // [Action.WorkHandover]: {
    //     type: Action.WorkHandover,
    //     text: __('交接工作'),
    //     disabled: false,
    // },
    [Action.ManagePwd]: {
        type: Action.ManagePwd,
        text: __('管控密码'),
        disabled: false,
    },
}

/**
 * 导入用户组织操作
 */
export const ImportActions = {
    [Action.Line]: {
        type: Action.Line,
    },
    [Action.ImportOrg]: {
        type: Action.ImportOrg,
        text: __('导出导入用户组织'),
        disabled: false,
    },
    [Action.ImportDomain]: {
        type: Action.ImportDomain,
        text: __('导入域用户组织'),
        disabled: false,
    },
    [Action.ImportThirdOrg]: {
        type: Action.ImportThirdOrg,
        text: __('导入第三方用户组织'),
        disabled: false,
    },
}

/**
 * 冻结用户管理
 */
export const FreezeUserActions = {
    [Action.Line]: {
        type: Action.Line,
    },
    [Action.FreezeUser]: {
        type: Action.FreezeUser,
        text: __('冻结用户'),
        disabled: false,
    },
    [Action.UnfreezeUser]: {
        type: Action.UnfreezeUser,
        text: __('解冻用户'),
        disabled: false,
    },
}

/**
 * 菜单分组
 */
export const MenuGroups = {
    [MenuGroupName.DepMgnt]: {
        name: MenuGroupName.DepMgnt,
        text: __('部门管理'),
        fallback: depMgmImg,
        actions: [],
    },
    [MenuGroupName.UserMgnt]: {
        name: MenuGroupName.UserMgnt,
        text: __('用户管理'),
        fallback: userMgmImg,
        actions: [],
    },
    [MenuGroupName.ImportOrg]: {
        name: MenuGroupName.ImportOrg,
        text: __('导入用户组织'),
        fallback: exportMgmImg,
        actions: [],
    },
    [MenuGroupName.FreezeUser]: {
        name: MenuGroupName.FreezeUser,
        text: __('用户冻结管理'),
        fallback: FreezeUserImg,
        actions: [],
    },
}

/**
 * 用户角色列表
 */
const userRole = [SystemRoleType.OrgManager, SystemRoleType.OrgAudit, SystemRoleType.PortalManager]

/**
 * 获取是否显示设置角色菜单
 * @param triSystemStatus 是否开启三权分立
 */
export const getIsShowSetRole = async (triSystemStatus: boolean): Promise<boolean | SystemRoleType> => {
    /**
     * 当前用户不存在未被屏蔽的角色，则屏蔽[设置角色]菜单
     */
    try {
        const disabledRoles = await getConfidentialConfig('disabled_roles')
        const disabledRoleIds = transformAdmins(disabledRoles)
        /**
         * 通过difference(a, b).length === 0判断a是否为b的子集
         * 1、开启三权分立，并且屏蔽了组织管理员、组织审计员，则屏蔽[设置角色]菜单
         * 2、屏蔽了用户角色列表所有的可见角色，则屏蔽[设置角色]菜单
         */
        if (triSystemStatus && difference(userRole, disabledRoleIds).length === 0) {
            return false
        }
        if (difference([...userRole, SystemRoleType.Supper], disabledRoleIds).length === 0) {
            return false
        }
        if (!disabledRoleIds.includes(SystemRoleType.OrgManager)) {
            return SystemRoleType.OrgManager
        }
        return true
    } catch {
        return false
    }
}

/**
 * 获取是否显示管控密码菜单
 * @param triSystemStatus 是否开启三权分立
 */
export const getIsShowManagePwd = async (): Promise<boolean> => {
    // 当前登录用户的角色id
    const { user: { roles } } = session.get('isf.userInfo')
    const roleIds = roles.map((role) => role.id)

    // 获取屏蔽角色信息，并转换为id
    let disabledPwdAdmins
    try {
        disabledPwdAdmins = await getConfidentialConfig('disabled_pwd_controllers_admins')
    } catch {
        disabledPwdAdmins = ['sys_admin', 'audit_admin', 'org_audit']
    }
    const disabledPwdAdminIds = transformAdmins(disabledPwdAdmins)

    /**
     * 通过difference(a, b).length === 0判断a是否为b的子集
     * 1、开启三权分立，并且屏蔽了用户角色列表所有的可见角色，则屏蔽[管控密码]菜单
     * 2、屏蔽了用户角色列表所有的可见角色，则屏蔽[管控密码]菜单
     */
    if (difference(roleIds, disabledPwdAdminIds).length === 0) {
        return false
    }
    return true
}

/**
 * 获取是否显示启用/禁用用户菜单
 */
export const getIsShowEnableAndDisableUser = async (): Promise<boolean> => {
    // 当前登录用户的角色id
    const { user: { roles } } = session.get('isf.userInfo')
    const roleIds = roles.map((role) => role.id)

    // 获取屏蔽角色信息，并转换为id
    let disabledEADUserAdmins
    try {
        disabledEADUserAdmins = await getConfidentialConfig('disabled_account_control_admins')
    } catch {
        disabledEADUserAdmins = ['sys_admin']
    }

    const disabledEADUserAdminsIds = transformAdmins(disabledEADUserAdmins)

    /**
     * 通过difference(a, b).length === 0判断a是否为b的子集
     * 1、开启三权分立，并且屏蔽了用户角色列表所有的可见角色，则屏蔽[启用/禁用用户]菜单
     * 2、屏蔽了用户角色列表所有的可见角色，则屏蔽[启用/禁用用户]菜单
     */
    if (difference(roleIds, disabledEADUserAdminsIds).length === 0) {
        return false
    }
    return true
}

/**
 * 根据角色获取菜单
 * @param roles 角色信息
 * @param triSystemStatus 是否开启三权分立
 * @param freezeStatus 是否开启冻结用户功能
 * @param enableThirdImport 是否开启导入第三方用户组织
 * @param isShowManagePwd 是否显示密码管控菜单
 * @param isShowSetRole 是否显示设置角色菜单
 * @param isShowEnableAndDisableUser 是否显示启用/禁用用户菜单
 * @return 菜单
 */
export function getMenuGroupsByRole({ roles, triSystemStatus, freezeStatus, enableThirdImport, isShowManagePwd, isShowSetRole, isShowEnableAndDisableUser }) {
    const roleIds = roles.map((role) => role.id)

    switch (true) {
        case roleIds.includes(SystemRoleType.Supper):
            return [
                {
                    ...MenuGroups[MenuGroupName.DepMgnt],
                    actions: [
                        DepActions[Action.CreateOrg],
                        DepActions[Action.EditOrg],
                        DepActions[Action.DelOrg],
                        DepActions[Action.Line],
                        DepActions[Action.CreateDep],
                        DepActions[Action.EditDep],
                        DepActions[Action.DelDep],
                        DepActions[Action.Line],
                        DepActions[Action.MoveDep],
                        DepActions[Action.AddUsersToDep],
                    ],
                },
                {
                    ...MenuGroups[MenuGroupName.UserMgnt],
                    actions: [
                        UserActions[Action.CreateUser],
                        UserActions[Action.EditUser],
                        UserActions[Action.DelUser],
                        UserActions[Action.Line],
                        // UserActions[Action.SetExpiration],
                        UserActions[Action.MoveUser],
                        UserActions[Action.RemoveUser],
                        isShowEnableAndDisableUser ? UserActions[Action.EnableUser] : undefined,
                        isShowEnableAndDisableUser ? UserActions[Action.DisableUser] : undefined,
                        UserActions[Action.Line],
                        ...isShowSetRole ? [UserActions[Action.SetRole]] : [],
                        isShowSetRole ? UserActions[Action.Line] : undefined,
                        UserActions[Action.ProductLicense],
                        isShowManagePwd ? UserActions[Action.Line]: undefined,
                        isShowManagePwd ? UserActions[Action.ManagePwd] : undefined,
                    ].filter((action) => action !== undefined),
                },
                {
                    ...MenuGroups[MenuGroupName.ImportOrg],
                    actions: [
                        ImportActions[Action.ImportOrg],
                        ImportActions[Action.ImportDomain],
                        enableThirdImport ? ImportActions[Action.ImportThirdOrg] : undefined,
                    ].filter(((action) => action !== undefined)),
                },
                freezeStatus ? {
                    ...MenuGroups[MenuGroupName.FreezeUser],
                    actions: [
                        FreezeUserActions[Action.FreezeUser],
                        FreezeUserActions[Action.UnfreezeUser],
                    ],
                } : undefined,
            ].filter((group) => group !== undefined)

        case roleIds.includes(SystemRoleType.Admin):
            return [
                {
                    ...MenuGroups[MenuGroupName.DepMgnt],
                    actions: [
                        DepActions[Action.CreateOrg],
                        DepActions[Action.EditOrg],
                        DepActions[Action.DelOrg],
                        DepActions[Action.Line],
                        DepActions[Action.CreateDep],
                        DepActions[Action.EditDep],
                        DepActions[Action.DelDep],
                        DepActions[Action.Line],
                        DepActions[Action.MoveDep],
                        DepActions[Action.AddUsersToDep],
                    ],
                },
                {
                    ...MenuGroups[MenuGroupName.UserMgnt],
                    actions: [
                        UserActions[Action.CreateUser],
                        UserActions[Action.EditUser],
                        UserActions[Action.DelUser],
                        UserActions[Action.Line],
                        // UserActions[Action.SetExpiration],
                        UserActions[Action.MoveUser],
                        UserActions[Action.RemoveUser],
                        isShowEnableAndDisableUser ? UserActions[Action.EnableUser] : undefined,
                        isShowEnableAndDisableUser ? UserActions[Action.DisableUser] : undefined,
                        UserActions[Action.Line],
                        UserActions[Action.ProductLicense],
                        isShowManagePwd ? UserActions[Action.Line]: undefined,
                        isShowManagePwd ? UserActions[Action.ManagePwd] : undefined,
                    ].filter(((action) => action !== undefined)),
                },
                {
                    ...MenuGroups[MenuGroupName.ImportOrg],
                    actions: [
                        ImportActions[Action.ImportOrg],
                        ImportActions[Action.ImportDomain],
                        enableThirdImport ? ImportActions[Action.ImportThirdOrg] : undefined,
                    ].filter((action) => action !== undefined),
                },
                freezeStatus ? {
                    ...MenuGroups[MenuGroupName.FreezeUser],
                    actions: [
                        FreezeUserActions[Action.FreezeUser],
                        FreezeUserActions[Action.UnfreezeUser],
                    ],
                } : undefined,
            ].filter((group) => group !== undefined)

        case roleIds.includes(SystemRoleType.Securit):
            return [
                {
                    ...MenuGroups[MenuGroupName.UserMgnt],
                    actions: [
                        UserActions[Action.EditUser],
                        UserActions[Action.Line],
                        isShowEnableAndDisableUser ? UserActions[Action.EnableUser] : undefined,
                        isShowEnableAndDisableUser ? UserActions[Action.DisableUser] : undefined,
                        UserActions[Action.Line],
                        ...isShowSetRole ? [UserActions[Action.SetRole]] : [],
                        UserActions[Action.WorkHandover],
                        UserActions[Action.Line],
                        isShowManagePwd ? UserActions[Action.ManagePwd] : undefined,
                    ].filter((action) => action !== undefined),
                },
            ]

        case roleIds.includes(SystemRoleType.OrgManager):
            return [
                {
                    ...MenuGroups[MenuGroupName.DepMgnt],
                    actions: [
                        DepActions[Action.CreateDep],
                        DepActions[Action.EditDep],
                        DepActions[Action.DelDep],
                        DepActions[Action.Line],
                        DepActions[Action.MoveDep],
                        DepActions[Action.AddUsersToDep],
                    ],
                },
                {
                    ...MenuGroups[MenuGroupName.UserMgnt],
                    actions: [
                        UserActions[Action.CreateUser],
                        UserActions[Action.EditUser],
                        UserActions[Action.DelUser],
                        UserActions[Action.Line],
                        // UserActions[Action.SetExpiration],
                        UserActions[Action.MoveUser],
                        UserActions[Action.RemoveUser],
                        isShowEnableAndDisableUser ? UserActions[Action.EnableUser] : undefined,
                        isShowEnableAndDisableUser ? UserActions[Action.DisableUser] : undefined,
                        UserActions[Action.Line],
                        ...isShowSetRole ? [UserActions[Action.SetRole]] : [],
                        UserActions[Action.WorkHandover],
                        !triSystemStatus && isShowManagePwd ? UserActions[Action.Line] : undefined,
                        !triSystemStatus && isShowManagePwd ? UserActions[Action.ManagePwd] : undefined,
                    ].filter((action) => action !== undefined),
                },
                {
                    ...MenuGroups[MenuGroupName.ImportOrg],
                    actions: [
                        ImportActions[Action.ImportOrg],
                        ImportActions[Action.ImportDomain],
                        enableThirdImport ? ImportActions[Action.ImportThirdOrg] : undefined,
                    ].filter((action) => action !== undefined),
                },
                freezeStatus ? {
                    ...MenuGroups[MenuGroupName.FreezeUser],
                    actions: [
                        FreezeUserActions[Action.FreezeUser],
                        FreezeUserActions[Action.UnfreezeUser],
                    ],
                } : undefined,
            ].filter((group) => group !== undefined)

        default:
            return []
    }
}

/**
 * 根据选中的组织部门改变操作的状态
 */
export function changeActionStatusByDep(menuGroup: ReadonlyArray<MenuGroup>, selectedDep: Core.ShareMgnt.ncTDepartmentInfo): ReadonlyArray<MenuGroup> {
    if (selectedDep) {
        switch (true) {
            // 选中未分配组或所有用户
            case selectedDep.data.id === SpecialDep.Unassigned || selectedDep.data.id === SpecialDep.AllUsers:
                return menuGroup.map((group) => {
                    if (group && group.name !== MenuGroupName.FreezeUser) {
                        return {
                            ...group,
                            actions: group.actions.map((action) => {
                                switch (action.type) {
                                    case Action.CreateUser:
                                    case Action.MoveUser:
                                        return selectedDep.data.id === SpecialDep.Unassigned ?
                                            changeActionStatus(action, false)
                                            : changeActionStatus(action, true)

                                    case Action.CreateOrg:

                                    case Action.EditUser:
                                    case Action.DelUser:
                                    case Action.SetExpiration:
                                    case Action.EnableUser:
                                    case Action.DisableUser:
                                        return changeActionStatus(action, false)

                                    case Action.EditOrg:
                                    case Action.DelOrg:
                                    case Action.CreateDep:
                                    case Action.EditDep:
                                    case Action.DelDep:
                                    case Action.MoveDep:
                                    case Action.AddUsersToDep:

                                    case Action.RemoveUser:
                                    case Action.SetRole:
                                    case Action.ProductLicense:
                                    case Action.WorkHandover:
                                    case Action.ManagePwd:

                                    case Action.ImportOrg:
                                    case Action.ImportDomain:
                                    case Action.ImportThirdOrg:
                                        return changeActionStatus(action, true)

                                    default:
                                        return action
                                }
                            }),
                        }
                    }

                    return group
                })

            // 选中组织
            case selectedDep.data.is_root:
                return menuGroup.map((group) => {
                    if (group && group.name !== MenuGroupName.FreezeUser) {
                        return {
                            ...group,
                            actions: group.actions.map((action) => {
                                switch (action.type) {
                                    case Action.CreateOrg:
                                    case Action.EditOrg:
                                    case Action.DelOrg:
                                    case Action.CreateDep:
                                    case Action.AddUsersToDep:
                                    case Action.CreateUser:
                                    case Action.EditUser:
                                    case Action.DelUser:
                                    case Action.SetExpiration:
                                    case Action.MoveUser:
                                    case Action.RemoveUser:
                                    case Action.EnableUser:
                                    case Action.DisableUser:
                                    case Action.ImportOrg:
                                    case Action.ImportDomain:
                                    case Action.ImportThirdOrg:
                                        return changeActionStatus(action, false)

                                    case Action.EditDep:
                                    case Action.DelDep:
                                    case Action.MoveDep:

                                    case Action.SetRole:
                                    case Action.ProductLicense:
                                    case Action.WorkHandover:
                                    case Action.ManagePwd:
                                        return changeActionStatus(action, true)

                                    default:
                                        return action
                                }
                            }),
                        }
                    }

                    return group
                })

            // 选中部门
            default:
                return menuGroup.map((group) => {
                    if (group && group.name !== MenuGroupName.FreezeUser) {
                        return {
                            ...group,
                            actions: group.actions.map((action) => {
                                switch (action.type) {
                                    case Action.CreateOrg:
                                    case Action.CreateDep:
                                    case Action.EditDep:
                                    case Action.DelDep:
                                    case Action.MoveDep:
                                    case Action.AddUsersToDep:

                                    case Action.CreateUser:
                                    case Action.EditUser:
                                    case Action.DelUser:
                                    case Action.SetExpiration:
                                    case Action.MoveUser:
                                    case Action.RemoveUser:
                                    case Action.EnableUser:
                                    case Action.DisableUser:

                                    case Action.ImportOrg:
                                    case Action.ImportDomain:
                                    case Action.ImportThirdOrg:
                                        return changeActionStatus(action, false)

                                    case Action.EditOrg:
                                    case Action.DelOrg:

                                    case Action.SetRole:
                                    case Action.ProductLicense:
                                    case Action.WorkHandover:
                                    case Action.ManagePwd:
                                        return changeActionStatus(action, true)

                                    default:
                                        return action
                                }
                            }),
                        }
                    }

                    return group
                })
        }
    } else {
        return menuGroup.map((group) => {
            return (
                {
                    ...group,
                    actions: group.actions.map((action) => ({ ...action, disabled: true })),
                }
            )
        })
    }
}

/**
 * 根据选中的用户数改变操作状态
 */
export function changeActionStatusBySelectedUsers(menuGroup: ReadonlyArray<MenuGroup>, selectedUsers: ReadonlyArray<Core.ShareMgnt.ncTUsrmGetUserInfo>, selectedDep: Core.ShareMgnt.ncTDepartmentInfo): ReadonlyArray<MenuGroup> {
    const usersLen = selectedUsers.length

    return menuGroup.map((group) => {
        if (group && group.name === MenuGroupName.UserMgnt) {
            return usersLen === 0 ?
                changeActionStatusByDep([group], selectedDep)[0]
                : {
                    ...group,
                    actions: group.actions.map((action) => {
                        switch (action.type) {
                            case Action.EditUser:
                            case Action.WorkHandover:
                            case Action.ProductLicense:
                                return changeActionStatus(action, false)
                            case Action.SetRole:
                            case Action.ManagePwd:
                                return usersLen === 1 ?
                                    changeActionStatus(action, false)
                                    : changeActionStatus(action, true)
                            default:
                                return action
                        }
                    }),
                }
        }

        return group
    })
}

/**
 * 改变操作的状态
 */
export function changeActionStatus(action: any, disabled: boolean): any {
    return {
        ...action,
        disabled,
    }
}

/**
 * 格式化数据
 */
export function formatDatas (data) {
    const formatDatas = data.map((cur) => (
        {
            directDeptInfo: {
                departmentIds: formatDirectDeptInfo(cur.parent_deps).departmentIds,
                departmentNames: formatDirectDeptInfo(cur.parent_deps).departmentNames,
                type: cur.parent_deps.type,
            },
            id:cur.id,
            user: {
                code: cur.code,
                createTime: cur.created_at,
                csfLevel: cur.csf_level,
                csfLevel2: cur.csf_level2,
                departmentIds: formatDirectDeptInfo(cur.parent_deps).departmentIds,
                departmentNames: formatDirectDeptInfo(cur.parent_deps).departmentNames,
                departmentCodes: formatDirectDeptInfo(cur.parent_deps).departmentCodes,
                displayName: cur.name,
                freezeStatus: cur.frozen,
                loginName: cur.account,
                priority: cur.priority,
                position:cur.position,
                remark: cur.remark,
                roles: cur.roles,
                status: !cur.enabled,
                managerDisplayName: cur.manager ? cur.manager.name : '',
                userType: cur.auth_type,
            },
        }
    ))
    return [...formatDatas]
}

/**
 * 格式化部门信息
 */
export function formatDirectDeptInfo(data) {
    const result = {
        departmentIds: [],
        departmentNames: [],
        departmentCodes: [],
        departmentPaths: [],
    };

    data.forEach((group) => {
        if (group.length > 0) {
            const lastDept = group[group.length - 1];
            if(lastDept.id) {
                result.departmentIds = [...result.departmentIds, lastDept.id ]
            }

            if(lastDept.name) {
                result.departmentNames = [...result.departmentNames, lastDept.name ]
            }

            if(lastDept.code) {
                result.departmentCodes = [...result.departmentCodes, lastDept.code ]
            }
            const path = group.map((dept) => dept.name).join('/');
            result.departmentPaths = [...result.departmentPaths, path]
        }
    });

    return result;
}

/**
 * 格式化部门编码
 */
export function formatCode(codes: string[]) {
    return codes.filter((code) => code)
}

/**
 * 格式化上级部门信息
 */
export function formatDepInfo(parent_deps) {
    const result = { code: '', name: '', departmentPaths: '' };
    if(parent_deps) {
        const path = parent_deps.map((dept) => dept.name).join('/');
        result.departmentPaths = path
        if(parent_deps && parent_deps[parent_deps.length - 1]) {
            result.code = parent_deps[parent_deps.length - 1].code
            result.name = parent_deps[parent_deps.length - 1].name
        }
    }
    return result
}