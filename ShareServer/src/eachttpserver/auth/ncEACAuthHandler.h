#ifndef __NC_EAC_AUTH_HANDLER_H__
#define __NC_EAC_AUTH_HANDLER_H__

#include <biginteger/InfInt.h>

#include <acssharemgnt/public/ncIACSShareMgnt.h>
#include <acsprocessor/public/ncIACSTokenManager.h>
#include <acsprocessor/public/ncIACSDeviceManager.h>
#include <acsprocessor/public/ncIACSConfManager.h>
#include <acsprocessor/public/ncIACSMessageManager.h>
#include <acsprocessor/public/ncIACSPolicyManager.h>
#include "drivenadapter/public/policyEngineInterface.h"

#include "ncEACHttpServerUtil.h"

#define RESET_PASSWORD_VERIFICATION_CODE_CHANNEL   _T ("authentication/v1/reset-pwd-verification-code")

//账户登录凭证信息
struct ncAccountCredential
{
    String account;
    String originPassword;
    ncTUserLoginOption option;
};

//第三方登录凭证信息
struct ncThirdPartyCredential
{
    String params;
};

class ncEACAuthHandler
{
public:
    ncEACAuthHandler(ncIACSTokenManager *acsTokenManager,
                     ncIACSShareMgnt *acsShareMgnt,
                     ncIACSDeviceManager *acsDeviceManager,
                     ncIACSConfManager *acsConfManager,
                     ncIACSMessageManager *acsMessageManager,
                     ncIACSPolicyManager *acsPolicyManager,
                     policyEngineInterface *policyEngine);
    ~ncEACAuthHandler (void);

    void doAuthRequestHandler (brpc::Controller* cntl);
    void doAuth2RequestHandler (brpc::Controller* cntl);
    void setExpires (int64);

    /***
     * 获取身份凭证（明文账号和（Base64（RSA（明文）））
     */
    void GetNew (brpc::Controller* cntl);

    /***
     * 控制台管理员登录
     */
    void ConsoleLogin (brpc::Controller* cntl);

    /***
     * 获取身份凭证（第三方凭证信息）
     */
    void GetByThirdParty (brpc::Controller* cntl);

    /***
     * 记录登录可观测性日志
     */
    void LoginLog(brpc::Controller* cntl);

protected:

    /***
     * 获取认证配置信息
     */
    void GetConfig (brpc::Controller* cntl, const String& userId);

    /***
     * 获取登录配置信息
     */
    void GetLoginConfigs(brpc::Controller* cntl, const String& userId);

    /***
     * 获取服务端配置信息
     */
    void GetServerConfig (brpc::Controller* cntl, const String& userId);


    /***
     * 根据NTLMV1算法进行验证
     */
    void GetByNTLMV1 (brpc::Controller* cntl, const String& userId);

    /***
     * 获取身份凭证（西电ticket）
     */
    void GetByTicket (brpc::Controller* cntl, const String& userId);

    /***
     * 获取身份凭证（windows ad会话凭证）
     */
    void GetByADSession (brpc::Controller* cntl, const String& userId);

    /***
    * 修改密码
    */
    void ModifyPassword(brpc::Controller* cntl, const String& userId);

    /***
    * 二次安全设备认证
    */
    void ValidateSecurityDevice(brpc::Controller* cntl, const String& userId);

    /*
    从头信息中的获取Section:
    */
    string GetSection (String& language_config);

    /***
    * 检查卸载口令是否正确
    */
    void CheckUninstallPwd(brpc::Controller* cntl, const String& userId);

    /***
    * 检查退出口令是否正确
    */
    void CheckExitPwd(brpc::Controller* cntl, const String& userId);

    /***
    * v2版本统一的登录协议
    */
    void Login(brpc::Controller* cntl, const String& userId);

    /***
    * 获取验证码(参数中存在有效uuid时会删除该验证码并创建一个新的)
    */
    void GetVcode(brpc::Controller* cntl, const String& userId);

    /***
    * 发送短信验证码
    */
    void SendSms (brpc::Controller* cntl, const String& fakeUserId);

    /***
    * 发送双因子认证短信验证码
    */
    void SendAuthVcode (brpc::Controller* cntl, const String& fakeUserId);

    /***
    * 发送验证码(邮箱/短信)
    */
    void SendVcode (brpc::Controller* cntl, const String& fakeUserId);

    /***
    * 根据账户名和发送方式发送验证码(邮箱/短信)
    */
    void SendPwdRetrevalVcode(brpc::Controller* cntl, const String& fakeUserId);

    /***
    * 激活账号
    */
    void SmsActivate (brpc::Controller* cntl, const String& fakeUserId);

    /***
    * 获取服务器时间
    */
    void ServerTime (brpc::Controller* cntl, const String& fakeUserId);

    /***
    * 鉴权，检查签名
    */
    void CheckSign (brpc::Controller* cntl, const String& bodyContent);

    /***
    * 检查登录参数
    */
    void CheckLoginParams (brpc::Controller* cntl, string &account, string &originPassword, bool &hasDeviceInfo, ncDeviceBaseInfo &baseInfo, JSON::Value &requestJson, ncTUserLoginOption &option, bool b2048 = false);

    /***
    * 检查管理员登录参数
    */
    void CheckConsoleLoginParams (brpc::Controller* cntl, bool &isAccountType, ncAccountCredential &accountCredential,
                                  ncThirdPartyCredential &thirdPartyCredential, ncDeviceBaseInfo &baseInfo, String& loginIp, JSON::Value &requestJson);
    /***
    * 检查登录IP
    */
    void CheckLoginIpParam(JSON::Value &requestJson, String& loginIp);

private:
    typedef void (ncEACAuthHandler::*ncMethodFunc) (brpc::Controller* cntl, const String&);
    map<String, ncMethodFunc>            _methodFuncs;
    map<String, ncMethodFunc>            _v2MethodFuncs;

    typedef String (ncEACAuthHandler::*ncThirdAuthFunc) (JSON::Value&, const ncTThirdPartyAuthConf&);
    map<String, ncThirdAuthFunc>            _thirdAuthFuncs;

    // RSA加解密
    string RSAEncrypt(const string& str, const char *path_key);
    string RSADecrypt(const string &str, const char *path_key);

    // Base64加密
    string Base64Encode(const string &plainText);

    // DES加解密
    string DESEncrypt(const string& in);
    string DESDecrypt(const string &in);

    // 根据ticket获取第三方用户的id
    String OAuthExecute (const String& tokenServer, const String& ticket, const String& service);

    // 根据xml结果解析第三方用户的id
    String ParseThirdUserId (const string& retXMlStr);

    // 根据xml结果解析第三方用户的id
    void ParseADSession (const String& session, String& account);

    // AnyShare 使用明文密码进行认证
    String AnySharePlain(JSON::Value& requestJson, const ncTThirdPartyAuthConf& config);

    // AnyShare 使用RSA对密码进行认证
    String AnyShareRSA(JSON::Value& requestJson, const ncTThirdPartyAuthConf& config);

    // 使用AnyShare颁发的appid和appkey进行认证
    String AnyShareSSO(JSON::Value& requestJson, const ncTThirdPartyAuthConf& config);

    // Windows ad单点登录
    String WindowsADSSO(JSON::Value& requestJson, const ncTThirdPartyAuthConf& config);

    // 佰勤科技认证
    String BQKJAuth(JSON::Value& requestJson, const ncTThirdPartyAuthConf& config);

    // 秉英OA认证
    String BeingAuth(JSON::Value& requestJson, const ncTThirdPartyAuthConf& config);

    // 思路认证
    String ThsAuth(JSON::Value& requestJson, const ncTThirdPartyAuthConf& config);

    // 龙创认证
    String LcsoftAuth(JSON::Value& requestJson, const ncTThirdPartyAuthConf& config);

    // 龙创认证
    String AISHUAuth(JSON::Value& requestJson, const ncTThirdPartyAuthConf& config);

    // 打印字节
    void printBytes(const string& str);

    // 解析登录时的deviceinfo
    bool parseDeviceInfo(JSON::Value& requestJson, ncDeviceBaseInfo& info);

    // 解析登录时的vcodeinfo
    bool parseVcodeInfo(JSON::Value& requestJson, ncTUserLoginOption &option);

    // url编码
    char *url_encode(char *str);

    // url解码
    char *url_decode(char *str);

    // 16进制转换为char
    char from_hex(char ch);

    // char转换为16进制
    char to_hex(char code);

    // 客户端类型
    String clientType2Str(ACSClientType clientType);

    // 登录时设备绑定响应
    void onUserLogin(brpc::Controller* cntl, const String& retUserId, const String& ip, bool hasDeviceInfo, ncDeviceBaseInfo& baseInfo);

    // 更新请求头中的MAC地址信息，记录对应token信息和日志
    void updateHTTPHeaderMACAddress(brpc::Controller* cntl, const String& macAddr);

    //更新请求头部的ip信息
    void updateHTTPHeaderIP(brpc::Controller* cntl, const String& ip);

    // 认证异常时记录日志
    void authenicaitonFailedLoginEvent(brpc::Controller* cntl,
                            const String& account,
                            const Exception& e,
                            bool hasDeviceInfo,
                            const ncDeviceBaseInfo& baseInfo);

    // 登录异常时记录日志
    void logFailedLoginEvent(brpc::Controller* cntl,
                            const String& account,
                            const Exception& e,
                            bool hasDeviceInfo,
                            const ncDeviceBaseInfo& baseInfo);

    // 管理控制台登录异常时记录日志
    void logFailedConsoleLoginEvent(brpc::Controller* cntl,
                            const string& osType,
                            const EHttpDetailException& e);

    // 计算expired time
    int getExpiredTime(const string& tokenType);

    // 分析account
    void parseAccount(JSON::Value& requestJson, String& account);

    // 获取默认权限过期时间
    int getDefaultPermExpiredDays();

    // 获取服务端密级配置
    JSON::Value getCsfLevelsConfig();

    // 客户端版本检查
    void checkClientVersion(ACSClientType clientType, const String& version);

    // 比较版本号，5.0.11.253 形式（主版本号，小版本号，修订号，构建号)
    int compareVersion(const String& version1, const String& version2);

    // 双因子认证
    void dualFactorAuth(ncTUserLoginOption &option, JSON::Value& requestJson, bool isOutNet);

    // 向邮箱和手机号发送验证码
    string sendEmailAndTelVcode (const string& email, const string& telnumber, const string& uuidIn);

    // 检查ip是否在外网
    bool checkIpIsOutNet(const string& ip);

    // ip转换为InfInt
    InfInt ipToInfInt(const string& ip);

private:
    nsCOMPtr<ncIACSTokenManager>        _acsTokenManager;   // token管理
    nsCOMPtr<ncIACSShareMgnt>           _acsShareMgnt;      // sharemgnt db查询
    nsCOMPtr<ncIACSDeviceManager>       _acsDeviceManager;   // 设备管理
    nsCOMPtr<ncIACSConfManager>         _acsConfManager;     // oem配置管理
    nsCOMPtr<ncIACSPolicyManager>       _acsPolicyManager;   // 策略管理
    nsCOMPtr<policyEngineInterface> _policyEngine;           // 策略引擎决策

    int64                               _expires;           // 默认的token有效期。

    set<string>                         _thirdConfigBlackList; //第三方配置信息过滤黑名单
    set<string>                         _adminIds;     // 系统管理员id列表

    set<string>                          _needAuthMethod;   // 需要鉴权验证的method方法
    nsCOMPtr<ncIACSMessageManager>      _acsMessageManager;      // 消息通知

    map<int, String>                    _accountStringTypeMap;
    map<ACSClientType, String>          _clientStringTypeMap;
    map<int, String>                    _visitorStringTypeMap;
    map<String, int>                    _clientIntTypeMap;
};

#endif  // __NC_EAC_AUTH_HANDLER_H__
