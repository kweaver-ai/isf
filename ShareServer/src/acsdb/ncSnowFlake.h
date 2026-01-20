/***************************************************************************************************
ncSnowFlake.h:
    Copyright (c) Eisoo Software, Inc.(2009 - 2013), All rights reserved

Purpose:
    snow flake unique id

Author:
    Young.yu@aishu.cn

Creating Time:
    2023-02-27
***************************************************************************************************/
#ifndef __NC_T_SNOW_FLAKE_H
#define __NC_T_SNOW_FLAKE_H

#include <abprec.h>

/*
 * SnowFlake
 */
class ncSnowFlake {
    AB_DECLARE_THREADSAFE_SINGLETON (ncSnowFlake)

public:
    ncSnowFlake();
    ~ncSnowFlake();

private:
    int m_nBitLenTime; //时间戳长度
    int m_nBitLenSequence; // 序列号长度
    int m_nBitLenMachineID; // 机器唯一码长度

    int64 m_nStartTime; // 开始时间戳 单位：10毫秒
	int64 m_nElapsedTime; // 间隔时间戳 单位：10毫秒
	int m_nSequence; // 序号
	int m_nMachineID; // 机器码

    ThreadMutexLock _sLock;
public:
    int64 NextID();
};

#endif // __NC_T_SNOW_FLAKE_H
