import { consolehttp } from '../../../openapiconsole';

/**
 * 查询第三方消息配置
 */
export const getThirdMessage: Core.APIs.Console.ThirdMessage.GetThirdMessage = (options?) => {
    return consolehttp('get', ['thirdparty-message-plugin', 'v1', 'configs'], undefined, {}, options)
}

/**
 * 添加第三方消息配置
 */
export const addThirdMessage: Core.APIs.Console.ThirdMessage.AddThirdMessage = ({ thirdparty_name, enabled, class_name, channels, config }, options?) => {
    return consolehttp('post', ['thirdparty-message-plugin', 'v1', 'configs'], { thirdparty_name, enabled, class_name, channels, config }, {}, options)
}

/**
 * 修改第三方消息配置
 */
export const editThirdMessage: Core.APIs.Console.ThirdMessage.EditThirdMessage = ({ id, thirdparty_name, enabled, class_name, channels, config }, options?) => {
    return consolehttp('put', ['thirdparty-message-plugin', 'v1', 'configs', id], { thirdparty_name, enabled, class_name, channels, config }, {}, options)
}

/**
 * 删除第三方消息配置
 */
export const deleteThirdMessage: Core.APIs.Console.ThirdMessage.DeleteThirdMessage = ({ id }, options?) => {
    return consolehttp('delete', ['thirdparty-message-plugin', 'v1', 'configs', id], undefined, {}, options)
}

/**
 * 上传第三方插件
 */
export const uploadThirdMessagePlugin: Core.APIs.Console.ThirdMessage.UploadThirdMessagePlugin = ({ id, data }, options?) => {
    return consolehttp('put', ['thirdparty-message-plugin', 'v1', 'plugins', id], data, {}, options)
}