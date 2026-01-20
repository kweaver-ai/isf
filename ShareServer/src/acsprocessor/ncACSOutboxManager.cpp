/***************************************************************************************************
ncACSOutboxManager.cpp:
    Copyright (c) Eisoo Software Inc. (2021), All rights reserved.

Purpose:
    acs outbox manager 接口

Author:
    Dylan.gao@aishu.cn

Creating Time:
    2021-06-18
***************************************************************************************************/
#include <abprec.h>

#include "acsprocessor.h"
#include "ncACSOutboxManager.h"

#define RETRY_INTERVAL        30000       // 当推送outbox出错时，等待30秒重新尝试

/* Implementation file */
NS_IMPL_QUERY_INTERFACE1(ncACSOutboxManager, ncIACSOutboxManager)

// protected
NS_IMETHODIMP_(nsrefcnt) ncACSOutboxManager::AddRef (void)
{
    return 1;
}

// protected
NS_IMETHODIMP_(nsrefcnt) ncACSOutboxManager::Release (void)
{
    return 1;
}

AB_DEFINE_THREADSAFE_SINGLETON_NO_POOL (ncACSOutboxManager)

ncACSOutboxManager::ncACSOutboxManager()
    : _pushOutboxThread (nullptr),
      _needNotify (true)
{
    NC_ACS_PROCESSOR_TRACE (_T("this: %p"), this);

    nsresult ret;
    _dbOutboxManager = do_CreateInstance (NC_DB_OUTBOX_MANAGER_CONTRACTID, &ret);
    if (NS_FAILED (ret)) {
        THROW_E (ACS_PROCESSOR, FAILED_TO_CREATE_OSS_CLIENT,
            _T("Failed to create db outbox manager: 0x%x"), ret);
    }
}

ncACSOutboxManager::~ncACSOutboxManager()
{
    NC_ACS_PROCESSOR_TRACE (_T("this: %p"), this);
}

void ncACSOutboxManager::pushOutbox ()
{
    bool retry = false;
    while (1) {
        try {
            if (retry){
                _notiSem.wait (RETRY_INTERVAL);
            }else {
                _notiSem.wait ();
            }
            _needNotifyLock.lock ();
            _needNotify = true;
            _needNotifyLock.unlock ();
            while (1)
            {
                // 若outbox表全部处理完，则推出循环
                if (!_dbOutboxManager->PushOutboxInfo ()){
                    retry = false;
                    break;
                }
            }
        }
        catch (Exception& e) {
            String errMsg;
            errMsg.format (_T("ncPushOutboxThread::run error: %s"), e.toString ().getCStr ());
            SystemLog::getInstance ()->log (__FILE__, __LINE__, ERROR_LOG_TYPE, errMsg.getCStr ());
            retry = true;
        }
        catch (...) {
            String errMsg = _T("ncPushOutboxThread::run : unknown error");
            SystemLog::getInstance ()->log (__FILE__, __LINE__, ERROR_LOG_TYPE, errMsg.getCStr ());
            retry = true;
        }
    }
}

/* [notxpcom] void StartPushOutboxThread (); */
NS_IMETHODIMP_(void) ncACSOutboxManager::StartPushOutboxThread()
{
    NC_ACS_PROCESSOR_TRACE (_T("this: %p, begin"), this);

    static ThreadMutexLock pushOutboxMutex;
    AutoLock<ThreadMutexLock> lock (&pushOutboxMutex);

    if (_pushOutboxThread == nullptr) {
        _pushOutboxThread.reset (new std::thread (&ncACSOutboxManager::pushOutbox, this));
    }

    NC_ACS_PROCESSOR_TRACE (_T("this: %p, end"), this);
}

/* [notxpcom] void NotifyPushOutboxThread (); */
NS_IMETHODIMP_(void) ncACSOutboxManager::NotifyPushOutboxThread()
{
    NC_ACS_PROCESSOR_TRACE (_T("this: %p, begin"), this);
    _needNotifyLock.lock ();
    if (_pushOutboxThread != nullptr && _needNotify) {
        SystemLog::getInstance ()->log (__FILE__, __LINE__, INFORMATION_LOG_TYPE,_T("notify ncPushOutboxThread"));
        _notiSem.notify ();
        _needNotify = false;
    }
    _needNotifyLock.unlock ();
    NC_ACS_PROCESSOR_TRACE (_T("this: %p, end"), this);
}
