import * as React from 'react';
import { TextArea, Panel } from '@/ui/ui.desktop'
import Dialog from '@/ui/Dialog2/ui.desktop';
import AdvancedConfigBase from './component.base';
import styles from './styles.view';
import __ from './locale';

export default class AdvancedConfig extends AdvancedConfigBase {
    render() {
        const { title, onRequestClose } = this.props
        const { config, isInvalidFormat } = this.state
        return (
            <Dialog
                width={440}
                title={__('高级配置')}
                onClose={() => onRequestClose()}
            >
                <Panel>
                    <div className={styles['panel-main']}>
                        <p className={styles['title']}>{title}</p>
                        <TextArea
                            value={config}
                            width={380}
                            height={180}
                            onChange={(value) => { this.handleConfigChange(value) }}
                        />
                        <div className={styles['tips']}>
                            {
                                isInvalidFormat
                                    ?
                                    __('参数格式配置错误，请您进行检查。')
                                    : null
                            }
                        </div>
                    </div>
                    <Panel.Footer>
                        <Panel.Button
                            theme='oem'
                            onClick={() => this.confirm()}
                        >
                            {__('确定')}
                        </Panel.Button>
                        <Panel.Button
                            onClick={() => onRequestClose()}
                        >
                            {__('取消')}
                        </Panel.Button>
                    </Panel.Footer>
                </Panel>
            </Dialog >
        )
    }
}