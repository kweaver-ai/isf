declare namespace core {
    namespace APIs {
        namespace Console {
            namespace LoginSecurityPolicy {
                /**
                 * 获取访问者网段绑定功能状态
                 */
                type GetPloicyInfo = Core.APIs.OpenAPI<{

                    /**
                     * 策略模式
                     */
                    mode: string;
                    /**
                     * 策略名称
                     */
                    name: string;

                }, Core.APIs.Console.PolicyState>;

                /**
                 * 获取密码强度信息
                 */
                type SetPwdStrengthMeter = Core.APIs.OpenAPI<{
                    /**
                     * 策略名称
                     */
                    name: string;

                    /**
                     * 设置策略的值
                     */
                    value: Core.APIs.Console.PasswordPolicy;
                }, void>

                /**
                 * 批量设置指定设备类型禁止登录状态
                 */
                type SetBatchOSTypeForbidLoginInfo = Core.APIs.OpenAPI<{
                    /**
                     * 策略名称
                     */
                    name: string;

                    /**
                     * 设置策略的值
                     */
                    value: Core.APIs.Console.OSTypeForbidLoginInfo;
                }, void>

                /**
                 * 设置系统保护等级
                 */
                type SetSystemProtectionLevels = Core.APIs.OpenAPI<{
                    /**
                     * 策略名称
                     */
                    name: string;

                    /**
                     * 设置策略的值
                     */
                    value: {
                        level: Core.APIs.Console.SysProtectionLevel,
                    };
                }, void>
            }
        }
    }
}