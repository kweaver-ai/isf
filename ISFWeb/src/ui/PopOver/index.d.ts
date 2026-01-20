declare namespace UI {
    namespace PopOver {

        interface Props extends React.Props<any> {

            /**
             * 锚点元素
             */
            anchor?: HTMLElement | null,
            /**
             * 锚点原点坐标
             */
            anchorOrigin?: [number | string, number | string],
            /**
             * 弹出层原点
             */
            targetOrigin?: [number | string, number | string],
            /**
             * 打开状态
             */
            open?: boolean,

            /**
             * 样式
             */
            style?: any;

            /**
             * 靠近浏览器边缘时向反方向展开
             */
            autoFix?: boolean

            /**
             * 监听页面位置变化
             */
            watch?: boolean

            /**
             * 弹出层class
             */
            className?: string

            /**
             * 用覆盖层冻结页面
             */
            freezable?: boolean

            trigger?: React.ReactElement<any>

            triggerEvent?: 'click' | 'mouseover'

            keepOpenWhenMouseOver?: boolean

            closeWhenMouseLeave?: boolean

            /**
             * 关闭弹出内容时隐藏弹出层
             */
            hideLayerWhenClose?: boolean

            /**
             * 点击弹出框触发
             * @param close 关闭函数
             */
            onRequestCloseWhenClick?: (close: () => void) => any

            /**
             * 点击弹出框外触发
             * @param close 关闭函数
             */
            onRequestCloseWhenBlur?: (close: () => void) => any

            /**
             * 点击覆盖层
             */
            onClickAway?: (e: React.SyntheticEvent<any>) => any

            /**
             * 鼠标移入
             */
            onMouseEnter?: (e: React.SyntheticEvent<any>) => any

            /**
             * 鼠标移出
             */
            onMouseLeave?: (e: React.SyntheticEvent<any>) => any

            /**
             * 单击
             */
            onClick?: (e: React.SyntheticEvent<any>) => any
        }
    }
}