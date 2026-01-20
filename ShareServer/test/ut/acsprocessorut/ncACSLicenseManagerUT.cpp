#include <abprec.h>
#include <gmock/gmock.h>

#include <acsprocessor/ncACSLicenseManager.h>

#include "ncACSLicenseManagerUT.h"

using namespace testing;

ncACSLicenseManagerUT::ncACSLicenseManagerUT ()
{
}

ncACSLicenseManagerUT::~ncACSLicenseManagerUT ()
{
}

void ncACSLicenseManagerUT::SetUp ()
{
}

void ncACSLicenseManagerUT::TearDown ()
{
}

TEST_F (ncACSLicenseManagerUT, GetLicenseInfo)
{
    nsresult ret;
    nsCOMPtr<ncIACSLicenseManager> acsLicenseManager = do_CreateInstance (NC_ACS_LICENSE_MANAGER_CONTRACTID, &ret);

    if (NS_FAILED (ret)) {
        printMessage2 (_T("Failed to create acs license manager: 0x%x"), ret);
    }

    map<String, String> licenseMaps;
    // AnyShare Enterprise
    // S6020主模块
    licenseMaps["S620B-7WLC8-JNSCX-HMJ2C-P32LS-9UK8L"] = "@anyshare/S6020/base/1/0/-1/-1";

    // RX4020-T1主模块
    licenseMaps["RX4T1-5EKNA-VZN4M-20LQC-N2ESD-QFHC6"] = "@anyshare/RX4020-T1/base/1/0/-1/-1";

    // RX4020-T2主模块
    licenseMaps["RX4T2-19DDD-NHRSC-RVT0D-K0K67-LUU6F"] = "@anyshare/RX4020-T2/base/1/0/-1/-1";

    // ASE-S主模块
    licenseMaps["ASSBA-B7J9J-L8PPU-U7WF8-KY6KP-02SX2"] = "@anyshare/ASE-S/base/1/0/-1/-1";

    // RX4020-T1节点模块
    licenseMaps["RX4NA-8DKVT-T7X57-4HPLL-GSZLM-21X5Q"] = "@anyshare/RX4020-T1-Node/node-agent/1";

    // RX4020-T2节点模块
    licenseMaps["R4N2A-8GY2L-94MUN-YZC4S-KVWUW-TRFYD"] = "@anyshare/RX4020-T2-Node/node-agent/1";

    // ASE-S节点模块
    licenseMaps["ASESN-EDUXL-939G3-UD4Q4-4TSUG-AWDLB"] = "@anyshare/ASE-S-Node/node-agent/1";

    // S6020升级模块
    licenseMaps["S62UG-2WYVB-RK46W-R9RV9-TF68S-R8RMH"] = "@anyshare/S6020-Upgrade/upgrade-agent";

    // S-100用户授权包
    licenseMaps["ASSUH-5YBR7-ZQ2ML-9QU28-PXKMP-Q3707"] = "@anyshare/S-100-User/user-agent/100";

    // S-1000用户授权包
    licenseMaps["ASSAH-7Q4HL-T9M6K-5JKBM-VASHU-BKSFV"] = "@anyshare/S-1000-User/user-agent/1000";

    // RX-100用户授权包
    licenseMaps["RXUAH-6J4QY-VTFWQ-C7U03-5EAQQ-TQPNE"] = "@anyshare/RX-100-User/user-agent/100";

    // RX-1000用户授权包
    licenseMaps["RXUAS-FHVZ1-26H0J-6WTQE-94M6X-A2JX9"] = "@anyshare/RX-1000-User/user-agent/1000";

    // RX-10000用户授权包
    licenseMaps["RXUAM-AYMPV-LCSF1-3RD0A-PE2X1-K5B14"] = "@anyshare/RX-10000-User/user-agent/10000";

    // RX站点用户授权包
    licenseMaps["RXSTU-64621-GT41V-0DJES-00KEE-U2A7M"] = "@anyshare/RX-Site-User/user-agent/-1";

    // ASE-100用户授权包
    licenseMaps["ASEUH-676DL-3HJKR-9B2D2-W2BBK-P53LA"] = "@anyshare/ASE-100-User/user-agent/100";

    // ASE-1000用户授权包
    licenseMaps["ASEUS-F1XA8-ZCA0C-2467Q-7Z2XR-MEUBK"] = "@anyshare/ASE-1000-User/user-agent/1000";

    // ASE-10000用户授权包
    licenseMaps["ASEUM-8ATRF-RH197-8HCVB-K2Z5F-JLJ8X"] = "@anyshare/ASE-10000-User/user-agent/10000";

    // ASE-站点用户授权包
    licenseMaps["ASUAG-1HWEY-EPPHZ-QZYVX-926GD-EKPMD"] = "@anyshare/ASE-Site-User/user-agent/-1";

    // EDMS
    // S6020-ND1主模块
    licenseMaps["A6NDB-0N2P9-T3KWX-D31DP-UXZX6-FVT3N"] = "@anyshare/S6020-ND1/base/1/0/-1/-1";

    // RX4020-ND1主模块
    licenseMaps["A6RXB-2VKK7-8HGPC-090FL-5GCYS-6PADQ"] = "@anyshare/RX4020-ND1/base/1/0/-1/-1";

    // NDE-S主模块
    licenseMaps["A6NSB-B8R0W-JVQUN-B8AE9-DGGLE-EM57A"] = "@anyshare/NDE-S/base/1/0/-1/-1";

    // RX4020-ND1节点模块
    licenseMaps["A6XNA-BYHJA-KVGKV-112YC-LKC0Z-R1VTB"] = "@anyshare/RX4020-ND1-Node/node-agent/1";

    // NDE-S节点模块
    licenseMaps["A6ESN-56DL4-JXX9P-Q1PX0-8TAM9-YG2C0"] = "@anyshare/NDE-S-Node/node-agent/1";

    // S6020-ND1升级模块
    licenseMaps["A6SUE-BULQ6-GR4Y2-QZ777-4NUT5-FZTMG"] = "@anyshare/S6020-ND1-Upgrade/upgrade-agent";

    // NDS-100用户授权包
    licenseMaps["NDSHA-0QEWQ-6QAPF-F2PZ3-4QQN4-2KA5B"] = "@anyshare/NDS-100-User/user-agent/100";

    // NDS-1000用户授权包
    licenseMaps["ADS1T-QLCZ9-HHTBN-0EYLM-8AYAE-Z3JKQ"] = "@anyshare/NDS-1000-User/user-agent/1000";

    // NDR-100用户授权包
    licenseMaps["NDRHA-0TAU5-1P5L7-G88E7-PMBCB-J5UYU"] = "@anyshare/NDR-100-User/user-agent/100";

    // NDR-1000用户授权包
    licenseMaps["NDRTA-DF7UN-MTC5Z-92JH8-HGFVP-KNKVX"] = "@anyshare/NDR-1000-User/user-agent/1000";

    // NDR-10000用户授权包
    licenseMaps["NDRMA-CZ2N8-G7LJD-L3QVH-VQYXA-635UJ"] = "@anyshare/NDR-10000-User/user-agent/10000";

    // NDR用户场地授权包
    licenseMaps["NDSUA-0S1T7-44SAG-ZSJF4-L3884-DC4PN"] = "@anyshare/NDR-Site-User/user-agent/-1";

    // NDE-100用户授权包
    licenseMaps["NEHUA-6237P-23C8M-JQ7DU-4BGCP-JTR4X"] = "@anyshare/NDE-100-User/user-agent/100";

    // NDE-1000用户授权包
    licenseMaps["NETUA-21CBW-AQN21-A8WZG-JVQJ0-4HGQL"] = "@anyshare/NDE-1000-User/user-agent/1000";

    // NDE-10000用户授权包
    licenseMaps["NEMUA-9KK93-LWD96-545KX-W9MDQ-KUW50"] = "@anyshare/NDE-10000-User/user-agent/10000";

    // NDE用户场地授权包
    licenseMaps["NESUA-4L74G-2JULT-S9GL2-X1UCG-H0WNP"] = "@anyshare/NDE-Site-User/user-agent/-1";

    // ASS
    // EX6020主模块
    licenseMaps["AEX62-8SETD-PNSA4-S4LGY-P4S05-DYFCM"] = "@anyshare/EX6020/base/1/0/-1/-1";

    // ASS-S主模块
    licenseMaps["ASSSB-ZD41Z-51JL3-AUSF4-K3YA6-Q21R4"] = "@anyshare/ASS-S/base/1/0/-1/-1";

    // EX-10用户授权包
    licenseMaps["AX10A-33KVJ-VVX7D-Y3KPS-2G29Q-2DGLR"] = "@anyshare/EX-10-User/user-agent/10";

    // EX-100用户授权包
    licenseMaps["AX100-WTCEX-VJW42-RFV5Y-BFC3M-FPF28"] = "@anyshare/EX-100-User/user-agent/100";

    // ASS-10用户授权包
    licenseMaps["AS10A-VRPC3-EB32T-0BX4F-51UHQ-BM610"] = "@anyshare/ASS-10-User/user-agent/10";

    // ASS-100用户授权包
    licenseMaps["AS100-Z26AA-350VK-KV01W-0WKHC-JZYQK"] = "@anyshare/ASS-100-User/user-agent/100";

    // 测试授权码
    // AnyShare 30天测试授权码
    licenseMaps["ASTT3-29FLS-F0AU8-JERDQ-6UMAN-CP3MV"] = "@anyshare/AS-30-Test/test/-1/-1/-1/30";

    // AnyShare 90天测试授权码
    licenseMaps["A7TT9-0XH78-LFPQN-VUWV3-2W0YZ-KTV33"] = "@anyshare/AS-90-Test/test/-1/-1/-1/90";

    // AnyShare 180天测试授权码
    licenseMaps["A7TT8-50ESA-GWJU4-WCYD8-DWH5D-4Y572"] = "@anyshare/AS-180-Test/test/-1/-1/-1/180";

    // AnyShare 360天测试授权码
    licenseMaps["A7TTY-2XD3C-LNQ0X-5AFDT-KS09U-AV9ML"] = "@anyshare/AS-360-Test/test/-1/-1/-1/360";

    // AnyShare大版本升级模块
    licenseMaps["A6UPB-8LW7N-09XZL-E9J1V-HS8WG-UJKY2"] = "@anyshare/AS5-Upgrade/base/0/0/-1/-1";

    // AnySare在线表格
    //AnyShare在线表格主模块
    licenseMaps["ECBSE-069LD-ESE9B-ZJD8P-DAEYX-0HJHZ"] = "@anyshare/AS-Excel/excel-base/0/-1";

    //在线表格50用户授权包
    licenseMaps["EC5UR-3FREH-TYXLZ-R689Z-CU0EK-ZSW0L"] = "@anyshare/Excel-50-User/excel-user-agent/50";

    //在线表格100用户授权包
    licenseMaps["ECHUR-479MP-KCX6T-ESQ2G-RG60G-90QKR"] = "@anyshare/Excel-100-User/excel-user-agent/100";

    //在线表格200用户授权包
    licenseMaps["ECTUR-15SED-REZZA-2B1RZ-2GZJJ-9NKG9"] = "@anyshare/Excel-200-User/excel-user-agent/200";

    //在线表格500用户授权包
    licenseMaps["ECFUR-7ARUQ-TSK87-TJYYR-2K8B5-XWLNH"] = "@anyshare/Excel-500-User/excel-user-agent/500";

    //在线表格1000用户授权包
    licenseMaps["ECSUR-0T7D2-22VAP-4JTX3-M8K2D-6ZYLK"] = "@anyshare/Excel-1000-User/excel-user-agent/1000";

    //在线表格2000用户授权包
    licenseMaps["ECTTU-EVSJY-K4GV9-6N9BC-3DG8R-YCJRK"] = "@anyshare/Excel-2000-User/excel-user-agent/2000";

    //在线表格5000用户授权包
    licenseMaps["ECFTU-32Q23-JCKXE-DL87A-HLZR9-YNSPD"] = "@anyshare/Excel-5000-User/excel-user-agent/5000";

    //在线表格10000用户授权包
    licenseMaps["ECFMU-A56F0-BMLGM-QRV24-3HANJ-4UV30"] = "@anyshare/Excel-10000-User/excel-user-agent/10000";

    //在线表格场地用户授权包
    licenseMaps["ECSTR-1KCZV-44WQA-GF7HC-PPGPD-ZLJN6"] = "@anyshare/Excel-Site-User/excel-user-agent/-1";


    for(map<String, String>::iterator iter = licenseMaps.begin(); iter != licenseMaps.end(); ++iter) {
        String license = iter->first;
        String parseInfo;
        acsLicenseManager->GetLicenseInfo(license, parseInfo);
        printMessage2(_T("%s->%s"), license.getCStr(), parseInfo.getCStr());

        ASSERT_EQ(parseInfo, iter->second);
    }
}

TEST_F (ncACSLicenseManagerUT, VerifyActiveCode)
{
    // 公司许可证授权切换，新的授权码可以激活
    nsresult ret;
    nsCOMPtr<ncIACSLicenseManager> acsLicenseManager = do_CreateInstance (NC_ACS_LICENSE_MANAGER_CONTRACTID, &ret);

    if (NS_FAILED (ret)) {
        printMessage2 (_T("Failed to create acs license manager: 0x%x"), ret);
    }

    // AnyShare Enterprise
    // S6020主模块
    ASSERT_EQ(acsLicenseManager->VerifyActiveCode("S620B-7WLC8-JNSCX-HMJ2C-P32LS-9UK8L", "005056AF7572", "0WDK3KGRF370"), 0);
    // RX4020-T1主模块
    ASSERT_EQ(acsLicenseManager->VerifyActiveCode("RX4T1-5EKNA-VZN4M-20LQC-N2ESD-QFHC6", "005056AF7572", "08PYD6ZVTD7Z"), 0);
    // RX4020-T2主模块
    ASSERT_EQ(acsLicenseManager->VerifyActiveCode("RX4T2-19DDD-NHRSC-RVT0D-K0K67-LUU6F", "005056AF7572", "0ZBS7ZA1NRVC"), 0);
    // ASE-S主模块
    ASSERT_EQ(acsLicenseManager->VerifyActiveCode("ASSBA-B7J9J-L8PPU-U7WF8-KY6KP-02SX2", "005056AF7572", "0W5QJME3HJ3Q"), 0);
    // RX4020-T1节点模块
    ASSERT_EQ(acsLicenseManager->VerifyActiveCode("RX4NA-8DKVT-T7X57-4HPLL-GSZLM-21X5Q", "005056AF7572", "0SLWLA6XD5RY"), 0);
    // RX4020-T2节点模块
    ASSERT_EQ(acsLicenseManager->VerifyActiveCode("R4N2A-8GY2L-94MUN-YZC4S-KVWUW-TRFYD", "005056AF7572", "0SVW1821J3B2"), 0);
    // ASE-S节点模块
    ASSERT_EQ(acsLicenseManager->VerifyActiveCode("ASESN-EDUXL-939G3-UD4Q4-4TSUG-AWDLB", "005056AF7572", "009WDW6BTFH2"), 0);
    // S6020升级模块
    ASSERT_EQ(acsLicenseManager->VerifyActiveCode("S62UG-2WYVB-RK46W-R9RV9-TF68S-R8RMH", "005056AF7572", "06LG3S8JFPHQ"), 0);
    // S-100用户授权包
    ASSERT_EQ(acsLicenseManager->VerifyActiveCode("ASSUH-5YBR7-ZQ2ML-9QU28-PXKMP-Q3707", "005056AF7572", "0KFYXEZJNDJ8"), 0);
    // S-1000用户授权包
    ASSERT_EQ(acsLicenseManager->VerifyActiveCode("ASSAH-7Q4HL-T9M6K-5JKBM-VASHU-BKSFV", "005056AF7572", "0YT0JCW3D9NA"), 0);
    // RX-100用户授权包
    ASSERT_EQ(acsLicenseManager->VerifyActiveCode("RXUAH-6J4QY-VTFWQ-C7U03-5EAQQ-TQPNE", "005056AF7572", "0EF4PUY15VLQ"), 0);
    // RX-1000用户授权包
    ASSERT_EQ(acsLicenseManager->VerifyActiveCode("RXUAS-FHVZ1-26H0J-6WTQE-94M6X-A2JX9", "005056AF7572", "025254SJJF56"), 0);
    // RX-10000用户授权包
    ASSERT_EQ(acsLicenseManager->VerifyActiveCode("RXUAM-AYMPV-LCSF1-3RD0A-PE2X1-K5B14", "005056AF7572", "0UHQTASJXLRC"), 0);
    // RX站点用户授权包
    ASSERT_EQ(acsLicenseManager->VerifyActiveCode("RXSTU-64621-GT41V-0DJES-00KEE-U2A7M", "005056AF7572", "0K3MNGKN99RQ"), 0);
    // ASE-100用户授权包
    ASSERT_EQ(acsLicenseManager->VerifyActiveCode("ASEUH-676DL-3HJKR-9B2D2-W2BBK-P53LA", "005056AF7572", "08RY92YTRF98"), 0);
    // ASE-1000用户授权包
    ASSERT_EQ(acsLicenseManager->VerifyActiveCode("ASEUS-F1XA8-ZCA0C-2467Q-7Z2XR-MEUBK", "005056AF7572", "0418VEGXRJN0"), 0);
    // ASE-10000用户授权包
    ASSERT_EQ(acsLicenseManager->VerifyActiveCode("ASEUM-8ATRF-RH197-8HCVB-K2Z5F-JLJ8X", "005056AF7572", "0MJGP40R5N3S"), 0);
    // ASE-站点用户授权包
    ASSERT_EQ(acsLicenseManager->VerifyActiveCode("ASUAG-1HWEY-EPPHZ-QZYVX-926GD-EKPMD", "005056AF7572", "06LG3S8JFPP4"), 0);
    // EDMS
    // S6020-ND1主模块
    ASSERT_EQ(acsLicenseManager->VerifyActiveCode("A6NDB-0N2P9-T3KWX-D31DP-UXZX6-FVT3N", "005056AF7572", "0EF4PUY1PLLQ"), 0);
    // RX4020-ND1主模块
    ASSERT_EQ(acsLicenseManager->VerifyActiveCode("A6RXB-2VKK7-8HGPC-090FL-5GCYS-6PADQ", "005056AF7572", "0GBAH6CF951Y"), 0);
    // NDE-S主模块
    ASSERT_EQ(acsLicenseManager->VerifyActiveCode("A6NSB-B8R0W-JVQUN-B8AE9-DGGLE-EM57A", "005056AF7572", "06LG3S8JFPHQ"), 0);
    // RX4020-ND1节点模块
    ASSERT_EQ(acsLicenseManager->VerifyActiveCode("A6XNA-BYHJA-KVGKV-112YC-LKC0Z-R1VTB", "005056AF7572", "0KFMXEMDN7NM"), 0);
    // NDE-S节点模块
    ASSERT_EQ(acsLicenseManager->VerifyActiveCode("A6ESN-56DL4-JXX9P-Q1PX0-8TAM9-YG2C0", "005056AF7572", "02D4L2Y9DV9E"), 0);
    // S6020-ND1升级模块
    ASSERT_EQ(acsLicenseManager->VerifyActiveCode("A6SUE-BULQ6-GR4Y2-QZ777-4NUT5-FZTMG", "005056AF7572", "0ZBS7ZA1NRVC"), 0);
    // NDS-100用户授权包
    ASSERT_EQ(acsLicenseManager->VerifyActiveCode("NDSHA-0QEWQ-6QAPF-F2PZ3-4QQN4-2KA5B", "005056AF7572", "0MJUX465NL5U"), 0);
    // NDS-1000用户授权包
    ASSERT_EQ(acsLicenseManager->VerifyActiveCode("ADS1T-QLCZ9-HHTBN-0EYLM-8AYAE-Z3JKQ", "005056AF7572", "0SLCLA01DDF6"), 0);
    // NDR-100用户授权包
    ASSERT_EQ(acsLicenseManager->VerifyActiveCode("NDRHA-0TAU5-1P5L7-G88E7-PMBCB-J5UYU", "005056AF7572", "009WBC29X3J4"), 0);
    // NDR-1000用户授权包
    ASSERT_EQ(acsLicenseManager->VerifyActiveCode("NDRTA-DF7UN-MTC5Z-92JH8-HGFVP-KNKVX", "005056AF7572", "0YNAHSC795HG"), 0);
    // NDR-10000用户授权包
    ASSERT_EQ(acsLicenseManager->VerifyActiveCode("NDRMA-CZ2N8-G7LJD-L3QVH-VQYXA-635UJ", "005056AF7572", "0418VEGX3XJS"), 0);
    // NDR用户场地授权包
    ASSERT_EQ(acsLicenseManager->VerifyActiveCode("NDSUA-0S1T7-44SAG-ZSJF4-L3884-DC4PN", "005056AF7572", "0G7C9A0LRDF0"), 0);
    // NDE-100用户授权包
    ASSERT_EQ(acsLicenseManager->VerifyActiveCode("NEHUA-6237P-23C8M-JQ7DU-4BGCP-JTR4X", "005056AF7572", "0W5KJMQTH15Q"), 0);
    // NDE-1000用户授权包
    ASSERT_EQ(acsLicenseManager->VerifyActiveCode("NETUA-21CBW-AQN21-A8WZG-JVQJ0-4HGQL", "005056AF7572", "0KFMXEMDN7NM"), 0);
    // NDE-10000用户授权包
    ASSERT_EQ(acsLicenseManager->VerifyActiveCode("NEMUA-9KK93-LWD96-545KX-W9MDQ-KUW50", "005056AF7572", "0SVW102XB3TQ"), 0);
    // NDE用户场地授权包
    ASSERT_EQ(acsLicenseManager->VerifyActiveCode("NESUA-4L74G-2JULT-S9GL2-X1UCG-H0WNP", "005056AF7572", "009UDS651VPZ"), 0);
    // ASS
    // EX6020主模块
    ASSERT_EQ(acsLicenseManager->VerifyActiveCode("AEX62-8SETD-PNSA4-S4LGY-P4S05-DYFCM", "005056AF7572", "0SLC3W8VFBHK"), 0);
    // ASS-S主模块
    ASSERT_EQ(acsLicenseManager->VerifyActiveCode("ASSSB-ZD41Z-51JL3-AUSF4-K3YA6-Q21R4", "005056AF7572", "0QR8FZGBPJB6"), 0);
    // EX-10用户授权包
    ASSERT_EQ(acsLicenseManager->VerifyActiveCode("AX10A-33KVJ-VVX7D-Y3KPS-2G29Q-2DGLR", "005056AF7572", "0YN27GE9NHRK"), 0);
    // EX-100用户授权包
    ASSERT_EQ(acsLicenseManager->VerifyActiveCode("AX100-WTCEX-VJW42-RFV5Y-BFC3M-FPF28", "005056AF7572", "0MJUX465NV5U"), 0);
    // ASS-10用户授权包
    ASSERT_EQ(acsLicenseManager->VerifyActiveCode("AS10A-VRPC3-EB32T-0BX4F-51UHQ-BM610", "005056AF7572", "0ANQ7AE7NJRQ"), 0);
    // ASS-100用户授权包
    ASSERT_EQ(acsLicenseManager->VerifyActiveCode("AS100-Z26AA-350VK-KV01W-0WKHC-JZYQK", "005056AF7572", "0CJWXM6BNF5G"), 0);
    // 测试授权码
    // AnyShare 30天测试授权码
    ASSERT_EQ(acsLicenseManager->VerifyActiveCode("ASTT3-29FLS-F0AU8-JERDQ-6UMAN-CP3MV", "005056AF7572", "02D4L2Y9DL9E"), 0);
    // AnyShare 90天测试授权码
    ASSERT_EQ(acsLicenseManager->VerifyActiveCode("A7TT9-0XH78-LFPQN-VUWV3-2W0YZ-KTV33", "005056AF7572", "0YT0JUATHBLC"), 0);
    // AnyShare 180天测试授权码
    ASSERT_EQ(acsLicenseManager->VerifyActiveCode("A7TT8-50ESA-GWJU4-WCYD8-DWH5D-4Y572", "005056AF7572", "0ZBS7ZA1NRVC"), 0);
    // AnyShare 360天测试授权码
    ASSERT_EQ(acsLicenseManager->VerifyActiveCode("A7TTY-2XD3C-LNQ0X-5AFDT-KS09U-AV9ML", "005056AF7572", "0CJWXK6BNF1G"), 0);
    // AnyShare大版本升级模块
    ASSERT_EQ(acsLicenseManager->VerifyActiveCode("A6UPB-8LW7N-09XZL-E9J1V-HS8WG-UJKY2", "005056AF0AF7", "0C32DW0TGSNP"), 0);
    // AnySare在线表格
    //AnyShare在线表格主模块
    ASSERT_EQ(acsLicenseManager->VerifyActiveCode("ECBSE-069LD-ESE9B-ZJD8P-DAEYX-0HJHZ", "005056AF7572", "06VU1Q2FBXTK"), 0);
    //在线表格50用户授权包
    ASSERT_EQ(acsLicenseManager->VerifyActiveCode("EC5UR-3FREH-TYXLZ-R689Z-CU0EK-ZSW0L", "005056AF7572", "0ZBS7YAXNP9Z"), 0);
    //在线表格100用户授权包
    ASSERT_EQ(acsLicenseManager->VerifyActiveCode("ECHUR-479MP-KCX6T-ESQ2G-RG60G-90QKR", "005056AF7572", "0MJUX62T7X7W"), 0);
    //在线表格200用户授权包
    ASSERT_EQ(acsLicenseManager->VerifyActiveCode("ECTUR-15SED-REZZA-2B1RZ-2GZJJ-9NKG9", "005056AF7572", "0WDYLZZPHDBU"), 0);
    //在线表格500用户授权包
    ASSERT_EQ(acsLicenseManager->VerifyActiveCode("ECFUR-7ARUQ-TSK87-TJYYR-2K8B5-XWLNH", "005056AF7572", "08PMB4KXV93E"), 0);
    //在线表格1000用户授权包
    ASSERT_EQ(acsLicenseManager->VerifyActiveCode("ECSUR-0T7D2-22VAP-4JTX3-M8K2D-6ZYLK", "005056AF7572", "0CJWXK2L737E"), 0);
    //在线表格2000用户授权包
    ASSERT_EQ(acsLicenseManager->VerifyActiveCode("ECTTU-EVSJY-K4GV9-6N9BC-3DG8R-YCJRK", "005056AF7572", "0CJWXK2L737E"), 0);
    //在线表格5000用户授权包
    ASSERT_EQ(acsLicenseManager->VerifyActiveCode("ECFTU-32Q23-JCKXE-DL87A-HLZR9-YNSPD", "005056AF7572", "06VENY4DLHXG"), 0);
    //在线表格10000用户授权包
    ASSERT_EQ(acsLicenseManager->VerifyActiveCode("ECFMU-A56F0-BMLGM-QRV24-3HANJ-4UV30", "005056AF7572", "08PYD2ZV1DXA"), 0);
    //在线表格场地用户授权包
    ASSERT_EQ(acsLicenseManager->VerifyActiveCode("ECSTR-1KCZV-44WQA-GF7HC-PPGPD-ZLJN6", "005056AF7572", "0ANQ7AE7NJRQ"), 0);
}
