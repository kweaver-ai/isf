import React from 'react';
import ReactDOM from 'react-dom';
import classnames from 'classnames';
import { isFunction, isEqual } from 'lodash';
import Text from '@/ui/Text/ui.desktop';
import { ClassName } from '@/ui/helper';
import { createEventDispatcher, SweetUIEvent } from '../../utils/event';
import EventListener from '../EventListener/index';
import Trigger from '../Trigger';
import View from '../View';
import SweetIcon from '../SweetIcon';
import { arrayTreeFilter } from './util';
import styles from './styles';
import __ from './locale';

const MinLevelNumber = 3;

const HiddenLevelWidth = 30;

interface CascaderProps {
    /**
     * 下拉框宽度
     */
    width: number;

    /**
     * 次级列表宽度
     */
    subWidth?: number;

    /**
     * 指定渲染item
     */
    renderItem: (option: any) => any;

    /**
     * 可选项数据源
     */
    options: ReadonlyArray<any>;

    /**
     * 已选项
     */
    value: Array<any>;

    /**
     * 下拉框占位文本
     */
    placeholder: string;

    /**
     * 级联菜单禁用状态
     */
    disabled?: boolean;

    /**
     * 级联菜单子项字段名
     */
    childrenKeyName?: string;

    /**
     * 选择即改变
     */
    isChangeOnSelect?: boolean;

    /**
     * 选择后展示的渲染函数
     */
    displayRender: (value: Array<any>) => string;

    /**
     * 自定义renderer
     */
    renderer?: ({ setPopupVisibleOnClick }: { setPopupVisibleOnClick: () => void }) => any;

    /**
     * 级联选择时的回调
     */
    onChange: (options: Array<any>) => void;

    /**
     * 菜单关闭之前触发
     */
    onBeforePopupClose: () => void;

    /**
     * 判断是否选中项
     */
    predictActiveOption: (args: { value: Array<any>; option: any; index: number }) => boolean;

    /**
     * 获取菜单项key的方法
     */
    getItemKey?: (item: any) => string;
}

interface CascaderState {
    /**
     * 选中项value集合
     */
    activeValue: Array<any>;

    /**
     * 菜单展开状态
     */
    visiable: boolean;

    selectedValue: Array<any>;
}

export default class Cascader extends React.Component<CascaderProps, CascaderState> {
    static defaultProps = {
        disabled: false,
        options: [],
        value: [],
        width: 200,
        isChangeOnSelect: false,
        displayRender: (value: Array<any>) => '',
        predictActiveOption: () => false,
        childrenKeyName: 'children',
        getItemKey: (option: any) => String(option.id),
    };

    preVisiable = false;

    trigger = null;

    menuRef = null;

    delayTimer = null;

    menuItems = {};

    constructor(props: CascaderProps, ...args: any[]) {
        super(props);
        let initialValue = [];
        if (Array.isArray(props.value) && props.value.length) {
            initialValue = arrayTreeFilter(
                this.props.options,
                (option, level) => this.props.predictActiveOption({ value: props.value, option, index: level }),
                {
                    childrenKeyName: props.childrenKeyName,
                },
            );
        }

        this.state = {
            activeValue: initialValue,
            visiable: false,
            selectedValue: initialValue,
        };

        this.preVisiable = this.state.visiable;
    }

    componentDidUpdate(prevProps: CascaderProps, prevState: CascaderState) {
        if (this.props.value !== prevProps.value) {

            const activeValue = arrayTreeFilter(
                this.props.options,
                (option, level) => this.props.predictActiveOption({ value: this.props.value, option, index: level }),
                {
                    childrenKeyName: this.props.childrenKeyName || 'children',
                },
            );

            this.setState({
                activeValue,
                selectedValue: activeValue,
            })
        }
    }

    componentWillUnmount() {
        this.clearDelayTimer();
    }

    clearDelayTimer = () => {
        if (this.delayTimer) {
            clearTimeout(this.delayTimer);
            this.delayTimer = null;
        }
    };

    handlePopupVisibleChange = (event: SweetUIEvent<boolean>) => {
        const visiable = event.detail;

        this.setState({
            visiable,
            activeValue: this.state.selectedValue,
        });

        if (!visiable) {
            this.preVisiable = false;
        }
    };

    saveMenuItem = (index) => (node) => {
        this.menuItems[index] = node;
        if (!this.preVisiable && this.state.visiable) {
            this.scrollActiveItemToView(this.menuItems);
        }
    };

    isActiveOption = (option: any, menuIndex: number) => {
        return isEqual(this.state.activeValue[menuIndex], option);
    };

    scrollActiveItemToView = (menuItems: Array<Element>) => {
        const optionsLength = this.getShowOptions().length;

        for (let i = 0; i < optionsLength; i++) {
            const itemComponent = menuItems[i];
            if (itemComponent) {
                const target = ReactDOM.findDOMNode(itemComponent);
                if (target && target.parentNode) {
                    target.parentNode.scrollTop = target.offsetTop;
                }
            }
        }
    };

    handleMenuSelect = (targetOption: any, menuIndex: number, close?: () => void) => {
        // Keep focused state for keyboard support
        const triggerNode = this.trigger.getRootDomNode();

        if (triggerNode && triggerNode.focus) {
            triggerNode.focus();
        }

        if (!targetOption || targetOption.disabled) {
            return;
        }

        let { activeValue } = this.state;

        activeValue = activeValue.slice(0, menuIndex + 1);
        activeValue[menuIndex] = targetOption;

        if (!targetOption[this.props.childrenKeyName] || !targetOption[this.props.childrenKeyName].length) {
            this.setState({
                selectedValue: activeValue,
            })
            this.handleChange(activeValue, close);
        } else if (this.props.isChangeOnSelect) {
            this.handleChange(activeValue)
        }

        this.setState({ activeValue });
    };

    handleChange = (options: Array<any>, close?: () => void) => {
        this.dispatchChangeEvent(options);
        if (isFunction(close)) {
            isFunction(this.props.onBeforePopupClose) && this.props.onBeforePopupClose()
            close()
        }
    };

    dispatchChangeEvent = createEventDispatcher(this.props.onChange);

    /**
     * 获取每一层级菜单项
     */
    getShowOptions = () => {
        let result = this.state.activeValue
            .map((activeOption) => activeOption[this.props.childrenKeyName])
            .filter((activeOption) => !!activeOption && activeOption.length > 0);
        result = [this.props.options, ...result];
        return result;
    };

    saveTrigger = (node) => {
        this.trigger = node;
    };

    saveMenuRef = (node) => {
        this.menuRef = node;

        this.preVisiable = this.state.visiable;
    };

    handleResize = () => {
        this.forceUpdate();
    };

    getHiddenMenuNumberWhenOverflow = (menuNumber: number) => {
        if (this.trigger) {
            const triggerNode = this.trigger.getRootDomNode();
            const { top, left, bottom, right } = triggerNode.getBoundingClientRect();
            const windowWidth = window.innerWidth;

            if (left + this.props.width * menuNumber > windowWidth) {
                // 不足以向右显示所有层级
                if (left + this.props.width * menuNumber + HiddenLevelWidth < right) {
                    // 足以向左显示所有层级则不做省略
                    return 0;
                } else {
                    // 不足以向左显示所有层级
                    if (left + this.props.width * MinLevelNumber + HiddenLevelWidth < windowWidth) {
                        // 可以向右展开最小层级数
                        // 省略中间的层级
                        const maxLevel = Math.floor((windowWidth - left - HiddenLevelWidth) / this.props.width);

                        return maxLevel >= MinLevelNumber ? menuNumber - maxLevel : 0;
                    } else if (left + this.props.width * MinLevelNumber + HiddenLevelWidth < right) {
                        // 可以向左展开最小层级数
                        const maxLevel = Math.floor((right - HiddenLevelWidth) / this.props.width);

                        return maxLevel >= MinLevelNumber ? menuNumber - maxLevel : 0;
                    } else {
                        // 都不能展开最小层级数
                        return 0;
                    }
                }
            } else {
                // 足以向右显示所有层级
                return 0;
            }
        }
    };

    getOption = (option: any, menuIndex: number, close: () => void) => {
        const onSelect = this.handleMenuSelect.bind(this, option, menuIndex, close);
        const hasChildren = option[this.props.childrenKeyName] && option[this.props.childrenKeyName].length > 0;
        let expandProps = {
            onClick: onSelect,
        };

        if (this.isActiveOption(option, menuIndex)) {
            expandProps.ref = this.saveMenuItem(menuIndex);
        }

        return (
            <li
                className={classnames(
                    styles['menu-item'],
                    { [styles['menu-item-active']]: this.isActiveOption(option, menuIndex) },
                    { [styles['menu-item-disabled']]: option.disabled },
                    { [ClassName.Color]: this.isActiveOption(option, menuIndex) },
                    ClassName.Color__Hover,
                )}
                key={this.props.getItemKey!(option)}
                title={option.title || this.props.renderItem(option)}
                {...expandProps}
                onMouseDown={(e) => e.preventDefault()}
            >
                {this.props.renderItem(option)}
                {hasChildren ? <SweetIcon className={styles['menu-item-expand-icon']} name={'arrowRight'} /> : null}
            </li>
        );
    };

    render() {
        const { width, placeholder, onBeforePopupClose, disabled, options, displayRender, renderer } = this.props;
        const { visiable, activeValue, selectedValue } = this.state;

        const result = this.getShowOptions();

        const ellipsisNumberOfLevels = this.getHiddenMenuNumberWhenOverflow(result.length) || 0;

        return (
            <EventListener target="window" onResize={this.handleResize}>
                <Trigger
                    ref={this.saveTrigger}
                    anchorOrigin={['left', 'bottom']}
                    alignOrigin={['left', 'top']}
                    renderer={({ setPopupVisibleOnClick }) => {
                        if (renderer && isFunction(renderer)) {
                            return renderer({ setPopupVisibleOnClick, activeValue })
                        }
                        return (
                            <View
                                key={'cascaderTrigger'}
                                inline={true}
                                className={classnames(
                                    styles['select'],
                                    { [styles['disabled']]: disabled },
                                    { [styles['focus']]: visiable },
                                )}
                                style={{ width: `${width}px` }}
                                onClick={disabled ? undefined : setPopupVisibleOnClick}
                                title={displayRender(selectedValue)}
                            >
                                <View
                                    inline={true}
                                    className={classnames(styles['label'], { [styles['disabled']]: disabled })}
                                >
                                    <Text>{displayRender(selectedValue)}</Text>
                                </View>
                                <SweetIcon
                                    name={'arrowDown'}
                                    className={classnames(styles['arrow-down'], { [styles['disabled-icon']]: disabled })}
                                />
                            </View>
                        )
                    }}
                    onBeforePopupClose={onBeforePopupClose}
                    onPopupVisibleChange={this.handlePopupVisibleChange}
                >
                    {({ close }) =>
                        options && options.length > 0 ? (
                            <View className={styles['cascader']} onMounted={this.saveMenuRef}>
                                {(
                                    ellipsisNumberOfLevels > 0 ?
                                        [
                                            result[0],
                                            null,
                                            ...result.slice(Math.min(result.length - 1, ellipsisNumberOfLevels + 1)),
                                        ]
                                        : result
                                ).map(
                                    (options, menuIndex) =>
                                        Array.isArray(options) ?
                                            (
                                                <ul
                                                    className={styles['menu']}
                                                    key={
                                                        menuIndex === 0 ?
                                                            0
                                                            : ellipsisNumberOfLevels === 0 ?
                                                                menuIndex
                                                                : (menuIndex + ellipsisNumberOfLevels - 1)
                                                    }
                                                    style={{
                                                        width:
                                                            this.props.subWidth && !options.some((option) => this.props.options.some(({ id }) => id === option.id)) ?
                                                                `${this.props.subWidth}px`
                                                                :
                                                                menuIndex === options.length - 1
                                                                    ? `${this.props.width}px`
                                                                    : `${this.props.width - 1}px`,
                                                    }}
                                                >
                                                    {options.map((option) =>
                                                        this.getOption(
                                                            option,
                                                            menuIndex === 0
                                                                ? 0
                                                                : ellipsisNumberOfLevels === 0
                                                                    ? menuIndex
                                                                    : menuIndex + ellipsisNumberOfLevels - 1,
                                                            close,
                                                        ),
                                                    )}
                                                </ul>
                                            )
                                            : (
                                                <View
                                                    key={'back'}
                                                    inline={true}
                                                    className={classnames(styles['menu'], styles['hidden-menu'])}
                                                    style={{ width: `${HiddenLevelWidth}px` }}
                                                    onClick={this.handleMenuSelect.bind(
                                                        this,
                                                        activeValue[ellipsisNumberOfLevels],
                                                        ellipsisNumberOfLevels,
                                                    )}
                                                    title={__('返回')}
                                                >
                                                    ...
                                                </View>
                                            ),
                                )}
                            </View>
                        ) : null}
                </Trigger>
            </EventListener>
        );
    }
}
