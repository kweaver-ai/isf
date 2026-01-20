#ifndef __ACS_DATA_STORE_H__
#define __ACS_DATA_STORE_H__

// 语言资源加载器
//extern IResourceLoader* ncACSDataStoreResLoader;
//#define LOAD_STRING(strID)                        \
    ncACSDataStoreResLoader->loadString (strID)

#include <ncRequestIDManager.h>

#define ACS_DATA_STORE        _T("acsdatastore")

#ifdef __WINDOWS__
#define NC_ACS_DATA_STORE_TRACE(...)            TRACE_EX2 (ACS_DATA_STORE, __VA_ARGS__)
#else
#define NC_ACS_DATA_STORE_TRACE(args...)        TRACE_REQID (ACS_DATA_STORE, args)
#endif

#include <umm/umm.h>

//
// 定义acsdatastore使用的内存池
//
NC_DECLARE_UMM_ALLOCATOR (acsDataStorePoolAllocator);

// 内部错误
#define ACS_DATA_STORE_INTERNAL_ERR                0x00000001L            // 内部错误
#define ACS_DATA_STORE_INVALID_ARGUMENT            0x00000002L            // 无效参数
#define ACS_DATA_STORE_INVALID_OPERATION           0x00000003L            // 非法操作

#endif // __ACS_DATA_STORE_H__
