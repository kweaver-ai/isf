#include <abprec.h>
#include <ncutil/ncBusinessDate.h>

#include <dboperator/public/ncIDBOperator.h>
#include <tedbcm/public/ncITEDBCManager.h>

#include "acsdb.h"
#include "ncDBMessageManager.h"
#include "common/util.h"

/* Implementation file */
NS_IMPL_QUERY_INTERFACE1 (ncDBMessageManager, ncIDBMessageManager)

// protected
NS_IMETHODIMP_(nsrefcnt) ncDBMessageManager::AddRef (void)
{
    return 1;
}

// protected
NS_IMETHODIMP_(nsrefcnt) ncDBMessageManager::Release (void)
{
    return 1;
}

AB_DEFINE_THREADSAFE_SINGLETON_NO_POOL (ncDBMessageManager)

ncDBMessageManager::ncDBMessageManager()
{
    NC_ACS_DB_TRACE (_T("this: %p"), this);
}

ncDBMessageManager::~ncDBMessageManager()
{
    NC_ACS_DB_TRACE (_T("this: %p"), this);
}

/* [notxpcom] void AddMessage ([const] in StringRef msgid, [const] in StringRef content, [const] in StringVecRef tousers, [const] in StringRef taskId); */
NS_IMETHODIMP_(void) ncDBMessageManager::AddMessage(const String & msgid, const String & content, const vector<String> & tousers, const String & taskId)
{
    NC_ACS_DB_TRACE (_T("this: %p begin"), this);

    nsCOMPtr<ncIDBOperator> dbOper = getter_AddRefs (GetDBOperator ());

    String dbName = Util::getDBName("anyshare");
    try {
        dbOper->StartTransaction ();

        String strSql;
        int64 stamp = BusinessDate::getCurrentTime ();
        strSql.format (_T("insert into %s.t_message values('%s', '%s', %lld, '%s')"),
                            dbName.getCStr(),
                            dbOper->EscapeEx(msgid).getCStr (),
                            dbOper->EscapeEx(content).getCStr (),
                            stamp,
                            dbOper->EscapeEx(taskId).getCStr());
        dbOper->Execute (strSql);

        for (size_t i = 0; i < tousers.size (); ++i) {
            strSql.format (_T("insert into %s.t_message_usermap values('%s', '%s', 0, 0)"),
                dbName.getCStr(), dbOper->EscapeEx(msgid).getCStr (), dbOper->EscapeEx(tousers[i]).getCStr ());
            dbOper->Execute (strSql);
        }
        dbOper->Commit ();
    }
    catch (...) {
        dbOper->Rollback ();
    }

    NC_ACS_DB_TRACE (_T("this: %p end"), this);
}

/* [notxpcom] void GetMessageByUserId ([const] in StringRef userid, in int64 stamp, in int limit, in dbMessageResultVecRef msgs); */
NS_IMETHODIMP_(void) ncDBMessageManager::GetMessageByUserId(const String & userid, int64 stamp, int limit, vector<dbMessageResult> & msgs)
{
    NC_ACS_DB_TRACE (_T("this: %p, userid: %s, limit:%d begin"), this, userid.getCStr (), limit);

    nsCOMPtr<ncIDBOperator> dbOper = getter_AddRefs (GetDBOperator ());
    msgs.clear ();

    String dbName = Util::getDBName("anyshare");
    String strSql;
    strSql.format (_T("select f_content, f_create_stamp, f_status, f_msg_id, f_channel from %s.t_message_usermap ")
                   _T("join %s.t_message using (f_msg_id) ")
                   _T("where f_user_id = '%s' and (f_create_stamp > %lld or f_read_stamp > %lld)")
                   _T("order by f_create_stamp desc ")
                   _T("limit %d"),
                   dbName.getCStr(), dbName.getCStr(), dbOper->EscapeEx(userid).getCStr (), stamp, stamp, limit);
    ncDBRecords results;
    dbOper->Select (strSql, results);

    if (!results.empty ()) {
        msgs.reserve (results.size());

        dbMessageResult tmpInfo;
        for (int i = int(results.size ()) - 1; i >= 0; --i) {
            tmpInfo.msgContent = results[i][0];
            tmpInfo.msgStamp = Int64::getValue (results[i][1]);
            tmpInfo.msgStatus = Int::getValue (results[i][2]);
            tmpInfo.msgId = results[i][3];

            String strChannel = results[i][4];
            if (strChannel == "" || strChannel == "doc-share/v1/share-with-users-on" || strChannel == "doc-share/v1/share-with-users-off" ||
                strChannel == "doc-share/v1/owner-set" || strChannel == "doc-share/v1/owner-removed" ||
                strChannel == "document/v1/moved-to-quarantine" || strChannel == "document/v1/restored-from-quarantine") {
                    msgs.push_back (tmpInfo);
                }
        }
    }

    NC_ACS_DB_TRACE (_T("this: %p, userid: %s, limit:%d end"), this, userid.getCStr (), limit);
}

/* [notxpcom] void ReadMessageByUserId ([const] in StringRef userid, in int64 stamp); */
NS_IMETHODIMP_(void) ncDBMessageManager::ReadMessageByUserId(const String & userid, int64 stamp)
{
    NC_ACS_DB_TRACE (_T("this: %p, userid: %s, stamp:%lld begin"), this, userid.getCStr (), stamp);

    nsCOMPtr<ncIDBOperator> dbOper = getter_AddRefs (GetDBOperator ());

    String dbName = Util::getDBName("anyshare");
    String strSql;
    strSql.format (_T("delete t_message_usermap from %s.t_message_usermap ")
                   _T("join %s.t_message using (f_msg_id) ")
                   _T("where f_user_id = '%s' and f_create_stamp <= %lld"),
                   dbName.getCStr(), dbName.getCStr(), dbOper->EscapeEx(userid).getCStr (), stamp);
    dbOper->Execute (strSql);

    NC_ACS_DB_TRACE (_T("this: %p, userid: %s, stamp:%lld end"), this, userid.getCStr (), stamp);
}

/* [notxpcom] void ReadMessageByIds ([const] in StringRef userid, [const] in StringVecRef msgids); */
NS_IMETHODIMP_(void) ncDBMessageManager::ReadMessageByIds(const String & userid, const vector<String> & msgids)
{
    NC_ACS_DB_TRACE (_T("this: %p, userid: %s, msgids size: %d begin"), this, userid.getCStr (), (int)msgids.size ());


    if (msgids.size () > 0) {
        nsCOMPtr<ncIDBOperator> dbOper = getter_AddRefs (GetDBOperator ());

        String groupStr;
        for (size_t i = 0; i < msgids.size (); ++i) {
            groupStr.append ("\'", 1);
            groupStr.append (dbOper->EscapeEx(msgids[i]));
            groupStr.append ("\'", 1);
            groupStr.append (",", 1);
        }
        groupStr.remove(groupStr.getLength() - 1);

        String dbName = Util::getDBName("anyshare");
        String strSql;
        strSql.format (_T("update %s.t_message_usermap set f_status=1, f_read_stamp = %lld ")
                       _T("where f_user_id = '%s' and f_msg_id in (%s)"),
                       dbName.getCStr(), BusinessDate::getCurrentTime (), dbOper->EscapeEx(userid).getCStr (), groupStr.getCStr ());
        dbOper->Execute (strSql);
    }
    NC_ACS_DB_TRACE (_T("this: %p, userid: %s, msgids size: %d end"), this, userid.getCStr (), (int)msgids.size ());

}

/* [notxpcom] void ReadMessageByTaskId ([const] in StringRef userid, [const] in StringRef taskId); */
NS_IMETHODIMP_(void) ncDBMessageManager::ReadMessageByTaskId(const String & userid, const String & taskId)
{
    NC_ACS_DB_TRACE (_T("this: %p, userid: %s, taskId:%s begin"), this, userid.getCStr (), taskId.getCStr());

    nsCOMPtr<ncIDBOperator> dbOper = getter_AddRefs (GetDBOperator ());

    String dbName = Util::getDBName("anyshare");
    String strSql;
    strSql.format (_T("select f_msg_id from %s.t_message ")
                   _T("join %s.t_message_usermap using (f_msg_id) ")
                   _T("where f_task_id = '%s' and f_user_id = '%s'"),
                   dbName.getCStr(), dbName.getCStr(), taskId.getCStr (), dbOper->EscapeEx(userid).getCStr ());
    ncDBRecords results;
    dbOper->Select (strSql, results);

    if (!results.empty ()) {
        for (size_t i = 0; i < results.size(); i++) {
            strSql.format (_T("update %s.t_message_usermap set f_status=1, f_read_stamp = %lld ")
               _T("where f_user_id = '%s' and f_msg_id = '%s' "),
               dbName.getCStr(), BusinessDate::getCurrentTime (), dbOper->EscapeEx(userid).getCStr (), dbOper->EscapeEx(results[i][0]).getCStr());
            dbOper->Execute (strSql);
        }
    }

    NC_ACS_DB_TRACE (_T("this: %p, userid: %s, taskId:%s end"), this, userid.getCStr (), taskId.getCStr());
}

/* [notxpcom] void BatchReadMessageByTaskId ([const] in StringRef taskId); */
NS_IMETHODIMP_(void) ncDBMessageManager::BatchReadMessageByTaskId(const String & taskId)
{
    NC_ACS_DB_TRACE (_T("this: %p, taskId:%s begin"), this, taskId.getCStr());

    nsCOMPtr<ncIDBOperator> dbOper = getter_AddRefs (GetDBOperator ());

    String dbName = Util::getDBName("anyshare");
    String strSql;
    strSql.format (_T("select f_msg_id from %s.t_message where f_task_id = '%s'"),
                   dbName.getCStr(), taskId.getCStr ());
    ncDBRecords results;
    dbOper->Select (strSql, results);

    if (!results.empty ()) {
        for (size_t i = 0; i < results.size(); i++) {
            strSql.format (_T("update %s.t_message_usermap set f_status=1, f_read_stamp = %lld where f_msg_id = '%s'"),
               dbName.getCStr(), BusinessDate::getCurrentTime (), dbOper->EscapeEx(results[i][0]).getCStr());
            dbOper->Execute (strSql);
        }
    }

    NC_ACS_DB_TRACE (_T("this: %p, taskId:%s end"), this, taskId.getCStr());
}

/* [notxpcom] void ReadMessageByMsgId ([const] in StringRef msgId); */
NS_IMETHODIMP_(void) ncDBMessageManager::ReadMessageByMsgId(const String & msgId)
{
    NC_ACS_DB_TRACE (_T("this: %p, msgId:%s begin"), this, msgId.getCStr());

    String dbName = Util::getDBName("anyshare");
    nsCOMPtr<ncIDBOperator> dbOper = getter_AddRefs (GetDBOperator ());
    String strSql;
    strSql.format (_T("update %s.t_message_usermap set f_status=1, f_read_stamp = %lld where f_msg_id = '%s'"),
        dbName.getCStr(), BusinessDate::getCurrentTime (), dbOper->EscapeEx(msgId).getCStr());
    dbOper->Execute (strSql);

    NC_ACS_DB_TRACE (_T("this: %p, msgId:%s end"), this, msgId.getCStr());
}

/* [notxpcom] void ReadMessageByMsgIdAndReceiverIds ([const] in StringRef msgId, [const] in StringVecRef receiverIds); */
NS_IMETHODIMP_(void) ncDBMessageManager::ReadMessageByMsgIdAndReceiverIds(const String & msgId, const vector<String> & receiverIds)
{
    NC_ACS_DB_TRACE (_T("this: %p, msgId: %s begin"), this, msgId.getCStr ());

    if (receiverIds.size() == 0) {
        return;
    }
    String dbName = Util::getDBName("anyshare");
    nsCOMPtr<ncIDBOperator> dbOper = getter_AddRefs (GetDBOperator ());
    String groupStr;
    for (size_t i = 0; i < receiverIds.size (); ++i) {
        groupStr.append ("\'", 1);
        groupStr.append (dbOper->EscapeEx(receiverIds[i]));
        groupStr.append ("\'", 1);
        groupStr.append (",", 1);
    }
    groupStr.remove(groupStr.getLength() - 1);

    String strSql;
    strSql.format (_T("update %s.t_message_usermap set f_status=1, f_read_stamp = %lld ")
                    _T("where f_msg_id = '%s' and f_user_id in (%s)"),
                    dbName.getCStr(), BusinessDate::getCurrentTime (), dbOper->EscapeEx(msgId).getCStr (), groupStr.getCStr ());
    dbOper->Execute (strSql);

    NC_ACS_DB_TRACE (_T("this: %p, msgId: %s end"), this, msgId.getCStr ());
}

/* [notxpcom] void DelMessage (in int64 stamp); */
NS_IMETHODIMP_(void) ncDBMessageManager::DelMessage(int64 stamp)
{
    NC_ACS_DB_TRACE (_T("this: %p, stamp:%lld begin"), this, stamp);

    nsCOMPtr<ncIDBOperator> dbOper = getter_AddRefs (GetDBOperator ());

    String dbName = Util::getDBName("anyshare");
    String strSql;
    strSql.format (_T("delete from %s.t_message where f_create_stamp <= %lld"), dbName.getCStr(), stamp);
    dbOper->Execute (strSql);

    NC_ACS_DB_TRACE (_T("this: %p, stamp:%lld end"), this, stamp);
}

ncIDBOperator* ncDBMessageManager::GetDBOperator ()
{
    nsCOMPtr<ncIDBOperator> dbOper = getter_AddRefs (ncACSDBGetDBOperator ());
    NS_ADDREF (dbOper.get ());
    return dbOper;
}
