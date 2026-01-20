declare namespace UI {
    namespace Wizard {
        interface Props extends React.Props<any> {
            title?: string; // 向导对话框标题栏

            onFinish?: () => any; // 点击完成时触发

            onCancel?: () => any; // 点击取消时触发
        }
    }
}