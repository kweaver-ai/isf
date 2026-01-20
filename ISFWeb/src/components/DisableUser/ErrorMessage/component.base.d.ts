declare namespace Console {
    namespace DisableUser {
        namespace ErrorMessage {
            interface Props extends React.Props<void> {
                // 提示信息
                errorType: number;

                // 确定事件
                onConfirm: () => void;
            }
        }
    }

}