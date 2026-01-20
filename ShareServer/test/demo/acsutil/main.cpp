#include <abprec.h>
#include <ncutil/ncutil.h>

#include "ncACSUtil.h"

void initDlls (void)
{
    static AppContext appCtx (_T("acs_util"));
    AppContext::setInstance (&appCtx);

    AppSettings* appSettings = AppSettings::getCFLAppSettings ();
    LibManager::getInstance ()->initLibs (appSettings, &appCtx, 0);

    ::ncInitXPCOM ();
}

int main (int argc, char** argv)
{
    initDlls ();

    ncACSUtil debugger;

    if (argc != 2) {
        debugger.getUsage ();
        return -1;
    }

    try {
        debugger.Execute (argv[1]);
    }
    catch (Exception& e) {
        printMessage2 (_T("Exception: %s"), e.toString ().getCStr ());
        debugger.getUsage ();
    }

    return 0;
}
