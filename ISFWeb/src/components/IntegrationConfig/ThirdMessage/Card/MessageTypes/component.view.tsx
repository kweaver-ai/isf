import * as React from 'react';
import { noop } from 'lodash';
import { UIIcon } from '@/ui/ui.desktop';
import { ValidateBox } from '@/sweet-ui';
import { ValidateMessage, MessageConfigItem } from '../../helper';
import styles from './styles.view.css';
import __ from './locale';

interface MessageTypesProps {
    /**
     * 配置项
     */
    messageConfigItem: MessageConfigItem;

    /**
     * 改变输入值
     */
    onRequestConfigValueChange: (value: string) => void;

    /**
     * 获取ref
     */
    onRequestRef: (ref) => void;

    /**
     * 删除
     */
    onRequestDelete: () => void;
}

// 完全受控组件，所有的状态都从父组件接收
const MessageTypes: React.FunctionComponent<MessageTypesProps> = React.memo(({
    messageConfigItem,
    onRequestConfigValueChange = noop,
    onRequestRef = noop,
    onRequestDelete = noop,
}) => {
    return (
        <div>
            <div className={styles['item-value']} ref={(ref) => onRequestRef(ref)}>
                <div className={styles['tips']}>{__('值：')}</div>
                <div className={styles['mark']}>{'*'}</div>
                {
                    <div className={styles['validate']}>
                        <ValidateBox
                            width={710}
                            placeholder={__('请在此处输入URI路径')}
                            value={messageConfigItem.configValue}
                            validateState={messageConfigItem.valueValidateStatus}
                            validateMessages={ValidateMessage}
                            onValueChange={({ detail: value }) => onRequestConfigValueChange(value)}
                        />
                    </div>
                }
            </div>
            <div className={styles['item-delete']}>
                <UIIcon
                    code={'\uf014'}
                    size={12}
                    color={'#505050'}
                    className={styles['ui-icon']}
                    onClick={onRequestDelete}
                />
            </div>
        </div>
    )
})

export default MessageTypes;