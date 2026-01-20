#include "common/config.hpp"
#include <abprec.h>
#include <ncutil/ncutil.h>
#include <dataapi/dataapi.h>
#include <teacpserver/public/ncITEACPServer.h>
#include <eachttpserver/eachttpserver.h>
#include <eachttpserver/ncEACHttpServer.h>
#include <eachttpserver/ncEACInnerHttpServer.h>
#include "gen-cpp/EThriftException_types.h"

#define BOOST_SPIRIT_THREADSAFE
#include <boost/property_tree/ptree.hpp>
#include <boost/property_tree/ini_parser.hpp>
using namespace boost::property_tree;

void initDlls ()
{
    static AppContext appCtx (_T("eacpdaemon"));
    ncLanguageManager::InitLangSetting("/sysvol/conf/service_conf/language.conf");

    AppContext::setInstance (&appCtx);
    AppSettings* appSettings = AppSettings::getCFLAppSettings ();
    LibManager::getInstance ()->initLibs (appSettings, &appCtx, 0);
    // ELOGGER_SET_ERROR_PROFILER (DATA_API_ERR_PROVIDER_NAME);

    ::ncInitXPCOM ();
}

int main (int argc, char** argv)
{
    try {
        initDlls ();
        SystemLog::getInstance ()->setAppName("eacpdaemon");

        // bthread并发限制，0表示不限制
        int maxConcurrency = 50;
        try {
            ptree pt;
            boost::property_tree::ini_parser::read_ini("/sysvol/conf/service_conf/app_default.conf", pt);
            maxConcurrency = pt.get<int>("ShareServer.worker_num", 50);
            SystemLog::getInstance ()->log (__FILE__, __LINE__, INFORMATION_LOG_TYPE,
                                            _T("EACP deamon started with max concurrency: %d"), maxConcurrency);
        }
        catch (...) {
        }
        if (maxConcurrency < 0) {
            printMessage2 (_T("Incorrect counts of eacp worker, allowed range is (0, INT_MAX], current setting: %d."), maxConcurrency);
            return 1;
        }

        Config::getInstance();
        ncEACHttpServer httpServer(maxConcurrency);
        ncEACInnerHttpServer innerHttpServer(maxConcurrency);
        try {
            httpServer.Start ();
            innerHttpServer.Start ();
        }
        catch (Exception& e) {
            printMessage2 (_T("Failed to start eacp server:%s"), e.toFullString ().getCStr ());
            SystemLog::getInstance ()->log (__FILE__, __LINE__, ERROR_LOG_TYPE,
                                            _T("EACP deamon started failed:%s"),e.toFullString ().getCStr ());
            return 1;
        }
        catch (ncTException& e) {
            Exception expt (e.fileName.c_str (),
                            e.codeLine,
                            e.expMsg.c_str (),
                            e.errID,
                            SIMPLE_ERROR_PROVIDER (e.errProvider.c_str ()));
            printMessage2(_T("Failed to start eacp server: %s"), expt.toString().getCStr());
            SystemLog::getInstance ()->log (__FILE__, __LINE__, ERROR_LOG_TYPE,
                                            _T("EACP deamon started failed:%s"), expt.toString().getCStr());
            return 1;
        }
        catch (...) {
            printMessage2 (_T("Failed to start HttpServer"));
            return 1;
        }

        nsresult ret;
        nsCOMPtr<ncITEACPServer> tserver = do_CreateInstance (NC_THRIFT_EACP_SERVER_CONTRACTID, &ret);
        if (NS_FAILED (ret)) {
            printMessage2 (_T("Failed to create thrift eacp server: 0x%x"), ret);
            return 1;
        }

        tserver->Start ();
    }
    catch (Exception& e) {
        printMessage2 (_T("Failed to start eacp server: %s"), e.toFullString ().getCStr ());
    }
    catch (...) {
        printMessage2 (_T("Failed to start eacp server: unknown error"));
    }

    return 0;
}
