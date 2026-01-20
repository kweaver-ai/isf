#ifndef __ACS_SHAREMGNT_H
#define __ACS_SHAREMGNT_H

#include "acssharemgnterr.h"
#include <brpc/server.h>
#include <ncRequestIDManager.h>

// 语言资源加载器
extern IResourceLoader* ncACSShareMgntResLoader;
#define LOAD_STRING(strID)                        \
    ncACSShareMgntResLoader->loadString (strID)

#define NC_ACS_SHAREMGNT_TRACE TRACEPRINTF

#define SHAREMGNT_DB_NAME _T("sharemgnt_db")

// 获取数据库连接
class ncIDBOperator;
ncIDBOperator* ncACSShareMgntGetDBOperator ();

#endif // __ACS_SHAREMGNT_H
