#ifndef __NC_EAC_ThirdDep_HANDLER_H__
#define __NC_EAC_ThirdDep_HANDLER_H__

#include "gen-cpp/ncTShareMgnt.h"
#include "gen-cpp/ShareMgnt_constants.h"

class ncEACThirdDepHandler
{
public:
    ncEACThirdDepHandler ();
    ~ncEACThirdDepHandler ();
    typedef void (ncEACThirdDepHandler::*ncMethodFunc) (brpc::Controller*, ncIntrospectInfo &);
    map<String, ncMethodFunc>            _methodFuncs;

protected:
    /***
     * 添加部门
     */
    void CreateDep (brpc::Controller* cntl, ncIntrospectInfo &info);
    void EditDep (brpc::Controller* cntl, ncIntrospectInfo &info);
    void DeleteDep (brpc::Controller* cntl, ncIntrospectInfo &info);
    void GetDepById (brpc::Controller* cntl, ncIntrospectInfo &info);
    void GetDepByThirdId (brpc::Controller* cntl, ncIntrospectInfo &info);
    void GetDepByName (brpc::Controller* cntl, ncIntrospectInfo &info);
    void MoveDep (brpc::Controller* cntl, ncIntrospectInfo &info);
    void AddUsersToDep (brpc::Controller* cntl, ncIntrospectInfo &info);
    void MoveUsersToDep (brpc::Controller* cntl, ncIntrospectInfo &info);
    void RemoveUsersFromDep (brpc::Controller* cntl, ncIntrospectInfo &info);
    void GetSubDepsByDepId (brpc::Controller* cntl, ncIntrospectInfo &info);
    void GetSubUsersByDepId (brpc::Controller* cntl, ncIntrospectInfo &info);
    void SetManager (brpc::Controller* cntl, ncIntrospectInfo &info);
    void CancelManager (brpc::Controller* cntl, ncIntrospectInfo &info);


private:

    JSON::Value ConvertDepartmentInfo (ncTUsrmDepartmentInfo & depInfo, bool needThirdId, bool showDepartPath);
    void ConvertOrganizationName (ncTUsrmDepartmentInfo & depInfo, ncTUsrmOrganizationInfo &orgInfo);
};

#endif  // __NC_EAC_CONFIG_HANDLER_H__
