declare namespace UI {
    namespace SimpleDialog {
        interface Props extends React.Props<any> {
            /**
             * 确认对话框时执行
             */
            onConfirm: () => any;
            onClose?: () => any;
        }
    }
}