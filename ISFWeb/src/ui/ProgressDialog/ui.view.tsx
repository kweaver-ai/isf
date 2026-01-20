import React from 'react';
import ProgressBar from '../ProgressBar/ui.desktop';
import Panel from '../Panel/ui.desktop';
import Text from '../Text/ui.desktop';
import styles from './styles.desktop';
import __ from './locale';

export default function ProgressDialogView({
    detailTemplate,
    item,
    progress,
    prohandleCancel,
}) {
    return (
        <Panel>
            <Panel.Main>
                <Text className={styles['dialog']}>
                    {
                        detailTemplate(item)
                    }
                </Text>
                <ProgressBar value={progress} />
            </Panel.Main>
            <Panel.Footer>
                <Panel.Button onClick={prohandleCancel}>{__('取消')}</Panel.Button>
            </Panel.Footer>
        </Panel>
    )
}