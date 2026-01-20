/****************************************************************************************************
ncEACPolicyHandler.cpp
     Copyright (c) Eisoo Software, Inc.(2009 - 2010), All rights reserved.

Purpose:
    ncEACPolicyHandler

Author:
    Sunshine.tang@aishu.cn

Created Time:
    2021-02-03
****************************************************************************************************/
#include "eachttpserver.h"
#include "ncEACPolicyHandler.h"


// public
ncEACPolicyHandler::ncEACPolicyHandler (ncIACSPolicyManager* acsPolicyManager)
    : _acsPolicyManager (acsPolicyManager)
{
    NC_EAC_HTTP_SERVER_TRACE (_T("this: %p"), this);

    NC_EAC_HTTP_SERVER_CHECK_ARGUMENT_NULL (acsPolicyManager);
    _accountTypeMap = {{"other", ACSAccountType::OTHER}, {"id_card", ACSAccountType::ID_CARD}};
    _clientTypeMap = {{"unknown", ACSClientType::UNKNOWN}, {"ios", ACSClientType::IOS}, {"android", ACSClientType::ANDROID},
                                            {"windows_phone", ACSClientType::WINDOWS_PHONE}, {"windows", ACSClientType::WINDOWS}, {"mac_os", ACSClientType::MAC_OS},
                                            {"web", ACSClientType::WEB}, {"mobile_web", ACSClientType::MOBILE_WEB}, {"nas", ACSClientType::NAS},
                                            {"console_web", ACSClientType::CONSOLE_WEB}, {"deploy_web", ACSClientType::DEPLOY_WEB}, {"linux", ACSClientType::LINUX}, {"app", ACSClientType::APP}};
}

// public
ncEACPolicyHandler::~ncEACPolicyHandler (void)
{
    NC_EAC_HTTP_SERVER_TRACE (_T("this: %p"), this);
}

// public
void
ncEACPolicyHandler::CheckPolicy (brpc::Controller* cntl)
{
    NC_EAC_HTTP_SERVER_TRY
        NC_EAC_HTTP_SERVER_TRACE (_T("[ncEACPolicyHandler::CheckPolicy] this: %p, cntl: %p begin"), this, cntl);

        cntl->http_response().set_content_type("application/json; charset=UTF-8");
        // 判断请求方法是否为POST
        if (cntl->http_request ().method () != brpc::HTTP_METHOD_POST){
            THROW_E (EACHTTP_SERVER_ERR_PROVIDER_NAME, EACHTTP_SERVER_METHOD_NOT_IMPLEMENTED,
                 LOAD_STRING (_T("IDS_EACHTTP_METHOD_INVALID")));
        }
        NC_EAC_HTTP_SERVER_TRACE (_T("this: %p, cntl:%p begin"), this, cntl);

        // 获取实名用户信息
        JSON::Value requestJson;
        ncPolicyCheckInfo policyInfo;
        parseUserInfo(cntl, policyInfo, requestJson);
        String userId = requestJson["userid"].s().c_str();

        // 用户策略检测
        _acsPolicyManager->CheckPolicy (policyInfo);

        string body;
        ncHttpSendReply (cntl, brpc::HTTP_STATUS_OK, "ok", body);

        NC_EAC_HTTP_SERVER_TRACE (_T("[ncEACPolicyHandler::CheckPolicy] this: %p, cntl: %p, userId: %s end"), this, cntl, userId.getCStr ());
    NC_EAC_HTTP_SERVER_CATCH (cntl)
}

void
ncEACPolicyHandler::parseUserInfo (brpc::Controller* cntl, ncPolicyCheckInfo& policyInfo, JSON::Value& requestJson)
{
    // 获取http请求的content
    string bodyBuffer = cntl->request_attachment ().to_string ();

    try {
        JSON::Reader::read (requestJson, bodyBuffer.c_str (), bodyBuffer.size ());
    }
    catch (Exception& e) {
        THROW_E (EAC_HTTP_SERVER, EACHTTP_INVALID_REQUEST_JSON_FORMAT,
            LOAD_STRING (_T("IDS_EACHTTP_INVALID_JSON")));
    }

    // 检测 BODY PARAM
    map<String, JsonValueDesc> user;
    user["id"] = JsonValueDesc (JSON::STRING, true);
    user["priority"] = JsonValueDesc (JSON::INTEGER, true);
    user["enabled"] = JsonValueDesc (JSON::BOOLEAN, true);
    user["type"] = JsonValueDesc (JSON::STRING, true);

    map<String, JsonValueDesc> ext;
    ext["account_type"] = JsonValueDesc(JSON::STRING, true);
    ext["client_type"] = JsonValueDesc(JSON::STRING, true);
    ext["login_ip"] = JsonValueDesc(JSON::STRING, true);
    ext["udid"] = JsonValueDesc(JSON::STRING, true);

    map<String, JsonValueDesc> requestObj;
    requestObj["user"] = JsonValueDesc (JSON::OBJECT, true, &user);
    requestObj["client_id"] = JsonValueDesc(JSON::STRING, true);
    requestObj["ip"] = JsonValueDesc(JSON::STRING, true);
    requestObj["ext"] = JsonValueDesc (JSON::OBJECT, true, &ext);

    JsonValueDesc requestValueDesc = JsonValueDesc (JSON::OBJECT, true, &requestObj);
    CheckRequestParameters ("body", requestJson, requestValueDesc);

    // 读取json格式数据
    policyInfo.userId = requestJson["user"]["id"].s().c_str();
    INVALID_USER_ID(policyInfo.userId);
    policyInfo.clientId = requestJson["client_id"].s().c_str();
    if (!ncOIDUtil::IsGUID (policyInfo.clientId)){
        THROW_E (EAC_HTTP_SERVER, EACHTTP_INVALID_PARAMETER, LOAD_STRING (_T("IDS_EACHTTP_INVALID_ARG_VALUE")), "client_id");
    }
    policyInfo.enabled = requestJson["user"]["enabled"].b();
    policyInfo.priority = requestJson["user"]["priority"].i();
    String userType = requestJson["user"]["type"].s().c_str();
    if (userType != "user") {
        THROW_E (EAC_HTTP_SERVER, EACHTTP_INVALID_PARAMETER, LOAD_STRING (_T("IDS_EACHTTP_INVALID_ARG_VALUE")), "user.type");
    }

    if (_accountTypeMap.count(requestJson["ext"]["account_type"].s()) == 0){
        THROW_E (EAC_HTTP_SERVER, EACHTTP_INVALID_PARAMETER, LOAD_STRING (_T("IDS_EACHTTP_INVALID_ARG_VALUE")), "account_type");
    }
    policyInfo.accountType = _accountTypeMap[requestJson["ext"]["account_type"].s()];

    if (_clientTypeMap.count(requestJson["ext"]["client_type"].s()) == 0){
        THROW_E (EAC_HTTP_SERVER, EACHTTP_INVALID_PARAMETER, LOAD_STRING (_T("IDS_EACHTTP_INVALID_ARG_VALUE")), "client_type");
    }
    policyInfo.clientType = _clientTypeMap[requestJson["ext"]["client_type"].s()];

    policyInfo.ip = requestJson["ip"].s().c_str();
    if (_clientTypeMap[requestJson["ext"]["client_type"].s()] != ACSClientType::APP) {
        INVALID_IP_ADDRESS(policyInfo.ip);
    }

    policyInfo.loginIp = requestJson["ext"]["login_ip"].s().c_str();
    if (_clientTypeMap[requestJson["ext"]["client_type"].s()] != ACSClientType::APP) {
        INVALID_IP_ADDRESS(policyInfo.loginIp);
    }
    policyInfo.udid = requestJson["ext"]["udid"].s().c_str();
}
