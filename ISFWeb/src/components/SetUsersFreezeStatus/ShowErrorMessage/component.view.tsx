// eslint-disable-next-line spaced-comment
/// <reference path="./component.base.d.ts" />

import * as React from 'react';
import { noop } from 'lodash'
import { getErrorMessage } from '@/core/exception';
import MessageDialog from '@/ui/MessageDialog/ui.desktop';
import { Status } from '../helper';
import __ from './locale';

/**
 * 显示错误弹窗
 * @param onConfrim  // 确定事件
 * @param Message // 提示信息
 */
export default function ErrorMessage({ errorType, userInfo, freezeStatus, errorInfo, onConfirm = noop }: Components.SetUsersFreezeStatus.ShowErrorMessage.Props) {

    switch (errorType) {
        case Status.CURRENT_USER_INCLUDED:
            return (
                <MessageDialog onConfirm={onConfirm}>
                    {
                        freezeStatus ?
                            __('您无法冻结自身账号。') :
                            __('您无法解冻自身账号。')
                    }
                </MessageDialog>
            );
        default:
            return (
                <MessageDialog onConfirm={onConfirm}>
                    {getErrorMessage(errorInfo.error.errID)}
                </MessageDialog>
            )

    }
}