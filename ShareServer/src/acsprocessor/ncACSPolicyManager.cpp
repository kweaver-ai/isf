/***************************************************************************************************
ncACSPolicyManager.cpp:
    Copyright (c) Eisoo Software Inc. (2009 - 2013), All rights reserved.

Purpose:
    acs policy manager 接口

Author:
    xu.zhi@aishu.cn

Creating Time:
    2020-7-30
***************************************************************************************************/
#include <abprec.h>
#include <arpa/inet.h>
#include <ncutil/ncBusinessDate.h>
#include <biginteger/InfInt.h>
#include <boost/date_time/posix_time/posix_time.hpp>
#include <ehttpclient/public/ncIEHTTPClient.h>

#include "acsprocessor.h"
#include "ncACSPolicyManager.h"
#include "ncACSProcessorUtil.h"
#include "ncACSDeviceManager.h"
#include "acsServiceAccessConfig.h"

using namespace boost::posix_time;
using namespace boost::gregorian;

/* Implementation file */
NS_IMPL_QUERY_INTERFACE1 (ncACSPolicyManager, ncIACSPolicyManager)

// protected
NS_IMETHODIMP_(nsrefcnt) ncACSPolicyManager::AddRef (void)
{
    return 1;
}

// protected
NS_IMETHODIMP_(nsrefcnt) ncACSPolicyManager::Release (void)
{
    return 1;
}

AB_DEFINE_THREADSAFE_SINGLETON_NO_POOL (ncACSPolicyManager)

ncRefreshTokenThread* ncACSPolicyManager::_srefreshTokenThread = NULL;
ncActiveRecordThread* ncACSPolicyManager::_sactiveRecordThread = NULL;

ncACSPolicyManager::ncACSPolicyManager(): _ossClientPtr (0)
{
    NC_ACS_PROCESSOR_TRACE (_T("this: %p"), this);

    nsresult ret;

    _acsShareMgnt = do_CreateInstance (NC_ACS_SHARE_MGNT_CONTRACTID, &ret);
    if (NS_FAILED (ret)) {
        THROW_E (ACS_PROCESSOR, FAILED_TO_CREATE_ACS_SHAREMGNT,
            _T("Failed to create acs sharemgnt: 0x%x"), ret);
    }
    _policyEngine = do_CreateInstance(POLICY_ENGINE_CONTRACTID, &ret);
    if (NS_FAILED(ret))
    {
        THROW_E(ACS_PROCESSOR, FAILED_TO_CREATE_ACS_DEVICE_MANAGER,
                _T("Failed to create acs policy engine: 0x%x"), ret);
    }

    vector<String> vecClientTypeStr{"unknown", "ios", "android", "windows_phone", "windows",
                                    "mac_os", "web", "mobile_web", "nas", "console_web", "deploy_web", "linux", "app"};
    for (int i = 0; i < vecClientTypeStr.size(); i++)
    {
        _clientStringTypeMap.insert(make_pair(static_cast<ACSClientType>(i), vecClientTypeStr[i]));
    }

    String administrativeHost;
    String administrativePort;
    getClientByIdUrl.format(_T("http://%s:%s/clients/"),administrativeHost.getCStr(),administrativePort.getCStr());

    checkIpUrl.format (_T("http://%s:%d/api/policy-management/v1/network/allow"),
        AcsServiceAccessConfig::getInstance()->policyMgntHost.getCStr (), AcsServiceAccessConfig::getInstance()->policyMgntPort);

    //
    // token 刷新线程创建组件对象时则启动。
    //
    if (_srefreshTokenThread == NULL) {
        _srefreshTokenThread = new ncRefreshTokenThread ();
        _srefreshTokenThread->start ();
    }

    //
    // 启动活跃记录线程
    //
    if (_sactiveRecordThread == NULL) {
        _sactiveRecordThread = new ncActiveRecordThread ();
        _sactiveRecordThread->start ();
    }
}

ncACSPolicyManager::~ncACSPolicyManager()
{
    NC_ACS_PROCESSOR_TRACE (_T("this: %p"), this);
}

void ncACSPolicyManager::createOSSClient ()
{
    if (nullptr == _ossClientPtr.get ()) {
        nsresult ret;
        nsCOMPtr<ncIOSSClient> ossClient = do_CreateInstance (NC_OSS_CLIENT_CONTRACTID, &ret);
        if (NS_FAILED (ret)){
            THROW_E (ACS_PROCESSOR, FAILED_TO_CREATE_XPCOM_INSTANCE,
                     _T("Failed to create OSSClient: %s(0x%x)"),
                     String::toString ((int64)ret).getCStr (), (int64)ret);
        }
        _ossClientPtr.reset (new nsCOMPtr<ncIOSSClient>(ossClient));
    }
}

/* [notxpcom] bool CheckIp ([const] in StringRef accessorId, [const] in StringRef ip); */
NS_IMETHODIMP_(bool) ncACSPolicyManager::CheckIp (const String & accessorId, const String & ip)
{
    std::string content;
    std::string type;

    InfInt result;
    if (ip.find(_T(":")) == String::NO_POSITION) {
        struct in_addr addr;
        inet_pton(AF_INET, ip.getCStr (), &addr);
        result = ntohl(addr.s_addr);
        type = "ipv4";
    } else {
        struct in6_addr addr;
        inet_pton(AF_INET6, ip.getCStr (), &addr);
        result = InfInt(ntohl(addr.s6_addr32[0])) * InfInt("79228162514264337593543950336") + \
                InfInt(ntohl(addr.s6_addr32[1])) * InfInt("18446744073709551616") + \
                InfInt(ntohl(addr.s6_addr32[2])) * InfInt("4294967296") + \
                InfInt(ntohl(addr.s6_addr32[3]));
        type = "ipv6";
    }

    content = "{\"accessor_id\":\"" + toSTLString(accessorId) + "\"," + \
                            "\"ip\":" + result.toString() + "," + \
                            "\"ip_type\":\"" + type + "\"" + \
                        +"}";

    createOSSClient ();
    ncOSSResponse res;
    vector<string> inHeaders;

    try {
        (*_ossClientPtr)->Post (checkIpUrl.getCStr (), content, inHeaders, 30, res);
    }
    catch(Exception& e) {
        SystemLog::getInstance ()->log (__FILE__, __LINE__, ERROR_LOG_TYPE, _T("Could not connect to index server"));
        throw;
    }

    if (res.code != 200) {
        SystemLog::getInstance ()->log (__FILE__, __LINE__, ERROR_LOG_TYPE, _T("Server internal error, code:%d, cause:%s."),
                                        res.code, res.body.c_str ());
        THROW_E (ACS_PROCESSOR, ACSHTTP_SERVER_INTERNAL_ERR, _T("Code:%d. Cause:%s."), res.code, res.body.c_str ());
    }

    JSON::Value response;
    JSON::Reader::read(response, res.body.c_str (), res.body.length ());

    return response["result"].b ();
}

/* [notxpcom] bool CheckPolicy ([const] in ncPolicyCheckInfoRef policyInfo); */
NS_IMETHODIMP_(bool) ncACSPolicyManager::CheckPolicy(const ncPolicyCheckInfo & policyInfo)
{
    bool enabled = policyInfo.enabled;
    int userpriority = policyInfo.priority;
    int accountType = static_cast<int>(policyInfo.accountType);
    int clientType = static_cast<int>(policyInfo.clientType);
    String userid = policyInfo.userId;
    String clientid = policyInfo.clientId;
    String loginip = policyInfo.loginIp;
    String udid = policyInfo.udid;
    String ip = policyInfo.ip;

    // 如果客户环境为app，则跳过ip和udid检查
    if (static_cast<ACSClientType>(clientType) != ACSClientType::APP) {

        //  如果此ip与登录ip不一致
        if (!((ip == "127.0.0.1") || (ip == "localhost")) && (ip != loginip)) {
            // 判断此ip是否被限制登录, 用于处理网络切换时强制下线
            if (!CheckIp (userid, ip)) {
                THROW_E (ACS_PROCESSOR, ERR_RESTRICTED_LOGIN_IP, "此ip被限制登录");
            }
            // 判断是否发生内外网切换
            if (ncACSProcessorUtil::getInstance()->CheckSwitchingNetwork (loginip, ip)) {
                // 判断是否开启内外网切换退出功能
                String auto_logout = _acsShareMgnt->GetShareMgntConfig("switch_network_auto_logout");
                if (auto_logout.compareIgnoreCase("1") == 0 ){
                    THROW_E (ACS_PROCESSOR, LOGOUT_BECAUSE_SWITCH_NETWORK, "您所在的网络环境已改变，请重新登录。");
                    }
            }
        }

        nsresult ret;
        nsCOMPtr<ncIACSDeviceManager> _acsDeviceManager = do_CreateInstance (NC_ACS_DEVICE_MANAGER_CONTRACTID, &ret);
        if (NS_FAILED (ret)) {
            THROW_E (ACS_PROCESSOR, FAILED_TO_CREATE_ACS_DEVICE_MANAGER,
                _T("Failed to create acs device manager: 0x%x"), ret);
        }
        ncDeviceInfo deviceInfo;
        if (!udid.isEmpty () && _acsDeviceManager->GetDeviceByUDID(userid, udid, deviceInfo)) {
            // 设备是否需要数据擦除
            if(deviceInfo.eraseFlag == 1) {
                THROW_E (ACS_PROCESSOR, ERR_DEVICE_ERASED, "该设备需要数据擦除");
            }

            // 设备是否被禁用了
            if(deviceInfo.disableFlag == 1) {
                THROW_E (ACS_PROCESSOR, ERR_DEVICE_DISABLED, "该设备已被禁用");
            }
        }

        // 设备是否未绑定
        map<String, vector<String> > udidsMap;
        _acsDeviceManager->GetBindDevicesWithAllUser(userid, udidsMap);
        if(udidsMap.size() != 0) {
            if (udid.isEmpty ())
                THROW_E (ACS_PROCESSOR, ERR_DEVICE_NOT_BINDED, "该设备未绑定");

            udid.toUpper();
            bool binded = false;
            // 查找个人账号绑定设备
            auto iter = udidsMap.find (userid);
            if (iter != udidsMap.end ()) {
                if (find(iter->second.begin(), iter->second.end(), udid) != iter->second.end()) {
                    binded = true;
                }
            }
            // 查找所有用户绑定设备
            iter = udidsMap.find (toCFLString(g_ShareMgnt_constants.NCT_ALL_USER_GROUP));
            if (iter != udidsMap.end ()) {
                if (find(iter->second.begin(), iter->second.end(), udid) != iter->second.end()) {
                    binded= true;
                }
            }
            if (!binded)
                THROW_E (ACS_PROCESSOR, ERR_DEVICE_NOT_BINDED, "该设备未绑定");
        }

        // 设备IP被限制登录
        if (!CheckIp (userid, loginip)) {
                THROW_E (ACS_PROCESSOR, ERR_RESTRICTED_LOGIN_IP, "此ip被限制登录");
        }

        // 设备类型被禁止登录
        if (_policyEngine->Audit_ClientRestriction(_clientStringTypeMap[policyInfo.clientType].getCStr()))
        {
                THROW_E(ACS_PROCESSOR, ERR_LOGIN_OSTYPE_FORBID, "管理员已禁止此类客户端登录");
        }

        // 禁止身份证号登录
        vector<String> sharemgntKeys;
        map<String, String> sharemgntKvMap;
        sharemgntKeys.push_back("id_card_login_status");
        _acsShareMgnt->BatchGetConfig(sharemgntKeys, sharemgntKvMap);
        int idCardLoginStatus = Int::getValue(sharemgntKvMap["id_card_login_status"]);
        if ((idCardLoginStatus != 1) && (accountType == 1)) {
            THROW_E (ACS_PROCESSOR, ERR_LOGIN_IDCARD_FORBID, "ID number login has been closed by your administrator, please log in again.");
        }
    }

    // 用户已被禁用
    if(!enabled){
        THROW_E (ACS_PROCESSOR, ERR_DISABLED_USER, LOAD_STRING (_T("IDS_DISABLED_USER")));
    }

    ptime localCur (BusinessDate::getLocalTime ());
    date localCurDay = localCur.date();
    time_duration localCurTime = localCur.time_of_day();

    String curTimeStr;
    curTimeStr.format (_T("%04d-%02d-%02d %02d:%02d:%02d"),
        (int)localCurDay.year (), (int)localCurDay.month (), (int)localCurDay.day (),
        (int)localCurTime.hours (), (int)localCurTime.minutes (), (int)localCurTime.seconds ());

#ifndef __UT__
    // 所有的token都会刷新用户活跃时间
    _srefreshTokenThread->pushRefreshInfo(userid, curTimeStr, true);
#endif

    // 向活跃记录线程推送用户ID和时间
    curTimeStr.format (_T("%04d-%02d-%02d"),
        (int)localCurDay.year (), (int)localCurDay.month (), (int)localCurDay.day ());
    _sactiveRecordThread->pushActiveUserInfo(userid, curTimeStr);

    NC_ACS_PROCESSOR_TRACE (_T("this: %p, userId: %s end, ok"), this, userid.getCStr ());
    return true;
}
