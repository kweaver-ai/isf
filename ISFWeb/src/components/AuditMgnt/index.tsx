import * as React from 'react';
import { useState, useCallback, useRef, useContext, useEffect } from 'react';
import { includes } from 'lodash';
import classNames from 'classnames';
import { Text } from '@/ui/ui.desktop';
import { isEqual } from 'lodash';
import DataReport from './DataReport';
import {
    BizGroupItem,
    DefaultSelectedDataReportInfo,
} from './type';
import styles from './styles.view.css';
import __ from './locale';
import AppConfigContext from '@/core/context/AppConfigContext';
import { Button } from 'antd';
import { SystemRoleType } from '@/core/role/role';
import session from '@/util/session';
import LogPolicyIcon from '../../icons/logpolicy.svg'
import LogPolicy from './RetentionLog';

const AuditMgnt: React.FC = () => {
    const dataReportRef = useRef<{ reloadPage: () => Promise<void> }>(null);
    const { oemColor } = useContext(AppConfigContext);
    const [hoverIndex, setHoverIndex] = useState(0)
    const showPolicyBtn = session.get('isf.userInfo')?.user?.roles?.some(({ id }) => [SystemRoleType.Supper, SystemRoleType.Securit].includes(id))
    const [showPolicyConfig, setShowPolicyConfig] = useState(false);
    const source = [
        {
            "id": 'audit-log',
            "name": __("审计日志"),
            "children": [
                {
                    "id": "operation",
                    "name": __('操作日志'),
                },
                {
                    "id": "management",
                    "name": __('管理日志'),
                },
                {
                    "id": "login",
                    "name": __('访问日志'),
                },
                {
                    "id": "history-operation",
                    "name": __('历史操作日志'),
                },
                {
                    "id": "history-management",
                    "name": __('历史管理日志'),
                },
                {
                    "id": "history-login",
                    "name": __('历史访问日志'),
                },
            ]
        }
    ]
    const [bizGroupList] = useState<BizGroupItem[]>(source);
    const [expandedBizGroups, setExpandeBizGroups] = useState<ReadonlyArray<string>>(["audit-log"]);
    const [selectedDataReportInfo, setSelectedDataReportInfo] = useState<{
        parentBizGroupId: string;
        parentBizGroupName: string;
        selectedDataReportId: string;
        selectedDataReportName: string;
    }>(DefaultSelectedDataReportInfo);
    const {
        selectedDataReportId,
        selectedDataReportName,
    } = selectedDataReportInfo;

    const updateExpandedBizGroups = useCallback(async (bizGroupId: string) => {
        if (includes(expandedBizGroups, bizGroupId)) {
            setExpandeBizGroups(expandedBizGroups.filter((id) => id !== bizGroupId));
        } else {
            setExpandeBizGroups([
                ...expandedBizGroups,
                bizGroupId,
            ]);
        }
    }, [expandedBizGroups, bizGroupList]);

    return (
        <div className={styles['container']}>
            <div className={styles['content']}>
                <div className={styles['biz-group-container']}>
                    <div className={styles['biz-group-list']} style={{ height: showPolicyBtn ? 'calc(100% - 102px)': '100%' }}>
                        {
                            bizGroupList.map(({ id: bizGroupId, name: bizGroupName, children }) => {
                                return (
                                    <div
                                        key={bizGroupId}
                                        className={styles['biz-group-item']}
                                    >
                                        <div
                                            onClick={() => { updateExpandedBizGroups(bizGroupId) }}
                                            className={classNames(
                                                styles['biz-group-item-header'],
                                                {
                                                    [styles['up-arrow']]: includes(expandedBizGroups, bizGroupId),
                                                },
                                            )} 
                                        >
                                            <div className={styles['biz-group-item-header-name']}>
                                                <Text>{bizGroupName}</Text>
                                            </div>
                                        </div>
                                        <div className={styles['biz-group-item-content']}>
                                            <div
                                                className={styles['content-list']}
                                                style={{ display: includes(expandedBizGroups, bizGroupId) ? 'block' : 'none' }}
                                            >
                                                {
                                                    children.map(({ id: dataReportId, name }, index) => {
                                                        return (
                                                            <div
                                                                key={dataReportId}
                                                                className={classNames(
                                                                    styles['list-item'],
                                                                    {
                                                                        [styles['selected']]: isEqual(dataReportId, selectedDataReportId),
                                                                    },
                                                                )}
                                                                style={{
                                                                    backgroundColor:
                                                                    isEqual(dataReportId, selectedDataReportId) || index === hoverIndex
                                                                        ? oemColor.colorPrimaryBg
                                                                        : "transparent",
                                                                }}  
                                                                onClick={() => {
                                                                    setSelectedDataReportInfo({
                                                                        parentBizGroupId: bizGroupId,
                                                                        parentBizGroupName: bizGroupName,
                                                                        selectedDataReportId: dataReportId,
                                                                        selectedDataReportName: name,
                                                                    });
                                                                }}
                                                                onMouseEnter={() => {
                                                                    setHoverIndex(index)
                                                                }}
                                                                onMouseLeave={() => {
                                                                    setHoverIndex(undefined)
                                                                }}
                                                            >
                                                                <div className={styles['item-name']}>
                                                                    <Text>{name}</Text>
                                                                </div>
                                                            </div>
                                                        )
                                                    })
                                                }
                                            </div>
                                        </div>
                                    </div>
                                )
                            })
                        }
                    </div>
                    {
                        showPolicyBtn &&
                        <div className={styles["policy-btn"]}>
                            <Button 
                                className={styles["btn"]}
                                type={"default"} 
                                icon={<LogPolicyIcon style={{ width: 14, height: 14 }}/>}
                                onClick={() => {
                                    setShowPolicyConfig(true);
                                }}
                            >
                                {__('日志转存策略')}
                            </Button>
                        </div>    
                    }
                </div>
                <div className={styles['data-report-container']}>
                    <DataReport
                        ref={dataReportRef}
                        dataReportInfo={{
                            id: selectedDataReportId,
                            name: selectedDataReportName,
                        }}
                    />
                </div>
            </div>
            {
                showPolicyConfig && 
                <LogPolicy onCancel={() => {
                    setShowPolicyConfig(false);
                }}/>
            }
        </div>
    )
}

export default AuditMgnt;