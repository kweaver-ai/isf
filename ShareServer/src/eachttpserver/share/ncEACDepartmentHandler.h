#ifndef __NC_EAC_DEPARTMENT_HANDLER_H__
#define __NC_EAC_DEPARTMENT_HANDLER_H__

#include <acssharemgnt/public/ncIACSShareMgnt.h>

class ncEACDepartmentHandler
{
public:
    ncEACDepartmentHandler (ncIACSShareMgnt* acsShareMgnt);
    ~ncEACDepartmentHandler (void);

    void doDepRequestHandler (brpc::Controller* cntl);

protected:
    /***
     * 获取部门基本信息
     */
    void GetBasicInfo (brpc::Controller* cntl, const String& userId);

    /***
     * 获取用户所能访问的根部门信息
     */
    void GetRoots (brpc::Controller* cntl, const String& userId);

    /***
     * 获取子部门信息
     */
    void GetSubDeps (brpc::Controller* cntl, const String& userId);

    /***
     * 获取部门下的子用户信息
     */
    void GetSubUsers (brpc::Controller* cntl, const String& userId);

    /***
     * 搜索
     */
    void Search (brpc::Controller* cntl, const String& userId);

    /***
     * 获取搜索数目
     */
    void SearchCount (brpc::Controller* cntl, const String& userId);

private:
    typedef void (ncEACDepartmentHandler::*ncMethodFunc) (brpc::Controller*, const String&);
    map<String, ncMethodFunc>            _methodFuncs;

private:
    nsCOMPtr<ncIACSShareMgnt>            _acsShareMgnt;            // 查询sharemgnt管理
};

#endif  // __NC_EAC_DEPARTMENT_HANDLER_H__
