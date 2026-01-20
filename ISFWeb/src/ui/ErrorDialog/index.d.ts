declare namespace UI {
    namespace ErrorDialog {
        interface Props extends React.Props<any> {
            /**
             * 角色
             */
            role?: string;
            /**
             * 确定错误时触发
             */
            onConfirm?: () => any;
        }
    }
}