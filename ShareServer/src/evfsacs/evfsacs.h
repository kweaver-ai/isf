#ifndef __EVFS_ACS_H
#define __EVFS_ACS_H

#include "evfsacserr.h"
#include <ncRequestIDManager.h>

// 语言资源加载器
extern IResourceLoader* ncEVFSACSResLoader;
#define LOAD_STRING(strID)                        \
    ncEVFSACSResLoader->loadString (strID)

#define THROW_EVFS_ACS_ERROR(errMsg, errId)            \
    THROW_E (EVFS_ACS, errId, errMsg.getCStr ());

#ifdef __WINDOWS__
#define NC_EVFS_ACS_TRACE(...)            TRACE_EX2 (EVFS_ACS, __VA_ARGS__)
#else
#define NC_EVFS_ACS_TRACE(args...)        TRACE_REQID (EVFS_ACS, args)
#endif

#endif // __EVFS_ACS_H
