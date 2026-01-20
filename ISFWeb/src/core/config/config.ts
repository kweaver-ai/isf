import { get } from 'lodash';
import { getConfig as getAuthConfig, getLoginConfig as getLoginConfigNoAuth } from '../apis/eachttp/auth1/auth1';

/**
 * 获取配置文件-服务器配置信息
 * @param [item] {string} 要获取的指定项，使用.进行深度搜索，如get('thirdauth.id')
 * @return {Promise} 返回item 值，如果未指定item ，则返回整个JSON对象
 */
export function getConfig(item?): Promise<Core.APIs.EACHTTP.Config | any> {
    return getAuthConfig().then((config) => item !== undefined ? get(config, item) : config);
}

/**
 * 获取配置文件-登录配置信息(无鉴权)
 * @param [item] {string} 要获取的指定项，使用.进行深度搜索，如get('thirdauth.id')
 * @return {Promise} 返回item 值，如果未指定item ，则返回整个JSON对象
 */
export function getLoginConfig(item?): Promise<Core.APIs.EACHTTP.Config | any> {
    return getLoginConfigNoAuth().then((config) => item !== undefined ? get(config, item) : config);
}

/**
 * 获取OEM配置
 * @param [item] {string} 返回指定项
 */
export function getOEMConfig(item?) {
    return getConfig('oemconfig').then((config = {}) => {
        return item !== undefined ? get(config, item) : config;
    });
}

/**
 * 获取第三方认证配置
 * @param [item] {string} 返回指定项
 */
export function getThirdAuth() {
    return getLoginConfig('thirdauth');
}