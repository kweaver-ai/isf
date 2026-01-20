#include "eachttpserver.h"
#include "ncEACPKIHandler.h"
#include "ncEACHttpServerUtil.h"

#include <ehttpclient/public/ncIEHTTPClient.h>

// 保证线程安全
#define BOOST_SPIRIT_THREADSAFE
#include <boost/property_tree/ptree.hpp>
#include <boost/foreach.hpp>
#include <boost/property_tree/xml_parser.hpp>
#include <boost/property_tree/ini_parser.hpp>
using namespace boost::property_tree;

// public
ncEACPKIHandler::ncEACPKIHandler (ncIACSShareMgnt* acsShareMgnt)
    : _acsShareMgnt (acsShareMgnt),
        _expires (24*3600)
{
    NC_EAC_HTTP_SERVER_TRACE (_T("this: %p"), this);

    _methodFuncs.insert (pair<String, ncMethodFunc>(_T("original"), &ncEACPKIHandler::onOriginal));
    _methodFuncs.insert (pair<String, ncMethodFunc>(_T("authen"), &ncEACPKIHandler::onAuthen));
}

// public
ncEACPKIHandler::~ncEACPKIHandler (void)
{
    NC_EAC_HTTP_SERVER_TRACE (_T("this: %p"), this);
}

// public
void
ncEACPKIHandler::doPKIRequestHandler (brpc::Controller* cntl)
{
    NC_EAC_HTTP_SERVER_TRY

        // 获取query string
        String method;
        ncHttpGetParams (cntl, method);

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

        // 消息处理
        ncMethodFunc func = iter->second;
        (this->*func) (cntl, "");

    NC_EAC_HTTP_SERVER_CATCH (cntl)
}

// protected
void
ncEACPKIHandler::onOriginal (brpc::Controller* cntl, const String& userId)
{
    NC_EAC_HTTP_SERVER_TRACE (_T("this: %p, cntl: %p, userId: %s begin"), this, cntl, userId.getCStr ());

    // process
    string original = getOriginal();

    JSON::Value replyJson;
    replyJson["original"] = original.c_str();

    // reply
    string body;
    JSON::Writer::write (replyJson.o (), body);
    ncHttpSendReply (cntl, brpc::HTTP_STATUS_OK, "ok", body);

    NC_EAC_HTTP_SERVER_TRACE (_T("this: %p, cntl: %p, userId: %s end"), this, cntl, userId.getCStr ());
}

// protected
void
ncEACPKIHandler::onAuthen (brpc::Controller* cntl, const String& userId)
{
    NC_EAC_HTTP_SERVER_TRACE (_T("this: %p, cntl: %p, userId: %s begin"), this, cntl, userId.getCStr ());

    // 获取http请求的content
    string bodyBuffer = cntl->request_attachment ().to_string ();

    // 获取original,detach
    JSON::Value requestJson;
    try {
        JSON::Reader::read (requestJson, bodyBuffer.c_str (), bodyBuffer.size ());
    }
    catch (Exception& e) {
        THROW_E (EAC_HTTP_SERVER, EACHTTP_INVALID_REQUEST_JSON_FORMAT,
            LOAD_STRING (_T("IDS_EACHTTP_INVALID_JSON")));
    }

    string original = requestJson["original"].s ();
    if (original.empty ()) {
        THROW_E (EAC_HTTP_SERVER, EACHTTP_ORIGINAL_INVALID,
            LOAD_STRING (_T("IDS_EACHTTP_ORIGINAL_INVALID")));
    }

    string detach = requestJson["detach"].s ();
    if (detach.empty ()) {
        THROW_E (EAC_HTTP_SERVER, EACHTTP_DETACH_INVALID,
                LOAD_STRING (_T("IDS_EACHTTP_DETACH_INVALID")));
    }

    // process
    string thirdId = getThirdId(original, detach);

    // 先根据第三方id查找，再根据帐号查找
    ncACSUserInfo userInfo;
    bool ret = _acsShareMgnt->GetUserInfoByThirdId (toCFLString(thirdId), userInfo);
    int accountType = 0;
    if (ret == false) {
        ret = _acsShareMgnt->GetUserInfoByAccount(toCFLString(thirdId), userInfo, accountType);

        if(ret == false) {
            THROW_E (EAC_HTTP_SERVER, EACHTTP_NOT_IMPORT_TO_ANYSHARE,
                LOAD_STRING (_T("IDS_EACHTTP_NOT_IMPORT_TO_ANYSHARE")));
        }
    }

    // 用户是否启用
    ret = _acsShareMgnt->IsUserEnabled(userInfo.id);
    if(ret == false) {
        THROW_E (EAC_HTTP_SERVER, EACHTTP_USER_DISABLED,
            LOAD_STRING (_T("IDS_USER_DISABLED")));
    }

#ifndef __UT__
    // 更新用户最近请求时间记录
    _acsShareMgnt->UpdateUserLastRequestTime(userInfo.id);
#endif

    // 回复
    JSON::Value replyJson;
    replyJson["userid"] = userInfo.id.getCStr ();
    replyJson["expires"] = _expires;

    string body;
    JSON::Writer::write (replyJson.o (), body);
    ncHttpSendReply (cntl, brpc::HTTP_STATUS_OK, "ok", body);

    String msg;
    msg.format (ncEACHttpServerLoader, _T("IDS_LOGIN_SUCCESS"), userInfo.account.getCStr ());
    ncEACHttpServerUtil::Log (cntl, userInfo.id, ncTokenVisitorType::REALNAME, ncTLogType::NCT_LT_LOGIN, ncTLogLevel::NCT_LL_INFO,
        ncTLoginType::NCT_CLT_LOGIN_IN, msg, "");

    NC_EAC_HTTP_SERVER_TRACE (_T("this: %p, cntl:%p, userId:%s end"), this, cntl, userId.getCStr ());
}

string ncEACPKIHandler::getOriginal()
{
    // http://$ip:$port/MessageService
    String serverAddress;
    String appId;
    getPKIServerInfo(serverAddress, appId);

    String url;
    url.format(_T("%s/MessageService"), serverAddress.getCStr());

    // build xml request
    string xmlTmp =
    "<?xml version=\"1.0\" encoding=\"UTF-8\"?>"
        "<message>"
        "<head>"
            "<version>1.1</version>"
            "<serviceType>OriginalService</serviceType>"
        "</head>"
        "<body>"
            "<appId>%s</appId>"
        "</body>"
    "</message>";
    String content;
    content.format(xmlTmp.c_str(), appId.getCStr());

    printMessage2(_T("%s"), content.getCStr());

    nsresult ret;
    nsCOMPtr<ncIEHTTPClient> httpClient = do_CreateInstance (NC_EHTTP_CLIENT_CONTRACTID, &ret);
    if (NS_FAILED (ret)) {
        THROW_E (EAC_HTTP_SERVER, FAILED_TO_CREATE_EHTTP_CLIENT,
                    "Failed to create ehttpclient instance: 0x%x", ret);
    }

    ncEHTTPResponse response;
    httpClient->Post (toSTLString(url), toSTLString(content), "application/xml", 30, response);

    printMessage2(_T("%s"), response.body.c_str());

    return parseOriginal(response.body);
}

string ncEACPKIHandler::parseOriginal (const string& retXMlStr)
{
    // Create an empty property tree object
    iptree ptAll;

    // Read the XML config string into the property tree. Catch any exception
    try {
        stringstream ss;
        ss << retXMlStr;
        read_xml(ss, ptAll);
    }
    catch (xml_parser_error &e) {
        THROW_E (EAC_HTTP_SERVER, INVALID_XML_FORMAT,
            LOAD_STRING (_T("IDS_INVALID_XML_FORMAT")));
    }

    try {
        /* 失败消息
        <message>
            <head>
                <version>1.1</version>
                <serviceType>AuthenService</serviceType>
                <messageState>true</messageState>
                <messageCode>-1</messageCode>
                <messageDesc>证书认证失败</messageDesc>
            </head>
            <body />
        </message>
        */
        boost::property_tree::iptree ptMessage = ptAll.get_child ("message");
        boost::property_tree::iptree ptHead = ptMessage.get_child ("head");
        string messageState = ptHead.get<string> ("messageState");

        if (messageState == "true") {
            string messageCode = ptHead.get<string> ("messageCode");
            string messageDesc = ptHead.get<string> ("messageDesc");
            THROW_E (EAC_HTTP_SERVER, CANT_AUTHENTICATE_TICKET, _T("%s: %s"), messageCode.c_str(), messageDesc.c_str());
        }

        /* 成功消息
        <?xml version="1.0" encoding="UTF-8"?>
        <message>
            <head><version>1.1</version>
                <serviceType>OriginalService</serviceType>
                <messageState></messageState>
                <messageCode></messageCode>
                <messageDesc></messageDesc>
            </head>
            <body>
                <original>85457440056748513228193848141804</original>
            </body>
        </message>
        */
        boost::property_tree::iptree ptBody = ptMessage.get_child ("body");
        string original = ptBody.get<string> ("original");
        return original;
    }
    catch (ptree_error& e) {
        THROW_E (EAC_HTTP_SERVER, INVALID_XML_FORMAT, e.what());
    }
}

string ncEACPKIHandler::getThirdId (const string& original, const string& detach)
{
    // http://$ip:$port/MessageService
    String serverAddress;
    String appId;
    getPKIServerInfo(serverAddress, appId);

    String url;
    url.format(_T("%s/MessageService"), serverAddress.getCStr());

    // build xml request
    string xmlTmp =
        "<?xml version=\"1.0\" encoding=\"UTF-8\"?>"
        "<message>"
            "<head>"
                "<version>1.1</version>"
                "<serviceType>AuthenService</serviceType>"
            "</head>"
            "<body>"
                "<appId>%s</appId>"
                "<authen>"
                    "<authCredential authMode=\"cert\">"
                        "<detach>%s</detach>"
                        "<original>%s</original>"
                    "</authCredential>"
                "</authen>"
                "<accessControl>false</accessControl>"
                "<attributes attributeType=\"all\"></attributes>"
            "</body>"
        "</message>";
    String content;

    // 对original进行base64加密，并去掉末尾的\n
    string base64Original = Base64Encode(original);
    base64Original.erase(std::remove(base64Original.begin(), base64Original.end(), '\n'), base64Original.end());

    content.format(xmlTmp.c_str(), appId.getCStr(), detach.c_str(), base64Original.c_str());

    printMessage2(_T("%s"), content.getCStr());

    nsresult ret;
    nsCOMPtr<ncIEHTTPClient> httpClient = do_CreateInstance (NC_EHTTP_CLIENT_CONTRACTID, &ret);
    if (NS_FAILED (ret)) {
        THROW_E (EAC_HTTP_SERVER, FAILED_TO_CREATE_EHTTP_CLIENT,
                    "Failed to create ehttpclient instance: 0x%x", ret);
    }

    ncEHTTPResponse response;
    httpClient->Post (toSTLString (url), toSTLString(content), "application/xml", 30, response);

    printMessage2(_T("%s"), response.body.c_str());

    map<string, string> userInfo;
    parseUserInfo(response.body, userInfo);

    for(map<string, string>::iterator iter = userInfo.begin(); iter != userInfo.end(); ++iter) {
        printMessage2(_T("%s=%s"), iter->first.c_str(), iter->second.c_str());
    }

    // "X509Certificate.SubjectDN" = "cn=董剑飞, T=123456789987654321, dc=jt, dc=cn" -> "T=123456789987654321"
    // "X509Certificate.SubjectDN" = "CN=刘俊,E=liujun05@yn.csg.cn,OU=06,OU=05,OU=CSG,O=SERC,C=CN"
    // "X509Certificate.SubjectDN" = "CN=李翔宇 120101198909233072,OU=00,OU=00,OU=24,L=84,L=00,ST=12,C=CN"
    String subjectDN;
    if (userInfo.find("X509Certificate.SubjectDN") != userInfo.end ()) {
        subjectDN = toCFLString(userInfo["X509Certificate.SubjectDN"]);
    }
    else if (userInfo.find("dnname") != userInfo.end ()) {
        subjectDN = toCFLString(userInfo["dnname"]);
    }

    vector<String> dnSplitStrs;
    subjectDN.split(",", dnSplitStrs);
    if(dnSplitStrs.size() <= 1) {
        THROW_E (EAC_HTTP_SERVER, INVALID_XML_FORMAT, "invalid X509Certificate.SubjectDN: %s", subjectDN.getCStr());
    }

    String cnStr = dnSplitStrs[0];
    // 先查看CN名字中是否包含身份证号
    vector<String> cnSpliStrs;
    cnStr.split("=", cnSpliStrs);
    if(cnSpliStrs.size() <= 1) {
        THROW_E (EAC_HTTP_SERVER, INVALID_XML_FORMAT, "invalid cn: %s", cnStr.getCStr());
    }

    String name = cnSpliStrs[1];
    vector<String> nameStrs;
    name.split(" ", nameStrs);

    String account;
    if(nameStrs.size() == 1) {
        // "X509Certificate.SubjectDN" = "cn=董剑飞, T=123456789987654321, dc=jt, dc=cn" -> "T=123456789987654321"
        // "X509Certificate.SubjectDN" = "CN=刘俊,E=liujun05@yn.csg.cn,OU=06,OU=05,OU=CSG,O=SERC,C=CN"
        String value = dnSplitStrs[1];

        // "T=123456789987654321" -> "123456789987654321"
        // "E=liujun05@yn.csg.cn" -> "liujun05@yn.csg.cn"
        vector<String> tmpStrs;
        value.split("=", tmpStrs);
        if(tmpStrs.size() <= 1) {
            THROW_E (EAC_HTTP_SERVER, INVALID_XML_FORMAT, "invalid account str: %s", value.getCStr());
        }
        account = tmpStrs[1];
    }
    else if(nameStrs.size() == 2 && nameStrs[1].getLength() >= 18) {
        // "X509Certificate.SubjectDN" = "CN=李翔宇 120101198909233072,OU=00,OU=00,OU=24,L=84,L=00,ST=12,C=CN"
        account = nameStrs[1];
    }
    else {
        THROW_E (EAC_HTTP_SERVER, INVALID_XML_FORMAT, "invalid X509Certificate.SubjectDN: %s", subjectDN.getCStr());
    }

    printMessage2(_T("account = %s"), account.getCStr());
    return toSTLString(account);
}

void ncEACPKIHandler::parseUserInfo (const string& retXMlStr, map<string, string>& userInfo)
{
    // Create an empty property tree object
    iptree ptAll;

    // Read the XML config string into the property tree. Catch any exception
    try {
        stringstream ss;
        ss << retXMlStr;
        read_xml(ss, ptAll);
    }
    catch(std::exception const&  ex) {
        THROW_E (EAC_HTTP_SERVER, INVALID_XML_FORMAT, ex.what());
    }

    /* 失败消息
    <message>
        <head>
            <version>1.1</version>
            <serviceType>AuthenService</serviceType>
            <messageState>true</messageState>
            <messageCode>-1</messageCode>
            <messageDesc>证书认证失败</messageDesc>
        </head>
        <body />
    </message>
    */
    boost::property_tree::iptree ptMessage = ptAll.get_child ("message");
    boost::property_tree::iptree ptHead = ptMessage.get_child ("head");
    string messageState = ptHead.get<string> ("messageState");

    if (messageState == "true") {
        string messageCode = ptHead.get<string> ("messageCode");
        string messageDesc = ptHead.get<string> ("messageDesc");
        String msg;
        msg.format(_T("%s(%s)"), messageDesc.c_str(), messageCode.c_str());
        THROW_E (EAC_HTTP_SERVER, CANT_AUTHENTICATE_TICKET, msg.getCStr());
    }

    /* 成功消息
    string retXMlStr = "<?xml version=\"1.0\" encoding=\"UTF-8\" ?>"
    "<message>"
        "<head>"
            "<version>1.1</version>"
            "<serviceType>AuthenService</serviceType>"
            "<messageState>false</messageState>"
        "</head>"
        "<body>"
            "<authResultSet allFailed=\"false\">"
                "<authResult authMode=\"password\" success=\"true\" />"
            "</authResultSet>"
            "<accessControlResult>Deny</accessControlResult>"
            "<attributes attributeType=\"all\">"
                "<attr name=\"UMS.UserID\" namespace=\"http://www.jit.com.cn/ums/ns/user\">111111</attr>"
                "<attr name=\"UMS.Username\" namespace=\"http://www.jit.com.cn/ums/ns/user\">test</attr>"
                "<attr name=\"UMS.LogonName\" namespace=\"http://www.jit.com.cn/ums/ns/user\">admin</attr>"
                "<attr name=\"privilege\" namespace=\"http://www.jit.com.cn/pmi/pms/ns/privilege\">管理员</attr>"
                "<attr name=\"role\" namespace=\"http://www.jit.com.cn/pmi/pms/ns/role\">经理</attr>"
                "<attr name=\"性别\" namespace=\"http://www.jit.com.cn/ums/ns/user\">男</attr>"
                "<attr name=\"职务\" namespace=\"http://www.jit.com.cn/ums/ns/user\">工程师</attr>"
                "<attr name=\"身份证\" namespace=\"http://www.jit.com.cn/ums/ns/user\">110110110110110110</attr>"
                "<attr name=\"部门\" namespace=\"http://www.jit.com.cn/ums/ns/user\">产品部</attr>"
            "</attributes>"
        "</body>"
    "</message>";
    */

    // Read the XML config string into the property tree. Catch any exception
    try {
        boost::property_tree::iptree ptBody = ptMessage.get_child ("body");
        boost::property_tree::iptree ptAttributes = ptBody.get_child ("attributes");

        BOOST_FOREACH(const iptree::value_type& v1, ptAttributes) {
            if(v1.first == "<xmlattr>") {
                continue;
            }

            const boost::property_tree::iptree& ptAttr = v1.second;
            string name = ptAttr.get<string>("<xmlattr>.name");
            string value = ptAttr.data();

            userInfo[name] = value;
        }
    }
    catch (ptree_error& e) {
        THROW_E (EAC_HTTP_SERVER, INVALID_XML_FORMAT, e.what());
    }
}

string ncEACPKIHandler::Base64Encode (const string& encPassword)
{
    string base64ed;
    base64ed.resize (encPassword.size ()*2, 0);
    int outl = 0;

    EVP_ENCODE_CTX* ctx = new EVP_ENCODE_CTX;
    EVP_EncodeInit(ctx);
    EVP_EncodeUpdate(ctx, (unsigned char *)&base64ed[0],
        &outl,
        (unsigned char *)encPassword.c_str (),
        (int)encPassword.size ());
    EVP_EncodeFinal(ctx,(unsigned char *)(&base64ed[0]+outl) ,&outl);

    delete ctx;
    ctx = NULL;

    return base64ed;
}

string ncEACPKIHandler::Base64Decode (const string& password)
{
    EVP_ENCODE_CTX* ctx = new EVP_ENCODE_CTX;
    EVP_DecodeInit(ctx);

    string decoded;
    decoded.resize (password.size (), 0);
    int outl = 0;

    EVP_DecodeUpdate(ctx, (unsigned char *)&decoded[0],
        &outl,
        (unsigned char *)password.c_str (),
        (int)password.size ());

    EVP_DecodeFinal(ctx,(unsigned char *)(&decoded[0] + outl) ,&outl);

    delete ctx;
    ctx = NULL;

    return decoded;
}

void ncEACPKIHandler::getPKIServerInfo(String& serverHost, String& appId)
{
    ncTThirdPartyAuthConf config = ncEACHttpServerUtil::Usrm_GetThirdPartyAuth ();

    // 检查第三方认证是否开启
    if(config.enabled == false || config.thirdPartyId != "jit") {
        THROW_E (EAC_HTTP_SERVER, EACHTTP_THIRD_AUTH_NOT_OPEN, LOAD_STRING (_T("IDS_EACHTTP_THIRD_AUTH_NOT_OPEN")));
    }

    // 获取json配置中的服务器地址和appid
    JSON::Value configJson;
    JSON::Reader::read (configJson, config.config.c_str(), config.config.length());

    serverHost = toCFLString(configJson["authServer"].s());
    appId = toCFLString(configJson["appId"].s());
}
