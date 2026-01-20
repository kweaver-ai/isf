import React from 'react';
import classnames from 'classnames';
import { createEventDispatcher, SweetUIEvent } from '../../utils/event';
import Portal from '../Portal';
import View from '../View';
import SweetIcon from '../SweetIcon';
import styles from './styles';
import AppConfigContext from '@/core/context/AppConfigContext';

interface DrawerProps {
    /**
     * 角色
     */
    role?: string;

    /**
     * Drawer的位置
     */
    position?: 'top' | 'bottom' | 'left' | 'right';

    /**
     * Drawer展开状态，当需要Drawer为受控组件时使用
     */
    open?: boolean;

    /**
     * Drawer初始展开状态，当需要Drawer为非受控组件时使用
     */
    defaultOpen?: boolean;

    /**
     * 关闭时销毁Drawer里的子元素
     */
    destroyOnClose?: boolean;

    /**
     * 是否显示遮罩层
     */
    showMask?: boolean;

    /**
     * 点击遮罩层是否允许关闭
     */
    maskClosable?: boolean;

    /**
     * 标题
     */
    title: string | React.ReactNode;

    /**
     * 自定义Footer
     */
    footer?: React.ReactNode;

    /**
     * Drawer大小，当position为'left' | 'right'时设置的是width，否则设置height
     * 默认设置width
     */
    size?: string | number;

    /**
     * Drawer-container的z-index值
     */
    zIndex?: number;

    /**
     * 抽屉内容className
     */
    drawerBodyClassName?: string;

    /**
     * 点击Drawer外部是否关闭抽屉
     * todo 判断是否抽屉外部的方案待改进，暂时调整为点击mask判断是否关闭抽屉
     * 引入的新问题：没有mask的场景下如何判断关闭条件，待解决
     */
    // canOutsideClickClose?: boolean;

    /**
     * 获取Drawer挂载的DOM节点
     * todo 非body时处理mask位置
     */
    getContainer?: () => React.ReactNode;

    /**
     * 关闭Drawer时的回调，如果传入`open`属性，该方法为必传项
     */
    onDrawerClose?: (event: SweetUIEvent<boolean>) => void;
}

interface DrawerState {
    hasEverOpened: boolean;
    open: boolean;
}

export default class Drawer extends React.Component<DrawerProps, DrawerState> {
    static contextType = AppConfigContext;
    static defaultProps = {
        position: 'right',
        getContainer: null,
        defaultOpen: false,
        showMask: true,
        size: '50%',
        zIndex: 20,
        // canOutsideClickClose: true,
        destroyOnClose: false,
        maskClosable: true,
    };

    constructor(props: DrawerProps, ...args: any[]) {
        super(props);
        this.state = {
            hasEverOpened: typeof props.open !== 'undefined' ? props.open : !!props.defaultOpen,
            open: typeof props.open !== 'undefined' ? props.open : !!props.defaultOpen,
        };
    }

    drawerContainer: React.ReactNode | null = null;

    drawer: HTMLDivElement | null = null;

    mask: HTMLElement | null = null;

    closeTimer = null

    static getDerivedStateFromProps(nextProps: DrawerProps, prevState: DrawerState) {
        if ('open' in nextProps && nextProps.open !== prevState.open) {
            return {
                open: nextProps.open,
            }
        }

        return null;
    }

    componentDidUpdate(prevProps: DrawerProps, prevState: DrawerState) {
        if (this.state.open !== prevState.open) {
            this.toggleDrawerOpen(this.state.open);
        }
    }

    toggleDrawerOpen = (open: boolean) => {
        if (open) {
            if (!this.drawerContainer) {
                this.setState({
                    hasEverOpened: true,
                    open: true,
                });
            } else {
                if (!this.props.destroyOnClose) {
                    this.setState({
                        open: true,
                    });
                    this.drawerContainer.style.display = ''
                    this.drawer && (this.drawer as HTMLDivElement).setAttribute('class', classnames(styles['drawer'], styles[`${this.props.position}`]))
                }
            }
            this.toggleBodyScroll(true);
        } else {
            this.delayToClose()
        }
    };

    toggleBodyScroll = (open: boolean) => {
        if (open) {
            document.body.style.overflow = 'hidden';
            document.body.style.paddingRight = '17px';
        } else {
            document.body.style.overflow = '';
            document.body.style.paddingRight = '';
        }
    };

    delayToClose = () => {
        if (this.closeTimer) {
            clearTimeout(this.closeTimer)
        }
        this.mask &&
            (this.mask as HTMLDivElement).setAttribute('class', classnames(styles['mask'], styles['mask-out']));
        this.drawer &&
            (this.drawer as HTMLDivElement).setAttribute(
                'class',
                classnames(styles['drawer'], styles[`${this.props.position}-out`]),
            );
        this.closeTimer = setTimeout(() => {
            this.closeDrawer();
        }, 200);
    };

    /**
     * 关闭抽屉
     */
    closeDrawer = () => {
        if (this.props.destroyOnClose) {
            this.destory();
        } else {
            this.hide();
        }
        this.toggleBodyScroll(false);
    };

    hide = () => {
        this.setState({
            open: false,
        });
        if (this.drawerContainer) {
            this.drawerContainer.style.display = 'none'
        }

        this.dispatchDrawerClose(false);
    };

    destory = () => {
        this.setState({
            hasEverOpened: false,
            open: false,
        });
        this.drawerContainer = null;
        this.dispatchDrawerClose(false);
    };

    dispatchDrawerClose = createEventDispatcher(this.props.onDrawerClose);

    getContainer = () => {
        const drawerContainer = document.createElement('div');
        const rootElement = (typeof this.props.getContainer === 'function' && this.props.getContainer()) || this.context && this.context.element || document.body

        rootElement.appendChild(drawerContainer as HTMLDivElement);
        (drawerContainer as HTMLDivElement).setAttribute('class', styles['drawer-container']);
        (drawerContainer as HTMLDivElement).setAttribute('style', `z-index: ${this.props.zIndex}`);
        (drawerContainer as HTMLDivElement).setAttribute('role', this.props.role)
        this.drawerContainer = drawerContainer;

        return drawerContainer;
    };

    saveDrawer = (node: HTMLDivElement) => {
        this.drawer = node;
    };

    saveMask = (node: HTMLElement) => {
        this.mask = node;
    };

    handleClickMask = (e: MouseEvent) => {
        if (this.state.open && this.props.maskClosable) {
            this.delayToClose();
        }
    }

    render() {
        const { showMask, size, position, drawerBodyClassName } = this.props;
        const { hasEverOpened, open } = this.state;

        return hasEverOpened ? (
            <Portal getContainer={this.getContainer}>
                {showMask ? (
                    <View
                        onMounted={this.saveMask}
                        className={classnames(styles['mask'], { [styles['mask-in']]: open })}
                        onClick={this.handleClickMask}
                    />
                ) : null}
                <div
                    className={classnames(styles['drawer'], { [styles[position!]]: open })}
                    style={{
                        [position === 'left' || position === 'right' ? 'width' : 'height']: size,
                    }}
                    ref={this.saveDrawer}
                >
                    {this.renderHeader()}
                    <View className={classnames(styles['drawer-body'], drawerBodyClassName)}>{this.props.children}</View>
                    {this.props.footer ? this.renderFooter() : null}
                </div>
            </Portal>
        ) : null;
    }

    private renderHeader = () => {
        return (
            <View className={styles['drawer-header']}>
                <View inline={true} className={styles['drawer-title']}>
                    {this.props.title}
                </View>
                <View onClick={this.delayToClose} className={styles['close-icon']}>
                    <SweetIcon name={'x'} size={16} />
                </View>
            </View>
        );
    };

    private renderFooter = () => {
        return <View className={styles['drawer-footer']}>{this.props.footer}</View>;
    };
}
