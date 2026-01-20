#ifndef __NC_ACS_COMMON_H__
#define __NC_ACS_COMMON_H__

///////////////////////////////////////////////////////////////////////////////////////////////////
// 公共类型

// 密级属性需要左移的位数
#define CSF_LEVEL_LEFT_SHIFT_NUM 24

// 权限值，锁定状态转为GNS对象的attr属性
enum ncDocAttr {
    ACS_ATTR_READ_ONLY          =       1,         // 只读(windows cifs文件属性)
    ACS_ATTR_LOCKED             =       2,         // 是否被锁定
    ACS_ATTR_ALLOW_DISPLAY      =       4,         // 是否允许显示
    ACS_ATTR_DENY_DISPLAY       =       8,         // 是否拒绝显示
    ACS_ATTR_ALLOW_PREVIEW      =       16,        // 是否允许预览
    ACS_ATTR_DENY_PREVIEW       =       32,         // 是否拒绝预览
    ACS_ATTR_ALLOW_READ         =       64,         // 是否允许下载
    ACS_ATTR_DENY_READ          =       128,         // 是否拒绝下载
    ACS_ATTR_ALLOW_CREATE       =       256,         // 是否允许新建
    ACS_ATTR_DENY_CREATE        =       512,         // 是否拒绝新建
    ACS_ATTR_ALLOW_EDIT         =       1024,         // 是否允许修改
    ACS_ATTR_DENY_EDIT          =       2048,         // 是否拒绝修改
    ACS_ATTR_ALLOW_DELETE       =       4096,         // 是否允许删除
    ACS_ATTR_DENY_DELETE        =       8192,         // 是否拒绝删除
};

// 访问者类型
enum class ncVisitorType {
    REALNAME  = 1,       // 实名用户
    ANONYMOUS = 4,       // 匿名用户
    BUSINESS  = 6,       // 业务系统
};

// 登录账号类型
enum class ACSAccountType {
    OTHER   = 0,
    ID_CARD = 1,
};

// 设备类型
enum class ACSClientType {
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

// 获取权限结果
struct ncAccessPerm {
    int                     allowValue;        // 允许的权限值
    int                     denyValue;         // 拒绝的权限值

    ncAccessPerm()
        : allowValue(0),
            denyValue(0)
    {
    }

    ncAccessPerm(int aValue, int dValue)
        : allowValue(aValue),
            denyValue(dValue)
    {
    }
};

#endif // __NC_ACS_COMMON_H__
