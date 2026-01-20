/****************************************************************************************************
eacServiceAccessConfig.cpp:
    EacServiceAccessConfig source file.
    Copyright (c) Eisoo Software, Inc.(2021), All rights reserved.

Purpose:
    Source file to implement interface EacServiceAccessConfig.

Author:
    Will.lv

Creating Time:
    2021-04-28
****************************************************************************************************/
#include "eacServiceAccessConfig.h"

#include <boost/property_tree/ptree.hpp>
#include <boost/property_tree/ini_parser.hpp>
using namespace boost::property_tree;

AB_DEFINE_THREADSAFE_SINGLETON_NO_POOL(EacServiceAccessConfig)

EacServiceAccessConfig::EacServiceAccessConfig()
: sharemgntHost()
, sharemgntPort(0)
, eacpThriftHost()
, eacpThriftPort(0)
{
    try {
        ptree pt;
        boost::property_tree::ini_parser::read_ini("/sysvol/conf/service_conf/service_access.conf", pt);
        sharemgntHost =   toCFLString (pt.get<string>("sharemgnt.host"));
        sharemgntPort =   pt.get<int>("sharemgnt.port");
        SystemLog::getInstance ()->log (__FILE__, __LINE__, INFORMATION_LOG_TYPE,
                                    _T("Info : read  sharemgnt thrift configuration, host: %s, port: %d ;")
                                    , sharemgntHost.getCStr (), sharemgntPort);

        eacpThriftHost = toCFLString (pt.get<string>("eacp.thriftHost"));
        eacpThriftPort = pt.get<int>("eacp.thriftPort");
        SystemLog::getInstance ()->log (__FILE__, __LINE__, INFORMATION_LOG_TYPE,
                                    _T("Info : read  eacp thrift configuration, host: %s, port: %d ;")
                                    , eacpThriftHost.getCStr (), eacpThriftPort);
    }
    catch (ptree_error& e) {
        SystemLog::getInstance ()->log (__FILE__, __LINE__, ERROR_LOG_TYPE,
                                    "Get /sysvol/conf/service_conf/service_access.conf failed");
    }
}

EacServiceAccessConfig::~EacServiceAccessConfig()
{

}
