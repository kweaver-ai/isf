#include "eachttpserver.h"
#include "ncEACContactorHandler.h"
#include "ncEACHttpServerUtil.h"
#include <ehttpserver/ncEHttpUtil.h>
#include <ethriftutil/ncThriftClient.h>
#include "eacServiceAccessConfig.h"
// public
ncEACContactorHandler::ncEACContactorHandler (ncIACSShareMgnt* acsShareMgnt)
    : _acsShareMgnt (acsShareMgnt)
{
    NC_EAC_HTTP_SERVER_TRACE (_T("this: %p"), this);

    if (_acsShareMgnt == NULL) {
        nsresult ret;
        _acsShareMgnt = do_CreateInstance (NC_ACS_SHARE_MGNT_CONTRACTID, &ret);
        if (NS_FAILED (ret)) {
            THROW_E(EAC_HTTP_SERVER, EACHTTP_SERVER_INIT_ACS_SHAREMGNT_ERR,
                LOAD_STRING (_T("IDS_EACHTTP_ACS_SHAREMGNT_INIT_ERROR")), ret);
        }
    }
    _methodFuncs.insert (pair<String, ncMethodFunc>(_T("get"), &ncEACContactorHandler::GetContactors));
    _methodFuncs.insert (pair<String, ncMethodFunc>(_T("search"), &ncEACContactorHandler::Search));
    _methodFuncs.insert (pair<String, ncMethodFunc>(_T("searchcount"), &ncEACContactorHandler::SearchCount));
    _methodFuncs.insert (pair<String, ncMethodFunc>(_T("addgroup"), &ncEACContactorHandler::AddGroup));
    _methodFuncs.insert (pair<String, ncMethodFunc>(_T("editgroup"), &ncEACContactorHandler::EditGroup));
    _methodFuncs.insert (pair<String, ncMethodFunc>(_T("getgroups"), &ncEACContactorHandler::GetGroup));
    _methodFuncs.insert (pair<String, ncMethodFunc>(_T("addpersons"), &ncEACContactorHandler::AddPersons));
    _methodFuncs.insert (pair<String, ncMethodFunc>(_T("searchpersons"), &ncEACContactorHandler::SearchPersons));
    _methodFuncs.insert (pair<String, ncMethodFunc>(_T("deletepersons"), &ncEACContactorHandler::DeletePersons));
    _methodFuncs.insert (pair<String, ncMethodFunc>(_T("getpersons"), &ncEACContactorHandler::GetPersonFromGroup));

}

// public
ncEACContactorHandler::~ncEACContactorHandler (void)
{
    NC_EAC_HTTP_SERVER_TRACE (_T("this: %p"), this);
}

// public
void
ncEACContactorHandler::doContactorRequestHandler (brpc::Controller* cntl)
{
    NC_EAC_HTTP_SERVER_TRY

        // 获取query string
        String method;
        String tokenId;
        String userId;
        ncHttpGetParams (cntl, method, tokenId);
        // method是否设置
        if (method.isEmpty ()) {
            THROW_E (EAC_HTTP_SERVER, EACHTTP_QUERY_METHOD_INVALID,
                     LOAD_STRING (_T("IDS_EACHTTP_QUERY_METHOD_INVALID")));
        }

        // method是否支持
        map<String, ncMethodFunc>::iterator iter = _methodFuncs.find (method);
        if (iter == _methodFuncs.end ()) {
            THROW_E (EAC_HTTP_SERVER, EACHTTP_QUERY_METHOD_INVALID,
                LOAD_STRING (_T("IDS_EACHTTP_QUERY_METHOD_NOT_SUPPORT")));
        }

        // token验证
        ncCheckTokenInfo checkTokenInfo;
        checkTokenInfo.tokenId = tokenId;
        checkTokenInfo.ip = ncEACHttpServerUtil::GetForwardedIp(cntl);
        ncIntrospectInfo introspectInfo;
        if (CheckToken (checkTokenInfo, introspectInfo) == false) {
            THROW_E (EAC_HTTP_SERVER, EACHTTP_SERVER_ACCESS_TOKEN_ERR,
                    LOAD_STRING (_T("IDS_EACHTTP_ERR_MSG_ACCESS_TOKEN_ERR")));
        }

        // 获取该token对应的userId
        userId = introspectInfo.userId;

        // 消息处理
        ncMethodFunc func = iter->second;
        (this->*func) (cntl, userId);

    NC_EAC_HTTP_SERVER_CATCH (cntl)
}

// protected
void
ncEACContactorHandler::GetGroup (brpc::Controller* cntl, const String& userId)
{
    NC_EAC_HTTP_SERVER_TRACE (_T("this: %p, cntl: %p, userId: %s begin"), this, cntl, userId.getCStr ());

    vector<ncTPersonGroup> groupInfos;
    try {
        ncThriftClient<ncTShareMgntClient> shareMgntClient (EacServiceAccessConfig::getInstance()->sharemgntHost, EacServiceAccessConfig::getInstance()->sharemgntPort);
        shareMgntClient->Usrm_GetPersonGroups (groupInfos, toSTLString(userId));
    }
    catch (ncTException & e) {
        HandlencTException(e);
    }
    catch (TException & e) {
        THROW_E (EAC_HTTP_SERVER, EACHTTP_SERVER_UNKNOWN_ERR, e.what ());
    }

    // 回复
    JSON::Object replyJson;
    JSON::Array& groupsJson = replyJson["groups"].a ();

    for (size_t i = 0; i < groupInfos.size (); ++i) {
        groupsJson.push_back (JSON::OBJECT);

        JSON::Object& tmpObj = groupsJson.back ().o ();
        tmpObj["id"] = groupInfos[i].groupId;
        tmpObj["groupname"] = groupInfos[i].groupName;
        tmpObj["count"] = groupInfos[i].personCount;
    }

    string body;
    JSON::Writer::write (replyJson, body);
    ncHttpSendReply (cntl, brpc::HTTP_STATUS_OK, "ok", body);

    NC_EAC_HTTP_SERVER_TRACE (_T("this: %p, cntl: %p, userId: %s end"), this, cntl, userId.getCStr ());
}

// protected
void
ncEACContactorHandler::GetContactors (brpc::Controller* cntl, const String& userId)
{
    NC_EAC_HTTP_SERVER_TRACE (_T("this: %p, cntl: %p, userId: %s begin"), this, cntl, userId.getCStr ());

    // 获取http请求的content
    string bodyBuffer = cntl->request_attachment ().to_string ();

    // 获取docid,name
    JSON::Value requestJson;
    try {
        JSON::Reader::read (requestJson, bodyBuffer.c_str (), bodyBuffer.size ());
    }
    catch (Exception& e) {
        THROW_E (EAC_HTTP_SERVER, EACHTTP_INVALID_REQUEST_JSON_FORMAT,
                 LOAD_STRING (_T("IDS_EACHTTP_INVALID_JSON")));
    }

    String groupId = requestJson["groupid"].s ().c_str ();
    INVALID_GROUP_ID(groupId);

    vector<ncACSUserInfo> userInfos;
    _acsShareMgnt->GetContactors (groupId, userInfos);

    // 回复
    JSON::Value replyJson;
    JSON::Array& usersJson = replyJson["userinfos"].a ();

    for (size_t i = 0; i < userInfos.size (); ++i) {
        // 用户被禁用时界面上不显示
        if (userInfos[i].enableStatus == false) {
            continue;
        }

        usersJson.push_back (JSON::OBJECT);
        JSON::Object& tmpObj = usersJson.back ().o ();

        tmpObj["userid"] = userInfos[i].id.getCStr ();
        tmpObj["account"] = userInfos[i].account.getCStr ();
        tmpObj["name"] = userInfos[i].visionName.getCStr ();
        tmpObj["mail"] = userInfos[i].email.getCStr ();
        tmpObj["csflevel"] = userInfos[i].csfLevel;
    }

    string body;
    JSON::Writer::write (replyJson.o (), body);
    ncHttpSendReply (cntl, brpc::HTTP_STATUS_OK, "ok", body);

    NC_EAC_HTTP_SERVER_TRACE (_T("this: %p, cntl: %p, userId: %s end"), this, cntl, userId.getCStr ());

}

// protected
void
ncEACContactorHandler::Search (brpc::Controller* cntl, const String& userId)
{
    NC_EAC_HTTP_SERVER_TRACE (_T("this: %p, cntl: %p, userId: %s begin"), this, cntl, userId.getCStr ());

    // 获取http请求的content
    string bodyBuffer = cntl->request_attachment ().to_string ();

    JSON::Value requestJson;
    try {
        JSON::Reader::read (requestJson, bodyBuffer.c_str (), bodyBuffer.size ());
    }
    catch (Exception& e) {
        THROW_E (EAC_HTTP_SERVER, EACHTTP_INVALID_REQUEST_JSON_FORMAT,
            LOAD_STRING (_T("IDS_EACHTTP_INVALID_JSON")));
    }

    String key = requestJson["key"].s ().c_str ();
    if (key.isEmpty ()) {
        THROW_E (EAC_HTTP_SERVER, EACHTTP_INVALID_SEARCH_KEY,
            LOAD_STRING (_T("IDS_EACHTTP_INVALID_SEARCH_KEY")));
    }

    int start;
    if (requestJson["start"].type() == JSON::NIL) {
        start = 0;
    } else {
        start = requestJson["start"].i();
        if(start < 0) {
            THROW_E (EAC_HTTP_SERVER, INVALID_PAGE_START_VALUE,
                "invalid start value, must bigger than 0");
        }
    }

    int limit;
    if (requestJson["limit"].type() == JSON::NIL) {
        limit = 10;
    } else {
        limit = requestJson["limit"].i();
        if(limit < -1) {
            THROW_E (EAC_HTTP_SERVER, INVALID_PAGE_LIMIT_VALUE,
                "invalid limit value, must bigger than -1");
        }
    }

    // 进行搜索
    vector<ncACSUserInfo> userInfos;
    vector<ncGroupInfo> groupInfos;
    _acsShareMgnt->SearchContactGroup (userId, key, start, limit, userInfos, groupInfos);

    // 获取用户的直属联系组信息
    vector<String> userIds;
    for(size_t i = 0; i < userInfos.size (); ++i) {
        userIds.push_back (userInfos[i].id);
    }

    map<String, ncGroupInfo> infoMap;
    _acsShareMgnt->GetBelongGroupByIdBatch (userId, userIds, infoMap);

    // 回复
    JSON::Value replyJson;
    JSON::Array& userJson = replyJson["userinfos"].a ();
    for (size_t i = 0; i < userInfos.size (); ++i) {
        userJson.push_back (JSON::OBJECT);
        JSON::Object& tmpObj = userJson.back ().o ();
        tmpObj["userid"] = userInfos[i].id.getCStr ();
        tmpObj["account"] = userInfos[i].account.getCStr ();
        tmpObj["name"] = userInfos[i].visionName.getCStr ();
        tmpObj["mail"] = userInfos[i].email.getCStr ();
        tmpObj["csflevel"] = userInfos[i].csfLevel;

        ncGroupInfo info = infoMap[userInfos[i].id];
        tmpObj["groupid"] = info.id.getCStr ();
        tmpObj["groupname"] = info.groupName.getCStr ();
    }

    JSON::Array& groupsJson = replyJson["groups"].a ();
    for (size_t i = 0; i < groupInfos.size (); ++i) {
        groupsJson.push_back (JSON::OBJECT);

        JSON::Object& tmpObj = groupsJson.back ().o ();
        tmpObj["id"] = groupInfos[i].id.getCStr ();
        tmpObj["createrid"] = groupInfos[i].createrId.getCStr ();
        tmpObj["groupname"] = groupInfos[i].groupName.getCStr ();
        tmpObj["count"] = groupInfos[i].count;
    }

    string body;
    JSON::Writer::write (replyJson.o (), body);

    ncHttpSendReply (cntl, brpc::HTTP_STATUS_OK, "ok", body);

    NC_EAC_HTTP_SERVER_TRACE (_T("this: %p, cntl: %p, userId: %s end"), this, cntl, userId.getCStr ());

}

// protected
void
ncEACContactorHandler::SearchCount (brpc::Controller* cntl, const String& userId)
{
    NC_EAC_HTTP_SERVER_TRACE (_T("this: %p, cntl: %p, userId: %s begin"), this, cntl, userId.getCStr ());

    // 获取http请求的content
    string bodyBuffer = cntl->request_attachment ().to_string ();

    JSON::Value requestJson;
    try {
        JSON::Reader::read (requestJson, bodyBuffer.c_str (), bodyBuffer.size ());
    }
    catch (Exception& e) {
        THROW_E (EAC_HTTP_SERVER, EACHTTP_INVALID_REQUEST_JSON_FORMAT,
            LOAD_STRING (_T("IDS_EACHTTP_INVALID_JSON")));
    }

    String key = requestJson["key"].s ().c_str ();
    if (key.isEmpty ()) {
        THROW_E (EAC_HTTP_SERVER, EACHTTP_INVALID_SEARCH_KEY,
            LOAD_STRING (_T("IDS_EACHTTP_INVALID_SEARCH_KEY")));
    }

    // 获取搜索数目
    int count = _acsShareMgnt->SearchContactGroupCount (userId, key, true);

    // 回复
    JSON::Value replyJson;
    replyJson["count"] = count;

    string body;
    JSON::Writer::write (replyJson.o (), body);

    ncHttpSendReply (cntl, brpc::HTTP_STATUS_OK, "ok", body);

    NC_EAC_HTTP_SERVER_TRACE (_T("this: %p, cntl: %p, userId: %s end"), this, cntl, userId.getCStr ());

}

// protected
void
ncEACContactorHandler::AddGroup (brpc::Controller* cntl, const String& userId)
{

    NC_EAC_HTTP_SERVER_TRACE (_T("this: %p, cntl: %p, userId: %s begin"), this, cntl, userId.getCStr ());

    // 获取http请求的content
    string bodyBuffer = cntl->request_attachment ().to_string ();
    JSON::Value requestJson;
    try {
        JSON::Reader::read (requestJson, bodyBuffer.c_str (), bodyBuffer.size ());
    }
    catch (Exception& e) {
        THROW_E (EAC_HTTP_SERVER, EACHTTP_INVALID_REQUEST_JSON_FORMAT,
                 LOAD_STRING (_T("IDS_EACHTTP_INVALID_JSON")));
    }
    String groupName = requestJson["groupname"].s ().c_str ();

    string retGroupId ;
    try {
        ncThriftClient<ncTShareMgntClient> shareMgntClient (EacServiceAccessConfig::getInstance()->sharemgntHost, EacServiceAccessConfig::getInstance()->sharemgntPort);
        shareMgntClient->Usrm_CreatePersonGroup (retGroupId, toSTLString(userId), toSTLString(groupName));
    }
    catch (ncTException & e) {
        HandlencTException(e);
    }
    catch (TException & e) {
        THROW_E (EAC_HTTP_SERVER, EACHTTP_SERVER_UNKNOWN_ERR, e.what ());
    }

    // 回复
    JSON::Value replyJson;
    replyJson["groupid"] = retGroupId.c_str();
    string body;
    JSON::Writer::write (replyJson.o (), body);
    ncHttpSendReply (cntl, brpc::HTTP_STATUS_OK, "ok", body);

    NC_EAC_HTTP_SERVER_TRACE (_T("this: %p, cntl: %p, userId: %s end"), this, cntl, userId.getCStr ());
}

// protected
void
ncEACContactorHandler::EditGroup (brpc::Controller* cntl, const String& userId)
{
    NC_EAC_HTTP_SERVER_TRACE (_T("this: %p, cntl: %p, userId: %s begin"), this, cntl, userId.getCStr ());

    // 获取http请求的content
    string bodyBuffer = cntl->request_attachment ().to_string ();
    JSON::Value requestJson;
    try {
        JSON::Reader::read (requestJson, bodyBuffer.c_str (), bodyBuffer.size ());
    }
    catch (Exception& e) {
        THROW_E (EAC_HTTP_SERVER, EACHTTP_INVALID_REQUEST_JSON_FORMAT,
                 LOAD_STRING (_T("IDS_EACHTTP_INVALID_JSON")));
    }
    String groupId = requestJson["groupid"].s ().c_str ();
    String newName = requestJson["newname"].s ().c_str ();

    try {
        ncThriftClient<ncTShareMgntClient> shareMgntClient (EacServiceAccessConfig::getInstance()->sharemgntHost, EacServiceAccessConfig::getInstance()->sharemgntPort);
        shareMgntClient->Usrm_EditPersonGroup (toSTLString(userId), toSTLString(groupId), toSTLString(newName));
    }
    catch (ncTException & e) {
        HandlencTException(e);
    }
    catch (TException & e) {
        THROW_E (EAC_HTTP_SERVER, EACHTTP_SERVER_UNKNOWN_ERR, e.what ());
    }

    // 回复
    JSON::Value replyJson;
    string body;
    JSON::Writer::write (replyJson.o (), body);
    ncHttpSendReply (cntl, brpc::HTTP_STATUS_OK, "ok", body);

    NC_EAC_HTTP_SERVER_TRACE (_T("this: %p, cntl: %p, userId: %s end"), this, cntl, userId.getCStr ());
}

// protected
void
ncEACContactorHandler::AddPersons (brpc::Controller* cntl, const String& userId)
{

    NC_EAC_HTTP_SERVER_TRACE (_T("this: %p, cntl: %p, userId: %s begin"), this, cntl, userId.getCStr ());

    // 获取http请求的content
    string bodyBuffer = cntl->request_attachment ().to_string ();
    JSON::Value requestJson;
    try {
        JSON::Reader::read (requestJson, bodyBuffer.c_str (), bodyBuffer.size ());
    }
    catch (Exception& e) {
        THROW_E (EAC_HTTP_SERVER, EACHTTP_INVALID_REQUEST_JSON_FORMAT,
                 LOAD_STRING (_T("IDS_EACHTTP_INVALID_JSON")));
    }
    String groupId = requestJson["groupid"].s ().c_str ();
    vector<string> userIds;
    JSON::Object::const_iterator it = requestJson.o ().find ("userids");
    if (it != requestJson.o ().end ()) {
        const JSON::Array& userIdArray = it->second.a ();
        for (auto iter = userIdArray.begin (); iter != userIdArray.end (); ++iter) {
            string reqUserId = *iter;
            if (!reqUserId.empty ()) {
                userIds.push_back (reqUserId);
            }
        }
    }
    requestJson.clear ();

    try {
        ncThriftClient<ncTShareMgntClient> shareMgntClient (EacServiceAccessConfig::getInstance()->sharemgntHost, EacServiceAccessConfig::getInstance()->sharemgntPort);
        shareMgntClient->Usrm_AddPersonById (toSTLString(userId), userIds, toSTLString(groupId));
    }
    catch (ncTException & e) {
        HandlencTException(e);
    }
    catch (TException & e) {
        THROW_E (EAC_HTTP_SERVER, EACHTTP_SERVER_UNKNOWN_ERR, e.what ());
    }

    // 回复
    JSON::Value replyJson;
    string body;
    JSON::Writer::write (replyJson.o (), body);
    ncHttpSendReply (cntl, brpc::HTTP_STATUS_OK, "ok", body);

    // 记录审计日志
    String msg;
    String exMsg;
    if (userIds.size() == 0) {
        return;
    }
    // 批量获取用户信息
    vector<String> uids;
    for (size_t i = 0; i < userIds.size (); ++i) {
        uids.push_back(toCFLString(userIds[i]));
    }
    map<String, ncACSUserInfo> userInfoMap;
    _acsShareMgnt->GetUserInfoByIdBatch (uids, userInfoMap);

    // 获取组名称
    String groupName = _acsShareMgnt->GetNameByAccessorId(groupId, ncIOCAccesorType::IOC_USER_GROUP);
    map<String, ncACSUserInfo>::iterator user_iter = userInfoMap.begin();
    for(; user_iter != userInfoMap.end(); ++user_iter) {
        msg.format(LOAD_STRING("IDS_ADD_CONTACTOR_PERSON"),
                                user_iter->second.visionName.getCStr (),
                                groupName.getCStr());

        ncEACHttpServerUtil::Log (cntl, userId, ncTokenVisitorType::REALNAME, ncTLogType::NCT_LT_OPEARTION, ncTLogLevel::NCT_LL_INFO,
            ncTDocOperType::NCT_DOT_SET, msg, exMsg);
    }

    NC_EAC_HTTP_SERVER_TRACE (_T("this: %p, cntl: %p, userId: %s end"), this, cntl, userId.getCStr ());
}

// protected
void
ncEACContactorHandler::DeletePersons (brpc::Controller* cntl, const String& userId)
{

    NC_EAC_HTTP_SERVER_TRACE (_T("this: %p, cntl: %p, userId: %s begin"), this, cntl, userId.getCStr ());
    // 获取http请求的content
    string bodyBuffer = cntl->request_attachment ().to_string ();
    JSON::Value requestJson;
    try {
        JSON::Reader::read (requestJson, bodyBuffer.c_str (), bodyBuffer.size ());
    }
    catch (Exception& e) {
        THROW_E (EAC_HTTP_SERVER, EACHTTP_INVALID_REQUEST_JSON_FORMAT,
                 LOAD_STRING (_T("IDS_EACHTTP_INVALID_JSON")));
    }
    String groupId = requestJson["groupid"].s ().c_str ();

    vector<string> userIds;
    JSON::Object::const_iterator it = requestJson.o ().find ("userids");
    if (it != requestJson.o ().end ()) {
        const JSON::Array& userIdArray = it->second.a ();
        for (auto iter = userIdArray.begin (); iter != userIdArray.end (); ++iter) {
            string reqUserId = *iter;
            if (!reqUserId.empty ()) {
                userIds.push_back (reqUserId);
            }
        }
    }
    requestJson.clear ();

    try {
        ncThriftClient<ncTShareMgntClient> shareMgntClient (EacServiceAccessConfig::getInstance()->sharemgntHost, EacServiceAccessConfig::getInstance()->sharemgntPort);
        shareMgntClient->Usrm_DelPerson (toSTLString(userId), userIds, toSTLString(groupId));
    }
    catch (ncTException & e) {
        HandlencTException(e);
    }
    catch (TException & e) {
        THROW_E (EAC_HTTP_SERVER, EACHTTP_SERVER_UNKNOWN_ERR, e.what ());
    }

    // 回复
    JSON::Value replyJson;
    string body;
    JSON::Writer::write (replyJson.o (), body);
    ncHttpSendReply (cntl, brpc::HTTP_STATUS_OK, "ok", body);

    // 记录审计日志
    String msg;
    String exMsg;
    if (userIds.size() == 0) {
        return;
    }
    // 批量获取用户信息
    vector<String> uids;
    for (size_t i = 0; i < userIds.size (); ++i) {
        uids.push_back(toCFLString(userIds[i]));
    }
    map<String, ncACSUserInfo> userInfoMap;
    _acsShareMgnt->GetUserInfoByIdBatch (uids, userInfoMap);

    // 获取联系人组id对应的信息
    // 获取组名称
    String groupName = _acsShareMgnt->GetNameByAccessorId(groupId, ncIOCAccesorType::IOC_USER_GROUP);
    map<String, ncACSUserInfo>::iterator user_iter = userInfoMap.begin();
    for(; user_iter != userInfoMap.end(); ++user_iter) {
        msg.format(LOAD_STRING("IDS_DELETE_CONTACTOR_PERSON"),
                                user_iter->second.visionName.getCStr (),
                                groupName.getCStr());

        ncEACHttpServerUtil::Log (cntl, userId, ncTokenVisitorType::REALNAME, ncTLogType::NCT_LT_OPEARTION, ncTLogLevel::NCT_LL_WARN,
            ncTDocOperType::NCT_DOT_SET, msg, exMsg);
    }

    NC_EAC_HTTP_SERVER_TRACE (_T("this: %p, cntl: %p, userId: %s end"), this, cntl, userId.getCStr ());
}

// protected
void
ncEACContactorHandler::SearchPersons (brpc::Controller* cntl, const String& userId)
{
    NC_EAC_HTTP_SERVER_TRACE (_T("this: %p, cntl: %p, userId: %s begin"), this, cntl, userId.getCStr ());
    // 获取http请求的content
    string bodyBuffer = cntl->request_attachment ().to_string ();
    JSON::Value requestJson;
    try {
        JSON::Reader::read (requestJson, bodyBuffer.c_str (), bodyBuffer.size ());
    }
    catch (Exception& e) {
        THROW_E (EAC_HTTP_SERVER, EACHTTP_INVALID_REQUEST_JSON_FORMAT,
                 LOAD_STRING (_T("IDS_EACHTTP_INVALID_JSON")));
    }
    String key = requestJson["key"].s ().c_str ();
    vector<ncTSearchPersonGroup> retInfo;
    try {
        ncThriftClient<ncTShareMgntClient> shareMgntClient (EacServiceAccessConfig::getInstance()->sharemgntHost, EacServiceAccessConfig::getInstance()->sharemgntPort);
        shareMgntClient->Usrm_SearchPersonFromGroupByName (retInfo, toSTLString(userId), toSTLString(key));
    }
    catch (ncTException & e) {
        HandlencTException(e);
    }
    catch (TException & e) {
        THROW_E (EAC_HTTP_SERVER, EACHTTP_SERVER_UNKNOWN_ERR, e.what ());
    }

    // 回复
    JSON::Value replyJson;
    JSON::Array& usersJson = replyJson["userinfos"].a ();
    for (int i = 0; i < retInfo.size(); ++i) {
        usersJson.push_back (JSON::OBJECT);
        JSON::Object& tmpObj = usersJson.back ().o ();
        tmpObj["userid"] = retInfo[i].userId;
        tmpObj["account"] = retInfo[i].loginName;
        tmpObj["name"] = retInfo[i].displayName;
        tmpObj["groupid"] = retInfo[i].groupId;
        tmpObj["groupname"] = retInfo[i].groupName;
    }
    string body;
    JSON::Writer::write (replyJson.o (), body);

    ncHttpSendReply (cntl, brpc::HTTP_STATUS_OK, "ok", body);

    NC_EAC_HTTP_SERVER_TRACE (_T("this: %p, cntl: %p, userId: %s end"), this, cntl, userId.getCStr ());
}

// protected
void
ncEACContactorHandler::GetPersonFromGroup (brpc::Controller* cntl, const String& userId)
{
    NC_EAC_HTTP_SERVER_TRACE (_T("this:%p, cntl: %p, userId: %s begin"), this, cntl, userId.getCStr ());

    // 获取http请求的content
    string bodyBuffer = cntl->request_attachment ().to_string ();

    JSON::Value requestJson;
    try {
        JSON::Reader::read (requestJson, bodyBuffer.c_str (), bodyBuffer.size ());
    }
    catch (Exception& e) {
        THROW_E (EAC_HTTP_SERVER, EACHTTP_INVALID_REQUEST_JSON_FORMAT,
            LOAD_STRING (_T("IDS_EACHTTP_INVALID_JSON")));
    }

    vector<ncTUsrmGetUserInfo> userInfos;
    string groupId = requestJson["groupid"].s ();
    int32_t start = requestJson["start"].i ();
    int32_t limit = requestJson["limit"].i ();

    try {
        ncThriftClient<ncTShareMgntClient> shareMgntClient (EacServiceAccessConfig::getInstance()->sharemgntHost, EacServiceAccessConfig::getInstance()->sharemgntPort);
        shareMgntClient->Usrm_GetPersonFromGroup (userInfos, userId.getCStr (), groupId, start, limit);
    }
    catch (ncTException & e) {
        HandlencTException(e);
    }
    catch (TException & e) {
        THROW_E (EAC_HTTP_SERVER, EACHTTP_SERVER_UNKNOWN_ERR, e.what ());
    }

    // reply
    JSON::Value replyJson;
    JSON::Array& usersJson = replyJson["userinfos"].a ();
    for (size_t i = 0; i < userInfos.size (); ++i) {
        usersJson.push_back (JSON::OBJECT);
        JSON::Object& tmpObj = usersJson.back ().o ();
        tmpObj["userid"] = userInfos[i].id.c_str ();
        tmpObj["username"] = userInfos[i].user.displayName.c_str ();
        tmpObj["email"] = userInfos[i].user.email.c_str ();

        JSON::Array& jArray = tmpObj["departname"].a ();
        for (size_t j = 0; j < userInfos[i].user.departmentNames.size (); ++j) {
            jArray.push_back (userInfos[i].user.departmentNames[j].c_str ());
        }
    }

    string body;
    JSON::Writer::write (replyJson.o (), body);
    ncHttpSendReply (cntl, brpc::HTTP_STATUS_OK, "ok", body);

    NC_EAC_HTTP_SERVER_TRACE (_T("this: %p, cntl: %p, userId: %s end"), this, cntl, userId.getCStr ());
}

/*
 * 处理thrift 联系人组接口调用异常
*/
void ncEACContactorHandler::HandlencTException(ncTException & e)
{
    if (e.errID == ncTShareMgntError::NCT_INVALID_GROUP_NAME) {
        THROW_E (EAC_HTTP_SERVER, EACHTTP_INVALID_GROUP_NAME, e.expMsg.c_str ());
    }
    else if (e.errID == ncTShareMgntError::NCT_GROUP_HAS_EXIST) {
        THROW_E (EAC_HTTP_SERVER, EACHTTP_GROUP_HAS_EXIST, e.expMsg.c_str ());
    }
    else if (e.errID == ncTShareMgntError::NCT_GROUP_NOT_EXIST) {
        THROW_E (EAC_HTTP_SERVER, EACHTTP_GROUP_NOT_EXIST, e.expMsg.c_str ());
    }
    else if (e.errID == ncTShareMgntError::NCT_USER_NOT_EXIST) {
        THROW_E (EAC_HTTP_SERVER, EACHTTP_USER_NOT_EXIST, e.expMsg.c_str ());
    }
    else if (e.errID == ncTShareMgntError::NCT_CANNOT_OPERATE_TMP_GROUP) {
        THROW_E (EAC_HTTP_SERVER, EACHTTP_CANNOT_OPERATE_TMP_GROUP, e.expMsg.c_str ());
    }
    else if (e.errID == ncTShareMgntError::NCT_USER_NOT_IN_PERM_SOCPE) {
        THROW_E (EAC_HTTP_SERVER, EACHTTP_USER_NOT_IN_PERM_SOCPE, e.expMsg.c_str ());
    }
    else if (e.errID == ncTShareMgntError::NCT_CANNOT_ADD_SELF_TO_GROUP) {
        THROW_E (EAC_HTTP_SERVER, EACHTTP_CANNOT_ADD_SELF_TO_GROUP, e.expMsg.c_str ());
    }
    else if (e.errID == ncTShareMgntError::NCT_CONTACT_NOT_EXIST) {
        THROW_E (EAC_HTTP_SERVER, EACHTTP_CONTACT_NOT_EXIST, e.expMsg.c_str ());
    }
    else if (e.errID == ncTShareMgntError::NCT_CONTACT_HASE_EXIST) {
        THROW_E (EAC_HTTP_SERVER, EACHTTP_CONTACT_HASE_EXIST, e.expMsg.c_str ());
    }
    else {
        THROW_E (EAC_HTTP_SERVER, EACHTTP_SERVER_UNKNOWN_ERR, e.expMsg.c_str ());
    }
}
