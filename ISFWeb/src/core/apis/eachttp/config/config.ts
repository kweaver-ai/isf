import { get as lodashGet } from 'lodash'
import { eachttp, CacheableOpenAPIFactory } from '../../../openapi/openapi';
import { CacheableConsoleAPIFactory } from '../../../openapiconsole'
import { Config, ConfigWithoutAuth } from './helper'

/**
 * 涉密配置接口（需要token验证）
 */
const getConfidentialConfigCache: Core.APIs.EACHTTP.Config.GetConfidentialConfigCache = CacheableConsoleAPIFactory('get', ['confidential', 'v1', 'configuration-login'], 60 * 1000 * 5)

/**
 * 涉密配置接口（不需要token验证）
 */
const getConfidentialConfigCacheWithoutAuth: Core.APIs.EACHTTP.Config.GetConfidentialConfigCache = CacheableConsoleAPIFactory('get', ['confidential', 'v1', 'configuration-logout'], 60 * 1000 * 5)

/**
 * 获取配置信息
 */
export const get = CacheableOpenAPIFactory(eachttp, 'config', 'get', { expires: 60 * 1000 }).bind(this)

/**
 * 获取涉密配置
 * @withoutAuth 是否不需要token验证
 */
export const getConfidential = async (withoutAuth = false, item?): Promise<any> => {
    let config: any

    try {
        const confidentialConfigCache = {
            "disabled_roles": [],
            "deployment_console_link_disabled": false,
            "bound_device_disabled": false,
            "reset_via_SMS_disabled": false,
            "domain_management_disabled": false,
            "third_party_authentication_disabled": false,
            "third_party_messages_disabled": false,
            "sensitive_word_control_disabled": false,
            "share_with_anyone_disabled": false,
            "disabled_login_auth": [],
            "pwd_policy": {
                "weak_pwd_disabled": false,
                "max_expire_time": -1,
                "min_strong_pwd_length": 8,
                "max_err_count": 99,
                "min_lock_time": 10
            },
            "application_management_disabled": false,
            "disabled_endpoint_types": [],
            "disabled_doclib_admins": [
                "security",
                "portal_admin",
                "audit_admin",
                "org_audit"
            ],
            "shared_text_disabled": false,
            "knowledgelib_disabled": false,
            "log_storage_space_disabled": false,
            "disabled_pwd_controllers_admins": [
                "sys_admin",
                "audit_admin",
                "org_audit"
            ],
            "disabled_accesslog_admins": [
                "sys_admin",
                "org_admin"
            ],
            "subnet_limit_disabled": false,
            "system_security_disabled": false,
            "doclib_icon_disabled": false,
            "disabled_account_control_admins": [
                "security"
            ],
            "client_console_link_disabled": false,
            "kjz_disabled": true,
            "data_dict_disabled": false,
            "protal_management_disabled": false,
            "report_center_disabled": false,
            "feedback_disabled": false,
            "doc_security_policy_disabled": false,
            "protection_level_init_disabled": true,
            "intelligent_search_disabled": false,
            "tag_management_disabled": false,
            "different_devices_login_disabled": false,
            "lockdown_criteria_disabled": false,
            "permission_request_disabled": false,
            "modify_folder_properties_disabled": true,
            "only_file_rename": false
        }

        const confidentialConfigCacheWidthoutAuth = {
            "internet_identity": {
                "user_agreement_disabled": false,
                "privacy_policy_disabled": false
            },
            "disabled_languages": []
        }
        config = withoutAuth
            ? confidentialConfigCacheWidthoutAuth
            : confidentialConfigCache
    } catch (error) {
        config = withoutAuth
            ? ConfigWithoutAuth
            : Config
    }

    config = item !== undefined ? lodashGet(config, item) : config

    return config
}

/**
 * 获取涉密配置信息（需要token验证）
 * @param [item] {string} 要获取的指定项，使用.进行深度搜索，如get('disabledRoles')
 * @return {Promise} 返回item 值
 */
export const getConfidentialConfig = (item?: string): Promise<Core.APIs.EACHTTP.Config | any> => {
    return getConfidential(false, item)
}

/**
 * 获取涉密配置信息（不需要token验证）
 */
export const getConfidentialConfigWithoutAuth = (item?: string): Promise<Core.APIs.EACHTTP.Config | any> => {
    return getConfidential(true, item)
}