#ifndef __NC_T_EACP_SERVER_H
#define __NC_T_EACP_SERVER_H

#include <server/TThreadPoolServer.h>

#include "./public/ncITEACPServer.h"

using namespace apache::thrift::server;

/* ncTEACPServer */
class ncTEACPServer : public ncITEACPServer
{
public:
    NS_DECL_ISUPPORTS
    NS_DECL_NCITEACPSERVER

    ncTEACPServer();

private:
    ~ncTEACPServer();
    int32_t _EACPPort;
protected:
};

#endif // __NC_T_EACP_SERVER_H
