declare namespace UI {
    namespace MenuItem {
        interface Props extends React.Props<void> {
            /**
             * 点击时触发
             */
            onClick: () => any;
        }

        interface Component extends React.FunctionComponent<Props> {
        }
    }
}