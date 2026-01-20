#ifndef __ENTRY_DOC_ACS_H
#define __ENTRY_DOC_ACS_H

#include "entrydocacserr.h"
#include <ncRequestIDManager.h>

#define THROW_ENTRY_DOC_ACS_ERROR(errMsg, errId)            \
    THROW_E (ENTRY_DOC_ACS, errId, errMsg.getCStr ());

#ifdef __WINDOWS__
#define NC_ENTRY_DOC_ACS_TRACE(...)            TRACE_EX2 (ENTRY_DOC_ACS, __VA_ARGS__)
#else
#define NC_ENTRY_DOC_ACS_TRACE(args...)        TRACE_REQID (ENTRY_DOC_ACS, args)
#endif

#endif // __ENTRY_DOC_ACS_H
