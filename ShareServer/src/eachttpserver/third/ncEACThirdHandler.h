#ifndef __NC_EAC_Third_HANDLER_H__
#define __NC_EAC_Third_HANDLER_H__

#include <acsprocessor/public/ncIACSMessageManager.h>

class ncEACThirdUserHandler;
class ncEACThirdOrgHandler;
class ncEACThirdDepHandler;
class ncIntrospectInfo;

class ncEACThirdHandler
{
public:
    ncEACThirdHandler ();
    ~ncEACThirdHandler ();

    void doThirdRequestHandler (brpc::Controller* cntl);

private:
    typedef void (ncEACThirdUserHandler::*ncUserMethodFunc) (brpc::Controller*, ncIntrospectInfo &);
    typedef void (ncEACThirdOrgHandler::*ncOrgMethodFunc) (brpc::Controller*, ncIntrospectInfo &);
    typedef void (ncEACThirdDepHandler::*ncDepMethodFunc) (brpc::Controller*, ncIntrospectInfo &);

    ncEACThirdUserHandler*               _userHandler;
    ncEACThirdOrgHandler*                _orgHandler;
    ncEACThirdDepHandler*                _depHandler;
};

#endif  // __NC_EAC_CONFIG_HANDLER_H__
