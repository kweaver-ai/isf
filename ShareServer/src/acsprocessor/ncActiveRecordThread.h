# ifndef __NC_ACTIVE_RECORD_THREAD_H__
# define __NC_ACTIVE_RECORD_THREAD_H__

#include <acsdb/public/ncIDBTokenManager.h>

class ncActiveRecordThread : public Thread
{
public:
    ncActiveRecordThread ();
    ~ncActiveRecordThread ();

public:
    virtual void run ();
    void pushActiveUserInfo (const String& userId, const String& time);

private:
    map<String, String>             _ActiveUserInfos;            // 需要记录的活跃用户信息
    ThreadMutexLock                 _mutex;                      // 信息锁，支持多线程访问
    int64                           _sleepMilliSecond;           // 记录活跃用户的时间间隔

    // dbActiveUserManager
    nsCOMPtr<ncIDBTokenManager>        _dbTokenManager;
};

#endif // __NC_ACTIVE_RECORD_THREAD_H__
