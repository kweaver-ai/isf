import * as React from 'react';
import { ModalDialog2 } from '@/sweet-ui';
import { ProgressBar } from '@/ui/ui.desktop';
import styles from './styles.view';
import __ from './locale';

interface ProgressProps {
    /**
     * 标题
     */
    title?: React.ReactNode;

    /**
     * 进度百分比
     */
    rate: number;

    /**
     * 进度提示
     */
    rateTip?: React.ReactNode;

    /**
     * 备注
     */
    note?: React.ReactNode;
}

const Progress = React.memo(function Progress({ title, rate, rateTip, note }: ProgressProps) {

    if (rate) {
        return (
            <ModalDialog2
                title={title}
                buttons={[]}
            >
                <div className={styles['progress']}>
                    <div className={styles['edit-tip']}>
                        {rateTip}
                    </div>
                    <ProgressBar
                        value={rate}
                        width={350}
                        height={16}
                        progressBackground={'#9abbef'}
                    />
                    {
                        note ?
                            note
                            :
                            <div className={styles['progress-tip']}>
                                {__('注：等待过程中，切勿关闭此页面，关闭后任务会中断。')}
                            </div>
                    }

                </div>
            </ModalDialog2>
        )
    }
    return null
})

export default Progress;