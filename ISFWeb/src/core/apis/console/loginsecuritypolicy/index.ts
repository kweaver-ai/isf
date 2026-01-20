import { consolehttp } from '../../../openapiconsole'

/**
 * 获取策略信息
 */
export const getPloicyInfo: Core.APIs.Console.LoginSecurityPolicy.GetPloicyInfo = ({ mode, name }, options?) => {
    return consolehttp('get', ['policy-management', 'v1', 'general'], undefined, { mode, name }, options)
}

/**
 * 设置密码策略
 */
export const setPwdStrengthMeter: Core.APIs.Console.LoginSecurityPolicy.SetPwdStrengthMeter = ({ name, value }, options?) => {
    return consolehttp('put', ['policy-management', 'v1', 'general', 'password_strength_meter', 'value'], [{ name, value }], {}, options)
}

/**
 * 批量设置指定设备类型禁止登录状态
 */
export const setBatchOSTypeForbidLoginInfo: Core.APIs.Console.LoginSecurityPolicy.SetBatchOSTypeForbidLoginInfo = ({ name, value }, options?) => {
    return consolehttp('put', ['policy-management', 'v1', 'general', 'client_restriction', 'value'], [{ name, value }], {}, options)
}

/**
 * 设置系统保护等级
 */
export const setSystemProtectionLevels: Core.APIs.Console.LoginSecurityPolicy.SetSystemProtectionLevels = ({ name, value }, options?) => {
    return consolehttp('put', ['policy-management', 'v1', 'general', 'system_protection_levels', 'value'], [{ name, value }], {}, options)
}
