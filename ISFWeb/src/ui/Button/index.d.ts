declare namespace UI {
    namespace Button {
        interface Props extends React.Props<any> {
            /**
             * 角色
             */
            role?: string;

            /**
             * 按钮类型
             */
            type?: 'button' | 'submit' | 'reset';

            /**
             * 按钮主题，默认为'regular'
             */
            theme?: 'oem' | 'regular' | 'dark';

            /**
             * 是否禁用
             */
            disabled?: boolean;

            /**
             * FontIcon code
             */
            icon?: string;

            /**
             * 最小宽度
             */
            minWidth?: number | string;

            /**
             * 宽度
             */
            width?: number | string;

            /**
             * 样式class
             */
            className?: string;

            /**
             * 点击按钮触发
             */
            onClick?: (event: React.MouseEvent<HTMLButtonElement>) => any;

            /**
             * 按下按钮触发
             */
            onMouseDown?: (event: React.MouseEvent<HTMLButtonElement>) => any;

            /**
             * 图标
             */
            fallback?: string;

            /**
             * 图标的大小
             */
            size: number;
        }

        interface Component extends React.FunctionComponent<UI.Button.Props> {
        }
    }
}