#include <abprec.h>
#include <ncutil/ncBusinessDate.h>
#include <dataapi/ncJson.h>
#include <dataapi/ncGNSUtil.h>
#include <algorithm>
#include <arpa/inet.h>
#include <boost/date_time/posix_time/posix_time.hpp>
#include <openssl/rsa.h>
#include <openssl/pem.h>
#include <openssl/err.h>
#include <openssl/evp.h>
#include <openssl/des.h>

#include "gen-cpp/ncTShareMgnt.h"
#include "gen-cpp/ShareMgnt_constants.h"

#include "acssharemgnt.h"
#include "ncACSShareMgnt.h"
#include "common/util.h"

using namespace boost::posix_time;
using namespace boost::gregorian;

// 匿名用户ID 同 EFAST\EApp\EVFS\src\evfs\public\ncIEVFSLinkHandler.idl:54
#define NC_EVFS_NAME_IOC_ANONYMOUS_ID                 ("12345678-0000-0000-0000-000000000000")

/* Implementation file */
NS_IMPL_QUERY_INTERFACE1 (ncACSShareMgnt, ncIACSShareMgnt)

// protected
NS_IMETHODIMP_(nsrefcnt) ncACSShareMgnt::AddRef (void)
{
    return 1;
}

// protected
NS_IMETHODIMP_(nsrefcnt) ncACSShareMgnt::Release (void)
{
    return 1;
}

AB_DEFINE_THREADSAFE_SINGLETON_NO_POOL (ncACSShareMgnt)

ncACSShareMgnt::ncACSShareMgnt()
{
    // 不允许登录的管理员帐号
    _adminIds.insert(g_ShareMgnt_constants.NCT_USER_ADMIN);
    _adminIds.insert(g_ShareMgnt_constants.NCT_USER_AUDIT);
    _adminIds.insert(g_ShareMgnt_constants.NCT_USER_SYSTEM);
    _adminIds.insert(g_ShareMgnt_constants.NCT_USER_SECURIT);
    NC_ACS_SHAREMGNT_TRACE (_T("this: %p"), this);
}

ncACSShareMgnt::~ncACSShareMgnt()
{
    NC_ACS_SHAREMGNT_TRACE (_T("this: %p"), this);
}

/* [notxpcom] String GetNameByAccessorId ([const] in StringRef accessorId, in ncIOCAccesorType accessorType); */
NS_IMETHODIMP_(String) ncACSShareMgnt::GetNameByAccessorId(const String& accessorId, ncIOCAccesorType accessorType)
{
    NC_ACS_SHAREMGNT_TRACE (_T("this: %p"), this);

    nsCOMPtr<ncIDBOperator> dbOper = getter_AddRefs (GetDBOperator ());

    String escAccessorId = dbOper->EscapeEx(accessorId);

    String strSql;
    String sharementDBName = Util::getDBName("sharemgnt_db");
    String userManagementDBName = Util::getDBName("user_management");
    if (accessorType == IOC_USER) {
        strSql.format (_T("select f_display_name from %s.t_user where f_user_id = '%s'"),
                       sharementDBName.getCStr(), escAccessorId.getCStr ());
    }
    else if (accessorType == IOC_DEPARTMENT) {
        strSql.format (_T("select f_name from %s.t_department where f_department_id = '%s'"),
                       sharementDBName.getCStr(), escAccessorId.getCStr ());
    }
    else if (accessorType == IOC_USER_GROUP) {
        strSql.format (_T("select f_group_name from %s.t_person_group where f_group_id = '%s'"),
            sharementDBName.getCStr(), escAccessorId.getCStr ());
    }
    else if (accessorType == IOC_GROUP)
    {
        strSql.format(_T("select f_group_name from %s.t_group where f_group_id = '%s'"),
            userManagementDBName.getCStr(), escAccessorId.getCStr());
    }
    else {
        NC_ACS_SHAREMGNT_TRACE (_T("GetNameByAccessorId(%s) end, not found"), accessorId.getCStr ());
        return accessorId;
    }

    ncDBRecords results;
    dbOper->Select (strSql, results);

    if (results.size () == 0) {
        NC_ACS_SHAREMGNT_TRACE (_T("GetNameByAccessorId(%s) end, not found"), accessorId.getCStr ());
        return accessorId;
    }

    NC_ACS_SHAREMGNT_TRACE (_T("this: %p"), this);

    return results[0][0];
}

/* [notxpcom] String GetDirectBelongDepartmentIds ([const] in StringRef userId, in StringVecRef departIds); */
NS_IMETHODIMP_(void) ncACSShareMgnt::GetDirectBelongDepartmentIds(const String & userId, vector<String> & departIds)
{
    NC_ACS_SHAREMGNT_TRACE (_T("[BEGIN]this: %p"), this);

    nsCOMPtr<ncIDBOperator> dbOper = getter_AddRefs (GetDBOperator ());

    departIds.clear();
    String dbName = Util::getDBName("sharemgnt_db");
    String strSql;
    strSql.format (_T("select f_department_id from %s.t_user_department_relation where f_user_id = '%s'"),
        dbName.getCStr(), dbOper->EscapeEx(userId).getCStr ());

    ncDBRecords results;
    dbOper->Select (strSql, results);

    // 如果为-1代表为未分配用户，则过滤掉
    for (int64 i = 0; i < results.size(); i++) {
        if (results[i][0] != "-1") {
            departIds.push_back (results[i][0]);
        }
    }

    NC_ACS_SHAREMGNT_TRACE (_T("[END]this: %p"), this);
}

/* [notxpcom] String GetDirectBelongOrganIds ([const] in StringRef userId, in StringVecRef organIds); */
NS_IMETHODIMP_(void) ncACSShareMgnt::GetDirectBelongOrganIds(const String & userId, vector<String> & organIds)
{
    NC_ACS_SHAREMGNT_TRACE (_T("this: %p"), this);

    nsCOMPtr<ncIDBOperator> dbOper = getter_AddRefs (GetDBOperator ());

    organIds.clear();

    String dbName = Util::getDBName("sharemgnt_db");
    String strSql;
    strSql.format (_T("select f_ou_id from %s.t_ou_user where f_user_id = '%s'"),
        dbName.getCStr(), dbOper->EscapeEx(userId).getCStr ());

    ncDBRecords results;
    dbOper->Select (strSql, results);

    for (int64 i = 0; i < results.size(); i++) {
        organIds.push_back (results[i][0]);
    }

    NC_ACS_SHAREMGNT_TRACE (_T("this: %p"), this);
}


/* [notxpcom] String GetAllBelongDepartmentIds ([const] in StringRef userId, in StringVecRef departIds); */
NS_IMETHODIMP_(void) ncACSShareMgnt::GetAllBelongDepartmentIds(const String & userId, vector<String> & departIds)
{
    NC_ACS_SHAREMGNT_TRACE (_T("[BEGIN]this: %p"), this);

    nsCOMPtr<ncIDBOperator> dbOper = getter_AddRefs (GetDBOperator ());

    departIds.clear();

    vector<String> tmpIds;
    GetDirectBelongDepartmentIds (userId, tmpIds);

    String dbName = Util::getDBName("sharemgnt_db");
    String strSql;
    while (tmpIds.size () > 0) {
        for (size_t i = 0; i < tmpIds.size (); ++i) {
            departIds.push_back (tmpIds[i]);
        }

        String groupStr = GenerateGroupStr (tmpIds);

        strSql.format (_T("select f_parent_department_id from %s.t_department_relation where f_department_id in (%s);"),
                        dbName.getCStr(), groupStr.getCStr ());

        ncDBRecords results;
        dbOper->Select (strSql, results);

        tmpIds.clear ();
        for (size_t i = 0; i < results.size(); i++) {
            tmpIds.push_back(results[i][0]);
        }
    }

    // 去除掉重复id
    RemoveDuplicateStrs (departIds);

    NC_ACS_SHAREMGNT_TRACE (_T("[END]this: %p"), this);
}

/* [notxpcom] void GetAllOrgInfo ([const] in StringRef userId, in OrganizationVecRef organs); */
NS_IMETHODIMP_(void) ncACSShareMgnt::GetAllOrgInfo(const String & userId, vector<ncOrganizationInfo> & organs)
{
    NC_ACS_SHAREMGNT_TRACE (_T("this: %p"), this);
    organs.clear();

    nsCOMPtr<ncIDBOperator> dbOper = getter_AddRefs (GetDBOperator ());

    String dbName = Util::getDBName("sharemgnt_db");
    String strSql;
    strSql.format(_T("select f_department_id,f_name from %s.t_department where f_is_enterprise = 1 ")
                  _T("order by f_priority, upper(f_name)"), dbName.getCStr());

    ncDBRecords results;
    dbOper->Select (strSql, results);

    for (size_t i = 0; i < results.size(); i++) {
        ncOrganizationInfo info;
        info.id = std::move(results[i][0]);
        info.name = std::move(results[i][1]);

        organs.push_back (info);
    }

    NC_ACS_SHAREMGNT_TRACE (_T("this: %p"), this);
}

/* [notxpcom] bool IsAdminId ([const] in StringRef userId); */
NS_IMETHODIMP_(bool) ncACSShareMgnt::IsAdminId(const String & userId)
{
    NC_ACS_SHAREMGNT_TRACE (_T("this: %p"), this);

    if (_adminIds.count(userId.getCStr()) != 0) {
        return true;
    }

    NC_ACS_SHAREMGNT_TRACE (_T("this: %p"), this);
    return false;
}

/* [notxpcom] bool IsUndistirbutedUser ([const] in StringRef userId); */
NS_IMETHODIMP_(bool) ncACSShareMgnt::IsUndistirbutedUser(const String & userId)
{
    NC_ACS_SHAREMGNT_TRACE (_T("this: %p"), this);

    nsCOMPtr<ncIDBOperator> dbOper = getter_AddRefs (GetDBOperator ());

    String dbName = Util::getDBName("sharemgnt_db");
    String strSql;
    strSql.format (_T("select f_user_id from %s.t_ou_user where f_user_id = '%s'"),
                       dbName.getCStr(), dbOper->EscapeEx(userId).getCStr ());

    ncDBRecords results;
    dbOper->Select (strSql, results);

    for(int i = 0; i < results.size(); ++i){
        if(results[i][0] == userId)
            return false;
    }

    NC_ACS_SHAREMGNT_TRACE (_T("this: %p"), this);

    return true;
}

/* [notxpcom] bool IsCustomDocManager ([const] in StringRef userId); */
NS_IMETHODIMP_(bool) ncACSShareMgnt::IsCustomDocManager(const String & userId)
{
    NC_ACS_SHAREMGNT_TRACE (_T("this: %p"), this);

    nsCOMPtr<ncIDBOperator> dbOper = getter_AddRefs (GetDBOperator ());

    // 只有超级管理员、系统管理员或组织管理员可以管理文档库
    std::vector<String> roleIds;
    GetUserRoleIds(userId, roleIds);
    if (roleIds.size() > 0) {
        if(std::find(roleIds.begin(), roleIds.end(), toCFLString(g_ShareMgnt_constants.NCT_SYSTEM_ROLE_SUPPER)) != roleIds.end()) {
            return true;
        }
        if(std::find(roleIds.begin(), roleIds.end(), toCFLString(g_ShareMgnt_constants.NCT_SYSTEM_ROLE_ADMIN)) != roleIds.end()) {
            return true;
        }
        if(std::find(roleIds.begin(), roleIds.end(), toCFLString(g_ShareMgnt_constants.NCT_SYSTEM_ROLE_ORG_MANAGER)) != roleIds.end()) {
            return true;
        }
    }
    NC_ACS_SHAREMGNT_TRACE (_T("this: %p"), this);

    return false;
}

/* [notxpcom] ncACSUserType GetUserType ([const] in StringRef userId); */
NS_IMETHODIMP_(ncACSUserType) ncACSShareMgnt::GetUserType(const String & userId)
{
    NC_ACS_SHAREMGNT_TRACE (_T("this: %p"), this);

    nsCOMPtr<ncIDBOperator> dbOper = getter_AddRefs (GetDBOperator ());

    String dbName = Util::getDBName("sharemgnt_db");
    String strSql;
    strSql.format (_T("select f_auth_type from %s.t_user where f_user_id = '%s';"),
                       dbName.getCStr(), dbOper->EscapeEx(userId).getCStr ());

    ncDBRecords results;
    dbOper->Select (strSql, results);

    int userType = USER_TYPE_NONE;
    if (results.size () == 1) {
        userType = Int::getValue (results[0][0]);
    }

    if ( (userType < USER_TYPE_LOCAL) || (userType > USER_TYPE_THIRD)) {
        userType = USER_TYPE_NONE;
    }

    NC_ACS_SHAREMGNT_TRACE (_T("this: %p"), this);

    return (ncACSUserType)userType;
}

/* [notxpcom] void GetUserOSSId ([const] in StringRef userId, in StringRef OSSId); */
NS_IMETHODIMP_(void) ncACSShareMgnt::GetUserOSSId(const String & userId, String & ossId)
{
    NC_ACS_SHAREMGNT_TRACE (_T("this: %p"), this);

    nsCOMPtr<ncIDBOperator> dbOper = getter_AddRefs (GetDBOperator ());

    String dbName = Util::getDBName("sharemgnt_db");
    String strSql;
    strSql.format (_T("select f_oss_id from %s.t_user where f_user_id = '%s';"),
                       dbName.getCStr(), dbOper->EscapeEx(userId).getCStr ());

    ncDBRecords results;
    dbOper->Select (strSql, results);

    int userType = USER_TYPE_NONE;
    if (results.size () == 1) {
        ossId = std::move(results[0][0]);
    }

    NC_ACS_SHAREMGNT_TRACE (_T("this: %p"), this);
}

/* [notxpcom] void GetDepartmentOSSId ([const] in StringRef departmentId, in StringRef OSSId); */
NS_IMETHODIMP_(void) ncACSShareMgnt::GetDepartmentOSSId(const String & departmentId, String & ossId)
{
    NC_ACS_SHAREMGNT_TRACE (_T("this: %p"), this);

    nsCOMPtr<ncIDBOperator> dbOper = getter_AddRefs (GetDBOperator ());

    String dbName = Util::getDBName("sharemgnt_db");
    String strSql;
    strSql.format (_T("select f_oss_id from %s.t_department where f_department_id = '%s';"),
                       dbName.getCStr(), dbOper->EscapeEx(departmentId).getCStr ());

    ncDBRecords results;
    dbOper->Select (strSql, results);

    if (results.size () == 1) {
        ossId = std::move(results[0][0]);
    }

    NC_ACS_SHAREMGNT_TRACE (_T("this: %p"), this);
}

/* [notxpcom] void GetAllGroups ([const] in StringRef createrId, in ncGroupInfoVecRef groups); */
NS_IMETHODIMP_(void) ncACSShareMgnt::GetAllGroups(const String& createrId, vector<ncGroupInfo> & groups)
{
    NC_ACS_SHAREMGNT_TRACE (_T("this: %p"), this);

    nsCOMPtr<ncIDBOperator> dbOper = getter_AddRefs (GetDBOperator ());

    groups.clear();

    String dbName = Util::getDBName("sharemgnt_db");
    String strSql;
    strSql.format(_T("select f_group_id,f_user_id,f_group_name,f_person_count ")
                  _T("from %s.t_person_group where f_user_id = '%s' ")
                  _T("order by upper(f_group_name) "),
                  dbName.getCStr(), dbOper->EscapeEx(createrId).getCStr(), LOAD_STRING(_T("IDS_TEMP_CONTACTOR_GROUP_NAME")).getCStr());

    ncDBRecords results;
    dbOper->Select (strSql, results);

    for (int64 i = 0; i < results.size(); i++) {
        ncGroupInfo info;
        info.id = std::move(results[i][0]);
        info.createrId = std::move(results[i][1]);
        info.groupName = std::move(results[i][2]);
        info.count = Int::getValue (results[i][3]);

        groups.push_back (info);
    }

    NC_ACS_SHAREMGNT_TRACE (_T("this: %p"), this);
}

/* [notxpcom] void GetContactors ([const] in StringRef groupId, in ncACSUserInfoVecRef userInfos); */
NS_IMETHODIMP_(void) ncACSShareMgnt::GetContactors(const String & groupId, vector<ncACSUserInfo> & userInfos)
{
    NC_ACS_SHAREMGNT_TRACE (_T("this: %p"), this);

    nsCOMPtr<ncIDBOperator> dbOper = getter_AddRefs (GetDBOperator ());

    userInfos.clear();

    String dbName = Util::getDBName("sharemgnt_db");
    String strSql;
    strSql.format(_T("select f_user_id,f_login_name,f_display_name,f_mail_address,f_status,f_csf_level,")
                  _T("f_auto_disable_status from %s.t_user where f_user_id in ")
                  _T("(select f_user_id from t_contact_person where f_group_id = '%s') ")
                  _T("order by f_priority, upper(f_display_name)"),
                  dbName.getCStr(), dbOper->EscapeEx(groupId).getCStr());

    ncDBRecords results;
    dbOper->Select (strSql, results);

    ncACSUserInfo userInfo;
    for (size_t i = 0; i < results.size (); ++i) {
        userInfo.id = std::move(results[i][0]);
        userInfo.account = std::move(results[i][1]);
        userInfo.visionName = std::move(results[i][2]);
        userInfo.email = std::move(results[i][3]);
        userInfo.enableStatus = (results[i][4] == "0") && (results[i][6] == "0");
        userInfo.csfLevel = Int::getValue(results[i][5]);

        userInfos.push_back (userInfo);
    }

    NC_ACS_SHAREMGNT_TRACE (_T("this: %p"), this);
}

/* [notxpcom] void GetAllBelongGroups ([const] in StringRef userId, in StringVecRef groupIds); */
NS_IMETHODIMP_(void) ncACSShareMgnt::GetAllBelongGroups(const String& userId,  vector<String> & groupIds)
{
    NC_ACS_SHAREMGNT_TRACE (_T("[BEGIN]this: %p"), this);

    nsCOMPtr<ncIDBOperator> dbOper = getter_AddRefs (GetDBOperator ());

    groupIds.clear();

    String dbName = Util::getDBName("sharemgnt_db");
    String strSql;
    strSql.format (_T("select f_group_id from %s.t_contact_person where f_user_id = '%s'"),
                      dbName.getCStr(), dbOper->EscapeEx(userId).getCStr ());

    ncDBRecords results;
    dbOper->Select (strSql, results);

    for (int64 i = 0; i < results.size(); i++) {
        groupIds.push_back (results[i][0]);
    }

    NC_ACS_SHAREMGNT_TRACE (_T("[END]this: %p"), this);
}

/* [notxpcom] void GetUserDisplayName ([const] in StringRef userId, in StringRef name);*/
NS_IMETHODIMP_(void) ncACSShareMgnt::GetUserDisplayName(const String& userId, String& name)
{
    NC_ACS_SHAREMGNT_TRACE (_T("this: %p"), this);

    name = userId;

    // 匿名用户处理
    if (userId.compare (NC_EVFS_NAME_IOC_ANONYMOUS_ID) == 0) {
        name = LOAD_STRING (_T("IDS_ACS_SHARE_MGNT_ANONYMOUS_NAME"));
        return ;
    }

    // 内外网数据交换用户处理
    if (userId.compare (NC_EVFS_NAME_IOC_DATAEXCHANGE_ID) == 0) {
        name = LOAD_STRING (_T("IDS_ACS_SHARE_MGNT_DATAEXCHANGE_NAME"));
        return ;
    }

    nsCOMPtr<ncIDBOperator> dbOper = getter_AddRefs (GetDBOperator ());

    String dbName = Util::getDBName("sharemgnt_db");
    String strSql;
    strSql.format (_T("select f_display_name from %s.t_user where f_user_id = '%s'"), dbName.getCStr(), dbOper->EscapeEx(userId).getCStr ());

    ncDBRecords results;
    dbOper->Select (strSql, results);

    if (results.size () == 1) {
        name = std::move(results[0][0]);
    }

    NC_ACS_SHAREMGNT_TRACE (_T("this: %p"), this);
}

/* [notxpcom] void GetUserName ([const] in StringRef userId, in StringRef displayName, in StringRef account);*/
NS_IMETHODIMP_(void) ncACSShareMgnt::GetUserName(const String& userId, String& displayName, String& account)
{
    NC_ACS_SHAREMGNT_TRACE (_T("this: %p"), this);

    displayName = userId;
    account = userId;

    // 匿名用户处理
    if (userId.compare (NC_EVFS_NAME_IOC_ANONYMOUS_ID) == 0) {
        displayName = LOAD_STRING (_T("IDS_ACS_SHARE_MGNT_ANONYMOUS_NAME"));
        account = LOAD_STRING (_T("IDS_ACS_SHARE_MGNT_ANONYMOUS_NAME"));
        return ;
    }

    // 内外网数据交换用户处理
    if (userId.compare (NC_EVFS_NAME_IOC_DATAEXCHANGE_ID) == 0) {
        displayName = LOAD_STRING (_T("IDS_ACS_SHARE_MGNT_DATAEXCHANGE_NAME"));
        account = LOAD_STRING (_T("IDS_ACS_SHARE_MGNT_DATAEXCHANGE_NAME"));
        return ;
    }

    // "所有用户"显示名
    if (userId.compare (toCFLString(g_ShareMgnt_constants.NCT_ALL_USER_GROUP)) == 0) {
        displayName = LOAD_STRING (_T("IDS_ALL_USER_GROUP_NAME"));
        account = LOAD_STRING (_T("IDS_ALL_USER_GROUP_NAME"));
        return ;
    }

    nsCOMPtr<ncIDBOperator> dbOper = getter_AddRefs (GetDBOperator ());

    String dbName = Util::getDBName("sharemgnt_db");
    String strSql;
    strSql.format (_T("select f_display_name, f_login_name from %s.t_user where f_user_id = '%s'"), dbName.getCStr(), dbOper->EscapeEx(userId).getCStr ());

    ncDBRecords results;
    dbOper->Select (strSql, results);

    if (results.size () == 1) {
        displayName = std::move(results[0][0]);
        account = std::move(results[0][1]);
    }

    NC_ACS_SHAREMGNT_TRACE (_T("this: %p"), this);
}

/* [notxpcom] bool CheckDisplayNameIsExist ([const] in StringRef displayName, in StringRef name);*/
NS_IMETHODIMP_(bool) ncACSShareMgnt::CheckDisplayNameIsExist(const String& displayName)
{
    NC_ACS_SHAREMGNT_TRACE (_T("this: %p"), this);

    nsCOMPtr<ncIDBOperator> dbOper = getter_AddRefs (GetDBOperator ());

    String dbName = Util::getDBName("sharemgnt_db");
    String strSql;
    strSql.format (_T("select f_user_id from %s.t_user where f_display_name = '%s'"), dbName.getCStr(), dbOper->EscapeEx(displayName).getCStr ());

    ncDBRecords results;
    dbOper->Select (strSql, results);

    bool bExist = false;
    if (results.size () == 1) {
        bExist = true;
    }

    NC_ACS_SHAREMGNT_TRACE (_T("this: %p"), this);

    return bExist;
}

/* [notxpcom] bool CheckDisplayNameIsExist ([const] in StringRef displayName, in StringRef name);*/
NS_IMETHODIMP_(void) ncACSShareMgnt::GetUserIdByDisplayName(const String& displayName, vector<String>& userId)
{
    NC_ACS_SHAREMGNT_TRACE (_T("this: %p"), this);

    nsCOMPtr<ncIDBOperator> dbOper = getter_AddRefs (GetDBOperator ());

    String dbName = Util::getDBName("sharemgnt_db");
    userId.clear();
    String strSql;
    strSql.format (_T("select f_user_id from %s.t_user where f_display_name = '%s'"), dbName.getCStr(), dbOper->EscapeEx(displayName).getCStr ());

    ncDBRecords results;
    dbOper->Select (strSql, results);

    for(size_t i = 0; i < results.size(); ++i) {
        userId.push_back(results[i][0]);
    }
    NC_ACS_SHAREMGNT_TRACE (_T("this: %p"), this);
}

/* [notxpcom] bool GetUserInfoById ([const] in StringRef userId, in ncACSUserInfoRef userInfo); */
NS_IMETHODIMP_(bool) ncACSShareMgnt::GetUserInfoById(const String & userId, ncACSUserInfo & userInfo)
{
    NC_ACS_SHAREMGNT_TRACE (_T("this: %p"), this);

    nsCOMPtr<ncIDBOperator> dbOper = getter_AddRefs (GetDBOperator ());

    String dbName = Util::getDBName("sharemgnt_db");
    String strSql;
    strSql.format (_T("select f_login_name,f_display_name,f_mail_address,f_status,f_agreed_to_terms_of_use,f_tel_number,f_csf_level,")
                   _T("f_user_document_read_status,f_auto_disable_status,f_priority from %s.t_user where f_user_id = '%s'"),
                       dbName.getCStr(), dbOper->EscapeEx(userId).getCStr ());

    ncDBRecords results;
    dbOper->Select (strSql, results);

    bool ret = false;
    if (results.size () == 1) {
        userInfo.id = userId;
        userInfo.account = std::move(results[0][0]);
        userInfo.visionName = std::move(results[0][1]);
        userInfo.email = std::move(results[0][2]);
        userInfo.csfLevel = Int::getValue (results[0][6]);
        userInfo.documentReadStatus = Int64::getValue (results[0][7]);
        userInfo.enableStatus = (results[0][3] == "0") && (results[0][8] == "0");
        userInfo.isAgreedToTermsOfUse = results[0][4] == "1";
        ret = true;

        userInfo.telNumber = std::move(results[0][5]);
        userInfo.priority = Int64::getValue (results[0][9]);
    }

    NC_ACS_SHAREMGNT_TRACE (_T("this: %p"), this);

    return ret;
}

 /* [notxpcom] bool GetAppInfoById ([const] in StringRef appId, in ncACSAppInfoRef appInfo); */
NS_IMETHODIMP_(bool) ncACSShareMgnt::GetAppInfoById(const String & appId, ncACSAppInfo & appInfo)
{
    NC_ACS_SHAREMGNT_TRACE (_T("this: %p"), this);

    nsCOMPtr<ncIDBOperator> dbOper = getter_AddRefs (GetDBOperator ());

    String dbName = Util::getDBName("user_management");
    String strSql;
    strSql.format (_T("select f_name from %s.t_app where f_id = '%s'"),
                       dbName.getCStr(), dbOper->EscapeEx(appId).getCStr ());

    ncDBRecords results;
    dbOper->Select (strSql, results);

    bool ret = false;
    if (results.size () == 1) {
        appInfo.id = appId;
        appInfo.name = std::move(results[0][0]);
        ret = true;
    }

    NC_ACS_SHAREMGNT_TRACE (_T("this: %p"), this);

    return ret;
}

string ncACSShareMgnt::DESEncrypt(const string& plainText)
{
    string cipherText;
    int cipherTextLength = 0;
    if (plainText.length() % 8 == 0) {
        cipherTextLength = plainText.length();
    }
    else {
        cipherTextLength = plainText.length() + (8 - plainText.length() % 8);
    }
    cipherText.assign(cipherTextLength, '*');

    DES_cblock key = {'E', 'a', '8', 'e', 'k', '&', 'a', 'h'};
    DES_cblock ivec = {'E', 'a', '8', 'e', 'k', '&', 'a', 'h'};
    DES_key_schedule keysched;

    DES_set_odd_parity(&key);
    if (DES_set_key_checked((C_Block *)key, &keysched)) {
        throw Exception("Unable to set key schedule");
    }

    DES_ncbc_encrypt((unsigned char*)plainText.c_str(), (unsigned char*)cipherText.c_str(), plainText.length(), &keysched, &ivec, DES_ENCRYPT);

    return cipherText;
}

string ncACSShareMgnt::Base64Encode(const string& input)
{
    BIO * bmem = NULL;
    BIO * b64 = NULL;
    BUF_MEM * bptr = NULL;

    b64 = BIO_new(BIO_f_base64());
    BIO_set_flags(b64, BIO_FLAGS_BASE64_NO_NL);
    bmem = BIO_new(BIO_s_mem());
    b64 = BIO_push(b64, bmem);
    BIO_write(b64, (char*)input.c_str(), input.length());
    BIO_flush(b64);
    BIO_get_mem_ptr(b64, &bptr);

    string buffer;
    buffer.assign(bptr->data, bptr->length);

    BIO_free_all(b64);

    return buffer;
}

/* [notxpcom] bool GetUserInfoByAccount ([const] in StringRef account, in ncACSUserInfoRef userInfo, in intRef accountType); */
NS_IMETHODIMP_(bool) ncACSShareMgnt::GetUserInfoByAccount(const String & account, ncACSUserInfo & userInfo, int& accountType)
{
    NC_ACS_SHAREMGNT_TRACE (_T("this: %p"), this);

    accountType = 0;

    nsCOMPtr<ncIDBOperator> dbOper = getter_AddRefs (GetDBOperator ());

    String dbName = Util::getDBName("sharemgnt_db");
    String strSql;
    strSql.format (_T("select f_user_id,f_login_name,f_display_name,f_mail_address,f_status,f_auto_disable_status, f_tel_number from %s.t_user where f_login_name = '%s'"),
                dbName.getCStr(), dbOper->EscapeEx(account).getCStr ());
    ncDBRecords results;
    dbOper->Select (strSql, results);

    if (results.size () == 0){
        String strSql_idcard;
        ncDBRecords results_idcard;
        strSql_idcard.format (_T("select f_value from %s.t_sharemgnt_config where f_key = 'id_card_login_status';"), dbName.getCStr());
        dbOper->Select (strSql_idcard, results_idcard);

        if (results_idcard.size() == 1 && Int::getValue (results_idcard[0][0]) == 1){
            string des_account = DESEncrypt(toSTLString(account) );
            string base64_account = Base64Encode(des_account);
            strSql.format (_T("select f_user_id,f_login_name,f_display_name,f_mail_address,f_status,f_auto_disable_status, f_tel_number from %s.t_user where f_idcard_number = '%s' limit 1"),
                            dbName.getCStr(), dbOper->EscapeEx(toCFLString(base64_account)).getCStr ());
            dbOper->Select (strSql, results);
            if (results.size () == 1) {
                accountType = 1;
            }
        }
    }

    bool ret = false;
    if (results.size () == 1) {
        userInfo.id = std::move(results[0][0]);
        userInfo.account = std::move(results[0][1]);
        userInfo.visionName = std::move(results[0][2]);
        userInfo.email = std::move(results[0][3]);
        userInfo.enableStatus = (results[0][4] == "0") && (results[0][5] == "0");
        userInfo.telNumber = std::move(results[0][6]);

        ret = true;
    }

    NC_ACS_SHAREMGNT_TRACE (_T("this: %p"), this);

    return ret;
}

/* [notxpcom] int GetAccountType ([const] in StringRef account); */
NS_IMETHODIMP_(int) ncACSShareMgnt::GetAccountType(const String& account)
{
    NC_ACS_SHAREMGNT_TRACE (_T("this: %p"), this);

    int accountType = 0;

    nsCOMPtr<ncIDBOperator> dbOper = getter_AddRefs (GetDBOperator ());

    String dbName = Util::getDBName("sharemgnt_db");
    String strSql_idcard;
    strSql_idcard.format (_T("select f_value from %s.t_sharemgnt_config where f_key = 'id_card_login_status';"), dbName.getCStr());

    ncDBRecords results_idcard;
    dbOper->Select (strSql_idcard, results_idcard);

    if (results_idcard.size() == 1 && Int::getValue (results_idcard[0][0]) == 1){
        string des_account = DESEncrypt(toSTLString(account) );
        string base64_account = Base64Encode(des_account);
        String strSql;
        strSql.format (_T("select f_user_id from %s.t_user where f_idcard_number = '%s' limit 1"),
                        dbName.getCStr(), dbOper->EscapeEx(toCFLString(base64_account)).getCStr ());
        ncDBRecords results;
        dbOper->Select (strSql, results);
        if (results.size () == 1) {
            accountType = 1;
        }
    }

    NC_ACS_SHAREMGNT_TRACE (_T("this: %p"), this);

    return accountType;
}

/* [notxpcom] bool IsUserEnabled ([const] in StringRef userId); */
NS_IMETHODIMP_(bool) ncACSShareMgnt::IsUserEnabled(const String & userId)
{
    NC_ACS_SHAREMGNT_TRACE (_T("this: %p"), this);

    nsCOMPtr<ncIDBOperator> dbOper = getter_AddRefs (GetDBOperator ());

    String dbName = Util::getDBName("sharemgnt_db");
    String strSql;
    strSql.format (_T("select f_status from %s.t_user where f_user_id = '%s'"),
                       dbName.getCStr(), dbOper->EscapeEx(userId).getCStr ());

    ncDBRecords results;
    dbOper->Select (strSql, results);

    bool ret = false;
    if (results.size () == 1) {
        ret = Int::getValue(results[0][0]) == 0;
    }

    NC_ACS_SHAREMGNT_TRACE (_T("this: %p"), this);

    return ret;
}

/* [notxpcom] bool GetUserInfoByThirdId ([const] in StringRef thirdId, in ncACSUserInfoRef userInfo); */
NS_IMETHODIMP_(bool) ncACSShareMgnt::GetUserInfoByThirdId(const String & thirdId, ncACSUserInfo & userInfo)
{
    NC_ACS_SHAREMGNT_TRACE (_T("this: %p"), this);

    nsCOMPtr<ncIDBOperator> dbOper = getter_AddRefs (GetDBOperator ());

    String dbName = Util::getDBName("sharemgnt_db");
    String strSql;
    strSql.format (_T("select f_user_id,f_login_name,f_display_name,f_mail_address from %s.t_user where f_third_party_id = '%s'"),
                       dbName.getCStr(), dbOper->EscapeEx(thirdId).getCStr ());

    ncDBRecords results;
    dbOper->Select (strSql, results);

    bool ret = false;
    if (results.size () == 1) {
        userInfo.id = std::move(results[0][0]);
        userInfo.account = std::move(results[0][1]);
        userInfo.visionName = std::move(results[0][2]);
        userInfo.email = std::move(results[0][3]);

        ret = true;
    }

    NC_ACS_SHAREMGNT_TRACE (_T("this: %p"), this);

    return ret;
}

/* [notxpcom] void GetDepartInfoById ([const] in StringRef departId, in ncACSDepartInfoRef departInfo); */
NS_IMETHODIMP_(bool) ncACSShareMgnt::GetDepartInfoById(const String & departId, ncACSDepartInfo & departInfo)
{
    NC_ACS_SHAREMGNT_TRACE (_T("this: %p"), this);

    nsCOMPtr<ncIDBOperator> dbOper = getter_AddRefs (GetDBOperator ());

    String dbName = Util::getDBName("sharemgnt_db");
    String strSql;
    strSql.format (_T("select f_department_id,f_name from %s.t_department where f_department_id = '%s';"),
                       dbName.getCStr(), dbOper->EscapeEx(departId).getCStr ());

    ncDBRecords results;
    dbOper->Select (strSql, results);

    bool ret = false;
    if (results.size () == 1) {
        departInfo.id = std::move(results[0][0]);
        departInfo.name = std::move(results[0][1]);

        ret = true;
    }

    NC_ACS_SHAREMGNT_TRACE (_T("this: %p"), this);

    return ret;
}

/* [notxpcom] void GetOrgInfoByDeptId ([const] in StringRef departId, in OrganizationVecRef organs);*/
NS_IMETHODIMP_(void) ncACSShareMgnt::GetOrgInfoByDeptId(const String & departId, vector<ncOrganizationInfo>& organs)
{
    NC_ACS_SHAREMGNT_TRACE (_T("this: %p"), this);
    organs.clear();

    nsCOMPtr<ncIDBOperator> dbOper = getter_AddRefs (GetDBOperator ());

    String dbName = Util::getDBName("sharemgnt_db");
    String strSql;
    strSql.format(_T("select f_department_id,f_name from %s.t_department ")
                  _T("where f_department_id in ")
                  _T("(select f_ou_id from %s.t_ou_department where f_department_id = '%s') ")
                  _T("order by f_priority, upper(f_name)"),
                  dbName.getCStr(), dbName.getCStr(), dbOper->EscapeEx(departId).getCStr());

    ncDBRecords results;
    dbOper->Select (strSql, results);

    for (size_t i = 0; i < results.size(); i++) {
        ncOrganizationInfo info;
        info.id = std::move(results[i][0]);
        info.name = std::move(results[i][1]);

        organs.push_back (info);
    }

    NC_ACS_SHAREMGNT_TRACE (_T("this: %p"), this);
}

/* [notxpcom] void GetOrgIdByUserId ([const] in StringRef userId, in StringVecRef orgIds);*/
NS_IMETHODIMP_(void) ncACSShareMgnt::GetOrgIdByUserId(const String & userId, vector<String>& orgIds)
{
    NC_ACS_SHAREMGNT_TRACE (_T("this: %p"), this);
    orgIds.clear();

    nsCOMPtr<ncIDBOperator> dbOper = getter_AddRefs (GetDBOperator ());

    String dbName = Util::getDBName("sharemgnt_db");
    String strSql;
    strSql.format ( _T("(select f_ou_id from %s.t_ou_user where f_user_id = '%s') "), dbName.getCStr(), dbOper->EscapeEx(userId).getCStr ());

    ncDBRecords results;
    dbOper->Select (strSql, results);

    for (size_t i = 0; i < results.size(); i++) {
        orgIds.push_back (results[i][0]);
    }

    NC_ACS_SHAREMGNT_TRACE (_T("this: %p"), this);
}


/* [notxpcom] void GetUserInfoByIdBatch ([const] in StringVecRef userIds, in ncACSUserInfoMapRef userInfoMap); */
NS_IMETHODIMP_(void) ncACSShareMgnt::GetUserInfoByIdBatch(const vector<String> & userIds, map<String, ncACSUserInfo> & userInfoMap)
{
    NC_ACS_SHAREMGNT_TRACE (_T("[BEGIN]this: %p"), this);

    userInfoMap.clear ();

    if (userIds.size () == 0) {
        return;
    }

    nsCOMPtr<ncIDBOperator> dbOper = getter_AddRefs (GetDBOperator ());

    String groupStr = GenerateGroupStr (userIds);

    String dbName = Util::getDBName("sharemgnt_db");
    String strSql;
    strSql.format (_T("select f_user_id,f_login_name,f_display_name,f_mail_address,f_csf_level from %s.t_user where f_user_id in (%s);"),
        dbName.getCStr(), groupStr.getCStr ());

    ncDBRecords results;
    dbOper->Select (strSql, results);

    ncACSUserInfo userInfo;
    for (size_t i = 0; i < results.size (); ++i) {
        userInfo.id = std::move(results[i][0]);
        userInfo.account = std::move(results[i][1]);
        userInfo.visionName = std::move(results[i][2]);
        userInfo.email = std::move(results[i][3]);
        userInfo.csfLevel = Int::getValue(results[i][4]);

        userInfoMap.insert (pair<String, ncACSUserInfo>(userInfo.id, userInfo));
    }

    NC_ACS_SHAREMGNT_TRACE (_T("[END]this: %p"), this);
}

/* [notxpcom] void GetDepartInfoByIdBatch ([const] in StringVecRef departIds, in ncACSDepartInfoMapRef departInfoMap); */
NS_IMETHODIMP_(void) ncACSShareMgnt::GetDepartInfoByIdBatch(const vector<String> & departIds, map<String, ncACSDepartInfo> & departInfoMap)
{
    NC_ACS_SHAREMGNT_TRACE (_T("this: %p"), this);

    departInfoMap.clear ();

    if (departIds.size () == 0) {
        return;
    }

    nsCOMPtr<ncIDBOperator> dbOper = getter_AddRefs (GetDBOperator ());

    String groupStr = GenerateGroupStr (departIds);

    String dbName = Util::getDBName("sharemgnt_db");
    String strSql;
    strSql.format (_T("select f_department_id,f_name from %s.t_department where f_department_id in (%s);"),
        dbName.getCStr(), groupStr.getCStr ());

    ncDBRecords results;
    dbOper->Select (strSql, results);

    ncACSDepartInfo departInfo;
    for (size_t i = 0; i < results.size (); ++i) {
        departInfo.id = std::move(results[i][0]);
        departInfo.name = std::move(results[i][1]);

        departInfoMap.insert (pair<String, ncACSDepartInfo>(departInfo.id, departInfo));
    }

    NC_ACS_SHAREMGNT_TRACE (_T("this: %p"), this);
}

/* [notxpcom] void GetGroupInfoByIdBatch ([const] in StringVecRef groupIds, in ncGroupInfoMapRef groupInfoMap); */
NS_IMETHODIMP_(void) ncACSShareMgnt::GetGroupInfoByIdBatch(const vector<String> & groupIds, map<String, ncGroupInfo> & groupInfoMap)
{
    NC_ACS_SHAREMGNT_TRACE (_T("this: %p"), this);

    groupInfoMap.clear ();

    if (groupIds.size () == 0) {
        return;
    }

    nsCOMPtr<ncIDBOperator> dbOper = getter_AddRefs (GetDBOperator ());

    String groupStr = GenerateGroupStr (groupIds);

    String dbName = Util::getDBName("sharemgnt_db");
    String strSql;
    strSql.format (_T("select f_group_id,f_user_id,f_group_name,f_person_count from %s.t_person_group where f_group_id in (%s);"),
        dbName.getCStr(), groupStr.getCStr ());

    ncDBRecords results;
    dbOper->Select (strSql, results);

    ncGroupInfo groupInfo;
    for (size_t i = 0; i < results.size (); ++i) {
        groupInfo.id = std::move(results[i][0]);
        groupInfo.createrId = std::move(results[i][1]);
        groupInfo.groupName = std::move(results[i][2]);
        groupInfo.count = Int64::getValue (results[i][3]);;

        groupInfoMap.insert (pair<String, ncGroupInfo>(groupInfo.id, groupInfo));
    }

    NC_ACS_SHAREMGNT_TRACE (_T("this: %p"), this);
}

/* [notxpcom] void GetParentDepartIds ([const] in StringRef depId, in StringSetRef parentDepIds); */
NS_IMETHODIMP_(void) ncACSShareMgnt::GetParentDepartIds(const String & depId, set<String> & parentDepIds)
{
    NC_ACS_SHAREMGNT_TRACE (_T("this: %p"), this);
    parentDepIds.clear();

    nsCOMPtr<ncIDBOperator> dbOper = getter_AddRefs (GetDBOperator ());

    parentDepIds.clear();

    String strSql;
    String tmpId = depId;

    String dbName = Util::getDBName("sharemgnt_db");
    while (1) {
        String escTmpId;
        strSql.format (_T("select f_parent_department_id from %s.t_department_relation where f_department_id = '%s'"),
                        dbName.getCStr(), dbOper->EscapeEx(tmpId).getCStr ());

        ncDBRecords results;
        dbOper->Select (strSql, results);

        if(results.size () > 0) {
            parentDepIds.insert(results[0][0]);
            tmpId = results[0][0];
        }
        else {
            break;
        }
    }

    NC_ACS_SHAREMGNT_TRACE (_T("this: %p"), this);
}

/* [notxpcom] void GetSubDeps ([const] in StringRef depId, in ncACSDepartInfoVectorRef depInfos); */
NS_IMETHODIMP_(void) ncACSShareMgnt::GetSubDeps(const String& depId, vector<ncACSDepartInfo>& depInfos)
{
    NC_ACS_SHAREMGNT_TRACE (_T("begin this: %p, depId: %s"), this, depId.getCStr ());

    depInfos.clear ();
    nsCOMPtr<ncIDBOperator> dbOper = getter_AddRefs (GetDBOperator ());

    String dbName = Util::getDBName("sharemgnt_db");
    String strSql;
    strSql.format(_T("select f_department_id, f_name from %s.t_department ")
                  _T("where f_department_id in ")
                  _T("(select f_department_id from %s.t_department_relation where f_parent_department_id= '%s') ")
                  _T("order by f_priority, upper(f_name) "),
                  dbName.getCStr(), dbName.getCStr(), dbOper->EscapeEx(depId).getCStr());

    ncDBRecords results;
    dbOper->Select (strSql, results);

    ncACSDepartInfo depInfo;
    for (size_t i = 0; i < results.size (); ++i) {
        depInfo.id = std::move(results[i][0]);
        depInfo.name = std::move(results[i][1]);
        depInfos.push_back (depInfo);
    }

    NC_ACS_SHAREMGNT_TRACE (_T("end this: %p, depId: %s"), this, depId.getCStr ());
}

/* [notxpcom] bool GetDeptRootPath ([const] in string deptId, in StringVecRef path) */
NS_IMETHODIMP_(void) ncACSShareMgnt::GetDeptRootPath(const String& deptId, vector<String>& path)
{
    NC_ACS_SHAREMGNT_TRACE (_T("this: %p"), this);
    path.clear();

    nsCOMPtr<ncIDBOperator> dbOper = getter_AddRefs (GetDBOperator ());
    String strSql;

    String tmpId = deptId;
    path.push_back(deptId);

    String dbName = Util::getDBName("sharemgnt_db");
    while (1) {
        strSql.format (_T("select f_parent_department_id from %s.t_department_relation where f_department_id = '%s'"),
                        dbName.getCStr(), dbOper->EscapeEx(tmpId).getCStr ());

        ncDBRecords results;
        dbOper->Select (strSql, results);

        if(results.size () > 0) {
            path.push_back(results[0][0]);
            tmpId = results[0][0];
        }
        else {
            break;
        }
    }

    reverse(path.begin(),path.end());

    NC_ACS_SHAREMGNT_TRACE (_T("this: %p"), this);
}

/* [notxpcom] bool GetParentDeptRootPathName ([const] in StringRef deptId, in StringRef path) */
NS_IMETHODIMP_(void) ncACSShareMgnt::GetParentDeptRootPathName(const String& deptId, String& path)
{
    NC_ACS_SHAREMGNT_TRACE (_T("this: %p"), this);

    path.clear();

    nsCOMPtr<ncIDBOperator> dbOper = getter_AddRefs (GetDBOperator ());
    String strSql;

    vector<String> deptNames;
    GetParentDeptPath(deptId, deptNames);

    for (size_t i = 0; i < deptNames.size(); i++) {
        path += deptNames[i].getCStr ();
        if (i != deptNames.size () - 1) {
            path += "/";
        }
    }

    NC_ACS_SHAREMGNT_TRACE (_T("this: %p"), this);
}

/* [notxpcom] bool GetUserRootPath ([const] in string userId, in StringVecVecRef pathIds) */
NS_IMETHODIMP_(void) ncACSShareMgnt::GetUserRootPath(const String& userId, vector<vector<String> >& paths)
{
    NC_ACS_SHAREMGNT_TRACE (_T("this: %p"), this);

    for(size_t i = 0; i < paths.size(); i++)
        paths[i].clear();

    vector<String> directDeptIds;
    GetDirectBelongDepartmentIds (userId, directDeptIds);

    for(size_t i = 0; i < directDeptIds.size(); i++){
        vector<String> tmpPath;
        GetDeptRootPath(directDeptIds[i], tmpPath);

        paths.push_back(tmpPath);
    }

    NC_ACS_SHAREMGNT_TRACE (_T("this: %p"), this);
}


/* [notxpcom] void GetUsrIdsOutOfPermScope ([const] in StringRef userId, [const] in StringVecRef checkUserIds, in StringVecRef outScopeUserIds); */
NS_IMETHODIMP_(void) ncACSShareMgnt::GetUsrIdsOutOfPermScope(const String & userId, const vector<String>& checkUserIds, vector<String>& outScopeUserIds)
{
    NC_ACS_SHAREMGNT_TRACE (_T("this: %p"), this);

    outScopeUserIds.clear();

    if (IsAdminId(userId))
        return;

    // 获取用户的权限范围
    map<String, ncPermScopeObjInfo> scopeInfos;
    GetUserPermScopeInfos(userId, scopeInfos);

    for(size_t i = 0; i < checkUserIds.size(); ++i) {
        const String& tmpId = checkUserIds[i];

        if (userId == tmpId)
            continue;

        bool bInScope = false;

        // 用户在权限范围列表中
        if (scopeInfos.find(tmpId) != scopeInfos.end())
            bInScope = true;
        else {
            vector<vector<String> > pathVecIds;
            GetUserRootPath(tmpId, pathVecIds);

            // 用户到根组织路径中的某个节点在权限范围列表中，则用户在权限范围内
            for (size_t j = 0; j < pathVecIds.size(); j++) {
                vector<String>& path = pathVecIds[j];
                for (size_t k = 0; k < path.size(); k++) {
                    if (scopeInfos.find(path[k]) != scopeInfos.end()) {
                        bInScope = true;
                        break;
                    }
                }
            }
        }

        if (!bInScope)
            outScopeUserIds.push_back(tmpId);
    }

    NC_ACS_SHAREMGNT_TRACE (_T("this: %p"), this);
}

/* [notxpcom] bool CheckUsrInPermScope ([const] in StringRef userId, [const] in StringRef checkUserId);*/
NS_IMETHODIMP_(bool) ncACSShareMgnt::CheckUsrInPermScope(const String & userId, const String& checkUserId)
{
    NC_ACS_SHAREMGNT_TRACE (_T("this: %p"), this);

    // 检查权限范围是否开启
    if (!GetPermShareLimitStatus()) {
        return true;
    }

    // 获取用户的权限范围
    map<String, ncPermScopeObjInfo> scopeInfos;
    GetUserPermScopeInfos(userId, scopeInfos);

    bool bInScope = false;
    if (userId == checkUserId)
        bInScope = true;
    else {
        // 用户在权限范围列表中
        if (scopeInfos.find(checkUserId) != scopeInfos.end())
            bInScope = true;
        else {
            vector<vector<String> > pathVecIds;
            GetUserRootPath(checkUserId, pathVecIds);

            // 用户到根组织路径中的某个节点在权限范围列表中，则用户在权限范围内
            for (size_t j = 0; j < pathVecIds.size(); j++) {
                vector<String>& path = pathVecIds[j];
                for (size_t k = 0; k < path.size(); k++) {
                    if (scopeInfos.find(path[k]) != scopeInfos.end()) {
                        bInScope = true;
                        break;
                    }
                }
            }
        }
    }

    NC_ACS_SHAREMGNT_TRACE (_T("this: %p"), this);

    return bInScope;
}

/* [notxpcom] bool CheckDeptInPermScope ([const] in StringRef userId, [const] in StringRef checkDeptId);*/
NS_IMETHODIMP_(bool) ncACSShareMgnt::CheckDeptInPermScope(const String & userId, const String& checkDeptId)
{
    NC_ACS_SHAREMGNT_TRACE (_T("this: %p"), this);

    // 获取用户的权限范围
    map<String, ncPermScopeObjInfo> scopeInfos;
    GetUserPermScopeInfos(userId, scopeInfos);

    bool bInScope = false;
    if (scopeInfos.find(checkDeptId) != scopeInfos.end())
        bInScope = true;
    else {
        // 获取部门的路径
        vector<String> depPath;
        GetDeptRootPath(checkDeptId, depPath);

        // 如果部门路径上某个节点在在范围列表中，部门在范围内
        for (size_t i = 0; i < depPath.size(); i++){
            if (scopeInfos.find(depPath[i]) != scopeInfos.end()) {
                bInScope = true;
                break;
             }
        }
    }

    NC_ACS_SHAREMGNT_TRACE (_T("this: %p"), this);
    return bInScope;
}

/* [notxpcom] void GetOrgIdsByScopeInfo ([const] in ncObjIdScopeInfoMapRef scopeInfos, in StringSetRef orgInfos); */
NS_IMETHODIMP_(void) ncACSShareMgnt::GetOrgIdsByScopeInfo(const map<String, ncPermScopeObjInfo>& scopeInfos, vector<ncOrganizationInfo>& orgInfos)
{
    NC_ACS_SHAREMGNT_TRACE (_T("this: %p"), this);

    // 根据共享范围对象获取根组织
    set<String> tmpIds;
    vector<ncOrganizationInfo> tmpOrgInfos;
    map<String, ncOrganizationInfo> tmpMapOrgInfos;
    for(map<String, ncPermScopeObjInfo>::const_iterator iter = scopeInfos.begin(); iter != scopeInfos.end(); ++iter){
        const ncPermScopeObjInfo& objInfo = iter->second;
        tmpOrgInfos.clear();

        if(objInfo.type == IOC_USER)
            GetOrgInfoByDeptId(objInfo.parentId, tmpOrgInfos);
        else
            GetOrgInfoByDeptId(objInfo.id, tmpOrgInfos);

        for(size_t j = 0; j < tmpOrgInfos.size(); j++){
            ncOrganizationInfo& orgInfo = tmpOrgInfos[j];

            // 去重
            if (tmpIds.find(orgInfo.id) != tmpIds.end())
                continue;

            tmpMapOrgInfos.insert(pair<String, ncOrganizationInfo> (orgInfo.id, orgInfo));
            tmpIds.insert(orgInfo.id);
        }
    }

    // 对根组织id按照权重进行重排序
    vector<String> departIds;
    if (tmpIds.size () > 0) {
        String groupStr = GenerateGroupStrBySet(tmpIds);
        SortObjIdsWithPriority (groupStr, IOC_DEPARTMENT, departIds);
    }

    for(size_t i = 0; i < departIds.size(); i++) {
        orgInfos.push_back(tmpMapOrgInfos[departIds[i]]);
    }

    NC_ACS_SHAREMGNT_TRACE (_T("this: %p"), this);
}


/* [notxpcom] void GetScopeOrgInfo ([const] in StringRef userId, in ncOrgaInfoPairVecRef organs); */
NS_IMETHODIMP_(void) ncACSShareMgnt::GetScopeOrgInfo(const String & userId, vector<pair<ncOrganizationInfo,bool> > & organs)
{
    NC_ACS_SHAREMGNT_TRACE (_T("this: %p"), this);

    // 获取用户的权限范围
    map<String, ncPermScopeObjInfo> scopeInfos;
    GetUserPermScopeInfos(userId, scopeInfos);

    // 根据共享范围对象获取根组织
    vector<ncOrganizationInfo> allOrgInfos;
    GetOrgIdsByScopeInfo(scopeInfos, allOrgInfos);

    // 检查组织是否可配置
    for (size_t i = 0; i < allOrgInfos.size(); i++){
        // 如果组织id在权限范围列表内，则可配置
        bool isConfigable = false;
        map<String, ncPermScopeObjInfo>::iterator mapIter = scopeInfos.find(allOrgInfos[i].id);
        if (mapIter != scopeInfos.end() && mapIter->second.type == IOC_DEPARTMENT)
            isConfigable = true;

        organs.push_back(pair<ncOrganizationInfo,bool>(allOrgInfos[i], isConfigable));
    }

    NC_ACS_SHAREMGNT_TRACE (_T("this: %p"), this);
}

/* [notxpcom] void GetScopeSubDeps ([const] in StringRef userId,[const] in StringRef depId, in ncACSDepartInfoPairVectorRef depInfos); */
NS_IMETHODIMP_(void) ncACSShareMgnt::GetScopeSubDeps(const String& userId, const String& depId, vector<pair<ncACSDepartInfo,bool> >& depInfos)
{
    NC_ACS_SHAREMGNT_TRACE (_T("begin this: %p, depId: %s"), this, depId.getCStr ());

    depInfos.clear();

    // 获取用户的权限范围
    map<String, ncPermScopeObjInfo> scopeInfos;
    GetUserPermScopeInfos(userId, scopeInfos);

    bool bGetAll = false;
   // 如果部门路径上某个节点在在范围列表中，则获取所有
    vector<String> curPath;
    GetDeptRootPath(depId, curPath);
    for (size_t i = 0; i < curPath.size(); i++){
        map<String, ncPermScopeObjInfo>::iterator mapIter = scopeInfos.find(curPath[i]);
        if (mapIter != scopeInfos.end() && mapIter->second.type == IOC_DEPARTMENT) {
            bGetAll = true;
            break;
         }
    }

    if (!bGetAll) {
        // 获取权限范围对象的路径，并且查找路径上的子部门
        set<String> subDepIds;
        for ( map<String, ncPermScopeObjInfo>::iterator mapIter = scopeInfos.begin(); mapIter != scopeInfos.end(); ++mapIter){
            ncPermScopeObjInfo& objInfo = mapIter->second;
            vector<String> pathIds;

            if (objInfo.type == IOC_DEPARTMENT) {
                GetDeptRootPath(objInfo.id, pathIds);
            } else {
                GetDeptRootPath(objInfo.parentId, pathIds);
            }

            // 路径上存在depId，且depId存在子部门
            vector<String>::iterator iter = find(pathIds.begin(), pathIds.end(), depId);
            if (iter != pathIds.end() && ++iter != pathIds.end()) {
                subDepIds.insert(*iter);
            }
        }

        // 对子部门id按照权重值进行重排序
        vector<String> departIds;
        if (subDepIds.size () > 0) {
            String groupStr = GenerateGroupStrBySet (subDepIds);
            SortObjIdsWithPriority (groupStr, IOC_DEPARTMENT, departIds);
        }

        for (size_t i = 0; i < departIds.size (); i++) {
            ncACSDepartInfo deptInfo;
            String subId = departIds[i];
            GetDepartInfoById (subId, deptInfo);

            // 如果子部门id在权限范围的id列表中，则子部门可配置，否则不可配置
            if (scopeInfos.find(subId) != scopeInfos.end())
                depInfos.push_back(pair<ncACSDepartInfo,bool>(deptInfo, true));
            else
                depInfos.push_back(pair<ncACSDepartInfo,bool>(deptInfo, false));
        }
    }

    if (bGetAll) {
        vector<ncACSDepartInfo> subDepInfos;
        GetSubDeps(depId, subDepInfos);

        for (size_t i = 0; i < subDepInfos.size(); i++)
            depInfos.push_back(pair<ncACSDepartInfo,bool>(subDepInfos[i], true));
    }

    NC_ACS_SHAREMGNT_TRACE (_T("end this: %p, depId: %s"), this, depId.getCStr ());
}

/* [notxpcom] void GetScopeSubUsers ([const] in StringRef userId, [const] in StringRef depId, in ncACSUserInfoVecRef userInfos); */
NS_IMETHODIMP_(void) ncACSShareMgnt::GetScopeSubUsers(const String& userId, const String& depId, vector<ncACSUserInfo>& userInfos)
{
    NC_ACS_SHAREMGNT_TRACE (_T("begin this: %p, depId: %s"), this, depId.getCStr ());

    userInfos.clear ();

    // 获取用户的权限范围
    map<String, ncPermScopeObjInfo> scopeInfos;
    GetUserPermScopeInfos(userId, scopeInfos);

    // 获取部门的路径
    vector<String> depPath;
    GetDeptRootPath(depId, depPath);

    // 检查部门是否可配置
    bool isConfigable = false;
    for (size_t i = 0; i < depPath.size(); i++){
        map<String, ncPermScopeObjInfo>::iterator mapIter = scopeInfos.find(depPath[i]);
        if (mapIter != scopeInfos.end() && mapIter->second.type == IOC_DEPARTMENT){
            isConfigable = true;
            break;
        }
    }

    // 部门可配置，则获取所有子用户
    if (isConfigable){
        GetSubUsers(depId, userInfos);
    } else {
        // 部门不可配置，则在权限范围对象的路径上查找子用户
        set<String> subUserIds;
        for (map<String, ncPermScopeObjInfo>::iterator iter = scopeInfos.begin(); iter != scopeInfos.end(); ++iter){
            ncPermScopeObjInfo& objInfo = iter->second;
            if (objInfo.type == IOC_USER){
                // depId在用户路径上，且存在子用户
                if (objInfo.parentId == depId) {
                    subUserIds.insert(objInfo.id);
                }
            }
        }

        // 将子用户id按照权重值进行重排序
        vector<String> userIds;
        if (subUserIds.size () > 0) {
            String groupStr = GenerateGroupStrBySet (subUserIds);
            SortObjIdsWithPriority (groupStr, IOC_USER, userIds);
        }

        for (size_t i = 0; i < userIds.size (); i++) {
            ncACSUserInfo userInfo;
            GetUserInfoById(userIds[i], userInfo);
            userInfos.push_back(userInfo);
        }
    }

    NC_ACS_SHAREMGNT_TRACE (_T("end this: %p, depId: %s"), this, depId.getCStr ());
}

/* [notxpcom] void SearchScopeOrganization ([const] in StringRef userId, [const] in StringRef key, in int start, in int limit, in ncACSUserInfoVecRef userInfos, in ncACSDepartInfoVectorRef departInfos); */
NS_IMETHODIMP_(void) ncACSShareMgnt::SearchScopeOrganization(const String & userId, const String & key, int start, int limit, vector<ncACSUserInfo>& userInfos, vector<ncACSDepartInfo>& departInfos)
{
    NC_ACS_SHAREMGNT_TRACE (_T("begin this: %p"), this);
    userInfos.clear();
    departInfos.clear();

    if (key == String::EMPTY || limit == 0) {
        return;
    }

    if (limit == -1) {
        limit = INT_MAX;
    }

    nsCOMPtr<ncIDBOperator> dbOper = getter_AddRefs (GetDBOperator ());

    // 获取用户的权限范围
    map<String, ncPermScopeObjInfo> scopeInfos;
    GetUserPermScopeInfos(userId, scopeInfos);

    // 根据权限范围获取组织信息
    vector<ncOrganizationInfo> orgInfos;
    GetOrgIdsByScopeInfo(scopeInfos, orgInfos);

    vector<String> tmpOuIds;
    for (size_t i = 0; i < orgInfos.size (); ++i)
        tmpOuIds.push_back (orgInfos[i].id);

    if (tmpOuIds.empty())
        return;

    String groupStr = GenerateGroupStr (tmpOuIds);

    // 获取搜索配置信息
    String searchUserConfig = GetShareMgntConfig("search_user_config");
    JSON::Value searchUserConfigJson;
    JSON::Reader::read (searchUserConfigJson, searchUserConfig.getCStr (), searchUserConfig.getLength ());
    bool exact_search_user = searchUserConfigJson["exactSearch"].b ();
    int searchRange = searchUserConfigJson["searchRange"].i ();
    int searchResults = searchUserConfigJson["searchResults"].i ();

    // 连接查询t_ou_user和t_user表
    String strSql;
    String escKey = dbOper->EscapeEx(key);
    String escLikeKey = dbOper->EscapeLikeClause(key);
    String searchStr = GenerateSearchStr (exact_search_user, searchRange, escKey, escLikeKey);
    String dbName = Util::getDBName("sharemgnt_db");
    strSql.format (_T("select distinct t_user.f_user_id,t_user.f_login_name,t_user.f_display_name,t_user.f_mail_address,t_user.f_csf_level, t_user.f_priority ")
                    _T("from %s.t_user, %s.t_ou_user ")
                    _T("where t_ou_user.f_ou_id in (%s) and ")
                    _T("t_user.f_status = '0' and t_user.f_auto_disable_status = '0' and ")
                    _T("t_user.f_user_id = t_ou_user.f_user_id and %s ;"),
                    dbName.getCStr(), dbName.getCStr(),
                    groupStr.getCStr (), searchStr.getCStr ());
    ncDBRecords results;
    dbOper->Select (strSql, results);

    int count = 0;
    for (size_t i = 0; i < results.size (); ++i) {
        String userId = results[i][0];
        bool isView = false;
        ncACSUserInfo userInfo;

        // 用户在权限范围列表中，则添加
        if (scopeInfos.find(userId) != scopeInfos.end())
        {
            // 判断用户是否未分配
            vector<String> tmpIds;
            GetDirectBelongDepartmentIds(userId, tmpIds);
            if (0 != tmpIds.size())
            {
                userInfo.belongDepartId = tmpIds[0];
            }
            isView = true;
        }
        else {
            vector<vector<String> > pathVecIds;
            GetUserRootPath(userId, pathVecIds);

            // 用户到根组织路径中的某个节点在权限范围列表中，则添加此用户
            bool hasFind = false;
            for (size_t j = 0; j < pathVecIds.size(); j++) {
                vector<String>& path = pathVecIds[j];
                for (size_t k = 0; k < path.size(); k++) {
                    if (scopeInfos.find(path[k]) != scopeInfos.end()) {
                        isView = true;
                        hasFind = true;
                        userInfo.belongDepartId = path[path.size() - 1];
                        break;
                    }
                }
                if (hasFind) {
                    break;
                }
            }
        }

        if (isView){
            if (count >= start) {
                userInfo.id = std::move(results[i][0]);
                userInfo.account = (searchResults == SHOW_LOGIN_AND_DISPLAY) ? std::move(results[i][1]) : "";
                userInfo.visionName = std::move(results[i][2]);
                userInfo.email = std::move(results[i][3]);
                userInfo.csfLevel = Int::getValue(results[i][4]);
                userInfos.push_back (userInfo);
            }
            count++;
        }
        if (count == start + limit) {
            break;
        }
    }

    bool only_share_to_user = Int::getValue(GetShareMgntConfig("only_share_to_user")) == 1;

    // 模糊搜索且可以共享给组织和部门时, 连接查询t_ou_department和t_department表
    if (!exact_search_user && !only_share_to_user) {
        if (userInfos.size () < limit) {
            limit -= (int)userInfos.size ();
            start = (userInfos.size () > 0) ? 0 : (start - SearchScopeOrganizationCount(userId, key, false));

            strSql.format(_T("select t_department.f_department_id,t_department.f_name ")
                          _T("from %s.t_department, %s.t_ou_department ")
                          _T("where ")
                          _T("t_ou_department.f_ou_id in (%s) and ")
                          _T("t_department.f_department_id = t_ou_department.f_department_id and ")
                          _T("(t_department.f_name = '%s' or t_department.f_name like '%%%s%%') ")
                          _T("order by t_department.f_priority, ")
                          _T("case when t_department.f_name = '%s' then 0 when t_department.f_name like '%%%s%%' then 1  else 2 end, ")
                          _T("upper(t_department.f_name), t_department.f_department_id ")
                          _T("limit %d, %d;"),
                          dbName.getCStr(), dbName.getCStr(),
                          groupStr.getCStr(), escKey.getCStr(), escLikeKey.getCStr(), escKey.getCStr(), escLikeKey.getCStr(), start, limit);

            dbOper->Select (strSql, results);

            ncACSDepartInfo depInfo;
            for (size_t i = 0; i < results.size (); ++i) {
                String deptId = results[i][0];
                bool isView = false;
                if (scopeInfos.find(deptId) != scopeInfos.end())
                    isView = true;
                else {
                    // 获取部门的路径
                    vector<String> depPath;
                    GetDeptRootPath(deptId, depPath);

                    for (size_t j = 0; j < depPath.size(); j++) {
                        if (scopeInfos.find(depPath[j]) != scopeInfos.end()){
                            isView = true;
                            break;
                        }
                    }
                }

                if (isView && departInfos.size() < limit) {
                    depInfo.id = results[i][0];
                    depInfo.name = results[i][1];
                    String parentPath("");
                    GetParentDeptRootPathName(deptId, parentPath);
                    depInfo.path = parentPath.getCStr ();
                    departInfos.push_back (depInfo);
                }
            }
        }
    }

    NC_ACS_SHAREMGNT_TRACE (_T("end this: %p"), this);
}

/* [notxpcom] int SearchScopeOrganizationCount ([const] in StringRef userId, [const] in StringRef key, in bool searchDepart); */
NS_IMETHODIMP_(int) ncACSShareMgnt::SearchScopeOrganizationCount(const String & userId, const String & key, bool searchDepart)
{
    NC_ACS_SHAREMGNT_TRACE (_T("begin this: %p"), this);

    nsCOMPtr<ncIDBOperator> dbOper = getter_AddRefs (GetDBOperator ());

    // 获取用户的权限范围
    map<String, ncPermScopeObjInfo> scopeInfos;
    GetUserPermScopeInfos(userId, scopeInfos);

    // 根据权限范围获取组织信息
    vector<ncOrganizationInfo> orgInfos;
    GetOrgIdsByScopeInfo(scopeInfos, orgInfos);

    vector<String> tmpOuIds;
    for (size_t i = 0; i < orgInfos.size (); ++i)
        tmpOuIds.push_back (orgInfos[i].id);

    if (tmpOuIds.empty())
        return 0;

    String groupStr = GenerateGroupStr (tmpOuIds);

    // 获取搜索配置信息
    String searchUserConfig = GetShareMgntConfig("search_user_config");
    JSON::Value searchUserConfigJson;
    JSON::Reader::read (searchUserConfigJson, searchUserConfig.getCStr (), searchUserConfig.getLength ());
    bool exact_search_user = searchUserConfigJson["exactSearch"].b ();
    int searchRange = searchUserConfigJson["searchRange"].i ();
    int searchResults = searchUserConfigJson["searchResults"].i ();

    // 连接查询t_ou_user和t_user表
    String strSql;
    String escKey = dbOper->EscapeEx(key);
    String escLikeKey = dbOper->EscapeLikeClause(key);
    String searchStr = GenerateSearchStr (exact_search_user, searchRange, escKey, escLikeKey);
    String dbName = Util::getDBName("sharemgnt_db");
    strSql.format (_T("select distinct t_user.f_user_id, ")
                    _T("t_user.f_priority, t_user.f_login_name, t_user.f_display_name ")
                    _T("from %s.t_user, %s.t_ou_user ")
                    _T("where t_ou_user.f_ou_id in (%s) and ")
                    _T("t_user.f_status = '0' and t_user.f_auto_disable_status = '0' and ")
                    _T("t_user.f_user_id = t_ou_user.f_user_id and %s;"),
                    dbName.getCStr(), dbName.getCStr(),
                    groupStr.getCStr (), searchStr.getCStr ());
    ncDBRecords results;
    dbOper->Select (strSql, results);

    int count = 0;
    ncACSUserInfo userInfo;
    for (size_t i = 0; i < results.size (); ++i) {
        String userId = results[i][0];
        // 用户在权限范围列表中，则添加
        if (scopeInfos.find(userId) != scopeInfos.end()) {
            count++;
            break;
        }

        vector<vector<String> > pathVecIds;
        GetUserRootPath(userId, pathVecIds);

        // 用户到根组织路径中的某个节点在权限范围列表中，则添加此用户
        bool hasFind = false;
        for (size_t j = 0; j < pathVecIds.size(); j++) {
            vector<String>& path = pathVecIds[j];
            for (size_t k = 0; k < path.size(); k++) {
                if (scopeInfos.find(path[k]) != scopeInfos.end()) {
                    count++;
                    hasFind = true;
                    break;
                }
            }
            if (hasFind) {
                break;
            }
        }
    }

    // 模糊搜索时, 连接查询t_ou_department和t_department表
    if (!exact_search_user && searchDepart) {
        strSql.format (_T("select t_department.f_department_id ")
            _T("from %s.t_department, %s.t_ou_department ")
            _T("where ")
            _T("t_ou_department.f_ou_id in (%s) and ")
            _T("t_department.f_department_id = t_ou_department.f_department_id and ")
            _T("(t_department.f_name = '%s' or t_department.f_name like '%%%s%%');"),
            dbName.getCStr(), dbName.getCStr(),
            groupStr.getCStr (), escKey.getCStr (), escLikeKey.getCStr ());

        dbOper->Select (strSql, results);
        for (size_t i = 0; i < results.size (); ++i) {
            String deptId = results[i][0];

            if (scopeInfos.find(deptId) != scopeInfos.end())
                count++;
            else {
                // 获取部门的路径
                vector<String> depPath;
                GetDeptRootPath(deptId, depPath);

                for (size_t j = 0; j < depPath.size(); j++) {
                    if (scopeInfos.find(depPath[j]) != scopeInfos.end()){
                        count++;
                        break;
                    }
                }
            }
        }
    }

    NC_ACS_SHAREMGNT_TRACE (_T("end this: %p"), this);

    return count;
}

/* [notxpcom] void GetSubUsers ([const] in StringRef depId, in ncACSUserInfoVecRef userInfos); */
NS_IMETHODIMP_(void) ncACSShareMgnt::GetSubUsers(const String& depId, vector<ncACSUserInfo>& userInfos)
{
    NC_ACS_SHAREMGNT_TRACE (_T("begin this: %p, depId: %s"), this, depId.getCStr ());

    userInfos.clear ();
    nsCOMPtr<ncIDBOperator> dbOper = getter_AddRefs (GetDBOperator ());

    String dbName = Util::getDBName("sharemgnt_db");
    String strSql;
    strSql.format(_T("select f_user_id, f_login_name, f_display_name, f_mail_address, f_status, f_csf_level, f_auto_disable_status ")
                  _T("from %s.t_user where f_user_id in ")
                  _T("(select f_user_id from %s.t_user_department_relation where f_department_id= '%s') ")
                  _T("order by f_priority, upper(f_display_name);"),
                  dbName.getCStr(), dbName.getCStr(), dbOper->EscapeEx(depId).getCStr());

    ncDBRecords results;
    dbOper->Select (strSql, results);

    ncACSUserInfo userInfo;
    for (size_t i = 0; i < results.size (); ++i) {
        userInfo.id = std::move(results[i][0]);
        userInfo.account = std::move(results[i][1]);
        userInfo.visionName = std::move(results[i][2]);
        userInfo.email = std::move(results[i][3]);
        userInfo.enableStatus = (results[i][4] == "0") && (results[i][6] == "0");
        userInfo.csfLevel = Int::getValue(results[i][5]);

        userInfos.push_back (userInfo);
    }

    NC_ACS_SHAREMGNT_TRACE (_T("end this: %p, depId: %s"), this, depId.getCStr ());
}

/* [notxpcom] void GetAllUser (in ncACSUserInfoVecRef userInfos); */
NS_IMETHODIMP_(void) ncACSShareMgnt::GetAllUser(vector<ncACSUserInfo>& userInfos)
{
    NC_ACS_SHAREMGNT_TRACE (_T("begin this: %p"), this);

    userInfos.clear ();
    nsCOMPtr<ncIDBOperator> dbOper = getter_AddRefs (GetDBOperator ());
    String strSql;
    String dbName = Util::getDBName("sharemgnt_db");
    strSql.format (_T("select f_user_id, f_login_name, f_display_name, f_mail_address from %s.t_user"), dbName.getCStr());

    ncDBRecords results;
    dbOper->Select (strSql, results);
    ncACSUserInfo userInfo;
    for (size_t i = 0; i < results.size (); ++i) {
        userInfo.id = std::move(results[i][0]);
        userInfo.account = std::move(results[i][1]);
        userInfo.visionName = std::move(results[i][2]);
        userInfo.email = std::move(results[i][3]);

        userInfos.push_back (userInfo);
    }

    NC_ACS_SHAREMGNT_TRACE (_T("end this: %p"), this);
}

/* [notxpcom] void SearchGivenOrganization ([const] in StringVecRef orgIds, [const] in StringRef key, in int start, in int limit, in ncACSUserInfoVecRef userInfos, in ncACSDepartInfoVectorRef departInfos); */
NS_IMETHODIMP_(void) ncACSShareMgnt::SearchGivenOrganization(const vector<String>& orgIds, const String & key, int start, int limit, vector<ncACSUserInfo> & userInfos, vector<ncACSDepartInfo> & departInfos)
{
    NC_ACS_SHAREMGNT_TRACE (_T("begin this: %p"), this);

    userInfos.clear ();
    departInfos.clear ();

    if (key == String::EMPTY) {
        return;
    }

    if (limit == -1) {
        limit = INT_MAX;
    }

    nsCOMPtr<ncIDBOperator> dbOper = getter_AddRefs (GetDBOperator ());

    if (orgIds.size () == 0) {
        return;
    }

    String groupStr = GenerateGroupStr (orgIds);

    // 获取搜索配置信息
    String searchUserConfig = GetShareMgntConfig("search_user_config");
    JSON::Value searchUserConfigJson;
    JSON::Reader::read (searchUserConfigJson, searchUserConfig.getCStr (), searchUserConfig.getLength ());
    bool exact_search_user = searchUserConfigJson["exactSearch"].b ();
    int searchRange = searchUserConfigJson["searchRange"].i ();
    int searchResults = searchUserConfigJson["searchResults"].i ();

    // 连接查询t_ou_user和t_user表
    String strSql;
    String escKey = dbOper->EscapeEx(key);
    String escLikeKey = dbOper->EscapeLikeClause(key);
    String searchStr = GenerateSearchStr (exact_search_user, searchRange, escKey, escLikeKey);
    String dbName = Util::getDBName("sharemgnt_db");
    strSql.format (_T("select distinct t_user.f_user_id,t_user.f_login_name,t_user.f_display_name,t_user.f_mail_address,t_user.f_csf_level, t_user.f_priority ")
                    _T("from %s.t_user, %s.t_ou_user ")
                    _T("where t_ou_user.f_ou_id in (%s) and ")
                    _T("t_user.f_status = '0' and t_user.f_auto_disable_status = '0' and ")
                    _T("t_user.f_user_id = t_ou_user.f_user_id and %s ")
                    _T("limit %d,%d;"),
                    dbName.getCStr(), dbName.getCStr(),
                    groupStr.getCStr (), searchStr.getCStr (), start, limit);
    ncDBRecords results;
    dbOper->Select (strSql, results);

    ncACSUserInfo userInfo;
    for (size_t i = 0; i < results.size (); ++i) {
        userInfo.id = std::move(results[i][0]);
        userInfo.account = (searchResults == SHOW_LOGIN_AND_DISPLAY) ? std::move(results[i][1]) : "";
        userInfo.visionName = std::move(results[i][2]);
        userInfo.email = std::move(results[i][3]);
        userInfo.csfLevel = Int::getValue(results[i][4]);

        userInfos.push_back (userInfo);
    }

    bool only_share_to_user = Int::getValue(GetShareMgntConfig("only_share_to_user")) == 1;

    // 模糊搜索且可以共享给组织和部门时, 连接查询t_ou_department和t_department表
    if (!exact_search_user && !only_share_to_user) {
        if (userInfos.size () < limit) {
            limit = limit - (int)userInfos.size ();
            start = (userInfos.size () > 0) ? 0 : (start - SearchGivenOrganizationCount(orgIds, key, false));

            strSql.format(_T("select t_department.f_department_id,t_department.f_name ")
                          _T("from %s.t_department, %s.t_ou_department ")
                          _T("where ")
                          _T("t_ou_department.f_ou_id in (%s) and ")
                          _T("t_department.f_department_id = t_ou_department.f_department_id and ")
                          _T("(t_department.f_name like '%%%s%%') ")
                          _T("order by t_department.f_priority, ")
                          _T("case when t_department.f_name = '%s' then 0 when t_department.f_name like '%%%s%%' then 1  else 2 end, ")
                          _T("upper(t_department.f_name), t_department.f_department_id ")
                          _T("limit %d,%d;"),
                          dbName.getCStr(), dbName.getCStr(),
                          groupStr.getCStr(), escLikeKey.getCStr(), escKey.getCStr(), escLikeKey.getCStr(), start, limit);

            dbOper->Select (strSql, results);

            ncACSDepartInfo depInfo;
            for (size_t i = 0; i < results.size (); ++i) {
                depInfo.id = std::move(results[i][0]);
                depInfo.name = std::move(results[i][1]);
                String parentPath("");
                GetParentDeptRootPathName(depInfo.id, parentPath);
                depInfo.path = parentPath.getCStr ();
                departInfos.push_back (depInfo);
            }
        }
    }

    NC_ACS_SHAREMGNT_TRACE (_T("end this: %p"), this);
}

/* [notxpcom] int SearchGivenOrganizationCount ([const] in StringVecRef orgIds, [const] in StringRef key, in bool searchDepart); */
NS_IMETHODIMP_(int) ncACSShareMgnt::SearchGivenOrganizationCount(const vector<String>& orgIds, const String & key, bool searchDepart)
{
    NC_ACS_SHAREMGNT_TRACE (_T("begin this: %p"), this);

    nsCOMPtr<ncIDBOperator> dbOper = getter_AddRefs (GetDBOperator ());

    if (orgIds.size () == 0) {
        return 0;
    }

    String groupStr = GenerateGroupStr (orgIds);

    // 获取搜索配置信息
    String searchUserConfig = GetShareMgntConfig("search_user_config");
    JSON::Value searchUserConfigJson;
    JSON::Reader::read (searchUserConfigJson, searchUserConfig.getCStr (), searchUserConfig.getLength ());
    bool exact_search_user = searchUserConfigJson["exactSearch"].b ();
    int searchRange = searchUserConfigJson["searchRange"].i ();
    int searchResults = searchUserConfigJson["searchResults"].i ();

    // 连接查询t_ou_user和t_user表
    String strSql;
    String escKey = dbOper->EscapeEx(key);
    String escLikeKey = dbOper->EscapeLikeClause(key);
    String searchStr = GenerateSearchStr (exact_search_user, searchRange, escKey, escLikeKey);
    String dbName = Util::getDBName("sharemgnt_db");
    strSql.format (_T("select count(distinct t_user.f_user_id) ")
                    _T("from %s.t_user, %s.t_ou_user ")
                    _T("where t_ou_user.f_ou_id in (%s) and ")
                    _T("t_user.f_status = '0' and t_user.f_auto_disable_status = '0' and ")
                    _T("t_user.f_user_id = t_ou_user.f_user_id and %s;"),
                    dbName.getCStr(), dbName.getCStr(),
                    groupStr.getCStr (), searchStr.getCStr ());
    ncDBRecords results;
    dbOper->Select (strSql, results);

    int count = Int::getValue(results[0][0]);

    // 模糊搜索时, 连接查询t_ou_department和t_department表
    if (!exact_search_user && searchDepart) {
        strSql.format (_T("select count(t_department.f_department_id) ")
            _T("from %s.t_department, %s.t_ou_department ")
            _T("where ")
            _T("t_ou_department.f_ou_id in (%s) and ")
            _T("t_department.f_department_id = t_ou_department.f_department_id and ")
            _T("(t_department.f_name like '%%%s%%');"),
            dbName.getCStr(), dbName.getCStr(),
            groupStr.getCStr (), escLikeKey.getCStr ());

        dbOper->Select (strSql, results);

        count += Int::getValue(results[0][0]);
    }

    NC_ACS_SHAREMGNT_TRACE (_T("end this: %p"), this);

    return count;
}

/* [notxpcom] void SearchAllOrganization ([const] in StringRef userId, [const] in StringRef key, in int start, in int limit, in ncACSUserInfoVecRef userInfos, in ncACSDepartInfoVectorRef departInfos); */
NS_IMETHODIMP_(void) ncACSShareMgnt::SearchAllOrganization(const String & userId, const String & key, int start, int limit, vector<ncACSUserInfo> & userInfos, vector<ncACSDepartInfo> & departInfos)
{
    NC_ACS_SHAREMGNT_TRACE (_T("begin this: %p"), this);

    userInfos.clear ();
    departInfos.clear ();

    if (key == String::EMPTY) {
        return;
    }

    if (limit == -1) {
        limit = INT_MAX;
    }

    nsCOMPtr<ncIDBOperator> dbOper = getter_AddRefs (GetDBOperator ());

    // 判断用户是否未分配
    vector<String> tmpIds;
    GetDirectBelongDepartmentIds (userId, tmpIds);
    if (0 == tmpIds.size ()) {
        return;
    }

    String dbName = Util::getDBName("sharemgnt_db");
    // 获取用户所属的组织id
    String strSql;
    strSql.format (_T("select f_department_id from %s.t_department where f_is_enterprise = 1;"), dbName.getCStr());

    ncDBRecords results;
    dbOper->Select (strSql, results);

    if (results.size () == 0) {
        return;
    }

    vector<String> tmpOuIds;
    for (size_t i = 0; i < results.size (); ++i) {
        tmpOuIds.push_back (results[i][0]);
    }

    SearchGivenOrganization(tmpOuIds, key, start, limit, userInfos, departInfos);

    NC_ACS_SHAREMGNT_TRACE (_T("end this: %p"), this);
}

/* [notxpcom] int SearchAllOrganizationCount ([const] in StringRef userId, [const] in StringRef key); */
NS_IMETHODIMP_(int) ncACSShareMgnt::SearchAllOrganizationCount(const String & userId, const String & key)
{
    NC_ACS_SHAREMGNT_TRACE (_T("begin this: %p"), this);

    nsCOMPtr<ncIDBOperator> dbOper = getter_AddRefs (GetDBOperator ());

    // 判断用户是否未分配
    vector<String> tmpIds;
    GetDirectBelongDepartmentIds (userId, tmpIds);
    if (0 == tmpIds.size ()) {
        return 0;
    }

    String dbName = Util::getDBName("sharemgnt_db");
    // 获取用户所属的组织id
    String strSql;
    strSql.format (_T("select f_department_id from %s.t_department where f_is_enterprise = 1;"), dbName.getCStr());

    ncDBRecords results;
    dbOper->Select (strSql, results);

    if (results.size () == 0) {
        return 0;
    }

    vector<String> tmpOuIds;
    for (size_t i = 0; i < results.size (); ++i) {
        tmpOuIds.push_back (results[i][0]);
    }

    int count = SearchGivenOrganizationCount(tmpOuIds, key, true);

    NC_ACS_SHAREMGNT_TRACE (_T("end this: %p"), this);

    return count;
}

/* [notxpcom] void SearchContactGroup ([const] in StringRef userId, [const] in StringRef key, in int start, in int limit, in ncACSUserInfoVecRef userInfos, in ncGroupInfoVecRef groupInfos); */
NS_IMETHODIMP_(void) ncACSShareMgnt::SearchContactGroup(const String & userId, const String & key, int start, int limit, vector<ncACSUserInfo> & userInfos, vector<ncGroupInfo> & groupInfos)
{
    NC_ACS_SHAREMGNT_TRACE (_T("begin this: %p"), this);

    userInfos.clear ();
    groupInfos.clear ();

    if (key == String::EMPTY) {
        return;
    }

    if (limit == -1) {
        limit = INT_MAX;
    }

    nsCOMPtr<ncIDBOperator> dbOper = getter_AddRefs (GetDBOperator ());
    String escUserId = dbOper->EscapeEx(userId);

    String dbName = Util::getDBName("sharemgnt_db");
    // 获取用户所创建的联系人组id
    String strSql;
    strSql.format (_T("select f_group_id from %s.t_person_group where f_user_id = '%s';"), dbName.getCStr(), escUserId.getCStr ());

    ncDBRecords results;
    dbOper->Select (strSql, results);

    if (results.size () == 0) {
        return;
    }

    vector<String> groupIds;
    for (size_t i = 0; i < results.size (); ++i) {
        groupIds.push_back (results[i][0]);
    }

    // 获取搜索配置信息
    int searchRange = 3;
    int searchResults = 3;
    String searchUserConfig = GetShareMgntConfig("search_user_config");
    JSON::Value searchUserConfigJson;
    JSON::Reader::read (searchUserConfigJson, searchUserConfig.getCStr (), searchUserConfig.getLength ());
    bool exact_search_user = searchUserConfigJson["exactSearch"].b ();
    searchRange = searchUserConfigJson["searchRange"].i ();
    searchResults = searchUserConfigJson["searchResults"].i ();

    String groupStr = GenerateGroupStr (groupIds);
    String escKey = dbOper->EscapeEx(key);
    String escLikeKey = dbOper->EscapeLikeClause(key);
    String searchStr = GenerateSearchStr (exact_search_user, searchRange, escKey, escLikeKey);

    // 连接查询t_contact_person和t_user表
    strSql.format (_T("select distinct t_user.f_user_id,t_user.f_login_name,t_user.f_display_name,t_user.f_mail_address,t_user.f_csf_level,t_user.f_priority ")
        _T("from %s.t_user, %s.t_contact_person ")
        _T("where ")
        _T("t_contact_person.f_group_id in (%s) and ")
        _T("t_user.f_status = '0' and t_user.f_auto_disable_status = '0' and ")
        _T("t_user.f_user_id = t_contact_person.f_user_id and %s ")
        _T("limit %d,%d;"),
        dbName.getCStr(), dbName.getCStr(),
        groupStr.getCStr (), searchStr.getCStr (), start, limit);

    dbOper->Select (strSql, results);

    ncACSUserInfo userInfo;
    for (size_t i = 0; i < results.size (); ++i) {
        userInfo.id = std::move(results[i][0]);
        userInfo.account = (searchResults == SHOW_LOGIN_AND_DISPLAY) ? std::move(results[i][1]) : "";
        userInfo.visionName = std::move(results[i][2]);
        userInfo.email = std::move(results[i][3]);
        userInfo.csfLevel = Int::getValue(results[i][4]);

        userInfos.push_back (userInfo);
    }

    bool only_share_to_user = Int::getValue(GetShareMgntConfig("only_share_to_user")) == 1;

    // 模糊搜索且可以共享给联系人组时, 连接查询t_person_group表
    if (!exact_search_user && !only_share_to_user) {
        if (userInfos.size () < limit) {
            limit -= (int)userInfos.size ();
            start = (userInfos.size () > 0) ? 0 : (start - SearchContactGroupCount(userId, key, false));

            // 搜索联系人组
            strSql.format(_T("select f_group_id,f_user_id,f_group_name,f_person_count from %s.t_person_group ")
                          _T("where f_user_id = '%s' ")
                          _T("and f_group_name like '%%%s%%' ")
                          _T("order by case when f_group_name = '%s' then 0 when f_group_name like '%%%s%%' then 1 else 2 end, upper(f_group_name) ")
                          _T("limit %d,%d;"),
                          dbName.getCStr(),
                          escUserId.getCStr(), escLikeKey.getCStr(), escKey.getCStr(), escLikeKey.getCStr(), start, limit);
            dbOper->Select (strSql, results);

            ncGroupInfo info;
            for (int64 i = 0; i < results.size(); i++) {
                info.id = std::move(results[i][0]);
                info.createrId = std::move(results[i][1]);
                info.groupName = std::move(results[i][2]);
                info.count = Int::getValue(results[i][3]);

                groupInfos.push_back (info);
            }
        }
    }

    NC_ACS_SHAREMGNT_TRACE (_T("end this: %p"), this);
}

/* [notxpcom] int SearchContactGroupCount ([const] in StringRef userId, [const] in StringRef key, in bool searchGroup); */
NS_IMETHODIMP_(int) ncACSShareMgnt::SearchContactGroupCount(const String & userId, const String & key, bool searchGroup)
{
    NC_ACS_SHAREMGNT_TRACE (_T("begin this: %p"), this);

    nsCOMPtr<ncIDBOperator> dbOper = getter_AddRefs (GetDBOperator ());
    String escUserId = dbOper->EscapeEx(userId);

    String dbName = Util::getDBName("sharemgnt_db");
    // 获取用户所创建的联系人组id
    String strSql;
    strSql.format (_T("select f_group_id from %s.t_person_group where f_user_id = '%s';"), dbName.getCStr(), escUserId.getCStr ());

    ncDBRecords results;
    dbOper->Select (strSql, results);

    if (results.size () == 0) {
        return 0;
    }

    vector<String> groupIds;
    for (size_t i = 0; i < results.size (); ++i) {
        groupIds.push_back (results[i][0]);
    }

    // 获取搜索配置信息
    int searchRange = 3;
    int searchResults = 3;
    String searchUserConfig = GetShareMgntConfig("search_user_config");
    JSON::Value searchUserConfigJson;
    JSON::Reader::read (searchUserConfigJson, searchUserConfig.getCStr (), searchUserConfig.getLength ());
    bool exact_search_user = searchUserConfigJson["exactSearch"].b ();
    searchRange = searchUserConfigJson["searchRange"].i ();
    searchResults = searchUserConfigJson["searchResults"].i ();

    String groupStr = GenerateGroupStr (groupIds);
    String escKey = dbOper->EscapeEx(key);
    String escLikeKey = dbOper->EscapeLikeClause(key);
    String searchStr = GenerateSearchStr (exact_search_user, searchRange, escKey, escLikeKey);

    // 连接查询t_contact_person和t_user表
    strSql.format (_T("select count(distinct t_user.f_user_id) ")
        _T("from %s.t_user, %s.t_contact_person ")
        _T("where ")
        _T("t_contact_person.f_group_id in (%s) and ")
        _T("t_user.f_status = '0' and t_user.f_auto_disable_status = '0' and ")
        _T("t_user.f_user_id = t_contact_person.f_user_id and %s;"),
        dbName.getCStr(), dbName.getCStr(),
        groupStr.getCStr (), searchStr.getCStr ());

    dbOper->Select (strSql, results);

    int count = Int::getValue(results[0][0]);

    // 模糊搜索时, 查询t_person_group表
    if (!exact_search_user && searchGroup) {
        // 搜索联系人组
        strSql.format (_T("select count(f_group_id) from %s.t_person_group ")
                        _T("where f_user_id = '%s' ")
                        _T("and f_group_name like '%%%s%%';"),
                        dbName.getCStr(), escUserId.getCStr (), escLikeKey.getCStr ());
        dbOper->Select (strSql, results);

        count += Int::getValue(results[0][0]);

    }

    NC_ACS_SHAREMGNT_TRACE (_T("end this: %p"), this);

    return count;
}

NS_IMETHODIMP_(void) ncACSShareMgnt::GetManageDepIds (const String& userId, vector<String>& departIds)
{
    NC_ACS_SHAREMGNT_TRACE (_T("begin this: %p"), this);

    GetAllManageDepIds(userId, departIds);

    NC_ACS_SHAREMGNT_TRACE (_T("end this: %p"), this);
}

/* [notxpcom] void GetAuditSupervisoryUserIds ([const] in StringRef userId, in StringVecRef userIds); */
NS_IMETHODIMP_(void) ncACSShareMgnt::GetAuditSupervisoryUserIds (const String& userId, vector<String>& userIds)
{
    userIds.clear ();
    vector<String> departIds;
    nsCOMPtr<ncIDBOperator> dbOper = getter_AddRefs (GetDBOperator ());

    String strSql;
    ncDBRecords results;

    String dbName = Util::getDBName("sharemgnt_db");
    // 获取用户管理的部门id
    strSql.format (_T("select f_department_id from %s.t_department_audit_person where f_user_id = '%s'"),
                        dbName.getCStr(), dbOper->EscapeEx(userId).getCStr ());
    dbOper->Select (strSql, results);

    vector<String> tmpIds;
    for (size_t i = 0; i < results.size (); ++i) {
        tmpIds.push_back (results[i][0]);
    }

    while (tmpIds.size () > 0) {
        for (size_t i = 0; i < tmpIds.size (); ++i) {
            departIds.push_back (tmpIds[i]);
        }

        String groupStr = GenerateGroupStr (tmpIds);

        strSql.format (_T("select f_department_id from %s.t_department_relation where f_parent_department_id in (%s);"),
                            dbName.getCStr(), groupStr.getCStr ());

        ncDBRecords results;
        dbOper->Select (strSql, results);

        tmpIds.clear ();
        for (size_t i = 0; i < results.size(); i++) {
            tmpIds.push_back(results[i][0]);
        }
    }
    if (departIds.size() == 0) {
        return;
    }
    // 去除掉重复id
    RemoveDuplicateStrs (departIds);

    String groupStr = GenerateGroupStr (departIds);

    // 连接查询t_user表和t_user_department_relation
    strSql.format (_T("select distinct u.f_user_id ")
        _T("from %s.t_user as u inner join %s.t_user_department_relation as r on u.f_user_id = r.f_user_id ")
        _T("where r.f_department_id in (%s) "),
        dbName.getCStr(), dbName.getCStr(),
        groupStr.getCStr ());
    dbOper->Select (strSql, results);

    for (size_t i = 0; i < results.size (); ++i) {
        userIds.push_back (results[i][0]);
    }
    NC_ACS_SHAREMGNT_TRACE (_T("this: %p"), this);
}

void ncACSShareMgnt::GetAllManageDepIds (const String& userId, vector<String>& departIds)
{
    departIds.clear ();
    nsCOMPtr<ncIDBOperator> dbOper = getter_AddRefs (GetDBOperator ());

    String strSql;
    ncDBRecords results;

    String dbName = Util::getDBName("sharemgnt_db");
    // 获取用户管理的部门id
    strSql.format (_T("select f_department_id from %s.t_department_responsible_person where f_user_id = '%s'"),
                        dbName.getCStr(), dbOper->EscapeEx(userId).getCStr ());
    dbOper->Select (strSql, results);

    vector<String> tmpIds;
    for (size_t i = 0; i < results.size (); ++i) {
        tmpIds.push_back (results[i][0]);
    }

    while (tmpIds.size () > 0) {
        for (size_t i = 0; i < tmpIds.size (); ++i) {
            departIds.push_back (tmpIds[i]);
        }

        String groupStr = GenerateGroupStr (tmpIds);

        strSql.format (_T("select f_department_id from %s.t_department_relation where f_parent_department_id in (%s);"),
                            dbName.getCStr(), groupStr.getCStr ());

        ncDBRecords results;
        dbOper->Select (strSql, results);

        tmpIds.clear ();
        for (size_t i = 0; i < results.size(); i++) {
            tmpIds.push_back(results[i][0]);
        }
    }

    // 去除掉重复id
    RemoveDuplicateStrs (departIds);

    NC_ACS_SHAREMGNT_TRACE (_T("this: %p"), this);
}

/* [notxpcom] void GetUserPermScopeInfos ([const] in StringRef userId, in ncPermScopeObjInfoVectorRef scopeInfos); */
void ncACSShareMgnt::GetUserPermScopeInfos (const String& userId, map<String, ncPermScopeObjInfo>& scopeInfos)
{
    NC_ACS_SHAREMGNT_TRACE (_T("begin this: %p"), this);

    scopeInfos.clear();

    String strSql;
    ncDBRecords results;
    nsCOMPtr<ncIDBOperator> dbOper = getter_AddRefs (GetDBOperator ());

    vector<String> strategyIds;
    String groupStr;

    String dbName = Util::getDBName("sharemgnt_db");
    // 获取共享者为用户的权限策略ID
    strSql.format(_T("select distinct f_strategy_id from %s.t_perm_share_strategy ")
                  _T("where f_obj_id = '%s' and f_sharer_or_scope = 1;"),
                    dbName.getCStr(), dbOper->EscapeEx(userId).getCStr ());
    dbOper->Select (strSql, results);

    for (size_t i = 0; i < results.size(); i++) {
        strategyIds.push_back(results[i][0]);
    }

    // 共享者：用户的父部门
    if (strategyIds.empty()) {
        // 获取用户的所有父部门ID
        vector<String> parentIds;
        GetAllBelongDepartmentIds(userId, parentIds);
        groupStr = GenerateGroupStr (parentIds);

        // 如果用户属于未分配组，返回空
        if (parentIds.size() == 0) {
            return;
        }

        // 获取共享者ID
        vector<String> sharerIds;
        strSql.format(_T("select distinct f_obj_id from %s.t_perm_share_strategy ")
                      _T("where f_obj_id in (%s) and f_sharer_or_scope = 1;"),
                        dbName.getCStr(), groupStr.getCStr ());
        dbOper->Select (strSql, results);

        for (size_t i = 0; i < results.size(); i++) {
            sharerIds.push_back(results[i][0]);
        }

        vector<String> parentPaths;
        vector<String>::iterator iter;

        for (iter = sharerIds.begin(); iter != sharerIds.end(); iter++) {
            vector<String> tmpPath;
            // 获取部门的跟组织路径
            GetDeptRootPath(*iter, tmpPath);
            // 跟组织路径中删除自己
            tmpPath.pop_back ();

            parentPaths.insert(parentPaths.end(), tmpPath.begin(), tmpPath.end());
        }

        // 跟组织路径上的部门，过滤掉
        for (iter = sharerIds.begin(); iter != sharerIds.end();) {
            if (find(parentPaths.begin(), parentPaths.end(), *iter) != parentPaths.end()) {
                sharerIds.erase(iter);
            }
            else {
                ++iter;
            }
        }

        // 获取相关的权限范围策略id
        if (sharerIds.size() > 0) {
            groupStr = GenerateGroupStr (sharerIds);
            strSql.format(_T("select distinct f_strategy_id from %s.t_perm_share_strategy ")
                          _T("where f_obj_id in (%s) and f_sharer_or_scope = 1;"),
                            dbName.getCStr(), groupStr.getCStr ());
            dbOper->Select (strSql, results);

            for (size_t i = 0; i < results.size (); i++) {
                strategyIds.push_back(results[i][0]);
            }
        }
    }

    // 获取用户的所有的权限范围对象
    if (strategyIds.size() > 0) {
        groupStr = GenerateGroupStr (strategyIds);

        strSql.format(_T("select f_obj_id, f_parent_id, f_obj_type from %s.t_perm_share_strategy where f_strategy_id in (%s) and f_sharer_or_scope = 2;"),
                        dbName.getCStr(), groupStr.getCStr ());

        dbOper->Select (strSql, results);

        for (size_t i = 0; i < results.size(); i++) {

            ncPermScopeObjInfo objInfo;
            objInfo.id = results[i][0];
            objInfo.parentId = results[i][1];
            if ( Int::getValue(results[i][2]) == 1){
                objInfo.type = IOC_USER;
            } else {
                objInfo.type = IOC_DEPARTMENT;
            }

            // 如果范围对象中用户和父部门关系不存在，则不添加
            if (objInfo.type == IOC_USER) {
                vector<String> directDeptIds;
                GetDirectBelongDepartmentIds(objInfo.id, directDeptIds);
                vector<String>::iterator deptIter = find(directDeptIds.begin(), directDeptIds.end(), objInfo.parentId);
                if (deptIter == directDeptIds.end())
                    continue;

                scopeInfos.insert(pair<String, ncPermScopeObjInfo> (objInfo.id, objInfo));
            }

            scopeInfos.insert(pair<String, ncPermScopeObjInfo> (objInfo.id + objInfo.parentId, objInfo));
        }
    }

    if (GetDefaulStrategySuperimStatus() || !scopeInfos.size ()) {
        // 如果启用所有用户的直属部门的权限范围，则添加直属部门ID到返现范围对象中
        strSql.format(_T("select f_status from %s.t_perm_share_strategy where f_strategy_id = '-1' and f_sharer_or_scope = 1;"), dbName.getCStr());
        dbOper->Select (strSql, results);

        if (results.size () == 1 && Int::getValue(results[0][0]) == 1) {
            vector<String> directDeptIds;
            GetDirectBelongDepartmentIds(userId, directDeptIds);

            for (size_t i = 0; i < directDeptIds.size(); i++){
                ncPermScopeObjInfo objInfo;
                objInfo.id = directDeptIds[i];
                objInfo.type = IOC_DEPARTMENT;
                scopeInfos.insert(pair<String, ncPermScopeObjInfo> (objInfo.id, objInfo));
            }
        }

        // 如果启用所有用户的直属组织的权限范围，则添加直属组织ID到返现范围对象中
        strSql.format(_T("select f_status from %s.t_perm_share_strategy where f_strategy_id = '-2' and f_sharer_or_scope = 1;"), dbName.getCStr());
        dbOper->Select (strSql, results);

        if (results.size () == 1 && Int::getValue(results[0][0]) == 1) {
            vector<String> directOrganIds;
            GetDirectBelongOrganIds(userId, directOrganIds);

            for (size_t i = 0; i < directOrganIds.size(); i++){
                ncPermScopeObjInfo objInfo;
                objInfo.id = directOrganIds[i];
                objInfo.type = IOC_DEPARTMENT;
                scopeInfos.insert(pair<String, ncPermScopeObjInfo> (objInfo.id, objInfo));
            }
        }
    }

    NC_ACS_SHAREMGNT_TRACE (_T("end this: %p"), this);
}

/* [notxpcom] void SearchDepartment ([const] in StringRef userId, [const] in StringRef key, in ncACSUserInfoVecRef userInfos); */
NS_IMETHODIMP_(void) ncACSShareMgnt::SearchDepartment(const String & userId, const String & key, vector<ncACSUserInfo> & userInfos)
{
    NC_ACS_SHAREMGNT_TRACE (_T("begin this: %p"), this);

    userInfos.clear ();
    if (key == String::EMPTY) {
        return;
    }

    nsCOMPtr<ncIDBOperator> dbOper = getter_AddRefs (GetDBOperator ());

    String escKey = dbOper->EscapeLikeClause(key);

    String dbName = Util::getDBName("sharemgnt_db");
    String strSql;
    ncDBRecords results;
    if (this->IsAdminId(userId)) {
        // admin搜索所有用户
        strSql.format(_T("select f_user_id,f_login_name,f_display_name,f_mail_address ")
                      _T("from %s.t_user ")
                      _T("where ")
                      _T("(f_login_name like '%%%s%%' or f_display_name like '%%%s%%') and ")
                      _T("f_user_id not in ('%s', '%s', '%s', '%s') ")
                      _T("order by t_user.f_priority, upper(t_user.f_display_name) ")
                      _T("limit 0,10;"),
                      dbName.getCStr(),
                      escKey.getCStr(), escKey.getCStr(),
                      g_ShareMgnt_constants.NCT_USER_ADMIN.c_str(), g_ShareMgnt_constants.NCT_USER_AUDIT.c_str(),
                      g_ShareMgnt_constants.NCT_USER_SYSTEM.c_str(), g_ShareMgnt_constants.NCT_USER_SECURIT.c_str());
    }
    else {
        // 用户组织管理员向下遍历搜索用户
        vector<String> departIds;
        GetAllManageDepIds (userId, departIds);

        if (departIds.size () == 0) {
            return;
        }

        String groupStr = GenerateGroupStr (departIds);

        // 连接查询t_user表和t_user_department_relation
        strSql.format(_T("select distinct t_user.f_user_id,t_user.f_login_name,t_user.f_display_name,t_user.f_mail_address, t_user.f_priority ")
                      _T("from %s.t_user, %s.t_user_department_relation ")
                      _T("where ")
                      _T("t_user_department_relation.f_department_id in (%s) and ")
                      _T("t_user.f_user_id = t_user_department_relation.f_user_id and ")
                      _T("(t_user.f_login_name like '%%%s%%' or t_user.f_display_name like '%%%s%%') and ")
                      _T("t_user.f_user_id not in ('%s', '%s', '%s', '%s') ")
                      _T("order by t_user.f_priority, upper(t_user.f_display_name) ")
                      _T("limit 0,10;"),
                      dbName.getCStr(), dbName.getCStr(),
                      groupStr.getCStr(), escKey.getCStr(), escKey.getCStr(),
                      g_ShareMgnt_constants.NCT_USER_AUDIT.c_str(), g_ShareMgnt_constants.NCT_USER_AUDIT.c_str(),
                      g_ShareMgnt_constants.NCT_USER_SYSTEM.c_str(), g_ShareMgnt_constants.NCT_USER_SECURIT.c_str());
    }

    dbOper->Select (strSql, results);

    ncACSUserInfo userInfo;
    for (size_t i = 0; i < results.size (); ++i) {
        userInfo.id = std::move(results[i][0]);
        userInfo.account = std::move(results[i][1]);
        userInfo.visionName = std::move(results[i][2]);
        userInfo.email = std::move(results[i][3]);

        userInfos.push_back (userInfo);
    }

    NC_ACS_SHAREMGNT_TRACE (_T("end this: %p"), this);
}

/* [notxpcom] void GetBelongDepartByIdBatch ([const] in StringVecRef userIds, in ncACSUserIdDepartMapRef infoMap); */
NS_IMETHODIMP_(void) ncACSShareMgnt::GetBelongDepartByIdBatch(const vector<String> & userIds, map<String, ncACSDepartInfo> & infoMap)
{
    NC_ACS_SHAREMGNT_TRACE (_T("begin this: %p"), this);

    infoMap.clear ();
    nsCOMPtr<ncIDBOperator> dbOper = getter_AddRefs (GetDBOperator ());

    if (userIds.size () == 0) {
        return;
    }

    String groupStr = GenerateGroupStr (userIds);

    String dbName = Util::getDBName("sharemgnt_db");
    String strSql;
    strSql.format (_T("select t_user_department_relation.f_user_id, t_user_department_relation.f_department_id, t_department.f_name ")
                    _T("from %s.t_user_department_relation, %s.t_department ")
                    _T("where ")
                    _T("f_user_id in (%s) and ")
                    _T("t_user_department_relation.f_department_id = t_department.f_department_id ")
                    _T("order by t_user_department_relation.f_relation_id;"),
                        dbName.getCStr(), dbName.getCStr(),
                        groupStr.getCStr ());

    ncDBRecords results;
    dbOper->Select (strSql, results);

    ncACSDepartInfo depInfo;
    for (size_t i = 0; i < results.size (); ++i) {
        depInfo.id = std::move(results[i][1]);
        depInfo.name = std::move(results[i][2]);

        infoMap.insert (pair<String, ncACSDepartInfo> (results[i][0], depInfo));
    }

    NC_ACS_SHAREMGNT_TRACE (_T("end this: %p"), this);
}

/* [notxpcom] void GetBelongGroupByIdBatch ([const] in StringRef createrId, [const] in StringVecRef userIds, in ncGroupInfoMapRef infoMap); */
NS_IMETHODIMP_(void) ncACSShareMgnt::GetBelongGroupByIdBatch(const String & createrId, const vector<String> & userIds, map<String, ncGroupInfo> & infoMap)
{
    NC_ACS_SHAREMGNT_TRACE (_T("begin this: %p"), this);

    infoMap.clear ();
    nsCOMPtr<ncIDBOperator> dbOper = getter_AddRefs (GetDBOperator ());

    if (userIds.size () == 0) {
        return;
    }

    String groupStr = GenerateGroupStr (userIds);

    String escCreaterId;
    dbOper->Escape (createrId, escCreaterId);

    String dbName = Util::getDBName("sharemgnt_db");
    String strSql;
    strSql.format (_T("select t_contact_person.f_user_id,t_person_group.f_group_id,t_person_group.f_group_name,")
                    _T("t_person_group.f_user_id,t_person_group.f_person_count ")
                    _T("from %s.t_contact_person, %s.t_person_group ")
                    _T("where ")
                    _T("t_person_group.f_user_id = '%s' and ")
                    _T("t_contact_person.f_group_id = t_person_group.f_group_id and ")
                    _T("t_contact_person.f_user_id in (%s) ")
                    _T("order by t_contact_person.f_id;"),
                        dbName.getCStr(), dbName.getCStr(),
                        escCreaterId.getCStr (), groupStr.getCStr ());

    ncDBRecords results;
    dbOper->Select (strSql, results);

    ncGroupInfo groupInfo;
    for (size_t i = 0; i < results.size (); ++i) {
        groupInfo.id = std::move(results[i][1]);
        groupInfo.groupName = std::move(results[i][2]);
        groupInfo.createrId = std::move(results[i][3]);
        groupInfo.count = Int64::getValue (results[i][4]);;

        // 这里使用insert，属于多个联系人组的话，后面的就不会插入了
        infoMap.insert (pair<String, ncGroupInfo>(std::move(results[i][0]), groupInfo));
    }

    NC_ACS_SHAREMGNT_TRACE (_T("end this: %p"), this);
}

/* [notxpcom] String ExtLogin ([const] in StringRef appId, [const] in StringRef account, [const] in StringRef key, [const] in StringVecRef params, in intRef accountType); */
NS_IMETHODIMP_(String) ncACSShareMgnt::ExtLogin(const String & appId, const String & account, const String & key, const vector<String> & params, int& accountType)
{
    NC_ACS_SHAREMGNT_TRACE (_T("begin this: %p"), this);

    accountType = 0;

    nsCOMPtr<ncIDBOperator> dbOper = getter_AddRefs (GetDBOperator ());

    String escAccount = dbOper->EscapeEx(account);
    String escAppId = dbOper->EscapeEx(appId);

    String dbName = Util::getDBName("sharemgnt_db");
    String strSql;
    strSql.format (_T("select f_user_id, f_status from %s.t_user where f_login_name = '%s'"),
                dbName.getCStr(), escAccount.getCStr ());

    ncDBRecords results;
    dbOper->Select (strSql, results);

    if (results.size () == 0){
        String strSql_idcard;
        ncDBRecords results_idcard;
        strSql_idcard.format (_T("select f_value from %s.t_sharemgnt_config where f_key = 'id_card_login_status';"), dbName.getCStr());
        dbOper->Select (strSql_idcard, results_idcard);

        if (results_idcard.size() == 1 && Int::getValue (results_idcard[0][0]) == 1){
            string des_account = DESEncrypt(toSTLString(account) );
            string base64_account = Base64Encode(des_account);
            strSql.format (_T("select f_user_id, f_status from %s.t_user where f_idcard_number = '%s' limit 1"),
                            dbName.getCStr(), dbOper->EscapeEx(toCFLString(base64_account)).getCStr ());
            dbOper->Select (strSql, results);
            if (results.size () == 1) {
                accountType = 1;
            }
        }
    }

    if (results.size () == 0) {
        return String::EMPTY;
    }

    String userId = results[0][0];
    // 管理员账号不允许登录
    if (_adminIds.count(userId.getCStr()) != 0) {
        return String::EMPTY;
    }

    int status = Int::getValue(results[0][1]);
    if (status != 0) {
        THROW_E (ACS_SHAREMGNT, SER_DISABLED, "User is disabled.");
    }

    // 检查appId是否存在
    strSql.format (_T("select f_app_key from %s.t_third_auth_info where f_app_id = '%s' and f_enabled = 1;"),
        dbName.getCStr(), escAppId.getCStr ());
    dbOper->Select (strSql, results);
    if(results.size () == 0) {
        return String::EMPTY;
    }

    String appKey = results[0][0];

    //中国人民公安大学
    /*
    参数：params[0]:time 时间戳字符串，格式:‘"%Y-%m-%d %H:%M:%S"’，‘2015-10-29 15:04:32’
    */
    if (appId == _T("ppsuc")) {
        String time = params[0];
        String comStr = account + appKey + time;
        String md5Str = genMD5String2 (comStr);
        if (md5Str.compareIgnoreCase (key) == 0) {
            return userId;
        }
    }

    //标准认证
    else {
        String comStr = appId + appKey + account;
        String md5Str = genMD5String2 (comStr);
        if (md5Str.compareIgnoreCase (key) == 0) {
            return userId;
        }
    }

    return String::EMPTY;
}

/* [notxpcom] int GetUserCSFLevel ([const] in StringRef userId, [const] in ncVisitorType visitorType); */
NS_IMETHODIMP_(int) ncACSShareMgnt::GetUserCSFLevel(const String & userId, const ncVisitorType visitorType)
{
    NC_ACS_SHAREMGNT_TRACE (_T("[BEGIN]this: %p"), this);

    // 匿名用户或内外网数据交换用户返回最高密级
    if (ncVisitorType::ANONYMOUS == visitorType || ncVisitorType::BUSINESS == visitorType || userId.compare (NC_EVFS_NAME_IOC_DATAEXCHANGE_ID) == 0) {
        return 0x7FFF;
    }

    nsCOMPtr<ncIDBOperator> dbOper = getter_AddRefs (GetDBOperator ());
    String dbName = Util::getDBName("sharemgnt_db");
    String strSql;
    strSql.format (_T("select f_csf_level from %s.t_user where f_user_id = '%s'"),
                       dbName.getCStr(), dbOper->EscapeEx(userId).getCStr ());

    ncDBRecords results;
    dbOper->Select (strSql, results);

    NC_ACS_SHAREMGNT_TRACE (_T("[END]this: %p"), this);

    if (results.size () == 1) {
        return Int::getValue(results[0][0]);
    }
    else {
        THROW_E (ACS_SHAREMGNT, USER_ID_NOT_EXISTS, "User id doesn't exist.");
    }
}

/* [notxpcom] int GetUserCSFLevel2 ([const] in StringRef userId, [const] in ncVisitorType visitorType); */
NS_IMETHODIMP_(int) ncACSShareMgnt::GetUserCSFLevel2(const String & userId, const ncVisitorType visitorType)
{
    NC_ACS_SHAREMGNT_TRACE (_T("[BEGIN]this: %p"), this);

    // 匿名用户或内外网数据交换用户返回最高密级
    if (ncVisitorType::ANONYMOUS == visitorType || ncVisitorType::BUSINESS == visitorType || userId.compare (NC_EVFS_NAME_IOC_DATAEXCHANGE_ID) == 0) {
        return 0x7FFF;
    }

    nsCOMPtr<ncIDBOperator> dbOper = getter_AddRefs (GetDBOperator ());
    String dbName = Util::getDBName("sharemgnt_db");
    String strSql;
    strSql.format (_T("select f_csf_level2 from %s.t_user where f_user_id = '%s'"),
                       dbName.getCStr(), dbOper->EscapeEx(userId).getCStr ());

    ncDBRecords results;
    dbOper->Select (strSql, results);

    NC_ACS_SHAREMGNT_TRACE (_T("[END]this: %p"), this);

    if (results.size () == 1) {
        return Int::getValue(results[0][0]);
    }
    else {
        THROW_E (ACS_SHAREMGNT, USER_ID_NOT_EXISTS, "User id doesn't exist.");
    }
}

/* [notxpcom] bool GetPermShareLimitStatus ( ); */
NS_IMETHODIMP_(bool) ncACSShareMgnt::GetPermShareLimitStatus()
{
    NC_ACS_SHAREMGNT_TRACE (_T("this: %p"), this);
    String dbName = Util::getDBName("sharemgnt_db");
    String strSql;
    strSql.format (_T("select f_value from %s.t_sharemgnt_config where f_key = 'perm_share_status';"), dbName.getCStr());

    nsCOMPtr<ncIDBOperator> dbOper = getter_AddRefs (GetDBOperator ());
    ncDBRecords results;
    dbOper->Select (strSql, results);

    bool status = true;
    if (results.size() == 1 && Int::getValue (results[0][0]) == 0)
        status = false;

    NC_ACS_SHAREMGNT_TRACE (_T("end this: %p"), this);

    return status;
}

/* [notxpcom] bool GetFindShareLimitStatus ( ); */
NS_IMETHODIMP_(bool) ncACSShareMgnt::GetFindShareLimitStatus()
{
    NC_ACS_SHAREMGNT_TRACE (_T("this: %p"), this);
    String dbName = Util::getDBName("sharemgnt_db");
    String strSql;
    strSql.format (_T("select f_value from %s.t_sharemgnt_config where f_key = 'find_share_status';"), dbName.getCStr());

    nsCOMPtr<ncIDBOperator> dbOper = getter_AddRefs (GetDBOperator ());
    ncDBRecords results;
    dbOper->Select (strSql, results);

    bool status = true;
    if (results.size() == 1 && Int::getValue (results[0][0]) == 0)
        status = false;

    NC_ACS_SHAREMGNT_TRACE (_T("end this: %p"), this);

    return status;
}

/* [notxpcom] bool GetLinkShareLimitStatus ( ); */
NS_IMETHODIMP_(bool) ncACSShareMgnt::GetLinkShareLimitStatus()
{
    NC_ACS_SHAREMGNT_TRACE (_T("this: %p"), this);
    String dbName = Util::getDBName("sharemgnt_db");
    String strSql;
    strSql.format (_T("select f_value from %s.t_sharemgnt_config where f_key = 'link_share_status';"), dbName.getCStr());

    nsCOMPtr<ncIDBOperator> dbOper = getter_AddRefs (GetDBOperator ());
    ncDBRecords results;
    dbOper->Select (strSql, results);

    bool status = true;
    if (results.size() == 1 && Int::getValue (results[0][0]) == 0)
        status = false;

    NC_ACS_SHAREMGNT_TRACE (_T("end this: %p"), this);

    return status;
}

/* [notxpcom] bool IsUserFindEnabled ([const] in StringRef userId); */
NS_IMETHODIMP_(bool) ncACSShareMgnt::IsUserFindEnabled(const String& userId)
{
    NC_ACS_SHAREMGNT_TRACE (_T("this: %p"), this);

    String dbName = Util::getDBName("sharemgnt_db");
    bool findEnable = false;
    if (!GetFindShareLimitStatus()) {
        findEnable = true;
    }
    else {
        // 用户
        vector<String> pathIds;
        GetAllBelongDepartmentIds(userId, pathIds);
        pathIds.push_back(userId);

        String strSql;
        strSql.format (_T("select distinct f_sharer_id from %s.t_find_share_strategy;"), dbName.getCStr());

        nsCOMPtr<ncIDBOperator> dbOper = getter_AddRefs (GetDBOperator ());
        ncDBRecords results;
        dbOper->Select (strSql, results);

        set<String> sharerIds;
        for (size_t i = 0; i < results.size(); ++i) {
            sharerIds.insert(results[i][0]);
        }

        for (size_t i = 0; i < pathIds.size(); ++i) {
            if (sharerIds.find(pathIds[i]) != sharerIds.end()) {
                findEnable = true;
                break;
            }
        }
    }

    NC_ACS_SHAREMGNT_TRACE (_T("end this: %p"), this);

    return findEnable;
}

/* [notxpcom] bool IsUserLinkEnabled ([const] in StringRef userId); */
NS_IMETHODIMP_(bool) ncACSShareMgnt::IsUserLinkEnabled(const String& userId)
{
    NC_ACS_SHAREMGNT_TRACE (_T("this: %p"), this);

    bool linkEnable = false;

    String dbName = Util::getDBName("sharemgnt_db");
    if (!GetLinkShareLimitStatus()) {
        linkEnable = true;
    }
    else {
        // 用户
        vector<String> pathIds;
        GetAllBelongDepartmentIds(userId, pathIds);
        pathIds.push_back(userId);

        String strSql;
        strSql.format (_T("select distinct f_sharer_id from %s.t_link_share_strategy;"), dbName.getCStr());

        nsCOMPtr<ncIDBOperator> dbOper = getter_AddRefs (GetDBOperator ());
        ncDBRecords results;
        dbOper->Select (strSql, results);

        set<String> sharerIds;
        for (size_t i = 0; i < results.size(); ++i) {
            sharerIds.insert(results[i][0]);
        }

        for (size_t i = 0; i < pathIds.size(); ++i) {
            if (sharerIds.find(pathIds[i]) != sharerIds.end()) {
                linkEnable = true;
                break;
            }
        }
    }

    NC_ACS_SHAREMGNT_TRACE (_T("end this: %p"), this);

    return linkEnable;
}

/* [notxpcom] [notxpcom] void GetAllUserContacts (in StringMapVecRef userContactIds);*/
NS_IMETHODIMP_(void) ncACSShareMgnt::GetAllUserContactIds(map<String, vector<String> >& userContactIds)
{
    NC_ACS_SHAREMGNT_TRACE (_T("this: %p"), this);

    userContactIds.clear ();

    nsCOMPtr<ncIDBOperator> dbOper = getter_AddRefs (GetDBOperator ());

    String dbName = Util::getDBName("sharemgnt_db");
    String strSql;
    strSql.format (_T("select t_person_group.f_user_id as user_id, t_contact_person.f_user_id as contact_id ")
                   _T("from %s.t_person_group join %s.t_contact_person ")
                   _T("on t_person_group.f_group_id = t_contact_person.f_group_id;"), dbName.getCStr(), dbName.getCStr());

    ncDBRecords results;
    dbOper->Select (strSql, results);

    for (size_t i = 0; i < results.size (); ++i) {
        String& userId = results[i][0];
        map<String, vector<String> >::iterator iter = userContactIds.find(userId);
        if( iter != userContactIds.end()) {
            vector<String>& contactIds = iter->second;
            contactIds.push_back(results[i][1]);
        }
        else {
            vector<String> contactIds;
            contactIds.push_back(results[i][1]);
            userContactIds.insert(pair<String, vector<String> >(userId, contactIds));
        }
    }

    NC_ACS_SHAREMGNT_TRACE (_T("this: %p"), this);
}

/* [notxpcom] void DeleteContactsByPatch ([const] in StringRef userId, [const] in StringVecRef contactIds);*/
NS_IMETHODIMP_(void) ncACSShareMgnt::DeleteContactsByPatch(const String& userId, const vector<String>& contactIds)
{
    NC_ACS_SHAREMGNT_TRACE (_T("this: %p"), this);

    if (contactIds.empty())
        return;

    nsCOMPtr<ncIDBOperator> dbOper = getter_AddRefs (GetDBOperator ());

    String groupStr = GenerateGroupStr (contactIds);

    String dbName = Util::getDBName("sharemgnt_db");
    String strSql, updateSql;
    strSql.format (_T("delete from %s.t_contact_person where f_user_id in (%s) ")
                   _T("and f_group_id in (select f_group_id from %s.t_person_group where f_user_id = '%s');"),
                   dbName.getCStr(), groupStr.getCStr (),
                   dbName.getCStr(), dbOper->EscapeEx(userId).getCStr());

    ncDBRecords results;
    dbOper->Execute (strSql);

    // 更新联系人组中联系人个数
    strSql.format (_T("select count(f_user_id), f_group_id from %s.t_contact_person group by f_group_id;"), dbName.getCStr());
    dbOper->Select (strSql, results);

    for (size_t i = 0; i < results.size (); ++i) {
        updateSql.format (_T("update %s.t_person_group set f_person_count = %d where f_group_id = '%s'"),
                        dbName.getCStr(), Int64::getValue (results[i][0]), dbOper->EscapeEx(results[i][1]).getCStr ());

        dbOper->Execute (updateSql);
    }

    NC_ACS_SHAREMGNT_TRACE (_T("this: %p"), this);
}

/* [notxpcom] bool GetLeakProofStatus ( ); */
NS_IMETHODIMP_(bool) ncACSShareMgnt::GetLeakProofStatus()
{
    NC_ACS_SHAREMGNT_TRACE (_T("this: %p"), this);
    String dbName = Util::getDBName("sharemgnt_db");
    String strSql;
    strSql.format (_T("select f_value from %s.t_sharemgnt_config where f_key = 'leak_proof_status';"), dbName.getCStr());

    nsCOMPtr<ncIDBOperator> dbOper = getter_AddRefs (GetDBOperator ());
    ncDBRecords results;
    dbOper->Select (strSql, results);

    bool status = true;
    if (results.size() == 1 && Int::getValue (results[0][0]) == 0)
        status = false;

    NC_ACS_SHAREMGNT_TRACE (_T("end this: %p"), this);

    return status;
}

/* [notxpcom] int GetLeakProofPerm ([const] in StringRef userId);*/
NS_IMETHODIMP_(int) ncACSShareMgnt::GetLeakProofPerm(const String& userId)
{
    NC_ACS_SHAREMGNT_TRACE (_T("this: %p"), this);

    nsCOMPtr<ncIDBOperator> dbOper = getter_AddRefs (GetDBOperator ());

    String dbName = Util::getDBName("sharemgnt_db");
    if (!GetLeakProofStatus()) {
        NC_ACS_SHAREMGNT_TRACE (_T("end this: %p"), this);
        return 3;
    }
    else {
        int value = 0;
        String strSql;

        // 获取用户的防泄密策略
        strSql.format (_T("select f_perm_value from %s.t_leak_proof_strategy where f_accessor_id = '%s'"), dbName.getCStr(), dbOper->EscapeEx(userId).getCStr ());

        ncDBRecords results;
        dbOper->Select (strSql, results);

        if (results.size() == 1) {
            value = Int::getValue(results[0][0]);
        }

        if (!value) {
            // 获取用户的所有所属部门
            vector<String> accessorIds;
            GetAllBelongDepartmentIds(userId, accessorIds);

            if (accessorIds.size () > 0) {
                String groupStr = GenerateGroupStr(accessorIds);

                strSql.format (_T("select f_accessor_id, f_perm_value from %s.t_leak_proof_strategy where f_accessor_id in (%s)"), dbName.getCStr(), groupStr.getCStr());
                dbOper->Select (strSql, results);

                map<String, int> accessorPerms;
                for (size_t i = 0; i < results.size(); ++i) {
                    accessorPerms.insert (pair<String, int> (results[i][0], Int::getValue(results[i][1])));
                }

                vector<String> parentPaths;
                for (map<String, int>::iterator iter = accessorPerms.begin();iter != accessorPerms.end(); iter++) {
                    // 获取跟组织路径
                    vector<String> tmpPath;
                    GetDeptRootPath(iter->first, tmpPath);
                    // 跟组织路径中删除自己
                    tmpPath.pop_back ();

                    parentPaths.insert(parentPaths.end(), tmpPath.begin(), tmpPath.end());
                }

                // 跟组织路径上的部门，直接过滤掉
                for (map<String, int>::iterator iter = accessorPerms.begin(); iter != accessorPerms.end();) {
                    if (find(parentPaths.begin(), parentPaths.end(), iter->first) != parentPaths.end()) {
                        iter = accessorPerms.erase(iter);
                    }
                    else {
                        ++iter;
                    }
                }

                // 计算权限值
                if (accessorPerms.size() > 0) {
                    for(map<String, int>::iterator iter = accessorPerms.begin();iter != accessorPerms.end(); iter++) {
                        value |= iter->second;
                    }
                }
            }
        }

        NC_ACS_SHAREMGNT_TRACE (_T("end this: %p"), this);
        return value;
    }
}

/* [notxpcom] int GetClearCacheInterval();*/
NS_IMETHODIMP_(int) ncACSShareMgnt::GetClearCacheInterval( )
{
    NC_ACS_SHAREMGNT_TRACE (_T("this: %p"), this);

    String dbName = Util::getDBName("sharemgnt_db");
    String strSql;
    strSql.format (_T("select `f_value` FROM %s.t_sharemgnt_config WHERE `f_key` = 'clear_cache_interval'"), dbName.getCStr());

    nsCOMPtr<ncIDBOperator> dbOper = getter_AddRefs (GetDBOperator ());
    ncDBRecords results;
    dbOper->Select (strSql, results);

    int value = -1;
    if (results.size() == 1)
        value = Int::getValue(results[0][0]);

    NC_ACS_SHAREMGNT_TRACE (_T("end this: %p"), this);
    return value;
}

/* [notxpcom] int GetClearCacheInterval();*/
NS_IMETHODIMP_(int64) ncACSShareMgnt::GetClearCacheSize( )
{
    NC_ACS_SHAREMGNT_TRACE (_T("this: %p"), this);

    String dbName = Util::getDBName("sharemgnt_db");
    String strSql;
    strSql.format (_T("select `f_value` FROM %s.t_sharemgnt_config WHERE `f_key` = 'clear_cache_size'"), dbName.getCStr());

    nsCOMPtr<ncIDBOperator> dbOper = getter_AddRefs (GetDBOperator ());
    ncDBRecords results;
    dbOper->Select (strSql, results);

    int64 value = -1;
    if (results.size() == 1)
        value = Int64::getValue(results[0][0]);

    NC_ACS_SHAREMGNT_TRACE (_T("end this: %p"), this);
    return value;
}

/* [notxpcom] int GetPwdControl ([const] in StringRef userId); */
NS_IMETHODIMP_(int) ncACSShareMgnt::GetPwdControl(const String & userId)
{
    NC_ACS_SHAREMGNT_TRACE (_T("this: %p"), this);

    nsCOMPtr<ncIDBOperator> dbOper = getter_AddRefs (GetDBOperator ());

    String dbName = Util::getDBName("sharemgnt_db");
    String strSql;
    strSql.format (_T("select f_pwd_control from %s.t_user where f_user_id = '%s'"), dbName.getCStr(), dbOper->EscapeEx(userId).getCStr ());

    ncDBRecords results;
    dbOper->Select (strSql, results);

    NC_ACS_SHAREMGNT_TRACE (_T("this: %p"), this);

    if (results.size () == 1) {
        return Int::getValue(results[0][0]);
    }
    else {
        THROW_E (ACS_SHAREMGNT, USER_ID_NOT_EXISTS, "User id doesn't exist.");
    }
}

/* [notxpcom] int GetUserAuthType ([const] in StringRef userId); */
NS_IMETHODIMP_(int) ncACSShareMgnt::GetUserAuthType(const String & userId)
{
    NC_ACS_SHAREMGNT_TRACE (_T("this: %p"), this);

    nsCOMPtr<ncIDBOperator> dbOper = getter_AddRefs (GetDBOperator ());

    String dbName = Util::getDBName("sharemgnt_db");
    String strSql;
    strSql.format (_T("select f_auth_type from %s.t_user where f_user_id = '%s'"), dbName.getCStr(), dbOper->EscapeEx(userId).getCStr ());

    ncDBRecords results;
    dbOper->Select (strSql, results);

    NC_ACS_SHAREMGNT_TRACE (_T("this: %p"), this);

    if (results.size () == 1) {
        return Int::getValue(results[0][0]);
    }
    else {
        THROW_E (ACS_SHAREMGNT, USER_ID_NOT_EXISTS, "User id doesn't exist.");
    }
}

/* [notxpcom] bool GetHideClientCacheStatus ( ); */
NS_IMETHODIMP_(bool) ncACSShareMgnt::GetHideClientCacheStatus()
{
    NC_ACS_SHAREMGNT_TRACE (_T("this: %p"), this);
    String dbName = Util::getDBName("sharemgnt_db");
    String strSql;
    strSql.format (_T("select f_value from %s.t_sharemgnt_config where f_key = 'hide_client_cache_setting';"), dbName.getCStr());

    nsCOMPtr<ncIDBOperator> dbOper = getter_AddRefs (GetDBOperator ());
    ncDBRecords results;
    dbOper->Select (strSql, results);

    bool status = true;
    if (results.size() == 1 && Int::getValue (results[0][0]) == 0)
        status = false;

    NC_ACS_SHAREMGNT_TRACE (_T("end this: %p"), this);

    return status;
}

/* [notxpcom] bool GetClearClientCacheStatus ( ); */
NS_IMETHODIMP_(bool) ncACSShareMgnt::GetClearClientCacheStatus()
{
    NC_ACS_SHAREMGNT_TRACE (_T("this: %p"), this);

    String dbName = Util::getDBName("sharemgnt_db");
    String strSql;
    strSql.format (_T("select f_value from %s.t_sharemgnt_config where f_key = 'force_clear_client_cache';"), dbName.getCStr());

    nsCOMPtr<ncIDBOperator> dbOper = getter_AddRefs (GetDBOperator ());
    ncDBRecords results;
    dbOper->Select (strSql, results);

    bool status = true;
    if (results.size() == 1 && Int::getValue (results[0][0]) == 0)
        status = false;

    NC_ACS_SHAREMGNT_TRACE (_T("end this: %p"), this);

    return status;
}

/* [notxpcom] void GetMutiTenantStatus ();*/
NS_IMETHODIMP_(bool) ncACSShareMgnt::GetMutiTenantStatus()
{
    NC_ACS_SHAREMGNT_TRACE (_T("this: %p"), this);

    String dbName = Util::getDBName("sharemgnt_db");
    String strSql;
    strSql.format ( _T("SELECT `f_value` FROM %s.t_sharemgnt_config WHERE `f_key` = 'multi_tenant';"), dbName.getCStr());

    nsCOMPtr<ncIDBOperator> dbOper = getter_AddRefs (GetDBOperator ());
    ncDBRecords results;
    dbOper->Select (strSql, results);

    bool status = false;
    if (results.size() == 1 && Int::getValue (results[0][0]) == 1)
        status = true;

    NC_ACS_SHAREMGNT_TRACE (_T("this: %p"), this);
    return status;
}

/* [notxpcom] String GetShareMgntConfig ([const] in StringRef key);*/
NS_IMETHODIMP_(String) ncACSShareMgnt::GetShareMgntConfig(const String &key)
{
    NC_ACS_SHAREMGNT_TRACE (_T("this: %p"), this);

    nsCOMPtr<ncIDBOperator> dbOper = getter_AddRefs (GetDBOperator ());

    String dbName = Util::getDBName("sharemgnt_db");
    String strSql;
    strSql.format ( _T("SELECT `f_value` FROM %s.t_sharemgnt_config WHERE `f_key` = '%s';"), dbName.getCStr(), dbOper->EscapeEx(key).getCStr ());

    ncDBRecords results;
    dbOper->Select (strSql, results);

    if (results.size() != 1) {
        THROW_E (ACS_SHAREMGNT, KEY_CONFIG_NOT_EXISTS, "key config doesn't exist.");
    }

    NC_ACS_SHAREMGNT_TRACE (_T("this: %p"), this);
    return results[0][0];
}

/* [notxpcom] void BatchGetConfig (in VectorStringRef keys, in StringMapRef kvMap);*/
NS_IMETHODIMP_(void) ncACSShareMgnt::BatchGetConfig(vector<String>& keys, map<String, String>& kvMap)
{
    NC_ACS_SHAREMGNT_TRACE (_T("this: %p begin"), this);

    kvMap.clear();

    nsCOMPtr<ncIDBOperator> dbOper = getter_AddRefs (GetDBOperator ());

    String dbName = Util::getDBName("sharemgnt_db");
    String groupStr = GenerateGroupStr (keys);
    String strSql;
    strSql.format (_T("select f_key, f_value from %s.t_sharemgnt_config where f_key in (%s)"), dbName.getCStr(), groupStr.getCStr ());

    ncDBRecords results;
    dbOper->Select (strSql, results);

    for (size_t i = 0; i < results.size (); ++i) {
        kvMap[results[i][0]] = results[i][1];
    }

    NC_ACS_SHAREMGNT_TRACE (_T("this: %p end"), this);
}

/* [notxpcom] bool GetCustomConfigOfString ([const] in StringRef key, in StringRef value); */
NS_IMETHODIMP_(bool) ncACSShareMgnt::GetCustomConfigOfString(const String & key, String & value)
{
    NC_ACS_SHAREMGNT_TRACE (_T("this: %p"), this);

    nsCOMPtr<ncIDBOperator> dbOper = getter_AddRefs (GetDBOperator ());
    String dbName = Util::getDBName("sharemgnt_db");
    bool ret = false;
    String strSql;
    strSql.format (_T("select f_value from %s.t_sharemgnt_config where f_key = '%s'"), dbName.getCStr(), dbOper->EscapeEx(key).getCStr());
    ncDBRecords results;

    dbOper->Select (strSql, results);
    if (results.size () == 1) {
        value = std::move(results[0][0]);
        ret = true;
    }

    NC_ACS_SHAREMGNT_TRACE (_T("end this: %p"), this);
    return ret;
}

/* [notxpcom] bool GetNetDocsLimitStatus ();*/
NS_IMETHODIMP_(bool) ncACSShareMgnt::GetNetDocsLimitStatus()
{
    NC_ACS_SHAREMGNT_TRACE (_T("this: %p"), this);
    return (Int::getValue(GetShareMgntConfig("enable_net_docs_limit")) == 1);
}

/*[notxpcom] bool GetMessagePluginStatus ();*/
NS_IMETHODIMP_(bool) ncACSShareMgnt::GetMessagePluginStatus()
{
    NC_ACS_SHAREMGNT_TRACE (_T("this: %p"), this);
    nsCOMPtr<ncIDBOperator> dbOper = getter_AddRefs (GetDBOperator ());
    String strSql;
    ncDBRecords results;
    String dbName = Util::getDBName("sharemgnt_db");
    strSql.format (_T("SELECT `f_id` FROM %s.`t_third_party_auth` WHERE `f_plugin_type` = 1 AND `f_enable` = 1"), dbName.getCStr());
    dbOper->Select (strSql, results);

    if (results.size() != 0) {
        return true;
    }
    else {
        return false;
    }
}

/* [notxpcom] bool CheckNetDocLimit ([const] in StringRef cid, [const] in StringRef accessIp);*/
NS_IMETHODIMP_(bool) ncACSShareMgnt::CheckNetDocLimit(const String &cid, const String &accessIp)
{
    NC_ACS_SHAREMGNT_TRACE (_T("this: %p"), this);
    // 根据cid获取所有可以访问的ip网段
    nsCOMPtr<ncIDBOperator> dbOper = getter_AddRefs (GetDBOperator ());

    String dbName = Util::getDBName("sharemgnt_db");
    String strSql;
    ncDBRecords results;
    strSql.format (_T("SELECT DISTINCT `f_ip`, `f_sub_net_mask` FROM %s.t_net_docs_limit_info WHERE `f_doc_id` = '%s'"),
                    dbName.getCStr(), dbOper->EscapeEx(cid).getCStr());
    dbOper->Select (strSql, results);

    if (results.size() == 0) {
        return true;
    }

    for (size_t i = 0; i < results.size(); ++i) {
        if (IsSameNetworkSegment(results[i][0], results[i][1], accessIp)) {
            return true;
        }
    }
    return false;
}

/* [notxpcom] void filterByNetDocLimit (in SetStringRef docIdSet, [const] in StringRef accessIp);*/
NS_IMETHODIMP_(void) ncACSShareMgnt::FilterByNetDocLimit(set<String> &docIdSet, const String &accessIp)
{
    NC_ACS_SHAREMGNT_TRACE (_T("this: %p"), this);

    // 获取所有文档库id
    String cid;
    set<String> allCidSet;
    for(set<String>::iterator idIt = docIdSet.begin (); idIt != docIdSet.end (); ++idIt) {
        if (ncGNSUtil::IsCIDGNS (*idIt) == false) {
            cid = ncGNSUtil::GetCIDPath (*idIt);
        }
        else {
            cid = *idIt;
        }
        allCidSet.insert(cid);
    }

    if (allCidSet.size() == 0) {
        return;
    }

    // 根据cid获取所有可以访问的ip网段
    nsCOMPtr<ncIDBOperator> dbOper = getter_AddRefs (GetDBOperator ());
    String strSql;
    ncDBRecords results;
    String groupStr = GenerateGroupStrBySet(allCidSet);
    String dbName = Util::getDBName("sharemgnt_db");
    strSql.format (_T("SELECT DISTINCT `f_ip`, `f_sub_net_mask`, `f_doc_id` FROM %s.t_net_docs_limit_info WHERE `f_doc_id` in (%s)"), dbName.getCStr(), groupStr.getCStr());
    dbOper->Select (strSql, results);
    if (results.size() == 0) {
        return;
    }

    // 获取可以访问的文档库id
    set<String> canAccessCidSet;
    // 获取配置过网段的文档库
    set<String> hasNetLimitCidSet;
    for (size_t i = 0; i < results.size(); ++i) {
        if (IsSameNetworkSegment(results[i][0], results[i][1], accessIp)) {
            canAccessCidSet.insert(results[i][2]);
        }
        hasNetLimitCidSet.insert(results[i][2]);
    }

    // 去掉受到网段限制的文档库id
    for(set<String>::iterator idIt = docIdSet.begin (); idIt != docIdSet.end ();) {
        if (ncGNSUtil::IsCIDGNS (*idIt) == false) {
            cid = ncGNSUtil::GetCIDPath (*idIt);
        }
        else {
            cid = *idIt;
        }
        if (hasNetLimitCidSet.find(cid) != hasNetLimitCidSet.end() && canAccessCidSet.find(cid) == canAccessCidSet.end()) {
            #if (__GNUC__ == 4 && __GNUC_MINOR__ > 4)
                idIt = docIdSet.erase (idIt);
            #else
                docIdSet.erase (idIt++);
            #endif
        } else {
            ++idIt;
        }
    }

}

/* [notxpcom] void UpdateUserLastRequestTime ();*/
NS_IMETHODIMP_(void) ncACSShareMgnt::UpdateUserLastRequestTime(const String& userId)
{
    NC_ACS_SHAREMGNT_TRACE (_T("this: %p"), this);

    nsCOMPtr<ncIDBOperator> dbOper = getter_AddRefs (GetDBOperator ());

    ptime localCur (BusinessDate::getLocalTime ());
    date localCurDay = localCur.date();
    time_duration localCurTime = localCur.time_of_day();

    String curTimeStr;
    curTimeStr.format (_T("%04d-%02d-%02d %02d:%02d:%02d"),
        (int)localCurDay.year (), (int)localCurDay.month (), (int)localCurDay.day (),
        (int)localCurTime.hours (), (int)localCurTime.minutes (), (int)localCurTime.seconds ());

    String dbName = Util::getDBName("sharemgnt_db");
    //更新用户最近请求时间（所有平台） 和客户端用户最近请求时间（只考虑客户端）
    String strSql;
    strSql.format (_T("update %s.t_user set `f_last_request_time` = '%s', `f_last_client_request_time` = '%s', `f_activate_status` = '1' ")
                   _T("where f_user_id = '%s'"),
                   dbName.getCStr(),
                   curTimeStr.getCStr (), curTimeStr.getCStr (), dbOper->EscapeEx(userId).getCStr());

    ncDBRecords results;
    dbOper->Execute (strSql);

    NC_ACS_SHAREMGNT_TRACE (_T("this: %p"), this);
}

/* [notxpcom] void AgreedToTermsOfUse ([const] in StringRef userId);*/
NS_IMETHODIMP_(void) ncACSShareMgnt::AgreedToTermsOfUse(const String& userId)
{
    NC_ACS_SHAREMGNT_TRACE (_T("this: %p"), this);

    nsCOMPtr<ncIDBOperator> dbOper = getter_AddRefs (GetDBOperator ());

    String dbName = Util::getDBName("sharemgnt_db");
    String strSql;
    strSql.format (_T("update %s.t_user set `f_agreed_to_terms_of_use` = '1' where f_user_id = '%s'"),
                   dbName.getCStr(), dbOper->EscapeEx(userId).getCStr());

    ncDBRecords results;
    dbOper->Execute (strSql);

    NC_ACS_SHAREMGNT_TRACE (_T("this: %p"), this);
}

/* [notxpcom] void UpdateUserDocumentReadStatus ([const] in StringRef userId, in int leftShiftNum, );*/
NS_IMETHODIMP_(void) ncACSShareMgnt::UpdateUserDocumentReadStatus(const String& userId, int leftShiftNum)
{
    NC_ACS_SHAREMGNT_TRACE (_T("this: %p"), this);

    nsCOMPtr<ncIDBOperator> dbOper = getter_AddRefs (GetDBOperator ());

    String dbName = Util::getDBName("sharemgnt_db");
    int64 leftShiftResult = 1 << leftShiftNum;
    String strSql;
    strSql.format (_T("update %s.t_user set `f_user_document_read_status` = f_user_document_read_status|%lld ")
                   _T("where f_user_id = '%s'"),
                   dbName.getCStr(),
                   leftShiftResult,
                   dbOper->EscapeEx(userId).getCStr());

    ncDBRecords results;
    dbOper->Execute (strSql);

    NC_ACS_SHAREMGNT_TRACE (_T("this: %p"), this);
}

bool ncACSShareMgnt::IsSameNetworkSegment(const String &ip,const String &netMask,const String &accessIp) {

    uint32_t uip = inet_addr(ip.getCStr());
    uint32_t unetMask = inet_addr(netMask.getCStr());
    uint32_t uaccessIp = inet_addr(accessIp.getCStr());
    if ((uip & unetMask) == (uaccessIp & unetMask)) {
        return true;
    }

    return false;
}

String ncACSShareMgnt::GenerateGroupStr (const vector<String>& strs)
{
    nsCOMPtr<ncIDBOperator> dbOper = getter_AddRefs (GetDBOperator ());
    String groupStr;
    for (size_t i = 0; i < strs.size (); ++i) {
        groupStr.append ("\'", 1);
        groupStr.append (dbOper->EscapeEx(strs[i]));
        groupStr.append ("\'", 1);

        if (i != (strs.size () -1)) {
            groupStr.append (",", 1);
        }
    }

    return groupStr;
}

/* 调用该方法前，必须确保参数无sql注入的情况 */
String ncACSShareMgnt::GenerateGroupStrWithOutEscapeEx (const vector<String>& strs)
{
    nsCOMPtr<ncIDBOperator> dbOper = getter_AddRefs (GetDBOperator ());
    String groupStr;
    for (size_t i = 0; i < strs.size (); ++i) {
        groupStr.append ("\'", 1);
        groupStr.append (strs[i]);
        groupStr.append ("\'", 1);

        if (i != (strs.size () -1)) {
            groupStr.append (",", 1);
        }
    }

    return groupStr;
}

String ncACSShareMgnt::GenerateGroupStrBySet (const set<String>& strs)
{
    nsCOMPtr<ncIDBOperator> dbOper = getter_AddRefs (GetDBOperator ());
    String groupStr;
    set<String>::iterator iter = strs.begin ();
    for (; iter != strs.end (); ) {
        groupStr.append ("\'", 1);
        groupStr.append (dbOper->EscapeEx(*iter));
        groupStr.append ("\'", 1);

        if (++iter != strs.end ()) {
            groupStr.append (",", 1);
        }
    }

    return groupStr;
}

void ncACSShareMgnt::RemoveDuplicateStrs (vector<String>& strs)
{
    // 先进行排序
    sort (strs.begin(), strs.end());

    // 在删除掉相邻重复的
    vector<String>::iterator pos = unique (strs.begin(), strs.end());

    // 删除掉最后无效的条目
    strs.erase (pos, strs.end());
}

ncIDBOperator* ncACSShareMgnt::GetDBOperator ()
{
    nsCOMPtr<ncIDBOperator> dbOper = getter_AddRefs (ncACSShareMgntGetDBOperator ());
    NS_ADDREF (dbOper.get ());

    return dbOper;
}

/* [notxpcom] bool GetDefaulStrategySuperimStatus ();*/
NS_IMETHODIMP_(bool) ncACSShareMgnt::GetDefaulStrategySuperimStatus()
{
    NC_ACS_SHAREMGNT_TRACE (_T("this: %p"), this);
    return (Int::getValue(GetShareMgntConfig("default_strategy_superim_status")) == 1);
}

/* [notxpcom] bool IsUserFreeze ([const] in StringRef userId);*/
NS_IMETHODIMP_(bool) ncACSShareMgnt::IsUserFreeze (const String& userId)
{
    // 获取用户冻结状态
    NC_ACS_SHAREMGNT_TRACE (_T("begin this: %p"), this);

    String dbName = Util::getDBName("sharemgnt_db");
    nsCOMPtr<ncIDBOperator> dbOper = getter_AddRefs (GetDBOperator ());
    ncDBRecords results;
    String strSql;
    strSql.format (_T("SELECT `f_freeze_status` FROM %s.`t_user` WHERE `f_user_id` = '%s'"), dbName.getCStr(), dbOper->EscapeEx(userId).getCStr ());
    dbOper->Select (strSql, results);

    bool freezeStatus = false;
    if (results.size() == 1 && Int::getValue (results[0][0]) == 1)
        freezeStatus = true;

    NC_ACS_SHAREMGNT_TRACE (_T("end this: %p"), this);
    return freezeStatus;
}

/* [notxpcom] bool GetFreezeStatus ();*/
NS_IMETHODIMP_(bool) ncACSShareMgnt::GetFreezeStatus ()
{
    NC_ACS_SHAREMGNT_TRACE (_T("this: %p"), this);
    return (Int::getValue(GetShareMgntConfig("enable_freeze")) == 1);
}

/* [notxpcom] void GetMailAddress ([const] in StringRef reveiverId, [const] in StringRef tbName, [const] in StringRef fieldName, in StringRef email); */
NS_IMETHODIMP_(void) ncACSShareMgnt::GetMailAddress(const String & receiverId, const String & tbName, const String & fieldName, String & email)
{
    NC_ACS_SHAREMGNT_TRACE (_T("begin this: %p"), this);

    nsCOMPtr<ncIDBOperator> dbOper = getter_AddRefs (GetDBOperator ());
    String dbName = Util::getDBName("sharemgnt_db");
    String strSql;
    strSql.format (_T("select f_mail_address from %s.%s where %s = '%s';"),
                   dbName.getCStr(), tbName.getCStr (), fieldName.getCStr (), dbOper->EscapeEx(receiverId).getCStr ());

    ncDBRecords results;
    dbOper->Select (strSql, results);

    if (results.size () == 1) {
        email = std::move(results[0][0]);
    }
    else {
        email = String::EMPTY;
    }

    NC_ACS_SHAREMGNT_TRACE (_T("end this: %p"), this);
}

/* [notxpcom] bool GetWaterMarkStrategy (); */
NS_IMETHODIMP_(bool) ncACSShareMgnt::GetWaterMarkStrategy()
{
    NC_ACS_SHAREMGNT_TRACE (_T("begin this: %p"), this);

    nsCOMPtr<ncIDBOperator> dbOper = getter_AddRefs (GetDBOperator ());
    String dbName = Util::getDBName("sharemgnt_db");
    String strSql;
    strSql.format (_T("select f_config from %s.t_watermark_config;"), dbName.getCStr());

    ncDBRecords results;
    dbOper->Select (strSql, results);

    // 用户名水印和自定义水印至少有一个需要开启
    if (results.size () != 1) {
        return false;
    }
    else
    {
        String watermarkConfig = std::move(results[0][0]);
        JSON::Value watermarkConfigJson;
        JSON::Reader::read (watermarkConfigJson, watermarkConfig.getCStr (), watermarkConfig.getLength ());
        if (!watermarkConfigJson["text"]["enabled"].b () &&
            !watermarkConfigJson["user"]["enabled"].b () &&
            !watermarkConfigJson["date"]["enabled"].b ()) {
            return false;
        }
    }

    NC_ACS_SHAREMGNT_TRACE (_T("end this: %p"), this);

    return true;
}

/* [notxpcom] void GetDownloadWatermarkDocs(in StringIntMapRef WatermarkDocsMap); */
NS_IMETHODIMP_(void) ncACSShareMgnt::GetDownloadWatermarkDocs(map<String, int> & WatermarkDocsMap)
{
    NC_ACS_SHAREMGNT_TRACE (_T("begin this: %p"), this);

    nsCOMPtr<ncIDBOperator> dbOper = getter_AddRefs (GetDBOperator ());
    String dbName = Util::getDBName("sharemgnt_db");
    String strSql;
    strSql.format (_T("select f_obj_id, f_watermark_type from %s.t_watermark_doc;"), dbName.getCStr());

    WatermarkDocsMap.clear ();
    ncDBRecords results;
    dbOper->Select (strSql, results);
    for (size_t i = 0; i < results.size (); ++i) {
        WatermarkDocsMap.insert (pair<String, int>(std::move(results[i][0]), Int::getValue(results[i][1])));
    }

    NC_ACS_SHAREMGNT_TRACE (_T("end this: %p"), this);
}

/* [notxpcom] void GetDownloadWatermarkDocTypes(in IntMapRef WatermarkDocTypesMap); */
NS_IMETHODIMP_(void) ncACSShareMgnt::GetDownloadWatermarkDocTypes(map<int, int> & WatermarkDocTypesMap)
{
    NC_ACS_SHAREMGNT_TRACE (_T("begin this: %p"), this);

    nsCOMPtr<ncIDBOperator> dbOper = getter_AddRefs (GetDBOperator ());
    String dbName = Util::getDBName("sharemgnt_db");
    String strSql;
    strSql.format (_T("select f_for_user_doc, f_for_custom_doc, f_for_archive_doc from %s.t_watermark_config;"), dbName.getCStr());

    WatermarkDocTypesMap.clear ();
    ncDBRecords results;
    dbOper->Select (strSql, results);
    if (results.size() == 1) {

        // 1：个人文档，3：自定义文档库，4: 共享文档（使用分享者个人文档配置），5：归档库，
        // DocType取值：2:下载水印，3:预览水印 + 下载水印
        WatermarkDocTypesMap.insert (pair <int, int>(1, Int::getValue(results[0][0])));
        WatermarkDocTypesMap.insert (pair <int, int>(3, Int::getValue(results[0][1])));
        WatermarkDocTypesMap.insert (pair <int, int>(4, Int::getValue(results[0][0])));
        WatermarkDocTypesMap.insert (pair <int, int>(5, Int::getValue(results[0][2])));
    }

    NC_ACS_SHAREMGNT_TRACE (_T("end this: %p"), this);
}

void ncACSShareMgnt::SortObjIdsWithPriority (const String& groupStr, const size_t objType, vector<String>& objIds)
{
    nsCOMPtr<ncIDBOperator> dbOper = getter_AddRefs (GetDBOperator ());
    String dbName = Util::getDBName("sharemgnt_db");
    String deptSql, userSql;
    deptSql.format(_T("select f_department_id from %s.t_department where f_department_id in (%s) ")
                   _T("order by f_priority, upper(f_name)"),
                   dbName.getCStr(), groupStr.getCStr());
    userSql.format(_T("select f_user_id from %s.t_user where f_user_id in (%s) ")
                   _T("order by f_priority, upper(f_display_name)"),
                   dbName.getCStr(), groupStr.getCStr());

    ncDBRecords results;
    if (objType == IOC_DEPARTMENT) {
        dbOper->Select (deptSql, results);
    }
    else if (objType == IOC_USER) {
        dbOper->Select (userSql, results);
    }

    for (size_t i = 0; i < results.size(); i++) {
        objIds.push_back(results[i][0]);
    }
}

/* [notxpcom] bool GetRealNameAuthStatus ();*/
NS_IMETHODIMP_(bool) ncACSShareMgnt::GetRealNameAuthStatus ()
{
    NC_ACS_SHAREMGNT_TRACE (_T("this: %p"), this);
    return (Int::getValue(GetShareMgntConfig("enable_real_name_auth")) == 1);
}

/* [notxpcom] bool IsUserRealNameAuth ([const] in StringRef userId);*/
NS_IMETHODIMP_(bool) ncACSShareMgnt::IsUserRealNameAuth (const String& userId)
{
    // 获取用户实名认证状态
    NC_ACS_SHAREMGNT_TRACE (_T("begin this: %p"), this);

    String dbName = Util::getDBName("sharemgnt_db");
    nsCOMPtr<ncIDBOperator> dbOper = getter_AddRefs (GetDBOperator ());
    ncDBRecords results;
    String strSql;
    strSql.format (_T("SELECT `f_real_name_auth_status` FROM %s.`t_user` WHERE `f_user_id` = '%s'"), dbName.getCStr(), dbOper->EscapeEx(userId).getCStr ());
    dbOper->Select (strSql, results);

    bool status = false;
    if (results.size() == 1 && Int::getValue (results[0][0]) == 1)
        status = true;

    NC_ACS_SHAREMGNT_TRACE (_T("end this: %p"), this);
    return status;
}

/* [notxpcom] bool GetFileCrawlStatus ();*/
NS_IMETHODIMP_(bool) ncACSShareMgnt::GetFileCrawlStatus ()
{
    // 获取文档抓取策略开关状态
    NC_ACS_SHAREMGNT_TRACE (_T("this: %p"), this);

    String dbName = Util::getDBName("sharemgnt_db");
    String strSql;
    strSql.format (_T("select f_value from %s.t_sharemgnt_config where f_key = 'file_crawl_status';"), dbName.getCStr());

    nsCOMPtr<ncIDBOperator> dbOper = getter_AddRefs (GetDBOperator ());
    ncDBRecords results;
    dbOper->Select (strSql, results);

    bool status = true;
    if (results.size() == 1 && Int::getValue (results[0][0]) == 0) {
        status = false;
    }

    NC_ACS_SHAREMGNT_TRACE (_T("end this: %p"), this);

    return status;
}

/* [notxpcom] bool isFileCrawlStrategy([const] in StringRef userId, [const] in StringRef docId);*/
NS_IMETHODIMP_(bool) ncACSShareMgnt::IsFileCrawlStrategy (const String& userId, const String& docId)
{
    // 判断是否存在文档抓取策略
    NC_ACS_SHAREMGNT_TRACE (_T("this: %p, userId: %s, docId: docId: %s"), this, userId.getCStr (), docId.getCStr ());

    String dbName = Util::getDBName("sharemgnt_db");
    nsCOMPtr<ncIDBOperator> dbOper = getter_AddRefs (GetDBOperator ());
    ncDBRecords results;
    String strSql;
    strSql.format (_T("select f_strategy_id from %s.t_file_crawl_strategy where f_user_id = '%s' and f_doc_id = '%s'"),
                   dbName.getCStr(),
                   dbOper->EscapeEx(userId).getCStr (),
                   dbOper->EscapeEx(ncGNSUtil::GetCIDPath (docId)).getCStr ());
    dbOper->Select (strSql, results);

    NC_ACS_SHAREMGNT_TRACE (_T("end this: %p"), this);

    return (results.size() == 1);
}

String ncACSShareMgnt::GenerateSearchStr (const bool& exact_search_user, const int& searchRange, const String escKey, const String escLikeKey)
{
    String searchStr;
    if (exact_search_user) {
        // 默认可以搜索用户名或者显示名
        searchStr.format (_T("(t_user.f_login_name = '%s' or t_user.f_display_name = '%s')"), escKey.getCStr (), escKey.getCStr ());

        if (searchRange == LOGIN_NAME) {
            searchStr.format (_T("t_user.f_login_name = '%s'"), escKey.getCStr ());
        }
        else if (searchRange == DISPLAY_NAME) {
            searchStr.format (_T("t_user.f_display_name = '%s'"), escKey.getCStr ());
        }
    }
    else {
        // 默认可以搜索用户名或者显示名
        searchStr.format(_T("(t_user.f_login_name = '%s' or t_user.f_display_name = '%s' or t_user.f_login_name like '%%%s%%' ESCAPE '\\\\' or t_user.f_display_name like '%%%s%%' ESCAPE '\\\\') ")
                         _T("order by t_user.f_priority, case when t_user.f_login_name = '%s' then 0 when t_user.f_display_name = '%s' then 1 when t_user.f_login_name like '%%%s%%' ESCAPE '\\\\' then 2 else 3 end, ")
                         _T("upper(t_user.f_display_name)"),
                         escKey.getCStr(), escKey.getCStr(), escLikeKey.getCStr(), escLikeKey.getCStr(), escKey.getCStr(), escKey.getCStr(), escLikeKey.getCStr());

        if (searchRange == LOGIN_NAME) {
            searchStr.format(_T("(t_user.f_login_name = '%s' or t_user.f_login_name like '%%%s%%' ESCAPE '\\\\') ")
                             _T("order by t_user.f_priority, case when t_user.f_login_name = '%s' then 0 else 1 end, ")
                             _T("upper(t_user.f_display_name)"),
                             escKey.getCStr(), escLikeKey.getCStr(), escKey.getCStr());
        }
        else if (searchRange == DISPLAY_NAME) {
            searchStr.format(_T("(t_user.f_display_name = '%s' or t_user.f_display_name like '%%%s%%' ESCAPE '\\\\') ")
                             _T("order by t_user.f_priority, case when t_user.f_display_name = '%s' then 0 else 1 end, ")
                             _T("upper(t_user.f_display_name)"),
                             escKey.getCStr(), escLikeKey.getCStr(), escKey.getCStr());
        }
    }

    return searchStr;
}

/* [notxpcom] bool IsDepartmentExist([const] in StringRef departId); */
NS_IMETHODIMP_(bool) ncACSShareMgnt::IsDepartmentExist (const String& departId)
{
    NC_ACS_SHAREMGNT_TRACE (_T("begin this: %p"), this);

    nsCOMPtr<ncIDBOperator> dbOper = getter_AddRefs (GetDBOperator ());
    String dbName = Util::getDBName("sharemgnt_db");
    String strSql;
    strSql.format (_T("select f_department_id from %s.t_department where f_department_id = '%s';"),
                   dbName.getCStr(), dbOper->EscapeEx(departId).getCStr ());

    ncDBRecords results;
    dbOper->Select (strSql, results);

    bool isExist = false;
    if (results.size() > 0) {
        isExist = true;
    }

    NC_ACS_SHAREMGNT_TRACE (_T("end this: %p"), this);

    return isExist;
}

/* [notxpcom] void GetSubDepartIds(in StringVecRef subDepartIds, [const] in StringVecRef departId); */
NS_IMETHODIMP_(void) ncACSShareMgnt::GetSubDepartIds(vector<String>& subDepartIds, const vector<String>& departIds)
{
    NC_ACS_SHAREMGNT_TRACE (_T("begin this: %p"), this);

    nsCOMPtr<ncIDBOperator> dbOper = getter_AddRefs (GetDBOperator ());

    if (departIds.size() == 0) {
        return;
    }

    String dbName = Util::getDBName("sharemgnt_db");
    String strSql;
    String groupStr = GenerateGroupStrWithOutEscapeEx (departIds);
    strSql.format (_T("select f_department_id from %s.t_department_relation where f_parent_department_id in (%s);"),
                   dbName.getCStr(), groupStr.getCStr ());

    ncDBRecords results;
    dbOper->Select (strSql, results);

    for (size_t i = 0; i < results.size(); i++) {
        subDepartIds.push_back(std::move(results[i][0]));
    }

    NC_ACS_SHAREMGNT_TRACE (_T("end this: %p"), this);
}

NS_IMETHODIMP_(void) ncACSShareMgnt::GetDeepestSubDepartIds(vector<String>& relateDepartIds, const vector<String>& tmpDepartIds)
{
    queue<String> tmpIds;

    // 获取第一级子部门
    vector<String> subDepartIds;
    GetSubDepartIds(subDepartIds, tmpDepartIds);
    for (size_t j = 0; j < subDepartIds.size(); j++) {
        tmpIds.push(subDepartIds[j]);
    }

    // 按照广度遍历
    String tmpId;
    while (tmpIds.size () > 0) {
        // 获取子部门
        tmpId = tmpIds.front();
        vector<String> subDepartIds;
        vector<String> tmpDepIds;
        tmpDepIds.push_back(tmpId);
        GetSubDepartIds(subDepartIds, tmpDepIds);
        for (size_t j = 0; j < subDepartIds.size(); j++) {
            tmpIds.push(subDepartIds[j]);
        }
        // 是最下级部门，则添加
        if (subDepartIds.size() == 0) {
            relateDepartIds.push_back(tmpId);
        }
        tmpIds.pop();
    }
}

NS_IMETHODIMP_(void) ncACSShareMgnt::GetAllSubDepartIds(vector<String>& relateDepartIds, const vector<String>& tmpDepartIds)
{
    // 获取第一级子部门
    vector<String> subDepartIds;
    GetSubDepartIds(subDepartIds, tmpDepartIds);

    // 按照广度添加子部门
    while (subDepartIds.size() > 0) {
        for (size_t i = 0; i < subDepartIds.size(); i++) {
            relateDepartIds.push_back(subDepartIds[i]);
        }
        vector<String> tmpDepartIds;
        GetSubDepartIds(tmpDepartIds, subDepartIds);
        subDepartIds = tmpDepartIds;
    }
}

NS_IMETHODIMP_(void) ncACSShareMgnt::GetSubDepartIdsByLevel(vector<String>& relateDepartIds, const vector<String>& tmpDepartIds, const int level)
{
    NC_ACS_SHAREMGNT_TRACE (_T("begin this: %p"), this);

    vector<String> tmpIds = tmpDepartIds;
    int tmpLevel = level;
    while (tmpLevel > 0) {
        vector<String> subDepartIds;
        GetSubDepartIds(subDepartIds, tmpIds);
        if (subDepartIds.size() > 0) {
            tmpIds = subDepartIds;
        } else {
            // 没有指定层级的子部门，返回
            return;
        }
        tmpLevel--;
    }

    if (tmpIds.size() > 0) {
        relateDepartIds = tmpIds;
    }

    NC_ACS_SHAREMGNT_TRACE (_T("end this: %p"), this);
}

/* [notxpcom] void GetRelateDepartIds(in StringVecRef relateDepartIds, [const] in StringVecRef tmpDepartIds, [const] in bool includeCurDepart, [const] in int level); */
NS_IMETHODIMP_(void) ncACSShareMgnt::GetRelateDepartIds(vector<String>& relateDepartIds, const vector<String>& tmpDepartIds, const bool includeCurDepart, const int level)
{
    NC_ACS_SHAREMGNT_TRACE (_T("begin this: %p"), this);

    if (tmpDepartIds.size() == 0) {
        return;
    }

    if (level == -2) {
        GetDeepestSubDepartIds(relateDepartIds, tmpDepartIds);
    }

    if (level == -1) {
        GetAllSubDepartIds(relateDepartIds, tmpDepartIds);
    }

    if (level > 0) {
        GetSubDepartIdsByLevel(relateDepartIds, tmpDepartIds, level);
    }

    if (includeCurDepart) {
        for (size_t i = 0; i < tmpDepartIds.size(); i++) {
            relateDepartIds.push_back(tmpDepartIds[i]);
        }
    }

    NC_ACS_SHAREMGNT_TRACE (_T("end this: %p"), this);
}

/* [notxpcom] void GetDocInfoOfDeparts(in ncDepDocInfoVecRef depDocInfos, [const] in StringVecRef departIds, [const] in dbOwnerInfoVecRef ownerInfoVec); */
NS_IMETHODIMP_(void) ncACSShareMgnt::GetDocInfoOfDeparts(vector<ncDepDocInfo>& depDocInfos, const vector<String>& departIds, const vector<dbOwnerInfo>& ownerInfoVec)
{
    NC_ACS_SHAREMGNT_TRACE (_T("begin this: %p"), this);

    nsCOMPtr<ncIDBOperator> dbOper = getter_AddRefs (GetDBOperator ());
    String dbName = Util::getDBName("sharemgnt_db");
    for (size_t i = 0; i < departIds.size(); i++) {
        ncDepDocInfo info;
        String strSql;
        if (departIds[i].isEmpty ()) {
            continue;
        }

        // 获取部门名、站点id、站点名
        strSql.format (_T("select f_name, f_oss_id from %s.t_department where t_department.f_department_id = '%s';"),
                       dbName.getCStr(), dbOper->EscapeEx(departIds[i]).getCStr ());

        ncDBRecords results;
        dbOper->Select (strSql, results);

        if (results.size () > 0) {
            info.name = std::move(results[0][0]);
            info.ossId = std::move(results[0][1]);
        }
        else {
            continue;
        }

        // 获取所有者id
        if (ownerInfoVec.size () == 0) {
            strSql.format(_T("select f_user_id from %s.t_department_responsible_person where f_department_id = '%s';"),
                          dbName.getCStr(), dbOper->EscapeEx(departIds[i]).getCStr ());
            ncDBRecords results;
            dbOper->Select (strSql, results);
            for (size_t j = 0; j < results.size(); j++) {
                dbOwnerInfo ownerInfo;
                ownerInfo.ownerId = std::move(results[j][0]);
                ownerInfo.ownerType = IOC_USER;
                info.ownerInfos.push_back (ownerInfo);
            }
        } else {
            info.ownerInfos = ownerInfoVec;
        }

        info.id = departIds[i];
        depDocInfos.push_back(info);
    }

    NC_ACS_SHAREMGNT_TRACE (_T("end this: %p"), this);
}

/* [notxpcom] void GetParentDeptPath ([const] in StringRef deptId, in StringVecRef deptNames) */
NS_IMETHODIMP_(void) ncACSShareMgnt::GetParentDeptPath(const String& deptId, vector<String>& deptNames)
{
    NC_ACS_SHAREMGNT_TRACE (_T("this: %p"), this);

    deptNames.clear();

    nsCOMPtr<ncIDBOperator> dbOper = getter_AddRefs (GetDBOperator ());
    String strSql;

    String tmpId = deptId;
    String dbName = Util::getDBName("sharemgnt_db");
    while (1) {
        if (tmpId.isEmpty ()) {
            break;
        }

        strSql.format (_T("select t_department.f_name, t_department.f_department_id ")
                       _T("from %s.t_department, %s.t_department_relation ")
                       _T("where t_department.f_department_id = t_department_relation.f_parent_department_id ")
                       _T("and t_department_relation.f_department_id = '%s';"),
                       dbName.getCStr(), dbName.getCStr(),
                       dbOper->EscapeEx(tmpId).getCStr ());
        ncDBRecords results;
        dbOper->Select (strSql, results);
        if (results.size () > 0) {
            deptNames.push_back(results[0][0].getCStr ());
            tmpId = results[0][1];
        }
        else {
            break;
        }
    }

    reverse(deptNames.begin(), deptNames.end());

    NC_ACS_SHAREMGNT_TRACE (_T("this: %p"), this);
}

/* [notxpcom] void GetUserRoleIds ([const] in StringRef userId, in StringVecRef roleIds); */
NS_IMETHODIMP_(void) ncACSShareMgnt::GetUserRoleIds(const String& userId, vector<String>& roleIds)
{
    NC_ACS_SHAREMGNT_TRACE (_T("this: %p"), this);

    roleIds.clear();

    String dbName = Util::getDBName("sharemgnt_db");
    nsCOMPtr<ncIDBOperator> dbOper = getter_AddRefs (GetDBOperator ());
    String strSql;
    strSql.format (_T("select f_role_id from %s.t_user_role_relation "
                     "where f_user_id = '%s' "),
                     dbName.getCStr(), dbOper->EscapeEx(userId).getCStr ());

    ncDBRecords results;
    dbOper->Select (strSql, results);

    for (size_t i = 0; i < results.size(); i++) {
        roleIds.push_back(std::move(results[i][0]));
    }

    NC_ACS_SHAREMGNT_TRACE (_T("this: %p"), this);
}

/* [notxpcom] void GetUserIdsByRoleId ([const] in StringRef roleId, in StringVecRef userIds); */
NS_IMETHODIMP_(void) ncACSShareMgnt::GetUserIdsByRoleId(const String& roleId, vector<String>& userIds)
{
    NC_ACS_SHAREMGNT_TRACE (_T("this: %p"), this);

    userIds.clear();

    String dbName = Util::getDBName("sharemgnt_db");
    nsCOMPtr<ncIDBOperator> dbOper = getter_AddRefs (GetDBOperator ());
    String strSql;
    strSql.format (_T("select f_user_id from %s.t_user_role_relation "
                     "where f_role_id = '%s' "),
                     dbName.getCStr(), dbOper->EscapeEx(roleId).getCStr ());

    ncDBRecords results;
    dbOper->Select (strSql, results);

    for (size_t i = 0; i < results.size(); i++) {
        userIds.push_back(std::move(results[i][0]));
    }

    NC_ACS_SHAREMGNT_TRACE (_T("this: %p"), this);
}

/* [notxpcom] void GetUserRole ([const] in StringRef userId, in ncRoleInfoVecRef roleInfos); */
NS_IMETHODIMP_(void) ncACSShareMgnt::GetUserRole(const String& userId, vector<ncRoleInfo>& roleInfos)
{
    NC_ACS_SHAREMGNT_TRACE (_T("begin this: %p"), this);

    nsCOMPtr<ncIDBOperator> dbOper = getter_AddRefs (GetDBOperator ());
    roleInfos.clear();

    String dbName = Util::getDBName("sharemgnt_db");
    String strSql;
    strSql.format(_T("select t.f_role_id, t.f_name from %s.`t_role` as t "
                     "inner join %s.t_user_role_relation as r on t.f_role_id = r.f_role_id "
                     "where r.f_user_id =  '%s' ORDER BY t.f_priority;"),
                     dbName.getCStr(), dbName.getCStr(), dbOper->EscapeEx(userId).getCStr ());
    ncDBRecords results;
    dbOper->Select (strSql, results);
    for (size_t j = 0; j < results.size(); j++) {
        ncRoleInfo info;
        info.id = std::move(results[j][0]);
        info.name = std::move(results[j][1]);
        roleInfos.push_back(info);
    }

    NC_ACS_SHAREMGNT_TRACE (_T("end this: %p"), this);
}

/* [notxpcom] bool IsAdminRole ([const] in StringRef userId); */
NS_IMETHODIMP_(bool) ncACSShareMgnt::IsAdminRole(const String& userId)
{
    NC_ACS_SHAREMGNT_TRACE (_T("this: %p"), this);

    String dbName = Util::getDBName("sharemgnt_db");
    nsCOMPtr<ncIDBOperator> dbOper = getter_AddRefs (GetDBOperator ());
    String strSql;
    strSql.format (_T("select f_user_id from %s.t_user_role_relation "
                     "where f_user_id = '%s' and f_role_id in ('%s', '%s') "),
                     dbName.getCStr(),
                     dbOper->EscapeEx(userId).getCStr (), g_ShareMgnt_constants.NCT_SYSTEM_ROLE_SUPPER.c_str(), g_ShareMgnt_constants.NCT_SYSTEM_ROLE_ADMIN.c_str());

    ncDBRecords results;
    dbOper->Select (strSql, results);

    if (results.size() > 0) {
        return true;
    }

    NC_ACS_SHAREMGNT_TRACE (_T("this: %p"), this);
    return false;
}

/* [notxpcom] bool IsSecuritRole ([const] in StringRef userId); */
NS_IMETHODIMP_(bool) ncACSShareMgnt::IsSecuritRole(const String& userId)
{
    NC_ACS_SHAREMGNT_TRACE (_T("this: %p"), this);

    String dbName = Util::getDBName("sharemgnt_db");
    nsCOMPtr<ncIDBOperator> dbOper = getter_AddRefs (GetDBOperator ());
    String strSql;
    strSql.format (_T("select f_user_id from %s.t_user_role_relation "
                     "where f_user_id = '%s' and f_role_id = '%s'"),
                     dbName.getCStr(),
                     dbOper->EscapeEx(userId).getCStr (), g_ShareMgnt_constants.NCT_SYSTEM_ROLE_SECURIT.c_str());

    ncDBRecords results;
    dbOper->Select (strSql, results);

    if (results.size() > 0) {
        return true;
    }

    NC_ACS_SHAREMGNT_TRACE (_T("this: %p"), this);
    return false;
}

/* [notxpcom] bool GetUserInfoByTelNumber ([const] in StringRef emailaddress, in ncACSUserInfoRef userInfo); */
NS_IMETHODIMP_(bool) ncACSShareMgnt::GetUserInfoByTelNumber(const String & telNumber, ncACSUserInfo & userInfo)
{

    NC_ACS_SHAREMGNT_TRACE (_T("this: %p"), this);

    nsCOMPtr<ncIDBOperator> dbOper = getter_AddRefs (GetDBOperator ());

    String dbName = Util::getDBName("sharemgnt_db");
    bool ret = false;
    String strSql;
    strSql.format (_T("select f_user_id, f_login_name, f_display_name, f_mail_address, f_status, f_tel_number, f_auth_type, f_pwd_control, f_third_party_attr, f_third_party_id from %s.t_user ")
                   _T("where f_tel_number = '%s'"),
                   dbName.getCStr(), dbOper->EscapeEx(telNumber).getCStr());
    ncDBRecords results;
    dbOper->Select (strSql, results);
    if (results.size () == 1) {
        userInfo.id = std::move(results[0][0]);
        userInfo.account = std::move(results[0][1]);
        userInfo.visionName = std::move(results[0][2]);
        userInfo.email = std::move(results[0][3]);
        userInfo.telNumber = std::move(results[0][5]);
        userInfo.authType = Int::getValue(results[0][6]);
        userInfo.pwdControl = Int::getValue(results[0][7]);
        userInfo.thirdAttr = std::move(results[0][8]);
        userInfo.thirdId = std::move(results[0][9]);
        userInfo.enableStatus = results[0][4] == "0";
        ret = true;
    }
    NC_ACS_SHAREMGNT_TRACE (_T("end this: %p"), this);

    return ret;
}

/* [notxpcom] bool GetUserInfoByEmail ([const] in StringRef telnumber, in ncACSUserInfoRef userInfo); */
NS_IMETHODIMP_(bool) ncACSShareMgnt::GetUserInfoByEmail(const String & emailaddress, ncACSUserInfo & userInfo)
{

    NC_ACS_SHAREMGNT_TRACE (_T("this: %p"), this);

    nsCOMPtr<ncIDBOperator> dbOper = getter_AddRefs (GetDBOperator ());

    String dbName = Util::getDBName("sharemgnt_db");
    bool ret = false;
    String strSql;
    strSql.format (_T("select f_user_id, f_login_name, f_display_name, f_mail_address, f_status, f_tel_number, f_auth_type, f_pwd_control from %s.t_user ")
                   _T("where f_mail_address = '%s'"), dbName.getCStr(), dbOper->EscapeEx(emailaddress).getCStr());
    ncDBRecords results;
    dbOper->Select (strSql, results);
    if (results.size () == 1) {
        userInfo.id = std::move(results[0][0]);
        userInfo.account = std::move(results[0][1]);
        userInfo.visionName = std::move(results[0][2]);
        userInfo.email = std::move(results[0][3]);
        userInfo.telNumber = std::move(results[0][5]);
        userInfo.authType = Int::getValue(results[0][6]);
        userInfo.pwdControl = Int::getValue(results[0][7]);
        userInfo.enableStatus = results[0][4] == "0";
        ret = true;
    }
    NC_ACS_SHAREMGNT_TRACE (_T("end this: %p"), this);

    return ret;
}

/* [notxpcom] bool OEM_GetConfigByOption ([const] in StringRef section, [const] in StringRef option, in StringRef value); */
NS_IMETHODIMP_(bool) ncACSShareMgnt::OEM_GetConfigByOption(const String & section, const String & option, String & value)
{
    NC_ACS_SHAREMGNT_TRACE (_T("this: %p"), this);

    nsCOMPtr<ncIDBOperator> dbOper = getter_AddRefs (GetDBOperator ());
    bool ret = false;
    String strSql;
    String dbName = Util::getDBName("sharemgnt_db");
    strSql.format (_T("select f_value from %s.t_oem_config where f_section = '%s' and f_option = '%s'"),
                    dbName.getCStr(), dbOper->EscapeEx(section).getCStr(), dbOper->EscapeEx(option).getCStr());
    ncDBRecords results;
    dbOper->Select (strSql, results);
    if (results.size () == 1) {
        value = std::move(results[0][0]);
        ret = true;
    }
    NC_ACS_SHAREMGNT_TRACE (_T("end this: %p"), this);
    return ret;
}

/* [notxpcom] void BatchUpdateUserLastRequestTime ([const] in vector<acsRefreshInfoRef> &infos); */
NS_IMETHODIMP_(void)
ncACSShareMgnt::BatchUpdateUserLastRequestTime(const vector<ncACSRefreshInfo> &infos)
{
    NC_ACS_SHAREMGNT_TRACE(_T("this: %p, infos size: %d begin"), this, (int)infos.size());

    nsCOMPtr<ncIDBOperator> dbOper = getter_AddRefs(GetDBOperator());

    String dbName = Util::getDBName("sharemgnt_db");
    try
    {
        dbOper->StartTransaction ();

        String strSql;
        for (size_t i = 0; i < infos.size(); ++i)
        {
            String lastTime = dbOper->EscapeEx(infos[i].lastRequestTime);
            if(infos[i].bUpdateClientTime)
            {
                strSql.format(_T("update %s.t_user set f_last_request_time = '%s', f_last_client_request_time = '%s' where f_user_id = '%s';"),
                          dbName.getCStr(), lastTime.getCStr(), lastTime.getCStr(), dbOper->EscapeEx(infos[i].userId).getCStr());
            }
            else
            {
                strSql.format(_T("update %s.t_user set f_last_request_time = '%s' where f_user_id = '%s';"),
                          dbName.getCStr(), lastTime.getCStr(), dbOper->EscapeEx(infos[i].userId).getCStr());
            }

            dbOper->Execute(strSql);
        }

        dbOper->Commit ();
    }
    catch (Exception &e)
    {
        dbOper->Rollback ();
    }
    catch (...)
    {
        dbOper->Rollback ();
    }

    NC_ACS_SHAREMGNT_TRACE(_T("this: %p, infos size: %d end"), this, (int)infos.size());
}

/* [notxpcom] bool GetManagerQuota ([const] in StringRef userId, in ncManagerQuotaRef departIds); */
NS_IMETHODIMP_(bool) ncACSShareMgnt::GetManagerQuota(const String & userId, ncManagerQuota & managerQuota)
{
    NC_ACS_SHAREMGNT_TRACE (_T("[BEGIN]this: %p"), this);

    nsCOMPtr<ncIDBOperator> dbOper = getter_AddRefs (GetDBOperator ());

    String dbName = Util::getDBName("sharemgnt_db");
    String strSql;
    strSql.format (_T("select f_limit_user_space, f_allocated_limit_user_space, f_limit_doc_space, f_allocated_limit_doc_space from %s.t_manager_limit_space where f_manager_id = '%s'"), dbName.getCStr(), dbOper->EscapeEx(userId).getCStr ());
    ncDBRecords results;
    dbOper->Select (strSql, results);
    if (results.size () == 0){
        return false;
    }
    managerQuota.totalUserQuota = Int64::getValue (results[0][0]);
    managerQuota.allocatedUserQuota = Int64::getValue (results[0][1]);
    managerQuota.totalDocQuota = Int64::getValue (results[0][2]);
    managerQuota.allocatedDocQuota = Int64::getValue (results[0][3]);

    NC_ACS_SHAREMGNT_TRACE (_T("[END]this: %p"), this);
}

/* [notxpcom] void GetUserIdsBydepIds ([const] in StringVecRef departIds, in StringVecRef userIds); */
NS_IMETHODIMP_(void) ncACSShareMgnt::GetUserIdsBydepIds(const vector<String>& departIds, vector<String>& userIds)
{
    NC_ACS_SHAREMGNT_TRACE (_T("[BEGIN]this: %p"), this);

    if (departIds.size () == 0){
        return;
    }
    nsCOMPtr<ncIDBOperator> dbOper = getter_AddRefs (GetDBOperator ());

    // 去除EscapeEx方法，departmentid不存在sql注入的问题
    String tmpDepartIds;
    for(size_t i = 0; i < departIds.size(); ++i) {
        tmpDepartIds.append ("\'", 1);
        tmpDepartIds.append (departIds[i]);
        tmpDepartIds.append ("\'", 1);
        if (i != (departIds.size() - 1)){
            tmpDepartIds.append (",");
        }
    }

    String dbName = Util::getDBName("sharemgnt_db");
    String strSql;
    strSql.format (_T("select distinct f_user_id from %s.t_user_department_relation where f_department_id in (%s)"), dbName.getCStr(), tmpDepartIds.getCStr ());
    ncDBRecords results;
    dbOper->Select (strSql, results);
    for (size_t i = 0; i < results.size(); i++) {
        userIds.push_back(std::move (results[i][0]));
    }

    NC_ACS_SHAREMGNT_TRACE (_T("[END]this: %p"), this);
}

/* [notxpcom] void UpdateUserActivateStatus ([const] in StringRef userId); */
NS_IMETHODIMP_(void) ncACSShareMgnt::UpdateUserActivateStatus(const String& userId)
{
    NC_ACS_SHAREMGNT_TRACE (_T("this: %p"), this);

    nsCOMPtr<ncIDBOperator> dbOper = getter_AddRefs (GetDBOperator ());

    String dbName = Util::getDBName("sharemgnt_db");
    // 更新用户激活状态为已激活
    String strSql;
    strSql.format (_T("update %s.t_user set `f_activate_status` = '1' where f_user_id = '%s'"), dbName.getCStr(), dbOper->EscapeEx(userId).getCStr());

    dbOper->Execute (strSql);

    NC_ACS_SHAREMGNT_TRACE (_T("this: %p"), this);
}

/* [notxpcom] int GetAppOrgPerm ([const] in StringRef appID, [const] in ncAppPermOrgTypeRef orgType);*/
NS_IMETHODIMP_(int) ncACSShareMgnt::GetAppOrgPerm(const String& appID, const ncAppPermOrgType& orgType)
{
    NC_ACS_SHAREMGNT_TRACE (_T("[BEGIN]this: %p"), this);

    nsCOMPtr<ncIDBOperator> dbOper = getter_AddRefs (GetDBOperator ());

    String dbName = Util::getDBName("user_management");
    String strSql;
    strSql.format (_T("select f_perm_value from %s.t_org_perm_app where f_app_id = '%s' and f_org_type = %d"), dbName.getCStr(), dbOper->EscapeEx(appID).getCStr (), orgType);
    ncDBRecords results;
    dbOper->Select (strSql, results);

    // 判断是否具有权限
    int perm = 0;
    if (results.size() == 1) {
        perm = Int::getValue(results[0][0]);
    }

    NC_ACS_SHAREMGNT_TRACE (_T("this: %p"), this);
    return perm;
}
