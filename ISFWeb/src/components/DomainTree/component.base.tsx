import * as React from 'react';
import { noop } from 'lodash';
import { Message2 as Message, Toast } from '@/sweet-ui';
import { usrmGetAllDomains, usrmExpandDomainNode, getDomainById } from '@/core/thrift/sharemgnt/sharemgnt';
import { NodeType, getNodeType, convertPath, filterSelect } from './helper';
import __ from './locale';

interface DomainTreeProps {
    /**
     * 渲染某个域时需要的id
     */
    domainId: number;

    /**
     * 选中项
     */
    selection: ReadonlyArray<any>;

    onSelectionChange: (selections: any, search?: boolean) => void;

    /**
     * 跳转域认证集成
     */
    doRedirectDomain: () => void;
}

interface DomainTreeState {
    /**
     * 树加载数据
     */
    treeData: ReadonlyArray<any>;

    /**
     * 加载中
     */
    loading: boolean;
}

export default class DomainTreeBase extends React.PureComponent<DomainTreeProps, DomainTreeState> {
    static defaultProps = {
        selection: [],
        onSelectionChange: noop,
        doRedirectDomain: noop,
    }

    state: DomainTreeState = {
        treeData: [],
        loading: true,
    }

    /**
        * 取消所有选择节点
        */
    cancelSelections = () => this.ref ? (this.ref as any).cancelSelections() : null

    /**
     * 加载树数据
     */
    async componentDidMount() {
        const { domainId } = this.props;

        if (domainId) {
            const { id } = await getDomainById([domainId])
            const roots = await usrmGetAllDomains()

            this.setState({
                treeData: roots.filter((item) => item.id === id),
                loading: false,
            })
        } else {
            const roots = await usrmGetAllDomains()

            this.setState({
                treeData: roots,
                loading: false,
            })
        }
    }

    /**
     * 点击获取子节点
     */
    protected getNodeChildren = async (node): Promise<ReadonlyArray<any>> => {
        let nodePath, domainInfo;

        if (getNodeType(node) === NodeType.Domain) {
            nodePath = node.name.indexOf('=') !== -1 ? node.name : node.name.split('.').map((dc) => 'DC=' + dc).join(',')
            domainInfo = {
                ncTUsrmDomainInfo: {
                    ...node,
                    config: {
                        ncTUsrmDomainConfig: node.config,
                    },
                },
            }
        } else {
            nodePath = node.pathName;
            domainInfo = node.domainInfo;
        }

        try {
            const { domainId } = this.props;
            const { ous, users } = await usrmExpandDomainNode([domainInfo, nodePath])

            return domainId ? ous.map((item) => {
                return {
                    ...item,
                    domainInfo,
                    parentNode: node,
                }
            }) : [...ous, ...users].map((item) => {
                return {
                    ...item,
                    domainInfo,
                    parentNode: node,
                }
            })
        } catch (ex) {
            ex.error && Message.error({
                message: ex.error.errMsg,
            })

            return []
        }
    }

    /**
     * 处理选中
     */
    protected handleSelect = (selections: ReadonlyArray<any>, select: any, id: string): void => {
        const { onSelectionChange, selection } = this.props;
        let currentIsSelect, parentNode;

        if (select) {
            if (select.ipAddress) {
                currentIsSelect = selection.some((item) => select.ipAddress.indexOf(item.ipAddress) !== -1)
            } else {
                const { pathName, ouPath, parentOUPath, objectGUID } = select;

                selection.some((item) => {
                    if (((parentOUPath || ouPath).indexOf(item.pathName) !== -1)
                        || ((convertPath(parentOUPath || ouPath).split('/')[0] === item.name) && item.ipAddress)
                    ) {
                        parentNode = item

                        return true
                    }
                })
                currentIsSelect = selection.some((item) => this.props.domainId ? pathName === item.pathName : objectGUID === item.objectGUID)
            }

            if (currentIsSelect) {
                Toast.open(__('该对象已存在'))
                this.ref.toggleSelect(id)
            } else if (parentNode) {
                Toast.open(__('您选择的对象已包含在已选部门“${name}”里', { name: parentNode.name }))
                this.ref.toggleSelect(id)
            } else {
                if (selection.some((item) => item.pathName || item.ouPath)) {
                    const { ou, newSec } = filterSelect(selection, select)

                    onSelectionChange([...ou, ...newSec, ...selections])
                } else {
                    onSelectionChange([...selection, ...selections])
                }
            }
        }
    }

    /**
     * 处理搜索选中
     */
    protected handleSelectResult = async (select: any): Promise<void> => {
        try {
            const { treeData } = this.state,
                { selection, onSelectionChange } = this.props,
                { pathName, ouPath, parentOUPath, objectGUID } = select,
                [rootNodeName, secondNodeName] = convertPath(parentOUPath || ouPath).split('/'),
                [rootNode] = treeData.filter((item) => item.name.toLowerCase() === rootNodeName.toLowerCase()),
                currentIsSelect = selection.some((item) => this.props.domainId ? pathName === item.pathName : objectGUID === item.objectGUID),
                domainInfo = {
                    ncTUsrmDomainInfo: {
                        ...rootNode,
                        config: {
                            ncTUsrmDomainConfig: rootNode.config,
                        },
                    },
                };
            let higherLeverNode, parentNode;

            if (secondNodeName) {
                // ouPath返回的内容的格式应该是 "OU=xx, DC=xx, CN=xx"，没有其他情况;
                // 所以先判断是否有等号，对于没有等号的值，直接舍弃。
                // 说明：此改动是因为后端有个bug：当用户名包含逗号时，返回结果中ouPath包含了一段错误的文字。
                // 导致前端处理数据出错，该bug后端修复改动范围较大，所以放在前端修复，去除掉错误的文字。
                const [rootNodePath, ...secondNodePath] = (parentOUPath || ouPath).split(',').filter((item) => item.includes('='));
                const { ous } = await usrmExpandDomainNode([domainInfo, secondNodePath.join()]);

                [parentNode] = ous.filter((item) => item.name === rootNodePath.split('=')[1]);
            }

            select = {
                ...select,
                domainInfo,
                parentNode: secondNodeName ? {
                    ...parentNode,
                    domainInfo:
                    {
                        ncTUsrmDomainInfo: {
                            ...rootNode,
                            config: {
                                ncTUsrmDomainConfig: rootNode.config,
                            },
                        },
                    },
                    parentNode: rootNode,
                } : rootNode,
            }
            selection.some((item) => {
                if (((parentOUPath || ouPath).indexOf(item.pathName) !== -1)
                    || ((rootNodeName === item.name) && item.ipAddress)
                ) {
                    higherLeverNode = item

                    return true
                }
            })

            if (currentIsSelect) {
                Toast.open(__('该对象已存在'))
            } else if (higherLeverNode) {
                Toast.open(__('您选择的对象已包含在已选部门“${name}”里', { name: higherLeverNode.name }))
            } else {
                if (selection.some((item) => item.pathName || item.ouPath)) {
                    const { ou, newSec } = filterSelect(selection, select)

                    onSelectionChange([...ou, ...newSec, select], true)
                } else {
                    onSelectionChange([...selection, select], true)
                }
            }
        } catch (ex) {
            ex.error && Message.error({
                message: ex.error.errMsg,
            })
        }
    }
}