#include "common/json.hpp"
#include "eachttpserver.h"
#include "ncEACUserHandler.h"

#include <ehttpserver/ncEHttpUtil.h>
#include "ncEACHttpServerUtil.h"
#include <sstream>

// public
ncEACUserHandler::ncEACUserHandler (ncIACSShareMgnt* acsShareMgnt)
    : _acsShareMgnt (acsShareMgnt)
{
    NC_EAC_HTTP_SERVER_TRACE (_T("this: %p"), this);

    NC_EAC_HTTP_SERVER_CHECK_ARGUMENT_NULL (acsShareMgnt);

    _methodFuncs.insert (pair<String, ncMethodFunc>(_T("get"), &ncEACUserHandler::Get));
    _methodFuncs.insert (pair<String, ncMethodFunc>(_T("getbasicinfo"), &ncEACUserHandler::GetBasicInfo));
    _methodFuncs.insert (pair<String, ncMethodFunc>(_T("agreedtotermsofuse"), &ncEACUserHandler::AgreedToTermsOfUse));
    _methodFuncs.insert (pair<String, ncMethodFunc>(_T("edit"), &ncEACUserHandler::EditUserInfo));
}

// public
ncEACUserHandler::~ncEACUserHandler (void)
{
    NC_EAC_HTTP_SERVER_TRACE (_T("this: %p"), this);
}

// public
void
ncEACUserHandler::doUserRequestHandler (brpc::Controller* cntl)
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

        // 消息处理
        ncMethodFunc func = iter->second;
        (this->*func) (cntl, introspectInfo);

    NC_EAC_HTTP_SERVER_CATCH (cntl)
}

// protected
void
ncEACUserHandler::Get (brpc::Controller* cntl, const ncIntrospectInfo& info)
{
    NC_EAC_HTTP_SERVER_TRY

        String userId = info.userId;
        NC_EAC_HTTP_SERVER_TRACE (_T("this: %p, cntl: %p, userId: %s begin"), this, cntl, userId.getCStr ());

        //判断token类型
        if (info.visitorType == ncTokenVisitorType::BUSINESS) {
            ncACSAppInfo appInfo;
            _acsShareMgnt->GetAppInfoById (userId, appInfo);

            JSON::Value replyJson;
            replyJson["type"] = "app";
            replyJson["id"] = appInfo.id.getCStr ();
            replyJson["name"] = appInfo.name.getCStr ();

            string body;
            JSON::Writer::write (replyJson.o (), body);
            ncHttpSendReply (cntl, brpc::HTTP_STATUS_OK, "ok", body);
            return;

        } else if (info.visitorType != ncTokenVisitorType::REALNAME) {
            THROW_E (EAC_HTTP_SERVER, EACHTTP_INVALID_PARAMETER,
                LOAD_STRING (_T("IDS_EACHTTP_INVALID_ARG_VALUE")), "token type");
        }

        // 获取用户信息
        ncACSUserInfo userInfo;
        _acsShareMgnt->GetUserInfoById (userId, userInfo);

        // 如果用户使用客户端登录并且非普通用户，则报错401
        if (info.clientType != ncClientType::CONSOLE_WEB
            && info.clientType != ncClientType::DEPLOY_WEB)
        {
            if (userId.getCStr() == g_ShareMgnt_constants.NCT_USER_ADMIN
                || userId.getCStr() == g_ShareMgnt_constants.NCT_USER_AUDIT
                || userId.getCStr() == g_ShareMgnt_constants.NCT_USER_SYSTEM
                || userId.getCStr() == g_ShareMgnt_constants.NCT_USER_SECURIT
                || userInfo.account == "admin"
                || userInfo.account == "audit"
                || userInfo.account == "system"
                || userInfo.account == "security")
            {
                THROW_E (EAC_HTTP_SERVER, EACHTTP_LOGIN_OSTYPE_IS_FORBID,
                    LOAD_STRING("IDS_LOGIN_OSTYPE_IS_FORBID"));
            }
        }

        int csfLevel = _acsShareMgnt->GetUserCSFLevel(userId, ncVisitorType::REALNAME);
        int csfLevel2 = _acsShareMgnt->GetUserCSFLevel2(userId, ncVisitorType::REALNAME);
        String csfLevelName = getCsfLevelName(csfLevel);
        String csfLevelName2 = getCsfLevel2Name(csfLevel2);
        int leakproofValue = _acsShareMgnt->GetLeakProofPerm(userId);
        int pwdcontrol = _acsShareMgnt->GetPwdControl(userId);
        int userauthtype = _acsShareMgnt->GetUserAuthType(userId);

        // 判断是否是组织管理员
        bool isManager = _acsShareMgnt->IsCustomDocManager(userId);

        bool isAuditor = false;

        // 回复
        JSON::Value replyJson;
        replyJson["type"] = "user";
        replyJson["userid"] = userInfo.id.getCStr ();
        replyJson["account"] = userInfo.account.getCStr ();
        replyJson["name"] = userInfo.visionName.getCStr ();
        replyJson["csflevel"] = csfLevel;
        replyJson["csflevel_name"] = csfLevelName.getCStr();
        replyJson["csflevel2"] = csfLevel2;
        replyJson["csflevel2_name"] = csfLevelName2.getCStr();
        replyJson["leakproofvalue"] = leakproofValue;
        replyJson["pwdcontrol"] = pwdcontrol;
        replyJson["usertype"] = userauthtype;
        replyJson["ismanager"] = isManager;
        replyJson["freezestatus"] = _acsShareMgnt->IsUserFreeze(userId);
        replyJson["agreedtotermsofuse"] = userInfo.isAgreedToTermsOfUse;

        // 兴业银行处理
        // 如果打开，则只显示前3位和后4位，中间多余的不显示，否则显示全部
        // 需要考虑空和不够长的问题
        String tempValue;
        bool result = _acsShareMgnt->GetCustomConfigOfString("show_hide_telnumber", tempValue);
        if (result) {
            if (tempValue.compare("1") == 0) {
                int telNumberLength = userInfo.telNumber.getLength();
                if (telNumberLength > 7) {
                    String telNumber = userInfo.telNumber.leftString(3) + "****" + userInfo.telNumber.rightString(4);
                    replyJson["telnumber"] = telNumber.getCStr();
                } else {
                    replyJson["telnumber"] = userInfo.telNumber.getCStr();
                }
            } else {
                replyJson["telnumber"] = userInfo.telNumber.getCStr();
            }
        } else {
            replyJson["telnumber"] = userInfo.telNumber.getCStr();
        }
        
        // 用户是否需要实名认证
        if (_acsShareMgnt->GetRealNameAuthStatus()) {
            replyJson["needrealnameauth"] = !_acsShareMgnt->IsUserRealNameAuth(userId);
        } else {
            replyJson["needrealnameauth"] = false;
        }

        String tempEmailConfigValue;
        String tempEmailValue = userInfo.email.getCStr();
        int emailLength = tempEmailValue.getLength();
        bool resultEmail = _acsShareMgnt->GetCustomConfigOfString("show_hide_email", tempEmailConfigValue);
        if (resultEmail && tempEmailConfigValue.compare("1") == 0 && emailLength > 1) {
            // 找到@符号的位置
            int atPos = tempEmailValue.find("@");        
            if (atPos > 0) {
                String tempEmail = tempEmailValue.leftString(1) + "****" +tempEmailValue.rightString(emailLength - atPos);
                replyJson["mail"] = tempEmail.getCStr();
            } else {
                replyJson["mail"] = tempEmailValue.getCStr();
            }
        } else {
            replyJson["mail"] = userInfo.email.getCStr();
        }

        // 获取用户的直属部门信息
        vector<String> directDeptIds;
        _acsShareMgnt->GetDirectBelongDepartmentIds(userId, directDeptIds);

        vector<ncACSDepartInfo> directDeptInfos;
        for ( size_t i = 0; i < directDeptIds.size(); ++i) {
            bool b_success = false;
            ncACSDepartInfo deptInfo;
            b_success = _acsShareMgnt->GetDepartInfoById(directDeptIds[i], deptInfo);
            if (b_success)
                directDeptInfos.push_back(deptInfo);
        }

        JSON::Array& depJson = replyJson["directdepinfos"].a ();
        for ( size_t i = 0; i < directDeptInfos.size(); ++i) {
            depJson.push_back(JSON::OBJECT);
            JSON::Object& tmpObj = depJson.back ().o ();
            tmpObj["depid"] = directDeptInfos[i].id.getCStr ();
            tmpObj["name"] = directDeptInfos[i].name.getCStr ();
        }

        JSON::Array& roleTypesJson = replyJson["roletypes"].a ();
        if(isAuditor) {
            roleTypesJson.push_back(ncACSRoleType::RT_AUDITOR);
        }

        // 获取用户的角色信息
        vector<ncRoleInfo> roleInfos;
        _acsShareMgnt->GetUserRole(userId, roleInfos);
        JSON::Array& roleJson = replyJson["roleinfos"].a ();
        for ( size_t i = 0; i < roleInfos.size(); ++i) {
            roleJson.push_back(JSON::OBJECT);
            JSON::Object& tmpObj = roleJson.back ().o ();
            tmpObj["id"] = roleInfos[i].id.getCStr ();
            tmpObj["name"] = roleInfos[i].name.getCStr ();
        }

        // 判断是否开启二次安全设备认证
        bool needsecondauth = false;

#ifndef __UT__
        ncTThirdPartyAuthConf config = ncEACHttpServerUtil::Usrm_GetThirdPartyAuth ();
        if (config.enabled) {
            JSON::Value configJson;
            JSON::Reader::read (configJson, config.config.c_str(), config.config.length());
            needsecondauth = configJson.get<bool> ("SecurityDevice", false);

            // 如果是城建设计院则需要检测内外网IP,内网时不需要进行二次安全设备认证
            if (config.thirdPartyId == _T("cjsjy")) {
                string realip = ncEACHttpServerUtil::GetForwardedIp(cntl).getCStr ();
                if(this->isLAN(realip)) {
                    needsecondauth = false;
                }
            }
        }
#endif
        replyJson["needsecondauth"] = needsecondauth;

        string body;
        JSON::Writer::write (replyJson.o (), body);
        ncHttpSendReply (cntl, brpc::HTTP_STATUS_OK, "ok", body);

        NC_EAC_HTTP_SERVER_TRACE (_T("this: %p, cntl: %p, userId: %s end"), this, cntl, userId.getCStr ());

    NC_EAC_HTTP_SERVER_CATCH (cntl)
}

// protected
void
ncEACUserHandler::GetBasicInfo (brpc::Controller* cntl, const ncIntrospectInfo& info)
{
    NC_EAC_HTTP_SERVER_TRY

        String userId = info.userId;
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

        String reqUserId = requestJson["userid"].s ().c_str ();
        INVALID_USER_ID(reqUserId);

        // 检查用户id是否存在
        ncACSUserInfo userInfo;
        bool ret = _acsShareMgnt->GetUserInfoById (reqUserId, userInfo);
        if (ret == false) {
            THROW_E (EAC_HTTP_SERVER, EACHTTP_USER_ID_NOT_EXIST,
                LOAD_STRING (_T("IDS_EACHTTP_USER_NOT_EXIST")));
        }

        // 回复
        JSON::Value replyJson;
        JSON::Array& depJson = replyJson["directdepinfos"].a ();

        // 检查用户是否需要屏蔽组织架构信息
        bool hideOu = false;
#ifndef __UT__
        hideOu = ncEACHttpServerUtil::HideOum_Check(userId.getCStr ());
#endif
        // 检查是否需要屏蔽用户信息，值不为0，屏蔽
        bool hideUser = _acsShareMgnt->GetShareMgntConfig("hide_user_info").compare("0") != 0;

        if (!hideOu && !hideUser) {
            // 获取用户的直属部门信息
            vector<String> directDeptIds;
            _acsShareMgnt->GetDirectBelongDepartmentIds(reqUserId, directDeptIds);

            vector<ncACSDepartInfo> directDeptInfos;
            for (size_t i = 0; i < directDeptIds.size(); ++i) {
                bool b_success = false;
                ncACSDepartInfo deptInfo;
                b_success = _acsShareMgnt->GetDepartInfoById(directDeptIds[i], deptInfo);
                if (b_success)
                    directDeptInfos.push_back(deptInfo);
            }

            for (size_t i = 0; i < directDeptInfos.size(); ++i) {
                depJson.push_back(JSON::OBJECT);
                JSON::Object& tmpObj = depJson.back ().o ();
                String id = directDeptInfos[i].id;
                String name = directDeptInfos[i].name;
                String depPath = name;
                if (!id.isEmpty ()) {
                    String depParentPath;
                    _acsShareMgnt->GetParentDeptRootPathName(id, depParentPath);
                    if (!depParentPath.isEmpty ()) {
                        depPath = depParentPath + '/' + name;
                    }
                }
                tmpObj["deppath"] = depPath.getCStr ();
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
ncEACUserHandler::AgreedToTermsOfUse (brpc::Controller* cntl, const ncIntrospectInfo& info)
{

    String userId = info.userId;
    NC_EAC_HTTP_SERVER_TRACE (_T("this: %p, cntl: %p, userId: %s begin"), this, cntl, userId.getCStr ());

    ncACSUserInfo userInfo;
    if(_acsShareMgnt->GetUserInfoById(userId, userInfo) == false) {
        THROW_E (EAC_HTTP_SERVER, EACHTTP_USER_NOT_EXIST,
                    "User not exists.");
    }

    // 同意用户协议
    _acsShareMgnt->AgreedToTermsOfUse (userId);

    // 回复
    JSON::Value replyJson;
    replyJson["result"] = "ok";

    string body;
    JSON::Writer::write (replyJson.o (), body);
    ncHttpSendReply (cntl, brpc::HTTP_STATUS_OK, "ok", body);

    NC_EAC_HTTP_SERVER_TRACE (_T("this: %p, cntl: %p, userId: %s end"), this, cntl, userId.getCStr ());

}

void
ncEACUserHandler::EditUserInfo (brpc::Controller* cntl, const ncIntrospectInfo& info)
{
    String userId = info.userId;
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

    // 获取用户信息
    ncACSUserInfo originalUserInfo;
    _acsShareMgnt->GetUserInfoById (userId, originalUserInfo);

    String msg;
    String exmsg;
    map<String, String> userinfo;
    bool emailExist = false;
    String emailAddress;
    bool telephoneExist = false;
    String telNumber;
    if (requestJson.o().find(_T("email")) != requestJson.o().end())
    {
        emailExist = true;
        emailAddress = toCFLString(requestJson["email"].s());
        if (!emailAddress.isEmpty())
        {
            // 解密
            string decodeEmailAddress(ncEACHttpServerUtil::Base64Decode(requestJson["email"].s().c_str()));
            emailAddress = toCFLString(ncEACHttpServerUtil::RSADecrypt2048(decodeEmailAddress));
        }
    }
    else if (requestJson.o().find(_T("emailaddress")) != requestJson.o().end())
    {
        emailExist = true;
        emailAddress = requestJson["emailaddress"].s().c_str();
    }
    if (emailExist)
    {
        userinfo.insert(pair<String, String>(_T("emailaddress"), emailAddress));
        if (emailAddress.isEmpty())
        {
            msg.format(LOAD_STRING(_T("IDS_EDIT_EMAIL_EMPTY")));
        }
        else
        {
            msg.format(LOAD_STRING(_T("IDS_EDIT_EMAIL")), emailAddress.getCStr());
        }
        exmsg.format(LOAD_STRING(_T("IDS_ORIGINAL_EMAIL")), originalUserInfo.email.getCStr());
    }

    if (requestJson.o().find(_T("displayname")) != requestJson.o().end())
    {
        String displayName = requestJson["displayname"].s().c_str();
        if (displayName.isEmpty())
        {
            THROW_E(EAC_HTTP_SERVER, EACHTTP_DISPLAYNAME_NOT_NULL,
                    LOAD_STRING(_T("IDS_DISPLAYNAME_NOT_NULL")));
        }
        else
        {
            userinfo.insert(pair<String, String>(_T("displayname"), displayName));
            msg.format(LOAD_STRING(_T("IDS_EDIT_DISPALYNAME")), displayName.getCStr());
            exmsg.format(LOAD_STRING(_T("IDS_ORIGINAL_DISPALYNAME")), originalUserInfo.visionName.getCStr());
        }
    }
    if (requestJson.o().find(_T("telephone")) != requestJson.o().end())
    {
        telephoneExist = true;
        telNumber = toCFLString(requestJson["telephone"].s());
        if (!telNumber.isEmpty())
        {
            // 解密
            string decodeTelNumber(ncEACHttpServerUtil::Base64Decode(requestJson["telephone"].s().c_str()));
            telNumber = toCFLString(ncEACHttpServerUtil::RSADecrypt2048(decodeTelNumber));
        }
    }
    else if (requestJson.o().find(_T("telnumber")) != requestJson.o().end())
    {
        telephoneExist = true;
        telNumber = requestJson["telnumber"].s().c_str();
    }

    if (telephoneExist)
    {
        if (telNumber.isEmpty())
        {
            msg.format(LOAD_STRING(_T("IDS_EDIT_TELNUMBER_EMPTY")));
        }
        else
        {
            msg.format(LOAD_STRING(_T("IDS_EDIT_TELNUMBER")), telNumber.getCStr());
        }
        exmsg.format(LOAD_STRING(_T("IDS_ORIGINAL_TELNUMBER")), originalUserInfo.telNumber.getCStr());
        userinfo.insert (pair<String,String>(_T("telnumber"), telNumber));
    }

    ncEACHttpServerUtil::SetUserInfo(userId, userinfo);
    ncHttpSendReply(cntl, brpc::HTTP_STATUS_OK, "ok", "");

    ncEACHttpServerUtil::Log(cntl, userId, ncTokenVisitorType::REALNAME, ncTLogType::NCT_LT_MANAGEMENT, ncTLogLevel::NCT_LL_WARN,
                            ncTManagementType::NCT_MNT_SET, msg, exmsg);

    NC_EAC_HTTP_SERVER_TRACE(_T("this: %p, cntl: %p, userId: %s end"), this, cntl, userId.getCStr());
}

// private:
bool
ncEACUserHandler::isLAN (const string& realip)
{

    /*
    以下IP范围为内网:
    A类 10.0.0.0     --  10.255.255.255
    B类 172.16.0.0   --  172.31.255.255
    C类 192.168.0.0  --  192.168.255.255
    */

    istringstream ipstream(realip);
    int ip[2];
    for(int i = 0; i < 2; i++) {
        string temp;
        getline(ipstream,temp,'.');
        istringstream t(temp);
        t >> ip[i];
    }
    if ((ip[0] == 10) || (ip[0] == 172 && ip[1] >= 16 && ip[1] <= 31) || (ip[0] == 192 && ip[1] == 168)) {
        return true;
    }
    return false;
}

// 获取服务端密级2配置
String ncEACUserHandler::getCsfLevel2Name(int csfLevel2)
{
    map<int, String> csflevels;
    // {"公开": 51}
    String csflevel2EnumStr = _acsShareMgnt->GetShareMgntConfig("csf_level2_enum");
    // 密级枚举为空，直接返回，防止下面转换出错
    if (csflevel2EnumStr.isEmpty()) {
        return "";
    }

    nlohmann::json csflevel2Enum = nlohmann::json::parse(csflevel2EnumStr.getCStr());
    for (const auto& it : csflevel2Enum.items()) {
        int level = it.value().get<int>();
        String name = it.key().c_str();
        csflevels[level] = name;
    }

    auto it = csflevels.find(csfLevel2);
    if (it != csflevels.end()) {
        return it->second;
    } else {
        return "";
    }

    return "";
}

// 获取服务端密级配置
String ncEACUserHandler::getCsfLevelName(int csfLevel)
{
    map<int, String> csflevels;
    // {"非密": 5, "内部": 6, "秘密": 7, "机密": 8}
    String csflevelEnumStr = _acsShareMgnt->GetShareMgntConfig("csf_level_enum");
    // 密级枚举为空，直接返回，防止下面转换出错
    if (csflevelEnumStr.isEmpty()) {
        return "";
    }

    nlohmann::json csflevelEnum = nlohmann::json::parse(csflevelEnumStr.getCStr());
    for (const auto& it : csflevelEnum.items()) {
        int level = it.value().get<int>();
        String name = it.key().c_str();
        csflevels[level] = name;
    }

    auto it = csflevels.find(csfLevel);
    if (it != csflevels.end()) {
        return it->second;
    } else {
        return "";
    }

    return "";
}
