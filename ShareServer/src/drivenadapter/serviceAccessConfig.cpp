/****************************************************************************************************
serviceAccessConfig.cpp:
    ServiceAccessConfig source file.
    Copyright (c) Eisoo Software, Inc.(2021), All rights reserved.

Purpose:
    Source file to implement interface ServiceAccessConfig.

Author:
    Will.lv

Creating Time:
    2021-04-28
****************************************************************************************************/
#include "serviceAccessConfig.h"

#include <boost/property_tree/ptree.hpp>
#include <boost/property_tree/ini_parser.hpp>
using namespace boost::property_tree;

AB_DEFINE_THREADSAFE_SINGLETON_NO_POOL(ServiceAccessConfig)

ServiceAccessConfig::ServiceAccessConfig()
: deployHost()
, deployPost(0)
, hydraAdminHost()
, hydraAdminPort(0)
, userManagePrivateHost()
, userManagePrivatePort(0)
, policyEngineHost()
, policyEnginePort(0)
, policyMgntHost()
, policyMgntPort(0)
, authenticationPrivateHost()
, authenticationPrivatePort(0)
{
    try {
        ptree pt;
        boost::property_tree::ini_parser::read_ini("/sysvol/conf/service_conf/service_access.conf", pt);
        deployHost = toCFLString (pt.get<std::string> ("deploy-service.HttpHost"));
        deployPost = pt.get<int> ("deploy-service.HttpPort");
        SystemLog::getInstance ()->log (__FILE__, __LINE__, INFORMATION_LOG_TYPE,
                                    _T("Info : read deploy configuration, host: %s, port: %d"), deployHost.getCStr (), deployPost);

        hydraAdminHost = toCFLString(pt.get<std::string> ("hydra.administrativeHost"));
        hydraAdminPort = pt.get<int> ("hydra.administrativePort");
        SystemLog::getInstance ()->log (__FILE__, __LINE__, INFORMATION_LOG_TYPE,
                                    _T("Info : read hydra admin configuration, host: %s, port: %d"), hydraAdminHost.getCStr (), hydraAdminPort);

        userManagePrivateHost = toCFLString(pt.get<std::string> ("user-management.privateHost"));
        userManagePrivatePort = pt.get<int> ("user-management.privatePort");
        SystemLog::getInstance ()->log (__FILE__, __LINE__, INFORMATION_LOG_TYPE,
                                    _T("Info : read user-management private configuration, host: %s, port: %d"), userManagePrivateHost.getCStr (), userManagePrivatePort);

        try {
            ossgatewayPrivateProtocol =  toCFLString (pt.get<string>("ossgateway.privateProtocol"));
            ossgatewayPrivateHost =  toCFLString (pt.get<string>("ossgateway.privateHost"));
            ossgatewayPrivatePort =  pt.get<int>("ossgateway.privatePort");
            SystemLog::getInstance ()->log (__FILE__, __LINE__, INFORMATION_LOG_TYPE,
                                                _T("Info : read  ossgateway private configuration, protcol: %s, host: %s, port: %d ;")
                                                , ossgatewayPrivateProtocol.getCStr(), ossgatewayPrivateHost.getCStr (), ossgatewayPrivatePort);
        } catch (ptree_error& e) {
            SystemLog::getInstance ()->log (__FILE__, __LINE__, ERROR_LOG_TYPE,
                                        "ServiceAccessConfig get ossgateway error");
        }

        policyEngineHost = toCFLString(pt.get<std::string> ("proton-policy-engine.host"));
        policyEnginePort = pt.get<int> ("proton-policy-engine.port");
        SystemLog::getInstance ()->log (__FILE__, __LINE__, INFORMATION_LOG_TYPE,
                                    _T("Info : read proton-policy-engine configuration, host: %s, port: %d"), policyEngineHost.getCStr (), policyEnginePort);

        policyMgntHost = toCFLString (pt.get<std::string> ("policy-management.host"));
        policyMgntPort = pt.get<int> ("policy-management.port");
        SystemLog::getInstance ()->log (__FILE__, __LINE__, INFORMATION_LOG_TYPE,
                                            _T("Info : read  policy-management configuration, host: %s, port: %d ;")
                                            , policyMgntHost.getCStr (), policyMgntPort);

        authenticationPrivateHost = toCFLString(pt.get<std::string> ("authentication.privateHost"));
        authenticationPrivatePort = pt.get<int> ("authentication.privatePort");
        SystemLog::getInstance ()->log (__FILE__, __LINE__, INFORMATION_LOG_TYPE,
                                    _T("Info : authentication private configuration, host: %s, port: %d"), authenticationPrivateHost.getCStr (), authenticationPrivatePort);
    }
    catch (ptree_error& e) {
        SystemLog::getInstance ()->log (__FILE__, __LINE__, ERROR_LOG_TYPE,
                                    "Get /sysvol/conf/service_conf/service_access.conf failed");
    }
}

ServiceAccessConfig::~ServiceAccessConfig()
{

}
