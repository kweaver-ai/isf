declare namespace Core {

    /**
     * 已登录的用户详细信息，包含tokenid、userid
     */
    namespace User {
        interface UserInfo extends Core.APIs.EACHTTP.UserInfo {
            /**
             * userid
             */
            userid: string;

            /**
             * tokenid
             */
            tokenid: string;
        }
    }
}