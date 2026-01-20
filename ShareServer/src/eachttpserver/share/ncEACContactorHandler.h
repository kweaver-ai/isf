#ifndef __NC_EAC_CONTACTOR_HANDLER_H__
#define __NC_EAC_CONTACTOR_HANDLER_H__

#include <acssharemgnt/public/ncIACSShareMgnt.h>

#include "gen-cpp/ncTShareMgnt.h"
#include "gen-cpp/ShareMgnt_constants.h"

class ncEACContactorHandler
{
public:
    ncEACContactorHandler (ncIACSShareMgnt* acsShareMgnt);
    ~ncEACContactorHandler (void);

    void doContactorRequestHandler (brpc::Controller* cntl);

protected:
    /***
     * 获取分组下所有联系人
     */
    void GetContactors (brpc::Controller* cntl, const String& userId);

    /***
     * 获取所有分组
     */
    void GetGroup (brpc::Controller* cntl, const String& userId);

    /***
     * 获取搜索用户信息数目
     */
    void SearchCount (brpc::Controller* cntl, const String& userId);

    /***
     * 增加联系人组
     */
    void AddGroup (brpc::Controller* cntl, const String& userId);

    /***
     * 编辑联系人组
     */
    void EditGroup (brpc::Controller* cntl, const String& userId);

    /***
     * 增加联系人
     */
    void AddPersons (brpc::Controller* cntl, const String& userId);

    /***
     * 删除联系人
     */
    void DeletePersons (brpc::Controller* cntl, const String& userId);

    /***
     * 搜索用户信息
     */
    void Search (brpc::Controller* cntl, const String& userId);

    /***
     * 搜索联系人信息
     */
    void SearchPersons (brpc::Controller* cntl, const String& userId);

    /***
     * 分页获取联系人组的用户信息
     */
    void GetPersonFromGroup (brpc::Controller* cntl, const String& userId);

private:
    typedef void (ncEACContactorHandler::*ncMethodFunc) (brpc::Controller*, const String&);
    map<String, ncMethodFunc>            _methodFuncs;
    static void HandlencTException(ncTException & e);

private:
    nsCOMPtr<ncIACSShareMgnt>            _acsShareMgnt;            // 查询sharemgnt管理
};

#endif  // __NC_EAC_USER_HANDLER_H__
