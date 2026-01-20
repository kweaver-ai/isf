export declare namespace UserGroup {
    /**
     * 用户组信息
     */
    interface GroupInfo {
        /**
         * 用户组id
         */
        id: string;

        /**
         * 用户组名称
         */
        name: string;

        /**
         * 用户组备注
         */
        notes: string;
    }

    /**
     * 成员信息
     */
    interface MemberInfo {
        /**
         * 成员显示名
         */
        name: string;

        /**
         * 成员id
         */
        id: string;

        /**
         * 成员类型, user, department
         */
        type: string;

        /**
         * 成员直属部门
         */
        department_names: ReadonlyArray<string>;
    }

    /**
     * 错误信息
     */
    interface ErrorInfo {
        /**
         * 错误码
         */
        code: number | string;

        /**
         * 错误信息
         */
        message: string;

        /**
         * 错误原因
         */
        cause: string;

        /**
         * 详情
         */
        detail?: any;
    }
}

/**
 * 默认页码
 */
export const DefaultPage = 1;

/**
 * 限制条数
 */
export const Limit = 50;

/**
 * 成员类型
 */
export enum MemberType {
    /**
     * 用户
     */
    User = 'user',

    /**
     * 部门/组织
     */
    Dep = 'department',
}