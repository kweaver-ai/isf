import * as React from 'react';
import { noop, uniqBy } from 'lodash';
import { Message2 } from '@/sweet-ui'
import { ExtraRoot } from '@/core/organization';
import { searchDepartmentOfUsers } from '@/core/thrift/sharemgnt/sharemgnt'
import { NodeType, FormatedNodeInfo } from '@/core/organization';
import WebComponent from '../webcomponent';
import __ from './locale';

interface Props {
    /**
     * zIndex
    */
    zIndex: number;

    // 点击确定事件
    onConfirm: (data: Array<Record<string, any>>) => any;
    // 点击取消事件
    onCancel: () => any;
    // 当前管理员id
    userid: string;
    // 是否加载用户
    selectType?: Array<NodeType>;

    // dialog 头信息
    title: string;

    // 初始值
    data?: Array<any>;

    // 数据转换内部数据结构
    converterIn?: (x) => FormatedNodeInfo;

    // 数据转换外部数据结构
    convererOut?: (node: FormatedNodeInfo) => any;

    /**
     * 是否为单选
     */
    isSingleChoice: boolean;

    // 是否显示未分配组
    isShowUndistributed?: boolean;

    /**
     * 是否只显示被禁用的用户
     */
    isShowDisabledUsers?: boolean;

    /**
     * 额外根节点
     */
    extraRoots: ReadonlyArray<ExtraRoot>;

    /**
     * 是否是级联树
     */
    isCascadeTree: boolean;

    /**
     * 已选列表是否只显示用户
     */
    isOnlyShowUser?: boolean;

    /**
     * 提示
     */
    tip?: string;

    /**
     * 提示的样式
     */
    tipStyle?: Record<string, any>;

    /**
     * 禁用节点时触发
     */
    getNodeStatus?: (node) => any;

    /**
     * 是否显示设置当前登陆用户为文档库所有者
     */
    isShowSetLoginUser?: boolean;

    /**
     * 用户信息
     */
    userInfo?: {name: string; type: NodeType | any; id: string};
}

interface State {
    // 新加的部门
    data: FormatedNodeInfo[];

    /**
     * 是否正在会根据组织或部门id获取用户
     */
    isGetUsersLoading: boolean;
}

export default class OrganizationPickerBase extends WebComponent<Props, State> {

    static defaultProps = {
        onConfirm: noop,
        onCancel: noop,
        userid: '',
        selectType: [NodeType.ORGANIZATION, NodeType.DEPARTMENT],
        title: __('添加部门'),
        tip: '',
        tipStyle: {},
        data: [],
        converterIn: (x) => x,
        isSingleChoice: false,
        isCascadeTree: false,
        isShowDisabledUsers: true,
        isShowUndistributed: false,
        isOnlyShowUser: false,
        getNodeStatus: () => ({ disabled: false }),
    }

    state = {
        data: this.props.data.map(this.props.converterIn),
        isGetUsersLoading: false,
    }

    /**
     * 用来存储级联树的数据
     */
    departmentTreeData: HTMLInputElement | null;

    /**
     * 选择共享者
     * @param value 共享者
     */
    protected selectDep = async (value) => {
        let addition: FormatedNodeInfo[] = [value]

        if (this.props.isOnlyShowUser && value.type !== NodeType.USER) {
            this.setState({ isGetUsersLoading: true })
            const users = await searchDepartmentOfUsers([value.id || value.departmentId, '', 0, -1])

            addition = [
                ...users.reduce((prev, userInfo) => (
                    (!this.props.isShowDisabledUsers && userInfo.user.status === 1) ?
                        prev
                        : [
                            ...prev,
                            {
                                type: NodeType.USER,
                                id: userInfo.id,
                                name: userInfo.user.displayName,
                                account: userInfo.user.loginName,
                                parent_path: value.parent_path ? value.parent_path + '/' + value.name : value.name,
                                original: userInfo,
                            },
                        ]
                ), []),
            ]
        }

        this.setState(({ data }) => ({
            isGetUsersLoading: false,
            data: this.props.isSingleChoice ?
                addition
                : uniqBy([...data, ...addition], 'id'),
        }))
    }

    /**
     * 将级联树数据加入已选列表
     */
    protected addTreeData = async (): Promise<void> => {
        const selected = await this.departmentTreeData.getSelections()

        if (selected && selected.length) {
            let addition: any[] = []
            if (this.props.isOnlyShowUser) {
                for (let item of selected) {
                    if (item.type === NodeType.USER) {
                        addition = [
                            ...addition,
                            item,
                        ]
                    } else {
                        try {
                            this.setState({ isGetUsersLoading: true })
                            const users = await searchDepartmentOfUsers([item.id, '', 0, -1])
                            addition = [
                                ...addition,
                                ...users.reduce((prev, userInfo) => (
                                    (!this.props.isShowDisabledUsers && userInfo.user.status === 1) ?
                                        prev
                                        : [
                                            ...prev,
                                            {
                                                type: NodeType.USER,
                                                id: userInfo.id,
                                                name: userInfo.user.displayName,
                                                account: userInfo.user.loginName,
                                                parent_path: item.parent_path ? item.parent_path + '/' + item.name : item.name,
                                                original: userInfo,
                                            },
                                        ]
                                ), []),
                            ]
                        } catch (ex) {
                            if (ex && ex.message) {
                                Message2.info({ message: ex.message })
                            }
                        }
                    }
                }
            } else {
                addition = [...selected]
            }
            
            this.setState(({ data }) => ({
                data: uniqBy([...data, ...addition], 'id'),
                isGetUsersLoading: false,
            }), () => {
                this.departmentTreeData.cancelSelections()
            });
        }
    }

    /**
     * 删除已选部门
     * @param dep 部门
     */
    deleteSelectDep(item: FormatedNodeInfo) {
        this.setState({
            data: this.state.data.filter((value) => value.id !== item.id),
        })
    }

    /**
     * 清空已选择部门
     */
    clearSelectDep() {
        this.setState({
            data: [],
        })
    }

    /**
     * 取消本次操作
     */
    cancelAddDep() {
        this.clearSelectDep();
        this.props.onCancel();
    }

    /**
     * 确定本次操作
     */
    confirmAddDep() {
        const { data } = this.state
        const { convererOut, onConfirm } = this.props

        // 如果传递了convererOut，则对data做convererOut处理；反之，data不做任何处理。
        onConfirm(convererOut ? data.map(convererOut) : data)
    }
}