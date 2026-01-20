declare namespace UI {
    namespace ToolBar {
        interface Props extends React.Props<void> {
        }

        interface Component extends React.FunctionComponent<Props> {
            Button: UI.Button.Component
        }
    }
}