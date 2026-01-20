import * as React from 'react'
import classnames from 'classnames'
import { UIIcon, Expand } from '@/ui/ui.desktop'
import styles from './styles.view';

const DomainExpandItem = ({ title, disabled, isExpand, onExpandItem, showIcon = true, children }) => {

    return (
        <div className={styles['expand-info']}>
            <div
                className={classnames(styles['expand-header'], { [styles['disabled']]: disabled })}
                onClick={() => onExpandItem(!isExpand)}
            >
                <div className={styles['expand-right']}>
                    {
                        showIcon ? <UIIcon
                            role={'ui-uiicon'}
                            className={classnames(styles['expand-icon'], { [styles['disabled']]: disabled })}
                            size={16}
                            code={isExpand ? '\uf04c' : '\uf04e'}
                        /> : null
                    }
                </div>
                <div className={styles['expand-left']}>{title}</div>
            </div>
            <Expand role={'ui-expand'} open={isExpand}>
                <div className={classnames({ [styles['content']]: showIcon })}>{children}</div>
            </Expand>
        </div >
    )
}

export default DomainExpandItem;
