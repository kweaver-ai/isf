/***************************************************************************************************
authentication.cpp:
    Copyright (c) Eisoo Software Inc. (2020), All rights reserved.

Purpose:
    authentication 实现

Author:
    Yuanbin.yan@aishu.cn

Creating Time:
    2023-08-28
***************************************************************************************************/
#include <abprec.h>
#include "authentication.h"
#include "drivenadapter.h"
#include <dataapi/ncJson.h>
#include <ossclient/public/ncIOSSClient.h>
#include <ossclient/ossclient.h>
#include <dataapi/dataapi.h>
#include <ncutil/ncBusinessDate.h>

#include "serviceAccessConfig.h"

/* Implementation file */
NS_IMPL_QUERY_INTERFACE1 (authentication, authenticationInterface)

// protected
NS_IMETHODIMP_(nsrefcnt) authentication::AddRef (void)
{
    return 1;
}

// protected
NS_IMETHODIMP_(nsrefcnt) authentication::Release (void)
{
    return 1;
}

AB_DEFINE_THREADSAFE_SINGLETON_NO_POOL (authentication)

authentication::authentication (void): _ossClientPtr (0)
{
    NC_DRIVEN_ADAPTER_TRACE (_T("this: %p"), this);

    _loginUrl.format(_T("http://%s:%d/api/authentication/v1/client-account-auth"),
        ServiceAccessConfig::getInstance()->authenticationPrivateHost.getCStr(), ServiceAccessConfig::getInstance()->authenticationPrivatePort);

    _auditLogUrl.format(_T("http://%s:%d/api/authentication/v1/audit-log"),
        ServiceAccessConfig::getInstance()->authenticationPrivateHost.getCStr(), ServiceAccessConfig::getInstance()->authenticationPrivatePort);
}

authentication::~authentication (void)
{
    NC_DRIVEN_ADAPTER_TRACE (_T("this: %p"), this);
}

void authentication::createOSSClient ()
{
    if (nullptr == _ossClientPtr.get ()) {
        nsresult ret;
        nsCOMPtr<ncIOSSClient> ossClient = do_CreateInstance (NC_OSS_CLIENT_CONTRACTID, &ret);
        if (NS_FAILED (ret)){
            THROW_E (ACS_DRIVEN_ADAPTER, FAILED_TO_CREATE_XPCOM_INSTANCE,
                     _T("Failed to create OSSClient: %s(0x%x)"),
                     String::toString ((int64)ret).getCStr (), (int64)ret);
        }
        _ossClientPtr.reset (new nsCOMPtr<ncIOSSClient>(ossClient));
    }
}

/* [notxpcom] int ClientLogin ([const] in stlStringRef account, [const] in stlStringRef password, [const] in ncTUserLoginOptionRef option, in ncJSONRef JSONValue); */
NS_IMETHODIMP_(int) authentication::ClientLogin(const string& account, const string& password, const ncTUserLoginOption& option, JSON::Value& responseJson)
{
    NC_DRIVEN_ADAPTER_TRACE (_T("[BEGIN] this: %p"), this);
    std::string content;
    vector<string> headers;
    JSON::Value requestJson;
    requestJson["method"] = "GET";
    requestJson["account"] = account.c_str();
    requestJson["password"] = password.c_str();
    requestJson["option"]["uuid"] = option.uuid.c_str();
    requestJson["option"]["vcode"] = option.vcode.c_str();
    requestJson["option"]["vcodeType"] = (int)option.vcodeType;
    JSON::Writer::write(requestJson.o(), content);

    ncOSSResponse response;
    createOSSClient();
    (*_ossClientPtr)->Post (_loginUrl.getCStr (), content, headers, 30, response);

    JSON::Reader::read (responseJson, response.body.c_str (), response.body.size ());

    return response.code;
    NC_DRIVEN_ADAPTER_TRACE (_T("[END] this: %p"), this);
}

/* [notxpcom] void AuditLog ([const] in StringRef userId, in ncTokenVisitorTypeRef typ, in ncTLogTypeRef logType, in ncTLogLevelRef level, in int opType, [const] in StringRef msg, [const] in StringRef exmsg, [const] in StringRef ip, [const] in StringRef mac, [const] in StringRef userAgent); */
NS_IMETHODIMP_(void) authentication::AuditLog(const String & userId, ncTokenVisitorType & typ, ncTLogType & logType, ncTLogLevel & level, int opType, const String & msg,
                                const String & exmsg, const String & ip, const String & mac, const String & userAgent)
{
    NC_DRIVEN_ADAPTER_TRACE (_T("[BEGIN] this: %p"), this);
    // 选择topic
    string topic;
    switch (logType)    {
    case ncTLogType::NCT_LT_LOGIN:
        topic = "as.audit_log.log_login";
        break;
    case ncTLogType::NCT_LT_MANAGEMENT:
        topic = "as.audit_log.log_management";
        break;
    case ncTLogType::NCT_LT_OPEARTION:
        topic = "as.audit_log.log_operation";
        break;
    default:
        THROW_E (ACS_DRIVEN_ADAPTER, INVALID_PARAMETER_VALUES, "Log Unknown log type");
        break;
    }

    // 整合消息内容
    JSON::Value message;
    handleMessage(userId, typ, level, opType, msg, exmsg, ip, mac, userAgent, message);

    // 整理数据
    std::string content;
    vector<string> headers;
    JSON::Value requestJson;
    requestJson["topic"] = topic;
    requestJson["message"] = message;
    JSON::Writer::write(requestJson.o(), content);

    ncOSSResponse response;
    createOSSClient();
    (*_ossClientPtr)->Post (_auditLogUrl.getCStr (), content, headers, 30, response);

    if (response.code != 204)
    {
        if (response.code == 0) {
            THROW_E (ACS_DRIVEN_ADAPTER, DRIVEN_HTTP_SERVER_INTERNAL_ERR, _T("Could not connect to server"));
        }
        else{
            JSON::Value errorJson;
            errorJson["code"] = response.code;
            errorJson["body"] = response.body;
            throw errorJson;
        }
    }

    NC_DRIVEN_ADAPTER_TRACE (_T("[END] this: %p"), this);
}

void authentication::handleMessage(const String& userId, ncTokenVisitorType typ, ncTLogLevel level, int opType, const String& msg,
                        const String& exmsg, const String& ip, const String& mac, const String& userAgent, JSON::Value & logMsgs)
{
    ncObjectID id;
    id.GenerateObjectID ();

    string additionalInfo;
    JSON::Value additionalInfoItems;
    logMsgs["user_id"] = userId.getCStr();
    switch (typ) {
    case ncTokenVisitorType::REALNAME:
        logMsgs["user_type"] = "authenticated_user";
        break;
    case ncTokenVisitorType::ANONYMOUS:
        logMsgs["user_type"] = "anonymous_user";
        logMsgs["user_name"] = DRIVEN_LOAD_STRING (_T("IDS_EACHTTP_ANONYMOUS_NAME"));

        // 增加anyrobot文件建模信息
        additionalInfoItems["user_account"] = DRIVEN_LOAD_STRING (_T("IDS_EACHTTP_ANONYMOUS_NAME"));
        JSON::Writer::write (additionalInfoItems.o (), additionalInfo);

        logMsgs["additional_info"] = additionalInfo.c_str();
        break;
    case ncTokenVisitorType::BUSINESS:
        logMsgs["user_type"] = "app";
        break;
    default:
        THROW_E (ACS_DRIVEN_ADAPTER, INVALID_PARAMETER_VALUES, "Log Unknown visitor type");
        break;
    }
    logMsgs["level"] = level;
    logMsgs["date"] = BusinessDate::getCurrentTime ();
    logMsgs["ip"] = ip.getCStr();
    logMsgs["mac"] = mac.getCStr();
    logMsgs["msg"] = msg.getCStr();
    logMsgs["ex_msg"] = exmsg.getCStr();
    logMsgs["user_agent"] = userAgent.getCStr();
    logMsgs["op_type"] = opType;
    logMsgs["out_biz_id"] = id.GetString().getCStr();
}