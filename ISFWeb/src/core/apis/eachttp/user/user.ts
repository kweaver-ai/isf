import { eachttp, CacheableOpenAPIFactory } from '../../../openapi/openapi';

/**
 * 获取用户信息
 */
export const get: Core.APIs.EACHTTP.User.GetUser = CacheableOpenAPIFactory(eachttp, 'user', 'get', { expires: 60 * 1000 })

/**
 * 编辑用户信息
 */
export const edit: Core.APIs.EACHTTP.User.Edit = function ({ emailaddress, displayname, telnumber }, options?) {
    return eachttp('user', 'edit', { emailaddress, displayname, telnumber }, options);
}
