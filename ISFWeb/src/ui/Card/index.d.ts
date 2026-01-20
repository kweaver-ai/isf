declare namespace UI {
    namespace Card {
        interface Props extends React.Props<any> {
            /**
             * className
             */
            className?: any;

            /**
             * 宽度
             */
            width?: number | string;

            /**
             * 高度
             */
            height?: number | string;
        }

        interface State {
            /**
             * 宽度
             */
            width: number | string;

            /**
             * 高度
             */
            height: number | string;
        }
    }
}