#ifndef __NC_CLEAN_PERM_THREAD_H__
#define __NC_CLEAN_PERM_THREAD_H__

#include <acsdb/public/ncIDBPermManager.h>

/*
 * ncCleanPermThread
 */
class ncCleanPermThread : public Thread
{
public:
    ncCleanPermThread ();
    ~ncCleanPermThread ();

public:
    virtual void run ();

private:
    int64                               _cleanInterval;
    nsCOMPtr<ncIDBPermManager>          _dbPermManager;
};

#endif // __NC_CLEAN_PERM_THREAD_H__
