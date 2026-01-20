import * as React from 'react';
import { noop, findIndex, map } from 'lodash';
import { getNetworkList, addNetwork, getNetworkInfo, editNetwork, deleteNetwork } from '@/core/apis/console/networkRestriction'
import { manageLog, Level, ManagementOps } from '@/core/log/log';
import { PublicErrorCode } from '@/core/apis/openapiconsole/errorcode'
import { Message } from '@/sweet-ui';
import { ListTipStatus } from '../../../ListTipComponent/helper';
import WebComponent from '../../../webcomponent';
import { IpVersion, OperateType } from '../../NetSegment/helper';
import { PageSize, NetType, getNetBindErrorMessage, getIdFromLocaltion, DefaultPage, NetInfo } from '../helper';
import __ from './locale';

interface NetListProps extends React.Props<void> {
    /**
     * 列表是否灰化
     */
    isEnabled: boolean;

    /**
     * 服务状态
     */
    serverStatus: boolean;

    /**
     * 选中列表的某个网段
     */
    onSelectNet: (info: NetInfo) => void;

    /**
     * 设置访问者列表禁用状态
     */
    onDisableVisitorList: (status: boolean) => void;
}

interface NetListState {
    /**
     * 网段列表
     */
    netList: ReadonlyArray<NetInfo>;

    /**
     * 搜索框输入值
     */
    searchKey: string;

    /**
     * 列表页数
     */
    page: number;

    /**
     * 是否编辑或添加网段
     */
    isEditNet: boolean;

    /**
     * 网段总数
     */
    netsCount: number;

    /**
     * 正在编辑/添加的网段
     */
    editingNet: NetInfo;

    /**
     * 操作类型
     */
    operateType: OperateType;

    /**
     * 选中的网段
     */
    selection: NetInfo;

    /**
     * 列表提示状态
     */
    listTipStatus: ListTipStatus;
}

const initialNet: NetInfo = {
    id: '',
    name: '',
    netType: NetType.Range,
    originIP: '',
    endIP: '',
    ip: '',
    mask: '',
    ipVersion: IpVersion.Ipv4,
}

export default class NetListBase extends WebComponent<NetListProps, NetListState> {

    static defaultProps = {
        isEnabled: false,
        serverStatus: false,
        onSelectNet: noop,
        onDisableVisitorList: noop,
    }

    state: NetListState = {
        netList: [],
        searchKey: '',
        page: DefaultPage,
        isEditNet: false,
        netsCount: 0,
        editingNet: { ...initialNet },
        operateType: OperateType.Add,
        selection: null,
        listTipStatus: ListTipStatus.Loading,
    }

    currentPage: number = DefaultPage;

    dataGrid = {
        changeParams: noop,
    };

    componentDidMount() {
        this.getNetList();
    }

    componentDidUpdate(prevProps, prevState) {
        const { isEnabled, serverStatus } = this.props;
        if (prevProps.isEnabled !== isEnabled || prevProps.serverStatus !== serverStatus) {
            this.setState({
                searchKey: '',
            })
            this.getNetList();
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
     * @param key 搜索关键在
     * @param page 指定页
     * @param selectedIndex 选中的下标项
     */
    protected async getNetList({ key = '', page = DefaultPage } = {}, { selectedIndex = 0 } = {}) {
        try {
            if (this.props.serverStatus) {
                this.setState({
                    listTipStatus: ListTipStatus.Loading,
                })
                const { count, data } = await getNetworkList(
                    {
                        key_word: key,
                        offset: (page - 1) * PageSize,
                        limit: PageSize,
                    },
                )
                const netList = this.formateGettedNets(data);

                this.setState({
                    page,
                    netList: netList,
                    netsCount: count,
                    selection: netList[selectedIndex],
                    listTipStatus: netList.length < 1 ?
                        (key ? ListTipStatus.NoSearchResults : ListTipStatus.Empty)
                        : ListTipStatus.None,
                }, () => {
                    const { listTipStatus } = this.state;
                    if (listTipStatus === ListTipStatus.NoSearchResults || listTipStatus === ListTipStatus.Empty) {
                        this.props.onDisableVisitorList(true);
                    } else {
                        this.props.onDisableVisitorList(false);
                    }
                    this.props.onSelectNet(netList[selectedIndex]);
                })
                this.resetParams({ page });
                this.currentPage = page;

            } else {
                this.setState({
                    listTipStatus: ListTipStatus.Empty,
                })
            }
        } catch (error) {
            this.setState({
                listTipStatus: ListTipStatus.LoadFailed,
            })
            this.props.onDisableVisitorList(true);
            getNetBindErrorMessage(error);
        }
    }
    /**
     * 新增网段，打开新增窗口
     * @param net 编辑的网段
     */
    protected updateNet(operateType: OperateType, net?: NetInfo) {
        this.setState({
            operateType,
            isEditNet: true,
            editingNet: net ? net : { ...initialNet },
        }, () => {
            const { selection } = this.state;
            this.setState({
                selection: net ? net : selection,
            })
            this.props.onSelectNet(net ? net : selection);
        })
    }

    /**
   * 删除网段
   * @param net 删除的数据
   */
    protected async deleteNet(e: Event, net: NetInfo) {
        e.stopPropagation();
        this.setState({
            selection: net,
        })

        if (await Message.confirm({ message: __('删除该网段后，其绑定的访问者将不能继续在此网段进行登录，确认要执行此操作吗？') })) {

            try {
                const { netList, searchKey, netsCount } = this.state;
                const index = findIndex(netList, (item) => (item.id === net.id));
                // 不存在的网段直接删除
                try {
                    await deleteNetwork({ id: net.id });
                } catch (error) {
                    if (error.code !== PublicErrorCode.NotFound) {
                        throw (error);
                    }
                }
                // 当前页为最后一页
                if (this.currentPage === parseInt(String((netsCount - 1) / PageSize)) + 1) {
                    let newNetList = netList.filter((item) => item.id !== net.id);
                    if (netList.length === 1) {
                        if (this.currentPage === 1) {
                            // 当前页为首页，且只有一条数据
                            this.setState({
                                netList: [],
                                netsCount: 0,
                                selection: null,
                                listTipStatus: this.state.searchKey ? ListTipStatus.NoSearchResults : ListTipStatus.Empty,
                            }, () => {
                                this.props.onSelectNet(null)
                            })
                        } else {
                            // 当前页不是首页，获取上一页数据，且选中上一页第一条数据
                            await this.getNetList({ key: searchKey, page: this.currentPage - 1 })
                        }

                    } else if (index === netList.length - 1) {
                        // 当前页超过一条数据，删除当前页最后一条数据，选择该条数据的上一条
                        this.setState({
                            netList: newNetList,
                            netsCount: netsCount - 1,
                            selection: newNetList[index - 1],
                        }, () => {
                            this.props.onSelectNet(this.state.selection)
                        })

                    } else {
                        // 当前页超过一条数据，删除的数据不是当前页最后一条时，选择该条数据的下一条
                        this.setState({
                            netList: newNetList,
                            netsCount: netsCount - 1,
                            selection: newNetList[index],
                        }, () => {
                            this.props.onSelectNet(this.state.selection)
                        })
                    }

                } else {
                    // 当前页不是最后一页,重新获取该页数据，且选中项的下标对应删除项的下标
                    await this.getNetList({ key: searchKey, page: this.currentPage }, { selectedIndex: index })
                }
                // 记录日志
                manageLog(
                    ManagementOps.DELETE,
                    __('删除 IP网段“${segment}” 成功', { segment: net.netType === NetType.Range ? net.originIP + '-' + net.endIP : net.ip + '/' + net.mask }),
                    net.name ? __('网段名称“${name}”', { name: net.name }) : null,
                    Level.WARN,
                )

            } catch (error) {
                getNetBindErrorMessage(error, net)
            }
        }
    }

    /**
     * 更新网段
     * @param data 更新的数据
     */
    protected async handleRequestEditSuccess(net: NetInfo) {
        try {
            // 保存添加
            if (!net.id) {
                const id = getIdFromLocaltion(
                    (await addNetwork(this.fromateAddNetParam(net))).getResponseHeader('Location'),
                );
                // 重新获取新增的网段信息
                const addNet = await this.getLastedNet(id);
                this.setState({
                    isEditNet: false,
                    searchKey: '',
                }, async () => {
                    await this.getNetList({}, { selectedIndex: 1 });
                    this.props.onSelectNet(addNet);
                });
                // 记录日志
                manageLog(
                    ManagementOps.ADD,
                    __('添加 IP网段“${segment}” 成功', { segment: addNet.netType === NetType.Range ? addNet.originIP + '-' + addNet.endIP : addNet.ip + '/' + addNet.mask }),
                    addNet.name ? __('网段名称“${name}”', { name: addNet.name }) : null,
                    Level.WARN,
                )

            } else {
                let originList = [...this.state.netList];

                await editNetwork(this.fromateAddNetParam(net))
                // 重新获取编辑的网段信息
                const editNet = await this.getLastedNet(net.id);
                // 更新列表
                const editNetIndex = findIndex(originList, (item) => (item.id === editNet.id))
                originList[editNetIndex] = editNet;

                this.setState({
                    netList: originList,
                    isEditNet: false,
                    editingNet: null,
                    selection: editNet,
                })
                // 记录日志
                manageLog(
                    ManagementOps.SET,
                    __('设置 IP网段“${segment}” 成功', { segment: editNet.netType === NetType.Range ? editNet.originIP + '-' + editNet.endIP : editNet.ip + '/' + editNet.mask }),
                    net.name ? __('网段名称“${name}”', { name: editNet.name }) : null,
                    Level.WARN,
                )
            }
        } catch (error) {
            getNetBindErrorMessage(error, net);
            if (error.code === PublicErrorCode.NotFound) {
                // 刷新当前页数据
                this.setState({
                    isEditNet: false,
                })
                this.getNetList({ key: this.state.searchKey, page: this.currentPage })
            }
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
     */
    protected loadSearchResult() {
        this.getNetList({ key: this.state.searchKey });
    }

    /**
     * 选中列表项
     * @param detail 选中数据
     */
    protected selectedNet({ detail }: { detail: NetInfo }) {
        const { selection } = this.state;
        this.setState({
            selection: detail,
        }, () => {
            // 点击同一项，不取消选中
            if (!detail) {
                this.setState({
                    selection,
                })
            }
        })
        detail && this.props.onSelectNet(detail)
    }

    /**
     * 手动触发页码改变
     * @param page 页码
     */
    protected handlePageChange(page: number) {
        this.getNetList({ key: this.state.searchKey, page })
        this.setState({ page })
    }

    /**
     * 关闭添加和编辑网段弹框
     */
    protected handleRequestEditCancel() {
        const { selection } = this.state;
        this.setState({
            isEditNet: false,
            editingNet: null,
            selection: null,
        }, () => {
            this.setState({
                selection,
            })
        })
    }

    /**
     * 改变表格参数
     */
    protected resetParams(params = { page: DefaultPage }) {
        this.dataGrid && this.dataGrid.changeParams(params);
    }

    /**
     * 格式化保存的网段参数
     * @param net 保存的网段
     */
    protected fromateAddNetParam(net) {
        let saveNet = {};

        if (net.netType === NetType.Range) {
            saveNet = {
                id: net.id,
                name: net.name,
                start_ip: net.originIP,
                end_ip: net.endIP,
                net_type: net.netType,
                ip_type: net.ipVersion,
            }
        } else {
            saveNet = {
                id: net.id,
                name: net.name,
                ip_address: net.ip,
                netmask: net.mask,
                net_type: net.netType,
                ip_type: net.ipVersion,
            }
        }

        return saveNet;
    }

    /**
     * 格式化获取到的网段数据
     * @param data 需要格式化的数据
     */
    protected formateGettedNets(data: ReadonlyArray<any>): Array<NetInfo> {
        return map(data, (item) => {
            return {
                id: item.id || '',
                name: item.name || '',
                originIP: item.start_ip || '',
                endIP: item.end_ip || '',
                ip: item.ip_address || '',
                mask: item.netmask || '',
                netType: item.net_type,
                ipVersion: item.ip_type,
            }
        })
    }

    /**
     * 获取单个最新的网段信息
     */
    protected async getLastedNet(netId: string): Promise<NetInfo> {
        // 重新获取最新的网段信息
        const latestNet = await getNetworkInfo({ id: netId });
        let {
            id,
            name,
            start_ip: originIP,
            end_ip: endIP,
            ip_address: ip,
            netmask: mask,
            net_type: netType,
            ip_type: ipVersion,
        } = latestNet;
        return { id, name, originIP, endIP, ip, mask, netType, ipVersion };
    }
}