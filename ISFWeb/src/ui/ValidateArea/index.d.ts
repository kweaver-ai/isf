declare namespace UI {
    namespace ValidateArea {
        interface Props extends UI.TextBox.Props {
            // 验证状态
            validateState: any,

            // 验证信息
            validateMessages: {
                [key: string]: string,
            },

            // 气泡位置，'top' or 'right'
            align?: 'top' | 'right',

            /**
             * 宽度
             */
            width?: string | number;

            /**
             * 宽度
             */
            height?: string | number;

            /**
             * 占位提示
             */
            placeholder?: string;
        }
        interface State {
            /**
             * 当前对话框是否聚焦
             */
            focus: boolean;

            /**
             * 鼠标是否悬浮当前对话框
             */
            hover: boolean;
        }
    }
}