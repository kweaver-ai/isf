import * as React from 'react';
import FlexBox from '@/ui/FlexBox/ui.desktop';
import Button from '@/ui/Button/ui.desktop';
import SelectUserOrDep from '../SelectUserOrDep/component.view';
import SetUsersFreezeStatusBase from './component.base';
import ShowSelected from './ShowSelected/component.view';
import ShowErrorMessage from './ShowErrorMessage/component.view';
import Dialog from '@/ui/Dialog2/ui.desktop';
import Panel from '@/ui/Panel/ui.desktop';
import { Status } from './helper';
import __ from './locale';
import styles from './styles.view';
import { Spin } from 'antd';
import { LoadingOutlined } from '@ant-design/icons';

export default class SetUsersFreezeStatus extends SetUsersFreezeStatusBase {

    render() {
        return (
            <div>
                {
                    this.state.status === Status.NORMAL ?
                        (
                            <Dialog
                                title={this.props.freezeStatus ? __('冻结用户') : __('解冻用户')}
                                buttons={[]}
                            >
                                <Panel>
                                    <Panel.Main>
                                        <FlexBox>
                                            <FlexBox.Item>
                                                <div className={styles['selected-tip']}>
                                                    {
                                                        this.props.freezeStatus ?
                                                            __('请将您想要冻结的用户添加到右侧列表：') :
                                                            __('请将您想要解冻的用户添加到右侧列表：')
                                                    }
                                                </div>
                                                <div className={styles['select-tree']}>
                                                    <SelectUserOrDep
                                                        userid={this.props.userid}
                                                        onSelected={(selected) => { this.getSelected(selected) }}
                                                    />
                                                </div>
                                            </FlexBox.Item>
                                            <FlexBox.Item>
                                                <div className={styles['select-button-lab']}>
                                                    <Button
                                                        onClick={this.clearSelected.bind(this)}
                                                        className={styles['select-button']}
                                                        disabled={!this.state.selecteds.length}
                                                    >
                                                        {
                                                            __('清空')
                                                        }
                                                    </Button>
                                                </div>
                                                <div className={styles['selected-data']}>
                                                    <ShowSelected
                                                        selecteds={this.state.selecteds}
                                                        onDelete={this.deleteSelected.bind(this)}
                                                    />
                                                </div>
                                            </FlexBox.Item>
                                        </FlexBox>
                                    </Panel.Main>
                                    <Panel.Footer>
                                        <Panel.Button
                                            onClick={this.setFreezeStatus.bind(this)}
                                            disabled={!this.state.selecteds.length}
                                        >
                                            {__('确定')}
                                        </Panel.Button>
                                        <Panel.Button onClick={this.props.onComplete} >
                                            {__('取消')}
                                        </Panel.Button>
                                    </Panel.Footer>
                                </Panel>
                            </Dialog>
                        ) :
                        null
                }

                {
                    this.state.status === Status.LOADING ?
                        this.props.freezeStatus ?
                            <Spin size='large' tip={__('正在冻结用户，请稍后……')} fullscreen indicator={<LoadingOutlined style={{ color: '#fff' }} spin />}/> :
                            <Spin size='large' tip={__('正在解冻用户，请稍后……')} fullscreen indicator={<LoadingOutlined style={{ color: '#fff' }} spin />}/> :
                        null
                }

                {
                    this.state.status === Status.ERROR || this.state.status === Status.CURRENT_USER_INCLUDED ?
                        <ShowErrorMessage
                            errorType={this.state.status}
                            errorInfo={this.errorInfo}
                            freezeStatus={this.props.freezeStatus}
                            onConfirm={this.props.onComplete}
                        /> :
                        null

                }
            </div>
        )
    }
}