import * as React from 'react';
import { noop } from 'lodash';
import Dialog from '@/ui/Dialog2/ui.desktop';
import { Panel, TextArea, UIIcon } from '@/ui/ui.desktop';
import { ConfigItem, transformArrayToObject } from '../helper';
import styles from './styles.view';
import __ from './locale';

interface ResetInvalidConfig {
    /**
     * 需要还原的默认参数（非法修改的默认参数）
     */
    invalidConfig: ReadonlyArray<ConfigItem>;

    /**
     * 确认弹窗
     */
    onRequestConfirm: () => void;

    /**
     * 关闭弹窗
     */
    onRequestClose: () => void;
}

const ResetInvalidConfig: React.FunctionComponent<ResetInvalidConfig> = function ResetInvalidConfig({
    invalidConfig = [],
    onRequestConfirm = noop,
    onRequestClose = noop,
}) {
    return (
        <Dialog
            width={400}
            title={__('高级配置')}
            onClose={() => onRequestClose()}
        >
            <Panel>
                <div className={styles['panel-main']}>
                    <div className={styles['icon']}>
                        <UIIcon
                            code={'\uf076'}
                            color={'#5a8cb4'}
                            size={40}
                        />
                    </div>
                    <div className={styles['message']}>
                        <p className={styles['tips']}>{__('以下参数的默认配置不可更改或参数类型错误，已还原，其他参数编辑成功。')}</p>
                        <TextArea
                            width={282}
                            height={96}
                            disabled={true}
                            value={invalidConfig.length === 0 ? '' : JSON.stringify(transformArrayToObject(invalidConfig), null, 4)}
                        />
                    </div>
                </div>
                <Panel.Footer>
                    <Panel.Button
                        onClick={() => onRequestConfirm()}
                    >
                        {__('确定')}
                    </Panel.Button>
                </Panel.Footer>
            </Panel>
        </Dialog >
    )
}

export default ResetInvalidConfig