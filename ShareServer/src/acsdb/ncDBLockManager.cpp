#include <abprec.h>
#include <ncutil/ncutil.h>
#include <dataapi/ncGNSUtil.h>

#include "acsdb.h"
#include "ncDBLockManager.h"
#include "common/util.h"

/* Implementation file */
NS_IMPL_QUERY_INTERFACE1 (ncDBLockManager, ncIDBLockManager)

// protected
NS_IMETHODIMP_(nsrefcnt) ncDBLockManager::AddRef (void)
{
    return 1;
}

// protected
NS_IMETHODIMP_(nsrefcnt) ncDBLockManager::Release (void)
{
    return 1;
}

AB_DEFINE_THREADSAFE_SINGLETON_NO_POOL (ncDBLockManager)

ncDBLockManager::ncDBLockManager ()
{
    NC_ACS_DB_TRACE (_T("this: %p"), this);
}

ncDBLockManager::~ncDBLockManager ()
{
    NC_ACS_DB_TRACE (_T("~this: %p"), this);
}

/*[notxpcom] void Delete ([const] in StringRef fileId);*/
NS_IMETHODIMP_(void) ncDBLockManager::Delete (const String& fileId)
{
    NC_ACS_DB_TRACE (_T("this: %p, fileId: %s, begin"), this, fileId.getCStr ());

    nsCOMPtr<ncIDBOperator> dbOper = getter_AddRefs (GetDBOperator ());

    String dbName = Util::getDBName("anyshare");
    String sql;
    sql.format (_T("delete from %s.t_lock where f_doc_id = '%s'"), dbName.getCStr(), dbOper->EscapeEx(fileId).getCStr ());

    dbOper->Execute (sql);

    NC_ACS_DB_TRACE (_T("this: %p, fileId: %s, end"), this, fileId.getCStr ());
}

/*[notxpcom] void DeleteSubs ([const] in StringRef dirId);*/
NS_IMETHODIMP_(void) ncDBLockManager::DeleteSubs (const String& dirId)
{
    NC_ACS_DB_TRACE (_T("this: %p, dirId: %s, begin"), this, dirId.getCStr ());

    nsCOMPtr<ncIDBOperator> dbOper = getter_AddRefs (GetDBOperator ());

    String dbName = Util::getDBName("anyshare");
    String escDirId = dbOper->EscapeEx(dirId);
    String sql;
    sql.format (_T("delete from %s.t_lock where f_doc_id like '%s/%%' or f_doc_id = '%s'"),
        dbName.getCStr(), escDirId.getCStr (), escDirId.getCStr ());

    dbOper->Execute (sql);

    NC_ACS_DB_TRACE (_T("this: %p, dirId: %s, end"), this, dirId.getCStr ());
}

/*[notxpcom] void DeleteByUserId ([const] in StringRef userId);*/
NS_IMETHODIMP_(void) ncDBLockManager::DeleteByUserId (const String& userId)
{
    NC_ACS_DB_TRACE (_T("this: %p, userId: %s, begin"), this, userId.getCStr ());

    nsCOMPtr<ncIDBOperator> dbOper = getter_AddRefs (GetDBOperator ());

    String dbName = Util::getDBName("anyshare");
    String sql;
    sql.format (_T("delete from %s.t_lock where f_user_id = '%s'"), dbName.getCStr(), dbOper->EscapeEx(userId).getCStr ());
    dbOper->Execute (sql);

    NC_ACS_DB_TRACE (_T("this: %p, userId: %s, end"), this, userId.getCStr ());
}

ncIDBOperator* ncDBLockManager::GetDBOperator ()
{
    nsCOMPtr<ncIDBOperator> dbOper = getter_AddRefs (ncACSDBGetDBOperator ());
    NS_ADDREF (dbOper.get ());
    return dbOper;
}
