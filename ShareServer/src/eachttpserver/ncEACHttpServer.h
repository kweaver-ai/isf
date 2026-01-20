#ifndef __NC_EAC_HTTP_SERVER_H__
#define __NC_EAC_HTTP_SERVER_H__

#if PRAGMA_ONCE
#pragma once
#endif

#include "./protocol/eachttpserver.pb.h"

#define REGISTER_HTTPSERVICE(callMethod, serviceName)                           \
    virtual void serviceName(google::protobuf::RpcController* cntl_base,        \
                             const HttpRequest* request,                        \
                             HttpResponse* response,                            \
                             google::protobuf::Closure* done) {                 \
        brpc::ClosureGuard done_guard(done);                                    \
        brpc::Controller* cntl = static_cast<brpc::Controller*>(cntl_base);     \
        callMethod(cntl);                                                       \
    }

class ncIACSTokenManager;
class ncIACSShareMgnt;
class ncIACSPermManager;
class ncIACSDeviceManager;
class ncIACSConfManager;
class ncIACSMessageManager;
class ncIACSPolicyManager;
class policyEngineInterface;
class ncIACSOutboxManager;

class ncEACAuthHandler;
class ncEACDepartmentHandler;
class ncEACUserHandler;
class ncEACContactorHandler;
class ncEACCAuthHandler;
class ncEACPKIHandler;
class ncEACDeviceHandler;
class ncEACMessageHandler;
class ncEACConfigHandler;
class ncEACThirdHandler;

class ncEACHttpServer : public HttpService {
public:
    ncEACHttpServer(int32_t maxConcurrency=0);
    virtual ~ncEACHttpServer();

    // 健康检查
    REGISTER_HTTPSERVICE(OnHealthRequest, Health)

    // Ping
    REGISTER_HTTPSERVICE(OnPingRequest, Ping)

    // 身份认证 V1
    REGISTER_HTTPSERVICE(OnAuth1Request, Auth1_GetConfig)
    REGISTER_HTTPSERVICE(OnAuth1Request, Auth1_GetByNTLMV1)
    REGISTER_HTTPSERVICE(OnAuth1Request, Auth1_GetByTicket)
    REGISTER_HTTPSERVICE(OnAuth1Request, Auth1_GetByADSession)
    REGISTER_HTTPSERVICE(OnAuth1Request, Auth1_ModifyPassword)
    REGISTER_HTTPSERVICE(OnAuth1Request, Auth1_ValidateSecurityDevice)
    REGISTER_HTTPSERVICE(OnAuth1Request, Auth1_CheckUninstallPwd)
    REGISTER_HTTPSERVICE(OnAuth1Request, Auth1_CheckExitPwd)
    REGISTER_HTTPSERVICE(OnAuth1Request, Auth1_GetVcode)
    REGISTER_HTTPSERVICE(OnAuth1Request, Auth1_SendSms)
    REGISTER_HTTPSERVICE(OnAuth1Request, Auth1_SendVcode)
    REGISTER_HTTPSERVICE(OnAuth1Request, Auth1_SendPwdRetrievalVCode)
    REGISTER_HTTPSERVICE(OnAuth1Request, Auth1_SmsActivate)
    REGISTER_HTTPSERVICE(OnAuth1Request, Auth1_ServerTime)
    REGISTER_HTTPSERVICE(OnAuth1Request, Auth1_SendAuthVcode)

    // 身份认证 V2
    REGISTER_HTTPSERVICE(OnAuth2Request, Auth2_GetConfig)
    REGISTER_HTTPSERVICE(OnAuth2Request, Auth2_Login)
    REGISTER_HTTPSERVICE(OnAuth2Request, Auth2_ModifyPassword)
    REGISTER_HTTPSERVICE(OnAuth2Request, Auth2_ValidateSecurityDevice)
    REGISTER_HTTPSERVICE(OnAuth2Request, Auth2_CheckUninstallPwd)
    REGISTER_HTTPSERVICE(OnAuth2Request, Auth2_CheckExitPwd)

    // 获取配置
    REGISTER_HTTPSERVICE(OnAuth1Request, Auth1_Configs)
    REGISTER_HTTPSERVICE(OnAuth1Request, Auth1_LoginConfigs)

    // 部门管理
    REGISTER_HTTPSERVICE(OnDepartmentRequest, Department_GetBasicInfo)
    REGISTER_HTTPSERVICE(OnDepartmentRequest, Department_GetRoots)
    REGISTER_HTTPSERVICE(OnDepartmentRequest, Department_GetSubDeps)
    REGISTER_HTTPSERVICE(OnDepartmentRequest, Department_GetSubUsers)
    REGISTER_HTTPSERVICE(OnDepartmentRequest, Department_Search)
    REGISTER_HTTPSERVICE(OnDepartmentRequest, Department_SearchCount)

    // 用户管理
    REGISTER_HTTPSERVICE(OnUserRequest, User_Get)
    REGISTER_HTTPSERVICE(OnUserRequest, User_GetBasicInfo)
    REGISTER_HTTPSERVICE(OnUserRequest, User_AgreedToTermsOfUse)
    REGISTER_HTTPSERVICE(OnUserRequest, User_Edit)

    // 联系人管理
    REGISTER_HTTPSERVICE(OnContactorRequest, Contactor_GetContactors)
    REGISTER_HTTPSERVICE(OnContactorRequest, Contactor_Search)
    REGISTER_HTTPSERVICE(OnContactorRequest, Contactor_SearchCount)
    REGISTER_HTTPSERVICE(OnContactorRequest, Contactor_AddGroup)
    REGISTER_HTTPSERVICE(OnContactorRequest, Contactor_EditGroup)
    REGISTER_HTTPSERVICE(OnContactorRequest, Contactor_GetGroup)
    REGISTER_HTTPSERVICE(OnContactorRequest, Contactor_AddPersons)
    REGISTER_HTTPSERVICE(OnContactorRequest, Contactor_SearchPersons)
    REGISTER_HTTPSERVICE(OnContactorRequest, Contactor_DeletePersons)
    REGISTER_HTTPSERVICE(OnContactorRequest, Contactor_GetPersons)

    // CA认证
    REGISTER_HTTPSERVICE(OnCARequest, Ca_Get)

    // PKI认证
    REGISTER_HTTPSERVICE(OnPKIRequest, Pki_onOriginal)
    REGISTER_HTTPSERVICE(OnPKIRequest, Pki_onAuthen)

    // 登录设备管理
    REGISTER_HTTPSERVICE(OnDeviceRequest, Device_onList)
    REGISTER_HTTPSERVICE(OnDeviceRequest, Device_onDisable)
    REGISTER_HTTPSERVICE(OnDeviceRequest, Device_onEnable)
    REGISTER_HTTPSERVICE(OnDeviceRequest, Device_onErase)
    REGISTER_HTTPSERVICE(OnDeviceRequest, Device_onGetStatus)
    REGISTER_HTTPSERVICE(OnDeviceRequest, Device_onEraseSuc)

    // 消息通知
    REGISTER_HTTPSERVICE(OnMessageRequest, Message_Get)
    REGISTER_HTTPSERVICE(OnMessageRequest, Message_Read)
    REGISTER_HTTPSERVICE(OnMessageRequest, Message_Read2)
    REGISTER_HTTPSERVICE(OnMessageRequest, Message_SendMail)

    // 配置管理
    REGISTER_HTTPSERVICE(OnConfigRequest, Config_Get)
    REGISTER_HTTPSERVICE(OnConfigRequest, Config_GetOEMConfigBySection)
    REGISTER_HTTPSERVICE(OnConfigRequest, Config_GetDocWatermarkConfig)
    REGISTER_HTTPSERVICE(OnConfigRequest, Config_GetFileCrawlConfig)
    REGISTER_HTTPSERVICE(OnConfigRequest, Config_SetQuickStartStatus)

    // 第三方接入接口 用户管理
    REGISTER_HTTPSERVICE(OnThirdRequest, Third_CreateUser)
    REGISTER_HTTPSERVICE(OnThirdRequest, Third_EditUser)
    REGISTER_HTTPSERVICE(OnThirdRequest, Third_DeleteUser)
    REGISTER_HTTPSERVICE(OnThirdRequest, Third_GetUserById)
    REGISTER_HTTPSERVICE(OnThirdRequest, Third_GetUserByThirdId)
    REGISTER_HTTPSERVICE(OnThirdRequest, Third_GetUserByName)
    REGISTER_HTTPSERVICE(OnThirdRequest, Third_GetAllUser)
    REGISTER_HTTPSERVICE(OnThirdRequest, Third_GetAllUserCount)

    // 第三方接入接口 组织管理
    REGISTER_HTTPSERVICE(OnThirdRequest, Third_CreateOrg)
    REGISTER_HTTPSERVICE(OnThirdRequest, Third_EditOrg)
    REGISTER_HTTPSERVICE(OnThirdRequest, Third_DeleteOrg)
    REGISTER_HTTPSERVICE(OnThirdRequest, Third_GetAllOrg)
    REGISTER_HTTPSERVICE(OnThirdRequest, Third_GetOrgById)
    REGISTER_HTTPSERVICE(OnThirdRequest, Third_GetOrgByName)
    REGISTER_HTTPSERVICE(OnThirdRequest, Third_GetSubDepByOrgId)
    REGISTER_HTTPSERVICE(OnThirdRequest, Third_GetSubUserByOrgId)

    // 第三方接入接口 部门管理
    REGISTER_HTTPSERVICE(OnThirdRequest, Third_CreateDep)
    REGISTER_HTTPSERVICE(OnThirdRequest, Third_EditDep)
    REGISTER_HTTPSERVICE(OnThirdRequest, Third_DeleteDep)
    REGISTER_HTTPSERVICE(OnThirdRequest, Third_GetDepById)
    REGISTER_HTTPSERVICE(OnThirdRequest, Third_GetDepByThirdId)
    REGISTER_HTTPSERVICE(OnThirdRequest, Third_GetDepByName)
    REGISTER_HTTPSERVICE(OnThirdRequest, Third_MoveDep)
    REGISTER_HTTPSERVICE(OnThirdRequest, Third_AddUsersToDep)
    REGISTER_HTTPSERVICE(OnThirdRequest, Third_MoveUsersToDep)
    REGISTER_HTTPSERVICE(OnThirdRequest, Third_RemoveUsersFromDep)
    REGISTER_HTTPSERVICE(OnThirdRequest, Third_GetSubDepsByDepId)
    REGISTER_HTTPSERVICE(OnThirdRequest, Third_GetSubUsersByDepId)
    REGISTER_HTTPSERVICE(OnThirdRequest, Third_SetManager)
    REGISTER_HTTPSERVICE(OnThirdRequest, Third_CancelManager)


    /***
     * 处理健康检查请求
     */
    void OnHealthRequest (brpc::Controller* cntl);

    /***
     * 处理ping请求
     */
    void OnPingRequest (brpc::Controller* cntl);

    /***
     * 处理auth1请求
     */
    void OnAuth1Request (brpc::Controller* cntl);

    /***
     * 处理auth2请求
     */
    void OnAuth2Request (brpc::Controller* cntl);

    /***
     * 处理entrydoc请求
     */
    void OnEntryDocRequest (brpc::Controller* cntl);

    /***
     * 处理managedoc请求
     */
    void OnManageDocRequest (brpc::Controller* cntl);

    /***
     * 处理quota请求
     */
    void OnQuotaRequest (brpc::Controller* cntl);

    /***
     * 处理department请求
     */
    void OnDepartmentRequest (brpc::Controller* cntl);

    /***
     * 处理user请求
     */
    void OnUserRequest (brpc::Controller* cntl);

    /***
     * 处理contactor请求
     */
    void OnContactorRequest (brpc::Controller* cntl);

    /***
     * 处理CA请求
     */
    void OnCARequest (brpc::Controller* cntl);

    /***
     * 处理PKI请求
     */
    void OnPKIRequest (brpc::Controller* cntl);

    /***
     * 处理device请求
     */
    void OnDeviceRequest (brpc::Controller* cntl);

    /***
     * 处理message请求
     */
    void OnMessageRequest (brpc::Controller* cntl);

    /***
     * 处理配置请求
     */
    void OnConfigRequest (brpc::Controller* cntl);

    /***
     * 第三方用户组织管理
     */
    void OnThirdRequest (brpc::Controller* cntl);

    void Start ();

private:
    int32_t                             _maxConcurrency;
    brpc::Server*                       _server;
    ThreadMutexLock                     _mutex;
    nsCOMPtr<ncIACSTokenManager>        _acsToken;
    nsCOMPtr<ncIACSShareMgnt>           _acsShareMgnt;
    nsCOMPtr<ncIACSPermManager>         _acsPermManager;
    nsCOMPtr<ncIACSDeviceManager>       _acsDeviceManager;
    nsCOMPtr<ncIACSConfManager>         _acsConfManager;
    nsCOMPtr<ncIACSMessageManager>      _acsMessageManager;
    nsCOMPtr<ncIACSPolicyManager>       _acsPolicyManager;
    nsCOMPtr<policyEngineInterface>     _policyEngine;
    nsCOMPtr<ncIACSOutboxManager>       _acsOutboxManager;

    ncEACAuthHandler*                   _authHandler;
    ncEACAuthHandler*                   _auth1Handler;
    ncEACAuthHandler*                   _auth2Handler;
    ncEACDepartmentHandler*             _depHandler;
    ncEACUserHandler*                   _userHandler;
    ncEACContactorHandler*              _contactorHandler;
    ncEACCAuthHandler*                  _CAHandler;
    ncEACPKIHandler*                    _pkiHandler;
    ncEACDeviceHandler*                 _deviceHandler;
    ncEACMessageHandler*                _messageHandler;
    ncEACConfigHandler*                 _configHandler;
    ncEACThirdHandler*                  _thirdHandler;
    int32_t                             _httpPort;
    int32_t                             _brpcServerPort;
    int32_t                             _eacpThriftInnerPort;
};

#endif // End __NC_EAC_HTTP_SERVER_H__
