/***************************************************************************************************
common.h:
    Copyright (c) Eisoo Software Inc. (2021), All rights reserved.

Purpose:
    驱动适配器公共头文件

Author:
    Sunshine.tang@aishu.cn

Creating Time:
    2021-02-24
***************************************************************************************************/
#ifndef __NC_COMMON_H__
#define __NC_COMMON_H__

///////////////////////////////////////////////////////////////////////////////////////////////////
// 公共类型

// 登录账号类型
enum class ncAccountType {
    OTHER   = 0,
    ID_CARD = 1,
};

// 设备类型
enum class ncClientType {
    UNKNOWN       = 0,
    IOS           = 1,
    ANDROID       = 2,
    WINDOWS_PHONE = 3,
    WINDOWS       = 4,
    MAC_OS        = 5,
    WEB           = 6,
    MOBILE_WEB    = 7,
    NAS           = 8,
    CONSOLE_WEB   = 9,
    DEPLOY_WEB    = 10,
    LINUX         = 11,
    APP           = 12,
};

// 访问者类型
enum class ncTokenVisitorType {
    REALNAME  = 1,       // 实名用户
    ANONYMOUS = 4,       // 匿名用户
    BUSINESS  = 6,       // 业务系统
};

#endif // __NC_COMMON_H__
