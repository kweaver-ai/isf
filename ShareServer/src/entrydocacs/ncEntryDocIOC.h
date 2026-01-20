#ifndef __NC_ENTRY_DOC_IOC_H
#define __NC_ENTRY_DOC_IOC_H

#include <acsprocessor/public/ncIACSCommon.h>
#include <acsprocessor/public/ncIACSPermManager.h>
#include <acsprocessor/public/ncIACSConfManager.h>
#include <acssharemgnt/public/ncIACSShareMgnt.h>
#include <acsprocessor/public/ncIACSOwnerManager.h>
#include <acsdb/public/ncIDBOwnerManager.h>

#include <entrydocioc/ncIEntryDocIOC.h>

/* Header file */
class ncEntryDocIOC : public ncIEntryDocIOC
{
public:
    NS_DECL_ISUPPORTS
    NS_DECL_NCIENTRYDOCIOC

    ncEntryDocIOC();

private:
    ~ncEntryDocIOC();

protected:
    nsCOMPtr<ncIACSPermManager>        _acsPermManager;
    nsCOMPtr<ncIACSShareMgnt>          _acsShareMgnt;
    nsCOMPtr<ncIACSConfManager>        _acsConfManager;
    nsCOMPtr<ncIACSOwnerManager>       _acsOwnerManager;
    nsCOMPtr<ncIDBOwnerManager>        _dbOwnerManager;
};

#endif // __NC_ENTRY_DOC_IOC_H
