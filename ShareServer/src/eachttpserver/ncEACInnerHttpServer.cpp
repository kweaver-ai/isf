#include "eachttpserver.h"

#include <acsprocessor/public/ncIACSTokenManager.h>
#include <acssharemgnt/public/ncIACSShareMgnt.h>
#include <acsprocessor/public/ncIACSDeviceManager.h>
#include <acsprocessor/public/ncIACSConfManager.h>
#include <acsprocessor/public/ncIACSMessageManager.h>
#include <acsprocessor/public/ncIACSPolicyManager.h>

#include "ncEACInnerHttpServer.h"
#include "./auth/ncEACAuthHandler.h"

#include <ethriftutil/ncThriftClient.h>
#include "gen-cpp/ncTEACP.h"
#include "gen-cpp/EACP_constants.h"
#include <boost/property_tree/ptree.hpp>
#include <boost/property_tree/ini_parser.hpp>

#define APP_DEFAULT_FILE               _T("/sysvol/conf/service_conf/app_default.conf")
#define EACP_ENV_HOST _T("EACP_ENV_HOST")

ncEACInnerHttpServer::ncEACInnerHttpServer(int32_t maxConcurrency)
    : _server (NULL)
    , _mutex ()
    , _maxConcurrency (maxConcurrency)
    , _acsToken (NULL)
    , _acsShareMgnt (NULL)
    , _acsDeviceManager (NULL)
    , _acsConfManager (NULL)
    , _acsMessageManager (NULL)
    , _acsPermManager (NULL)
    , _acsPolicyManager (NULL)
    , _policyEngine (NULL)
    , _auth1Handler (NULL)
    , _messageHandler (NULL)
    , _privatePort (0)
    , _brpcInnerPort (0)
    , _policyHandler (NULL)
{}

ncEACInnerHttpServer::~ncEACInnerHttpServer()
{
    if (_server) {
        delete _server;
        _server = NULL;
    }

    if (_auth1Handler) {
        delete _auth1Handler;
        _auth1Handler = NULL;
    }


    if (_policyHandler) {
        delete _policyHandler;
        _policyHandler = NULL;
    }

    if (_messageHandler) {
        delete _messageHandler;
        _messageHandler = NULL;
    }
}

void ncEACInnerHttpServer::Start ()
{
    _server = new brpc::Server ();
    boost::property_tree::ptree pt;
    boost::property_tree::ini_parser::read_ini (APP_DEFAULT_FILE, pt);
    _privatePort = pt.get<int> ("ShareServer.private_port");
    _brpcInnerPort = pt.get<int> ("ShareServer.brpcinner_port");
    _eacpThriftInnerPort = pt.get<int> ("ShareServer.eacp_port");
    brpc::ServerOptions options;
    options.internal_port = _brpcInnerPort;
    options.idle_timeout_sec = 60;
    options.max_concurrency = _maxConcurrency;
    options.has_builtin_services = pt.get<bool> ("ShareServer.has_builtin_services", false);
    brpc::fLS::FLAGS_rpcz_database_dir = BRPC_ANALYSE_DATA_PATH;
    brpc::fLI::FLAGS_http_verbose_max_body_length = 5120;
    logging::fLI::FLAGS_minloglevel = 3;

    // brpc内置服务默认关闭
    string ready = "";
    string alive = "";
    if (!options.has_builtin_services) {
        // 只在内置服务关闭时才增加健康检查
        ready = "/health/ready                             =>      Health,";
        alive = "/health/alive                             =>      Health,";
    }

    if (_server->AddService (this, brpc::SERVER_DOESNT_OWN_SERVICE,
        // 健康检查
        ready+
        alive+

        // 内容分析及检索依赖接口
        "/api/eacp/v1/permissions                          =>      Permissions,"
        "/api/eacp/v1/csflevel                             =>      UserCSFLevel,"
        "/api/eacp/v1/latestsharetime                      =>      LatestShareTime,"

        // 身份认证
        "/api/eacp/v1/auth1/getnew                         =>      Auth1_GetNew,"
        "/api/eacp/v1/auth1/consolelogin                   =>      Auth1_ConsoleLogin,"
        "/api/eacp/v1/auth1/getbythirdparty                =>      Auth1_GetByThirdParty,"

        // 登录日志
        "/api/eacp/v1/auth1/login-log                         =>   Auth1_LoginLog,"

        // 用户策略检查
        "/api/eacp/v1/policy/check                         =>      Policy_Check"
        ) != 0) {

        THROW_E (EAC_HTTP_SERVER, EACHTTP_SERVER_UNKNOWN_ERR,
                 _T("Failed to add services."));
    }

    nsresult result;
    AutoLock<ThreadMutexLock> lock (&_mutex);

    if (_acsToken == 0 ) {
        _acsToken = do_CreateInstance (NC_ACS_TOKEN_MANAGER_CONTRACTID, &result);
        if (NS_FAILED (result))
            THROW_E (EAC_HTTP_SERVER, EACHTTP_SERVER_INIT_ACS_TOKEN_MANAGER_ERR,
                     LOAD_STRING (_T("IDS_EACHTTP_ACSTOKEN_INIT_ERROR")), result);

    }

    if (_acsShareMgnt == 0) {
        _acsShareMgnt = do_CreateInstance (NC_ACS_SHARE_MGNT_CONTRACTID, &result);
        if (NS_FAILED (result)) {
            THROW_E(EAC_HTTP_SERVER, EACHTTP_SERVER_INIT_ACS_SHAREMGNT_ERR,
                LOAD_STRING (_T("IDS_EACHTTP_ACS_SHAREMGNT_INIT_ERROR")), result);
        }
    }

    if (_acsDeviceManager == 0) {
        _acsDeviceManager = do_CreateInstance (NC_ACS_DEVICE_MANAGER_CONTRACTID, &result);
        if (NS_FAILED (result)) {
            THROW_E (EAC_HTTP_SERVER, EACHTTP_SERVER_INIT_ACS_DEVICE_MANAGER_ERR,
                _T("Failed to create acs device manager: 0x%x"), result);
        }
    }

    if (_acsConfManager == 0) {
        _acsConfManager = do_CreateInstance (NC_ACS_CONF_MANAGER_CONTRACTID, &result);
        if (NS_FAILED (result)) {
            THROW_E (EAC_HTTP_SERVER, EACHTTP_SERVER_INIT_ACS_CONF_MANAGER_ERR,
                _T("Failed to create acs conf manager: 0x%x"), result);
        }
    }

    if (_acsMessageManager == 0) {
        _acsMessageManager = do_CreateInstance (NC_ACS_MESSAGE_MANAGER_CONTRACTID, &result);
        if (NS_FAILED (result)) {
            THROW_E (EAC_HTTP_SERVER, EACHTTP_SERVER_INIT_ACS_MESSAGE_MANAGER_ERR,
                _T("Failed to create acs message manager: 0x%x"), result);
        }
    }

    if (_acsPermManager == 0) {
        _acsPermManager = do_CreateInstance (NC_ACS_PERM_MANAGER_CONTRACTID, &result);
        if (NS_FAILED (result)) {
            THROW_E(EAC_HTTP_SERVER, EACHTTP_SERVER_INIT_ACS_PERM_MANAGER_ERR,
                _T("Failed to create acs message manager: 0x%x"), result);
        }
    }

    if (_acsPolicyManager == 0) {
        _acsPolicyManager = do_CreateInstance (NC_ACS_POLICY_MANAGER_CONTRACTID, &result);
        if (NS_FAILED (result)) {
            THROW_E (EAC_HTTP_SERVER, EACHTTP_SERVER_INIT_ACS_POLICY_MANAGER_ERR,
                _T("Failed to create acs policy manager: 0x%x"), result);
        }
    }

    if (_policyEngine == 0) {
        _policyEngine = do_CreateInstance (POLICY_ENGINE_CONTRACTID, &result);
        if (NS_FAILED (result)) {
            THROW_E (EAC_HTTP_SERVER, EACHTTP_XPCOM_CREATE_INSTANCE_FAILED,
                _T("Failed to create acs policy engine: 0x%x"), result);
        }
    }

    _auth1Handler = new ncEACAuthHandler(_acsToken, _acsShareMgnt, _acsDeviceManager, _acsConfManager, _acsMessageManager, _acsPolicyManager, _policyEngine);
    _policyHandler = new ncEACPolicyHandler(_acsPolicyManager);
    _messageHandler = new ncEACMessageHandler (_acsMessageManager, _acsPermManager, _acsShareMgnt);

    string envHost = toSTLString (Environment::getEnvVariable (EACP_ENV_HOST));
    if (envHost.find(':') == string::npos) {
        if (_server->Start (_privatePort, &options) != 0) {
            THROW_E (EAC_HTTP_SERVER, EACHTTP_SERVER_UNKNOWN_ERR,
                    _T("Failed to start acs http service."));
        }
    }else {
        if (_server->Start (("[::0]:" + String::toString(_privatePort)).getCStr (), &options) != 0) {
            THROW_E (EAC_HTTP_SERVER, EACHTTP_SERVER_UNKNOWN_ERR,
                    (_T("Failed to start acs http service.")));
        }
    }
}

void
ncEACInnerHttpServer::OnHealthRequest (brpc::Controller* cntl)
{
    NC_EAC_HTTP_SERVER_TRY
        // 检测EACP Thrift Sever是否正常
        try {
            ncThriftClient<ncTEACPClient> eacpClient ("localhost", _eacpThriftInnerPort, 3000, 3000, 3000);
            eacpClient->EACP_ThriftServerPing();
        }
        catch (ncTException & e) {
            SystemLog::getInstance ()->log (__FILE__, __LINE__, ERROR_LOG_TYPE,
                                        _T("Ping thrift server error: %s"), e.expMsg.c_str ());
            THROW_E (EAC_HTTP_SERVER, EACHTTP_SERVER_UNKNOWN_ERR, e.expMsg.c_str ());
        }
        catch (TException & e) {
            SystemLog::getInstance ()->log (__FILE__, __LINE__, ERROR_LOG_TYPE,
                                        _T("Ping thrift server error: %s"), e.what ());
            THROW_E (EAC_HTTP_SERVER, EACHTTP_SERVER_UNKNOWN_ERR, e.what ());
        }
        catch (...) {
            SystemLog::getInstance ()->log (__FILE__, __LINE__, ERROR_LOG_TYPE,
                                        _T("Ping thrift server error: Unkown Excepiton."));
            THROW_E (EAC_HTTP_SERVER, EACHTTP_SERVER_UNKNOWN_ERR, _T("Eacp Thrift API Server Error"));
        }

        cntl->http_response().set_content_type("application/json; charset=UTF-8");
        ncHttpSendReply (cntl, brpc::HTTP_STATUS_OK, "ok", "ok");

    NC_EAC_HTTP_SERVER_CATCH (cntl)
}

void
ncEACInnerHttpServer::OnPermRequest (brpc::Controller* cntl)
{
    NC_EAC_HTTP_SERVER_TRY
        NC_EAC_HTTP_SERVER_TRACE (_T("[ncEACInnerHttpServer::OnPermRequest] this: %p, cntl: %p begin"), this, cntl);

        string bodyBuffer = cntl->request_attachment ().to_string ();
        JSON::Value reqVal;
        JSON::Reader::read (reqVal, bodyBuffer.c_str (), bodyBuffer.size ());

        String userId = reqVal["userid"].s ().c_str ();
        JSON::Array& gnsArr = reqVal["gns"].a ();
        String userIp = ncHttpGetIP (cntl);

        set<String> docIds;
        std::map<String, String> readIdMap;
        for (auto iter : gnsArr) {
            String cid = ncGNSUtil::GetCIDPath (iter.s ().c_str ());
            docIds.insert (cid);
            readIdMap.insert (make_pair (iter.s ().c_str (), cid));
        }
        // 过滤受到网段限制的文档库
        // 这里为实名用户
        bool isAnonymous = false;
        if (!isAnonymous && _acsConfManager->GetNetDocsLimitStatus()) {
            _acsShareMgnt->FilterByNetDocLimit(docIds, userIp);
        }

        JSON::Array resArr;
        for (auto iter : docIds) {
            for (auto iterMap : readIdMap) {
                if (0 != iterMap.second.compare (iter)) continue;
                int perm = 0;
                ncAccessPerm result;
                try {
                    result = _acsPermManager->GetPermission(userId, iterMap.first);
                    perm = result.allowValue & (ACS_AP_PREVIEW | ACS_AP_READ);
                }
                catch (Exception& e) {
                    SystemLog::getInstance ()->log (__FILE__, __LINE__, ERROR_LOG_TYPE,
                                                    _T("GetPermission failed, ERROR: %s"), e.toFullString ().getCStr ());
                    continue;
                }
                catch (...) {
                    SystemLog::getInstance ()->log (__FILE__, __LINE__, ERROR_LOG_TYPE,
                                                    _T("GetPermission failed, ERROR: Unknow"));
                    continue;
                }
                // 返回 预览或者下载权限的gns
                if (perm) {
                    resArr.push_back (JSON_MOVE (iterMap.first.getCStr ()));
                }
            }
        }

        string resBodyStr;
        JSON::Writer::write (resArr, resBodyStr);
        ncHttpSendReply (cntl, brpc::HTTP_STATUS_OK, "ok", resBodyStr);

        NC_EAC_HTTP_SERVER_TRACE (_T("[ncEACInnerHttpServer::OnPermRequest] this: %p, cntl: %p, userId: %s end"), this, cntl, userId.getCStr ());
    NC_EAC_HTTP_SERVER_CATCH (cntl)
}

void
ncEACInnerHttpServer::OnUserCSFLevelRequest (brpc::Controller* cntl)
{
    NC_EAC_HTTP_SERVER_TRY
        NC_EAC_HTTP_SERVER_TRACE (_T("[ncEACInnerHttpServer::OnUserCSFLevelRequest] this: %p, cntl: %p begin"), this, cntl);

        String userId, visitorType;
        ncHttpGetQueryString (cntl, "userid", userId);
        ncHttpGetQueryString (cntl, "visitortype", visitorType);

        ncVisitorType type;
        if (0 == visitorType.compare ("realname")) {
            type = ncVisitorType::REALNAME;
        }
        else if (0 == visitorType.compare ("anonymous")) {
            type = ncVisitorType::ANONYMOUS;
        }
        else if (0 == visitorType.compare ("business")) {
            type = ncVisitorType::BUSINESS;
        }
        else {
            THROW_E (EAC_HTTP_SERVER, EACHTTP_QUERY_METHOD_INVALID,
                LOAD_STRING (_T("IDS_EACHTTP_VISITORTYPE_INVALID")), userId.getCStr (), visitorType.getCStr ());
        }

        int level = _acsShareMgnt->GetUserCSFLevel(userId, type);

        JSON::Value resVal;
        resVal["csflevel"].i () =  level;

        string resBodyStr;
        JSON::Writer::write (resVal.o (), resBodyStr);
        ncHttpSendReply (cntl, brpc::HTTP_STATUS_OK, "ok", resBodyStr);

        NC_EAC_HTTP_SERVER_TRACE (_T("[ncEACInnerHttpServer::OnUserCSFLevelRequest] this: %p, cntl: %p, userId: %s end"), this, cntl, userId.getCStr ());
    NC_EAC_HTTP_SERVER_CATCH (cntl)
}

void
ncEACInnerHttpServer::OnLatestShareTimeRequest (brpc::Controller* cntl)
{
    NC_EAC_HTTP_SERVER_TRY
        NC_EAC_HTTP_SERVER_TRACE (_T("[ncEACInnerHttpServer::OnLatestShareTimeRequest] this: %p, cntl: %p begin"), this, cntl);

        string bodyBuffer = cntl->request_attachment ().to_string ();
        JSON::Value reqVal;
        JSON::Reader::read (reqVal, bodyBuffer.c_str (), bodyBuffer.size ());

        String userId = reqVal["userid"].s ().c_str ();
        String visitorType = reqVal["visitortype"].s ().c_str ();
        JSON::Array& gnsArr = reqVal["gns"].a ();

        ncVisitorType type;
        if (0 == visitorType.compare ("realname")) {
            type = ncVisitorType::REALNAME;
        }
        else if (0 == visitorType.compare ("anonymous")) {
            type = ncVisitorType::ANONYMOUS;
        }
        else if (0 == visitorType.compare ("business")) {
            type = ncVisitorType::BUSINESS;
        }
        else {
            THROW_E (EAC_HTTP_SERVER, EACHTTP_QUERY_METHOD_INVALID,
                LOAD_STRING (_T("IDS_EACHTTP_VISITORTYPE_INVALID")), userId.getCStr (), visitorType.getCStr ());
        }

        JSON::Value resVal;
        for (auto iter = gnsArr.begin (); iter != gnsArr.end (); ++iter) {
            String gnsStr (iter->s ().c_str ());
            int64 shareTime = 0;
            set<String> subShareTimeSet;

            try {
                // 检查是否有共享审核权限
                shareTime = _acsPermManager->GetLatestShareTime(userId, gnsStr, type, subShareTimeSet);
            }
            catch (Exception& e) {
                SystemLog::getInstance ()->log (__FILE__, __LINE__, ERROR_LOG_TYPE,
                                                _T("GetLatestShareTime failed, ERROR: %s"), e.toFullString ().getCStr ());
                continue;
            }
            catch (...) {
                SystemLog::getInstance ()->log (__FILE__, __LINE__, ERROR_LOG_TYPE,
                                                _T("GetLatestShareTime failed, ERROR: Unknow"));
                continue;
            }

            JSON::Object& resShareObj = resVal[gnsStr.getCStr ()].o ();
            resShareObj["share_time"] = shareTime;
            JSON::Array& resShareFilesArr =  resShareObj["share_files"].a ();
            for (auto iter : subShareTimeSet) {
                resShareFilesArr.push_back (iter.getCStr ());
            }
        }

        string resBodyStr;
        JSON::Writer::write (resVal.o (), resBodyStr);
        ncHttpSendReply (cntl, brpc::HTTP_STATUS_OK, "ok", resBodyStr);

        NC_EAC_HTTP_SERVER_TRACE (_T("[ncEACInnerHttpServer::OnLatestShareTimeRequest] this: %p, cntl: %p, userId: %s end"), this, cntl, userId.getCStr ());
    NC_EAC_HTTP_SERVER_CATCH (cntl)
}
