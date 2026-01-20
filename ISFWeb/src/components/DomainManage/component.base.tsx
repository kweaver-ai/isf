import * as React from 'react'
import { usrmGetAllDomains, setDomainSyncStatus, startSync, getDomainConfig, setDomainStatus, deleteDomain } from '@/core/thrift/sharemgnt/sharemgnt'
import { manageLog, Level, ManagementOps } from '@/core/log'
import { ErrorCode } from '@/core/thrift/sharemgnt/errcode'
import { Message2 as Message, Toast } from '@/sweet-ui'
import { ListTipStatus } from '../ListTipComponent/helper'
import WebComponent from '../webcomponent'
import { ActionType, DomainType } from './helper'
import __ from './locale'

interface DomainManageState {
    /**
     * loading状态
     */
    listTipStatus: ListTipStatus;

    /**
     * 域列表
     */
    domainList: ReadonlyArray<Core.ShareMgnt.ncTUsrmDomainInfo>;

    /**
     * 渲染类型
     */
    actionType: ActionType;

    /**
     * 选中项
     */
    selection: Core.ShareMgnt.ncTUsrmDomainInfo;

    /**
     * 展开关键字
     */
    expandedKeys: ReadonlyArray<number>;
}

export default class DomainManageBase extends WebComponent<any, DomainManageState> {

    state: DomainManageState = {
        /**
         * 列表提示
         */
        listTipStatus: ListTipStatus.Loading,

        /**
         * 域列表数据
         */
        domainList: [],

        /**
         * 动作类型
         */
        actionType: ActionType.None,

        /**
         * 选中项
         */
        selection: {},

        /**
         * 展开项的id
         */
        expandedKeys: [],
    }

    /**
     * 编辑的域
     */
    editDomain = {}

    async componentDidMount() {
        await this.loadDomainList()
    }

    /**
     * 加载数据
     */
    private loadDomainList = async (): Promise<void> => {
        this.setState({
            listTipStatus: ListTipStatus.Loading,
        })

        try {
            const domainList = await usrmGetAllDomains()

            this.setState({
                domainList: this.formatDomainList(domainList),
                listTipStatus: domainList.length < 1 ? ListTipStatus.Empty : ListTipStatus.None,
            })
        } catch (ex) {
            this.setState({
                listTipStatus: ListTipStatus.LoadFailed,
            })
        }
    }

    /**
     *  处理列表数据成树状结构
     */
    private formatDomainList = (list: ReadonlyArray<Core.ShareMgnt.ncTUsrmDomainInfo>): ReadonlyArray<Core.ShareMgnt.ncTUsrmDomainInfo> => {
        let mainDomainId, mainDomain, subdomains = [], maindomains = {}, domainslist = [];

        list.forEach((domain) => {
            switch (domain.type) {
                case DomainType.Primary:
                    domain.children = [];
                    maindomains[domain.id] = domain;
                    domainslist = [...domainslist, domain];
                    break;

                case DomainType.Sub:
                    domain.children = [];
                    maindomains[domain.id] = domain;
                    subdomains = [...subdomains, domain];
                    break;

                case DomainType.Trust:
                    subdomains = [...subdomains, domain];
            }
        })

        subdomains.forEach((subdomain) => {
            mainDomainId = subdomain.parentId.toString();
            mainDomain = maindomains[mainDomainId];
            if (mainDomain) {
                if (!mainDomain.children) {
                    mainDomain.children = [];
                }
                mainDomain.children = [...mainDomain.children, subdomain]
            }
        })

        return domainslist
    }

    /**
     *  定期同步开关
     */
    protected changeSyncStatus = async (syncStatus: boolean, id: number, name: string): Promise<void> => {
        const { domainList } = this.state;

        try {
            await setDomainSyncStatus([id, syncStatus ? 0 : -1])

            this.setState({
                domainList: domainList.map((item) => item.id === id ? { ...item, syncStatus: syncStatus ? 0 : 1 } : item),
            })

            if (syncStatus) {
                manageLog(
                    ManagementOps.SET,
                    __('开启 域控 “${name}” 定期同步 成功', { name }),
                    '',
                    Level.INFO,
                )

                startSync([id.toString(), 1])
            } else {
                manageLog(
                    ManagementOps.SET,
                    __('关闭 域控 “${name}” 定期同步 成功', { name }),
                    '',
                    Level.INFO,
                )
            }
        } catch (ex) {
            this.handleError(ex, name)
        }
    }

    /**
     * 调用启用、禁用域接口函数
     */
    private handleSetDomainStatus = async (status: boolean, id: number, name: string): Promise<void> => {
        try {
            await setDomainStatus([id, status])

            this.loadDomainList()

            manageLog(
                ManagementOps.SET,
                status ? __('启用 域控 “${name}” 成功', { name }) : __('禁用 域控 “${name}” 成功', { name }),
                '',
                Level.INFO,
            )
        } catch (ex) {
            this.handleError(ex, name)
        }
    }

    /**
     *  启用或禁用 域
     */
    protected changeDomainStatus = async (detail: boolean, id: number, name: string): Promise<void> => {
        if (!detail) {
            if (await Message.alert({ message: __('禁用域控制器将导致已导入的用户无法验证，您确定要执行此操作吗？'), showCancelIcon: true })) {
                this.handleSetDomainStatus(detail, id, name)
            }
        } else {
            if (await Message.alert({ message: __('是否确定启用“${name}”？', { name }), showCancelIcon: true })) {
                await this.handleSetDomainStatus(detail, id, name)
                Toast.open(__('域控制器已启用'))
            }
        }
    }

    /**
     * 改变渲染抽屉
     */
    protected changeActionType = (type: ActionType): void => {
        this.setState({
            actionType: type,
        })
    }

    /**
     * 立即同步单个域
     */
    protected syncNow = async (id: number, autoSync: boolean, name: string): Promise<void> => {
        try {
            await startSync([id.toString(), autoSync])
            await getDomainConfig([id])
            await Message.info({ message: __('域控信息将开始同步，请稍后从日志中查询同步结果。') })
        } catch (ex) {
            this.handleError(ex, name)
        }
    }

    /**
     * 编辑单个域
     */
    protected setDomain = (e, editDomain: any): void => {
        e.stopPropagation();
        this.editDomain = editDomain;
        this.setState({
            actionType: ActionType.Edit,
        })
    }

    /**
     * 删除 单个域
     */
    protected deleteDomain = async (id: number, name: string): Promise<void> => {
        if (await Message.alert({ message: __('删除域控制器将导致已导入用户无法验证，您确定要执行此操作吗？'), showCancelIcon: true })) {

            try {
                await deleteDomain([id])
                manageLog(
                    ManagementOps.DELETE,
                    __('删除 域控 “${name}” 成功', { name }),
                    '',
                    Level.WARN,
                )
                this.loadDomainList()
            } catch (ex) {
                switch (ex.error.errID) {
                    case ErrorCode.DomainNotExists:
                        this.loadDomainList()
                        break
                }
            }
        }
    }

    /**
     * 选中项
     */
    protected changeSelection = ({ detail }: { detail: Core.ShareMgnt.ncTUsrmDomainInfo }): void => {
        this.setState({
            selection: detail ? detail : {},
        })
    }

    /**
     * 新建或者编辑域控信息完成回调
     */
    protected setDomainInfoFinish = (id?: number): void => {
        this.loadDomainList()
        id && this.setState({
            expandedKeys: [...this.state.expandedKeys, id],
        })
    }

    /**
     * 展开子域
     */
    protected expand = ({ expandedKeys }) => {
        this.setState({
            expandedKeys,
        })
    }

    /**
     * 处理域不存在回调
     */
    protected handleDomainError = (name: string): void => {
        this.setState({
            actionType: ActionType.None,
        }, async () => {
            if (await Message.alert({ message: __('域控制器 “${name}” 已不存在。', { name }) })) {
                this.loadDomainList()
            }
        })
    }

    /**
     * 错误处理
     */
    private handleError = async (ex, name: string): Promise<void> => {
        switch (ex.error.errID) {
            case ErrorCode.DomainNotExists:
                if (await Message.alert({ message: __('域控制器 “${name}” 已不存在。', { name }) })) {
                    this.loadDomainList()
                }
                break
        }
    }
}