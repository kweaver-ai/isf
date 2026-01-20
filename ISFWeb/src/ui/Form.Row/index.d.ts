declare namespace UI {
    namespace FormRow {
        interface Props extends React.Props<void> {
            /**
             * classname
             */
            className?: string;
        }

        interface Element extends React.ReactElement<Props> {
        }

        interface Component extends React.FunctionComponent<Props> {
            (props: Props): Element;
        }
    }
}