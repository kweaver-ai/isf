import * as React from 'react';
import { Icon, Centered } from '@/ui/ui.desktop';
import { Spin } from '@/sweet-ui'
import { ListTipStatus, ListTipMessage, ListTipStatusMapImg } from './helper';
import styles from './styles.view';
import __ from './locale';

interface ListTipComponentProps {
    /**
     * 列表提示状态
     */
    listTipStatus: ListTipStatus;

    /**
     * 是否在弹窗中显示
     */
    isInDialog?: boolean;

    /**
     * 图片大小
     */
    imgSize?: number;

    /**
     * 图片src
     */
    listTipStatusMapImg?: Record<string, any>;

    /**
     * 提示语
     */
    listTipMessage?: Record<string, any>;

    /**
     * 额外提示语
     */
    extraListTipMessage?: Record<string, any>;
}
/**
 * 提示组件
 */
export const getTipComponent = ({ img, size }, text?, extraText?) => {
    return (
        <Centered role={'ui-centered'}>
            <Icon
                role={'ui-icon'}
                url={img}
                size={size}
            />
            <p className={styles['tip-text']}>{text}</p>
            {extraText && <p className={styles['tip-text']}>{extraText}</p>}
        </Centered>
    )
}

/**
 * 列表加载提示信息
 */
const ListTipComponent: React.FunctionComponent<ListTipComponentProps> = React.memo(({
    listTipStatus,
    listTipMessage,
    extraListTipMessage,
    isInDialog = false,
    imgSize,
    listTipStatusMapImg,
}) => {
    const text = listTipMessage && listTipStatus in listTipMessage ? listTipMessage[listTipStatus] : ListTipMessage[listTipStatus]

    const img = listTipStatusMapImg && listTipStatus in listTipStatusMapImg ? listTipStatusMapImg[listTipStatus] : ListTipStatusMapImg[listTipStatus];

    // 是否在第二行展示额外文本
    const extraText = extraListTipMessage && listTipStatus in extraListTipMessage ? extraListTipMessage[listTipStatus] : '';

    switch (listTipStatus) {
        case ListTipStatus.Loading:
            return getTipComponent({ img, size: imgSize || 32 })

        case ListTipStatus.LightLoading:
            return (
                <Centered role={'ui-centered'}>
                    <Spin />
                </Centered>
            )

        case ListTipStatus.Empty:
            return getTipComponent({ img, size: imgSize || (isInDialog ? 40 : 64) }, text, extraText)

        case ListTipStatus.NoSearchResults:
            return getTipComponent({ img, size: imgSize || (isInDialog ? 40 : 64) }, text, extraText)

        case ListTipStatus.LoadFailed:
            return getTipComponent({ img, size: imgSize || (isInDialog ? 40 : 64) }, text, extraText)

        case ListTipStatus.NoSyncPlan:
            return getTipComponent({ img, size: imgSize || (isInDialog ? 40 : 128) }, text, extraText)

        case ListTipStatus.NoExamine:
            return getTipComponent({ img, size: imgSize || (isInDialog ? 40 : 128) }, text, extraText)

        case ListTipStatus.OrgEmpty:
            return getTipComponent({ img, size: imgSize || (isInDialog ? 40 : 64) }, text, extraText)

        case ListTipStatus.ClientOrgEmpty:
            return getTipComponent({ img: img, size: imgSize || 96 }, text, extraText)

        case ListTipStatus.ClientNoSearch:
            return getTipComponent({ img, size: imgSize || (isInDialog ? 64 : 128) }, text, extraText)

        case ListTipStatus.ClientLoadFaild:
            return getTipComponent({ img, size: imgSize || (isInDialog ? 64 : 128) }, text, extraText)

        case ListTipStatus.NoDocFlow:
            return getTipComponent({ img, size: imgSize || (isInDialog ? 64 : 128) }, text, extraText)

        default:
            return null;
    }
})

export default ListTipComponent;