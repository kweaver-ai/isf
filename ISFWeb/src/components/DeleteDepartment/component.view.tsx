import * as React from 'react';
import { ModalDialog2, SweetIcon } from '@/sweet-ui';
import { FlexBox } from '@/ui/ui.desktop';
import { Status } from './component.base';
import DeleteDepartmentBase from './component.base';
import __ from './locale';
import styles from './styles.view';
import { Spin } from 'antd';
import { LoadingOutlined } from '@ant-design/icons';

export default class DeleteDepartment extends DeleteDepartmentBase {

    render() {
        const { status } = this.state;

        return (
            <div className={styles['container']}>
                {
                    status === Status.Normal ?
                        <ModalDialog2
                            role={'sweetui-modaldialog2'}
                            width={420}
                            icons={[{
                                icon: <SweetIcon name={'x'} size={16} role={'sweetui-sweeticon'} />,
                                onClick: this.props.onRequestCancelDeleteDep,
                            }]}
                            buttons={[
                                {
                                    text: __('确定'),
                                    theme: 'oem',
                                    onClick: this.confirmDeleteDepartment.bind(this),
                                },
                                {
                                    text: __('取消'),
                                    theme: 'regular',
                                    onClick: this.props.onRequestCancelDeleteDep,
                                },
                            ]}
                        >
                            <FlexBox role={'ui-flexbox'}>
                                <FlexBox.Item align="left top" role={'ui-flexbox.item'}>
                                    <SweetIcon
                                        role={'sweetui-sweeticon'}
                                        name="notice"
                                        size={40}
                                        color="#F39422"
                                    />
                                </FlexBox.Item>
                                <FlexBox.Item className={styles['container']} role={'ui-flexbox.item'}>
                                    <div className={styles['select-dialog']}>
                                        <label>
                                            {
                                                __('此操作将删除部门“${depName}”及其下所有的子部门，该部门下所有的用户成员将会转移至未分配组中，确认要执行此操作吗？', { depName: this.props.dep.name })
                                            }
                                        </label>
                                    </div>
                                </FlexBox.Item>
                            </FlexBox>
                        </ModalDialog2> : null
                }
                {
                    status === Status.Loading ?
                        <Spin size='large' tip={__('正在删除部门，请稍候……')} fullscreen indicator={<LoadingOutlined style={{ color: '#fff' }} spin />}/> :
                        null
                }
            </div>
        )
    }
}