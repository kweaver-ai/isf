#include <iconv.h>
#include <algorithm>
#include <abprec.h>
#include <ncutil/ncBusinessDate.h>
#include <boost/regex.hpp>
#include <dataapi/ncGNSUtil.h>
#include <ethriftutil/ncThriftClient.h>

#include "gen-cpp/ncTEACP.h"
#include "gen-cpp/EThriftException_types.h"

#include "gen-cpp/ncTEVFS.h"
#include "gen-cpp/EVFS_constants.h"
#include "gen-cpp/ncTShareMgnt.h"
#include "gen-cpp/ShareMgnt_constants.h"

#include <dataapi/ncJson.h>

#include <boost/property_tree/ptree.hpp>
#include <boost/property_tree/ini_parser.hpp>
using namespace boost::property_tree;

#include "teacpserver.h"
#include "ncTEACPHandler.h"
#include "eachttpserver/eacServiceAccessConfig.h"

#define THRIFT_EVFS_DEFAULT_IP            "localhost"
#define PUSH_MESSAGE_TO_THIRD_PARTY       "push_message_to_third_party"
#define ENABLE_SEND_SHARE_MAIL            "enable_send_share_mail"
#define ENABLE_MESSAGE_NOTIFY             "enable_message_notify"
#define CUSTOME_APPLICATION_CONFIG        "custome_application_config"
#define DEFAULT_EXPIRED_INTERVAL          180
#define NONE_EXPIRED                      -1
#define ALL_DEPART_LEVEL                  -1
#define MIN_DOC_NAME_LEVEL                1
#define MAX_DOC_NAME_LEVEL                4
#define SERVICE_ACCESS_FILE               _T("/sysvol/conf/service_conf/service_access.conf")

#define THROW_EACP_T_EXCEPTION(tErrProviderName, tErrId, tErrdetail, args...)   \
    do {                                                                        \
        ncTException ___te___;                                                  \
        ___te___.expType = ncTExpType::NCT_FATAL;                               \
        ___te___.fileName = __FILE__;                                           \
        ___te___.codeLine = __LINE__;                                           \
        ___te___.errProvider = tErrProviderName;                                \
        ___te___.errID = tErrId;                                                \
        ___te___.errDetail = tErrdetail;                                        \
        String ___tErrMsg___;                                                   \
        ___tErrMsg___.format (args);                                            \
        ___te___.expMsg = toSTLString (___tErrMsg___);                          \
        throw ___te___;                                                         \
    } while (false)

ncTEACPHandler::ncTEACPHandler (): _ossClientPtr (0), _nsqManager (0)
{
    NC_T_EACP_SERVER_TRACE (_T("this: %p"), this);

    nsresult ret;

    _dbLockManager = do_CreateInstance (NC_DB_LOCK_MANAGER_CONTRACTID, &ret);
    if (NS_FAILED (ret)) {
        THROW_E (T_EACP_SERVER_ERR_PROVIDER_NAME, FAILED_TO_CREATE_DB_LOCK_MANAGER, _T("Failed to create db lock manager: 0x%x"), ret);
    }

    _acsPermManager = do_CreateInstance (NC_ACS_PERM_MANAGER_CONTRACTID, &ret);
    if (NS_FAILED (ret)) {
        THROW_E(T_EACP_SERVER_ERR_PROVIDER_NAME, FAILED_TO_CREATE_ACS_PERM_MANANGER, _T("Failed to create acs perm manager: 0x%x"), ret);
    }

    _acsTokenManager = do_CreateInstance (NC_ACS_TOKEN_MANAGER_CONTRACTID, &ret);
    if (NS_FAILED (ret)) {
        THROW_E(T_EACP_SERVER_ERR_PROVIDER_NAME, FAILED_TO_CREATE_ACS_TOKEN_MANANGER, _T("Failed to create acs token manager: 0x%x"), ret);
    }

    _acsOwnerManager = do_CreateInstance (NC_ACS_OWNER_MANAGER_CONTRACTID, &ret);
    if (NS_FAILED (ret)) {
        THROW_E(T_EACP_SERVER_ERR_PROVIDER_NAME, FAILED_TO_CREATE_ACS_OWNER_MANANGER, _T("Failed to create acs owner manager: 0x%x"), ret);
    }

    _acsShareMgnt = do_CreateInstance (NC_ACS_SHARE_MGNT_CONTRACTID, &ret);
    if (NS_FAILED (ret)) {
        THROW_E(T_EACP_SERVER_ERR_PROVIDER_NAME, FAILED_TO_CREATE_ACS_SHAREMGNT, _T("Failed to create acs sharemgnt manager: 0x%x"), ret);
    }

    _acsLicenseManager = do_CreateInstance (NC_ACS_LICENSE_MANAGER_CONTRACTID, &ret);
    if (NS_FAILED (ret)) {
        THROW_E(T_EACP_SERVER_ERR_PROVIDER_NAME, FAILED_TO_CREATE_ACS_LICENSE_MANAGER, _T("Failed to create acs license manager: 0x%x"), ret);
    }

    _acsLockManager = do_CreateInstance (NC_ACS_LOCK_MANAGER_CONTRACTID, &ret);
    if (NS_FAILED (ret)) {
        THROW_E(T_EACP_SERVER_ERR_PROVIDER_NAME, FAILED_TO_CREATE_ACS_LOCK_MANAGER, _T("Failed to create acs lock manager: 0x%x"), ret);
    }

    _acsDeviceManager = do_CreateInstance (NC_ACS_DEVICE_MANAGER_CONTRACTID, &ret);
    if (NS_FAILED (ret)) {
        THROW_E(T_EACP_SERVER_ERR_PROVIDER_NAME, FAILED_TO_CREATE_ACS_DEVICE_MANAGER, _T("Failed to create acs device manager: 0x%x"), ret);
    }

    _dbConfManager = do_CreateInstance (NC_DB_CONF_MANAGER_CONTRACTID, &ret);
    if (NS_FAILED (ret)) {
        THROW_E (T_EACP_SERVER_ERR_PROVIDER_NAME, FAILED_TO_CREATE_DB_CONF_MANANGER, _T("Failed to create db conf manager: 0x%x"), ret);
    }

    _acsMessageManager = do_CreateInstance (NC_ACS_MESSAGE_MANAGER_CONTRACTID, &ret);
    if (NS_FAILED (ret)) {
        THROW_E (T_EACP_SERVER_ERR_PROVIDER_NAME, FAILED_TO_CREATE_ACS_MESSAGE_MANANGER, _T("Failed to create acs message manager: 0x%x"), ret);
    }

    try {
        ptree pt;
        boost::property_tree::ini_parser::read_ini("/sysvol/conf/service_conf/service_access.conf", pt);

        string host = pt.get<string>("doc-share.privateHost");
        int port = pt.get<int>("doc-share.privatePort");
        _docShareAddr.format(_T("%s:%d"), host.c_str (), port);

        SystemLog::getInstance ()->log(__FILE__, __LINE__, INFORMATION_LOG_TYPE,
                                    "doc-share.server_addr = %s", _docShareAddr.getCStr ());
    }
    catch (ptree_error& e) {
        SystemLog::getInstance ()->log (__FILE__, __LINE__, ERROR_LOG_TYPE,
                                        "Get doc-share.server_addr failed");
    }
}

ncTEACPHandler::ncTEACPHandler(ncIDBLockManager* dbLockManager,
                               ncIACSPermManager* acsPermManager,
                               ncIACSTokenManager* acsTokenManager,
                               ncIACSOwnerManager* acsOwnerManager,
                               ncIACSShareMgnt* acsShareMgnt)
    : _dbLockManager (dbLockManager),
    _acsPermManager (acsPermManager),
    _acsTokenManager (acsTokenManager),
    _acsOwnerManager (acsOwnerManager),
    _acsShareMgnt (acsShareMgnt),
    _ossClientPtr (0),
    _nsqManager (0)
{
}

ncTEACPHandler::~ncTEACPHandler ()
{
    NC_T_EACP_SERVER_TRACE (_T("this: %p"), this);
}

void ncTEACPHandler::createOSSClient ()
{
    if (nullptr == _ossClientPtr.get ()) {
        nsresult ret;
        nsCOMPtr<ncIOSSClient> ossClient = do_CreateInstance (NC_OSS_CLIENT_CONTRACTID, &ret);
        if (NS_FAILED (ret)){
            THROW_E (T_EACP_SERVER_ERR_PROVIDER_NAME, FAILED_TO_CREATE_OSS_CLIENT,
                     _T("Failed to create OSSClient: 0x%x"), ret);
        }
        _ossClientPtr.reset (new nsCOMPtr<ncIOSSClient>(ossClient));
    }
}

void ncTEACPHandler::createNsqManager ()
{
    if (!_nsqManager) {
        AutoLock<ThreadMutexLock> autoLock (&_nsqManagerLock);
        if (!_nsqManager) {
            nsresult ret;
            _nsqManager = do_CreateInstance (NSQ_CONTRACTID, &ret);
            if (NS_FAILED (ret)) {
                THROW_E (T_EACP_SERVER_ERR_PROVIDER_NAME, ncTEACPError::NCT_UNKNOWN_ERROR,
                    _T("Failed to create nsq instance: 0x%x"), ret);
            }
        }
    }
}

void ncTEACPHandler::EACP_OnDeleteUser(const std::string& userId)
{
    NC_T_EACP_SERVER_TRACE (_T("this: %p, userId: %s begin"), this, userId.c_str ());

    try {
        // 删除用户时，不删除对应的个人文档，删除动作由关闭文档负责

        // 删除用户id对应的所有者信息
        _acsOwnerManager->DeleteOwnerByUserId (toCFLString (userId));

        // 删除用户对应的token信息
        _acsTokenManager->DeleteTokenByUserId (toCFLString (userId));

        // 删除用户id对应的权限配置
        _acsPermManager->DeleteCustomPermByUserId (toCFLString (userId));

        // 删除用户对应的文件锁
        _dbLockManager->DeleteByUserId (toCFLString (userId));

    }
    catch (Exception& e) {
        ncTException te;
        convertException (e, te);
        throw te;
    }

    NC_T_EACP_SERVER_TRACE (_T("this: %p, userId: %s end"), this, userId.c_str ());
}

void ncTEACPHandler::EACP_OnDeleteDepartment(const std::vector<std::string> & departmentIds)
{
    NC_T_EACP_SERVER_TRACE (_T("this: %p, departmentIds size: %d begin"), this, (int)departmentIds.size ());

    try {
        for (size_t i = 0; i < departmentIds.size (); ++i) {
            // 删除用户id对应的权限配置
            String depId = toCFLString (departmentIds[i]);
            _acsPermManager->DeleteCustomPermByUserId (depId);
        }
    }
    catch (Exception& e) {
        ncTException te;
        convertException (e, te);
        throw te;
    }

    NC_T_EACP_SERVER_TRACE (_T("this: %p, departmentIds size: %d end"), this, (int)departmentIds.size ());
}

void ncTEACPHandler::EACP_GetLicenseInfo(std::string& _return, const std::string& license)
{
    NC_T_EACP_SERVER_TRACE (_T("this: %p begin"), this);;

    try {
        String tmpLicense = toCFLString (license);
        String info;

        int ret = _acsLicenseManager->GetLicenseInfo (tmpLicense, info);

        if (ret == 0) {
            _return = toSTLString (info);
        }
        else if (ret == 3){
            throw Exception (LOAD_STRING (_T("IDS_INVALID_LICENSE")));
        }
        else if (ret == 4) {
            throw Exception (LOAD_STRING (_T("IDS_INVALID_ACTIVE_STRING")));
        }
        else {
            throw Exception (LOAD_STRING (_T("IDS_UNKNOWN_ERROR")));
        }
    }
    catch (Exception& e) {
        ncTException te;
        convertException (e, te);
        throw te;
    }

    NC_T_EACP_SERVER_TRACE (_T("this: %p end"), this);
}

int32_t
ncTEACPHandler::EACP_VerifyActiveCode(const std::string& license, const std::string& machineCode, const std::string& activeCode)
{
    NC_T_EACP_SERVER_TRACE (_T("this: %p begin"), this);;

    try {
        String tmpLicense = toCFLString (license);
        String tmpMachineCode = toCFLString(machineCode);
        String tmpActiveCode = toCFLString(activeCode);

        int ret = _acsLicenseManager->VerifyActiveCode (tmpLicense, tmpMachineCode, tmpActiveCode);

        return ret;
    }
    catch (Exception& e) {
        ncTException te;
        convertException (e, te);
        throw te;
    }

    NC_T_EACP_SERVER_TRACE (_T("this: %p end"), this);
}

void ncTEACPHandler::EACP_GetAutolockConfig(ncTAutolockConfig& _return)
{
    NC_T_EACP_SERVER_TRACE (_T("this: %p begin"), this);

    try {
        _return.isEnable = _acsLockManager->IsAutolockEnabled();
        _return.expiredInterval = _acsLockManager->GetExpiredInterval();
    }
    catch (Exception& e) {
        ncTException te;
        convertException (e, te);
        throw te;
    }

    NC_T_EACP_SERVER_TRACE (_T("this: %p end"), this);
}

void ncTEACPHandler::EACP_SetAutolockConfig(const ncTAutolockConfig& config)
{
    NC_T_EACP_SERVER_TRACE (_T("this: %p begin"), this);

    try {
        if (config.isEnable && config.expiredInterval < DEFAULT_EXPIRED_INTERVAL && config.expiredInterval != NONE_EXPIRED) {
            THROW_E (T_EACP_SERVER_ERR_PROVIDER_NAME, ncTEACPError::NCT_INVALID_EXPIRED_INTERVAL,
                LOAD_STRING("IDS_INVALID_EXPIRED_INTERVAL").getCStr());
        }
        _acsLockManager->SetAutolockConfig(config.isEnable, config.expiredInterval);
    }
    catch (Exception& e) {
        ncTException te;
        convertException (e, te);
        throw te;
    }

    NC_T_EACP_SERVER_TRACE (_T("this: %p end"), this);
}

void ncTEACPHandler::EACP_GetDevicesByUserId(std::vector<ncTLoginDeviceInfo> & _return, const std::string& userId, const int32_t start, const int32_t limit)
{
    NC_T_EACP_SERVER_TRACE (_T("this: %p begin"), this);

    try {
        // 检查索引
        if (start < 0)
            return;
        if (limit == 0)
            return;

        if (limit < -1)
            THROW_E (T_EACP_SERVER_ERR_PROVIDER_NAME, ncTEACPError::NCT_INVALID_LIMIT_VALUE,
                     LOAD_STRING (_T("IDS_INVALID_LIMIT_VALUE")).getCStr ());

        vector<ncDeviceInfo> deviceInfos;
        _acsDeviceManager->GetDevicesByUserIdAndUdid(toCFLString(userId), String::EMPTY, deviceInfos, start, limit, false);

        ncTLoginDeviceInfo tmpInfo;
        for(size_t i = 0; i < deviceInfos.size(); ++i) {
            tmpInfo.baseInfo.udid = toSTLString(deviceInfos[i].baseInfo.udid);
            tmpInfo.baseInfo.name = toSTLString(deviceInfos[i].baseInfo.name);
            tmpInfo.baseInfo.osType = static_cast<int>(deviceInfos[i].baseInfo.clientType);
            tmpInfo.baseInfo.deviceType = toSTLString(deviceInfos[i].baseInfo.deviceType);
            tmpInfo.baseInfo.lastLoginIp = toSTLString(deviceInfos[i].baseInfo.lastLoginIp);
            tmpInfo.baseInfo.lastLoginTime = deviceInfos[i].baseInfo.lastLoginTime;
            tmpInfo.eraseFlag = deviceInfos[i].eraseFlag;
            tmpInfo.lastEraseTime = deviceInfos[i].lastEraseTime;
            tmpInfo.disableFlag = deviceInfos[i].disableFlag;
            tmpInfo.bindFlag = deviceInfos[i].bindFlag;
            tmpInfo.loginFlag = deviceInfos[i].loginFlag;

            _return.push_back(tmpInfo);
        }
    }
    catch (Exception& e) {
        ncTException te;
        convertException (e, te);
        throw te;
    }

    NC_T_EACP_SERVER_TRACE (_T("this: %p end"), this);
}

void ncTEACPHandler::EACP_AddDevice(const std::string& userId, const std::string& udid, const int32_t osType)
{
    NC_T_EACP_SERVER_TRACE (_T("this: %p begin"), this);
    try {
        ncDeviceBaseInfo baseInfo;
        baseInfo.udid = toCFLString(udid);
        baseInfo.clientType = static_cast<ACSClientType>(osType);
        baseInfo.lastLoginTime = -1;

        _acsDeviceManager->AddDevice(toCFLString(userId), baseInfo);
    }
    catch (Exception& e) {
        ncTException te;
        convertException (e, te);
        throw te;
    }

    NC_T_EACP_SERVER_TRACE (_T("this: %p end"), this);
}

void ncTEACPHandler::EACP_AddDevices(std::vector<string> & _return, const std::string& userId, const vector<std::string>& udids, const int32_t osType)
{
    NC_T_EACP_SERVER_TRACE (_T("this: %p begin"), this);
    try {
        vector<ncDeviceBaseInfo> baseInfos;
        ncDeviceBaseInfo baseInfo;
        String checkedAddr;
        set<String> udidsToBeAdded;
        for (auto iter = udids.begin (); iter != udids.end (); iter++) {
            String udid = toCFLString(*iter);
            if (udid.isEmpty ()) { // 过滤空字符串
                continue;
            }

            try {
                checkedAddr = checkMacAddr (udid);
                checkedAddr.toUpper();
                if (udidsToBeAdded.find (checkedAddr) != udidsToBeAdded.end ()) { // 去重
                    continue;
                }
                baseInfo.udid = checkedAddr;
                baseInfo.clientType = static_cast<ACSClientType>(osType);
                baseInfo.lastLoginTime = -1;
                baseInfos.push_back (baseInfo);
                udidsToBeAdded.insert (checkedAddr);
            }
            catch (Exception& e) {
                _return.push_back (*iter);
            }
        }
        _acsDeviceManager->AddDevices (toCFLString(userId), baseInfos);
    }
    catch (Exception& e) {
        ncTException te;
        convertException (e, te);
        throw te;
    }
    NC_T_EACP_SERVER_TRACE (_T("this: %p end"), this);
}

void ncTEACPHandler::EACP_SearchDevicesByUserIdAndUdid(std::vector<ncTLoginDeviceInfo>& _return, const std::string& userId, const std::string& udid, const int32_t start, const int32_t limit)
{
    NC_T_EACP_SERVER_TRACE (_T("this: %p begin"), this);

    try {
        // 检查索引
        if (start < 0)
            return;
        if (limit == 0)
            return;

        if (limit < -1)
            THROW_E (T_EACP_SERVER_ERR_PROVIDER_NAME, ncTEACPError::NCT_INVALID_LIMIT_VALUE,
                     LOAD_STRING (_T("IDS_INVALID_LIMIT_VALUE")).getCStr ());

        vector<ncDeviceInfo> deviceInfos;
        _acsDeviceManager->GetDevicesByUserIdAndUdid(toCFLString(userId), toCFLString(udid), deviceInfos, start, limit, false);

        ncTLoginDeviceInfo tmpInfo;
        for(size_t i = 0; i < deviceInfos.size(); ++i) {
            tmpInfo.baseInfo.udid = toSTLString(deviceInfos[i].baseInfo.udid);
            tmpInfo.baseInfo.name = toSTLString(deviceInfos[i].baseInfo.name);
            tmpInfo.baseInfo.osType = static_cast<int>(deviceInfos[i].baseInfo.clientType);
            tmpInfo.baseInfo.deviceType = toSTLString(deviceInfos[i].baseInfo.deviceType);
            tmpInfo.baseInfo.lastLoginIp = toSTLString(deviceInfos[i].baseInfo.lastLoginIp);
            tmpInfo.baseInfo.lastLoginTime = deviceInfos[i].baseInfo.lastLoginTime;
            tmpInfo.eraseFlag = deviceInfos[i].eraseFlag;
            tmpInfo.lastEraseTime = deviceInfos[i].lastEraseTime;
            tmpInfo.disableFlag = deviceInfos[i].disableFlag;
            tmpInfo.bindFlag = deviceInfos[i].bindFlag;
            tmpInfo.loginFlag = deviceInfos[i].loginFlag;

            _return.push_back(tmpInfo);
        }
    }
    catch (Exception& e) {
        ncTException te;
        convertException (e, te);
        throw te;
    }

    NC_T_EACP_SERVER_TRACE (_T("this: %p end"), this);
}

void ncTEACPHandler::EACP_SearchDevices(std::vector<string>& _return, const std::string& key, const int32_t start, const int32_t limit)
{
    NC_T_EACP_SERVER_TRACE (_T("this: %p begin"), this);

    try {
        if (key.empty()) {
            return;
        }

         // 检查索引
        if (start < 0)
            return;
        if (limit == 0)
            return;

        if (limit < -1)
            THROW_E (T_EACP_SERVER_ERR_PROVIDER_NAME, ncTEACPError::NCT_INVALID_LIMIT_VALUE,
                     LOAD_STRING (_T("IDS_INVALID_LIMIT_VALUE")).getCStr ());

        vector<ncDeviceInfo> deviceInfos;
        _acsDeviceManager->GetDevicesByUserIdAndUdid(String::EMPTY, toCFLString(key), deviceInfos, start, limit, false);

        for(size_t i = 0; i < deviceInfos.size(); ++i) {
            string udid = toSTLString(deviceInfos[i].baseInfo.udid);
            if (find(_return.begin(), _return.end(), udid) == _return.end()) {
                _return.push_back(udid);
            }
        }
    }
    catch (Exception& e) {
        ncTException te;
        convertException (e, te);
        throw te;
    }

    NC_T_EACP_SERVER_TRACE (_T("this: %p end"), this);
}

void ncTEACPHandler::EACP_SearchUserByDeviceUdid(std::vector<string>& _return, const std::string& udid, const int32_t start, const int32_t limit)
{
    NC_T_EACP_SERVER_TRACE (_T("this: %p begin"), this);

    try {
        // 检查索引
        if (start < 0)
            return;
        if (limit == 0)
            return;

        if (limit < -1)
            THROW_E (T_EACP_SERVER_ERR_PROVIDER_NAME, ncTEACPError::NCT_INVALID_LIMIT_VALUE,
                     LOAD_STRING (_T("IDS_INVALID_LIMIT_VALUE")).getCStr ());

        String UDID = removeBlankAndDot (toCFLString(udid));
        if (UDID.isEmpty ()) {
            THROW_E (T_EACP_SERVER_ERR_PROVIDER_NAME, ncTEACPError::NCT_MAC_ADDR_INVALID,
                     LOAD_STRING("IDS_INVALID_MAC_ADDR").getCStr (), UDID.getCStr ());
        }
        vector<String> users;
        _acsDeviceManager->GetUserNameByUdid(UDID, users, start, limit);

        for (auto user = users.begin(); user != users.end(); user++) {
            _return.push_back(user->getCStr());
        }
    }
    catch (Exception& e) {
        ncTException te;
        convertException (e, te);
        throw te;
    }

    NC_T_EACP_SERVER_TRACE (_T("this: %p end"), this);
}

void ncTEACPHandler::EACP_DeleteDevices(const std::string& userId, const vector<std::string>& udids)
{
    NC_T_EACP_SERVER_TRACE (_T("this: %p begin"), this);
    try {
        set<String> udidsToBeDeleted;
        vector<String> inUdids;
        String checkedAddr;
        for (auto iter = udids.begin (); iter != udids.end (); iter++) {
            String udid = toCFLString(*iter);
            if (udid.isEmpty ()) { // 过滤空字符串
                continue;
            }

            checkedAddr = removeBlankAndDot (udid);
            if (checkedAddr.isEmpty ()) {
                THROW_E (T_EACP_SERVER_ERR_PROVIDER_NAME, ncTEACPError::NCT_MAC_ADDR_INVALID,
                         LOAD_STRING("IDS_INVALID_MAC_ADDR").getCStr (), checkedAddr.getCStr ());
            }
            if (udidsToBeDeleted.find (checkedAddr) != udidsToBeDeleted.end ()) { // 去重
                continue;
            }
            inUdids.push_back(checkedAddr);
            udidsToBeDeleted.insert (checkedAddr);
        }
        _acsDeviceManager->DeleteDevices(toCFLString(userId), inUdids);
    }
    catch (Exception& e) {
        ncTException te;
        convertException (e, te);
        throw te;
    }

    NC_T_EACP_SERVER_TRACE (_T("this: %p end"), this);
}

void ncTEACPHandler::EACP_BindDevice(const std::string& userId, const std::string& udid)
{
    NC_T_EACP_SERVER_TRACE (_T("this: %p begin"), this);
    try {
        vector<String> udids;
        udids.push_back(toCFLString(udid));
        _acsDeviceManager->SetBindStatus(toCFLString(userId), udids);
    }
    catch (Exception& e) {
        ncTException te;
        convertException (e, te);
        throw te;
    }

    NC_T_EACP_SERVER_TRACE (_T("this: %p end"), this);
}

void ncTEACPHandler::EACP_UnbindDevice(const std::string& userId, const std::string& udid)
{
    NC_T_EACP_SERVER_TRACE (_T("this: %p begin"), this);
    try {
        _acsDeviceManager->SetUnbindStatus(toCFLString(userId), toCFLString(udid));
    }
    catch (Exception& e) {
        ncTException te;
        convertException (e, te);
        throw te;
    }

    NC_T_EACP_SERVER_TRACE (_T("this: %p end"), this);
}

void ncTEACPHandler::EACP_DisableDevice(const std::string& userId, const std::string& udid)
{
    NC_T_EACP_SERVER_TRACE (_T("this: %p begin"), this);
    try {
        _acsDeviceManager->SetDisableStatus(toCFLString(userId), toCFLString(udid));
    }
    catch (Exception& e) {
        ncTException te;
        convertException (e, te);
        throw te;
    }

    NC_T_EACP_SERVER_TRACE (_T("this: %p end"), this);
}

void ncTEACPHandler::EACP_EnableDevice(const std::string& userId, const std::string& udid)
{
    NC_T_EACP_SERVER_TRACE (_T("this: %p begin"), this);
    try {
        _acsDeviceManager->SetEnableStatus(toCFLString(userId), toCFLString(udid));
    }
    catch (Exception& e) {
        ncTException te;
        convertException (e, te);
        throw te;
    }

    NC_T_EACP_SERVER_TRACE (_T("this: %p end"), this);
}

void ncTEACPHandler::EACP_ClearPermOutOfScope( )
{
    NC_T_EACP_SERVER_TRACE (_T("this: %p begin"), this);
    try {
        // 开启系统权限共享限制开启才可以清除不在范围的权限和联系人
        if (!_acsShareMgnt->GetPermShareLimitStatus())
            return;

        //1. 清除权限边界范围外的其他权限
        vector<ncOwnerPermInfo> ownerPermInfos;
        _acsPermManager->GetAllCustomPermOwnerInfos(ownerPermInfos);

        for (size_t i = 0; i < ownerPermInfos.size(); ++i) {
            ncOwnerPermInfo& tPermInfo = ownerPermInfos[i];
            if (tPermInfo.accessorType == ACS_CONTACTOR || tPermInfo.accessorType == ACS_GROUP){
                continue;
            }

            String& accessorId = tPermInfo.accessorId;

            // 检查文档访问者是否在所有者的权限范围内
            bool bInScope = false;
            for (size_t j = 0; j < tPermInfo.ownerIds.size(); ++j){
                String& ownerId = tPermInfo.ownerIds[j];

                if (tPermInfo.accessorType == ACS_USER) {
                    if (_acsShareMgnt->CheckUsrInPermScope(ownerId, accessorId)) {
                        bInScope = true;
                        break;
                    }
                }
                else {
                    if (_acsShareMgnt->CheckDeptInPermScope(ownerId, accessorId)) {
                        bInScope = true;
                        break;
                    }
                }
            }

            // 访问者不在所有者的权限范围内，则清除此条权限
            if (!bInScope)
                _acsPermManager->DeleteCustomPermByDocUserId(tPermInfo.docId, accessorId);
        }

        // 3. 清除权限范围外的联系人
        map<String, vector<String> > allContacts;
        _acsShareMgnt->GetAllUserContactIds(allContacts);

        map<String, vector<String> >::iterator iter;
        for (iter = allContacts.begin(); iter != allContacts.end(); ++iter) {
            String userId = iter->first;
            vector<String>& contactIds = iter->second;

            // 获取不在用户权限范围内的联系人
            vector<String> outContactIds;
            _acsShareMgnt->GetUsrIdsOutOfPermScope(userId, contactIds, outContactIds);

            // 删除不在用户权限范围的联系人
            _acsShareMgnt->DeleteContactsByPatch(userId, outContactIds);
        }

    }
    catch (Exception& e) {
        ncTException te;
        convertException (e, te);
        throw te;
    }
    NC_T_EACP_SERVER_TRACE (_T("this: %p end"), this);
}

void ncTEACPHandler::EACP_ClearLinkOutOfScope( )
{
    NC_T_EACP_SERVER_TRACE (_T("this: %p begin"), this);
    try {
        // 系统外链共享限制开启的状态下才可以清除用户的发现共享
        if (!_acsShareMgnt->GetLinkShareLimitStatus())
            return;

        // 获取所有拥有外链共享的用户ID
        createOSSClient ();
        String url;
        url.format ("http://%s/api/doc-share/v1/links/owners", _docShareAddr.getCStr ());
        vector<string> headers;
        ncOSSResponse res;
        (*_ossClientPtr)->Get (url.getCStr (), headers, 30L, res);
        if (200 != res.code) {
            THROW_E (T_EACP_SERVER_ERR_PROVIDER_NAME, -1, _T("Request failed (url:%s, retcode: %d, msg: %s)."),
                url.getCStr (), res.code, res.body.c_str ());
        }

        JSON::Value responseJson;
        try {
            JSON::Reader::read (responseJson, res.body.c_str (), res.body.size ());
        }
        catch (Exception& e) {
            THROW_E (T_EACP_SERVER_ERR_PROVIDER_NAME, ncTEACPError::NCT_INVALID_JSON,
                LOAD_STRING (_T("IDS_INVALID_JSON")).getCStr ());
        }

        // 检查用户外链共享是否开启， 未开启，则清除已有的外链共享
        JSON::Array& userIds = responseJson.a ();
        createNsqManager();
        ncNSQEventType topic = ncNSQEventType::NSQ_CORE_USER_ANONYMOUS_OUT_OF_SCOPE;
        for (size_t i = 0; i < userIds.size (); ++i) {
            String userId = userIds[i].s ().c_str ();
            if(!_acsShareMgnt->IsUserLinkEnabled(userId)){
                NSQMsg nsqMsg;
                nsqMsg.userId = userId;
                _nsqManager->PublishNSQMessage (topic, nsqMsg);
            }
        }
    }
    catch (Exception& e) {
        ncTException te;
        convertException (e, te);
        throw te;
    }
    NC_T_EACP_SERVER_TRACE (_T("this: %p end"), this);
}

void ncTEACPHandler::EACP_ClearFindOutOfScope( )
{
    NC_T_EACP_SERVER_TRACE (_T("this: %p begin"), this);
    NC_T_EACP_SERVER_TRACE (_T("this: %p end"), this);
}

void ncTEACPHandler::convertException (const Exception& e, ncTException& te)
{
    te.expType = ncTExpType::NCT_FATAL;
    te.codeLine = e.getCodeLine ();
    te.errID = e.getErrorId ();
    te.fileName = toSTLString (e.getFileName ());
    te.expMsg = toSTLString (e.getMessage ());
    te.errProvider = toSTLString (e.getErrorProviderName ());
}

void ncTEACPHandler::removeDuplicateStrs (vector<String>& strs)
{
    // 先进行排序
    sort (strs.begin(), strs.end());

    // 在删除掉相邻重复的
    vector<String>::iterator pos = unique (strs.begin(), strs.end());

    // 删除掉最后无效的条目
    strs.erase (pos, strs.end());
}

String ncTEACPHandler::removeBlankAndDot (const String& str)
{
    if (str.isEmpty ())
        return str;

    const char* ptr = str.getCStr ();
    const size_t size = str.getLength ();
    size_t preLen = 0;

    while ((*ptr == _T(' ')) && preLen < size) {
        ++preLen, ++ptr;
    }

    size_t postLen = size;
    ptr = str.getCStr () + size - 1;
    while ( (*ptr == _T(' ') || *ptr == _T('.')) && (postLen > preLen)) {
        --postLen, --ptr;
    }

    return String(str.getCStr () + preLen, postLen - preLen);
}

int ncTEACPHandler::getUTF8StringLength(const String& str)
{
    int cur = 0;
    int length = 0;
    while (cur < str.getLength ()) {
        char c = str[cur];
        if (c == 0) {
            break;
        }
        if (c > 0) {
            cur++;
        }
        while (c < 0) {
            cur++;
            c = c << 1;
        }
        length++;
    }

    return length;
}

String ncTEACPHandler::checkMacAddr(const String& macAddr)
{
    String trimStr = removeBlankAndDot (macAddr);

    size_t length = getUTF8StringLength (trimStr);

    if (length != 17)
        THROW_E (T_EACP_SERVER_ERR_PROVIDER_NAME, ncTEACPError::NCT_MAC_ADDR_INVALID,
            LOAD_STRING("IDS_INVALID_MAC_ADDR").getCStr (), trimStr.getCStr ());

    vector<String> sqlitMacAddr;
    trimStr.split (_T("-"), sqlitMacAddr);

    for (auto iter = sqlitMacAddr.begin (); iter != sqlitMacAddr.end (); iter++) {
        if (iter->getLength () != 2)
            THROW_E (T_EACP_SERVER_ERR_PROVIDER_NAME, ncTEACPError::NCT_MAC_ADDR_INVALID,
                LOAD_STRING("IDS_INVALID_MAC_ADDR").getCStr (), trimStr.getCStr ());
        for (int i = 0; i < 2; i++) {
            auto ch = (*iter)[i];
            if ((ch >= 48 && ch <= 57) || (ch >= 65 && ch <=70) || (ch >= 97 && ch <= 102))
                continue;
            THROW_E (T_EACP_SERVER_ERR_PROVIDER_NAME, ncTEACPError::NCT_MAC_ADDR_INVALID,
                LOAD_STRING("IDS_INVALID_MAC_ADDR").getCStr (), trimStr.getCStr ());
        }
    }

    return trimStr;
}

bool ncTEACPHandler::isValidURL (const String& url)
{
    boost::regex re ("(http|https):\\/\\/(\\w+\\.)*(\\w*)\\/{0,1}([\\w\\d]+\\/{0,1})+");
    return boost::regex_match (toSTLString (url), re);
}

int ncTEACPHandler::convertUtf8ToGbk(string& str)
{
    transform(str.begin(), str.end(), str.begin(), ::tolower);
    const char * rsIn = str.c_str();
    int rsInLength = str.length();
    char buff[str.length() + 128];
    memset(buff, 0, sizeof(buff));

    char *rsOut = buff;
    int rsOutLength = str.length() + 128;

    size_t iLeftRoomLen, iLeftInLen, iOutLen;

    iconv_t stCvt;
    stCvt = iconv_open("gbk", "utf-8");
    if (stCvt == 0)
        return -1;

    iLeftInLen = rsInLength;
    iLeftRoomLen = iLeftInLen * 4 + 1;
    iOutLen = iLeftRoomLen;

    char * pszWorkingBuffer = new char[iLeftRoomLen];
    if (pszWorkingBuffer == NULL)
        return -1;
    char * pszOutBuf = pszWorkingBuffer;
    memset(pszWorkingBuffer,0,iLeftRoomLen);
    int iRet;
    char *pInBuf = (char *)rsIn;
    while (iLeftInLen > 0)
    {
        iRet = iconv(stCvt, &pInBuf, &iLeftInLen, &pszWorkingBuffer, &iLeftRoomLen);
        if (iRet == (int)((size_t)-1))
        {
            if (errno == EILSEQ)
            {
                iLeftInLen -= 3;
                pInBuf += 3;
            }
            else
            {
                iconv_close(stCvt);
                delete[] pszOutBuf;
                return -2;
            }
        }
    }
    iconv_close(stCvt);
    pszOutBuf[iOutLen - iLeftRoomLen] = 0;
    rsOutLength = iOutLen - iLeftRoomLen;
    memcpy(rsOut, pszOutBuf, rsOutLength);
    delete[] pszOutBuf;
    str = rsOut;
    return 0;
}

void ncTEACPHandler::EACP_ClearTokenByUserId(const std::string& userId)
{
    NC_T_EACP_SERVER_TRACE (_T("this: %p, begin"), this);

    try {
        _acsTokenManager->DeleteTokenByUserId(toCFLString(userId));
    }
    catch (Exception& e) {
        ncTException te;
        convertException (e, te);
        throw te;
    }

    NC_T_EACP_SERVER_TRACE (_T("this: %p end"), this);
}

void ncTEACPHandler::EACP_ClearAllInvitationInfo ()
{
    NC_T_EACP_SERVER_TRACE (_T("this: %p, begin"), this);

    NC_T_EACP_SERVER_TRACE (_T("this: %p end"), this);
}

void ncTEACPHandler::EACP_ClearUserInvitationInfo (const std::string& userId)
{
    NC_T_EACP_SERVER_TRACE (_T("this: %p, begin"), this);

    NC_T_EACP_SERVER_TRACE (_T("this: %p end"), this);
}

void ncTEACPHandler::EACP_SetPushMessagesConfig (const std::string& appId)
{
    NC_T_EACP_SERVER_TRACE (_T("this: %p, begin"), this);

    _dbConfManager->SetConfig ( toCFLString (PUSH_MESSAGE_TO_THIRD_PARTY), toCFLString (appId));

    NC_T_EACP_SERVER_TRACE (_T("this: %p end"), this);
}

void ncTEACPHandler::EACP_GetPushMessagesConfig (std::string& appId)
{
    NC_T_EACP_SERVER_TRACE (_T("this: %p, begin"), this);

    String result = _dbConfManager->GetConfig (toCFLString (PUSH_MESSAGE_TO_THIRD_PARTY));
    appId = toSTLString (result);

    NC_T_EACP_SERVER_TRACE (_T("this: %p end"), this);
}

bool ncTEACPHandler::EACP_CheckTokenId (const ncTCheckTokenInfo& tokenInfo)
{
    NC_T_EACP_SERVER_TRACE (_T("this: %p begin"), this);
    NC_T_EACP_SERVER_TRACE (_T("this: %p end"), this);
    return false;
    // NC_T_EACP_SERVER_TRACE (_T("this: %p, begin"), this);

    // ncCheckTokenInfo checkTokenInfo;
    // checkTokenInfo.tokenId = toCFLString(tokenInfo.tokenId);
    // checkTokenInfo.ip = toCFLString(tokenInfo.ip);
    // ncIntrospectInfo introspectInfo;

    // NC_T_EACP_SERVER_TRACE (_T("this: %p end"), this);

    // return CheckToken (checkTokenInfo, introspectInfo);
}

String ncTEACPHandler::IsFileToStr(bool isFile)
{
    if(isFile) {
        return LOAD_STRING(_T("IDS_FILE"));
    }
    else {
        return LOAD_STRING(_T("IDS_DIR"));
    }
}

String ncTEACPHandler::getPathName(const String& path)
{
    String name;
    size_t pos = path.findLastOf (_T('/'));
    if (pos != String::NO_POSITION)
        name = path.subString (pos + 1);

    return name;
}

void ncTEACPHandler::EACP_SetSendShareMailStatus (bool status)
{
    NC_T_EACP_SERVER_TRACE (_T("this: %p, begin"), this);

    try {
        _dbConfManager->SetConfig (toCFLString (ENABLE_SEND_SHARE_MAIL), status ? "true" : "false");
    }
    catch (Exception& e) {
        ncTException te;
        convertException (e, te);
        throw te;
    }

    NC_T_EACP_SERVER_TRACE (_T("this: %p end"), this);
}

bool ncTEACPHandler::EACP_GetSendShareMailStatus ()
{
    NC_T_EACP_SERVER_TRACE (_T("this: %p, begin"), this);

    bool status;
    try {
        status = _dbConfManager->GetConfig (toCFLString (ENABLE_SEND_SHARE_MAIL)) == "true";
    }
    catch (Exception& e) {
        ncTException te;
        convertException (e, te);
        throw te;
    }

    NC_T_EACP_SERVER_TRACE (_T("this: %p end"), this);

    return status;
}

void ncTEACPHandler::EACP_SetCustomApplicationConfig (const string& appConfig)
{
    NC_T_EACP_SERVER_TRACE (_T("this: %p, begin"), this);

    try {
        // 保证配置时合法Json
        JSON::Value appConfigJson;
        try {
            JSON::Reader::read (appConfigJson, appConfig.c_str (), appConfig.size ());
        }
        catch (Exception& e) {
            THROW_E (T_EACP_SERVER_ERR_PROVIDER_NAME, ncTEACPError::NCT_INVALID_JSON,
                LOAD_STRING (_T("IDS_INVALID_JSON")).getCStr ());
        }

        _dbConfManager->SetConfig (toCFLString(CUSTOME_APPLICATION_CONFIG), toCFLString (appConfig));
    }
    catch (Exception& e) {
        ncTException te;
        convertException (e, te);
        throw te;
    }

    NC_T_EACP_SERVER_TRACE (_T("this: %p end"), this);
}

void ncTEACPHandler::EACP_GetCustomApplicationConfig (string& appConfig)
{
    NC_T_EACP_SERVER_TRACE (_T("this: %p, begin"), this);

    try {
        appConfig = toSTLString (_dbConfManager->GetConfig (toCFLString(CUSTOME_APPLICATION_CONFIG)));
    }
    catch (Exception& e) {
        ncTException te;
        convertException (e, te);
        throw te;
    }

    NC_T_EACP_SERVER_TRACE (_T("this: %p end"), this);
}

void ncTEACPHandler::EACP_SetMessageNotifyStatus (bool status)
{
    NC_T_EACP_SERVER_TRACE (_T("this: %p, begin"), this);

    try {
        _dbConfManager->SetConfig (toCFLString (ENABLE_MESSAGE_NOTIFY), status ? "true" : "false");
    }
    catch (Exception& e) {
        ncTException te;
        convertException (e, te);
        throw te;
    }

    NC_T_EACP_SERVER_TRACE (_T("this: %p end"), this);
}

bool ncTEACPHandler::EACP_GetMessageNotifyStatus ()
{
    NC_T_EACP_SERVER_TRACE (_T("this: %p, begin"), this);

    bool status;
    try {
        status = _dbConfManager->GetConfig (toCFLString (ENABLE_MESSAGE_NOTIFY)) == "true";
    }
    catch (Exception& e) {
        ncTException te;
        convertException (e, te);
        throw te;
    }

    NC_T_EACP_SERVER_TRACE (_T("this: %p end"), this);

    return status;
}

void ncTEACPHandler::EACP_ThriftServerPing ()
{
    NC_T_EACP_SERVER_TRACE (_T("this: %p, begin"), this);

    NC_T_EACP_SERVER_TRACE (_T("this: %p end"), this);
    return;
}
