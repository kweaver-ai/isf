declare namespace UI {
    namespace FormLabel {
        interface Props extends React.Props<void> {
            /**
             * 对齐方式
             */
            align?: string;

            /**
             * classname
             */
            className?: string;

            /**
             * 是否显示冒号
             */
            colon?: boolean;
        }

        interface Element extends React.ReactElement<Props> {
        }

        interface Component extends React.FunctionComponent<Props> {
            (props: Props): Element;
        }
    }
}