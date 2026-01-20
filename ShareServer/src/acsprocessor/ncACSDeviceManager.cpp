#include <abprec.h>
#include <ncutil/ncBusinessDate.h>
#include <dataapi/ncJson.h>
#include "acsprocessor.h"
#include "ncACSDeviceManager.h"
#include "gen-cpp/ShareMgnt_constants.h"
#include "ncACSProcessorUtil.h"
#include <drivenadapter/public/common.h>

/* Implementation file */
NS_IMPL_QUERY_INTERFACE1 (ncACSDeviceManager, ncIACSDeviceManager)

// protected
NS_IMETHODIMP_(nsrefcnt) ncACSDeviceManager::AddRef (void)
{
    return 1;
}

// protected
NS_IMETHODIMP_(nsrefcnt) ncACSDeviceManager::Release (void)
{
    return 1;
}

AB_DEFINE_THREADSAFE_SINGLETON_NO_POOL (ncACSDeviceManager)

ncACSDeviceManager::ncACSDeviceManager()
{
    NC_ACS_PROCESSOR_TRACE("");

    nsresult ret;
    _dbDeviceManager = do_CreateInstance (NC_DB_DEVICE_MANAGER_CONTRACTID, &ret);
    if (NS_FAILED (ret)) {
        THROW_E (ACS_PROCESSOR, FAILED_TO_CREATE_DB_DEVICE_MANANGER,
            _T("Failed to create db device manager: 0x%x"), ret);
    }

    _dbTokenManager = do_CreateInstance (NC_DB_TOKEN_MANAGER_CONTRACTID, &ret);
    if (NS_FAILED (ret)) {
        THROW_E (ACS_PROCESSOR, FAILED_TO_CREATE_DB_TOKEN_MANANGER,
            _T("Failed to create db token manager: 0x%x"), ret);
    }

    _acsShareMgnt = do_CreateInstance (NC_ACS_SHARE_MGNT_CONTRACTID, &ret);
    if (NS_FAILED (ret)) {
        THROW_E (ACS_PROCESSOR, FAILED_TO_CREATE_ACS_SHAREMGNT,
            _T("Failed to create acs sharemgnt: 0x%x"), ret);
    }

    _acsTokenManager = do_CreateInstance (NC_ACS_TOKEN_MANAGER_CONTRACTID, &ret);
    if (NS_FAILED (ret)) {
        THROW_E(ACS_PROCESSOR, FAILED_TO_CREATE_ACS_TOKEN_MANANGER,
            _T("Failed to create acs tokenmanager: 0x%x"), ret);
    }

    _hydra = do_CreateInstance (HYDRA_CONTRACTID, &ret);
    if (NS_FAILED (ret)) {
        THROW_E(ACS_PROCESSOR, FAILED_TO_CREATE_DRIVENADAPTER_MANANGER,
            _T("Failed to create hydra adapter instance: 0x%x"), ret);
    }
}

ncACSDeviceManager::~ncACSDeviceManager()
{
  /* destructor code */
}

/* [notxpcom] void GetDevicesByUserIdAndUdid ([const] in StringRef userId, [const] in StringRef udid, in ncDeviceInfoVecRef deviceInfos, in int start, in int limit, in bool authFlag); */
NS_IMETHODIMP_(void) ncACSDeviceManager::GetDevicesByUserIdAndUdid(const String & userId, const String& udid, vector<ncDeviceInfo> & deviceInfos, int start, int limit, bool authFlag)
{
    NC_ACS_PROCESSOR_TRACE("userId: %s begin", userId.getCStr());
    deviceInfos.clear();

    vector<dbDeviceInfo> dbInfos;
    _dbDeviceManager->GetDevicesByUserIdAndUdid(userId, udid, dbInfos, start, limit);

    vector<ncTokenIntrospectInfo> tokenInfos;
    set<String> udids;
    if (!authFlag) {
        if (!userId.isEmpty ()) {
            _hydra->GetConsentInfo(userId, tokenInfos);
        }

        for (auto iter = tokenInfos.begin (); iter != tokenInfos.end (); ++iter) {
            udids.insert(iter->udid);
        }
    }

    ncDeviceInfo tmpInfo;
    for(size_t i = 0; i < dbInfos.size(); ++i) {
        tmpInfo.baseInfo.udid = dbInfos[i].baseInfo.udid;
        tmpInfo.baseInfo.name = dbInfos[i].baseInfo.name;
        tmpInfo.baseInfo.clientType = static_cast<ACSClientType>(dbInfos[i].baseInfo.osType);
        tmpInfo.baseInfo.deviceType = dbInfos[i].baseInfo.deviceType;
        tmpInfo.baseInfo.lastLoginIp = dbInfos[i].baseInfo.lastLoginIp;
        tmpInfo.baseInfo.lastLoginTime = dbInfos[i].baseInfo.lastLoginTime;

        tmpInfo.eraseFlag = dbInfos[i].eraseFlag;
        tmpInfo.lastEraseTime = dbInfos[i].lastEraseTime;
        tmpInfo.disableFlag = dbInfos[i].disableFlag;
        tmpInfo.bindFlag = dbInfos[i].bindFlag;
        if (!authFlag) {
            tmpInfo.loginFlag = udids.count(tmpInfo.baseInfo.udid);
        }

        deviceInfos.push_back(tmpInfo);
    }

    NC_ACS_PROCESSOR_TRACE("userId: %s end", userId.getCStr());
}

/* [notxpcom] void GetLoginedMobileDevices ([const] in StringRef userId, in ncDeviceInfoVecRef deviceInfos); */
NS_IMETHODIMP_(void) ncACSDeviceManager::GetLoginedMobileDevices(const String & userId, vector<ncDeviceInfo> & deviceInfos)
{
    NC_ACS_PROCESSOR_TRACE("userId: %s begin", userId.getCStr());
    deviceInfos.clear();

    vector<dbDeviceInfo> dbInfos;
    _dbDeviceManager->GetDevicesByUserIdAndUdid(userId, String::EMPTY, dbInfos, 0, -1);

    vector<ncTokenIntrospectInfo> tokenInfos;
    if (!userId.isEmpty ()) {
        _hydra->GetConsentInfo(userId, tokenInfos);
    }

    set<String> udids;
    for (auto iter = tokenInfos.begin (); iter != tokenInfos.end (); ++iter) {
        udids.insert(iter->udid);
    }

    // 获取用户登录过的移动设备
    ncDeviceInfo tmpInfo;
    for(size_t i = 0; i < dbInfos.size(); ++i) {
        if(dbInfos[i].baseInfo.lastLoginTime == -1) {
            continue;
        }
        if(dbInfos[i].baseInfo.osType < static_cast<int>(ACSClientType::IOS) || dbInfos[i].baseInfo.osType > static_cast<int>(ACSClientType::WINDOWS_PHONE)) {
            continue;
        }

        tmpInfo.baseInfo.udid = dbInfos[i].baseInfo.udid;
        tmpInfo.baseInfo.name = dbInfos[i].baseInfo.name;
        tmpInfo.baseInfo.clientType = static_cast<ACSClientType>(dbInfos[i].baseInfo.osType);
        tmpInfo.baseInfo.deviceType = dbInfos[i].baseInfo.deviceType;
        tmpInfo.baseInfo.lastLoginIp = dbInfos[i].baseInfo.lastLoginIp;
        tmpInfo.baseInfo.lastLoginTime = dbInfos[i].baseInfo.lastLoginTime;

        tmpInfo.eraseFlag = dbInfos[i].eraseFlag;
        tmpInfo.lastEraseTime = dbInfos[i].lastEraseTime;
        tmpInfo.disableFlag = dbInfos[i].disableFlag;
        tmpInfo.bindFlag = dbInfos[i].bindFlag;

        tmpInfo.loginFlag = udids.count(tmpInfo.baseInfo.udid);

        deviceInfos.push_back(tmpInfo);
    }

    NC_ACS_PROCESSOR_TRACE("userId: %s end", userId.getCStr());
}

/* [notxpcom] void RecordDevice ([const] in StringRef userId, [const] in ncDeviceBaseInfoRef baseInfo); */
NS_IMETHODIMP_(void) ncACSDeviceManager::RecordDevice(const String & userId, const ncDeviceBaseInfo & baseInfo)
{
    NC_ACS_PROCESSOR_TRACE("userId: %s, udid: %s begin", userId.getCStr(), baseInfo.udid.getCStr());

    dbDeviceBaseInfo dbInfo;
    dbInfo.udid = baseInfo.udid;
    dbInfo.name = baseInfo.name;
    dbInfo.osType = static_cast<int>(baseInfo.clientType);
    dbInfo.deviceType = baseInfo.deviceType;
    dbInfo.lastLoginIp = baseInfo.lastLoginIp;
    dbInfo.lastLoginTime = baseInfo.lastLoginTime;

    dbDeviceInfo retInfo;
    if(_dbDeviceManager->GetDeviceByUDID(userId, baseInfo.udid, retInfo)) {
        _dbDeviceManager->UpdateDevice(userId, dbInfo);
    }
    else {
        _dbDeviceManager->AddDevice(userId, dbInfo);
    }

    String allUserGroupId = toCFLString(g_ShareMgnt_constants.NCT_ALL_USER_GROUP);
    dbDeviceInfo allUserRetInfo;
    if(_dbDeviceManager->GetDeviceByUDID(allUserGroupId, baseInfo.udid, allUserRetInfo)) {
        _dbDeviceManager->UpdateDevice(allUserGroupId, dbInfo);
    }

    NC_ACS_PROCESSOR_TRACE("userId: %s, udid: %s end", userId.getCStr(), baseInfo.udid.getCStr());
}

/* [notxpcom] void AddDevice ([const] in StringRef userId, [const] in ncDeviceBaseInfoRef baseInfo); */
NS_IMETHODIMP_(void) ncACSDeviceManager::AddDevice(const String & userId, const ncDeviceBaseInfo & baseInfo)
{
    NC_ACS_PROCESSOR_TRACE("userId: %s, udid: %s begin", userId.getCStr(), baseInfo.udid.getCStr());

    dbDeviceBaseInfo dbInfo;
    dbInfo.udid = baseInfo.udid;
    dbInfo.name = baseInfo.name;
    dbInfo.osType = static_cast<int>(baseInfo.clientType);
    dbInfo.deviceType = baseInfo.deviceType;
    dbInfo.lastLoginIp = baseInfo.lastLoginIp;
    dbInfo.lastLoginTime = baseInfo.lastLoginTime;

    dbDeviceInfo retInfo;
    if(_dbDeviceManager->GetDeviceByUDID(userId, baseInfo.udid, retInfo)) {
        THROW_E (ACS_PROCESSOR, ERR_DEVICE_UDID_EXISTS, LOAD_STRING("IDS_DEVICE_UDID_EXISTS"));
    }

    _dbDeviceManager->AddDevice(userId, dbInfo);

    vector<String> udids;
    udids.push_back(baseInfo.udid);
    SetBindStatus(userId, udids);

    NC_ACS_PROCESSOR_TRACE("userId: %s, udid: %s end", userId.getCStr(), baseInfo.udid.getCStr());
}

/*[notxpcom] void AddDevices([const] in StringRef userId, [const] in dbDeviceBaseInfoVecRef baseInfos);*/
NS_IMETHODIMP_(void) ncACSDeviceManager::AddDevices(const String & userId, const vector<ncDeviceBaseInfo> & baseInfos)
{
    NC_ACS_PROCESSOR_TRACE("userId: %s begin", userId.getCStr());

    // 获取已配置的设备ID
    vector<String> addedDeviceUdids;
    vector<dbDeviceInfo> addedDeviceInfos;
    _dbDeviceManager->GetDevicesByUserIdAndUdid (userId, String::EMPTY, addedDeviceInfos, 0, -1);
    for (auto iter = addedDeviceInfos.begin(); iter != addedDeviceInfos.end(); iter++) {
        String udid = iter->baseInfo.udid;
        udid.toUpper();
        addedDeviceUdids.push_back(udid);
    }

    // 过滤输入参数，只添加未配置过的设备
    dbDeviceBaseInfo dbInfo;
    vector<dbDeviceBaseInfo> dbInfos;
    vector<String> newUdids;
    for (auto iter = baseInfos.begin (); iter != baseInfos.end (); iter++) {
        if (std::find(addedDeviceUdids.begin(), addedDeviceUdids.end(), iter->udid) == addedDeviceUdids.end()) {
            dbInfo.udid = iter->udid;
            dbInfo.name = iter->name;
            dbInfo.osType = static_cast<int>(iter->clientType);
            dbInfo.deviceType = iter->deviceType;
            dbInfo.lastLoginIp = iter->lastLoginIp;
            dbInfo.lastLoginTime = iter->lastLoginTime;
            dbInfos.push_back (dbInfo);
            newUdids.push_back(iter->udid);
        }
    }

    _dbDeviceManager->AddDevices (userId, dbInfos);

    SetBindStatus (userId, newUdids);

    // 获取用户显示名
    String displayName;
    String account;
    _acsShareMgnt->GetUserName (userId, displayName, account);

    String msg;
    for (auto iter = newUdids.begin(); iter != newUdids.end(); iter++) {
        msg.format (LOAD_STRING (_T("IDS_USER_BIND_MAC_ADDR_MSG")), displayName.getCStr (), (*iter).getCStr ());

        ncACSProcessorUtil::getInstance ()->Log (toCFLString(g_ShareMgnt_constants.NCT_USER_ADMIN),
                                                ncTokenVisitorType::REALNAME,
                                                ncTLogType::NCT_LT_MANAGEMENT,
                                                ncTLogLevel::NCT_LL_INFO,
                                                ncTManagementType::NCT_MNT_CREATE,
                                                msg,"", "127.0.0.1");
    }

    NC_ACS_PROCESSOR_TRACE("userId: %s end", userId.getCStr());
}

/*[notxpcom] void GetUserNameByUdid([const] in StringRef udid, in VecStringRef users, in int start, in int limit);*/
NS_IMETHODIMP_(void) ncACSDeviceManager::GetUserNameByUdid(const String& udid, vector<String>& users, int start, int limit)
{
    NC_ACS_PROCESSOR_TRACE("udid: %s begin", udid.getCStr());

    _dbDeviceManager->GetUserNameByUdid (udid, users, start, limit);

    NC_ACS_PROCESSOR_TRACE("udid: %s end", udid.getCStr());
}

/* [notxpcom] void DeleteDevices ([const] in StringRef userId, [const] in VecStringRef udids); */
NS_IMETHODIMP_(void) ncACSDeviceManager::DeleteDevices(const String & userId, const vector<String> & udids)
{
    NC_ACS_PROCESSOR_TRACE("userId: %s begin", userId.getCStr());

    _dbDeviceManager->DeleteDevices(userId, udids);

    NC_ACS_PROCESSOR_TRACE("userId: %s end", userId.getCStr());
}

/* [notxpcom] bool GetDeviceByUDID ([const] in StringRef userId, [const] in StringRef udid, in ncDeviceInfoRef deviceInfo); */
NS_IMETHODIMP_(bool) ncACSDeviceManager::GetDeviceByUDID(const String & userId, const String & udid, ncDeviceInfo & deviceInfo)
{
    dbDeviceInfo dbInfo;
    bool ret = _dbDeviceManager->GetDeviceByUDID(userId, udid, dbInfo);
    if(ret) {
        deviceInfo.baseInfo.udid = dbInfo.baseInfo.udid;
        deviceInfo.baseInfo.name = dbInfo.baseInfo.name;
        deviceInfo.baseInfo.clientType = static_cast<ACSClientType>(dbInfo.baseInfo.osType);
        deviceInfo.baseInfo.deviceType = dbInfo.baseInfo.deviceType;
        deviceInfo.baseInfo.lastLoginIp = dbInfo.baseInfo.lastLoginIp;
        deviceInfo.baseInfo.lastLoginTime = dbInfo.baseInfo.lastLoginTime;

        deviceInfo.eraseFlag = dbInfo.eraseFlag;
        deviceInfo.lastEraseTime = dbInfo.lastEraseTime;
        deviceInfo.disableFlag = dbInfo.disableFlag;
        deviceInfo.bindFlag = dbInfo.bindFlag;

        // 上层调用均未使用loginFlag，因此注释该值获取，避免调用耗时大量增加
        // deviceInfo.loginFlag = _acsTokenManager->HasTokenByUDID(userId, udid);
    }

    return ret;
}

/* [notxpcom] void SetEraseStatus ([const] in StringRef userId, [const] in StringRef udid); */
NS_IMETHODIMP_(void) ncACSDeviceManager::SetEraseStatus(const String & userId, const String & udid)
{
    NC_ACS_PROCESSOR_TRACE("userId: %s, udid: %s begin", userId.getCStr(), udid.getCStr());

    _dbDeviceManager->SetEraseStatus(userId, udid);

    NC_ACS_PROCESSOR_TRACE("userId: %s, udid: %s begin", userId.getCStr(), udid.getCStr());
}

/* [notxpcom] void SetEraseSucInfo([const] in StringRef userId, [const] in StringRef udid, in int64 date); */
NS_IMETHODIMP_(void) ncACSDeviceManager::SetEraseSucInfo(const String & userId, const String & udid, int64 date)
{
    NC_ACS_PROCESSOR_TRACE("userId: %s, udid: %s begin", userId.getCStr(), udid.getCStr());
    if(userId.isEmpty() || udid.isEmpty()) {
        return;
    }

    _dbDeviceManager->SetEraseSucInfo(userId, udid, date);

    vector<ncTokenIntrospectInfo> tokenInfos;
    _hydra->GetConsentInfo(userId, tokenInfos);

    for(size_t i = 0; i < tokenInfos.size(); ++i) {
        if(tokenInfos[i].udid == udid){
            _hydra->DeleteConsentAndLogin(tokenInfos[i].clientId, userId);
        }
    }
    NC_ACS_PROCESSOR_TRACE("userId: %s, udid: %s begin", userId.getCStr(), udid.getCStr());
}

/* [notxpcom] void SetDisableStatus ([const] in StringRef userId, [const] in StringRef udid); */
NS_IMETHODIMP_(void) ncACSDeviceManager::SetDisableStatus(const String & userId, const String & udid)
{
    NC_ACS_PROCESSOR_TRACE("userId: %s, udid: %s begin", userId.getCStr(), udid.getCStr());

    _dbDeviceManager->SetDisableStatus(userId, udid);

    NC_ACS_PROCESSOR_TRACE("userId: %s, udid: %s begin", userId.getCStr(), udid.getCStr());
}

/* [notxpcom] void SetEnableStatus ([const] in StringRef userId, [const] in StringRef udid); */
NS_IMETHODIMP_(void) ncACSDeviceManager::SetEnableStatus(const String & userId, const String & udid)
{
    NC_ACS_PROCESSOR_TRACE("userId: %s, udid: %s begin", userId.getCStr(), udid.getCStr());

    _dbDeviceManager->SetEnableStatus(userId, udid);

    NC_ACS_PROCESSOR_TRACE("userId: %s, udid: %s begin", userId.getCStr(), udid.getCStr());
}

/* [notxpcom] void SetBindStatus ([const] in StringRef userId, [const] in VecStringRef udids); */
NS_IMETHODIMP_(void) ncACSDeviceManager::SetBindStatus(const String & userId, const vector<String> & udids)
{
    NC_ACS_PROCESSOR_TRACE("userId: %s begin", userId.getCStr());

    _dbDeviceManager->SetBindStatus(userId, udids);

    NC_ACS_PROCESSOR_TRACE("userId: %s begin", userId.getCStr());
}

/* [notxpcom] void SetUnbindStatus ([const] in StringRef userId, [const] in StringRef udid); */
NS_IMETHODIMP_(void) ncACSDeviceManager::SetUnbindStatus(const String & userId, const String & udid)
{
    NC_ACS_PROCESSOR_TRACE("userId: %s, udid: %s begin", userId.getCStr(), udid.getCStr());

    _dbDeviceManager->SetUnbindStatus(userId, udid);

    NC_ACS_PROCESSOR_TRACE("userId: %s, udid: %s begin", userId.getCStr(), udid.getCStr());
}

/* [notxpcom] void GetBindDevices([const] in StringRef userId, in StringIntMapRef udidsMap); */
NS_IMETHODIMP_(void) ncACSDeviceManager::GetBindDevices(const String & userId, map<String, int>& udidsMap)
{
    NC_ACS_PROCESSOR_TRACE("userId: %s begin", userId.getCStr());

    _dbDeviceManager->GetBindDevices(userId, udidsMap);

    NC_ACS_PROCESSOR_TRACE("userId: %s begin", userId.getCStr());
}

/* [notxpcom] void GetBindDevicesWithAllUser([const] in StringRef userIds, in StringVecMapRef udidsMap); */
NS_IMETHODIMP_(void) ncACSDeviceManager::GetBindDevicesWithAllUser(const String & userId, map<String, vector<String> >& udidsMap)
{
    NC_ACS_PROCESSOR_TRACE("userId: %s begin", userId.getCStr());

    _dbDeviceManager->GetBindDevicesWithAllUser(userId, udidsMap);

    NC_ACS_PROCESSOR_TRACE("userId: %s begin", userId.getCStr());
}

/* [notxpcom] void ForceLogOff ([const] in StringRef userId); */
NS_IMETHODIMP_(void) ncACSDeviceManager::ForceLogOff(const String & userId)
{
    NC_ACS_PROCESSOR_TRACE("userId: %s begin", userId.getCStr());

    vector<ncTokenIntrospectInfo> tokenInfos;
    _hydra->GetConsentInfo(userId, tokenInfos);

    for(size_t i = 0; i < tokenInfos.size(); ++i){
        if(static_cast<ACSClientType>(tokenInfos[i].clientType) == ACSClientType::WINDOWS){
            _hydra->DeleteConsentAndLogin(tokenInfos[i].clientId, userId);
        }
    }

    NC_ACS_PROCESSOR_TRACE("userId: %s begin", userId.getCStr());
}

/* [notxpcom] void AllForceLogOff ([const] in StringRef userId); */
NS_IMETHODIMP_(void) ncACSDeviceManager::AllForceLogOff(const String & userId)
{
    NC_ACS_PROCESSOR_TRACE("userId: %s begin", userId.getCStr());

    vector<ncTokenIntrospectInfo> tokenInfos;
    _hydra->GetConsentInfo(userId, tokenInfos);

    for(size_t i = 0; i < tokenInfos.size(); ++i){
        if(static_cast<ACSClientType>(tokenInfos[i].clientType) == ACSClientType::WINDOWS ||
        static_cast<ACSClientType>(tokenInfos[i].clientType) == ACSClientType::MAC_OS ||
        static_cast<ACSClientType>(tokenInfos[i].clientType) == ACSClientType::WEB ||
        static_cast<ACSClientType>(tokenInfos[i].clientType) == ACSClientType::LINUX) {
            _hydra->DeleteConsentAndLogin(tokenInfos[i].clientId, userId);
        }
    }

    NC_ACS_PROCESSOR_TRACE("userId: %s begin", userId.getCStr());
}

/*[notxpcom] bool UserHasDeviceInfo ([const] in StringRef userId, in bool binded);*/
NS_IMETHODIMP_(bool) ncACSDeviceManager::UserHasDeviceInfo (const String& userId, bool binded)
{
    NC_ACS_PROCESSOR_TRACE("userId: %s begin", userId.getCStr ());

    bool exists = _dbDeviceManager->UserHasDeviceInfo (userId, binded);

    NC_ACS_PROCESSOR_TRACE("userId: %s begin", userId.getCStr ());

    return exists;
}
