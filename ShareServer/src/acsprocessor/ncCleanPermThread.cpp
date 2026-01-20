#include <abprec.h>
#include <ncutil/ncBusinessDate.h>

#define BOOST_SPIRIT_THREADSAFE
#include <boost/property_tree/ptree.hpp>
#include <boost/property_tree/ini_parser.hpp>
using namespace boost::property_tree;

#include "acsprocessor.h"
#include "ncCleanPermThread.h"
#include "ncACSProcessorUtil.h"

#ifdef _DEBUG
#define CLEAN_INTERVAL        60000       // Debug版清理权限的时间间隔：1分钟
#else
#define CLEAN_INTERVAL        600000      // Release版清理权限信息的时间间隔：10分钟
#endif

ncCleanPermThread::ncCleanPermThread ()
{
    NC_ACS_PROCESSOR_TRACE (_T("this: %p"), this);

    nsresult ret;
    _dbPermManager = do_CreateInstance (NC_DB_PERM_MANAGER_CONTRACTID, &ret);
    if (NS_FAILED (ret)) {
        THROW_E (ACS_PROCESSOR, FAILED_TO_CREATE_DB_PERM_MANANGER,
            _T("Failed to create db perm manager: 0x%x"), ret);
    }

    _cleanInterval = CLEAN_INTERVAL;
    try {
        ptree pt;
        boost::property_tree::ini_parser::read_ini("/sysvol/conf/service_conf/app_default.conf", pt);

        try {
            int64 cleanInterval;
            cleanInterval = pt.get<int64>("ShareServer.clean_interval");
            if (cleanInterval > 0) {
                _cleanInterval = cleanInterval;
            }
            SystemLog::getInstance ()->log(__FILE__, __LINE__, INFORMATION_LOG_TYPE,
                                "_cleanInterval = %lld", _cleanInterval);
        }
        catch (ptree_error& e) {
        }
    }
    catch (ptree_error& e) {
    }

}

ncCleanPermThread::~ncCleanPermThread ()
{
    NC_ACS_PROCESSOR_TRACE (_T("this: %p"), this);
}

void ncCleanPermThread::run ()
{
    while (1) {
        try {
            //
            // 清理过期的权限信息，并更新otag。
            //
            vector<dbCustomPermInfo> infos;
            _dbPermManager->GetExpirePermInfos (BusinessDate::getCurrentTime (), infos);

            for (size_t i = 0; i < infos.size (); ++i) {
                _dbPermManager->DeleteCustomPerm(infos[i].id);
                // 发送权限变更NSQ
                ncACSProcessorUtil::getInstance ()->SendPermChangeNSQ (infos[i].docId);
            }
        }
        catch (Exception& e) {
            String errMsg;
            errMsg.format (_T("ncCleanPermThread::run error: %s"), e.toString ().getCStr ());

            SystemLog::getInstance ()->log (__FILE__, __LINE__, ERROR_LOG_TYPE, errMsg.getCStr ());
            printMessage2 (errMsg.getCStr ());
        }
        catch (...) {
            String errMsg = _T("ncCleanPermThread::run : unknown error");

            SystemLog::getInstance ()->log (__FILE__, __LINE__, ERROR_LOG_TYPE, errMsg.getCStr ());
            printMessage2 (errMsg.getCStr ());
        }

        sleep (_cleanInterval);
    }
}
