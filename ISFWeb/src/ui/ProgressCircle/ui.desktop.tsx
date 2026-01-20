import React from 'react';
import classnames from 'classnames'
import Mask from '../Mask/ui.desktop'
import Icon from '../Icon/ui.desktop'
import Centered from '../Centered/ui.desktop';
import styles from './styles.desktop';
import darkLoading from './assets/images/dark.gif'
import lightLoading from './assets/images/light.gif'

const ProgressCircle: React.FunctionComponent<UI.ProgressCircle.Props> = function ProgressCircle({
    role,
    detail,
    showMask,
    theme,
    fixedPositioned,
}) {
    return (
        <div
            role={role}
            className={styles['container']}
        >
            {
                showMask ?
                    <Mask />
                    : null
            }
            <div className={classnames({ [styles['position-fixed']]: fixedPositioned, [styles['position-static']]: !fixedPositioned })}>
                <Centered>
                    <div className={classnames(styles['loading-box'], { [styles['grey']]: !showMask })} >
                        <Icon url={theme === 'dark' ? darkLoading : lightLoading} />
                        {
                            detail ?
                                (
                                    <div className={styles['loading-message']}>
                                        {detail}
                                    </div>
                                ) :
                                null
                        }
                    </div>
                </Centered>
            </div>
        </div>
    )
}

ProgressCircle.defaultProps = {
    detail: '',
    showMask: true,
    theme: 'light',
    fixedPositioned: true,
}

export default ProgressCircle