#include <abprec.h>
#include <ncutil/ncBusinessDate.h>
#include <dataapi/ncJson.h>
#include <dataapi/ncGNSUtil.h>
#include "acsprocessor.h"
#include "ncMessageUtil.h"
#include "ncACSProcessorUtil.h"
#include "ncACSConfManager.h"
#include <boost/date_time/posix_time/posix_time.hpp>
#include <boost/date_time/local_time_adjustor.hpp>
#include <boost/date_time/c_local_time_adjustor.hpp>

ncMessageUtil::ncMessageUtil ()
{
    NC_ACS_PROCESSOR_TRACE (_T("this: %p"), this);

    nsresult ret;
    _acsShareMgnt = do_CreateInstance (NC_ACS_SHARE_MGNT_CONTRACTID, &ret);
    if (NS_FAILED (ret)) {
        THROW_E (ACS_PROCESSOR, FAILED_TO_CREATE_ACS_SHAREMGNT,
            _T("Failed to create acs sharemgnt: 0x%x"), ret);
    }

    _acsOwnerManager = do_CreateInstance (NC_ACS_OWNER_MANAGER_CONTRACTID, &ret);
    if (NS_FAILED (ret)) {
        THROW_E (ACS_PROCESSOR, FAILED_TO_CREATE_ACS_OWNER_MANANGER,
            _T("Failed to create acs owner manager: 0x%x"), ret);
    }

    _userManager = do_CreateInstance (USER_MANAGEMENT_CONTRACTID, &ret);
    if (NS_FAILED (ret)) {
        THROW_E (ACS_PROCESSOR, FAILED_TO_CREATE_ACS_OWNER_MANANGER,
            _T("Failed to create usermanagement instance: 0x%x"), ret);
    }
}

ncMessageUtil::~ncMessageUtil ()
{
    NC_ACS_PROCESSOR_TRACE (_T("this: %p"), this);
}

AB_DEFINE_THREADSAFE_SINGLETON_NO_POOL (ncMessageUtil)

void ncMessageUtil::GetAllUsers (const String & departId, vector<String> & userlist)
{
    std::vector<ncACSUserInfo> users;
    _acsShareMgnt->GetSubUsers (departId, users);
    for (size_t i = 0; i < users.size (); ++i)
        userlist.push_back (std::move (users[i].id));

    std::vector<ncACSDepartInfo> departs;
    _acsShareMgnt->GetSubDeps (departId, departs);
    for (size_t i = 0; i < departs.size (); ++i)
        GetAllUsers(departs[i].id, userlist);
}

String ncMessageUtil::AccessorTypeToStr (int accessorType)
{
    static map<int, String> sAccessorTypeMap;
    static ThreadMutexLock valueLock;

    if (sAccessorTypeMap.empty ()) {
        AutoLock<ThreadMutexLock> lock (&valueLock);
        if (sAccessorTypeMap.empty ()) {
            sAccessorTypeMap.insert (pair<int,String>(1, _T("user")));
            sAccessorTypeMap.insert (pair<int,String>(2, _T("department")));
            sAccessorTypeMap.insert (pair<int,String>(3, _T("contactor")));
        }
    }

    map<int, String>::iterator iter = sAccessorTypeMap.find (accessorType);
    if (iter != sAccessorTypeMap.end ()) {
        return iter->second;
    }
    else {
        return _T("unknown");
    }
}

void ncMessageUtil::RemoveDuplicateStrs (vector<String>& strs)
{
    // 先进行排序
    sort (strs.begin(), strs.end());

    // 在删除掉相邻重复的
    vector<String>::iterator pos = unique (strs.begin(), strs.end());

    // 删除掉最后无效的条目
    strs.erase (pos, strs.end());
}

void ncMessageUtil::CalcMsgReceivers (shared_ptr<acsMessageInfo> msg, std::vector<String>& receivers)
{
    receivers.clear();
    // 如果是配置权限或所有者，需要给访问者发送消息 (简单消息发送给访问者)
    if (msg->msgType == ACS_SHARE_OPEN_MSG || msg->msgType == ACS_SHARE_CLOSE_MSG ||
        msg->msgType == ACS_OWNER_SET_MSG || msg->msgType == ACS_OWNER_UNSET_MSG ||
        msg->msgType == ACS_SIMPLE_MSG) {

        if (msg->accessorType == 1) {
            // user
            receivers.push_back (msg->accessorId);
        }
        else if (msg->accessorType == 2) {
            // depart
            GetAllUsers (msg->accessorId, receivers);
        }
        else if (msg->accessorType == 3) {
            // contactor
            std::vector<ncACSUserInfo> userInfos;
            _acsShareMgnt->GetContactors (msg->accessorId, userInfos);
            for (size_t i = 0; i < userInfos.size (); ++i) {
                receivers.push_back (std::move (userInfos[i].id));
            }
        }
    }
    else if(msg->msgType == ACS_ANTIVIRUS_MSG) {
        acsMessageAntivirusInfo *antivirusMsg = static_cast<acsMessageAntivirusInfo*> (msg.get ());
        receivers.push_back (antivirusMsg->accessorId);
    }
    // 如果是隔离区消息，需要给文件所属文档库所有者发送消息
    else if(msg->msgType == ACS_QUARANTINE_MSG || msg->msgType == ACS_QUARANTINE_APPEAL_MSG
            || msg->msgType == ACS_APPEAL_APPROVE_MSG || msg->msgType == ACS_APPEAL_VOTE_MSG) {
        acsMessageQuarantineInfo *quarantineMsg = static_cast<acsMessageQuarantineInfo*> (msg.get ());
        _acsOwnerManager->GetOwnerIds (ncGNSUtil::GetCIDPath (quarantineMsg->docId), receivers);
    }
    // 如果是文件到期提醒，接收者在结构体中
    else if (msg->msgType == ACS_DOC_REMIND_MSG) {
        acsMessageDocRemindInfo *docDueMsg = static_cast<acsMessageDocRemindInfo*> (msg.get ());
        receivers = docDueMsg->receivers;
    }
}

// private
String ncMessageUtil::getUserNameById (const String& userId){
    try {
        ncACSUserInfo user;
        _acsShareMgnt->GetUserInfoById (userId, user);
        // 若t_user表找不到，去应用账户表查找
        if (user.visionName.isEmpty ()){
            AppInfo appInfo;
            _userManager->GetAppInfo (userId, appInfo);
            return appInfo.name;
        }else {
            return user.visionName;
        }
    }
    catch (Exception& e) {
        SystemLog::getInstance ()->log (__FILE__, __LINE__, ERROR_LOG_TYPE, _T("getUserName error: %s"), e.toString().getCStr());
        return "";
    }
    catch (JSON::Value& e) {
        SystemLog::getInstance ()->log (__FILE__, __LINE__, ERROR_LOG_TYPE, _T("getUserName error: %s"), e["body"].s ().c_str ());
        return "";
    }
    catch (...) {
        SystemLog::getInstance ()->log (__FILE__, __LINE__, ERROR_LOG_TYPE, _T("getUserName error: Unknown error"));
        return "";
    }
}