/***************************************************************************************************
pluginMessage.h:
    Copyright (c) Eisoo Software, Inc.(2022), All rights reserved

Purpose:
    pluginMessage access control

Author:
    Dylan.gao@aishu.cn

Creating Time:
    2022-03-25
***************************************************************************************************/

#ifndef __PLUGIN_MESSAGE_H
#define __PLUGIN_MESSAGE_H

#include <abprec.h>

#include "public/pluginMessageInterface.h"
#include <ncMQClient.h>

/* Header file */
class pluginMessage : public pluginMessageInterface
{
    AB_DECLARE_THREADSAFE_SINGLETON (pluginMessage)

public:
    NS_DECL_ISUPPORTS
    NS_DECL_PLUGINMESSAGEINTERFACE

    pluginMessage();
    ~pluginMessage();

private:
    boost::shared_ptr<ncMQClient>   _mqClient;

};

#endif // __PLUGIN_MESSAGE_H
