import * as React from 'react'
import { noop } from 'lodash';
import { Toast, Message2 } from '@/sweet-ui';
import { deleteDep } from '@/core/department';
import { Dep } from '@/core/user';
import WebComponent from '../webcomponent';
import __ from './locale';

/*
* 删除部门操作当前状态
*/
export enum Status {
    /**
     * 默认状态
    */
    Normal,

    /**
     * 加载中
    */
    Loading,

    /**
     * 删除错误
    */
    Error,
}

interface DeleteDepartmentProps extends React.Props<void> {
    /**
     * 选择的部门
    */
    dep: Dep;

    /**
     * 当前登录的用户
    */
    userid: string;

    /**
     * 取消删除部门
    */
    onRequestCancelDeleteDep: () => any;

    /**
     * 删除部门成功
    */
    onDeleteDepSuccess: (depInfo: Dep) => any;
}

interface DeleteDepartmentState {
    /*
    * 删除部门操作当前状态
    */
    status: Status;
}

export default class DeleteDepartmentBase extends WebComponent<DeleteDepartmentProps, DeleteDepartmentState> {
    static defaultProps = {
        dep: null,
        userid: '',
        onRequestCancelDeleteDep: noop,
        onDeleteDepSuccess: noop,
    }

    state = {
        status: Status.Normal,
    }

    /**
     * 确定删除部门操作
     */
    protected async confirmDeleteDepartment(): Promise<void> {
        this.setState({
            status: Status.Loading,
        })
        const { dep, onDeleteDepSuccess, onRequestCancelDeleteDep } = this.props;

        try {
            await deleteDep(dep.id)

            this.setState({
                status: Status.Normal,
            })
            Toast.open(__('删除成功'))
            onDeleteDepSuccess(dep);
        } catch (error) {
            const { message, description } = error

            this.setState({
                status: Status.Error,
            })

            if (message || description) {
                Message2.info({ message: description || message })
            }

            // 删除失败则取消删除
            onRequestCancelDeleteDep()
        }
    }
}