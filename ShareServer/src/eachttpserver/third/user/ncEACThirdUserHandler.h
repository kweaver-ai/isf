#ifndef __NC_EAC_ThirdUser_HANDLER_H__
#define __NC_EAC_ThirdUser_HANDLER_H__

#include "gen-cpp/ncTShareMgnt.h"
#include "gen-cpp/ShareMgnt_constants.h"
#include "gen-cpp/ncTEFAST.h"
#include "gen-cpp/EFAST_constants.h"

class ncEACThirdUserHandler
{
public:
    ncEACThirdUserHandler ();
    ~ncEACThirdUserHandler ();
    typedef void (ncEACThirdUserHandler::*ncMethodFunc) (brpc::Controller*, ncIntrospectInfo &);
    map<String, ncMethodFunc>            _methodFuncs;
    static JSON::Value ConvertUserInfo (ncTUsrmGetUserInfo & userinfo, bool needThirdId);

protected:
    /***
     * 添加用户
     */
    void CreateUser (brpc::Controller* cntl, ncIntrospectInfo &info);
    void EditUser (brpc::Controller* cntl, ncIntrospectInfo &info);
    void DeleteUser (brpc::Controller* cntl, ncIntrospectInfo &info);
    void GetUserById (brpc::Controller* cntl, ncIntrospectInfo &info);
    void GetUserByThirdId (brpc::Controller* cntl, ncIntrospectInfo &info);
    void GetUserByName (brpc::Controller* cntl, ncIntrospectInfo &info);
    void GetAllUser (brpc::Controller* cntl, ncIntrospectInfo &info);
    void GetAllUserCount (brpc::Controller* cntl, ncIntrospectInfo &info);

private:
    String ConvertUserType(ncTUsrmGetUserInfo & userinfo);
    String ConvertCsfLevel(ncTUsrmGetUserInfo & userinfo);
    String ConvertPwdControl(ncTUsrmGetUserInfo & userinfo);

};

#endif  // __NC_EAC_CONFIG_HANDLER_H__
