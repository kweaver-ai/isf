import * as React from 'react'
import { noop } from 'lodash';
import { Message2, Toast } from '@/sweet-ui';
import { Dep } from '@/core/user';
import { deleteDep } from '@/core/department';
import WebComponent from '../webcomponent';
import __ from './locale';

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

interface DeleteOrganizationProps extends React.Props<void> {
    /**
     * 选择的部门
    */
    dep: Dep;

    /**
     * 当前登录的用户
    */
    userid: string;

    /**
     * 取消删除组织
    */
    onRequestCancelDeleteOrg: () => void;

    /**
     * 删除组织成功
    */
    onDeleteOrgSuccess: (orgInfo: Dep) => void;
}

interface DeleteOrganizationState {
    /*
    * 删除组织操作当前状态
    */
    status: number;

    /*
    * 二次确认删除
    */
    confirmDelete: boolean;
}

export default class DeleteOrganizationBase extends WebComponent<DeleteOrganizationProps, DeleteOrganizationState> {
    static defaultProps = {
        dep: null,
        userid: '',
        onRequestCancelDeleteOrg: noop,
        onDeleteOrgSuccess: noop,
    }

    state = {
        status: Status.Normal,
        confirmDelete: false,
    }

    /**
     * 确定删除组织
     */
    protected async confirmDeleteOrganization() {
        this.setState({
            status: Status.Loading,
        })
        const { dep, onDeleteOrgSuccess, onRequestCancelDeleteOrg } = this.props;
        try {
            await deleteDep(dep.id)

            this.setState({
                status: Status.Normal,
            })
            Toast.open(__('删除成功'))
            onDeleteOrgSuccess(dep);
        } catch (error) {
            const { message, description } = error

            this.setState({
                status: Status.Error,
            })

            if (message || description) {
                Message2.info({ message: description || message })
            }

            // 删除失败则取消删除
            onRequestCancelDeleteOrg()
        }
    }

    /**
    * 二次确认删除组织
    */
    protected confirmDelete(confirmDelete: boolean) {
        this.setState({
            confirmDelete,
        })
    }
}