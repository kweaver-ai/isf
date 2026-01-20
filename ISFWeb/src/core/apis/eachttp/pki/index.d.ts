declare namespace Core {
    namespace APIs {
        namespace EACHTTP {
            namespace PKI {

                /**
                 * 获取original
                 */
                type Original = Core.APIs.OpenAPI<void, Core.APIs.EACHTTP.OriginalInfo>

                /**
                 * 使用PKI认证
                 */
                type Authen = Core.APIs.OpenAPI<{
                    /**
                     * 通过pki?method=original获取的值
                     */
                    original: string;

                    /**
                    * 使用key与original计算后的值
                    */
                    detach: string;
                }, Core.APIs.EACHTTP.AuthenInfo>
            }
        }
    }
}