/* eslint-disable spaced-comment */
/// <reference path="./component.base.d.ts" />
/// <reference path="@/core/thrift/sharemgnt/sharemgnt.d.ts"/>
/* eslint-enable spaced-comment */

import { uniqBy } from 'lodash';
import { getDepartmentById, setUserFreezeStatus } from '@/core/thrift/sharemgnt/sharemgnt';
import { manageLog, Level, ManagementOps } from '@/core/log';
import { NodeType, getNodeType } from '@/core/organization';
import { Range, getSeletedUsers } from '../helper';
import WebComponent from '../webcomponent';
import { Status } from './helper';
import __ from './locale';

export default class SetUsersFreezeStatusBase extends WebComponent<Components.SetUsersFreezeStatusBase.Props, any> {

    static defaultProps = {
        freezeStatus: false,
    }

    state: Components.SetUsersFreezeStatusBase.State = {
        selecteds: [],
        status: Status.NORMAL,
    }

    // 错误信息
    errorInfo = null;

    /**
     * 获取选中的用户
     * @param selected //选中项
     */
    async getSelected(selected): Promise<void> {
        this.setState({
            selecteds: uniqBy([...this.state.selecteds, (await this.convererData(selected))], 'id'),
        })
    }

    /**
     * 清空已选中项
     */
    clearSelected(): void {
        this.setState({
            selecteds: [],
        })
    }

    /**
     * 获取选中的用户
     */
    async getSelectedUser(): Promise<Array<Core.ShareMgnt.ncTUsrmUserInfo>> {
        let users = [];
        for (let selected of this.state.selecteds) {
            if (selected.type === NodeType.USER) {
                users = uniqBy([...users, selected.data], 'id')
            } else {
                users = uniqBy([...users, ...(await getSeletedUsers(Range.DEPARTMENT_DEEP, selected))], 'id')
            }
        }
        return users
    }

    /**
     * 转换数据
     * @param selected 当前选中的项
     */
    async convererData(selected): Promise<Components.SetUsersFreezeStatusBase.Data> {
        return {
            id: selected.id || selected.departmentId,
            name: selected.name || selected.displayName || selected.departmentName || (selected.user && selected.user.displayName),
            type: getNodeType(selected),
            parentName: (await this.getSelectedParent(selected)).parentDepartName || this.getUserParentsName(selected),
            data: await this.getSelectedParent(selected),
        }
    }

    /**
     * 获取选中用户的父部门
     * @param selected 当前选中项
     */
    async getSelectedParent(selected): Promise<Core.ShareMgnt.ncTUsrmUserInfo | Core.ShareMgnt.ncTUsrmDepartmentInfo | any> {
        if (getNodeType(selected) === NodeType.USER) {
            return selected
        } else if (getNodeType(selected) === NodeType.DEPARTMENT) {
            return (await getDepartmentById(selected.id || selected.departmentId))
        } else {
            return {
                ...selected,
                parentDepartName: '---',
            }
        }
    }

    /**
     * 删除选中项
     * @param willDeleteData
     */
    deleteSelected(willDeleteData): void {
        this.setState({
            selecteds: this.state.selecteds.filter((data) => {
                return data.id !== willDeleteData.id;
            }),
        })
    }

    /**
     * 设置冻结状态
     */
    async setFreezeStatus(): Promise<void> {
        this.setState({
            status: Status.LOADING,
        })
        let selectedUsers = await this.getSelectedUser();
        if (this.checkCurrentUser(selectedUsers)) {
            for (let selectedUser of selectedUsers) {
                try {
                    setUserFreezeStatus(selectedUser.id, this.props.freezeStatus);
                    if (this.props.freezeStatus) {
                        manageLog(ManagementOps.SET,
                            __('冻结 用户“${displayName}（${loginName}）”成功', {
                                displayName: selectedUser.displayName ? selectedUser.displayName : selectedUser.user.displayName,
                                loginName: selectedUser.loginName ? selectedUser.loginName : selectedUser.user.loginName,
                            }),
                            null,
                            Level.INFO,
                        )
                    } else {
                        manageLog(ManagementOps.SET,
                            __('解冻 用户“${displayName}（${loginName}）”成功', {
                                displayName: selectedUser.displayName ? selectedUser.displayName : selectedUser.user.displayName,
                                loginName: selectedUser.loginName ? selectedUser.loginName : selectedUser.user.loginName,
                            }),
                            null,
                            Level.INFO,
                        )
                    }

                }
                catch (ex) {
                    this.setState({
                        status: Status.ERROR,
                    })
                    this.errorInfo = ex;
                }
            }
            this.props.onSuccess();
        } else {
            this.setState({
                status: Status.CURRENT_USER_INCLUDED,
            })
        }
    }

    /**
     * 检查是否存在当前登录用户
     * @param selectedUsers
     */
    checkCurrentUser(selectedUsers): boolean {
        for (let selectedUser of selectedUsers) {
            if (selectedUser.id === this.props.userid) {
                return false
            }
        }
        return true
    }

    /**
     * 获取用户信息
     */
    getUserParentsName(selected): string {
        return selected.departmentNames ?
            selected.departmentNames.join(',') :
            selected.user.departmentNames.join(',')
    }

}
