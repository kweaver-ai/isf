import { consolehttp } from '../../../openapiconsole'
import {
    SearchInOrgTree,
} from './types'

/**
 * 获取策略列表
 */
export const searchInOrgTree: SearchInOrgTree = (query, options?) => {
    return consolehttp('get', ['user-management', 'v1', 'console', 'search-in-org-tree'], undefined, query, options)
}