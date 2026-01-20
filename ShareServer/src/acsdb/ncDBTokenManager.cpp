#include <abprec.h>

#include "acsdb.h"
#include "ncDBTokenManager.h"
#include "gen-cpp/ShareMgnt_constants.h"
#include "common/util.h"

const int REFRESH_TOKEN_EXPIRES_TIME = 5184000;
const int VERSION_V2 = 2;

/* Implementation file */
NS_IMPL_QUERY_INTERFACE1 (ncDBTokenManager, ncIDBTokenManager)

// protected
NS_IMETHODIMP_(nsrefcnt) ncDBTokenManager::AddRef (void)
{
    return 1;
}

// protected
NS_IMETHODIMP_(nsrefcnt) ncDBTokenManager::Release (void)
{
    return 1;
}

AB_DEFINE_THREADSAFE_SINGLETON_NO_POOL (ncDBTokenManager)

ncDBTokenManager::ncDBTokenManager()
{
    NC_ACS_DB_TRACE (_T("this: %p"), this);
}

ncDBTokenManager::~ncDBTokenManager()
{
    NC_ACS_DB_TRACE (_T("this: %p"), this);
}

/* [notxpcom] int SaveActiveUser (StringMapRef activeUserInfos); */
NS_IMETHODIMP_(void) ncDBTokenManager::SaveActiveUser(map<String, String>& activeUserInfos)
{
    NC_ACS_DB_TRACE (_T("this: %p begin"), this);

    nsCOMPtr<ncIDBOperator> dbOper = getter_AddRefs (GetDBOperator ());

    String dbName = Util::getDBName("anyshare");
    String strSql;
    for (map<String, String>::iterator iter = activeUserInfos.begin(); iter != activeUserInfos.end(); ++iter) {
        strSql.format (_T("select `f_user_id` from %s.t_active_user_info where `f_user_id` = '%s' and `f_time` = '%s'"),
            dbName.getCStr(), dbOper->EscapeEx(iter->first).getCStr (), dbOper->EscapeEx(iter->second).getCStr ());
        ncDBRecords results;
        dbOper->Select (strSql, results);
        if (!results.size()) {
            strSql.format (_T("insert into %s.t_active_user_info (f_user_id, f_time) values('%s', '%s')"),
                dbName.getCStr(), dbOper->EscapeEx(iter->first).getCStr (), dbOper->EscapeEx(iter->second).getCStr ());
            dbOper->Execute (strSql);
        }
    }

    NC_ACS_DB_TRACE (_T("this: %p end"), this);
}

ncIDBOperator* ncDBTokenManager::GetDBOperator ()
{
    nsCOMPtr<ncIDBOperator> dbOper = getter_AddRefs (ncACSDBGetDBOperator ());
    NS_ADDREF (dbOper.get ());
    return dbOper;
}
