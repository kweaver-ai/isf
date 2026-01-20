/***************************************************************************************************
ncSnowFlake.cpp:
    Copyright (c) Eisoo Software, Inc.(2009 - 2013), All rights reserved

Purpose:
    snow flake unique id

Author:
    Young.yu@aishu.cn

Creating Time:
    2023-02-27
***************************************************************************************************/
#include "ncSnowFlake.h"
#include "acsdb.h"
#include <arpa/inet.h>
#include <ncutil/ncBusinessDate.h>

AB_DEFINE_THREADSAFE_SINGLETON_NO_POOL (ncSnowFlake);

ncSnowFlake::ncSnowFlake()
:m_nBitLenTime(39)
,m_nBitLenSequence(8)
,m_nBitLenMachineID(16)
{
    NC_ACS_DB_TRACE (_T("this: %p begin"), this);

    // 时间戳以10毫秒为单位, 默认事件2014年9月1日 参考sonyflake
    m_nStartTime = 140950080000;

    // 获取最大序号
    m_nSequence = (1<<m_nBitLenSequence - 1);

    // 默认为0
    m_nElapsedTime = 0;

    // 根据pod的ip获取machineid
    string ip = toSTLString (Environment::getEnvVariable ("PORT_IP"));
    int64 result = 0;
     if (ip.find(_T(":")) == String::NO_POSITION) {
        struct in_addr addr;
        inet_pton(AF_INET, ip.c_str(), &addr);
        result = ntohl(addr.s_addr);
    } else {
        struct in6_addr addr;
        inet_pton(AF_INET6, ip.c_str(), &addr);
        result = ntohl(addr.s6_addr32[3]);
    }
    m_nMachineID = 0xFFFF & result;

    NC_ACS_DB_TRACE (_T("this: %p end"), this);
}

ncSnowFlake::~ncSnowFlake()
{
    NC_ACS_DB_TRACE (_T("this: %p"), this);
}

int64 ncSnowFlake::NextID()
{
    NC_ACS_DB_TRACE (_T("this: %p begin"), this);
    const int maskSequence = (1<<m_nBitLenSequence) - 1;

    AutoLock<ThreadMutexLock> lock (&_sLock);

    // 获取序号
    int64 current = BusinessDate::getCurrentTime () / 10000 - m_nStartTime;
    if (m_nElapsedTime < current) {
        m_nElapsedTime = current;
        m_nSequence = 0;
    } else {
        m_nSequence = (m_nSequence + 1) & maskSequence;
        if (m_nSequence == 0) {
            m_nElapsedTime++;
            int64 nOverTime = m_nElapsedTime - current;
            Thread::sleep(nOverTime * 10 - (BusinessDate::getCurrentTime () / 1000) % 10);
        }
    }

    // 判断是否超过时间范围
    int64 maxTime = (int64(1)<<m_nBitLenTime);
    if (m_nElapsedTime > maxTime) {
        return 0;
    }

    int64 nID = (int64(m_nElapsedTime)<<(m_nBitLenSequence+m_nBitLenMachineID)) |
		(int64(m_nSequence)<<m_nBitLenMachineID) |
		(int64(m_nMachineID));
    NC_ACS_DB_TRACE (_T("this: %p end"), this);
    return nID;
}
