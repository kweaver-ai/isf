import * as React from 'react';
import { noop } from 'lodash';
import session from '@/util/session';
import { addUsersToDep, editUserOss, getOrgDepartmentById } from '@/core/thrift/sharemgnt/sharemgnt';
import { getObjectStorageInfoById } from '@/core/apis/console/ossgateway'
import { ErrorCode } from '@/core/thrift/sharemgnt/errcode';
import { manageLog, Level, ManagementOps } from '@/core/log';
import { Message2, Toast } from '@/sweet-ui';
import { SystemRoleType, SystemAccountId } from '@/core/role/role';
import { displayUserOssInfo } from '@/core/oss/oss';
import { FormatedNodeInfo } from '@/core/organization';
import WebComponent from '../webcomponent';
import { Status } from './helper';
import __ from './locale';

interface AddUsersToDepProps {
    /**
     * 用户id
     */
    userid: string;

    /**
     * 目标部门
     */
    targetDep: Core.ShareMgnt.ncTDepartmentInfo;

    /**
     * 取消添加用户至部门
     */
    onRequestCancel: () => void;

    /**
     * 添加用户到部门完成
     */
    onRequestSuccess: () => void;

    /**
     * 部门不存在，移除部门
     */
    onRequestRemoveDep: (targetDep: Core.ShareMgnt.ncTDepartmentInfo) => void;
}

interface AddUsersToDepState {
    /**
     * 添加的用户
     */
    users: ReadonlyArray<FormatedNodeInfo>;

    /**
     * 当前渲染界面
     */
    renderStatus: Status;
}

export default class AddUsersToDepBase extends WebComponent<AddUsersToDepProps, AddUsersToDepState> {
    static defaultProps = {
        onRequestCancel: noop,
        onRequestSuccess: noop,
        onRequestRemoveDep: noop,
    }

    state = {
        users: [],
        renderStatus: Status.Config,
    }

    userid = this.props.userid || session.get('isf.userid');

    ossInfo: Core.ShareMgnt.ncTUsrmOSSInfo

    isShowUndistributed = session.get('isf.userInfo').user.roles.some(({ id }) => [SystemRoleType.Admin, SystemRoleType.Supper].includes(id))
        ||
        [SystemAccountId.Admin, SystemAccountId.Securit].includes(this.userid);

    /**
     * 不存在的用户
     */
    notExistUsers = []

    async componentDidMount() {
        let ossInfo = this.props.targetDep.ossInfo

        if (!ossInfo) {
            ({ ossInfo } = await getOrgDepartmentById(this.props.targetDep.id))
        }

        if (ossInfo.ossId) {
            const ossData = await getObjectStorageInfoById(ossInfo.ossId)
            this.ossInfo = {
                enabled: ossData.enabled,
                ossId: ossData.id,
                ossName: ossData.name,
            }
        } else {
            this.ossInfo = ossInfo
        }
    }

    /**
     * 添加用户
     */
    protected addUsers = (users: ReadonlyArray<FormatedNodeInfo>): void => {
        this.setState({
            users,
        })
    }

    /**
     * 确认添加用户至部门
     */
    protected confrimAddUsers = async (): Promise<void> => {

        const { users } = this.state

        const loginUser = users.find((user) => user.id === this.userid)

        if (loginUser) {
            Message2.info({ message: __('无法添加用户“${userName}”，此用户为当前登录账号，请重新选择。', { userName: loginUser.name }) })
        } else {
            this.setState({
                renderStatus: Status.None,
            })
            
            await this.requestAddUsers(users, false)

        }
    }

    /**
     * 请求添加用户至部门
     */
    private async requestAddUsers(users: ReadonlyArray<FormatedNodeInfo>, isChangeOss: boolean): Promise<void> {
        let successUsers = []

        try {
            this.setState({
                renderStatus: Status.Adding,
            })

            successUsers = (await addUsersToDep([users.map((user) => user.id), this.props.targetDep.id]))
                .map((uid) => users.find((u) => u.id === uid))

            if (successUsers.length !== users.length) {
                this.notExistUsers = [
                    ...this.notExistUsers,
                    ...users.filter((user) => !successUsers.some((sUser) => sUser.id === user.id)),
                ]
            }

        } catch (ex) {
            if (ex && ex.error && ex.error.errID) {
                switch (ex.error.errID) {
                    case ErrorCode.DepOrOrgNotExist:
                        this.setState({
                            renderStatus: Status.None,
                        })

                        if (await Message2.info({ message: __('添加失败，部门“${depName}”不存在。', { depName: this.props.targetDep.name }) })) {
                            this.props.onRequestRemoveDep(this.props.targetDep)
                        }

                        return

                    case ErrorCode.UserNotExist:
                        this.notExistUsers = [...this.notExistUsers, ex.error.detail]

                        break

                    default:
                        break
                }
            }
        }

        await this.logAddUserSuccess(successUsers)

        if (isChangeOss) {
            this.setState({
                renderStatus: Status.ChangeOSSing,
            })

            await this.changeOss(successUsers, this.ossInfo)
        }

        this.setState({
            renderStatus: Status.None,
        })

        if (this.notExistUsers.length) {
            await Message2.info({ message: __('无法添加用户“${userName}”，该用户不存在。', { userName: this.notExistUsers.map((user) => user.name).join(',') }) })
        }

        Toast.open(__('操作完成，本次添加${num}个用户', { num: (users.length - this.notExistUsers.length) }))

        this.props.onRequestSuccess()
    }

    /**
     * 记录添加用户至部门的日志
     */
    private async logAddUserSuccess(users: ReadonlyArray<FormatedNodeInfo>): Promise<void> {
        for (let user of users) {
            await manageLog(
                ManagementOps.COPY,
                __(
                    '添加用户 “${displayName}(${loginName})” 到部门 “${depName}” 成功',
                    {
                        displayName: user.name,
                        loginName: user.account,
                        depName: this.props.targetDep.name,
                    },
                ),
                '',
                Level.INFO,
            )
        }
    }

    /**
     * 修改存储位置
     */
    private async changeOss(users: ReadonlyArray<FormatedNodeInfo>, ossInfo: Core.ShareMgnt.ncTUsrmOSSInfo): Promise<void> {
        for (let user of users) {
            try {
                await editUserOss(user.id, ossInfo.ossId || '')

                await manageLog(
                    ManagementOps.SET,
                    __(
                        '编辑用户 "${displayName}(${loginName})"的存储位置 成功',
                        {
                            displayName: user.name,
                            loginName: user.account,
                        },
                    ),
                    __('存储位置 “${storage}”', { storage: displayUserOssInfo(ossInfo) }),
                    Level.INFO,
                )

            } catch (ex) {
                if (ex && ex.error && ex.error.errID) {
                    if (ex.error.errID === ErrorCode.UserNotExist) {
                        this.notExistUsers = [...this.notExistUsers, user]

                        continue
                    } else {
                        await Message2.info({ message: __('目标部门的存储位置已不可用，无法替换，仍保留用户当前的存储位置。') })

                        break
                    }
                }
            }

        }
    }
}