declare namespace UI {
    namespace WizardStep {
        interface Props extends React.Props<any> {
            /**
             * 标题
             */
            title: string;

            /**
             * 当前页是否显示
             */
            active?: boolean;

            /**
             * 进入时触发
             */
            onEnter?: () => any;

            /**
             * 离开前触发，如果返回false则不允许离开
             */
            onBeforeLeave?: () => any;

            /**
             * 离开时触发
             */
            onLeave?: () => any;
        }
    }
}