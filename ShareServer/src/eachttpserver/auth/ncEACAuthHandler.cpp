#include <ncutil/ncBusinessDate.h>

#include "eachttpserver.h"
#include "ncEACAuthHandler.h"
#include "gen-cpp/ShareMgnt_constants.h"
#include <ehttpserver/ncEHttpUtil.h>
#include <ethriftutil/ncThriftClient.h>
#include <ehttpclient/public/ncIEHTTPClient.h>
#include <openssl/rsa.h>
#include <openssl/pem.h>
#include <openssl/err.h>
#include <openssl/evp.h>
#include <openssl/des.h>
#include "eacServiceAccessConfig.h"
#include <boost/property_tree/ptree.hpp>
#include <boost/algorithm/string.hpp>
#include <boost/property_tree/xml_parser.hpp>
using namespace boost::property_tree;

// public
ncEACAuthHandler::ncEACAuthHandler (ncIACSTokenManager* acsTokenManager,
                                    ncIACSShareMgnt* acsShareMgnt,
                                    ncIACSDeviceManager* acsDeviceManager,
                                    ncIACSConfManager* acsConfManager,
                                    ncIACSMessageManager* acsMessageManager,
                                    ncIACSPolicyManager* acsPolicyManager,
                                    policyEngineInterface* policyEngine)
    : _acsTokenManager (acsTokenManager),
        _acsShareMgnt (acsShareMgnt),
        _acsDeviceManager(acsDeviceManager),
        _acsConfManager(acsConfManager),
        _acsMessageManager (acsMessageManager),
        _acsPolicyManager (acsPolicyManager),
        _policyEngine (policyEngine),
        _expires (3600)
{
    NC_EAC_HTTP_SERVER_TRACE (_T("this: %p"), this);

    // auth v1版本支持的协议
    _methodFuncs.insert (pair<String, ncMethodFunc>(_T("getconfig"), &ncEACAuthHandler::GetConfig));
    _methodFuncs.insert (pair<String, ncMethodFunc>(_T("getbyntlmv1"), &ncEACAuthHandler::GetByNTLMV1));
    _methodFuncs.insert (pair<String, ncMethodFunc>(_T("getbyticket"), &ncEACAuthHandler::GetByTicket));
    _methodFuncs.insert (pair<String, ncMethodFunc>(_T("getbyadsession"), &ncEACAuthHandler::GetByADSession));
    _methodFuncs.insert (pair<String, ncMethodFunc>(_T("modifypassword"), &ncEACAuthHandler::ModifyPassword));
    _methodFuncs.insert (pair<String, ncMethodFunc>(_T("validatesecuritydevice"), &ncEACAuthHandler::ValidateSecurityDevice));
    _methodFuncs.insert (pair<String, ncMethodFunc>(_T("checkuninstallpwd"), &ncEACAuthHandler::CheckUninstallPwd));
    _methodFuncs.insert (pair<String, ncMethodFunc>(_T("checkexitpwd"), &ncEACAuthHandler::CheckExitPwd));
    _methodFuncs.insert (pair<String, ncMethodFunc>(_T("getvcode"), &ncEACAuthHandler::GetVcode));
    _methodFuncs.insert (pair<String, ncMethodFunc>(_T("sendsms"), &ncEACAuthHandler::SendSms));
    _methodFuncs.insert (pair<String, ncMethodFunc>(_T("sendvcode"), &ncEACAuthHandler::SendVcode));
    _methodFuncs.insert (pair<String, ncMethodFunc>(_T("smsactivate"), &ncEACAuthHandler::SmsActivate));
    _methodFuncs.insert (pair<String, ncMethodFunc>(_T("servertime"), &ncEACAuthHandler::ServerTime));
    _methodFuncs.insert (pair<String, ncMethodFunc>(_T("sendauthvcode"), &ncEACAuthHandler::SendAuthVcode));
    _methodFuncs.insert (pair<String, ncMethodFunc>(_T("pwd-retrieval-vcode"), &ncEACAuthHandler::SendPwdRetrevalVcode));
    _methodFuncs.insert (pair<String, ncMethodFunc>(_T("configs"), &ncEACAuthHandler::GetServerConfig));
    _methodFuncs.insert (pair<String, ncMethodFunc>(_T("login-configs"), &ncEACAuthHandler::GetLoginConfigs));

    // auth v2版本支持的协议
    _v2MethodFuncs.insert (pair<String, ncMethodFunc>(_T("getconfig"), &ncEACAuthHandler::GetConfig));
    _v2MethodFuncs.insert (pair<String, ncMethodFunc>(_T("login"), &ncEACAuthHandler::Login));
    _v2MethodFuncs.insert (pair<String, ncMethodFunc>(_T("modifypassword"), &ncEACAuthHandler::ModifyPassword));
    _v2MethodFuncs.insert (pair<String, ncMethodFunc>(_T("validatesecuritydevice"), &ncEACAuthHandler::ValidateSecurityDevice));
    _v2MethodFuncs.insert (pair<String, ncMethodFunc>(_T("checkuninstallpwd"), &ncEACAuthHandler::CheckUninstallPwd));
    _v2MethodFuncs.insert (pair<String, ncMethodFunc>(_T("checkexitpwd"), &ncEACAuthHandler::CheckExitPwd));

    _thirdAuthFuncs.insert(pair<String, ncThirdAuthFunc>(_T("anyshare_plain"), &ncEACAuthHandler::AnySharePlain));
    _thirdAuthFuncs.insert(pair<String, ncThirdAuthFunc>(_T("anyshare_rsa"), &ncEACAuthHandler::AnyShareRSA));
    _thirdAuthFuncs.insert(pair<String, ncThirdAuthFunc>(_T("anyshare"), &ncEACAuthHandler::AnyShareSSO));
    _thirdAuthFuncs.insert(pair<String, ncThirdAuthFunc>(_T("windows_ad_sso"), &ncEACAuthHandler::WindowsADSSO));
    _thirdAuthFuncs.insert(pair<String, ncThirdAuthFunc>(_T("bqkj"), &ncEACAuthHandler::BQKJAuth));
    _thirdAuthFuncs.insert(pair<String, ncThirdAuthFunc>(_T("being"), &ncEACAuthHandler::BeingAuth));
    _thirdAuthFuncs.insert(pair<String, ncThirdAuthFunc>(_T("ths"), &ncEACAuthHandler::ThsAuth));
    _thirdAuthFuncs.insert(pair<String, ncThirdAuthFunc>(_T("lcsoft"), &ncEACAuthHandler::LcsoftAuth));
    _thirdAuthFuncs.insert(pair<String, ncThirdAuthFunc>(_T("aishu"), &ncEACAuthHandler::AISHUAuth));


    //第三方配置信息过滤关键字
    _thirdConfigBlackList.insert("host");
    _thirdConfigBlackList.insert("port");
    _thirdConfigBlackList.insert("pwd");
    _thirdConfigBlackList.insert("password");
    _thirdConfigBlackList.insert("database");
    _thirdConfigBlackList.insert("db");
    _thirdConfigBlackList.insert("user");
    _thirdConfigBlackList.insert("sid");

    // 不允许登录的管理员帐号
    _adminIds.insert(g_ShareMgnt_constants.NCT_USER_ADMIN);
    _adminIds.insert(g_ShareMgnt_constants.NCT_USER_AUDIT);
    _adminIds.insert(g_ShareMgnt_constants.NCT_USER_SYSTEM);
    _adminIds.insert(g_ShareMgnt_constants.NCT_USER_SECURIT);

    // 需要鉴权验证的method方法
    _needAuthMethod.insert("configs");

    //http response和request的枚举参数需要进行枚举与string互转，初始化对应map
    vector<String> vecClientTypeStr{"unknown", "ios", "android", "windows_phone","windows",
     "mac_os", "web", "mobile_web", "nas","console_web", "deploy_web", "linux", "app"};

     for( int i = 0; i < vecClientTypeStr.size(); i++ )
     {
        _clientStringTypeMap.insert(make_pair(static_cast<ACSClientType>(i), vecClientTypeStr[i]));
        _clientIntTypeMap.insert(make_pair(vecClientTypeStr[i], i));
     }

    _accountStringTypeMap = map<int, String>{make_pair(0,"other"), make_pair(1,"id_card")};
    _visitorStringTypeMap = map<int, String>{make_pair(1,"realname"), make_pair(4,"anonymous"),make_pair(6,"business")};

}

// public
ncEACAuthHandler::~ncEACAuthHandler (void)
{
    NC_EAC_HTTP_SERVER_TRACE (_T("this: %p"), this);
}

// public
void
ncEACAuthHandler::doAuthRequestHandler (brpc::Controller* cntl)
{
    NC_EAC_HTTP_SERVER_TRY

        // 获取query string
        String method;
        String tokenId;
        String userId;
        ncHttpGetParams (cntl, method, tokenId);

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

        if (_needAuthMethod.find (method.getCStr ()) != _needAuthMethod.end ()) {
            // token验证
            ncCheckTokenInfo checkTokenInfo;
            checkTokenInfo.tokenId = tokenId;
            checkTokenInfo.ip = ncEACHttpServerUtil::GetForwardedIp(cntl);
            ncIntrospectInfo introspectInfo;
            if (CheckToken (checkTokenInfo, introspectInfo) == false) {
                THROW_E (EAC_HTTP_SERVER, EACHTTP_SERVER_ACCESS_TOKEN_ERR,
                    LOAD_STRING (_T("IDS_EACHTTP_ERR_MSG_ACCESS_TOKEN_ERR")));
            }

            // 获取该token对应的userid
            userId = introspectInfo.userId;
        }

        // 消息处理
        ncMethodFunc func = iter->second;
        (this->*func) (cntl, userId);

    NC_EAC_HTTP_SERVER_CATCH (cntl)
}

// public
void
ncEACAuthHandler::doAuth2RequestHandler (brpc::Controller* cntl)
{
    NC_EAC_HTTP_SERVER_TRY

        // 获取query string
        String method;
        String tokenId;
        String userId;
        ncHttpGetParams (cntl, method, tokenId);

        // method是否设置
        if (method.isEmpty ()) {
            THROW_E (EAC_HTTP_SERVER, EACHTTP_QUERY_METHOD_INVALID,
                     LOAD_STRING (_T("IDS_EACHTTP_QUERY_METHOD_INVALID")));
        }

        // method是否支持
        map<String, ncMethodFunc>::iterator iter = _v2MethodFuncs.find (method);
        if (iter == _v2MethodFuncs.end ()) {
            THROW_E (EAC_HTTP_SERVER, EACHTTP_QUERY_METHOD_INVALID,
                LOAD_STRING (_T("IDS_EACHTTP_QUERY_METHOD_NOT_SUPPORT")));
        }

        if (_needAuthMethod.find (method.getCStr ()) != _needAuthMethod.end ()){
            // token验证
            ncCheckTokenInfo checkTokenInfo;
            checkTokenInfo.tokenId = tokenId;
            checkTokenInfo.ip = ncEACHttpServerUtil::GetForwardedIp(cntl);
            ncIntrospectInfo introspectInfo;
            if (CheckToken (checkTokenInfo, introspectInfo) == false) {
                THROW_E (EAC_HTTP_SERVER, EACHTTP_SERVER_ACCESS_TOKEN_ERR,
                    LOAD_STRING (_T("IDS_EACHTTP_ERR_MSG_ACCESS_TOKEN_ERR")));
            }

            // 获取该token对应的userId
            userId = introspectInfo.userId;
        }

        // 消息处理
        ncMethodFunc func = iter->second;
        (this->*func) (cntl, userId);

    NC_EAC_HTTP_SERVER_CATCH (cntl)
}


// public
void
ncEACAuthHandler::setExpires (int64 expires)
{
    _expires = expires;
}

string
ncEACAuthHandler::GetSection (String& languageConfig)
{
    string language_config = toSTLString (languageConfig);
    string section("shareweb_");
    int N = language_config.length();
    for (int i = 0; i < N; i++)
    {
        if (language_config[i] >= 'A' && language_config[i] <= 'Z') {
            language_config[i] += ('a' - 'A');
        }
    }

    if (language_config.find("zh-tw") != string::npos || language_config.find("zh_tw") != string::npos) {
        section += "zh-tw";
    }
    else if (language_config.find("en-us") != string::npos || language_config.find("en_us") != string::npos) {
        section += "en-us";
    }
    else {
        section += "zh-cn";
    }
    return section;
}

void ncEACAuthHandler::GetLoginConfigs(brpc::Controller *cntl, const String &fakeUserId)
{
    NC_EAC_HTTP_SERVER_TRY
    // 判断请求方法是否为GET
    if (cntl->http_request().method() != brpc::HTTP_METHOD_GET)
    {
        THROW_E(EACHTTP_SERVER_ERR_PROVIDER_NAME, EACHTTP_SERVER_METHOD_NOT_IMPLEMENTED,
                LOAD_STRING(_T("IDS_EACHTTP_METHOD_INVALID")));
    }

    JSON::Value replyJson;
    // 从anyshare t_conf表中批量获取配置
    vector<String> keys;
    keys.push_back("oem_remember_pass");

    map<String, String> kvMap;
    _acsConfManager->BatchGetConfig(keys, kvMap);

    // 返回oemconfig
    replyJson["oemconfig"]["rememberpass"] = kvMap["oem_remember_pass"].compareIgnoreCase("true") == 0;

    // sharemgnt t_conf表中批量获取配置
    vector<String> sharemgntKeys;
    map<String, String> sharemgntKvMap;
    sharemgntKeys.push_back("windows_ad_sso");
    sharemgntKeys.push_back("enable_secret_mode");
    sharemgntKeys.push_back("vcode_login_config");
    sharemgntKeys.push_back("strong_pwd_status");
    sharemgntKeys.push_back("strong_pwd_length");
    sharemgntKeys.push_back("vcode_server_status");
    sharemgntKeys.push_back("dualfactor_auth_server_status");
    _acsShareMgnt->BatchGetConfig(sharemgntKeys, sharemgntKvMap);

    // windows ad单点登录
    replyJson["windows_ad_sso"]["is_enabled"] = sharemgntKvMap["windows_ad_sso"].compareIgnoreCase("1") == 0;
    // 开启涉密模式
    replyJson["enable_secret_mode"] = sharemgntKvMap["enable_secret_mode"].compareIgnoreCase("1") == 0;
    // 开启强密码配置
    replyJson["enable_strong_pwd"] = sharemgntKvMap["strong_pwd_status"].compareIgnoreCase("1") == 0;

    // 获取全部配置
    ncTAllConfig allConfig;
    ncEACHttpServerUtil::GetAllConfig(allConfig);

    // 第三方认证配置
    ncTThirdPartyAuthConf &config = allConfig.thirdPartyAuthConf;
    if (config.enabled)
    {
        JSON::Value configJson;
        JSON::Reader::read(configJson, config.config.c_str(), config.config.length());
        replyJson["thirdauth"]["id"] = config.thirdPartyId;

        for (auto iter = configJson.o().begin(); iter != configJson.o().end(); ++iter)
        {
            // 过滤黑名单中信息
            if (_thirdConfigBlackList.find(iter->first) != _thirdConfigBlackList.end())
                continue;
            replyJson["thirdauth"]["config"][iter->first] = iter->second;
        }
    }

    // 获取验证码配置
    string vcode_config = sharemgntKvMap["vcode_login_config"].getCStr();
    JSON::Value vcodeJson;
    try
    {
        JSON::Reader::read(vcodeJson, vcode_config.c_str(), vcode_config.length());
    }
    catch (Exception &e)
    {
        THROW_E(EAC_HTTP_SERVER, EACHTTP_INVALID_REQUEST_JSON_FORMAT,
                LOAD_STRING(_T("IDS_EACHTTP_INVALID_JSON")));
    }
    replyJson["vcode_login_config"]["isenable"] = vcodeJson["isEnable"];
    replyJson["vcode_login_config"]["passwderrcnt"] = vcodeJson["passwdErrCnt"];

    // 获取发送验证码服务器开关配置
    string vcode_server_status = sharemgntKvMap["vcode_server_status"].getCStr();
    JSON::Value sendVcodeTypeStatusJson;
    try
    {
        JSON::Reader::read(sendVcodeTypeStatusJson, vcode_server_status.c_str(), vcode_server_status.length());
    }
    catch (Exception &e)
    {
        THROW_E(EAC_HTTP_SERVER, EACHTTP_INVALID_REQUEST_JSON_FORMAT,
                LOAD_STRING(_T("IDS_EACHTTP_INVALID_JSON")));
    }
    replyJson["vcode_server_status"]["send_vcode_by_sms"] = sendVcodeTypeStatusJson["send_vcode_by_sms"];
    replyJson["vcode_server_status"]["send_vcode_by_email"] = sendVcodeTypeStatusJson["send_vcode_by_email"];

    // 获取双因子验证开关配置
    string mfa_server_status = sharemgntKvMap["dualfactor_auth_server_status"].getCStr();
    JSON::Value dualauthJson;
    try
    {
        JSON::Reader::read(dualauthJson, mfa_server_status.c_str(), mfa_server_status.length());
    }
    catch (Exception &e)
    {
        THROW_E(EAC_HTTP_SERVER, EACHTTP_INVALID_REQUEST_JSON_FORMAT,
                LOAD_STRING(_T("IDS_EACHTTP_INVALID_JSON")));
    }
    replyJson["dualfactor_auth_server_status"]["auth_by_sms"] = dualauthJson["auth_by_sms"];
    replyJson["dualfactor_auth_server_status"]["auth_by_email"] = dualauthJson["auth_by_email"];
    replyJson["dualfactor_auth_server_status"]["auth_by_OTP"] = dualauthJson["auth_by_OTP"];
    replyJson["dualfactor_auth_server_status"]["auth_by_Ukey"] = dualauthJson["auth_by_Ukey"];

    // 如果是内网 则不开启 短信
    String ip = ncEACHttpServerUtil::GetForwardedIp(cntl);
    if (checkIpIsOutNet(ip.getCStr()) == false) {
        replyJson["dualfactor_auth_server_status"]["auth_by_sms"] = false;
    }

    // 获取强密码最小长度
    replyJson["strong_pwd_length"] = Int::getValue(sharemgntKvMap["strong_pwd_length"]);

    string body;
    JSON::Writer::write(replyJson.o(), body);
    ncHttpSendReply(cntl, brpc::HTTP_STATUS_OK, "ok", body);

    NC_EAC_HTTP_SERVER_CATCH (cntl)
}

void ncEACAuthHandler::GetServerConfig(brpc::Controller *cntl, const String &fakeUserId)
{
    NC_EAC_HTTP_SERVER_TRY
    // 判断请求方法是否为GET
    if (cntl->http_request().method() != brpc::HTTP_METHOD_GET)
    {
        THROW_E(EACHTTP_SERVER_ERR_PROVIDER_NAME, EACHTTP_SERVER_METHOD_NOT_IMPLEMENTED,
                LOAD_STRING(_T("IDS_EACHTTP_METHOD_INVALID")));
    }

    JSON::Value replyJson;
    // 从anyshare t_conf表中批量获取配置
    vector<String> keys;
    keys.push_back("oem_client_logout_time");
    keys.push_back("internal_link_prefix");
    keys.push_back("oem_max_pass_expired_days");

    map<String, String> kvMap;
    _acsConfManager->BatchGetConfig(keys, kvMap);

    // 返回oemconfig
    replyJson["oemconfig"]["clientlogouttime"] = Int::getValue(kvMap["oem_client_logout_time"]);
    replyJson["oemconfig"]["maxpassexpireddays"] = Int::getValue(kvMap["oem_max_pass_expired_days"]);
    // 内链地址的前缀
    replyJson["internal_link_prefix"] = kvMap["internal_link_prefix"].getCStr();

    // sharemgnt t_conf表中批量获取配置
    vector<String> sharemgntKeys;
    map<String, String> sharemgntKvMap;
    sharemgntKeys.push_back("hide_client_cache_setting");
    sharemgntKeys.push_back("force_clear_client_cache");
    sharemgntKeys.push_back("only_share_to_user");
    _acsShareMgnt->BatchGetConfig(sharemgntKeys, sharemgntKvMap);

    // 开启客户端https连接
    replyJson["oemconfig"]["clearcache"] = sharemgntKvMap["force_clear_client_cache"].compareIgnoreCase("1") == 0;
    replyJson["oemconfig"]["hidecachesetting"] = sharemgntKvMap["hide_client_cache_setting"].compareIgnoreCase("1") == 0;
    // 是否只允许共享给用户
    replyJson["only_share_to_user"] = sharemgntKvMap["only_share_to_user"].compareIgnoreCase("1") == 0;
    replyJson["csf_level_enum"] = getCsfLevelsConfig();

    // 获取全部配置
    ncTAllConfig allConfig;
    ncEACHttpServerUtil::GetAllConfig(allConfig);

    // 获取第三方标密系统配置
    ncTThirdCSFSysConfig &thirdCSFSysConfig = allConfig.thirdCSFSysConfig;
    if (thirdCSFSysConfig.isEnabled == true)
    {
        JSON::Value _thirdCSFSysConfig;
        _thirdCSFSysConfig["id"] = thirdCSFSysConfig.id.c_str();
        _thirdCSFSysConfig["only_upload_classified"] = thirdCSFSysConfig.only_upload_classified;
        _thirdCSFSysConfig["only_share_classified"] = thirdCSFSysConfig.only_share_classified;
        _thirdCSFSysConfig["auto_match_doc_classfication"] = thirdCSFSysConfig.auto_match_doc_classfication;
        replyJson["third_csfsys_config"] = _thirdCSFSysConfig;
    }

    // 获取文件最大标签数配置
    replyJson["tag_max_num"] = allConfig.tag_max_num;

    // 是否配置SMTP服务器
    ncTSmtpSrvConf &smtpSrvConfig = allConfig.smtpSrvConfig;
    replyJson["smtp_server_exists"] = !smtpSrvConfig.server.empty();

    string body;
    JSON::Writer::write(replyJson.o(), body);
    ncHttpSendReply(cntl, brpc::HTTP_STATUS_OK, "ok", body);

    NC_EAC_HTTP_SERVER_CATCH (cntl)
}

void ncEACAuthHandler::GetConfig (brpc::Controller* cntl, const String& fakeUserId)
{
    NC_EAC_HTTP_SERVER_TRACE (_T("this: %p, cntl:%p, userId:%s begin"), this, cntl, fakeUserId.getCStr ());

    JSON::Value replyJson;
    // 从anyshare t_conf表中批量获取配置
    vector<String> keys;
    keys.push_back("oem_indefinite_perm");
    keys.push_back("oem_allow_owner");
    keys.push_back("oem_remember_pass");
    keys.push_back("oem_max_pass_expired_days");
    keys.push_back("oem_allow_auth_low_csf_user");
    keys.push_back("oem_client_logout_time");
    keys.push_back("oem_enable_file_transfer_limit");
    keys.push_back("oem_enable_onedrive");
    keys.push_back("client_manual_login");
    keys.push_back("entrydoc_view_config");
    keys.push_back("enable_chaojibiaoge");
    keys.push_back("enable_qhdj");
    keys.push_back("auto_lock_remind");
    keys.push_back("internal_link_prefix");
    keys.push_back("show_knowledge_page");
    keys.push_back("enable_message_notify");
    keys.push_back("cad_plugin_threshold");
    keys.push_back("custome_application_config");

    map<String, String> kvMap;
    _acsConfManager->BatchGetConfig(keys, kvMap);

    // 返回oemconfig
    replyJson["oemconfig"]["indefiniteperm"] = kvMap["oem_indefinite_perm"].compareIgnoreCase("true") == 0;
    replyJson["oemconfig"]["allowowner"] = kvMap["oem_allow_owner"].compareIgnoreCase("true") == 0;
    replyJson["oemconfig"]["rememberpass"] = kvMap["oem_remember_pass"].compareIgnoreCase("true") == 0;
    replyJson["oemconfig"]["maxpassexpireddays"] = Int::getValue(kvMap["oem_max_pass_expired_days"]);
    replyJson["oemconfig"]["allowauthlowcsfuser"] = kvMap["oem_allow_auth_low_csf_user"].compareIgnoreCase("true") == 0;
    replyJson["oemconfig"]["clientlogouttime"] = Int::getValue(kvMap["oem_client_logout_time"]);
    replyJson["oemconfig"]["enablefiletransferlimit"] = kvMap["oem_enable_file_transfer_limit"].compareIgnoreCase("true") == 0;
    replyJson["oemconfig"]["enableonedrive"] = kvMap["oem_enable_onedrive"].compareIgnoreCase("true") == 0;
    replyJson["oemconfig"]["enableclientmanuallogin"] = kvMap["client_manual_login"].compareIgnoreCase("true") == 0;
    replyJson["oemconfig"]["defaultpermexpireddays"] = getDefaultPermExpiredDays();
    // 入口文档视图配置
    replyJson["entrydoc_view_config"] = Int::getValue(kvMap["entrydoc_view_config"]);
    // 开启外部应用配置
    replyJson["extapp"]["enable_chaojibiaoge"] = kvMap["enable_chaojibiaoge"].compareIgnoreCase("true") == 0;
    // 开启秦淮电教馆功能定制按钮
    replyJson["extapp"]["enable_qhdj"] = kvMap["enable_qhdj"].compareIgnoreCase("true") == 0;
    // 增加文件锁提醒配置
    replyJson["auto_lock_remind"] = kvMap["auto_lock_remind"].compareIgnoreCase("true") == 0;
    // 内链地址的前缀
    replyJson["internal_link_prefix"] = kvMap["internal_link_prefix"].getCStr ();
    // 显示知识主页开关
    replyJson["show_knowledge_page"] = Int::getValue(kvMap["show_knowledge_page"]);
    // 是否启用消息通知
    replyJson["enable_message_notify"] = kvMap["enable_message_notify"].compareIgnoreCase("true") == 0;
    // 浩辰CAD使用大图插件的临界值
    replyJson["cad_plugin_threshold"] = Int::getValue(kvMap["cad_plugin_threshold"]);
    // 唯一定制化的应用配置
    JSON::Value appConfigJson;
    String appConfig = kvMap["custome_application_config"].getCStr ();
    if (!appConfig.isEmpty ()) {
        try {
            JSON::Reader::read (appConfigJson, appConfig.getCStr (), appConfig.getLength ());
        }
        catch (Exception& e) {
            THROW_E (EAC_HTTP_SERVER, EACHTTP_INVALID_REQUEST_JSON_FORMAT,
                LOAD_STRING (_T("IDS_EACHTTP_INVALID_JSON")));
        }
        replyJson["custome_application_config"] = appConfigJson;
    }

    // sharemgnt t_conf表中批量获取配置
    vector<String> sharemgntKeys;
    map<String, String> sharemgntKvMap;
    sharemgntKeys.push_back("force_clear_client_cache");
    sharemgntKeys.push_back("hide_client_cache_setting");
    sharemgntKeys.push_back("client_https");
    sharemgntKeys.push_back("windows_ad_sso");
    sharemgntKeys.push_back("enable_secret_mode");
    sharemgntKeys.push_back("invitation_share_status");
    sharemgntKeys.push_back("third_pwd_modify_url");
    sharemgntKeys.push_back("vcode_login_config");
    sharemgntKeys.push_back("limit_rate_status");
    sharemgntKeys.push_back("strong_pwd_status");
    sharemgntKeys.push_back("limit_rate_config");
    sharemgntKeys.push_back("strong_pwd_length");
    sharemgntKeys.push_back("file_crawl_status");
    sharemgntKeys.push_back("id_card_login_status");
    sharemgntKeys.push_back("enable_set_folder_security_level");
    sharemgntKeys.push_back("only_share_to_user");
    sharemgntKeys.push_back("vcode_server_status");
    sharemgntKeys.push_back("enable_exit_pwd");
    sharemgntKeys.push_back("enable_outlink_watermark");
    sharemgntKeys.push_back("dualfactor_auth_server_status");
    _acsShareMgnt->BatchGetConfig(sharemgntKeys, sharemgntKvMap);

    // 开启客户端https连接
    replyJson["oemconfig"]["clearcache"] = sharemgntKvMap["force_clear_client_cache"].compareIgnoreCase("1") == 0;
    replyJson["oemconfig"]["hidecachesetting"] = sharemgntKvMap["hide_client_cache_setting"].compareIgnoreCase("1") == 0;
    replyJson["oemconfig"]["enableshareaudit"] = false;
    replyJson["oemconfig"]["enablecsflevel"] = false;
    replyJson["oemconfig"]["enablehttplinkaudit"] = false;
    replyJson["https"] = sharemgntKvMap["client_https"].compareIgnoreCase("1") == 0;
    // windows ad单点登录
    replyJson["windows_ad_sso"]["is_enabled"] = sharemgntKvMap["windows_ad_sso"].compareIgnoreCase("1") == 0;
    // 开启涉密模式
    replyJson["enable_secret_mode"] = sharemgntKvMap["enable_secret_mode"].compareIgnoreCase("1") == 0;
    // 开启共享邀请
    replyJson["enable_invitation_share"] = sharemgntKvMap["invitation_share_status"].compareIgnoreCase("1") == 0;
    // 获取服务端密级配置
    replyJson["csf_level_enum"] = getCsfLevelsConfig();
    // 获取禁用客户端类型配置,此值没有地方用到，默认为0
    replyJson["forbid_ostype"] = "0";
    // 第三方用户密码修改地址
    replyJson["third_pwd_modify_url"] = sharemgntKvMap["third_pwd_modify_url"].getCStr();
    // 开启限速配置
    replyJson["enable_limit_rate"] = sharemgntKvMap["limit_rate_status"].compareIgnoreCase("1") == 0;
    // 开启强密码配置
    replyJson["enable_strong_pwd"] = sharemgntKvMap["strong_pwd_status"].compareIgnoreCase("1") == 0;
    // 是否允许用户设置外链水印
    replyJson["enable_outlink_watermark"] = sharemgntKvMap["enable_outlink_watermark"].compareIgnoreCase("1") == 0;

    // 获取全部配置
    ncTAllConfig allConfig ;
    ncEACHttpServerUtil::GetAllConfig (allConfig);

    // 获取第三方标密系统配置
    ncTThirdCSFSysConfig &thirdCSFSysConfig = allConfig.thirdCSFSysConfig;
    if (thirdCSFSysConfig.isEnabled == true)
    {
        JSON::Value _thirdCSFSysConfig;
        _thirdCSFSysConfig["id"] = thirdCSFSysConfig.id.c_str();
        _thirdCSFSysConfig["only_upload_classified"] = thirdCSFSysConfig.only_upload_classified;
        _thirdCSFSysConfig["only_share_classified"] = thirdCSFSysConfig.only_share_classified;
        _thirdCSFSysConfig["auto_match_doc_classfication"] = thirdCSFSysConfig.auto_match_doc_classfication;
        replyJson["third_csfsys_config"] = _thirdCSFSysConfig;
    }
    // 第三方认证配置
    ncTThirdPartyAuthConf &config = allConfig.thirdPartyAuthConf;
    if(config.enabled) {
        JSON::Value configJson;
        JSON::Reader::read (configJson, config.config.c_str(), config.config.length());
        replyJson["thirdauth"]["id"] = config.thirdPartyId;

    for(auto iter = configJson.o().begin(); iter != configJson.o().end(); ++iter) {
        //过滤黑名单中信息
        if(_thirdConfigBlackList.find(iter->first) != _thirdConfigBlackList.end())
            continue;
        replyJson["thirdauth"]["config"][iter->first] = iter->second;
        }
    }
    // owas的url地址
    ncTThirdPartyToolConfig &toolConfig = allConfig.toolOfficeConfig;
    replyJson["oemconfig"]["owasurl"] = "";
    ncTThirdPartyToolConfig &toolWOPIConfig = allConfig.toolWOPIConfig;
    replyJson["oemconfig"]["wopiurl"] = "";
    if(toolConfig.enabled && toolWOPIConfig.enabled) {
        replyJson["oemconfig"]["owasurl"] = toolConfig.url;
        replyJson["oemconfig"]["wopiurl"] = toolWOPIConfig.url;
    }
    // CAD预览开关
    ncTThirdPartyToolConfig &toolCADConfig = allConfig.toolCADConfig;
    replyJson["oemconfig"]["cadpreview"] = false;
    if(toolCADConfig.enabled) {
        replyJson["oemconfig"]["cadpreview"] = true;
        string toolName = toolCADConfig.thirdPartyToolName;
        replyJson["oemconfig"]["cadtool"] = toolName;
        if (toolName.compare("hc") == 0) {  // 使用浩辰CAD时，返回转码服务器url
            replyJson["oemconfig"]["cadurl"] = toolCADConfig.url;
        }
    }
    //gd/sep配置开关
    ncTThirdPartyToolConfig &toolSurSenConfig = allConfig.toolSurSenConfig;
    replyJson["oemconfig"]["sursenpreview"] = false;
    if(toolSurSenConfig.enabled) {
        replyJson["oemconfig"]["sursenpreview"] = true;
    }
    replyJson["oemconfig"]["enableuseragreement"] = allConfig.enableuseragreement;

    // 获取文件最大标签数配置
    replyJson["tag_max_num"] = allConfig.tag_max_num;
    // 开启文件提取码
    // replyJson["enable_link_access_code"] = ncEACHttpServerUtil::GetLinkAccessCodeStatus();
    // 该功能已弃用，未保持接口兼容性暂时返回false
    replyJson["enable_link_access_code"] = false;
    // 开启文件到期提醒
    // replyJson["enable_doc_due_remind"] = ncEACHttpServerUtil::GetDocDueRemindStatus();
    // 该功能已弃用，未保持接口兼容性暂时返回false
    replyJson["enable_doc_due_remind"] = false;

    // 获取服务端版本、API版本 都改为固定值, 以便兼容
    replyJson["server_version"] = "7.0.1.5-20210430-el7.x86_64-1734";
    replyJson["api_version"] = "7.0.0";

    // 获取验证码配置
    string vcode_config = sharemgntKvMap["vcode_login_config"].getCStr();
    JSON::Value vcodeJson;
    try {
        JSON::Reader::read (vcodeJson, vcode_config.c_str(), vcode_config.length ());
    }
    catch (Exception& e) {
        THROW_E (EAC_HTTP_SERVER, EACHTTP_INVALID_REQUEST_JSON_FORMAT,
            LOAD_STRING (_T("IDS_EACHTTP_INVALID_JSON")));
    }
    replyJson["vcode_login_config"]["isenable"] = vcodeJson["isEnable"];
    replyJson["vcode_login_config"]["passwderrcnt"] = vcodeJson["passwdErrCnt"];

    // 获取发送验证码服务器开关配置
    string vcode_server_status = sharemgntKvMap["vcode_server_status"].getCStr();
    JSON::Value sendVcodeTypeStatusJson;
    try {
        JSON::Reader::read (sendVcodeTypeStatusJson, vcode_server_status.c_str(), vcode_server_status.length ());
    }
    catch (Exception& e) {
        THROW_E (EAC_HTTP_SERVER, EACHTTP_INVALID_REQUEST_JSON_FORMAT,
            LOAD_STRING (_T("IDS_EACHTTP_INVALID_JSON")));
    }
    replyJson["vcode_server_status"]["send_vcode_by_sms"] = sendVcodeTypeStatusJson["send_vcode_by_sms"];
    replyJson["vcode_server_status"]["send_vcode_by_email"] = sendVcodeTypeStatusJson["send_vcode_by_email"];

    // 获取双因子验证开关配置
    string mfa_server_status = sharemgntKvMap["dualfactor_auth_server_status"].getCStr();
    JSON::Value dualauthJson;
    try {
        JSON::Reader::read(dualauthJson, mfa_server_status.c_str(), mfa_server_status.length ());
    }
    catch (Exception& e) {
        THROW_E (EAC_HTTP_SERVER, EACHTTP_INVALID_REQUEST_JSON_FORMAT,
            LOAD_STRING (_T("IDS_EACHTTP_INVALID_JSON")));
    }
    replyJson["dualfactor_auth_server_status"]["auth_by_sms"] = dualauthJson["auth_by_sms"];
    replyJson["dualfactor_auth_server_status"]["auth_by_email"] = dualauthJson["auth_by_email"];
    replyJson["dualfactor_auth_server_status"]["auth_by_OTP"] = dualauthJson["auth_by_OTP"];
    replyJson["dualfactor_auth_server_status"]["auth_by_Ukey"] = dualauthJson["auth_by_Ukey"];

    // 获取限速配置
    string limit_config = sharemgntKvMap["limit_rate_config"].getCStr();
    JSON::Value limitRateJson;
    try {
        JSON::Reader::read (limitRateJson, limit_config.c_str(), limit_config.length ());
    }
    catch (Exception& e) {
        THROW_E (EAC_HTTP_SERVER, EACHTTP_INVALID_REQUEST_JSON_FORMAT,
            LOAD_STRING (_T("IDS_EACHTTP_INVALID_JSON")));
    }
    replyJson["limit_rate_config"]["isenabled"] = limitRateJson["isEnabled"];
    replyJson["limit_rate_config"]["limittype"] = limitRateJson["limitType"];

    // 获取强密码最小长度
    replyJson["strong_pwd_length"] = Int::getValue(sharemgntKvMap["strong_pwd_length"]);

    // 文件抓取策略
    replyJson["file_crawl_status"] = sharemgntKvMap["file_crawl_status"].compareIgnoreCase("1") == 0;

    // 是否允许设置文件夹密级
    replyJson["enable_set_folder_security_level"] = sharemgntKvMap["enable_set_folder_security_level"].compareIgnoreCase("1") == 0;

    // 是否只允许共享给用户
    replyJson["only_share_to_user"] = sharemgntKvMap["only_share_to_user"].compareIgnoreCase("1") == 0;

    // 身份证登陆策略
    replyJson["id_card_login_status"] = sharemgntKvMap["id_card_login_status"].compareIgnoreCase("1") == 0;

    // 退出口令开关
    replyJson["enable_exit_pwd"] = sharemgntKvMap["enable_exit_pwd"].compareIgnoreCase("1") == 0;

    //是否配置SMTP服务器
    ncTSmtpSrvConf &smtpSrvConfig = allConfig.smtpSrvConfig;
    replyJson["smtp_server_exists"] = !smtpSrvConfig.server.empty ();

    string body;
    JSON::Writer::write (replyJson.o (), body);
    ncHttpSendReply (cntl, brpc::HTTP_STATUS_OK, "ok", body);

    NC_EAC_HTTP_SERVER_TRACE (_T("this: %p, cntl:%p, userId:%s end"), this, cntl, fakeUserId.getCStr ());
}

// protected
void
ncEACAuthHandler::GetByNTLMV1 (brpc::Controller* cntl, const String& fakeUserId)
{
    NC_EAC_HTTP_SERVER_TRACE (_T("this: %p, cntl:%p, userId:%s begin"), this, cntl, fakeUserId.getCStr ());

    // 获取http请求的content
    string bodyBuffer = cntl->request_attachment ().to_string ();

    // 获取account,password,clientversion
    JSON::Value requestJson;
    try {
        JSON::Reader::read (requestJson, bodyBuffer.c_str (), bodyBuffer.size ());
    }
    catch (Exception& e) {
        THROW_E (EAC_HTTP_SERVER, EACHTTP_INVALID_REQUEST_JSON_FORMAT,
            LOAD_STRING (_T("IDS_EACHTTP_INVALID_JSON")));
    }

    string account = requestJson["account"].s ();
    if (account.empty ()) {
        THROW_E (EAC_HTTP_SERVER, EACHTTP_ACCOUNT_INVALID,
            LOAD_STRING (_T("IDS_EACHTTP_ACCOUNT_INVALID")));
    }

    string challenge = requestJson["challenge"].s ();
    if (challenge.empty ()) {
        THROW_E (EAC_HTTP_SERVER, EACHTTP_PASSWORD_INVALID,
                LOAD_STRING (_T("IDS_EACHTTP_PASSWORD_INVALID")));
    }

    string password = requestJson["password"].s ();
    if (challenge.empty ()) {
        THROW_E (EAC_HTTP_SERVER, EACHTTP_PASSWORD_INVALID,
                LOAD_STRING (_T("IDS_EACHTTP_PASSWORD_INVALID")));
    }

    string domain = requestJson["domain"].s ();

    // 进行验证
    ncTNTLMResponse ntlmResp;
    if (password.length () > 48) {
        ncEACHttpServerUtil::Usrm_UserLoginByNTLMV2 (ntlmResp, account, domain, challenge, password);
    } else {
        ncEACHttpServerUtil::Usrm_UserLoginByNTLMV1 (ntlmResp, account, challenge, password);
    }

    bool ret = _acsPolicyManager->CheckIp (toCFLString (ntlmResp.userId), ncEACHttpServerUtil::GetForwardedIp(cntl));
    if(!ret) {
        THROW_E (EAC_HTTP_SERVER, EACHTTP_LOGIN_IP_IS_RESTRICTED,
            LOAD_STRING("IDS_LOGIN_IP_IS_RESTRICTED"));
    }

#ifndef __UT__
    // 更新用户最近请求时间记录
    _acsShareMgnt->UpdateUserLastRequestTime(toCFLString(ntlmResp.userId));
#endif

    // 回复
    JSON::Value replyJson;
    replyJson["user_id"] = ntlmResp.userId;
    replyJson["context"]["visitor_type"] =  _visitorStringTypeMap[1].getCStr();
    replyJson["context"]["udid"] = "";
    replyJson["context"]["login_ip"] = ncEACHttpServerUtil::GetForwardedIp(cntl).getCStr();
    replyJson["context"]["client_type"] = _clientStringTypeMap[ACSClientType::NAS].getCStr();
    replyJson["context"]["account_type"] = _accountStringTypeMap[0].getCStr();
    replyJson["sesskey"] = ntlmResp.sessKey;

    string body;
    JSON::Writer::write (replyJson.o (), body);
    ncHttpSendReply (cntl, brpc::HTTP_STATUS_OK, "ok", body);

    String msg;
    String nasRes;
    nasRes.format (ncEACHttpServerLoader, _T("IDS_NAS_GATEWAY"));   //获取到 “NAS网关” 国际化资源
    msg.format (ncEACHttpServerLoader, _T("IDS_LOGIN_SUCCESS"), nasRes.getCStr ());
    ncEACHttpServerUtil::Log (cntl, toCFLString(ntlmResp.userId), ncTokenVisitorType::REALNAME, ncTLogType::NCT_LT_LOGIN, ncTLogLevel::NCT_LL_INFO,
        ncTLoginType::NCT_CLT_LOGIN_IN, msg, "", true /* logForwardedIp */);

    ncEACHttpServerUtil::LoginLog(toCFLString(ntlmResp.userId), "", ACSClientType::NAS, ncEACHttpServerUtil::GetForwardedIp(cntl));

    NC_EAC_HTTP_SERVER_TRACE (_T("this: %p, cntl:%p, userId:%s end"), this, cntl, fakeUserId.getCStr ());
}

// protected
void
ncEACAuthHandler::CheckLoginParams (brpc::Controller* cntl, string &account, string &originPassword, bool &hasDeviceInfo, ncDeviceBaseInfo &baseInfo, JSON::Value &requestJson, ncTUserLoginOption &option, bool b2048)
{
    // 获取http请求的content
    string bodyBuffer = cntl->request_attachment ().to_string ();

    // 获取account,password
    try {
        JSON::Reader::read (requestJson, bodyBuffer.c_str (), bodyBuffer.size ());
    }
    catch (Exception& e) {
        THROW_E (EAC_HTTP_SERVER, EACHTTP_INVALID_REQUEST_JSON_FORMAT,
            LOAD_STRING (_T("IDS_EACHTTP_INVALID_JSON")));
    }

    account = requestJson["account"].s ();
    if (account.empty ()) {
        THROW_E (EAC_HTTP_SERVER, EACHTTP_ACCOUNT_INVALID,
            LOAD_STRING (_T("IDS_EACHTTP_ACCOUNT_INVALID")));
    }

    string password = requestJson["password"].s ();
    if (password.empty ()) {
        THROW_E (EAC_HTTP_SERVER, EACHTTP_PASSWORD_INVALID,
            LOAD_STRING (_T("IDS_EACHTTP_PASSWORD_INVALID")));
    }

    // 获取登录的设备信息
    hasDeviceInfo = parseDeviceInfo(requestJson, baseInfo);


    // 这里password有可能尾部有\0
    string decodePwd(ncEACHttpServerUtil::Base64Decode(password.c_str()));
    if (b2048) {
        originPassword = ncEACHttpServerUtil::RSADecrypt2048(decodePwd);
    } else {
        originPassword = ncEACHttpServerUtil::RSADecrypt(decodePwd);
    }

    // 获取验证码信息
    option.vcode = requestJson["vcode"]["content"].s ();
    option.__isset.vcode = true;
    option.uuid = requestJson["vcode"]["id"].s ();
    option.__isset.uuid = true;
    option.isModify = false;
    option.__isset.isModify = false;
}

// protected
void
ncEACAuthHandler::CheckLoginIpParam(JSON::Value &requestJson, String& loginIp)
{
    //ip信息获取
    if(requestJson.o().find("ip") == requestJson.o().end())
    {
        THROW_E (EAC_HTTP_SERVER, EACHTTP_INVALID_PARAMETER,
            LOAD_STRING (_T("IDS_EACHTTP_PARAMETERS_IS_REQUIRED")), "ip");
    }
    loginIp =  toCFLString(requestJson["ip"].s ());
    if(loginIp.isEmpty ())
    {
        THROW_E (EAC_HTTP_SERVER, EACHTTP_INVALID_PARAMETER,
            LOAD_STRING (_T("IDS_EACHTTP_PARAMETERS_IS_INVALID")), "ip");
    }
}

// protected
void
ncEACAuthHandler::CheckConsoleLoginParams (brpc::Controller* cntl, bool &isAccountType, ncAccountCredential &accountCredential,
                                  ncThirdPartyCredential &thirdPartyCredential, ncDeviceBaseInfo &baseInfo, String& loginIp, JSON::Value &requestJson)
{
// 获取http请求的content
    string bodyBuffer = cntl->request_attachment ().to_string ();

    // 获取account,password
    try {
        JSON::Reader::read (requestJson, bodyBuffer.c_str (), bodyBuffer.size ());
    }
    catch (Exception& e) {
        THROW_E (EAC_HTTP_SERVER, EACHTTP_INVALID_REQUEST_JSON_FORMAT,
            LOAD_STRING (_T("IDS_EACHTTP_INVALID_JSON")));
    }

    //检查登录IP
    CheckLoginIpParam(requestJson, loginIp);
    //更新报文头的ip
    updateHTTPHeaderIP(cntl, loginIp);

    //credential 参数存在性判断
    if(requestJson.o().find("credential") == requestJson.o().end())
    {
        THROW_E (EAC_HTTP_SERVER, EACHTTP_INVALID_PARAMETER,
            LOAD_STRING (_T("IDS_EACHTTP_PARAMETERS_IS_REQUIRED")), "credential");
    }

    //credential type参数有效性和存在性判断以及参数赋值
    JSON::Value credentialJson = requestJson["credential"];
    if(credentialJson.o().find("type") == credentialJson.o().end())
    {
        THROW_E (EAC_HTTP_SERVER, EACHTTP_INVALID_PARAMETER,
            LOAD_STRING (_T("IDS_EACHTTP_PARAMETERS_IS_REQUIRED")), "credential type");
    }
    string type = credentialJson["type"].s ();

    if( type == "account" )
    {
        isAccountType = true;

        //account参数有效性和存在性判断
        if(credentialJson.o().find("account") == credentialJson.o().end())
        {
            THROW_E (EAC_HTTP_SERVER, EACHTTP_INVALID_PARAMETER,
                LOAD_STRING (_T("IDS_EACHTTP_PARAMETERS_IS_REQUIRED")), "account");
        }
        accountCredential.account = toCFLString(credentialJson["account"].s ());
        if (accountCredential.account.isEmpty ()){
            THROW_E (EAC_HTTP_SERVER, EACHTTP_INVALID_PARAMETER,
                LOAD_STRING (_T("IDS_EACHTTP_PARAMETERS_IS_INVALID")), "account");
        }

        //控制台登录不在此位置进行rsa解密 而是在thrift服务内进行解密
        //password参数有效性和存在性判断
        if(credentialJson.o().find("password") == credentialJson.o().end())
        {
            THROW_E (EAC_HTTP_SERVER, EACHTTP_INVALID_PARAMETER,
                LOAD_STRING (_T("IDS_EACHTTP_PARAMETERS_IS_REQUIRED")), "password");
        }
        accountCredential.originPassword = toCFLString(credentialJson["password"].s ());
        if (accountCredential.originPassword.isEmpty ()){
            THROW_E (EAC_HTTP_SERVER, EACHTTP_INVALID_PARAMETER,
                LOAD_STRING (_T("IDS_EACHTTP_PARAMETERS_IS_INVALID")), "password");
        }

        parseVcodeInfo(credentialJson, accountCredential.option);

        //设置登录ip信息
        accountCredential.option.loginIp = loginIp.getCStr();
        accountCredential.option.__isset.loginIp = true;
    }
    else if(type == "third_party")
    {
        isAccountType = false;
        if(credentialJson.o().find("params") == credentialJson.o().end())
        {
            THROW_E (EAC_HTTP_SERVER, EACHTTP_INVALID_PARAMETER,
                LOAD_STRING (_T("IDS_EACHTTP_PARAMETERS_IS_REQUIRED")), "credential params");
        }
        if(credentialJson["params"].type() != JSON::OBJECT)
        {
            THROW_E (EAC_HTTP_SERVER, EACHTTP_INVALID_PARAMETER,
                LOAD_STRING (_T("IDS_EACHTTP_PARAMETERS_IS_INVALID")), "params");
        }
        string params;
        JSON::Value paramsJSON;
        paramsJSON["params"] = credentialJson["params"];
        JSON::Writer::write (paramsJSON.o (), params);
        thirdPartyCredential.params = toCFLString(params);
    }
    else
    {
        THROW_E (EAC_HTTP_SERVER, EACHTTP_INVALID_PARAMETER,
            LOAD_STRING (_T("IDS_EACHTTP_PARAMETERS_IS_INVALID")), "credential type");
    }

    // 获取登录的设备信息
    parseDeviceInfo(requestJson, baseInfo);
}

// public
void
ncEACAuthHandler::ConsoleLogin (brpc::Controller* cntl)
{
    NC_EAC_HTTP_SERVER_TRY
    cntl->http_response().set_content_type("application/json; charset=UTF-8");
    // 判断请求方法是否为POST
    if (cntl->http_request ().method () != brpc::HTTP_METHOD_POST){
        THROW_E (EACHTTP_SERVER_ERR_PROVIDER_NAME, EACHTTP_SERVER_METHOD_NOT_IMPLEMENTED,
                 LOAD_STRING (_T("IDS_EACHTTP_METHOD_INVALID")));
    }
    NC_EAC_HTTP_SERVER_TRACE (_T("this: %p, cntl:%p begin"), this, cntl);

    // 参数获取与检查
    bool isAccountType;
    ncDeviceBaseInfo baseInfo;
    JSON::Value requestJson;
    ncAccountCredential accountCredential;
    ncThirdPartyCredential thirdPartyCredential;
    String loginIp;
    CheckConsoleLoginParams(cntl, isAccountType, accountCredential, thirdPartyCredential, baseInfo, loginIp, requestJson);

    string retUserId;
    string strOSType = clientType2Str(baseInfo.clientType).getCStr();
    int accountType;
    if(isAccountType)
    {
        //账户密码登录 验证
        ncTUsrmAuthenType::type nAuthenType = ncTUsrmAuthenType::NCT_AUTHEN_TYPE_MANAGER; //默认管理员登录
        try {
            ncEACHttpServerUtil::Usrm_Login (retUserId, accountCredential.account.getCStr (),
                                        accountCredential.originPassword.getCStr (),
                                        nAuthenType, accountCredential.option, strOSType);
        }
        catch(EHttpDetailException& e) {
            logFailedConsoleLoginEvent(cntl, strOSType, e);
            throw;
        }
        catch (Exception& e) {
            throw;
        }
        accountType = _acsShareMgnt->GetAccountType(accountCredential.account);
    }
    else
    {
        //第三方凭证验证
        String strAccount = ncEACHttpServerUtil::Usrm_LoginConsoleByThirdPartyNew(thirdPartyCredential.params);

        ncACSUserInfo userInfo;
        if(_acsShareMgnt->GetUserInfoByAccount(strAccount, userInfo, accountType) == false) {
            THROW_E (EAC_HTTP_SERVER, EACHTTP_NOT_IMPORT_TO_ANYSHARE,
                LOAD_STRING (_T("IDS_EACHTTP_NOT_IMPORT_TO_ANYSHARE")));
        }
        retUserId = userInfo.id.getCStr();
    }

    bool ret = _acsPolicyManager->CheckIp (toCFLString(retUserId), loginIp);
    if(!ret) {
        THROW_E (EAC_HTTP_SERVER, EACHTTP_LOGIN_IP_IS_RESTRICTED,
            LOAD_STRING("IDS_LOGIN_IP_IS_RESTRICTED"));
    }

    // 更新用户激活状态
    _acsShareMgnt->UpdateUserActivateStatus(toCFLString(retUserId));

    // 回复
    JSON::Value replyJson;
    replyJson["user_id"] = retUserId;
    replyJson["context"]["visitor_type"] = _visitorStringTypeMap[1].getCStr();
    replyJson["context"]["udid"] = baseInfo.udid.getCStr ();
    replyJson["context"]["login_ip"] = loginIp.getCStr();
    replyJson["context"]["client_type"] = _clientStringTypeMap[baseInfo.clientType].getCStr();
    replyJson["context"]["account_type"] = _accountStringTypeMap[accountType].getCStr();

    string body;
    JSON::Writer::write (replyJson.o (), body);
    ncHttpSendReply (cntl, brpc::HTTP_STATUS_OK, "ok", body);

    //登录成功 记录日志
    ncEACHttpServerUtil::Log (cntl, toCFLString(retUserId), ncTokenVisitorType::REALNAME, ncTLogType::NCT_LT_LOGIN, ncTLogLevel::NCT_LL_INFO,
        ncTLoginType::NCT_CLT_AUTHENICATION, LOAD_STRING (_T("IDS_AUTHENICATION_SUCCESS")), "");

    NC_EAC_HTTP_SERVER_TRACE (_T("this: %p, cntl:%p, userId:%s end"), this, cntl, retUserId.c_str());
    NC_EAC_HTTP_SERVER_CATCH (cntl)
}

// public
void
ncEACAuthHandler::GetNew (brpc::Controller* cntl)
{
    NC_EAC_HTTP_SERVER_TRY
    cntl->http_response().set_content_type("application/json; charset=UTF-8");
    // 判断请求方法是否为POST
    if (cntl->http_request ().method () != brpc::HTTP_METHOD_POST){
        THROW_E (EACHTTP_SERVER_ERR_PROVIDER_NAME, EACHTTP_SERVER_METHOD_NOT_IMPLEMENTED,
                 LOAD_STRING (_T("IDS_EACHTTP_METHOD_INVALID")));
    }
    NC_EAC_HTTP_SERVER_TRACE (_T("this: %p, cntl:%p begin"), this, cntl);

    // 参数检查
    string account;
    string originPassword;
    bool hasDeviceInfo;
    ncDeviceBaseInfo baseInfo;
    JSON::Value requestJson;
    ncTUserLoginOption option;
    String loginIp;
    CheckLoginParams(cntl, account, originPassword, hasDeviceInfo, baseInfo, requestJson, option, true);
    CheckLoginIpParam(requestJson, loginIp);
    if (!option.uuid.empty()) {
        option.vcodeType = ncTVcodeType::IMAGE_VCODE;
    }

    //更新报文头的ip
    updateHTTPHeaderIP(cntl, loginIp);

    // 检查是否是外网ip
    bool isOutNet = checkIpIsOutNet(loginIp.getCStr());

    // 双因子认证
    dualFactorAuth(option, requestJson, isOutNet);

    // 进行验证
    string retUserId;
    try {
        ncEACHttpServerUtil::ClientLogin(retUserId, account, originPassword, option);
    }
    catch(Exception& e) {
        authenicaitonFailedLoginEvent(cntl, toCFLString(account), e, hasDeviceInfo, baseInfo);
        throw;
    }

    onUserLogin(cntl, toCFLString(retUserId), loginIp, hasDeviceInfo, baseInfo);

    int accountType = _acsShareMgnt->GetAccountType(toCFLString(account));

    // 更新用户激活状态
    _acsShareMgnt->UpdateUserLastRequestTime(toCFLString(retUserId));

    // 回复
    JSON::Value replyJson;
    replyJson["user_id"] = retUserId;
    replyJson["context"]["visitor_type"] = _visitorStringTypeMap[1].getCStr();
    replyJson["context"]["udid"] = baseInfo.udid.getCStr ();
    replyJson["context"]["login_ip"] = loginIp.getCStr();
    replyJson["context"]["client_type"] = _clientStringTypeMap[baseInfo.clientType].getCStr();
    replyJson["context"]["account_type"] = _accountStringTypeMap[accountType].getCStr();

    string body;
    JSON::Writer::write (replyJson.o (), body);
    ncHttpSendReply (cntl, brpc::HTTP_STATUS_OK, "ok", body);

    ncEACHttpServerUtil::Log (cntl, toCFLString(retUserId), ncTokenVisitorType::REALNAME, ncTLogType::NCT_LT_LOGIN, ncTLogLevel::NCT_LL_INFO,
        ncTLoginType::NCT_CLT_AUTHENICATION, LOAD_STRING (_T("IDS_AUTHENICATION_SUCCESS")), "");

    NC_EAC_HTTP_SERVER_TRACE (_T("this: %p, cntl:%p, userId:%s end"), this, cntl, retUserId.c_str());
    NC_EAC_HTTP_SERVER_CATCH (cntl)
}

void ncEACAuthHandler::GetByThirdParty (brpc::Controller* cntl)
{
    NC_EAC_HTTP_SERVER_TRY
    cntl->http_response().set_content_type("application/json; charset=UTF-8");
    // 判断请求方法是否为POST
    if (cntl->http_request ().method () != brpc::HTTP_METHOD_POST){
        THROW_E (EACHTTP_SERVER_ERR_PROVIDER_NAME, EACHTTP_SERVER_METHOD_NOT_IMPLEMENTED,
                 LOAD_STRING (_T("IDS_EACHTTP_METHOD_INVALID")));
    }
    NC_EAC_HTTP_SERVER_TRACE (_T("this: %p, cntl:%p begin"), this, cntl);

    // 获取http请求的content
    string bodyBuffer = cntl->request_attachment ().to_string ();

    // 获取ticket
    JSON::Value requestJson;
    try {
        JSON::Reader::read (requestJson, bodyBuffer.c_str (), bodyBuffer.size ());
    }
    catch (Exception& e) {
        THROW_E (EAC_HTTP_SERVER, EACHTTP_INVALID_REQUEST_JSON_FORMAT,
            LOAD_STRING (_T("IDS_EACHTTP_INVALID_JSON")));
    }

    // 获取第三方认证函数
    string thirdPartyId = requestJson["thirdpartyid"].s ();

    ncTThirdPartyAuthConf config;
    if(thirdPartyId != "anyshare" && thirdPartyId != "aishu") {
        // 检查第三方认证是否开启
        config = ncEACHttpServerUtil::Usrm_GetThirdPartyAuth ();
        if(thirdPartyId == "windows_ad_sso") {
            // 检查是否开启了域用户自动登录
            if(ncEACHttpServerUtil::Usrm_GetADSSOStatus() == false) {
                THROW_E (EAC_HTTP_SERVER, EACHTTP_AD_SSO_NOT_ENABLED,
                    LOAD_STRING (_T("IDS_AD_SSO_NOT_ENABLED")));
            }
        }
        else {
            if(config.thirdPartyId != thirdPartyId) {
                THROW_E (EAC_HTTP_SERVER, EACHTTP_THIRD_AUTH_NOT_OPEN, LOAD_STRING (_T("IDS_EACHTTP_THIRD_AUTH_NOT_OPEN")));
            }
        }
    }

    // 获取登录IP信息
    String loginIp;
    CheckLoginIpParam(requestJson, loginIp);

    //更新报文头的ip
    updateHTTPHeaderIP(cntl, loginIp);

    // 获取登录的设备信息
    ncDeviceBaseInfo baseInfo;
    bool hasDeviceInfo = parseDeviceInfo(requestJson, baseInfo);

    // 将clienttype信息加入params中
    if (hasDeviceInfo) {
        requestJson["params"]["clienttype"] = static_cast<int>(baseInfo.clientType);
    }

    // 执行认证
    map<String, ncThirdAuthFunc>::iterator iter = _thirdAuthFuncs.find (thirdPartyId.c_str());

    // 转给ShareMgnt进行认证
    String account;
    if(iter == _thirdAuthFuncs.end()) {
        string params;
        requestJson["deviceinfo"]["X-Real-IP"] = loginIp.getCStr();
        requestJson["deviceinfo"]["name"] = requestJson["device"]["name"];
        requestJson["deviceinfo"]["ostype"] = static_cast<int>(baseInfo.clientType);
        requestJson["deviceinfo"]["devicetype"] = requestJson["device"]["description"];
        requestJson["deviceinfo"]["udid"] = baseInfo.udid.getCStr ();
        requestJson["deviceinfo"]["udids"] = requestJson["device"]["udids"];
        JSON::Writer::write (requestJson.o (), params);

        account = ncEACHttpServerUtil::Usrm_ValidateThirdParty(params);
    }
    // 在eacp内部认证
    else {
        ncThirdAuthFunc func = iter->second;
        account = (this->*func) (requestJson, config);
    }

    ncACSUserInfo userInfo;
    int accountType = 0;
    if(_acsShareMgnt->GetUserInfoByAccount(account, userInfo, accountType) == false) {
        THROW_E (EAC_HTTP_SERVER, EACHTTP_NOT_IMPORT_TO_ANYSHARE,
            LOAD_STRING (_T("IDS_EACHTTP_NOT_IMPORT_TO_ANYSHARE")));
    }

    // 用户是否启用
    if(userInfo.enableStatus == false) {
        THROW_E (EAC_HTTP_SERVER, EACHTTP_USER_DISABLED,
            LOAD_STRING (_T("IDS_USER_DISABLED")));
    }

    // 绑定设备管理
    onUserLogin(cntl, userInfo.id, loginIp, hasDeviceInfo, baseInfo);

#ifndef __UT__
    // 更新用户最近请求时间记录
    _acsShareMgnt->UpdateUserLastRequestTime(userInfo.id);
#endif

    // 回复
    JSON::Value replyJson;
    replyJson["user_id"] = userInfo.id.getCStr ();
    replyJson["context"]["visitor_type"] = _visitorStringTypeMap[1].getCStr();
    replyJson["context"]["udid"] = baseInfo.udid.getCStr ();
    replyJson["context"]["login_ip"] = loginIp.getCStr();
    replyJson["context"]["client_type"] = _clientStringTypeMap[baseInfo.clientType].getCStr();
    replyJson["context"]["account_type"] = _accountStringTypeMap[accountType].getCStr();

    string body;
    JSON::Writer::write (replyJson.o (), body);
    ncHttpSendReply (cntl, brpc::HTTP_STATUS_OK, "ok", body);

    ncEACHttpServerUtil::Log (cntl, userInfo.id, ncTokenVisitorType::REALNAME, ncTLogType::NCT_LT_LOGIN, ncTLogLevel::NCT_LL_INFO,
        ncTLoginType::NCT_CLT_AUTHENICATION, LOAD_STRING (_T("IDS_AUTHENICATION_SUCCESS")), "");

    NC_EAC_HTTP_SERVER_TRACE (_T("this: %p, cntl:%p, userId:%s end"), this, cntl, userInfo.id.getCStr ());
    NC_EAC_HTTP_SERVER_CATCH (cntl)
}

void ncEACAuthHandler::GetByTicket (brpc::Controller* cntl, const String& fakeUserId)
{
    NC_EAC_HTTP_SERVER_TRACE (_T("this: %p, cntl:%p, userId:%s begin"), this, cntl, fakeUserId.getCStr ());

    // 获取http请求的content
    string bodyBuffer = cntl->request_attachment ().to_string ();

    // 获取ticket
    JSON::Value requestJson;
    try {
        JSON::Reader::read (requestJson, bodyBuffer.c_str (), bodyBuffer.size ());
    }
    catch (Exception& e) {
        THROW_E (EAC_HTTP_SERVER, EACHTTP_INVALID_REQUEST_JSON_FORMAT,
            LOAD_STRING (_T("IDS_EACHTTP_INVALID_JSON")));
    }

    String ticket = requestJson["ticket"].s ().c_str ();
    if (ticket.isEmpty ()) {
        THROW_E (EAC_HTTP_SERVER, EACHTTP_TICKET_INVALID,
            LOAD_STRING (_T("IDS_EACHTTP_TICKET_INVALID")));
    }

    String service = requestJson["service"].s ().c_str ();
    if (service.isEmpty ()) {
        THROW_E (EAC_HTTP_SERVER, EACHTTP_SERVICE_INVALID,
            LOAD_STRING (_T("IDS_EACHTTP_SERVICE_INVALID")));
    }

    // 检查是否开启了oauth
    ncTThirdPartyAuthConf config = ncEACHttpServerUtil::Usrm_GetThirdPartyAuth ();
    if(config.thirdPartyId != "wisedu" && config.thirdPartyId != "wisedu_sync") {
        THROW_E (EAC_HTTP_SERVER, EACHTTP_THIRD_AUTH_NOT_OPEN,
            LOAD_STRING (_T("IDS_EACHTTP_THIRD_AUTH_NOT_OPEN")));
    }
    JSON::Value configJson;
    JSON::Reader::read (configJson, config.config.c_str(), config.config.length());

    // 进行OAuth认证
    String thirdId = OAuthExecute (toCFLString(configJson["validateServer"].s()), ticket, service);

    // 获取anyshare用户信息
    ncACSUserInfo userInfo;
    bool ret = _acsShareMgnt->GetUserInfoByThirdId (thirdId, userInfo);
    if (ret == false) {
        THROW_E (EAC_HTTP_SERVER, EACHTTP_NOT_IMPORT_TO_ANYSHARE,
            LOAD_STRING (_T("IDS_EACHTTP_NOT_IMPORT_TO_ANYSHARE")));
    }

    // 用户是否启用
    ret = _acsShareMgnt->IsUserEnabled(userInfo.id);
    if(ret == false) {
        THROW_E (EAC_HTTP_SERVER, EACHTTP_USER_DISABLED,
            LOAD_STRING (_T("IDS_USER_DISABLED")));
    }

#ifndef __UT__
    // 更新用户最近请求时间记录
    _acsShareMgnt->UpdateUserLastRequestTime(userInfo.id);
#endif


    // 回复
    JSON::Value replyJson;
    replyJson["user_id"] = userInfo.id.getCStr ();
    replyJson["context"]["visitor_type"] = _visitorStringTypeMap[1].getCStr();
    replyJson["context"]["udid"] = "";
    replyJson["context"]["login_ip"] = ncEACHttpServerUtil::GetForwardedIp(cntl).getCStr();
    replyJson["context"]["client_type"] = _clientStringTypeMap[ACSClientType::UNKNOWN].getCStr();
    replyJson["context"]["account_type"] = _accountStringTypeMap[0].getCStr();

    string body;
    JSON::Writer::write (replyJson.o (), body);
    ncHttpSendReply (cntl, brpc::HTTP_STATUS_OK, "ok", body);

    String msg;
    msg.format (ncEACHttpServerLoader, _T("IDS_LOGIN_SUCCESS"), "Unknown");
    ncEACHttpServerUtil::Log (cntl, userInfo.id, ncTokenVisitorType::REALNAME, ncTLogType::NCT_LT_LOGIN, ncTLogLevel::NCT_LL_INFO,
        ncTLoginType::NCT_CLT_LOGIN_IN, msg, "");

    NC_EAC_HTTP_SERVER_TRACE (_T("this: %p, cntl:%p, userId:%s end"), this, cntl, fakeUserId.getCStr ());
}

void ncEACAuthHandler::GetByADSession (brpc::Controller* cntl, const String& fakeUserId)
{
    NC_EAC_HTTP_SERVER_TRACE (_T("this: %p, cntl:%p, userId:%s begin"), this, cntl, fakeUserId.getCStr ());

    // 获取http请求的content
    string bodyBuffer = cntl->request_attachment ().to_string ();

    // 获取adsession
    JSON::Value requestJson;
    try {
        JSON::Reader::read (requestJson, bodyBuffer.c_str (), bodyBuffer.size ());
    }
    catch (Exception& e) {
        THROW_E (EAC_HTTP_SERVER, EACHTTP_INVALID_REQUEST_JSON_FORMAT,
            LOAD_STRING (_T("IDS_EACHTTP_INVALID_JSON")));
    }

    // 检查是否开启了域用户自动登录
    if(ncEACHttpServerUtil::Usrm_GetADSSOStatus() == false) {
        THROW_E (EAC_HTTP_SERVER, EACHTTP_AD_SSO_NOT_ENABLED,
            LOAD_STRING (_T("IDS_AD_SSO_NOT_ENABLED")));
    }

    // 获取session
    String session = requestJson["adsession"].s ().c_str ();
    if (session.isEmpty ()) {
        THROW_E (EAC_HTTP_SERVER, EACHTTP_INVALID_AD_SESSION,
            LOAD_STRING (_T("IDS_INVALID_AD_SESSION")));
    }

    String account;
    ParseADSession(session, account);

    ncACSUserInfo userInfo;
    int accountType = 0;
    if(_acsShareMgnt->GetUserInfoByAccount(account, userInfo, accountType) == false) {
        THROW_E (EAC_HTTP_SERVER, EACHTTP_INVALID_AD_SESSION,
            LOAD_STRING (_T("IDS_INVALID_AD_SESSION")));
    }

    // 系统管理员不允许登录
    if(_adminIds.count(userInfo.id.getCStr())) {
        THROW_E(EAC_HTTP_SERVER, EACHTTP_FORBIDDEN_LOGIN,
            "%s are not allowed to log in.", account.getCStr());;
    }

    // 用户是否启用
    if(userInfo.enableStatus == false) {
        THROW_E (EAC_HTTP_SERVER, EACHTTP_USER_DISABLED,
            LOAD_STRING (_T("IDS_USER_DISABLED")));
    }

#ifndef __UT__
    // 更新用户最近请求时间记录
    _acsShareMgnt->UpdateUserLastRequestTime(userInfo.id);
#endif

    // 回复
    JSON::Value replyJson;
    replyJson["user_id"] = userInfo.id.getCStr ();
    replyJson["context"]["visitor_type"] = _visitorStringTypeMap[1].getCStr();
    replyJson["context"]["udid"] = "";
    replyJson["context"]["login_ip"] = ncEACHttpServerUtil::GetForwardedIp(cntl).getCStr();
    replyJson["context"]["client_type"] = _clientStringTypeMap[ACSClientType::UNKNOWN].getCStr();
    replyJson["context"]["account_type"] = _accountStringTypeMap[accountType].getCStr();

    string body;
    JSON::Writer::write (replyJson.o (), body);
    ncHttpSendReply (cntl, brpc::HTTP_STATUS_OK, "ok", body);

    String msg;
    msg.format (ncEACHttpServerLoader, _T("IDS_LOGIN_SUCCESS"), "Unknown");
    ncEACHttpServerUtil::Log (cntl, userInfo.id, ncTokenVisitorType::REALNAME, ncTLogType::NCT_LT_LOGIN, ncTLogLevel::NCT_LL_INFO,
        ncTLoginType::NCT_CLT_LOGIN_IN, msg, "");

    ncEACHttpServerUtil::LoginLog(userInfo.id, "", ACSClientType::UNKNOWN, ncEACHttpServerUtil::GetForwardedIp(cntl));

    NC_EAC_HTTP_SERVER_TRACE (_T("this: %p, cntl:%p, userId:%s end"), this, cntl, fakeUserId.getCStr ());
}

void ncEACAuthHandler::LoginLog(brpc::Controller* cntl)
{
    NC_EAC_HTTP_SERVER_TRY
    cntl->http_response().set_content_type("application/json; charset=UTF-8");
    // 判断请求方法是否为POST
    if (cntl->http_request ().method () != brpc::HTTP_METHOD_POST){
        THROW_E (EACHTTP_SERVER_ERR_PROVIDER_NAME, EACHTTP_SERVER_METHOD_NOT_IMPLEMENTED,
                 LOAD_STRING (_T("IDS_EACHTTP_METHOD_INVALID")));
    }
    NC_EAC_HTTP_SERVER_TRACE (_T("this: %p, cntl:%p begin"), this, cntl);

    // 获取http请求的content
    string bodyBuffer = cntl->request_attachment ().to_string ();

    // 获取id,udid,client-type
    JSON::Value requestJson;
    try {
        JSON::Reader::read (requestJson, bodyBuffer.c_str (), bodyBuffer.size ());
    }
    catch (Exception& e) {
        THROW_E (EAC_HTTP_SERVER, EACHTTP_INVALID_REQUEST_JSON_FORMAT,
            LOAD_STRING (_T("IDS_EACHTTP_INVALID_JSON")));
    }

    string id = requestJson["id"].s ();
    if (id.empty ()) {
        THROW_E (EAC_HTTP_SERVER, EACHTTP_INVALID_PARAMETER,
            LOAD_STRING (_T("IDS_EACHTTP_PARAMETERS_IS_REQUIRED")), "id");
    }

    string strClientType = requestJson["client_type"].s ();
    if (strClientType.empty ()) {
        THROW_E (EAC_HTTP_SERVER, EACHTTP_INVALID_PARAMETER,
            LOAD_STRING (_T("IDS_EACHTTP_PARAMETERS_IS_REQUIRED")), "client_type");
    }

    ACSClientType clientType;
    auto iter = _clientIntTypeMap.find(toCFLString(strClientType));
    if(iter != _clientIntTypeMap.end())
    {
        clientType = static_cast<ACSClientType>(iter->second);
    }
    else
    {
        THROW_E (EAC_HTTP_SERVER, EACHTTP_INVALID_PARAMETER,
            LOAD_STRING (_T("IDS_INVALID_OS_TYPE")));
    }

    string ip = requestJson["ip"].s ();
    string udid = requestJson["udid"].s ();

    ncEACHttpServerUtil::LoginLog (toCFLString(id), toCFLString(udid), clientType, toCFLString(ip));

    string body;
    ncHttpSendReply (cntl, brpc::HTTP_STATUS_NO_CONTENT, "ok", body);
    NC_EAC_HTTP_SERVER_TRACE (_T("this: %p, cntl:%p, userId:%s end"), this, cntl, id.c_str ());
    NC_EAC_HTTP_SERVER_CATCH (cntl)
}

// protected
void
ncEACAuthHandler::ModifyPassword(brpc::Controller* cntl, const String& userId)
{
    NC_EAC_HTTP_SERVER_TRACE (_T("this: %p, cntl:%p begin"), this, cntl);

    // 获取http请求的content
    string bodyBuffer = cntl->request_attachment ().to_string ();
    String bodyContent (bodyBuffer.c_str (), bodyBuffer.size ());

    // 签名认证
    CheckSign (cntl, bodyContent);

    // 获取account, oldPwd, newPwd
    JSON::Value requestJson;
    try {
        JSON::Reader::read (requestJson, bodyBuffer.c_str (), bodyBuffer.size ());
    }
    catch (Exception& e) {
        THROW_E (EAC_HTTP_SERVER, EACHTTP_INVALID_REQUEST_JSON_FORMAT,
            LOAD_STRING (_T("IDS_EACHTTP_INVALID_JSON")));
    }

    string account (requestJson["account"].s ());
    string oldPwd;
    string bs64_oldpwd;
    if (!requestJson["oldpwd"].s ().empty())
    {
        bs64_oldpwd = ncEACHttpServerUtil::Base64Decode(requestJson["oldpwd"].s());
        oldPwd = ncEACHttpServerUtil::RSADecrypt(bs64_oldpwd);
    }
    string bs64_newpwd(ncEACHttpServerUtil::Base64Decode(requestJson["newpwd"].s()));
    string newPwd ( ncEACHttpServerUtil::RSADecrypt(bs64_newpwd) );

    // 获取验证码信息
    ncTUserModifyPwdOption option;
    option.uuid = requestJson["vcodeinfo"]["uuid"].s ();
    option.__isset.uuid = true;
    option.vcode = requestJson["vcodeinfo"]["vcode"].s ();
    option.__isset.vcode = true;
    option.isForgetPwd = requestJson["isforgetpwd"].b ();
    option.__isset.isForgetPwd = true;
    ncACSUserInfo userInfo;
    userInfo.account = toCFLString(account);
    string emailaddress;
    string telnumber;
    if(option.isForgetPwd)
    {
        // 忘记密码时获取邮箱或手机号
        emailaddress = requestJson["emailaddress"].s();
        telnumber = requestJson["telnumber"].s();
        // 邮箱或手机都为空或都不为空或uuid、vcode为空，抛参数异常
        if((emailaddress.empty() && telnumber.empty())    ||
            (!emailaddress.empty() && !telnumber.empty()) ||
            option.uuid.empty()                           ||
            option.vcode.empty())
        {
            THROW_E (EAC_HTTP_SERVER, EACHTTP_INVALID_PARAM_VALUE,
                LOAD_STRING (_T("IDS_EACHTTP_ERR_MSG_AGR_ERR")));
        }

        // 去除密码和邮箱是否和用户匹配的检查，出现场景很极限，除非在发送验证码之后修改了用户的手机和邮箱，才会出现这个报错
        // 而且出现这个问题 用户会觉得发送验证码没问题，但是修改密码的时候又提示绑定错误的情况很迷惑
        // if(!telnumber.empty() && _acsShareMgnt->GetUserInfoByTelNumber(toCFLString(telnumber), userInfo) == false){
        //     THROW_E (EAC_HTTP_SERVER, EACHTTP_PHONE_HAS_NOT_BEEN_BOUND,
        //         LOAD_STRING (_T("IDS_EACHTTP_PHONE_HAS_NOT_BEEN_BOUND")));
        // }
        // else if(!emailaddress.empty() && _acsShareMgnt->GetUserInfoByEmail(toCFLString(emailaddress), userInfo) == false){
        //     THROW_E (EAC_HTTP_SERVER, EACHTTP_EMAIL_ADDRESS_HAS_NOT_BEEN_BOUND,
        //             LOAD_STRING (_T("IDS_EACHTTP_EMAIL_ADDRESS_HAS_NOT_BEEN_BOUND")));
        // }
    }

    // 禁用的用户不能修改密码
    int accountType = 0;
    if(_acsShareMgnt->GetUserInfoByAccount(userInfo.account.getCStr(), userInfo, accountType) == false) {
        THROW_E (EAC_HTTP_SERVER, EACHTTP_INVALID_USER_OR_PASSWORD, "Invalid user or password.");
    }
    // 由于admin的禁用字段表示权责分离是否开启，导致这里要忽略admin
    if(userInfo.id.compare(toCFLString(g_ShareMgnt_constants.NCT_USER_ADMIN)) != 0  && userInfo.enableStatus == false) {
        THROW_E (EAC_HTTP_SERVER, EACHTTP_USER_DISABLED, LOAD_STRING (_T("IDS_USER_DISABLED")));
    }

    if (_acsShareMgnt->IsAdminId(userInfo.id) && newPwd == "eisoo.com") {
        THROW_E (EAC_HTTP_SERVER, EACHTTP_NEW_PASSWORD_SHOULD_NOT_BE_DEFAULT_PASSWORD, "New password should not be default password.");
    }

    try {
        ncThriftClient<ncTShareMgntClient> shareMgntClient (EacServiceAccessConfig::getInstance()->sharemgntHost, EacServiceAccessConfig::getInstance()->sharemgntPort);
        shareMgntClient->Usrm_ModifyPassword (userInfo.account.getCStr(), oldPwd.c_str (), newPwd.c_str (), option);
    }
    catch (ncTException& e) {
        // 获取详细错误信息
        JSON::Value errDetail;
        try {
            if (!e.errDetail.empty ()) {
                JSON::Reader::read (errDetail, e.errDetail.c_str (), e.errDetail.length ());
            }
        }
        catch (Exception& e) {
            THROW_E (EAC_HTTP_SERVER, EACHTTP_INVALID_REQUEST_JSON_FORMAT,
                LOAD_STRING (_T("IDS_EACHTTP_INVALID_JSON")));
        }

        if (e.errID == ncTShareMgntError::NCT_INVALID_PASSWORD) {
            THROW_E (EAC_HTTP_SERVER, EACHTTP_INVALID_PASSWORD, e.expMsg.c_str ());
        }
        else if (e.errID == ncTShareMgntError::NCT_INVALID_STRONG_PASSWORD) {
            THROW_E (EAC_HTTP_SERVER, EACHTTP_INVALID_STRONG_PASSWORD, e.expMsg.c_str ());
        }
        else if (e.errID == ncTShareMgntError::NCT_USER_NOT_EXIST ) {
            THROW_E (EAC_HTTP_SERVER, EACHTTP_USER_NOT_EXIST, e.expMsg.c_str ());
        }
        else if (e.errID == ncTShareMgntError::NCT_CHECK_PASSWORD_FAILED) {
            THROW_HTTP_DETAIL_E (EAC_HTTP_SERVER, errDetail, EACHTTP_INVALID_USER_OR_PASSWORD, e.expMsg.c_str ());
        }
        else if (e.errID == ncTShareMgntError::NCT_CANNOT_MODIFY_NONLOCAL_USER_PASSWORD) {
            THROW_E (EAC_HTTP_SERVER, EACHTTP_NOT_LOCAL_USER, e.expMsg.c_str ());
        }
        else if (e.errID == ncTShareMgntError::NCT_PWD_FIRST_FAILED) {
            THROW_E (EAC_HTTP_SERVER, EACHTTP_PWD_FAILED_FIRSTLY, e.expMsg.c_str ());
        }
        else if (e.errID == ncTShareMgntError::NCT_PWD_SECOND_FAILED) {
            THROW_E (EAC_HTTP_SERVER, EACHTTP_PWD_FAILED_SECONDLY, e.expMsg.c_str ());
        }
        else if (e.errID == ncTShareMgntError::NCT_ACCOUNT_LOCKED) {
            JSON::Value detailJson;
            detailJson["remainlockTime"] = ncEACHttpServerUtil::ParseLockTime(e.expMsg.c_str ());
            detailJson["isShowStatus"] = errDetail["isShowStatus"];
            THROW_HTTP_DETAIL_E(EAC_HTTP_SERVER, detailJson, EACHTTP_ACCOUNT_LOCKED, e.expMsg.c_str ());
        }
        else if(e.errID == ncTShareMgntError::NCT_WRONG_PASSWORD) {
            THROW_E(EAC_HTTP_SERVER, EACHTTP_WRONG_PASSWORD, e.expMsg.c_str ());
        }
        else if(e.errID == ncTShareMgntError::NCT_CHECK_VCODE_IS_NULL){
            THROW_HTTP_DETAIL_E(EAC_HTTP_SERVER, errDetail, EACHTTP_CHECK_VCODE_IS_NULL, e.expMsg.c_str ());
        }
        else if(e.errID == ncTShareMgntError::NCT_CHECK_VCODE_IS_TIMEOUT){
            THROW_HTTP_DETAIL_E(EAC_HTTP_SERVER, errDetail, EACHTTP_CHECK_VCODE_IS_TIMEOUT, e.expMsg.c_str ());
        }
        else if(e.errID == ncTShareMgntError::NCT_CHECK_VCODE_IS_WRONG){
            THROW_HTTP_DETAIL_E(EAC_HTTP_SERVER, errDetail, EACHTTP_CHECK_VCODE_IS_WRONG, e.expMsg.c_str ());
        }
        else if(e.errID == ncTShareMgntError::NCT_CHECK_VCODE_MORE_THAN_THE_LIMIT){
            THROW_HTTP_DETAIL_E(EAC_HTTP_SERVER,errDetail, EACHTTP_CHECK_VCODE_MORE_THAN_THE_LIMIT, e.expMsg.c_str ());
        }
        else if(e.errID == ncTShareMgntError::NCT_CANNOT_MODIFY_CONTROL_PASSWORD) {
            THROW_E(EAC_HTTP_SERVER, EACHTTP_CANNOT_MODIFY_CONTROL_PASSWORD, e.expMsg.c_str ());
        }
        else if(e.errID == ncTShareMgntError::NCT_PWD_THIRD_FAILED) {
            // 密码连续输错3次，记录日志
            ncACSUserInfo userInfo;
            int accountType = 0;
            _acsShareMgnt->GetUserInfoByAccount(toCFLString(account), userInfo, accountType);
            String msg;
            msg.format (ncEACHttpServerLoader, _T("IDS_MODIFY_PASSWORD_FAILED_AND_LOCKED"));

            ncEACHttpServerUtil::Log (cntl, userInfo.id, ncTokenVisitorType::REALNAME, ncTLogType::NCT_LT_MANAGEMENT, ncTLogLevel::NCT_LL_WARN,
               ncTManagementType::NCT_MNT_SET, msg, LOAD_STRING("IDS_INPUT_WRONG_PASSWORD_FOR_THREE_TIMES"));

            if (e.errID == ncTShareMgntError::NCT_PWD_THIRD_FAILED) {
                JSON::Value detailJson;
                detailJson["remainlockTime"] = ncEACHttpServerUtil::ParseLockTime(e.expMsg.c_str ());
                detailJson["isShowStatus"] = errDetail["isShowStatus"];
                THROW_HTTP_DETAIL_E(EAC_HTTP_SERVER, detailJson, EACHTTP_PWD_FAILED_THIRDLY, e.expMsg.c_str ());
            }
        }
        else if (e.errID == ncTShareMgntError::NCT_PASSWORD_IS_INITIAL) {
            THROW_E (EAC_HTTP_SERVER, EACHTTP_NEW_PASSWORD_SHOULD_NOT_BE_DEFAULT_PASSWORD, "New password should not be default password.");
        }
        else {
            THROW_E (EAC_HTTP_SERVER, EACHTTP_MODIFY_PASSWORD, e.expMsg.c_str ());
        }
    }

    // 密码修改完毕，需要清理掉用户对应的token
    _acsTokenManager->DeleteTokenByUserId (userInfo.id);

    ncHttpSendReply (cntl, brpc::HTTP_STATUS_OK, "ok", "");
    String msg;
    String exmsg;
    if(option.isForgetPwd){
        // 忘记密码修改成功日志
        msg.format (ncEACHttpServerLoader, _T("IDS_RESET_PASSWORD_SUCCESS"));
        if(!telnumber.empty())
        {
            exmsg.format(LOAD_STRING(_T("IDS_RESET_PASSWORD_SUCCESS_BY_SMS")));
        }
        else if(!emailaddress.empty())
        {
            exmsg.format(LOAD_STRING(_T("IDS_RESET_PASSWORD_SUCCESS_BY_EMAIL")));
        }
    }
    else{
        msg.format (ncEACHttpServerLoader, _T("IDS_MODIFY_PASSWORD_SUCCESS"), userInfo.visionName.getCStr ());
    }

    ncEACHttpServerUtil::Log (cntl, userInfo.id, ncTokenVisitorType::REALNAME, ncTLogType::NCT_LT_MANAGEMENT, ncTLogLevel::NCT_LL_WARN,
       ncTManagementType::NCT_MNT_SET, msg, exmsg);

    NC_EAC_HTTP_SERVER_TRACE (_T("this: %p, cntl:%p, userId:%s end"), this, cntl, userId.getCStr());
}

// protected
void ncEACAuthHandler::ValidateSecurityDevice (brpc::Controller* cntl, const String& fakeUserId)
{
    NC_EAC_HTTP_SERVER_TRACE (_T("this: %p, cntl:%p begin"), this, cntl);

    // 获取http请求的content
    string bodyBuffer = cntl->request_attachment ().to_string ();
    JSON::Value requestJson;
    try {
        JSON::Reader::read (requestJson, bodyBuffer.c_str (), bodyBuffer.size ());
    }
    catch (Exception& e) {
        THROW_E (EAC_HTTP_SERVER, EACHTTP_INVALID_REQUEST_JSON_FORMAT,
            LOAD_STRING (_T("IDS_EACHTTP_INVALID_JSON")));
    }

    string params;
    JSON::Writer::write (requestJson.o (), params);

    bool result = false;

    try {
        result = ncEACHttpServerUtil::Usrm_ValidateSecurityDevice(params);
    }
    catch(Exception& e) {
        THROW_E (EAC_HTTP_SERVER, EACHTTP_THIRD_SERVER_UNAVALIABLE, LOAD_STRING (_T("IDS_INVALID_THIRD_PARTY_TICKET")));
    }

    // 记录安全设备认证日志
    string account = requestJson["params"]["account"].s ();
    ncACSUserInfo userInfo;
    int accountType = 0;
    _acsShareMgnt->GetUserInfoByAccount(toCFLString(account), userInfo, accountType);
    if (result) {
        ncEACHttpServerUtil::Log(cntl, userInfo.id, ncTokenVisitorType::REALNAME, ncTLogType::NCT_LT_LOGIN, ncTLogLevel::NCT_LL_INFO,
                                 ncTLoginType::NCT_CLT_LOGIN_IN, LOAD_STRING (_T("IDS_SECURITY_DEVICE_AUTH_SUCCESS")), "");
    }
    else {
        ncEACHttpServerUtil::Log(cntl, userInfo.id, ncTokenVisitorType::REALNAME, ncTLogType::NCT_LT_LOGIN, ncTLogLevel::NCT_LL_WARN,
                                 ncTLoginType::NCT_CLT_LOGIN_IN, LOAD_STRING (_T("IDS_SECURITY_DEVICE_AUTH_FAILED")), "");
    }

    JSON::Value replyJson;
    replyJson["result"] = result;
    string body;
    JSON::Writer::write (replyJson.o (), body);
    ncHttpSendReply (cntl, brpc::HTTP_STATUS_OK, "ok", body);
}

String ncEACAuthHandler::OAuthExecute (const String& tokenServer, const String& ticket, const String& service)
{
    String url;
    url.format (_T("%s?ticket=%s&service=%s"), tokenServer.getCStr (), ticket.getCStr (), service.getCStr ());

    nsresult ret;
    nsCOMPtr<ncIEHTTPClient> httpClient = do_CreateInstance (NC_EHTTP_CLIENT_CONTRACTID, &ret);
    if (NS_FAILED (ret)) {
        THROW_E (EAC_HTTP_SERVER, FAILED_TO_CREATE_EHTTP_CLIENT,
                    "Failed to create ehttpclient instance: 0x%x", ret);
    }

    ncEHTTPResponse response;
    httpClient->Get (toSTLString (url), 30, response);

    String userId = ParseThirdUserId (response.body);

    return userId;
}

String ncEACAuthHandler::ParseThirdUserId (const string& retXMlStr)
{
    // Create an empty property tree object
    ptree ptAll;

    // Read the XML config string into the property tree. Catch any exception
    try {
        stringstream ss;
        ss << retXMlStr;
        read_xml(ss, ptAll);
    }
    catch (xml_parser_error &e) {
        THROW_E (EAC_HTTP_SERVER, INVALID_XML_FORMAT,
            LOAD_STRING (_T("IDS_INVALID_XML_FORMAT")));
    }

    boost::property_tree::ptree ptResponse = ptAll.get_child ("cas:serviceResponse");
    /* 失败消息
    <cas:serviceResponse xmlns:cas='http://www.yale.edu/tp/cas'>
        <cas:authenticationFailure code='INVALID_TICKET'>
            ticket &#039;SST-1530-TZOn1Kbo5afbE5n3oyVV-yBMM-ids1-1407559557579&#039; not recognized
        </cas:authenticationFailure>
    </cas:serviceResponse>
    */
    ptree::const_assoc_iterator itFail = ptResponse.find ("cas:authenticationFailure");
    if (itFail != ptResponse.not_found()) {
        string error = ptResponse.get<string> ("cas:authenticationFailure");
        THROW_E (EAC_HTTP_SERVER, CANT_AUTHENTICATE_TICKET,
            LOAD_STRING (_T("IDS_CANT_AUTHENTICATE_TICKET")), error.c_str ());
    }

    /* 成功消息
    <cas:serviceResponse xmlns:cas='http://www.yale.edu/tp/cas'>
        <cas:authenticationSuccess>
            <cas:user>test</cas:user>
            <cas:attributes>
                <cas:containerId>test</cas:containerId>
                <cas:bin duserlist></cas:binduserlist>
                <cas:user_name>test</cas:user_name>
            </cas:attributes>
        </cas:authenticationSuccess>
    </cas:serviceResponse>
    */
    ptree::const_assoc_iterator itSuc = ptResponse.find ("cas:authenticationSuccess");
    if (itSuc != ptResponse.not_found()) {

        ptree p1 = ptResponse.get_child ("cas:authenticationSuccess");
        string thirdId = p1.get<string> ("cas:user");

        return toCFLString (thirdId);
    }

    THROW_E (EAC_HTTP_SERVER, INVALID_XML_FORMAT,
        LOAD_STRING (_T("IDS_INVALID_XML_FORMAT")));
}

void ncEACAuthHandler::ParseADSession (const String& session, String& account)
{
    String content;

    // 解密原文
    try {
        string sess = toSTLString(session);
        string decodeContent(ncEACHttpServerUtil::Base64Decode(sess));
        string sourceContent (ncEACHttpServerUtil::RSADecrypt(decodeContent));

        // 解析account和key
        content = toCFLString(sourceContent);
    }
    catch(Exception& e) {
        THROW_E (EAC_HTTP_SERVER, EACHTTP_INVALID_AD_SESSION, LOAD_STRING (_T("IDS_INVALID_AD_SESSION")));
    }

    // 检查原文格式是否正确
    vector<String> strs;
    content.split('\\', strs);

    if(strs.size() != 2) {
        THROW_E (EAC_HTTP_SERVER, EACHTTP_INVALID_AD_SESSION, LOAD_STRING (_T("IDS_INVALID_AD_SESSION")));
    }

    String tmpStr = "E1s)o.C0MieD1ievohl";
    if(strs[1] != tmpStr) {
        THROW_E (EAC_HTTP_SERVER, EACHTTP_INVALID_AD_SESSION, LOAD_STRING (_T("IDS_INVALID_AD_SESSION")));
    }

    account = strs[0];
}

string ncEACAuthHandler::RSAEncrypt (const string& plainText, const char *path_key)
{
    string cipherText;
    RSA *p_rsa;
    FILE *file;
    int rsa_len = 0;

    if ((file = fopen (path_key,"r")) == NULL) {
        throw Exception (_T("open key file error"));
    }
    if ((p_rsa = PEM_read_RSA_PUBKEY (file, NULL, NULL, NULL)) == NULL) {
        throw Exception (_T("PEM_read_RSA_PUBKEY error"));
    }

    rsa_len = RSA_size (p_rsa);
    cipherText.resize (rsa_len, 0);

    int num = RSA_public_encrypt (plainText.length(), (unsigned char *)plainText.c_str(), (unsigned char*)cipherText.c_str(), p_rsa, RSA_PKCS1_PADDING);
    if (num != rsa_len){
        RSA_free (p_rsa);
        fclose (file);
        throw Exception (_T("RSA_public_encrypt error"));
    }

    RSA_free (p_rsa);
    fclose (file);

    return cipherText;
}

string ncEACAuthHandler::RSADecrypt (const string& cipherText, const char *path_key)
{
    // printf("Before RSADecrypt: \n");
    // printBytes(cipherText);

    static ThreadMutexLock sLock;
    static RSA* p_rsa = NULL;
    AutoLock<ThreadMutexLock> lock (&sLock);
    if(p_rsa == NULL) {
        FILE *file;
        if ((file = fopen (path_key,"r")) == NULL) {
            throw Exception(_T("fopen error"));
        }
        if ((p_rsa = PEM_read_RSAPrivateKey (file, NULL, NULL, NULL)) == NULL) {
            fclose (file);
            throw Exception(_T("PEM_read_RSAPrivateKey error"));
        }
        fclose (file);
    }

    string plainText;
    int rsa_len = RSA_size (p_rsa);
    plainText.resize (rsa_len, 0);

    if (RSA_private_decrypt (cipherText.length(), (unsigned char *)cipherText.c_str(), (unsigned char*)plainText.c_str(), p_rsa, RSA_PKCS1_PADDING)<0) {
        THROW_E (EAC_HTTP_SERVER, EACHTTP_INVALID_PARAMETER, "RSA_private_decrypt error");
    }

    // 去除末尾的\0
    string ret(plainText.c_str());

    // printf("After RSADecrypt: %s\n", ret.c_str());
    // printBytes(ret);
    return ret;
}

string ncEACAuthHandler::Base64Encode(const string& input)
{
    BIO * bmem = NULL;
    BIO * b64 = NULL;
    BUF_MEM * bptr = NULL;

    b64 = BIO_new(BIO_f_base64());
    bmem = BIO_new(BIO_s_mem());
    b64 = BIO_push(b64, bmem);
    BIO_write(b64, (char*)input.c_str(), input.length());
    BIO_flush(b64);
    BIO_get_mem_ptr(b64, &bptr);

    string buffer;
    buffer.assign(bptr->data, bptr->length);

    BIO_free_all(b64);

    return buffer;
}

string ncEACAuthHandler::DESEncrypt(const string& plainText)
{
    string cipherText;
    int cipherTextLength = 0;
    if (plainText.length() % 8 == 0) {
        cipherTextLength = plainText.length();
    }
    else {
        cipherTextLength = plainText.length() + (8 - plainText.length() % 8);
    }
    cipherText.assign(cipherTextLength, '*');

    DES_cblock key = {'E', 't', 'm', 'B', '8', '?', 's', 'f'};
    DES_cblock ivec = {'E', 't', 'm', 'B', '8', '?', 's', 'f'};
    DES_key_schedule keysched;

    DES_set_odd_parity(&key);
    if (DES_set_key_checked((C_Block *)key, &keysched)) {
        throw Exception("Unable to set key schedule");
    }

    DES_ncbc_encrypt((unsigned char*)plainText.c_str(), (unsigned char*)cipherText.c_str(), plainText.length(), &keysched, &ivec, DES_ENCRYPT);

    return cipherText;
}

string ncEACAuthHandler::DESDecrypt(const string& cipherText)
{
    string plainText;
    plainText.assign(cipherText.length(), '*');

    DES_cblock key = {'E', 't', 'm', 'B', '8', '?', 's', 'f'};
    DES_cblock ivec = {'E', 't', 'm', 'B', '8', '?', 's', 'f'};
    DES_key_schedule keysched;

    DES_set_odd_parity(&key);
    if (DES_set_key_checked((C_Block *)key, &keysched)) {
        throw Exception("Unable to set key schedule");
    }

    DES_ncbc_encrypt((unsigned char*)cipherText.c_str(), (unsigned char*)plainText.c_str(), cipherText.length(), &keysched, &ivec, DES_DECRYPT);

    // 去除末尾的\0
    string ret(plainText.c_str());
    return ret;
}

void ncEACAuthHandler::printBytes(const string& str)
{
    for(int i = 0; i < str.length(); ++i) {
        printf(" [%02x]", (unsigned char)str[i]);
    }
    printf("\n");
}

String ncEACAuthHandler::BQKJAuth(JSON::Value& requestJson, const ncTThirdPartyAuthConf& config)
{
    // 获取ticket
    string ticket;
    try {
        ticket = requestJson["params"]["ticket"].s();
    }
    catch (Exception& e) {
        THROW_E(EAC_HTTP_SERVER, EACHTTP_INVALID_THIRD_PARTY_TICKET,
            LOAD_STRING("IDS_INVALID_THIRD_PARTY_TICKET"), ticket.c_str());
    }
    if(ticket == "") {
        THROW_E(EAC_HTTP_SERVER, EACHTTP_INVALID_THIRD_PARTY_TICKET,
            LOAD_STRING("IDS_INVALID_THIRD_PARTY_TICKET"), ticket.c_str());
    }

    // http://58.30.20.150:8180
    string authServer;
    try {
        JSON::Value configJson;
        JSON::Reader::read (configJson, config.config.c_str(), config.config.length());

        authServer = configJson["authServer"].s();
    }
    catch(Exception& e) {
        THROW_E (EAC_HTTP_SERVER, FAILED_TO_PARSE_THIRD_AUTH_CONFIG,
                    LOAD_STRING("IDS_FAILED_TO_PARSE_THIRD_AUTH_CONFIG"));
    }

    nsresult ret;
    nsCOMPtr<ncIEHTTPClient> httpClient = do_CreateInstance (NC_EHTTP_CLIENT_CONTRACTID, &ret);
    if (NS_FAILED (ret)) {
        THROW_E (EAC_HTTP_SERVER, FAILED_TO_CREATE_EHTTP_CLIENT,
                    "Failed to create ehttpclient instance: 0x%x", ret);
    }

    // http://58.30.20.150:8180/val.sso?ticketId=11223345
    string authURL = authServer + "/val.sso?ticketId=" + ticket;

    ncEHTTPResponse response;
    httpClient->Get(authURL, 30, response);
    if(response.body == "true\r\n") {
        // 获取用户帐号信息
        String tmpStr = toCFLString(ticket);
        vector<String> splitStrs;
        tmpStr.split('@', splitStrs);

        String account;
        if(splitStrs.size() == 3) {
            account = splitStrs[2];
        }

        return account;
    }
    else if(response.body == "false\r\n") {
        THROW_E(EAC_HTTP_SERVER, EACHTTP_INVALID_THIRD_PARTY_TICKET,
            LOAD_STRING("IDS_INVALID_THIRD_PARTY_TICKET"), ticket.c_str());
    }
    else {
        THROW_E(EAC_HTTP_SERVER, FAILED_TO_EXECUTE_THIRD_AUTH,
            LOAD_STRING("IDS_FAILED_TO_EXECUTE_THIRD_AUTH"));
    }
}

String ncEACAuthHandler::BeingAuth(JSON::Value& requestJson, const ncTThirdPartyAuthConf& config)
{
    // 获取ticket, target, suffix
    string ticket, target, suffix;
    try {
        ticket = requestJson["params"]["ticket"].s();
        target = requestJson["params"]["target"].s();
        suffix = requestJson["params"]["suffix"].s();

        char* tmpTicket = (char*)ticket.c_str();
        ticket = url_decode(tmpTicket);
    }
    catch (Exception& e) {
        THROW_E(EAC_HTTP_SERVER, EACHTTP_INVALID_THIRD_PARTY_TICKET,
            LOAD_STRING("IDS_INVALID_THIRD_PARTY_TICKET"), ticket.c_str());
    }
    if(ticket == "") {
        THROW_E(EAC_HTTP_SERVER, EACHTTP_INVALID_THIRD_PARTY_TICKET,
            LOAD_STRING("IDS_INVALID_THIRD_PARTY_TICKET"), ticket.c_str());
    }

    // 获取服务器地址
    // authServer: http://121.22.8.50:18080/secure
    string authServer;
    try {
        JSON::Value configJson;
        JSON::Reader::read (configJson, config.config.c_str(), config.config.length());

        authServer = configJson["authServer"].s();
    }
    catch(Exception& e) {
        THROW_E (EAC_HTTP_SERVER, FAILED_TO_PARSE_THIRD_AUTH_CONFIG,
            LOAD_STRING("IDS_FAILED_TO_PARSE_THIRD_AUTH_CONFIG"));
    }

    // 验证获取用户名
    // http://121.22.8.50:18080/secure/%s?TARGET=%s&appcode=anyshare
    String checkURL;
    string baseUrl = authServer + "/%s?TARGET=%s&appcode=anyshare";
    checkURL.format(baseUrl.c_str(), suffix.c_str(), target.c_str());

    // post body
    string xmlTmp =
        "<SOAP-ENV:Envelope xmlns:SOAP-ENV=\"http://schemas.xmlsoap.org/soap/envelope/\"><SOAP-ENV:Header />"
            "<SOAP-ENV:Body>"
                "<samlp:Request xmlns:samlp=\"urn:oasis:names:tc:SAML:1.0:protocol\""
                    "MajorVersion=\"1\" MinorVersion=\"1\" RequestID=\"_a9f7450396ca0fc31abf52c8db871eff\" IssueInstant=\"%s\">"
                    "<samlp:AssertionArtifact>%s</samlp:AssertionArtifact>"
                "</samlp:Request>"
            "</SOAP-ENV:Body>"
        "</SOAP-ENV:Envelope>";

    // 设置时间
    time_t rawtime;
    struct tm* TM;
    time(&rawtime);
    TM = localtime(&rawtime);

    // 时间格式为:%Y-%M-%dT%H:%m:%SZ
    string timeTemplate("%04d-%02d-%02dT%02d:%02d:%02dZ");
    String IssueInstant;
    IssueInstant.format(timeTemplate.c_str(), 1900 + TM->tm_year, TM->tm_mon, TM->tm_mday, TM->tm_hour, TM->tm_min, TM->tm_sec);

    String content;
    content.format(xmlTmp.c_str(),IssueInstant.getCStr(), ticket.c_str());

    nsresult ret;
    nsCOMPtr<ncIEHTTPClient> httpClient = do_CreateInstance (NC_EHTTP_CLIENT_CONTRACTID, &ret);
    if (NS_FAILED (ret)) {
        THROW_E (EAC_HTTP_SERVER, FAILED_TO_CREATE_EHTTP_CLIENT,
            "Failed to create ehttpclient instance: 0x%x", ret);
    }

    ncEHTTPResponse response;
    httpClient->Post (toSTLString (checkURL), toSTLString(content), "application/xml", 30, response);

    string& resBody = response.body;
    //printMessage2(_T("---------------------------------------------"));
    //printMessage2(_T("response body: %s"), response.body.c_str());

    ptree ptAll;
    // Read the XML config string into the property tree. Catch any exception
    try {
        stringstream ss;
        ss << resBody;
        read_xml(ss, ptAll);
    }
    catch (xml_parser_error &e) {
        THROW_E (EAC_HTTP_SERVER, INVALID_XML_FORMAT,
            LOAD_STRING (_T("IDS_INVALID_XML_FORMAT")));
    }

    // 成功消息
    /*
    <?xml version="1.0" encoding="UTF-8"?>
    <SOAP-ENV:Envelope xmlns:SOAP-ENV="http://schemas.xmlsoap.org/soap/envelope/">
        <SOAP-ENV:Body>
            <saml1p:Response xmlns:saml1p="urn:oasis:names:tc:SAML:1.0:protocol"
                IssueInstant="2014-08-21T13:16:39.789Z" MajorVersion="1"
                MinorVersion="1" Recipient="http://192.168.5.112:8080/client/ssoauth"
                ResponseID="_5c231d5f6fc050eaaedc0a382adba7b1">
                <saml1p:Status>
                    <saml1p:StatusCode Value="saml1p:Success" />
                </saml1p:Status>
                <saml1:Assertion xmlns:saml1="urn:oasis:names:tc:SAML:1.0:assertion"
                    AssertionID="_c38ca75182041b5165d21636bcee1c54" IssueInstant="2014-08-21T13:16:39.789Z"
                    Issuer="localhost" MajorVersion="1" MinorVersion="1">
                    <saml1:Conditions NotBefore="2014-08-21T13:16:39.789Z"
                        NotOnOrAfter="2014-08-21T13:17:09.789Z">
                        <saml1:AudienceRestrictionCondition>
                            <saml1:Audience>http://192.168.5.112:8080/client/ssoauth</saml1:Audience>
                        </saml1:AudienceRestrictionCondition>
                    </saml1:Conditions>
                    <saml1:AuthenticationStatement
                        AuthenticationInstant="2014-08-21T13:16:33.549Z"
                        AuthenticationMethod="urn:oasis:names:tc:SAML:1.0:am:unspecified">
                        <saml1:Subject>
                            <saml1:NameIdentifier>admin</saml1:NameIdentifier>
                            <saml1:SubjectConfirmation>
                                <saml1:ConfirmationMethod>urn:oasis:names:tc:SAML:1.0:cm:artifact</saml1:ConfirmationMethod>
                            </saml1:SubjectConfirmation>
                        </saml1:Subject>
                    </saml1:AuthenticationStatement>
                </saml1:Assertion>
            </saml1p:Response>
        </SOAP-ENV:Body>
    </SOAP-ENV:Envelope>
    */


    // 解析xml消息
    string account;
    try{
        boost::property_tree::ptree ptEnvelope = ptAll.get_child ("SOAP-ENV:Envelope");
        boost::property_tree::ptree ptBody = ptEnvelope.get_child("SOAP-ENV:Body");
        boost::property_tree::ptree ptResponse = ptBody.get_child("saml1p:Response");
        boost::property_tree::ptree ptStatus = ptResponse.get_child("saml1p:Status");
        boost::property_tree::ptree ptStatusCode = ptStatus.get_child("saml1p:StatusCode");

        string statusCode = ptStatusCode.get<string>("<xmlattr>.Value");
        string deniedStr("saml1p:RequestDenied");
        if (statusCode == deniedStr){
            THROW_E(EAC_HTTP_SERVER, EACHTTP_INVALID_THIRD_PARTY_TICKET,
                LOAD_STRING("IDS_INVALID_THIRD_PARTY_TICKET"), ticket. c_str());
        }

        boost::property_tree::ptree ptAssertion = ptResponse.get_child("saml1:Assertion");
        boost::property_tree::ptree ptAuth = ptAssertion.get_child("saml1:AuthenticationStatement");
        boost::property_tree::ptree ptSubject = ptAuth.get_child("saml1:Subject");
        boost::property_tree::ptree ptName = ptSubject.get_child("saml1:NameIdentifier");
        account = ptName.data();
        //printMessage2(_T("account: %s"), account.c_str());
    }
    catch (ptree_error& e) {
        THROW_E (EAC_HTTP_SERVER, INVALID_XML_FORMAT, e.what());
    }

    return account.c_str();
}

bool ncEACAuthHandler::parseVcodeInfo(JSON::Value& requestJson, ncTUserLoginOption &option)
{
    bool hasVcode = false;
    // 获取验证码信息
    if(requestJson.o().find("vcode") != requestJson.o().end())
    {
        JSON::Value vcodeJson = requestJson["vcode"];

        if(vcodeJson.o().find("content") == vcodeJson.o().end())
        {
            THROW_E (EAC_HTTP_SERVER, EACHTTP_INVALID_PARAMETER,
                    LOAD_STRING (_T("IDS_EACHTTP_PARAMETERS_IS_REQUIRED")), "vcode content");
        }
        option.vcode = vcodeJson["content"].s ();

        if(vcodeJson.o().find("id") == vcodeJson.o().end())
        {
            THROW_E (EAC_HTTP_SERVER, EACHTTP_INVALID_PARAMETER,
                LOAD_STRING (_T("IDS_EACHTTP_PARAMETERS_IS_REQUIRED")), "vcode id");
        }
        option.uuid = vcodeJson["id"].s ();

        hasVcode = true;
    }

    option.__isset.vcode = true;
    option.__isset.uuid = true;
    option.__isset.vcode = true;
    option.__isset.uuid = true;
    option.isModify = false;
    option.__isset.isModify = true;

    return hasVcode;
}

bool ncEACAuthHandler::parseDeviceInfo(JSON::Value& requestJson, ncDeviceBaseInfo& info)
{
    //参数存在性检查
    if(requestJson.o().find("device") == requestJson.o().end())
    {
        THROW_E (EAC_HTTP_SERVER, EACHTTP_DEVICE_INFO_INVALID,
            LOAD_STRING (_T("IDS_EACHTTP_DEVICE_INFO_INVALID")));
    }

    JSON::Value deviceJson = requestJson["device"];
    if(deviceJson.o().find("client_type") == deviceJson.o().end())
    {
        THROW_E (EAC_HTTP_SERVER, EACHTTP_DEVICE_INFO_INVALID,
            LOAD_STRING (_T("IDS_EACHTTP_DEVICE_INFO_INVALID")));
    }

    //参数有效性检查
    try {
        //获取客户端类型
        String strClientType = toCFLString(requestJson["device"]["client_type"].s());
        auto iter = _clientIntTypeMap.find(strClientType);
        if(iter != _clientIntTypeMap.end())
        {
            info.clientType = static_cast<ACSClientType>(iter->second);
        }
        else
        {
            THROW_E (EAC_HTTP_SERVER, EACHTTP_DEVICE_INFO_INVALID,
                LOAD_STRING (_T("IDS_INVALID_OS_TYPE")));
        }

        //获取设备名称和设备描述
        if(deviceJson.o().find("name") != deviceJson.o().end())
        {
            if(requestJson["device"]["name"].type() == JSON::STRING)
            {
                info.name = toCFLString(requestJson["device"]["name"].s());
            }
            else
            {
                THROW_E (EAC_HTTP_SERVER, EACHTTP_INVALID_PARAMETER,
                LOAD_STRING (_T("IDS_EACHTTP_PARAMETERS_IS_INVALID")), "device name");
            }
        }

        if(deviceJson.o().find("description") != deviceJson.o().end())
        {
            if(requestJson["device"]["description"].type() == JSON::STRING)
            {
                info.deviceType = toCFLString(requestJson["device"]["description"].s());
            }
            else
            {
                THROW_E (EAC_HTTP_SERVER, EACHTTP_INVALID_PARAMETER,
                    LOAD_STRING (_T("IDS_EACHTTP_PARAMETERS_IS_INVALID")), "device descrption");
            }
        }

        // 解析设备所有mac标识字段
        if(deviceJson.o().find("udids") != deviceJson.o().end())
        {
            if (requestJson["device"]["udids"].type() == JSON::ARRAY) {
                int cnt = requestJson["device"]["udids"].a().size();
                for (int i = 0; i < cnt; i++) {
                    info.udids.push_back(toCFLString(requestJson["device"]["udids"][i].s()));
                }

                if(cnt > 0)
                {
                    info.udid = info.udids[0];
                }
            }
            else
            {
                THROW_E (EAC_HTTP_SERVER, EACHTTP_INVALID_PARAMETER,
                    LOAD_STRING (_T("IDS_EACHTTP_PARAMETERS_IS_INVALID")), "device udids");
            }
        }
    }
    catch(Exception&) {
        THROW_E (EAC_HTTP_SERVER, EACHTTP_DEVICE_INFO_INVALID,
            LOAD_STRING (_T("IDS_EACHTTP_DEVICE_INFO_INVALID")));
    }

    return true;
}

char ncEACAuthHandler::from_hex(char ch)
{
    return isdigit(ch) ? ch - '0' : tolower(ch) - 'a' + 10;

}

char ncEACAuthHandler::to_hex(char code)
{
    static char hex[] = "0123456789abcdef";
    return hex[code & 15];
}

char *ncEACAuthHandler::url_encode(char *str)
{
    char *pstr = str;
    char* buf = (char*)malloc(strlen(str) * 3 + 1);
    char *pbuf = buf;
    while (*pstr) {
        if (isalnum(*pstr) || *pstr == '-' || *pstr == '_' || *pstr == '.' || *pstr == '~')
            *pbuf++ = *pstr;
        else if (*pstr == ' ')
            *pbuf++ = '+';
        else
            *pbuf++ = '%', *pbuf++ = to_hex(*pstr >> 4), *pbuf++ = to_hex(*pstr & 15);
        pstr++;
    }
    *pbuf = '\0';
    return buf;
}

char *ncEACAuthHandler::url_decode(char *str)
{
    char *pstr = str;
    char *buf = (char*)malloc(strlen(str) + 1);
    char *pbuf = buf;
    while (*pstr) {
        if (*pstr == '%') {
            if (pstr[1] && pstr[2]) {
                *pbuf++ = from_hex(pstr[1]) << 4 | from_hex(pstr[2]);
                pstr += 2;
            }
        } else if (*pstr == '+') {
            *pbuf++ = ' ';
        } else {
            *pbuf++ = *pstr;
        }
        pstr++;
    }
    *pbuf = '\0';
    return buf;
}

String ncEACAuthHandler::ThsAuth(JSON::Value& requestJson, const ncTThirdPartyAuthConf& config)
{
    // 获取ticket
    string ticket;
    try {
        ticket = requestJson["params"]["ticket"].s();
    }
    catch (Exception& e) {
        THROW_E(EAC_HTTP_SERVER, EACHTTP_INVALID_THIRD_PARTY_TICKET,
            LOAD_STRING("IDS_INVALID_THIRD_PARTY_TICKET"), ticket.c_str());
    }
    if(ticket == "") {
        THROW_E(EAC_HTTP_SERVER, EACHTTP_INVALID_THIRD_PARTY_TICKET,
            LOAD_STRING("IDS_INVALID_THIRD_PARTY_TICKET"), ticket.c_str());
    }

    string authServer;
    try {
        JSON::Value configJson;
        JSON::Reader::read (configJson, config.config.c_str(), config.config.length());

        authServer = configJson["authServer"].s();
    }
    catch(Exception& e) {
        THROW_E (EAC_HTTP_SERVER, FAILED_TO_PARSE_THIRD_AUTH_CONFIG,
                    LOAD_STRING("IDS_FAILED_TO_PARSE_THIRD_AUTH_CONFIG"));
    }

    // 验证获取用户名
    string checkURL = authServer + "/services/SSOService";

    // post body
    string xmlTmp =
        "<SOAP-ENV:Envelope SOAP-ENV:encodingStyle=\"http://schemas.xmlsoap.org/soap/encoding/\" xmlns:SOAP-ENC=\"http://schemas.xmlsoap.org/soap/encoding/\" xmlns:SOAP-ENV=\"http://schemas.xmlsoap.org/soap/envelope/\" "
            "xmlns:ns0=\"http://schemas.xmlsoap.org/soap/encoding/\" xmlns:ns1=\"http://schemas.xmlsoap.org/soap/envelope/\" xmlns:ns2=\"http://services.sso.platform.ths.com\" xmlns:xsi=\"http://www.w3.org/2001/XMLSchema-instance\"><SOAP-ENV:Header/>"
            "<ns1:Body>"
                "<ns2:decodeToken>"
                    "<token xsi:type=\"ns0:string\">%s</token>"
                "</ns2:decodeToken>"
            "</ns1:Body>"
        "</SOAP-ENV:Envelope>";

    String content;
    content.format(xmlTmp.c_str(), ticket.c_str());

    nsresult ret;
    nsCOMPtr<ncIEHTTPClient> httpClient = do_CreateInstance (NC_EHTTP_CLIENT_CONTRACTID, &ret);
    if (NS_FAILED (ret)) {
        THROW_E (EAC_HTTP_SERVER, FAILED_TO_CREATE_EHTTP_CLIENT,
            "Failed to create ehttpclient instance: 0x%x", ret);
    }

    httpClient->AddHeader("Soapaction","decodeToken");

    ncEHTTPResponse response;
    httpClient->Post (checkURL, toSTLString(content), "application/xml", 30, response);

    string& resBody = response.body;

    ptree ptAll;
    // Read the XML config string into the property tree. Catch any exception
    try {
        stringstream ss;
        ss << resBody;
        read_xml(ss, ptAll);
    }
    catch (xml_parser_error &e) {
        THROW_E (EAC_HTTP_SERVER, INVALID_XML_FORMAT,
            LOAD_STRING (_T("IDS_INVALID_XML_FORMAT")));
    }

    // 解析xml消息
    string account;
    try{

        boost::property_tree::ptree ptEnvelope = ptAll.get_child ("soapenv:Envelope");
        boost::property_tree::ptree ptBody = ptEnvelope.get_child("soapenv:Body");
        boost::property_tree::ptree ptResponse = ptBody.get_child("ns1:decodeTokenResponse");

        boost::property_tree::ptree ptName = ptResponse.get_child("decodeTokenReturn");
        account = ptName.data();
    }
    catch (ptree_error& e) {
        THROW_E (EAC_HTTP_SERVER, INVALID_XML_FORMAT, e.what());
    }

    return account.c_str();
}

String ncEACAuthHandler::LcsoftAuth(JSON::Value& requestJson, const ncTThirdPartyAuthConf& config)
{
    // 获取ticket
    string ticket;
    try {
        ticket = requestJson["params"]["ticket"].s();
    }
    catch (Exception& e) {
        THROW_E(EAC_HTTP_SERVER, EACHTTP_INVALID_THIRD_PARTY_TICKET,
            LOAD_STRING("IDS_INVALID_THIRD_PARTY_TICKET"), ticket.c_str());
    }
    if(ticket == "") {
        THROW_E(EAC_HTTP_SERVER, EACHTTP_INVALID_THIRD_PARTY_TICKET,
            LOAD_STRING("IDS_INVALID_THIRD_PARTY_TICKET"), ticket.c_str());
    }

    // 获取service
    string service;
    try {
        service = requestJson["params"]["service"].s();
    }
    catch (Exception& e) {
        THROW_E(EAC_HTTP_SERVER, EACHTTP_SERVICE_INVALID,
            LOAD_STRING("IDS_EACHTTP_SERVICE_INVALID"), service.c_str());
    }
    if(service == "") {
        THROW_E(EAC_HTTP_SERVER, EACHTTP_SERVICE_INVALID,
            LOAD_STRING("IDS_EACHTTP_SERVICE_INVALID"), service.c_str());
    }

    string authServer;
    try {
        JSON::Value configJson;
        JSON::Reader::read (configJson, config.config.c_str(), config.config.length());

        authServer = configJson["authServer"].s();
    }
    catch(Exception& e) {
        THROW_E (EAC_HTTP_SERVER, FAILED_TO_PARSE_THIRD_AUTH_CONFIG,
                    LOAD_STRING("IDS_FAILED_TO_PARSE_THIRD_AUTH_CONFIG"));
    }

    // 验证获取用户名
    String url;
    url.format ("%s/cas/serviceValidate?service=%s&ticket=%s", authServer.c_str(), service.c_str(), ticket.c_str());

    nsresult ret;
    nsCOMPtr<ncIEHTTPClient> httpClient = do_CreateInstance (NC_EHTTP_CLIENT_CONTRACTID, &ret);
    if (NS_FAILED (ret)) {
        THROW_E (EAC_HTTP_SERVER, FAILED_TO_CREATE_EHTTP_CLIENT,
            "Failed to create ehttpclient instance: 0x%x", ret);
    }

    ncEHTTPResponse response;
    httpClient->Get (toSTLString (url), 30, response);

    string& resBody = response.body;

    ptree ptAll;
    // Read the XML config string into the property tree. Catch any exception
    try {
        stringstream ss;
        ss << resBody;
        read_xml(ss, ptAll);
    }
    catch (xml_parser_error &e) {
        THROW_E (EAC_HTTP_SERVER, INVALID_XML_FORMAT,
            LOAD_STRING (_T("IDS_INVALID_XML_FORMAT")));
    }

    // 解析xml消息
    string account;
    try {
        boost::property_tree::ptree ptResponse = ptAll.get_child ("cas:serviceResponse");
        boost::property_tree::ptree ptResult = ptResponse.get_child("cas:authenticationSuccess");

        boost::property_tree::ptree ptName = ptResult.get_child("cas:user");
        account = ptName.data();
    }
    catch (ptree_error& e) {
        THROW_E (EAC_HTTP_SERVER, INVALID_XML_FORMAT, e.what());
    }

    return account.c_str();
}

String ncEACAuthHandler::WindowsADSSO(JSON::Value& requestJson, const ncTThirdPartyAuthConf& config)
{
    // 获取ticket
    string ticket;
    try {
        ticket = requestJson["params"]["ticket"].s();
    }
    catch (Exception& e) {
        THROW_E(EAC_HTTP_SERVER, EACHTTP_INVALID_THIRD_PARTY_TICKET,
            LOAD_STRING("IDS_INVALID_THIRD_PARTY_TICKET"), ticket.c_str());
    }
    if(ticket == "") {
        THROW_E(EAC_HTTP_SERVER, EACHTTP_INVALID_THIRD_PARTY_TICKET,
            LOAD_STRING("IDS_INVALID_THIRD_PARTY_TICKET"), ticket.c_str());
    }

    String account;
    ParseADSession(toCFLString(ticket), account);

    return account;
}

String ncEACAuthHandler::AnyShareSSO(JSON::Value& requestJson, const ncTThirdPartyAuthConf& config)
{
    String appid;
    try {
        appid = requestJson["params"]["appid"].s ().c_str ();
    }
    catch (Exception& e) {
        THROW_E(EAC_HTTP_SERVER, EACHTTP_INVALID_THIRD_PARTY_TICKET,
            LOAD_STRING("IDS_INVALID_THIRD_PARTY_TICKET"));
    }

    //中国人民公安大学
    if (appid == _T("ppsuc")) {
        String account, time, key;
        try{
            account = requestJson["params"]["un"].s ().c_str ();
            time = requestJson["params"]["time"].s ().c_str ();
            key = requestJson["params"]["verify"].s ().c_str ();
        }
        catch (Exception& e) {
            THROW_E(EAC_HTTP_SERVER, EACHTTP_INVALID_THIRD_PARTY_TICKET,
                LOAD_STRING("IDS_INVALID_THIRD_PARTY_TICKET"));
        }

        vector<String> params;
        params.push_back(time);

        //进行验证
        int accountType = 0;
        String userId = _acsShareMgnt->ExtLogin (appid, account, key, params, accountType);
        if (userId == String::EMPTY) {
            THROW_E(EAC_HTTP_SERVER, EACHTTP_INVALID_THIRD_PARTY_TICKET,
                LOAD_STRING("IDS_INVALID_THIRD_PARTY_TICKET"));
        }
        return account;
    }

    //标准认证
    else {
        String account, key;
        try {
            account = requestJson["params"]["account"].s ().c_str ();
            key = requestJson["params"]["key"].s ().c_str ();
        }
        catch (Exception& e) {
            THROW_E(EAC_HTTP_SERVER, EACHTTP_INVALID_THIRD_PARTY_TICKET,
                LOAD_STRING("IDS_INVALID_THIRD_PARTY_TICKET"));
        }

        // 进行验证
        vector<String> params;
        int accountType = 0;
        String userId = _acsShareMgnt->ExtLogin (appid, account, key, params, accountType);
        if (userId == String::EMPTY) {
            THROW_E(EAC_HTTP_SERVER, EACHTTP_INVALID_THIRD_PARTY_TICKET,
                LOAD_STRING("IDS_INVALID_THIRD_PARTY_TICKET"));
        }
        return account;
    }
}

String ncEACAuthHandler::AnySharePlain(JSON::Value& requestJson, const ncTThirdPartyAuthConf& config)
{
    string account = requestJson["params"]["account"].s ();
    if (account.empty ()) {
        THROW_E (EAC_HTTP_SERVER, EACHTTP_ACCOUNT_INVALID,
            LOAD_STRING (_T("IDS_EACHTTP_ACCOUNT_INVALID")));
    }

    string password = requestJson["params"]["password"].s ();
    if (password.empty ()) {
        THROW_E (EAC_HTTP_SERVER, EACHTTP_PASSWORD_INVALID,
                LOAD_STRING (_T("IDS_EACHTTP_PASSWORD_INVALID")));
    }

    // 进行验证
    string retUserId;
    ncTUserLoginOption option;
    ncEACHttpServerUtil::Usrm_UserLogin (retUserId, account, password, option);

    // 由于 Usrm_UserLogin 支持域帐号免后缀登录
    // 故account和实际匹配的帐号可能不一样，这里要根据userid重新获取一下
    ncACSUserInfo userInfo;
    if(_acsShareMgnt->GetUserInfoById(retUserId.c_str(), userInfo) == false) {
        THROW_E (EAC_HTTP_SERVER, EACHTTP_USER_NOT_EXIST,
                    "User not exists.");
    }

    return userInfo.account;
}

String ncEACAuthHandler::AISHUAuth(JSON::Value& requestJson, const ncTThirdPartyAuthConf& config)
{
    string account = requestJson["params"]["account"].s ();
    if (account.empty ()) {
        THROW_E (EAC_HTTP_SERVER, EACHTTP_ACCOUNT_INVALID,
            LOAD_STRING (_T("IDS_EACHTTP_ACCOUNT_INVALID")));
    }

    return toCFLString(account);
}

String ncEACAuthHandler::AnyShareRSA(JSON::Value& requestJson, const ncTThirdPartyAuthConf& config)
{
    string account = requestJson["params"]["account"].s ();
    if (account.empty ()) {
        THROW_E (EAC_HTTP_SERVER, EACHTTP_ACCOUNT_INVALID,
            LOAD_STRING (_T("IDS_EACHTTP_ACCOUNT_INVALID")));
    }

    string password = requestJson["params"]["password"].s ();
    if (password.empty ()) {
        THROW_E (EAC_HTTP_SERVER, EACHTTP_PASSWORD_INVALID,
            LOAD_STRING (_T("IDS_EACHTTP_PASSWORD_INVALID")));
    }

    // 这里password有可能尾部有\0
    string decodePwd(ncEACHttpServerUtil::Base64Decode(password.c_str()));
    string originPassword (ncEACHttpServerUtil::RSADecrypt(decodePwd));

    // 进行验证
    string retUserId;
    ncTUserLoginOption option;
    ncEACHttpServerUtil::Usrm_UserLogin (retUserId, account, originPassword, option);

    // 由于 Usrm_UserLogin 支持域帐号免后缀登录
    // 故account和实际匹配的帐号可能不一样，这里要根据userid重新获取一下
    ncACSUserInfo userInfo;
    if(_acsShareMgnt->GetUserInfoById(retUserId.c_str(), userInfo) == false) {
        THROW_E (EAC_HTTP_SERVER, EACHTTP_USER_NOT_EXIST,
                    "User not exists.");
    }

    return userInfo.account;
}

String ncEACAuthHandler::clientType2Str(ACSClientType clientType)
{
    switch (clientType) {
        case ACSClientType::UNKNOWN:
            return "Unknown";
        case ACSClientType::IOS:
            return "iOS";
        case ACSClientType::ANDROID:
            return "Android";
        case ACSClientType::WINDOWS_PHONE:
            return "WindowsPhone";
        case ACSClientType::WINDOWS:
            return "Windows";
        case ACSClientType::MAC_OS:
            return "MacOSX";
        case ACSClientType::WEB:
            return "Web";
        case ACSClientType::MOBILE_WEB:
            return "MobileWeb";
        case ACSClientType::NAS:
            return LOAD_STRING(_T("IDS_NAS_GATEWAY")); //获取到 “NAS网关” 国际化资源
        case ACSClientType::CONSOLE_WEB:
            return LOAD_STRING(_T("IDS_CONSOLE_WEB"));
        case ACSClientType::DEPLOY_WEB:
            return LOAD_STRING(_T("IDS_DEPLOY_WEB"));
        case ACSClientType::LINUX:
            return "Linux";
        default:
            return "Unknown";
    }
}

void ncEACAuthHandler::onUserLogin(brpc::Controller* cntl, const String& retUserId, const String& ip, bool hasDeviceInfo, ncDeviceBaseInfo& baseInfo)
{
    // 获取用户的设备信息
    vector<ncDeviceInfo> deviceInfos;
    _acsDeviceManager->GetDevicesByUserIdAndUdid (retUserId, String::EMPTY, deviceInfos, 0, -1, true);

    vector<String> LimitUdids;
    vector<String> bindUdids;

    for (size_t i = 0; i < deviceInfos.size(); ++i) {

        // 设备id比较时不区分大小写
        deviceInfos[i].baseInfo.udid.toUpper ();

        // 获取用户禁用设备信息
        if (deviceInfos[i].disableFlag == 1) {
            LimitUdids.push_back (deviceInfos[i].baseInfo.udid);
        }

        // 获取用户绑定设备信息
        if (deviceInfos[i].bindFlag == 1) {
            bindUdids.push_back (deviceInfos[i].baseInfo.udid);
        }
    }

    // 兴业银行处理 如果配置开启，则所有客户端登录时都执行这个逻辑，包括windows mac_os web和linux
    String tempValue;
    bool AllForceLogOff = false;
    bool result = _acsShareMgnt->GetCustomConfigOfString("all_client_force_log_off", tempValue);
    if (result) {
        if (tempValue.compare("1") == 0) {
            AllForceLogOff = true;
        }
    } 

    // 登录请求中带了设备信息
    if(hasDeviceInfo && !baseInfo.name.isEmpty() && !baseInfo.deviceType.isEmpty() && !baseInfo.udids.empty()) {
        // 1.设备禁用检查，udids中有一个被禁用则禁用
        for (size_t i = 0; i < baseInfo.udids.size(); ++i) {
            baseInfo.udids[i].toUpper ();
            auto iter = find(LimitUdids.begin (), LimitUdids.end (), baseInfo.udids[i]);
            if (iter != LimitUdids.end()) {
                THROW_E (EAC_HTTP_SERVER, EACHTTP_DEVICE_DISABLED,
                    LOAD_STRING("IDS_DEVICE_DISABLED"));
            }
        }

        // 2. 如果绑定了设备，则进行设备绑定检查。
        bool isBindDevice = false;
        if (bindUdids.size() > 0) {
            // 该设备不在绑定列表中（设备id不区分大小写），则不允许登录
            for (size_t i = 0; i < baseInfo.udids.size(); ++i) {
                auto iter = find(bindUdids.begin (), bindUdids.end (), baseInfo.udids[i]);
                if (iter != bindUdids.end ()) {
                    // 更新请求头中的MAC地址信息，记录对应日志
                    updateHTTPHeaderMACAddress (cntl, baseInfo.udids[i]);
                    // 更新字段，在token记录对应信息
                    baseInfo.udid = baseInfo.udids[i];
                    isBindDevice = true;
                    break;
                }
            }
        }

        // 3. 如果有“所有用户”绑定的设备，进行设备绑定检查
        if (!isBindDevice) {
            for (auto iter = baseInfo.udids.begin(); iter != baseInfo.udids.end(); iter++) {
                String tmpUdid = *iter;
                ncDeviceInfo info;
                bool infoExsit = _acsDeviceManager->GetDeviceByUDID(toCFLString(g_ShareMgnt_constants.NCT_ALL_USER_GROUP), tmpUdid, info);
                if(infoExsit && info.bindFlag) {
                    // 更新请求头中的MAC地址信息，记录对应日志
                    updateHTTPHeaderMACAddress (cntl, tmpUdid);
                    // 更新字段，在token记录对应信息
                    baseInfo.udid = tmpUdid;
                    isBindDevice = true;
                    break;
                }
            }
        }

        // 如果没有绑定的设备
        bool allUserBindDeviceExist = _acsDeviceManager->UserHasDeviceInfo (toCFLString(g_ShareMgnt_constants.NCT_ALL_USER_GROUP), true);
        if (bindUdids.empty() && !allUserBindDeviceExist) {
            isBindDevice = true;
        }

        if(!isBindDevice) {
            THROW_E (EAC_HTTP_SERVER, EACHTTP_DEVICE_NOT_BINDED,
                LOAD_STRING("IDS_DEVICE_NOT_BINDED"));
        }

        bool ret = _acsPolicyManager->CheckIp (retUserId, ip);
        if(!ret) {
            THROW_E (EAC_HTTP_SERVER, EACHTTP_LOGIN_IP_IS_RESTRICTED,
                LOAD_STRING("IDS_LOGIN_IP_IS_RESTRICTED"));
        }

        // 登录成功，记录设备信息
        baseInfo.lastLoginIp = ip;
        baseInfo.lastLoginTime = BusinessDate::getCurrentTime ();
        _acsDeviceManager->RecordDevice(retUserId, baseInfo);

        if(AllForceLogOff) {
            _acsDeviceManager->AllForceLogOff(retUserId);
        } else if(baseInfo.clientType == ACSClientType::WINDOWS && ncEACHttpServerUtil::GetLoginStrategyStatus()) {
            _acsDeviceManager->ForceLogOff(retUserId);
        }
    }
    else {
        // 如果绑定了其它设备，则不允许登录
        bool allUserBindDeviceExist = _acsDeviceManager->UserHasDeviceInfo (toCFLString(g_ShareMgnt_constants.NCT_ALL_USER_GROUP), true);
        if(bindUdids.size() > 0 || allUserBindDeviceExist) {
            THROW_E (EAC_HTTP_SERVER, EACHTTP_DEVICE_NOT_BINDED,
                LOAD_STRING("IDS_DEVICE_NOT_BINDED"));
        }

        bool ret = _acsPolicyManager->CheckIp (retUserId, ip);
        if(!ret) {
            THROW_E (EAC_HTTP_SERVER, EACHTTP_LOGIN_IP_IS_RESTRICTED,
                LOAD_STRING("IDS_LOGIN_IP_IS_RESTRICTED"));
        }

        if(AllForceLogOff) {
            _acsDeviceManager->AllForceLogOff(retUserId);
        }
    }

    // 如果存在设备类型，则需检查设备类型是否被禁用
    if (hasDeviceInfo) {
        bool ret = _policyEngine->Audit_ClientRestriction(_clientStringTypeMap[baseInfo.clientType].getCStr());
        if(ret) {
            THROW_E (EAC_HTTP_SERVER, EACHTTP_LOGIN_OSTYPE_IS_FORBID,
                LOAD_STRING("IDS_LOGIN_OSTYPE_IS_FORBID"));
        }
    }
}

// protected
void ncEACAuthHandler::updateHTTPHeaderMACAddress(brpc::Controller* cntl, const String& macAddr)
{
    if (macAddr.isEmpty ()) {
        return;
    }

    if (cntl->http_request ().GetHeader ("X-Request-MAC") != NULL) {
        cntl->http_request().RemoveHeader("X-Request-MAC");
    }

    cntl->http_request().SetHeader("X-Request-MAC", macAddr.getCStr ());
}

// protected
void ncEACAuthHandler::updateHTTPHeaderIP(brpc::Controller* cntl, const String& ip)
{
    if (ip.isEmpty ()) {
        return;
    }

    if (cntl->http_request ().GetHeader ("X-Forwarded-For") != NULL) {
        cntl->http_request().RemoveHeader("X-Forwarded-For");
    }

    cntl->http_request().SetHeader("X-Forwarded-For", ip.getCStr ());
}

// protected
void ncEACAuthHandler::CheckUninstallPwd (brpc::Controller* cntl, const String& fakeUserId)
{
    NC_EAC_HTTP_SERVER_TRACE (_T("this: %p, cntl:%p begin"), this, cntl);

    // 获取http请求的content
    string bodyBuffer = cntl->request_attachment ().to_string ();
    JSON::Value requestJson;
    try {
        JSON::Reader::read (requestJson, bodyBuffer.c_str (), bodyBuffer.size ());
    }
    catch (Exception& e) {
        THROW_E (EAC_HTTP_SERVER, EACHTTP_INVALID_REQUEST_JSON_FORMAT,
            LOAD_STRING (_T("IDS_EACHTTP_INVALID_JSON")));
    }

    String uninstallPwd = requestJson["uninstallpwd"].s ().c_str ();;
    ncEACHttpServerUtil::CheckUninstallPwd(uninstallPwd);

    JSON::Value replyJson;
    replyJson["result"] = "ok";

    string body;
    JSON::Writer::write (replyJson.o (), body);
    ncHttpSendReply (cntl, brpc::HTTP_STATUS_OK, "ok", body);
}

// protected
void ncEACAuthHandler::CheckExitPwd (brpc::Controller* cntl, const String& fakeUserId)
{
    NC_EAC_HTTP_SERVER_TRACE (_T("this: %p, cntl:%p begin"), this, cntl);

    // 获取http请求的content
    string bodyBuffer = cntl->request_attachment ().to_string ();
    JSON::Value requestJson;
    try {
        JSON::Reader::read (requestJson, bodyBuffer.c_str (), bodyBuffer.size ());
    }
    catch (Exception& e) {
        THROW_E (EAC_HTTP_SERVER, EACHTTP_INVALID_REQUEST_JSON_FORMAT,
            LOAD_STRING (_T("IDS_EACHTTP_INVALID_JSON")));
    }

    String uninstallPwd = requestJson["exitpwd"].s ().c_str ();;
    ncEACHttpServerUtil::CheckExitPwd(uninstallPwd);

    JSON::Value replyJson;
    replyJson["result"] = "ok";

    string body;
    JSON::Writer::write (replyJson.o (), body);
    ncHttpSendReply (cntl, brpc::HTTP_STATUS_OK, "ok", body);
}

// protected
void ncEACAuthHandler::Login (brpc::Controller* cntl, const String& fakeUserId)
{
    NC_EAC_HTTP_SERVER_TRACE (_T("this: %p, cntl:%p, userId:%s begin"), this, cntl, fakeUserId.getCStr ());

    // 获取http请求的content
    string bodyBuffer = cntl->request_attachment ().to_string ();

    JSON::Value requestJson;
    try {
        JSON::Reader::read (requestJson, bodyBuffer.c_str (), bodyBuffer.size ());
    }
    catch (Exception& e) {
        THROW_E (EAC_HTTP_SERVER, EACHTTP_INVALID_REQUEST_JSON_FORMAT,
            LOAD_STRING (_T("IDS_EACHTTP_INVALID_JSON")));
    }

    // 登录类型
    string grantType = requestJson["grant_type"].s ();

    // token类型
    string tokenType = requestJson["token_type"].s ();
    int expiredTime = getExpiredTime(tokenType);

    // 获取登录的设备信息
    ncDeviceBaseInfo baseInfo;
    baseInfo.clientType = ACSClientType::UNKNOWN;
    bool hasDeviceInfo = parseDeviceInfo(requestJson, baseInfo);

    // 将clienttype信息加入params中
    if (hasDeviceInfo) {
        requestJson["params"]["clienttype"] = requestJson["deviceinfo"]["ostype"].i();
    }

    // 预先分析出帐号
    String account;
    parseAccount(requestJson, account);

    // 获取第三方配置信息
    ncACSUserInfo userInfo;
    ncTThirdPartyAuthConf config = ncEACHttpServerUtil::Usrm_GetThirdPartyAuth ();
    map<String, ncThirdAuthFunc>::iterator iter = _thirdAuthFuncs.find (grantType.c_str());
    int accountType = 0;
    try {
        if(iter == _thirdAuthFuncs.end()) {
            // 转给ShareMgnt进行认证
            string params;
            requestJson["deviceinfo"]["X-Real-IP"] = ncEACHttpServerUtil::GetForwardedIp(cntl).getCStr();
            JSON::Writer::write (requestJson.o (), params);

            account = ncEACHttpServerUtil::Usrm_ValidateThirdParty(params);
        }
        else {
            // 在eacp内部认证
            ncThirdAuthFunc func = iter->second;
            account = (this->*func) (requestJson, config);
        }

        // 用户是否存在
        if(_acsShareMgnt->GetUserInfoByAccount(account, userInfo, accountType) == false) {
            THROW_E (EAC_HTTP_SERVER, EACHTTP_NOT_IMPORT_TO_ANYSHARE,
                LOAD_STRING (_T("IDS_EACHTTP_NOT_IMPORT_TO_ANYSHARE")));
        }

        // 系统管理员不允许登录
        if(_adminIds.count(userInfo.id.getCStr())) {
            THROW_E(EAC_HTTP_SERVER, EACHTTP_FORBIDDEN_LOGIN,
                "%s are not allowed to log in.", account.getCStr());;
        }

        // 用户是否启用
        if(userInfo.enableStatus == false) {
            THROW_E (EAC_HTTP_SERVER, EACHTTP_USER_DISABLED,
                LOAD_STRING (_T("IDS_USER_DISABLED")));
        }
    }
    catch(Exception& e) {
        logFailedLoginEvent(cntl, account, e, hasDeviceInfo, baseInfo);
        throw;
    }

    // 绑定设备管理
    onUserLogin(cntl, userInfo.id, ncEACHttpServerUtil::GetForwardedIp(cntl), hasDeviceInfo, baseInfo);

#ifndef __UT__
    // 更新用户最近请求时间记录
    _acsShareMgnt->UpdateUserLastRequestTime(userInfo.id);
#endif

    // 回复
    JSON::Value replyJson;
    replyJson["userid"] = userInfo.id.getCStr ();
    replyJson["expires_in"] = expiredTime;

    string body;
    JSON::Writer::write (replyJson.o (), body);
    ncHttpSendReply (cntl, brpc::HTTP_STATUS_OK, "ok", body);

    String msg;
    if(hasDeviceInfo) {
        msg.format (ncEACHttpServerLoader, _T("IDS_LOGIN_SUCCESS"), clientType2Str(baseInfo.clientType).getCStr());
    }
    else {
        msg.format (ncEACHttpServerLoader, _T("IDS_LOGIN_SUCCESS"), "Unknown");
    }
    ncEACHttpServerUtil::Log (cntl, userInfo.id, ncTokenVisitorType::REALNAME, ncTLogType::NCT_LT_LOGIN, ncTLogLevel::NCT_LL_INFO,
        ncTLoginType::NCT_CLT_LOGIN_IN, msg, "");

    NC_EAC_HTTP_SERVER_TRACE (_T("this: %p, cntl:%p, userId:%s end"), this, cntl, fakeUserId.getCStr ());
}

// protected
void ncEACAuthHandler::GetVcode (brpc::Controller* cntl, const String& fakeUserId)
{
    NC_EAC_HTTP_SERVER_TRACE (_T("this: %p, cntl:%p begin"), this, cntl);

    // 获取http请求的content
    string bodyBuffer = cntl->request_attachment ().to_string ();
    JSON::Value requestJson;
    try {
        JSON::Reader::read (requestJson, bodyBuffer.c_str (), bodyBuffer.size ());
    }
    catch (Exception& e) {
        THROW_E (EAC_HTTP_SERVER, EACHTTP_INVALID_REQUEST_JSON_FORMAT,
                    LOAD_STRING (_T("IDS_EACHTTP_INVALID_JSON")));
    }

    // 获取验证码
    string uuid = requestJson["uuid"];
    ncTVcodeCreateInfo vcodeInfo;
    ncEACHttpServerUtil::Usrm_CreateVcodeInfo(vcodeInfo, uuid, ncTVcodeType::IMAGE_VCODE);

    // 回复
    JSON::Value replyJson;
    replyJson["uuid"] = vcodeInfo.uuid;
    replyJson["vcode"] = vcodeInfo.vcode;

    string body;
    JSON::Writer::write (replyJson.o (), body);
    ncHttpSendReply (cntl, brpc::HTTP_STATUS_OK, "ok", body);

    NC_EAC_HTTP_SERVER_TRACE (_T("this: %p, cntl:%p end"), this, cntl);
}

int ncEACAuthHandler::getExpiredTime(const string& tokenType)
{
    int expiredTime = 3600;
    if(tokenType == "middle-lived") {
        expiredTime = 3 * 24 * 3600;
    }
    else if(tokenType == "long-lived") {
        expiredTime = 30 * 24 * 3600;
    }

    return expiredTime;
}

int ncEACAuthHandler::getDefaultPermExpiredDays()
{

    bool indefinite_perm = _acsConfManager->GetConfig("oem_indefinite_perm").compareIgnoreCase("true") == 0;
    int defaultPermExpiredDays = Int::getValue(_acsConfManager->GetConfig("oem_default_perm_expired_days"));
    if (defaultPermExpiredDays == 0 || defaultPermExpiredDays < -1 || defaultPermExpiredDays > 3650) {
        THROW_E (EAC_HTTP_SERVER, EACHTTP_INVALID_OEM_CONFIG, "Invalid oem_default_perm_expired_days: %d, must in [-1, [1, 3650]].", defaultPermExpiredDays);
    }
    if (defaultPermExpiredDays == -1 && !indefinite_perm) {
        THROW_E (EAC_HTTP_SERVER, EACHTTP_INVALID_OEM_CONFIG, "oem_default_perm_expired_days is -1 while oem_indefinite_perm is false");
    }

    return defaultPermExpiredDays;
}

JSON::Value ncEACAuthHandler::getCsfLevelsConfig()
{
    map<string, int32_t> csflevels;
    ncEACHttpServerUtil::GetCSFLevels(csflevels);
    JSON::Value _csflevels;
    for (auto iter = csflevels.begin ();iter != csflevels.end ();iter++)
    {
        _csflevels[iter->first] = iter->second;
    }
    return _csflevels;
}

void ncEACAuthHandler::authenicaitonFailedLoginEvent(brpc::Controller* cntl,
                                            const String& account,
                                            const Exception& e,
                                            bool hasDeviceInfo,
                                            const ncDeviceBaseInfo& baseInfo)
{
    // 用户名或密码错误，记录日志
    ncACSUserInfo userInfo;
    int accountType = 0;
    _acsShareMgnt->GetUserInfoByAccount(account, userInfo, accountType);

    if(e.getErrorId() == EACHTTP_INVALID_USER_OR_PASSWORD) {
        ncEACHttpServerUtil::Log (cntl, userInfo.id, ncTokenVisitorType::REALNAME, ncTLogType::NCT_LT_LOGIN, ncTLogLevel::NCT_LL_WARN,
            ncTLoginType::NCT_CLT_AUTHENICATION, LOAD_STRING("IDS_AUTHENICATION_FAILED"), LOAD_STRING("IDS_EACHTTP_INVALID_ACCOUNT_OR_PASSWORD"));
    }
    // 用户被禁用，记录日志
    else if(e.getErrorId() == EACHTTP_USER_DISABLED) {
        ncEACHttpServerUtil::Log (cntl, userInfo.id, ncTokenVisitorType::REALNAME, ncTLogType::NCT_LT_LOGIN, ncTLogLevel::NCT_LL_WARN,
            ncTLoginType::NCT_CLT_AUTHENICATION, LOAD_STRING("IDS_AUTHENICATION_FAILED"), LOAD_STRING("IDS_USER_DISABLED"));
    }

    // 如果是第三次密码输入错误，需要记录日志
    else if(e.getErrorId() == EACHTTP_PWD_FAILED_THIRDLY) {
        String exMsg;
        exMsg.format (ncEACHttpServerLoader, _T("IDS_INPUT_WRONG_PASSWORD_FOR_THREE_TIMES"), _acsConfManager->GetPasswordErrCnt());
        ncEACHttpServerUtil::Log (cntl, userInfo.id, ncTokenVisitorType::REALNAME, ncTLogType::NCT_LT_LOGIN, ncTLogLevel::NCT_LL_WARN,
            ncTLoginType::NCT_CLT_AUTHENICATION, LOAD_STRING("IDS_AUTHENICATION_FAILED"), exMsg);
    }
}

void ncEACAuthHandler::logFailedLoginEvent(brpc::Controller* cntl,
                                            const String& account,
                                            const Exception& e,
                                            bool hasDeviceInfo,
                                            const ncDeviceBaseInfo& baseInfo)
{
    // 用户名或密码错误，记录日志
    ncACSUserInfo userInfo;
    int accountType = 0;
    _acsShareMgnt->GetUserInfoByAccount(account, userInfo, accountType);

    if(e.getErrorId() == EACHTTP_INVALID_USER_OR_PASSWORD) {
        String msg;
        if(hasDeviceInfo) {
            msg.format (ncEACHttpServerLoader, _T("IDS_LOGIN_FAILED"), clientType2Str(baseInfo.clientType).getCStr());
        }
        else {
            msg.format (ncEACHttpServerLoader, _T("IDS_LOGIN_FAILED"), "Unknown");
        }

        ncEACHttpServerUtil::Log (cntl, userInfo.id, ncTokenVisitorType::REALNAME, ncTLogType::NCT_LT_LOGIN, ncTLogLevel::NCT_LL_WARN,
            ncTLoginType::NCT_CLT_LOGIN_IN, msg, LOAD_STRING("IDS_EACHTTP_INVALID_ACCOUNT_OR_PASSWORD"));
    }
    // 用户被禁用，记录日志
    else if(e.getErrorId() == EACHTTP_USER_DISABLED) {
        String msg;
        if(hasDeviceInfo) {
            msg.format (ncEACHttpServerLoader, _T("IDS_LOGIN_FAILED"), clientType2Str(baseInfo.clientType).getCStr());
        }
        else {
            msg.format (ncEACHttpServerLoader, _T("IDS_LOGIN_FAILED"), "Unknown");
        }

        ncEACHttpServerUtil::Log (cntl, userInfo.id, ncTokenVisitorType::REALNAME, ncTLogType::NCT_LT_LOGIN, ncTLogLevel::NCT_LL_WARN,
            ncTLoginType::NCT_CLT_LOGIN_IN, msg, LOAD_STRING("IDS_USER_DISABLED"));
    }

    // 如果是第三次密码输入错误，需要记录日志
    else if(e.getErrorId() == EACHTTP_PWD_FAILED_THIRDLY) {
        String msg;
        if(hasDeviceInfo) {
            msg.format (ncEACHttpServerLoader, _T("IDS_LOGIN_FAILED_AND_LOCKED"), clientType2Str(baseInfo.clientType).getCStr());
        }
        else {
            msg.format (ncEACHttpServerLoader, _T("IDS_LOGIN_FAILED_AND_LOCKED"), "Unknown");
        }

        String exMsg;
        exMsg.format (ncEACHttpServerLoader, _T("IDS_INPUT_WRONG_PASSWORD_FOR_THREE_TIMES"), _acsConfManager->GetPasswordErrCnt());
        ncEACHttpServerUtil::Log (cntl, userInfo.id, ncTokenVisitorType::REALNAME, ncTLogType::NCT_LT_LOGIN, ncTLogLevel::NCT_LL_WARN,
            ncTLoginType::NCT_CLT_LOGIN_IN, msg, exMsg);
    }
}

void ncEACAuthHandler::logFailedConsoleLoginEvent(brpc::Controller* cntl,
                                            const string& osType,
                                            const EHttpDetailException& e)
{
    // 获取详细错误信息
    JSON::Value errDetail = e.getDetail();

    // 用户名或密码错误，记录日志
    if(e.getErrorId() == EACHTTP_INVALID_USER_OR_PASSWORD) {
        String strUserID = toCFLString(errDetail["id"]);
        ncEACHttpServerUtil::Log (cntl, strUserID, ncTokenVisitorType::REALNAME, ncTLogType::NCT_LT_LOGIN, ncTLogLevel::NCT_LL_WARN,
            ncTLoginType::NCT_CLT_AUTHENICATION, LOAD_STRING("IDS_AUTHENICATION_FAILED"), LOAD_STRING("IDS_INCORRECT_USERNAME_OR_PASSWORD"));
    }

    // 如果是第三次密码输入错误，需要记录日志
    else if(e.getErrorId() == EACHTTP_PWD_FAILED_THIRDLY) {
        String strUserID = toCFLString(errDetail["id"]);
        String exMsg;
        exMsg.format (ncEACHttpServerLoader, _T("IDS_WRONG_PASSWORD_MANY_TIMES"), _acsConfManager->GetPasswordErrCnt());
        ncEACHttpServerUtil::Log (cntl, strUserID, ncTokenVisitorType::REALNAME, ncTLogType::NCT_LT_LOGIN, ncTLogLevel::NCT_LL_WARN,
            ncTLoginType::NCT_CLT_AUTHENICATION, LOAD_STRING("IDS_AUTHENICATION_FAILED"), exMsg);
    }
}

// 分析account
void ncEACAuthHandler::parseAccount(JSON::Value& requestJson, String& account)
{
    bool hasDeviceInfo = false;
    if(requestJson.o().find("params") != requestJson.o().end()) {
        JSON::Value& params = requestJson["params"];
        if(params.o().find("account") != params.o().end()) {
            account = requestJson["params"]["account"].s().c_str();
        }
    }
}

int ncEACAuthHandler::compareVersion(const String& version1, const String& version2)
{
    vector<String> vectorVersion1;
    vector<String> vectorVersion2;
    version1.split('.', vectorVersion1);
    version2.split('.', vectorVersion2);

    if (vectorVersion1.size() != 4 || vectorVersion2.size() != 4) {
        THROW_E(EAC_HTTP_SERVER, EACHTTP_ARGUMENT_INVALID,LOAD_STRING("IDS_VERSION_INVALID"));
    }

    Regex versionRegex (_T("^[0-9]{1,8}$"));
    for (int i = 0; i < 4; ++i) {
        if (!versionRegex.match (vectorVersion1[i]) || !versionRegex.match (vectorVersion2[i])) {
            THROW_E(EAC_HTTP_SERVER, EACHTTP_ARGUMENT_INVALID,LOAD_STRING("IDS_VERSION_INVALID"));
        }
    }

    int result = 0;
    for (int i = 0; i < 4; ++i) {
        result = atoi(vectorVersion1[i].getCStr()) - atoi(vectorVersion2[i].getCStr());
        if (result != 0) {
            break;
        }
    }

    return result;
}

void ncEACAuthHandler::checkClientVersion(ACSClientType clientType, const String& version)
{
    String limitVersion;

    if (clientType == ACSClientType::IOS) {
        limitVersion = _acsConfManager->GetConfig("ios_limit_version");
    }
    else if (clientType == ACSClientType::ANDROID) {
        limitVersion = _acsConfManager->GetConfig("andriod_limit_version");
    }
    else {
        return;
    }

    if (!limitVersion.isEmpty() && compareVersion(version, limitVersion) < 0) {
        THROW_E (EAC_HTTP_SERVER, EACHTTP_CLIENT_LOW_VERSION, LOAD_STRING(_T("IDS_CLIENT_LOW_VERSION")));
    }
}

void ncEACAuthHandler::SendVcode (brpc::Controller* cntl, const String& fakeUserId)
{
    NC_EAC_HTTP_SERVER_TRACE (_T("this: %p, cntl:%p begin"), this, cntl);
    string bodyBuffer = cntl->request_attachment ().to_string();
    // 返回uuid
    string retuuid;

    JSON::Value requestJson;
    try {
        JSON::Reader::read (requestJson, bodyBuffer.c_str (), bodyBuffer.size ());
    }
    catch (Exception& e) {
        THROW_E (EAC_HTTP_SERVER, EACHTTP_INVALID_REQUEST_JSON_FORMAT,
                    LOAD_STRING (_T("IDS_EACHTTP_INVALID_JSON")));
    }

    string email = requestJson["emailaddress"].s();
    string telnumber = requestJson["telnumber"].s();
    string uuidIn = requestJson["uuid"].s();

    // 邮箱或手机都为空或都不为空，抛参数异常
    if((email.empty() && telnumber.empty()) || (!email.empty() && !telnumber.empty()))
    {
        THROW_E (EAC_HTTP_SERVER, EACHTTP_INVALID_PARAM_VALUE,
            LOAD_STRING (_T("IDS_EACHTTP_ERR_MSG_AGR_ERR")));
    }

    // 发送验证码
    retuuid = sendEmailAndTelVcode(email, telnumber, uuidIn);

    // 回复
    JSON::Value replyJson;
    replyJson["uuid"] = retuuid;
    string body;
    JSON::Writer::write (replyJson.o (), body);
    ncHttpSendReply (cntl, brpc::HTTP_STATUS_OK, "ok", body);

    NC_EAC_HTTP_SERVER_TRACE (_T("this: %p, cntl:%p end"), this, cntl);
}

void ncEACAuthHandler::SendPwdRetrevalVcode (brpc::Controller* cntl, const String& fakeUserId)
{
    NC_EAC_HTTP_SERVER_TRACE (_T("this: %p, cntl:%p begin"), this, cntl);
    string bodyBuffer = cntl->request_attachment ().to_string();
    // 返回uuid
    string retuuid;

    JSON::Value requestJson;
    try {
        JSON::Reader::read (requestJson, bodyBuffer.c_str (), bodyBuffer.size ());
    }
    catch (Exception& e) {
        THROW_E (EAC_HTTP_SERVER, EACHTTP_INVALID_REQUEST_JSON_FORMAT,
                    LOAD_STRING (_T("IDS_EACHTTP_INVALID_JSON")));
    }

    // 参数获取
    string account = requestJson["account"].s();
    string type = requestJson["type"].s();
    if (account == "" || (type != "telephone" && type != "email"))
    {
        THROW_E (EAC_HTTP_SERVER, EACHTTP_INVALID_PARAM_VALUE,
                    LOAD_STRING (_T("IDS_EACHTTP_ERR_MSG_AGR_ERR")));
    }

    // 获取用户信息，如果用户不存在，报错参数错误
    ncACSUserInfo userInfo;
    int accountType;
    bool ret = _acsShareMgnt->GetUserInfoByAccount(toCFLString(account), userInfo, accountType);
    if (ret == false)
    {
        THROW_E (EAC_HTTP_SERVER, EACHTTP_USER_NOT_EXIST,
                    LOAD_STRING (_T("IDS_EACHTTP_USER_NOT_EXIST")));
    }

    string email;
    string telnumber;
    if (type == "telephone") {
        if (userInfo.telNumber == "") {
            THROW_E (EAC_HTTP_SERVER, EACHTTP_PHONE_HAS_NOT_BEEN_BOUND,
                LOAD_STRING (_T("IDS_EACHTTP_PHONE_HAS_NOT_BEEN_BOUND")));
        }
        telnumber = toSTLString(userInfo.telNumber);
    } else if(type == "email") {
        if (userInfo.email == "") {
             THROW_E (EAC_HTTP_SERVER, EACHTTP_EMAIL_ADDRESS_HAS_NOT_BEEN_BOUND,
                LOAD_STRING (_T("IDS_EACHTTP_EMAIL_ADDRESS_HAS_NOT_BEEN_BOUND")));
        }
        email = toSTLString(userInfo.email);
    }

    // 发送验证码
    retuuid = sendEmailAndTelVcode(email, telnumber, "");

    // 回复
    JSON::Value replyJson;
    replyJson["uuid"] = retuuid;
    string body;
    JSON::Writer::write (replyJson.o (), body);
    ncHttpSendReply (cntl, brpc::HTTP_STATUS_OK, "ok", body);

    NC_EAC_HTTP_SERVER_TRACE (_T("this: %p, cntl:%p end"), this, cntl);
}

string ncEACAuthHandler::sendEmailAndTelVcode (const string& email, const string& telnumber, const string& uuidIn)
{
    string retuuid;

    ncACSUserInfo userInfo;
    if(!telnumber.empty()){
        Regex versionRegex (_T("^[0-9]{1,20}$"));
        if (!versionRegex.match(toCFLString(telnumber)))
            THROW_E(EAC_HTTP_SERVER, EACHTTP_INVALID_TEL_NUMBER,LOAD_STRING("IDS_EACHTTP_INVALID_TEL_NUMBER"));
        // 手机号是否绑定
        if(_acsShareMgnt->GetUserInfoByTelNumber(toCFLString(telnumber), userInfo) == false)
            THROW_E (EAC_HTTP_SERVER, EACHTTP_PHONE_HAS_NOT_BEEN_BOUND,
                LOAD_STRING (_T("IDS_EACHTTP_PHONE_HAS_NOT_BEEN_BOUND")));
    }
    else if(!email.empty()){
        // 邮箱合法性校验
        Regex versionRegex (_T("^[a-zA-Z0-9_\\.\\-]+@[a-zA-Z0-9\\-_]+(\\.[a-zA-Z0-9\\-_]+)+$"));
        if (!versionRegex.match (toCFLString(email)) || email.length() > 100)
            THROW_E(EAC_HTTP_SERVER, EACHTTP_INVALID_EMAIL,LOAD_STRING("IDS_EACHTTP_INVALID_EMAIL"));
        // 邮箱是否绑定
        if(_acsShareMgnt->GetUserInfoByEmail(toCFLString(email), userInfo) == false)
            THROW_E (EAC_HTTP_SERVER, EACHTTP_EMAIL_ADDRESS_HAS_NOT_BEEN_BOUND,
                LOAD_STRING (_T("IDS_EACHTTP_EMAIL_ADDRESS_HAS_NOT_BEEN_BOUND")));
    }

    // 获取验证码发送服务器开关
    bool send_vcode_by_email_status = false;
    bool send_vcode_by_sms_status = false;
    try {
        String tmpString = _acsConfManager->GetVcodeServerStatus();
        string vcode_server_status = toSTLString(tmpString);

        JSON::Value sendVcodeTypeStatusJson;
        JSON::Reader::read (sendVcodeTypeStatusJson, vcode_server_status.c_str (), vcode_server_status.size ());
        send_vcode_by_email_status = sendVcodeTypeStatusJson["send_vcode_by_email"];
        send_vcode_by_sms_status = sendVcodeTypeStatusJson["send_vcode_by_sms"];
    }catch (TException & e) {
        THROW_E (EAC_HTTP_SERVER, EACHTTP_SERVER_UNKNOWN_ERR, e.what ());
    }

    // 判断用户是否被禁用
    if (!userInfo.enableStatus)
    {
        THROW_E (EAC_HTTP_SERVER, EACHTTP_USER_DISABLED, "Invalid user.");
    }
    // 判断用户是否非本地用户
    if (userInfo.authType != 1)
    {
        THROW_E (EAC_HTTP_SERVER, EACHTTP_NOT_LOCAL_USER, "External user cannot modify password.");
    }
    // 判断用户是否为管控用户
    if (userInfo.pwdControl == 1)
    {
        THROW_E (EAC_HTTP_SERVER, EACHTTP_CANNOT_MODIFY_CONTROL_PASSWORD, "The user has been set Change Password Not Allowed.");
    }

    // 获取验证码
    ncTVcodeCreateInfo vcodeInfo;
    ncEACHttpServerUtil::Usrm_CreateVcodeInfo(vcodeInfo, uuidIn, ncTVcodeType::NUM_VCODE);
    retuuid = vcodeInfo.uuid;

    if(!telnumber.empty() && send_vcode_by_sms_status){
        // 发送短信验证码
        std::vector<std::shared_ptr<acsMessage>> msgList;
        std::shared_ptr<acsMessage> msgptr(new acsMessage());
        // channel
        msgptr->channel = RESET_PASSWORD_VERIFICATION_CODE_CHANNEL;
        // content
        JSON::Object payloadObj;
        payloadObj["code"] = vcodeInfo.vcode.c_str();
        std::string payload;
        JSON::Writer::write (payloadObj, payload);
        msgptr->content = std::move (toCFLString (payload));
        // receivers
        vector<messageReceiver> receivers;
        messageReceiver receiver;
        receiver.id = userInfo.id;
        receiver.account = userInfo.account;
        receiver.name = userInfo.visionName;
        receiver.email = userInfo.email;
        receiver.telephone = userInfo.telNumber;
        receiver.thirdAttr = userInfo.thirdAttr;
        receiver.thirdId = userInfo.thirdId;
        receivers.push_back(receiver);
        msgptr->receivers = receivers;

        msgList.push_back (msgptr);
        _acsMessageManager->AddPluginMessage (msgList);
    }
    else if(!email.empty() && send_vcode_by_email_status){
        // 获取OEM配置信息
        String tmpTitle = _acsConfManager->GetProductName();
        // OEM product
        string title = toSTLString(tmpTitle);
        // 验证码内容
        String content;
        tmpTitle.format (LOAD_STRING(_T("IDS_VCODE_EMAIL_TITLE")), title.c_str());
        content.format (LOAD_STRING(_T("IDS_VCODE_EMAIL_CONTENT")), userInfo.visionName.getCStr(), userInfo.account.getCStr(), vcodeInfo.vcode.c_str());
        // 发送邮箱验证码
        std::vector<string> mailto;
        mailto.push_back(toSTLString(userInfo.email));
        string tmpcontent = toSTLString(content);
        title = toSTLString(tmpTitle);
        ncEACHttpServerUtil::SendMail(mailto, title, tmpcontent);
    }
    else{
        THROW_E (EAC_HTTP_SERVER, EACHTTP_SEND_VCODE_SERVER_NOT_AVAILABLE,
            LOAD_STRING (_T("IDS_EACHTTP_SEND_VCODE_SERVER_NOT_AVAILABLE")));
    }

    // 回复
    return retuuid;
}

void ncEACAuthHandler::SendAuthVcode (brpc::Controller* cntl, const String& fakeUserId)
{
    NC_EAC_HTTP_SERVER_TRACE (_T("this: %p, cntl:%p, userId:%s begin"), this, cntl, fakeUserId.getCStr ());

    // 参数检查
    string account;
    string originPassword;
    bool hasDeviceInfo;
    ncDeviceBaseInfo baseInfo;
    JSON::Value requestJson;
    ncTUserLoginOption option;
    CheckLoginParams(cntl, account, originPassword, hasDeviceInfo, baseInfo, requestJson, option);

    // 获取双因子认证验证码服务的开关
    bool MFASMSServerStatus = false;
    try{
        String tmpString = _acsConfManager->GetAuthVcodeServerStatus();
        string auth_vcode_server_status = toSTLString(tmpString);

        JSON::Value SendAuthVcodeTypeStatusJson;
        JSON::Reader::read(SendAuthVcodeTypeStatusJson, auth_vcode_server_status.c_str(), auth_vcode_server_status.size());
        MFASMSServerStatus = SendAuthVcodeTypeStatusJson["auth_by_sms"];
    }
    catch (TException & e) {
        THROW_E (EAC_HTTP_SERVER, EACHTTP_SERVER_UNKNOWN_ERR, e.what ());
    }

    if (!MFASMSServerStatus){
        THROW_E (EAC_HTTP_SERVER, EACHTTP_SEND_VCODE_SERVER_NOT_AVAILABLE,
            LOAD_STRING (_T("IDS_EACHTTP_SEND_VCODE_SERVER_NOT_AVAILABLE")));
    }

    // 判断是否配置短信验证服务器
    bool enableStatus = ncEACHttpServerUtil::GetThirdAuthTypeStatus(ncTMFAType::SMSAuth);
    if(!enableStatus){
        THROW_E (EAC_HTTP_SERVER, EACHTTP_MFA_SMS_SERVER_NOT_SET,
                LOAD_STRING (_T("IDS_EACHTTP_MFA_SMS_SERVER_NOT_SET")));
    }

    // 进行验证
    string retUserId;
    try {
        ncEACHttpServerUtil::Usrm_UserLogin (retUserId, account, originPassword, option);
    }
    catch(Exception& e) {
        logFailedLoginEvent(cntl, toCFLString(account), e, hasDeviceInfo, baseInfo);
        throw;
    }

    onUserLogin(cntl, toCFLString(retUserId), ncEACHttpServerUtil::GetForwardedIp(cntl), hasDeviceInfo, baseInfo);

    // 生成验证码并发送
    ncTReturnInfo retInfo;
    string oldTelnum = requestJson["oldtelnum"].s ();
    ncEACHttpServerUtil::SendAuthVcode (retInfo, retUserId, ncTVcodeType::DAUL_AUTH_VCODE, oldTelnum);

    // 回复
    JSON::Value replyJson;

    replyJson["authway"] = retInfo.telNumber;
    replyJson["sendinterval"] = retInfo.sendInterval;
    replyJson["isduplicatesended"] = retInfo.isDuplicateSended;

    string body;
    JSON::Writer::write (replyJson.o (), body);
    ncHttpSendReply (cntl, brpc::HTTP_STATUS_OK, "ok", body);

    NC_EAC_HTTP_SERVER_TRACE (_T("this: %p, cntl:%p end"), this, cntl);
}

void ncEACAuthHandler::SendSms (brpc::Controller* cntl, const String& fakeUserId)
{
    NC_EAC_HTTP_SERVER_TRACE (_T("this: %p, cntl:%p begin"), this, cntl);

    // 获取http请求的content
    string bodyBuffer = cntl->request_attachment ().to_string ();
    JSON::Value requestJson;
    try {
        JSON::Reader::read (requestJson, bodyBuffer.c_str (), bodyBuffer.size ());
    }
    catch (Exception& e) {
        THROW_E (EAC_HTTP_SERVER, EACHTTP_INVALID_REQUEST_JSON_FORMAT,
                    LOAD_STRING (_T("IDS_EACHTTP_INVALID_JSON")));
    }

    string account = requestJson["account"].s ();
    if (account.empty ()) {
        THROW_E (EAC_HTTP_SERVER, EACHTTP_ACCOUNT_INVALID,
            LOAD_STRING (_T("IDS_EACHTTP_ACCOUNT_INVALID")));
    }

    string password = requestJson["password"].s ();
    if (password.empty ()) {
        THROW_E (EAC_HTTP_SERVER, EACHTTP_PASSWORD_INVALID,
            LOAD_STRING (_T("IDS_EACHTTP_PASSWORD_INVALID")));
    }

    string telNumber = requestJson["tel_number"].s ();
    if (telNumber.empty ()) {
        THROW_E (EAC_HTTP_SERVER, EACHTTP_INVALID_TEL_NUMBER,
            LOAD_STRING (_T("IDS_EACHTTP_INVALID_TEL_NUMBER")));
    }

    ncEACHttpServerUtil::SMSSendVcode(account, password, telNumber);

    // 回复
    string body;
    ncHttpSendReply (cntl, brpc::HTTP_STATUS_OK, "ok", body);

    NC_EAC_HTTP_SERVER_TRACE (_T("this: %p, cntl:%p end"), this, cntl);
}

void ncEACAuthHandler::SmsActivate (brpc::Controller* cntl, const String& fakeUserId)
{
    NC_EAC_HTTP_SERVER_TRACE (_T("this: %p, cntl:%p begin"), this, cntl);

    // 获取http请求的content
    string bodyBuffer = cntl->request_attachment ().to_string ();
    JSON::Value requestJson;
    try {
        JSON::Reader::read (requestJson, bodyBuffer.c_str (), bodyBuffer.size ());
    }
    catch (Exception& e) {
        THROW_E (EAC_HTTP_SERVER, EACHTTP_INVALID_REQUEST_JSON_FORMAT,
                    LOAD_STRING (_T("IDS_EACHTTP_INVALID_JSON")));
    }

    string account = requestJson["account"].s ();
    if (account.empty ()) {
        THROW_E (EAC_HTTP_SERVER, EACHTTP_ACCOUNT_INVALID,
            LOAD_STRING (_T("IDS_EACHTTP_ACCOUNT_INVALID")));
    }

    string password = requestJson["password"].s ();
    if (password.empty ()) {
        THROW_E (EAC_HTTP_SERVER, EACHTTP_PASSWORD_INVALID,
            LOAD_STRING (_T("IDS_EACHTTP_PASSWORD_INVALID")));
    }

    string telNumber = requestJson["tel_number"].s ();
    if (telNumber.empty ()) {
        THROW_E (EAC_HTTP_SERVER, EACHTTP_INVALID_TEL_NUMBER,
            LOAD_STRING (_T("IDS_EACHTTP_INVALID_TEL_NUMBER")));
    }

    string mailAddress = requestJson["mail_address"].s ();
    if (mailAddress.empty ()) {
        THROW_E (EAC_HTTP_SERVER, EACHTTP_INVALID_EMAIL,
            LOAD_STRING (_T("IDS_EACHTTP_INVALID_EMAIL")));
    }

    string verifyCode = requestJson["verify_code"].s ();
    if (mailAddress.empty ()) {
        THROW_E (EAC_HTTP_SERVER, EACHTTP_SMS_VERIFY_CODE_ERROR,
            LOAD_STRING (_T("IDS_EACHTTP_SMS_VERIFY_CODE_ERROR")));
    }

    // 进行激活
    string retUserId;
    ncEACHttpServerUtil::SMSActivate (retUserId, account, password, telNumber, mailAddress, verifyCode);

#ifndef __UT__
    // 更新用户最近请求时间记录
    _acsShareMgnt->UpdateUserLastRequestTime(toCFLString(retUserId));
#endif

    // 回复
    JSON::Value replyJson;
    replyJson["userid"] = retUserId;
    replyJson["expires"] = _expires;

    string body;
    JSON::Writer::write (replyJson.o (), body);
    ncHttpSendReply (cntl, brpc::HTTP_STATUS_OK, "ok", body);

    NC_EAC_HTTP_SERVER_TRACE (_T("this: %p, cntl:%p end"), this, cntl);
}

void ncEACAuthHandler::ServerTime (brpc::Controller* cntl, const String& fakeUserId)
{
    NC_EAC_HTTP_SERVER_TRACE (_T("this:%p, cntl:%p begin"), this, cntl);

    // 回复
    JSON::Value replyJson;
    replyJson["time"] = BusinessDate::getCurrentTime ();
    string body;
    JSON::Writer::write (replyJson.o (), body);
    ncHttpSendReply (cntl, brpc::HTTP_STATUS_OK, "ok", body);

    NC_EAC_HTTP_SERVER_TRACE (_T("this:%p, cntl:%p end"), this, cntl);
}

void ncEACAuthHandler::CheckSign (brpc::Controller* cntl, const String& bodyContent)
{
    bool needCheck = _acsConfManager->GetConfig("enable_eacp_check_sign").compareIgnoreCase("true") == 0;

    if (needCheck) {
        const string* pUserId = cntl->http_request ().uri ().GetQuery ("userid");
        const string* pSign = cntl->http_request ().uri ().GetQuery ("sign");
        String userId = pUserId ? URLDecode (toCFLString (*pUserId)) : String::EMPTY;
        String sign = pSign ? URLDecode (toCFLString (*pSign)) : String::EMPTY;

        if (sign.isEmpty ()) {
            THROW_E (EAC_HTTP_SERVER, EACHTTP_MODIFY_PASSWORD, LOAD_STRING (_T("IDS_SIGN_NOT_SET")));
        }

        if (sign.compareIgnoreCase (genMD5String2 (bodyContent + userId + "eisoo.com")) != 0) {
            THROW_E (EAC_HTTP_SERVER, EACHTTP_MODIFY_PASSWORD, LOAD_STRING (_T("IDS_SIGN_VERIFY_FAILED")));
        }
    }
}

void ncEACAuthHandler::dualFactorAuth(ncTUserLoginOption &option, JSON::Value& requestJson, bool isOutNet)
{
    // 获取双因子认证验证码服务的开关
    bool MFASMSServerStatus = false;
    bool MFAOTPServerStatus = false;
    try{
        String tmpString = _acsConfManager->GetAuthVcodeServerStatus();
        string auth_vcode_server_status = toSTLString(tmpString);

        JSON::Value SendAuthVcodeTypeStatusJson;
        JSON::Reader::read(SendAuthVcodeTypeStatusJson, auth_vcode_server_status.c_str(), auth_vcode_server_status.size());
        MFASMSServerStatus = SendAuthVcodeTypeStatusJson["auth_by_sms"];
        MFAOTPServerStatus = SendAuthVcodeTypeStatusJson["auth_by_OTP"];
        if (!MFAOTPServerStatus && !MFASMSServerStatus){
            return;
        }
    }
    catch (TException & e) {
        THROW_E (EAC_HTTP_SERVER, EACHTTP_SERVER_UNKNOWN_ERR, e.what ());
    }

    // 获取双因子认证验证信息(暂时只支持短信验证和动态密码)
    string valicode = requestJson["dualfactorauthinfo"]["validcode"]["vcode"].s ();
    string OTP = requestJson["dualfactorauthinfo"]["OTP"]["OTP"].s ();
    boost::replace_all(valicode, " ", "");
    boost::replace_all(OTP, " ", "");

    if (MFASMSServerStatus && isOutNet) {
        // 判断是否配置短信验证服务器
        bool enableStatus = ncEACHttpServerUtil::GetThirdAuthTypeStatus(ncTMFAType::SMSAuth);
        if(!enableStatus){
            THROW_E (EAC_HTTP_SERVER, EACHTTP_MFA_SMS_SERVER_NOT_SET,
                    LOAD_STRING (_T("IDS_EACHTTP_MFA_SMS_SERVER_NOT_SET")));
        }

        if (valicode.empty()) {
            THROW_E (EAC_HTTP_SERVER, EACHTTP_CHECK_VCODE_IS_NULL,
                LOAD_STRING (_T("IDS_EACHTTP_CHECK_VCODE_IS_NULL")));
        }
        option.vcode = valicode;
        option.vcodeType = ncTVcodeType::DAUL_AUTH_VCODE;
        option.__isset.vcodeType = true;
    }

    if (MFAOTPServerStatus) {
        // 判断是否配置动态密保服务器
        bool enableStatus = ncEACHttpServerUtil::GetThirdAuthTypeStatus(ncTMFAType::OTPAuth);
        if(!enableStatus){
            THROW_E (EAC_HTTP_SERVER, EACHTTP_OTP_SERVER_NOT_SET,
                    LOAD_STRING (_T("IDS_EACHTTP_OTP_SERVER_NOT_SET")));
        }


        if (OTP.empty()) {
            THROW_E (EAC_HTTP_SERVER, EACHTTP_CHECK_OTP_IS_NULL,
                LOAD_STRING (_T("IDS_EACHTTP_CHECK_OTP_IS_NULL")));
        }

        option.vcode = OTP;
        option.vcodeType = ncTVcodeType::DAUL_AUTH_OTP;
        option.__isset.OTP = true;
    }
}

bool ncEACAuthHandler::checkIpIsOutNet(const string& ip)
{
    // 获取内部网段配置ip信息
    String tmpString = _acsConfManager->GetConfig("login_sms_net_config");
    if (tmpString.isEmpty()) {
        SystemLog::getInstance ()->log (__FILE__, __LINE__, INFORMATION_LOG_TYPE,
                                            _T("Warning : read login_sms_net_config configuration failed, config is empty ;"));
        return true;
    }

    // 将json字符串转换成JSON::Value
    // 格式类似{"enabled":"true","ip_range":[["192.168.1.1","192.168.1.100"],["192.168.2.1","192.168.2.100"]]}
    // 支持ipv6
    JSON::Value outNetIpList;
    try {
        JSON::Reader::read(outNetIpList, tmpString.getCStr(), tmpString.getLength());
    }
    catch (Exception& e) {
        SystemLog::getInstance ()->log (__FILE__, __LINE__, INFORMATION_LOG_TYPE,
                                            _T("Warning : read login_sms_net_config configuration failed, login_sms_net_config is not json format, config: %s ;"), tmpString.getCStr ());
        return true;
    }

    // 范围判断
    try {
        // 判断配置是否开启，如果开启，则判断ip范围
        if (outNetIpList["enabled"].b() == true) {
            // 判断ip是否在配置范围内,支持多个网段，支持ipv6
            JSON::Array& ipRange = outNetIpList["ip_range"].a ();
            for (int i = 0; i < ipRange.size(); i++) {
                JSON::Array& ipRangeItem = ipRange[i].a ();
                if (ipRangeItem.size() != 2) {
                    // 格式不正确,跳过
                    SystemLog::getInstance ()->log (__FILE__, __LINE__, INFORMATION_LOG_TYPE,
                                            _T("Warning : read login_sms_net_config configuration failed, ip_range format error , size is not 2, config: %s ;"), tmpString.getCStr ());
                    continue;
                }

                string ip1 = ipRangeItem[0].s();
                string ip2 = ipRangeItem[1].s();

                // 判断ip1和ip2是否是ipv4或ipv6
                if ((ip1.find(_T(":")) == String::NO_POSITION && ip2.find(_T(":")) == String::NO_POSITION && ip.find(_T(":")) == String::NO_POSITION) 
                    or (ip1.find(_T(":")) != String::NO_POSITION && ip2.find(_T(":")) != String::NO_POSITION && ip.find(_T(":")) != String::NO_POSITION)) {
                        InfInt ip1Int = ipToInfInt(ip1);
                        InfInt ip2Int = ipToInfInt(ip2);
                        InfInt ipInt = ipToInfInt(ip);
                        if (ipInt >= ip1Int && ipInt <= ip2Int) {
                            return false;
                        }
                } else {
                    SystemLog::getInstance ()->log (__FILE__, __LINE__, INFORMATION_LOG_TYPE,
                                            _T("Warning : read login_sms_net_config configuration failed, ip_range format error , ip1 or ip2 is not ipv4 or ipv6, config: %s ;"), tmpString.getCStr ());
                }
            }
        } else {
            SystemLog::getInstance ()->log (__FILE__, __LINE__, INFORMATION_LOG_TYPE,
                                            _T("Warning : read login_sms_net_config configuration success, but enabled is not true, config: %s ;"), tmpString.getCStr ());
        }
    }
    catch (Exception& e) {
        SystemLog::getInstance ()->log (__FILE__, __LINE__, INFORMATION_LOG_TYPE,
                                            _T("Warning : read login_sms_net_config configuration failed, config: %s ;"), tmpString.getCStr ());
    }

    return true;
}

InfInt ncEACAuthHandler::ipToInfInt(const string& ip)
{
    if (ip.find(_T(":")) == String::NO_POSITION) {
        struct in_addr addr;
        inet_pton(AF_INET, ip.c_str(), &addr);
        return ntohl(addr.s_addr);
    }
    else {
        struct in6_addr addr;
        inet_pton(AF_INET6, ip.c_str(), &addr);
        return InfInt(ntohl(addr.s6_addr32[0])) * InfInt("79228162514264337593543950336") + \
                InfInt(ntohl(addr.s6_addr32[1])) * InfInt("18446744073709551616") + \
                InfInt(ntohl(addr.s6_addr32[2])) * InfInt("4294967296") + \
                InfInt(ntohl(addr.s6_addr32[3]));
    }
}