declare namespace UI {
    namespace SelectOption {
        interface Props extends React.Props<any> {
            /**
             * 选项值
             */
            value: any;

            /*
             * 当前是否选中
             */
            selected?: boolean;

            /**
             * 当前是否被禁用
             */
            disabled?: boolean;

            /**
             * 选中时触发
             */
            onSelect?: (selected: any) => any;
        }

        interface Element extends React.ReactElement<Props> {
        }

        interface Component extends React.FunctionComponent<Props> {
            (props: Props): Element;
        }
    }
}