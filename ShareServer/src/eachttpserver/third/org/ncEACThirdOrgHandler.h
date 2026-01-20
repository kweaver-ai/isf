#ifndef __NC_EAC_ThirdOrg_HANDLER_H__
#define __NC_EAC_ThirdOrg_HANDLER_H__

#include "gen-cpp/ncTShareMgnt.h"
#include "gen-cpp/ShareMgnt_constants.h"


class ncEACThirdOrgHandler
{
public:
    ncEACThirdOrgHandler ();
    ~ncEACThirdOrgHandler ();
    typedef void (ncEACThirdOrgHandler::*ncMethodFunc) (brpc::Controller*, ncIntrospectInfo &);
    map<String, ncMethodFunc>            _methodFuncs;
    static JSON::Value ConvertDepInfo (ncTDepartmentInfo & depInfo);

protected:
    /***
     * 添加组织
     */
    void CreateOrg (brpc::Controller* cntl, ncIntrospectInfo &info);
    void EditOrg (brpc::Controller* cntl, ncIntrospectInfo &info);
    void DeleteOrg (brpc::Controller* cntl, ncIntrospectInfo &info);
    void GetAllOrg (brpc::Controller* cntl, ncIntrospectInfo &info);
    void GetOrgById (brpc::Controller* cntl, ncIntrospectInfo &info);
    void GetOrgByName (brpc::Controller* cntl, ncIntrospectInfo &info);
    void GetSubDepByOrgId (brpc::Controller* cntl, ncIntrospectInfo &info);
    void GetSubUserByOrgId (brpc::Controller* cntl, ncIntrospectInfo &info);

private:
    JSON::Value ConvertOrgInfo (ncTRootOrgInfo & orgInfo, bool needThirdId);
    void ConvertOrganizationInfo (ncTUsrmOrganizationInfo &orgInfo, ncTRootOrgInfo & _return);
};

#endif  // __NC_EAC_CONFIG_HANDLER_H__
