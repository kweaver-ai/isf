/***************************************************************************************************
ncDBOutboxManager.cpp:
    Copyright (c) Eisoo Software Inc. (2021), All rights reserved.

Purpose:
    db outbox manager

Author:
    Dylan.gao@aishu.cn

Creating Time:
    2021-06-24
***************************************************************************************************/
#include <abprec.h>
#include <dataapi/ncJson.h>
#include <dboperator/public/ncIDBOperator.h>
#include <tedbcm/public/ncITEDBCManager.h>
#include <ncutil/ncBusinessDate.h>

#include "acsdb.h"
#include "ncDBOutboxManager.h"
#include "common/util.h"

/* Implementation file */
NS_IMPL_QUERY_INTERFACE1 (ncDBOutboxManager, ncIDBOutboxManager)

// protected
NS_IMETHODIMP_(nsrefcnt) ncDBOutboxManager::AddRef (void)
{
    return 1;
}

// protected
NS_IMETHODIMP_(nsrefcnt) ncDBOutboxManager::Release (void)
{
    return 1;
}

AB_DEFINE_THREADSAFE_SINGLETON_NO_POOL (ncDBOutboxManager)

ncDBOutboxManager::ncDBOutboxManager()
{
    NC_ACS_DB_TRACE (_T("this: %p"), this);

    nsresult ret;
    _nsqManager = do_CreateInstance (NSQ_CONTRACTID, &ret);
    if (NS_FAILED (ret)) {
        THROW_E (ACS_DB, FAILED_TO_CREATE_DB_OPERATOR_POOL_FACTORY,
            _T("Failed to create nsq instance: 0x%x"), ret);
    }

    _outboxToNSQMap.insert(pair<ncOutboxType, ncNSQEventType>(OUTBOX_CANCEL_APPLICATION, NSQ_DOC_SHARE_CANCEL));
    _outboxToNSQMap.insert(pair<ncOutboxType, ncNSQEventType>(OUTBOX_PERM_CHANGE, NSQ_SHARE_PERM_CHANGE));
    _outboxToNSQMap.insert(pair<ncOutboxType, ncNSQEventType>(OUTBOX_REALNAME_APPLY, NSQ_DOC_SHARE_REALNAME_APPLY));
    _outboxToNSQMap.insert(pair<ncOutboxType, ncNSQEventType>(OUTBOX_ANONYMOUS_APPLY, NSQ_DOC_SHARE_ANONYMOUS_APPLY));
}

ncDBOutboxManager::~ncDBOutboxManager()
{
    NC_ACS_DB_TRACE (_T("this: %p"), this);
}

void ncDBOutboxManager::sendNSQ (const String& message)
{
    JSON::Value messageJson;
    JSON::Reader::read (messageJson, message.getCStr (), message.getLength ());

    ncOutboxType type = static_cast<ncOutboxType>(messageJson["type"].i ());
    if(!_outboxToNSQMap.count(type)){
        throw Exception("invalid outboxType");
    }

    ncNSQEventType topic = _outboxToNSQMap[type];

    NSQMsg nsqMsg;
    // 申请id
    nsqMsg.applyId = messageJson["content"]["apply_id"].s ().c_str ();
    nsqMsg.userId = messageJson["content"]["user_id"].s ().c_str ();
    // 文档信息
    nsqMsg.docId = messageJson["content"]["doc_id"].s ().c_str ();
    nsqMsg.docPath = messageJson["content"]["doc_path"].s ().c_str ();
    nsqMsg.isFile = messageJson["content"]["is_file"].b ();
    nsqMsg.docLibType = messageJson["content"]["doc_lib_type"].i ();
    nsqMsg.docCsfLevel = messageJson["content"]["doc_csf_level"].i ();
    // 申请类型
    nsqMsg.applyType  = static_cast<ncNSQApplyType>(messageJson["content"]["apply_type"].i ());
    // 权限类型
    nsqMsg.operation = messageJson["content"]["operation"].i ();
    // 访问者信息
    nsqMsg.accessorId = messageJson["content"]["accessor_id"].s ().c_str ();
    nsqMsg.accessorType = messageJson["content"]["accessor_type"].i ();
    nsqMsg.accessorName = messageJson["content"]["accessor_name"].s ().c_str ();
    // 过期时间
    nsqMsg.expiresAt = messageJson["content"]["expires_at"].i ();
    // 处理权限值
    nsqMsg.allowValue = messageJson["content"]["allow_value"].i ();
    nsqMsg.denyValue = messageJson["content"]["deny_value"].i ();
    // 继承变更
    nsqMsg.inherit = messageJson["content"]["inherit"].b ();
    // 匿名共享
    nsqMsg.accessLimit = messageJson["content"]["access_limit"].i ();
    nsqMsg.linkId = messageJson["content"]["link_id"].s ().c_str ();
    nsqMsg.password = messageJson["content"]["password"].s ().c_str ();
    nsqMsg.title = messageJson["content"]["title"].s ().c_str ();
    // 匿名共享
    nsqMsg.conflictApplyId = messageJson["content"]["conflict_apply_id"].s ().c_str ();
    // 取消的申请ID集合
    vector<String> cancelApplyIds;
    JSON::Array cancelApplyIdsArr = messageJson["content"]["cancel_apply_ids"].a ();
    for (size_t i = 0; i < cancelApplyIdsArr.size(); ++i)
    {
        cancelApplyIds.push_back (cancelApplyIdsArr[i].s ().c_str ());
    }
    nsqMsg.cancelApplyIds = cancelApplyIds;

    // 发送NSQ消息
    _nsqManager->PublishNSQMessage (topic, nsqMsg);
}

/* [notxpcom] bool PushOutboxInfo (); */
NS_IMETHODIMP_(bool) ncDBOutboxManager::PushOutboxInfo()
{
    NC_ACS_DB_TRACE (_T("this: %p begin"), this);
    String dbName = Util::getDBName("anyshare");
    nsCOMPtr<ncIDBOperator> dbOper = getter_AddRefs (ncACSDBGetDBOperator ());
    try {
        dbOper->StartTransaction ();
        String strSql;
        // 使用 select ... for update 来加锁，防止事务处理时其他进程读取数据进行处理，间接实现了分布式锁
        strSql.format("SELECT f_id, f_message FROM %s.t_eacp_outbox ORDER BY f_create_time ASC LIMIT 1 For UPDATE", dbName.getCStr());

        ncDBRecords records;
        dbOper->Select(strSql, records);
        if (records.size() > 0) {
            SystemLog::getInstance ()->log (__FILE__, __LINE__, INFORMATION_LOG_TYPE,
                                            _T("Info : start push outbox message: f_id is %ld;")
                                            , Int64::getValue(records[0][0]));
            sendNSQ (records[0][1]);
            // 推送成功，删除相应记录
            strSql.format("DELETE FROM %s.t_eacp_outbox WHERE f_id = '%ld'", dbName.getCStr(), Int64::getValue(records[0][0]));
            dbOper->Execute(strSql);
            dbOper->Commit ();
            SystemLog::getInstance ()->log (__FILE__, __LINE__, INFORMATION_LOG_TYPE,
                                            _T("Info : push outbox message success"));
            return true;
        }else {
            dbOper->Commit ();
            SystemLog::getInstance ()->log (__FILE__, __LINE__, INFORMATION_LOG_TYPE, _T("Info : finish push all outbox message"));
            return false;
        }

    }
    catch(Exception&) {
        SystemLog::getInstance ()->log (__FILE__, __LINE__, INFORMATION_LOG_TYPE,_T("error has happened: rollback"));
        dbOper->Rollback ();
        throw;
    }

    NC_ACS_DB_TRACE (_T("this: %p end"), this);
}

/* [notxpcom] void AddOutboxInfo ([const] in StringRef message); */
NS_IMETHODIMP_(void) ncDBOutboxManager::AddOutboxInfo(const String& message)
{
    NC_ACS_DB_TRACE (_T("this: %p begin"), this);

    nsCOMPtr<ncIDBOperator> dbOper = getter_AddRefs (ncACSDBGetDBOperator ());

    String dbName = Util::getDBName("anyshare");
    String strSql;
    strSql.format (_T("INSERT INTO %s.t_eacp_outbox (f_message, f_create_time) values ('%s', %lld)"),
                      dbName.getCStr(), message.getCStr (), BusinessDate::getCurrentTime ());
    dbOper->Execute (strSql);

    NC_ACS_DB_TRACE (_T("this: %p end"), this);
}
