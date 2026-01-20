#include "eachttpserver.h"
#include "ncEACCAuthHandler.h"
#include "ncEACHttpServerUtil.h"

// public
ncEACCAuthHandler::ncEACCAuthHandler()
{
    NC_EAC_HTTP_SERVER_TRACE (_T("this: %p"), this);

    _methodFuncs.insert (pair<String, ncMethodFunc>(_T("get"), &ncEACCAuthHandler::Get));
}

// public
ncEACCAuthHandler::~ncEACCAuthHandler (void)
{
    NC_EAC_HTTP_SERVER_TRACE (_T("this: %p"), this);
}

// public
void
ncEACCAuthHandler::doCARequestHandler (brpc::Controller* cntl)
{
    NC_EAC_HTTP_SERVER_TRY

        // 获取query string
        String method;
        ncHttpGetParams (cntl, method);

        // method是否设置
        if (method.isEmpty ()) {
            THROW_E (EAC_HTTP_SERVER, EACHTTP_QUERY_METHOD_INVALID,
                     LOAD_STRING (_T("IDS_EACHTTP_QUERY_METHOD_INVALID")));
        }

        // method是否支持
        map<String, ncMethodFunc>::iterator iter = _methodFuncs.find (method);
        if (iter == _methodFuncs.end ()) {
            THROW_E (EAC_HTTP_SERVER, EACHTTP_QUERY_METHOD_INVALID,
                     LOAD_STRING (_T("IDS_EACHTTP_QUERY_METHOD_NOT_SUPPORT")));
        }

        // 消息处理
        ncMethodFunc func = iter->second;
        (this->*func) (cntl, "");

    NC_EAC_HTTP_SERVER_CATCH (cntl)
}

// protected
void
ncEACCAuthHandler::Get (brpc::Controller* cntl, const String& userId)
{
    NC_EAC_HTTP_SERVER_TRACE (_T("this: %p, cntl: %p, userId: %s begin"), this, cntl, userId.getCStr ());

    ncTThirdPartyAuthConf config = ncEACHttpServerUtil::Usrm_GetThirdPartyAuth ();

    // 回复
    JSON::Value replyJson;
    if(config.thirdPartyId == "ideabank") {
        /*
            {
                "appId": "hengruitest",
                "appKey": "752794749052377055235260",
                "authServer": "61.177.144.130"
            }
         */
        JSON::Value configJson;
        JSON::Reader::read (configJson, config.config.c_str(), config.config.length());

        replyJson["enable"] = true;
        replyJson["vendor"] = "ideabank";
        replyJson["description"] = "ideabank";
        replyJson["server"] = configJson["authServer"].s();
        replyJson["appid"] = configJson["appId"].s();
        replyJson["appkey"] = configJson["appKey"].s();
    }
    else {
        replyJson["enable"] = false;
    }

    string body;
    JSON::Writer::write (replyJson.o (), body);
    ncHttpSendReply (cntl, brpc::HTTP_STATUS_OK, "ok", body);

    NC_EAC_HTTP_SERVER_TRACE (_T("this: %p, cntl: %p, userId: %s end"), this, cntl, userId.getCStr ());
}
