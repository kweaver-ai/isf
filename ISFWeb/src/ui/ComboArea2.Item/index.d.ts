declare namespace UI {
    namespace ComboAreaItem {

        interface Props extends React.Props<any> {
            /**
             * 只读
             */
            readOnly?: boolean;

            /**
             * 禁用
             */
            disabled?: boolean;

            /**
             * 数据
             */
            data: any;

            /**
             * 样式
             */
            className: string;

            /**
             * chip的样式
             */
            chipClassName: string;

            /**
             * 移除chip时触发
             */
            removeChip: (data: any) => void;
        }

        interface State {

        }
    }
}