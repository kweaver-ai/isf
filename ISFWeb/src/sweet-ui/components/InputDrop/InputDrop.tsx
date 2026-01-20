import React from 'react'
import classnames from 'classnames';
import ValidateBox from '../ValidateBox';
import SweetIcon from '../SweetIcon';
import PopMenu from '../PopMenu'
import Select2 from '../Select2'
import styles from './styles'

interface InputDropProps {
    /**
     * role
     */
    role?: string;

    /**
     * 节点
     */
    element?: any;

    /**
     * 宽度
     */
    width: string;

    /**
     * value值
     */
    value: string;

    /**
     * 占位符
     */
    placeholder?: string;

    /**
     * 验证状态
     */
    validateState: any;

    /**
     * 验证信息
     */
    validateMessages: {
        [key: string]: string;
    };

    /**
     * 是否禁用
     */
    disabled?: boolean;

    /**
     * 下拉框选项
     */
    options: Array<number | string>;

    /**
     * 输入验证函数
     */
    validator?: (value: string) => boolean;

    /**
     * 值改变函数
     */
    onValueChange: (value: string) => void;

    /**
     * 失焦函数
     */
    onBlur?: () => void;
}

const IntputDrop: React.FunctionComponent<InputDropProps> = React.memo(({
    role,
    width,
    value,
    placeholder,
    validateState,
    validateMessages,
    disabled,
    options,
    element,
    validator,
    onValueChange,
    onBlur,
}) => {

    return (
        <div className={styles['input-drop']} >
            <ValidateBox
                role={role}
                width={width}
                value={value}
                placeholder={placeholder}
                validator={validator}
                onValueChange={({ detail }) => onValueChange(detail)}
                validateMessages={validateMessages}
                validateState={validateState}
                disabled={disabled}
                onBlur={onBlur}
            />
            <PopMenu
                freeze={false}
                anchorOrigin={['left', 'bottom']}
                alignOrigin={['left', 'top']}
                triggerEvent={'click'}
                element={element}
                onRequestCloseWhenClick={(close) => close()}
                trigger={({ setPopupVisibleOnClick }) =>
                    <div>

                        <div onClick={() => { disabled ? undefined : setPopupVisibleOnClick() }}>
                            <SweetIcon
                                role={role}
                                className={classnames(
                                    styles['arrow'],
                                    {
                                        [styles['arrow-active']]: validateState in validateMessages,
                                    },
                                )}
                                name={'arrowDown'}
                                size={16}
                            />
                        </div>
                    </div>
                }
            >
                <div style={{ width: width ? width : 200 }}>
                    {
                        options.map((option) => (
                            <Select2.Option
                                role={role}
                                key={option}
                                value={option}
                                selected={value === option}
                                disabled={disabled}
                                onClick={() => onValueChange(String(option))}
                            >
                                {option}
                            </Select2.Option>
                        ))
                    }
                </div>
            </PopMenu >
        </div >
    )
})

export default IntputDrop