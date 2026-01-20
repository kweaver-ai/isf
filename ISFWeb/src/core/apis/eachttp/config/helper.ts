/**
 * 通用模式配置
 */
export const Config = {
    disabled_roles: [],
    deployment_console_link_disabled: false,
    bound_device_disabled: false,
    reset_via_SMS_disabled: false,
    domain_management_disabled: false,
    third_party_authentication_disabled: false,
    third_party_messages_disabled: false,
    sensitive_word_control_disabled: false,
    share_with_anyone_disabled: false,
    disabled_login_auth: [],
    pwd_policy: {
        weak_pwd_disabled: false,
        max_expire_time: -1,
        min_strong_pwd_length: 8,
        max_err_count: 99,
        min_lock_time: 10,
    },
    application_management_disabled: false,
    disabled_endpoint_types: [],
    disabled_doclib_admins: [
        'security',
        'audit_admin',
        'org_audit',
    ],
    shared_text_disabled: false,
    knowledgelib_disabled: false,
    log_storage_space_disabled: false,
    disabled_pwd_controllers_admins: [
        'sys_admin',
        'audit_admin',
        'org_audit',
    ],
    disabled_accesslog_admins: [
        'sys_admin',
        'org_admin',
    ],
    subnet_limit_disabled: false,
    system_security_disabled: false,
    doclib_icon_disabled: false,
    disabled_account_control_admins: [
        'security',
    ],
    client_console_link_disabled: false,
    kjz_disabled: true,
    data_dict_disabled: false,
    protal_management_disabled: false,
    report_center_disabled: false,
    feedback_disabled: false,
    doc_security_policy_disabled: false,
    permission_request_disabled: false,
    modify_folder_properties_disabled: true,
    only_file_rename: false,
    protection_level_init_disabled: true,
    intelligent_search_disabled: false,
    tag_management_disabled: false,
    different_devices_login_disabled: false,
    lockdown_criteria_disabled: false,
}

export const ConfigWithoutAuth = {
    internet_identity: {
        user_agreement_disabled: false,
        privacy_policy_disabled: false,
    },
    disabled_languages: [],
}