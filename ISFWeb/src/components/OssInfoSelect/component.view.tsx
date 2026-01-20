import * as React from 'react';
import classnames from 'classnames';
import { UIIcon } from '@/ui/ui.desktop';
import { ValidateSelect } from '@/sweet-ui';
import OssInfoSelectBase from './component.base';
import { LocationType } from './helper';
import styles from './styles.view';
import __ from './locale';

export default class OssInfoSelect extends OssInfoSelectBase {

    render() {
        const { ossInfo, ossInfos } = this.state;
        const { width, type, validateState, validateMessages, disabled } = this.props;

        return (
            <div>
                <div className={styles['select']}>
                    <ValidateSelect
                        role={'sweetui-validateselect'}
                        className={classnames({ [styles['select-box']]: type === LocationType.Department || type === LocationType.Organization })}
                        freeze={false}
                        popupZIndex={10000}
                        value={ossInfo && ossInfo.hasOwnProperty('ossId') ? ossInfo.ossId : '-1'}
                        menuWidth={width}
                        selectorWidth={width}
                        validateState={validateState}
                        validateMessages={validateMessages}
                        disabled={disabled}
                        onChange={({ detail }) => this.updateSelectedOss(detail)}
                        onBlur={this.props.onBlur}
                    >
                        {
                            ossInfos.map((oss, index) => (
                                <ValidateSelect.Option
                                    role={'sweetui-validateSelect.option'}
                                    value={oss.ossId}
                                    key={index}
                                >
                                    {disabled ? '' : oss.displayName}
                                </ValidateSelect.Option>
                            ))
                        }
                    </ValidateSelect>
                </div>
                <div className={styles['explain']}>
                    <UIIcon
                        role={'ui-uiicon'}
                        code={'\uf055'}
                        size={16}
                        color={'#999'}
                        title={
                            <div className={styles['explain-content']}>
                                {__('存储位置：用户或文档库指定使用的对象存储服务')}
                            </div>
                        }
                    />
                </div>
            </div>

        )
    }
}