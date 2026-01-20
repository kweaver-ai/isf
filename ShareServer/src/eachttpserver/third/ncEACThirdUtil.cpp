#include "eachttpserver.h"

#include <third/ncEACThirdUtil.h>
#include <ethriftutil/ncThriftClient.h>
#include "eacServiceAccessConfig.h"

/*
 * 处理sharemgnt调用异常
*/
void ncEACThirdUtil::HandleShareMgntException(ncTException & e)
{

}

bool ncEACThirdUtil::CheckAppOrgPerm (const String& appID, const ncAppPermOrgType& orgType, const ncAppOrgPermValue& perm)
{

    nsresult ret;
    nsCOMPtr<ncIACSShareMgnt> _acsShareMgnt = do_CreateInstance (NC_ACS_SHARE_MGNT_CONTRACTID, &ret);
    if (NS_FAILED (ret)) {
        THROW_E (EAC_HTTP_SERVER, FAILED_TO_CREATE_ACS_SHAREMGNT,
            _T("Faild to create acs sharemgnt: 0x%x"), ret);
    }

    // 获取权限
    int perm1 = _acsShareMgnt->GetAppOrgPerm (appID, orgType);

    // 检查权限
    if ((perm&perm1) != 0) {
        return true;
    }

    return false;
}
