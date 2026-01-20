#ifndef __PLATFORM_DATAENGINE_EFAST_COMMON_CONFIG_H_
#define __PLATFORM_DATAENGINE_EFAST_COMMON_CONFIG_H_
#include <abprec.h>
#include "nsISupports.h"
#include <boost/property_tree/ptree.hpp>
#include <boost/property_tree/ini_parser.hpp>
using namespace boost::property_tree;

class Config {
  private:
    Config() {
      try {
        ptree pt;
        boost::property_tree::ini_parser::read_ini("/sysvol/conf/service_conf/app_default.conf", pt);
        systemId = toCFLString(pt.get<std::string>("ShareServer.system_id", ""));
      }catch (ptree_error& e) {
        SystemLog::getInstance ()->log (__FILE__, __LINE__, ERROR_LOG_TYPE,
                                      _T("Error : read configuration information from file failed."));
      }
      SystemLog::getInstance ()->log (__FILE__, __LINE__, ERROR_LOG_TYPE,
                                      _T("read systemId success, systemId: %s"), systemId.getCStr());
    };
    ~Config() = default;

  public:
    Config(const Config &) = delete;
    Config &operator=(const Config &) = delete;

    /**
     * @brief 获取静态单例，c++11 线程安全
     */
    static Config &getInstance(){
      static Config helper;
      return helper;
    };

    String getSystemId() const{
      return systemId;
    };

  private:
    String systemId;
};

#endif // __PLATFORM_DATAENGINE_EFAST_COMMON_CONFIG_H_
