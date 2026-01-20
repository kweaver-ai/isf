#include <abprec.h>
#include "acsprocessor.h"
#include "ncPushMessageThread2.h"
#include "ncACSProcessorUtil.h"
#include "ncMessageUtil.h"

// json 解析，保证线程安全
#define BOOST_SPIRIT_THREADSAFE
#include <boost/property_tree/ptree.hpp>
#include <boost/property_tree/ini_parser.hpp>
using namespace boost::property_tree;

#define MESSAGE_QUEUE_MAX_SIZE            100000

ncPushMessageThread2::ncPushMessageThread2 ()
{
    NC_ACS_PROCESSOR_TRACE (_T("this: %p"), this);

    // 获取系统语言
    _isENSystem = false;
    try {
        ptree pt;
        boost::property_tree::ini_parser::read_ini("/sysvol/conf/service_conf/language.conf", pt);

        try {
            string language = pt.get<string>("LANG");
            _isENSystem = language.compare ("en_US") == 0;
        }
        catch (ptree_error& e) {
        }
    }
    catch (ptree_error& e) {
    }

    nsresult ret;

    _acsShareMgnt = do_CreateInstance (NC_ACS_SHARE_MGNT_CONTRACTID, &ret);
    if (NS_FAILED (ret)) {
        THROW_E (ACS_PROCESSOR, FAILED_TO_CREATE_ACS_SHAREMGNT,
            _T("Failed to create acs sharemgnt: 0x%x"), ret);
    }

    _dbConfManager = do_CreateInstance (NC_DB_CONF_MANAGER_CONTRACTID, &ret);
    if (NS_FAILED (ret)) {
        THROW_E (ACS_PROCESSOR, FAILED_TO_CREATE_DB_CONF_MANANGER,
            _T("Failed to create db conf manager: 0x%x"), ret);
    }

    _acsConfManager = do_CreateInstance(NC_ACS_CONF_MANAGER_CONTRACTID, &ret);
    if (NS_FAILED (ret)) {
        THROW_E (ACS_PROCESSOR, FAILED_TO_CREATE_ACS_CONF_MANANGER,
            _T("Failed to create acs conf manager: 0x%x"), ret);
    }

    _pluginMessage = do_CreateInstance (PLUGIN_MESSAGE_CONTRACTID, &ret);
    if (NS_FAILED (ret)) {
        THROW_E (ACS_PROCESSOR, FAILED_TO_CREATE_DRIVENADAPTER_MANANGER,
            _T("Failed to create pluginMessage instance: 0x%x"), ret);
    }
}

ncPushMessageThread2::~ncPushMessageThread2 ()
{
    NC_ACS_PROCESSOR_TRACE (_T("this: %p"), this);
}

void ncPushMessageThread2::AddMessage(const vector<shared_ptr<acsMessage>> & msgs)
{
    _queueMutex.lock ();

    // 检查消息数量，过多时直接丢弃
    if (_msgQueue.size () > MESSAGE_QUEUE_MAX_SIZE) {
        _queueMutex.unlock ();
        return;
    }

    for (size_t i = 0; i < msgs.size (); ++i) {
        _msgQueue.push (std::move (msgs[i]));
    }

    _queueMutex.unlock ();
    _queueSem.notify (msgs.size ());
}

void ncPushMessageThread2::run ()
{
    while (1) {
        try {
            _queueSem.wait ();

            shared_ptr<acsMessage> msg;
            _queueMutex.lock ();
            msg = std::move (_msgQueue.front ());
            _msgQueue.pop ();
            _queueMutex.unlock ();

            // 发送插件消息
            _pluginMessage->SendPluginMessage (*msg);
        }
        catch (Exception& e) {
            String errMsg;
            errMsg.format (_T("ncPushMessageThread2::run error: %s"), e.toString ().getCStr ());


            SystemLog::getInstance ()->log (__FILE__, __LINE__, ERROR_LOG_TYPE, errMsg.getCStr ());
            printMessage2 (errMsg.getCStr ());
        }
        catch (...) {
            String errMsg = _T("ncPushMessageThread2::run : unknown error");

            SystemLog::getInstance ()->log (__FILE__, __LINE__, ERROR_LOG_TYPE, errMsg.getCStr ());
            printMessage2 (errMsg.getCStr ());
        }
    }
}
