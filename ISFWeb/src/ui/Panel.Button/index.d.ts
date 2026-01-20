declare namespace UI {
    namespace PanelButton {
        interface Props extends UI.Button.Props {
            /**
             * 角色
             */
            role?: string;
        }

        interface Component extends React.FunctionComponent<Props> {
        }
    }
}