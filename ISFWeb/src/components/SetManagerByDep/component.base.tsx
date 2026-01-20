import * as React from 'react';
import { noop } from 'lodash';
import { usrmGetDepartResponsiblePerson, setUserRolemMember, deleteUserRolemMember, getRoleMemberDetail } from '@/core/thrift/sharemgnt/sharemgnt';
import { getUserInfo } from '@/core/thrift/user/user';
import { SystemRoleType } from '@/core/role/role';
import { QuotaType } from '@/core/quota';
import { ShareMgnt } from '@/core/thrift';
import { manageLog, Level, ManagementOps } from '@/core/log';
import WebComponent from '../webcomponent';
import __ from './locale';

export enum ValidateState {
    Normal,
    Empty,
    InvalidSpace,
}

export const ValidateMessages = {
    [ValidateState.Empty]: __('此项不允许为空。'),
    [ValidateState.InvalidSpace]: __('配额空间值为不超过 1000000 的正数，支持小数点后两位，请重新输入。'),
}

export default class SetManagerByDepBase extends WebComponent<Console.SetManagerByDep.Props, Console.SetManagerByDep.State> {
    static defaultProps = {
        departmentId: '',
        departmentName: '',
        userid: '',
        onSetSuccess: noop,
        onCancel: noop,
    }

    state = {
        isConfigManager: true,
        isAddingManager: false,
        managers: [],
        currentUser: null,
        isLimitUserSpace: false,
        limitUserSpace: '',
        limitUserSpaceState: ValidateState.Normal,
        isLimitDocSpace: false,
        limitDocSpace: '',
        limitDocSpaceState: ValidateState.Normal,
        errorStatus: null,
        isSetting: false,
        limitCheckDisable: {
            limitUserCheckDisable: false,
            limitDocCheckDisable: false,
        },
    }

    /**
     * 增加的用户
     */
    addedManagers = [];

    /**
     * 删除的组织管理员
     */
    deleteManagers = [];

    /**
     * 原始数据
     */
    originManagers = null;

    /**
     * 当前登录管理员的限额信息
     */
    limitSpaceInfo: {
        userSpace: number | string;
        docSpace: number | string;
    } = {
        userSpace: '',
        docSpace: '',
    };

    /**
     * 当前登录管理员的角色信息
     */
    roles: ReadonlyArray<any> = []

    async componentDidMount() {
        let managers = await usrmGetDepartResponsiblePerson([this.props.departmentId]);
        this.originManagers = managers;
        this.setState({
            managers,
        })

        // 获取并存储当前登录用户的限额信息、角色信息
        const { userid } = this.props
        try {
            const { user: { limitSpaceInfo: { limitDocSpace, limitUserSpace }, roles } } = await ShareMgnt('Usrm_GetUserInfo', [userid])
            this.limitSpaceInfo = {
                userSpace: limitUserSpace === -1 ? '' : (limitUserSpace / Math.pow(1024, 3)).toFixed(2),
                docSpace: limitDocSpace === -1 ? '' : (limitDocSpace / Math.pow(1024, 3)).toFixed(2),
            };
            this.roles = roles;

            // 如果是组织管理员，并且被限额，则限额复选框勾选并灰化，显示限额信息
            const { userSpace, docSpace } = this.limitSpaceInfo
            if (roles.some((item) => item.id === SystemRoleType.OrgManager) && (limitDocSpace !== -1 || limitUserSpace !== -1)) {
                this.setState({
                    limitUserSpace: userSpace,
                    limitDocSpace: docSpace,
                    isLimitUserSpace: limitUserSpace === -1 ? false : true,
                    isLimitDocSpace: limitDocSpace === -1 ? false : true,
                    limitCheckDisable: {
                        limitUserCheckDisable: limitUserSpace === -1 ? false : true,
                        limitDocCheckDisable: limitDocSpace === -1 ? false : true,
                    },
                })
            } else {
                this.setState({
                    limitUserSpace: userSpace,
                    limitDocSpace: docSpace,
                    isLimitUserSpace: false,
                    isLimitDocSpace: false,
                    limitCheckDisable: {
                        limitUserCheckDisable: false,
                        limitDocCheckDisable: false,
                    },
                })
            }
        } catch (error) {

        }
    }
    /**
     * 获取部门名
     * @param dep 部门
     * @return 返回部门名
     */
    protected userDataFormatter = (person: Core.ShareMgnt.ncTUsrmGetUserInfo): string => {
        return person.user.displayName
    }

    /**
     * 删除组织管理员
     */
    protected deleteManager = (data: Array<Core.ShareMgnt.ncTUsrmGetUserInfo>) => {
        for (let manager of this.state.managers) {
            if (!data.find((value) => manager.id === value.id)) {
                this.addedManagers = this.addedManagers.filter((addManager) => addManager.id !== manager.id)
                if (this.originManagers.find((originManager) => originManager.id === manager.id)) {
                    this.deleteManagers = [...this.deleteManagers, manager]
                }
            }
        }
        this.setState({
            managers: data,
        })

    }

    /**
     * 打开添加窗口
     */
    protected openAddManager = () => {
        this.setState({
            isAddingManager: true,
            currentUser: null,
            limitUserSpaceState: ValidateState.Normal,
            limitDocSpaceState: ValidateState.Normal,
        })
    }

    /**
     * 选择用户
     */
    protected selectUser = async (user) => {
        const userInfo = await getUserInfo(user.id)
        this.setState({
            currentUser: userInfo,
        })
    }

    /**
     * 取消编辑
     */
    protected cancelAddManager = () => {
        this.setState({
            isAddingManager: false,
            currentUser: null,
        })
    }

    /**
     * 确定增加管理员
     */
    protected onConfirmAddManager = () => {
        const checkDocResult = this.checkDocSpace();
        const checkUserResult = this.checkUserSpace()
        if (!checkDocResult || !checkUserResult) {
            return
        }
        const limitSpaceInfo = {
            limitUserSpace: this.state.isLimitUserSpace ? Math.ceil(Number(this.state.limitUserSpace) * Math.pow(1024, 3)) : -1,
            limitDocSpace: this.state.isLimitDocSpace ? Math.ceil(Number(this.state.limitDocSpace) * Math.pow(1024, 3)) : -1,
        }
        this.setState({
            managers: [...this.state.managers.filter((value) => {
                return value.id !== this.state.currentUser.id
            }), { ...this.state.currentUser, user: { ...this.state.currentUser.user, limitSpaceInfo: limitSpaceInfo } }],
            isAddingManager: false,
        })
        this.addedManagers = [...this.addedManagers, { ...this.state.currentUser, user: { ...this.state.currentUser.user, limitSpaceInfo: limitSpaceInfo } }];
    }

    /**
     * 确定保存组织管理员
     */
    protected onConfirmManager = async () => {
        this.setState({
            isConfigManager: false,
            isSetting: true,
        })
        try {
            await this.deletingManager()
            await this.addManagers()
            this.setState({
                isSetting: false,
            })
            this.props.onSetSuccess()
        } catch (ex) {
            this.setState({
                errorStatus: ex,
                isSetting: false,
            })
        }
    }

    /**
     * 增加的组织管理员
     */
    private async addManagers() {
        for (let manager of this.addedManagers) {
            let managerInfo = {
                departmentIds: [this.props.departmentId],
                departmentNames: [this.props.departmentName],
            }
            if (manager.user.roles && manager.user.roles.some((role) => role.id === SystemRoleType.OrgManager)) {
                let curManagerInfo = (await getRoleMemberDetail([this.props.userid, SystemRoleType.OrgManager, manager.id])).manageDeptInfo

                if (curManagerInfo) {
                    managerInfo = {
                        departmentIds: [...curManagerInfo.departmentIds.filter((id) => id !== this.props.departmentId), this.props.departmentId],
                        departmentNames: [...curManagerInfo.departmentNames.filter((name) => name !== this.props.departmentName), this.props.departmentName],
                    }
                }
            }

            let memberInfo = {
                userId: manager.id,
                displayName: manager.user.displayName,
                departmentIds: manager.user.departmentIds,
                departmentNames: manager.user.departmentNames,
                manageDeptInfo: {
                    ncTManageDeptInfo: {
                        departmentIds: managerInfo.departmentIds,
                        departmentNames: managerInfo.departmentNames,
                        limitUserSpaceSize: manager.user.limitSpaceInfo.limitUserSpace,
                        limitDocSpaceSize: manager.user.limitSpaceInfo.limitDocSpace,
                    },
                },
            }
            try {
                await setUserRolemMember([this.props.userid, SystemRoleType.OrgManager, memberInfo])
                manageLog(
                    ManagementOps.SET,
                    __('将 “${userName}” 设为组织管理员，管辖部门：“${departmentName}”', { userName: manager.user.displayName, departmentName: this.props.departmentName }),
                    null,
                    Level.INFO,
                )
            } catch (ex) {
                throw ex
            }
        }
    }

    /**
     * 删除的组织管理员
     */
    private async deletingManager() {
        const delManagers = this.deleteManagers.filter((delManager) => !this.addedManagers.some((addManager) => addManager.id === delManager.id))
        for (let manager of delManagers) {
            let managerInfo = {
                departmentIds: [],
                departmentNames: [],
            }

            let curManagerInfo = (await getRoleMemberDetail([this.props.userid, SystemRoleType.OrgManager, manager.id])).manageDeptInfo

            if (curManagerInfo) {
                managerInfo = {
                    departmentIds: [...curManagerInfo.departmentIds.filter((id) => id !== this.props.departmentId)],
                    departmentNames: [...curManagerInfo.departmentNames.filter((name) => name !== this.props.departmentName)],
                }
            }

            let memberInfo = {
                userId: manager.id,
                displayName: manager.user.displayName,
                departmentIds: manager.user.departmentIds,
                departmentNames: manager.user.departmentNames,
                manageDeptInfo: {
                    ncTManageDeptInfo: {
                        departmentIds: managerInfo.departmentIds,
                        departmentNames: managerInfo.departmentNames,
                        limitUserSpaceSize: manager.user.limitSpaceInfo.limitUserSpace,
                        limitDocSpaceSize: manager.user.limitSpaceInfo.limitDocSpace,
                    },
                },
            }
            try {
                if (managerInfo.departmentIds.length) {
                    await setUserRolemMember([this.props.userid, SystemRoleType.OrgManager, memberInfo])
                } else {
                    await deleteUserRolemMember([this.props.userid, SystemRoleType.OrgManager, manager.id])
                }
                manageLog(
                    ManagementOps.DELETE,
                    __('取消 “${userName}” 的组织管理员身份，管辖部门：“${departmentName}”', { userName: manager.user.displayName, departmentName: this.props.departmentName }),
                    null,
                    Level.INFO,
                )
            } catch (ex) {
                throw ex
            }
        }
    }

    /**
     * 关闭错误弹窗
     */
    protected closeError = () => {
        this.setState({
            errorStatus: null,
        })
        this.props.onCancel()
    }

    /**
     * 检查用户输入框的合法性
     */
    private checkUserSpace() {
        if (this.state.limitUserSpace === '' && this.state.isLimitUserSpace) {
            this.setState({
                limitUserSpaceState: ValidateState.Empty,
            })
            return false;
        } else if ((Number(this.state.limitUserSpace) <= 0 || Number(this.state.limitUserSpace) > 1000000) && this.state.isLimitUserSpace) {
            this.setState({
                limitUserSpaceState: ValidateState.InvalidSpace,
            })
            return false;
        }
        return true;
    }

    /**
     * 检查文档输入框的合法性
     */
    private checkDocSpace() {
        if (this.state.limitDocSpace === '' && this.state.isLimitDocSpace) {
            this.setState({
                limitDocSpaceState: ValidateState.Empty,
            })
            return false;
        } else if ((Number(this.state.limitDocSpace) <= 0 || Number(this.state.limitDocSpace) > 1000000) && this.state.isLimitDocSpace) {
            this.setState({
                limitDocSpaceState: ValidateState.InvalidSpace,
            })
            return false;
        }
        return true;
    }

    /**
     * 判断输入框的值是否是不超过 1000000 的正数，支持小数点后两位
     */
    protected isNumberPoint(input: any): boolean {
        return /^([1-9]\d{0,5}|0)(\.\d{0,2})?$|^1000000$/.test(String(input))
    }

    /**
     * 输入框的值为空且失焦时气泡提示
     */
    protected handleOnBlur = (type: QuotaType): void => {
        const { limitUserSpace, limitDocSpace } = this.state;

        switch (type) {
            case QuotaType.UserSpace:
                if (limitUserSpace === '') {
                    this.setState({
                        limitUserSpaceState: ValidateState.Empty,
                    })
                }
                break
            case QuotaType.DocSpace:
                if (limitDocSpace === '') {
                    this.setState({
                        limitDocSpaceState: ValidateState.Empty,
                    })
                }
                break

            default:
                break
        }
    }
}