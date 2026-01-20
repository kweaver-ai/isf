declare namespace UI {
    namespace TextBox {
        interface Props extends UI.TextInput.Props {
            /**
             * 宽度
             */
            width?: string | number;
        }

        interface State {
            /**
             * 当前对话框是否聚焦
             */
            focus: boolean;
        }

        interface Element extends React.ReactElement<Props> {
        }
    }
}