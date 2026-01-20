import React, { useState, useRef, useEffect } from 'react'
import {
    getLogStrategy,
    getPasswordStatus,
    updateLogStrategy,
    updatePasswordStatus,
} from '@/core/apis/console/auditlog'
import { Form } from '@/ui/ui.desktop'
import { message, Modal, Button, Select, TimePicker, InputNumber } from 'antd'
import dayjs, { Dayjs } from 'dayjs'
import __ from './locale'
import styles from './styles.view.css'

type CycleOption = {
    value: CycleUnit;
    label: string;
}

// 转存周期单位
enum CycleUnit {
    Day = 'day',
    Week = 'week',
    Month = 'month',
    Year = 'year',
}

// 转存格式
enum DumpFormat {
    CSV = 'csv',
    XML = 'xml',
}

const dumpFormatOptions = [
    { value: DumpFormat.CSV, label: 'CSV' },
    { value: DumpFormat.XML, label: 'XML' },
]

const encryptionOptions = [
    { value: true, label: __('是') },
    { value: false, label: __('否') },
]

const cycleUnitOptions: CycleOption[] = [
    { value: CycleUnit.Day, label: __('天') },
    { value: CycleUnit.Week, label: __('周') },
    { value: CycleUnit.Month, label: __('月') },
    { value: CycleUnit.Year, label: __('年') },
]

export default function LogPolicy({onCancel}) {
    const [transferCycle, setTransferCycle] = useState<number>(1)
    const [transferCycleUnit, setTransferCycleUnit] = useState<CycleUnit>(CycleUnit.Year)
    const [dumpTime, setDumpTime] = useState<string>('03:00:00')
    const [dumpFormat, setDumpFormat] = useState<DumpFormat>(DumpFormat.CSV)
    const [passwordStatus, setPasswordStatus] = useState<boolean>(false)
    const [isChange, setIsChange] = useState<boolean>(false)

    const prevConfigRef = useRef({
        transferCycle: 1,
        transferCycleUnit: CycleUnit.Year,
        dumpTime: '03:00:00',
        dumpFormat: DumpFormat.CSV,
        passwordStatus: false,
    })

    useEffect(() => {
        async function fetchLogRetentionPeriod() {
            try {
                const {
                    retention_period,
                    retention_period_unit,
                    dump_time,
                    dump_format,
                } = await getLogStrategy()

                const { status } = await getPasswordStatus()

                setTransferCycle(retention_period)
                setTransferCycleUnit(retention_period_unit)
                setDumpTime(dump_time)
                setDumpFormat(dump_format)
                setPasswordStatus(status)

                prevConfigRef.current = {
                    transferCycle: retention_period,
                    transferCycleUnit: retention_period_unit,
                    dumpTime: dump_time,
                    dumpFormat: dump_format,
                    passwordStatus: status,
                }
            } catch ({ description }) {
                description && message.error(description)
            }
        }

        fetchLogRetentionPeriod()
    }, [])

    async function handleSetTransferCycle() {
        try {
            await updateLogStrategy({
                field: ['retention_period', 'retention_period_unit', 'dump_time', 'dump_format'],
                retention_period: transferCycle,
                retention_period_unit: transferCycleUnit,
                dump_time: dumpTime,
                dump_format: dumpFormat,
            })

            await updatePasswordStatus({ status: passwordStatus })

            prevConfigRef.current = {
                transferCycle,
                transferCycleUnit,
                dumpTime,
                dumpFormat,
                passwordStatus,
            }

            message.success(__('保存成功'))
            setIsChange(false)
            onCancel()
        } catch ({ description }) {
            description && message.error(description)
        }
    }

    function handleTransferCycle(cycle: number ) {
        setTransferCycle(cycle)
        setIsChange(true)
    }

    function handleUnitChange(unit: CycleUnit) {
        setTransferCycleUnit(unit)
        setIsChange(true)
    }

    function handleDumpTimeChange(time: Dayjs) {
        setDumpTime(time.format('HH:mm:ss'))
        setIsChange(true)
    }

    function handleDumpFormatChange(format: DumpFormat) {
        setDumpFormat(format)
        setIsChange(true)
    }

    function handlePasswordStatusChange(status: string) {
        setPasswordStatus(status)
        setIsChange(true)
    }

    function handleCancel() {
        setTransferCycle(prevConfigRef.current.transferCycle)
        setTransferCycleUnit(prevConfigRef.current.transferCycleUnit)
        setDumpTime(prevConfigRef.current.dumpTime)
        setDumpFormat(prevConfigRef.current.dumpFormat)
        setPasswordStatus(prevConfigRef.current.passwordStatus)

        setIsChange(false)
        onCancel()
    }

    return (
        <Modal
           centered
           maskClosable={false}
           open={true}
           width={560}
           title="日志转存策略"
           onCancel={onCancel}
           footer={[
                <Button
                    key="submit"
                    type="primary"
                    onClick={handleSetTransferCycle}
                    disabled={!isChange || !transferCycle}
                >
                    {__('保存')}
                </Button>,
                <Button key="cancel" onClick={handleCancel}>
                    {__('取消')}
                </Button>,
            ]}
            getContainer={document.getElementById('isf-web-plugins') as HTMLElement}
        >
            <div className={styles['content']}>
                <p className={styles['text']}>
                    {__(
                        '所有日志记录到了如下设置的时长后，将自动转存为历史日志，且不会被删除',
                    )}
                </p>
                <Form role={'ui-form'}>
                    <Form.Row role={'ui-form.row'}>
                        <Form.Label role={'ui-form.label'} className={styles['label-width']}>
                            {__('日志转存周期：')}
                        </Form.Label>
                        <Form.Field role={'ui-form.field'} className={styles['field-wrapper']} isRequired>
                            <InputNumber
                                value={transferCycle}
                                min={1}
                                max={999999}
                                precision={0}
                                style={{ width: 200 }}
                                onChange={handleTransferCycle}
                            />
                            <Select
                                style={{ width: 80, marginLeft: 10 }}
                                value={transferCycleUnit}
                                onChange={handleUnitChange}
                            >
                                {cycleUnitOptions.map((option) => (
                                    <Select.Option
                                        key={option.value}
                                        value={option.value}
                                        selected={
                                            transferCycleUnit === option.value
                                        }
                                    >
                                        {option.label}
                                    </Select.Option>
                                ))}
                            </Select>
                        </Form.Field>
                    </Form.Row>
                    <Form.Row role={'ui-form.row'}>
                        <Form.Label role={'ui-form.label'} className={styles['label-width']}>
                            {__('日志转存时间：')}
                        </Form.Label>
                        <Form.Field role={'ui-form.field'} isRequired>
                            <TimePicker
                                allowClear={false}
                                format={'HH:mm:ss'}
                                style={{ width: 200 }}
                                value={dayjs(dumpTime, 'HH:mm:ss')}
                                onChange={handleDumpTimeChange}
                            />
                        </Form.Field>
                    </Form.Row>
                    <Form.Row role={'ui-form.row'}>
                        <Form.Label role={'ui-form.label'} className={styles['label-width']}>
                            {__('日志转存格式：')}
                        </Form.Label>
                        <Form.Field role={'ui-form.field'} isRequired>
                            <Select
                               style={{ width: 200 }}
                                value={dumpFormat}
                                onChange={handleDumpFormatChange}
                            >
                                {dumpFormatOptions.map((option) => (
                                    <Select.Option
                                        key={option.value}
                                        value={option.value}
                                        selected={
                                            dumpFormat === option.value
                                        }
                                    >
                                        {option.label}
                                    </Select.Option>
                                ))}
                            </Select>
                        </Form.Field>
                    </Form.Row>
                    <Form.Row role={'ui-form.row'}>
                        <Form.Label role={'ui-form.label'} className={styles['label-width']}>
                            {__('下载时是否加密：')}
                        </Form.Label>
                        <Form.Field role={'ui-form.field'} isRequired>
                            <Select
                                style={{ width: 200 }}
                                value={passwordStatus}
                                onChange={handlePasswordStatusChange}
                            >
                                {encryptionOptions.map((option) => (
                                    <Select.Option
                                        key={option.value}
                                        value={option.value}
                                        selected={
                                            passwordStatus === option.value
                                        }
                                    >
                                        {option.label}
                                    </Select.Option>
                                ))}
                            </Select>
                        </Form.Field>
                    </Form.Row>
                </Form>
            </div> 
        </Modal>
    )
}
