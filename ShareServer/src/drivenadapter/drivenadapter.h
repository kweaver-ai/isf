/***************************************************************************************************
drivenadapter.h:
    Copyright (c) Eisoo Software, Inc.(2020), All rights reserved

Purpose:
    drivenadapter

Author:
    Young.yu@aishu.cn

Creating Time:
    2020-11-17
***************************************************************************************************/
#ifndef __DRIVEN_ADAPETER_H
#define __DRIVEN_ADAPETER_H

#include "acsdrivenadaptererr.h"
#include <brpc/server.h>
#include <ncRequestIDManager.h>

// 语言资源加载器
extern IResourceLoader* ncDrivenAdapterLoader;
#define DRIVEN_LOAD_STRING(strID)                        \
    ncDrivenAdapterLoader->loadString (strID).getCStr ()

#define NC_DRIVEN_ADAPTER_TRACE TRACEPRINTF

#endif // __DRIVEN_ADAPETER_H
