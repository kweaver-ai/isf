/***************************************************************************************************
ncDBEDocManager.h:
    Copyright (c) Eisoo Software Inc. (2020), All rights reserved.

Purpose:
    NSQ manager

Author:
   Will.lv@aishu.cn

Creating Time:
    2021-03-15
***************************************************************************************************/
#ifndef __NSQ_H
#define __NSQ_H

#include "public/nsqInterface.h"
#include <dataapi/ncJson.h>
#include <ncMQClient.h>

// NSQ事件类
class ncNSQEvent
{
public:
    ncNSQEvent (ncNSQEventType type,
                String _nsqMsg)
        : eventType (type)
        , nsqMsg (_nsqMsg)
    {
    }

public:
    ncNSQEventType      eventType;        // NSQ事件的类型
    String              nsqMsg;           // 发送至NSQ的消息
};

typedef std::shared_ptr<ncNSQEvent> ncNSQEventSPtr;

/* Header file */
class nsq : public nsqInterface
{
    AB_DECLARE_THREADSAFE_SINGLETON (nsq)

public:
    NS_DECL_ISUPPORTS
    NS_DECL_NSQINTERFACE

    nsq ();
    ~nsq ();

private:
    String FormatNSQMsg (const ncNSQEventType& nsqType, const NSQMsg& nsqMsg);
    void Publish (const ncNSQEventSPtr eventSPtr);

    void   intToPermArray (int perm, JSON::Array& permArray);
    String docLibType2String(int docType);
    String applyType2String(const ncNSQApplyType applyType);
    String accessorType2String(int accessorType);
    String opType2String(int opType);
    String ticks2String(int64 ticks);
private:
    map<ncNSQEventType, String>           _eventTopicMap;
    boost::shared_ptr<ncMQClient>         _mqClient;
};

#endif // __NSQ_H
