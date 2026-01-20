import * as React from 'react';
import { Text, UIIcon, FlexBox } from '@/ui/ui.desktop';
import SetManagerByDep from '../SetManagerByDep/component.view';
import DisplayManagerBase from './component.base';
import styles from './styles.view';
import * as manager from './assets/manager.png';
import * as edit from './assets/edit.png'

export default class DisplayManager extends DisplayManagerBase {
    render() {
        return (
            <div>
                <div className={styles['container']}>
                    <FlexBox role={'ui-flexbox'}>
                        <FlexBox.Item role={'ui-flexbox.item'}>
                            <div>
                                <img src={manager} />
                            </div>
                        </FlexBox.Item>
                        <FlexBox.Item role={'ui-flexbox.item'}>
                            <div className={styles['managers']} >
                                <Text role={'ui-text'}>
                                    {
                                        this.state.manager.length ?
                                            this.state.manager.join(',') :
                                            '---'
                                    }
                                </Text>
                            </div>
                        </FlexBox.Item>
                        <FlexBox.Item role={'ui-flexbox.item'}>
                            <div className={styles['edited']}>
                                {
                                    !this.props.departmentId || this.props.departmentId === '-1' || this.props.departmentId === '-2' || !this.props.hasPermission ?
                                        null :
                                        (
                                            <UIIcon
                                                role={'ui-uiicon'}
                                                code={'\uf01c'}
                                                onClick={this.openEditDialog}
                                                size={20}
                                                fallback={edit}
                                            />
                                        )

                                }
                            </div>
                        </FlexBox.Item>
                    </FlexBox>
                </div>

                {
                    this.state.isEdited ?
                        (
                            <SetManagerByDep
                                departmentId={this.props.departmentId}
                                departmentName={this.props.departmentName}
                                userid={this.props.userid}
                                onCancel={this.cancelEditing}
                                onSetSuccess={this.updateManagers}
                            />
                        ) :
                        null
                }
            </div>
        )
    }
}