import * as React from 'react';
import { noop, has, trim } from 'lodash';
import { getUserGroups, getUserGroupById, deleteUserGroup } from '@/core/apis/console/usergroup';
import { UserManagementErrorCode } from '@/core/apis/openapiconsole/errorcode';
import { Message2, Toast } from '@/sweet-ui';
import WebComponent from '../../webcomponent';
import { ListTipStatus } from '../../ListTipComponent/helper';
import { UserGroup, DefaultPage, Limit } from '../helper';
import __ from './locale';

interface GroupGridProps extends React.Props<void> {
    /**
     * 选中项改变
     */
    onRequestSelectGroup: (selectedGroup: UserGroup.GroupInfo) => void;

    /**
     * 列表状态改变
     */
    onRequestChangeStatus: (listTipStatus: ListTipStatus) => void;
}

interface GroupGridState {
    /**
     * 列表数据
     */
    data: {
        /**
         * 用户组
         */
        groups: ReadonlyArray<UserGroup.GroupInfo>;

        /**
         * 总数
         */
        total: number;

        /**
         * 当前页码
         */
        page: number;
    };

    /**
     * 搜搜关键字
     */
    searchKey: string;

    /**
     * 列表提示
     */
    listTipStatus: ListTipStatus;

    /**
     * 选中项
     */
    selection: UserGroup.GroupInfo | null;

    /**
     * 新建/编辑的用户组信息
     */
    operatedGroup: UserGroup.GroupInfo | null;
}

export default class GroupGridBase extends WebComponent<GroupGridProps, GroupGridState> {
    static defaultProps = {
        onRequestSelectGroup: noop,
        onRequestChangeSearchKey: noop,
    }

    state = {
        data: {
            groups: [],
            total: 0,
            page: DefaultPage,
        },
        searchKey: '',
        listTipStatus: ListTipStatus.Loading,
        selection: null,
        operatedGroup: null,
    }

    /**
     * dataGrid的ref
     */
    dataGrid = {
        changeParams: noop,
    }

    /**
     * searchBox的ref
     */
    searchBox = {
        load: noop,
    }

    /**
     * 选中项的index值
     */
    selectionIndex: number = 0

    /**
     * 是否正在请求中
     */
    isRequesting = false

    componentDidMount() {
        this.searchBox.load()
    }

    /**
     * 在组件销毁后设置state，防止内存泄漏
     */
    componentWillUnmount() {
        this.setState = (state, callback) => {
            return
        }
    }

    /**
     * 获取用户组数据
     */
    protected getGroup = (keyword: string = '', param = { offset: 0 }): Promise<any> => {
        this.setState({
            listTipStatus: ListTipStatus.Loading,
        })

        let params = {
            ...param,
            limit: Limit,
        }

        keyword = trim(keyword)

        if (keyword) {
            params = {
                ...params,
                keyword,
            }
        }

        return getUserGroups(params)
    }

    /**
     * 渲染用户组数据
     */
    protected loadGroup = (gridInfo: { entries: ReadonlyArray<UserGroup.GroupInfo>; total_count: number }, selectionIndex: number = this.selectionIndex): void => {
        const { entries, total_count } = gridInfo

        this.selectionIndex = !entries.length ? 0 : selectionIndex < entries.length ? selectionIndex : entries.length - 1

        this.setState({
            data: {
                ...this.state.data,
                groups: entries,
                total: total_count,
            },
            listTipStatus: entries.length < 1 ?
                this.state.searchKey ?
                    ListTipStatus.NoSearchResults
                    : ListTipStatus.Empty
                : ListTipStatus.None,
            selection: entries.length ? entries[this.selectionIndex] : null,
        }, () => {
            this.props.onRequestSelectGroup(this.state.selection)

            this.props.onRequestChangeStatus(this.state.listTipStatus)
        })
    }

    /**
    * 加载失败
    */
    protected loadFailed = (ex: UserGroup.ErrorInfo): void => {
        this.setState({
            data: {
                ...this.state.data,
                groups: [],
            },
            selection: null,
            listTipStatus: ListTipStatus.LoadFailed,
        }, () => {
            this.props.onRequestSelectGroup(this.state.selection)

            this.props.onRequestChangeStatus(this.state.listTipStatus)
        })
    }

    /**
     * 更新列表
     */
    private updateGrid = async (param: any = { offset: 0 }, selectionIndex: number = this.selectionIndex): Promise<void> => {
        try {
            this.loadGroup(await this.getGroup(this.state.searchKey, param), selectionIndex)
        } catch (ex) {
            this.loadFailed(ex)
        }
    }

    /**
     * 改变搜索关键字
     */
    protected changeSearchKey = (searchKey: string): void => {
        if (searchKey !== this.state.searchKey) {
            this.setState({
                searchKey,
            })

            // 回到首页
            this.resetParams()
        }
    }

    /**
     * 点击新建用户组
     */
    protected addGroup = (): void => {
        this.setState({
            operatedGroup: {
                id: '',
                name: '',
                notes: '',
            },
        })
    }

    /**
     * 点击编辑用户组
     */
    protected editGroup = async (event: any, group: UserGroup.GroupInfo): Promise<void> => {
        event.stopPropagation()

        const index = this.state.data.groups.findIndex((item: UserGroup.GroupInfo) => (item.id === group.id))

        if (!this.isRequesting) {
            try {
                this.isRequesting = true

                const groupInfo = await getUserGroupById({ id: group.id })

                if (groupInfo.id === group.id) {
                    this.selectGroup({ detail: group })

                    this.setState({
                        operatedGroup: group,
                    })
                }

                this.isRequesting = false
            } catch (ex) {
                if (ex && ex.code) {
                    switch (ex.code) {
                        case UserManagementErrorCode.GroupNotFound:
                            await this.handleGroupNotExist(group, index)
                            break

                        default:
                            ex.description && Message2.info({ message: ex.description })
                    }
                }

                this.isRequesting = false
            }
        }
    }

    /**
     * 取消新建/编辑用户组
     */
    protected cancelSetGroup = async (group?: UserGroup.GroupInfo): Promise<void> => {
        this.setState({
            operatedGroup: null,
        })

        if (group) {
            if (await Message2.info({ message: __('用户组“${groupName}”已不存在', { groupName: group.name }) + __('。') })) {
                await this.updateCurrentPage()
            }
        }
    }

    /**
     * 新建/编辑成功
     */
    protected setSuccess = async (group: UserGroup.GroupInfo): Promise<void> => {
        const { operatedGroup } = this.state

        // 编辑
        if (operatedGroup.id) {
            this.setState({
                operatedGroup: null,
                selection: group,
                data: {
                    ...this.state.data,
                    groups: this.state.data.groups.map((g) => g.id === group.id ? group : g),
                },
            })

            Toast.open(__('编辑成功'))
        }
        // 新建
        else {
            this.setState({
                operatedGroup: null,
                searchKey: '',
            })

            this.resetParams()

            await this.updateGrid()

            Toast.open(__('新建成功'))
        }
    }

    /**
     * 删除用户组
     */
    protected deletegroup = async (event: any, group: UserGroup.GroupInfo): Promise<void> => {
        event.stopPropagation()

        const index = this.state.data.groups.findIndex((item: UserGroup.GroupInfo) => (item.id === group.id))

        if (!this.isRequesting) {
            try {
                this.isRequesting = true

                const groupInfo = await getUserGroupById({ id: group.id })

                if (groupInfo.id === group.id) {
                    this.selectGroup({ detail: group })

                    if (await Message2.alert({
                        message: __('删除该用户组，会导致该组的所有配置失效，确定要删除吗？'),
                        showCancelIcon: true,
                    })) {
                        await deleteUserGroup({ id: group.id })

                        await this.updateCurrentPage(index)

                        Toast.open(__('删除成功'))
                    }
                }

                this.isRequesting = false
            } catch (ex) {
                if (ex && ex.code) {
                    switch (ex.code) {
                        case UserManagementErrorCode.GroupNotFound:
                            await this.handleGroupNotExist(group, index)
                            break

                        default:
                            ex.description && Message2.info({ message: ex.description })
                    }
                }

                this.isRequesting = false
            }
        }
    }

    /**
     * 选中用户组
     */
    protected selectGroup = (event: any): void => {
        const selection = event.detail || null

        // 阻止去选中
        if (!selection && has(event, 'defaultPrevented')) {
            event.defaultPrevented = true
        } else {
            this.setState({
                selection,
            }, () => {
                this.selectionIndex = this.state.data.groups.findIndex((group) => group.id === this.state.selection.id)

                this.props.onRequestSelectGroup(this.state.selection)
            })
        }
    }

    /**
     * 手动触发页码改变
     */
    protected handlePageChange = async (page: number): Promise<void> => {
        this.setState({
            data: {
                ...this.state.data,
                page,
            },
        })

        await this.updateGrid({ offset: page ? (page - 1) * Limit : 0 }, 0)
    }

    /**
     * 重置页码
     */
    private resetParams(param: any = { page: DefaultPage }): void {
        this.setState({
            data: {
                ...this.state.data,
                page: param.page,
            },
        })

        this.selectionIndex = 0

        this.dataGrid.changeParams(param)
    }

    /**
     * 删除后刷新当前页
     */
    private async updateCurrentPage(selectionIndex: number = this.selectionIndex): Promise<void> {
        const { total, page } = this.state.data

        const totalPage = Math.ceil((total - 1) / Limit)

        if ((totalPage < page) && totalPage !== 0) {
            this.resetParams({ page: totalPage })

            await this.updateGrid({ offset: (totalPage - 1) * Limit }, 0)
        } else {
            await this.updateGrid({ offset: (page - 1) * Limit }, selectionIndex)
        }
    }

    /**
     * 处理用户组不存在时错误
     */
    public async handleGroupNotExist(group: UserGroup.GroupInfo, selectedIndex: number = this.selectionIndex): Promise<void> {
        Toast.open(__('用户组“${groupName}”已不存在', { groupName: group.name }))

        await this.updateCurrentPage(selectedIndex)
    }
}