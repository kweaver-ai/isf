import React from 'react';
import classnames from 'classnames';
import { isFunction, isEqual, trim } from 'lodash';
import SweetIcon from '../SweetIcon';
import Tag from '../Tag';
import styles from './styles';

interface ComboAreaProps {
    /**
     * 标签，支持数组或者对象数组，配合formatter函数显示标签的内容
     */
    value: ReadonlyArray<any>;

    /**
     * 禁用状态
     */
    disabled?: boolean;

    /**
     * 文本框宽度
     */
    width?: number;

    /**
     * 文本框高度
     */
    height?: number;

    /**
     * placeholder
     */
    placeholder: string;

    /**
     * css样式
     */
    style?: React.CSSProperties;

    /**
     * tag标签样式
     */
    tagClassName?: string;

    /**
     * 文本框状态
     */
    status?: 'normal' | 'error';

    /**
     * 输入框是否可编辑，默认不可编辑
     */
    editable?: boolean;

    /**
     * 格式化函数
     */
    formatter?: (o: any) => string;

    /**
     * 标签发生改变时触发，添加或者删除标签
     */
    onChange?: (value: ReadonlyArray<any>) => void;

    /**
     * 键盘输入时
     */
    onKeyDown?: (event: React.KeyboardEvent<Element>) => void;

    /**
     * 输入框聚焦时触发
     */
    onFocus?: (event: React.FocusEvent<Element>) => void;

    /**
     * 输入框失焦时触发
     */
    onBlur?: (event: React.FocusEvent<Element>) => void;

    /**
     * 鼠标进入文本域时触发
     */
    onMouseEnter?: (event: React.MouseEvent<HTMLElement>) => void;

    /**
     * 鼠标移除文本域时候触发
     */
    onMouseLeave?: (event: React.MouseEvent<HTMLElement>) => void;
}

interface ComboAreaState {
    /**
     * 标签
     */
    value: ReadonlyArray<any>;

    /**
     * 控制文本框聚焦样式
     */
    active: boolean;
}

export default class ComboArea extends React.Component<ComboAreaProps, ComboAreaState> {
    constructor(props) {
        super(props)
        this.state = {
            value: this.props.value,
            active: false,
        }
    }

    static defaultProps = {
        width: 200,
        height: 96,
        editable: false,
        status: 'normal',
        formatter: (val: any) => val,
    }

    /**
     * 外边框
     */
    wrapper: null | HTMLElement = null;

    /**
     * 输入框
     */
    input: null | HTMLElement = null;

    componentDidUpdate(prevState) {
        if (!isEqual(prevState.value, this.props.value)) {
            this.setState(() => ({
                value: this.props.value,
            }))
        }
    }

    /**
     * 输入框聚焦
     * @param {*} event React.FocusEvent<Element>
     * @memberof ComboArea
     */
    handleFocus(event: React.FocusEvent<Element>) {
        this.setState(() => ({
            active: true,
        }))
        this.input && this.input.focus();
        isFunction(this.props.onFocus) && this.props.onFocus(event);
    }

    /**
     * 输入框失焦
     * @param {React.FocusEvent<Element>} event 失焦
     * @memberof ComboArea
     */
    handleBlur(event: React.FocusEvent<Element>) {
        this.setState(() => ({
            active: false,
        }))
        this.input && this.input.blur();
        isFunction(this.props.onBlur) && this.props.onBlur(event);
    }

    /**
     * 删除标签
     * @param {*} deletedTagIndex 被删除的标签在数组value中的位置
     * @param {*} deletedTag 被删除的标签
     * @memberof ComboArea
     */
    handleDeleteTag(deletedTagIndex: number) {
        const newTags = this.state.value.filter((tag, index) => {
            return deletedTagIndex !== index
        })
        this.setState(() => ({
            value: newTags,
        }))
        isFunction(this.props.onChange) && this.props.onChange(newTags);

        // 解决IE11 删除标签输入框聚焦后点击外部不发生失焦事件
        this.wrapper && this.wrapper.focus()
    }

    /**
     * 输入框有内容输入时
     * @param {*} event React.FocusEvent<Element>
     * @memberof ComboArea
     */
    handleKeyDown(event: React.KeyboardEvent<Element>) {
        if (this.input) {
            if (event.keyCode === 13) {
                event.preventDefault()

                const inputValue = this.getInputValue(this.input)
                if (inputValue !== '') {
                    const newTags = [...this.state.value, inputValue]
                    this.setState(() => ({
                        value: newTags,
                    }))

                    isFunction(this.props.onChange) && this.props.onChange(newTags);
                    this.setInputValue('', this.input)
                }
            }

            // 当按下删除键且输入的内容为空的时候删除的是已经生成的图标
            if (event.keyCode === 8 && this.getInputValue(this.input) === '') {
                const newTags = this.state.value.filter((tag, index) => {
                    return index !== this.state.value.length - 1
                })
                this.setState(() => ({
                    value: newTags,
                }))
                isFunction(this.props.onChange) && this.props.onChange(newTags);
            }

            isFunction(this.props.onKeyDown) && this.props.onKeyDown(event);
        }
    }

    /**
     * 获取输入框的值
     */
    private getInputValue(domNode: HTMLElement): string {
        return 'innerText' in domNode ? trim(domNode.innerText) : trim(domNode.textContent)
    }

    /**
     * 设置输入框的值
     */
    private setInputValue(value: string, domNode: HTMLElement): void {
        'innerText' in domNode ? domNode.innerText = trim(value) : domNode.textContent = trim(value);
    }

    render() {
        const { width, height, style, placeholder, status, disabled, editable, formatter, onMouseEnter, onMouseLeave, tagClassName } = this.props;
        const { value, active } = this.state
        const maxWidth = `${width - 8 * 2 - 20}px` // 标签及输入框的最大宽度，20为预留的出现滚动条的宽度

        return (
            <div
                className={classnames(
                    styles['wrapper'],
                    {
                        [styles[`${status}`]]: status,
                        [styles[`${status}-active`]]: active && !disabled,
                        [styles['disabled']]: disabled,
                    },
                )}
                style={{ ...style, width, height }}
                tabIndex={1}
                ref={(node) => this.wrapper = node}
                onFocus={this.handleFocus.bind(this)}
                onBlur={this.handleBlur.bind(this)}
                onMouseEnter={onMouseEnter}
                onMouseLeave={onMouseLeave}
            >
                <div
                    className={styles['placeholder']}
                    placeholder={placeholder}
                >
                    {
                        value && value.length > 0 && value.map((tag, index) => {
                            return <Tag
                                key={index}
                                disabled={disabled}
                                className={classnames(styles['tag'], tagClassName ? tagClassName : null)}
                                style={{ maxWidth }}
                                closable={true}
                                onClose={this.handleDeleteTag.bind(this, index)}
                            >
                                {
                                    formatter(tag)
                                }
                            </Tag>
                        })
                    }
                    {/* {
                    !disabled && editable ?
                            <div
                                className={styles['input-wrapper']}
                                style={{ maxWidth }}
                                ref={node => this.input = node}
                                contentEditable={true}
                                placeholder={(!value || value.length === 0) ? placeholder : ''} // 可编辑状态，输入框中的placeholder生效，因为对于外层的div,子元素始终不为empty
                                onKeyDown={this.handleKeyDown.bind(this)}
                                onInput={this.handleInputValueChange.bind(this)}
                            />
                            : null
                   } */}
                </div>
                {
                    status === 'error' ?
                        <SweetIcon
                            name={'caution'}
                            size={16}
                            color={'#e60012'}
                            className={styles['caution']}
                        />
                        : null
                }
            </div >
        );
    }
}