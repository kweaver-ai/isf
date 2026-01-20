#ifndef __NC_DB_OUTBOX_MANAGER_H
#define __NC_DB_OUTBOX_MANAGER_H

#include "drivenadapter/public/nsqInterface.h"
#include <acsdb/public/ncIDBOutboxManager.h>

class ncIDBOperator;


/* Header file */
class ncDBOutboxManager : public ncIDBOutboxManager
{
    AB_DECLARE_THREADSAFE_SINGLETON (ncDBOutboxManager)

public:
    NS_DECL_ISUPPORTS
    NS_DECL_NCIDBOUTBOXMANAGER

    ncDBOutboxManager();
    ~ncDBOutboxManager();

private:
    void sendNSQ (const String& message);

    nsCOMPtr<nsqInterface>   _nsqManager;
    std::map<ncOutboxType,ncNSQEventType>        _outboxToNSQMap;
};

#endif // __NC_DB_OUTBOX_MANAGER_H
