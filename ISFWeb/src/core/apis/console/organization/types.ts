import { UserRole } from '../../../role/role';
import { OpenAPI } from '../../index';

export enum OrgType {
    User = 'user',
    Dep = 'department',
}

export interface SearchResInfo {
    id: string;
    name: string;
    account?: string;
    type: OrgType | any;
    parent_dep_paths?: string[];
    parent_dep_path?: string;
}

// ======================================================================= 函数 =============================================================

/**
 * 获取策略列表
 */
export type SearchInOrgTree = OpenAPI<{
    keyword: string;
    role: UserRole;
    type: OrgType[];
    offset?: number;
    limit?: number;

    /**
     * 是否过滤禁用用户
     */
    user_enabled?: boolean;

    /**
     * 是否过滤未分配用户
     */
    user_assigned?: boolean;
}, {
        users: {
            entries: SearchResInfo[];
            total_count: number;
        };
        departments: {
            entries: SearchResInfo[];
            total_count: number;
        };
    }>;