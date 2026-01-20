declare namespace UI {
    namespace FlexBoxItem {
        interface Props extends React.Props<any> {
            /**
             * 角色
             */
            role?: string;

            /**
             * 宽度，存在多个FlexBoxItem时会根据比例自适应
             */
            width?: string | number;

            /**
             * 对齐方式，使用空格隔开，如: "left center"
             */
            align?: string;
        }

        interface Component extends React.FunctionComponent<Props> {
        }

        interface Element extends React.ReactElement<Props> {
        }
    }
}