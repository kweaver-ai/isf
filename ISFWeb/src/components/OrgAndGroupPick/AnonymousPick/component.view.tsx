import * as React from 'react';
import { UIIcon } from '@/ui/ui.desktop';
import { CheckBox } from '@/sweet-ui';
import AnonymousPickBase from './component.base';
import styles from './styles.view';
import __ from './locale';

export default class AnonymousPick extends AnonymousPickBase {
    render() {
        const { disabled, isMult } = this.props;
        const { checkStatus } = this.state;

        return (
            <div className={styles['container']}>
                {
                    isMult ? (
                        <CheckBox
                            role={'sweetui-checkbox'}
                            disabled={disabled}
                            className={styles['check-box']}
                            checked={checkStatus}
                            onCheckedChange={({ detail }) => {
                                this.checkAnonymous(detail)
                            }}
                        />
                    ) : null
                }
                <div
                    className={styles['name']}
                    title={__('匿名用户')}
                    onClick={() => isMult ? this.checkAnonymous(!checkStatus) : this.selectAnonymous()}
                >
                    <UIIcon
                        role={'ui-uiicon'}
                        code={'\uf007'}
                        size={16}
                    />
                    <div className={styles['text-custom']}>{__('匿名用户')}</div>
                </div>
            </div>
        )
    }
}