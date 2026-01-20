declare namespace UI {
    namespace NumberBox {
        interface Props extends React.Props<any> {
            /**
             * 是否被禁用
             */
            disabled?: boolean;

            /**
             * 起始值
             */
            value: number | string;

            /**
             * class
             */
            className?: string;

            /**
             * 内联样式
             */
            style?: {
                [key: string]: string | number,
            };

            /**
             * 宽度
             */
            width?: number | string;

            /**
             * 输入不规范提示
             * [max, min] 超过最大值的提示，或者最小值的提示
             */
            ValidatorMessage?: {
                max?: string,
                min?: string,
            }

            /**
             * 步进
             */
            step?: number;

            /**
             * 上限
             */
            max?: number;

            /**
             * 下限
             */
            min?: number;

            /**
             * 是否自动聚焦
             */
            autoFocus?: boolean;

            /**
             * 输入限制函数
             */
            validator?: (value) => boolean;

            /**
             * 文本框变化时触发
             */
            onChange: (value: string | number) => any;

            /**
             * 聚焦时时触发
             * @param event 文本框对象
             */
            onFocus?: (event: FocusEvent) => any;

            /**
             * 失去焦点时触发
             * @param event 文本框对象
             */
            onBlur?: (event: FocusEvent) => any;

            /**
             * 键盘输入时触发
             */
            onKeyDown?: (event: KeyboardEvent) => any;

            /**
             * 点击时触发
             */
            onClick?: (event: MouseEvent) => any

        }

        interface State {
            /**
             * 输入值
             */
            value: string;

            /**
             * 文本框是否聚焦
             */
            focus?: boolean;

            /**
             * 错误提示原因
             */
            validateState?: number;

            /**
             * 验证数据
             */
            validateMessages: {
                [key: number]: string,
            }
        }
    }
}