declare namespace UI {
    namespace TriggerPopMenu {

        interface Props extends React.Props<any> {

            /**
             * 弹出层样式
             */
            popMenuClassName?: string;

            /**
             * 自定义title
             */
            title?: string;

            /**
             * 按钮名称
             */
            label: string;

            /**
             * 按钮样式
             */
            className?: string

            /**
             * 点击弹出框外触发
             * @param close 关闭函数
             */
            onRequestCloseWhenBlur: (close: () => void) => any;

            /**
             * 点击弹出框触发
             * @param close 关闭函数
             */
            onRequestCloseWhenClick?: (close: () => void) => any;

            /**
             * 按钮高亮延迟消失的时间
             */
            timeout?: number;

            /**
             * 按钮文字的限制长度
             */
            numberOfChars?: number;
        }

        interface State {

            /**
             * 点击状态
             */
            clickStatus: boolean;
        }
    }
}