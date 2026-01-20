/***************************************************************************************************
hydra.cpp:
    Copyright (c) Eisoo Software Inc. (2020), All rights reserved.

Purpose:
    hydra 实现

Author:
    Sunshine.tang@aishu.cn

Creating Time:
    2021-02-22
***************************************************************************************************/
#include <abprec.h>
#include "hydra.h"
#include "drivenadapter.h"
#include <dataapi/ncJson.h>
#include <ossclient/public/ncIOSSClient.h>
#include <ossclient/ossclient.h>

#include "serviceAccessConfig.h"

/* Implementation file */
NS_IMPL_QUERY_INTERFACE1 (hydra, hydraInterface)

// protected
NS_IMETHODIMP_(nsrefcnt) hydra::AddRef (void)
{
    return 1;
}

// protected
NS_IMETHODIMP_(nsrefcnt) hydra::Release (void)
{
    return 1;
}

AB_DEFINE_THREADSAFE_SINGLETON_NO_POOL (hydra)

hydra::hydra (void): _ossClientPtr (0)
{
    NC_DRIVEN_ADAPTER_TRACE (_T("this: %p"), this);

    consentSessionUrl.format(_T("http://%s:%d/admin/oauth2/auth/sessions/consent"),
        ServiceAccessConfig::getInstance()->hydraAdminHost.getCStr(), ServiceAccessConfig::getInstance()->hydraAdminPort);
    loginSessionUrl.format(_T("http://%s:%d/admin/oauth2/auth/sessions/login"),
        ServiceAccessConfig::getInstance()->hydraAdminHost.getCStr(), ServiceAccessConfig::getInstance()->hydraAdminPort);

    // visitorType 类型转换 string->ncTokenVisitorType
    _visitorTypeMap = {{"realname",ncTokenVisitorType::REALNAME},
                        {"anonymous",ncTokenVisitorType::ANONYMOUS},
                        {"business",ncTokenVisitorType::BUSINESS}};

    // accountType 类型转换 string->ncAccountType
    _accountTypeMap = {{"other", ncAccountType::OTHER}, {"id_card", ncAccountType::ID_CARD}};

    // clientType 类型转换 string->ncClientType
    _clientTypeMap = {{"unknown", ncClientType::UNKNOWN}, {"ios", ncClientType::IOS}, {"android", ncClientType::ANDROID},
                        {"windows_phone", ncClientType::WINDOWS_PHONE}, {"windows", ncClientType::WINDOWS}, {"mac_os", ncClientType::MAC_OS},
                        {"web", ncClientType::WEB}, {"mobile_web", ncClientType::MOBILE_WEB}, {"nas", ncClientType::NAS},
                        {"console_web", ncClientType::CONSOLE_WEB}, {"deploy_web", ncClientType::DEPLOY_WEB}, {"linux", ncClientType::LINUX}, {"app", ncClientType::APP}};
}

hydra::~hydra (void)
{
    NC_DRIVEN_ADAPTER_TRACE (_T("this: %p"), this);
}

void hydra::createOSSClient ()
{
    if (nullptr == _ossClientPtr.get ()) {
        nsresult ret;
        nsCOMPtr<ncIOSSClient> ossClient = do_CreateInstance (NC_OSS_CLIENT_CONTRACTID, &ret);
        if (NS_FAILED (ret)){
            THROW_E (ACS_DRIVEN_ADAPTER, FAILED_TO_CREATE_XPCOM_INSTANCE,
                     _T("Failed to create OSSClient: %s(0x%x)"),
                     String::toString ((int64)ret).getCStr (), (int64)ret);
        }
        _ossClientPtr.reset (new nsCOMPtr<ncIOSSClient>(ossClient));
    }
}

String hydra::UrlEncode3986 (const String &input)
{
    String escaped;

    int max = input.getLength();
    for(int i = 0; i < max; ++ i) {
        if ((48 <= input[i] && input[i] <= 57) ||    //0-9
            (65 <= input[i] && input[i] <= 90) ||    //abc...xyz
            (97 <= input[i] && input[i] <= 122) ||   //ABC...XYZ
            (input[i] == '_' || input[i] == '-' || input[i] == '~' || input[i] == '.')
            ) {
                escaped.append (input[i], 1);
        }
        else {
            escaped.append ("%");

            char dig1 = (input[i]&0xF0)>>4;
            char dig2 = (input[i]&0x0F);
            if ( 0 <= dig1 && dig1 <= 9) dig1 += 48;    //0,48inascii
            if (10 <= dig1 && dig1 <=15) dig1 += 65 - 10; //A,97inascii
            if ( 0 <= dig2 && dig2 <= 9) dig2 += 48;
            if (10 <= dig2 && dig2 <=15) dig2 += 65 - 10;

            String r;
            r.append (&dig1, 1);
            r.append (&dig2, 1);

            escaped.append (r);//converts char 255 to string "FF"
        }
    }

    return escaped;
}

/* [notxpcom] void IntrospectToken ([const] in StringRef tokenId, in ncTokenIntrospectInfoRef tokenInfo); */
NS_IMETHODIMP_(void) hydra::IntrospectToken(const String & tokenId, ncTokenIntrospectInfo & tokenInfo)
{
    NC_DRIVEN_ADAPTER_TRACE (_T("[BEGIN] this: %p"), this);

    String token = UrlEncode3986 (tokenId.getCStr ());
    std::string post;
    post.append ("token=");
    post.append (token.getCStr ());

    createOSSClient ();
    ncOSSResponse response;
    vector<string> inHeaders;
    String url;
    url.format(_T("http://%s:%d/admin/oauth2/introspect"),
        ServiceAccessConfig::getInstance()->hydraAdminHost.getCStr(), ServiceAccessConfig::getInstance()->hydraAdminPort);

    inHeaders.push_back ("Content-Type: application/x-www-form-urlencoded");
    (*_ossClientPtr)->Post (url.getCStr (), post, inHeaders, 30, response);
    // 请求失败时，原样抛出错误信息
    if (response.code != 200) {
        if (response.code == 0) {
            THROW_E (ACS_DRIVEN_ADAPTER, DRIVEN_HTTP_SERVER_INTERNAL_ERR, _T("Could not connect to server"));
            SystemLog::getInstance ()->log (__FILE__, __LINE__, ERROR_LOG_TYPE,
                                            _T("Could not connect to index server"));

        }
        else {
            THROW_E (ACS_DRIVEN_ADAPTER, DRIVEN_HTTP_SERVER_INTERNAL_ERR, _T("Code:%d. Cause:%s."), response.code, response.body.c_str());
            SystemLog::getInstance ()->log (__FILE__, __LINE__, ERROR_LOG_TYPE,
                                            _T("Server internal error, code:%d, cause:%s."),
                                            response.code, response.body.c_str ());

        }
    }

    // 封装返回结果
    JSON::Value oauth_token_value;
    JSON::Reader::read (oauth_token_value, response.body.c_str (), response.body.length ());

    // 令牌状态
    tokenInfo.active = oauth_token_value["active"].b ();
    // 无效令牌
    if (!tokenInfo.active){
        NC_DRIVEN_ADAPTER_TRACE (_T("[END] this: %p"), this);
        return;
    }

    // 用户ID
    tokenInfo.userId = oauth_token_value["sub"].s ().c_str ();
    // scope(权限范围)
    tokenInfo.scope = oauth_token_value["scope"].s ().c_str ();
    // 客户端ID
    tokenInfo.clientId = oauth_token_value["client_id"].s ().c_str ();
    // 客户端凭据模式
    if (tokenInfo.userId == tokenInfo.clientId) {
        tokenInfo.visitorType = ncTokenVisitorType::BUSINESS;
        NC_DRIVEN_ADAPTER_TRACE (_T("[END] this: %p"), this);
        return;
    }
    // 以下字段只在非客户端凭据模式时才存在
    // 访问者类型
    tokenInfo.visitorType = _visitorTypeMap[oauth_token_value["ext"]["visitor_type"].s ()];
    // 匿名用户
    if (tokenInfo.visitorType == ncTokenVisitorType::ANONYMOUS){
        NC_DRIVEN_ADAPTER_TRACE (_T("[END] this: %p"), this);
        return;
    }

    // 实名用户
    if (tokenInfo.visitorType == ncTokenVisitorType::REALNAME) {
        // 登陆IP
        tokenInfo.loginIp = oauth_token_value["ext"]["login_ip"].s ().c_str ();
        // 设备ID
        tokenInfo.udid = oauth_token_value["ext"]["udid"].s ().c_str ();
        // 登录账号类型
        tokenInfo.accountType = _accountTypeMap[oauth_token_value["ext"]["account_type"].s ()];
        // 设备类型
        tokenInfo.clientType = _clientTypeMap[oauth_token_value["ext"]["client_type"].s ()];
        NC_DRIVEN_ADAPTER_TRACE (_T("[END] this: %p"), this);
        return;
    }
}

/* [notxpcom] void DeleteConsentAndLogin ([const] in StringRef clientId, [const] in StringRef userId); */
NS_IMETHODIMP_(void) hydra::DeleteConsentAndLogin(const String & clientId, const String & userId)
{
    NC_DRIVEN_ADAPTER_TRACE ("userId: %s begin", userId.getCStr ());

    String deleteConsentUrl;
    if (clientId.isEmpty ())
        deleteConsentUrl.format (_T("%s?subject=%s&all=true"), consentSessionUrl.getCStr (), userId.getCStr ());
    else
        deleteConsentUrl.format (_T("%s?subject=%s&client=%s"), consentSessionUrl.getCStr (), userId.getCStr (),clientId.getCStr ());

    createOSSClient ();
    ncOSSResponse res;
    vector<string> inHeaders;

    (*_ossClientPtr)->Delete (deleteConsentUrl.getCStr (), inHeaders, 30, res);
    if (res.code != 204 && res.code != 404) {
        if (res.code == 0) {
            THROW_E (ACS_DRIVEN_ADAPTER, DRIVEN_HTTP_SERVER_INTERNAL_ERR, _T("Could not connect to server"));
            SystemLog::getInstance ()->log (__FILE__, __LINE__, ERROR_LOG_TYPE,
                                            _T("Could not connect to index server"));
        }
        else {
            THROW_E (ACS_DRIVEN_ADAPTER, DRIVEN_HTTP_SERVER_INTERNAL_ERR, _T("Code:%d. Cause:%s."), res.code, res.body.c_str ());
            SystemLog::getInstance ()->log (__FILE__, __LINE__, ERROR_LOG_TYPE,
                                            _T("Server internal error, code:%d, cause:%s."),
                                            res.code, res.body.c_str ());
        }
    }

    String deleteLoginUrl;
    deleteLoginUrl.format (_T("%s?subject=%s"), loginSessionUrl.getCStr (), userId.getCStr ());
    inHeaders.clear ();
    (*_ossClientPtr)->Delete (deleteLoginUrl.getCStr (), inHeaders, 30, res);
    if (res.code != 204 && res.code != 404) {
        if (res.code == 0) {
            THROW_E (ACS_DRIVEN_ADAPTER, DRIVEN_HTTP_SERVER_INTERNAL_ERR, _T("Could not connect to server"));
            SystemLog::getInstance ()->log (__FILE__, __LINE__, ERROR_LOG_TYPE,
                                            _T("Could not connect to index server"));
        }
        else {
            THROW_E (ACS_DRIVEN_ADAPTER, DRIVEN_HTTP_SERVER_INTERNAL_ERR, _T("Code:%d. Cause:%s."), res.code, res.body.c_str ());
            SystemLog::getInstance ()->log (__FILE__, __LINE__, ERROR_LOG_TYPE,
                                            _T("Server internal error, code:%d, cause:%s."),
                                            res.code, res.body.c_str ());
        }
    }

    NC_DRIVEN_ADAPTER_TRACE ("userId: %s begin", userId.getCStr ());
}

/* [notxpcom] void GetConsentInfo ([const] in StringRef userId, in ncTokenIntrospectInfoVecRef tokenInfos); */
NS_IMETHODIMP_(void) hydra::GetConsentInfo(const String & userId, vector<ncTokenIntrospectInfo> & tokenInfos)
{
    tokenInfos.clear ();

    String consentInfoUrl;
    consentInfoUrl.format (_T("%s?subject=%s"), consentSessionUrl.getCStr (), userId.getCStr ());

    createOSSClient ();
    ncOSSResponse res;
    vector<string> inHeaders;

    (*_ossClientPtr)->Get (consentInfoUrl.getCStr () ,inHeaders, 30, res);
    if (res.code != 200) {
        if (res.code == 0) {
            THROW_E (ACS_DRIVEN_ADAPTER, DRIVEN_HTTP_SERVER_INTERNAL_ERR, _T("Could not connect to server"));
            SystemLog::getInstance ()->log (__FILE__, __LINE__, ERROR_LOG_TYPE,
                                            _T("Could not connect to index server"));
        }
        else {
            THROW_E (ACS_DRIVEN_ADAPTER, DRIVEN_HTTP_SERVER_INTERNAL_ERR, _T("Code:%d. Cause:%s."), res.code, res.body.c_str ());
            SystemLog::getInstance ()->log (__FILE__, __LINE__, ERROR_LOG_TYPE,
                                            _T("Server internal error, code:%d, cause:%s."),
                                            res.code, res.body.c_str ());
        }
    }
    JSON::Value JconsentInfos;
    JSON::Reader::read (JconsentInfos, res.body.c_str (), res.body.length ());
    for(size_t i = 0; i < JconsentInfos.a ().size (); ++ i){
        ncTokenIntrospectInfo tokenInfo;
        tokenInfo.udid = toCFLString (JconsentInfos[i]["session"]["access_token"]["udid"].s ().c_str ());
        tokenInfo.clientId = toCFLString (JconsentInfos[i]["consent_request"]["client"]["client_id"].s ().c_str ());
        tokenInfo.userId = toCFLString (JconsentInfos[i]["consent_request"]["subject"].s ().c_str ());
        tokenInfo.clientType = _clientTypeMap[JconsentInfos[i]["session"]["access_token"]["client_type"].s ()];
        tokenInfos.push_back (tokenInfo);
    }
}
