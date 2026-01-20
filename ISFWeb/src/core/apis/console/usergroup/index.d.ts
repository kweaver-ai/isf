declare namespace Core {
    namespace APIs {
        namespace Console {
            namespace UserGroup {
                /**
                 * 获取用户组列表
                 */
                type GetUserGroups = Core.APIs.OpenAPI<
                    {
                        /**
                         * 排序结果方向，默认为降序
                         */
                        direction?: string;

                        /**
                         * 排序类型，默认按创建时间排序
                         */
                        sort?: string;

                        /**
                         * 获取数据起始下标
                         */
                        offset?: number;

                        /**
                         * 获取数据量
                         */
                        limit?: number;

                        /**
                         * 搜索关键字
                         */
                        keyword?: string;
                    },
                    {
                        /**
                         * 符合搜索条件的用户组
                         */
                        entries: ReadonlyArray<Core.APIs.Console.UserGroupInfo>;

                        /**
                         * 文档库总数
                         */
                        total_count: number;
                    }
                >;

                /**
                 * 根据id获取用户组详情
                 */
                type GetUserGroupById = Core.APIs.OpenAPI<
                    {
                        /**
                         * 用户组id
                         */
                        id: string;
                    },
                    Core.APIs.Console.UserGroupInfo
                >;

                /**
                 * 创建用户组
                 */
                type CreateUserGroup = Core.APIs.OpenAPI<
                    {
                        /**
                         * 用户组名称
                         */
                        name: string;

                        /**
                         * 备注
                         */
                        notes?: string;

                        /**
                         * 初始组成员所在组ID数组，默认为空，如果不为空，则会在创建组之后，将他们的组成员加入到新建组中
                         */
                        group_ids_of_members?: ReadonlyArray<string>;
                    },
                    {
                        /**
                         * 创建的用户组的id
                         */
                        id: string;
                    }
                >;

                /**
                 * 编辑用户组
                 */
                type EditUserGroup = Core.APIs.OpenAPI<
                    {
                        /**
                         * id
                         */
                        id: string;

                        /**
                         * 修改的内容名称
                         */
                        fields: ReadonlyArray<string>;

                        /**
                         * 修改后的用户组名称
                         */
                        name: string;

                        /**
                         * 修改后的用户组备注
                         */
                        notes: string;
                    },
                    void
                >;

                /**
                 * 删除用户组
                 */
                type DeleteUserGroup = Core.APIs.OpenAPI<
                    {
                        /**
                         * 用户组id
                         */
                        id: string;
                    },
                    void
                >;

                /**
                 * 获取用户组成员列表
                 */
                type GetGroupMembers = Core.APIs.OpenAPI<
                    {
                        /**
                         * 用户组id
                         */
                        id: string;

                        /**
                         * 排序结果方向，默认为降序
                         */
                        direction?: string;

                        /**
                         * 排序类型，默认按创建时间排序
                         */
                        sort?: string;

                        /**
                         * 获取数据起始下标
                         */
                        offset?: number;

                        /**
                         * 获取数据量
                         */
                        limit?: number;

                        /**
                         * 搜索关键字
                         */
                        keyword?: string;
                    },
                    {
                        /**
                         * 符合搜索条件的成员
                         */
                        entries: ReadonlyArray<Core.APIs.Console.GroupMemberInfo>;

                        /**
                         * 文档库总数
                         */
                        total_count: number;
                    }
                >;

                /**
                 * 获取用户组成员(搜索精确匹配用户名)
                 */
                type GetGroupMembersByUserMatch = Core.APIs.OpenAPI<
                    {
                        /**
                         * 用户组id
                         */
                        group_id: string;

                        /**
                         * 搜索字段（显示名）
                         */
                        key: string;
                    },
                    {
                        /**
                         * 用户条目列表
                         */
                        entries: ReadonlyArray<{
                            /**
                             * 用户id
                             */
                            id: string;

                            /**
                             * 用户名
                             */
                            name: string;

                            /**
                             * 类型
                             */
                            type: string;

                            /**
                             * 父部门信息，描述多个父部门的层级关系，每个父部门星际数组内第一个对象是跟部门，最后一个是直接父部门
                             */
                            parent_deps: ReadonlyArray<ReadonlyArray<Core.APIs.Console.ParentDeps>>;

                            /**
                             * 所处组成员信息
                             */
                            group_members: ReadonlyArray<Core.APIs.Console.GroupMemberInfo>;
                        }>;

                        /**
                         * 总条目数
                         */
                        total_count: number;
                    }
                >;

                /**
                 * 添加成员
                 */
                type AddGroupMembers = Core.APIs.OpenAPI<
                    {
                        /**
                         * 用户组id
                         */
                        id: string;

                        /**
                         * 成员信息
                         */
                        members: ReadonlyArray<Core.APIs.Console.GroupMemberInfo>;
                    },
                    void
                >;

                /**
                 * 删除成员
                 */
                type DeleteGroupMembers = Core.APIs.OpenAPI<
                    {
                        /**
                         * 用户组id
                         */
                        id: string;

                        /**
                         * 成员信息
                         */
                        members: ReadonlyArray<Core.APIs.Console.GroupMemberInfo>;
                    },
                    void
                >;
            }
        }
    }
}