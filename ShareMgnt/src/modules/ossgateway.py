from src.common.lib import raise_exception
from src.common.sharemgnt_logger import ShareMgnt_Log
from src.driven.service_access.ossgateway_config import OssgatewayDriven
from ShareMgnt.ttypes import (ncTShareMgntError,
                              ncTUsrmOSSInfo)


def get_oss_info(oss_id):
    """
    根据对象存储id获取对象存储信息
    """
    oss_info = ncTUsrmOSSInfo()
    ossgateway_driven = OssgatewayDriven()
    evfs_oss_info = None
    try:
        code, evfs_oss_info = ossgateway_driven.get_storage_info_by_id(oss_id)
        if code != 200:
            ShareMgnt_Log(
                f'get_storage_info_by_id failed: {code},{evfs_oss_info}')
            raise_exception(exp_msg=_("IDS_GET_OSS_INFO_FAILD"),
                            exp_num=ncTShareMgntError.NCT_NO_AVAILABLE_OSS)
        oss_info.ossId = evfs_oss_info["id"]
        oss_info.ossName = evfs_oss_info["name"]
        oss_info.enabled = evfs_oss_info["enabled"]
        return oss_info
    except Exception as ex:
        pass
