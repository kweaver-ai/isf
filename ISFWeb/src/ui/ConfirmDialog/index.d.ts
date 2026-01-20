declare namespace UI {
    namespace ConfirmDialog {
        interface Props extends React.Props<void> {
            /**
             * 角色
             */
            role?: string;

            /**
             * 执行确认操作
             */
            onConfirm: () => any;

            /**
             * 执行取消操作
             */
            onCancel: () => any;

            /**
             * 提示文字
             */
            title?: string;

            /**
             * 确认按钮是否灰化
             */
            confirmBtnDisable?: boolean;
        }
    }
}