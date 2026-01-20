import * as React from 'react';
import { TextArea, Panel } from '@/ui/ui.desktop'
import Dialog from '@/ui/Dialog2/ui.desktop';
import AdvancedConfigBase from './component.base';
import styles from './styles.view.css';
import __ from './locale';

export default class AdvancedConfig extends AdvancedConfigBase {
    render() {
        const { onRequestClose } = this.props
        const { config, isInvalidFormat } = this.state

        return (
            <Dialog
                width={440}
                title={__('高级配置')}
                onClose={() => onRequestClose()}
            >
                <Panel>
                    <div className={styles['panel-main']}>
                        <p className={styles['title']}>{__('消息插件参数配置：')}</p>
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
                                    <p>{__('参数格式配置错误，请您进行检查。')}</p>
                                    : null
                            }
                        </div>
                    </div>
                    <Panel.Footer>
                        <Panel.Button
                            onClick={() => this.confirm()}
                            width='auto'
                        >
                            {__('确定')}
                        </Panel.Button>
                        <Panel.Button
                            onClick={() => onRequestClose()}
                            width='auto'
                        >
                            {__('取消')}
                        </Panel.Button>
                    </Panel.Footer>
                </Panel>
            </Dialog >
        )
    }
}