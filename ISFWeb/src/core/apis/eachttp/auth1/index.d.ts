declare namespace Core {
    namespace APIs {
        namespace EACHTTP {
            namespace Auth1 {

                /**
                 * 验证码数据
                 */
                type VcodeInfo = {
                    /**
                     * 验证码唯一标识 
                     */
                    uuid: string;

                    /**
                     * 验证码内容
                     */
                    vcode: string;

                    /**
                     * 时候是修改密码
                     */
                    ismodif: boolean;
                }

                /**
                 * 获取getconfig配置
                 */
                type GetConfig = Core.APIs.OpenAPI<void | null, Core.APIs.EACHTTP.Config>

                /**
                * 修改用户密码
                */
                type ModifyPassword = Core.APIs.OpenAPI<{
                    /**
                     * 用户登录名
                     */
                    account: string;

                    /**
                     * 邮箱
                     */
                    emailaddress?: string;

                    /**
                     * 手机号
                     */
                    telnumber?: string,

                    /**
                     * 用户旧密码
                     */
                    oldpwd: string;

                    /**
                     * 用户新密码
                     */
                    newpwd: string;

                    /**
                     * 是否忘记密码
                     */
                    isforgetpwd: boolean;

                    /**
                     * 验证码
                     */
                    vcodeinfo?: VcodeInfo
                }, void>

                /**
                 * 获取验证码
                 */
                type GetVcode = Core.APIs.OpenAPI<{
                    /**
                     * 验证码标识
                     */
                    uuid: string;
                },
                    Core.APIs.EACHTTP.VcodeInfo>

                /**
                 * 获取服务器时间
                 */
                type ServerTime = Core.APIs.EACHTTP<void, Core.APIs.EACHTTP.ServerTime>
            }

        }
    }
}