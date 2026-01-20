#include "eachttpserver.h"
#include "ncEACDepartmentHandler.h"
#include <dataapi/ncGNSUtil.h>
#include "ncEACHttpServerUtil.h"

#include <ehttpserver/ncEHttpUtil.h>

// public
ncEACDepartmentHandler::ncEACDepartmentHandler (ncIACSShareMgnt* acsShareMgnt)
    : _acsShareMgnt (acsShareMgnt)
{
    NC_EAC_HTTP_SERVER_TRACE (_T("this: %p"), this);

    if (_acsShareMgnt == NULL) {
        nsresult ret;
        _acsShareMgnt = do_CreateInstance (NC_ACS_SHARE_MGNT_CONTRACTID, &ret);
        if (NS_FAILED (ret)) {
            THROW_E(EAC_HTTP_SERVER, EACHTTP_SERVER_INIT_ACS_SHAREMGNT_ERR,
                LOAD_STRING (_T("IDS_EACHTTP_ACS_SHAREMGNT_INIT_ERROR")), ret);
        }
    }

    _methodFuncs.insert (pair<String, ncMethodFunc>(_T("getbasicinfo"), &ncEACDepartmentHandler::GetBasicInfo));
    _methodFuncs.insert (pair<String, ncMethodFunc>(_T("getroots"), &ncEACDepartmentHandler::GetRoots));
    _methodFuncs.insert (pair<String, ncMethodFunc>(_T("getsubdeps"), &ncEACDepartmentHandler::GetSubDeps));
    _methodFuncs.insert (pair<String, ncMethodFunc>(_T("getsubusers"), &ncEACDepartmentHandler::GetSubUsers));
    _methodFuncs.insert (pair<String, ncMethodFunc>(_T("search"), &ncEACDepartmentHandler::Search));
    _methodFuncs.insert (pair<String, ncMethodFunc>(_T("searchcount"), &ncEACDepartmentHandler::SearchCount));
}

// public
ncEACDepartmentHandler::~ncEACDepartmentHandler (void)
{
    NC_EAC_HTTP_SERVER_TRACE (_T("this: %p"), this);
}

// public
void
ncEACDepartmentHandler::doDepRequestHandler (brpc::Controller* cntl)
{
    NC_EAC_HTTP_SERVER_TRY

        // 获取query string
        String method;
        String tokenId;
        String userId;
        ncHttpGetParams (cntl, method, tokenId);
        // method是否设置
        if (method.isEmpty ()) {
            THROW_E (EAC_HTTP_SERVER, EACHTTP_QUERY_METHOD_INVALID,
                    LOAD_STRING (_T("IDS_EACHTTP_QUERY_METHOD_INVALID")));
        }

        // method是否支持
        map<String, ncMethodFunc>::iterator iter = _methodFuncs.find (method);
        if (iter == _methodFuncs.end ()) {
            THROW_E (EAC_HTTP_SERVER, EACHTTP_QUERY_METHOD_INVALID,
                        LOAD_STRING (_T("IDS_EACHTTP_QUERY_METHOD_NOT_SUPPORT")));
        }

        // token验证
        ncCheckTokenInfo checkTokenInfo;
        checkTokenInfo.tokenId = tokenId;
        checkTokenInfo.ip = ncEACHttpServerUtil::GetForwardedIp(cntl);
        ncIntrospectInfo introspectInfo;
        if (CheckToken (checkTokenInfo, introspectInfo) == false) {
            THROW_E (EAC_HTTP_SERVER, EACHTTP_SERVER_ACCESS_TOKEN_ERR,
                    LOAD_STRING (_T("IDS_EACHTTP_ERR_MSG_ACCESS_TOKEN_ERR")));
        }

        // 获取该token对应的userId
        userId = introspectInfo.userId;

        // 消息处理
        ncMethodFunc func = iter->second;
        (this->*func) (cntl, userId);

    NC_EAC_HTTP_SERVER_CATCH (cntl)
}

// protected
void
ncEACDepartmentHandler::GetBasicInfo (brpc::Controller* cntl, const String& userId)
{
    NC_EAC_HTTP_SERVER_TRY

        NC_EAC_HTTP_SERVER_TRACE (_T("this: %p, cntl: %p, userId: %s begin"), this, cntl, userId.getCStr ());

        // 获取http请求的content
        string bodyBuffer = cntl->request_attachment ().to_string ();

        // 获取depid
        JSON::Value requestJson;
        try {
            JSON::Reader::read (requestJson, bodyBuffer.c_str (), bodyBuffer.size ());
        }
        catch (Exception& e) {
            THROW_E (EAC_HTTP_SERVER, EACHTTP_INVALID_REQUEST_JSON_FORMAT,
                LOAD_STRING (_T("IDS_EACHTTP_INVALID_JSON")));
        }

        String depId = requestJson["depid"].s ().c_str ();
        INVALID_DEPARTMENT_ID(depId);

        // 检查部门id是否存在
        ncACSDepartInfo departInfo;
        bool ret = _acsShareMgnt->GetDepartInfoById (depId, departInfo);
        if (ret == false) {
            THROW_E (EAC_HTTP_SERVER, EACHTTP_DEPART_ID_NOT_EXISTS,
                LOAD_STRING (_T("IDS_DEPART_ID_NOT_EXISTS")));
        }

        // 回复
        JSON::Value replyJson;
        replyJson["name"] = departInfo.name.getCStr ();

        // 检查用户是否需要屏蔽组织架构信息
        bool hideOu = false;
#ifndef __UT__
        hideOu = ncEACHttpServerUtil::HideOum_Check(userId.getCStr ());
#endif
        String depParentPath;
        if (!hideOu && !departInfo.id.isEmpty ()) {
            _acsShareMgnt->GetParentDeptRootPathName(departInfo.id, depParentPath);
        }
        replyJson["path"] = depParentPath.getCStr ();


        string body;
        JSON::Writer::write (replyJson.o (), body);
        ncHttpSendReply (cntl, brpc::HTTP_STATUS_OK, "ok", body);

        NC_EAC_HTTP_SERVER_TRACE (_T("this: %p, cntl: %p, userId: %s end"), this, cntl, userId.getCStr ());

    NC_EAC_HTTP_SERVER_CATCH (cntl)
}

// protected
void
ncEACDepartmentHandler::GetRoots (brpc::Controller* cntl, const String& userId)
{
    NC_EAC_HTTP_SERVER_TRY

        NC_EAC_HTTP_SERVER_TRACE (_T("this: %p, cntl: %p, userId: %s begin"), this, cntl, userId.getCStr ());

        // 回复
        JSON::Value replyJson;
        JSON::Array& depJson = replyJson["depinfos"].a ();

        // 检查用户是否需要屏蔽组织架构信息
        bool hideOu = false;
#ifndef __UT__
        hideOu = ncEACHttpServerUtil::HideOum_Check(userId.getCStr ());
#endif
        if (hideOu == false) {
            // 获取根组织信息
            if (_acsShareMgnt->GetPermShareLimitStatus()) {
                vector<pair<ncOrganizationInfo,bool>> organizeInfos;
                _acsShareMgnt->GetScopeOrgInfo (userId, organizeInfos);

                for (size_t i = 0; i < organizeInfos.size (); ++i) {
                    depJson.push_back (JSON::OBJECT);

                    JSON::Object& tmpObj = depJson.back ().o ();
                    tmpObj["depid"] = organizeInfos[i].first.id.getCStr ();
                    tmpObj["name"] = organizeInfos[i].first.name.getCStr ();
                    tmpObj["isconfigable"] = organizeInfos[i].second;
                }
            }
            else
            {
                if (!_acsShareMgnt->IsUndistirbutedUser(userId)) {
                    vector<ncOrganizationInfo> organizeInfos;
                    _acsShareMgnt->GetAllOrgInfo (userId, organizeInfos);

                    for (size_t i = 0; i < organizeInfos.size (); ++i) {
                        depJson.push_back (JSON::OBJECT);

                        JSON::Object& tmpObj = depJson.back ().o ();
                        tmpObj["depid"] = organizeInfos[i].id.getCStr ();
                        tmpObj["name"] = organizeInfos[i].name.getCStr ();
                        tmpObj["isconfigable"] = true;
                    }
                }
            }
        }
        string body;
        JSON::Writer::write (replyJson.o (), body);
        ncHttpSendReply (cntl, brpc::HTTP_STATUS_OK, "ok", body);

        NC_EAC_HTTP_SERVER_TRACE (_T("this: %p, cntl: %p, userId: %s end"), this, cntl, userId.getCStr ());

    NC_EAC_HTTP_SERVER_CATCH (cntl)
}

// protected
void
ncEACDepartmentHandler::GetSubDeps (brpc::Controller* cntl, const String& userId)
{
    NC_EAC_HTTP_SERVER_TRY

        NC_EAC_HTTP_SERVER_TRACE (_T("this: %p, cntl: %p, userId: %s begin"), this, cntl, userId.getCStr ());

        // 获取http请求的content
        string bodyBuffer = cntl->request_attachment ().to_string ();

        // 获取docid,name
        JSON::Value requestJson;
        try {
            JSON::Reader::read (requestJson, bodyBuffer.c_str (), bodyBuffer.size ());
        }
        catch (Exception& e) {
            THROW_E (EAC_HTTP_SERVER, EACHTTP_INVALID_REQUEST_JSON_FORMAT,
                LOAD_STRING (_T("IDS_EACHTTP_INVALID_JSON")));
        }

        String depId = requestJson["depid"].s ().c_str ();
        INVALID_DEPARTMENT_ID(depId);

        // 检查部门id是否存在
        ncACSDepartInfo departInfo;
        bool ret = _acsShareMgnt->GetDepartInfoById (depId, departInfo);
        if (ret == false) {
            THROW_E (EAC_HTTP_SERVER, EACHTTP_DEPART_ID_NOT_EXISTS,
                LOAD_STRING (_T("IDS_DEPART_ID_NOT_EXISTS")));
        }

        // 回复
        JSON::Value replyJson;
        JSON::Array& depJson = replyJson["depinfos"].a ();

        if (_acsShareMgnt->GetPermShareLimitStatus()) {
            vector<pair<ncACSDepartInfo,bool> > depInfos;
            _acsShareMgnt->GetScopeSubDeps(userId, depId, depInfos);

            for (size_t i = 0; i < depInfos.size(); i++) {
                depJson.push_back (JSON::OBJECT);

                JSON::Object& tmpObj = depJson.back ().o ();
                tmpObj["depid"] = depInfos[i].first.id.getCStr ();
                tmpObj["name"] = depInfos[i].first.name.getCStr ();
                tmpObj["isconfigable"] = depInfos[i].second;
            }

        } else {
            vector<ncACSDepartInfo> depInfos;
            _acsShareMgnt->GetSubDeps (depId, depInfos);

            for (size_t i = 0; i < depInfos.size (); ++i) {
                depJson.push_back (JSON::OBJECT);

                JSON::Object& tmpObj = depJson.back ().o ();
                tmpObj["depid"] = depInfos[i].id.getCStr ();
                tmpObj["name"] = depInfos[i].name.getCStr ();
                tmpObj["isconfigable"] = true;
            }
        }

        string body;
        JSON::Writer::write (replyJson.o (), body);
        ncHttpSendReply (cntl, brpc::HTTP_STATUS_OK, "ok", body);

        NC_EAC_HTTP_SERVER_TRACE (_T("this: %p, cntl: %p, userId: %s end"), this, cntl, userId.getCStr ());

    NC_EAC_HTTP_SERVER_CATCH (cntl)
}

// protected
void
ncEACDepartmentHandler::GetSubUsers (brpc::Controller* cntl, const String& userId)
{
    NC_EAC_HTTP_SERVER_TRY

        NC_EAC_HTTP_SERVER_TRACE (_T("this: %p, cntl: %p, userId: %s begin"), this, cntl, userId.getCStr ());

        // 获取http请求的content
        string bodyBuffer = cntl->request_attachment ().to_string ();

        // 获取docid
        JSON::Value requestJson;
        try {
            JSON::Reader::read (requestJson, bodyBuffer.c_str (), bodyBuffer.size ());
        }
        catch (Exception& e) {
            THROW_E (EAC_HTTP_SERVER, EACHTTP_INVALID_REQUEST_JSON_FORMAT,
                LOAD_STRING (_T("IDS_EACHTTP_INVALID_JSON")));
        }

        String depId = requestJson["depid"].s ().c_str ();
        INVALID_DEPARTMENT_ID(depId)

        // 检查部门id是否存在
        ncACSDepartInfo departInfo;
        bool ret = _acsShareMgnt->GetDepartInfoById (depId, departInfo);
        if (ret == false) {
            THROW_E (EAC_HTTP_SERVER, EACHTTP_DEPART_ID_NOT_EXISTS,
                LOAD_STRING (_T("IDS_DEPART_ID_NOT_EXISTS")));
        }

        // 回复
        JSON::Value replyJson;
        JSON::Array& userJson = replyJson["userinfos"].a ();

        // 检查是否需要屏蔽用户信息
        if (_acsShareMgnt->GetShareMgntConfig("hide_user_info").compare("0") == 0) {
            vector<ncACSUserInfo> userInfos;
            if (_acsShareMgnt->GetPermShareLimitStatus())
                _acsShareMgnt->GetScopeSubUsers(userId, depId, userInfos);
             else
                _acsShareMgnt->GetSubUsers (depId, userInfos);

            for (size_t i = 0; i < userInfos.size (); ++i) {
                // 用户被禁用时界面上不显示
                if (userInfos[i].enableStatus == false) {
                    continue;
                }

                userJson.push_back (JSON::OBJECT);

                JSON::Object& tmpObj = userJson.back ().o ();
                tmpObj["userid"] = userInfos[i].id.getCStr ();
                tmpObj["account"] = userInfos[i].account.getCStr ();
                tmpObj["name"] = userInfos[i].visionName.getCStr ();
                tmpObj["mail"] = userInfos[i].email.getCStr ();
                tmpObj["csflevel"] = userInfos[i].csfLevel;
            }
        }
        string body;
        JSON::Writer::write (replyJson.o (), body);
        ncHttpSendReply (cntl, brpc::HTTP_STATUS_OK, "ok", body);

        NC_EAC_HTTP_SERVER_TRACE (_T("this: %p, cntl: %p, userId: %s end"), this, cntl, userId.getCStr ());

    NC_EAC_HTTP_SERVER_CATCH (cntl)
}

// protected
void
ncEACDepartmentHandler::Search (brpc::Controller* cntl, const String& userId)
{
    NC_EAC_HTTP_SERVER_TRY

        NC_EAC_HTTP_SERVER_TRACE (_T("this: %p, cntl: %p, userId: %s begin"), this, cntl, userId.getCStr ());

        // 获取http请求的content
        string bodyBuffer = cntl->request_attachment ().to_string ();

        JSON::Value requestJson;
        try {
            JSON::Reader::read (requestJson, bodyBuffer.c_str (), bodyBuffer.size ());
        }
        catch (Exception& e) {
            THROW_E (EAC_HTTP_SERVER, EACHTTP_INVALID_REQUEST_JSON_FORMAT,
                LOAD_STRING (_T("IDS_EACHTTP_INVALID_JSON")));
        }

        String key = requestJson["key"].s ().c_str ();
        if (key.isEmpty ()) {
            THROW_E (EAC_HTTP_SERVER, EACHTTP_INVALID_SEARCH_KEY,
                LOAD_STRING (_T("IDS_EACHTTP_INVALID_SEARCH_KEY")));
        }

        int start;
        if (requestJson["start"].type() == JSON::NIL) {
            start = 0;
        } else {
            start = requestJson["start"].i();
            if(start < 0) {
                THROW_E (EAC_HTTP_SERVER, INVALID_PAGE_START_VALUE,
                    "invalid start value, must bigger than 0");
            }
        }

        int limit;
        if (requestJson["limit"].type() == JSON::NIL) {
            limit = 10;
        } else {
            limit = requestJson["limit"].i();
            if(limit < -1) {
                THROW_E (EAC_HTTP_SERVER, INVALID_PAGE_LIMIT_VALUE,
                    "invalid limit value, must bigger than -1");
            }
        }

        // 进行搜索
        vector<ncACSUserInfo> userInfos;
        vector<ncACSDepartInfo> departInfos;
        bool limitPermShare = _acsShareMgnt->GetPermShareLimitStatus();
        if (limitPermShare)
            _acsShareMgnt->SearchScopeOrganization(userId, key, start, limit, userInfos, departInfos);
        else
            _acsShareMgnt->SearchAllOrganization (userId, key, start, limit, userInfos, departInfos);

        // 检查用户是否需要屏蔽组织架构信息
        bool hideOu = false;
#ifndef __UT__
        hideOu = ncEACHttpServerUtil::HideOum_Check(userId.getCStr ());
#endif
        // 检查是否需要屏蔽用户信息，值不为0，屏蔽
        bool hideUser = _acsShareMgnt->GetShareMgntConfig("hide_user_info").compare("0") != 0;

        map<String, ncACSDepartInfo> infoMap;
        if (!hideOu && !hideUser) {
            // 获取用户的直属部门信息
            if (!limitPermShare) {
                vector<String> userIds;
                for(size_t i = 0; i < userInfos.size (); ++i) {
                    userIds.push_back (userInfos[i].id);
                }
                _acsShareMgnt->GetBelongDepartByIdBatch (userIds, infoMap);
            }
            else {
                vector<String> departIds;
                for (size_t i = 0; i < userInfos.size (); ++i) {
                    departIds.push_back (userInfos[i].belongDepartId);
                }
                _acsShareMgnt->GetDepartInfoByIdBatch (departIds, infoMap);
            }
        }

        // 回复
        JSON::Value replyJson;
        JSON::Array& userJson = replyJson["userinfos"].a ();
        for (size_t i = 0; i < userInfos.size (); ++i) {
            userJson.push_back (JSON::OBJECT);

            JSON::Object& tmpObj = userJson.back ().o ();
            tmpObj["userid"] = userInfos[i].id.getCStr ();
            tmpObj["account"] = userInfos[i].account.getCStr ();
            tmpObj["name"] = userInfos[i].visionName.getCStr ();
            tmpObj["mail"] = userInfos[i].email.getCStr ();
            tmpObj["csflevel"] = userInfos[i].csfLevel;

            if (!hideOu && !hideUser) {
                ncACSDepartInfo info = limitPermShare ? infoMap[userInfos[i].belongDepartId] : infoMap[userInfos[i].id];
                tmpObj["depid"] = info.id.getCStr ();
                tmpObj["depname"] = info.name.getCStr ();
                String depPath = info.name;
                if (!info.id.isEmpty ()) {
                    String depParentPath;
                    _acsShareMgnt->GetParentDeptRootPathName(info.id, depParentPath);
                    if (!depParentPath.isEmpty ()) {
                        depPath = depParentPath + '/' + info.name;
                    }
                }
                tmpObj["deppath"] = depPath.getCStr ();
            }
            else {
                tmpObj["depid"] = "";
                tmpObj["depname"] = "";
                tmpObj["deppath"] = "";
            }
        }

        JSON::Array& depJson = replyJson["depinfos"].a ();
        for (size_t i = 0; i < departInfos.size (); ++i) {
            depJson.push_back (JSON::OBJECT);

            JSON::Object& tmpObj = depJson.back ().o ();
            tmpObj["depid"] = departInfos[i].id.getCStr ();
            tmpObj["name"] = departInfos[i].name.getCStr ();
            if (!hideOu) {
                tmpObj["path"] = departInfos[i].path.getCStr ();
            }
            else {
                tmpObj["path"] = "";
            }
        }

        string body;
        JSON::Writer::write (replyJson.o (), body);
        ncHttpSendReply (cntl, brpc::HTTP_STATUS_OK, "ok", body);

        NC_EAC_HTTP_SERVER_TRACE (_T("this: %p, cntl: %p, userId: %s end"), this, cntl, userId.getCStr ());

    NC_EAC_HTTP_SERVER_CATCH (cntl)
}

// protected
void
ncEACDepartmentHandler::SearchCount (brpc::Controller* cntl, const String& userId)
{
    NC_EAC_HTTP_SERVER_TRY

        NC_EAC_HTTP_SERVER_TRACE (_T("this: %p, cntl: %p, userId: %s begin"), this, cntl, userId.getCStr ());

        // 获取http请求的content
        string bodyBuffer = cntl->request_attachment ().to_string ();

        JSON::Value requestJson;
        try {
            JSON::Reader::read (requestJson, bodyBuffer.c_str (), bodyBuffer.size ());
        }
        catch (Exception& e) {
            THROW_E (EAC_HTTP_SERVER, EACHTTP_INVALID_REQUEST_JSON_FORMAT,
                LOAD_STRING (_T("IDS_EACHTTP_INVALID_JSON")));
        }

        String key = requestJson["key"].s ().c_str ();
        if (key.isEmpty ()) {
            THROW_E (EAC_HTTP_SERVER, EACHTTP_INVALID_SEARCH_KEY,
                LOAD_STRING (_T("IDS_EACHTTP_INVALID_SEARCH_KEY")));
        }

        // 进行搜索
        int searchCount;
        if (_acsShareMgnt->GetPermShareLimitStatus()) {
            searchCount = _acsShareMgnt->SearchScopeOrganizationCount(userId, key, true);
        } else {
            searchCount = _acsShareMgnt->SearchAllOrganizationCount (userId, key);
        }

        // 回复
        JSON::Value replyJson;
        replyJson["count"] = searchCount;

        string body;
        JSON::Writer::write (replyJson.o (), body);
        ncHttpSendReply (cntl, brpc::HTTP_STATUS_OK, "ok", body);

        NC_EAC_HTTP_SERVER_TRACE (_T("this: %p, cntl: %p, userId: %s end"), this, cntl, userId.getCStr ());

    NC_EAC_HTTP_SERVER_CATCH (cntl)
}
