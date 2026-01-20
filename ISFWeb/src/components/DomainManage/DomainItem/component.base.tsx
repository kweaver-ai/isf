import * as React from 'react'
import { noop } from 'lodash'
import WebComponent from '../../webcomponent'
import { ActionType, DomainType } from '../helper'

/**
 * 展开项
 */
export enum ExpandItem {
    /**
     * 域控信息
     */
    DomainInfo,

    /**
     * 定期同步设置
     */
    SyncSetting,

    /**
     * 关键字设置
     */
    KeySetting,

    /**
     * 备用域控制信息
     */
    SpareDomain,

}

interface DomainItemProps {
    /**
     * 渲染类型
     */
    actionType: ActionType;

    /**
     * 选中项
     */
    selection: Core.ShareMgnt.ncTUsrmDomainInfo;

    /**
     * 编辑项
     */
    editDomain: Core.ShareMgnt.ncTUsrmDomainInfo;

    /**
     * 编辑域控信息成功回调
     */
    onEditDomainInfoSuccess: (id?: number) => void;

    /**
     * 主域不存在
     */
    onDomainError: (name: string) => void;

    /**
     * 关闭按钮
     */
    onRequestCancel: () => void;
}

interface DomainItemState {
    /**
     * 其他三项是否可见
     */
    restShow: boolean;

    /**
     * 展开选项
     */
    expandItem: {
        /**
         * 域控信息是否展开
         */
        infoExpand: boolean;

        /**
         * 定期同步设置信息是否展开
         */
        syncExpand: boolean;

        /**
         * 关键字配置是否展开
         */
        keyExpand: boolean;

        /**
         * 备用域信息是否展开
         */
        spareExpand: boolean;
    };

    /**
     * 域控信息是否编辑中
     */
    isDomainInfoEditStatus: boolean;

    /**
     * 定期同步信息是否编辑中
     */
    isSyncSettingEditStatus: boolean;

    /**
     * 同步关键字是否编辑中
     */
    isKeywordSettingEditStatus: boolean;

    /**
     * 备用域是否编辑中
     */
    isSpareDomainEditStatus: boolean;
}

export default class DomainItemBase extends WebComponent<DomainItemProps, DomainItemState> {
    static defaultProps = {
        /**
         * 渲染类型
         */
        actionType: ActionType.Add,

        /**
         * 选中项
         */
        selection: {},

        /**
         * 编辑项
         */
        editDomain: {},

        /**
         * 编辑域控信息回调或者重新加载列表
         */
        onEditDomainInfoSuccess: noop,

        /**
         * 域不存在
         */
        onDomainError: noop,

        /**
         * 关闭按钮
         */
        onRequestCancel: noop,
    }

    state: DomainItemState = {
        /**
         * 其他三项是否可见
         */
        restShow: this.props.actionType === ActionType.Edit && this.props.editDomain.type === DomainType.Primary ? true : false,

        /**
         * 展开选项
         */
        expandItem: {
            /**
             * 域控信息
             */
            infoExpand: true,

            /**
             * 同步设置信息
             */
            syncExpand: false,

            /**
             * 关键字信息
             */
            keyExpand: false,

            /**
             * 备用域信息
             */
            spareExpand: false,
        },

        /**
         * 域控信息保存取消按钮是否显示
         */
        isDomainInfoEditStatus: this.props.actionType === ActionType.Add || this.props.actionType === ActionType.Edit && this.props.editDomain.type !== DomainType.Primary,

        /**
         * 定期同步设置保存取消是否显示
         */
        isSyncSettingEditStatus: false,

        /**
         * 同步关键字保存取消是否显示
         */
        isKeywordSettingEditStatus: false,

        /**
         * 备用域保存取消是否显示
         */
        isSpareDomainEditStatus: false,
    }

    /**
     * 存储新建完成后域控信息
     */
    domainInfo = {
        id: -1,
        name: '',
        domainName: '',
        password: '',
        useSSL: false,
        ipAddress: '',
        domainPort: 389,
    }

    /**
     * 获取子元素
     */
    protected onRef = (ref: object): void => {
        this.child = ref
    }

    /**
     * 展开或收起
     */
    protected expandItem = (type: ExpandItem, isExpand: boolean): void => {
        const { isDomainInfoEditStatus, isSyncSettingEditStatus, isKeywordSettingEditStatus, isSpareDomainEditStatus } = this.state;

        if (!isDomainInfoEditStatus && !isSyncSettingEditStatus && !isKeywordSettingEditStatus && !isSpareDomainEditStatus) {
            switch (type) {
                case ExpandItem.DomainInfo:
                    this.setState({
                        expandItem: {
                            infoExpand: isExpand,
                            syncExpand: false,
                            keyExpand: false,
                            spareExpand: false,
                        },
                    })
                    break;

                case ExpandItem.SyncSetting:
                    this.setState({
                        expandItem: {
                            infoExpand: false,
                            syncExpand: isExpand,
                            keyExpand: false,
                            spareExpand: false,
                        },
                    })
                    break;

                case ExpandItem.KeySetting:
                    this.setState({
                        expandItem: {
                            infoExpand: false,
                            syncExpand: false,
                            keyExpand: isExpand,
                            spareExpand: false,
                        },
                    })
                    break;

                case ExpandItem.SpareDomain:
                    this.setState({
                        expandItem: {
                            infoExpand: false,
                            syncExpand: false,
                            keyExpand: false,
                            spareExpand: isExpand,
                        },
                    })
                    break;
            }
        }

    }

    /**
     * 新建或编辑域控信息完毕获取域信息
     */
    protected getDoaminInfo = ({ id, name, adminName, password, useSSL, ipAddress, domainPort, isDomainInfoEditStatus }) => {
        const { editDomain, selection, actionType } = this.props
        this.setState({
            isDomainInfoEditStatus,
            restShow: actionType === ActionType.Add ? !selection.type : editDomain.type && editDomain.type === DomainType.Primary,
        })

        this.domainInfo = {
            id,
            name,
            domainName: adminName,
            password,
            useSSL,
            ipAddress,
            domainPort,
        }

        this.props.selection.id ? this.props.onEditDomainInfoSuccess(this.props.selection.id) : this.props.onEditDomainInfoSuccess()

        actionType === ActionType.Add ? !selection.type ? null : this.props.onRequestCancel() : editDomain.type && editDomain.type === DomainType.Primary ? null : this.props.onRequestCancel()
    }

    /**
     * 编辑定期同步设置完成回调
     */
    protected getSynSetInfo = () => {
        this.props.onEditDomainInfoSuccess()
    }

    /**
     * 新建子域时点击新建（调用子组件的方法）
     */
    protected saveDomainInfo = () => {
        this.child.saveDomainInfo()
    }

    /**
     * 设置域控信息编辑状态
     */
    protected changeDomainInfoEditStatus = (status: boolean): void => {
        this.setState({
            isDomainInfoEditStatus: status,
        })
    }

    /**
     * 设置同步信息编辑状态
     */
    protected changeSyncInfoEditStatus = (status: boolean): void => {
        this.setState({
            isSyncSettingEditStatus: status,
        })
    }

    /**
     * 设置同步关键字编辑状态
     */
    protected changeKeySettingEditStatus = (status: boolean): void => {
        this.setState({
            isKeywordSettingEditStatus: status,
        })
    }

    /**
     * 设置备用域编辑状态
     */
    protected changeSpareDomainEditStatus = (status: boolean): void => {
        this.setState({
            isSpareDomainEditStatus: status,
        })
    }

    /**
     * 处理域错误
     */
    protected handleDomainError = (name: string): void => {
        this.props.onDomainError(name)
    }

    /**
     * 关闭按钮
     */
    protected cancel = (): void => {
        this.props.onRequestCancel()
    }
}