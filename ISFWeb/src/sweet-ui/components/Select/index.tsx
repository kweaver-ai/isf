import React from 'react';
import classnames from 'classnames';
import { reduce } from 'lodash';
import Text from '@/ui/Text/ui.desktop';
import { createEventDispatcher, SweetUIEvent } from '../../utils/event';
import View from '../View';
import Trigger from '../Trigger';
import SweetIcon from '../SweetIcon';
import Menu from '../Menu';
import SelectOption from './Option';
import AppConfigContext from '@/core/context/AppConfigContext';
import styles from './styles';

interface SelectProps {
    /**
     * 禁用状态
     */
    disabled?: boolean;

    /**
     * width，包含盒模型的padding和border
     */
    width?: number;

    /**
     * 下拉菜单最大高度
     */
    maxMenuHeight?: number;

    /**
     * 显示内容
     */
    value: string;

    /**
     * className
     */
    className?: string;

    placeholder?: string; // TODO

    /**
     * 是否只读
     */
    readOnly?: boolean;

    /**
     * 弹出层展开时是否冻结滚动条
     */
    freeze?: boolean;

    /**
     *  下拉框层级
     */
    popupZIndex?: number;

    /**
     * 选中时触发
     */
    onChange: (event: SweetUIEvent<string>) => void;

    /**
     * 下拉选项展开状态变化时触发
     */
    onPopupVisibleChange: (event: SweetUIEvent<boolean>) => void;

    role?: string;
}

interface SelectState {
    value: any;
    visiable: boolean;
    label: string | React.ReactNode;
}

class Select extends React.Component<SelectProps, SelectState> {
    static contextType = AppConfigContext;

    static defaultProps = {
        anchorOrigin: ['left', 'bottom'],
        alignOrigin: ['left', 'top'],
        disabled: false,
        freeze: true,
        value: '',
        width: 120,
        maxMenuHeight: 150,
        readOnly: true,
    };

    state = {
        value: this.props.value,
        label: '',
        visiable: false,
    };

    componentDidMount() {
        this.updateSelected(this.props.children)
    }

    componentDidUpdate(prevProps: SelectProps, prevState: SelectState) {
        if (this.props.children !== prevProps.children || this.state.value !== prevState.value) {
            this.updateSelected(this.props.children);
        }
    }

    private updateSelected(options) {
        const optionsList = Array.isArray(options) ? options : [options];
        const selected = reduce(optionsList, (prev, option) => {
            const match = (option.props.value === this.props.value || option.props.selected) ? option : null;
            // 如果在Options中找到匹配项
            if (match) {
                // 如果之前已经有匹配，但是是由select属性匹配到的
                if (prev) {
                    // 此次匹配是value匹配，使用此次匹配
                    if (prev.props.selected) {
                        return match;
                    }
                    // 否则认为上一次的匹配是value匹配，抛弃此次匹配
                    else {
                        return prev;
                    }
                }
                // 如果之前没有匹配，使用此次匹配
                else {
                    return match;
                }
            } else {
                return prev;
            }
        }, null);

        if (selected) {
            this.setState({
                value: selected.props.value,
                label: selected.props.children,
            });
        }
    }

    handlePopupVisibleChange = (event: SweetUIEvent<boolean>) => {
        this.setState({
            visiable: event.detail,
        });
        this.dispatchPopupVisibleChangeEvent(event.detail);
    };

    dispatchPopupVisibleChangeEvent = createEventDispatcher(this.props.onPopupVisibleChange);

    handleSelect = (e, item, close) => {
        if (!item.props.disabled) {
            if (typeof item.props.onClick === 'function') {
                item.props.onClick(e);
            }
            if (!e.defaultPrevented) {
                this.setState({
                    value: item.props.value,
                    label: item.props.children,
                });

                this.dispatchSelectEvent(item.props.value)
                close()
            }
        }
    };

    dispatchSelectEvent = createEventDispatcher(this.props.onChange);

    render() {
        const { children, disabled, width, maxMenuHeight, className, readOnly, placeholder, role, element } = this.props;
        const { value, visiable, label } = this.state;
        const rootElement = element ||  this.context && this.context.element
        return (
            <Trigger
                role={role}
                element={rootElement}
                renderer={({ setPopupVisibleOnClick, role }) => (
                    <View
                        key={'selectTrigger'}
                        inline={true}
                        className={classnames(
                            styles['select'],
                            { [styles['disabled']]: disabled },
                            { [styles['focus']]: visiable },
                            className,
                        )}
                        style={{ width: `${width}px` }}
                        onClick={disabled ? undefined : setPopupVisibleOnClick}
                        role={role}
                    >
                        {/* {placeholder ? (
							<View
								style={{
									display: value ? 'none' : 'block'
								}}
								className={styles['placeholder']}
							>
								{placeholder}
							</View>
						) : null} */}
                        <View inline={true} className={classnames(styles['label'], { [styles['disabled']]: disabled })}>
                            <Text>{label}</Text>
                        </View>
                        <SweetIcon
                            name={'arrowDown'}
                            className={classnames(styles['arrow-down'], { [styles['disabled-icon']]: disabled })}
                        />
                    </View>
                )}
                onPopupVisibleChange={this.handlePopupVisibleChange}
                freeze={this.props.freeze}
                popupZIndex={this.props.popupZIndex}
            >
                {({ close }) => (
                    <Menu width={width} maxHeight={maxMenuHeight} className={styles['select-menu']}>
                        {React.Children.map(children, (child, index) =>
                            React.cloneElement(child, {
                                selected: child.props.value === value,
                                onClick: (e) => this.handleSelect(e, child, close),
                            }),
                        )}
                    </Menu>
                )}
            </Trigger>
        );
    }
}

Select.Option = SelectOption;

export default Select;
