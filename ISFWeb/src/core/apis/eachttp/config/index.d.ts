declare namespace Core {
    namespace APIs {
        namespace EACHTTP {
            namespace Config {
                /**
                * 获取配置信息
                */
                type Get = Core.APIs.OpenAPI<{}, Core.APIs.EACHTTP.CacheConfigs>

                /**
                 * 涉密配置接口
                 */
                type GetConfidentialConfigCache = Core.APIs.OpenAPI<{}, Core.APIs.EACHTTP.CacheConfigs>

                /**
                 * 获取OEM配置
                 */
                type GetOemConfigBySection = Core.APIs.OpenAPI<{
                    /**
                     * 格式类似：shareweb_zh-cn
                     */
                    section: string;
                }, Core.APIs.EACHTTP.OEMInfo>
            }
        }
    }
}