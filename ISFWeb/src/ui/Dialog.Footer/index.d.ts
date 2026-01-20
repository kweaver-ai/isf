declare namespace UI {
    namespace DialogFooter {
        interface Props extends React.Props<void> {
        }

        interface Element extends React.ReactElement<Props> {
        }

        interface StatelessComponent extends React.FunctionComponent<Props> {
            (props: Props): Element;
        }
    }
}