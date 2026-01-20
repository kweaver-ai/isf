declare namespace UI {
    namespace Control {
        interface Props extends React.Props<void> {
            /**
             * className
             */
            className?: string;

            /**
             * role
             */
            role?: string;

            /**
             * style
             */
            style?: {
                [key: string]: string | number | void;
            };

            /**
             * width，包含盒模型的padding和border
             */
            width?: number | string;

            /**
             * height，包含盒模型的padding和border
             */
            height?: number | string;

            /**
             * minHeight，包含盒模型的padding和border
             */
            minHeight?: number | string;

            /**
             * maxHeight，包含盒模型的padding和border
             */
            maxHeight?: number | string;

            /**
             * 是否默认聚焦
             */
            focus?: boolean;

            /**
             * 是否禁用
             */
            disabled?: boolean;

            /**
             * 点击触发
             */
            onClick?: (event: MouseEvent) => any;

            /**
             * 失焦
             */
            onBlur?: (event: MouseEvent) => any;

            /**
             * 鼠标悬浮
             */
            onMouseOver?: (event: MouseEvent) => any;

            /**
             * 鼠标移除
             */
            onMouseLeave?: (event: MouseEvent) => any;
        }

        interface Element extends React.ReactElement<Props> {
        }
    }
}