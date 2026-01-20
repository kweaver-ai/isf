import { EACP } from '../thrift';
/**
 * 设置消息开关状态
 * @param msgNotifyStatus
 */
export const setMessageNotifyStatus: Core.EACP.SetMessageNotifyStatus = function (msgNotifyStatus) {
    return EACP('EACP_SetMessageNotifyStatus', [msgNotifyStatus]);
}

/**
 * 获取消息开关
 */
export const getMessageNotifyStatus: Core.EACP.GetMessageNotifyStatus = function () {
    return EACP('EACP_GetMessageNotifyStatus');
}

/**
 * 清除超出实名共享的权限配置
 */
export const clearPermOutOfScope: Core.EACP.ClearPermOutOfScope = function () {
    return EACP('EACP_ClearPermOutOfScope');
}

/**
 *  清除超出范围的历史匿名共享
 */
export const clearLinkOutOfScope: Core.EACP.ClearLinkOutOfScope = function () {
    return EACP('EACP_ClearLinkOutOfScope');
}