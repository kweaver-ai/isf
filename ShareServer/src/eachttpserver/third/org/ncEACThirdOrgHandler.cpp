#include "eachttpserver.h"
#include "ncEACThirdOrgHandler.h"
#include "ncEACHttpServerUtil.h"
#include "../ncEACThirdUtil.h"
#include "../user/ncEACThirdUserHandler.h"
#include <ehttpserver/ncEHttpUtil.h>
#include <ethriftutil/ncThriftClient.h>
#include "eacServiceAccessConfig.h"

ncEACThirdOrgHandler::ncEACThirdOrgHandler ()
{
    NC_EAC_HTTP_SERVER_TRACE (_T("this: %p"), this);
    _methodFuncs.insert (pair<String, ncMethodFunc>(_T("createorg"), &ncEACThirdOrgHandler::CreateOrg));
    _methodFuncs.insert (pair<String, ncMethodFunc>(_T("editorg"), &ncEACThirdOrgHandler::EditOrg));
    _methodFuncs.insert (pair<String, ncMethodFunc>(_T("deleteorg"), &ncEACThirdOrgHandler::DeleteOrg));
    _methodFuncs.insert (pair<String, ncMethodFunc>(_T("getallorg"), &ncEACThirdOrgHandler::GetAllOrg));
    _methodFuncs.insert (pair<String, ncMethodFunc>(_T("getorgbyid"), &ncEACThirdOrgHandler::GetOrgById));
    _methodFuncs.insert (pair<String, ncMethodFunc>(_T("getorgbyname"), &ncEACThirdOrgHandler::GetOrgByName));
    _methodFuncs.insert (pair<String, ncMethodFunc>(_T("getsubdepsbyorgid"), &ncEACThirdOrgHandler::GetSubDepByOrgId));
    _methodFuncs.insert (pair<String, ncMethodFunc>(_T("getsubusersbyorgid"), &ncEACThirdOrgHandler::GetSubUserByOrgId));
}

ncEACThirdOrgHandler::~ncEACThirdOrgHandler (void)
{
    NC_EAC_HTTP_SERVER_TRACE (_T("this: %p"), this);
}

void
ncEACThirdOrgHandler::CreateOrg (brpc::Controller* cntl, ncIntrospectInfo &info)
{
    // 添加组织信息
    NC_EAC_HTTP_SERVER_TRACE (_T("this: %p, cntl: %p, begin"), this, cntl);

    // 检查是否为应用账户，且有权限
    if (info.visitorType != ncTokenVisitorType::BUSINESS
        || !ncEACThirdUtil::CheckAppOrgPerm (info.userId, ncAppPermOrgType::DEPARTMENT, ncAppOrgPermValue::MODIFY)) {
        THROW_E (EAC_HTTP_SERVER, EACHTTP_USER_AUTH_ERROR,
                    LOAD_STRING (_T("IDS_EACHTTP_ERR_MSG_NOT_AUTHORIZED")));
    }

    string bodyBuffer = cntl->request_attachment ().to_string ();

    // 获取参数
    JSON::Value requestJson;
    try {
        JSON::Reader::read (requestJson, bodyBuffer.c_str (), bodyBuffer.size ());
    }
    catch (Exception& e) {
        THROW_E (EAC_HTTP_SERVER, EACHTTP_INVALID_REQUEST_JSON_FORMAT,
            LOAD_STRING (_T("IDS_EACHTTP_INVALID_JSON")));
    }
    ncTAddOrgParam addOrgInfo;
    addOrgInfo.orgName = requestJson["orgName"].s();

    if (requestJson["priority"].type() != JSON::NIL)
        addOrgInfo.__set_priority(requestJson["priority"].i ());
    if (requestJson["thirdId"].type() != JSON::NIL)
        addOrgInfo.__set_thirdId(requestJson["thirdId"].s ());
    if (requestJson["manager"].type() != JSON::NIL)
    {
        JSON::Object& managerInfo = requestJson["manager"].o ();
        if (managerInfo["type"].s () != "user") 
        {
            THROW_E (EAC_HTTP_SERVER, EACHTTP_INVALID_PARAM_VALUE,
                LOAD_STRING (_T("IDS_EACHTTP_ERR_MSG_AGR_ERR")));
        }
        
        addOrgInfo.__set_managerID(managerInfo["id"].s ());
    }

    //调用sharemgnt服务
    string retOrgId;
    ncTUsrmOrganizationInfo orgInfo;
    try {
        ncThriftClient<ncTShareMgntClient> shareMgntClient (EacServiceAccessConfig::getInstance()->sharemgntHost, EacServiceAccessConfig::getInstance()->sharemgntPort);
        shareMgntClient->Usrm_CreateOrganization (retOrgId, addOrgInfo);

        // 根据id获取组织详细信息
        shareMgntClient->Usrm_GetOrganizationById (orgInfo, retOrgId);
    }
    catch (ncTException & e) {
        THROW_E (EAC_HTTP_SERVER, EACHTTP_CREATER_ORG_ERROR, e.expMsg.c_str ());
    }
    catch (TException & e) {
        printf(e.what ());
        THROW_E (EAC_HTTP_SERVER, EACHTTP_SERVER_UNKNOWN_ERR, e.what ());
    }

    //返回结果
    JSON::Value replyJson;
    replyJson["orgId"] = retOrgId.c_str();
    string body;
    JSON::Writer::write (replyJson.o (), body);
    ncHttpSendReply (cntl, brpc::HTTP_STATUS_OK, "ok", body);

    // 记录审计日志
    try {
        String msg,exmsg;
        // 创建 组织 %s 成功
        msg.format (ncEACHttpServerLoader, _T("IDS_EACHTTP_APP_CREATER_ORG_SUCCESS"), addOrgInfo.orgName.c_str());
        exmsg.format (ncEACHttpServerLoader, _T("IDS_EACHTTP_APP_CREATER_ORG_SUCCESS_EXMSG"), orgInfo.ossInfo.ossName.c_str());

        ncEACHttpServerUtil::Log (cntl, info.userId, info.visitorType, ncTLogType::NCT_LT_MANAGEMENT, ncTLogLevel::NCT_LL_INFO,
                                 ncTManagementType::NCT_MNT_CREATE, msg.getCStr(), exmsg.getCStr());
    }
    catch (...){
    }

    NC_EAC_HTTP_SERVER_TRACE (_T("this: %p, cntl: %p, end"), this, cntl);
}

void
ncEACThirdOrgHandler::EditOrg (brpc::Controller* cntl, ncIntrospectInfo &info)
{
    // 编辑组织信息
    NC_EAC_HTTP_SERVER_TRACE (_T("this: %p, cntl: %p, begin"), this, cntl);

    // 检查是否为应用账户，且有权限
    if (info.visitorType != ncTokenVisitorType::BUSINESS
        || !ncEACThirdUtil::CheckAppOrgPerm (info.userId, ncAppPermOrgType::DEPARTMENT, ncAppOrgPermValue::MODIFY)) {
        THROW_E (EAC_HTTP_SERVER, EACHTTP_USER_AUTH_ERROR,
                    LOAD_STRING (_T("IDS_EACHTTP_ERR_MSG_NOT_AUTHORIZED")));
    }

    string bodyBuffer = cntl->request_attachment ().to_string ();

    // 获取参数
    JSON::Value requestJson;
    try {
        JSON::Reader::read (requestJson, bodyBuffer.c_str (), bodyBuffer.size ());
    }
    catch (Exception& e) {
        THROW_E (EAC_HTTP_SERVER, EACHTTP_INVALID_REQUEST_JSON_FORMAT,
            LOAD_STRING (_T("IDS_EACHTTP_INVALID_JSON")));
    }
    ncTEditDepartParam editOrgInfo;
    editOrgInfo.departId = requestJson["orgId"].s();

    if (requestJson["orgName"].type() != JSON::NIL)
        editOrgInfo.__set_departName(requestJson["orgName"].s ());
    if (requestJson["priority"].type() != JSON::NIL)
        editOrgInfo.__set_priority(requestJson["priority"].i ());
    if (requestJson["manager"].type() != JSON::NIL)
    {
        JSON::Object& managerInfo = requestJson["manager"].o ();
        if (managerInfo["type"].s () != "user") 
        {
            THROW_E (EAC_HTTP_SERVER, EACHTTP_INVALID_PARAM_VALUE,
                LOAD_STRING (_T("IDS_EACHTTP_ERR_MSG_AGR_ERR")));
        }
        
        editOrgInfo.__set_managerID(managerInfo["id"].s ());
    }

    //调用sharemgnt服务
    ncTUsrmOrganizationInfo oldOrgInfo;
    ncTUsrmOrganizationInfo newOrgInfo;
    try {
        ncThriftClient<ncTShareMgntClient> shareMgntClient (EacServiceAccessConfig::getInstance()->sharemgntHost, EacServiceAccessConfig::getInstance()->sharemgntPort);
        // 根据id获取组织详细信息
        shareMgntClient->Usrm_GetOrganizationById (oldOrgInfo, editOrgInfo.departId);

        shareMgntClient->Usrm_EditOrganization (editOrgInfo);

        // 根据id获取组织详细信息
        shareMgntClient->Usrm_GetOrganizationById (newOrgInfo, editOrgInfo.departId);
    }
    catch (ncTException & e) {
        THROW_E (EAC_HTTP_SERVER, EACHTTP_EDIT_ORG_ERROR, e.expMsg.c_str ());
    }
    catch (TException & e) {
        printf(e.what ());
        THROW_E (EAC_HTTP_SERVER, EACHTTP_SERVER_UNKNOWN_ERR, e.what ());
    }

    // 返回结果
    string body = "";
    ncHttpSendReply (cntl, brpc::HTTP_STATUS_OK, "ok", body);

    // 记录审计日志
    try {
        String msg,exmsg;
        // 编辑 组织 %s 成功
        msg.format (ncEACHttpServerLoader, _T("IDS_EACHTTP_APP_EDIT_ORG_SUCCESS"), newOrgInfo.organizationName.c_str());
        exmsg.format (ncEACHttpServerLoader, _T("IDS_EACHTTP_APP_EDIT_ORG_SUCCESS_EXMSG"), oldOrgInfo.organizationName.c_str(), newOrgInfo.ossInfo.ossName.c_str());

        ncEACHttpServerUtil::Log (cntl, info.userId, info.visitorType, ncTLogType::NCT_LT_MANAGEMENT, ncTLogLevel::NCT_LL_INFO,
                                 ncTManagementType::NCT_MNT_SET, msg.getCStr(), exmsg.getCStr());
    }
    catch (...){
    }

    NC_EAC_HTTP_SERVER_TRACE (_T("this: %p, cntl: %p, end"), this, cntl);
}

void
ncEACThirdOrgHandler::DeleteOrg (brpc::Controller* cntl, ncIntrospectInfo &info)
{
    // 删除组织信息
    NC_EAC_HTTP_SERVER_TRACE (_T("this: %p, cntl: %p, begin"), this, cntl);

    // 检查是否为应用账户，且有权限
    if (info.visitorType != ncTokenVisitorType::BUSINESS
        || !ncEACThirdUtil::CheckAppOrgPerm (info.userId, ncAppPermOrgType::DEPARTMENT, ncAppOrgPermValue::MODIFY)) {
        THROW_E (EAC_HTTP_SERVER, EACHTTP_USER_AUTH_ERROR,
                    LOAD_STRING (_T("IDS_EACHTTP_ERR_MSG_NOT_AUTHORIZED")));
    }

    string bodyBuffer = cntl->request_attachment ().to_string ();

    // 获取参数
    JSON::Value requestJson;
    try {
        JSON::Reader::read (requestJson, bodyBuffer.c_str (), bodyBuffer.size ());
    }
    catch (Exception& e) {
        THROW_E (EAC_HTTP_SERVER, EACHTTP_INVALID_REQUEST_JSON_FORMAT,
            LOAD_STRING (_T("IDS_EACHTTP_INVALID_JSON")));
    }

    string orgId = requestJson["orgId"].s();
    ncTUsrmOrganizationInfo orgInfo;
    //调用sharemgnt服务
    try {
        ncThriftClient<ncTShareMgntClient> shareMgntClient (EacServiceAccessConfig::getInstance()->sharemgntHost, EacServiceAccessConfig::getInstance()->sharemgntPort);
        shareMgntClient->Usrm_GetOrganizationById (orgInfo, orgId);

        nsresult ret;
        nsCOMPtr<userManagementInterface> userManager = do_CreateInstance (USER_MANAGEMENT_CONTRACTID, &ret);
        if (NS_FAILED (ret)) {
            THROW_E (EAC_HTTP_SERVER, EACHTTP_SERVER_INIT_USER_MANAGEMENT_ERR,
                _T("Failed to create usermanagement instance: 0x%x"), ret);
        }
        userManager->DeleteDepart (toCFLString(orgId));
    }
    catch (ncTException & e) {
        THROW_E (EAC_HTTP_SERVER, EACHTTP_DELETE_ORG_ERROR, e.expMsg.c_str ());
    }
    catch (TException & e) {
        printf(e.what ());
        THROW_E (EAC_HTTP_SERVER, EACHTTP_SERVER_UNKNOWN_ERR, e.what ());
    }

    // 返回结果
    string body = "";
    ncHttpSendReply (cntl, brpc::HTTP_STATUS_OK, "ok", body);

    // 记录审计日志
    try {
        String msg,exmsg;
        // 删除 组织 %s 成功
        msg.format (ncEACHttpServerLoader, _T("IDS_EACHTTP_APP_DELETE_ORG_SUCCESS"), orgInfo.organizationName.c_str());

        ncEACHttpServerUtil::Log (cntl, info.userId, info.visitorType, ncTLogType::NCT_LT_MANAGEMENT, ncTLogLevel::NCT_LL_INFO,
                                 ncTManagementType::NCT_MNT_DELETE, msg.getCStr(), exmsg.getCStr());
    }
    catch (...){
    }

    NC_EAC_HTTP_SERVER_TRACE (_T("this: %p, cntl: %p, end"), this, cntl);
}

JSON::Value
ncEACThirdOrgHandler::ConvertOrgInfo (ncTRootOrgInfo & orgInfo, bool needThirdId)
{
    //将用户信息转为json格式返回
    JSON::Value replyJson;
    replyJson["orgId"] = orgInfo.id.c_str();
    replyJson["orgName"] = orgInfo.name.c_str();
    replyJson["subDepCnt"] = orgInfo.subDepartmentCount;
    replyJson["subUserCnt"] = orgInfo.subUserCount;
    if (needThirdId) {
        replyJson["thirdId"] = orgInfo.thirdId;
    }

    // 仅用于向上兼容，返回任意一个管理员信息
    if (!orgInfo.responsiblePersons.empty ()) {
        replyJson["orgManagerId"] = orgInfo.responsiblePersons[0].id.c_str();
        replyJson["orgManagerName"] = orgInfo.responsiblePersons[0].user.displayName.c_str();
    }
    else {
        replyJson["orgManagerId"] = "";
        replyJson["orgManagerName"] = "";
    }

    JSON::Array& orgManagers = replyJson["orgManagerInfos"].a ();
    for(int i = 0; i < orgInfo.responsiblePersons.size(); ++i) {
        orgManagers.push_back (JSON::OBJECT);
        JSON::Object& tmpObj = orgManagers.back ().o ();
        tmpObj["orgManagerId"] = orgInfo.responsiblePersons[i].id.c_str();
        tmpObj["orgManagerName"] = orgInfo.responsiblePersons[i].user.displayName.c_str();
    }

    return replyJson;
}

JSON::Value
ncEACThirdOrgHandler::ConvertDepInfo (ncTDepartmentInfo & depInfo)
{
    //将用户信息转为jason格式返回
    JSON::Value replyJson;
    replyJson["depId"] = depInfo.id.c_str();
    replyJson["depName"] = depInfo.name.c_str();
    replyJson["subDepCnt"] = depInfo.subDepartmentCount;
    replyJson["subUserCnt"] = depInfo.subUserCount;

    // 仅用于向上兼容，返回任意一个管理员信息
    if (!depInfo.responsiblePersons.empty ()) {
        replyJson["depManagerId"] = depInfo.responsiblePersons[0].id.c_str();
        replyJson["depManagerName"] = depInfo.responsiblePersons[0].user.displayName.c_str();
    }
    else {
        replyJson["depManagerId"] = "";
        replyJson["depManagerName"] = "";
    }

    JSON::Array& depManagers = replyJson["depManagerInfos"].a ();
    for(int i = 0; i < depInfo.responsiblePersons.size(); ++i) {
        depManagers.push_back (JSON::OBJECT);
        JSON::Object& tmpObj = depManagers.back ().o ();
        tmpObj["depManagerId"] = depInfo.responsiblePersons[i].id.c_str();
        tmpObj["depManagerName"] = depInfo.responsiblePersons[i].user.displayName.c_str();
    }
    return replyJson;
}

void ncEACThirdOrgHandler::ConvertOrganizationInfo (ncTUsrmOrganizationInfo &orgInfo, ncTRootOrgInfo & _return)
{
    // 将ncTUsrmOrganizationInfo结构体转为ncTRootOrgInfo类型

    _return.id = orgInfo.organizationId;
    _return.name = orgInfo.organizationName;
    _return.responsiblePersons = orgInfo.responsiblePersons;
    _return.ossInfo = orgInfo.ossInfo;
    _return.thirdId = orgInfo.thirdId;
}


void
ncEACThirdOrgHandler::GetAllOrg (brpc::Controller* cntl, ncIntrospectInfo &info)
{
    // 获取所有组织信息
    NC_EAC_HTTP_SERVER_TRACE (_T("this: %p, cntl: %p, begin"), this, cntl);

    // 检查是否为应用账户，且有权限
    if (info.visitorType != ncTokenVisitorType::BUSINESS
        || !ncEACThirdUtil::CheckAppOrgPerm (info.userId, ncAppPermOrgType::DEPARTMENT, ncAppOrgPermValue::READ)) {
        THROW_E (EAC_HTTP_SERVER, EACHTTP_USER_AUTH_ERROR,
                    LOAD_STRING (_T("IDS_EACHTTP_ERR_MSG_NOT_AUTHORIZED")));
    }

    string bodyBuffer = cntl->request_attachment ().to_string ();

    // 调用sharemgnt服务处理请求
    vector<ncTRootOrgInfo> retOrgInfos;
    try {
        ncThriftClient<ncTShareMgntClient> shareMgntClient (EacServiceAccessConfig::getInstance()->sharemgntHost, EacServiceAccessConfig::getInstance()->sharemgntPort);
        shareMgntClient->Usrm_GetSupervisoryRootOrg (retOrgInfos, g_ShareMgnt_constants.NCT_USER_ADMIN);
    }
    catch (ncTException & e) {
        THROW_E (EAC_HTTP_SERVER, EACHTTP_GET_ORG_ERROR, e.expMsg.c_str ());
    }
    catch (TException & e) {
        printf(e.what ());
        THROW_E (EAC_HTTP_SERVER, EACHTTP_SERVER_UNKNOWN_ERR, e.what ());
    }

    //结果处理
    JSON::Value replyJson;
    JSON::Array& orgJson = replyJson["orgInfos"].a ();
    for (int i = 0; i < retOrgInfos.size(); ++i) {
        orgJson.push_back (JSON::OBJECT);
        JSON::Object& tmpObj = orgJson.back ().o ();
        tmpObj = ConvertOrgInfo(retOrgInfos[i], false);
    }
    string body;
    JSON::Writer::write (replyJson.o (), body);
    ncHttpSendReply (cntl, brpc::HTTP_STATUS_OK, "ok", body);

    NC_EAC_HTTP_SERVER_TRACE (_T("this: %p, cntl: %p, end"), this, cntl);
}

void
ncEACThirdOrgHandler::GetOrgById (brpc::Controller* cntl, ncIntrospectInfo &info)
{
    // 通过组织id获取组织信息
    NC_EAC_HTTP_SERVER_TRACE (_T("this: %p, cntl: %p, begin"), this, cntl);

    // 检查是否为应用账户，且有权限
    if (info.visitorType != ncTokenVisitorType::BUSINESS
        || !ncEACThirdUtil::CheckAppOrgPerm (info.userId, ncAppPermOrgType::DEPARTMENT, ncAppOrgPermValue::READ)) {
        THROW_E (EAC_HTTP_SERVER, EACHTTP_USER_AUTH_ERROR,
                    LOAD_STRING (_T("IDS_EACHTTP_ERR_MSG_NOT_AUTHORIZED")));
    }

    string bodyBuffer = cntl->request_attachment ().to_string ();

    // 获取参数
    JSON::Value requestJson;
    try {
        JSON::Reader::read (requestJson, bodyBuffer.c_str (), bodyBuffer.size ());
    }
    catch (Exception& e) {
        THROW_E (EAC_HTTP_SERVER, EACHTTP_INVALID_REQUEST_JSON_FORMAT,
            LOAD_STRING (_T("IDS_EACHTTP_INVALID_JSON")));
    }

    string orgId = requestJson["orgId"].s();

    //调用sharemgnt服务
    ncTUsrmOrganizationInfo orgInfo;
    try {
        ncThriftClient<ncTShareMgntClient> shareMgntClient (EacServiceAccessConfig::getInstance()->sharemgntHost, EacServiceAccessConfig::getInstance()->sharemgntPort);
        shareMgntClient->Usrm_GetOrganizationById (orgInfo, orgId);
    }
    catch (ncTException & e) {
        THROW_E (EAC_HTTP_SERVER, EACHTTP_GET_ORG_ERROR, e.expMsg.c_str ());
    }
    catch (TException & e) {
        printf(e.what ());
        THROW_E (EAC_HTTP_SERVER, EACHTTP_SERVER_UNKNOWN_ERR, e.what ());
    }

    //返回结果
    ncTRootOrgInfo ret;
    ConvertOrganizationInfo(orgInfo, ret);
    JSON::Value replyJson = ConvertOrgInfo(ret, true);
    string body;
    JSON::Writer::write (replyJson.o (), body);
    ncHttpSendReply (cntl, brpc::HTTP_STATUS_OK, "ok", body);

    // 记录审计日志

    NC_EAC_HTTP_SERVER_TRACE (_T("this: %p, cntl: %p, end"), this, cntl);
}


void
ncEACThirdOrgHandler::GetOrgByName (brpc::Controller* cntl, ncIntrospectInfo &info)
{
    // 通过组织名获取组织信息
    NC_EAC_HTTP_SERVER_TRACE (_T("this: %p, cntl: %p, begin"), this, cntl);

    // 检查是否为应用账户，且有权限
    if (info.visitorType != ncTokenVisitorType::BUSINESS
        || !ncEACThirdUtil::CheckAppOrgPerm (info.userId, ncAppPermOrgType::DEPARTMENT, ncAppOrgPermValue::READ)) {
        THROW_E (EAC_HTTP_SERVER, EACHTTP_USER_AUTH_ERROR,
                    LOAD_STRING (_T("IDS_EACHTTP_ERR_MSG_NOT_AUTHORIZED")));
    }

    string bodyBuffer = cntl->request_attachment ().to_string ();

    // 获取参数
    JSON::Value requestJson;
    try {
        JSON::Reader::read (requestJson, bodyBuffer.c_str (), bodyBuffer.size ());
    }
    catch (Exception& e) {
        THROW_E (EAC_HTTP_SERVER, EACHTTP_INVALID_REQUEST_JSON_FORMAT,
            LOAD_STRING (_T("IDS_EACHTTP_INVALID_JSON")));
    }

    string orgName = requestJson["orgName"].s();

    //调用sharemgnt服务
    ncTUsrmOrganizationInfo orgInfo;
    try {
        ncThriftClient<ncTShareMgntClient> shareMgntClient (EacServiceAccessConfig::getInstance()->sharemgntHost, EacServiceAccessConfig::getInstance()->sharemgntPort);
        shareMgntClient->Usrm_GetOrganizationByName (orgInfo, orgName);
    }
    catch (ncTException & e) {
        THROW_E (EAC_HTTP_SERVER, EACHTTP_GET_ORG_ERROR, e.expMsg.c_str ());
    }
    catch (TException & e) {
        printf(e.what ());
        THROW_E (EAC_HTTP_SERVER, EACHTTP_SERVER_UNKNOWN_ERR, e.what ());
    }

    //返回结果
    ncTRootOrgInfo ret;
    ConvertOrganizationInfo(orgInfo, ret);
    JSON::Value replyJson = ConvertOrgInfo(ret, false);
    string body;
    JSON::Writer::write (replyJson.o (), body);
    ncHttpSendReply (cntl, brpc::HTTP_STATUS_OK, "ok", body);

    // 记录审计日志
    NC_EAC_HTTP_SERVER_TRACE (_T("this: %p, cntl: %p, end"), this, cntl);
}

void
ncEACThirdOrgHandler::GetSubDepByOrgId (brpc::Controller* cntl, ncIntrospectInfo &info)
{
    // 通过组织id获取子部门信息
    NC_EAC_HTTP_SERVER_TRACE (_T("this: %p, cntl: %p, begin"), this, cntl);

    // 检查是否为应用账户，且有权限
    if (info.visitorType != ncTokenVisitorType::BUSINESS
        || !ncEACThirdUtil::CheckAppOrgPerm (info.userId, ncAppPermOrgType::DEPARTMENT, ncAppOrgPermValue::READ)) {
        THROW_E (EAC_HTTP_SERVER, EACHTTP_USER_AUTH_ERROR,
                    LOAD_STRING (_T("IDS_EACHTTP_ERR_MSG_NOT_AUTHORIZED")));
    }

    string bodyBuffer = cntl->request_attachment ().to_string ();

    // 获取参数
    JSON::Value requestJson;
    try {
        JSON::Reader::read (requestJson, bodyBuffer.c_str (), bodyBuffer.size ());
    }
    catch (Exception& e) {
        THROW_E (EAC_HTTP_SERVER, EACHTTP_INVALID_REQUEST_JSON_FORMAT,
            LOAD_STRING (_T("IDS_EACHTTP_INVALID_JSON")));
    }

    string orgId = requestJson["orgId"].s();

    //调用sharemgnt服务
    vector<ncTDepartmentInfo> retDepInfos;
    try {
        ncThriftClient<ncTShareMgntClient> shareMgntClient (EacServiceAccessConfig::getInstance()->sharemgntHost, EacServiceAccessConfig::getInstance()->sharemgntPort);
        shareMgntClient->Usrm_GetSubDepartments (retDepInfos, orgId);
    }
    catch (ncTException & e) {
        THROW_E (EAC_HTTP_SERVER, EACHTTP_GET_SUB_DEP_ERROR, e.expMsg.c_str ());
    }
    catch (TException & e) {
        printf(e.what ());
        THROW_E (EAC_HTTP_SERVER, EACHTTP_SERVER_UNKNOWN_ERR, e.what ());
    }

    //结果处理
    JSON::Value replyJson;
    JSON::Array& depJson = replyJson["subDepInfos"].a ();
    for (int i = 0; i < retDepInfos.size(); ++i) {
        depJson.push_back (JSON::OBJECT);
        JSON::Object& tmpObj = depJson.back ().o ();
        tmpObj = ConvertDepInfo(retDepInfos[i]);
    }
    string body;
    JSON::Writer::write (replyJson.o (), body);
    ncHttpSendReply (cntl, brpc::HTTP_STATUS_OK, "ok", body);

    NC_EAC_HTTP_SERVER_TRACE (_T("this: %p, cntl: %p, end"), this, cntl);
}


void
ncEACThirdOrgHandler::GetSubUserByOrgId (brpc::Controller* cntl, ncIntrospectInfo &info)
{
    // 通过组织id获取子用户信息
    NC_EAC_HTTP_SERVER_TRACE (_T("this: %p, cntl: %p, begin"), this, cntl);

    // 检查是否为应用账户，且有权限
    if (info.visitorType != ncTokenVisitorType::BUSINESS
        || !ncEACThirdUtil::CheckAppOrgPerm (info.userId, ncAppPermOrgType::DEPARTMENT, ncAppOrgPermValue::READ)) {
        THROW_E (EAC_HTTP_SERVER, EACHTTP_USER_AUTH_ERROR,
                    LOAD_STRING (_T("IDS_EACHTTP_ERR_MSG_NOT_AUTHORIZED")));
    }

    string bodyBuffer = cntl->request_attachment ().to_string ();

    // 获取参数
    JSON::Value requestJson;
    try {
        JSON::Reader::read (requestJson, bodyBuffer.c_str (), bodyBuffer.size ());
    }
    catch (Exception& e) {
        THROW_E (EAC_HTTP_SERVER, EACHTTP_INVALID_REQUEST_JSON_FORMAT,
            LOAD_STRING (_T("IDS_EACHTTP_INVALID_JSON")));
    }

    string orgId = requestJson["orgId"].s();
    int start = requestJson["start"].i();
    int limit = requestJson["limit"].i();

    //调用sharemgnt服务
    vector<ncTUsrmGetUserInfo> retUserInfos;
    try {
        ncThriftClient<ncTShareMgntClient> shareMgntClient (EacServiceAccessConfig::getInstance()->sharemgntHost, EacServiceAccessConfig::getInstance()->sharemgntPort);
        shareMgntClient->Usrm_GetDepartmentOfUsers (retUserInfos, orgId, start, limit);
    }
    catch (ncTException & e) {
        THROW_E (EAC_HTTP_SERVER, EACHTTP_GET_SUB_USER_ERROR, e.expMsg.c_str ());
    }
    catch (TException & e) {
        printf(e.what ());
        THROW_E (EAC_HTTP_SERVER, EACHTTP_SERVER_UNKNOWN_ERR, e.what ());
    }

    //结果处理
    JSON::Value replyJson;
    JSON::Array& depJson = replyJson["subUserInfos"].a ();
    for (int i = 0; i < retUserInfos.size(); ++i) {
        depJson.push_back (JSON::OBJECT);
        JSON::Object& tmpObj = depJson.back ().o ();
        tmpObj = ncEACThirdUserHandler::ConvertUserInfo(retUserInfos[i], false);
    }
    string body;
    JSON::Writer::write (replyJson.o (), body);
    ncHttpSendReply (cntl, brpc::HTTP_STATUS_OK, "ok", body);

    NC_EAC_HTTP_SERVER_TRACE (_T("this: %p, cntl: %p, end"), this, cntl);
}
