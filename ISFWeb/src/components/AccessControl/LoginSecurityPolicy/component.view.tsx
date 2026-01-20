import * as React from 'react';
import WebComponent from '../../webcomponent';
import LoginPolicy from './LoginPolicy/component.view';
import PasswordPolicy from './PasswordPolicy/component.view';
import styles from './styles.view';

export default class LoginSecurityPolicy extends WebComponent<any, any> {
    render() {
        const { navigate } = this.props

        return (
            <div className={styles['container']}>
                <PasswordPolicy navigate={navigate}/>
                <LoginPolicy navigate={navigate}/>
            </div>
        )
    }
}