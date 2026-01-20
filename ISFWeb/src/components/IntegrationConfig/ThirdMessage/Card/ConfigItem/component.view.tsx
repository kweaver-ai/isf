import * as React from 'react';
import { findIndex, noop } from 'lodash';
import SelectMenu2 from '@/ui/SelectMenu2/ui.desktop';
import { TextBox, UIIcon } from '@/ui/ui.desktop';
import { ValidateBox } from '@/sweet-ui';
import { Switch } from '@/sweet-ui';
import { NumberBox } from '@/sweet-ui';
import { ValidateMessage, Types, ConfigItem as Config } from '../../helper';
import styles from './styles.view.css';
import __ from './locale';

const typeList = [
    { name: __('字符'), key: Types.StringType },
    { name: __('开关'), key: Types.BooleanType },
    { name: __('数字'), key: Types.NumberType },
]

interface ConfigItemProps {
    /**
     * 一条配置信息
     */
    configItem: Config;

    /**
     * 配置项名称改变
     * @param {*} value 输入的值，名称全部是字符串类型
     */
    onRequestConfigNameChange: (value: string) => void;

    /**
     * 配置项类型改变
     * @param {*} type 选择的类型
     */
    onRequestConfigTypeChange: (type: Types) => void;

    /**
     * 配置项值改变
     * @param {*} value 输入的值，根据选择的类型来定输入值的类型
     * @param {*} type 选择的类型
     */
    onRequestConfigValueChange: (value: any, type: Types) => void;

    /**
     * 删除配置项
     */
    onRequestDelete: () => void;
}

// 完全受控组件，所有的状态都从父组件接收
const ConfigItem: React.FunctionComponent<ConfigItemProps> = function ConfigItem({
    configItem,
    onRequestConfigNameChange = noop,
    onRequestConfigTypeChange = noop,
    onRequestConfigValueChange = noop,
    onRequestDelete = noop,
}) {
    return (
        <div className={styles['container']}>
            {
                <div className={styles['item-name']}>
                    <div className={styles['tips']}>{__('名称：')}</div>
                    <div className={styles['mark']}>{'*'}</div>
                    <div className={styles['validate']}>
                        <ValidateBox
                            // ref={(ref) => this.validateBox = ref}
                            width={160}
                            value={configItem.configName}
                            validateState={configItem.nameValidateStatus}
                            validateMessages={ValidateMessage}
                            onValueChange={({ detail: value }) => onRequestConfigNameChange(value)}
                        />
                    </div>
                </div>
            }
            {
                <div className={styles['item-type']}>
                    <div className={styles['tips']}>{__('类型：')}</div>
                    <SelectMenu2
                        candidateItems={typeList}
                        selectValue={typeList[findIndex(typeList, (item) => item.key === configItem.configType)]}
                        className={styles['select-menu']}
                        onSelect={(type) => onRequestConfigTypeChange(type.key)}
                        numberOfChars={20}
                    />
                </div>
            }

            <div className={styles['item-value']}>
                <div className={styles['tips']}>{__('值：')}</div>
                {
                    configItem.configType === Types.StringType
                        ?
                        // 只做是否为空效验
                        <div className={styles['validate']}>
                            <ValidateBox
                                width={300}
                                value={configItem.configValue}
                                validateState={configItem.valueValidateStatus}
                                validateMessages={ValidateMessage}
                                onValueChange={({ detail: value }) => onRequestConfigValueChange(value, Types.StringType)}
                            />
                        </div>
                        :
                        configItem.configType === Types.NumberType
                            ?
                            <div className={styles['validate']}>
                                <NumberBox
                                    width={300}
                                    value={configItem.configValue}
                                    onValueChange={({ detail }) => onRequestConfigValueChange(detail, Types.NumberType)}
                                />
                            </div>
                            :
                            configItem.configType === Types.BooleanType
                                ?
                                <div className={styles['switch']}>
                                    <Switch
                                        checked={configItem.configValue}
                                        onChange={({ detail }) => onRequestConfigValueChange(detail, Types.BooleanType)}
                                    />
                                </div>
                                :
                                <div className={styles['validate']}>
                                    <TextBox
                                        width={300}
                                        disabled={true}
                                        placeholder={__('未知类型无法显示，您可以进入高级配置查看')}
                                    />
                                </div>
                }
            </div>
            <div className={styles['item-delete']}>
                {
                    !configItem.defaultConfig
                        ?
                        <UIIcon
                            code={'\uf014'}
                            size={12}
                            color={'#505050'}
                            className={styles['uiicon']}
                            onClick={() => onRequestDelete()}
                        />
                        : null
                }
            </div>
        </div>
    )
}

export default ConfigItem