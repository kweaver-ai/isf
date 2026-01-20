#include <abprec.h>

#include "acsprocessor.h"
#include "ncActiveRecordThread.h"

// json 解析，保证线程安全
#define BOOST_SPIRIT_THREADSAFE
#include <boost/property_tree/ptree.hpp>
#include <boost/property_tree/ini_parser.hpp>
using namespace boost::property_tree;

#ifdef _DEBUG
#define RECORD_INTERVAL        30000       // Debug版记录活跃用户的时间间隔：30秒
#else
#define RECORD_INTERVAL        600000      // Release版记录活跃用户的时间间隔：10分钟
#endif

ncActiveRecordThread::ncActiveRecordThread ()
{
    NC_ACS_PROCESSOR_TRACE (_T("this: %p"), this);

    nsresult ret;
    _dbTokenManager = do_CreateInstance (NC_DB_TOKEN_MANAGER_CONTRACTID, &ret);
    if (NS_FAILED (ret)) {
        THROW_E (ACS_PROCESSOR, FAILED_TO_CREATE_DB_TOKEN_MANANGER,
            _T("Failed to create db token manager: 0x%x"), ret);
    }

    _sleepMilliSecond = RECORD_INTERVAL;

    // 记录活跃用户的时间间隔
    try {
        ptree pt;
        boost::property_tree::ini_parser::read_ini("/sysvol/conf/service_conf/app_default.conf", pt);

        try {
            _sleepMilliSecond = pt.get<int64>("ShareServer.active_record_interval");
        }
        catch (ptree_error& e) {
        }
    }
    catch (ptree_error& e) {
    }

    SystemLog::getInstance ()->log(__FILE__, __LINE__, INFORMATION_LOG_TYPE,
                                   "_sleepMilliSecond = %lld", _sleepMilliSecond);
}

ncActiveRecordThread::~ncActiveRecordThread ()
{
    NC_ACS_PROCESSOR_TRACE (_T("this: %p"), this);
}

void ncActiveRecordThread::pushActiveUserInfo (const String& userId, const String& time)
{
    NC_ACS_PROCESSOR_TRACE (_T("this: %p, userId: %s, time: %s"),
        this, userId.getCStr (), time.getCStr ());

    AutoLock<ThreadMutexLock> lock (&_mutex);
    _ActiveUserInfos[userId] = time;
}

void ncActiveRecordThread::run()
{
    map<String, String> tmpActiveUserInfos;

    while(1) {
        try {
            // 记录活跃用户信息
            tmpActiveUserInfos.clear ();
            {
                AutoLock<ThreadMutexLock> lock (&_mutex);
                tmpActiveUserInfos.swap (_ActiveUserInfos);
                _ActiveUserInfos.clear ();
            }

            if (!tmpActiveUserInfos.empty ()) {
                _dbTokenManager->SaveActiveUser(tmpActiveUserInfos);
            }
        }
        catch (Exception& e) {
            NC_ACS_PROCESSOR_TRACE (_T("this: %p, error: %s"),
                this, e.toString ().getCStr ());
        }
        catch (...) {
            NC_ACS_PROCESSOR_TRACE (_T("this: %p, error: unknown"), this);
        }

        // 记录后sleep
        sleep (_sleepMilliSecond);
    }
}
