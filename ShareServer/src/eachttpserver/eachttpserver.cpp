#include <abprec.h>
#include <syswrap/syswrap.h>
#include <ncErrCodeConverter.h>
#include <ncErrCodeMsg.h>

#include "eachttpserver.h"


//
// 资源处理对象
//
IResourceLoader* ncEACHttpServerLoader = 0;

void ncCreateEACHttpServerMoResourceLoader (const AppSettings *appSettings,
                                            const AppContext *appCtx)
{
    if (ncEACHttpServerLoader == 0)
        ANY_NEW1_THROW (ncEACHttpServerLoader,
                        MoResourceLoader,
                        ::getResourceFileName (_T("eachttpserver"),
                                               appSettings,
                                               appCtx,
                                               AB_RESOURCE_MO_EXT_NAME));
}

String ncHttpGetIP (brpc::Controller* cntl)
{
    String ip;
    ncHttpGetHeader (cntl, "X-Forwarded-For", ip);
    if (!ip.isEmpty ()) {
        vector<String> ips;
        ip.split (",", ips);
        ip = ips[0];
    }

    if (ip.isEmpty ()) {
        ip = butil::endpoint2str (cntl->remote_side ()).c_str ();
        int index = ip.rfind (":");
        ip = ip.subString (0, index);

        int start = (ip.find ("[") == String::NO_POSITION) ? 0 : ip.find ("[") + 1;
        int end = (ip.find ("]") == String::NO_POSITION) ? ip.getLength () : ip.find ("]") - start;
        ip = ip.subString (start, end);
    }

    if ((ip == "127.0.0.1") || (ip == "localhost")) {
        String webIp;
        ncHttpGetHeader (cntl, "X-Real-IP", webIp);
        if (!webIp.isEmpty ()) {
            ip = webIp;
        }
    }

    return ip;
}

void ncHttpGetHeader (brpc::Controller* cntl, const String& key, String& value)
{
    const string* pValue = cntl->http_request ().GetHeader (toSTLString (key));
    value = pValue ? toCFLString (*pValue) : String::EMPTY;
}

void ncHttpGetParams (brpc::Controller* cntl, String& method, String& token, String& userId)
{
    try {
        String path = toCFLString(cntl->http_request ().uri ().path ());
        vector<String> tmpNodes;
        vector<String> nodes;
        path.trim ('/').split ('/', tmpNodes);
        for (auto i = 0; i < tmpNodes.size(); i++) {
            if (!tmpNodes[i].isEmpty ()) {
                nodes.push_back (tmpNodes[i]);
            }
        }
        method = nodes[4];

        const string* pTokenId = cntl->http_request ().GetHeader ("Authorization");
        if (pTokenId) {
            String tmp = toCFLString(*pTokenId);
            if (tmp.startsWith ("Bearer ")) {
                token = tmp.subString (7);
            }
            else {
                token = String::EMPTY;
            }
        }
        else {
            pTokenId = cntl->http_request ().uri ().GetQuery ("tokenid");
            token = pTokenId ? URLDecode (toCFLString (*pTokenId)) : String::EMPTY;
        }

        const string* pUserId = cntl->http_request ().uri ().GetQuery ("userid");
        userId = pUserId ? URLDecode (toCFLString (*pUserId)) : String::EMPTY;
    }
    catch (Exception& e) {
        THROW_E (EACHTTP_SERVER_ERR_PROVIDER_NAME,
                 EACHTTP_URI_FORMAT_ERR,
                 LOAD_STRING (_T("IDS_EACHTTP_URI_FORMAT_ERR")));
    }
    catch (...) {
        throw;
    }
}

void ncHttpGetParams (brpc::Controller* cntl, String& method, String& token)
{
    try {
        String path = toCFLString(cntl->http_request ().uri ().path ());
        vector<String> tmpNodes;
        vector<String> nodes;
        path.trim ('/').split ('/', tmpNodes);
        for (auto i = 0; i < tmpNodes.size(); i++) {
            if (!tmpNodes[i].isEmpty ()) {
                nodes.push_back (tmpNodes[i]);
            }
        }
        method = nodes[4];

        const string* pTokenId = cntl->http_request ().GetHeader ("Authorization");
        if (pTokenId) {
            String tmp = toCFLString(*pTokenId);
            if (tmp.startsWith ("Bearer ")) {
                token = tmp.subString (7);
            }
            else {
                token = String::EMPTY;
            }
        }
        else {
            pTokenId = cntl->http_request ().uri ().GetQuery ("tokenid");
            token = pTokenId ? URLDecode (toCFLString (*pTokenId)) : String::EMPTY;
        }

    }
    catch (Exception& e) {
        THROW_E (EACHTTP_SERVER_ERR_PROVIDER_NAME,
                 EACHTTP_URI_FORMAT_ERR,
                 LOAD_STRING (_T("IDS_EACHTTP_URI_FORMAT_ERR")));
    }
    catch (...) {
        throw;
    }
}

void ncHttpGetParams (brpc::Controller* cntl, String& method)
{
    try {
        String path = toCFLString(cntl->http_request ().uri ().path ());
        vector<String> tmpNodes;
        vector<String> nodes;
        path.trim ('/').split ('/', tmpNodes);
        for (auto i = 0; i < tmpNodes.size(); i++) {
            if (!tmpNodes[i].isEmpty ()) {
                nodes.push_back (tmpNodes[i]);
            }
        }
        method = nodes[4];
    }
    catch (Exception& e) {
        THROW_E (EACHTTP_SERVER_ERR_PROVIDER_NAME,
                 EACHTTP_URI_FORMAT_ERR,
                 LOAD_STRING (_T("IDS_EACHTTP_URI_FORMAT_ERR")));
    }
    catch (...) {
        throw;
    }
}

String URLDecode (const String &in)
{
    String out;
    for (size_t i = 0; i < in.getLength ();) {
        if ((char)in[i] == '%') {
            char c = 0;
            for (int ii = 1; ii < 3; ++ii) {
                c = c << 4;
                if (isdigit (in[i+ii])) {
                    c |= (int)in[i+ii] - 48;
                }
                else {
                    c |= (int)in[i+ii] - 55;
                }
            }
            out += c;
            i += 3;
        }
        else if ((char)in[i] == '+') {
            out += ' ';
            ++i;
        }
        else {
            out += in[i];
            ++i;
        }
    }
    return out;
}

void ncHttpGetQueryString (brpc::Controller* cntl, const String& key, String& value)
{
    try {
        const string* pValue = cntl->http_request ().uri ().GetQuery (key.getCStr ());
        value = pValue ? URLDecode (toCFLString (*pValue)) : String::EMPTY;
    }
    catch (Exception& e) {
        THROW_E (EACHTTP_SERVER_ERR_PROVIDER_NAME,
                 EACHTTP_URI_FORMAT_ERR,
                 LOAD_STRING (_T("IDS_EACHTTP_URI_FORMAT_ERR")));
    }
    catch (...) {
        throw;
    }
}

void ncHttpGetToken (brpc::Controller* cntl, String& token)
{
    try {
        const string* pTokenId = cntl->http_request ().GetHeader ("Authorization");
        if (pTokenId) {
            String tmp = toCFLString(*pTokenId);
            if (tmp.startsWith ("Bearer ")) {
                token = tmp.subString (7);
            }
            else {
                token = String::EMPTY;
            }
        }
        else {
            pTokenId = cntl->http_request ().uri ().GetQuery ("tokenid");
            token = pTokenId ? URLDecode (toCFLString (*pTokenId)) : String::EMPTY;
        }
    }
    catch (Exception& e) {
        THROW_E (EACHTTP_SERVER_ERR_PROVIDER_NAME,
                 EACHTTP_URI_FORMAT_ERR,
                 LOAD_STRING (_T("IDS_EACHTTP_URI_FORMAT_ERR")));
    }
    catch (...) {
        throw;
    }
}

void ncHttpReplyException (brpc::Controller* cntl, Exception& e)
{
    NC_EAC_HTTP_SERVER_TRACE (_T("[ERROR] %s"), e.toString ().getCStr ());

    /* 将Exception转化为
        500 Server Internal Error
        {
            "code": 12,
            "message": "Server Internal Error",
            "cause": "创建acstoken组件失败，provider: acsdb，line: 989, file: ncACSTokenManager.cpp"
        }
    */
    ncHttpErrCodeMsg httpErrMsg = ncErrCodeConverter::getInstance ()->ConvToHttpErrCode (
                                                            ::toSTLString (e.getErrorProviderName ()),
                                                            e.getErrorId ());

    JSON::Value exceptJson;
    exceptJson["code"] = httpErrMsg.errCode;
    exceptJson["message"] = httpErrMsg.errMsg.c_str ();
    exceptJson["cause"] = e.toFullString ().getCStr ();

    string body;
    JSON::Writer::write (exceptJson.o (), body);

    ncHttpSendReply (cntl, httpErrMsg.statusCode, httpErrMsg.statusMsg.c_str (), body);
}

void ncHttpReplyException (brpc::Controller* cntl, JSON::Value& e)
{
    NC_EAC_HTTP_SERVER_TRACE (_T("[ERROR] code: %d, body: %s"), e["code"].i (), e["body"].s ().c_str ());

    /*
     * 将其他服务抛出的错误原样转发
    */
    string statusMsg = "";
    ncHttpSendReply (cntl, e["code"].i (), statusMsg.c_str (), e["body"].s ());
}

void ncHttpReplyDetailException (brpc::Controller* cntl, EHttpDetailException& e)
{
    NC_EAC_HTTP_SERVER_TRACE (_T("[ERROR] %s"), e.toString ().getCStr ());

    /* 将Exception转化为
        500 Server Internal Error
        {
            "code": 12,
            "message": "Server Internal Error",
            "cause": "创建acstoken组件失败，provider: acsdb，line: 989, file: ncACSTokenManager.cpp"
        }
    */
    ncHttpErrCodeMsg httpErrMsg = ncErrCodeConverter::getInstance ()->ConvToHttpErrCode (
                                                            ::toSTLString (e.getErrorProviderName ()),
                                                            e.getErrorId ());

    JSON::Value exceptJson;
    exceptJson["code"] = httpErrMsg.errCode;
    exceptJson["message"] = httpErrMsg.errMsg.c_str ();
    exceptJson["cause"] = e.toFullString ().getCStr ();
    exceptJson["detail"] = e.getDetail();

    string body;
    JSON::Writer::write (exceptJson.o (), body);

    ncHttpSendReply (cntl, httpErrMsg.statusCode, httpErrMsg.statusMsg.c_str (), body);
}

void ncHttpSendReply (brpc::Controller* cntl, int code, const char *reason, const string& body)
{
    NC_EAC_HTTP_SERVER_TRY
        cntl->http_response ().set_status_code (code);
        cntl->response_attachment ().append (body);
        cntl->http_response ().AppendHeader ("Content-Length", to_string (body.length ()));
    NC_EAC_HTTP_SERVER_CATCH (cntl)
}

void CheckRequestParameters (const String& key, const JSON::Value & jsonV, JsonValueDesc & jsonValueDesc)
{
    if (jsonV.type () != jsonValueDesc.type){
        THROW_E (EACHTTP_SERVER_ERR_PROVIDER_NAME, EACHTTP_INVALID_PARAMETER,
                 LOAD_STRING (_T("IDS_EACHTTP_DATA_TYPE_IS_ERROR")), key.getCStr (), JSON::get_type_name (jsonValueDesc.type));
    }
    else if (jsonValueDesc.type == JSON::OBJECT){
        for (auto iter = jsonValueDesc.valueDescPtr->begin (); iter != jsonValueDesc.valueDescPtr->end (); iter++) {
            auto tmpIter = jsonV.o ().find (iter->first.getCStr ());
            String newKey;
            newKey.format (_T("%s.%s"), key.getCStr (), iter->first.getCStr ());
            if (tmpIter != jsonV.o ().end ()) {
                iter->second.isExist = true;
                CheckRequestParameters (newKey, tmpIter->second, iter->second);
            }
            else {
                if (iter->second.isRequired) {
                    THROW_E (EACHTTP_SERVER_ERR_PROVIDER_NAME,
                            EACHTTP_INVALID_PARAMETER,
                            LOAD_STRING (_T("IDS_EACHTTP_PARAMETERS_IS_REQUIRED")),
                            newKey.getCStr ());
                }
            }
        }
    }
    else if (jsonValueDesc.type == JSON::ARRAY){
        JSON::Array jsonArray = jsonV.a ();
        for (size_t i = 0; i < jsonArray.size (); ++i) {
            auto iter = jsonValueDesc.valueDescPtr->find ("element");
            if (iter != jsonValueDesc.valueDescPtr->end ()){
                String newKey;
                newKey.format (_T("%s[%d]"), key.getCStr (), i);
                CheckRequestParameters (newKey, jsonArray[i], iter->second);
            }
        }
    }
}

bool CheckToken (const ncCheckTokenInfo & checkTokenInfo, ncIntrospectInfo & introspectInfo)
{
    nsresult ret;
    nsCOMPtr<userManagementInterface> userManager = do_CreateInstance (USER_MANAGEMENT_CONTRACTID, &ret);
    if (NS_FAILED (ret)) {
        THROW_E (EAC_HTTP_SERVER, EACHTTP_SERVER_INIT_USER_MANAGEMENT_ERR,
            _T("Failed to create usermanagement instance: 0x%x"), ret);
    }

    nsCOMPtr<hydraInterface> hydra = do_CreateInstance (HYDRA_CONTRACTID, &ret);
    if (NS_FAILED (ret)) {
        THROW_E (EAC_HTTP_SERVER, EACHTTP_SERVER_INIT_HYDRA_ERR,
            _T("Failed to create hydra adapter instance: 0x%x"), ret);
    }

    nsCOMPtr<ncIACSPolicyManager> policyManager = do_CreateInstance (NC_ACS_POLICY_MANAGER_CONTRACTID, &ret);
    if (NS_FAILED (ret)) {
        THROW_E (EAC_HTTP_SERVER, EACHTTP_SERVER_INIT_ACS_POLICY_MANAGER_ERR,
            _T("Failed to create policyManager instance: 0x%x"), ret);
    }

    ncTokenIntrospectInfo tmpInfo;
    hydra->IntrospectToken(checkTokenInfo.tokenId, tmpInfo);
    // 无效token
    if (!tmpInfo.active) {
        return false;
    }
    introspectInfo.userId = tmpInfo.userId;
    introspectInfo.scope = tmpInfo.scope;
    introspectInfo.clientId = tmpInfo.clientId;
    introspectInfo.visitorType = tmpInfo.visitorType;
    introspectInfo.clientType = tmpInfo.clientType;
    // 业务系统或匿名用户
    if (tmpInfo.visitorType == ncTokenVisitorType::BUSINESS || tmpInfo.visitorType == ncTokenVisitorType::ANONYMOUS) {
        return true;
    }
    // 实名用户
    if (tmpInfo.visitorType == ncTokenVisitorType::REALNAME) {
        // 获取用户信息
        UserInfo userInfo;
        userManager->GetUserInfo(tmpInfo.userId, userInfo);
        // 获取用户角色信息
        introspectInfo.roleIds = set<ncUserRoleType>(userInfo.roles.begin(), userInfo.roles.end());

        // 用户存在，获取用户策略相关信息
        ncPolicyCheckInfo policyInfo;
        policyInfo.userId = tmpInfo.userId;
        policyInfo.priority = userInfo.priority;
        policyInfo.enabled = userInfo.enabled;
        policyInfo.clientId = tmpInfo.clientId;
        policyInfo.ip = checkTokenInfo.ip;
        policyInfo.accountType = static_cast<ACSAccountType>(tmpInfo.accountType);
        policyInfo.clientType = static_cast<ACSClientType>(tmpInfo.clientType);
        policyInfo.loginIp = tmpInfo.loginIp;
        policyInfo.udid = tmpInfo.udid;

        // 用户策略检测
        policyManager->CheckPolicy(policyInfo);
        return true;
    }
    return false;
}

/**
 * Initialize, close, install, or uninstall the lib dataapi.dll.
 */
class ncEACHttpServerlibrary : public ISharedLibrary
{
public:
    ncEACHttpServerlibrary (void)
    {
    }

    virtual ~ncEACHttpServerlibrary (void)
    {
    }

    /**
     * Initialize efshttpserver library and its resource file.
     * This method should be called when application is starting.
     *
     * @param appSettings    Object to get application settings, this object
     *                        saves the default language setting of the app.
     * @param appCtx        Appplication context object, specifiy configuration
     *                        path, resource file path.
     *
     * @throw SharedLibraryException if failed to initialize.
     */
    virtual void onInitLibrary (const AppSettings* appSettings, const AppContext* appCtx)
    {
        if (appSettings == 0)
            throw SharedLibraryException (_T("The object of the application settings is null."));

        if (appCtx == 0)
            throw SharedLibraryException (_T("The object of the application context is null."));

        if (ncEACHttpServerLoader == 0) {
            try {
                ncCreateEACHttpServerMoResourceLoader (appSettings, appCtx);
            }
            catch (Exception& e) {
                throw SharedLibraryException (e.getMessage (),
                                              e.getErrorId (),
                                              e.getErrorProvider ());
            }
        }
    }

    /**
     * Close dataapi library and release allocated resource.
     * This method should be called when application exit.
     */
    virtual void onCloseLibrary (void) AB_NOTHROW
    {
        if (ncEACHttpServerLoader!= 0) {
            delete ncEACHttpServerLoader;
            ncEACHttpServerLoader = 0;
        }
    }

    /**
     * Install dataapi library, can be ignored.
     *
     * @param appSettings    Object to get application settings, this object
     *                        saves the default language setting of the app.
     * @param appCtx        Appplication context object, specifiy configuration
     *                        path, resource file path.
     *
     * @throw SharedLibraryException if failed to install.
     */
    virtual void onInstall (const AppSettings* appSettings, const AppContext* appCtx)
    {
    }

    /**
     * Uninstall dataapi library, can be ignored.
     */
    virtual void onUninstall (void) AB_NOTHROW
    {
    }

    /**
     * Get the library name, return "efshttpserver".
     */
    virtual const tchar_t* getLibName (void) const
    {
        return EAC_HTTP_SERVER;
    }

private:
}; // class ncEACHttpServerlibrary

DEFINE_LIBRARY (ncEACHttpServerlibrary);
