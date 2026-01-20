/**
 * HTTPConfig
 */
export interface HTTPConfig {
    /**
     * 开放http协议
     */
    protocol: string;

    /**
     * 开放api域名
     */
    host: string;

    /**
     * 开放api端口
     */
    port: string;

    /**
     * 访问前缀
     */
    prefix: string;

    /**
     * token
     */
    getToken: () => string;

    /**
     * 刷新token
     */
    refreshToken?: () => Promise<OAuth2Token | null | undefined>;

    /**
     * token 过期
     */
    onTokenExpired?: () => void;

    /**
     * 服务器无法连接
     */
    onServerError?: () => void;

    /**
     * 网络无法连接
     */
    onNetworkError?: () => void;

    /**
     * 超时
     */
    onTimeout?: () => void;
}

interface OAuth2Token {
    access_token: string;
    expires_in?: number;
    id_token: string;
    refresh_token?: string;
    scope?: string;
    token_type?: string;
}