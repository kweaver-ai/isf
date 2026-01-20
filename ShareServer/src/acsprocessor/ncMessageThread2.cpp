#include <abprec.h>
#include <ncutil/ncBusinessDate.h>
#include <dataapi/ncJson.h>

#include "acsprocessor.h"
#include "ncMessageThread2.h"
#include "ncMessageUtil.h"
#include "ncACSProcessorUtil.h"
#include <boost/date_time/posix_time/posix_time.hpp>
#include <boost/date_time/local_time_adjustor.hpp>
#include <boost/date_time/c_local_time_adjustor.hpp>

#ifdef _DEBUG
#define CLEAN_INTERVAL        60000       // Debug版清理的间间隔：1分钟
#else
#define CLEAN_INTERVAL        600000      // Release版清理的时间间隔：10分钟
#endif

ncMessageThread2::ncMessageThread2 ()
{
    NC_ACS_PROCESSOR_TRACE (_T("this: %p"), this);

    nsresult ret;
    _dbMessageManager = do_CreateInstance (NC_DB_MESSAGE_MANAGER_CONTRACTID, &ret);
    if (NS_FAILED (ret)) {
        THROW_E (ACS_PROCESSOR, FAILED_TO_CREATE_DB_MESSAGE_MANANGER,
            _T("Failed to create db message manager: 0x%x"), ret);
    }

    _acsShareMgnt = do_CreateInstance (NC_ACS_SHARE_MGNT_CONTRACTID, &ret);
    if (NS_FAILED (ret)) {
        THROW_E (ACS_PROCESSOR, FAILED_TO_CREATE_ACS_SHAREMGNT,
            _T("Failed to create acs sharemgnt: 0x%x"), ret);
    }
}

ncMessageThread2::~ncMessageThread2 ()
{
    NC_ACS_PROCESSOR_TRACE (_T("this: %p"), this);
}

void ncMessageThread2::AddMessage (const vector<shared_ptr<acsMessage>> & msgs)
{
    _queueMutex.lock ();
    int count = 0;
    for (size_t i = 0; i < msgs.size (); ++i)
    {
        count++;
        _msgQueue.push (std::move (msgs[i]));
    }
    _queueMutex.unlock ();
    if (count > 0)
    {
        _queueSem.notify (count);
    }
}

void ncMessageThread2::run ()
{
    while (1) {
        try {
            bool gotit = _queueSem.wait (CLEAN_INTERVAL);

            if (!gotit) {
                // 做清除操作
                //_dbMessageManager->DelMessage (BusinessDate::getCurrentTime () - 3 * 30 * Date::ticksPerDay);

                continue;
            }

            shared_ptr<acsMessage> msg;
            _queueMutex.lock ();
            msg = std::move (_msgQueue.front ());
            _msgQueue.pop ();
            _queueMutex.unlock ();

            vector<String> receiverIds;
            for (size_t i = 0; i < msg->receivers.size (); ++i){
                receiverIds.push_back (msg->receivers[i].id);
            }
            if (!receiverIds.empty ()) {
                ncMessageUtil::getInstance ()->RemoveDuplicateStrs (receiverIds);
                _dbMessageManager->AddMessage (msg->msgId, msg->content, receiverIds, "");
            }
        }
        catch (Exception& e) {
            String errMsg;
            errMsg.format (_T("ncMessageThread2::run error: %s"), e.toString ().getCStr ());


            SystemLog::getInstance ()->log (__FILE__, __LINE__, ERROR_LOG_TYPE, errMsg.getCStr ());
            printMessage2 (errMsg.getCStr ());
        }
        catch (...) {
            String errMsg = _T("ncMessageThread2::run : unknown error");

            SystemLog::getInstance ()->log (__FILE__, __LINE__, ERROR_LOG_TYPE, errMsg.getCStr ());
            printMessage2 (errMsg.getCStr ());
        }
    }
}
