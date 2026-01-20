#include <abprec.h>
#include <dataapi/dataapi.h>
#include <dboperator/public/ncIDBOperator.h>
#include <tedbcm/public/ncITEDBCManager.h>

#include "acsdbut.h"
#include "acsdb/public/ncIDBOwnerManager.h"
#include "acsdb/public/ncIDBPermManager.h"
#include "acsdb/public/ncIDBTokenManager.h"

//////////////////////////////////////////////////////////////////////////
//
// ncDBUTApi
//
// 获取DBOperator
ncIDBOperator*
ncDBUTApi::GetDBOperator (void)
{
    // 创建新的数据库连接
    nsresult ret;
    nsCOMPtr<ncIDBOperator> localDBOper = do_CreateInstance (NC_MYSQL_OPERATOR_CID, &ret);
    if (NS_FAILED (ret)) {
        String error;
        error.format (_T("Failed to create db operator: 0x%x"), ret);
        throw Exception (error);
    }

    ncDBConnectionInfo connInfo;
    connInfo.ip = _T("127.0.0.1");
    connInfo.port = 3306;
    connInfo.user = _T("root");
    connInfo.password = _T("eisoo.com");
    connInfo.db = _T("anyshare");

    localDBOper->Connect (connInfo);

    NS_ADDREF (localDBOper.get ());
    return localDBOper;
}

ncIDBOperator*
ncDBUTApi::GenerateAnyshareTB (const String& tbName)
{
    // 创建连接到anyshare数据库的实例
    nsresult ret;
    nsCOMPtr<ncIDBOperator> anyshareOperator = do_CreateInstance (NC_MYSQL_OPERATOR_CID, &ret);
    if (NS_FAILED (ret)){
        String error;
        error.format (_T("Failed to create anyshare operator: 0x%x"), ret);
        throw Exception (error);
    }
    ncDBConnectionInfo anyshareInfo;
    anyshareInfo.ip = _T("127.0.0.1");
    anyshareInfo.port = 3306;
    anyshareInfo.user = _T("root");
    anyshareInfo.password = _T("eisoo.com");
    anyshareInfo.db = _T("anyshare");
    anyshareOperator->Connect (anyshareInfo);

    String strSql;
    // 创建数据库 anyshare
    strSql = _T("create DATABASE if not exists anyshare CHARACTER SET = 'utf8mb4' COLLATE = 'utf8mb4_bin'");
    anyshareOperator->Execute (strSql);
    // 选择数据库 anyshare
    strSql = _T("use anyshare");
    anyshareOperator->Execute (strSql);
    // 创建表
    strSql = FindCreateTBSQL (tbName);
    anyshareOperator->Execute (strSql);

    NS_ADDREF (anyshareOperator.get ());
    return anyshareOperator.get ();
}

void
ncDBUTApi::CleanTable (ncIDBOperator* dbOperator, const String& tableName)
{
    String SQL;
    SQL.format (_T("TRUNCATE TABLE %s"), tableName.getCStr ());
    dbOperator->Execute (SQL);
}

String
ncDBUTApi::FindCreateTBSQL (const String& tb)
{
    if (tb == NC_TB_T_ACS_DOC){
        return ncDBDocUTApi::GetCreateTBSQL ();
    }
    else if (tb == NC_TB_T_ACS_OWNER){
        return ncDBOwnerUTApi::GetCreateTBSQL ();
    }
    else if (tb == NC_TB_T_ACS_CUSTOM_PERM){
        return ncDBPermUTApi::GetCreateTBSQL ();
    }
    else if (tb == NC_TB_T_ACS_ACCESS_TOKEN){
        return ncDBTokenUTApi::GetCreateTBSQL ();
    }
    else if (tb == NC_TB_T_ACS_FINDER){
        return ncDBFinderUTApi::GetCreateTBSQL ();
    }
    else{}
    return String ();
}

//////////////////////////////////////////////////////////////////////////
//
// ncDBDocUTApi
//
// 获取建表的SQL语句 (性能极低，仅在UT中使用)
const String
ncDBDocUTApi::GetCreateTBSQL (void)
{
    String sql_t_acs_doc =                                                        \
            _T("CREATE TABLE IF NOT EXISTS `t_acs_doc` (")                        \
            _T("`f_doc_id` char(40) NOT NULL,")    \
            _T("`f_doc_type` tinyint(4) NOT NULL,") \
            _T("`f_type_name` char(128) NOT NULL,") \
            _T("`f_status` int(11) DEFAULT '1',") \
            _T("`f_create_time` bigint(20) NOT NULL DEFAULT '0',") \
            _T("`f_delete_time` bigint(20) DEFAULT '0',") \
            _T("`f_obj_id` char(40) NOT NULL,") \
            _T("`f_name` char(128) NOT NULL,") \
            _T("`f_creater_id` char(40) NOT NULL,") \
            _T("`f_site_id` char(150) NOT NULL DEFAULT '',") \
            _T("PRIMARY KEY (`f_doc_id`),") \
            _T("KEY `t_doc_f_doc_id_index` (`f_doc_id`) USING BTREE,") \
            _T("KEY `t_doc_f_doc_type_index` (`f_doc_type`) USING BTREE,") \
            _T("KEY `t_doc_f_obj_id_index` (`f_obj_id`) USING BTREE,") \
            _T("KEY `t_doc_f_name_index` (`f_name`) USING BTREE") \
            _T(") ENGINE=InnoDB;");
    return sql_t_acs_doc;
}

//////////////////////////////////////////////////////////////////////////
//
// ncDBFinderUTApi
//
const String
ncDBFinderUTApi::GetCreateTBSQL (void)
{
    String sql_t_finder =                                            \
        _T("CREATE TABLE IF NOT EXISTS `t_finder` (") \
        _T("`f_primary_id` bigint(20) NOT NULL AUTO_INCREMENT,")    \
        _T("`f_doc_id` text NOT NULL,") \
        _T("`f_user_id` char(40) NOT NULL,") \
        _T("`f_accessors_id` text,") \
        _T("`f_is_allowed` tinyint(4) NOT NULL DEFAULT '1',") \
        _T("PRIMARY KEY (`f_primary_id`),") \
        _T("KEY `t_finder_f_doc_id_index` (`f_doc_id`(120)) USING BTREE") \
        _T(") ENGINE=InnoDB;");
    return sql_t_finder;
}

bool
ncDBFinderUTApi::EqualDBFinderInfo (const ncDbFinderInfo& info1, const ncDbFinderInfo& info2)
{
    // TODO
    return true;
}

void
ncDBFinderUTApi::InsertFinderInfo1 (ncIDBOperator* anyshareOperator, ncDbFinderInfo& info1)
{
    // TODO
}

void
ncDBFinderUTApi::InsertFinderInfo2 (ncIDBOperator* anyshareOperator, ncDbFinderInfo& info1, ncDbFinderInfo& info2)
{
    // TODO
}


//////////////////////////////////////////////////////////////////////////
//
// ncDBOwnerUTApi
//
const String
ncDBOwnerUTApi::GetCreateTBSQL (void)
{
    String sql_t_acs_owner = _T("CREATE TABLE IF NOT EXISTS `t_acs_owner` (")
        _T("`f_primary_id` bigint(20) NOT NULL AUTO_INCREMENT,")
        _T("`f_gns_path` text NOT NULL,")
        _T("`f_owner_id` char(40) NOT NULL,")
        _T("`f_owner_name` varchar(150) NOT NULL,")
        _T("`f_type` tinyint(4) NOT NULL,")
        _T("`f_modify_time` bigint(20) NOT NULL DEFAULT '0',")
        _T("`f_deletable` tinyint(1) NOT NULL,")
        _T("PRIMARY KEY (`f_primary_id`),")
        _T("KEY `t_owner_f_gns_path_index` (`f_gns_path`(120)) USING BTREE,")
        _T("KEY `t_owner_f_owner_id_index` (`f_owner_id`) USING BTREE")
        _T(") ENGINE=InnoDB;");
    return sql_t_acs_owner;
}

bool
ncDBOwnerUTApi::EqualDBOwnerInfo (const dbOwnerInfo& info1, const dbOwnerInfo& info2)
{
    // TODO
    return true;
}

void
ncDBOwnerUTApi::InsertOwnerInfo1 (ncIDBOperator* anyshareOperator, dbOwnerInfo& info1)
{
    // TODO
}

void
ncDBOwnerUTApi::InsertOwnerInfo2 (ncIDBOperator* anyshareOperator, dbOwnerInfo& info1, dbOwnerInfo& info2)
{
    // TODO
}


//////////////////////////////////////////////////////////////////////////
//
// ncDBPermUTApi
//
const String
ncDBPermUTApi::GetCreateTBSQL (void)
{
    String sql_t_acs_custom_perm = _T("CREATE TABLE IF NOT EXISTS `t_acs_custom_perm` (")
        _T("`f_primary_id` bigint(20) NOT NULL AUTO_INCREMENT,")
        _T("`f_doc_id` text NOT NULL,")
        _T("`f_accessor_id` char(40) NOT NULL,")
        _T("`f_accessor_type` tinyint(4) NOT NULL,")
        _T("`f_accessor_name` varchar(150) NOT NULL,")
        _T("`f_type` tinyint(4) NOT NULL DEFAULT '1',")
        _T("`f_perm_value` int(11) NOT NULL DEFAULT '1',")
        _T("`f_source` tinyint(4) NOT NULL DEFAULT '1',")
        _T("`f_end_time` bigint(20) NOT NULL ,")
        _T("`f_modify_time` bigint(20) NOT NULL DEFAULT '0',")
        _T("`f_create_time` bigint(20) NOT NULL,")
        _T("PRIMARY KEY (`f_primary_id`),")
        _T("KEY `t_perm_f_doc_id_index` (`f_doc_id`(120)) USING BTREE,")
        _T("KEY `t_perm_f_accessor_id_index` (`f_accessor_id`) USING BTREE,")
        _T("KEY `t_perm_f_end_time_index` (`f_end_time`) USING BTREE")
        _T(") ENGINE=InnoDB;");
    return sql_t_acs_custom_perm;
}

bool
ncDBPermUTApi::EqualDBPermInfo (const dbCustomPermInfo& info1, const dbCustomPermInfo& info2)
{
    // TODO
    return true;
}

void
ncDBPermUTApi::InsertPermInfo1 (ncIDBOperator* anyshareOperator, dbCustomPermInfo& info1)
{
    // TODO
}

void
ncDBPermUTApi::InsertPermInfo2 (ncIDBOperator* anyshareOperator, dbCustomPermInfo& info1, dbCustomPermInfo& info2)
{
    // TODO
}


//////////////////////////////////////////////////////////////////////////
//
// ncDBTokenUTApi
//
const String
ncDBTokenUTApi::GetCreateTBSQL (void)
{
    String sql_acs_access_token = "CREATE TABLE IF NOT EXISTS `t_acs_access_token` ("
        "`f_token_id` char(40) NOT NULL,"
        "`f_user_id` char(40) NOT NULL,"
        "`f_udid` char(40) NOT NULL,"
        "`f_create_time` datetime NOT NULL,"
        "`f_last_request_time` datetime NOT NULL,"
        "`f_expires` bigint(20) NOT NULL,"
        "`f_flag` int(11) NOT NULL DEFAULT '0',"
        "`f_login_ip` char(15) NOT NULL,"
        "`f_os_type` tinyint(4) NOT NULL DEFAULT '0',"
        "`f_version` tinyint(3) unsigned DEFAULT '1',"
        "PRIMARY KEY (`f_token_id`),"
        "KEY `t_token_f_user_id_index` (`f_user_id`) USING BTREE,"
        "KEY `t_token_f_last_request_tiime` (`f_last_request_time`) USING BTREE,"
        "KEY `t_token_f_udid_index` (`f_udid`) USING BTREE,"
        "KEY `t_token_f_os_type` (`f_os_type`),"
        "KEY `t_token_f_login_ip` (`f_login_ip`)"
        ") ENGINE=InnoDB;";
    return sql_acs_access_token;
}

bool
ncDBTokenUTApi::EqualDBTokenInfo (const dbTokenInfo& info1, const dbTokenInfo& info2)
{
    // TODO
    return true;
}

void
ncDBTokenUTApi::InsertTokenInfo1 (ncIDBOperator* anyshareOperator, dbTokenInfo& info1)
{
    // TODO
}

void
ncDBTokenUTApi::InsertTokenInfo2 (ncIDBOperator* anyshareOperator, dbTokenInfo& info1, dbTokenInfo& info2)
{
    // TODO
}
