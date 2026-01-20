declare namespace UI {
    namespace Dialog {
        interface Refs {
            container: React.ReactHTMLElement<HTMLDivElement>
        }

        interface Props extends React.Props<Refs> {
            /**
             * 宽度
             */
            width?: number | string;

            /**
             * 高度
             */
            height?: number | string;

            /**
             * 隐藏对话框
             */
            hide?: boolean;

            /**
             * 对话框是否可拖动
             */
            draggable?: boolean;

            /**
             * 对话框缩放时触发
             */
            onResize?: (size: { width: number | string, height: number | string }) => any;
        }

        interface State {

            /**
             * 相对视窗顶定位距离
             */
            top?: number | string;

            /**
             * 相对视窗左定位距离
             */
            left?: number | string;
        }

        interface Element extends React.ComponentElement<Props, any> {
        }
    }
}