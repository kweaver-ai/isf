import * as React from 'react';
import { Form, Button } from '@/ui/ui.desktop';
import { ValidateBox, TextBox } from '@/sweet-ui';
import SetUpKeywordSyncBase from './component.base';
import { ValidateMessage } from './helper';
import styles from './styles.view';
import __ from './locale';

export default class SetUpKeywordSync extends SetUpKeywordSyncBase {

    render() {
        const {
            keywordInput: {
                departNameKeys,
                departThirdIdKeys,
                loginNameKeys,
                displayNameKeys,
                emailKeys,
                userThirdIdKeys,
                groupKeys,
                subOuFilter,
                subUserFilter,
                baseFilter,
                statusKeys,
                idcardNumberKeys,
                telNumberKeys,
            },
            isEditStatus,
            validateStatus: {
                departNameKeysValidateStatus,
                departThirdIdKeysValidateStatus,
                loginNameKeysValidateStatus,
                displayNameKeysValidateStatus,
                emailKeysValidateStatus,
                userThirdIdKeysValidateStatus,
                subOuFilterValidateStatus,
                subUserFilterValidateStatus,
                baseFilterValidateStatus,
            },
        } = this.state;

        return (
            <div className={styles['container']}>
                <Form role={'ui-form'} className={styles['content']}>
                    <Form.Row role={'ui-form.row'}>
                        <Form.Label role={'ui-form.label'} className={styles['lable']}>
                            {__('子部门搜索Filter：')}
                        </Form.Label>
                        <Form.Field role={'ui-form.field'}>
                            <span className={styles['mark']}>*</span>
                            <ValidateBox
                                role={'sweetui-validatebox'}
                                width={320}
                                value={subOuFilter}
                                validateState={subOuFilterValidateStatus}
                                validateMessages={ValidateMessage}
                                onBlur={() => this.handleValidate('subOuFilter')}
                                onValueChange={({ detail }) => this.handleInputChange(detail, 'subOuFilter')}
                            />
                        </Form.Field>
                    </Form.Row>
                    <Form.Row role={'ui-form.row'}>
                        <Form.Label role={'ui-form.label'} className={styles['lable']}>
                            {__('子用户搜索Filter：')}
                        </Form.Label>
                        <Form.Field role={'ui-form.field'}>
                            <span className={styles['mark']}>*</span>
                            <ValidateBox
                                role={'sweetui-validatebox'}
                                width={320}
                                value={subUserFilter}
                                validateState={subUserFilterValidateStatus}
                                validateMessages={ValidateMessage}
                                onBlur={() => this.handleValidate('subUserFilter')}
                                onValueChange={({ detail }) => this.handleInputChange(detail, 'subUserFilter')}
                            />
                        </Form.Field>
                    </Form.Row>
                    <Form.Row role={'ui-form.row'}>
                        <Form.Label role={'ui-form.label'} className={styles['lable']}>
                            {__('部门和用户Filter：')}
                        </Form.Label>
                        <Form.Field role={'ui-form.field'}>
                            <span className={styles['mark']}>*</span>
                            <ValidateBox
                                role={'sweetui-validatebox'}
                                width={320}
                                value={baseFilter}
                                validateState={baseFilterValidateStatus}
                                validateMessages={ValidateMessage}
                                onBlur={() => this.handleValidate('baseFilter')}
                                onValueChange={({ detail }) => this.handleInputChange(detail, 'baseFilter')}
                            />
                        </Form.Field>
                    </Form.Row>
                    <Form.Row role={'ui-form.row'}>
                        <Form.Label role={'ui-form.label'} className={styles['lable']}>
                            {__('部门名关键字：')}
                        </Form.Label>
                        <Form.Field role={'ui-form.field'}>
                            <span className={styles['mark']}>*</span>
                            <ValidateBox
                                role={'sweetui-validatebox'}
                                width={320}
                                value={departNameKeys}
                                validateState={departNameKeysValidateStatus}
                                validateMessages={ValidateMessage}
                                onBlur={() => this.handleValidate('departNameKeys')}
                                onValueChange={({ detail }) => this.handleInputChange(detail, 'departNameKeys')}
                            />
                        </Form.Field>
                    </Form.Row>
                    <Form.Row role={'ui-form.row'}>
                        <Form.Label role={'ui-form.label'} className={styles['lable']}>
                            {__('部门ID关键字：')}
                        </Form.Label>
                        <Form.Field role={'ui-form.field'}>
                            <span className={styles['mark']}>*</span>
                            <ValidateBox
                                role={'sweetui-validatebox'}
                                width={320}
                                value={departThirdIdKeys}
                                validateState={departThirdIdKeysValidateStatus}
                                validateMessages={ValidateMessage}
                                onBlur={() => this.handleValidate('departThirdIdKeys')}
                                onValueChange={({ detail }) => this.handleInputChange(detail, 'departThirdIdKeys')}
                            />
                        </Form.Field>
                    </Form.Row>
                    <Form.Row role={'ui-form.row'}>
                        <Form.Label role={'ui-form.label'} className={styles['lable']}>
                            {__('登录名关键字：')}
                        </Form.Label>
                        <Form.Field role={'ui-form.field'}>
                            <span className={styles['mark']}>*</span>
                            <ValidateBox
                                role={'sweetui-validatebox'}
                                width={320}
                                value={loginNameKeys}
                                validateState={loginNameKeysValidateStatus}
                                validateMessages={ValidateMessage}
                                onBlur={() => this.handleValidate('loginNameKeys')}
                                onValueChange={({ detail }) => this.handleInputChange(detail, 'loginNameKeys')}
                            />
                        </Form.Field>
                    </Form.Row>
                    <Form.Row role={'ui-form.row'}>
                        <Form.Label role={'ui-form.label'} className={styles['lable']}>
                            {__('显示名关键字：')}
                        </Form.Label>
                        <Form.Field role={'ui-form.field'}>
                            <span className={styles['mark']}>*</span>
                            <ValidateBox
                                role={'sweetui-validatebox'}
                                width={320}
                                value={displayNameKeys}
                                validateState={displayNameKeysValidateStatus}
                                validateMessages={ValidateMessage}
                                onBlur={() => this.handleValidate('displayNameKeys')}
                                onValueChange={({ detail }) => this.handleInputChange(detail, 'displayNameKeys')}
                            />
                        </Form.Field>
                    </Form.Row>
                    <Form.Row role={'ui-form.row'}>
                        <Form.Label role={'ui-form.label'} className={styles['lable']}>
                            {__('邮箱关键字：')}
                        </Form.Label>
                        <Form.Field role={'ui-form.field'}>
                            <span className={styles['mark']}>*</span>
                            <ValidateBox
                                role={'sweetui-validatebox'}
                                width={320}
                                value={emailKeys}
                                validateState={emailKeysValidateStatus}
                                validateMessages={ValidateMessage}
                                onBlur={() => this.handleValidate('emailKeys')}
                                onValueChange={({ detail }) => this.handleInputChange(detail, 'emailKeys')}
                            />
                        </Form.Field>
                    </Form.Row>
                    <Form.Row role={'ui-form.row'}>
                        <Form.Label role={'ui-form.label'} className={styles['lable']}>
                            {__('用户ID关键字：')}
                        </Form.Label>
                        <Form.Field role={'ui-form.field'}>
                            <span className={styles['mark']}>*</span>
                            <ValidateBox
                                role={'sweetui-validatebox'}
                                width={320}
                                value={userThirdIdKeys}
                                validateState={userThirdIdKeysValidateStatus}
                                validateMessages={ValidateMessage}
                                onBlur={() => this.handleValidate('userThirdIdKeys')}
                                onValueChange={({ detail }) => this.handleInputChange(detail, 'userThirdIdKeys')}
                            />
                        </Form.Field>
                    </Form.Row>
                    <Form.Row role={'ui-form.row'}>
                        <Form.Label role={'ui-form.label'} className={styles['lable']}>
                            {__('身份证号关键字：')}
                        </Form.Label>
                        <Form.Field role={'ui-form.field'}>
                            <TextBox
                                role={'sweetui-textbox'}
                                width={320}
                                value={idcardNumberKeys}
                                onValueChange={({ detail }) => this.handleInputChange(detail, 'idcardNumberKeys')}
                            />
                        </Form.Field>
                    </Form.Row>
                    <Form.Row role={'ui-form.row'}>
                        <Form.Label role={'ui-form.label'} className={styles['lable']}>
                            {__('手机号码关键字：')}
                        </Form.Label>
                        <Form.Field role={'ui-form.field'}>
                            <TextBox
                                role={'sweetui-textbox'}
                                width={320}
                                value={telNumberKeys}
                                onValueChange={({ detail }) => this.handleInputChange(detail, 'telNumberKeys')}
                            />
                        </Form.Field>
                    </Form.Row>
                    <Form.Row role={'ui-form.row'}>
                        <Form.Label role={'ui-form.label'} className={styles['lable']}>
                            {__('安全组关键字：')}
                        </Form.Label>
                        <Form.Field role={'ui-form.field'}>
                            <TextBox
                                role={'sweetui-textbox'}
                                width={320}
                                value={groupKeys}
                                onValueChange={({ detail }) => this.handleInputChange(detail, 'groupKeys')}
                            />
                        </Form.Field>
                    </Form.Row>
                    <Form.Row role={'ui-form.row'}>
                        <Form.Label role={'ui-form.label'} className={styles['lable']}>
                            {__('禁用关键字：')}
                        </Form.Label>
                        <Form.Field role={'ui-form.field'}>
                            <TextBox
                                role={'sweetui-textbox'}
                                width={320}
                                value={statusKeys}
                                onValueChange={({ detail }) => this.handleInputChange(detail, 'statusKeys')}
                            />
                        </Form.Field>
                    </Form.Row>
                    <Form.Row role={'ui-form.row'}>
                        <Form.Label role={'ui-form.label'} className={styles['lable']}>

                        </Form.Label>
                        <Form.Field role={'ui-form.field'}>
                            {
                                isEditStatus &&
                                (
                                    <div className={styles['footer']}>
                                        <Button
                                            className={styles['button-wrapper']}
                                            onClick={this.handleRequestSaveKeyword}
                                            minWidth={80}
                                        >
                                            {__('保存')}
                                        </Button>
                                        <Button
                                            className={styles['button-wrapper']}
                                            onClick={this.handleCancelEdit}
                                            minWidth={80}
                                        >
                                            {__('取消')}
                                        </Button>
                                    </div>
                                )
                            }
                        </Form.Field>
                    </Form.Row>

                </Form>
            </div>
        )
    }
}