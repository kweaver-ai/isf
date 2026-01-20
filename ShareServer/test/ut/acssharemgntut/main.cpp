#include <abprec.h>
#include <ncutil/ncutil.h>
#include <gtest/gtest.h>
#include <gmock/gmock.h>
#include <dboperator/public/ncIDBOperator.h>

void InitDlls ()
{
    // 加载 cfl 库
    static AppContext appCtx (_T("acs_sharemgnt_ut"));
    AppContext::setInstance (&appCtx);
    AppSettings* appSettings = AppSettings::getCFLAppSettings ();
    LibManager::getInstance ()->initLibs (appSettings, &appCtx, 0);

    // xpcom 核心库初始化
    ::ncInitXPCOM ();

    // 开启cfl的异常输出
    abEnableOutputError ();
}

void CreateShareMgntDB (ncIDBOperator* dbOper)
{
    String strSql = _T("drop database if exists sharemgnt_db");
    dbOper->Execute (strSql);

    strSql = _T("create DATABASE if not exists sharemgnt_db CHARACTER SET = 'utf8mb4' COLLATE = 'utf8mb4_bin'");
    dbOper->Execute (strSql);
}

void CreateShareMgntTables (ncIDBOperator* dbOper)
{
    String strSql;

    // select sharemgnt
    strSql = _T("use sharemgnt_db;");
    dbOper->Execute (strSql);

    // t_department
    strSql = _T("CREATE TABLE if not exists `t_department` (")    \
        _T("`f_department_id` char(40) NOT NULL,")    \
        _T("`f_auth_type` tinyint(4) NOT NULL DEFAULT '1',")    \
        _T("`f_priority` smallint(6) NOT NULL DEFAULT '999',")    \
        _T("`f_name` char(100) NOT NULL,")    \
        _T("`f_domain_object_guid` char(100) NOT NULL DEFAULT '',")    \
        _T("`f_domain_path` char(100) NOT NULL DEFAULT '',")    \
        _T("`f_is_enterprise` tinyint(4) NOT NULL DEFAULT '0',")    \
        _T("`f_remote_dep_id` text,")    \
        _T("`f_mail_address` char(150) NOT NULL DEFAULT '',")    \
        _T("PRIMARY KEY (`f_department_id`),")    \
        _T("UNIQUE KEY `f_department_id_index` (`f_department_id`),")    \
        _T("KEY `f_name_index` (`f_name`),")    \
        _T("KEY `f_is_enterprise_index` (`f_is_enterprise`),")    \
        _T("KEY `f_mail_address_index` (`f_mail_address`)")    \
        _T(") ENGINE=InnoDB;");
    dbOper->Execute (strSql);

    // t_department_relation
    strSql = _T("CREATE TABLE if not exists `t_department_relation` (")    \
        _T("`f_relation_id` bigint(20) NOT NULL AUTO_INCREMENT,")    \
        _T("`f_department_id` char(40) NOT NULL,")    \
        _T("`f_parent_department_id` char(40) NOT NULL,")    \
        _T("PRIMARY KEY (`f_relation_id`),")    \
        _T("UNIQUE KEY `f_department_id_index` (`f_department_id`),")    \
        _T("KEY `f_parent_department_id_index` (`f_parent_department_id`)")    \
        _T(") ENGINE=InnoDB;");
    dbOper->Execute (strSql);

    // t_user
    strSql = _T("CREATE TABLE if not exists `t_user` (")
        _T("`f_user_id` char(40) NOT NULL,")
        _T("`f_login_name` char(150) NOT NULL,")
        _T("`f_display_name` char(150) NOT NULL,")
        _T("`f_password` char(50) DEFAULT '',")
        _T("`f_mail_address` char(150) NOT NULL DEFAULT '',")
        _T("`f_auth_type` tinyint(4) NOT NULL DEFAULT '1',")
        _T("`f_status` tinyint(4) NOT NULL DEFAULT '0',")
        _T("`f_domain_object_guid` char(100) DEFAULT '',")
        _T("`f_third_party_id` char(40) DEFAULT NULL,")
        _T("`f_third_party_depart_id` char(255) DEFAULT NULL,")
        _T("`f_priority` smallint(6) NOT NULL DEFAULT '999',")
        _T("`f_csf_level` tinyint(4) NOT NULL DEFAULT '1',")
        _T("`f_freeze_status` tinyint(4) NOT NULL DEFAULT '0',")
        _T("`f_agreed_to_terms_of_use` tinyint(4) NOT NULL DEFAULT '0',")
        _T("`f_tel_number` char(40) DEFAULT NULL,")
        _T("`f_is_activate` tinyint(4) NOT NULL DEFAULT '0',")
        _T("`f_pwd_control` tinyint(1) NOT NULL DEFAULT '0',")
        _T("`f_user_document_read_status` bigint(20) DEFAULT '0',")
        _T("`f_auto_disable_status` tinyint(4) NOT NULL DEFAULT '0',")
        _T("`f_last_request_time` datetime DEFAULT NULL,")
        _T("`f_last_client_request_time` datetime NOT NULL DEFAULT now(),")
        _T("`f_third_party_attr` varchar(255) NOT NULL DEFAULT '',")
        _T("PRIMARY KEY (`f_user_id`),")
        _T("UNIQUE KEY `f_login_name` (`f_login_name`),")
        _T("KEY `f_mail_address_index` (`f_mail_address`),")
        _T("KEY `f_domain_object_guid` (`f_domain_object_guid`),")
        _T("UNIQUE KEY `f_tel_number_index` (`f_tel_number`)")
        _T(") ENGINE=InnoDB;");
    dbOper->Execute (strSql);

    // t_user_department_relation
    strSql = _T("CREATE TABLE if not exists `t_user_department_relation` (")    \
        _T("`f_relation_id` bigint(20) NOT NULL AUTO_INCREMENT,")    \
        _T("`f_user_id` char(40) NOT NULL,")    \
        _T("`f_department_id` char(40) NOT NULL,")    \
        _T("PRIMARY KEY (`f_relation_id`),")    \
        _T("KEY `f_user_id_index` (`f_user_id`),")    \
        _T("KEY `f_department_id_index` (`f_department_id`)")    \
        _T(") ENGINE=InnoDB;");
    dbOper->Execute (strSql);

    // t_person_group
    strSql = _T("CREATE TABLE IF NOT EXISTS `t_person_group` (")    \
        _T("`f_group_id` char(40) NOT NULL,")    \
        _T("`f_user_id` char(40) NOT NULL,")    \
        _T("`f_group_name` char(100) NOT NULL,")    \
        _T("`f_person_count` bigint(20) NOT NULL DEFAULT '0',")    \
        _T("PRIMARY KEY (`f_group_id`),")    \
        _T("KEY `f_user_id` (`f_user_id`)")    \
        _T(") ENGINE=InnoDB;");
    dbOper->Execute (strSql);

    // t_contact_person
    strSql = _T("CREATE TABLE IF NOT EXISTS `t_contact_person` (")    \
        _T("`f_id` int(11) NOT NULL AUTO_INCREMENT,")    \
        _T("`f_group_id` char(40) NOT NULL,")    \
        _T("`f_user_id` char(40) NOT NULL,")    \
        _T("PRIMARY KEY (`f_id`),")    \
        _T("KEY `f_group_id` (`f_group_id`,`f_user_id`),")    \
        _T("KEY `f_user_id` (`f_user_id`)")    \
        _T(") ENGINE=InnoDB;");
    dbOper->Execute (strSql);

    // t_ou_user
    strSql = _T("CREATE TABLE IF NOT EXISTS `t_ou_user` (")    \
        _T("`f_id` int(11) NOT NULL AUTO_INCREMENT,")    \
        _T("`f_user_id` char(40) NOT NULL,")    \
        _T("`f_ou_id` char(40) NOT NULL,")    \
        _T("PRIMARY KEY (`f_id`),")    \
        _T("KEY `f_user_id_index` (`f_user_id`),")    \
        _T("KEY `f_ou_id_index` (`f_ou_id`)")    \
        _T(") ENGINE=InnoDB;");
    dbOper->Execute (strSql);

    // t_ou_department
    strSql = _T("CREATE TABLE IF NOT EXISTS `t_ou_department` (")    \
        _T("`f_id` int(11) NOT NULL AUTO_INCREMENT,")    \
        _T("`f_department_id` char(40) NOT NULL,")    \
        _T("`f_ou_id` char(40) NOT NULL,")    \
        _T("PRIMARY KEY (`f_id`),")    \
        _T("KEY `f_department_id_index` (`f_department_id`),")    \
        _T("KEY `f_ou_id_index` (`f_ou_id`)")    \
        _T(") ENGINE=InnoDB;");
    dbOper->Execute (strSql);

    // t_third_party_auth
    strSql = _T("CREATE TABLE IF NOT EXISTS `t_third_party_auth` (")    \
        _T("`f_id` int(11) NOT NULL AUTO_INCREMENT,")    \
        _T("`f_app_id` varchar(50) NOT NULL,")    \
        _T("`f_app_name` varchar(128) NOT NULL DEFAULT '',")    \
        _T("`f_enable` tinyint(1) NOT NULL DEFAULT 0,")    \
        _T("`f_config` text,")    \
        _T("PRIMARY KEY (`f_id`)")    \
        _T(") ENGINE=InnoDB;");
    dbOper->Execute (strSql);

    // t_third_auth_info
    strSql = _T("CREATE TABLE IF NOT EXISTS `t_third_auth_info` (")    \
        _T("`f_app_id` varchar(50) NOT NULL,")    \
        _T("`f_app_key` char(36) NOT NULL,")    \
        _T("`f_enabled` tinyint(1) NOT NULL DEFAULT '1',")    \
        _T("PRIMARY KEY (`f_app_id`)")    \
        _T(") ENGINE=InnoDB;");
    dbOper->Execute (strSql);

    // t_department_responsible_person
    strSql = _T("CREATE TABLE IF NOT EXISTS `t_department_responsible_person` (")    \
        _T("`f_department_id` char(40) NOT NULL,")    \
        _T("`f_user_id` char(40) NOT NULL,")    \
        _T("UNIQUE KEY `responsible_person_depart_index` (`f_user_id`,`f_department_id`)")    \
        _T(") ENGINE=InnoDB;");
    dbOper->Execute (strSql);

    // t_sharemgnt_config
    strSql = _T("CREATE TABLE IF NOT EXISTS `t_sharemgnt_config` (")    \
        _T("`f_key` char(32) NOT NULL,")    \
        _T("`f_value` varchar(1024) NOT NULL,")    \
        _T("PRIMARY KEY (`f_key`)")    \
        _T(") ENGINE=InnoDB;");
    dbOper->Execute (strSql);

    // t_user_role_relation
    strSql = _T("CREATE TABLE IF NOT EXISTS `t_user_role_relation` (")  \
        _T("`f_user_id` char(40) NOT NULL,")    \
        _T("`f_role_id` char(40) NOT NULL,")    \
        _T("PRIMARY KEY (`f_user_id`, `f_role_id`)")    \
        _T(") ENGINE=InnoDB;");
    dbOper->Execute (strSql);

    strSql =  _T("CREATE TABLE IF NOT EXISTS `t_oem_config` (")
        _T("`f_section` char(32) NOT NULL,")    \
        _T("`f_option` char(32) NOT NULL,")    \
        _T("`f_value` mediumblob NOT NULL,")    \
        _T(" UNIQUE KEY `f_index_section_option` (`f_section`,`f_option`) USING BTREE")    \
        _T(") ENGINE=InnoDB;");
    dbOper->Execute (strSql);
}

void InitShareMgntDB ()
{
    ncDBConnectionInfo connInfo;
    connInfo.ip = _T("127.0.0.1");
    connInfo.port = 3306;
    connInfo.user = _T("root");
    connInfo.password = _T("eisoo.com");

    nsresult ret;
    nsCOMPtr<ncIDBOperator> dbOper = do_CreateInstance (NC_MYSQL_OPERATOR_CID, &ret);
    if (NS_FAILED (ret)) {
        String error;
        error.format (_T("Failed to create db operator: 0x%x"), ret);

        throw Exception (error);
    }
    dbOper->Connect (connInfo);

    CreateShareMgntDB (dbOper);

    CreateShareMgntTables (dbOper);
}

int main(int argc, char** argv)
{
    try {
        // 初始化基础库
        InitDlls ();

        // 初始化 gtest 参数
        testing::InitGoogleTest(&argc, argv);

        // 初始化 gmock 参数
        testing::InitGoogleMock(&argc, argv);

        InitShareMgntDB ();

        // 运行测试
        return RUN_ALL_TESTS ();
    }
    catch (Exception& e) {
        printMessage2 (_T("Test Error: %s"), e.toFullString ().getCStr ());
        return 1;
    }
    catch (...) {
        printMessage2 (_T("Test Error: Unknown."));
        return 1;
    }

    return 0;
}
