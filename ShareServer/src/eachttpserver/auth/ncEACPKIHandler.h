#ifndef __NC_EAC_PKI_HANDLER_H__
#define __NC_EAC_PKI_HANDLER_H__

#include <acssharemgnt/public/ncIACSShareMgnt.h>

class ncEACPKIHandler
{
public:
    ncEACPKIHandler (ncIACSShareMgnt* acsShareMgnt);
    ~ncEACPKIHandler (void);

    void doPKIRequestHandler (brpc::Controller* cntl);

protected:
    /***
     * 获取original
     */
    void onOriginal (brpc::Controller* cntl, const String& userId);

    /***
     * 进行证书认证
     */
    void onAuthen (brpc::Controller* cntl, const String& userId);

private:
    typedef void (ncEACPKIHandler::*ncMethodFunc) (brpc::Controller*, const String&);
    map<String, ncMethodFunc>            _methodFuncs;

    string getOriginal();
    string parseOriginal (const string& retXMlStr);

    string getThirdId (const string& original, const string& detach);
    void parseUserInfo (const string& retXMlStr, map<string, string>& userInfo);
    string Base64Encode (const string& encPassword);
    string Base64Decode (const string& password);

    void getPKIServerInfo(String& serverHost, String& appId);

private:
    nsCOMPtr<ncIACSShareMgnt>           _acsShareMgnt;      // sharemgnt db查询

    int64                               _expires;           // 默认的token有效期
};

#endif  // __NC_EAC_PKI_HANDLER_H__
