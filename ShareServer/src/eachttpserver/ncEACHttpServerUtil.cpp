#include <abprec.h>
#include <ncutil/ncPerformanceProfilerPrec.h>
#include <ncutil/ncBusinessDate.h>
#include <boost/regex.hpp>
#include <butil/endpoint.h>
#include "ncEACHttpServerUtil.h"
#include <evfsioc/ncIEVFSNameIOC.h>
#include <ehttpserver/ncEHttpUtil.h>
#include <ethriftutil/ncThriftClient.h>

#include <boost/thread/locks.hpp>
#include <boost/thread/shared_mutex.hpp>
#include <acsprocessor/public/ncIACSPermManager.h>
#include <acssharemgnt/public/ncIACSShareMgnt.h>
#include "eacServiceAccessConfig.h"
#include "drivenadapter/public/nsqInterface.h"

#define PREFIXSTR _T("AnyShare://")

#define GNS_PREFIX                _T("gns://")
#define GNS_PREFIX_LENGTH        6


// 该函数获取容器内部IP,请勿使用
// 使用GetForwardedIp 可获取http报头X-Forwarded-For字段中的IP
String ncEACHttpServerUtil::GetRealIp (brpc::Controller* cntl)
{
    String ip = butil::endpoint2str (cntl->remote_side ()).c_str ();
    int index = ip.rfind (":");
    ip = ip.subString (0, index);

    int start = (ip.find ("[") == String::NO_POSITION) ? 0 : ip.find ("[") + 1;
    int end = (ip.find ("]") == String::NO_POSITION) ? ip.getLength () : ip.find ("]") - start;
    ip = ip.subString (start, end);

    if ((ip == "127.0.0.1") || (ip == "localhost")) {
        String webIp;
        ncHttpGetHeader (cntl, "X-Real-IP", webIp);
        if (webIp.isEmpty ()) {
            ncHttpGetHeader (cntl, "X-Forwarded-For", webIp);
            if (!webIp.isEmpty ()) {
                vector<String> webIps;
                webIp.split (",", webIps);
                webIp = webIps[0];
            }
        }
        if (!webIp.isEmpty ()) {
            ip = webIp;
        }
    }

    NC_EAC_HTTP_SERVER_TRACE (_T("ip: %s"), ip.getCStr ());
    return ip;
}

String ncEACHttpServerUtil::GetForwardedIp (brpc::Controller* cntl)
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

    NC_EAC_HTTP_SERVER_TRACE (_T("ip: %s"), ip.getCStr ());
    return ip;
}

void ncEACHttpServerUtil::Usrm_Login(string& retUserId, const string& account, const string& password, const ncTUsrmAuthenType::type authenType, const ncTUserLoginOption& option, const string &osType)
{
#ifndef __UT__
    try {
        ncThriftClient<ncTShareMgntClient> shareMgntClient (EacServiceAccessConfig::getInstance()->sharemgntHost, EacServiceAccessConfig::getInstance()->sharemgntPort);
        shareMgntClient->Usrm_Login (retUserId, account, password, authenType, option, osType);
    }
    catch (ncTException & e) {
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

        if (e.errID == ncTShareMgntError::NCT_INVALID_ACCOUNT_OR_PASSWORD) {
            THROW_HTTP_DETAIL_E (EAC_HTTP_SERVER, errDetail, EACHTTP_INVALID_USER_OR_PASSWORD, e.expMsg.c_str ());
        }
        else if (e.errID == ncTShareMgntError::NCT_USER_DISABLED) {
            THROW_HTTP_DETAIL_E (EAC_HTTP_SERVER, errDetail, EACHTTP_USER_DISABLED, e.expMsg.c_str ());
        }
        else if (e.errID == ncTShareMgntError::NCT_THIRD_PARTY_AUTH_NOT_OPEN) {
            THROW_E (EAC_HTTP_SERVER, EACHTTP_THIRD_AUTH_NOT_OPEN, e.expMsg.c_str ());
        }
        else if (e.errID == ncTShareMgntError::NCT_NEED_THIRD_OAUTH) {
            THROW_E (EAC_HTTP_SERVER, EACHTTP_OAUTH_NEEDED, e.expMsg.c_str ());
        }
        else if (e.errID == ncTShareMgntError::NCT_FORBIDDEN_LOGIN) {
            THROW_HTTP_DETAIL_E (EAC_HTTP_SERVER, errDetail, EACHTTP_NORMAL_FORBIDDEN_LOGIN, e.expMsg.c_str ());
        }
        else if (e.errID == ncTShareMgntError::NCT_DOMAIN_NOT_EXIST) {
            THROW_E (EAC_HTTP_SERVER, EACHTTP_DOMAIN_NOT_EXIST, e.expMsg.c_str ());
        }
        else if (e.errID == ncTShareMgntError::NCT_DOMAIN_DISABLE) {
            THROW_E (EAC_HTTP_SERVER, EACHTTP_DOMAIN_DISABLE, e.expMsg.c_str ());
        }
        else if (e.errID == ncTShareMgntError::NCT_DOMAIN_USER_NOT_EXIST) {
            THROW_E (EAC_HTTP_SERVER, EACHTTP_DOMAIN_USER_NOT_EXIST, e.expMsg.c_str ());
        }
        else if(e.errID == ncTShareMgntError::NCT_DOMAIN_SERVER_UNAVAILABLE){
            THROW_E(EAC_HTTP_SERVER, EACHTTP_DOMAIN_SERVER_UNAVAILABLE, e.expMsg.c_str ());
        }
        else if(e.errID == ncTShareMgntError::NCT_PRODUCT_NOT_AUTHORIZED){
            THROW_HTTP_DETAIL_E (EAC_HTTP_SERVER, errDetail, EACHTTP_USER_NOT_AUTORIZED, e.expMsg.c_str ());
        }
        else if(e.errID == ncTShareMgntError::NCT_PRODUCT_HAS_EXPIRED){
            THROW_HTTP_DETAIL_E (EAC_HTTP_SERVER, errDetail, EACHTTP_USER_HAS_EXPIRED, e.expMsg.c_str ());
        }
        else if(e.errID == ncTShareMgntError::NCT_PASSWORD_EXPIRE){
            THROW_HTTP_DETAIL_E (EAC_HTTP_SERVER, errDetail, EACHTTP_PASSWORD_EXPIRE, e.expMsg.c_str ());
        }
        else if(e.errID == ncTShareMgntError::NCT_PASSWORD_NOT_SAFE){
            THROW_HTTP_DETAIL_E (EAC_HTTP_SERVER, errDetail, EACHTTP_PASSWORD_NOT_SAFE, e.expMsg.c_str ());
        }
        else if(e.errID == ncTShareMgntError::NCT_PASSWORD_IS_INITIAL){
            THROW_HTTP_DETAIL_E (EAC_HTTP_SERVER, errDetail, EACHTTP_PASSWORD_IS_INITIAL, e.expMsg.c_str ());
        }
        else if(e.errID == ncTShareMgntError::NCT_PWD_FIRST_FAILED){
            THROW_E(EAC_HTTP_SERVER, EACHTTP_PWD_FAILED_FIRSTLY, e.expMsg.c_str ());
        }
        else if(e.errID == ncTShareMgntError::NCT_PWD_SECOND_FAILED){
            THROW_E(EAC_HTTP_SERVER, EACHTTP_PWD_FAILED_SECONDLY, e.expMsg.c_str ());
        }
        else if(e.errID == ncTShareMgntError::NCT_PWD_THIRD_FAILED){
            JSON::Value detailJson;
            detailJson["remainlockTime"] = errDetail["remainlockTime"];
            detailJson["isShowStatus"] = errDetail["isShowStatus"];
            detailJson["id"] = errDetail["id"];
            THROW_HTTP_DETAIL_E(EAC_HTTP_SERVER, detailJson, EACHTTP_PWD_FAILED_THIRDLY, e.expMsg.c_str ());
        }
        else if(e.errID == ncTShareMgntError::NCT_ACCOUNT_LOCKED){
            // parse lock tim
            JSON::Value detailJson;
            detailJson["remainlockTime"] = errDetail["remainlockTime"];
            detailJson["isShowStatus"] = errDetail["isShowStatus"];
            THROW_HTTP_DETAIL_E(EAC_HTTP_SERVER, detailJson, EACHTTP_ACCOUNT_LOCKED, e.expMsg.c_str ());
        }
        else if(e.errID == ncTShareMgntError::NCT_CANNOT_CONNECT_THID_PARTY_SERVER){
            THROW_E(EAC_HTTP_SERVER, EACHTTP_THIRD_SERVER_UNAVALIABLE, e.expMsg.c_str ());
        }
        else if(e.errID == ncTShareMgntError::NCT_USER_NOT_EXIST){
            THROW_E(EAC_HTTP_SERVER, EACHTTP_USER_NOT_EXIST, e.expMsg.c_str ());
        }
        else if(e.errID == ncTShareMgntError::NCT_WRONG_PASSWORD){
            THROW_E(EAC_HTTP_SERVER, EACHTTP_WRONG_PASSWORD, e.expMsg.c_str ());
        }
        else if(e.errID == ncTShareMgntError::NCT_CONTROLED_PASSWORD_EXPIRE){
            THROW_E(EAC_HTTP_SERVER, EACHTTP_CONTROLED_PASSWORD_EXPIRE, e.expMsg.c_str ());
        }
        else if(e.errID == ncTShareMgntError::NCT_CANNOT_LOGIN_SLAVE_SITE){
            THROW_HTTP_DETAIL_E (EAC_HTTP_SERVER, errDetail, EACHTTP_CANNOT_LOGIN_SLAVE_SITE, e.expMsg.c_str ());
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
        else if (e.errID == ncTShareMgntError::NCT_INSUFFICIENT_SYSTEM_RESOURCES) {
            THROW_HTTP_DETAIL_E (EAC_HTTP_SERVER, errDetail, EACHTTP_INSUFFICIENT_SYSTEM_RESOURCES, e.expMsg.c_str ());
        }
        else if (e.errID == ncTShareMgntError::NCT_USER_NOT_ACTIVATE) {
            THROW_E(EAC_HTTP_SERVER, EACHTTP_USER_NOT_ACTIVATE, e.expMsg.c_str());
        }
        else if (e.errID == ncTShareMgntError::NCT_OTP_WRONG) {
            THROW_E(EAC_HTTP_SERVER, EACHTTP_OTP_WRONG, e.expMsg.c_str());
        }
        else if (e.errID == ncTShareMgntError::NCT_OTP_TIMEOUT) {
            THROW_E(EAC_HTTP_SERVER, EACHTTP_OTP_TIMEOUT, e.expMsg.c_str());
        }
        else if (e.errID == ncTShareMgntError::NCT_OTP_TOO_MANY_WRONG_TIME) {
            THROW_E(EAC_HTTP_SERVER, EACHTTP_OTP_TOO_MANY_WRONG_TIME, e.expMsg.c_str());
        }
        else if (e.errID == ncTShareMgntError::NCT_CHECK_VCODE_MORE_THAN_THE_LIMIT) {
            THROW_E(EAC_HTTP_SERVER, EACHTTP_CHECK_VCODE_MORE_THAN_THE_LIMIT, e.expMsg.c_str());
        }
        else if(e.errID == ncTShareMgntError::NCT_MFA_OTP_SERVER_ERROR){
            THROW_E(EAC_HTTP_SERVER, EACHTTP_MFA_OTP_SERVER_ERROR, e.expMsg.c_str ());
        }
        else if(e.errID == ncTShareMgntError::NCT_THIRD_PLUGIN_INTER_ERROR){
            THROW_E(EAC_HTTP_SERVER, EACHTTP_THIRD_PLUGIN_INTER_ERROR, e.expMsg.c_str ());
        }
        else if(e.errID == ncTShareMgntError::NCT_MFA_CONFIG_ERROR){
            THROW_E(EAC_HTTP_SERVER, EACHTTP_MFA_CONFIG_ERROR, e.expMsg.c_str ());
        }
        else if(e.errID == ncTShareMgntError::NCT_FAILED_THIRD_CONFIG){
            THROW_E(EAC_HTTP_SERVER, EACHTTP_FAILED_THIRD_CONFIG, e.expMsg.c_str ());
        }
        else if(e.errID == ncTShareMgntError::NCT_USER_LOGIN_IP_RESTRICTED)
        {
            THROW_E(EAC_HTTP_SERVER, EACHTTP_LOGIN_IP_IS_RESTRICTED, e.expMsg.c_str ());
        }
        else if(e.errID == ncTShareMgntError::NCT_ACCOUNT_CANNOT_LOGIN_IN_SECRET_NODE)
        {
            THROW_E(EAC_HTTP_SERVER, EACHTTP_ACCOUNT_CANNOT_LOGIN_IN_SECRET_MODE, e.expMsg.c_str ());
        }
        else if(e.errID == ncTShareMgntError::NCT_INVALID_PARAMTER)
        {
            THROW_E(EAC_HTTP_SERVER, EACHTTP_INVALID_PARAM_VALUE, e.expMsg.c_str ());
        }
        else {
            THROW_HTTP_DETAIL_E (EAC_HTTP_SERVER, errDetail, USER_LOGIN_ERROR, e.expMsg.c_str ());
        }
    }
    catch (TException & e) {
        printf(e.what ());
        THROW_E (EAC_HTTP_SERVER, USER_LOGIN_ERROR, e.what ());
    }
#endif
}

void ncEACHttpServerUtil::Usrm_UserLogin (string& retUserId, const string& account, const string& password, const ncTUserLoginOption& option)
{
#ifndef __UT__
    try {
        ncThriftClient<ncTShareMgntClient> shareMgntClient (EacServiceAccessConfig::getInstance()->sharemgntHost, EacServiceAccessConfig::getInstance()->sharemgntPort);
        shareMgntClient->Usrm_UserLogin (retUserId, account, password, option);
    }
    catch (ncTException & e) {
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

        if (e.errID == ncTShareMgntError::NCT_INVALID_ACCOUNT_OR_PASSWORD) {
            THROW_HTTP_DETAIL_E (EAC_HTTP_SERVER, errDetail, EACHTTP_INVALID_USER_OR_PASSWORD, e.expMsg.c_str ());
        }
        else if (e.errID == ncTShareMgntError::NCT_USER_DISABLED) {
            THROW_HTTP_DETAIL_E (EAC_HTTP_SERVER, errDetail, EACHTTP_USER_DISABLED, e.expMsg.c_str ());
        }
        else if (e.errID == ncTShareMgntError::NCT_THIRD_PARTY_AUTH_NOT_OPEN) {
            THROW_E (EAC_HTTP_SERVER, EACHTTP_THIRD_AUTH_NOT_OPEN, e.expMsg.c_str ());
        }
        else if (e.errID == ncTShareMgntError::NCT_NEED_THIRD_OAUTH) {
            THROW_E (EAC_HTTP_SERVER, EACHTTP_OAUTH_NEEDED, e.expMsg.c_str ());
        }
        else if (e.errID == ncTShareMgntError::NCT_FORBIDDEN_LOGIN) {
            THROW_HTTP_DETAIL_E (EAC_HTTP_SERVER, errDetail, EACHTTP_FORBIDDEN_LOGIN, e.expMsg.c_str ());
        }
        else if (e.errID == ncTShareMgntError::NCT_DOMAIN_NOT_EXIST) {
            THROW_E (EAC_HTTP_SERVER, EACHTTP_DOMAIN_NOT_EXIST, e.expMsg.c_str ());
        }
        else if (e.errID == ncTShareMgntError::NCT_DOMAIN_DISABLE) {
            THROW_E (EAC_HTTP_SERVER, EACHTTP_DOMAIN_DISABLE, e.expMsg.c_str ());
        }
        else if (e.errID == ncTShareMgntError::NCT_DOMAIN_USER_NOT_EXIST) {
            THROW_E (EAC_HTTP_SERVER, EACHTTP_DOMAIN_USER_NOT_EXIST, e.expMsg.c_str ());
        }
        else if(e.errID == ncTShareMgntError::NCT_DOMAIN_SERVER_UNAVAILABLE){
            THROW_E(EAC_HTTP_SERVER, EACHTTP_DOMAIN_SERVER_UNAVAILABLE, e.expMsg.c_str ());
        }
        else if(e.errID == ncTShareMgntError::NCT_PRODUCT_NOT_AUTHORIZED){
            THROW_HTTP_DETAIL_E (EAC_HTTP_SERVER, errDetail, EACHTTP_USER_NOT_AUTORIZED, e.expMsg.c_str ());
        }
        else if(e.errID == ncTShareMgntError::NCT_PRODUCT_HAS_EXPIRED){
            THROW_HTTP_DETAIL_E (EAC_HTTP_SERVER, errDetail, EACHTTP_USER_HAS_EXPIRED, e.expMsg.c_str ());
        }
        else if(e.errID == ncTShareMgntError::NCT_PASSWORD_EXPIRE){
            THROW_HTTP_DETAIL_E (EAC_HTTP_SERVER, errDetail, EACHTTP_PASSWORD_EXPIRE, e.expMsg.c_str ());
        }
        else if(e.errID == ncTShareMgntError::NCT_PASSWORD_NOT_SAFE){
            THROW_HTTP_DETAIL_E (EAC_HTTP_SERVER, errDetail, EACHTTP_PASSWORD_NOT_SAFE, e.expMsg.c_str ());
        }
        else if(e.errID == ncTShareMgntError::NCT_PASSWORD_IS_INITIAL){
            THROW_HTTP_DETAIL_E (EAC_HTTP_SERVER, errDetail, EACHTTP_PASSWORD_IS_INITIAL, e.expMsg.c_str ());
        }
        else if(e.errID == ncTShareMgntError::NCT_PWD_FIRST_FAILED){
            THROW_E(EAC_HTTP_SERVER, EACHTTP_PWD_FAILED_FIRSTLY, e.expMsg.c_str ());
        }
        else if(e.errID == ncTShareMgntError::NCT_PWD_SECOND_FAILED){
            THROW_E(EAC_HTTP_SERVER, EACHTTP_PWD_FAILED_SECONDLY, e.expMsg.c_str ());
        }
        else if(e.errID == ncTShareMgntError::NCT_PWD_THIRD_FAILED){
            JSON::Value detailJson;
            detailJson["remainlockTime"] = errDetail["remainlockTime"];
            detailJson["isShowStatus"] = errDetail["isShowStatus"];
            THROW_HTTP_DETAIL_E(EAC_HTTP_SERVER, detailJson, EACHTTP_PWD_FAILED_THIRDLY, e.expMsg.c_str ());
        }
        else if(e.errID == ncTShareMgntError::NCT_ACCOUNT_LOCKED){
            // parse lock tim
            JSON::Value detailJson;
            detailJson["remainlockTime"] = errDetail["remainlockTime"];
            detailJson["isShowStatus"] = errDetail["isShowStatus"];
            THROW_HTTP_DETAIL_E(EAC_HTTP_SERVER, detailJson, EACHTTP_ACCOUNT_LOCKED, e.expMsg.c_str ());
        }
        else if(e.errID == ncTShareMgntError::NCT_CANNOT_CONNECT_THID_PARTY_SERVER){
            THROW_E(EAC_HTTP_SERVER, EACHTTP_THIRD_SERVER_UNAVALIABLE, e.expMsg.c_str ());
        }
        else if(e.errID == ncTShareMgntError::NCT_USER_NOT_EXIST){
            THROW_E(EAC_HTTP_SERVER, EACHTTP_USER_NOT_EXIST, e.expMsg.c_str ());
        }
        else if(e.errID == ncTShareMgntError::NCT_WRONG_PASSWORD){
            THROW_E(EAC_HTTP_SERVER, EACHTTP_WRONG_PASSWORD, e.expMsg.c_str ());
        }
        else if(e.errID == ncTShareMgntError::NCT_CONTROLED_PASSWORD_EXPIRE){
            THROW_E(EAC_HTTP_SERVER, EACHTTP_CONTROLED_PASSWORD_EXPIRE, e.expMsg.c_str ());
        }
        else if(e.errID == ncTShareMgntError::NCT_CANNOT_LOGIN_SLAVE_SITE){
            THROW_HTTP_DETAIL_E (EAC_HTTP_SERVER, errDetail, EACHTTP_CANNOT_LOGIN_SLAVE_SITE, e.expMsg.c_str ());
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
        else if (e.errID == ncTShareMgntError::NCT_INSUFFICIENT_SYSTEM_RESOURCES) {
            THROW_HTTP_DETAIL_E (EAC_HTTP_SERVER, errDetail, EACHTTP_INSUFFICIENT_SYSTEM_RESOURCES, e.expMsg.c_str ());
        }
        else if (e.errID == ncTShareMgntError::NCT_USER_NOT_ACTIVATE) {
            THROW_E(EAC_HTTP_SERVER, EACHTTP_USER_NOT_ACTIVATE, e.expMsg.c_str());
        }
        else if (e.errID == ncTShareMgntError::NCT_OTP_WRONG) {
            THROW_E(EAC_HTTP_SERVER, EACHTTP_OTP_WRONG, e.expMsg.c_str());
        }
        else if (e.errID == ncTShareMgntError::NCT_OTP_TIMEOUT) {
            THROW_E(EAC_HTTP_SERVER, EACHTTP_OTP_TIMEOUT, e.expMsg.c_str());
        }
        else if (e.errID == ncTShareMgntError::NCT_OTP_TOO_MANY_WRONG_TIME) {
            THROW_E(EAC_HTTP_SERVER, EACHTTP_OTP_TOO_MANY_WRONG_TIME, e.expMsg.c_str());
        }
        else if (e.errID == ncTShareMgntError::NCT_CHECK_VCODE_MORE_THAN_THE_LIMIT) {
            THROW_E(EAC_HTTP_SERVER, EACHTTP_CHECK_VCODE_MORE_THAN_THE_LIMIT, e.expMsg.c_str());
        }
        else if(e.errID == ncTShareMgntError::NCT_MFA_OTP_SERVER_ERROR){
            THROW_E(EAC_HTTP_SERVER, EACHTTP_MFA_OTP_SERVER_ERROR, e.expMsg.c_str ());
        }
        else if(e.errID == ncTShareMgntError::NCT_THIRD_PLUGIN_INTER_ERROR){
            THROW_E(EAC_HTTP_SERVER, EACHTTP_THIRD_PLUGIN_INTER_ERROR, e.expMsg.c_str ());
        }
        else if(e.errID == ncTShareMgntError::NCT_MFA_CONFIG_ERROR){
            THROW_E(EAC_HTTP_SERVER, EACHTTP_MFA_CONFIG_ERROR, e.expMsg.c_str ());
        }
        else if(e.errID == ncTShareMgntError::NCT_FAILED_THIRD_CONFIG){
            THROW_E(EAC_HTTP_SERVER, EACHTTP_FAILED_THIRD_CONFIG, e.expMsg.c_str ());
        }
        else {
            THROW_HTTP_DETAIL_E (EAC_HTTP_SERVER, errDetail, USER_LOGIN_ERROR, e.expMsg.c_str ());
        }
    }
    catch (TException & e) {
        printf(e.what ());
        THROW_E (EAC_HTTP_SERVER, USER_LOGIN_ERROR, e.what ());
    }
#endif
}

void ncEACHttpServerUtil::ClientLogin (string& retUserId, const string& account, const string& password, const ncTUserLoginOption& option)
{
#ifndef __UT__
    nsresult ret;
    nsCOMPtr<authenticationInterface> authentication = do_CreateInstance (AUTHENTICATION_CONTRACTID, &ret);
    if (NS_FAILED (ret)) {
        THROW_E (EAC_HTTP_SERVER, EACHTTP_SERVER_INIT_AUTHENTICATION_ERR,
            _T("Failed to create authentication instance: 0x%x"), ret);
    }

    JSON::Value response;
    int code = authentication->ClientLogin (account, password, option, response);
    if (code == brpc::HTTP_STATUS_OK) {
        // 客户端认证成功
        SystemLog::getInstance ()->log (__FILE__, __LINE__, INFORMATION_LOG_TYPE, "Debug: statusCode = %d, userID = %s", code, response["user_id"].s().c_str());
        retUserId = response["user_id"].s();
    } else if (code == brpc::HTTP_STATUS_BAD_REQUEST) {
        string message = response["message"].s();
        SystemLog::getInstance ()->log (__FILE__, __LINE__, INFORMATION_LOG_TYPE, "Debug: statusCode = %d, message = %s", code, message.c_str());
        THROW_E (EAC_HTTP_SERVER, EACHTTP_INVALID_PARAMETER, message.c_str ());
    } else if (code == brpc::HTTP_STATUS_UNAUTHORIZED) {
        // 客户端认证失败，需要分析具体的错误原因。
        int errID = response["code"].i() % (int)1e6;
        string message = response["message"].s();

        // 获取详细错误信息
        JSON::Value errDetail;
        if(response.o().find("detail") != response.o().end()) {
            errDetail = response["detail"];
        }

        if (errID == ncTShareMgntError::NCT_INVALID_ACCOUNT_OR_PASSWORD) {
            THROW_HTTP_DETAIL_E (EAC_HTTP_SERVER, errDetail, EACHTTP_INVALID_USER_OR_PASSWORD, message.c_str ());
        }
        else if (errID == ncTShareMgntError::NCT_USER_DISABLED) {
            THROW_HTTP_DETAIL_E (EAC_HTTP_SERVER, errDetail, EACHTTP_USER_DISABLED, message.c_str ());
        }
        else if (errID == ncTShareMgntError::NCT_THIRD_PARTY_AUTH_NOT_OPEN) {
            THROW_E (EAC_HTTP_SERVER, EACHTTP_THIRD_AUTH_NOT_OPEN, message.c_str ());
        }
        else if (errID == ncTShareMgntError::NCT_NEED_THIRD_OAUTH) {
            THROW_E (EAC_HTTP_SERVER, EACHTTP_OAUTH_NEEDED, message.c_str ());
        }
        else if (errID == ncTShareMgntError::NCT_FORBIDDEN_LOGIN) {
            THROW_HTTP_DETAIL_E (EAC_HTTP_SERVER, errDetail, EACHTTP_FORBIDDEN_LOGIN, message.c_str ());
        }
        else if (errID == ncTShareMgntError::NCT_DOMAIN_NOT_EXIST) {
            THROW_E (EAC_HTTP_SERVER, EACHTTP_DOMAIN_NOT_EXIST, message.c_str ());
        }
        else if (errID == ncTShareMgntError::NCT_DOMAIN_DISABLE) {
            THROW_E (EAC_HTTP_SERVER, EACHTTP_DOMAIN_DISABLE, message.c_str ());
        }
        else if (errID == ncTShareMgntError::NCT_DOMAIN_USER_NOT_EXIST) {
            THROW_E (EAC_HTTP_SERVER, EACHTTP_DOMAIN_USER_NOT_EXIST, message.c_str ());
        }
        else if(errID == ncTShareMgntError::NCT_DOMAIN_SERVER_UNAVAILABLE){
            THROW_E(EAC_HTTP_SERVER, EACHTTP_DOMAIN_SERVER_UNAVAILABLE, message.c_str ());
        }
        else if(errID == ncTShareMgntError::NCT_PRODUCT_NOT_AUTHORIZED){
            THROW_HTTP_DETAIL_E (EAC_HTTP_SERVER, errDetail, EACHTTP_USER_NOT_AUTORIZED, message.c_str ());
        }
        else if(errID == ncTShareMgntError::NCT_PRODUCT_HAS_EXPIRED){
            THROW_HTTP_DETAIL_E (EAC_HTTP_SERVER, errDetail, EACHTTP_USER_HAS_EXPIRED, message.c_str ());
        }
        else if(errID == ncTShareMgntError::NCT_PASSWORD_EXPIRE){
            THROW_HTTP_DETAIL_E (EAC_HTTP_SERVER, errDetail, EACHTTP_PASSWORD_EXPIRE, message.c_str ());
        }
        else if(errID == ncTShareMgntError::NCT_PASSWORD_NOT_SAFE){
            THROW_HTTP_DETAIL_E (EAC_HTTP_SERVER, errDetail, EACHTTP_PASSWORD_NOT_SAFE, message.c_str ());
        }
        else if(errID == ncTShareMgntError::NCT_PASSWORD_IS_INITIAL){
            THROW_HTTP_DETAIL_E (EAC_HTTP_SERVER, errDetail, EACHTTP_PASSWORD_IS_INITIAL, message.c_str ());
        }
        else if(errID == ncTShareMgntError::NCT_PWD_FIRST_FAILED){
            THROW_E(EAC_HTTP_SERVER, EACHTTP_PWD_FAILED_FIRSTLY, message.c_str ());
        }
        else if(errID == ncTShareMgntError::NCT_PWD_SECOND_FAILED){
            THROW_E(EAC_HTTP_SERVER, EACHTTP_PWD_FAILED_SECONDLY, message.c_str ());
        }
        else if(errID == ncTShareMgntError::NCT_PWD_THIRD_FAILED){
            JSON::Value detailJson;
            detailJson["remainlockTime"] = errDetail["remainlockTime"];
            detailJson["isShowStatus"] = errDetail["isShowStatus"];
            THROW_HTTP_DETAIL_E(EAC_HTTP_SERVER, detailJson, EACHTTP_PWD_FAILED_THIRDLY, _T(message.c_str ()), to_string(errDetail["remainlockTime"].i ()).c_str());
        }
        else if(errID == ncTShareMgntError::NCT_ACCOUNT_LOCKED){
            // parse lock time
            JSON::Value detailJson;
            detailJson["remainlockTime"] = errDetail["remainlockTime"];
            detailJson["isShowStatus"] = errDetail["isShowStatus"];
            THROW_HTTP_DETAIL_E(EAC_HTTP_SERVER, detailJson, EACHTTP_ACCOUNT_LOCKED, _T(message.c_str ()), to_string(errDetail["remainlockTime"].i ()).c_str());
        }
        else if(errID == ncTShareMgntError::NCT_CANNOT_CONNECT_THID_PARTY_SERVER){
            THROW_E(EAC_HTTP_SERVER, EACHTTP_THIRD_SERVER_UNAVALIABLE, message.c_str ());
        }
        else if(errID == ncTShareMgntError::NCT_USER_NOT_EXIST){
            THROW_E(EAC_HTTP_SERVER, EACHTTP_USER_NOT_EXIST, message.c_str ());
        }
        else if(errID == ncTShareMgntError::NCT_WRONG_PASSWORD){
            THROW_E(EAC_HTTP_SERVER, EACHTTP_WRONG_PASSWORD, message.c_str ());
        }
        else if(errID == ncTShareMgntError::NCT_CONTROLED_PASSWORD_EXPIRE){
            THROW_E(EAC_HTTP_SERVER, EACHTTP_CONTROLED_PASSWORD_EXPIRE, message.c_str ());
        }
        else if(errID == ncTShareMgntError::NCT_CANNOT_LOGIN_SLAVE_SITE){
            THROW_HTTP_DETAIL_E (EAC_HTTP_SERVER, errDetail, EACHTTP_CANNOT_LOGIN_SLAVE_SITE, message.c_str ());
        }
        else if(errID == ncTShareMgntError::NCT_CHECK_VCODE_IS_NULL){
            THROW_HTTP_DETAIL_E(EAC_HTTP_SERVER, errDetail, EACHTTP_CHECK_VCODE_IS_NULL, message.c_str ());
        }
        else if(errID == ncTShareMgntError::NCT_CHECK_VCODE_IS_TIMEOUT){
            THROW_HTTP_DETAIL_E(EAC_HTTP_SERVER, errDetail, EACHTTP_CHECK_VCODE_IS_TIMEOUT, message.c_str ());
        }
        else if(errID == ncTShareMgntError::NCT_CHECK_VCODE_IS_WRONG){
            THROW_HTTP_DETAIL_E(EAC_HTTP_SERVER, errDetail, EACHTTP_CHECK_VCODE_IS_WRONG, message.c_str ());
        }
        else if (errID == ncTShareMgntError::NCT_INSUFFICIENT_SYSTEM_RESOURCES) {
            THROW_HTTP_DETAIL_E (EAC_HTTP_SERVER, errDetail, EACHTTP_INSUFFICIENT_SYSTEM_RESOURCES, message.c_str ());
        }
        else if (errID == ncTShareMgntError::NCT_USER_NOT_ACTIVATE) {
            THROW_E(EAC_HTTP_SERVER, EACHTTP_USER_NOT_ACTIVATE, message.c_str());
        }
        else if (errID == ncTShareMgntError::NCT_OTP_WRONG) {
            THROW_E(EAC_HTTP_SERVER, EACHTTP_OTP_WRONG, message.c_str());
        }
        else if (errID == ncTShareMgntError::NCT_OTP_TIMEOUT) {
            THROW_E(EAC_HTTP_SERVER, EACHTTP_OTP_TIMEOUT, message.c_str());
        }
        else if (errID == ncTShareMgntError::NCT_OTP_TOO_MANY_WRONG_TIME) {
            THROW_E(EAC_HTTP_SERVER, EACHTTP_OTP_TOO_MANY_WRONG_TIME, message.c_str());
        }
        else if (errID == ncTShareMgntError::NCT_CHECK_VCODE_MORE_THAN_THE_LIMIT) {
            THROW_E(EAC_HTTP_SERVER, EACHTTP_CHECK_VCODE_MORE_THAN_THE_LIMIT, message.c_str());
        }
        else if(errID == ncTShareMgntError::NCT_MFA_OTP_SERVER_ERROR){
            THROW_E(EAC_HTTP_SERVER, EACHTTP_MFA_OTP_SERVER_ERROR, message.c_str ());
        }
        else if(errID == ncTShareMgntError::NCT_THIRD_PLUGIN_INTER_ERROR){
            THROW_E(EAC_HTTP_SERVER, EACHTTP_THIRD_PLUGIN_INTER_ERROR, message.c_str ());
        }
        else if(errID == ncTShareMgntError::NCT_MFA_CONFIG_ERROR){
            THROW_E(EAC_HTTP_SERVER, EACHTTP_MFA_CONFIG_ERROR, message.c_str ());
        }
        else if(errID == ncTShareMgntError::NCT_FAILED_THIRD_CONFIG){
            THROW_E(EAC_HTTP_SERVER, EACHTTP_FAILED_THIRD_CONFIG, message.c_str ());
        }
        else {
            THROW_HTTP_DETAIL_E (EAC_HTTP_SERVER, errDetail, USER_LOGIN_ERROR, message.c_str ());
        }
    } else {
        // 未知错误
        if (code == 0) {
            THROW_E (EAC_HTTP_SERVER, EACHTTP_SERVER_UNKNOWN_ERR, _T("Could not connect to server"));
        } else {
            THROW_E (EAC_HTTP_SERVER, EACHTTP_SERVER_UNKNOWN_ERR, _T("Code:%d. Cause:%s."), code, response["cause"].s().c_str());
        }
    }

    return;
#endif
}

bool ncEACHttpServerUtil::Usrm_GetADSSOStatus ()
{
    try {
        ncThriftClient<ncTShareMgntClient> shareMgntClient (EacServiceAccessConfig::getInstance()->sharemgntHost, EacServiceAccessConfig::getInstance()->sharemgntPort);
        return shareMgntClient->Usrm_GetADSSOStatus();
    }
    catch (ncTException & e) {
        THROW_E (EAC_HTTP_SERVER, EACHTTP_GET_AD_SSO_STATUS_ERROR,
            LOAD_STRING (_T("IDS_EACHTTP_GET_AD_SSO_STATUS_ERROR")), e.expMsg.c_str ());
    }
    catch (TException & e) {
        THROW_E (EAC_HTTP_SERVER, EACHTTP_GET_AD_SSO_STATUS_ERROR,
            LOAD_STRING (_T("IDS_EACHTTP_GET_AD_SSO_STATUS_ERROR")), e.what ());
    }
}

ncTThirdPartyAuthConf ncEACHttpServerUtil::Usrm_GetThirdPartyAuth ()
{
    try {
        ncThriftClient<ncTShareMgntClient> shareMgntClient (EacServiceAccessConfig::getInstance()->sharemgntHost, EacServiceAccessConfig::getInstance()->sharemgntPort);

        ncTThirdPartyAuthConf retConf;
        shareMgntClient->Usrm_GetThirdPartyAuth(retConf);

        return retConf;
    }
    catch (ncTException & e) {
        THROW_E (EAC_HTTP_SERVER, EACHTTP_GET_THIRD_PARTY_AUTH_ERROR,
            LOAD_STRING (_T("IDS_EACHTTP_GET_THIRD_PARTY_AUTH_ERROR")), e.expMsg.c_str ());
    }
    catch (TException & e) {
        THROW_E (EAC_HTTP_SERVER, EACHTTP_GET_THIRD_PARTY_AUTH_ERROR,
            LOAD_STRING (_T("IDS_EACHTTP_GET_THIRD_PARTY_AUTH_ERROR")), e.what ());
    }
}

bool ncEACHttpServerUtil::Usrm_ValidateSecurityDevice (string& params)
{
    try {
        ncThriftClient<ncTShareMgntClient> shareMgntClient (EacServiceAccessConfig::getInstance()->sharemgntHost, EacServiceAccessConfig::getInstance()->sharemgntPort);
        return shareMgntClient->Usrm_ValidateSecurityDevice(params);
    }
    catch (ncTException & e) {
        THROW_E (EAC_HTTP_SERVER, EACHTTP_VALIDATE_THIRD_TOKEN_ERROR, "%s", e.expMsg.c_str ());
    }
    catch (TException & e) {
        THROW_E (EAC_HTTP_SERVER, EACHTTP_VALIDATE_THIRD_TOKEN_ERROR, "%s", e.what ());
    }
}

void ncEACHttpServerUtil::Log (brpc::Controller* cntl, const String& userId, ncTokenVisitorType typ, ncTLogType logType, ncTLogLevel level,
                                int opType, const String& msg, const String& exmsg, bool logForwardedIp /* false */)
{
#ifndef __UT__

    nsresult ret;
    // 初始化nsq
    nsCOMPtr<authenticationInterface> authentication = do_CreateInstance (AUTHENTICATION_CONTRACTID, &ret);
    if (NS_FAILED (ret)) {
        THROW_E (EAC_HTTP_SERVER, EACHTTP_SERVER_INIT_AUTHENTICATION_ERR,
            _T("Failed to create authentication instance: 0x%x"), ret);
    }

    // 构建msg
    String macAddress;
    String userAgent;
    String ip = GetForwardedIp (cntl);
    ncHttpGetHeader (cntl, "X-Request-MAC", macAddress);
    ncHttpGetHeader (cntl, "User-Agent", userAgent);

    authentication->AuditLog(userId, typ, logType, level, opType, msg, exmsg, ip, macAddress, userAgent);

#endif
}

void ncEACHttpServerUtil::LoginLog (const String& userId, const String& udid, ACSClientType clientType, const String& ip)
{
#ifndef __UT__

    nsresult ret;
    nsCOMPtr<userManagementInterface> userManager = do_CreateInstance (USER_MANAGEMENT_CONTRACTID, &ret);
    if (NS_FAILED (ret)) {
        THROW_E (EAC_HTTP_SERVER, EACHTTP_SERVER_INIT_USER_MANAGEMENT_ERR,
            _T("Failed to create usermanagement instance: 0x%x"), ret);
    }

    nsCOMPtr<stdLogInterface> stdLog = do_CreateInstance (STD_LOG_CONTRACTID, &ret);
    if (NS_FAILED (ret)) {
        THROW_E (EAC_HTTP_SERVER, EACHTTP_XPCOM_CREATE_INSTANCE_FAILED,
            _T("Failed to create stdLog instance: 0x%x"), ret);
    }

    // 获取用户信息
    UserInfo userInfo;
    userManager->GetUserInfo (userId, userInfo);

    ncOperator actor;
    actor.clientType = ncClientType(clientType);
    actor.id = userId;
    actor.name = userInfo.name;
    actor.ip  = ip;
    actor.udid = udid;
    actor.departmentPaths = std::move(userInfo.vecDepartIDs);
    actor.departmentNames = std::move(userInfo.vecDepartNames);

    stdLog->Log (actor);

#endif
}

String ncEACHttpServerUtil::Usrm_LoginConsoleByThirdPartyNew(String& params)
{
    try {
        ncThriftClient<ncTShareMgntClient> shareMgntClient (EacServiceAccessConfig::getInstance()->sharemgntHost, EacServiceAccessConfig::getInstance()->sharemgntPort);
        string account;
        string strParams = params.getCStr ();
        shareMgntClient->Usrm_LoginConsoleByThirdPartyNew(account, strParams);

        return account.c_str();
    }
    catch (ncTException & e) {
        if(e.errID == ncTShareMgntError::NCT_USER_NOT_EXIST){
            THROW_E(EAC_HTTP_SERVER, EACHTTP_NOT_IMPORT_TO_ANYSHARE, e.expMsg.c_str ());
        }
        else if(e.errID == ncTShareMgntError::NCT_FORBIDDEN_LOGIN)
        {
            THROW_E (EAC_HTTP_SERVER, EACHTTP_NORMAL_FORBIDDEN_LOGIN, e.expMsg.c_str ());
        }
        else if(e.errID == ncTShareMgntError::NCT_USER_DISABLED)
        {
            THROW_E (EAC_HTTP_SERVER, EACHTTP_USER_DISABLED, e.expMsg.c_str ());
        }
        else
        {
            THROW_E (EAC_HTTP_SERVER, EACHTTP_VALIDATE_THIRD_TOKEN_ERROR, "%s", e.expMsg.c_str ());
        }
    }
    catch (TException & e) {
        THROW_E (EAC_HTTP_SERVER, EACHTTP_VALIDATE_THIRD_TOKEN_ERROR, "%s", e.what ());
    }
}

String ncEACHttpServerUtil::Usrm_ValidateThirdParty (string& params)
{
    try {
        ncThriftClient<ncTShareMgntClient> shareMgntClient (EacServiceAccessConfig::getInstance()->sharemgntHost, EacServiceAccessConfig::getInstance()->sharemgntPort);
        string account;
        shareMgntClient->Usrm_ValidateThirdParty(account, params);

        return account.c_str();
    }
    catch (ncTException & e) {
        if (e.errID == ncTShareMgntError::NCT_FAILED_THIRD_CONFIG) {
            THROW_E (EAC_HTTP_SERVER, EACHTTP_FAILED_THIRD_CONFIG, e.expMsg.c_str ());
        } else if (e.errID == ncTShareMgntError::NCT_INVALID_ACCOUNT_OR_PASSWORD) {
            THROW_E (EAC_HTTP_SERVER, EACHTTP_INVALID_USER_OR_PASSWORD, e.expMsg.c_str ());
        } else {
            THROW_E (EAC_HTTP_SERVER, EACHTTP_VALIDATE_THIRD_TOKEN_ERROR, "%s", e.expMsg.c_str ());
        }
    }
    catch (TException & e) {
        THROW_E (EAC_HTTP_SERVER, EACHTTP_VALIDATE_THIRD_TOKEN_ERROR, "%s", e.what ());
    }
}

String ncEACHttpServerUtil::IsFileToStr(bool isFile)
{
    if(isFile) {
        return LOAD_STRING("IDS_FILE");
    }
    else {
        return LOAD_STRING("IDS_DIR");
    }
}


int ncEACHttpServerUtil::GetCustomConfigOfInt64(const string& key)
{
    try {
        ncThriftClient<ncTShareMgntClient> shareMgntClient (EacServiceAccessConfig::getInstance()->sharemgntHost, EacServiceAccessConfig::getInstance()->sharemgntPort);
        return shareMgntClient->GetCustomConfigOfInt64(key);
    }
    catch (Exception& e){
        THROW_E (EAC_HTTP_SERVER, GET_TAG_MAX_NUM_ERROR, _T("%s"), e.toString ().getCStr ());
    }
    catch (ncTException & e) {
        THROW_E (EAC_HTTP_SERVER, GET_TAG_MAX_NUM_ERROR, _T("%s"), e.expMsg.c_str ());
    }
    catch (TException & e) {
        THROW_E (EAC_HTTP_SERVER, GET_TAG_MAX_NUM_ERROR, _T("%s"), e.what ());
    }
}

bool ncEACHttpServerUtil::GetLoginStrategyStatus ()
{
    try {
        ncThriftClient<ncTShareMgntClient> shareMgntClient (EacServiceAccessConfig::getInstance()->sharemgntHost, EacServiceAccessConfig::getInstance()->sharemgntPort);
        return shareMgntClient->GetLoginStrategyStatus();
    }
    catch (ncTException & e) {
        THROW_E (EAC_HTTP_SERVER, EACHTTP_GET_LOGIN_STRATEGY_STATUS_ERROR,
            LOAD_STRING (_T("IDS_EACHTTP_GET_LOGIN_STRATEGY_STATUS_ERROR")), e.expMsg.c_str ());
    }
    catch (TException & e) {
        THROW_E (EAC_HTTP_SERVER, EACHTTP_GET_LOGIN_STRATEGY_STATUS_ERROR,
            LOAD_STRING (_T("IDS_EACHTTP_GET_LOGIN_STRATEGY_STATUS_ERROR")), e.what ());
    }
}
//调用sharmgnt thrift接口发送邮件
void ncEACHttpServerUtil::SendMail(vector<string>& mailto, string& subject, string& content)
{
    try {
        ncThriftClient<ncTShareMgntClient> shareMgntClient (EacServiceAccessConfig::getInstance()->sharemgntHost, EacServiceAccessConfig::getInstance()->sharemgntPort);
        shareMgntClient->SMTP_SendEmail(mailto, subject, content);
    }
    catch (ncTException & e) {
        if ((e.errID == ncTShareMgntError::NCT_SMTP_RECIPIENT_MAIL_ILLEGAL) || (e.errID == ncTShareMgntError::NCT_INVALID_EMAIL)) {
            THROW_E (EAC_HTTP_SERVER, EACHTTP_SMTP_RECIPIENT_MAIL_ILLEGAL,
            LOAD_STRING (_T("IDS_EACHTTP_SMTP_RECIPIENT_MAIL_ILLEGAL")), e.expMsg.c_str ());
        }
        else if (e.errID == ncTShareMgntError::NCT_SMTP_SERVER_NOT_SET) {
            THROW_E (EAC_HTTP_SERVER, EACHTTP_SMTP_SERVER_NOT_SET,
            LOAD_STRING (_T("IDS_EACHTTP_SMTP_SERVER_NOT_SET")), e.expMsg.c_str ());
        }
        else if ((e.errID == ncTShareMgntError::NCT_SMTP_SERVER_NOT_AVAILABLE) ||
                 (e.errID == ncTShareMgntError::NCT_SMTP_LOGIN_FAILED) ||
                 (e.errID == ncTShareMgntError::NCT_SMTP_AUTHENTICATION_METHOD_NOT_FOUND)) {
            THROW_E (EAC_HTTP_SERVER, EACHTTP_SMTP_SERVER_NOT_AVAILABLE,
            LOAD_STRING (_T("IDS_EACHTTP_SMTP_SERVER_NOT_AVAILABLE")), e.expMsg.c_str ());
        }
        else if (e.errID == ncTShareMgntError::NCT_SMTP_SEND_FAILED) {
            THROW_E (EAC_HTTP_SERVER, EACHTTP_SMTP_SEND_FAILED,
            LOAD_STRING (_T("IDS_EACHTTP_SMTP_SEND_FAILED")), e.expMsg.c_str ());
        }
        else {
            THROW_E (EAC_HTTP_SERVER, EACHTTP_SERVER_UNKNOWN_ERR, e.expMsg.c_str ());
        }
    }
}

bool ncEACHttpServerUtil::OEM_GetConfigByOption(string section, string option)
{
    try {
        ncThriftClient<ncTShareMgntClient> shareMgntClient (EacServiceAccessConfig::getInstance()->sharemgntHost, EacServiceAccessConfig::getInstance()->sharemgntPort);
        string _return;
        shareMgntClient->OEM_GetConfigByOption(_return, section, option);
        return (_return == "true");
    }
    catch (ncTException & e) {
        THROW_E (EAC_HTTP_SERVER, EACHTTP_GET_OEM_CONFIG_ERROR,
            LOAD_STRING (_T("IDS_EACHTTP_GET_OEM_CONFIG_ERROR")), e.expMsg.c_str ());
    }
    catch (TException & e) {
        THROW_E (EAC_HTTP_SERVER, EACHTTP_GET_OEM_CONFIG_ERROR,
            LOAD_STRING (_T("IDS_EACHTTP_GET_OEM_CONFIG_ERROR")), e.what ());
    }
}

String ncEACHttpServerUtil::GenerateGroupStr (const vector<String>& strs)
{
    String groupStr;
    for (size_t i = 0; i < strs.size (); ++i) {
        groupStr.append ("\"", 1);
        groupStr.append (strs[i]);
        groupStr.append ("\"", 1);

        if (i != (strs.size () -1)) {
            groupStr.append (",", 1);
        }
    }

    return groupStr;
}


String ncEACHttpServerUtil::EndTimeToStr (int64 ticks)
{
    if(ticks == -1) {
        return LOAD_STRING("IDS_UNLIMITED");
    }
    else {
        Date d(ticks);
        String str = DateFormat::getLocalInstance ()-> format(&d, "yyyy/MM/dd HH:mm");
        return str;
    }
}

String ncEACHttpServerUtil::AccessorTypeToResStr (int accessorType)
{
    if(accessorType == ACS_USER) {
        return LOAD_STRING("IDS_ACS_USER");
    }
    else if(accessorType == ACS_DEPARTMENT) {
        return LOAD_STRING("IDS_ACS_DEPARTMENT");
    }
    else if(accessorType == ACS_CONTACTOR) {
        return LOAD_STRING("IDS_ACS_CONTACTOR");
    }
    else if(accessorType == ACS_GROUP) {
        return LOAD_STRING("IDS_ACS_GROUP");
    }
    else {
        return _T("unknown");
    }
}

String ncEACHttpServerUtil::ConvertPermToStr (int permValue)
{
    String permStr;
    permStr += _T(" ");
    if (permValue & ncAtomPermValue::ACS_AP_DISPLAY) {
        permStr += LOAD_STRING (_T("IDS_DISPLAY"));
        permStr += _T("/");
    }
    if (permValue & ncAtomPermValue::ACS_AP_READ) {
        permStr += LOAD_STRING (_T("IDS_READ"));
        permStr += _T("/");
    }
    if (permValue & ncAtomPermValue::ACS_AP_CREATE) {
        permStr += LOAD_STRING (_T("IDS_CREATE"));
        permStr += _T("/");
    }
    if (permValue & ncAtomPermValue::ACS_AP_EDIT) {
        permStr += LOAD_STRING (_T("IDS_EDIT"));
        permStr += _T("/");
    }
    if (permValue & ncAtomPermValue::ACS_AP_DELETE) {
        permStr += LOAD_STRING (_T("IDS_DELETE"));
        permStr += _T("/");
    }

    if(!permStr.isEmpty()) {
        permStr.remove(permStr.getLength() - 1, 1);
    }

    return permStr;
}

String ncEACHttpServerUtil::GetFormatCSFLevelStr (const int& CSFLevel)
{
    map<string, int32_t> csflevels;
    GetCSFLevels (csflevels);
    String CSFLevelStr;
    for (auto iter = csflevels.begin(); iter != csflevels.end(); ++iter) {
        if (CSFLevel == iter->second) {
            CSFLevelStr.format("(%s)", iter->first.c_str ());
        }
    }
    return CSFLevelStr;
}

void ncEACHttpServerUtil::Usrm_UserLoginByNTLMV1 (ncTNTLMResponse& ntlmResp, const string& account, const string& challenge, const string& password)
{
    try {
        ncThriftClient<ncTShareMgntClient> shareMgntClient (EacServiceAccessConfig::getInstance()->sharemgntHost, EacServiceAccessConfig::getInstance()->sharemgntPort);
        shareMgntClient->Usrm_UserLoginByNTLMV1 (ntlmResp, account, challenge, password);
    }
    catch (ncTException & e) {
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
        if (e.errID == ncTShareMgntError::NCT_INVALID_ACCOUNT_OR_PASSWORD) {
            THROW_E (EAC_HTTP_SERVER, EACHTTP_INVALID_USER_OR_PASSWORD, e.expMsg.c_str ());
        }
        else if (e.errID == ncTShareMgntError::NCT_USER_DISABLED) {
            THROW_E (EAC_HTTP_SERVER, EACHTTP_USER_DISABLED, e.expMsg.c_str ());
        }
        else if (e.errID == ncTShareMgntError::NCT_THIRD_PARTY_AUTH_NOT_OPEN) {
            THROW_E (EAC_HTTP_SERVER, EACHTTP_THIRD_AUTH_NOT_OPEN, e.expMsg.c_str ());
        }
        else if (e.errID == ncTShareMgntError::NCT_NEED_THIRD_OAUTH) {
            THROW_E (EAC_HTTP_SERVER, EACHTTP_OAUTH_NEEDED, e.expMsg.c_str ());
        }
        else if (e.errID == ncTShareMgntError::NCT_FORBIDDEN_LOGIN) {
            THROW_E (EAC_HTTP_SERVER, EACHTTP_FORBIDDEN_LOGIN, e.expMsg.c_str ());
        }
        else if (e.errID == ncTShareMgntError::NCT_DOMAIN_NOT_EXIST) {
            THROW_E (EAC_HTTP_SERVER, EACHTTP_DOMAIN_NOT_EXIST, e.expMsg.c_str ());
        }
        else if (e.errID == ncTShareMgntError::NCT_DOMAIN_DISABLE) {
            THROW_E (EAC_HTTP_SERVER, EACHTTP_DOMAIN_DISABLE, e.expMsg.c_str ());
        }
        else if (e.errID == ncTShareMgntError::NCT_DOMAIN_USER_NOT_EXIST) {
            THROW_E (EAC_HTTP_SERVER, EACHTTP_DOMAIN_USER_NOT_EXIST, e.expMsg.c_str ());
        }
        else if(e.errID == ncTShareMgntError::NCT_DOMAIN_SERVER_UNAVAILABLE){
            THROW_E(EAC_HTTP_SERVER, EACHTTP_DOMAIN_SERVER_UNAVAILABLE, e.expMsg.c_str ());
        }
        else if(e.errID == ncTShareMgntError::NCT_PRODUCT_NOT_AUTHORIZED){
            THROW_E(EAC_HTTP_SERVER, EACHTTP_USER_NOT_AUTORIZED, e.expMsg.c_str ());
        }
        else if(e.errID == ncTShareMgntError::NCT_PRODUCT_HAS_EXPIRED){
            THROW_E(EAC_HTTP_SERVER, EACHTTP_USER_HAS_EXPIRED, e.expMsg.c_str ());
        }
        else if(e.errID == ncTShareMgntError::NCT_PASSWORD_EXPIRE){
            THROW_E(EAC_HTTP_SERVER, EACHTTP_PASSWORD_EXPIRE, e.expMsg.c_str ());
        }
        else if(e.errID == ncTShareMgntError::NCT_PASSWORD_NOT_SAFE){
            THROW_E(EAC_HTTP_SERVER, EACHTTP_PASSWORD_NOT_SAFE, e.expMsg.c_str ());
        }
        else if(e.errID == ncTShareMgntError::NCT_PASSWORD_IS_INITIAL){
            THROW_E(EAC_HTTP_SERVER, EACHTTP_PASSWORD_IS_INITIAL, e.expMsg.c_str ());
        }
        else if(e.errID == ncTShareMgntError::NCT_PWD_FIRST_FAILED){
            THROW_E(EAC_HTTP_SERVER, EACHTTP_PWD_FAILED_FIRSTLY, e.expMsg.c_str ());
        }
        else if(e.errID == ncTShareMgntError::NCT_PWD_SECOND_FAILED){
            THROW_E(EAC_HTTP_SERVER, EACHTTP_PWD_FAILED_SECONDLY, e.expMsg.c_str ());
        }
        else if(e.errID == ncTShareMgntError::NCT_PWD_THIRD_FAILED){
            // parse lock tim
            JSON::Value detailJson;
            detailJson["remainlockTime"] = errDetail["remainlockTime"];
            detailJson["isShowStatus"] = errDetail["isShowStatus"];
            THROW_HTTP_DETAIL_E(EAC_HTTP_SERVER, detailJson, EACHTTP_PWD_FAILED_THIRDLY, e.expMsg.c_str ());
        }
        else if(e.errID == ncTShareMgntError::NCT_ACCOUNT_LOCKED){
            // parse lock tim
            JSON::Value detailJson;
            detailJson["remainlockTime"] = errDetail["remainlockTime"];
            detailJson["isShowStatus"] = errDetail["isShowStatus"];
            THROW_HTTP_DETAIL_E(EAC_HTTP_SERVER, detailJson, EACHTTP_ACCOUNT_LOCKED, e.expMsg.c_str ());
        }
        else if(e.errID == ncTShareMgntError::NCT_CANNOT_CONNECT_THID_PARTY_SERVER){
            THROW_E(EAC_HTTP_SERVER, EACHTTP_THIRD_SERVER_UNAVALIABLE, e.expMsg.c_str ());
        }
        else if(e.errID == ncTShareMgntError::NCT_USER_NOT_EXIST){
            THROW_E(EAC_HTTP_SERVER, EACHTTP_USER_NOT_EXIST, e.expMsg.c_str ());
        }
        else if(e.errID == ncTShareMgntError::NCT_WRONG_PASSWORD){
            THROW_E(EAC_HTTP_SERVER, EACHTTP_WRONG_PASSWORD, e.expMsg.c_str ());
        }
        else if(e.errID == ncTShareMgntError::NCT_CONTROLED_PASSWORD_EXPIRE){
            THROW_E(EAC_HTTP_SERVER, EACHTTP_CONTROLED_PASSWORD_EXPIRE, e.expMsg.c_str ());
        }
        else if(e.errID == ncTShareMgntError::NCT_CANNOT_LOGIN_SLAVE_SITE){
            THROW_E(EAC_HTTP_SERVER, EACHTTP_CANNOT_LOGIN_SLAVE_SITE, e.expMsg.c_str ());
        }
        else {
            THROW_E (EAC_HTTP_SERVER, USER_LOGIN_ERROR, e.expMsg.c_str ());
        }
    }
    catch (TException & e) {
        printf(e.what ());
        THROW_E (EAC_HTTP_SERVER, USER_LOGIN_ERROR, e.what ());
    }
}

void ncEACHttpServerUtil::Usrm_UserLoginByNTLMV2 (ncTNTLMResponse& ntlmResp, const string& account, const string& domain, const string& challenge, const string& password)
{
    try {
        ncThriftClient<ncTShareMgntClient> shareMgntClient (EacServiceAccessConfig::getInstance()->sharemgntHost, EacServiceAccessConfig::getInstance()->sharemgntPort);
        shareMgntClient->Usrm_UserLoginByNTLMV2 (ntlmResp, account, domain, challenge, password);
    }
    catch (ncTException & e) {
        if (e.errID == ncTShareMgntError::NCT_INVALID_ACCOUNT_OR_PASSWORD) {
            THROW_E (EAC_HTTP_SERVER, EACHTTP_INVALID_USER_OR_PASSWORD, e.expMsg.c_str ());
        }
        else if (e.errID == ncTShareMgntError::NCT_USER_DISABLED) {
            THROW_E (EAC_HTTP_SERVER, EACHTTP_USER_DISABLED, e.expMsg.c_str ());
        }
        else if (e.errID == ncTShareMgntError::NCT_THIRD_PARTY_AUTH_NOT_OPEN) {
            THROW_E (EAC_HTTP_SERVER, EACHTTP_THIRD_AUTH_NOT_OPEN, e.expMsg.c_str ());
        }
        else if (e.errID == ncTShareMgntError::NCT_NEED_THIRD_OAUTH) {
            THROW_E (EAC_HTTP_SERVER, EACHTTP_OAUTH_NEEDED, e.expMsg.c_str ());
        }
        else if (e.errID == ncTShareMgntError::NCT_FORBIDDEN_LOGIN) {
            THROW_E (EAC_HTTP_SERVER, EACHTTP_FORBIDDEN_LOGIN, e.expMsg.c_str ());
        }
        else if (e.errID == ncTShareMgntError::NCT_DOMAIN_NOT_EXIST) {
            THROW_E (EAC_HTTP_SERVER, EACHTTP_DOMAIN_NOT_EXIST, e.expMsg.c_str ());
        }
        else if (e.errID == ncTShareMgntError::NCT_DOMAIN_DISABLE) {
            THROW_E (EAC_HTTP_SERVER, EACHTTP_DOMAIN_DISABLE, e.expMsg.c_str ());
        }
        else if (e.errID == ncTShareMgntError::NCT_DOMAIN_USER_NOT_EXIST) {
            THROW_E (EAC_HTTP_SERVER, EACHTTP_DOMAIN_USER_NOT_EXIST, e.expMsg.c_str ());
        }
        else if(e.errID == ncTShareMgntError::NCT_DOMAIN_SERVER_UNAVAILABLE){
            THROW_E(EAC_HTTP_SERVER, EACHTTP_DOMAIN_SERVER_UNAVAILABLE, e.expMsg.c_str ());
        }
        else if(e.errID == ncTShareMgntError::NCT_PRODUCT_NOT_AUTHORIZED){
            THROW_E(EAC_HTTP_SERVER, EACHTTP_USER_NOT_AUTORIZED, e.expMsg.c_str ());
        }
        else if(e.errID == ncTShareMgntError::NCT_PRODUCT_HAS_EXPIRED){
            THROW_E(EAC_HTTP_SERVER, EACHTTP_USER_HAS_EXPIRED, e.expMsg.c_str ());
        }
        else if(e.errID == ncTShareMgntError::NCT_PASSWORD_EXPIRE){
            THROW_E(EAC_HTTP_SERVER, EACHTTP_PASSWORD_EXPIRE, e.expMsg.c_str ());
        }
        else if(e.errID == ncTShareMgntError::NCT_PASSWORD_NOT_SAFE){
            THROW_E(EAC_HTTP_SERVER, EACHTTP_PASSWORD_NOT_SAFE, e.expMsg.c_str ());
        }
        else if(e.errID == ncTShareMgntError::NCT_PASSWORD_IS_INITIAL){
            THROW_E(EAC_HTTP_SERVER, EACHTTP_PASSWORD_IS_INITIAL, e.expMsg.c_str ());
        }
        else if(e.errID == ncTShareMgntError::NCT_PWD_FIRST_FAILED){
            THROW_E(EAC_HTTP_SERVER, EACHTTP_PWD_FAILED_FIRSTLY, e.expMsg.c_str ());
        }
        else if(e.errID == ncTShareMgntError::NCT_PWD_SECOND_FAILED){
            THROW_E(EAC_HTTP_SERVER, EACHTTP_PWD_FAILED_SECONDLY, e.expMsg.c_str ());
        }
        else if(e.errID == ncTShareMgntError::NCT_PWD_THIRD_FAILED){
            JSON::Value detailJson;
            detailJson["remainlockTime"] = ParseLockTime(e.expMsg.c_str ());
            THROW_HTTP_DETAIL_E(EAC_HTTP_SERVER, detailJson, EACHTTP_PWD_FAILED_THIRDLY, e.expMsg.c_str ());
        }
        else if(e.errID == ncTShareMgntError::NCT_ACCOUNT_LOCKED){
            // parse lock tim
            JSON::Value detailJson;
            detailJson["remainlockTime"] = ParseLockTime(e.expMsg.c_str ());
            THROW_HTTP_DETAIL_E(EAC_HTTP_SERVER, detailJson, EACHTTP_ACCOUNT_LOCKED, e.expMsg.c_str ());
        }
        else if(e.errID == ncTShareMgntError::NCT_CANNOT_CONNECT_THID_PARTY_SERVER){
            THROW_E(EAC_HTTP_SERVER, EACHTTP_THIRD_SERVER_UNAVALIABLE, e.expMsg.c_str ());
        }
        else if(e.errID == ncTShareMgntError::NCT_USER_NOT_EXIST){
            THROW_E(EAC_HTTP_SERVER, EACHTTP_USER_NOT_EXIST, e.expMsg.c_str ());
        }
        else if(e.errID == ncTShareMgntError::NCT_WRONG_PASSWORD){
            THROW_E(EAC_HTTP_SERVER, EACHTTP_WRONG_PASSWORD, e.expMsg.c_str ());
        }
        else if(e.errID == ncTShareMgntError::NCT_CONTROLED_PASSWORD_EXPIRE){
            THROW_E(EAC_HTTP_SERVER, EACHTTP_CONTROLED_PASSWORD_EXPIRE, e.expMsg.c_str ());
        }
        else if(e.errID == ncTShareMgntError::NCT_CANNOT_LOGIN_SLAVE_SITE){
            THROW_E(EAC_HTTP_SERVER, EACHTTP_CANNOT_LOGIN_SLAVE_SITE, e.expMsg.c_str ());
        }
        else {
            THROW_E (EAC_HTTP_SERVER, USER_LOGIN_ERROR, e.expMsg.c_str ());
        }
    }
    catch (TException & e) {
        printf(e.what ());
        THROW_E (EAC_HTTP_SERVER, USER_LOGIN_ERROR, e.what ());
    }
}

ncTThirdPartyToolConfig ncEACHttpServerUtil::GetThirdPartyToolConfig(const string& thirdPartyToolId)
{
    ncTThirdPartyToolConfig config;
    try {
        ncThriftClient<ncTShareMgntClient> shareMgntClient (EacServiceAccessConfig::getInstance()->sharemgntHost, EacServiceAccessConfig::getInstance()->sharemgntPort);
        shareMgntClient->GetThirdPartyToolConfig(config, thirdPartyToolId);
        return config;
    }
    catch (ncTException & e) {
        THROW_E (EAC_HTTP_SERVER, GET_THIRD_PARTY_TOOL_CONFIG_ERROR, e.expMsg.c_str ());
    }
    catch (TException & e) {
        THROW_E (EAC_HTTP_SERVER, GET_THIRD_PARTY_TOOL_CONFIG_ERROR, e.what ());
    }
}

bool ncEACHttpServerUtil::Secretm_GetStatus()
{
    try {
        ncThriftClient<ncTShareMgntClient> shareMgntClient (EacServiceAccessConfig::getInstance()->sharemgntHost, EacServiceAccessConfig::getInstance()->sharemgntPort);
        return shareMgntClient->Secretm_GetStatus();
    }
    catch (ncTException & e) {
        THROW_E (EAC_HTTP_SERVER, EACHTTP_GET_SECRET_STATUS_ERROR,
            LOAD_STRING (_T("IDS_EACHTTP_GET_SECRET_STATUS_ERROR")), e.expMsg.c_str ());
    }
    catch (TException & e) {
        THROW_E (EAC_HTTP_SERVER, EACHTTP_GET_SECRET_STATUS_ERROR,
            LOAD_STRING (_T("IDS_EACHTTP_GET_SECRET_STATUS_ERROR")), e.what ());
    }
}

void ncEACHttpServerUtil::GetCSFLevels(map<string, int32_t> & csflevels)
{
    try {
        ncThriftClient<ncTShareMgntClient> shareMgntClient (EacServiceAccessConfig::getInstance()->sharemgntHost, EacServiceAccessConfig::getInstance()->sharemgntPort);
        shareMgntClient->GetCSFLevels(csflevels);
    }
    catch (ncTException & e) {
        THROW_E (EAC_HTTP_SERVER, EACHTTP_GET_CSF_LEVELS_ERROR,
            LOAD_STRING (_T("IDS_EACHTTP_GET_CSF_LEVELS_ERROR")), e.expMsg.c_str ());
    }
    catch (TException & e) {
        THROW_E (EAC_HTTP_SERVER, EACHTTP_GET_CSF_LEVELS_ERROR,
            LOAD_STRING (_T("IDS_EACHTTP_GET_CSF_LEVELS_ERROR")), e.what ());
    }

}

void ncEACHttpServerUtil::CheckUninstallPwd(const String& uninstallPwd)
{
    try {
        ncThriftClient<ncTShareMgntClient> shareMgntClient (EacServiceAccessConfig::getInstance()->sharemgntHost, EacServiceAccessConfig::getInstance()->sharemgntPort);
        shareMgntClient->OEM_CheckUninstallPwd(uninstallPwd.getCStr());
    }
    catch (ncTException & e) {
        if (e.errID == ncTShareMgntError::NCT_UNINSTALL_PWD_NOT_ENABLED) {
            THROW_E (EAC_HTTP_SERVER, EACHTTP_UNINSTALL_PWD_NOT_ENABLED, e.expMsg.c_str ());
        }
        else if (e.errID == ncTShareMgntError::NCT_UNINSTALL_PWD_INCORRECT) {
            THROW_E (EAC_HTTP_SERVER, EACHTTP_UNINSTALL_PWD_INCORRECT, e.expMsg.c_str ());
        }
        else {
            THROW_E (EAC_HTTP_SERVER, EACHTTP_SERVER_UNKNOWN_ERR, e.expMsg.c_str ());
        }
    }
    catch (TException & e) {
        THROW_E (EAC_HTTP_SERVER, EACHTTP_SERVER_UNKNOWN_ERR, e.what ());
    }
}

void ncEACHttpServerUtil::CheckExitPwd(const String& exitPwd)
{
    try {
        ncThriftClient<ncTShareMgntClient> shareMgntClient (EacServiceAccessConfig::getInstance()->sharemgntHost, EacServiceAccessConfig::getInstance()->sharemgntPort);
        bool enableExitPwd = shareMgntClient->GetCustomConfigOfBool("enable_exit_pwd");
        if (!enableExitPwd) {
            // 未开启功能
            THROW_E (EAC_HTTP_SERVER, EACHTTP_EXIT_PWD_NOT_ENABLED, LOAD_STRING(_T("IDS_EXIT_PWD_NOT_ENABLED")));
        }

        string pwd;
        shareMgntClient->GetCustomConfigOfString(pwd, "exit_pwd");

        if (0 != exitPwd.compare(toCFLString(pwd))) {
            // 密码错误
            THROW_E (EAC_HTTP_SERVER, EACHTTP_EXIT_PWD_INCORRECT, LOAD_STRING(_T("IDS_EIXT_PWD_INCORRECT")));
        }
    }
    catch (ncTException & e) {
            THROW_E (EAC_HTTP_SERVER, EACHTTP_SERVER_UNKNOWN_ERR, e.expMsg.c_str ());
    }
    catch (TException & e) {
        THROW_E (EAC_HTTP_SERVER, EACHTTP_SERVER_UNKNOWN_ERR, e.what ());
    }
}

void ncEACHttpServerUtil::GetThirdCSFSysConfig(ncTThirdCSFSysConfig & config)
{
    try {
        ncThriftClient<ncTShareMgntClient> shareMgntClient (EacServiceAccessConfig::getInstance()->sharemgntHost, EacServiceAccessConfig::getInstance()->sharemgntPort);
        shareMgntClient->GetThirdCSFSysConfig(config);
    }
    catch (ncTException & e) {
        THROW_E (EAC_HTTP_SERVER, EACHTTP_GET_THIRDCSFSYS_CONFIG_ERROR,
            LOAD_STRING (_T("IDS_EACHTTP_GET_THIRDCSFSYS_CONFIG_ERROR")), e.expMsg.c_str ());
    }
    catch (TException & e) {
        THROW_E (EAC_HTTP_SERVER, EACHTTP_GET_THIRDCSFSYS_CONFIG_ERROR,
            LOAD_STRING (_T("IDS_EACHTTP_GET_THIRDCSFSYS_CONFIG_ERROR")), e.what ());
    }
}

bool ncEACHttpServerUtil::GetShareDocStatus(int docType, int linkType)
{
    try {
        ncThriftClient<ncTShareMgntClient> shareMgntClient (EacServiceAccessConfig::getInstance()->sharemgntHost, EacServiceAccessConfig::getInstance()->sharemgntPort);
        return shareMgntClient->GetShareDocStatus(static_cast<ncTDocType::type>(docType), static_cast<ncTTemplateType::type>(linkType));
    }
    catch (ncTException & e) {
        THROW_E (EAC_HTTP_SERVER, EACHTTP_GET_SHARE_DOC_STATUS_ERROR,
            LOAD_STRING (_T("IDS_EACHTTP_GET_SHARE_DOC_STATUS_ERROR")), e.expMsg.c_str ());
    }
    catch (TException & e) {
        THROW_E (EAC_HTTP_SERVER, EACHTTP_GET_SHARE_DOC_STATUS_ERROR,
            LOAD_STRING (_T("IDS_EACHTTP_GET_SHARE_DOC_STATUS_ERROR")), e.what ());
    }
}

int ncEACHttpServerUtil::ParseLockTime(const char *str)
{
    boost::regex reg("\\d+");
    boost::cmatch what;
    string result;
    if (boost::regex_search(str, what, reg)) {
        result = string(what[0]);
        return atoi(result.c_str());
    } else {
        return 0;
    }
}

bool ncEACHttpServerUtil::HideOum_Check(const string& userId)
{
    try {
        ncThriftClient<ncTShareMgntClient> shareMgntClient (EacServiceAccessConfig::getInstance()->sharemgntHost, EacServiceAccessConfig::getInstance()->sharemgntPort);
        return shareMgntClient->HideOum_Check(userId);
    }
    catch (ncTException & e) {
        THROW_E (EAC_HTTP_SERVER, HIDE_OU_CHECK_ERROR, e.expMsg.c_str ());
    }
    catch (TException & e) {
        THROW_E (EAC_HTTP_SERVER, HIDE_OU_CHECK_ERROR, e.what ());
    }
}

void ncEACHttpServerUtil::Usrm_CreateVcodeInfo(ncTVcodeCreateInfo& vcodeInfo, const string& uuid, const ncTVcodeType::type vcodeType)
{
    try {
        ncThriftClient<ncTShareMgntClient> shareMgntClient (EacServiceAccessConfig::getInstance()->sharemgntHost, EacServiceAccessConfig::getInstance()->sharemgntPort);
        return shareMgntClient->Usrm_CreateVcodeInfo(vcodeInfo, uuid, vcodeType);
    }
    catch (ncTException & e) {
        THROW_E (EAC_HTTP_SERVER, EACHTTP_CREATE_VCODE_ERROR, _T("%s"), e.expMsg.c_str ());
    }
    catch (TException & e) {
        THROW_E (EAC_HTTP_SERVER, EACHTTP_CREATE_VCODE_ERROR, _T("%s"), e.what ());
    }
    catch (Exception& e){
        THROW_E (EAC_HTTP_SERVER, EACHTTP_CREATE_VCODE_ERROR, _T("%s"), e.toString ().getCStr ());
    }
}

void ncEACHttpServerUtil::SetUserInfo(const String& userId, map<String, String> &userinfo)
{
    try {
        ncTEditUserParam param;
        param.id = toSTLString(userId);
        map<String,String>::iterator emailAddress = userinfo.find("emailaddress");
        if (emailAddress != userinfo.end())
        {
            param.__isset.email = true;
            param.email = toSTLString(emailAddress->second);
        }
        map<String,String>::iterator displayName = userinfo.find("displayname");
        if (displayName != userinfo.end())
        {
            param.__isset.displayName = true;
            param.displayName = toSTLString(displayName->second);
        }
        map<String,String>::iterator telNumber = userinfo.find("telnumber");
        if (telNumber != userinfo.end())
        {
            param.__isset.telNumber = true;
            param.telNumber = toSTLString(telNumber->second);
        }
        ncThriftClient<ncTShareMgntClient> shareMgntClient (EacServiceAccessConfig::getInstance()->sharemgntHost, EacServiceAccessConfig::getInstance()->sharemgntPort);
        shareMgntClient->Usrm_EditUser(param, "");
    }
    catch (ncTException & e) {
        if (e.errID == ncTShareMgntError::NCT_USER_NOT_EXIST) {
            THROW_E (EAC_HTTP_SERVER, EACHTTP_USER_NOT_EXIST, _T("%s"), e.expMsg.c_str ());
        }
        else if(e.errID == ncTShareMgntError::NCT_INVALID_EMAIL) {
            THROW_E (EAC_HTTP_SERVER, EACHTTP_INVALID_EMAIL_ADDRESS, _T("%s"), e.expMsg.c_str ());
        }
        else if (e.errID == ncTShareMgntError::NCT_DUPLICATED_EMALI) {
            THROW_E (EAC_HTTP_SERVER, EACHTTP_DUPLICATED_EMALI_ADDRESS, _T("%s"), e.expMsg.c_str ());
        }
        else if(e.errID == ncTShareMgntError::NCT_INVALID_DISPLAY_NAME){
            THROW_E (EAC_HTTP_SERVER, EACHTTP_INVALID_DISPLAY_NAME, _T("%s"), e.expMsg.c_str ());
        }
        else if(e.errID == ncTShareMgntError::NCT_DUPLICATED_DISPLAY_NAME){
            THROW_E (EAC_HTTP_SERVER, EACHTTP_DUPLICATED_DISPLAY_NAME, _T("%s"), e.expMsg.c_str ());
        }
        else if(e.errID == ncTShareMgntError::NCT_INVALID_TEL_NUMBER){
            THROW_E (EAC_HTTP_SERVER, EACHTTP_INVALID_TEL_NUMBER, _T("%s"), e.expMsg.c_str ());
        }
        else if(e.errID == ncTShareMgntError::NCT_TEL_NUMBER_EXISTS){
            THROW_E (EAC_HTTP_SERVER, EACHTTP_TEL_NUMBER_EXISTS, _T("%s"), e.expMsg.c_str ());
        }
        else {
            THROW_E (EAC_HTTP_SERVER, EACHTTP_EDIT_USER_INFO_ERROR, _T("%s"), e.expMsg.c_str ());
        }
    }
    catch (TException & e) {
        THROW_E (EAC_HTTP_SERVER, EACHTTP_EDIT_USER_INFO_ERROR, _T("%s"), e.what ());
    }
}

void ncEACHttpServerUtil::SMSSendVcode(const string& account, const string& password, const string& telNumber)
{
    try {
        ncThriftClient<ncTShareMgntClient> shareMgntClient (EacServiceAccessConfig::getInstance()->sharemgntHost, EacServiceAccessConfig::getInstance()->sharemgntPort);
        return shareMgntClient->SMS_SendVcode(account, password, telNumber);
    }
    catch (ncTException & e) {
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
        if (e.errID == ncTShareMgntError::NCT_USER_NOT_EXIST ) {
            THROW_E (EAC_HTTP_SERVER, EACHTTP_USER_NOT_EXIST, e.expMsg.c_str ());
        }
        else if (e.errID == ncTShareMgntError::NCT_INVALID_ACCOUNT_OR_PASSWORD) {
            THROW_HTTP_DETAIL_E (EAC_HTTP_SERVER, errDetail, EACHTTP_INVALID_USER_OR_PASSWORD, e.expMsg.c_str ());
        }
        else if (e.errID == ncTShareMgntError::NCT_INVALID_TEL_NUMBER) {
            THROW_E (EAC_HTTP_SERVER, EACHTTP_INVALID_TEL_NUMBER, e.expMsg.c_str ());
        }
        else if (e.errID == ncTShareMgntError::NCT_TEL_NUMBER_EXISTS) {
            THROW_E (EAC_HTTP_SERVER, EACHTTP_TEL_NUMBER_EXISTS, e.expMsg.c_str ());
        }
        else if (e.errID == ncTShareMgntError::NCT_USER_IS_ACTIVATE) {
            THROW_E (EAC_HTTP_SERVER, EACHTTP_USER_IS_ACTIVATE, e.expMsg.c_str ());
        }
        else if (e.errID == ncTShareMgntError::NCT_SMS_ACTIVATE_DISABLED) {
            THROW_E (EAC_HTTP_SERVER, EACHTTP_SMS_ACTIVATE_DISABLED, e.expMsg.c_str ());
        }
        else {
            THROW_E (EAC_HTTP_SERVER, EACHTTP_SEND_VERIFY_CODE_FAIL, e.expMsg.c_str ());
        }
    }
    catch (TException & e) {
        THROW_E (EAC_HTTP_SERVER, EACHTTP_SEND_VERIFY_CODE_FAIL, _T("%s"), e.what ());
    }
    catch (Exception& e){
        THROW_E (EAC_HTTP_SERVER, EACHTTP_SEND_VERIFY_CODE_FAIL, _T("%s"), e.toString ().getCStr ());
    }
}

void ncEACHttpServerUtil::SMSActivate(string& retUserId, const string& account, const string& password, const string& telNumber, const string& mailAddress, const string& verifyCode)
{
    try {
        ncThriftClient<ncTShareMgntClient> shareMgntClient (EacServiceAccessConfig::getInstance()->sharemgntHost, EacServiceAccessConfig::getInstance()->sharemgntPort);
        return shareMgntClient->SMS_Activate(retUserId, account, password, telNumber, mailAddress, verifyCode);
    }
    catch (ncTException & e) {
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
        if (e.errID == ncTShareMgntError::NCT_USER_NOT_EXIST ) {
            THROW_E (EAC_HTTP_SERVER, EACHTTP_USER_NOT_EXIST, e.expMsg.c_str ());
        }
        else if (e.errID == ncTShareMgntError::NCT_INVALID_ACCOUNT_OR_PASSWORD) {
            THROW_HTTP_DETAIL_E (EAC_HTTP_SERVER, errDetail, EACHTTP_INVALID_USER_OR_PASSWORD, e.expMsg.c_str ());
        }
        else if (e.errID == ncTShareMgntError::NCT_INVALID_TEL_NUMBER) {
            THROW_E (EAC_HTTP_SERVER, EACHTTP_INVALID_TEL_NUMBER, e.expMsg.c_str ());
        }
        else if (e.errID == ncTShareMgntError::NCT_TEL_NUMBER_EXISTS) {
            THROW_E (EAC_HTTP_SERVER, EACHTTP_TEL_NUMBER_EXISTS, e.expMsg.c_str ());
        }
        else if (e.errID == ncTShareMgntError::NCT_USER_IS_ACTIVATE) {
            THROW_E (EAC_HTTP_SERVER, EACHTTP_USER_IS_ACTIVATE, e.expMsg.c_str ());
        }
        else if (e.errID == ncTShareMgntError::NCT_INVALID_EMAIL) {
            THROW_E (EAC_HTTP_SERVER, EACHTTP_SMTP_RECIPIENT_MAIL_ILLEGAL,
            LOAD_STRING (_T("IDS_EACHTTP_SMTP_RECIPIENT_MAIL_ILLEGAL")), e.expMsg.c_str ());
        }
        else if (e.errID == ncTShareMgntError::NCT_DUPLICATED_EMALI) {
            THROW_E (EAC_HTTP_SERVER, EACHTTP_EMAIL_EXISTS, e.expMsg.c_str ());
        }
        else if (e.errID == ncTShareMgntError::NCT_SMS_VERIFY_CODE_ERROR) {
            THROW_E (EAC_HTTP_SERVER, EACHTTP_SMS_VERIFY_CODE_ERROR, e.expMsg.c_str ());
        }
        else if (e.errID == ncTShareMgntError::NCT_SMS_VERIFY_CODE_TIMEOUT) {
            THROW_E (EAC_HTTP_SERVER, EACHTTP_SMS_VERIFY_CODE_TIMEOUT, e.expMsg.c_str ());
        }
        else if (e.errID == ncTShareMgntError::NCT_SMS_ACTIVATE_DISABLED) {
            THROW_E (EAC_HTTP_SERVER, EACHTTP_SMS_ACTIVATE_DISABLED, e.expMsg.c_str ());
        }
        else {
            THROW_E (EAC_HTTP_SERVER, EACHTTP_USER_ACTIVATE_ERROR, _T("%s"), e.expMsg.c_str ());
        }
    }
    catch (TException & e) {
        THROW_E (EAC_HTTP_SERVER, EACHTTP_USER_ACTIVATE_ERROR, _T("%s"), e.what ());
    }
    catch (Exception& e){
        THROW_E (EAC_HTTP_SERVER, EACHTTP_USER_ACTIVATE_ERROR, _T("%s"), e.toString ().getCStr ());
    }
}

inline bool CheckIsUTF8Tail (unsigned char b)
{
    return (b & 0xC0) == 0x80;
}

size_t ncEACHttpServerUtil::CheckIsTextUTF8 (const String& str)
{
    unsigned char* text = (unsigned char*)str.getCStr ();
    size_t size = str.getLength ();
    size_t i = 0;
    size_t ok_count = 0;
    while (i < size) {
        // 1 byte
        if (text[i] < 0x80) {
            if(text[i] <= 0x10 && (text[i] <= 7 || text[i] >= 0x0E))
                return 0;
            ++i;
            ++ok_count;
        }
        else if (text[i] <= 0xDF) {
            // 2 bytes
            if(++i < size && CheckIsUTF8Tail (text[i])) {
                ++i;
                ++ok_count;
            }
            else if (size != i) {
                return 0;
            }
        }
        else if (text[i] <= 0xEF) {
            // 3 bytes
            if ((++i < size && CheckIsUTF8Tail (text[i])) && (++i < size && CheckIsUTF8Tail (text[i]))) {
                ++i;
                ++ok_count;
            }
            else if (size != i) {
                return 0;
            }
        }
        else {
            // 4 bytes is checked as invalid
            return 0;
        }
    }

    return ok_count;
}

void ncEACHttpServerUtil::CheckOsType(ACSClientType clientType)
{
    // 检查设备类型
    if(clientType < ACSClientType::UNKNOWN || clientType > ACSClientType::LINUX) {
        THROW_E (EAC_HTTP_SERVER, EACHTTP_DEVICE_INFO_INVALID,
            LOAD_STRING (_T("IDS_INVALID_OS_TYPE")));
    }
}

void ncEACHttpServerUtil::SendAuthVcode (ncTReturnInfo &retInfo, string& userId, const ncTVcodeType::type vcodeType, const string oldTelnum)
{
    // 生成并发送双因子短信验证码
    try {

        ncThriftClient<ncTShareMgntClient> shareMgntClient (EacServiceAccessConfig::getInstance()->sharemgntHost, EacServiceAccessConfig::getInstance()->sharemgntPort);
        shareMgntClient->Usrm_SendAuthVcode(retInfo, userId, vcodeType, oldTelnum);
    }
    catch (ncTException & e) {
        if (e.errID == ncTShareMgntError::NCT_MFA_SMS_SERVER_ERROR) {
            THROW_E(EAC_HTTP_SERVER, EACHTTP_MFA_SMS_SERVER_ERROR, e.expMsg.c_str ());
        }
        else if (e.errID == ncTShareMgntError::NCT_PHONE_NUMBER_HAS_BEEN_CHANGED) {
            THROW_E(EAC_HTTP_SERVER, EACHTTP_PHONE_NUMBER_HAS_BEEN_CHANGED, e.expMsg.c_str ());
        }
        else if (e.errID == ncTShareMgntError::NCT_USER_HAS_NOT_BOUND_PHONE) {
            THROW_E(EAC_HTTP_SERVER, EACHTTP_USER_HAS_NOT_BOUND_PHONE, e.expMsg.c_str ());
        }
        else if(e.errID == ncTShareMgntError::NCT_THIRD_PLUGIN_INTER_ERROR){
            THROW_E(EAC_HTTP_SERVER, EACHTTP_THIRD_PLUGIN_INTER_ERROR, e.expMsg.c_str ());
        }
        else if(e.errID == ncTShareMgntError::NCT_MFA_CONFIG_ERROR){
            THROW_E(EAC_HTTP_SERVER, EACHTTP_MFA_CONFIG_ERROR, e.expMsg.c_str ());
        }
        else if(e.errID == ncTShareMgntError::NCT_FAILED_THIRD_CONFIG){
            THROW_E(EAC_HTTP_SERVER, EACHTTP_FAILED_THIRD_CONFIG, e.expMsg.c_str ());
        }
        else {
            THROW_E (EAC_HTTP_SERVER, EACHTTP_SEND_VERIFY_CODE_FAIL, e.expMsg.c_str ());
        }
    }
    catch (TException & e) {
        THROW_E (EAC_HTTP_SERVER, EACHTTP_CREATE_VCODE_ERROR, _T("%s"), e.what ());
    }
    catch (Exception& e){
        THROW_E (EAC_HTTP_SERVER, EACHTTP_CREATE_VCODE_ERROR, _T("%s"), e.toString ().getCStr ());
    }
}

bool ncEACHttpServerUtil::GetThirdAuthTypeStatus (const ncTMFAType::type authType)
{
    // 判断双因子验证插件类型
    try {
        ncThriftClient<ncTShareMgntClient> shareMgntClient (EacServiceAccessConfig::getInstance()->sharemgntHost, EacServiceAccessConfig::getInstance()->sharemgntPort);
        return shareMgntClient->GetThirdAuthTypeStatus(authType);
    }
    catch (ncTException & e) {
        THROW_E (EAC_HTTP_SERVER, EACHTTP_THIRD_AUTH_NOT_OPEN, _T("%s"), e.expMsg.c_str ());
    }
    catch (TException & e) {
        THROW_E (EAC_HTTP_SERVER, EACHTTP_THIRD_AUTH_NOT_OPEN, _T("%s"), e.what ());
    }
    catch (Exception& e){
        THROW_E (EAC_HTTP_SERVER, EACHTTP_THIRD_AUTH_NOT_OPEN, _T("%s"), e.toString ().getCStr ());
    }
}

void ncEACHttpServerUtil::GetSmtpSrvConfig(ncTSmtpSrvConf& config)
{
    try {
        ncThriftClient<ncTShareMgntClient> shareMgntClient (EacServiceAccessConfig::getInstance()->sharemgntHost, EacServiceAccessConfig::getInstance()->sharemgntPort);
        shareMgntClient->SMTP_GetConfig(config);
    }
    catch (ncTException & e) {
        THROW_E (EAC_HTTP_SERVER, EACHTTP_GET_SMTPSRV_CONFIG_ERROR,
            LOAD_STRING (_T("IDS_EACHTTP_GET_SMTPSRV_CONFIG_ERROR")), e.expMsg.c_str ());
    }
    catch (TException & e) {
        THROW_E (EAC_HTTP_SERVER, EACHTTP_GET_SMTPSRV_CONFIG_ERROR,
            LOAD_STRING (_T("IDS_EACHTTP_GET_SMTPSRV_CONFIG_ERROR")), e.what ());
    }
}

void ncEACHttpServerUtil::GetAllConfig(ncTAllConfig& config)
{
    const int64_t timeOut = 300;

    static int64_t sUpdateConfigTime = 0;
    static ncTAllConfig allConfig;
    static boost::shared_mutex _configUpdateLock;

    int64 localTime = BusinessDate::getCurrentTime () / 1000000;
    int64_t nEscapeTime = 0;
    bool isEscape = false;

    {
        boost::shared_lock<boost::shared_mutex> lock(_configUpdateLock);
        nEscapeTime = localTime - sUpdateConfigTime;
        if (timeOut > nEscapeTime && 0 < nEscapeTime) {
            isEscape = false;
            config = allConfig;
        } else {
            isEscape = true;
        }
    }

    if (isEscape) {
        boost::lock_guard<boost::shared_mutex> lock(_configUpdateLock);
        localTime = BusinessDate::getCurrentTime () / 1000000;
        nEscapeTime = localTime - sUpdateConfigTime;
        if (timeOut > nEscapeTime && 0 < nEscapeTime) {
            config = allConfig;
        } else {
            try {
                // 使用临时变量获取配置后赋值给静态，避免由于返回值个别配置未设置导致静态变量内容未更新
                ncTAllConfig tmpConfig;
                ncThriftClient<ncTShareMgntClient> shareMgntClient (EacServiceAccessConfig::getInstance()->sharemgntHost, EacServiceAccessConfig::getInstance()->sharemgntPort);
                shareMgntClient->GetAllConfig(tmpConfig);
                allConfig = tmpConfig;
            }
            catch (ncTException & e) {
                THROW_E (EAC_HTTP_SERVER, EACHTTP_SERVER_UNKNOWN_ERR,
                    "Get all config error.", e.expMsg.c_str ());
            }
            catch (TException & e) {
                THROW_E (EAC_HTTP_SERVER, EACHTTP_SERVER_UNKNOWN_ERR,
                    "Get all config error.", e.what ());
            }
            config = allConfig;
            sUpdateConfigTime = localTime;
        }
    }
}

bool ncEACHttpServerUtil::CheckIsRFC3339 (const String& rfcTimeStr)
{
    bool ret = false;
    if (!rfcTimeStr.isEmpty ()){
        Regex __reg__1 (_T("^\\d{4}-\\d{1,2}-\\d{1,2}$"));
        Regex __reg__2 (_T("^\\d{4}-\\d{1,2}-\\d{1,2}T(0\\d{1}|1\\d{1}|2[0-3]):[0-5]\\d{1}:([0-5]\\d{1}).?(\\d{1,4})?Z$"));
        Regex __reg__3 (_T("^\\d{4}-\\d{1,2}-\\d{1,2}T(0\\d{1}|1\\d{1}|2[0-3]):[0-5]\\d{1}:([0-5]\\d{1})[-|+]{1}(0\\d{1}|1\\d{1}|2[0-3]):[0-5]\\d{1}$"));
        if (__reg__1.match (rfcTimeStr) ||
            __reg__2.match (rfcTimeStr) ||
            __reg__3.match (rfcTimeStr)) {
            ret = true;
        }
    }
    return ret;
}

int64 ncEACHttpServerUtil::RFC3339ToTimeStamp (String& rfcTimeStr)
{
    if (!CheckIsRFC3339(rfcTimeStr)) {
        THROW_E (EAC_HTTP_SERVER, EACHTTP_INVALID_TIME_FORMAT, "invalid time format");
    }

    size_t idxAddTime, idxLessTime, idxT;
    idxT = rfcTimeStr.find ("T");

    // 转时间格式：
    // 2019-10-31T10:00:00+08:00 >> 2019-10-31 10:00:00 带时区的以客户端时区为准
    // 2019-10-31T02:00:00Z >> 2019-10-31 02:00:00
    // 2019-10-31 >> 2019-10-31 00:00:00 不带时区信息的以服务器时区为准
    String rfcTimeJetLag;
    string rfcTime = toSTLString (rfcTimeStr);
    if (idxT != String::NO_POSITION) {
        String tmpRfcTimeStr;
        rfcTimeStr.replace (rfcTimeStr.find ("T"), 1, " ");
        tmpRfcTimeStr = rfcTimeStr.subString (0, 19);
        rfcTimeJetLag = rfcTimeStr.subString (19);
        rfcTime = toSTLString (tmpRfcTimeStr);
    }
    struct tm tm;
    memset (&tm, 0, sizeof (tm));
    sscanf (rfcTime.c_str (), "%d-%d-%d %d:%d:%d",
           &tm.tm_year, &tm.tm_mon, &tm.tm_mday,
           &tm.tm_hour, &tm.tm_min, &tm.tm_sec);

    tm.tm_year -= 1900;
    tm.tm_mon--;

    // 当时间格式为带时区的格式需转换成客户端标准时间格式
    // 2019-10-31T10:00:00+08:00 >> 2019-10-31T02:00:00Z
    // rfcTimeJetLag -> 08:00
    idxAddTime = rfcTimeJetLag.find ("+");
    idxLessTime = rfcTimeJetLag.find ("-");
    vector<String> subtime;
    int tmphour;
    int tmpmin;
    if (idxAddTime != String::NO_POSITION) {
        size_t pos = rfcTimeJetLag.findFirstOf("+");
        rfcTimeJetLag = rfcTimeJetLag.subString((pos == String::NO_POSITION) ? 0 : (pos + 1));
        rfcTimeJetLag.split (":", subtime);
        tmphour = atoi (subtime[0].getCStr ());
        tmpmin = atoi (subtime[1].getCStr ());
        if (tm.tm_min - tmpmin < 0) {
            tm.tm_hour--;
        }
        tm.tm_min = (tm.tm_min + 60 - tmpmin) % 60;
        if (tm.tm_hour - tmphour < 0) {
            tm.tm_mday--;
        }
        tm.tm_hour = (tm.tm_hour + 24 - tmphour) % 24;
    }
    else if (idxLessTime != String::NO_POSITION) {
        size_t pos = rfcTimeJetLag.findFirstOf("-");
        rfcTimeJetLag = rfcTimeJetLag.subString((pos == String::NO_POSITION) ? 0 : (pos + 1));
        rfcTimeJetLag.split (":", subtime);
        tmphour = atoi (subtime[0].getCStr ());
        tmpmin = atoi (subtime[1].getCStr ());
        if (tm.tm_min + tmpmin >= 60) {
            tm.tm_hour++;
        }
        tm.tm_min = (tm.tm_min + tmpmin) % 60;
        if (tm.tm_hour + tmphour >= 24) {
            tm.tm_mday++;
        }
        tm.tm_hour = (tm.tm_hour + tmphour) % 24;
    }
    if (idxT != String::NO_POSITION) {
        /*主要对带时区的时间格式
        * 获取服务器时区，mongodb存储时会对当前服务器时区进行时差增减，
        * 时间格式为标准时间时需要做对应的加减
        */
        time_t ts = 0;
        struct tm t1;
        char buf[16];
        ::localtime_r (&ts, &t1);
        ::strftime (buf, sizeof (buf), "%z", &t1);
        string timezone = buf;
        String timezoneStr = toCFLString (timezone);
        idxAddTime = timezoneStr.find ("+");
        idxLessTime = timezoneStr.find ("-");
        // 本地时区差示例（伊朗）：+0330
        int timezoneHour = atoi (timezoneStr.subString (1, 2).getCStr ());
        int timezoneMin = atoi (timezoneStr.subString (3).getCStr ());
        if (idxAddTime != String::NO_POSITION) {
            if (tm.tm_min + timezoneMin >= 60) {
                tm.tm_hour++;
            }
            tm.tm_min = (tm.tm_min + timezoneMin) % 60;
            if (tm.tm_hour + timezoneHour >= 24) {
                tm.tm_mday++;
            }
            tm.tm_hour = (tm.tm_hour + timezoneHour) % 24;
        }
        else if (idxLessTime != String::NO_POSITION) {
            if (tm.tm_min - timezoneMin < 0) {
                tm.tm_hour--;
            }
            tm.tm_min = (tm.tm_min + 60 - timezoneMin) % 60;
            if (tm.tm_hour - timezoneHour < 0) {
                tm.tm_mday--;
            }
            tm.tm_hour = (tm.tm_hour + 24 - timezoneHour) % 24;
        }
    }
    //转换时间戳
    time_t t = mktime (&tm);
    int64 timeStamp = 0 == t ? -1 : t * 1000000;
    return timeStamp;
}

void ncEACHttpServerUtil::GetVisitor (brpc::Controller* cntl, ncIntrospectInfo& introspectInfo)
{
    String tokenId;
    ncHttpGetToken (cntl, tokenId);

    // token验证
    ncCheckTokenInfo checkTokenInfo;
    checkTokenInfo.tokenId = tokenId;
    checkTokenInfo.ip = GetForwardedIp (cntl);
    if (!CheckToken (checkTokenInfo, introspectInfo)) {
        THROW_E (EAC_HTTP_SERVER, EACHTTP_SERVER_ACCESS_TOKEN_ERR,
            LOAD_STRING (_T("IDS_EACHTTP_ERR_MSG_ACCESS_TOKEN_ERR")));
    }
}

int ncEACHttpServerUtil::PermArrayToInt (const JSON::Array& permArray, const String& paramName)
{
    static map<String, ncAtomPermValue> permStrIntMap;
    static ThreadMutexLock valueLock;

    if (permStrIntMap.empty ()) {
        AutoLock<ThreadMutexLock> lock (&valueLock);
        if (permStrIntMap.empty ()) {
            permStrIntMap.insert (pair<String, ncAtomPermValue>(_T("display"), ncAtomPermValue::ACS_AP_DISPLAY));
            permStrIntMap.insert (pair<String, ncAtomPermValue>(_T("preview"), ncAtomPermValue::ACS_AP_PREVIEW));
            permStrIntMap.insert (pair<String, ncAtomPermValue>(_T("download"), ncAtomPermValue::ACS_AP_READ));
            permStrIntMap.insert (pair<String, ncAtomPermValue>(_T("create"), ncAtomPermValue::ACS_AP_CREATE));
            permStrIntMap.insert (pair<String, ncAtomPermValue>(_T("modify"), ncAtomPermValue::ACS_AP_EDIT));
            permStrIntMap.insert (pair<String, ncAtomPermValue>(_T("delete"), ncAtomPermValue::ACS_AP_DELETE));
        }
    }

    int perm = 0;
    for (size_t i = 0; i < permArray.size (); ++i){
        auto iter = permStrIntMap.find (permArray[i].s ().c_str ());
        if (iter == permStrIntMap.end ()) {
            THROW_E (EAC_HTTP_SERVER, EACHTTP_ARGUMENT_INVALID,
                     LOAD_STRING (_T("IDS_EACHTTP_INVALID_ARG_VALUE")), paramName.getCStr ());
        }

        perm |= iter->second;
    }

    return perm;
}

void ncEACHttpServerUtil::IntToPermArrayV1 (int perm, JSON::Array& permArray)
{
    if (perm & ncAtomPermValue::ACS_AP_DISPLAY) {
        permArray.push_back (_T("display"));
    }
    if (perm & ncAtomPermValue::ACS_AP_READ) {
        permArray.push_back (_T("read"));
    }
    if (perm & ncAtomPermValue::ACS_AP_CREATE) {
        permArray.push_back (_T("create"));
    }
    if (perm & ncAtomPermValue::ACS_AP_EDIT) {
        permArray.push_back (_T("modify"));
    }
    if (perm & ncAtomPermValue::ACS_AP_DELETE) {
        permArray.push_back (_T("delete"));
    }
}

void ncEACHttpServerUtil::IntToPermArray (int perm, JSON::Array& permArray)
{
    if (perm & ncAtomPermValue::ACS_AP_DISPLAY) {
        permArray.push_back (_T("display"));
    }
    if (perm & ncAtomPermValue::ACS_AP_PREVIEW) {
        permArray.push_back (_T("preview"));
    }
    if (perm & ncAtomPermValue::ACS_AP_READ) {
        permArray.push_back (_T("download"));
    }
    if (perm & ncAtomPermValue::ACS_AP_CREATE) {
        permArray.push_back (_T("create"));
    }
    if (perm & ncAtomPermValue::ACS_AP_EDIT) {
        permArray.push_back (_T("modify"));
    }
    if (perm & ncAtomPermValue::ACS_AP_DELETE) {
        permArray.push_back (_T("delete"));
    }
}

string ncEACHttpServerUtil::RSADecrypt2048(const string &cipherText)
{
    static ThreadMutexLock sLock;
    static RSA *p_rsa = NULL;

    static string prikey = "-----BEGIN RSA PRIVATE KEY-----\n\
MIIEpQIBAAKCAQEAsyOstgbYuubBi2PUqeVjGKlkwVUY6w1Y8d4k116dI2SkZI8f\n\
xcjHALv77kItO4jYLVplk9gO4HAtsisnNE2owlYIqdmyEPMwupaeFFFcg751oiTX\n\
JiYbtX7ABzU5KQYPjRSEjMq6i5qu/mL67XTkhvKwrC83zme66qaKApmKupDODPb0\n\
RRkutK/zHfd1zL7sciBQ6psnNadh8pE24w8O2XVy1v2bgSNkGHABgncR7seyIg81\n\
JQ3c/Axxd6GsTztjLnlvGAlmT1TphE84mi99fUaGD2A1u1qdIuNc+XuisFeNcUW6\n\
fct0+x97eS2eEGRr/7qxWmO/P20sFVzXc2bF1QIDAQABAoIBAACDungGYoJ87bLl\n\
DUQUqtl0CRxODoWEUwxUz0XIGYrzu84nJBf5GOs9Xv6i9YbNgJN2xkJrtTU7VUJF\n\
AfaSP4kZXqqAO9T1Id9zVc5oomuldSiLUwviwaMek1Yh9sFRqWNGGxBdd7Y1ckm8\n\
Roy+kHZ7xXqlIxOmdCC+7DgQMVgSV64wzQY8p7L9kTLIkeDodEolkUkGsreF9I9S\n\
kzlLjGU9flPt13319G0KSaQUWEpxF/UBr2gKJvQPQHSRzzl5HlRwznZkU4Hs6RID\n\
ue6E68ZJNMRn3FUAvLMCRw9C4PQQR/x/50WH4BXJ9veVIOIpTVCJedI0QZjbVuBk\n\
RPKHTMkCgYEA2XjGIw9Vp0qu/nCeo5Nk15xt/SJCn0jIhyRpckHtCidotkiZmFdU\n\
vUK7IwbAUPqEJcgmS/zwREV8Gff8S324C2RoDN4FxFtBMZgQjqV1zYqGLQSbTJUh\n\
GlpTe7jKVskuSPSf00OqqAIlYNtzZK3mWj8MadFD99Wo9gktXRAFdf0CgYEA0uBe\n\
wuE007XLqb8ANS+4U0CkexeVDkDzI9yXN2CB+L5wmJ/WsNF8iD53xHxpwZWRiizX\n\
ArBdhWL9yv4YkbryyD15NRSQhLanRcs0MqGh1GJJ9vpGzBjfJJ3Bw0hBfkwnf/C6\n\
nNzGjNWNTeNKwlcFaVhBADyGYZt9Len9YYFNKrkCgYEAmsn7BYNprOxciCAy2i0U\n\
Lt9Z7j3Pe757dK13HGtOQ9bvEie0o5ktaJSxzGmGw1y8aIQAtj9v6Lgob/dxrW3r\n\
bLhn0xjItA1b5ufciRu+MLFzdWF9BFJ1QGOgXkSWSJVji2wKwn28X18/qaQpizS3\n\
6+5KcJsRrLp4S78WedHogSUCgYEAomb5k8wtCv7vIoNefZeKtVMLWWEIAjozBmNU\n\
cel5L0A7Js+yX+p1pde2FTRbniK6O1fdHs0EuT1Lh5G5CkKXx27QcfisdAjXOgEM\n\
6hFguFgZ7oNBEt30vBZiqypyhfnQUc/rZ/L/VmcAtANgB9tM55x4Mt5p/7Hn7fxO\n\
j1EtRMECgYEAp2sI035BcCR2kFW1vC9eXLAPZ0anyy1/T1dEgFJ/ELqmGEMEWZKA\n\
9H1KH6YIkDdXabwfaSTRebaEescCxRtgmo5WEdZxw4Nz66SSomc24aD0iem7+VSl\n\
x2qRWdif0jHG8fOdMey3NrY7NF4xQTzuO9jDnLpBTwFg3o7QlywIBlM=\n\
-----END RSA PRIVATE KEY-----\n";

    AutoLock<ThreadMutexLock> lock(&sLock);
    if (p_rsa == NULL)
    {
        BIO *in = BIO_new_mem_buf((void *)prikey.c_str(), -1);
        if (in == NULL)
        {
            throw Exception(_T("BIO_new_mem_buf error"));
        }

        p_rsa = PEM_read_bio_RSAPrivateKey(in, NULL, NULL, NULL);
        BIO_free(in);
        if (p_rsa == NULL)
        {
            throw Exception(_T("PEM_read_bio_RSAPrivateKey error"));
            ;
        }
    }

    string plainText;
    int rsa_len = RSA_size(p_rsa);
    plainText.resize(rsa_len, 0);

    if (RSA_private_decrypt(cipherText.length(), (unsigned char *)cipherText.c_str(), (unsigned char *)plainText.c_str(), p_rsa, RSA_PKCS1_PADDING) < 0)
    {
        THROW_E(EAC_HTTP_SERVER, EACHTTP_INVALID_PARAMETER, "RSA_private_decrypt error");
    }

    // 去除末尾的\0
    string ret(plainText.c_str());
    return ret;
}

string ncEACHttpServerUtil::RSADecrypt(const string &cipherText)
{
    static ThreadMutexLock sLock;
    static RSA *p_rsa = NULL;

    static string prikey = "-----BEGIN RSA PRIVATE KEY-----\n\
MIICXgIBAAKBgQDB2fhLla9rMx+6LWTXajnK11Kdp520s1Q+TfPfIXI/7G9+L2YC\n\
4RA3M5rgRi32s5+UFQ/CVqUFqMqVuzaZ4lw/uEdk1qHcP0g6LB3E9wkl2FclFR0M\n\
+/HrWmxPoON+0y/tFQxxfNgsUodFzbdh0XY1rIVUIbPLvufUBbLKXHDPpwIDAQAB\n\
AoGBALCM/H6ajXFs1nCR903aCVicUzoS9qckzI0SIhIOPCfMBp8+PAJTSJl9/ohU\n\
YnhVj/kmVXwBvboxyJAmOcxdRPWL7iTk5nA1oiVXMer3Wby+tRg/ls91xQbJLVv3\n\
oGSt7q0CXxJpRH2oYkVVlMMlZUwKz3ovHiLKAnhw+jEsdL2BAkEA9hA97yyeA2eq\n\
f9dMu/ici99R3WJRRtk4NEI4WShtWPyziDg48d3SOzYmhEJjPuOo3g1ze01os70P\n\
ApE7d0qcyQJBAMmt+FR8h5MwxPQPAzjh/fTuTttvUfBeMiUDrIycK1I/L96lH+fU\n\
i4Nu+7TPOzExnPeGO5UJbZxrpIEUB7Zs8O8CQQCLzTCTGiNwxc5eMgH77kVrRudp\n\
Q7nv6ex/7Hu9VDXEUFbkdyULbj9KuvppPJrMmWZROw04qgNp02mayM8jeLXZAkEA\n\
o+PM/pMn9TPXiWE9xBbaMhUKXgXLd2KEq1GeAbHS/oY8l1hmYhV1vjwNLbSNrH9d\n\
yEP73TQJL+jFiONHFTbYXwJAU03Xgum5mLIkX/02LpOrz2QCdfX1IMJk2iKi9osV\n\
KqfbvHsF0+GvFGg18/FXStG9Kr4TjqLsygQJT76/MnMluw==\n\
-----END RSA PRIVATE KEY-----\n";

    AutoLock<ThreadMutexLock> lock(&sLock);
    if (p_rsa == NULL)
    {
        BIO *in = BIO_new_mem_buf((void *)prikey.c_str(), -1);
        if (in == NULL)
        {
            throw Exception(_T("BIO_new_mem_buf error"));
        }

        p_rsa = PEM_read_bio_RSAPrivateKey(in, NULL, NULL, NULL);
        BIO_free(in);
        if (p_rsa == NULL)
        {
            throw Exception(_T("PEM_read_bio_RSAPrivateKey error"));
            ;
        }
    }

    string plainText;
    int rsa_len = RSA_size(p_rsa);
    plainText.resize(rsa_len, 0);

    if (RSA_private_decrypt(cipherText.length(), (unsigned char *)cipherText.c_str(), (unsigned char *)plainText.c_str(), p_rsa, RSA_PKCS1_PADDING) < 0)
    {
        THROW_E(EAC_HTTP_SERVER, EACHTTP_INVALID_PARAMETER, "RSA_private_decrypt error");
    }

    // 去除末尾的\0
    string ret(plainText.c_str());
    return ret;
}

size_t ncEACHttpServerUtil::calcDecodeLength(const string &tmp)
{
    // Calculates the length of a decoded string
    size_t len = tmp.length();
    size_t padding = 0;

    if (tmp[len - 1] == '=' && tmp[len - 2] == '=') // last two chars are =
        padding = 2;
    else if (tmp[len - 1] == '=') // last char is =
        padding = 1;

    return (len * 3) / 4 - padding;
}

string ncEACHttpServerUtil::Base64Decode(const string &input)
{
    // printf("Before Base64Decode: %s\n", input.c_str());
    // printBytes(input);

    // 去除所有的\r和\n
    string tmp = input;
    boost::replace_all(tmp, "\n", "");
    boost::replace_all(tmp, "\r", "");

    string buffer;
    buffer.resize(calcDecodeLength(tmp));

    BIO *b64 = NULL;
    BIO *bmem = NULL;

    // 使用没有\n的密文进行解码
    b64 = BIO_new(BIO_f_base64());
    BIO_set_flags(b64, BIO_FLAGS_BASE64_NO_NL);
    bmem = BIO_new_mem_buf((char *)tmp.c_str(), tmp.length());
    bmem = BIO_push(b64, bmem);
    BIO_read(bmem, (char *)buffer.c_str(), tmp.length());
    BIO_free_all(bmem);

    // printf("After Base64Decode: \n");
    // printBytes(buffer);

    return buffer;
}
