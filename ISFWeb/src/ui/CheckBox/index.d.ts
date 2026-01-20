declare namespace UI {
    namespace CheckBox {
        interface Props extends React.Props<any> {
            /**
             * 复选框值
             */
            value?: any;

            /**
             * HTML id
             */
            id?: string;

            /**
             * 是否禁用
             */
            disabled?: boolean;

            /**
             * 初始选中状态
             */
            checked?: boolean;

            /**
             * 变化时触发
             * @param checked 是否勾选
             * @param value 复选框值
             */
            onChange?: (checked: boolean, value: any) => any;

            /**
             * 点击时触发
             * @param value 复选框值
             */
            onCheck?: (value: any) => any;

            /**
             * 取消勾选时触发
             * @param value 复选框值
             */
            onUncheck?: (value: any) => any;

            /**
             * 勾选时触发
             * @param event 鼠标事件
             */
            onClick?: (event: MouseEvent) => any;

            /**
             * 是否是半选状态
             */
            halfChecked?: boolean;

            /**
             * 自定义样式
             */
            className?: string;
        }

        interface State {
            /**
             * 当前选中状态
             */
            checked: boolean;
        }
    }
}