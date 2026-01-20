declare namespace Components {
    namespace SetUsersFreezeStatus {
        namespace ShowErrorMessage {
            interface Props extends React.Props<void> {

                // 错误状态
                errorType: number;

                // 选中的用户
                userInfo?: any;

                // 冻结或者解冻状态
                freezeStatus: boolean;

                // 错误详细信息
                errorInfo?: any;

                // 确定事件
                onConfirm: (value) => void;
            }
        }
    }

}