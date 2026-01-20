#include "eachttpserver.h"
#include "ncEACMessageHandler.h"
#include "ncEACHttpServerUtil.h"

ncEACMessageHandler::ncEACMessageHandler (ncIACSMessageManager* acsMessageManager, ncIACSPermManager* acsPermManager, ncIACSShareMgnt* acsShareMgnt)
        : _acsMessageManager(acsMessageManager)
        , _acsPermManager(acsPermManager)
        , _acsShareMgnt(acsShareMgnt)
{
    NC_EAC_HTTP_SERVER_TRACE (_T("this: %p"), this);

    NC_EAC_HTTP_SERVER_CHECK_ARGUMENT_NULL (acsMessageManager);
    NC_EAC_HTTP_SERVER_CHECK_ARGUMENT_NULL (acsPermManager);

    _methodFuncs.insert (pair<String, ncMethodFunc>(_T("get"), &ncEACMessageHandler::Get));
    _methodFuncs.insert (pair<String, ncMethodFunc>(_T("read"), &ncEACMessageHandler::Read));
    _methodFuncs.insert (pair<String, ncMethodFunc>(_T("read2"), &ncEACMessageHandler::Read2));
    _methodFuncs.insert (pair<String, ncMethodFunc>(_T("sendmail"), &ncEACMessageHandler::SendMail));
}

ncEACMessageHandler::~ncEACMessageHandler (void)
{
    NC_EAC_HTTP_SERVER_TRACE (_T("this: %p"), this);
}

void
ncEACMessageHandler::doMessageRequestHandler (brpc::Controller* cntl)
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

void
ncEACMessageHandler::Get (brpc::Controller* cntl, const String& userId)
{
    NC_EAC_HTTP_SERVER_TRACE (_T("this: %p, cntl: %p, userId: %s begin"), this, cntl, userId.getCStr ());

    // 获取http请求的content
    string bodyBuffer = cntl->request_attachment ().to_string ();

    // 获取stamp
    JSON::Value requestJson;
    if (!bodyBuffer.empty ()) {
        try {
            JSON::Reader::read (requestJson, bodyBuffer.c_str (), bodyBuffer.size ());
        }
        catch (Exception& e) {
            THROW_E (EAC_HTTP_SERVER, EACHTTP_INVALID_REQUEST_JSON_FORMAT,
                LOAD_STRING (_T("IDS_EACHTTP_INVALID_JSON")));
        }
    }
    else {
        requestJson["stamp"] = 0;
    }
    std::vector<acsMessageResult> msgs;
    _acsMessageManager->GetMessageByUserId (userId, requestJson["stamp"].i (), msgs);

    int64 stamp = 0;
    if (!msgs.empty ())
    {
        stamp = msgs.back ().msgStamp;
    }

    JSON::Value replyJson;
    replyJson["stamp"] = stamp;
    JSON::Array& msgsJson = replyJson["msgs"].a ();
    JSON::Value tmpObj;
    for (size_t i = 0; i < msgs.size (); ++i) {
        try {
            JSON::Reader::read (tmpObj, msgs[i].msgContent.getCStr (), msgs[i].msgContent.getLength ());
            tmpObj["isread"] = msgs[i].msgStatus ? true : false;
            tmpObj["id"] = msgs[i].msgId.getCStr ();
            msgsJson.push_back (std::move (tmpObj));
        } catch (Exception& e) {}
    }

    // 回复
    string body;
    JSON::Writer::write (replyJson.o (), body);
    ncHttpSendReply (cntl, brpc::HTTP_STATUS_OK, "ok", body);

    NC_EAC_HTTP_SERVER_TRACE (_T("this: %p, cntl: %p, userId: %s end"), this, cntl, userId.getCStr ());
}

void
ncEACMessageHandler::Read (brpc::Controller* cntl, const String& userId)
{
    NC_EAC_HTTP_SERVER_TRACE (_T("this: %p, cntl: %p, userId: %s begin"), this, cntl, userId.getCStr ());

    // 获取http请求的content
    string bodyBuffer = cntl->request_attachment ().to_string ();

    // 获取stamp
    JSON::Value requestJson;
    try {
        JSON::Reader::read (requestJson, bodyBuffer.c_str (), bodyBuffer.size ());
    }
    catch (Exception& e) {
        THROW_E (EAC_HTTP_SERVER, EACHTTP_INVALID_REQUEST_JSON_FORMAT,
            LOAD_STRING (_T("IDS_EACHTTP_INVALID_JSON")));
    }

    _acsMessageManager->ReadMessageByUserId (userId, requestJson["stamp"].i ());

    // 回复
    string body;
    ncHttpSendReply (cntl, brpc::HTTP_STATUS_OK, "ok", body);

    NC_EAC_HTTP_SERVER_TRACE (_T("this: %p, cntl: %p, userId: %s end"), this, cntl, userId.getCStr ());
}

void
ncEACMessageHandler::Read2 (brpc::Controller* cntl, const String& userId)
{
    NC_EAC_HTTP_SERVER_TRACE (_T("this: %p, cntl: %p, userId: %s begin"), this, cntl, userId.getCStr ());

    // 获取http请求的content
    string bodyBuffer = cntl->request_attachment ().to_string ();

    // 获取msgids
    JSON::Value requestJson;
    try {
        JSON::Reader::read (requestJson, bodyBuffer.c_str (), bodyBuffer.size ());
    }
    catch (Exception& e) {
        THROW_E (EAC_HTTP_SERVER, EACHTTP_INVALID_REQUEST_JSON_FORMAT,
            LOAD_STRING (_T("IDS_EACHTTP_INVALID_JSON")));
    }

    vector<String> msgIds;
    JSON::Array& jsonConfigs = requestJson["msgids"].a ();
    for (size_t i = 0; i < jsonConfigs.size (); ++i) {
        String tmp = jsonConfigs[i].s ().c_str ();
        if (!ncOIDUtil::IsGUID (tmp)){
            THROW_E (EAC_HTTP_SERVER, EACHTTP_ARGUMENT_INVALID,
                        LOAD_STRING (_T("IDS_EACHTTP_INVALID_ARG_VALUE")), "msgids");
        }
        msgIds.push_back(tmp);
    }

    _acsMessageManager->ReadMessageByIds (userId, msgIds);

    // 回复
    string body;
    ncHttpSendReply (cntl, brpc::HTTP_STATUS_OK, "ok", body);

    NC_EAC_HTTP_SERVER_TRACE (_T("this: %p, cntl: %p, userId: %s end"), this, cntl, userId.getCStr ());
}

void
ncEACMessageHandler::SendMail (brpc::Controller* cntl, const String& userId)
{
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

    string subject = requestJson["subject"].s ().c_str ();
    string content = requestJson["content"].s ().c_str ();

    vector<string> mailto;
    JSON::Array& jsonConfigs = requestJson["mailto"].a ();
    for (size_t i = 0; i < jsonConfigs.size (); ++i) {
        string tmp = jsonConfigs[i].s ().c_str ();
        mailto.push_back(tmp);
    }
    ncEACHttpServerUtil::SendMail(mailto, subject, content);

    // 回复
    string body;
    ncHttpSendReply (cntl, brpc::HTTP_STATUS_OK, "ok", body);

    NC_EAC_HTTP_SERVER_TRACE (_T("this: %p, cntl: %p, userId: %s end"), this, cntl, userId.getCStr ());
}

void
ncEACMessageHandler::SendMessages (brpc::Controller* cntl)
{
    NC_EAC_HTTP_SERVER_TRY
    NC_EAC_HTTP_SERVER_TRACE (_T("this: %p, cntl: %p begin"), this, cntl);

    cntl->http_response().set_content_type("application/json; charset=UTF-8");
    if (brpc::HTTP_METHOD_POST != cntl->http_request ().method ()) {
        THROW_E (EAC_HTTP_SERVER, EACHTTP_QUERY_METHOD_INVALID,
            LOAD_STRING (_T("IDS_EACHTTP_QUERY_METHOD_NOT_SUPPORT")));
    }

    // 获取http请求的content
    string bodyBuffer = cntl->request_attachment ().to_string ();

    // 获取BODY PARAM
    JSON::Value requestJson;
    try {
        JSON::Reader::read (requestJson, bodyBuffer.c_str (), bodyBuffer.size ());
    }
    catch (Exception& e) {
        THROW_E (EAC_HTTP_SERVER, EACHTTP_INVALID_REQUEST_JSON_FORMAT,
            LOAD_STRING (_T("IDS_EACHTTP_INVALID_JSON")));
    }

    // 检测BODY PARAM
    map<String, JsonValueDesc> anyObj;
    map<String, JsonValueDesc> receiverObj;
    receiverObj["id"] = JsonValueDesc (JSON::STRING, true);
    receiverObj["account"] = JsonValueDesc (JSON::STRING, true);
    receiverObj["name"] = JsonValueDesc (JSON::STRING, true);
    receiverObj["email"] = JsonValueDesc (JSON::STRING, true);
    receiverObj["telephone"] = JsonValueDesc (JSON::STRING, true);
    receiverObj["third_attr"] = JsonValueDesc (JSON::STRING, true);
    receiverObj["third_id"] = JsonValueDesc (JSON::STRING, true);
    map<String, JsonValueDesc> receiverArray;
    receiverArray["element"] = JsonValueDesc (JSON::OBJECT, true, &receiverObj);
    map<String, JsonValueDesc> requestObj;
    requestObj["channel"] = JsonValueDesc (JSON::STRING, true);
    requestObj["content"] = JsonValueDesc (JSON::OBJECT, true, &anyObj);
    requestObj["receivers"] = JsonValueDesc (JSON::ARRAY, true, &receiverArray);
    map<String, JsonValueDesc> requestArray;
    requestArray["element"] = JsonValueDesc (JSON::OBJECT, true, &requestObj);
    JsonValueDesc requestValueDesc = JsonValueDesc (JSON::ARRAY, true, &requestArray);
    CheckRequestParameters ("body", requestJson, requestValueDesc);

    // 封装参数
    std::vector<std::shared_ptr<acsMessage>> messages;
    JSON::Array& requestObjArr = requestJson.a ();
    for (size_t i = 0; i < requestObjArr.size (); ++i) {
        std::shared_ptr<acsMessage> msgptr(new acsMessage());
        // channel
        msgptr->channel = requestObjArr[i]["channel"].s ().c_str ();
        if (msgptr->channel.isEmpty ()){
            String field;
            field.format (_T("body[%d].channel"), i);
            THROW_E (EAC_HTTP_SERVER, EACHTTP_INVALID_PARAMETER,
                LOAD_STRING (_T("IDS_EACHTTP_PARAMETERS_IS_INVALID")), field.getCStr ());
        }

        // content
        if (requestObjArr[i]["content"].o ().empty ()){
            String field;
            field.format (_T("body[%d].content"), i);
            THROW_E (EAC_HTTP_SERVER, EACHTTP_INVALID_PARAMETER,
                LOAD_STRING (_T("IDS_EACHTTP_PARAMETERS_IS_INVALID")), field.getCStr ());
        }
        std::string contstr;
        JSON::Writer::write (requestObjArr[i]["content"].o (), contstr);
        msgptr->content = std::move (toCFLString (contstr));


        // receivers
        JSON::Array& receiversArr = requestObjArr[i]["receivers"].a ();
        if (receiversArr.size () == 0){
            String field;
            field.format (_T("body[%d].receivers"), i);
            THROW_E (EAC_HTTP_SERVER, EACHTTP_INVALID_PARAMETER,
                LOAD_STRING (_T("IDS_EACHTTP_PARAMETERS_IS_INVALID")), field.getCStr ());
        }
        vector<messageReceiver> receivers;
        for (size_t j = 0; j < receiversArr.size (); ++j) {
            JSON::Object& receiverJson = receiversArr[j].o ();
            messageReceiver receiver;
            receiver.id = receiverJson["id"].s ().c_str ();
            receiver.account = receiverJson["account"].s ().c_str ();
            receiver.name = receiverJson["name"].s ().c_str ();
            receiver.email = receiverJson["email"].s ().c_str ();
            receiver.telephone = receiverJson["telephone"].s ().c_str ();
            receiver.thirdAttr = receiverJson["third_attr"].s ().c_str ();
            receiver.thirdId = receiverJson["third_id"].s ().c_str ();
            INVALID_USER_ID (receiver.id);
            if (receiver.account.isEmpty ()){
                String field;
                field.format (_T("body[%d].receivers[%d].account"), i, j);
                THROW_E (EAC_HTTP_SERVER, EACHTTP_INVALID_PARAMETER,
                    LOAD_STRING (_T("IDS_EACHTTP_PARAMETERS_IS_INVALID")), field.getCStr ());
            }
            if (receiver.name.isEmpty ()){
                String field;
                field.format (_T("body[%d].receivers[%d].name"), i, j);
                THROW_E (EAC_HTTP_SERVER, EACHTTP_INVALID_PARAMETER,
                    LOAD_STRING (_T("IDS_EACHTTP_PARAMETERS_IS_INVALID")), field.getCStr ());
            }

            receivers.push_back(receiver);
        }
        msgptr->receivers = receivers;
        messages.push_back (msgptr);
    }

    // 发送消息
    _acsMessageManager->AddMessage2 (messages);

    // 回复
    JSON::Value replyJson;
    JSON::Array& resultJson = replyJson.a ();
    String messageIds;
    for (auto msgIt = messages.begin (); msgIt != messages.end (); ++msgIt) {
        resultJson.push_back ((*msgIt)->msgId.getCStr ());
        messageIds.append ((*msgIt)->msgId);
        if (msgIt+1 != messages.end ()) {
            messageIds.append (",", 1);
        }
    }
    string body;
    JSON::Writer::write (replyJson.a (), body);
    String headerLocation;
    headerLocation.format(_T("/api/eacp/v1/message/%s"), messageIds.getCStr ());
    cntl->http_response ().SetHeader ("Location", headerLocation.getCStr ());
    ncHttpSendReply (cntl, brpc::HTTP_STATUS_CREATED, "Created", body);

    NC_EAC_HTTP_SERVER_TRACE (_T("this: %p, cntl: %p end"), this, cntl);
    NC_EAC_HTTP_SERVER_CATCH (cntl)
}

void
ncEACMessageHandler::ReadMessage (brpc::Controller* cntl)
{
    NC_EAC_HTTP_SERVER_TRY

    // 获取http请求的PATH PARAMETERS
    String unresolvedPath = cntl->http_request ().unresolved_path ().c_str ();
    vector<String> pathParams;
    unresolvedPath.split ('/', pathParams);
    if (pathParams.size () == 1) {
        ReadMessageForAllReceivers (cntl);
    }else if (pathParams.size () == 3) {
        if (pathParams[1] != "receivers"){
            THROW_E (EAC_HTTP_SERVER, EACHTTP_URI_NOT_EXIST,
                 _T("can not find the uri"));
        }
        ReadMessageForSomeReceivers (cntl);
    }else{
        THROW_E (EAC_HTTP_SERVER, EACHTTP_URI_NOT_EXIST,
                 _T("can not find the uri"));
    }

    NC_EAC_HTTP_SERVER_CATCH (cntl)
}

void
ncEACMessageHandler::ReadMessageForAllReceivers (brpc::Controller* cntl)
{
    NC_EAC_HTTP_SERVER_TRACE (_T("this: %p, cntl: %p begin"), this, cntl);

    cntl->http_response().set_content_type("application/json; charset=UTF-8");
    if (brpc::HTTP_METHOD_PUT != cntl->http_request ().method ()) {
        THROW_E (EAC_HTTP_SERVER, EACHTTP_QUERY_METHOD_INVALID,
            LOAD_STRING (_T("IDS_EACHTTP_QUERY_METHOD_NOT_SUPPORT")));
    }

    // 获取消息id
    String messageId = cntl->http_request ().unresolved_path ().c_str ();
    INVALID_UUID_IN_URI (messageId);

    // 获取http请求的content
    string bodyBuffer = cntl->request_attachment ().to_string ();

    // 获取BODY PARAM
    JSON::Value requestJson;
    try {
        JSON::Reader::read (requestJson, bodyBuffer.c_str (), bodyBuffer.size ());
    }
    catch (Exception& e) {
        THROW_E (EAC_HTTP_SERVER, EACHTTP_INVALID_REQUEST_JSON_FORMAT,
            LOAD_STRING (_T("IDS_EACHTTP_INVALID_JSON")));
    }

    // 检测BODY PARAM
    map<String, JsonValueDesc> requestObj;
    requestObj["read"] = JsonValueDesc (JSON::BOOLEAN, true);
    JsonValueDesc requestValueDesc = JsonValueDesc (JSON::OBJECT, true, &requestObj);
    CheckRequestParameters ("body", requestJson, requestValueDesc);
    if (!requestJson["read"].b ()){
        THROW_E (EAC_HTTP_SERVER, EACHTTP_INVALID_PARAMETER,
            LOAD_STRING (_T("IDS_EACHTTP_PARAMETERS_IS_INVALID")), _T("read"));
    }

    // 读取消息
    _acsMessageManager->ReadMessageForAllReceivers (messageId);

    // 回复
    ncHttpSendReply (cntl, brpc::HTTP_STATUS_OK, "ok", "");

    NC_EAC_HTTP_SERVER_TRACE (_T("this: %p, cntl: %p end"), this, cntl);
}

void
ncEACMessageHandler::ReadMessageForSomeReceivers (brpc::Controller* cntl)
{
    NC_EAC_HTTP_SERVER_TRACE (_T("this: %p, cntl: %p begin"), this, cntl);

    cntl->http_response().set_content_type("application/json; charset=UTF-8");
    if (brpc::HTTP_METHOD_PUT != cntl->http_request ().method ()) {
        THROW_E (EAC_HTTP_SERVER, EACHTTP_QUERY_METHOD_INVALID,
            LOAD_STRING (_T("IDS_EACHTTP_QUERY_METHOD_NOT_SUPPORT")));
    }

    // 获取http请求的PATH PARAMETERS
    String unresolvedPath = cntl->http_request ().unresolved_path ().c_str ();
    vector<String> pathParams;
    unresolvedPath.split ('/', pathParams);

    // 获取消息id和接收者id组
    String messageId = pathParams[0];
    vector<String> receiverIds;
    pathParams[2].split (',', receiverIds);

    INVALID_UUID_IN_URI (messageId);
    for (size_t i = 0; i < receiverIds.size (); ++i) {
        INVALID_UUID_IN_URI (receiverIds[i]);
    }

    // 获取http请求的content
    string bodyBuffer = cntl->request_attachment ().to_string ();

    // 获取BODY PARAM
    JSON::Value requestJson;
    try {
        JSON::Reader::read (requestJson, bodyBuffer.c_str (), bodyBuffer.size ());
    }
    catch (Exception& e) {
        THROW_E (EAC_HTTP_SERVER, EACHTTP_INVALID_REQUEST_JSON_FORMAT,
            LOAD_STRING (_T("IDS_EACHTTP_INVALID_JSON")));
    }

    // 检测BODY PARAM
    map<String, JsonValueDesc> requestObj;
    requestObj["read"] = JsonValueDesc (JSON::BOOLEAN, true);
    JsonValueDesc requestValueDesc = JsonValueDesc (JSON::OBJECT, true, &requestObj);
    CheckRequestParameters ("body", requestJson, requestValueDesc);
    if (!requestJson["read"].b ()){
        THROW_E (EAC_HTTP_SERVER, EACHTTP_INVALID_PARAMETER,
            LOAD_STRING (_T("IDS_EACHTTP_PARAMETERS_IS_INVALID")), _T("read"));
    }

    // 读取消息
    _acsMessageManager->ReadMessageForSomeReceivers (messageId, receiverIds);

    // 回复
    ncHttpSendReply (cntl, brpc::HTTP_STATUS_OK, "ok", "");

    NC_EAC_HTTP_SERVER_TRACE (_T("this: %p, cntl: %p end"), this, cntl);
}