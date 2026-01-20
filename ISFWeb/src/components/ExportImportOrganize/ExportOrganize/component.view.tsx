import * as React from 'react';
import classnames from 'classnames';
import { Dialog2 as Dialog, Panel, FlexBox, Button, Title, UIIcon, Icon, ProgressBar, MessageDialog, Text } from '@/ui/ui.desktop';
import { isBrowser, Browser } from '@/util/browser';
import { NodeType, getDepName } from '@/core/organization';
import OrganizationTree from '../../OrganizationTree/component.view';
import SearchDep from '../../SearchDep/component.desktop';
import ExportOrganizeBase from './component.base';
import __ from './locale';
import styles from './styles.view';

/**
 * 判断是否为Safari浏览器，是时，添加空的伪元素，解决Safari浏览器下显示双tooltip
*/
const isSafari = isBrowser({ app: Browser.Safari });

export default class ExportOrganize extends ExportOrganizeBase {

    getDepname(member: object) {
        /**
         * 点击搜索的用户
         */
        if (member.departmentId) {
            return member.departmentName
            /**
             * 点击树中的用户
             */
        } else {
            return member.origianl.name
        }
    }

    /**
     * 导出的组织树以及已选列表
     */
    getMemberTemplate() {
        const { selectedMember } = this.state;
        const { userid } = this.props;

        return (
            <div className={styles['contain']}>
                <div className={styles['item']}>
                    <div className={styles['item-box']}>
                        <div className={styles['search-box']}>
                            <SearchDep
                                placeholder={__('搜索')}
                                onSelectDep={this.addMember}
                                userid={userid}
                                width={242}
                                selectType={[NodeType.ORGANIZATION, NodeType.DEPARTMENT]}
                            />
                        </div>
                        <div className={styles['organization-tree']}>
                            <OrganizationTree
                                userid={userid}
                                selectType={[NodeType.ORGANIZATION, NodeType.DEPARTMENT]}
                                onSelectionChange={this.addMember}
                            />
                        </div>
                    </div>
                </div>
                <div className={styles['item']}>
                    <div className={styles['item-box']}>
                        <div className={styles['search-box']}>
                            <div className={styles['select-content']}>
                                <FlexBox>
                                    <FlexBox.Item align="left middle">
                                        <label>
                                            {__('已选：')}
                                        </label>
                                    </FlexBox.Item>
                                    <FlexBox.Item align="right middle">
                                        <div>
                                            <Button
                                                onClick={this.clearSelectDep.bind(this)}
                                                disabled={!selectedMember.length}
                                            >
                                                {__('清空')}
                                            </Button>
                                        </div>
                                    </FlexBox.Item>
                                </FlexBox>
                            </div>
                            <div className={classnames(styles['organization-selected'], styles['select-content'])}>
                                <ul>
                                    {
                                        selectedMember.map((member) => (
                                            <li
                                                key={member.id}
                                                style={{ position: 'relative' }}
                                                className={styles['dep-item']}>
                                                <div className={styles['seleted-data']}>
                                                    <Title content={getDepName(member)}>
                                                        <div className={classnames(
                                                            styles['dep-name'],
                                                            {
                                                                [styles['safari']]: isSafari,
                                                            },
                                                        )}>
                                                            {member.name}
                                                        </div>
                                                    </Title>
                                                </div>
                                                <div className={styles['selected-data-del']}>
                                                    <UIIcon
                                                        size={13}
                                                        code={'\uf014'}
                                                        onClick={() => { this.deleteSelectDep(member) }}
                                                    />
                                                </div>
                                            </li>
                                        ))
                                    }
                                </ul>
                            </div>
                        </div>
                    </div>
                </div>
            </div>
        )
    }

    render() {
        const { exportStatus, selectedMember, progress, unExistMember } = this.state;
        const { onCancel } = this.props;

        return (
            <div>
                {
                    exportStatus ? null :
                        < Dialog
                            title={__('选择导出用户组织')}
                            onClose={onCancel}
                        >
                            <Panel>
                                <Panel.Main>
                                    <div className={styles['select-range-size']}>
                                        {
                                            this.getMemberTemplate()
                                        }
                                    </div>
                                </Panel.Main>
                                <Panel.Footer>
                                    <Panel.Button
                                        theme='oem'
                                        onClick={this.onSaveMember}
                                        disabled={!selectedMember.length}
                                    >
                                        {
                                            __('导出')
                                        }
                                    </Panel.Button>
                                    <Panel.Button onClick={onCancel}>
                                        {
                                            __('取消')
                                        }
                                    </Panel.Button>
                                </Panel.Footer>
                            </Panel>
                        </Dialog >
                }
                {
                    !exportStatus ?
                        null
                        :
                        <Dialog
                            title={__('导出用户组织')}
                            onClose={this.cancelExport.bind(this)}
                        >
                            <Panel>
                                <Panel.Main>
                                    <div className={styles['export-organize-progress']}>
                                        <div className={styles['organize-progress']}>
                                            <ProgressBar
                                                value={progress}
                                                height={20}
                                            />
                                            <span className={styles['organize-progress-status']}>
                                                {progress === 1 ? __('完成导出') : __('正在导出选中的用户组织，请稍后......')}
                                            </span>
                                        </div>
                                    </div>
                                </Panel.Main>
                                <Panel.Footer>
                                    <Panel.Button
                                        onClick={this.downloadFile.bind(this)}
                                        style={{ padding: '0 22px' }}
                                        disabled={progress === 1 ? false : true}
                                    >
                                        {__('下载已导出的文件')}
                                    </Panel.Button>
                                </Panel.Footer>
                            </Panel>
                        </Dialog>
                }
                {
                    unExistMember.length > 0 ?
                        <MessageDialog
                            onConfirm={this.closeErrorMessage}
                        >
                            <div className={styles['message-info']}>
                                {
                                    __('“')
                                }
                            </div>
                            <div className={styles['message-info-name']}>
                                <Text>
                                    {
                                        unExistMember.join(',')
                                    }
                                </Text>
                            </div>
                            <div className={styles['message-info']}>
                                {
                                    __('”部门已不存在')
                                }
                            </div>
                        </MessageDialog>
                        : null
                }
            </div>
        )
    }

}