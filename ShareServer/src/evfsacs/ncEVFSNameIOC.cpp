#include <abprec.h>
#include <dataapi/ncGNSUtil.h>
#include <ncutil/ncPerformanceProfilerPrec.h>

#include "evfsacs.h"
#include "ncEVFSNameIOC.h"

/* Implementation file */
NS_IMPL_QUERY_INTERFACE1 (ncEVFSNameIOC, ncIEVFSNameIOC)

// protected
NS_IMETHODIMP_(nsrefcnt) ncEVFSNameIOC::AddRef (void)
{
    return 1;
}

// protected
NS_IMETHODIMP_(nsrefcnt) ncEVFSNameIOC::Release (void)
{
    return 1;
}

AB_DEFINE_THREADSAFE_SINGLETON_NO_POOL (ncEVFSNameIOC)

ncEVFSNameIOC::ncEVFSNameIOC()
    : _acsShareMgnt (NULL)
{
    NC_EVFS_ACS_TRACE (_T("this: %p"), this);

    // 创建 acsprocessor
    nsresult ret;

    _acsShareMgnt = do_CreateInstance (NC_ACS_SHARE_MGNT_CONTRACTID, &ret);
    if (NS_FAILED (ret)) {
        THROW_E (EVFS_ACS, FAILED_TO_CREATE_ACS_SHAREMGNT,
            _T("Failed to create acs sharemgnt: 0x%x"), ret);
    }
}

ncEVFSNameIOC::~ncEVFSNameIOC()
{
    NC_EVFS_ACS_TRACE (_T("this: %p"), this);
}

/* [notxpcom] String ConvertUserName ([const] in StringRef userId); */
NS_IMETHODIMP_(String) ncEVFSNameIOC::ConvertUserName(const String & userId)
{
    NC_EVFS_ACS_TRACE (_T("this: %p"), this);

    NC_ADD_PPN_CODE_CLIP_BEGIN_EX_EE (acs_ConvertUserName, _T("acs_ConvertUserName"));
    String userName;
    _acsShareMgnt->GetUserDisplayName (userId, userName);
    NC_ADD_PPN_CODE_CLIP_END_EX_EE (acs_ConvertUserName);
    return userName;
}

/* [notxpcom] String ConvertUserName ([const] in StringRef userId); */
NS_IMETHODIMP_(void) ncEVFSNameIOC::ConvertUserNameBatch(const vector<String> & userIds, map<String, String> & userNameMap)
{
    NC_EVFS_ACS_TRACE (_T("this: %p"), this);

    map<String, ncACSUserInfo> userInfoMap;
    _acsShareMgnt->GetUserInfoByIdBatch (userIds, userInfoMap);
    for (int i = 0; i < userIds.size (); ++i) {
        map<String, ncACSUserInfo>::iterator findIt = userInfoMap.find (userIds[i]);
        if (findIt != userInfoMap.end ()) {
            userNameMap.insert (pair<String, String> (userIds[i], findIt->second.visionName));
        }
    }
}

/* [notxpcom] String ConvertUserName ([const] in StringRef userId); */
NS_IMETHODIMP_(void) ncEVFSNameIOC::ConvertDepNameBatch(const vector<String> & depIds, map<String, String> & depNameMap)
{
    NC_EVFS_ACS_TRACE (_T("this: %p"), this);

    map<String, ncACSDepartInfo> depInfoMap;
    _acsShareMgnt->GetDepartInfoByIdBatch (depIds, depInfoMap);
    for (int i = 0; i < depIds.size (); ++i) {
        map<String, ncACSDepartInfo>::iterator findIt = depInfoMap.find (depIds[i]);
        if (findIt != depInfoMap.end ()) {
            depNameMap.insert (pair<String, String> (depIds[i], findIt->second.name));
        }
    }
}

/* [notxpcom] String GetDirectBelongDepartNames ([const] in StringRef userId, in StringVecRef departNames); */
NS_IMETHODIMP_(void) ncEVFSNameIOC::GetDirectBelongDepartNames(const String & userId, vector<String> & departNames)
{
    NC_EVFS_ACS_TRACE (_T("this: %p"), this);

    departNames.clear ();

    vector<String> depIds;
    _acsShareMgnt->GetDirectBelongDepartmentIds (userId, depIds);

    map<String, ncACSDepartInfo> depInfoMap;
    _acsShareMgnt->GetDepartInfoByIdBatch (depIds, depInfoMap);
    for (auto iter = depInfoMap.begin (); iter != depInfoMap.end (); ++iter) {
        departNames.push_back (iter->second.name);
    }
}

/* [notxpcom] String GetUserOSSId ([const] in StringRef userId); */
NS_IMETHODIMP_(String) ncEVFSNameIOC::GetUserOSSId(const String & userId)
{
    NC_EVFS_ACS_TRACE (_T("this: %p"), this);

    NC_ADD_PPN_CODE_CLIP_BEGIN_EX_EE (acs_GetUserOSSId, _T("acs_GetUserOSSId"));

    String ossId;
    _acsShareMgnt->GetUserOSSId (userId, ossId);

    NC_ADD_PPN_CODE_CLIP_END_EX_EE (acs_GetUserOSSId);
    return ossId;
}

