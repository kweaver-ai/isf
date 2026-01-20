/***************************************************************************************************
usermanagement.cpp:
    Copyright (c) Eisoo Software Inc. (2020), All rights reserved.

Purpose:
    usermanagement 服务接口调用

Author:
    Young.yu@aishu.cn

Creating Time:
    2020-11-17
***************************************************************************************************/
#include <abprec.h>
#include "userManagement.h"

#include "drivenadapter.h"
#include <dataapi/ncJson.h>
#include <ossclient/public/ncIOSSClient.h>
#include <ossclient/ossclient.h>

#include "serviceAccessConfig.h"

/* Implementation file */
NS_IMPL_QUERY_INTERFACE1 (userManagement, userManagementInterface)

// protected
NS_IMETHODIMP_(nsrefcnt) userManagement::AddRef (void)
{
    return 1;
}

// protected
NS_IMETHODIMP_(nsrefcnt) userManagement::Release (void)
{
    return 1;
}

AB_DEFINE_THREADSAFE_SINGLETON_NO_POOL (userManagement)

userManagement::userManagement (): _ossClientPtr (0)
{
    NC_DRIVEN_ADAPTER_TRACE (_T("this: %p"), this);

    _userRoleTypeMap = {{"super_admin", ncUserRoleType::SUPER_ADMIN}, {"sys_admin", ncUserRoleType::SYS_ADMIN}, {"audit_admin", ncUserRoleType::AUDIT_ADMIN},
                            {"sec_admin", ncUserRoleType::SEC_ADMIN}, {"org_manager", ncUserRoleType::ORG_MANAGER}, {"org_audit", ncUserRoleType::ORG_AUDIT},
                            {"normal_user", ncUserRoleType::NORMAL_USER}};

    // http 服务设置获取
    _getAccessorIDsByDepartIDUrl.format(_T("http://%s:%d/api/user-management/v1/departments"),
        ServiceAccessConfig::getInstance()->userManagePrivateHost.getCStr(), ServiceAccessConfig::getInstance()->userManagePrivatePort);
    _getAccessorIDsByUserIDUrl.format(_T("http://%s:%d/api/user-management/v1/users"),
        ServiceAccessConfig::getInstance()->userManagePrivateHost.getCStr(), ServiceAccessConfig::getInstance()->userManagePrivatePort);
    _getOrgNamesByIDsUrl.format(_T("http://%s:%d/api/user-management/v1/names"),
        ServiceAccessConfig::getInstance()->userManagePrivateHost.getCStr(), ServiceAccessConfig::getInstance()->userManagePrivatePort);
    _deleteDepartUrl.format(_T("http://%s:%d/api/user-management/v1/departments"),
        ServiceAccessConfig::getInstance()->userManagePrivateHost.getCStr(), ServiceAccessConfig::getInstance()->userManagePrivatePort);
}

userManagement::~userManagement (void)
{
    NC_DRIVEN_ADAPTER_TRACE (_T("this: %p"), this);
}

void userManagement::createOSSClient ()
{
    if (nullptr == _ossClientPtr.get ()) {
        nsresult ret;
        nsCOMPtr<ncIOSSClient> ossClient = do_CreateInstance (NC_OSS_CLIENT_CONTRACTID, &ret);
        if (NS_FAILED (ret)){
            THROW_E (ACS_DRIVEN_ADAPTER, FAILED_TO_CREATE_XPCOM_INSTANCE,
                     _T("Failed to create OSSClient: %s(0x%x)"),
                     String::toString ((int64)ret).getCStr (), (int64)ret);
        }
        _ossClientPtr.reset (new nsCOMPtr<ncIOSSClient>(ossClient));
    }
}

/* [notxpcom] void GetOrgNameIDInfo ([const] in OrgIDInfoRef orgIDInfos, in OrgNameIDInfoRef info); */
NS_IMETHODIMP_(void) userManagement::GetOrgNameIDInfo(const ncOrgIDInfo & orgIDInfos, ncOrgNameIDInfo & info)
{
    NC_DRIVEN_ADAPTER_TRACE (_T("[BEGIN] this: %p"), this);

    std::string content;
    JSON::Value requestJson;
    requestJson["method"] = "GET";
    JSON::Array& userIdsParam = requestJson["user_ids"].a ();
    JSON::Array& departmentIdsParam = requestJson["department_ids"].a ();
    JSON::Array& contactorIdsParam = requestJson["contactor_ids"].a ();
    JSON::Array& groupIdsParam = requestJson["group_ids"].a ();

    for (size_t i = 0; i < orgIDInfos.vecUserIDs.size (); ++i) {
        userIdsParam.push_back (orgIDInfos.vecUserIDs[i].getCStr ());
    }

    for (size_t i = 0; i < orgIDInfos.vecDepartIDs.size (); ++i) {
        departmentIdsParam.push_back (orgIDInfos.vecDepartIDs[i].getCStr ());
    }

     for (size_t i = 0; i < orgIDInfos.vecContactorIDs.size (); ++i) {
        contactorIdsParam.push_back (orgIDInfos.vecContactorIDs[i].getCStr ());
    }

    for (size_t i = 0; i < orgIDInfos.vecGroupIDs.size (); ++i) {
        groupIdsParam.push_back (orgIDInfos.vecGroupIDs[i].getCStr ());
    }

    JSON::Writer::write(requestJson.o(), content);
    requestJson.clear();

    createOSSClient ();
    ncOSSResponse res;
    vector<string> inHeaders;

    (*_ossClientPtr)->Post(_getOrgNamesByIDsUrl.getCStr(), content, inHeaders, 30, res);
    if (res.code != 200)
    {
        if (res.code == 0) {
            THROW_E (ACS_DRIVEN_ADAPTER, DRIVEN_HTTP_SERVER_INTERNAL_ERR, _T("Could not connect to server"));
        }
        else{
            JSON::Value errorJson;
            errorJson["code"] = res.code;
            errorJson["body"] = res.body;
            throw errorJson;
        }
    }

    JSON::Value JconsentInfos;
    JSON::Reader::read(JconsentInfos, res.body.c_str(), res.body.length());
    JSON::Array userNames = JconsentInfos["user_names"].a ();
    for (size_t i = 0; i < userNames.size(); ++i)
    {
        String strID = toCFLString(userNames[i]["id"].s().c_str());
        String strName = toCFLString(userNames[i]["name"].s().c_str());
        info.mapUserInfo.insert(pair<String, String>(strID, strName));
    }
    JSON::Array departNames = JconsentInfos["department_names"].a ();
    for (size_t i = 0; i < departNames.size(); ++i)
    {
        String strID = toCFLString(departNames[i]["id"].s().c_str());
        String strName = toCFLString(departNames[i]["name"].s().c_str());
        info.mapDepartInfo.insert(pair<String, String>(strID, strName));
    }
    JSON::Array conatctorNames = JconsentInfos["contactor_names"].a ();
    for (size_t i = 0; i < conatctorNames.size(); ++i)
    {
        String strID = toCFLString(conatctorNames[i]["id"].s().c_str());
        String strName = toCFLString(conatctorNames[i]["name"].s().c_str());
        info.mapContactorInfo.insert(pair<String, String>(strID, strName));
    }
    JSON::Array groupNames = JconsentInfos["group_names"].a ();
    for (size_t i = 0; i < groupNames.size(); ++i)
    {
        String strID = toCFLString(groupNames[i]["id"].s().c_str());
        String strName = toCFLString(groupNames[i]["name"].s().c_str());
        info.mapGroupInfo.insert(pair<String, String>(strID, strName));
    }

    NC_DRIVEN_ADAPTER_TRACE (_T("[END] this: %p"), this);
}


/* [notxpcom] void GetAccessorIDsByUserID ([const] in StringRef userId, in SetStringRef accessorIds); */
NS_IMETHODIMP_(void) userManagement::GetAccessorIDsByUserID(const String & userId, set<String> & accessorIds)
{
    NC_DRIVEN_ADAPTER_TRACE (_T("[BEGIN] this: %p"), this);

    String tmpUrl;
    tmpUrl.format("%s/%s/accessor_ids", _getAccessorIDsByUserIDUrl.getCStr(), userId.getCStr());

    createOSSClient ();
    ncOSSResponse res;
    vector<string> inHeaders;

    (*_ossClientPtr)->Get(tmpUrl.getCStr(), inHeaders, 30, res);
    if (res.code != 200)
    {
        if (res.code == 0) {
            THROW_E (ACS_DRIVEN_ADAPTER, DRIVEN_HTTP_SERVER_INTERNAL_ERR, _T("Could not connect to server"));
        }
        else{
            JSON::Value errorJson;
            errorJson["code"] = res.code;
            errorJson["body"] = res.body;
            throw errorJson;
        }
    }

    JSON::Value JconsentInfos;
    JSON::Reader::read(JconsentInfos, res.body.c_str(), res.body.length());
    for (size_t i = 0; i < JconsentInfos.a().size(); ++i)
    {
        String tmp = toCFLString(JconsentInfos[i].s().c_str());
        accessorIds.insert(tmp);
    }

    NC_DRIVEN_ADAPTER_TRACE (_T("[END] this: %p"), this);
}

/* [notxpcom] void GetAccessorIDsByDepartID ([const] in StringRef depId, in SetStringRef accessorIds); */
NS_IMETHODIMP_(void) userManagement::GetAccessorIDsByDepartID(const String & depId, set<String> & accessorIds)
{
    NC_DRIVEN_ADAPTER_TRACE (_T("[BEGIN] this: %p"), this);

    String tmpUrl;
    tmpUrl.format("%s/%s/accessor_ids", _getAccessorIDsByDepartIDUrl.getCStr(), depId.getCStr());

    createOSSClient ();
    ncOSSResponse res;
    vector<string> inHeaders;

    (*_ossClientPtr)->Get(tmpUrl.getCStr(), inHeaders, 30, res);
    if (res.code != 200)
    {
        if (res.code == 0) {
            THROW_E (ACS_DRIVEN_ADAPTER, DRIVEN_HTTP_SERVER_INTERNAL_ERR, _T("Could not connect to server"));
        }
        else{
            JSON::Value errorJson;
            errorJson["code"] = res.code;
            errorJson["body"] = res.body;
            throw errorJson;
        }
    }

    JSON::Value JconsentInfos;
    JSON::Reader::read(JconsentInfos, res.body.c_str(), res.body.length());
    for (size_t i = 0; i < JconsentInfos.a().size(); ++i)
    {
        String tmp = toCFLString(JconsentInfos[i].s().c_str());
        accessorIds.insert(tmp);
    }

    NC_DRIVEN_ADAPTER_TRACE (_T("[END] this: %p"), this);
}

/* [notxpcom] void GetUserInfo ([const] in StringRef userId, in UserInfoRef userInfo); */
NS_IMETHODIMP_(void) userManagement::GetUserInfo(const String & userId, UserInfo & userInfo)
{
    NC_DRIVEN_ADAPTER_TRACE (_T("[BEGIN] this: %p"), this);

    createOSSClient ();
    String url;
    url.format (_T("http://%s:%d/api/user-management/v1/users/%s/roles,enabled,priority,name,account,email,telephone,third_attr,third_id,parent_deps"),
        ServiceAccessConfig::getInstance()->userManagePrivateHost.getCStr(), ServiceAccessConfig::getInstance()->userManagePrivatePort, userId.getCStr ());
    vector<string> headers;
    ncOSSResponse response;
    (*_ossClientPtr)->Get (url.getCStr (), headers, 30, response);
    // 请求失败时，原样抛出错误信息
    if (response.code != 200){
        if (response.code == 0) {
            THROW_E (ACS_DRIVEN_ADAPTER, DRIVEN_HTTP_SERVER_INTERNAL_ERR, _T("Could not connect to server"));
        }
        else{
            JSON::Value errorJson;
            errorJson["code"] = response.code;
            errorJson["body"] = response.body;
            throw errorJson;
        }
    }
    // 封装返回结果
    JSON::Value responseJson;
    JSON::Reader::read (responseJson, response.body.c_str (), response.body.size ());
    JSON::Array responseArr = responseJson.a ();
    userInfo.enabled = responseArr[0]["enabled"].b();
    userInfo.priority = responseArr[0]["priority"].i();
    userInfo.id = userId;
    userInfo.name = responseArr[0]["name"].s().c_str();
    userInfo.account = responseArr[0]["account"].s().c_str();
    userInfo.email = responseArr[0]["email"].s().c_str();
    userInfo.telephone = responseArr[0]["telephone"].s().c_str();
    userInfo.thirdAttr = responseArr[0]["third_attr"].s().c_str();
    userInfo.thirdId = responseArr[0]["third_id"].s().c_str();
    JSON::Array roles = responseArr[0]["roles"].a();
    for(auto role:roles){
        userInfo.roles.insert(_userRoleTypeMap[role.s()]);
    }
    JSON::Array departIds = responseArr[0]["parent_deps"].a ();
    vector<String> departIdsVector;
    vector<String> departNamesVector;
    String departPath;
    String departNamePath;
    String strID;
    String strName;
    for (size_t i = 0; i < departIds.size(); ++i)
    {
        departPath = "";
        departNamePath = "";
        strID = "";
        strName = "";
        departIdsVector.clear();
        departNamesVector.clear();
        JSON::Array departId = departIds[i].a();
        for(size_t j = 0; j < departId.size(); j++)
        {
            strID = toCFLString(departId[j]["id"].s().c_str());
            strName = toCFLString(departId[j]["name"].s().c_str());
            departIdsVector.push_back(strID);
            departNamesVector.push_back(strName);
        }
        for (auto it = departIdsVector.begin(); it != departIdsVector.end(); it++){
            departPath += *it + "/";
        }
        departPath = departPath.subString(0, departPath.getLength () -1);

        for (auto it = departNamesVector.begin(); it != departNamesVector.end(); it++){
            departNamePath += *it + "/";
        }
        departNamePath = departNamePath.subString(0, departNamePath.getLength () -1);

        userInfo.vecDepartIDs.push_back(departPath);
        userInfo.vecDepartNames.push_back(departNamePath);
    }

    NC_DRIVEN_ADAPTER_TRACE (_T("[END] this: %p"), this);
}

/* [notxpcom] void BatchGetUserInfo([const] in VecStringRef userIds, in VecUserInfoRef userInfos); */
NS_IMETHODIMP_(void) userManagement::BatchGetUserInfo(const vector<String>& userIds, vector<UserInfo>& userInfos)
{
    NC_DRIVEN_ADAPTER_TRACE (_T("[BEGIN] this: %p"), this);

    createOSSClient ();
    userInfos.clear ();
    if (userIds.size () == 0){
        THROW_E (ACS_DRIVEN_ADAPTER, INVALID_PARAMETER_VALUES, _T("the size of userIds can not be zero"));
    }
    // 如果大小为1，调用获得单一用户信息接口
    if (userIds.size () == 1){
        UserInfo userInfo;
        GetUserInfo (userIds[0], userInfo);
        userInfos.push_back (userInfo);
        return;
    }
    // 大小大于1，调用批量获取接口
    String ids;
    for (size_t i = 0; i < userIds.size (); ++i) {
        ids.append (userIds[i].getCStr ());
        if (i != (userIds.size ()-1)){
            ids.append (",", 1);
        }
    }
    String url;
    url.format (_T("http://%s:%d/api/user-management/v1/users/%s/name,account,email,telephone,third_attr,roles,enabled,priority,third_id"),
        ServiceAccessConfig::getInstance()->userManagePrivateHost.getCStr(), ServiceAccessConfig::getInstance()->userManagePrivatePort, ids.getCStr ());
    vector<string> headers;
    ncOSSResponse response;
    (*_ossClientPtr)->Get (url.getCStr (), headers, 30, response);
    // 请求失败时，原样抛出错误信息
    if (response.code != 200){
        if (response.code == 0) {
            THROW_E (ACS_DRIVEN_ADAPTER, DRIVEN_HTTP_SERVER_INTERNAL_ERR, _T("Could not connect to server"));
        }
        else{
            JSON::Value errorJson;
            errorJson["code"] = response.code;
            errorJson["body"] = response.body;
            throw errorJson;
        }
    }
    // 封装返回结果
    JSON::Value responseJson;
    JSON::Reader::read (responseJson, response.body.c_str (), response.body.size ());
    JSON::Array responseArr = responseJson.a ();
    for (size_t i = 0; i < responseArr.size(); i++){
        UserInfo userInfo;
        userInfo.enabled = responseArr[i]["enabled"].b();
        userInfo.priority = responseArr[i]["priority"].i();
        userInfo.id = responseArr[i]["id"].s().c_str();
        userInfo.name = responseArr[i]["name"].s().c_str();
        userInfo.account = responseArr[i]["account"].s().c_str();
        userInfo.email = responseArr[i]["email"].s().c_str();
        userInfo.telephone = responseArr[i]["telephone"].s().c_str();
        userInfo.thirdAttr = responseArr[i]["third_attr"].s().c_str();
        userInfo.thirdId = responseArr[i]["third_id"].s().c_str();
        JSON::Array roles = responseArr[i]["roles"].a();
        for(auto role:roles){
            userInfo.roles.insert(_userRoleTypeMap[role.s()]);
        }
        userInfos.push_back (userInfo);
    }

    NC_DRIVEN_ADAPTER_TRACE (_T("[END] this: %p"), this);
}

/* [notxpcom] void GetAppInfo([const] in StringRef appId, in AppInfoRef appInfo); */
NS_IMETHODIMP_(void) userManagement::GetAppInfo(const String & appId, AppInfo & appInfo)
{
    NC_DRIVEN_ADAPTER_TRACE (_T("[BEGIN] this: %p"), this);

    createOSSClient ();
    String url;
    url.format (_T("http://%s:%d/api/user-management/v1/apps/%s"),
        ServiceAccessConfig::getInstance()->userManagePrivateHost.getCStr (), ServiceAccessConfig::getInstance()->userManagePrivatePort, appId.getCStr ());
    vector<string> headers;
    ncOSSResponse response;
    (*_ossClientPtr)->Get (url.getCStr (), headers, 30, response);
    // 请求失败时，原样抛出错误信息
    if (response.code != 200){
        if (response.code == 0) {
            THROW_E (ACS_DRIVEN_ADAPTER, DRIVEN_HTTP_SERVER_INTERNAL_ERR, _T("Could not connect to server"));
        }
        else{
            JSON::Value errorJson;
            errorJson["code"] = response.code;
            errorJson["body"] = response.body;
            throw errorJson;
        }
    }
    // 封装返回结果
    JSON::Value responseJson;
    JSON::Reader::read (responseJson, response.body.c_str (), response.body.size ());
    appInfo.id = responseJson["id"].s().c_str ();
    appInfo.name = responseJson["name"].s().c_str ();

    NC_DRIVEN_ADAPTER_TRACE (_T("[END] this: %p"), this);
}

/* [notxpcom] void DeleteDepart ([const] in StringRef departID); */
NS_IMETHODIMP_(void) userManagement::DeleteDepart(const String & departID)
{
    NC_DRIVEN_ADAPTER_TRACE (_T("[BEGIN] this: %p"), this);

    createOSSClient ();
    String url;
    url.format (_T("%s/%s"), _deleteDepartUrl.getCStr (), departID.getCStr ());
    vector<string> headers;
    ncOSSResponse response;
    (*_ossClientPtr)->Delete (url.getCStr (), headers, 30, response);
    // 请求失败时，原样抛出错误信息
    if (response.code != 204){
        if (response.code == 0) {
            THROW_E (ACS_DRIVEN_ADAPTER, DRIVEN_HTTP_SERVER_INTERNAL_ERR, _T("Could not connect to server"));
        }
        else{
            JSON::Value errorJson;
            errorJson["code"] = response.code;
            errorJson["body"] = response.body;
            throw errorJson;
        }
    }

    NC_DRIVEN_ADAPTER_TRACE (_T("[END] this: %p"), this);
}
