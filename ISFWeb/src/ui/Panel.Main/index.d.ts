declare namespace UI {
    namespace PanelMain {
        interface Props extends React.Props<void> {
            /**
             * 角色
             */
            role?: string;
        }

        interface Component extends React.FunctionComponent<Props> {
        }
    }
}