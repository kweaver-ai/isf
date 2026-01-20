// eslint-disable-next-line spaced-comment
/// <reference path="./component.base.d.ts" />

import * as React from 'react';
import FlexBox from '@/ui/FlexBox/ui.desktop';
import UIIcon from '@/ui/UIIcon/ui.desktop';
import Text from '@/ui/Text/ui.desktop';
import styles from './styles.view';

export default function ShowSelected({ selecteds, onDelete }: Components.SetUsersFreezeStatus.ShowSelected.Props) {
    return (
        <div>
            {
                selecteds.map((selected) => {
                    return (
                        <div key={selected.id} className={styles['selected-item']}>
                            <FlexBox>
                                <FlexBox.Item>
                                    <div className={styles['selected-name']}>
                                        <Text>{selected.name}</Text>
                                    </div>
                                </FlexBox.Item>
                                <FlexBox.Item>
                                    <div className={styles['selected-parent-name']}>
                                        <Text>
                                            {
                                                selected.parentName
                                            }
                                        </Text>
                                    </div>
                                </FlexBox.Item>
                                <FlexBox.Item align="right">

                                    <div className={styles['selected-delete']}>
                                        <UIIcon size={13} code={'\uf014'} onClick={() => { onDelete(selected) }} />
                                    </div>
                                </FlexBox.Item>

                            </FlexBox>
                        </div>
                    )
                })
            }
        </div>
    )
}