import * as React from 'react';
import classnames from 'classnames';
import { Radio } from '@/sweet-ui';
import { Dialog2 as Dialog, Panel, Button, Form, ProgressBar, Text } from '@/ui/ui.desktop';
import MainScreenBase from './component.base';
import __ from './locale';
import styles from './styles.view';

/**
 * 覆盖/同步同名用户
 */
enum OperationType {
    /**
     * 覆盖
     */
    Cover = 1,
    /**
     * 同步
     */
    Synchronization = 2,
}

/**
 * 操作步骤
 */
enum OrganizeStep {
    /**
     * 第一步
     */
    First = 1,
    /**
     * 第二步
     */
    Second = 2,
    /**
     * 第三步
     */
    Third = 3,
}

export default class MainScreen extends MainScreenBase {

    impotdiabled() {
        const { packageFile, operationStatus, isDisable } = this.state;

        if (!isDisable) {
            return true
        } else if (packageFile && operationStatus !== null) {
            return false
        } else {
            return true
        }
    }
    render() {
        const { packageFile, operationStatus, progress, isDisable, isGetprogress } = this.state;
        const { onExportItem } = this.props;
        const packageName = packageFile ? packageFile.name : '';

        return (
            <Dialog
                title={__('导出导入用户组织')}
                onClose={() =>this.onCancel()}
            >
                <Panel>
                    <Panel.Main>
                        <div className={styles['organize-center']}>
                            <span className={styles['organize-step']}>{OrganizeStep.First}</span>
                            {__('导出用户组织模板，批量修改用户组织')}
                            <Button
                                className={styles['export-btn']}
                                onClick={onExportItem}
                                disabled={!isDisable}
                            >
                                {__('导出')}
                            </Button>
                            <div className={styles['organize-line']}></div>
                            <span className={styles['organize-step']}>{OrganizeStep.Second}</span>
                            {__('导入修改好用户组织表')}
                            <div className={styles['organize-box']}>
                                <Form.Row>
                                    <Form.Field>
                                        <Text className={packageName ? styles['text-normal'] : styles['text-unchose']}>
                                            { packageName ? packageName : __('未选择任何文件') }
                                        </Text>
                                    </Form.Field>
                                    <Form.Field>
                                        <div
                                            className={classnames(
                                                styles['btn-uploader-picker'],
                                                styles['btn'],
                                            )}
                                            ref={(select) => this.select = select}
                                        >
                                        </div>
                                    </Form.Field>
                                </Form.Row>
                                {
                                    !isDisable ?
                                        <div  className={styles['btn-mask']}></div>
                                        : null
                                }
                            </div>
                            <div className={styles['organize-line']}></div>
                            <span className={styles['organize-step']}>{OrganizeStep.Third}</span>
                            {__('导入过程中，如果发现当前系统存在同名用户的信息')}
                            <div className={styles['organize-check']}>
                                <label className={styles['approval-status']}>
                                    <Radio
                                        disabled={!isDisable}
                                        onChange={() => this.changeApprovalStatus(OperationType.Cover)}
                                        checked={operationStatus === OperationType.Cover ? true : false}
                                    >
                                        {__('覆盖同名用户')}
                                    </Radio>
                                    <span className={styles['organize-notice']}>
                                        {__('注：勾选覆盖同名用户的信息，则导入时直接将新用户的信息覆盖旧用户的信息')}
                                    </span>
                                </label>
                                <label>
                                    <Radio
                                        disabled={!isDisable}
                                        onChange={() => this.changeApprovalStatus(OperationType.Synchronization)}
                                        checked={operationStatus === OperationType.Synchronization ? true : false}
                                    >
                                        {__('跳过同名用户')}
                                    </Radio>
                                    <span className={styles['organize-notice']}>
                                        {__('注：勾选跳过同名用户的信息，则导入时直接忽略同名用户的信息')}
                                    </span>
                                </label>
                            </div>
                            <div className={styles['organize-import']}>
                                <Button
                                    className={styles['import-btn']}
                                    onClick={this.onImportItem}
                                    disabled={this.impotdiabled()}
                                >
                                    {__('导入')}
                                </Button>
                            </div>
                            {
                                progress !== 0 && isGetprogress ?
                                    <div className={styles['organize-progress']}>
                                        <ProgressBar
                                            value={progress}
                                            height={20}
                                        />
                                    </div>
                                    : null
                            }
                        </div>
                    </Panel.Main>
                </Panel>
            </Dialog>
        )
    }
}