import React from 'react';
import ReactDOM from 'react-dom';
import classnames from 'classnames';
import PopOver from '../PopOver';
import { AlignPosition } from '../PopOver/Locator';
import View from '../View';
import AppConfigContext from '@/core/context/AppConfigContext';
import styles from './styles';

/**
 * 气泡出现的位置
 */
export type Placement =
    | 'top'
    | 'left'
    | 'right'
    | 'bottom'
    | 'topLeft'
    | 'topRight'
    | 'bottomLeft'
    | 'bottomRight'
    | 'leftTop'
    | 'leftBottom'
    | 'rightTop'
    | 'rightBottom';

/**
 * 箭头宽度
 */
const ArrowWidth = Math.sqrt(6 * 6 * 2);

interface ValidateTipProps {
    /**
     * 用于手动控制浮层显隐
     */
    visible?: boolean;

    /**
     * 提示内容
     */
    content: string | React.ReactNode;

    /**
     * 气泡框位置
     */
    placement: Placement;

    /**
     * 气泡距操作项的距离，不包括箭头，单位 px
     */
    popoverDistance?: number;

    /**
     * 气泡样式
     */
    className?: string;

    /**
     * 提示状态
     */
    tipStatus?: 'normal' | 'error';

    /**
     * 获取浮层渲染父节点，默认渲染到body上
     */
    getContainer: () => HTMLElement;

    /**
     * 事件监听的目标元素，默认window
     */
    target?: HTMLElement;

    role?: string;
}

interface ValidateTipState {
    visible: boolean;
    arrowPlacement: { top: string | number; left: string | number };
    placement: string;
}

export default class ValidateTip extends React.Component<ValidateTipProps, ValidateTipState> {
    static contextType = AppConfigContext;
    static defaultProps = {
        placement: 'right',
        popoverDistance: 12,
        tipStatus: 'normal',
        getContainer: () => null,
    };

    constructor(props: ValidateTipProps, ...args: any) {
        super(props);

        this.state = {
            visible: !!props.visible,
            arrowPlacement: { top: -1, left: -1 },
            placement: props.placement,
        };
    }

    popover: HTMLElement | null = null;

    triggerElement: Element | null = null;

    arrowCoordinate: { top: number; left: number } = { top: -1, left: -1 };

    anchor: HTMLElement = document.body

    componentDidMount() {
        this.anchor = this.getRootDomNode()
    }

    static getDerivedStateFromProps(nextProps: ValidateTipProps) {
        if ('visible' in nextProps) {
            return { visible: nextProps.visible };
        }

        return null;
    }

    private getRootDomNode = () => {
        return ReactDOM.findDOMNode(this) || document.body;
    };

    private builtAlignConfig: (
        placement: string,
    ) => {
        origin: [AlignPosition, AlignPosition];
        offset: [AlignPosition, AlignPosition];
    } = (placement: string) => {
            switch (placement) {
                case 'top':
                    return {
                        origin: ['center', 'top'],
                        offset: ['center', 'bottom'],
                    };
                case 'topLeft':
                    return {
                        origin: ['left', 'top'],
                        offset: ['left', 'bottom'],
                    };
                case 'topRight':
                    return {
                        origin: ['right', 'top'],
                        offset: ['right', 'bottom'],
                    };
                case 'left':
                    return {
                        origin: ['left', 'center'],
                        offset: ['right', 'center'],
                    };
                case 'leftTop':
                    return {
                        origin: ['left', 'top'],
                        offset: ['right', 'top'],
                    };
                case 'leftBottom':
                    return {
                        origin: ['left', 'bottom'],
                        offset: ['right', 'bottom'],
                    };
                case 'right':
                    return {
                        origin: ['right', 'center'],
                        offset: ['left', 'center'],
                    };
                case 'rightTop':
                    return {
                        origin: ['right', 'top'],
                        offset: ['left', 'top'],
                    };
                case 'rightBottom':
                    return {
                        origin: ['right', 'bottom'],
                        offset: ['left', 'bottom'],
                    };
                case 'bottom':
                    return {
                        origin: ['center', 'bottom'],
                        offset: ['center', 'top'],
                    };
                case 'bottomLeft':
                    return {
                        origin: ['left', 'bottom'],
                        offset: ['left', 'top'],
                    };
                case 'bottomRight':
                    return {
                        origin: ['right', 'bottom'],
                        offset: ['right', 'top'],
                    };
                default:
                    return {
                        origin: ['right', 'bottom'],
                        offset: ['right', 'top'],
                    };
            }
        };

    private savePopOver = (node: HTMLElement) => {
        this.popover = node;
    };

    private handlePopupAlign = ({ x, y }: { x: number; y: number }) => {
        // 定位发生变化
        if ((this.popover && x !== this.state.arrowPlacement.left) || y !== this.state.arrowPlacement.top) {
            // todo 根据定位更新位置判断下一次箭头的placement，state更新箭头placement，从而更新箭头样式
            // const {clientWidth, clientHeight} = this.popover

            this.setState({
                arrowPlacement: { top: y, left: x },
            });
        }
    };

    render() {
        const { placement, popoverDistance, content, tipStatus } = this.props;
        const { arrowPlacement: { top, left } } = this.state;

        let visible = this.state.visible;

        if (!('visible' in this.props) && !content) {
            visible = false;
        }

        let popoverStyle = {},
            arrowStyle = {};

        switch (this.props.placement) {
            case 'top':
            case 'topLeft':
            case 'topRight':
                popoverStyle = {
                    paddingBottom: `${popoverDistance}px`,
                };
                arrowStyle = {
                    bottom: `${popoverDistance - ArrowWidth + 5.2}px`,
                    borderRight: tipStatus === 'error' ? '1px solid #E60012' : 'none',
                    borderBottom: tipStatus === 'error' ? '1px solid #E60012' : 'none',
                };
                break;

            case 'right':
            case 'rightTop':
            case 'rightBottom':
                popoverStyle = {
                    paddingLeft: `${popoverDistance}px`,
                };
                arrowStyle = {
                    left: `${popoverDistance - ArrowWidth + 5.2}px`,
                    borderLeft: tipStatus === 'error' ? '1px solid #E60012' : 'none',
                    borderBottom: tipStatus === 'error' ? '1px solid #E60012' : 'none',
                };
                break;

            case 'bottom':
            case 'bottomLeft':
            case 'bottomRight':
                popoverStyle = {
                    paddingTop: `${popoverDistance}px`,
                };
                arrowStyle = {
                    top: `${popoverDistance - ArrowWidth + 5.2}px`,
                    borderLeft: tipStatus === 'error' ? '1px solid #E60012' : 'none',
                    borderTop: tipStatus === 'error' ? '1px solid #E60012' : 'none',
                };
                break;

            case 'left':
            case 'leftTop':
            case 'leftBottom':
                popoverStyle = {
                    paddingRight: `${popoverDistance}px`,
                };
                arrowStyle = {
                    right: `${popoverDistance - ArrowWidth + 5.3}px`,
                    borderRight: tipStatus === 'error' ? '1px solid #E60012' : 'none',
                    borderTop: tipStatus === 'error' ? '1px solid #E60012' : 'none',
                };
                break;
        }

        const popup = (
            <View
                className={classnames(styles['popover'], this.props.className)}
                style={{
                    ...popoverStyle,
                    transformOrigin: `${left}px ${top}px `,
                }}
            >
                <View
                    className={classnames(styles['content'], {
                        [styles['error-tip']]: this.props.tipStatus === 'error',
                    })}
                    onMounted={this.savePopOver}
                >
                    {this.props.content}
                </View>
                <View
                    className={classnames(styles['arrow'], styles[this.props.placement])}
                    style={{
                        ...arrowStyle,
                        width: ArrowWidth,
                        height: ArrowWidth,
                    }}
                />
            </View>
        );

        return [
            React.cloneElement(this.props.children, {
                role: this.props.role,
                key: 'children',
            }),
            <PopOver
                key={'popover'}
                role={this.props.role}
                element={this.props.getContainer() || this.context?.element}
                target={this.props.target}
                open={visible}
                anchor={this.anchor}
                popup={popup}
                freeze={false}
                autoFix={false}
                anchorOrigin={this.builtAlignConfig(placement).origin}
                alignOrigin={this.builtAlignConfig(placement).offset}
                onPopupAlign={({ detail }) => this.handlePopupAlign(detail)}
            />,
        ];
    }
}
