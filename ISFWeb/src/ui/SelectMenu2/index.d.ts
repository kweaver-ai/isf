declare namespace UI {
    namespace SelectMenu {
        interface Props extends React.Props<any> {

            /**
             * 选择项
             * { name:string, ...otherProps }
             */
            selectValue: { name: string };

            /**
             * 标题文字
             */
            label: string;

            /**
             * 样式
             */
            className?: string;

            /**
             * 下拉菜单候选项
             * Array<{ name:string, ...otherProps }>
             */
            candidateItems: Array<{ name: string }>;

            /**
             * 选择下拉菜单项时触发
             * {
             *      name: string;
             *      ...otherProps
             * }
             */
            onSelect: (item: Object) => void;

            /**
             * 弹出层关闭之前执行的操作
             */
            onBeforePopupClose?: (event) => void;
        }

        interface State {
            /**
             * 选择项
             * {
             *      name: string;
             *      ...otherProps
             * }
             */
            selectValue: { name: string };

        }
    }
}