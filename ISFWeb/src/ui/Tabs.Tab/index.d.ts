declare namespace UI {
    namespace TabsTab {
        interface Props extends React.Props<any> {
            /**
             * role
             */
            role?: string;

            /**
             * 激活状态
             */
            active?: boolean;

            /**
             * 样式
             */
            style?: Object;

            /**
             * 激活时触发
             */
            onActive?: () => any;

            /**
             * className
             */
            className?: string;

        }

        interface Element extends React.FunctionComponent<Props> {
        }
    }
}