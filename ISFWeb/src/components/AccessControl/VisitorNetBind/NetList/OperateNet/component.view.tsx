import * as React from 'react';
import Dialog from '@/ui/Dialog2/ui.desktop';
import Panel from '@/ui/Panel/ui.desktop';
import Form from '@/ui/Form/ui.desktop';
import ValidateBox from '@/ui/ValidateBox/ui.desktop';
import { Button, Radio } from '@/sweet-ui';
import { NetType, ValidateMessages, ValidateState } from '../../helper';
import __ from './locale';
import OperateNetBase from './component.base';
import styles from './styles.view.css';

export default class OperateNet extends OperateNetBase {
    render() {
        const { onRequestEditCancel, net } = this.props;
        const { validateState, error, net: { name, originIP, endIP, ip, mask, netType } } = this.state;

        return (
            <Dialog
                title={net.id ? __('编辑网段') : __('添加网段')}
                onClose={onRequestEditCancel}
            >
                <Panel>
                    <Panel.Main>
                        <Form>
                            <Form.Row>
                                <Form.Label>
                                    <span>{__('网段名称：')}</span>
                                </Form.Label>
                                <Form.Field>
                                    <ValidateBox
                                        validateMessages={ValidateMessages}
                                        validateState={validateState.name}
                                        value={name}
                                        placeholder={__('可选填')}
                                        onChange={(name) => this.setValue('name', name)}
                                    />
                                </Form.Field>
                            </Form.Row>
                            <Form.Row>
                                <Form.Label>
                                    <Radio
                                        value={NetType.Range}
                                        checked={netType === NetType.Range}
                                        onChange={this.changeNetType.bind(this)}
                                    >
                                        <span >{__('起始IP：')}</span>
                                    </Radio>
                                </Form.Label>
                                <Form.Field>
                                    <ValidateBox
                                        className={styles['space']}
                                        validateMessages={ValidateMessages}
                                        validateState={validateState.originIP}
                                        value={netType === NetType.Range ? originIP : ''}
                                        disabled={netType !== NetType.Range}
                                        validator={(input) => this.validateInput(input)}
                                        onChange={(originIP) => this.setValue('originIP', originIP)}
                                    />
                                </Form.Field>
                            </Form.Row>
                            <Form.Row>
                                <Form.Label>
                                    <span className={styles['net-label']}>{__('终止IP：')}</span>
                                </Form.Label>
                                <Form.Field>
                                    <ValidateBox
                                        validateMessages={ValidateMessages}
                                        validateState={validateState.endIP}
                                        value={netType === NetType.Range ? endIP : ''}
                                        disabled={netType !== NetType.Range}
                                        validator={(input) => this.validateInput(input)}
                                        onChange={(endIP) => this.setValue('endIP', endIP)}
                                    />
                                </Form.Field>
                            </Form.Row>
                            <Form.Row>
                                <Form.Label>
                                    <Radio
                                        value={NetType.Mask}
                                        checked={netType === NetType.Mask}
                                        onChange={this.changeNetType.bind(this)}
                                    >
                                        <span >{__('IP地址：')}</span>
                                    </Radio>
                                </Form.Label>
                                <Form.Field>
                                    <ValidateBox
                                        className={styles['space']}
                                        validateMessages={ValidateMessages}
                                        validateState={validateState.ip}
                                        value={netType === NetType.Mask ? ip : ''}
                                        disabled={netType !== NetType.Mask}
                                        validator={(input) => this.validateInput(input)}
                                        onChange={(ip) => this.setValue('ip', ip)}
                                    />
                                </Form.Field>
                            </Form.Row>
                            <Form.Row>
                                <Form.Label>
                                    <span className={styles['net-label']}>{__('子网掩码：')}</span>
                                </Form.Label>
                                <Form.Field>
                                    <ValidateBox
                                        validateMessages={ValidateMessages}
                                        validateState={validateState.mask}
                                        value={netType === NetType.Mask ? mask : ''}
                                        disabled={netType !== NetType.Mask}
                                        validator={(input) => this.validateInput(input)}
                                        onChange={(mask) => this.setValue('mask', mask)}
                                    />
                                </Form.Field>
                            </Form.Row>
                        </Form>
                        <div className={styles['error']}>
                            {error === ValidateState.InvalidRange ? __('终止IP不能小于起始IP。') : ''}
                        </div>
                    </Panel.Main>
                    <Panel.Footer>
                        <Button
                            theme={'oem'}
                            width='auto'
                            className={styles['save-btn']}
                            onClick={() => this.saveNet()}
                        >
                            {__('确定')}
                        </Button>
                        <Button
                            width='auto'
                            onClick={onRequestEditCancel}
                        >
                            {__('取消')}
                        </Button>
                    </Panel.Footer>
                </Panel>
            </Dialog >
        )
    }

}