#ifndef __NC_REFRESH_TOKEN_THREAD_H__
#define __NC_REFRESH_TOKEN_THREAD_H__

#include <acssharemgnt/public/ncIACSShareMgnt.h>

/*
 * ncRefreshTokenThread
 */
class ncRefreshTokenThread : public Thread
{
public:
    ncRefreshTokenThread ();
    ~ncRefreshTokenThread ();

public:
    virtual void run ();
    void pushRefreshInfo (const String& userId, const String& lastRequestTime, const bool& bUpdateClientTime);

private:
    map<String, pair<String, bool>>        _refreshInfos;    // 需要刷新的token信息
    ThreadMutexLock            _mutex;    // token信息锁，支持多线程访问

    // acsShareMgnt
    nsCOMPtr<ncIACSShareMgnt>        _acsShareMgnt;
};

#endif // __NC_REFRESH_TOKEN_THREAD_H__
