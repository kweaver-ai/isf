declare namespace UI {
    namespace Fold {
        interface LabelProps {
            /**
             * 标签的className
             */
            className: any;
        }

        interface Props extends React.Props<void> {
            /**
             * className
             */
            className: any;

            /**
             * 标题
             */
            label: string | React.ReactElement<any>;

            /**
             * 展开状态
             */
            open: boolean;

            /**
             * 展开事件
             */
            onToggle: () => any;

            /**
             * 标签样式
             */
            labelProps: LabelProps;
        }

        interface State {
            /**
             * 展开状态
             */
            open: boolean
        }
    }
}