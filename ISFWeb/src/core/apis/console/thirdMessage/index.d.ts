declare namespace Core {
    namespace APIs {
        namespace Console {
            namespace ThirdMessage {
                /**
                 * 添加第三方消息配置
                 */
                type AddThirdMessage = Core.APIs.OpenAPI<{
                    /**
                     * 第三方app名字
                     */
                    thirdparty_name: string;

                    /**
                     * 第三方配置开关
                     */
                    enabled: boolean;

                    /**
                     * 消息插件类名
                     */
                    class_name: string;

                    /**
                     * 消息类型
                     */
                    channels: ReadonlyArray<string>;

                    /**
                     * 插件需要其他配置，透传给第三方插件
                     */
                    config: Record<string, any>;
                }, Core.APIs.Console.AddThirdMessageResult>;

                /**
                 * 删除第三方消息配置
                 */
                type DeleteThirdMessage = Core.APIs.OpenAPI<{
                    /**
                     * 插件id
                     */
                    id: string;
                }, void>;

                /**
                 * 修改第三方消息配置
                 */
                type EditThirdMessage = Core.APIs.OpenAPI<{
                    /**
                     * 插件id
                     */
                    id: string;

                    /**
                     * 第三方app名字
                     */
                    thirdparty_name: string;

                    /**
                     * 第三方配置开关
                     */
                    enabled: boolean;

                    /**
                     * 消息插件类名
                     */
                    class_name: string;

                    /**
                     * 消息类型
                     */
                    channels: ReadonlyArray<string>;

                    /**
                     * 插件需要其他配置，透传给第三方插件
                     */
                    config: Record<string, any>;
                }, void>;

                /**
                 * 查询第三方消息配置
                 */
                type GetThirdMessage = Core.APIs.OpenAPI<
                    any, Core.APIs.Console.GetThirdMessageResult>;

                /**
                 * 上传第三方插件
                 */
                type UploadThirdMessagePlugin = Core.APIs.OpenAPI<{
                    /**
                     * 插件id
                     */
                    id: string;

                    /**
                     * 上传插件数据
                     */
                    data: any;
                }, void>
            }
        }
    }
}