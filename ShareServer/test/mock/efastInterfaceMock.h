/*************************************************************************************************
efastInterfaceMock.h:
    Copyright (c) Eisoo Software, Inc.(2021), All rights reserved.

Purpose:
    drivenadapters efastInterface mock

Author:
    Dylan.gao@aishu.cn

Creating Time:
    2021-07-26
***************************************************************************************************/
#ifndef __EFAST_MOCK_H
#define __EFAST_MOCK_H

#include <gmock/gmock.h>
#include <drivenadapter/public/efastInterface.h>

class efastMock: public efastInterface
{
    XPCOM_OBJECT_MOCK (efastMock);

public:
};

#endif // End __EFAST_MOCK_H
