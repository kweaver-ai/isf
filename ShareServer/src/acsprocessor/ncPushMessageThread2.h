#ifndef __NC_PUSH_MESSAGE_THREAD2_H__
#define __NC_PUSH_MESSAGE_THREAD2_H__

#include <dataapi/ncJson.h>
#include <acsdb/public/ncIDBConfManager.h>
#include <acssharemgnt/public/ncIACSShareMgnt.h>
#include <acsprocessor/public/ncIACSMessageManager.h>
#include <acsprocessor/public/ncIACSConfManager.h>
#include <drivenadapter/public/pluginMessageInterface.h>

/*
 * ncPushMessageThread2
 */
class ncPushMessageThread2 : public Thread
{
public:
    ncPushMessageThread2 ();
    ~ncPushMessageThread2 ();

public:
    virtual void run ();

    void AddMessage(const vector<shared_ptr<acsMessage>> & msgs);

private:
    nsCOMPtr<ncIACSShareMgnt>           _acsShareMgnt;
    nsCOMPtr<ncIDBConfManager>          _dbConfManager;
    nsCOMPtr<ncIACSConfManager>         _acsConfManager;
    nsCOMPtr<pluginMessageInterface>    _pluginMessage;

    queue<shared_ptr<acsMessage>>       _msgQueue;      // 任务队列
    Semaphore                           _queueSem;      // 队列信号
    ThreadMutexLock                     _queueMutex;    // 队列锁

    bool                                _isENSystem;    // 系统语言是否为英文
};

#endif
