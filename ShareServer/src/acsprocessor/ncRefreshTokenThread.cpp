#include <abprec.h>

#include "acsprocessor.h"
#include "ncRefreshTokenThread.h"

#define REFRESH_INTERVAL       10000        // 10秒钟，刷新token的间隔时间
#define NUM_PERM_COMMIT        10000        // 一次提交10000条
#define RETRY_INTERVAL         1000        // 刷新失败时，过1秒后再尝试

ncRefreshTokenThread::ncRefreshTokenThread ()
{
    NC_ACS_PROCESSOR_TRACE (_T("this: %p"), this);

    nsresult ret;
    _acsShareMgnt = do_CreateInstance (NC_ACS_SHARE_MGNT_CONTRACTID, &ret);
    if (NS_FAILED (ret)) {
        THROW_E (ACS_PROCESSOR, FAILED_TO_CREATE_ACS_SHAREMGNT,
            _T("Faild to create acs sharemgnt: 0x%x"), ret);
    }
}

ncRefreshTokenThread::~ncRefreshTokenThread ()
{
    NC_ACS_PROCESSOR_TRACE (_T("this: %p"), this);
}

void ncRefreshTokenThread::pushRefreshInfo (const String& userId, const String& lastRequestTime, const bool& bUpdateClientTime)
{
    NC_ACS_PROCESSOR_TRACE (_T("this: %p, userId: %s, lastRequestTime: %s"),
        this, userId.getCStr (), lastRequestTime.getCStr ());

    AutoLock<ThreadMutexLock> lock (&_mutex);
    _refreshInfos[userId] = pair<String, bool>(lastRequestTime, bUpdateClientTime);
}

void ncRefreshTokenThread::run ()
{
    while (1)
    {
        // 一次刷新token过程
        while (1)
        {
            // 获取最多1W条刷新信息
            int count = 0;
            vector<ncACSRefreshInfo> tmpInfos;
            ncACSRefreshInfo tmp;
            {
                AutoLock<ThreadMutexLock> lock (&_mutex);
                map<String, pair<String, bool>>::iterator iter = _refreshInfos.begin ();
                for (; iter != _refreshInfos.end(); ++iter) {

                    tmp.userId = iter->first;
                    tie(tmp.lastRequestTime, tmp.bUpdateClientTime) = iter->second;
                    tmpInfos.push_back (tmp);

                    ++count;
                    if (count == NUM_PERM_COMMIT) {
                        break;
                    }
                }
            }

            if (count != 0)
            {
                // 刷新这些tokenid
                bool suc = false;
                try {
                    _acsShareMgnt->BatchUpdateUserLastRequestTime (tmpInfos);
                    suc = true;
                }
                catch (Exception& e) {
                    NC_ACS_PROCESSOR_TRACE (_T("this: %p, error: %s"),
                        this, e.toString ().getCStr ());
                }
                catch (...) {
                    NC_ACS_PROCESSOR_TRACE (_T("this: %p, error: unknown"),
                        this);
                }

                // 移除掉这些tokenid
                if (suc) {
                    {
                        AutoLock<ThreadMutexLock> lock (&_mutex);
                        for (size_t i = 0; i < tmpInfos.size (); ++i) {
                            _refreshInfos.erase (tmpInfos[i].userId);
                        }
                    }

                    NC_ACS_PROCESSOR_TRACE (_T("this: %p, success"),
                        this);
                }
                else {
                    // 如果刷新失败，sleep后再行尝试
                    sleep (RETRY_INTERVAL);
                    continue;
                }
            }

            // 没有需要刷新的数据了
            if (count < NUM_PERM_COMMIT) {
                break;
            }
        }

        // 刷新后sleep
        sleep (REFRESH_INTERVAL);
    }
}
