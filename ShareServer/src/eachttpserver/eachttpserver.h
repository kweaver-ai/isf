#ifndef __EAC_HTTP_SERVER__
#define __EAC_HTTP_SERVER__

#if AB_PRAGMA_ONCE
#pragma once
#endif

////////////////////////////////////////////////////////////////////////////////////////////////////
//
// 头文件
//
#include <abprec.h>

#ifdef __USING_UMM__

    #include <umm/umm.h>

#endif // __USE_UMM__

#ifdef __USING_NCUTIL__

    #include <ncutil/ncutil.h>
    #include <ncutil/ncMemoryLeakProfiler.h>

#endif // __USING_NCUTIL__

#include <eachttpserver/eachttpservererr.h>

#ifdef __USING_EDATAAPI__

    #include <dataapi/dataapi.h>
    #include <dataapi/ncGNSUtil.h>
    #include <dataapi/ncJson.h>

#endif // __USING_EDATAAPI__

#include <ehttpserver/ehttpserver.h>

#include <time.h>

#include <brpc/server.h>
#include <brpc/http_status_code.h>

#include <drivenadapter/public/userManagementInterface.h>
#include <drivenadapter/public/hydraInterface.h>

#include <acsprocessor/public/ncIACSPolicyManager.h>
////////////////////////////////////////////////////////////////////////////////////////////////////
//
// 公共宏

//
// 列举的reserve的size
//
#define NC_EAC_HTTP_SERVER_LIST_RESERVER_SIZE        (2560)
#define EAC_HTTP_SERVER                            _T("eachttpserver")


////////////////////////////////////////////////////////////////////////////////////////////////////
//
// 资源装载对象
//
extern IResourceLoader*     ncEACHttpServerLoader;

#define LOAD_STRING(strID)                                \
    ncEACHttpServerLoader->loadString (strID).getCStr ()

#define __reg__IPv4 Regex((_T("^(25[0-5]|2[0-4][0-9]|[01]?[0-9][0-9]?)(\\.(25[0-5]|2[0-4][0-9]|[01]?[0-9][0-9]?)){3}$")))
#define __reg__IPv6 Regex((_T("^(([0-9a-fA-F]{1,4}:){7,7}[0-9a-fA-F]{1,4}|([0-9a-fA-F]{1,4}:){1,7}:|([0-9a-fA-F]{1,4}:){1,6}:[0-9a-fA-F]{1,4}|([0-9a-fA-F]{1,4}:){1,5}(:[0-9a-fA-F]{1,4}){1,2}|([0-9a-fA-F]{1,4}:){1,4}(:[0-9a-fA-F]{1,4}){1,3}|([0-9a-fA-F]{1,4}:){1,3}(:[0-9a-fA-F]{1,4}){1,4}|([0-9a-fA-F]{1,4}:){1,2}(:[0-9a-fA-F]{1,4}){1,5}|[0-9a-fA-F]{1,4}:((:[0-9a-fA-F]{1,4}){1,6})|:((:[0-9a-fA-F]{1,4}){1,7}|:)|fe80:(:[0-9a-fA-F]{0,4}){0,4}%[0-9a-zA-Z]{1,}|::(ffff(:0{1,4}){0,1}:){0,1}((25[0-5]|(2[0-4]|1{0,1}[0-9]){0,1}[0-9])\\.){3,3}(25[0-5]|(2[0-4]|1{0,1}[0-9]){0,1}[0-9])|([0-9a-fA-F]{1,4}:){1,4}:((25[0-5]|(2[0-4]|1{0,1}[0-9]){0,1}[0-9])\\.){3,3}(25[0-5]|(2[0-4]|1{0,1}[0-9]){0,1}[0-9]))$")))

////////////////////////////////////////////////////////////////////////////////////////////////////
//
//  参数描述
//
struct JsonValueDesc {
public:
    JsonValueDesc ()
        : type (JSON::NIL)
        , isRequired (false)
        , isExist (false)
    {}

    JsonValueDesc (JSON::Type _type,
                bool _isRequired)
        : type (_type)
        , isRequired (_isRequired)
        , isExist (false)
    {}

    JsonValueDesc (JSON::Type _type,
                bool _isRequired,
                map<String, JsonValueDesc>* _valueDescPtr)
        : type (_type)
        , isRequired (_isRequired)
        , isExist (false)
        , valueDescPtr (_valueDescPtr)
    {}

public:
    JSON::Type                    type;                       // 参数类型
    bool                          isRequired;                 // 是否为必须参数
    bool                          isExist;                    // 是否存在
    map<String, JsonValueDesc>*   valueDescPtr;               // 下一层数据信息(类型为数组时key值为"element"，用来描述数组元素)
};

// check token info
struct ncCheckTokenInfo {
    String                tokenId;           // token id
    String                ip;                // 设备ip
};

// token内省信息
struct ncIntrospectInfo {
    String                userId;            // 用户id
    String                scope;             // 权限范围
    String                clientId;          // 客户端id
    ncTokenVisitorType    visitorType;       // 访问者类型
    set<ncUserRoleType>   roleIds;           // 角色id数组
    ncClientType          clientType;        // 设备类型
};
////////////////////////////////////////////////////////////////////////////////////////////////////
//
// eachttpserver TRACE
//
#include <ncRequestIDManager.h>

#define NC_EAC_HTTP_SERVER_TRACE TRACEPRINTF

#define NC_EAC_HTTP_SERVER_CHECK_ARGUMENT_NULL(arg)                                \
    do {                                                                    \
    if (arg == NULL) {                                                    \
        THROW_E (EAC_HTTP_SERVER, EACHTTP_ARGUMENT_INVALID, _T("Argument (%s) is NULL."), _T(#arg));\
    }                                                                    \
    }                                                                        \
    while (false)

//
// 检查 ID 是否合法 使用宏保留出错的代码行
//
// 检查 UUID 是否合法 (URI中)
#define INVALID_UUID_IN_URI(uuid)                                                               \
    if (uuid.isEmpty ()) {                                                                     \
        THROW_E (EAC_HTTP_SERVER, EACHTTP_URI_NOT_EXIST, _T("uuid is empty"));                 \
    }                                                                                               \
    if (!ncOIDUtil::IsGUID (uuid)){                                                            \
        THROW_E (EAC_HTTP_SERVER, EACHTTP_URI_NOT_EXIST, _T("uuid format error"));             \
    }


// 检查 TOKEN_ID 是否合法
#define INVALID_TOKEN_ID(tokenID)                                                                        \
    if (tokenID.isEmpty ()) {                                                                            \
        THROW_E (EAC_HTTP_SERVER, EACHTTP_TOKENID_FORMAT_ERR,                                            \
                    LOAD_STRING (_T("IDS_EACHTTP_TOKENID_FORMAT_ERR")), tokenID.getCStr ());                \
    }                                                                                                    \
    if (!ncOIDUtil::IsGUID (tokenID)){                                                                    \
        THROW_E (EAC_HTTP_SERVER, EACHTTP_TOKENID_FORMAT_ERR,                                        \
                    LOAD_STRING (_T("IDS_EACHTTP_TOKENID_FORMAT_ERR")));                \
    }

// 检查 USER_ID 是否合法
#define INVALID_USER_ID(userID)                                                                            \
    if (userID.isEmpty ()) {                                                                            \
        THROW_E (EAC_HTTP_SERVER, EACHTTP_USERID_FORMAT_ERR,                                            \
                    LOAD_STRING (_T("IDS_EACHTTP_USERID_FORMAT_ERR")), userID.getCStr ());                    \
    }                                                                                                    \
    if (!ncOIDUtil::IsGUID (userID)){                                                                    \
        THROW_E (EAC_HTTP_SERVER, EACHTTP_USERID_FORMAT_ERR,                                            \
                    LOAD_STRING (_T("IDS_EACHTTP_USERID_FORMAT_ERR")), userID.getCStr ());                \
    }

#define INVALID_GROUP_ID(groupID)                                                                        \
    if (groupID.isEmpty ()) {                                                                            \
        THROW_E (EAC_HTTP_SERVER, EACHTTP_INVALID_REQUEST_JSON_FORMAT,                                \
                    LOAD_STRING (_T("IDS_EACHTTP_INVALID_JSON")), groupID.getCStr ());                    \
    }                                                                                                    \
    if (!ncOIDUtil::IsGUID (groupID)){                                                                    \
        THROW_E (EAC_HTTP_SERVER, EACHTTP_GROUPID_FORMAT_ERR,                                        \
                    LOAD_STRING (_T("IDS_EACHTTP_GROUPID_FORMAT_ERR")), groupID.getCStr ());            \
    }

#define INVALID_DEPARTMENT_ID(depID)                                                                    \
    if (depID.isEmpty ()) {                                                                            \
        THROW_E (EAC_HTTP_SERVER, EACHTTP_INVALID_REQUEST_JSON_FORMAT,                                    \
                    LOAD_STRING (_T("IDS_EACHTTP_INVALID_JSON")), depID.getCStr ());                        \
    }                                                                                                    \
    if (!ncOIDUtil::IsGUID (depID)){                                                                    \
        THROW_E (EAC_HTTP_SERVER, EACHTTP_DEPARTMENTID_FORMAT_ERR,                    \
                    LOAD_STRING (_T("IDS_EACHTTP_DEPARTMENTID_FORMAT_ERR")));            \
    }

// 检查 OBJECT_ID 是否合法
#define INVALID_OBJECT_ID(objectID)                                                                        \
    if (!ncOIDUtil::IsObjectID (objectID)){                                                                \
        THROW_E (EAC_HTTP_SERVER, EACHTTP_OBJECTID_FORMAT_ERR,                        \
                    LOAD_STRING (_T("IDS_EACHTTP_OBJECTID_FORMAT_ERR")), objectID.getCStr ());            \
    }

// 检查 DOC_ID 是否合法
#define INVALID_DOC_ID(docID)                                                                            \
    if (docID.isEmpty ()) {                                                                            \
        THROW_E (EAC_HTTP_SERVER, EACHTTP_INVALID_REQUEST_JSON_FORMAT,                                    \
                    LOAD_STRING (_T("IDS_EACHTTP_INVALID_JSON")), docID.getCStr ());                        \
    }                                                                                                    \
    if (!ncOIDUtil::IsDocID (docID)){                                                                    \
        THROW_E (EAC_HTTP_SERVER, EACHTTP_DOCID_FORMAT_ERR,                            \
                    LOAD_STRING (_T("IDS_EACHTTP_DOCID_FORMAT_ERR")));                \
    }

// 检查 IP 是否合法
#define INVALID_IP_ADDRESS(__IP__)                                                                                                       \
    do {                                                                                                                                 \
        if (__IP__.isEmpty ()) {                                                                                                         \
            THROW_E (EAC_HTTP_SERVER, EACHTTP_INVALID_REQUEST_JSON_FORMAT,                                                               \
                    LOAD_STRING (_T("IDS_EACHTTP_INVALID_JSON")), __IP__.getCStr ());                                                    \
        }                                                                                                                                \
        if (!__reg__IPv4.match (__IP__) && !__reg__IPv6.match(__IP__)) {                                                                 \
            THROW_E (EAC_HTTP_SERVER, EACHTTP_INVALID_PARAMETER,                                                                         \
                     LOAD_STRING (_T("IDS_EACHTTP_INVALID_ARG_VALUE")), __IP__.getCStr ());                                              \
        }                                                                                                                                \
    } while (0)

// 截取并检查TOKEN_ID
inline void GET_VALID_TOKEN_ID (const String& authorization, String& tokenId)
{
    // Authorization: Bearer ACCESS_TOKEN
    Regex __reg__ (_T("^Bearer \\.?"));                                                                   \
    if (!__reg__.match (authorization)) {                                                                 \
            THROW_E (EAC_HTTP_SERVER, EACHTTP_TOKENID_FORMAT_ERR,                                         \
                     LOAD_STRING (_T("IDS_EACHTTP_TOKENID_FORMAT_ERR")));                                 \
    }

    tokenId = authorization.subString(7);
}

////////////////////////////////////////////////////////////////////////////////////////////////////
//
// 异常处理宏
//

#define NC_EAC_HTTP_SERVER_TRY \
    try {

#define NC_EAC_HTTP_SERVER_CATCH(cntl) \
    } \
    catch (EHttpDetailException& e) { \
        ncHttpReplyDetailException (cntl, e); \
    } \
    catch (Exception& e) { \
        ncHttpReplyException (cntl, e); \
    } \
    catch (JSON::Value& e) { \
        ncHttpReplyException (cntl, e); \
    } \
    catch (string& e) { \
        printMessage2 (_T("catch string: %s"), e.c_str ()); \
    } \
    catch (...){ \
        try { \
            THROW_E (EAC_HTTP_SERVER, EACHTTP_SERVER_UNKNOWN_ERR, \
                     LOAD_STRING (_T("IDS_EACHTTP_SERVER_UNKNOWN_ERROR"))); \
        } \
        catch (Exception& expt) { \
            ncHttpReplyException (cntl, expt); \
        } \
    }

////////////////////////////////////////////////////////////////////////////////////////////////////
//
// 全局函数
//
void ncCreateEACHttpServerMoResourceLoader (const AppSettings *appSettings,
                                            const AppContext *appCtx);

String ncHttpGetIP (brpc::Controller* cntl);
void ncHttpGetHeader (brpc::Controller* cntl, const String& key, String& value);
void ncHttpGetParams (brpc::Controller* cntl, String& method, String& token, String& userId);
void ncHttpGetParams (brpc::Controller* cntl, String& method, String& token);
void ncHttpGetParams (brpc::Controller* cntl, String& method);
void ncHttpGetQueryString (brpc::Controller* cntl, const String& key, String& value);
void ncHttpGetToken (brpc::Controller* cntl, String& token);

void ncHttpReplyException (brpc::Controller* cntl, Exception& e);
void ncHttpReplyException (brpc::Controller* cntl, JSON::Value& e);
void ncHttpReplyDetailException (brpc::Controller* cntl, EHttpDetailException& e);
void ncHttpSendReply (brpc::Controller* cntl, int code, const char *reason, const string& body);

void CheckRequestParameters (const String& key, const JSON::Value& jsonV, JsonValueDesc& jsonValueDesc);
bool CheckToken (const ncCheckTokenInfo& checkTokenInfo, ncIntrospectInfo& introspectInfo);

String URLDecode (const String &in);

////////////////////////////////////////////////////////////////////////////////////////////////////
//
// brpc 配置
//
#define BRPC_ANALYSE_DATA_PATH ("/sysvol/brpc/eacp")
// brpc gflags 配置
namespace brpc
{
    // brpc的trace存放位置
    DECLARE_string(rpcz_database_dir);

    // 允许在brpc内置服务页面开启http_verbose
    DECLARE_bool(http_verbose);
    BRPC_VALIDATE_GFLAG(http_verbose, PassValidate);

    // 允许在brpc内置服务页面设置http_verbose打印的body长度上限
    DECLARE_int32(http_verbose_max_body_length);
    BRPC_VALIDATE_GFLAG(http_verbose_max_body_length, PositiveInteger);
}

// brpc 记录日志的最低级别，低于该级别会被过滤，0=INFO 1=NOTICE 2=WARNING 3=ERROR 4=FATAL
namespace logging
{
    DECLARE_int32(minloglevel);
}
#endif // __EAC_HTTP_SERVER__
