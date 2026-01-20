#include <abprec.h>
#include <concurrency/ThreadManager.h>
#include <concurrency/ThreadFactory.h>
#include <protocol/TBinaryProtocol.h>
#include <transport/TServerSocket.h>
#include <transport/TTransportUtils.h>

#include "gen-cpp/ncTEACP.h"
#include "gen-cpp/EACP_constants.h"

#include "teacpserver.h"
#include "ncTEACPHandler.h"
#include "ncTEACPServer.h"
#include <boost/property_tree/ptree.hpp>
#include <boost/property_tree/ini_parser.hpp>

using namespace apache::thrift;
using namespace apache::thrift::protocol;
using namespace apache::thrift::transport;
using namespace apache::thrift::server;
using namespace apache::thrift::concurrency;
using namespace boost;

#define APP_DEFAULT_FILE               _T("/sysvol/conf/service_conf/app_default.conf")

/* Implementation file */
NS_IMPL_THREADSAFE_ISUPPORTS1(ncTEACPServer, ncITEACPServer)

ncTEACPServer::ncTEACPServer()
: _EACPPort (0)
{
    NC_T_EACP_SERVER_TRACE (_T("this: %p"), this);
}

ncTEACPServer::~ncTEACPServer()
{
    NC_T_EACP_SERVER_TRACE (_T("this: %p"), this);
}

/* [notxpcom] void Start (); */
NS_IMETHODIMP_(void) ncTEACPServer::Start()
{
    NC_T_EACP_SERVER_TRACE (_T("this: %p begin"), this);

    try {
        std::shared_ptr<TProtocolFactory> protocolFactory (new TBinaryProtocolFactory ());
        std::shared_ptr<ncTEACPIf> handler (new ncTEACPHandler ());
        std::shared_ptr<TProcessor> processor (new ncTEACPProcessor (handler));

         // 获取thrift端口配置
        boost::property_tree::ptree pt;
        boost::property_tree::ini_parser::read_ini (APP_DEFAULT_FILE, pt);
        _EACPPort = pt.get<int> ("ShareServer.eacp_port");
        int thriftWorkerNumber = pt.get<int>("ShareServer.eacp_thrift_worker_num", 10);
        SystemLog::getInstance ()->log (__FILE__, __LINE__, INFORMATION_LOG_TYPE, _T("Info: eacp_thrift_worker_num is: %d"), thriftWorkerNumber);
        std::shared_ptr<TServerTransport> serverTransport (new TServerSocket ("::", _EACPPort));
        std::shared_ptr<TTransportFactory> transportFactory (new TBufferedTransportFactory ());

        std::shared_ptr<ThreadManager> threadManager = ThreadManager::newSimpleThreadManager (thriftWorkerNumber);
        std::shared_ptr<ThreadFactory> threadFactory = std::shared_ptr<ThreadFactory>(new ThreadFactory ());
        threadManager->threadFactory (threadFactory);
        threadManager->start ();

        auto_ptr<TThreadPoolServer> server (new TThreadPoolServer (processor,
            serverTransport,
            transportFactory,
            protocolFactory,
            threadManager));

        server->serve();
    }
    catch (TException& e) {
        THROW_E (T_EACP_SERVER_ERR_PROVIDER_NAME, CFL_ERR_UNKNOWN, e.what ());
    }

    NC_T_EACP_SERVER_TRACE (_T("this: %p end"), this);
}

/* [notxpcom] void Stop (); */
NS_IMETHODIMP_(void) ncTEACPServer::Stop()
{
    NC_T_EACP_SERVER_TRACE (_T("this: %p"), this);
}
