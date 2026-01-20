import * as React from 'react';
import { noop, isFunction } from 'lodash';
import { Button, PopMenu } from '@/sweet-ui';
import { InlineButton } from '@/ui/ui.desktop';
import styles from './styles.view';

interface Props {
    /**
     * 操作按钮的key
     */
    triggerKey: any;

    /**
     * menu item
     */
    menus: ReadonlyArray<{
        /**
         * item 唯一标识
         */
        key?: any;

        /**
         * 图标
         */
        icon?: string;

        /**
         * 文字
         */
        text: string;

        /**
         * 主题色
         */
        theme?: string;
    }>;

    /**
     * popMenu关闭时
     */
    onRequestCloseWhenClick?: (close: () => void, e: React.MouseEvent) => void;

    /**
     * 点击trigger的时候触发
     */
    onRequestClickTrigger?: (e: React.MouseEvent) => void;

    /**
     * 点击item触发
     */
    onRequestClickMenu: (key: any) => void;
}

const OperateMenu: React.FunctionComponent<Props> = React.memo(({
    triggerKey,
    menus = [],
    onRequestCloseWhenClick,
    onRequestClickTrigger = noop,
    onRequestClickMenu = noop,
}) => {

    return (
        <PopMenu
            triggerEvent={'click'}
            freeze={false}
            anchorOrigin={['left', 'bottom']}
            alignOrigin={['left', 'top']}
            className={styles['ul']}
            onRequestCloseWhenClick={(close, e) => { isFunction(onRequestCloseWhenClick) ? onRequestCloseWhenClick(close, e) : e.stopPropagation(); close(); }}
            trigger={({ setPopupVisibleOnClick }) =>
                <InlineButton
                    key={triggerKey}
                    code={'\uf0d0'}
                    color={'#000'}
                    size={24}
                    className={styles['btn']}
                    onClick={(e) => { onRequestClickTrigger(e); setPopupVisibleOnClick(); }}
                />
            }
        >
            {
                menus.map(({ key, icon, text, theme }, index) => (
                    <PopMenu.Item
                        key={key || index}
                        className={styles['li']}
                        onClick={() => onRequestClickMenu(key)}
                    >
                        <Button
                            style={{ color: '#505050', opacity: 1 }}
                            icon={icon}
                            theme={theme || 'text'}
                            size={'auto'}
                        >
                            {text}
                        </Button>
                    </PopMenu.Item>
                ))
            }
        </PopMenu>
    )
})

export default OperateMenu