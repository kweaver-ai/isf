/****************************************************************************************************
util.h
  fnv header file.
  Copyright (c) Google Software, Inc.(2024), All rights reserved.

Purpose:
  Header file to implement util.

Author:
  Darren (Darren.su@aishu.cn)

Creating Time:
  2024-12-09
****************************************************************************************************/
#ifndef __PLATFORM_DATAENGINE_EFAST_COMMON_UTIL_H_
#define __PLATFORM_DATAENGINE_EFAST_COMMON_UTIL_H_
#include "common/config.hpp"
#include <abprec.h>
#include "nsISupports.h"

class Util {
  public:
    static String getDBName(String dbName){
      return Config::getInstance().getSystemId() + dbName;
    };
};

#endif // __PLATFORM_DATAENGINE_EFAST_COMMON_UTIL_H_
