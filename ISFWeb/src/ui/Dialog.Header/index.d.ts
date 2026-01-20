declare namespace UI {
    namespace DialogHeader {
        interface Props extends React.Props<any> {
            /**
             * 是否显示关闭按钮
             */
            closable?: boolean;

            /**
             * 点击关闭按钮时执行
             */
            onClose?: () => any;

            /**
             * 扩展右侧按钮
             */
            HeaderButtons?: Array<React.ReactElement<any>>;

        }
    }
}