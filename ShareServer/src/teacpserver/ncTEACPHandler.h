#ifndef __NC_T_EACP_HANDLER_H
#define __NC_T_EACP_HANDLER_H

#include <acsdb/public/ncIDBLockManager.h>
#include <acsdb/public/ncIDBConfManager.h>
#include <acssharemgnt/public/ncIACSShareMgnt.h>
#include <acsprocessor/public/ncIACSPermManager.h>
#include <acsprocessor/public/ncIACSTokenManager.h>
#include <acsprocessor/public/ncIACSOwnerManager.h>
#include <acsprocessor/public/ncIACSLockManager.h>
#include <acsprocessor/public/ncIACSLicenseManager.h>
#include <acsprocessor/public/ncIACSDeviceManager.h>
#include <acsprocessor/public/ncIACSMessageManager.h>
#include <ossclient/public/ncIOSSClient.h>
#include <boost/thread/tss.hpp>
#include "drivenadapter/public/nsqInterface.h"

#include "gen-cpp/ncTEACP.h"
#include "gen-cpp/EACP_constants.h"
/*
 * EACP Handler
 */
class ncTEACPHandler : virtual public ncTEACPIf {
public:
    ncTEACPHandler();
    ncTEACPHandler( ncIDBLockManager* dbLockManager,
                    ncIACSPermManager* acsPermManager,
                    ncIACSTokenManager* acsTokenManager,
                    ncIACSOwnerManager* acsOwnerManager,
                    ncIACSShareMgnt* acsShareMgnt);

    virtual ~ncTEACPHandler();

public:
    //用户，群组和部门删除操作
    virtual void EACP_OnDeleteUser(const std::string& userId);
    virtual void EACP_OnDeleteDepartment(const std::vector<std::string> & departmentIds);

    //文件锁操作
    virtual void EACP_GetAutolockConfig(ncTAutolockConfig& _return);
    virtual void EACP_SetAutolockConfig(const ncTAutolockConfig& config);

    //许可证管理
    virtual void EACP_GetLicenseInfo(std::string& _return, const std::string& license);
    virtual int32_t EACP_VerifyActiveCode(const std::string& license, const std::string& machineCode, const std::string& activeCode);

    //设备管理操作
    virtual void EACP_GetDevicesByUserId(std::vector<ncTLoginDeviceInfo> & _return, const std::string& userId, const int32_t start, const int32_t limit);
    virtual void EACP_AddDevice(const std::string& userId, const std::string& udid, const int32_t osType);
    virtual void EACP_AddDevices(std::vector<string>& _return, const std::string& userId, const vector<std::string>& udids, const int32_t osType);
    virtual void EACP_DeleteDevices(const std::string& userId, const vector<std::string>& udids);
    virtual void EACP_BindDevice(const std::string& userId, const std::string& udid);
    virtual void EACP_UnbindDevice(const std::string& userId, const std::string& udid);
    virtual void EACP_EnableDevice(const std::string& userId, const std::string& udid);
    virtual void EACP_DisableDevice(const std::string& userId, const std::string& udid);
    virtual void EACP_SearchDevicesByUserIdAndUdid(std::vector<ncTLoginDeviceInfo>& _return, const std::string& userId, const std::string& udid, const int32_t start, const int32_t limit);
    virtual void EACP_SearchUserByDeviceUdid(std::vector<string>& _return, const std::string& udid, const int32_t start, const int32_t limit);
    virtual void EACP_SearchDevices(std::vector<string>& _return, const std::string& key, const int32_t start, const int32_t limit);

    // 权限共享边界操作
    virtual void EACP_ClearPermOutOfScope();
    virtual void EACP_ClearLinkOutOfScope();
    virtual void EACP_ClearFindOutOfScope();

    // 清空用户的所有token（重置用户密码时）
    virtual void EACP_ClearTokenByUserId(const std::string& userId);

    // 清空所有用户共享邀请链接（关闭共享邀请时调用）
    virtual void EACP_ClearAllInvitationInfo();

    // 清空用户所有共享邀请链接（用户冻结时调用）
    virtual void EACP_ClearUserInvitationInfo(const std::string& userId);

    // 设置消息推送到的第三方ID
    virtual void EACP_SetPushMessagesConfig(const std::string& appId);

    // 获取消息推送到的第三方ID
    virtual void EACP_GetPushMessagesConfig(std::string& appId);

    // 检查用户token
    virtual bool EACP_CheckTokenId(const ncTCheckTokenInfo& tokenInfo);

    // 设置邮件通知分享开关状态
    virtual void EACP_SetSendShareMailStatus(bool status);

    // 获取邮件通知分享开关状态
    virtual bool EACP_GetSendShareMailStatus();

    // 设置消息通知状态
    virtual void EACP_SetMessageNotifyStatus(bool status);

    // 获取消息通知状态
    virtual bool EACP_GetMessageNotifyStatus();

    // 设置|获取定制化的应用配置
    virtual void EACP_SetCustomApplicationConfig (const string& appConfig);
    virtual void EACP_GetCustomApplicationConfig (string& appConfig);

    // eacp Thrift Server health check
    virtual void EACP_ThriftServerPing ();

private:
    bool isValidURL (const String& url);

    void convertException (const Exception& e, ncTException& te);
    void removeDuplicateStrs (vector<String>& strs);
    String removeBlankAndDot (const String& str);
    static int convertUtf8ToGbk(string& str);

    int getUTF8StringLength (const String& str);
    String checkMacAddr(const String& macAddr);

    static String IsFileToStr(bool isFile);

    String getPathName(const String& path);

    void createOSSClient ();
    void createNsqManager ();

private:
    nsCOMPtr<ncIDBLockManager>          _dbLockManager;
    nsCOMPtr<ncIACSPermManager>         _acsPermManager;
    nsCOMPtr<ncIACSTokenManager>        _acsTokenManager;
    nsCOMPtr<ncIACSOwnerManager>        _acsOwnerManager;
    nsCOMPtr<ncIACSShareMgnt>           _acsShareMgnt;
    nsCOMPtr<ncIACSLicenseManager>      _acsLicenseManager;
    nsCOMPtr<ncIACSLockManager>         _acsLockManager;
    nsCOMPtr<ncIACSDeviceManager>       _acsDeviceManager;
    nsCOMPtr<ncIDBConfManager>          _dbConfManager;
    nsCOMPtr<ncIACSMessageManager>      _acsMessageManager;
    boost::thread_specific_ptr<nsCOMPtr<ncIOSSClient>> _ossClientPtr;
    nsCOMPtr<nsqInterface>              _nsqManager;
    ThreadMutexLock                     _nsqManagerLock;

private:
    String                              _docShareAddr;
};

#endif // __NC_T_EACP_HANDLER_H
