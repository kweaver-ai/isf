#ifndef __NC_EAC_INNER_HTTP_SERVER_H__
#define __NC_EAC_INNER_HTTP_SERVER_H__

#if PRAGMA_ONCE
#pragma once
#endif

#include "./protocol/eacinnerhttpserver.pb.h"
#include "./auth/ncEACAuthHandler.h"
#include "./auth/ncEACPolicyHandler.h"
#include "./message/ncEACMessageHandler.h"

#define REGISTER_INNER_HTTPSERVICE(callMethod, serviceName)                           \
    virtual void serviceName(google::protobuf::RpcController* cntl_base,        \
                             const InnerHttpRequest* request,                   \
                             InnerHttpResponse* response,                       \
                             google::protobuf::Closure* done) {                 \
        brpc::ClosureGuard done_guard(done);                                    \
        brpc::Controller* cntl = static_cast<brpc::Controller*>(cntl_base);     \
        callMethod(cntl);                                                       \
    }

class ncIACSTokenManager;
class ncIACSShareMgnt;
class ncIACSDeviceManager;
class ncIACSConfManager;
class ncIACSMessageManager;
class ncIACSPermManager;
class policyEngineInterface;

class ncEACAuthHandler;
class ncEACPolicyHandler;

class ncEACInnerHttpServer : public InnerHttpService {
public:
    ncEACInnerHttpServer(int32_t maxConcurrency=0);
    virtual ~ncEACInnerHttpServer();

    // 内容分析及检索依赖接口
    REGISTER_INNER_HTTPSERVICE(OnPermRequest, Permissions)
    REGISTER_INNER_HTTPSERVICE(OnUserCSFLevelRequest, UserCSFLevel)
    REGISTER_INNER_HTTPSERVICE(OnLatestShareTimeRequest, LatestShareTime)

    // 健康检查
    REGISTER_INNER_HTTPSERVICE(OnHealthRequest, Health)

    // 身份认证
    REGISTER_INNER_HTTPSERVICE(_auth1Handler->GetNew, Auth1_GetNew)
    REGISTER_INNER_HTTPSERVICE(_auth1Handler->GetByThirdParty, Auth1_GetByThirdParty)
    REGISTER_INNER_HTTPSERVICE(_auth1Handler->ConsoleLogin, Auth1_ConsoleLogin)

    // 日志
    REGISTER_INNER_HTTPSERVICE(_auth1Handler->LoginLog, Auth1_LoginLog)

    // 用户策略检查
    REGISTER_INNER_HTTPSERVICE(_policyHandler->CheckPolicy, Policy_Check)

    /***
     * 处理健康检查请求
     */
    void OnHealthRequest (brpc::Controller* cntl);
    // 批量获取文件权限
    void OnPermRequest (brpc::Controller* cntl);
    // 获取用户密级
    void OnUserCSFLevelRequest (brpc::Controller* cntl);
    // 根据用户id获取生效的显示权限的最新时间，用来满足涉密需求
    void OnLatestShareTimeRequest (brpc::Controller* cntl);

    void Start ();

private:
    int32_t                             _maxConcurrency;
    brpc::Server*                       _server;
    ThreadMutexLock                     _mutex;
    nsCOMPtr<ncIACSTokenManager>        _acsToken;
    nsCOMPtr<ncIACSShareMgnt>           _acsShareMgnt;
    nsCOMPtr<ncIACSDeviceManager>       _acsDeviceManager;
    nsCOMPtr<ncIACSConfManager>         _acsConfManager;
    nsCOMPtr<ncIACSMessageManager>      _acsMessageManager;
    nsCOMPtr<ncIACSPermManager>         _acsPermManager;
    nsCOMPtr<ncIACSPolicyManager>       _acsPolicyManager;
    nsCOMPtr<policyEngineInterface>     _policyEngine;

    ncEACAuthHandler*                   _auth1Handler;
    ncEACPolicyHandler*                 _policyHandler;
    ncEACMessageHandler*                _messageHandler;
    int32_t                             _privatePort;
    int32_t                             _brpcInnerPort;
    int32_t                             _eacpThriftInnerPort;
};

#endif // End __NC_EAC_INNER_HTTP_SERVER_H__
