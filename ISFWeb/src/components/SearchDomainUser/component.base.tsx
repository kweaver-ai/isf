import * as React from 'react';
import { noop } from 'lodash';
import { Message2 as Message } from '@/sweet-ui';
import { usrmSearchDomainInfoByName } from '@/core/thrift/sharemgnt/sharemgnt';
import { ErrorCode } from '@/core/thrift/sharemgnt/errcode'
import WebComponent from '../webcomponent';
import __ from './locale';
/**
 * 域用户信息
 */
interface DomainUser {
    /**
     * 用户id
     */
    objectGUID: string;

    /**
     * 用户名称
     */
    displayName: string;

    /**
     * 用户路径
     */
    dnPath: string;

    /**
     * 邮箱
     */
    email: string;

    /**
     * 身份证
     */
    idcardNumber: string;

    /**
     * 登录名
     */
    loginName: string;

    /**
     * 用户父级路径
     */
    ouPath: string;
}

/**
 * 域用户部门信息
 */
interface DomainDep {
    /**
     * 部门id
     */
    objectGUID: string;

    /**
     * 是否导入组织下的所有子组织及用户
     */
    importAll: null;

    /**
     * 部门名
     */
    name: string;

    /**
     * 部门父级路径
     */
    parentOUPath: string;

    /**
     * 部门路径
     */
    pathName: string;

    /**
     * 组织负责人
     */
    rulerName: string;
}

interface Props {
    /**
     * 选择事件
     */
    onRequestSelect: (DomainUserOrDep: DomainUser | DomainDep) => void;

    /**
     * 检测 输入框的值是否发生改变
     */
    onValueChange: (value: string) => boolean;

    /**
    * 宽度
    */
    width: string | number;

    /**
     * 是否搜索框能够输入
     */
    disabled?: boolean;

    /**
     * 搜索框默认内容
     */
    placeholder?: string;

    /**
     * 搜索值
     */
    value?: string;

    /**
     * 域id
     */
    domainId?: string;
}

interface State {
    /**
     * 搜索结果
     */
    results: ReadonlyArray<DomainUser | DomainDep>;

    /**
     * 搜索关键字
     */
    searchKey: string;

    /**
     * 域异常
     */
    domainError: boolean;
}

export default class SearchDomainUserBase extends WebComponent<Props, State> {
    static defaultProps = {
        onRequestSelect: noop,
        onValueChange: noop,
        disabled: false,
        width: 200,
        placeholder: __('搜索'),
        value: '',
        domainId: '',
    }

    state = {
        results: [],
        searchKey: this.props.value || '',
        domainError: false,
    }

    /**
    * 取消的请求
    */
    cancelRequest = {
        abort: noop,
    }

    /**
    * 是否懒加载
    */
    lazyLoad: boolean = false;

    componentWillUnmount() {
        // 组件卸载时取消搜索请求
        this.cancelRequest.abort && this.cancelRequest.abort()
    }

    /**
     * 根据key获取域用户或部门
     * @return 域用户或部门数组
     */
    protected getDomainUserOrDepByKey = (key: string, offset: number = 0, limit = 10): Promise<{ ous: ReadonlyArray<DomainDep>; users: ReadonlyArray<DomainUser> }> => {
        if (key) {
            this.setState({
                domainError: false,
            })

            // 控制新请求数据是否是懒加载
            this.lazyLoad = offset === 0 ? false : true

            // 发起新的请求前取消上次请求
            this.cancelRequest.abort && this.cancelRequest.abort()
            // 将请求赋值给 cancelRequest ，以便做取消操作
            const request = usrmSearchDomainInfoByName([this.props.domainId || -1, key, offset, limit])
            this.cancelRequest = request

            return request
        }
    }

    /**
     * 搜索框失焦处理
     */
    protected handelOnBlur = () => {
        // 失焦后取消搜索请求
        this.cancelRequest.abort && this.cancelRequest.abort()
    }

    /**
     * 加载失败
     */
    protected loaderFailed = async (errorEvent): Promise<void> => {
        if (errorEvent.error) {
            const { error: { errID, errDetail, errMsg } } = errorEvent

            this.setState({
                domainError: true,
            })
            if (errID === ErrorCode.DomainUnavailable) {
                await Message.info({ message: __('域“${domain}”异常，连接LDAP服务器异常，请检查域控IP是否正确，或域控制器是否开启。', { domain: errDetail }) })
            } else if (errID === ErrorCode.RequestDataLarge) {
                await Message.info({ message: __('请求数据过大，请输入更精确的关键字重试。') })
            } else if (errMsg) {
                await Message.info({ message: errMsg })
            }
        }
    }

    /**
     * 获取搜索到的结果
     * @param results 域用户或部门数组
     */
    protected getSearchData = ({ ous, users }: { ous: ReadonlyArray<DomainDep>; users: ReadonlyArray<DomainUser> }): void => {
        const newResults = this.props.domainId ? ous : [...ous, ...users]

        this.setState({
            results: this.lazyLoad ? [...this.state.results, ...newResults] : newResults,
        })
    }

    /**
     * 选择搜索到的单个域用户或部门
     */
    protected selectItem = (domainUser: DomainUser | DomainDep): void => {
        this.props.onRequestSelect(domainUser);

        if (domainUser && (domainUser.name || domainUser.displayName)) {
            this.setState({
                searchKey: domainUser.name || domainUser.displayName,
            })
        }
    }

    /**
     * 关键词修改
     */
    protected handelChange = (searchKey): void => {
        this.setState({
            searchKey,
        })
        this.props.onValueChange(searchKey)

        if (!searchKey) {
            this.setState({
                results: [],
            })
        }
    }
}
