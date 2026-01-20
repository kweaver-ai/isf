#include <ncutil/ncBusinessDate.h>

#include "eachttpserver.h"
#include "ncEACDeviceHandler.h"
#include "ncEACHttpServerUtil.h"

// public
ncEACDeviceHandler::ncEACDeviceHandler(ncIACSDeviceManager* deviceManager)
    : _acsDeviceManager (deviceManager)
{
    NC_EAC_HTTP_SERVER_TRACE (_T("this: %p"), this);

    _methodFuncs.insert (pair<String, ncMethodFunc>(_T("list"), &ncEACDeviceHandler::onList));
    _methodFuncs.insert (pair<String, ncMethodFunc>(_T("disable"), &ncEACDeviceHandler::onDisable));
    _methodFuncs.insert (pair<String, ncMethodFunc>(_T("enable"), &ncEACDeviceHandler::onEnable));
    _methodFuncs.insert (pair<String, ncMethodFunc>(_T("erase"), &ncEACDeviceHandler::onErase));
    _methodFuncs.insert (pair<String, ncMethodFunc>(_T("getstatus"), &ncEACDeviceHandler::onGetStatus));
    _methodFuncs.insert (pair<String, ncMethodFunc>(_T("onerasesuc"), &ncEACDeviceHandler::onEraseSuc));
}

// public
ncEACDeviceHandler::~ncEACDeviceHandler (void)
{
    NC_EAC_HTTP_SERVER_TRACE (_T("this: %p"), this);
}

// public
void
ncEACDeviceHandler::doDeviceRequestHandler (brpc::Controller* cntl)
{
    NC_EAC_HTTP_SERVER_TRY

        // 获取query string
        String method;
        String tokenId;
        String userId;
        ncHttpGetParams (cntl, method, tokenId, userId);

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

        // onerasesuc不需要检查token
        if(method != "onerasesuc") {
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

// protected
void
ncEACDeviceHandler::onList (brpc::Controller* cntl, const String& userId)
{
    NC_EAC_HTTP_SERVER_TRACE (_T("this: %p, cntl: %p, userId: %s begin"), this, cntl, userId.getCStr ());

    // 获取用户的登录设备信息
    vector<ncDeviceInfo> deviceInfos;
    _acsDeviceManager->GetLoginedMobileDevices(userId, deviceInfos);

    JSON::Value replyJson;
    JSON::Array& devicesJson = replyJson["deviceinfos"].a ();
    for (size_t i = 0; i < deviceInfos.size (); ++i) {
        devicesJson.push_back (JSON::OBJECT);

        JSON::Object& tmpObj = devicesJson.back ().o ();

        tmpObj["name"] = deviceInfos[i].baseInfo.name.getCStr ();
        tmpObj["ostype"] = static_cast<int>(deviceInfos[i].baseInfo.clientType);
        tmpObj["devicetype"] = deviceInfos[i].baseInfo.deviceType.getCStr();
        tmpObj["udid"] = deviceInfos[i].baseInfo.udid.getCStr ();
        tmpObj["lastloginip"] = deviceInfos[i].baseInfo.lastLoginIp.getCStr ();
        tmpObj["lastlogintime"] = deviceInfos[i].baseInfo.lastLoginTime;
        tmpObj["eraseflag"] = deviceInfos[i].eraseFlag;
        tmpObj["lasterasetime"] = deviceInfos[i].lastEraseTime;
        tmpObj["disableflag"] = deviceInfos[i].disableFlag;
        tmpObj["loginflag"] = deviceInfos[i].loginFlag;
        tmpObj["bindflag"] = deviceInfos[i].bindFlag;
    }

    // reply
    string body;
    JSON::Writer::write (replyJson.o (), body);
    ncHttpSendReply (cntl, brpc::HTTP_STATUS_OK, "ok", body);

    NC_EAC_HTTP_SERVER_TRACE (_T("this: %p, cntl: %p, userId: %s end"), this, cntl, userId.getCStr ());
}

// protected
void
ncEACDeviceHandler::onDisable (brpc::Controller* cntl, const String& userId)
{
    NC_EAC_HTTP_SERVER_TRACE (_T("this: %p, cntl: %p, userId: %s begin"), this, cntl, userId.getCStr ());

    // 获取http请求的content
    string bodyBuffer = cntl->request_attachment ().to_string ();

    // 获取udid
    JSON::Value requestJson;
    try {
        JSON::Reader::read (requestJson, bodyBuffer.c_str (), bodyBuffer.size ());
    }
    catch (Exception& e) {
        THROW_E (EAC_HTTP_SERVER, EACHTTP_INVALID_REQUEST_JSON_FORMAT,
            LOAD_STRING (_T("IDS_EACHTTP_INVALID_JSON")));
    }

    String udid = ::toCFLString (requestJson["udid"].s ());
    if(udid.isEmpty()) {
        THROW_E (EAC_HTTP_SERVER, EACHTTP_INVALID_UDID,
            LOAD_STRING (_T("IDS_EACHTTP_INVALID_UDID")));
    }

    _acsDeviceManager->SetDisableStatus(userId, udid);

    // reply
    string body;
    ncHttpSendReply (cntl, brpc::HTTP_STATUS_OK, "ok", body);

    String msg;
    msg.format(LOAD_STRING("IDS_DISABLE_DEVICE_SUC"), udid.getCStr());
    ncEACHttpServerUtil::Log (cntl, userId, ncTokenVisitorType::REALNAME, ncTLogType::NCT_LT_OPEARTION, ncTLogLevel::NCT_LL_INFO,
        ncTDocOperType::NCT_DOT_DEVICE_MGM, msg, "");

    NC_EAC_HTTP_SERVER_TRACE (_T("this: %p, cntl: %p, userId: %s end"), this, cntl, userId.getCStr ());
}

// protected
void
ncEACDeviceHandler::onEnable (brpc::Controller* cntl, const String& userId)
{
    NC_EAC_HTTP_SERVER_TRACE (_T("this: %p, cntl: %p, userId: %s begin"), this, cntl, userId.getCStr ());

    // 获取http请求的content
    string bodyBuffer = cntl->request_attachment ().to_string ();

    // 获取udid
    JSON::Value requestJson;
    try {
        JSON::Reader::read (requestJson, bodyBuffer.c_str (), bodyBuffer.size ());
    }
    catch (Exception& e) {
        THROW_E (EAC_HTTP_SERVER, EACHTTP_INVALID_REQUEST_JSON_FORMAT,
            LOAD_STRING (_T("IDS_EACHTTP_INVALID_JSON")));
    }

    String udid = ::toCFLString (requestJson["udid"].s ());
    if(udid.isEmpty()) {
        THROW_E (EAC_HTTP_SERVER, EACHTTP_INVALID_UDID,
            LOAD_STRING (_T("IDS_EACHTTP_INVALID_UDID")));
    }

    _acsDeviceManager->SetEnableStatus(userId, udid);

    // reply
    string body;
    ncHttpSendReply (cntl, brpc::HTTP_STATUS_OK, "ok", body);

    String msg;
    msg.format(LOAD_STRING("IDS_ENABLE_DEVICE_SUC"), udid.getCStr());
    ncEACHttpServerUtil::Log (cntl, userId, ncTokenVisitorType::REALNAME, ncTLogType::NCT_LT_OPEARTION, ncTLogLevel::NCT_LL_INFO,
        ncTDocOperType::NCT_DOT_DEVICE_MGM, msg, "");

    NC_EAC_HTTP_SERVER_TRACE (_T("this: %p, cntl: %p, userId: %s end"), this, cntl, userId.getCStr ());
}

// protected
void
ncEACDeviceHandler::onErase (brpc::Controller* cntl, const String& userId)
{
    NC_EAC_HTTP_SERVER_TRACE (_T("this: %p, cntl: %p, userId: %s begin"), this, cntl, userId.getCStr ());

    // 获取http请求的content
    string bodyBuffer = cntl->request_attachment ().to_string ();

    // 获取udid
    JSON::Value requestJson;
    try {
        JSON::Reader::read (requestJson, bodyBuffer.c_str (), bodyBuffer.size ());
    }
    catch (Exception& e) {
        THROW_E (EAC_HTTP_SERVER, EACHTTP_INVALID_REQUEST_JSON_FORMAT,
            LOAD_STRING (_T("IDS_EACHTTP_INVALID_JSON")));
    }

    String udid = ::toCFLString (requestJson["udid"].s ());
    if(udid.isEmpty()) {
        THROW_E (EAC_HTTP_SERVER, EACHTTP_INVALID_UDID,
            LOAD_STRING (_T("IDS_EACHTTP_INVALID_UDID")));
    }

    _acsDeviceManager->SetEraseStatus(userId, udid);

    // reply
    string body;
    ncHttpSendReply (cntl, brpc::HTTP_STATUS_OK, "ok", body);

    String msg;
    msg.format(LOAD_STRING("IDS_SEND_ERASE_DEVICE_REQ_SUC"), udid.getCStr());
    ncEACHttpServerUtil::Log (cntl, userId, ncTokenVisitorType::REALNAME, ncTLogType::NCT_LT_OPEARTION, ncTLogLevel::NCT_LL_INFO,
        ncTDocOperType::NCT_DOT_DEVICE_MGM, msg, "");

    NC_EAC_HTTP_SERVER_TRACE (_T("this: %p, cntl: %p, userId: %s end"), this, cntl, userId.getCStr ());
}

// protected
void
ncEACDeviceHandler::onGetStatus (brpc::Controller* cntl, const String& userId)
{
    NC_EAC_HTTP_SERVER_TRACE (_T("this: %p, cntl: %p, userId: %s begin"), this, cntl, userId.getCStr ());

    // 获取http请求的content
    string bodyBuffer = cntl->request_attachment ().to_string ();

    // 获取udid
    JSON::Value requestJson;
    try {
        JSON::Reader::read (requestJson, bodyBuffer.c_str (), bodyBuffer.size ());
    }
    catch (Exception& e) {
        THROW_E (EAC_HTTP_SERVER, EACHTTP_INVALID_REQUEST_JSON_FORMAT,
            LOAD_STRING (_T("IDS_EACHTTP_INVALID_JSON")));
    }

    String udid = ::toCFLString (requestJson["udid"].s ());
    if(udid.isEmpty()) {
        THROW_E (EAC_HTTP_SERVER, EACHTTP_INVALID_UDID,
            LOAD_STRING (_T("IDS_EACHTTP_INVALID_UDID")));
    }

    ncDeviceInfo info;
    bool ret = _acsDeviceManager->GetDeviceByUDID(userId, udid, info);
    if(ret == false) {
        THROW_E (EAC_HTTP_SERVER, EACHTTP_INVALID_UDID,
            LOAD_STRING (_T("IDS_EACHTTP_INVALID_UDID")));
    }

    // reply
    JSON::Value replyJson;
    replyJson["eraseflag"] = info.eraseFlag;
    replyJson["disableflag"] = info.disableFlag;

    string body;
    JSON::Writer::write (replyJson.o (), body);
    ncHttpSendReply (cntl, brpc::HTTP_STATUS_OK, "ok", body);

    NC_EAC_HTTP_SERVER_TRACE (_T("this: %p, cntl: %p, userId: %s end"), this, cntl, userId.getCStr ());
}

// protected
void
ncEACDeviceHandler::onEraseSuc (brpc::Controller* cntl, const String& userId)
{
    NC_EAC_HTTP_SERVER_TRACE (_T("this: %p, cntl: %p, userId: %s begin"), this, cntl, userId.getCStr ());

    // 获取http请求的content
    string bodyBuffer = cntl->request_attachment ().to_string ();

    // 获取udid
    JSON::Value requestJson;
    try {
        JSON::Reader::read (requestJson, bodyBuffer.c_str (), bodyBuffer.size ());
    }
    catch (Exception& e) {
        THROW_E (EAC_HTTP_SERVER, EACHTTP_INVALID_REQUEST_JSON_FORMAT,
            LOAD_STRING (_T("IDS_EACHTTP_INVALID_JSON")));
    }

    String udid = ::toCFLString (requestJson["udid"].s ());
    if(udid.isEmpty()) {
        THROW_E (EAC_HTTP_SERVER, EACHTTP_INVALID_UDID,
            LOAD_STRING (_T("IDS_EACHTTP_INVALID_UDID")));
    }

    _acsDeviceManager->SetEraseSucInfo(userId, udid, BusinessDate::getCurrentTime ());

    // reply
    string body;
    ncHttpSendReply (cntl, brpc::HTTP_STATUS_OK, "ok", body);

    String msg;
    msg.format(LOAD_STRING("IDS_ERASE_DEVICE_SUC"), udid.getCStr());
    ncEACHttpServerUtil::Log (cntl, userId, ncTokenVisitorType::REALNAME, ncTLogType::NCT_LT_OPEARTION, ncTLogLevel::NCT_LL_INFO,
        ncTDocOperType::NCT_DOT_DEVICE_MGM, msg, "");

    NC_EAC_HTTP_SERVER_TRACE (_T("this: %p, cntl: %p, userId: %s end"), this, cntl, userId.getCStr ());
}
