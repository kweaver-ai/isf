import * as React from 'react';
import classnames from 'classnames';
import { ModalDialog2, SweetIcon, ValidateBox, Radio } from '@/sweet-ui';
import Panel from '@/ui/Panel/ui.desktop';
import Form from '@/ui/Form/ui.desktop';
import NetSegmentBase from './component.base';
import { Placeholder, IpVersion, InputKey, ValidateMessages, ValidateState, OperateType } from './helper';
import styles from './styles.view.css';
import __ from './locale';

export default class NetSegment extends NetSegmentBase {
    render() {
        const {
            validateState,
            ipError,
            netInfo,
            netInfo: {
                name,
                ipVersion,
                originIP,
                endIP,
                ip,
                mask,
                prefix,
                single,
            },
        } = this.state
        const { isShowNetName, operateType, isShowTitle } = this.props
        const isIpv4 = ipVersion === IpVersion.Ipv4

        return (
            <ModalDialog2
                title={
                    isShowTitle
                        ? operateType === OperateType.Add ? __('添加网段') : __('编辑网段')
                        : ''
                }
                icons={[{
                    icon: <SweetIcon role={'sweetui-sweeticon'} name={'x'} size={16} />,
                    onClick: this.props.onRequestCancel,
                }]}
                buttons={[
                    {
                        text: __('确定'),
                        theme: 'oem',
                        onClick: this.confirm,
                        size:'auto',
                    },
                    {
                        text: __('取消'),
                        theme: 'regular',
                        onClick: this.props.onRequestCancel,
                        size:'auto',
                    },
                ]}
            >
                <Panel role={'ui-panel'}>
                    <Form>
                        {
                            isShowNetName
                                ? (
                                    <Form.Row>
                                        <Form.Label className={styles['row-bottom']}>
                                            <span className={styles['row-bottom-title']}>
                                                {__('网段名称：')}
                                            </span>
                                        </Form.Label>
                                        <Form.Field className={styles['row-bottom']}>
                                            <ValidateBox
                                                key={'name'}
                                                type={'text'}
                                                role={'sweetui-validatebox'}
                                                width={300}
                                                value={name}
                                                placeholder={Placeholder.name}
                                                onBlur={() => this.checkFrom([InputKey.Name])}
                                                onValueChange={({ detail }) => this.setIpInfo(InputKey.Name, detail)}
                                                validateState={validateState.name}
                                                validateMessages={ValidateMessages}
                                            />
                                        </Form.Field>
                                    </Form.Row>
                                )
                                : null
                        }
                        <Form.Row>
                            <Form.Label>
                                {__('IP类型：')}
                            </Form.Label>
                            <Form.Field>
                                <div
                                    className={classnames(
                                        styles['net-ip'],
                                        { [styles['net-ip-checked']]: isIpv4 },
                                    )}
                                    onClick={() => this.ipVersionSelect(IpVersion.Ipv4)}
                                >
                                    <span className={styles['ip-text']}>
                                        {'IPv4'}
                                    </span>
                                </div>
                                <div
                                    className={classnames(
                                        styles['net-ip'],
                                        { [styles['net-ip-checked']]: !isIpv4 },
                                    )}
                                    onClick={() => this.ipVersionSelect(IpVersion.Ipv6)}
                                >
                                    <span className={styles['ip-text']}>
                                        {'IPv6'}
                                    </span>
                                </div>
                            </Form.Field>
                        </Form.Row>
                        <Form.Row>
                            <Form.Label>
                                <Radio
                                    value={false}
                                    checked={!single}
                                    onChange={({ detail: { value } }) => this.ipRangeSelect(value)}
                                >
                                    <span className={styles['color-strong']} >{__('起始IP：')}</span>
                                </Radio>
                            </Form.Label>
                            <Form.Field>
                                <ValidateBox
                                    type={'text'}
                                    key={'originIP'}
                                    role={'sweetui-validatebox'}
                                    width={300}
                                    disabled={single}
                                    value={single ? '' : originIP}
                                    placeholder={isIpv4 ? Placeholder.ipv4 : Placeholder.ipv6}
                                    validator={(input) => this.validateInput(input, isIpv4)}
                                    onBlur={() => this.checkFrom([InputKey.OriginIp])}
                                    onValueChange={({ detail }) => this.setIpInfo(InputKey.OriginIp, detail)}
                                    validateState={single ? ValidateState.OK : validateState.originIP}
                                    validateMessages={ValidateMessages}
                                />
                            </Form.Field>
                        </Form.Row>
                        <Form.Row>
                            <Form.Label>
                                <span className={styles['net-label']}>
                                    {__('终止IP：')}
                                </span>
                            </Form.Label>
                            <Form.Field>
                                <ValidateBox
                                    key={'endIP'}
                                    type={'text'}
                                    role={'sweetui-validatebox'}
                                    width={300}
                                    disabled={single}
                                    value={single ? '' : endIP}
                                    placeholder={isIpv4 ? Placeholder.ipv4 : Placeholder.ipv6}
                                    validator={(input) => this.validateInput(input, isIpv4)}
                                    onBlur={() => this.checkFrom([InputKey.EndIp])}
                                    onValueChange={({ detail }) => this.setIpInfo(InputKey.EndIp, detail)}
                                    validateState={single ? ValidateState.OK : validateState.endIP}
                                    validateMessages={ValidateMessages}
                                />
                            </Form.Field>
                        </Form.Row>
                        <Form.Row>
                            <Form.Label>
                                <Radio
                                    value={true}
                                    checked={single}
                                    onChange={({ detail: { value } }) => this.ipRangeSelect(value)}
                                >
                                    <span className={styles['color-strong']} >{__('单个IP：')}</span>
                                </Radio>
                            </Form.Label>
                            <Form.Field>
                                <ValidateBox
                                    key={'ip'}
                                    type={'text'}
                                    role={'sweetui-validatebox'}
                                    width={300}
                                    disabled={!single}
                                    value={single ? ip : ''}
                                    placeholder={isIpv4 ? Placeholder.ipv4 : Placeholder.ipv6}
                                    validator={(input) => this.validateInput(input, isIpv4)}
                                    onBlur={() => this.checkFrom([InputKey.Ip])}
                                    onValueChange={({ detail }) => this.setIpInfo(InputKey.Ip, detail)}
                                    validateState={single ? validateState.ip : ValidateState.OK}
                                    validateMessages={ValidateMessages}
                                />
                            </Form.Field>
                        </Form.Row>
                        <Form.Row>
                            {
                                isIpv4
                                    ? (
                                        <>
                                            <Form.Label>
                                                <span className={styles['net-label']}>
                                                    {__('子网掩码：')}
                                                </span>
                                            </Form.Label>
                                            <Form.Field>
                                                <ValidateBox
                                                    key={'mask'}
                                                    type={'text'}
                                                    role={'sweetui-validatebox'}
                                                    width={300}
                                                    disabled={!single}
                                                    placeholder={Placeholder.mask}
                                                    value={!single ? '' : mask}
                                                    onBlur={() => this.checkFrom([InputKey.Mask])}
                                                    validator={(input) => this.validateInput(input, isIpv4)}
                                                    onValueChange={({ detail }) => this.setIpInfo(InputKey.Mask, detail)}
                                                    validateState={single ? validateState.mask : ValidateState.OK}
                                                    validateMessages={ValidateMessages}
                                                />
                                            </Form.Field>
                                        </>
                                    ) : (
                                        <>
                                            <Form.Label>
                                                <span className={styles['net-label']}>
                                                    {__('前缀长度：')}
                                                </span>
                                            </Form.Label>
                                            <Form.Field>
                                                <ValidateBox
                                                    key={'prefix'}
                                                    type={'text'}
                                                    role={'sweetui-validatebox'}
                                                    width={300}
                                                    disabled={!single}
                                                    value={!single ? '' : prefix}
                                                    placeholder={Placeholder.prefix}
                                                    onBlur={() => this.checkFrom([InputKey.Prefix])}
                                                    validator={(input) => {
                                                        return !Number.isNaN(parseInt(input)) &&
                                                            /^(0|[1-9][0-9]*)$/.test(input)
                                                    }}
                                                    onValueChange={({ detail }) => this.setIpInfo(InputKey.Prefix, detail)}
                                                    validateState={validateState.prefix}
                                                    validateMessages={ValidateMessages}
                                                />
                                            </Form.Field>
                                        </>
                                    )
                            }
                        </Form.Row>
                    </Form>
                    <div className={styles['error']}>
                        {ipError ? __('终止IP不能小于起始IP。') : ''}
                    </div>
                </Panel>
            </ ModalDialog2 >
        )
    }
}