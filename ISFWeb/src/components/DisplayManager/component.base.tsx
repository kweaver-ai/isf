import * as React from 'react';
import { noop } from 'lodash';
import { usrmGetDepartResponsiblePerson } from '@/core/thrift/sharemgnt/sharemgnt'
import WebComponent from '../webcomponent';

export default class DiaplayManagerBase extends WebComponent<Console.DisplayManager.Props, Console.DisplayManager.State> {
    static DefaultProps = {
        departmentId: '',
        departmentName: '',
        userid: '',
        hasPermission: false,
        onComplete: noop,
    }

    state = {
        manager: [],
        isEdited: false,
    }

    async componentDidMount() {
        this.setState({
            manager: await this.getManager(),
        })
    }

    async componentDidUpdate(prevProps, prevState) {
        if (prevProps.departmentId !== this.props.departmentId) {
            this.setState({
                isEdited: false,
                manager: await this.getManager(),
            })
        }
    }

    /**
     * 获取组织管理员
     */
    private getManager = async () => {
        if (this.props.departmentId) {
            return (await usrmGetDepartResponsiblePerson([this.props.departmentId])).map((value) => value.user.displayName)
        }

        return []
    }

    /**
     * 打开编辑窗口
     */
    protected openEditDialog = () => {
        this.setState({
            isEdited: true,
        })
    }

    /**
     * 取消编辑
     */
    protected cancelEditing = () => {
        this.setState({
            isEdited: false,
        })
    }

    /**
     * 更新数据
     */
    protected updateManagers = async () => {
        this.setState({
            isEdited: false,
            manager: await this.getManager(),
        })

        this.props.onComplete()
    }

}