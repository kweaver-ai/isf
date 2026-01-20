#ifndef __NC_EAC_HTTP_SERVER_UTIL_H
#define __NC_EAC_HTTP_SERVER_UTIL_H

#include "gen-cpp/ncTShareMgnt.h"
#include "gen-cpp/ShareMgnt_constants.h"
#include <gen-cpp/ncTEVFS.h>
#include <gen-cpp/EVFS_constants.h>

#include <acsprocessor/public/ncIACSDeviceManager.h>
#include <drivenadapter/public/stdLogInterface.h>
#include <drivenadapter/public/authenticationInterface.h>
#include "eachttpserver.h"
#include <drivenadapter/public/nsqInterface.h>

struct ncGNSPathInfo {
    String path;
    String name;
    bool isFile;

    ncGNSPathInfo(): path(), name(), isFile(false){}
};


class ncEACHttpServerUtil
{
public:
    // header x-real-ip and X-Forwarded-For
    static String GetRealIp (brpc::Controller* cntl);

    // header X-Forwarded-For and x-real-ip
    static String GetForwardedIp (brpc::Controller* cntl);

    // 文档类型->(文件/文件夹
    static String IsFileToStr(bool isFile);

    // sharemgnt.thrift
    static void Usrm_UserLogin (string& retUserId, const string& account, const string& password, const ncTUserLoginOption& option);
    static void Usrm_Login(string& retUserId, const string& account, const string& password,const ncTUsrmAuthenType::type authenType, const ncTUserLoginOption& option, const string &osType);
    static bool Usrm_GetADSSOStatus ();
    static ncTThirdPartyAuthConf Usrm_GetThirdPartyAuth ();
    static String Usrm_LoginConsoleByThirdPartyNew(String& params);
    static String Usrm_ValidateThirdParty (string& params);
    static bool Usrm_ValidateSecurityDevice (string& params);
    static bool GetLoginStrategyStatus ();
    static void SendMail(vector<string>& mailto, string& subject, string& content);
    static bool OEM_GetConfigByOption(string section, string option);
    static void Usrm_UserLoginByNTLMV1 (ncTNTLMResponse& ntlmResp, const string& account, const string& challenge, const string& password);
    static void Usrm_UserLoginByNTLMV2 (ncTNTLMResponse& ntlmResp, const string& account, const string& domain, const string& challenge, const string& password);
    static ncTThirdPartyToolConfig GetThirdPartyToolConfig(const string& thirdPartyToolId);
    static bool Secretm_GetStatus();
    static void GetCSFLevels(map<string, int32_t> & csflevels);
    static void CheckUninstallPwd(const String& uninstallPwd);
    static void CheckExitPwd(const String& exitPwd);
    static void GetThirdCSFSysConfig(ncTThirdCSFSysConfig & config);
    static bool GetShareDocStatus(int docType, int linkType);
    static bool HideOum_Check(const string& userId);
    static void Usrm_CreateVcodeInfo(ncTVcodeCreateInfo& vcodeInfo, const string& uuid, const ncTVcodeType::type vcodeType);
    static void SetUserInfo (const String& userId, map<String, String> &userinfo);
    static void SMSSendVcode(const string& account, const string& password, const string& telNumber);
    static void SMSActivate(string& retUserId, const string& account, const string& password, const string& telNumber, const string& mailAddress, const string& verifyCode);
    static void SendAuthVcode (ncTReturnInfo &retInfo, string& userId, const ncTVcodeType::type vcodeType, const string oldTelnum);
    static bool GetThirdAuthTypeStatus (const ncTMFAType::type authType);
    static void GetSmtpSrvConfig(ncTSmtpSrvConf& config);
    static void GetAllConfig(ncTAllConfig& config);
    static int GetCustomConfigOfInt64(const string& key);

    // eacplog.thrift
    static void Log (brpc::Controller* cntl, const String& userId, ncTokenVisitorType typ, ncTLogType logType,
                     ncTLogLevel level, int opType, const String& msg, const String& exmsg, bool logForwardedIp = false);

    static void LoginLog (const String& userId, const String& udid, ACSClientType clientType, const String& ip);

    static void ClientLogin (string& retUserId, const string& account, const string& password, const ncTUserLoginOption& option);

    // ["abc", "def", "ghi"] -> "\"abc\",\"def\",\"ghi\""
    static String GenerateGroupStr (const vector<String>& strs);

    // 将截至时间转为str
    static String EndTimeToStr (int64 endTime);

    // 将访问者类型可读
    static String AccessorTypeToResStr (int accessorType);

    // 将权限转为可读str
    static String ConvertPermToStr (int linkPerm);

    // 从字符串中获取数字字符
    static int ParseLockTime (const char *str);

    // 格式化密级字符串
    static String GetFormatCSFLevelStr (const int& CSFLevel);

    // 检查是否为UTF8字符串并获取UTF8字符串长度
    static size_t CheckIsTextUTF8 (const String& str);

    // 检查设备类型
    static void CheckOsType(ACSClientType clientType);

    // 判断是否符合RFC3339格式
    static bool CheckIsRFC3339(const String& rfcTimeStr);

    // RFC3339转微秒时间戳
    static int64 RFC3339ToTimeStamp (String& rfcTimeStr);

    static void GetVisitor (brpc::Controller* cntl, ncIntrospectInfo& introspectInfo);

    // 权限数组转换为权限值
    static int PermArrayToInt (const JSON::Array& permArray, const String& paramName);

    // 权限值转换为权限数组
    static void IntToPermArray (int perm, JSON::Array& permArray);
    // 权限值转换为权限数组 带read兼容旧接口使用
    static void IntToPermArrayV1 (int perm, JSON::Array& permArray);
    // RSA解密
    static string RSADecrypt2048(const string &str);
    static string RSADecrypt(const string &cipherText);
    // Base64解密
    static string Base64Decode(const string &str);
    static size_t calcDecodeLength(const string &tmp);

private:
};

#endif  // __NC_EAC_HTTP_SERVER_UTIL_H
