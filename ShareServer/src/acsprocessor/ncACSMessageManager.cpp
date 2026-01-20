#include <abprec.h>

#include "acsprocessor.h"
#include "ncACSMessageManager.h"
#include "ncMessageThread2.h"
#include "ncPushMessageThread2.h"

#define PLUGIN_MESSAGES_READ  _T("message/v1/messages-read")

/* Implementation file */
NS_IMPL_QUERY_INTERFACE1(ncACSMessageManager, ncIACSMessageManager)

// protected
NS_IMETHODIMP_(nsrefcnt) ncACSMessageManager::AddRef (void)
{
    return 1;
}

// protected
NS_IMETHODIMP_(nsrefcnt) ncACSMessageManager::Release (void)
{
    return 1;
}

AB_DEFINE_THREADSAFE_SINGLETON_NO_POOL (ncACSMessageManager)

ncACSMessageManager::ncACSMessageManager()
    : _messageThread2 (NULL),
      _pushMessageThread2 (NULL)
{
    NC_ACS_PROCESSOR_TRACE (_T("this: %p"), this);

    nsresult ret;
    _dbMessageManager = do_CreateInstance (NC_DB_MESSAGE_MANAGER_CONTRACTID, &ret);
    if (NS_FAILED (ret)) {
        THROW_E (ACS_PROCESSOR, FAILED_TO_CREATE_DB_MESSAGE_MANANGER,
            _T("Failed to create db message manager: 0x%x"), ret);
    }

    _dbConfManager = do_CreateInstance (NC_DB_CONF_MANAGER_CONTRACTID, &ret);
    if (NS_FAILED (ret)) {
        THROW_E (ACS_PROCESSOR, FAILED_TO_CREATE_DB_CONF_MANANGER,
            _T("Failed to create db conf manager: 0x%x"), ret);
    }

    _acsConfManager = do_CreateInstance(NC_ACS_CONF_MANAGER_CONTRACTID, &ret);
    if (NS_FAILED (ret)) {
        THROW_E (ACS_PROCESSOR, FAILED_TO_CREATE_ACS_CONF_MANANGER,
            _T("Failed to create acs conf manager: 0x%x"), ret);
    }

    _userManagement = do_CreateInstance (USER_MANAGEMENT_CONTRACTID, &ret);
    if (NS_FAILED (ret)) {
        THROW_E (ACS_PROCESSOR, FAILED_TO_CREATE_DRIVENADAPTER_MANANGER,
            _T("Failed to create userManagement instance: 0x%x"), ret);
    }
}

ncACSMessageManager::~ncACSMessageManager()
{
    NC_ACS_PROCESSOR_TRACE (_T("this: %p"), this);
}

/* [notxpcom] void StartMessageThread2 (); */
NS_IMETHODIMP_(void) ncACSMessageManager::StartMessageThread2()
{
    NC_ACS_PROCESSOR_TRACE (_T("this: %p, begin"), this);

    if (_messageThread2 != NULL) {
        return;
    }

    static ThreadMutexLock mutex2;
    AutoLock<ThreadMutexLock> lock (&mutex2);

    if (_messageThread2 == NULL) {
        _messageThread2 = new ncMessageThread2 ();
        _messageThread2->start ();
    }

    NC_ACS_PROCESSOR_TRACE (_T("this: %p, end"), this);
}

/* [notxpcom] void StartPushMessageThread2 (); */
NS_IMETHODIMP_(void) ncACSMessageManager::StartPushMessageThread2()
{
    NC_ACS_PROCESSOR_TRACE (_T("this: %p, begin"), this);

    if (_pushMessageThread2 != NULL) {
        return;
    }

    static ThreadMutexLock pushMessageMutex2;
    AutoLock<ThreadMutexLock> lock (&pushMessageMutex2);

    if (_pushMessageThread2 == NULL) {
        _pushMessageThread2 = new ncPushMessageThread2 ();
        _pushMessageThread2->start ();
    }

    NC_ACS_PROCESSOR_TRACE (_T("this: %p, end"), this);
}

/* [notxpcom] void AddMessage2 ([const] in acsMessagePtrVecRef msgs); */
NS_IMETHODIMP_(void) ncACSMessageManager::AddMessage2(const vector<std::shared_ptr<acsMessage>> & msgs)
{
    NC_ACS_PROCESSOR_TRACE (_T("this: %p, begin"), this);

    // 设置所有消息的ID
    for (int i=0; i<msgs.size(); ++i) {
        Guid msgid;
        msgs[i]->msgId = std::move (msgid.toString ());
    }

    // 客户端消息
    StartMessageThread2 ();
    _messageThread2->AddMessage (msgs);

    // 插件消息
    StartPushMessageThread2 ();
    _pushMessageThread2->AddMessage (msgs);

    NC_ACS_PROCESSOR_TRACE (_T("this: %p, end"), this);
}

/* [notxpcom] void AddPluginMessage ([const] in acsMessagePtrVecRef msgs); */
NS_IMETHODIMP_(void) ncACSMessageManager::AddPluginMessage(const vector<std::shared_ptr<acsMessage>> & msgs)
{
    NC_ACS_PROCESSOR_TRACE (_T("this: %p, begin"), this);

    // 设置所有消息的ID
    for (int i=0; i<msgs.size(); ++i) {
        Guid msgid;
        msgs[i]->msgId = std::move (msgid.toString ());
    }

    // 插件消息
    StartPushMessageThread2 ();
    _pushMessageThread2->AddMessage (msgs);

    NC_ACS_PROCESSOR_TRACE (_T("this: %p, end"), this);
}

/* [notxpcom] void GetMessageByUserId ([const] in StringRef userid, in int64 stamp, in acsMessageResultVecRef msgs); */
NS_IMETHODIMP_(void) ncACSMessageManager::GetMessageByUserId(const String & userid, int64 stamp, vector<acsMessageResult> & msgs)
{
    NC_ACS_PROCESSOR_TRACE (_T("this: %p, userId: %s begin"), this, userid.getCStr ());

    std::vector<dbMessageResult> result;
    _dbMessageManager->GetMessageByUserId (userid, stamp, 300, result);

    if (!result.empty ()) {
        msgs.reserve (result.size ());

        acsMessageResult tmp;
        for (size_t i = 0; i < result.size (); ++i) {
            tmp.msgId = result[i].msgId;
            tmp.msgContent = std::move (result[i].msgContent);
            tmp.msgStamp = result[i].msgStamp;
            tmp.msgStatus = result[i].msgStatus;

            msgs.push_back (std::move (tmp));
        }
    }

    NC_ACS_PROCESSOR_TRACE (_T("this: %p, userId: %s end"), this, userid.getCStr ());
}

/* [notxpcom] void ReadMessageByUserId ([const] in StringRef userid, in int64 stamp); */
NS_IMETHODIMP_(void) ncACSMessageManager::ReadMessageByUserId(const String & userid, int64 stamp)
{
    _dbMessageManager->ReadMessageByUserId (userid, stamp);
}

/* [notxpcom] void ReadMessageByIds ([const] in StringRef userid, [const] in StringVecRef msgids); */
NS_IMETHODIMP_(void) ncACSMessageManager::ReadMessageByIds(const String & userid, const vector<String> & msgids)
{
    _dbMessageManager->ReadMessageByIds (userid, msgids);

    // 发送"消息已读"插件消息
    std::vector<std::shared_ptr<acsMessage>> messages;
    std::shared_ptr<acsMessage> msgptr(new acsMessage());
    msgptr->channel = PLUGIN_MESSAGES_READ;

    UserInfo userInfo;
    _userManagement->GetUserInfo (userid, userInfo);
    messageReceiver receiver;
    receiver.id = userInfo.id;
    receiver.account = userInfo.account;
    receiver.name = userInfo.name;
    receiver.email = userInfo.email;
    receiver.telephone = userInfo.telephone;
    receiver.thirdAttr = userInfo.thirdAttr;
    receiver.thirdId = userInfo.thirdId;

    vector<messageReceiver> receivers;
    receivers.push_back (std::move (receiver));
    msgptr->receivers = receivers;

    std::string contstr;
    JSON::Object payload;
    JSON::Array msgIdArr;
    for (auto It = msgids.begin (); It != msgids.end (); ++It) {
        msgIdArr.push_back((*It).getCStr());
    }
    payload["msg_ids"] = msgIdArr;
    JSON::Writer::write (payload, contstr);
    msgptr->content = std::move (toCFLString (contstr));

    messages.push_back (msgptr);
    AddPluginMessage (messages);
}

/* [notxpcom] void ReadMessageByTaskId ([const] in StringRef userid, [const] in StringRef taskId); */
NS_IMETHODIMP_(void) ncACSMessageManager::ReadMessageByTaskId(const String & userid, const String & taskId)
{
    _dbMessageManager->ReadMessageByTaskId (userid, taskId);
}

/* [notxpcom] void BatchReadMessageByTaskId ([const] in StringRef taskId); */
NS_IMETHODIMP_(void) ncACSMessageManager::BatchReadMessageByTaskId(const String & taskId)
{
    _dbMessageManager->BatchReadMessageByTaskId (taskId);
}

/* [notxpcom] void ReadMessageForAllReceivers ([const] in StringRef msgId); */
NS_IMETHODIMP_(void) ncACSMessageManager::ReadMessageForAllReceivers(const String & msgId)
{
    _dbMessageManager->ReadMessageByMsgId (msgId);
}

/* [notxpcom] void ReadMessageForSomeReceivers ([const] in StringRef msgId, [const] in StringVecRef receiverIds); */
NS_IMETHODIMP_(void) ncACSMessageManager::ReadMessageForSomeReceivers(const String & msgId, const vector<String> & receiverIds)
{
    _dbMessageManager->ReadMessageByMsgIdAndReceiverIds (msgId, receiverIds);
}
