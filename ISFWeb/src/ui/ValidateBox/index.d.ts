declare namespace UI {
    namespace ValidateBox {
        interface Props extends UI.TextBox.Props {
            /**
             * 角色
             */
            role?: string;

            // 验证状态
            validateState: any,

            // 验证信息
            validateMessages: {
                [key: string]: string,
            },

            // 气泡位置，'top' or 'right'
            align?: string,

            /**
             * 宽度
             */
            width?: string | number;

            /**
             * 是否只读
             */
            readOnly?: boolean;

        }
        interface State {
            /**
             * 当前对话框是否聚焦
             */
            focus: boolean;

            /**
             * 鼠标时候悬浮在当前对话框
             */
            hover: boolean;
        }
    }
}