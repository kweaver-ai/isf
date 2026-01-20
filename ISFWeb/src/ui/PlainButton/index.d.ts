declare namespace UI {
    namespace PlainButton {
        interface Props {
            /**
             * 类型
             */
            type: string;

            /**
             * 禁用状态
             */
            disabled: boolean;

            /**
             * 图标
             */
            fallback?: string;

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
             * 点击事件
             */
            onClick: Function;

            /**
             * 图标的大小
             */
            size: number;
        }
    }
}