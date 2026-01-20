declare namespace UI {
    namespace UIIcon {
        interface Props extends UI.FontIcon.Props {
            /**
             * 角色
             */
            role?: string;
        }

        interface Element extends React.ReactElement<Props> {

        }
    }
}