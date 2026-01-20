import * as React from 'react'
import classnames from 'classnames'
import { ModalDialog2, SweetIcon, Button } from '@/sweet-ui'
import { Panel, Title, UIIcon } from '@/ui/ui.desktop'
import { isBrowser, Browser } from '@/util/browser'
import DomainTree from '../../DomainTree/component.view'
import SyncObjectPickerBase from './component.base'
import styles from './styles.view';
import * as deleteIcon from './assets/delete.png'
import __ from './locale'

// 判断是否为Safari浏览器时，添加空的伪元素，解决Safari浏览器下显示双tooltip
const isSafari = isBrowser({ app: Browser.Safari });

export default class SyncObjectPicker extends SyncObjectPickerBase {
    /**
     * 接口返回的部门路径转换
     */
    convertPath(path: string) {
        let ou = [], dc = [], str = '';

        path.split(',').forEach((item) => {
            if (item.indexOf('OU=') === 0) {
                ou = [item.split('=')[1], ...ou]
            } else {
                dc = [...dc, item.split('=')[1]]
            }
        })

        dc.forEach((item, index) => {
            if (index === dc.length - 1) {
                str = str + item + '/'
            } else {
                str = str + item + '.'
            }
        })

        ou.forEach((item, index) => {
            if (index === ou.length - 1) {
                str = str + item
            } else {
                str = str + item + '/'
            }
        })

        return str
    }

    render() {
        const { data } = this.state,
            {
                domainId,
            } = this.props

        return (
            <div className={styles['container']}>
                <ModalDialog2
                    role={'sweetui-modaldialog2'}
                    zIndex={this.props.zIndex || 50}
                    title={__('选择同步源')}
                    width={616}
                    icons={[{
                        icon: <SweetIcon role={'sweetui-sweeticon'} name="x" size={16} />,
                        onClick: this.cancelAddDep.bind(this),
                    },
                    ]}
                    buttons={[
                        {
                            text: __('确定'),
                            theme: 'oem',
                            onClick: this.confirmAddDep,
                            disabled: !data.length,
                        },
                        {
                            text: __('取消'),
                            theme: 'regular',
                            onClick: this.cancelAddDep,
                        },
                    ]}
                >
                    <Panel role={'ui-panel'}>
                        {
                            <div>
                                <div className={classnames(
                                    styles['content-row'],
                                    styles['row-margin'],
                                )}>
                                    <div>
                                        {__('请选择同步源：')}
                                    </div>
                                    <div className={styles['org-picker-clear']}>
                                        <span>{__('已选：')}</span>

                                        <Button
                                            role={'sweetui-button'}
                                            className={styles['clear-text']}
                                            theme={'text'}
                                            disabled={!data.length}
                                            onClick={this.clearSelectDep}
                                        >
                                            {__('清空')}
                                        </Button>
                                    </div>
                                </div>
                                <div className={styles['content-row']}>
                                    <div className={styles['org-picker-tree']}>
                                        <DomainTree
                                            ref={(ref) => this.ref = ref}
                                            domainId={domainId}
                                            selection={data}
                                            onSelectionChange={this.selectDep}
                                        />
                                    </div>
                                    <div className={styles['org-picker-arrow']}>
                                        <UIIcon
                                            role={'ui-uiicon'}
                                            size={28}
                                            code={'\uf0f5'}
                                            color={'#757575'}
                                            onClick={this.addTreeData}
                                        />
                                    </div>
                                    <div className={styles['org-picker-selections']}>
                                        <ul>
                                            {
                                                data.map((sharer, index) => (
                                                    <li
                                                        key={index}
                                                        style={{ position: 'relative' }}
                                                        className={styles['selection']}
                                                    >
                                                        <Title role={'ui-title'} content={sharer.pathName ? this.convertPath(sharer.pathName) : sharer.name}>
                                                            <div className={classnames(
                                                                styles['dep-name'],
                                                                {
                                                                    [styles['safari']]: isSafari,
                                                                },
                                                            )}>
                                                                <span>{sharer.name}</span>
                                                            </div>
                                                        </Title>
                                                        <UIIcon
                                                            className={styles['icon-del']}
                                                            size={13}
                                                            code={'\uf014'}
                                                            fallback={deleteIcon}
                                                            onClick={() => { this.deleteSelectDep(sharer) }}
                                                        />
                                                    </li>
                                                ))
                                            }
                                        </ul>
                                    </div>
                                </div>
                            </div>
                        }
                    </Panel>
                </ModalDialog2>
            </div>
        )
    }

}