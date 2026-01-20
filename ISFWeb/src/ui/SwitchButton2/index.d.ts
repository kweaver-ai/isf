declare namespace UI {
    namespace SwitchButton2 {
        interface Props extends React.Props<any> {
            /**
             * 当前传入的值
             */
            value: any;

            /**
             * 禁用状态
             */
            disabled?: boolean;

            /**
             * 按钮当前状态
             */
            active: boolean;

            /**
             * 切换事件
             */
            onChange: () => any;
        }
    }
}