declare namespace UI {
    namespace Panel {
        interface Props extends React.Props<void> {
            /**
             * 角色
             */
            role?: string;
        }

        interface State {
        }

        interface Component extends React.FunctionComponent<Props> {
            Main: UI.PanelMain.Component;
            Footer: UI.PanelFooter.Component;
            Button: UI.PanelButton.Component;
        }
    }
}