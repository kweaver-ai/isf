#ifndef __NC_EAC_THIRD_UTIL_H
#define __NC_EAC_THIRD_UTIL_H

#include "gen-cpp/ncTShareMgnt.h"
#include "gen-cpp/ShareMgnt_constants.h"
#include <acssharemgnt/public/ncIACSShareMgnt.h>

class ncEACThirdUtil
{
public:

    // sharemgnt.thrift
    static void HandleShareMgntException(ncTException & e);

    static bool CheckAppOrgPerm (const String& appID, const ncAppPermOrgType& orgType, const ncAppOrgPermValue& perm);
};

#endif
