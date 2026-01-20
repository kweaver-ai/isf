#include <abprec.h>
#include <ncutil/ncutil.h>
#include <gtest/gtest.h>
#include <gmock/gmock.h>
#include <dboperator/public/ncIDBOperator.h>

void InitDlls ()
{
    // 加载 cfl 库
    static AppContext appCtx (_T("acs_db_ut"));
    AppContext::setInstance (&appCtx);
    AppSettings* appSettings = AppSettings::getCFLAppSettings ();
    LibManager::getInstance ()->initLibs (appSettings, &appCtx, 0);

    // xpcom 核心库初始化
    ::ncInitXPCOM ();

    // 开启cfl的异常输出
    abEnableOutputError ();
}

void CreateAnyShareDB ()
{
    ncDBConnectionInfo connInfo;
    connInfo.ip = _T("127.0.0.1");
    connInfo.port = 3306;
    connInfo.user = _T("root");
    connInfo.password = _T("eisoo.com");

    nsresult ret;
    nsCOMPtr<ncIDBOperator> dbOperator = do_CreateInstance (NC_MYSQL_OPERATOR_CID, &ret);
    if (NS_FAILED (ret)) {
        String error;
        error.format (_T("Failed to create db operator: 0x%x"), ret);

        throw Exception (error);
    }
    dbOperator->Connect (connInfo);

    String strSql = _T("drop database if exists anyshare");
    dbOperator->Execute (strSql);

    strSql = _T("create DATABASE if not exists anyshare CHARACTER SET = 'utf8mb4' COLLATE = 'utf8mb4_bin'");
    dbOperator->Execute (strSql);

    dbOperator->Execute ("use anyshare");

    // 创建 t_lock
    strSql = "CREATE TABLE IF NOT EXISTS `t_lock` ("
        "`f_primary_id` char(40) NOT NULL,"
        "`f_doc_id` text NOT NULL,"
        "`f_user_id` char(40) NOT NULL,"
        "`f_create_date` bigint(20) NOT NULL DEFAULT '-1',"
        "`f_refresh_date` bigint(20) NOT NULL DEFAULT '-1',"
        "`f_expire_time` bigint(20) NOT NULL DEFAULT '-2',"
        "PRIMARY KEY (`f_primary_id`),"
        "KEY `t_finder_f_doc_id_index` (`f_doc_id`(120)) USING BTREE,"
        "KEY `t_lock_f_refresh_date_index` (`f_refresh_date`) USING BTREE,"
        "KEY `t_lock_f_expire_time_index` (`f_expire_time`) USING BTREE"
        ") ENGINE=InnoDB;";
    dbOperator->Execute (strSql);

    // 创建 t_conf
    strSql = "CREATE TABLE IF NOT EXISTS `t_conf` ("
        "`f_key` char(32) NOT NULL,"
        "`f_value` char(255) NOT NULL,"
        "PRIMARY KEY (`f_key`),"
        "KEY `f_key_index` (`f_key`) USING BTREE"
        ") ENGINE=InnoDB;";
    dbOperator->Execute (strSql);

    // 创建 t_acs_access_token
    strSql = "CREATE TABLE IF NOT EXISTS `t_acs_access_token` ("
        "`f_token_id` char(40) NOT NULL,"
        "`f_user_id` char(40) NOT NULL,"
        "`f_udid` char(40) NOT NULL,"
        "`f_create_time` datetime NOT NULL,"
        "`f_last_request_time` datetime NOT NULL,"
        "`f_expires` bigint(20) NOT NULL,"
        "`f_flag` int(11) NOT NULL DEFAULT '0',"
        "`f_login_ip` char(15) NOT NULL,"
        "`f_os_type` tinyint(4) NOT NULL DEFAULT '0',"
        "`f_account_type` tinyint(4) NOT NULL DEFAULT '0',"
        "`f_version` tinyint(3) unsigned DEFAULT '1',"
        "PRIMARY KEY (`f_token_id`),"
        "KEY `t_token_f_user_id_index` (`f_user_id`) USING BTREE,"
        "KEY `t_token_f_last_request_tiime` (`f_last_request_time`) USING BTREE,"
        "KEY `t_token_f_udid_index` (`f_udid`) USING BTREE,"
        "KEY `t_token_f_os_type` (`f_os_type`),"
        "KEY `t_token_f_account_type` (`f_account_type`) USING BTREE,"
        "KEY `t_token_f_login_ip` (`f_login_ip`)"
        ") ENGINE=InnoDB;";
    dbOperator->Execute (strSql);

    // 创建 t_acs_refresh_token
    strSql = "CREATE TABLE IF NOT EXISTS `t_acs_refresh_token` ("
        "`f_refresh_token_id` char(40) NOT NULL,"
        "`f_token_id` char(40) NOT NULL,"
        "`f_user_id` char(40) NOT NULL,"
        "`f_udid` char(40) NOT NULL,"
        "`f_create_time` bigint(20) NOT NULL,"
        "`f_access_token_expires` int(11) NOT NULL,"
        "`f_refresh_token_expires` int(11) NOT NULL,"
        "`f_login_ip` char(15) NOT NULL,"
        "`f_os_type` tinyint(4) NOT NULL DEFAULT '0',"
        "`f_account_type` tinyint(4) NOT NULL DEFAULT '0',"
        "PRIMARY KEY (`f_refresh_token_id`),"
        "KEY `f_create_time_index` (`f_create_time`)"
    ") ENGINE=InnoDB;";
    dbOperator->Execute (strSql);

    strSql = "CREATE TABLE IF NOT EXISTS `t_acs_doc` ("
        "`f_doc_id` char(40) NOT NULL,"
        "`f_doc_type` tinyint(4) NOT NULL,"
        "`f_type_name` char(128) NOT NULL,"
        "`f_status` int(11) DEFAULT '1',"
        "`f_create_time` bigint(20) NOT NULL DEFAULT '0',"
        "`f_delete_time` bigint(20) DEFAULT '0',"
        "`f_deleter_id` char(40) NOT NULL DEFAULT '',"
        "`f_obj_id` char(40) NOT NULL,"
        "`f_name` char(128) NOT NULL,"
        "`f_creater_id` char(40) NOT NULL,"
        "`f_site_id` char(150) NOT NULL DEFAULT '',"
        "`f_relate_depart_id` char(40) NOT NULL DEFAULT '',"
        "`f_display_order` mediumint DEFAULT '-1',"
        "PRIMARY KEY (`f_doc_id`),"
        "KEY `t_doc_f_doc_type_index` (`f_doc_type`) USING BTREE,"
        "KEY `t_doc_f_obj_id_index` (`f_obj_id`) USING BTREE,"
        "KEY `t_doc_f_name_index` (`f_name`) USING BTREE,"
        "KEY `t_doc_f_type_name_index` (`f_type_name`) USING BTREE,"
        "KEY `t_doc_f_relate_depart_id_index` (`f_relate_depart_id`) USING BTREE,"
        "KEY `t_display_order_index` (`f_display_order`) USING BTREE"
    ") ENGINE=InnoDB;";
    dbOperator->Execute (strSql);
}

void CreateShareMgntDB ()
{
    ncDBConnectionInfo connInfo;
    connInfo.ip = _T("127.0.0.1");
    connInfo.port = 3306;
    connInfo.user = _T("root");
    connInfo.password = _T("eisoo.com");

    nsresult ret;
    nsCOMPtr<ncIDBOperator> dbOperator = do_CreateInstance (NC_MYSQL_OPERATOR_CID, &ret);
    if (NS_FAILED (ret)) {
        String error;
        error.format (_T("Failed to create db operator: 0x%x"), ret);

        throw Exception (error);
    }
    dbOperator->Connect (connInfo);

    String strSql = _T("create DATABASE if not exists sharemgnt_db CHARACTER SET = 'utf8mb4' COLLATE = 'utf8mb4_bin'");
    dbOperator->Execute (strSql);

    dbOperator->Execute ("use sharemgnt_db");

    // 创建 t_user
    strSql = "CREATE TABLE IF NOT EXISTS `t_user` ("
        "`f_user_id` char(40) NOT NULL,"
        "`f_login_name` char(150) NOT NULL,"
        "`f_display_name` char(150) NOT NULL,"
        "`f_password` char(32) NOT NULL,"
        "`f_des_password` char(128) DEFAULT '',"
        "`f_ntlm_password` char(32) DEFAULT '',"
        "`f_mail_address` char(150) NOT NULL,"
        "`f_auth_type` tinyint(4) NOT NULL DEFAULT '0',"
        "`f_status` tinyint(4) NOT NULL DEFAULT '0',"
        "`f_freeze_status` tinyint(4) NOT NULL DEFAULT '0',"
        "`f_pwd_timestamp` datetime,"
        "`f_pwd_error_latest_timestamp` datetime,"
        "`f_pwd_error_cnt` tinyint(4) NOT NULL DEFAULT '0',"
        "`f_domain_object_guid`char(100) DEFAULT '',"
        "`f_domain_path` char(255) DEFAULT '',"
        "`f_ldap_server_type` tinyint(4) NOT NULL DEFAULT '0',"
        "`f_third_party_id` char(255),"
        "`f_third_party_depart_id` varchar(255),"
        "`f_priority` smallint(6) NOT NULL DEFAULT '999',"
        "`f_csf_level` tinyint(4) NOT NULL DEFAULT '5',"
        "`f_pwd_control` tinyint(1) NOT NULL DEFAULT '0',"
        "`f_site_id` char(150),"
        "`f_create_time` datetime DEFAULT now(),"
        "`f_last_request_time` datetime DEFAULT now(),"
        "`f_last_client_request_time` datetime NOT NULL DEFAULT now(),"
        "`f_auto_disable_status` tinyint(4) NOT NULL DEFAULT '0',"
        "`f_agreed_to_terms_of_use` tinyint(4) NOT NULL DEFAULT '0',"
        "`f_real_name_auth_status` tinyint(4) NOT NULL DEFAULT '0',"
        "`f_tel_number` char(40) DEFAULT NULL,"
        "`f_is_activate` tinyint(4) NOT NULL DEFAULT '0',"
        "`f_activate_status` tinyint(4) NOT NULL DEFAULT '0',"
        "`f_third_party_attr` varchar(255) NOT NULL DEFAULT '',"
        "`f_user_document_read_status` bigint(20) DEFAULT '0',"
        "PRIMARY KEY (`f_user_id`),"
        "KEY `f_mail_address_index` (`f_mail_address`),"
        "KEY `f_domain_object_guid` (`f_domain_object_guid`),"
        "UNIQUE KEY `f_login_name` (`f_login_name`),"
        "KEY `f_display_name_index` (`f_display_name`),"
        "KEY `f_third_party_id_index` (`f_third_party_id`),"
        "KEY `f_tel_number_index` (`f_tel_number`)"
        ") ENGINE=InnoDB;";
    dbOperator->Execute (strSql);
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

        CreateAnyShareDB ();

        CreateShareMgntDB ();

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
