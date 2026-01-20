#ifndef __NC_EAC_USER_HANDLER_H__
#define __NC_EAC_USER_HANDLER_H__

#include <acsprocessor/public/ncIACSTokenManager.h>
#include <acssharemgnt/public/ncIACSShareMgnt.h>

class ncEACUserHandler
{
public:
    ncEACUserHandler (ncIACSShareMgnt* acsShareMgnt);
    ~ncEACUserHandler (void);

    void doUserRequestHandler (brpc::Controller* cntl);

protected:
    /***
     * 获取token对应的用户信息
     */
    void Get (brpc::Controller* cntl, const ncIntrospectInfo& info);

    /***
     * 获取用户基本信息
     */
    void GetBasicInfo (brpc::Controller* cntl, const ncIntrospectInfo& info);

    /***
     * 同意用户协议
     */
    void AgreedToTermsOfUse (brpc::Controller* cntl, const ncIntrospectInfo& info);

    /***
     * 获取所有用户信息
     */
    void GetAll (brpc::Controller* cntl, const String& userId);

    /***
     * 编辑用户信息
     */
    void EditUserInfo (brpc::Controller* cntl, const ncIntrospectInfo& info);

private:
    typedef void (ncEACUserHandler::*ncMethodFunc) (brpc::Controller*, const ncIntrospectInfo&);
    map<String, ncMethodFunc>            _methodFuncs;

    /***
     * 判断是否是内网
     */
    bool isLAN (const string& realip);

    // 获取服务端密级配置
    String getCsfLevelName(int csfLevel);

    // 获取服务端密级2配置
    String getCsfLevel2Name(int csfLevel2);

private:
    nsCOMPtr<ncIACSShareMgnt>            _acsShareMgnt;            // 查询sharemgnt管理
};

#endif  // __NC_EAC_USER_HANDLER_H__
