import * as React from 'react'
import { noop } from 'lodash';
import { manageLog, Level, ManagementOps } from '@/core/log';
import { removeUserFromDepartment, getDepParentPathById } from '@/core/thrift/sharemgnt/sharemgnt';
import WebComponent from '../webcomponent';
import { Range, getSeletedUsers } from '../helper';
import __ from './locale';

export enum Status {
    NORMAL,

    CURRENT_USER_INCLUDED, // 当前用户

    LOADING, // 加载中
}

interface DepInfo {
    id: string;
    name: string;
    parentPath: string;
}

interface Props {
    users?: Array<any>; // 选择的用户 * any 后续补充

    dep: DepInfo; // 选择的部门 * any 后续补充

    userid: string; // 当前登录的用户

    onComplete: () => any; // 移除结束的事件

    onSuccess: (range: Range) => {}; // 移除成功
}

interface State {
    selected: Range; // 选择移除的对象

    status: Status; // 移除的状态
}

export default class RemoveUserBase extends WebComponent<Props, State> {
    static defaultProps = {
        users: [],
        dep: null,
        userid: '',
        onComplete: noop,
        onSuccess: noop,
    }

    state = {
        selected: this.props.users.length ? Range.USERS : Range.DEPARTMENT,
        status: Status.NORMAL,
    }

    // 移动单个用户时选择的部门
    dep: DepInfo

    // 是否在移除单个用户
    singleUser: boolean = this.props.users && this.props.users.length === 1 && this.state.selected === Range.USERS

    componentDidMount() {
        this.setState({
            selected: this.props.users.length ? Range.USERS : Range.DEPARTMENT,
        })
    }

    /**
 * 检查是否存在不能被移除的用户
 */
    checkUser(users: Array<any>): boolean {
        return !users.some((value, index) => {
            if (value.id === this.props.userid) {
                this.setState({
                    status: Status.CURRENT_USER_INCLUDED,
                })
                return true
            }
        })
    }

    /**
 * 移除用户
 */
    async removeUsers(users: Array<any>) {
        for (let user of users) {
            let idsInfo: ReadonlyArray<{ id: string; path: string }> = []

            if (this.singleUser) {
                idsInfo = [{ id: this.dep.id, path: this.dep.path }]
            } else {
                const { user: userInfo, directDeptInfo } = user
                // 判断有无直属部门
                if (directDeptInfo) {
                    const [{ parentPath }] = await getDepParentPathById([directDeptInfo.departmentId])
                    const path = `${parentPath ? parentPath + '/' : ''}${directDeptInfo.departmentName}`

                    idsInfo = [{ id: directDeptInfo.departmentId, path }]
                } else {
                    if (userInfo && Array.isArray(userInfo.departmentIds)) {
                        const { departmentIds, departmentNames } = userInfo
                        // 获取当前默认部门信息
                        const parentPathList = await getDepParentPathById([...departmentIds, this.props.dep.id])
                        // 当前(进入)部门的路径
                        const depPath = `${parentPathList[parentPathList.length - 1].parentPath ? parentPathList[parentPathList.length - 1].parentPath + '/' : ''}${this.props.dep.name}`

                        departmentIds.forEach((id, index) => {
                            const { parentPath: parent } = parentPathList[index]
                            // 拼接用户所属部门路径
                            const path = `${parent ? parent + '/' : ''}${departmentNames[index]}`

                            if (path.startsWith(depPath)) {
                                idsInfo = [...idsInfo, { id, path }]
                            }
                        })
                    }
                }
            }

            idsInfo.forEach(async ({ id, path }) => {
                try {
                    let result = await removeUserFromDepartment(
                        [user.id],
                        id,
                    )
                    if (result && result.length === 0) {
                        await this.logRemoved(user, path);
                    }
                }
                catch (ex) {
                    if (ex.error.errID !== 20110) {
                        this.setState({
                            status: ex.error.errID,
                        })
                        throw (ex);
                    }
                }
            })
        }
    }

    /**
 * 确定事件
 */
    async confirmRemoveUsers() {
        this.setState({
            status: Status.LOADING,
        })
        const users = await getSeletedUsers(this.state.selected, this.props.dep, this.props.users);

        if (this.checkUser(users)) {
            await this.removeUsers(users);
            this.props.onSuccess(this.state.selected);
        }

    }

    logRemoved(user, depPath) {
        return manageLog(ManagementOps.REMOVE,
            __('移除用户 “${displayName}(${loginName})” 成功；移除部门：${depPath}', {
                displayName: user.user.displayName,
                loginName: user.user.loginName,
                depPath,
            }),
            null,
            Level.INFO,
        )
    }

    /**
     * 选择要删除的对象
     * @param value 部门、部门及其子部门、所选中的
     */
    onSelectedType(value) {
        this.singleUser = value === Range.USERS && this.props.users && this.props.users.length === 1

        this.setState({
            selected: value,
        })
    }

    /**
     * 操作单个用户时，选择用户的部门
     */
    protected onDepChange = (dep): void => {
        this.dep = dep
    }
}