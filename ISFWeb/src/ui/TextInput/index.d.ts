declare namespace UI {
    namespace TextInput {
        interface Props extends React.Props<any> {
            /**
             * 文本框类型
             */
            type?: 'text' | 'password' | 'search';

            /**
             * HTML id
             */
            id?: string;

            /**
             * class
             */
            className?: string;

            /**
             * 内联样式
             */
            style?: Object;

            /**
             * 初始值
             */
            value?: any;

            /**
             * 是否只读
             */
            readOnly?: boolean;

            /**
             * 自动聚焦
             */
            autoFocus?: boolean;

            /** 聚焦时选中 */
            selectOnFocus?: [number] | [number, number] | boolean

            /**
             * 是否禁用
             */
            disabled?: boolean;

            /**
             * 占位符
             */
            placeholder?: string;

            /**
             * 文本框是否必填
             */
            required?: boolean;

            /**
             * 文本框输入字符最大长度(包括空格在内的长度)
             */
            maxLength?: number;

            /**
             * 输入限制函数
             */
            validator?: (val: string) => boolean;

            /**
             * 输入变化事件
             */
            onChange?: (val: string) => any;

            /**
             * 聚焦事件
             */
            onFocus?: (event: FocusEvent) => any;

            /**
             * 失去焦点事件
             */
            onBlur?: (event: FocusEvent) => any;

            /**
             * 点击事件
             */
            onClick?: (event: MouseEvent) => any;

            /**
             * 回车事件
             */
            onEnter?: (event: KeyboardEvent) => any;

            /**
             * 键盘按下事件
             */
            onKeyDown?: (event: KeyboardEvent) => any;

            /**
             * 鼠标悬浮
             */
            onMouseover?: (event: KeyboardEvent) => any;

            /**
             * 鼠标移出悬浮
             */
            onMouseout?: (event: KeyboardEvent) => any;

        }

        interface State {
            /**
             * DOM上实际的value值
             */
            value: string;
        }
    }
}