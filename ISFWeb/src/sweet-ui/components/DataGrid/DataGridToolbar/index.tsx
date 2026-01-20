import React from 'react';
import { isFunction } from 'lodash';
import View from '../../View';
import CheckBox from '../../CheckBox';
import { ToolbarComponentProps } from '../../DataTable';
import __ from './locale'
import styles from './styles';

export interface DataGridToolbarProps extends ToolbarComponentProps {
    ToolbarComponent: React.FunctionComponent<any>;

    /**
     * 是否允许全选复选框
     */
    enableSelectAll?: boolean;

    /**
     * 当前是否是全选状态
     */
    isSelectAllChecked?: boolean;
}

const DataGridToolbar: React.FunctionComponent<DataGridToolbarProps> = function DataGridToolbar({
    data = [],
    selection,
    selectAll,
    enableSelectAll,
    isSelectAllChecked,
    ToolbarComponent,
}) {
    return (
        <View className={styles['root']}>
            {
                Array.isArray(selection) && enableSelectAll ? (
                    <View className={styles['select-area']}>
                        <CheckBox
                            className={styles['select-all']}
                            disabled={!data || !data.length}
                            checked={isSelectAllChecked}
                            onChange={(event) => selectAll(event.target.checked)}
                        />
                        {
                            !selection || !selection.length ?
                                <View className={styles['select-all-label']}>{__('全选')}</View> :
                                null
                        }
                    </View>
                ) : null
            }
            <View className={styles['tools']}>
                {
                    isFunction(ToolbarComponent) ?
                        <ToolbarComponent
                            {...{ data, selection }}
                        /> :
                        ToolbarComponent
                }
            </View>
        </View>
    );
};

export default DataGridToolbar;
