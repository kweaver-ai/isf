#include "eachttpserver.h"
#include "ncEACThirdHandler.h"
#include "user/ncEACThirdUserHandler.h"
#include "org/ncEACThirdOrgHandler.h"
#include "dep/ncEACThirdDepHandler.h"
#include <ehttpserver/ncEHttpUtil.h>
#include <third/ncEACThirdUtil.h>
#include "ncEACHttpServerUtil.h"

ncEACThirdHandler::ncEACThirdHandler ()
    :_userHandler(NULL)
    ,_orgHandler(NULL)
    ,_depHandler(NULL)
{
    NC_EAC_HTTP_SERVER_TRACE (_T("this: %p"), this);
    _userHandler = new ncEACThirdUserHandler();
    _orgHandler = new ncEACThirdOrgHandler();
    _depHandler = new ncEACThirdDepHandler();
}

ncEACThirdHandler::~ncEACThirdHandler (void)
{
    NC_EAC_HTTP_SERVER_TRACE (_T("this: %p"), this);
    if (_userHandler) {
        delete _userHandler;
        _userHandler = NULL;
    }
    if (_orgHandler) {
        delete _orgHandler;
        _orgHandler = NULL;
    }
    if (_depHandler) {
        delete _depHandler;
        _depHandler = NULL;
    }
}

void
ncEACThirdHandler::doThirdRequestHandler (brpc::Controller* cntl)
{
    NC_EAC_HTTP_SERVER_TRY

        // 获取query string
        String method;
        String tokenId;
        ncHttpGetParams (cntl, method, tokenId);

        // method是否设置
        if (method.isEmpty ()) {
            THROW_E (EAC_HTTP_SERVER, EACHTTP_QUERY_METHOD_INVALID,
                LOAD_STRING (_T("IDS_EACHTTP_QUERY_METHOD_INVALID")));
        }

        // token验证
        ncCheckTokenInfo checkTokenInfo;
        checkTokenInfo.tokenId = tokenId;
        checkTokenInfo.ip = ncEACHttpServerUtil::GetForwardedIp(cntl);
        ncIntrospectInfo introspectInfo;
        if (CheckToken (checkTokenInfo, introspectInfo) == false) {
            THROW_E (EAC_HTTP_SERVER, EACHTTP_SERVER_ACCESS_TOKEN_ERR,
                    LOAD_STRING (_T("IDS_EACHTTP_ERR_MSG_ACCESS_TOKEN_ERR")));
        }

        bool findMethod = true;

        // 在_userHandler中查找method是否支持
        if (findMethod) {
            map<String, ncUserMethodFunc>::iterator iter = _userHandler->_methodFuncs.find (method);
            if (iter != _userHandler->_methodFuncs.end ()) {
                ncUserMethodFunc func = iter->second;
                (_userHandler->*func) (cntl, introspectInfo);
                findMethod = false;
            }
        }

        // 在_orgHandler中查找method是否支持
        if (findMethod) {
            map<String, ncOrgMethodFunc>::iterator iter = _orgHandler->_methodFuncs.find (method);
            if (iter != _orgHandler->_methodFuncs.end ()) {
                ncOrgMethodFunc func = iter->second;
                (_orgHandler->*func) (cntl, introspectInfo);
                findMethod = false;
            }
        }

        // 在_depHandler中查找method是否支持
        if (findMethod) {
            map<String, ncDepMethodFunc>::iterator iter = _depHandler->_methodFuncs.find (method);
            if (iter != _depHandler->_methodFuncs.end ()) {
                ncDepMethodFunc func = iter->second;
                (_depHandler->*func) (cntl, introspectInfo);
                findMethod = false;
            }
        }

        if (findMethod) {
            THROW_E (EAC_HTTP_SERVER, EACHTTP_QUERY_METHOD_INVALID,
                LOAD_STRING (_T("IDS_EACHTTP_QUERY_METHOD_NOT_SUPPORT")));
        }

    NC_EAC_HTTP_SERVER_CATCH (cntl)
}
