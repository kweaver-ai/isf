/***************************************************************************************************
stdLog.cpp:
    Copyright (c) Eisoo Software Inc. (2020), All rights reserved.

Purpose:
    stdLog manager

Author:
    xu.zhi@aishu.cn

Creating Time:
    2022-08-02
***************************************************************************************************/
#include <abprec.h>
#include "stdLog.h"
#include "drivenadapter.h"
#include <dataapi/ncJson.h>

/* Implementation file */
NS_IMPL_QUERY_INTERFACE1 (stdLog, stdLogInterface)

// protected
NS_IMETHODIMP_(nsrefcnt) stdLog::AddRef (void)
{
    return 1;
}

// protected
NS_IMETHODIMP_(nsrefcnt) stdLog::Release (void)
{
    return 1;
}

AB_DEFINE_THREADSAFE_SINGLETON_NO_POOL (stdLog)

stdLog::stdLog (void)
{
    NC_DRIVEN_ADAPTER_TRACE (_T("this: %p"), this);

    // clientType 类型转换 ncClientType->String
    _clientTypeIntToStr[ncClientType::UNKNOWN] = _T("unknown");             // 未知类型
    _clientTypeIntToStr[ncClientType::WINDOWS] = _T("windows");             // 同步盘/富客户端
    _clientTypeIntToStr[ncClientType::IOS] = _T("ios");                     // IOS客户端
    _clientTypeIntToStr[ncClientType::ANDROID] = _T("android");             // Andriod客户端
    _clientTypeIntToStr[ncClientType::MAC_OS] = _T("mac_os");               // Mac客户端
    _clientTypeIntToStr[ncClientType::WEB] = _T("web");                     // 桌面Web客户端
    _clientTypeIntToStr[ncClientType::MOBILE_WEB] = _T("mobile_web");       // 移动Web客户端
    _clientTypeIntToStr[ncClientType::WINDOWS_PHONE] = _T("windows_phone");
    _clientTypeIntToStr[ncClientType::NAS] = _T("nas");
    _clientTypeIntToStr[ncClientType::CONSOLE_WEB] = _T("console_web");
    _clientTypeIntToStr[ncClientType::DEPLOY_WEB] = _T("deploy_web");
    _clientTypeIntToStr[ncClientType::LINUX] = _T("linux");
    _clientTypeIntToStr[ncClientType::APP] = _T("app");

    _mqClient = ncMQClient::GetConnectorFromFile("/sysvol/conf/service_conf/mq_config.yaml");
}

stdLog::~stdLog (void)
{
    NC_DRIVEN_ADAPTER_TRACE (_T("this: %p"), this);
}

// 获取描述信息
String stdLog::getLogDescription (const String& actorName, ncClientType clientType)
{
    String description;
    switch (clientType)
    {
    case ncClientType::IOS:
        description.format (DRIVEN_LOAD_STRING (_T("IDS_USER_LOGIN_SUCCESS")), actorName.getCStr (), "iOS");
        break;
    case ncClientType::ANDROID:
        description.format (DRIVEN_LOAD_STRING (_T("IDS_USER_LOGIN_SUCCESS")), actorName.getCStr (), "Android");
        break;
    case ncClientType::WINDOWS_PHONE:
        description.format (DRIVEN_LOAD_STRING (_T("IDS_USER_LOGIN_SUCCESS")), actorName.getCStr (), "WindowsPhone");
        break;
    case ncClientType::WINDOWS:
        description.format (DRIVEN_LOAD_STRING (_T("IDS_USER_LOGIN_SUCCESS")), actorName.getCStr (), "Windows");
        break;
    case ncClientType::MAC_OS:
        description.format (DRIVEN_LOAD_STRING (_T("IDS_USER_LOGIN_SUCCESS")), actorName.getCStr (), "MacOSX");
        break;
    case ncClientType::WEB:
        description.format (DRIVEN_LOAD_STRING (_T("IDS_USER_LOGIN_SUCCESS")), actorName.getCStr (), "Web");
        break;
    case ncClientType::MOBILE_WEB:
        description.format (DRIVEN_LOAD_STRING (_T("IDS_USER_LOGIN_SUCCESS")), actorName.getCStr (), "MobileWeb");
        break;
    case ncClientType::NAS:
        description.format (DRIVEN_LOAD_STRING (_T("IDS_USER_CONSOLE_LOGIN_SUCCESS")), actorName.getCStr (), DRIVEN_LOAD_STRING (_T("IDS_NAS_GATEWAY")));
        break;
    case ncClientType::CONSOLE_WEB:
        description.format (DRIVEN_LOAD_STRING (_T("IDS_USER_CONSOLE_LOGIN_SUCCESS")), actorName.getCStr (), DRIVEN_LOAD_STRING (_T("IDS_CONSOLE_WEB")));
        break;
    case ncClientType::DEPLOY_WEB:
        description.format (DRIVEN_LOAD_STRING (_T("IDS_USER_CONSOLE_LOGIN_SUCCESS")), actorName.getCStr (), DRIVEN_LOAD_STRING (_T("IDS_DEPLOY_WEB")));
        break;
    case ncClientType::LINUX:
        description.format (DRIVEN_LOAD_STRING (_T("IDS_USER_LOGIN_SUCCESS")), actorName.getCStr (), "Linux");
        break;
    default:
        description.format (DRIVEN_LOAD_STRING (_T("IDS_USER_LOGIN_SUCCESS")), actorName.getCStr (), "Unknown");
        break;
    }
    return description;
}


// 格式化业务日志
void stdLog::formatLog (const ncOperator& actor, JSON::Value& eventJson)
{
    eventJson.clear ();
    // 操作者信息
    eventJson["recorder"] = "AnyShare";
    JSON::Object& operatorJson = eventJson["operator"].o ();
    operatorJson["type"] = "authenticated_user";
    operatorJson["id"] = actor.id.getCStr ();
    operatorJson["name"] = actor.name.getCStr ();
    JSON::Object& agentJson = operatorJson["agent"].o ();
    agentJson["type"] = _clientTypeIntToStr[actor.clientType].getCStr ();
    agentJson["ip"] = actor.ip.getCStr ();
    agentJson["udid"] = actor.udid.getCStr ();

    JSON::Array& departPathArray = operatorJson["department_path"].a ();
    JSON::Object deparmentJson;
    for (int i = 0; i < actor.departmentPaths.size(); i++){
        deparmentJson["id_path"] = actor.departmentPaths[i].getCStr ();
        deparmentJson["name_path"] = actor.departmentNames[i].getCStr ();
        departPathArray.push_back(deparmentJson);
    }

    // 操作对象信息
    eventJson["operation"] = "login";
    eventJson["description"] = getLogDescription (actor.name, actor.clientType).getCStr ();
    JSON::Object& object = eventJson["detail"].o ();
    object["result"] = "success";
    object["reason"] = "";

    // 改造增加log_from信息
    JSON::Object& logForm = eventJson["log_from"].o ();
    logForm["package"] = "IdentifyAndAuthentication";
    JSON::Object& serviceInfo = logForm["service"].o ();
    serviceInfo["name"] = "eacp";

    return;
}

/* [notxpcom] void Log ([const] in ncOperatorRef actor); */
NS_IMETHODIMP_(void) stdLog::Log (const ncOperator& actor)
{
    NC_DRIVEN_ADAPTER_TRACE (_T("[BEGIN] this: %p"), this);
    JSON::Value eventJson;
    formatLog (actor, eventJson);

    string log;
    JSON::Writer::write (eventJson.o (), log);

    try {
        _mqClient->Pub(TELEMETRY_LOG_USER_OPERATION, log.c_str());
    } catch (Exception& e) {
        THROW_E (ACS_DRIVEN_ADAPTER, DRIVEN_HTTP_SERVER_INTERNAL_ERR,
                _T("MQ pub message failed, topic user_login, errorId %d, err %s"),
                e.getErrorId (), e.toFullString ().getCStr ());
    }

    NC_DRIVEN_ADAPTER_TRACE (_T("[END] this: %p"), this);
}
