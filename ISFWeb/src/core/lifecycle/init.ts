import { setupOpenApi } from '@/util/http'
import { setup as setupOpenApiConsole } from '../openapiconsole/index'
import { setup as setupThrift } from '../thrift/thrift'
import { apply } from '@/util/skin'
import i18n from '@/core/i18n'
import { initializeI18n } from '@/util/i18n';
import { apis } from '@dip/components/dist/dip-components.min.js';
import __ from './locale'
import { message } from 'antd';

const initData = ({
    theme,
    lang,
    protocol = 'https:',
    host = location.hostname,
    port = 443,
    prefix = '',
    getToken,
    refreshToken,
    onTokenExpired,
    container
}) => {
    i18n.setup({
        locale: lang,
    })
    initializeI18n(lang)
    const data = {
        protocol,
        host,
        port,
        prefix,
        getToken,
        refreshToken,
        onTokenExpired,
        theme,
        lang,
        onNetworkError: () => {
            message.warning(__('无法连接网络'))
        },
        onServerError: () => {
            message.warning(__('服务异常'))
        },
        onTimeout: () => {
            message.warning(__('您的请求已超时'))
        },
    }
    setupOpenApiConsole(data)
    setupOpenApi(data)
    setupThrift(data)
    
    // 设置dip-components所需的信息
    apis.setup({
        protocol,
        host,
        port,
        lang,
        prefix,
        getToken,
        refreshToken,
        onTokenExpired,
        theme,
        popupContainer: container,
    });

    // 设置主题色
    const removeThemeStyle = apply(theme)

    return {
        removeThemeStyle,
    }
}

export {
    initData,
}
