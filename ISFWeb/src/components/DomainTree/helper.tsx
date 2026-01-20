import * as React from 'react';

/**
 * 节点类型
 */
export enum NodeType {
    /**
     * 域
     */
    Domain,

    /**
     * 部门
     */
    OU,

    /**
     * 用户
     */
    User,
}

/**
 * 获取节点类型
 * @param node 节点
 * @return 返回节点类型
 */
export function getNodeType(node: any): NodeType {
    return node.objectGUID ?
        (node.pathName ? NodeType.OU : NodeType.User)
        : NodeType.Domain
}

/**
 * 判断节点是否是叶子节点
 * @param node 节点
 * @param selectType 可选用户范围
 * @return 返回是否是子节点
 */
export function isLeaf(node: any): boolean {
    switch (getNodeType(node)) {
        case NodeType.Domain:

        case NodeType.OU:
            return false

        case NodeType.User:
            return true;
    }
}

/**
 * 获取节点图标
 */
export function getNodeIcon(node: any): { code: string } {
    switch (getNodeType(node)) {
        case NodeType.Domain:
            return {
                code: '\uf008',
            }

        case NodeType.OU:
            return {
                code: '\uf009',
            }

        case NodeType.User:
            return {
                code: '\uf007',
            }
    }
}

/**
* 接口返回的部门路径转换
*/
export function convertPath(path: string): string {
    const { ou, dc } = path.split(',').reduce(({ ou, dc }, item) => {
        // ouPath返回的内容的格式应该是 "OU=xx, DC=xx, CN=xx"，没有其他情况;
        // 所以先判断是否有等号，对于没有等号的值，直接舍弃。
        // 说明：此改动是因为后端有个bug：当用户名包含逗号时，返回结果中ouPath包含了一段错误的文字。
        // 导致前端处理数据出错，该bug后端修复改动范围较大，所以放在前端修复，去除掉错误的文字。
        if (item.includes('=')) {
            if (item.indexOf('OU=') === 0 || item.indexOf('CN=') === 0) {
                return { ou: [item.split('=')[1], ...ou], dc }
            }

            return { ou, dc: [...dc, item.split('=')[1]] }
        }

        return { ou, dc }
    }, { ou: [], dc: [] })

    return dc.join('.') + '/' + ou.join('/')
}

/**
* 选中的上级部门过滤已选中的下级用户或部门
*/
export function filterSelect(selection: any, select: any): { ou: ReadonlyArray<any>; newSec: ReadonlyArray<any> } {
    const { ou, dc } = selection.reduce(({ ou, dc }, item) => {
        if (item.ipAddress) {
            return { ou: [item, ...ou], dc }
        }

        const currentPath = item.pathName || item.ouPath

        if (
            (select.ipAddress && convertPath(currentPath).split('/')[0] === select.name)
            ||
            (currentPath).indexOf(select.pathName) !== -1
        ) {
            return { ou, dc }
        }

        return { ou, dc: [...dc, item] }
    }, { ou: [], dc: [] })

    return { ou, newSec: dc }
}