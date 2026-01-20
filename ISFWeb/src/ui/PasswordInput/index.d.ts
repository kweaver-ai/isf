declare namespace UI {
    namespace PasswordInput {
        interface Props extends React.Props<void> {

            /**
             * 文本值
             */
            value?: string;

             /**
             * class
             */
            className?: string;

            /**
             * 是否只读
             */
            readOnly?: boolean;

            /**
             * 是否禁用
             */
            disabled?: boolean;

            /**
             * 是否允许空值
             */
            required?: boolean;

            /**
             * 占位提示
             */
            placeholder?: string;

            /**
             * 输入验证
             * @param value 文本值
             */
            validator?: (input: string) => boolean;

            /**
             * 文本值发生变化时触发
             * @param value 文本值
             */
            onChange?: (value: string) => any;

            /**
             * 聚焦时触发
             */
            onFocus?: () => any;

            /**
             * 失焦时触发
             */
            onBlur?: () => any;

            /**
             * 悬浮效果
             */
            onMouseover?: () => any;

            /**
             * 移除悬浮效果
             */
            onMouseout?: () => any;

            /**
             * 点击键盘按键触发
             */
            onKeyDown?: () => any;
        }

        interface State {
            /**
             * 文本值
             */
            value?: string;

            /**
             * 是否聚焦
             */
            focus?: boolean;

        }
    }
}