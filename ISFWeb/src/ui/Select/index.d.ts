declare namespace UI {
    namespace Select {

        /**
         * 子菜单属性
         */
        type Menu = {
            /**
             * 宽度
             */
            width?: string | number;

            /**
             * 最大宽度
             */
            maxHeight?: string | number;
        }

        interface Props extends React.Props<any> {
            /*
             * class
             */
            className?: string;
            /*
             * 初始值
             */
            value?: any;

            /*
             * 是否禁用
             */
            disabled?: boolean

            /*
             * 选中Option执行
             */
            onChange?: Function;

            /*
             * 下拉菜单
             */
            menu?: Menu;

            /**
             * children必须是SelectOption元素
             */
            children?: Array<UI.SelectOption.Element>

            /**
             * 长度
             */
            width?: number
        }

        interface State {
            /**
             * 初始值或选中项的值
             */
            value?: any;

            /**
             * 选中项的文本
             */
            text?: string;

            /**
             * 是否激活下拉菜单
             */
            active?: boolean;
        }

        interface Component extends React.Component<Props, State> {

        }
    }
}