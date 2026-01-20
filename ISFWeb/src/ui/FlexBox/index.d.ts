declare namespace UI {
    namespace FlexBox {

        interface Props extends React.Props<any> {
        }

        interface Element extends React.ReactElement<Props> {
        }

        interface Component extends React.FunctionComponent<Props> {
            (props: Props): Element;

            Item: UI.FlexBoxItem.Component;
        }
    }
}