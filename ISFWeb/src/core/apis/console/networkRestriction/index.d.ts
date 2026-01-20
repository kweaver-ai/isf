declare namespace core {
    namespace APIs {
        namespace Console {
            namespace NetworkRestriction {
                /**
                 * 获取访问者网段绑定功能状态
                 */
                type GetState = Core.APIs.OpenAPI<{
                    /**
                     * 策略名称
                     */
                    name: string;

                }, Core.APIs.Console.PolicyState>;

                /**
                 * 启用和关闭访问者网段绑定功能
                 */
                type SetState = Core.APIs.OpenAPI<{

                    /**
                     * 策略名称
                     */
                    name: string;
                    /**
                     * 策略状态
                     */
                    is_enabled: boolean;
                }, void>;


                /**
                 * 获取网段列表
                 */
                type GetNetworkList = Core.APIs.OpenAPI<{
                    /**
                     * 搜索关键字
                     */
                    key_word?: string;

                    /**
                     * 数据起始下标
                     */
                    offset?: number;

                    /**
                     * 数据量
                     */
                    limit?: number;

                }, Core.APIs.Console.NetworkList>;

                /**
                 * 新增网段
                 */
                type AddNetwork = Core.APIs.OpenAPI<{
                    /**
                    * 网段名称
                    */
                    name: string;

                    /**
                     * 起始IP
                     */
                    start_ip: string;

                    /**
                     * 终止IP
                     */
                    end_ip: string;

                    /**
                     * IP地址
                     */
                    ip_address: string;

                    /**
                     * 子网掩码
                     */
                    netmask: string;

                    /**
                     * IP选项
                     */
                    net_type: Core.APIs.Console.NetType;

                    /**
                     * IP版本
                     */
                    ip_type: IpVersion;

                }, void>

                /**
                 * Ip版本
                 */
                enum IpVersion {
                    /**
                     * ipv4
                     */
                    IPV4 = 'ipv4',
                    /**
                     * ipv6
                     */
                    IPV6 = 'ipv6',
                }

                /**
                 * 获取网段信息
                 */
                type GetNetworkInfo = Core.APIs.OpenAPI<{
                    /**
                     * 网段id
                     */
                    id: string;

                }, Core.APIs.Console.Network>

                /**
                 * 修改网段
                 */
                type EditNetwork = Core.APIs.OpenAPI<{
                    /**
                     * 网段id
                     */
                    id: string;

                    /**
                    * 网段名称
                    */
                    name: string;

                    /**
                     * 起始IP
                     */
                    start_ip: string;

                    /**
                     * 终止IP
                     */
                    end_ip: string;

                    /**
                     * IP地址
                     */
                    ip_address: string;

                    /**
                     * 子网掩码
                     */
                    netmask: string;

                    /**
                     * IP选项
                     */
                    net_type: Core.APIs.Console.NetType;

                    /**
                     * IP版本
                     */
                    ip_type: IpVersion;

                }, void>

                /**
                 * 删除网段
                 */
                type DeleteNetwork = Core.APIs.OpenAPI<{
                    /**
                     * 网段id
                     */
                    id: string;

                }, void>

                /**
                 * 获取网段绑定的访问者列表
                 */
                type GetAccessorsByNetwork = Core.APIs.OpenAPI<{

                    /**
                    * 网段id
                    */
                    id: string;

                    /**
                     * 搜索关键字
                     */
                    key_word?: string;

                    /**
                     * 数据起始下标
                     */
                    offset?: number;

                    /**
                     * 数据量
                     */
                    limit?: number;

                }, Core.APIs.Console.AccessorList>

                /**
                 * 向绑定的网段新增访问者
                 */
                type AddAccessorsByNetwork = Core.APIs.OpenAPI<{
                    /**
                     * 网段id
                     */
                    id: string;

                    /**
                     * 新增的访问者列表
                     */
                    accessorsList?: ReadonlyArray<Core.APIs.Console.AddAccessor>

                }, ReadonlyArray<Core.APIs.Console.Accessor>>

                /**
                 * 删除网段绑定的访问者
                 */
                type DeleteAccessorByNetwork = Core.APIs.OpenAPI<{
                    /**
                     * 网段id
                     */
                    id: string;

                    /**
                     * 访问者id,多个用`,`连接
                     */
                    accessor_id: string;

                }, void>

                /**
                 * 获取访问者已绑定的网段
                 */
                type GetNetworkListByAccessor = Core.APIs.OpenAPI<{
                    /**
                     * 访问者id
                     */
                    id: string;

                    /**
                     * 数据起始下标
                     */
                    offset?: number;

                    /**
                     * 数据量
                     */
                    limit?: number;

                }, Core.APIs.Console.NetworkList>
            }
        }
    }
}