import * as React from 'react'
import { noop, isEqual } from 'lodash'
import { Message2 } from '@/sweet-ui'
import {
    getDepartmentOfUsersCount,
    getDepartmentOfUsers,
    getAllUserCount,
    getAllUsers,
    setUserStatus,
    setUserFreezeStatus,
    editUserPriority,
} from '@/core/thrift/sharemgnt/sharemgnt'
import { manageLog, Level, ManagementOps } from '@/core/log';
import { ErrorCode } from '@/core/thrift/sharemgnt/errcode'
import { ListTipStatus } from '../../ListTipComponent/helper'
import WebComponent from '../../webcomponent'
import { EnableStatus, SpecialDep, formatDatas } from '../helper'
import __ from './locale'
import { getLevelConfig, searchUsers } from '@/core/apis/console/usermanagement';
import { getRoleTypes, UserRole } from '@/core/role/role';
import { LevelConfig } from '@/core/apis/console/usermanagement/types';
import { getAuthorizedProducts } from '@/core/apis/console/license';

interface UserGridProps extends React.Props<void> {
    /**
     * 登录用户id
     */
    userid: string;

    /**
     * 当前选中的部门节点
     */
    selectedDep: Core.ShareMgnt.ncTDepartmentInfo;

    /**
     * 是否开启冻结用户功能
     */
    freezeStatus: boolean;

    /**
     * 是否显示设置角色
     */
    isShowSetRole: boolean;

    /**
     * 获取是否显示启用/禁用用户
     */
    isShowEnableAndDisableUser: boolean;

    /**
     * 是否禁用空间站
     */
    isKjzDisabled: boolean;

    /**
     * 选中用户
     */
    onRequestSelectUsers: (selections: ReadonlyArray<Core.ShareMgnt.ncTUsrmGetUserInfo>) => any;

    /**
     * 所选部门已不存在，请求删除该部门
     */
    onRequestDelDepNode: (dep: Core.ShareMgnt.ncTDepartmentInfo) => any;
}

interface ColumnVisibility {
    [key: string]: boolean;
}

interface UserGridState {
    data: {
        /**
         * 列表用户
         */
        users: ReadonlyArray<Core.ShareMgnt.ncTUsrmGetUserInfo>;

        /**
         * 当前页码
         */
        page: number;

        /**
         * 总数
         */
        total: number;
        
        /**
         * 分页大小
         */
        pageSize: number;
    };

    /**
     * 列表提示
     */
    listTipStatus: ListTipStatus;

    /**
     * 选中的用户
     */
    selections: ReadonlyArray<Core.ShareMgnt.ncTUsrmGetUserInfo>;

    /**
     * 搜索关键字
     */
    searchKey: string;

    /**
     * 筛选条件
     */
    serarchField: string;

    /**
     * 权重值输入无效的用户ID
     */
    invalidPriorityID: string;

    /**
     * 是否显示设置用户有效期
     */
    isShowSetExpireTime: boolean;
    
    /**
     * 列可见性状态
     */
    columnVisibility: ColumnVisibility;
    
    /**
     * 上下文菜单是否可见
     */
    contextMenuVisible: boolean;
    
    /**
     * 上下文菜单位置
     */
    contextMenuPosition: { x: number; y: number };
}

export const DefaultPage = 1

export const Limit = 50

export default class UserGridBase extends WebComponent<UserGridProps, UserGridState> {
    static defaultProps = {
        selectedDep: null,
        freezeStatus: false,
        onRequestSelectUsers: noop,
        onRequestUpdateTree: noop,
    }

    state = {
        data: {
            users: [],
            page: DefaultPage,
            total: 0,
            pageSize: Limit,
        },
        listTipStatus: ListTipStatus.Loading,
        selections: [],
        searchKey: '',
        serarchField: 'name',
        invalidPriorityID: '',
        isShowSetExpireTime: false,
        columnVisibility: {} as ColumnVisibility,
        contextMenuVisible: false,
        contextMenuPosition: { x: 0, y: 0 } as { x: number; y: number },
    }
    
    /**
     * 不允许隐藏的列key列表
     */
    protected nonHideableColumns = ['displayName', 'loginName', 'code', 'directMangager', 'departmentCodes', 'position'];

    /**
     * searchBox的ref
     */
    searchBox = {
        load: noop,
    }

    /**
     * 所有密级
     */
    csfLevels: LevelConfig

    /**
     * 有效期过期的用户
     */
    expiredUserInfo: Core.ShareMgnt.ncTUsrmGetUserInfo = null

    async componentDidMount() {
        this.csfLevels = await getLevelConfig({fields: 'csf_level_enum,csf_level2_enum,show_csf_level2'})

        if (!this.props.selectedDep) {
            this.setState({
                searchKey: '',
                selections: [],
                data: {
                    ...this.state.data,
                    users: [],
                    total: 0,
                },
                listTipStatus: ListTipStatus.Empty,
            })
        } else {
            try {
                this.updateGrid({ offset: 0 })
            } catch (ex) {
                this.loadFailed(ex)
            }
        }

        document.addEventListener('click', this.closeContextMenu);
    }

    /**
     * 在组件销毁后设置state，防止内存泄漏
     */
    componentWillUnmount() {
        this.setState = (state, callback) => {
            return
        }

        document.removeEventListener('click', this.closeContextMenu);
    }

    async componentDidUpdate(prevProps) {
        if (!isEqual(this.props.selectedDep, prevProps.selectedDep)) {
            if (this.props.selectedDep) {
                this.setState({
                    searchKey: '',
                    selections: [],
                    data: {
                        ...this.state.data,
                        users: [],
                        total: 0,
                        page: 1,
                    },
                }, async () => {
                    await this.resetParams()
                   
                    await this.updateGrid({ offset: 0 })
                })
            } else {
                this.setState({
                    searchKey: '',
                    selections: [],
                    data: {
                        ...this.state.data,
                        users: [],
                        total: 0,
                    },
                    listTipStatus: ListTipStatus.Empty,
                }, async () => {
                    await this.resetParams()
                })
            }
        }
    }

    protected getRole = (roleTypes) => {
        if(roleTypes.includes(UserRole.Super)) {
            return UserRole.Super
        }else if(roleTypes.includes(UserRole.Admin)) {
            return UserRole.Admin
        }else if(roleTypes.includes(UserRole.Security)) {
            return UserRole.Security
        }else if(roleTypes.includes(UserRole.OrgManager)) {
            return UserRole.OrgManager
        }else {
            return UserRole.Super
        }
    }

    /**
     * 获取产品授权并合并到用户数据中
     */
    protected fetchProductsAndMerge = async (users: any[]): Promise<any[]> => {
        const user_ids = users.map((item) => item.id).filter(Boolean)
        if (!user_ids.length) {
            return users
        }
        
        try {
            const productData = await getAuthorizedProducts({ user_ids })
            return users.map((item) => {
                const cur = productData.find((productItem: any) => productItem.id === item.id)
                return {
                    ...item,
                    user: {
                        ...item.user,
                        products: cur?.products || [],
                    }
                }
            })
        } catch (error) {
            return users.map((item) => ({
                ...item,
                user: {
                    ...item.user,
                    products: [],
                }
            }))
        }
    }

    /**
     * 获取成员列表数据
     */
    protected getUsers = (keyword: string, param: any = { offset: 0 }, detail = this.state.serarchField): Promise<[number, ReadonlyArray<Core.ShareMgnt.ncTUsrmGetUserInfo>]> => {
        this.setState({
            listTipStatus: ListTipStatus.Loading,
        })

        const { id } = this.props.selectedDep
        const pageSize = this.state.data.pageSize || Limit

        if (!id) {
            return Promise.resolve([0, []])
        }

        // 获取用户数据的统一函数
        const fetchUserData = async (): Promise<[number, any[]]> => {
            if (keyword) {
                const roleTypes = getRoleTypes()
                const baseParams = { role: this.getRole(roleTypes), [detail || this.state.serarchField]: keyword, offset: param.offset, limit: pageSize }
                const params = id === '-2' ? baseParams : { ...baseParams, department_id: id }
                const data = await searchUsers(params)
                const users = formatDatas(data.entries)
                return [data.total_count, users]
            } else {
                if (id === SpecialDep.AllUsers) {
                    const [count, users] = await Promise.all([
                        getAllUserCount(),
                        getAllUsers([param.offset, pageSize]),
                    ])
                    return [count, users]
                } else {
                    const [count, users] = await Promise.all([
                        getDepartmentOfUsersCount([id]),
                        getDepartmentOfUsers([id, param.offset, pageSize]),
                    ])
                    return [count, users]
                }
            }
        }

        return fetchUserData().then(async ([count, users]) => {
            const usersWithProducts = await this.fetchProductsAndMerge(users)
            return [count, usersWithProducts]
        })
    }

    /**
     * 渲染成员列表数据
     */
    protected loadUsers = ([total, users]: [number, ReadonlyArray<Core.ShareMgnt.ncTUsrmGetUserInfo>]): void => {
        this.setState({
            data: {
                ...this.state.data,
                users,
                total,
            },
            selections: [],
            listTipStatus:
                users.length < 1 ?
                    this.state.searchKey ?
                        ListTipStatus.NoSearchResults
                        : ListTipStatus.Empty
                    : ListTipStatus.None,
        })
    }

    /**
     * 加载失败
     */
    protected loadFailed = (ex: any): void => {
        this.setState({
            data: {
                ...this.state.data,
                users: [],
                total: 0,
            },
            selections: [],
            listTipStatus: ListTipStatus.LoadFailed,
        })

        if (ex && ex.error && ex.error.errID && ex.error.errID === ErrorCode.DepOrOrgNotExist) {
            this.getErrMsg(ex)
        }
        if (ex && ex.code) {
            Message2.info({
                message: ex.message,
            })
        }
    }

    /**
     * 改变搜索关键字
     */
    protected changeSearchKey = (searchKey: string): void => {
        if (searchKey !== this.state.searchKey) {
            this.setState({
                searchKey,
            }, async () => {
                // 回到首页并触发搜索
                this.resetParams()
                await this.updateGrid({ offset: 0 })
            })
        }
    }

    /**
     * 改变筛选条件
     */
    protected changeFilter = async(detail) => {
        this.setState({
            serarchField: detail,
        })

        if(this.state.searchKey) {
            this.resetParams()
            this.getUsers(this.state.searchKey, { offset: this.state.data.page ? this.state.data.page * this.state.data.pageSize - this.state.data.pageSize : 0 }, detail).then(([total, users]) => {
                this.loadUsers([total, users])
            })
        }
    }

    /**
     * 设置权重值
     */
    protected setPriority = async (userInfo, priority): Promise<void> => {

        if (priority < 1 || priority > 999) {
            priority = 999

            this.setState({
                invalidPriorityID: userInfo.id,
            })
        }

        try {
            await editUserPriority([userInfo.id, parseInt(priority)])

            this.setState({
                invalidPriorityID: '',
            })

            this.updateUserInfo({ ...userInfo, user: { ...userInfo.user, priority } })
        } catch (ex) {
            this.getErrMsg(ex, userInfo)
        }
    }

    /**
     * 改变用户状态
     */
    protected changeStatus = async (userInfo: Core.ShareMgnt.ncTUsrmGetUserInfo, status: boolean): Promise<void> => {
        try {
            if (userInfo.id === this.props.userid) {
                await Message2.info({ message: status ? __('您无法启用自身账号。') : __('您无法禁用自身账号。') })
            } else {
                await setUserStatus(userInfo.id, status)

                this.updateUserInfo({ ...userInfo, user: { ...userInfo.user, status: status ? EnableStatus.Enabled : EnableStatus.Disabled } }, true)

                await manageLog(
                    ManagementOps.SET,
                    __(`${status ? '启用' : '禁用'}` + ' 用户“${name}” 成功', { name: `${userInfo.user.displayName}(${userInfo.user.loginName})` }),
                    '',
                    Level.WARN,
                )
            }
        } catch (ex) {
            this.getErrMsg(ex, userInfo)
        }
    }

    /**
     * 改变用户冻结状态
     */
    protected changeFreezeStatus = async (userInfo: Core.ShareMgnt.ncTUsrmGetUserInfo, freezeStatus: boolean): Promise<void> => {
        try {
            if (userInfo.id === this.props.userid) {
                await Message2.info({ message: freezeStatus ? __('您无法冻结自身账号。') : __('您无法解冻自身账号。') })
            } else {
                await setUserFreezeStatus(userInfo.id, freezeStatus)

                this.updateUserInfo({ ...userInfo, user: { ...userInfo.user, freezeStatus } })

                await manageLog(
                    ManagementOps.SET,
                    __(`${freezeStatus ? '冻结 ' : '解冻 '}` + '用户“${name}” 成功', { name: `${userInfo.user.displayName}(${userInfo.user.loginName})` }),
                    '',
                    Level.INFO,
                )
            }
        } catch (ex) {
            this.getErrMsg(ex, userInfo)
        }
    }

    /**
     * 选中用户
     */
    protected changeSelection = (selections: ReadonlyArray<Core.ShareMgnt.ncTUsrmGetUserInfo>): void => {
        this.setState({
            selections,
        })

        this.props.onRequestSelectUsers(selections)
    }

    /**
     * 手动触发页码改变
     */
    protected handlePageChange = async (page: number, pageSize: number, offset: number): Promise<void> => {
        this.setState({
            data: {
                ...this.state.data,
                page,
                pageSize,
            },
        }, async () => {
            await this.updateGrid({ offset })
        })
    }
    
    /**
     * 手动触发分页大小改变
     */
    protected handlePageSizeChange = async (pageSize: number): Promise<void> => {
        this.setState({
            data: {
                ...this.state.data,
                pageSize,
                page: DefaultPage,
            },
        }, async () => {
            await this.updateGrid({ offset: 0 })
        })
    }

    /**
     * 添加用户
     */
    public addUser = async (userInfo: Core.ShareMgnt.ncTUsrmGetUserInfo): Promise<void> => {
        const { data: { page, pageSize = Limit }, searchKey } = this.state

        if (page !== 1 || searchKey) {
            const [total, users] = await this.getUsers('')

            this.resetParams()

            this.setState({
                searchKey: '',
                data: {
                    ...this.state.data,
                    total,
                    users: [userInfo, ...users.filter((user) => user.id !== userInfo.id).slice(0, pageSize - 1)],
                    page: 1,
                },
                selections: [],
                listTipStatus: ListTipStatus.None,
            })
        } else {
            const usersWithProducts = await this.fetchProductsAndMerge([userInfo])
            const newUserWithProducts = usersWithProducts[0]
            
            this.setState({
                data: {
                    ...this.state.data,
                    users: [newUserWithProducts, ...this.state.data.users.slice(0, pageSize - 1)],
                    total: this.state.data.total + 1,
                    page: 1,
                },
                selections: [],
                listTipStatus: ListTipStatus.None,
            })
        }
    }

    /**
     * 更新指定用户信息
     */
    public updateUserInfo = async (userInfo: Core.ShareMgnt.ncTUsrmGetUserInfo, isNeedUpdateProducts: boolean = false): Promise<void> => {
        const { data, selections } = this.state
        // 如果需要更新产品信息，则先拉取最新产品数据再合并
        if (selections.length === 1 && isNeedUpdateProducts) {
            const usersWithProducts = await this.fetchProductsAndMerge([userInfo])
            userInfo = usersWithProducts[0]
        }

        this.setState({
            data: {
                ...data,
                users: data.users.map((user) => user.id === userInfo.id ? userInfo : user),
            },
            selections: selections.map((user) => user.id === userInfo.id ? userInfo : user),
        })
    }

    /**
     * 回到首页
     */
    public backToFirst = (): void => {
        this.resetParams()

        this.setState({
            searchKey: '',
        }, async () => {
            await this.updateGrid({ offset: 0 })
        })
    }

    /**
     * 刷新当前页
     */
    public async updateCurrentPage(delNumber: number = 0): Promise<void> {
        const { total, page, pageSize } = this.state.data

        const totalPage = Math.ceil((total - delNumber) / pageSize)

        if ((totalPage < page) && totalPage !== 0) {
            this.resetParams({ page: totalPage })

            await this.updateGrid({ offset: (totalPage - 1) * pageSize })
        } else {
            await this.updateGrid({ offset: (page - 1) * pageSize })
        }
    }

    /**
     * 更新列表
     */
    private updateGrid = async (param: any = {}): Promise<void> => {
        try {
            if (!this.csfLevels?.csf_level_enum?.length || !this.csfLevels?.csf_level2_enum?.length) {
                this.csfLevels = await getLevelConfig({ fields: 'csf_level_enum,csf_level2_enum,show_csf_level2' })
            }

            const offset = param.offset !== undefined ? param.offset : (this.state.data.page ? this.state.data.page * this.state.data.pageSize - this.state.data.pageSize : 0)
            this.loadUsers(await this.getUsers(this.state.searchKey, { ...param, offset }))
        } catch (ex) {
            this.loadFailed(ex)
        }
    }

    /**
     * 重置页码
     */
    private resetParams(param: any = { page: DefaultPage }): void {
        this.setState({
            data: {
                ...this.state.data,
                page: param.page,
                pageSize: this.state.data.pageSize || Limit,
            },
        })
    }

    /**
     * 处理上下文菜单显示
     */
    protected handleContextMenu = (e: React.MouseEvent, column: any) => {
        e.preventDefault();
        e.stopPropagation();
        
        // 计算菜单位置，确保不会超出窗口边界
        const menuWidth = 200;
        const menuHeight = 300;
        const windowWidth = window.innerWidth;
        const windowHeight = window.innerHeight;
        
        let x = e.clientX;
        let y = e.clientY;
        
        // 如果菜单会超出右边边界，调整x坐标
        if (x + menuWidth > windowWidth) {
            x = windowWidth - menuWidth;
        }
        
        // 如果菜单会超出底部边界，调整y坐标
        if (y + menuHeight > windowHeight) {
            y = windowHeight - menuHeight;
        }
        
        this.setState({
            contextMenuVisible: true,
            contextMenuPosition: { x, y },
        });
    };
    
    /**
     * 关闭上下文菜单
     */
    protected closeContextMenu = () => {
        this.setState({ contextMenuVisible: false });
    };
    
    /**
     * 切换列可见性
     */
    protected toggleColumnVisibility = (columnKey: string) => {
        this.setState(prevState => {
            // 获取当前列的可见性状态，默认为true（可见）
            const currentVisibility = prevState.columnVisibility[columnKey] !== false;
            // 切换可见性状态
            const newVisibility = !currentVisibility;
            // 更新状态
            const newColumnVisibility = {
                ...prevState.columnVisibility,
                [columnKey]: newVisibility,
            };
            return { columnVisibility: newColumnVisibility };
        });
    };
    
    /**
     * 全选所有列
     */
    protected selectAllColumns = (columns: any[]) => {
        this.setState(prevState => {
            const newColumnVisibility = { ...prevState.columnVisibility };
            
            // 获取所有可隐藏的列（有key且不在nonHideableColumns中的列）
            const hideableColumns = columns.filter(col => col.key && !this.nonHideableColumns.includes(col.key));
            
            // 将所有可隐藏的列设置为可见
            hideableColumns.forEach(col => {
                newColumnVisibility[col.key] = true;
            });
            
            return { columnVisibility: newColumnVisibility };
        });
    };
    
    /**
     * 取消选择所有可隐藏的列
     */
    protected unselectAllColumns = (columns: any[]) => {
        this.setState(prevState => {
            const newColumnVisibility = { ...prevState.columnVisibility };
            
            // 获取所有可隐藏的列（有key且不在nonHideableColumns中的列）
            const hideableColumns = columns.filter(col => col.key && !this.nonHideableColumns.includes(col.key));
            
            // 将所有可隐藏的列设置为不可见
            hideableColumns.forEach(col => {
                newColumnVisibility[col.key] = false;
            });
            
            return { columnVisibility: newColumnVisibility };
        });
    };
    
    /**
     * 检查是否所有可隐藏的列都被选中
     */
    protected isAllColumnsSelected = (columns: any[]) => {
        // 过滤出所有可隐藏的列（有key且不在nonHideableColumns中的列）
        const hideableColumns = columns.filter(col => col.key && !this.nonHideableColumns.includes(col.key));
        
        // 如果没有可隐藏的列，则认为全选
        if (hideableColumns.length === 0) return true;
        
        // 检查是否所有可隐藏的列都被选中
        return hideableColumns.every(col => this.state.columnVisibility[col.key] !== false);
    };
    
    /**
     * 过滤可见列
     */
    protected filterVisibleColumns = (columns: any[]) => {
        return columns.filter(col => {
            // 如果列没有key或key在非隐藏列表中，始终可见
            if (!col.key || this.nonHideableColumns.includes(col.key)) return true;
            // 其他列根据visibility状态决定是否可见
            return this.state.columnVisibility[col.key] !== false;
        });
    };
    
    /**
     * 获取错误信息
     */
    private async getErrMsg({ error }, userInfo: Core.ShareMgnt.ncTUsrmGetUserInfo | null = null): Promise<void> {
        if (error && error.errID) {
            const { selectedDep } = this.props

            switch (error.errID) {
                case ErrorCode.DepOrOrgNotExist:
                    if (await Message2.info({
                        message: __(`${selectedDep.is_root || selectedDep.id === SpecialDep.Unassigned || selectedDep === SpecialDep.AllUsers ? '组织' : '部门'}` + '“${name}”不存在。', { name: selectedDep.name }),
                    })) {
                        this.props.onRequestDelDepNode(this.props.selectedDep)
                    }

                    break

                case ErrorCode.EnableUserCountOverproof:
                    Message2.info({
                        message: __('启用用户 ${userName} 失败。启用用户数已达用户许可总数的上限。', { userName: (userInfo && userInfo.user.displayName) || '' }),
                    })
                    break

                case ErrorCode.CountPassDue:
                    if (await Message2.info({
                        message: __('该用户账号已过期，是否重新设置有效期限？'),
                    })) {
                        this.expiredUserInfo = userInfo
                        this.setState({
                            isShowSetExpireTime: true,
                        })
                    }
                    break

                default:
                    Message2.info({
                        message: error.errMsg,
                    })
                    break
            }
        }
    }
}