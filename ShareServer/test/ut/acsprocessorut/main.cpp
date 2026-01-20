#include <abprec.h>
#include <ncutil/ncutil.h>
#include <gtest/gtest.h>
#include <gmock/gmock.h>

void InitDlls ()
{
    // 加载 cfl 库
    static AppContext appCtx (_T("acs_processor_ut"));
    AppContext::setInstance (&appCtx);
    AppSettings* appSettings = AppSettings::getCFLAppSettings ();
    LibManager::getInstance ()->initLibs (appSettings, &appCtx, 0);

    // xpcom 核心库初始化
    ::ncInitXPCOM ();

    // 开启cfl的异常输出
    abEnableOutputError ();
}

int main(int argc, char** argv)
{
    try {
        // 初始化基础库
        InitDlls ();

        // 初始化 gtest 参数
        testing::InitGoogleTest(&argc, argv);

        // 初始化 gmock 参数
        testing::InitGoogleMock(&argc, argv);

        // 运行测试
        return RUN_ALL_TESTS ();
    }
    catch (Exception& e) {
        printMessage2 (_T("Test Error: %s"), e.toFullString ().getCStr ());
        return 1;
    }
    catch (...) {
        printMessage2 (_T("Test Error: Unknown."));
        return 1;
    }

    return 0;
}
