#include <abprec.h>
#include <ncutil/ncutil.h>
#include <ncutil/ncBusinessDate.h>
#include <dataapi/ncGNSUtil.h>

#include "acsdb.h"
#include "ncDBOwnerManager.h"
#include "common/util.h"

/* Implementation file */
NS_IMPL_QUERY_INTERFACE1 (ncDBOwnerManager, ncIDBOwnerManager)

// protected
NS_IMETHODIMP_(nsrefcnt) ncDBOwnerManager::AddRef (void)
{
    return 1;
}

// protected
NS_IMETHODIMP_(nsrefcnt) ncDBOwnerManager::Release (void)
{
    return 1;
}

AB_DEFINE_THREADSAFE_SINGLETON_NO_POOL (ncDBOwnerManager)

ncDBOwnerManager::ncDBOwnerManager()
{
    NC_ACS_DB_TRACE (_T("this: %p"), this);
}

ncDBOwnerManager::~ncDBOwnerManager()
{
    NC_ACS_DB_TRACE (_T("this: %p"), this);
}

/* [notxpcom] void AddOwner ([const] in dbOwnerInfoRef ownerInfo); */
NS_IMETHODIMP_(void) ncDBOwnerManager::AddOwner(const dbOwnerInfo & ownerInfo)
{
    NC_ACS_DB_TRACE (_T("this: %p, ownerInfo.docId: %s, ownerInfo.ownerId: %s begin"),
        this, ownerInfo.docId.getCStr (), ownerInfo.ownerId.getCStr ());

    nsCOMPtr<ncIDBOperator> dbOper = getter_AddRefs (GetDBOperator ());

    String dbName = Util::getDBName("anyshare");
    String strSql;
    strSql.format (_T("insert into %s.t_acs_owner (f_gns_path, f_owner_id, f_type, f_owner_name, f_modify_time, f_deletable) values ('%s', '%s', '%d', '%s', %lld, %d)"),
                    dbName.getCStr(),
                    dbOper->EscapeEx(ownerInfo.docId).getCStr (), dbOper->EscapeEx(ownerInfo.ownerId).getCStr (), ownerInfo.ownerType,
                    dbOper->EscapeEx(ownerInfo.ownerName).getCStr (), BusinessDate::getCurrentTime (), 0);
    dbOper->Execute (strSql);

    NC_ACS_DB_TRACE (_T("this: %p, ownerInfo.docId: %s, ownerInfo.ownerId: %s end"),
        this, ownerInfo.docId.getCStr (), ownerInfo.ownerId.getCStr ());
}

/* [notxpcom] void DeleteOwnerInfosByDocId ([const] in StringRef docId); */
NS_IMETHODIMP_(void) ncDBOwnerManager::DeleteOwnerInfosByDocId(const String & docId)
{
    NC_ACS_DB_TRACE (_T("this: %p, docId: %s begin"), this, docId.getCStr ());

    nsCOMPtr<ncIDBOperator> dbOper = getter_AddRefs (GetDBOperator ());
    String dbName = Util::getDBName("anyshare");
    String strSql;
    strSql.format (_T("delete from %s.t_acs_owner where f_gns_path = '%s'"), dbName.getCStr(), dbOper->EscapeEx(docId).getCStr ());
    dbOper->Execute (strSql);

    NC_ACS_DB_TRACE (_T("this: %p, docId: %s end"), this, docId.getCStr ());
}

/* [notxpcom] void GetInheritOwnerInfosByDocId ([const] in StringRef docId, in dbOwnerInfoVecRef ownerInfos, in bool onlyUser); */
NS_IMETHODIMP_(void) ncDBOwnerManager::GetInheritOwnerInfosByDocId(const String & docId, vector<dbOwnerInfo> & ownerInfos, bool onlyUser)
{
    NC_ACS_DB_TRACE (_T("this: %p, docId: %s begin"), this, docId.getCStr ());

    ownerInfos.clear ();
    nsCOMPtr<ncIDBOperator> dbOper = getter_AddRefs (GetDBOperator ());

    String dbName = Util::getDBName("anyshare");
    String groupStr;
    int depth = ncGNSUtil::GetPathDepth (docId);

    for (int i = 1; i <= depth; ++i) {
        String curDocId = ncGNSUtil::GetPathByDepth (docId, i);

        groupStr.append ("\'", 1);
        groupStr.append (dbOper->EscapeEx(curDocId));
        groupStr.append ("\'", 1);

        if (i != depth) {
            groupStr.append (",", 1);
        }
    }

    String whereClause;
    if (onlyUser) {
        whereClause = _T("and f_type = 1 ");
    }

    String strSql;
    strSql.format (_T("select f_gns_path,f_owner_id,f_owner_name,f_type,f_modify_time from %s.t_acs_owner where f_gns_path in(%s) %s ;"),
                    dbName.getCStr(), groupStr.getCStr (), whereClause.getCStr());

    ncDBRecords results;
    dbOper->Select (strSql, results);

    for (size_t i = 0; i < results.size (); ++i) {
        dbOwnerInfo info;
        info.docId = results[i][0];
        info.ownerId = results[i][1];
        info.ownerName = results[i][2];
        info.ownerType = Int::getValue(results[i][3]);
        info.modifyTime = Int64::getValue(results[i][4]);
        ownerInfos.push_back (info);
    }

    NC_ACS_DB_TRACE (_T("this: %p, docId: %s end, ret ownerInfos size: %d"), this, docId.getCStr (), (int)ownerInfos.size ());
}

/* [notxpcom] void GetOwnerInfosByUserId ([const] in StringRef userId, in dbOwnerInfoVecRef ownerInfos); */
NS_IMETHODIMP_(void) ncDBOwnerManager::GetOwnerInfosByUserId(const String & userId, vector<dbOwnerInfo> & ownerInfos)
{
    NC_ACS_DB_TRACE (_T("[BEGIN]this: %p, userId: %s"), this, userId.getCStr ());

    ownerInfos.clear ();

    nsCOMPtr<ncIDBOperator> dbOper = getter_AddRefs (GetDBOperator ());

    String dbName = Util::getDBName("anyshare");
    String strSql;
    strSql.format (_T("select f_gns_path from %s.t_acs_owner where f_owner_id = '%s'"),
                      dbName.getCStr(), dbOper->EscapeEx(userId).getCStr ());

    ncDBRecords results;
    dbOper->Select (strSql, results);

    for (size_t i = 0; i < results.size (); ++i) {
        dbOwnerInfo info;
        info.docId = results[i][0];
        info.ownerId = userId;
        ownerInfos.push_back (info);
    }

    NC_ACS_DB_TRACE (_T("[END]this: %p, userId: %s, ret ownerInfos size: %d"), this, userId.getCStr (), (int)ownerInfos.size ());
}

/* [notxpcom] void DeleteOwnerByFileId ([const] in StringRef fileId); */
NS_IMETHODIMP_(void) ncDBOwnerManager::DeleteOwnerByFileId(const String & fileId)
{
    NC_ACS_DB_TRACE (_T("this: %p, fileId: %s begin"), this, fileId.getCStr ());

    nsCOMPtr<ncIDBOperator> dbOper = getter_AddRefs (GetDBOperator ());

    String dbName = Util::getDBName("anyshare");
    String strSql;
    strSql.format (_T("delete from %s.t_acs_owner where f_gns_path = '%s'"),
                     dbName.getCStr(), dbOper->EscapeEx(fileId).getCStr ());
    dbOper->Execute (strSql);

    NC_ACS_DB_TRACE (_T("this: %p, fileId: %s end"), this, fileId.getCStr ());
}

/* [notxpcom] void DeleteOwnerByDirId ([const] in StringRef dirId); */
NS_IMETHODIMP_(void) ncDBOwnerManager::DeleteOwnerByDirId(const String & dirId)
{
    NC_ACS_DB_TRACE (_T("this: %p, dirId: %s begin"), this, dirId.getCStr ());

    nsCOMPtr<ncIDBOperator> dbOper = getter_AddRefs (GetDBOperator ());

    String dbName = Util::getDBName("anyshare");
    String escDirId = dbOper->EscapeEx(dirId);
    String strSql;
    strSql.format (_T("delete from %s.t_acs_owner where f_gns_path = '%s' or f_gns_path like '%s/%%'"),
                    dbName.getCStr(), escDirId.getCStr (), escDirId.getCStr ());

    dbOper->Execute (strSql);

    NC_ACS_DB_TRACE (_T("this: %p, dirId: %s end"), this, dirId.getCStr ());
}

/* [notxpcom] void DeleteOwnerByUserId ([const] in StringRef userId); */
NS_IMETHODIMP_(void) ncDBOwnerManager::DeleteOwnerByUserId(const String & userId)
{
    NC_ACS_DB_TRACE (_T("this: %p, userId: %s begin"), this, userId.getCStr ());

    nsCOMPtr<ncIDBOperator> dbOper = getter_AddRefs (GetDBOperator ());

    String dbName = Util::getDBName("anyshare");
    String strSql;
    strSql.format (_T("delete from %s.t_acs_owner where f_owner_id = '%s';"),
                    dbName.getCStr(), dbOper->EscapeEx(userId).getCStr ());

    dbOper->Execute (strSql);

    NC_ACS_DB_TRACE (_T("this: %p, userId: %s end"), this, userId.getCStr ());
}

/* [notxpcom] bool IsOwner ([const] in StringRef docId, [const] in StringRef userId); */
NS_IMETHODIMP_(bool) ncDBOwnerManager::IsOwner(const String & docId, const String & userId)
{
    NC_ACS_DB_TRACE (_T("this: %p, docId: %s, userId: %s begin"),
        this, docId.getCStr (), userId.getCStr ());

    nsCOMPtr<ncIDBOperator> dbOper = getter_AddRefs (GetDBOperator ());

    String groupStr;
    int depth = ncGNSUtil::GetPathDepth (docId);
    if (depth == 0) {
        return false;
    }
    for (int i = 1; i <= depth; ++i) {
        String curDocId = ncGNSUtil::GetPathByDepth (docId, i);

        groupStr.append ("\'", 1);
        groupStr.append (dbOper->EscapeEx(curDocId));
        groupStr.append ("\'", 1);

        if (i != depth) {
            groupStr.append (",", 1);
        }
    }

    String dbName = Util::getDBName("anyshare");
    String strSql;
    strSql.format (_T("select f_owner_id from %s.t_acs_owner where f_owner_id = '%s' and f_gns_path in(%s);"),
                    dbName.getCStr(), dbOper->EscapeEx(userId).getCStr (), groupStr.getCStr ());

    ncDBRecords results;
    dbOper->Select (strSql, results);

    NC_ACS_DB_TRACE (_T("this: %p, docId: %s, userId: %s end"),
        this, docId.getCStr (), userId.getCStr ());

    if (results.size () == 0) {
        return false;
    }
    else {
        return true;
    }
}

/* [notxpcom] void GetSubObjsByUserId ([const] in StringRef docId, [const] in StringRef userId, in StringVecRef subObjs); */
NS_IMETHODIMP_(void) ncDBOwnerManager::GetSubObjsByUserId(const String & docId, const String & userId, vector<String> & subObjs)
{
    NC_ACS_DB_TRACE (_T("[BEGIN]this: %p, userId: %s"), this, userId.getCStr ());

    subObjs.clear ();

    nsCOMPtr<ncIDBOperator> dbOper = getter_AddRefs (GetDBOperator ());

    String docIdCond;
    if (ncGNSUtil::GetPathDepth (docId) != 0) {
        docIdCond.format ("and f_gns_path like '%s/%%'", dbOper->EscapeEx(docId).getCStr ());
    }

    String dbName = Util::getDBName("anyshare");
    String strSql;
    strSql.format (_T("select f_gns_path from %s.t_acs_owner where f_owner_id = '%s' %s order by length(f_gns_path) asc"),
                    dbName.getCStr(), dbOper->EscapeEx(userId).getCStr (), docIdCond.getCStr ());

    ncDBRecords results;
    dbOper->Select (strSql, results);

    for (size_t i = 0; i < results.size (); ++i) {
        subObjs.push_back (results[i][0]);
    }

    NC_ACS_DB_TRACE (_T("[END]this: %p, userId: %s, ret subObjs size: %d"), this, userId.getCStr (), (int)subObjs.size ());
}

String ncDBOwnerManager::GenerateGroupStr (const vector<String>& strs)
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

ncIDBOperator* ncDBOwnerManager::GetDBOperator ()
{
    nsCOMPtr<ncIDBOperator> dbOper = getter_AddRefs (ncACSDBGetDBOperator ());
    NS_ADDREF (dbOper.get ());
    return dbOper;
}
