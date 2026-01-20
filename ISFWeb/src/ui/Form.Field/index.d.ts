declare namespace UI {
    namespace FormField {
        interface Props extends React.Props<void> {

            /**
             * 样式
             */
            className?: string;

            /**
             * 是否显示必填标识
             */
            isRequired?: boolean;
        }

        interface Element extends React.ReactElement<Props> {
        }

        interface Component extends React.FunctionComponent<Props> {
            (props: Props): Element
        }
    }
}