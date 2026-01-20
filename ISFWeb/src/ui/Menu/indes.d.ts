declare namespace UI {
    namespace Menu {
        interface Props extends React.Props<any> {
            /**
             * 宽度
             */
            width?: string | number;

            /**
             * 最大高度
             */
            maxHeight?: string | number;

        }

        interface Component extends React.FunctionComponent<Props> {
            Item: UI.MenuItem.Component;
        }
    }
}