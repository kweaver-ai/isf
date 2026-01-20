/***************************************************************************************************
ossGateway.cpp:
    Copyright (c) Eisoo Software Inc. (2020), All rights reserved.

Purpose:
    ossGateway 服务接口调用

Author:
    Young.yu@aishu.cn

Creating Time:
    2023-05-30
***************************************************************************************************/
#include <abprec.h>
#include "ossGateway.h"

#include "drivenadapter.h"
#include <dataapi/ncJson.h>
#include <ossclient/public/ncIOSSClient.h>
#include <ossclient/ossclient.h>

#include "serviceAccessConfig.h"

/* Implementation file */
NS_IMPL_QUERY_INTERFACE1 (ossGateway, ossGatewayInterface)

// protected
NS_IMETHODIMP_(nsrefcnt) ossGateway::AddRef (void)
{
    return 1;
}

// protected
NS_IMETHODIMP_(nsrefcnt) ossGateway::Release (void)
{
    return 1;
}

AB_DEFINE_THREADSAFE_SINGLETON_NO_POOL (ossGateway)

ossGateway::ossGateway (): _ossClientPtr (0)
{
    NC_DRIVEN_ADAPTER_TRACE (_T("this: %p"), this);

    // http 服务设置获取
    _getLocalStorageInfoUrl.format (_T("%s://%s:%d/api/ossgateway/v1/local-storages?enabled=true"),
        ServiceAccessConfig::getInstance()->ossgatewayPrivateProtocol.getCStr(),
        ServiceAccessConfig::getInstance()->ossgatewayPrivateHost.getCStr(),
        ServiceAccessConfig::getInstance()->ossgatewayPrivatePort);
    _uploadInfoUrl.format (_T("%s://%s:%d/api/ossgateway/v1/upload-info"),
        ServiceAccessConfig::getInstance()->ossgatewayPrivateProtocol.getCStr(),
        ServiceAccessConfig::getInstance()->ossgatewayPrivateHost.getCStr(),
        ServiceAccessConfig::getInstance()->ossgatewayPrivatePort);
    _uploadPartUrl.format (_T("%s://%s:%d/api/ossgateway/v1/uploadpart"),
        ServiceAccessConfig::getInstance()->ossgatewayPrivateProtocol.getCStr(),
        ServiceAccessConfig::getInstance()->ossgatewayPrivateHost.getCStr(),
        ServiceAccessConfig::getInstance()->ossgatewayPrivatePort);
    _completeUploadUrl.format (_T("%s://%s:%d/api/ossgateway/v1/completeupload"),
        ServiceAccessConfig::getInstance()->ossgatewayPrivateProtocol.getCStr(),
        ServiceAccessConfig::getInstance()->ossgatewayPrivateHost.getCStr(),
        ServiceAccessConfig::getInstance()->ossgatewayPrivatePort);
    _getDownloadInfoUrl.format (_T("%s://%s:%d/api/ossgateway/v1/download"),
        ServiceAccessConfig::getInstance()->ossgatewayPrivateProtocol.getCStr(),
        ServiceAccessConfig::getInstance()->ossgatewayPrivateHost.getCStr(),
        ServiceAccessConfig::getInstance()->ossgatewayPrivatePort);
}

ossGateway::~ossGateway (void)
{
    NC_DRIVEN_ADAPTER_TRACE (_T("this: %p"), this);
}

void ossGateway::createOSSClient ()
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

/*NS_IMETHOD_(void) GetLocalOSSInfo(vector<OSSInfo>& infos);*/
NS_IMETHODIMP_(void) ossGateway::GetLocalOSSInfo (vector<OSSInfo>& infos )
{
    NC_DRIVEN_ADAPTER_TRACE (_T("[BEGIN] this: %p"), this);

    createOSSClient ();

    vector<string> headers;
    ncOSSResponse response;
    (*_ossClientPtr)->Get (_getLocalStorageInfoUrl.getCStr (), headers, 30, response);
    // 请求失败时，原样抛出错误信息
    if (response.code != 200){
        if (response.code == 0) {
            SystemLog::getInstance ()->log (__FILE__, __LINE__, ERROR_LOG_TYPE,
                _T("GetLocalOSSInfo url: %s, connect error"), _getLocalStorageInfoUrl.getCStr());
            THROW_E (ACS_DRIVEN_ADAPTER, DRIVEN_HTTP_SERVER_INTERNAL_ERR, _T("Could not connect to server"));
        }
        else{
            JSON::Value errorJson;
            errorJson["code"] = response.code;
            errorJson["body"] = response.body;

            SystemLog::getInstance ()->log (__FILE__, __LINE__, ERROR_LOG_TYPE,
                _T("GetLocalOSSInfo url: %s, code: %d, body:%s"), _getLocalStorageInfoUrl.getCStr(), response.code, response.body.c_str());
            throw errorJson;
        }
    }

    // 获取站点信息
    JSON::Value JconsentInfos;
    JSON::Reader::read(JconsentInfos, response.body.c_str(), response.body.length());
    for (size_t i = 0; i < JconsentInfos.a().size(); ++i)
    {
        OSSInfo tempInfo;
        tempInfo.id = toCFLString(JconsentInfos[i]["id"].s().c_str());
        tempInfo.bDefault = JconsentInfos[i]["default"].b();
        infos.push_back(tempInfo);
    }

    NC_DRIVEN_ADAPTER_TRACE (_T("[END] this: %p"), this);
}

/*NS_IMETHOD_(void) UploadInfo(const String & objName, const String & ossID, OSSUploadInfo& info); */
NS_IMETHODIMP_(void) ossGateway::UploadInfo (const String& objName, const String& ossID, OSSUploadInfo& info)
{
    NC_DRIVEN_ADAPTER_TRACE (_T("[BEGIN] this: %p"), this);

    createOSSClient ();

    String tmpUrl;
    int64 nFileSize = 1024*1024*1024; // 此URL为自适应分片/普通上传，在文件大小为1G时一定为分片上传，如果有问题，与oss网关联系
    tmpUrl.format("%s/%s/%s?internal_request=true&file_size=%lld&request_method=PUT", _uploadInfoUrl.getCStr(), ossID.getCStr(), objName.getCStr(), nFileSize);

    vector<string> headers;
    ncOSSResponse response;
    (*_ossClientPtr)->Get (tmpUrl.getCStr (), headers, 30, response);
    // 请求失败时，原样抛出错误信息
    if (response.code != 200){
        if (response.code == 0) {
            SystemLog::getInstance ()->log (__FILE__, __LINE__, ERROR_LOG_TYPE,
                _T("UploadInfo url: %s, connect error"), tmpUrl.getCStr());
            THROW_E (ACS_DRIVEN_ADAPTER, DRIVEN_HTTP_SERVER_INTERNAL_ERR, _T("Could not connect to server"));
        }
        else{
            JSON::Value errorJson;
            errorJson["code"] = response.code;
            errorJson["body"] = response.body;

            SystemLog::getInstance ()->log (__FILE__, __LINE__, ERROR_LOG_TYPE,
                _T("UploadInfo url: %s, code: %d, body:%s"), tmpUrl.getCStr(), response.code, response.body.c_str());
            throw errorJson;
        }
    }

    JSON::Value JconsentInfos;
    JSON::Reader::read(JconsentInfos, response.body.c_str(), response.body.length());
    info.uploadId = toCFLString(JconsentInfos["upload_id"].s().c_str());
    info.partSize = atoi(JconsentInfos["partsize"].s().c_str());

    NC_DRIVEN_ADAPTER_TRACE (_T("[END] this: %p"), this);
}

/*NS_IMETHOD_(void) UploadPart(const String & objName, const String & ossID, const String & uploadID, const int partNum, RequestInfo & info);*/
NS_IMETHODIMP_(void) ossGateway::UploadPart (const String& objName, const String& ossID, const String& uploadID, const int partNum, RequestInfo& res)
{
    NC_DRIVEN_ADAPTER_TRACE (_T("[BEGIN] this: %p"), this);

    createOSSClient();

    String tmpUrl;
    tmpUrl.format (_T("%s/%s/%s?internal_request=true&part_id=%d&upload_id=%s"),
        _uploadPartUrl.getCStr(),
        ossID.getCStr(),
        objName.getCStr(),
        partNum,
        uploadID.getCStr());
    vector<string> headers;
    ncOSSResponse response;

    (*_ossClientPtr)->Get (tmpUrl.getCStr (), headers, 30, response);
    // 请求失败时，原样抛出错误信息
    if (response.code != 200){
        if (response.code == 0) {
            SystemLog::getInstance ()->log (__FILE__, __LINE__, ERROR_LOG_TYPE,
                _T("UploadPart url: %s, connect error"), tmpUrl.getCStr());
            THROW_E (ACS_DRIVEN_ADAPTER, DRIVEN_HTTP_SERVER_INTERNAL_ERR, _T("Could not connect to server"));
        }
        else{
            JSON::Value errorJson;
            errorJson["code"] = response.code;
            errorJson["body"] = response.body;

            SystemLog::getInstance ()->log (__FILE__, __LINE__, ERROR_LOG_TYPE,
                _T("UploadPart url: %s, code: %d, body:%s"), tmpUrl.getCStr(), response.code, response.body.c_str());
            throw errorJson;
        }
    }

    // 确认返回值
    String partNumStr = String::toString(partNum);
    JSON::Value jsonVal;
    JSON::Reader::read (jsonVal, response.body.c_str (), response.body.length ());
    res.url = jsonVal[partNumStr.getCStr()]["url"].s();
    res.method = jsonVal[partNumStr.getCStr()]["method"].s();

    String header;
    auto hds = jsonVal[partNumStr.getCStr()]["headers"].o();
    for(auto& el : hds){
        header.format(_T("%s: %s"),el.first.c_str(),el.second.s().c_str());
        res.headers.push_back(toSTLString(header));
    }

    NC_DRIVEN_ADAPTER_TRACE (_T("[END] this: %p"), this);
}

/*NS_IMETHOD_(void) CompleteUpload(const String & objName, const String & ossID, const String & uploadID, const map<int, UploadPartInfo> & multiPartInfo, RequestInfo & res);*/
NS_IMETHODIMP_(void) ossGateway::CompleteUpload (const String& objName, const String& ossID, const String& uploadID, const map<int, UploadPartInfo>& multiPartInfo, RequestInfo& res)
{
    NC_DRIVEN_ADAPTER_TRACE (_T("[BEGIN] this: %p"), this);

    createOSSClient();

    String tmpUrl;
    tmpUrl.format (_T("%s/%s/%s?internal_request=true&upload_id=%s"),
        _completeUploadUrl.getCStr(),
        ossID.getCStr(),
        objName.getCStr(),
        uploadID.getCStr());
    vector<string> headers;
    headers.push_back("Content-Type: application/json");
    string reqBody("{");
    int cnt = 0;
    for(auto&& el : multiPartInfo){
        cnt++;
        String jsonKV;
        if(cnt >= multiPartInfo.size()){
            jsonKV.format(_T("\"%d\":\"%s\""),el.first,toCFLString(el.second.etag).trim ('"').getCStr());
            reqBody.append(::toSTLString(jsonKV));
            break;
        }
        jsonKV.format(_T("\"%d\":\"%s\","),el.first, toCFLString(el.second.etag).trim ('"').getCStr());
        reqBody.append(::toSTLString(jsonKV));
    }
    reqBody.append("}");

    ncOSSResponse response;
    (*_ossClientPtr)->Post (tmpUrl.getCStr (), reqBody, headers, 30, response);
    // 请求失败时，原样抛出错误信息
    if (response.code != 200){
        if (response.code == 0) {
            SystemLog::getInstance ()->log (__FILE__, __LINE__, ERROR_LOG_TYPE,
                _T("CompleteUpload url: %s , body: %s, connect error"), tmpUrl.getCStr(), reqBody.c_str());
            THROW_E (ACS_DRIVEN_ADAPTER, DRIVEN_HTTP_SERVER_INTERNAL_ERR, _T("Could not connect to server"));
        }
        else{
            JSON::Value errorJson;
            errorJson["code"] = response.code;
            errorJson["body"] = response.body;

            SystemLog::getInstance ()->log (__FILE__, __LINE__, ERROR_LOG_TYPE,
                _T("CompleteUpload url: %s, body:%s, code: %d, body:%s"), tmpUrl.getCStr(), reqBody.c_str(), response.code, response.body.c_str());
            throw errorJson;
        }
    }

    // 确认返回值
    JSON::Value jsonVal;
    JSON::Reader::read (jsonVal, response.body.c_str (), response.body.length ());
    res.url = jsonVal["url"].s();
    res.method = jsonVal["method"].s();
    res.body = jsonVal["request_body"].s();

    String header;
    auto hds = jsonVal["headers"].o();
    for(auto& el : hds){
        header.format(_T("%s: %s"),el.first.c_str(),el.second.s().c_str());
        res.headers.push_back(toSTLString(header));
    }

    NC_DRIVEN_ADAPTER_TRACE (_T("[END] this: %p"), this);
}

/*NS_IMETHOD_(void) GetDownLoadInfo(const String & objName, const String & ossID, const String & fileName, const int64 expireTime, RequestInfo & res)*/
NS_IMETHODIMP_(void) ossGateway::GetDownLoadInfo (const String& objName, const String& ossID, const String& fileName, const int64 expireTime, RequestInfo& res)
{
    NC_DRIVEN_ADAPTER_TRACE (_T("[BEGIN] this: %p"), this);

    createOSSClient();

    String tmpUrl;
    tmpUrl.format (_T("%s/%s/%s?internal_request=true&type=query_string"),
        _getDownloadInfoUrl.getCStr(),
        ossID.getCStr(),
        objName.getCStr());

    if (fileName != "") {
        String snQuery;
        snQuery.format("&save_name=%s", UrlEncode3986(fileName).getCStr());
        tmpUrl.append(snQuery);
    }

    vector<string> headers;

    ncOSSResponse response;
    (*_ossClientPtr)->Get (tmpUrl.getCStr (), headers, 30, response);
    // 请求失败时，原样抛出错误信息
    if (response.code != 200){
        if (response.code == 0) {
            SystemLog::getInstance ()->log (__FILE__, __LINE__, ERROR_LOG_TYPE,
                _T("GetDownLoadInfo url: %s , connect error"), tmpUrl.getCStr());
            THROW_E (ACS_DRIVEN_ADAPTER, DRIVEN_HTTP_SERVER_INTERNAL_ERR, _T("Could not connect to server"));
        }
        else{
            JSON::Value errorJson;
            errorJson["code"] = response.code;
            errorJson["body"] = response.body;

            SystemLog::getInstance ()->log (__FILE__, __LINE__, ERROR_LOG_TYPE,
                _T("GetDownLoadInfo url: %s, code: %d, body:%s"), tmpUrl.getCStr(), response.code, response.body.c_str());
            throw errorJson;
        }
    }

    // 确认返回值
    JSON::Value jsonVal;
    JSON::Reader::read (jsonVal, response.body.c_str (), response.body.length ());
    res.url = jsonVal["url"].s();
    res.method = jsonVal["method"].s();

    String header;
    auto hds = jsonVal["headers"].o();
    for(auto& el : hds){
        header.format(_T("%s: %s"),el.first.c_str(),el.second.s().c_str());
        res.headers.push_back(toSTLString(header));
    }

    NC_DRIVEN_ADAPTER_TRACE (_T("[END] this: %p"), this);
}

String ossGateway::UrlEncode3986 (const String &input)
{
    String escaped;

    int max = input.getLength();
    for(int i = 0; i < max; ++ i) {
        if ((48 <= input[i] && input[i] <= 57) ||    //0-9
            (65 <= input[i] && input[i] <= 90) ||    //abc...xyz
            (97 <= input[i] && input[i] <= 122) ||   //ABC...XYZ
            (input[i] == '_' || input[i] == '-' || input[i] == '~' || input[i] == '.')
            ) {
                escaped.append (input[i], 1);
        }
        else {
            escaped.append ("%");

            char dig1 = (input[i]&0xF0)>>4;
            char dig2 = (input[i]&0x0F);
            if ( 0 <= dig1 && dig1 <= 9) dig1 += 48;    //0,48inascii
            if (10 <= dig1 && dig1 <=15) dig1 += 65 - 10; //A,97inascii
            if ( 0 <= dig2 && dig2 <= 9) dig2 += 48;
            if (10 <= dig2 && dig2 <=15) dig2 += 65 - 10;

            String r;
            r.append (&dig1, 1);
            r.append (&dig2, 1);

            escaped.append (r);//converts char 255 to string "FF"
        }
    }

    return escaped;
}