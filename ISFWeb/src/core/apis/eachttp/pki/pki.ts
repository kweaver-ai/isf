import { eachttp } from '../../../openapi/openapi';

/**
 * 获取original
 */
export const original: Core.APIs.EACHTTP.PKI.Original = function ({ } = {}, options?) {
    return eachttp('pki', 'original', {}, { userid: null, tokenid: null });
}

/**
 * 使用PKI认证
 */
export const authen: Core.APIs.EACHTTP.PKI.Authen = function ({ original, detach }, options?) {
    return eachttp('pki', 'authen', { original, detach }, { userid: null, tokenid: null });
}

