import { noop, trim } from 'lodash';
import { Message2, Toast } from '@/sweet-ui';
import { getGroupMembers, getGroupMembersByUserMatch, addGroupMembers, deleteGroupMembers } from '@/core/apis/console/usergroup';
import { UserManagementErrorCode } from '@/core/apis/openapiconsole/errorcode';
import { ListTipStatus } from '../../ListTipComponent/helper';
import WebComponent from '../../webcomponent';
import { UserGroup, DefaultPage, Limit } from '../helper';
import { MemberGridProps, MemberGridState, SearchType } from './type'
import __ from './locale';

export default class MemberGridBase extends WebComponent<MemberGridProps, MemberGridState> {
    static defaultProps = {
        onRequestUpdateGroup: noop,
        selectedGroup: null,
    }

    state = {
        data: {
            members: [],
            total: 0,
            page: 1,
        },
        searchKey: '',
        listTipStatus: ListTipStatus.Loading,
        selections: [],
        isShowAdd: false,
        searchBoxIsOnFocus: false,
        searchType: SearchType.UserGroupMemberName,
    }

    /**
     * 是否正在请求中
     */
    isRequesting = false

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

    componentDidMount() {
        if (this.searchBox && this.props.selectedGroup) {
            this.searchBox.load()
        }
    }

    /**
     * 在组件销毁后设置state，防止内存泄漏
     */
    componentWillUnmount() {
        this.setState = (state, callback) => {
            return
        }
    }

    async componentDidUpdate(prevProps) {
        if (this.props.selectedGroup !== prevProps.selectedGroup || this.props.groupStatus !== prevProps.groupStatus) {
            if (this.props.selectedGroup) {
                this.setState({
                    searchKey: '',
                    selections: [],
                    data: {
                        members: [],
                        total: 0,
                        page: 1,
                    },
                }, async () => {
                    await this.resetParams()

                    this.searchBox.load()
                })
            } else {
                this.setState({
                    searchKey: '',
                    selections: [],
                    data: {
                        ...this.state.data,
                        members: [],
                        total: 0,
                    },
                    listTipStatus: this.props.groupStatus,
                }, async () => {
                    await this.resetParams()
                })
            }
        }
    }

    /**
     * 获取成员列表数据
     */
    protected getMembers = async (keyword: string, param = { offset: 0 }): Promise<any> => {
        this.setState({
            listTipStatus: ListTipStatus.Loading,
        })

        if (this.props.selectedGroup) {
            const id = this.props.selectedGroup.id;
            const { searchType } = this.state
            keyword = trim(keyword)

            if (searchType === SearchType.UserDisplayNmae && !!keyword) {
                const { entries, total_count } = await getGroupMembersByUserMatch({ group_id: id, key: keyword, ...param, limit: Limit })

                return {
                    entries,
                    total_count,
                    id,
                }
            } else {
                let params = {
                    ...param,
                    limit: Limit,
                }

                if (keyword) {
                    params = {
                        ...params,
                        keyword,
                    }
                }

                return {
                    ...(await getGroupMembers(id, params)),
                    id,
                }
            }
        }

        return { entries: [], total_count: 0, id: '' }
    }

    /**
     * 渲染成员列表数据
     */
    protected loadMembers = (gridInfo: { entries: ReadonlyArray<UserGroup.MemberInfo>; total_count: number; id: string }): void => {
        const { entries, total_count, id } = gridInfo

        // 只渲染最后一次选择的用户组成员
        if (id === this.props.selectedGroup.id) {
            this.setState({
                data: {
                    ...this.state.data,
                    members: entries,
                    total: total_count,
                },
                selections: [],
                listTipStatus:
                    entries.length < 1 ?
                        this.state.searchKey ?
                            ListTipStatus.NoSearchResults
                            : ListTipStatus.Empty
                        : ListTipStatus.None,
            })
        }
    }

    /**
     * 加载失败
     */
    protected loadFailed = (ex: UserGroup.ErrorInfo): void => {
        if (ex && ex.code && ex.code === UserManagementErrorCode.GroupNotFound) {
            this.props.onRequestUpdateGroup()
        } else {
            this.setState({
                data: {
                    ...this.state.data,
                    members: [],
                    total: 0,
                },
                selections: [],
                listTipStatus: ListTipStatus.LoadFailed,
            })
        }
    }

    /**
     * 更新列表
     */
    private updateGrid = async (param = { offset: 0 }): Promise<void> => {
        try {
            this.loadMembers(await this.getMembers(this.state.searchKey, param))
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
     * 改变搜索类型
     */
    protected changeSearchType = (searchType: SearchType): void => {
        if (searchType !== this.state.searchType) {
            this.setState({
                searchType,
            }, () => {
                // 搜索关键字存在时，重新搜索
                if (!!this.state.searchKey) {
                    // 回到首页
                    this.resetParams()
                    this.updateGrid()
                }
            })
        }
    }

    /**
     * 是否显示添加成员的弹窗
     */
    protected changePickerStatus = (isShowAdd: boolean) => {
        this.setState({
            isShowAdd,
        })
    }

    /**
     * 确定添加成员
     */
    protected confirmAdd = async (members: ReadonlyArray<UserGroup.MemberInfo>): Promise<void> => {
        if (!this.isRequesting) {
            try {
                this.isRequesting = true

                await addGroupMembers(this.props.selectedGroup.id, { members })

                Toast.open(__('添加成员成功'))

                this.setState({
                    searchKey: '',
                    selections: [],
                    isShowAdd: false,
                })

                this.resetParams()

                await this.updateGrid()

                this.isRequesting = false
            } catch (ex) {
                if (ex && ex.code) {
                    switch (ex.code) {
                        case UserManagementErrorCode.GroupNotFound:
                            {
                                this.changePickerStatus(false)
                                this.props.onRequestUpdateGroup()
                            }
                            break

                        case UserManagementErrorCode.GroupMemberNotExisted:
                        case UserManagementErrorCode.DepartmentNotExisted:
                            Message2.info({
                                message: __(
                                    '用户组成员“${memberName}”已不存在。',
                                    { memberName: ex.detail.ids.map((id) => members.find((m) => m.id === id).name).join(', ') },
                                ),
                            })
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
     * 删除成员
     */
    protected deleteMembers = async (): Promise<void> => {
        const { selections } = this.state

        try {
            await deleteGroupMembers(this.props.selectedGroup.id, { members: selections })

            Toast.open(__('删除成员成功'))

            await this.updateCurrentPage(selections.length)

            this.setState({
                selections: [],
            })
        } catch (ex) {
            if (ex && ex.code) {
                switch (ex.code) {
                    case UserManagementErrorCode.GroupNotFound:
                        this.props.onRequestUpdateGroup()
                        break

                    default:
                        ex.description && Message2.info({ message: ex.description })
                }
            }
        }
    }

    /**
     * 选中成员
     */
    protected changeSelection = (selections: ReadonlyArray<UserGroup.MemberInfo>): void => {
        this.setState({
            selections,
        })
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

        await this.updateGrid({ offset: page ? page * Limit - Limit : 0 })
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

        this.dataGrid.changeParams(param)
    }

    /**
     * 删除后刷新当前页
     */
    private async updateCurrentPage(delNumber = 1): Promise<void> {
        const { total, page } = this.state.data

        const totalPage = Math.ceil((total - delNumber) / Limit)

        if ((totalPage < page) && totalPage !== 0) {
            this.resetParams({ page: totalPage })

            await this.updateGrid({ offset: (totalPage - 1) * Limit })
        } else {
            await this.updateGrid({ offset: (page - 1) * Limit })
        }
    }
}