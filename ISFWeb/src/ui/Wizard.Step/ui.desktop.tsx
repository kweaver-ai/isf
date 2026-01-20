import React from 'react';
import classnames from 'classnames';
import styles from './styles.desktop';

const WizardStep: React.FunctionComponent<UI.WizardStep.Props> = function WizardStep({ title, active, onEnter, onBeforeLeave, onLeave, children, role }) {
    return (
        <div role={role} className={classnames([styles['content']], { [styles['active']]: active })}>
            {
                children
            }
        </div>
    )
}

export default WizardStep;