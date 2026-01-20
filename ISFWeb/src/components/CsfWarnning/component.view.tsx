import * as React from 'react';
import { UIIcon } from '@/ui/ui.desktop';
import __ from './locale';
import styles from './styles.view';

const CsfWarnning: React.FunctionComponent<any> = React.memo(() => {
    return (
        <div className={styles['warn']}>
            <UIIcon
                role={'ui-uiicon'}
                code={'\uf055'}
                size={'16px'}
                title={
                    <div className={styles['text']} >
                        {__('新建用户密级：变更后，')}<br />
                        {__('仅对新建用户生效')}<br />
                    </div>
                }
                color={'#555'}
            />
        </div>
    )
})

export default CsfWarnning;
