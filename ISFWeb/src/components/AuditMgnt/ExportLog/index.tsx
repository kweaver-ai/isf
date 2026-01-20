import * as React from 'react';
import { timer } from '@/util/timer';
import { jsencrypt2048 } from '@/core/auth';
import Form from '@/ui/Form/ui.desktop';
import { exportHistoryLog, getDownloadResult, getPasswordStatus, getGenCompressFileStatus } from '@/core/apis/console/auditlog';
import styles from './styles.view.css';
import __ from './locale';
import { Button, Input, Modal, Spin } from 'antd';
import intl from 'react-intl-universal';
import { LoadingOutlined } from '@ant-design/icons';
export enum ValidateState {
    /**
   * 验证合法
   */
    Normal,

    /**
   * 验证不合法
   */
    Diff,
}

export enum ExportStatus {
    /**
   * 涉密子开关关闭
   */
    SWITCH_CLOSE,

    /**
   * 涉密子开关开启
   */
    SWITCH_OPEN,

    /**
   * 转圈圈组件正在加载中
   */
    LOADING,
}

/**
 * 日志类型
 */
export enum LogType {
    /**
   * 登录日志
   */
    NCT_LT_LOGIN = 10,

    /**
   * 管理日志
   */
    NCT_LT_MANAGEMENT = 11,

    /**
   * 操作日志
   */
    NCT_LT_OPEARTION = 12,
}

export enum LogDetail {
    /**
   * 活跃日志
   */
    Active,

    /**
   * 历史日志
   */
    History,
}

export default function ExportLog({ id, onExportComplete, onRequestCancel }) {
    const [password, setPassword] = React.useState('');
    const [passwordAgain, setPasswordAgain] = React.useState('');
    const [isSamePassword, setIsSamePassword] = React.useState(true);
    const [validateState, setValidateState] = React.useState(
        ValidateState.Normal,
    );
    const [exportStatus, setExportStatus] = React.useState(
        ExportStatus.SWITCH_CLOSE,
    );

    const onPasswordInputFirstChange = (value) => {
        setPassword(value);
        setValidateState(ValidateState.Normal);
        setIsSamePassword(true);
    };

    const onpasswordAgainChange = (value) => {
        setPasswordAgain(value);
        setIsSamePassword(true);
    };

    const isPasswordValidate = (value) => {
        if (!/[a-z]/.test(value)) {
            return false;
        }
        if (!/[A-Z]/.test(value)) {
            return false;
        }
        if (!/[0-9]/.test(value)) {
            return false;
        }
        if (value.length <  10 || value.length > 100) {
            return false;
        }
        if (!/^[\w~`!@#$%\-,\.]+$/i.test(value)) {
            return false;
        }
        return true;
    };

    const checkedForm = (): boolean => {
        if (!isPasswordValidate(password)) {
            setValidateState(ValidateState.Diff);
            return false;
        }
        if (password !== passwordAgain) {
            setIsSamePassword(false);
            return false;
        }
        return true;
    };

    const submitExport = async () => {
        if (checkedForm()) {
            downloadLog(true);
        }
    };

    const getCompressProgress = (taskId) => {
        let time = timer(() => {
            return getGenCompressFileStatus({ taskid: taskId }).then(({ status }) => {
                if (status) {
                    time();
                    getDownloadResult({ taskid: taskId }).then(({ url }) => {

                        onExportComplete(
                            url,
                        );
                    });
                }
            });
        }, 1000);
    };

    const downloadLog = async (needPassword?: boolean) => {
        setExportStatus(ExportStatus.LOADING);

        const { task_id } = await exportHistoryLog({
            obj_id: id,
            ...needPassword && { pwd: jsencrypt2048(password) },
        });

        getCompressProgress(task_id);
    };

    React.useEffect(() => {
        async function fetchData() {
            const { status } = await getPasswordStatus()
            // 子开关状态
            if (status) {
                setExportStatus(ExportStatus.SWITCH_OPEN);
            } else {
                setExportStatus(ExportStatus.LOADING);
                downloadLog()
            }
        }

        fetchData()
    }, [])

    return (
        <>
            {exportStatus === ExportStatus.SWITCH_OPEN ? (
                <Modal 
                    centered
                    maskClosable={false}
                    open={true}
                    width={480}
                    title={__('提示')} 
                    onCancel={onRequestCancel}
                    footer={[
                        <Button key="ok" type="primary" disabled={!(password && passwordAgain)} onClick={submitExport}>
                            {intl.get('ok')}
                        </Button>,
                        <Button key="cancel" onClick={onRequestCancel} >
                            {intl.get('cancel')}
                        </Button>,
                    ]}
                    getContainer={document.getElementById('isf-web-plugins') as HTMLElement}
                >
                    <div className={styles['export-log']}>
                        <div className={styles['input-tips']}>
                            {__('您导出的日志将会加密打包，请输入密码：')}
                        </div>
                        <div className={styles['content']}>
                            <Form>
                                <Form.Row>
                                    <Form.Label>{__('密码：')}</Form.Label>
                                    <Form.Field>
                                        <Input.Password
                                            style={{ width: 340 }}
                                            placeholder={__('请输入密码')}
                                            value={password}
                                            onChange={(e) => {
                                                onPasswordInputFirstChange(e.target.value);
                                            }}
                                            status={validateState === ValidateState.Diff || !isSamePassword ? 'error' : ''}
                                        />
                                    </Form.Field>
                                </Form.Row>
                                {
                                    validateState === ValidateState.Diff && 
                                    <Form.Row>
                                        <Form.Label></Form.Label>
                                        <Form.Field>
                                            <div className={styles['warn-font']}>
                                                {__('密码为${min}~100位，必须同时包含 大小写英文字母与数字，允许包含~!%#$@-_. 字符', { min: 10 },)}
                                            </div>
                                        </Form.Field>
                                    </Form.Row>
                                }
                                <Form.Row>
                                    <Form.Label>{__('确认密码：')}</Form.Label>
                                    <Form.Field>
                                        <Input.Password
                                            style={{ width: 340 }}
                                            placeholder={__('请再次输入密码')}
                                            value={passwordAgain}
                                            onChange={(e) => {
                                                onpasswordAgainChange(e.target.value);
                                            }}
                                            status={isSamePassword ? '' : 'error'}
                                        />
                                    </Form.Field>
                                </Form.Row>
                                {!isSamePassword ? (
                                    <Form.Row>
                                        <Form.Label></Form.Label>
                                        <Form.Field>
                                            <div className={styles['warn-font']}>
                                                {__('两次输入的密码不一致，请重新输入')}
                                            </div>
                                        </Form.Field>
                                    </Form.Row>
                                ) : null}
                            </Form>
                        </div>
                    </div>
                </Modal>
            ) : null}
            {exportStatus === ExportStatus.LOADING ? (
                <Spin size='large' tip={__('正在打包，请稍候…')} fullscreen indicator={<LoadingOutlined style={{ color: '#fff' }} spin />}/>
            ) : null}
        </>
    );
}
