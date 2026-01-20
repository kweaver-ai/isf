import { CacheableConsoleAPIFactory } from '../../../openapiconsole'

/**
 * 获取服务器配置信息
 */
export const getConfig: Core.APIs.EACHTTP.Auth1.GetConfig = CacheableConsoleAPIFactory('get', ['eacp', 'v1', 'auth1', 'configs'], 10 * 60 * 1000)

/**
 * 获取登录配置信息(无鉴权)
 */
export const getLoginConfig: Core.APIs.EACHTTP.Auth1.GetConfig = CacheableConsoleAPIFactory('get', ['eacp', 'v1', 'auth1', 'login-configs'], 10 * 60 * 1000)
