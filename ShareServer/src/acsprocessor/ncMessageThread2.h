#ifndef __NC_MESSAGE_THREAD2_H__
#define __NC_MESSAGE_THREAD2_H__

#include <dataapi/ncJson.h>
#include <acsdb/public/ncIDBMessageManager.h>
#include <acssharemgnt/public/ncIACSShareMgnt.h>
#include <acsprocessor/public/ncIACSMessageManager.h>

/*
 * ncMessageThread2
 */
class ncMessageThread2 : public Thread
{
public:
    ncMessageThread2 ();
    ~ncMessageThread2 ();

public:
    virtual void run ();
    void AddMessage(const vector<shared_ptr<acsMessage>> & msgs);

private:
    nsCOMPtr<ncIDBMessageManager>       _dbMessageManager;
    nsCOMPtr<ncIACSShareMgnt>           _acsShareMgnt;

    queue<shared_ptr<acsMessage>>       _msgQueue;          // 任务队列
    Semaphore                           _queueSem;          // 队列信号
    ThreadMutexLock                     _queueMutex;        // 队列锁
};

#endif
