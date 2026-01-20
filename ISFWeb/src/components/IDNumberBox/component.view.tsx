import * as React from 'react'
import IDNumberBase from './component.base';
import { ValidateBox } from '@/sweet-ui';
import { UIIcon } from '@/ui/ui.desktop';
import { ValidateState } from './component.base';
import __ from './locale'
import styles from './styles.view';
import * as edit from './assets/edit.png'

export default class IDNumber extends IDNumberBase {
    render() {
        return (
            <div className={styles['container']}>
                <label className={styles['dialog_form_label']} htmlFor="editUser-idcard">{__('身份证号：')}</label>
                <div className={styles['validatebox-wrapper']} name={'editUser-idcard'}>
                    <ValidateBox
                        role={'sweetui-validatebox'}
                        width={this.props.width}
                        value={this.state.showIDNumber}
                        disabled={this.state.status}
                        onValueChange={({ detail }) => this.handleChange(detail)}
                        validateState={this.state.IDCardStatus}
                        validateMessages={{ [ValidateState.Error]: __('请输入正确的身份证号。') }}
                    />
                </div>
                <div className={styles['icon-wrapper']}>
                    {
                        this.state.status ?
                            <UIIcon
                                role={'ui-uiicon'}
                                title={__('编辑')}
                                size={16}
                                code={'\uf085'}
                                color={'#505050'}
                                onClick={this.toggleEditStatus.bind(this, this.state.status)}
                                fallback={edit}
                            />
                            :
                            <UIIcon
                                role={'ui-uiicon'}
                                title={__('撤销')}
                                size={16}
                                code={'\uf017'}
                                color={'#505050'}
                                onClick={this.toggleEditStatus.bind(this, this.state.status)}
                            />
                    }

                </div>
            </div>
        )
    }
}
