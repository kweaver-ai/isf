import * as React from 'react'
import { map } from 'lodash';
import { ShareMgnt, EACP } from '@/core/thrift'
import { getNetworkListByAccessor } from '@/core/apis/console/networkRestriction'
import { ListTipStatus } from '../../ListTipComponent/helper';
import { NetInfo as VisitorNetInfo } from '../VisitorNetBind/helper'
import WebComponent from '../../webcomponent'
import __ from './locale';

export enum SearchField {
    User,
    DeviceId,
}

/**
 * 每页显示数据条数
 */
export const PageSize = 20;

/**
 * 默认开始页
 */
export const DefaultPage = 1;

export const deviceType = {
    '0': __('未知设备'),
    '1': __('iOS'),
    '2': __('Android'),
    '3': __('Windows Phone'),
    '4': __('Windows'),
    '5': __('Mac OS X'),
    '6': __('Web'),
}

export default class BindingQueryBase extends WebComponent<Console.BindingQuery.Props, Console.BindingQuery.State> {
    static defaultProps = {
        userid: '',
    }

    state = {
        searchField: SearchField.User,
        visitorNetInfos: [],
        deviceInfos: [],
        deviceIdUsers: [],
        searchResults: [],
        searchKey: '',
        listTipStatus: ListTipStatus.Empty,
        deviceListStatus: ListTipStatus.Empty,

    }

    loaders = {
        [SearchField.User]: this.searchUserLoader.bind(this),
        [SearchField.DeviceId]: this.searchDeviceIdLoader.bind(this),
    }

    /**
     * 默认加载条数
     */
    defaultLimit = 20

    /**
     * 选择的要搜索的用户
     */
    selectedUser = null;

    /**
     * 绑定访问者的网段总数
     */
    visitorNetCount = 0;

    /**
     * 存储正在选择的用户的id
     */
    currentSelectUserId: string = '';

    /**
     * 在组件销毁后设置state，防止内存泄漏
     */
    componentWillUnmount() {
        this.setState = (state, callback) => {
            return;
        };
    }

    /**
     * 用户、组织搜索
     */
    searchUserLoader(key, start = 0) {
        return key ? Promise.all([
            ShareMgnt('Usrm_SearchDepartments', [this.props.userid, key, start, this.defaultLimit]),
            ShareMgnt('Usrm_SearchSupervisoryUsers', [this.props.userid, key, start, this.defaultLimit]),
        ]).then(([departs, users]) => [...departs, ...users]) : Promise.resolve([])
    }

    /**
     * 搜索设备识别码
     */
    searchDeviceIdLoader(key: string, start: number = 0): Promise<ReadonlyArray<string>> {
        return key ? EACP('EACP_SearchDevices', [key, start, this.defaultLimit]) : Promise.resolve([])
    }

    /**
     * 搜索关键字
     */
    handleSearchChange(searchKey) {
        if (searchKey !== this.state.searchKey) {
            this.setState({
                searchKey,
                visitorNetInfos: [],
                deviceInfos: [],
                deviceIdUsers: [],
                searchResults: [],
                listTipStatus: ListTipStatus.Empty,
                deviceListStatus: ListTipStatus.Empty,
            })
            this.visitorNetCount = 0;
        }
    }

    /**
     * 搜索范围 访问者，文档库，设备识别码
     */
    handleChangeSearchField({ detail: searchField }) {
        this.setState({
            searchField,
            searchKey: '',
            visitorNetInfos: [],
            deviceInfos: [],
            deviceIdUsers: [],
            searchResults: [],
            listTipStatus: ListTipStatus.Empty,
            deviceListStatus: ListTipStatus.Empty,
        })
        this.visitorNetCount = 0;
    }

    /**
     * 搜索结果
     */
    handleSearchLoaded(searchResults) {
        this.setState({
            searchResults,
        })
    }

    /**
     * 搜索懒加载
     */
    protected lazyLoade = async (page: number, limit: number): Promise<void> => {
        const { searchField } = this.state

        this.setState({
            searchResults: [
                ...this.state.searchResults,
                ...(await this.loaders[searchField](this.state.searchKey, (page - 1) * limit)),
            ],
        })

        this.refs['auto-complete'].toggleActive(true)
    }

    /**
     * 设备列表懒加载回调
     * @param detail 分页信息
     */
    handleRequestLoadDevice = ({ detail }) => {
        const { deviceInfos } = this.state;
        const { limit, start } = detail;
        this.setState({
            deviceListStatus: ListTipStatus.Loading,
        }, async () => {

            const result = await EACP('EACP_GetDevicesByUserId', [this.currentSelectUserId, start, limit])

            this.setState({
                deviceInfos: [...deviceInfos, ...result],
                deviceListStatus: ListTipStatus.None,
            })
        })
    }

    /**
     * 查询用户绑定
     */
    async queryUser(item) {
        this.currentSelectUserId = item.id;

        this.setState({
            listTipStatus: ListTipStatus.Loading,
            deviceListStatus: ListTipStatus.Loading,
        })

        const device = await EACP('EACP_GetDevicesByUserId', [item.id, 0, this.defaultLimit]);

        const results = await getNetworkListByAccessor(
            {
                id: item.id || item.departmentId,
                offset: (DefaultPage - 1) * PageSize,
                limit: PageSize,
            },
        )
        const netList = this.formateGettedNets(results.data);
        this.visitorNetCount = results.count;
        this.setState({
            visitorNetInfos: netList,
            deviceInfos: device,
            searchKey: item.displayName || item.departmentName,
            listTipStatus: netList.length < 1 ?
                (this.state.searchKey ? ListTipStatus.NoSearchResults : ListTipStatus.Empty)
                : ListTipStatus.None,
            deviceListStatus: device.length < 1 ?
                (this.state.searchKey ? ListTipStatus.NoSearchResults : ListTipStatus.Empty)
                : ListTipStatus.None,
        }, () => {
            this.refs['auto-complete'].toggleActive(false)
        })
        this.selectedUser = item;
    }

    /**
     * 处理懒加载用户绑定的网段
     */
    protected handleLazyLoadNetList = async ({ detail }) => {
        if (this.visitorNetCount > this.state.visitorNetInfos.length) {
            const { start: offset, limit } = detail;
            this.setState({
                listTipStatus: ListTipStatus.Loading,
            })
            const more = await getNetworkListByAccessor(
                {
                    id: this.selectedUser.id || this.selectedUser.departmentId,
                    offset,
                    limit,
                },
            )
            const netList = this.formateGettedNets(more.data);
            this.setState({
                visitorNetInfos: [...this.state.visitorNetInfos, ...netList],
                listTipStatus: ListTipStatus.None,
            })
        }
    }

    /**
     * 渲染按设备识别码的查询结果
     */
    async renderDeviceId(udid: string): Promise<void> {
        this.setState({
            listTipStatus: ListTipStatus.Loading,
        })
        const users = await this.queryDeviceIdUsers(udid);
        this.setState({
            deviceIdUsers: users,
            searchKey: udid,
            listTipStatus: users.length < 1 ?
                (this.state.searchKey ? ListTipStatus.NoSearchResults : ListTipStatus.Empty)
                : ListTipStatus.None,
        }, () => {
            this.refs['auto-complete'].toggleActive(false)
        })
    }

    /**
     * 查询绑定该设备识别码的用户
     */
    queryDeviceIdUsers(udid: string, start: number = 0): Promise<ReadonlyArray<string>> {
        return EACP('EACP_SearchUserByDeviceUdid', [udid, start, this.defaultLimit])
    }

    /**
     * 按下enter
     */
    handleEnter(e, selectIndex: number) {
        if (selectIndex >= 0) {
            switch (this.state.searchField) {
                case SearchField.User: {
                    this.queryUser(this.state.searchResults[selectIndex])
                    break
                }
                case SearchField.DeviceId:
                    this.renderDeviceId(this.state.searchResults[selectIndex])
                    break
            }
        }
    }

    /**
     * 设备识别码的懒加载
     */
    async lazyLoadDeviceUdidUsers({ detail }): Promise<void> {
        this.setState({
            listTipStatus: ListTipStatus.Loading,
        })
        const users = await this.queryDeviceIdUsers(this.state.searchKey, detail.start)
        this.setState({
            listTipStatus: ListTipStatus.None,
            deviceIdUsers: [
                ...this.state.deviceIdUsers,
                ...users,
            ],
        })
    }

    /**
     * 格式化获取到的网段数据
     * @param data 需要格式化的数据
     */
    protected formateGettedNets(data: ReadonlyArray<any>): Array<VisitorNetInfo> {
        return map(data, (item) => {
            return {
                id: item.id,
                name: item.name,
                originIP: item.start_ip,
                endIP: item.end_ip,
                ip: item.ip_address,
                mask: item.netmask,
                netType: item.net_type,
            }
        })
    }
}