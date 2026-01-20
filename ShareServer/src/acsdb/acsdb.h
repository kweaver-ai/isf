#ifndef __ACS_DB_H
#define __ACS_DB_H

#include "acsdberr.h"
#include <brpc/server.h>
#include <ncRequestIDManager.h>

// 语言资源加载器
extern IResourceLoader* ncACSDBResLoader;
#define LOAD_STRING(strID)                        \
    ncACSDBResLoader->loadString (strID)

#define NC_ACS_DB_TRACE TRACEPRINTF

#define ANYSHARE_DB_NAME                    _T ("anyshare")

// 获取数据库连接
class ncIDBOperator;
ncIDBOperator* ncACSDBGetDBOperator (int timeout = 0);

#endif // __ACS_DB_H
