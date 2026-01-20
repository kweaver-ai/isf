import React from 'react'
import { createRoot, Root } from 'react-dom/client'
import classnames from 'classnames'
import { omit } from 'lodash'
import { bindEvent, unbindEvent } from '@/util/browser'
import styles from './styles.desktop'
import AppConfigContext from '@/core/context/AppConfigContext'

const isChildNode = (el: Node | null, target: Node | null): boolean => {
    if (target !== null && el !== null) {
        return el === target || isChildNode(el, target.parentNode)
    }
    return false
}

export default class PopOverBase extends React.Component<UI.PopOver.Props, any> {
    static contextType = AppConfigContext

    private root: Root | null = null
    private layer: HTMLDivElement | null = null
    private popContainer: HTMLDivElement | null = null
    private watchTimeoutId?: number
    private closeTimeoutId?: number
    private rendering = false

    static defaultProps = {
        anchorOrigin: [0, 0],
        targetOrigin: [0, 0],
        watch: false,
        autoFix: true,
        freezable: true,
        keepOpenWhenMouseOver: true,
        closeWhenMouseLeave: false,
        closeTimeout: 150,
        hideLayerWhenClose: false,
    }

    state = {
        open: false,
        anchor: null,
    }
  
    constructor(props, context) {
        super(props, context)
        this.handleClickAway = this.handleClickAway.bind(this)
        this.watch = this.watch.bind(this)
        this.close = this.close.bind(this)
    }

    componentDidMount() {
        const { open, anchor } = this.props

        this.setState({
            ...(open !== undefined ? { open } : {}),
            ...(anchor !== undefined ? { anchor } : {}),
        })
    }

    static getDerivedStateFromProps({ open, anchor }, prevState) {
        if (open !== undefined || anchor !== undefined) {
            return {
                open: open !== undefined ? open : prevState.open,
                anchor: anchor !== undefined ? anchor : prevState.anchor,
            }
        }
        return null
    }

    componentDidUpdate() {
        if (!this.rendering) {
            this.renderToLayer()
        }
    }

    componentWillUnmount() {
        this.unrenderToLayer()
    }

    /**
 * 关闭弹出内容
 */
    close() {
        if (typeof this.props.onClose === 'function') {
            this.props.onClose()
        } else {
            this.setState({
                open: false,
            })
        }
    }

    /**
 * 显示弹出内容
 */
    open() {
        if (typeof this.props.onOpen === 'function') {
            this.props.onOpen()
        } else {
            this.setState({
                open: true,
            })
        }
    }

    /**
 * 点击
 * @param e
 */
    async handleClick(e) {

        if (typeof this.props.onClick === 'function') {
            this.props.onClick(e)
        }

        if (typeof this.props.onRequestCloseWhenClick === 'function') {
            try {
                await new Promise((resolve) => {
                    (this.props.onRequestCloseWhenClick as (close: () => void) => void)(resolve)
                })
                this.close()
            } catch (e) { }
        }
    }

    /**
 * 点击popcontainer之外触发onClickAway
 */
    async handleClickAway(e) {
        const clickElement = e.target || e.srcElement
        if (this.state.open && this.popContainer && !isChildNode(this.popContainer, clickElement)) {

            if (typeof this.props.onClickAway === 'function') {
                this.props.onClickAway(e)
            }

            if (typeof this.props.onRequestCloseWhenBlur === 'function') {
                try {
                    await new Promise((resolve) => {
                        (this.props.onRequestCloseWhenBlur as (close: () => void) => void)(resolve)
                    })
                    this.close()
                } catch (e) { }
            }
        }
    }

    handleMouseEnter(e) {

        if (typeof this.props.onMouseEnter === 'function') {
            this.props.onMouseEnter(e)
        }
        if (this.props.triggerEvent === 'mouseover' && this.props.keepOpenWhenMouseOver) {
            clearTimeout(this.closeTimeoutId)
            this.open()
        }
    }

    handleMouseLeave(e) {

        if (typeof this.props.onMouseLeave === 'function') {
            this.props.onMouseLeave(e)
        }
        if (this.props.triggerEvent === 'mouseover' && this.props.closeWhenMouseLeave) {
            clearTimeout(this.closeTimeoutId)
            /**
         * 延时150ms，防止手抖导致弹出窗关闭
         */
            this.closeTimeoutId = setTimeout(this.close, this.props.closeTimeout)
        }
    }

    /** 点击trigger */
    handleTriggerClick(e, triggerProps) {
        this.setState({ anchor: e.currentTarget })
        this.open()
        if (typeof triggerProps.onClick === 'function') {
            triggerProps.onClick(e)
        }
    }

    /** 鼠标移入trigger */
    handleTriggerMouseEnter(e, triggerProps) {
        clearTimeout(this.closeTimeoutId)
        this.setState({ anchor: e.currentTarget })
        this.open()
        if (typeof triggerProps.onMouseEnter === 'function') {
            triggerProps.onMouseEnter(e)
        }
    }

    /**
 * 鼠标移出trigger
 */
    handleTriggerMouseLeave(e, triggerProps) {
        clearTimeout(this.closeTimeoutId)
        /**
    * 延时150ms，防止手抖导致弹出窗关闭
    */
        this.closeTimeoutId = setTimeout(this.close, this.props.closeTimeout)
        if (typeof triggerProps.onMouseLeave === 'function') {
            triggerProps.onMouseLeave(e)
        }
    }

    /**
   * 创建挂载节点，并渲染弹出内容
   */
    renderToLayer() {
        const { freezable, watch, className, children, popContainerClassName, ...otherProps } = this.props
        const rootElement = this.props.element || this.context?.element
    
        if (this.state.open) {
            if (!this.layer) {
                this.layer = document.createElement('div')
                const container = rootElement || document.body
                container.appendChild(this.layer)
        
                if (freezable) {
                    this.layer.className = classnames(styles.layer, { 
                        [styles['pop-absolute']]: !!rootElement 
                    }, className)
                }
            }

            if (!this.root) {
                this.root = createRoot(this.layer!)
            }

            const popContent = (
                <AppConfigContext.Provider value={this.context}>
                    <div
                        className={classnames(styles['pop-container'], { 
                            [styles['pop-absolute']]: !!rootElement 
                        }, popContainerClassName)}
                        {...omit(otherProps,
                            'anchorOrigin', 'targetOrigin', 'autoFix', 'triggerEvent',
                            'onRequestCloseWhenClick', 'onRequestCloseWhenBlur', 'keepOpenWhenMouseOver',
                            'closeWhenMouseLeave', 'closeTimeout', 'hideLayerWhenClose', 'onOpen', 'onClose',
                        )}
                        ref={ref => {
                            this.popContainer = ref
                            if (ref) this.setPosition() 
                        }}
                        onClick={this.handleClick.bind(this)}
                        onMouseEnter={this.handleMouseEnter.bind(this)}
                        onMouseLeave={this.handleMouseLeave.bind(this)}
                        data-test-scope="ui/PopOver"
                    >
                        {children}
                    </div>
                </AppConfigContext.Provider>
            )

            this.rendering = true

            this.root.render(popContent)

            setTimeout(() => {
                bindEvent(document, 'mousedown', this.handleClickAway)
                bindEvent(document, 'touchstart', this.handleClickAway)
                this.rendering = false
                this.props.hideLayerWhenClose && this.layer?.classList.add(styles.layer)
                if (watch) this.watch()
            })
        } else {
            this.unrenderToLayer()
        }
    }

    /**
   * 卸载弹出内容并删除挂载节点
   */
    unrenderToLayer() {
        clearTimeout(this.watchTimeoutId);
        const rootElement = this.props.element || this.context?.element;
    
        if (this.layer) {
            if (this.props.hideLayerWhenClose) {
                this.layer.style.display = 'none';
            } else {
        
                setTimeout(() => {
                    if (this.root) {
                        this.root.unmount();
                        this.root = null;
                    }

                    const container = rootElement || document.body;
                    if (container.contains(this.layer)) {
                        container.removeChild(this.layer);
                    }
          
                    unbindEvent(document, 'mousedown', this.handleClickAway);
                    unbindEvent(document, 'touchstart', this.handleClickAway);
                    this.layer = null;
                });
            }
        }
    }

    watch() {
        clearTimeout(this.watchTimeoutId)
        this.setPosition()
        this.watchTimeoutId = setTimeout(this.watch, 16)
    }

    /**
     * 根据定位信息计算弹出层位置
     */
    setPosition() {
        const { anchorOrigin, targetOrigin, autoFix } = this.props
        const { anchor } = this.state
        let windowInnerWidth = 0
        let windowInnerHeight = 0
        let elementRect = { top: 0, left: 0, right: 0, bottom: 0 }
        const rootElement = this.props.element || this.context?.element
        if (rootElement) {
            windowInnerWidth = rootElement.clientWidth
            windowInnerHeight = rootElement.clientHeight
            elementRect = rootElement.getBoundingClientRect()
        } else {
            windowInnerWidth = window.innerWidth || window.document.documentElement.clientWidth
            windowInnerHeight = window.innerHeight || window.document.documentElement.clientHeight
        }
        if (this.popContainer) {
            let anchorOriginX, anchorOriginY,
                targetOriginX, targetOriginY,
                anchorRect,
                { offsetHeight, offsetWidth } = this.popContainer,
                left = 0,
                top = 0
            if (anchor) {
                anchorRect = anchor.getBoundingClientRect()
                switch (anchorOrigin[1]) {
                    case 'top':
                        anchorOriginY = 0;
                        break
                    case 'bottom':
                        anchorOriginY = anchorRect.bottom - anchorRect.top
                        break
                    case 'middle':
                        anchorOriginY = (anchorRect.bottom - anchorRect.top) / 2
                        break
                    default:
                        anchorOriginY = anchorOrigin[1]
                        break
                }
                switch (anchorOrigin[0]) {
                    case 'left':
                        anchorOriginX = 0;
                        break
                    case 'right':
                        anchorOriginX = anchorRect.right - anchorRect.left
                        break
                    case 'center':
                        anchorOriginX = (anchorRect.right - anchorRect.left) / 2
                        break
                    default:
                        anchorOriginX = anchorOrigin[0]
                        break
                }
            } else {
                anchorRect = { top: 0, left: 0, right: 0, bottom: 0 }
                anchorOriginX = anchorOriginY = 0
                switch (anchorOrigin[0]) {
                    case 'left':
                        break
                    case 'right':
                        anchorRect.left = anchorRect.right = windowInnerWidth
                        break
                    case 'center':
                        anchorRect.left = anchorRect.right = windowInnerWidth / 2
                        break
                    default:
                        anchorRect.left = anchorRect.right = anchorOrigin[0]
                        break
                }
                switch (anchorOrigin[1]) {
                    case 'top':
                        break
                    case 'bottom':
                        anchorRect.top = anchorRect.bottom = windowInnerHeight
                        break
                    case 'middle':
                        anchorRect.top = anchorRect.bottom = windowInnerHeight / 2
                        break
                    default:
                        anchorRect.top = anchorRect.bottom = anchorOrigin[1]
                        break
                }
            }
            switch (targetOrigin[1]) {
                case 'top':
                    targetOriginY = 0;
                    break
                case 'bottom':
                    targetOriginY = offsetHeight
                    break
                case 'middle':
                    targetOriginY = offsetHeight / 2
                    break
                default:
                    targetOriginY = targetOrigin[1]
                    break
            }
            switch (targetOrigin[0]) {
                case 'left':
                    targetOriginX = 0;
                    break
                case 'right':
                    targetOriginX = offsetWidth
                    break
                case 'center':
                    targetOriginX = offsetWidth / 2
                    break
                default:
                    targetOriginX = targetOrigin[0]
                    break
            }
            left = anchorRect.left + anchorOriginX - targetOriginX - elementRect.left
            top = anchorRect.top + anchorOriginY - targetOriginY - elementRect.top
            if (autoFix) {
                if (
                    (left + offsetWidth > windowInnerWidth) && (targetOriginX < offsetWidth / 2) ||
                    (left < 0) && (targetOriginX > offsetWidth / 2)
                ) {
                    left = anchorRect.right - anchorOriginX - offsetWidth + targetOriginX - elementRect.left
                    if (left < 0) {
                        left = 0
                    }
                }
                if (
                    (top + offsetHeight > windowInnerHeight) && (targetOriginY < offsetHeight / 2) ||
                    (top < 0) && (targetOriginY > offsetHeight / 2)
                ) {
                    top = anchorRect.bottom - anchorOriginY - offsetHeight + targetOriginY - elementRect.top
                }

                if (top < 0) {
                    top = 0
                }
            }
            this.popContainer.style.left = `${Math.round(left)}px`
            this.popContainer.style.top = `${Math.round(top)}px`
        }
    }
}