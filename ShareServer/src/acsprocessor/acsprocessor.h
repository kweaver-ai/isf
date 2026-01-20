#ifndef __ACS_PROCESSOR_H
#define __ACS_PROCESSOR_H

#include "acsprocessorerr.h"
#include <brpc/server.h>
#include <ncRequestIDManager.h>

// 语言资源加载器
extern IResourceLoader* ncACSProcessorResLoader;
#define LOAD_STRING(strID)                        \
    ncACSProcessorResLoader->loadString (strID).getCStr ()

#define NC_ACS_PROCESSOR_TRACE TRACEPRINTF

// 内部错误
#define ACS_PROCESSOR_INTERNAL_ERR                0x00000001L            // 内部错误
#define ACS_PROCESSOR_INVALID_ARGUMENT            0x00000002L            // 无效参数
#define ACS_PROCESSOR_INVALID_OPERATION           0x00000003L            // 非法操作

#endif // __ACS_PROCESSOR_H
