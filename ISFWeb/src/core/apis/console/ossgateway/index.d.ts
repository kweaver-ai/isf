declare namespace Core {
    namespace APIs {
        namespace Console {
            namespace OssGateWay {
                /**
                 * 注册账户
                 */
                type RegisterAccount = Core.APIs.OpenAPI<
                    {},
                    {
                        /**
                         * 账户id
                         */
                        id: string;

                        /**
                         * 账户密钥
                         */
                        key: string;
                    }
                >;

                /**
                 * 获取账户
                 */
                type GetAccount = Core.APIs.OpenAPI<
                    {},
                    {
                        /**
                         * 账户id
                         */
                        id: string;
                    }
                >;

                /**
                 * 删除账户
                 */
                type DeleteAccount = Core.APIs.OpenAPI<
                    {},
                    void
                >;

                /**
                 * 获取当前站点的默认对象存储
                 */
                type GetDefaultStorage = Core.APIs.OpenAPI<
                    {},
                    { storage_id: string }
                >;

                /**
                 * 获取当前站点可界面管理的对象存储
                 */
                type GetObjectStorageInfoByApp = Core.APIs.OpenAPI<
                    { app: string; enabled?: boolean },
                    ReadonlyArray<OSSInfo>
                >;

                /**
                 * 对象存储服务信息
                 */
                type OSSInfo = {
                    /**
                     *  提供者
                     */
                    provider: string;

                    /**
                     * 是否启用，true 启用；false 禁用
                     */
                    enabled: boolean;

                    /**
                     * 对象存储的标识ID
                     */
                    id: string;

                    /**
                     * 对象存储的名字
                     */
                    name: string;

                    /**
                     * 是否是默认存储
                     */
                    default: boolean;
                }

                /**
                 * 获取指定对象存储的基本信息
                 */
                type GetObjectStorageInfoById = Core.APIs.OpenAPI<
                    string,
                    ncTOSSBaseInfo
                >;

                /**
                 * 对象存储服务基本信息
                 */
                type ncTOSSBaseInfo = {
                    /**
                     * 对象存储ID
                     */
                    id: string;

                    /**
                     * 对象存储名称
                     */
                    name: string;

                    /**
                     * 对象存储状态
                     */
                    enabled: boolean;

                    /**
                     * 是否是默认存储
                     */
                    default: boolean;
                }
            }
        }
    }
}