/***************************************************************************************************
ncTypeCommon.h:
    Copyright (c) AISHU Software Inc. (2022), All rights reserved.

Purpose:
    公共结构体

Author:
    Dylan.gao@aishu.cn

Creating Time:
    2022-4-23
***************************************************************************************************/
#ifndef __NC_TYPE_COMMON_H__
#define __NC_TYPE_COMMON_H__

// 消息接收者
struct messageReceiver {
    String   id;         // 接收者id
    String   account;    // 登录名
    String   name;       // 显示名
    String   email;      // 邮箱地址
    String   telephone;  // 电话号码
    String   thirdAttr;  // 第三方应用属性
    String   thirdId;    // 第三方应用id
};

// 消息
struct acsMessage {
    String              content;       // 消息内容
    vector<messageReceiver> receivers; // 接收者组
    String              msgId;         // 消息id
    String              channel;       // 消息类型
};

// 所有者信息
struct dbOwnerInfo {
    String                docId;        // 文档id
    String                ownerId;    // 所有者id
    String                ownerName;    // 所有者显示名
    int                   ownerType;    // 所有者类型
    int64                 modifyTime;   // 最后修改时间
};

#endif