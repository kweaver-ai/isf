import React, { useState, useEffect } from "react"
import intl from "react-intl-universal"
import { Checkbox, DatePicker, Form, Select, Tooltip } from "antd"
import { foreverExpire, formatPerm } from "../util"
import dayjs from "dayjs"
import TipIcon from "../../../icons/tip.svg"
import OperationPolicyIcon from "../../../icons/operation-policy.svg"
import styles from "./styles.css"
import { apis, components } from "@dip/components/dist/dip-components.min.js";
import { ObligationType, OperationObligationType, PolicyConfigListType } from "@/core/apis/console/authorization/type"
import { getObligationsType } from "@/core/apis/console/authorization"

export interface PermType {
    id: string;
    name?: string;
    description?: string;
    obligations?: ObligationType[];
}

export interface PermConfigProps {
    setPermPolicyRef: any;
    resourceType: any;
    operationType: string;
    operationConfig: PolicyConfigListType[];
    getValue: (obj: {operationType, operation: { allow: PermType[]; deny: PermType[]; }; expires_at: string}) => void; 
    permTip: string;
    setPermTip: (tip: string) => void;
    expireTip: string;
    setExpireTip: (tip: string) => void;
    curConfig?: { operation: { allow: PermType[]; deny: PermType[]; }; expires_at: string};
}

export const PermConfig = ({ setPermPolicyRef, resourceType, operationType, operationConfig, getValue, curConfig, permTip, setPermTip, expireTip, setExpireTip }: PermConfigProps) => {
    const [expires, setExpires] = useState(foreverExpire)
    const [permConfig, setPermConfig] = useState(operationConfig || []);
    const [obligationTypes, setObligationTypes] = useState([])

    // 允许权限处理
    const handleAllowChange = (e, item) => {
        const config = permConfig.map((cur: PolicyConfigListType) => {
            if (cur.id === item.id) {
                cur.allow = e.target.checked;
                if (e.target.checked) {
                    cur.deny = false;
                }
            }
            return cur;
        });
        setPermConfig(config);
        getValue({operationType, operation: formatPermValue(permConfig), expires_at: expires})
        setPermTip("");
    }

    // 拒绝权限处理
    const handleDenyChange = (e, item) => {
        const config = permConfig.map((cur: PolicyConfigListType) => {
            if (cur.id === item.id) {
                cur.deny = e.target.checked;
                if (e.target.checked) {
                    cur.allow = false;
                }
            }
            return cur;
        });
        setPermConfig(config);
        getValue({operationType, operation: formatPermValue(permConfig), expires_at: expires})
        setPermTip("");
    }

    // 重置权限
    const resetPerm = value => {
        if (!value) {
            const config = permConfig.map(item => ({
                ...item,
                allow: false,
                deny: false,
                obligations: [],
            }));
            setPermConfig(config);
            getValue({operationType, operation: { allow: [], deny: [] }, expires_at: expires})
        }
    }

    const handleExpiresChange = (date: dayjs.Dayjs | null) => {
        const time = date ? date.format('YYYY-MM-DDTHH:mm:ssZZ').replace(/([+-]\d{2})(\d{2})/, '$1:$2') : foreverExpire;
        setExpires(time);
        setExpireTip("");
        getValue({operationType, operation: formatPermValue(permConfig), expires_at: time})
    }

    // 设置权限策略
    const setPermPolicy = (item) => {
        const unmount = apis.mountComponent(
            components.OpertationPolicy,
            {
                resourceType: resourceType?.id,
                perm: item,
                getPermPolicy,
                onCancel: () => {
                    unmount();
                },
            },
            setPermPolicyRef.current,
        );
    }

    //获取已设置策略的权限
    const getPolicyPerm = (permConfig: PolicyConfigListType[]) => {
        return permConfig
            .filter((item: PolicyConfigListType) => item.obligations && item.obligations.length)
            .map(item => item.name);
    };
    
    //获取权限策略
    const getPermPolicy = (data: { id: string; obligations: any }) => {
        const config = permConfig.map(item => {
            if (item.id === data.id) {
                item.obligations = data.obligations;
            }
            return item;
        });
        setPermConfig(config);
        getValue({operationType, operation: formatPermValue(config), expires_at: expires})
    };
    
    //获取已设置的权限值
    const getPermValue = (permConfig: PolicyConfigListType[]) => {
        const allow = permConfig.filter((item: PolicyConfigListType) => item.allow);
        const deny = permConfig.filter((item: PolicyConfigListType) => item.deny);
        return { allow, deny };
    };
    
    //格式化为接口可用的权限值
    const formatPermValue = (permConfig: PolicyConfigListType[]) => {
        const allow = permConfig.filter(item => item.allow).map(item => ({ id: item.id, obligations: item.obligations }));
        const deny = permConfig.filter(item => item.deny).map(item => ({ id: item.id }));
        return { allow, deny };
    };

    useEffect(() => {
        getValue({operationType, operation: { allow: curConfig?.operation?.allow || [], deny: curConfig?.operation?.deny || [] }, expires_at: curConfig?.expires_at || foreverExpire})
        setPermTip("")
        setExpireTip("")
    }, []);

    // 获取当前资源类型的义务配置
    const getObligationConfig = async() => {
        try{
            const ids = operationConfig.map(item => item.id);
            const data = await getObligationsType({ resource_type_id: resourceType?.id, operation_ids: ids });
            const obligation_types = data.filter((item: OperationObligationType) => item?.obligation_types?.length);
            setObligationTypes(obligation_types);
        }catch(err) {
            setObligationTypes([]);
        }
    }

    // 初始化权限配置
    const initPermConfig = (operationConfig: PolicyConfigListType[], obligationTypes: OperationObligationType[]) => {
        if (operationConfig) {
            if(curConfig) {
                const permConfig = operationConfig.map(item => ({
                    ...item,
                    allow: curConfig?.operation?.allow.some((cur: PermType) => cur.id === item.id),
                    deny: curConfig?.operation?.deny.some((cur: PermType) => cur.id === item.id),
                    obligations: curConfig?.operation?.allow.find((cur: PermType) => cur.id === item.id)?.obligations || [],
                    obligation_types: obligationTypes.find((cur: OperationObligationType) => cur.operation_id === item.id)?.obligation_types || [],
                }));
                setPermConfig(permConfig)
                setExpires(curConfig?.expires_at || foreverExpire)

            }else {
                const permConfig = operationConfig.map(item => ({
                    ...item,
                    allow: false,
                    deny: false,
                    obligations: [],
                    obligation_types: obligationTypes.find((cur: OperationObligationType) => cur.operation_id === item.id)?.obligation_types || [],
                }));
                setPermConfig(permConfig)
                setExpires(foreverExpire)
            }
        }else {
            setPermConfig([])
        }
    }
    
    useEffect(() => { 
        if (operationConfig.length) {
            getObligationConfig();
        }
    }, [operationConfig]);

    useEffect(() => {
        initPermConfig(operationConfig, obligationTypes);
    }, [operationConfig, obligationTypes])

    return (
        <div className={styles["perm-config"]} id="perm-config">
            <Form>
                <Form.Item
                    label={intl.get('edit.operation.perm')}
                    required
                    className={styles['required-label']}
                    validateStatus={permTip ? 'error' : ''}
                    help={permTip ? <div>{permTip}</div> : ''}
                >
                    <div className={styles["perm-set"]}>
                        <Select
                            style={{ width: 'calc(100% - 32px)' }}
                            dropdownStyle={{ width: obligationTypes.length ? 526 : 420 }}
                            placeholder={intl.get('set.perm')}
                            allowClear={true}
                            onChange={resetPerm}
                            value={!getPermValue(permConfig).allow.length && !getPermValue(permConfig).deny.length
                                ? null
                                : formatPerm({ operation: getPermValue(permConfig) }, permConfig)}
                            getPopupContainer={() => document.getElementById('perm-config') as HTMLElement}
                            dropdownRender={() => {
                                return (
                                    <div className={styles['perm']} onClick={(e) => e.stopPropagation()}>
                                        <div className={styles['title']}>
                                            <div className={styles['title-label']}>{intl.get('permission')}</div>
                                            <div className={styles['title-box']}>{intl.get('allow')}</div>
                                            <div className={styles['title-box']}>{intl.get('deny')}</div>
                                            {
                                                obligationTypes.length ? (
                                                    <div className={styles['title-policy']}>
                                                        <span>{intl.get('operation.policy')}</span>
                                                        <div 
                                                            className={styles['policy-icon']}
                                                            title={[
                                                                intl.get('operation.policy.tip.first'),
                                                                intl.get('operation.policy.tip.second'),
                                                                intl.get('operation.policy.tip.third'),
                                                            ].join('\n')}>
                                                            <TipIcon style={{ width: '14px', height: '14px', color: 'rgba(0,0,0,.45)', marginLeft: '8px' }} />
                                                        </div>
                                                    </div>
                                                ) : null
                                            }
                                        </div>
                                        {permConfig.map(item => {
                                            return (
                                                <div className={styles['permission']} key={item.id}>
                                                    <div className={styles["label"]}>
                                                        <div className={styles["name"]} title={item.name}>
                                                            {item.name}
                                                        </div>
                                                        {
                                                            item.description ? 
                                                                <div className={styles["tip"]}>
                                                                    <Tooltip placement='top' title={item.description}><TipIcon style={{width: "14px", height: "14px", color: "rgba(0,0,0,.45)"}}/></Tooltip>
                                                                </div> : null
                                                        }
                                                    </div>
                                                    <div className={styles['checkbox']}>
                                                        <Checkbox
                                                            checked={item.allow}
                                                            onChange={e => handleAllowChange(e, item)}
                                                        />
                                                    </div>
                                                    <div className={styles['checkbox']}>
                                                        <Checkbox
                                                            checked={item.deny}
                                                            onChange={e => handleDenyChange(e, item)}
                                                        />
                                                    </div>
                                                    {
                                                        item?.obligation_types?.length ? (
                                                            <div className={styles['policy']}>
                                                                <a
                                                                    disabled={Boolean(!item.allow || item.deny)}
                                                                    onClick={() => {
                                                                        if (item.allow) {
                                                                            setPermPolicy(item)
                                                                        }
                                                                    }}
                                                                    title={!item.allow || item.deny ? intl.get('set.operation.policy.tip') : ''}
                                                                >
                                                                    {intl.get(item.obligations?.length ? 'setted' : 'set')}
                                                                </a>
                                                            </div>
                                                        ) : null
                                                    }
                                                </div>
                                            );
                                        })}
                                    </div>
                                );
                            }}
                        />
                        {getPolicyPerm(permConfig).length ? (
                            <div
                                className={styles['policy-icon']}
                                title={intl.get('perm.policy.tip', { perm: getPolicyPerm(permConfig).join('、') })}
                            >
                                <OperationPolicyIcon style={{ width: 16, height: 16 }} />
                            </div>
                        ) : null}
                    </div>
                </Form.Item>
                <Form.Item
                    label={intl.get('edit.expires')}
                    className={styles['label']}
                    validateStatus={expireTip ? 'error' : ''}
                    help={
                        expireTip ? (
                            <div style={{ width: '100%'}}>
                                {expireTip}
                            </div>
                        ) : (
                            ''
                        )
                    }
                >
                    <div onClick={(e) => e.stopPropagation()}>
                        <DatePicker
                            style={{ width: 'calc(100% - 32px)' }}
                            allowClear={true}
                            showNow={false}
                            showTime={{ format: 'HH:mm' }}
                            format="YYYY/MM/DD HH:mm"
                            value={expires === foreverExpire ? null : dayjs(expires)}
                            placeholder={expires === foreverExpire ? intl.get('forever.expire') : ''}
                            disabledDate={current => {
                                return current < dayjs().add(-1, 'day');
                            }}
                            onChange={handleExpiresChange}
                            getPopupContainer={trigger => trigger.parentNode as HTMLDivElement}
                        />
                    </div>
                </Form.Item>
            </Form>
        </div>
    )
}