#include "eachttpserver.h"
#include "ncEACConfigHandler.h"
#include <ethriftutil/ncThriftClient.h>
#include "ncEACHttpServerUtil.h"
#include "eacServiceAccessConfig.h"

ncEACConfigHandler::ncEACConfigHandler (ncIACSShareMgnt* acsShareMgnt)
        :  _acsShareMgnt (acsShareMgnt)
{
    NC_EAC_HTTP_SERVER_TRACE (_T("this: %p"), this);

    NC_EAC_HTTP_SERVER_CHECK_ARGUMENT_NULL (acsShareMgnt);

    _methodFuncs.insert (pair<String, ncMethodFunc>(_T("get"), &ncEACConfigHandler::Get));
    _methodFuncs.insert (pair<String, ncMethodFunc>(_T("getoemconfigbysection"), &ncEACConfigHandler::GetOEMConfigBySection));
    _methodFuncs.insert (pair<String, ncMethodFunc>(_T("getdocwatermarkconfig"), &ncEACConfigHandler::GetDocWatermarkConfig));
    _methodFuncs.insert (pair<String, ncMethodFunc>(_T("getfilecrawlconfig"), &ncEACConfigHandler::GetFileCrawlConfig));
    _methodFuncs.insert (pair<String, ncMethodFunc>(_T("setquickstartstatus"), &ncEACConfigHandler::SetQuickStartStatus));

    _methodWhiteList.insert ("getoemconfigbysection");
    _methodWhiteList.insert ("getdocwatermarkconfig");

}

ncEACConfigHandler::~ncEACConfigHandler (void)
{
    NC_EAC_HTTP_SERVER_TRACE (_T("this: %p"), this);
}

void
ncEACConfigHandler::doConfigRequestHandler (brpc::Controller* cntl)
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

        // 如果在白名单中则不进行token验证
        if (_methodWhiteList.find (method.getCStr ()) == _methodWhiteList.end ()) {
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
        }

        // 消息处理
        ncMethodFunc func = iter->second;
        (this->*func) (cntl, userId);

    NC_EAC_HTTP_SERVER_CATCH (cntl)
}

void
ncEACConfigHandler::Get (brpc::Controller* cntl, const String& userId)
{
    NC_EAC_HTTP_SERVER_TRACE (_T("this: %p, cntl: %p, userId: %s begin"), this, cntl, userId.getCStr ());

    // 获取http请求的content
    string bodyBuffer = cntl->request_attachment ().to_string ();

    // 2019-02-27: 协议兼容 --> 【ANR-4103】之前无请求参数
    int clientType = 0;
    if (!bodyBuffer.empty ()) {
        JSON::Value requestJson;
        try {
            JSON::Reader::read (requestJson, bodyBuffer.c_str (), bodyBuffer.size ());
        }
        catch (Exception& e) {
            THROW_E (EAC_HTTP_SERVER, EACHTTP_INVALID_REQUEST_JSON_FORMAT,
                LOAD_STRING (_T("IDS_EACHTTP_INVALID_JSON")));
        }

        clientType = requestJson.get<int> ("ostype", 0);
        ncEACHttpServerUtil::CheckOsType(static_cast<ACSClientType>(clientType));
    }

    int interval = _acsShareMgnt->GetClearCacheInterval();
    int64 size = _acsShareMgnt->GetClearCacheSize();
    int64 detectInterval = Int::getValue(_acsShareMgnt->GetShareMgntConfig("client_detect_interval"));

    JSON::Value replyJson;
    JSON::Array& cacheJson = replyJson["cache"].a ();
    cacheJson.push_back (JSON::OBJECT);
    JSON::Object& tmpObj = cacheJson.back ().o ();
    tmpObj["interval"] = interval;
    tmpObj["size"] = size;
    replyJson["detect_interval"] = detectInterval;

    ncTLocalSyncConfig localSyncConfig;
    try {
        ncThriftClient<ncTShareMgntClient> shareMgntClient (EacServiceAccessConfig::getInstance()->sharemgntHost, EacServiceAccessConfig::getInstance()->sharemgntPort);
        shareMgntClient->GetLocalSyncConfigByUserId(localSyncConfig, toSTLString(userId));
    }
    catch (ncTException & e) {
        THROW_E (EAC_HTTP_SERVER, EACHTTP_SERVER_UNKNOWN_ERR, e.expMsg.c_str ());
    }
    catch (TException & e) {
        THROW_E (EAC_HTTP_SERVER, EACHTTP_SERVER_UNKNOWN_ERR, e.what ());
    }
    replyJson["localsync"]["openstatus"] = localSyncConfig.openStatus;
    replyJson["localsync"]["deletestatus"] = localSyncConfig.deleteStatus;

    // 获取用户信息
    ncACSUserInfo userInfo;
    _acsShareMgnt->GetUserInfoById (userId, userInfo);
    // 快速入门文档阅读状态，使用比特范围 1-7，和osType参数的值对应
    replyJson["needquickstart"] = ((userInfo.documentReadStatus & (1<<clientType)) == 0) ? true : false;

    // 回复
    string body;
    JSON::Writer::write (replyJson.o (), body);
    ncHttpSendReply (cntl, brpc::HTTP_STATUS_OK, "ok", body);

    NC_EAC_HTTP_SERVER_TRACE (_T("this: %p, cntl: %p, userId: %s end"), this, cntl, userId.getCStr ());
}

void
ncEACConfigHandler::GetOEMConfigBySection (brpc::Controller* cntl, const String& fakeUserId)
{
    NC_EAC_HTTP_SERVER_TRACE (_T("this: %p, cntl: %p begin"), this, cntl);

    // 获取http请求的content
    string bodyBuffer = cntl->request_attachment ().to_string ();

    // 获取ticket
    JSON::Value requestJson;
    try {
        JSON::Reader::read (requestJson, bodyBuffer.c_str (), bodyBuffer.size ());
    }
    catch (Exception& e) {
        THROW_E (EAC_HTTP_SERVER, EACHTTP_INVALID_REQUEST_JSON_FORMAT,
            LOAD_STRING (_T("IDS_EACHTTP_INVALID_JSON")));
    }

    string section = requestJson["section"].s ();

    JSON::Value replyJson;

    vector<ncTOEMInfo> retOEMInfos;
    try {
        ncThriftClient<ncTShareMgntClient> shareMgntClient (EacServiceAccessConfig::getInstance()->sharemgntHost, EacServiceAccessConfig::getInstance()->sharemgntPort);
        shareMgntClient->OEM_GetConfigBySection (retOEMInfos, section);
    }
    catch (ncTException & e) {
        THROW_E (EAC_HTTP_SERVER, EACHTTP_GET_OEMCONFIG_ERROR, e.expMsg.c_str ());
    }
    catch (TException & e) {
        THROW_E (EAC_HTTP_SERVER, EACHTTP_SERVER_UNKNOWN_ERR, e.what ());
    }

    for (int i = 0; i < retOEMInfos.size(); ++i) {
        replyJson[retOEMInfos[i].option] = retOEMInfos[i].value;
    }

    string body;
    JSON::Writer::write (replyJson.o (), body);
    ncHttpSendReply (cntl, brpc::HTTP_STATUS_OK, "ok", body);

    NC_EAC_HTTP_SERVER_TRACE (_T("this: %p, cntl: %p end"), this, cntl);
}

void
ncEACConfigHandler::GetDocWatermarkConfig (brpc::Controller* cntl, const String& fakeUserId)
{
    NC_EAC_HTTP_SERVER_TRACE (_T("this: %p, cntl: %p, userId: %s begin"), this, cntl, fakeUserId.getCStr ());

    string body;
    try {
        ncThriftClient<ncTShareMgntClient> shareMgntClient (EacServiceAccessConfig::getInstance()->sharemgntHost, EacServiceAccessConfig::getInstance()->sharemgntPort);
        shareMgntClient->GetDocWatermarkConfig(body);
    }
    catch (ncTException & e) {
        THROW_E (EAC_HTTP_SERVER, EACHTTP_GET_WATERMARK_CONFIG_ERROR,
            LOAD_STRING (_T("IDS_EACHTTP_GET_WATERMARK_CONFIG_ERROR")), e.expMsg.c_str ());
    }
    catch (TException & e) {
        THROW_E (EAC_HTTP_SERVER, EACHTTP_GET_WATERMARK_CONFIG_ERROR,
            LOAD_STRING (_T("IDS_EACHTTP_GET_WATERMARK_CONFIG_ERROR")), e.what ());
    }

    ncHttpSendReply (cntl, brpc::HTTP_STATUS_OK, "ok", body);

    NC_EAC_HTTP_SERVER_TRACE (_T("this: %p, cntl:%p, userId: %s end"), this, cntl, fakeUserId.getCStr ());
}

void
ncEACConfigHandler::GetFileCrawlConfig (brpc::Controller* cntl, const String& userId)
{
    NC_EAC_HTTP_SERVER_TRACE (_T("this: %p, cntl: %p, userId: %s begin"), this, cntl, userId.getCStr ());

    int64 fileCrawStatus = Int::getValue(_acsShareMgnt->GetShareMgntConfig("file_crawl_status"));

    ncTFileCrawlConfig fileCrawlconfig;
    try {
        ncThriftClient<ncTShareMgntClient> shareMgntClient (EacServiceAccessConfig::getInstance()->sharemgntHost, EacServiceAccessConfig::getInstance()->sharemgntPort);
        shareMgntClient->GetFileCrawlConfigByUserId(fileCrawlconfig, toSTLString(userId));
    }
    catch (ncTException & e) {
        THROW_E (EAC_HTTP_SERVER, EACHTTP_GET_SITE_OFFICE_ONLINE_INFO_ERROR,
            LOAD_STRING (_T("IDS_EACHTTP_GET_SITE_OFFICE_ONLINE_INFO_ERROR")), e.expMsg.c_str ());
    }
    catch (TException & e) {
        THROW_E (EAC_HTTP_SERVER, EACHTTP_GET_SITE_OFFICE_ONLINE_INFO_ERROR,
            LOAD_STRING (_T("IDS_EACHTTP_GET_SITE_OFFICE_ONLINE_INFO_ERROR")), e.what ());
    }

    JSON::Value responseJson;
    responseJson["needfilecrawl"] = false;

    if (fileCrawStatus == 1 && fileCrawlconfig.strategyId != -1 && !fileCrawlconfig.docName.empty()) {
        responseJson["needfilecrawl"] = true;
        responseJson["filecrawltype"] = fileCrawlconfig.fileCrawlType;
        responseJson["filecrawldocid"] = fileCrawlconfig.docId;
    }

    // 回复
    string body;
    JSON::Writer::write (responseJson.o (), body);
    ncHttpSendReply (cntl, brpc::HTTP_STATUS_OK, "ok", body);

    NC_EAC_HTTP_SERVER_TRACE (_T("this: %p, cntl: %p, userId: %s end"), this, cntl, userId.getCStr ());
}

void
ncEACConfigHandler::SetQuickStartStatus (brpc::Controller* cntl, const String& userId)
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

    int clientType = requestJson.get<int> ("ostype", -1);
    ncEACHttpServerUtil::CheckOsType(static_cast<ACSClientType>(clientType));

    // 快速入门文档阅读状态，使用比特范围 1-7，和osType参数的值对应
    int leftShiftNum = clientType;
    _acsShareMgnt->UpdateUserDocumentReadStatus(userId, leftShiftNum);

    // 回复
    ncHttpSendReply (cntl, brpc::HTTP_STATUS_OK, "ok", "");

    NC_EAC_HTTP_SERVER_TRACE (_T("this: %p, cntl: %p, userId: %s end"), this, cntl, userId.getCStr ());
}
