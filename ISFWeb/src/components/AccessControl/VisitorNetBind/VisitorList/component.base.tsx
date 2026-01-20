import * as React from 'react';
import { map, noop, isEqual, filter } from 'lodash'
import * as PropTypes from 'prop-types';
import { getAccessorsByNetwork, addAccessorsByNetwork, deleteAccessorByNetwork } from '@/core/apis/console/networkRestriction'
import { ErrorCode } from '@/core/apis/openapiconsole/errorcode';
import { manageLog, Level, ManagementOps } from '@/core/log/log';
import { Message } from '@/sweet-ui';
import { NodeType } from '@/core/organization';
import { ListTipStatus } from '../../../ListTipComponent/helper';
import WebComponent from '../../../webcomponent';
import { PageSize, NetType, DefaultPage, getNetBindErrorMessage, NetInfo, Visitor } from '../helper';
import __ from './locale';

interface VisitorListProps extends React.Props<void> {
    /**
     * 访问者列表是否可用
     */
    isEnabled: boolean;

    /**
     * 选中的网段
     */
    selectedNet: NetInfo;
}

interface VisitorListState {
    /**
     * 搜索框输入的值
     */
    searchKey: string;

    /**
     * 访问者列表
     */
    visitorList: ReadonlyArray<Visitor>;

    /**
     * 访问者总数
     */
    visitorsCount: number;

    /**
     * 是否添加访问者
     */
    isAddVisitor: boolean;

    /**
     * 列表页码
     */
    page: number;

    /**
     * 列表提示状态
     */
    listTipStatus: ListTipStatus;
}

export default class VisitorListBase extends WebComponent<VisitorListProps, VisitorListState> {

    static defaultProps = {
        selectedNet: null,
    }

    static contextTypes = {
        toast: PropTypes.func,
    }
    state: VisitorListState = {
        searchKey: '',
        visitorList: [],
        visitorsCount: 0,
        isAddVisitor: false,
        page: DefaultPage,
        listTipStatus: ListTipStatus.Loading,
    }

    currentPage: number = DefaultPage;

    dataGrid = {
        changeParams: noop,
    }

    componentDidMount() {
        this.setState({
            listTipStatus: ListTipStatus.Empty,
        })
    }

    componentDidUpdate({ selectedNet }) {
        // 如果 props.selectedNet.id 没有发生变化，不用重新获取访问这里列表了
        if (selectedNet && this.props.selectedNet && selectedNet.id === this.props.selectedNet.id) {
            return
        }

        if (!isEqual(selectedNet, this.props.selectedNet)) {
            this.setState({
                searchKey: '',
            })
            this.getVisitorList({ selectedNet: this.props.selectedNet });
        }
    }

    /**
     * 在组件销毁后设置state，防止内存泄漏
     */
    componentWillUnmount() {
        this.setState = (state, callback) => {
            return;
        };
    }

    /**
     * 获取列表数据
     * @param page 指定页列表参数
     */
    protected async getVisitorList({ selectedNet = null, key = '', page = DefaultPage } = {}) {

        try {
            this.setState({
                page,
            })
            this.resetParams({ page });
            if (selectedNet) {
                this.setState({
                    listTipStatus: ListTipStatus.Loading,
                })
                // 调用列表接口
                const results = await getAccessorsByNetwork(
                    {
                        id: selectedNet.id,
                        key_word: key,
                        offset: (page - 1) * PageSize,
                        limit: PageSize,
                    },
                )
                // 获取访问者列表
                const visitorList = this.formateData(results.data);

                this.setState({
                    visitorList: visitorList,
                    visitorsCount: results.count,
                    listTipStatus: visitorList.length < 1 ?
                        (key ? ListTipStatus.NoSearchResults : ListTipStatus.Empty)
                        : ListTipStatus.None,
                })
                this.currentPage = page;
            } else {
                this.setState({
                    visitorList: [],
                    visitorsCount: 0,
                    listTipStatus: ListTipStatus.Empty,
                })
                this.currentPage = 0;
            }

        } catch (error) {
            this.setState({
                visitorList: [],
                visitorsCount: 0,
                listTipStatus: ListTipStatus.LoadFailed,
            })
            // 网段不存在
            if (error.code === ErrorCode.ResourceInaccessibleByPolicy && error.detail.notfound_params.indexOf('id') >= 0) {
                getNetBindErrorMessage(error)
            }
        }
    }

    /**
     * 打开添加访问者窗口
     */
    protected openVisitorsDialog() {
        this.setState({
            isAddVisitor: true,
        })
    }

    /**
     * 保存访问者
     * @param list 保存的访问者
     */
    protected async saveVisitors(list) {
        try {

            const { selectedNet } = this.props;
            const visitorsParam = this.formateSaveParam(list);

            const results = await addAccessorsByNetwork({
                id: selectedNet.id,
                accessorsList: visitorsParam,
            })

            let noExitVisitors = filter(results, (result) => {
                if (result.body.code === ErrorCode.ResourceInaccessibleByPolicy && result.body.detail.notfound_params.indexOf('id') >= 0) {
                    // 网段不存在
                    this.setState({
                        isAddVisitor: false,
                    })
                    throw (result.body);
                }
                // 访问者不存在
                return result.body.code === ErrorCode.ResourceInaccessibleByPolicy && result.body.detail.notfound_params.indexOf('accessor_id') >= 0;
            })

            this.setState({
                isAddVisitor: false,
                searchKey: '',
            }, async () => {
                this.getVisitorList({ selectedNet });

                const existVisitors = filter(list, (item) => {
                    return noExitVisitors.findIndex((cur) => cur.id === item.id) < 0;
                })

                // 记录日志
                const names = map(existVisitors, (item) => item.name);

                if(names.length) {
                    manageLog(
                        ManagementOps.SET,
                        __('IP网段绑定访问者“${visitors}”成功', { visitors: names.join('，') }),
                        selectedNet.id === 'public-net' ?
                            __('绑定的网段为“所有外网网段”')
                            : (
                                selectedNet.name ?
                                    __('绑定的网段为“${segment}”；网段名称“${name}”', {
                                        segment: selectedNet.netType === NetType.Range ?
                                            selectedNet.originIP + '-' + selectedNet.endIP
                                            : selectedNet.ip + '/' + selectedNet.mask,
                                        name: selectedNet.name,
                                    })
                                    : __('绑定的网段为“${segment}”', {
                                        segment: selectedNet.netType === NetType.Range ?
                                            selectedNet.originIP + '-' + selectedNet.endIP
                                            : selectedNet.ip + '/' + selectedNet.mask,
                                    })
                            ),
                        Level.WARN,
                    )
                }

                if (noExitVisitors.length) {
                    await Message.alert({
                        message: __('访问者“${names}”添加失败，该用户已不存在。',
                            {
                                names: map(noExitVisitors, ({ id }) => {
                                    return filter(list, (item) => item.id === id)[0].name
                                }).join('、'),
                            }),
                    })
                    this.context.toast(__('操作完成，本次添加${successCount}个访问者', { successCount: list.length - noExitVisitors.length }));
                } else {
                    this.context.toast(__('操作成功'));
                }
            })

        } catch (error) {
            getNetBindErrorMessage(error)
        }

    }

    /**
     * 删除访问者
     * @param visitor 删除的数据
     */
    protected async deleteVisitor(visitor) {
        try {

            const { selectedNet } = this.props;
            const { visitorList, visitorsCount, searchKey } = this.state;
            const afterDeleteList = visitorList.filter((item) => item.id !== visitor.id);

            // 删除访问者
            const results = await deleteAccessorByNetwork({
                id: selectedNet.id,
                accessor_id: visitor.id,
            });

            // 网段不存在
            if (results[0].body.code === ErrorCode.ResourceInaccessibleByPolicy && results[0].body.detail.notfound_params.indexOf('id') >= 0) {
                throw (results[0].body);
            }

            // 当前页是最后一页
            if (this.currentPage === parseInt(String((visitorsCount - 1) / PageSize)) + 1) {
                if (visitorList.length > 1) {
                    // 当前页超过一条
                    this.setState({
                        visitorList: afterDeleteList,
                        visitorsCount: visitorsCount - 1,
                    })

                } else if (this.currentPage === 1) {
                    // 当前页为第一页，且只有一条数据
                    this.setState({
                        visitorList: afterDeleteList,
                        visitorsCount: 0,
                        listTipStatus: ListTipStatus.Empty,
                    })

                } else {
                    // 当前页只有一条且不为第一页，获取上一页数据
                    this.getVisitorList({ selectedNet, key: searchKey, page: this.currentPage - 1 })
                }
            } else {
                // 当前页不是最后一页,重新获取该页数据
                this.getVisitorList({ selectedNet, key: searchKey, page: this.currentPage })
            }
            // 记录日志
            manageLog(
                ManagementOps.SET,
                __('解除绑定 IP网段与访问者“${visitor}”的绑定 成功', { visitor: visitor.name }),
                selectedNet.id === 'public-net' ?
                    __('解除绑定的网段为“所有外网网段”')
                    : (
                        selectedNet.name ?
                            __('解除绑定的网段为“${segment}”；网段名称“${name}”', {
                                segment: selectedNet.netType === NetType.Range ?
                                    selectedNet.originIP + '-' + selectedNet.endIP
                                    : selectedNet.ip + '/' + selectedNet.mask,
                                name: selectedNet.name,
                            })
                            : __('解除绑定的网段为“${segment}”', {
                                segment: selectedNet.netType === NetType.Range ?
                                    selectedNet.originIP + '-' + selectedNet.endIP
                                    : selectedNet.ip + '/' + selectedNet.mask,
                            })
                    ),
                Level.WARN,
            )

        } catch (error) {
            await getNetBindErrorMessage(error);
        }
    }

    /**
     * 搜索框输入值
     * @param searchKey 搜索框输入值
     */
    protected changeSearchKey(searchKey: string) {
        this.setState({
            searchKey,
        }, () => {
            this.resetParams();
        })
    }

    /**
     * 加载搜索结果
     * @param data 搜索结果
     */
    protected loadSearchResult(data: ReadonlyArray<object>) {
        const { selectedNet } = this.props;
        this.getVisitorList({ selectedNet, key: this.state.searchKey })
    }

    /**
     * 关闭添加和编辑网段弹框
     */
    protected closeVisitorsDialog() {
        this.setState({
            isAddVisitor: false,
        })
    }

    /**
     * 改变表格参数
     */
    protected resetParams(params: { page: number } = { page: DefaultPage }) {
        this.dataGrid.changeParams(params);
    }

    /**
     * 手动触发页码改变
     * @param page 页码
     */
    protected handlePageChange(page: number) {
        const { selectedNet } = this.props;

        this.setState({
            page,
        }, () => {
            this.getVisitorList({ selectedNet, key: this.state.searchKey, page })
        })
    }

    /**
     * 格式化数据
     * @param data 需要格式化的数据
     */
    protected formateData(data: ReadonlyArray<any>) {
        return map(data, ({ accessor_id: id, accessor_name: name, accessor_type: type }) => {

            return { id, name, type: type === 'department' ? NodeType.DEPARTMENT : NodeType.USER };

        })
    }

    /**
     * 格式化访问者保存数据
     * @param data 需要格式化的数据
     */
    protected formateSaveParam(data: ReadonlyArray<any>) {
        return map(data, ({ id: accessor_id, type: accessor_type }) => {
            return {
                accessor_id,
                accessor_type:
                    accessor_type === NodeType.DEPARTMENT || accessor_type === NodeType.ORGANIZATION ?
                        'department' :
                        'user',
            }
        })
    }
}