import * as React from 'react';
import { Panel, Form, Button, Dialog2 as Dialog, CheckBoxOption } from '@/ui/ui.desktop';
import { ValidateBox } from '@/sweet-ui'
import __ from './locale';
import styles from './styles.view';
import PwdManageBase, { UserType, ValidateState, systemOriginalPwd, hiddenPwd } from './component.base';

const getValidateMessages = (strongPwdLength: number): any => ({
    [ValidateState.Empty]: __('此输入项不允许为空。'),
    [ValidateState.StrongPwdError]: __('密码为 ${strongPwdLength}~100 位，必须同时包含 大小写英文字母、数字与半角特殊字符。', { strongPwdLength }),
    [ValidateState.WeakPwdError]: __('密码只能包含 英文 或 数字 或 空格 或 半角特殊字符，长度范围 6~100 个字符，请重新输入。'),
    [ValidateState.OrgPwdError]: __('您输入的密码不能为初始密码，请重新输入。'),
});

export default class PwdManage extends PwdManageBase {

    render() {
        const { userType, strongPwdLength } = this.userData;
        const { pwdControl, password, lockStatus, validateState } = this.state;

        return (
            this.props.selectedUser.id === this.props.userid ?
                <div></div>
                :
                <Dialog
                    role={'ui-dialog'}
                    title={__('管控密码')}
                    width={410}
                    onClose={() => { this.props.onRequestCancel() }}
                >
                    <Panel role={'ui-panel'}>
                        <Panel.Main role={'ui-panel.main'}>
                            <div className={styles['check']}>
                                <CheckBoxOption
                                    role={'ui-checkboxoption'}
                                    disabled={userType !== UserType.LocalUser}
                                    checked={pwdControl}
                                    onChange={() => this.toggleCheck()}
                                >
                                    {__('不允许用户自主修改密码')}
                                </CheckBoxOption>
                            </div>
                            <Form role={'ui-form'}>
                                <Form.Row role={'ui-form.row'}>
                                    <Form.Label role={'ui-form.label'}>{__('用户密码：')}</Form.Label>
                                    <Form.Field role={'ui-form.field'}>
                                        <div className={styles['password']}>
                                            <ValidateBox
                                                role={'ui-validatebox'}
                                                className={styles['validate-box']}
                                                width={150}
                                                disabled={!pwdControl}
                                                value={!this.isInputPwd && password === systemOriginalPwd ? hiddenPwd : password}
                                                onValueChange={({detail}) => this.updatePwd(detail)}
                                                validateMessages={getValidateMessages(strongPwdLength)}
                                                validateState={validateState}
                                            />
                                            {
                                                pwdControl ?
                                                    <Button
                                                        role={'ui-button'}
                                                        className={styles['btn']}
                                                        onClick={() => this.handlePwdRandomly()}
                                                    >
                                                        {__('随机密码')}
                                                    </Button>
                                                    :
                                                    <Button
                                                        role={'ui-button'}
                                                        disabled={password === systemOriginalPwd || userType !== UserType.LocalUser}
                                                        className={styles['btn']}
                                                        onClick={() => this.resetPwd()}
                                                    >
                                                        {__('重置密码')}
                                                    </Button>
                                            }
                                        </div>
                                    </Form.Field>
                                </Form.Row>
                                <div className={styles['tip-container']}>
                                    <p className={styles['tip']}>{__('注：用户密码生成后您将无法再次查看，请妥善保管好您的密码。')}</p>
                                </div>
                                <Form.Row role={'ui-form'}>
                                    <Form.Label role={'ui-form.label'}>{__('账号状态：')}</Form.Label>
                                    <Form.Field role={'ui-form.field'}>
                                        {
                                            lockStatus ?
                                                <span className={styles['locked']}>{__('锁定')}</span>
                                                :
                                                <span className={styles['unlocked']}>{__('正常')}</span>
                                        }
                                        <Button
                                            role={'ui-button'}
                                            className={styles['btn']}
                                            disabled={!lockStatus}
                                            onClick={() => this.unlockUser()}
                                        >
                                            {__('解锁')}
                                        </Button>
                                    </Form.Field>
                                </Form.Row>
                            </Form>
                        </Panel.Main>
                        <Panel.Footer role={'ui-panel.footer'}>
                            <Panel.Button
                                theme='oem'
                                role={'ui-panel.button'}
                                onClick={() => this.confirm()}
                            >
                                {__('确定')}
                            </Panel.Button>
                            <Panel.Button
                                role={'ui-panel.button'}
                                onClick={() => this.cancel()}
                            >
                                {__('取消')}
                            </Panel.Button>
                        </Panel.Footer>
                    </Panel>
                </Dialog>
        )
    }
}
