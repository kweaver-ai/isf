declare namespace UI {
    namespace FontIcon {
        interface Props extends React.Props<any> {
            /**
             * 角色
             */
            role?: string;

            /**
             * className
             */
            className?: string;

            /**
             * 是否禁用状态
             */
            disabled?: boolean;

            /**
             * 字体图标Unicode码
             */
            code: string;

            /**
             * 不支持字体图标时使用的png图片，base64编码
             */
            fallback?: string;

            /**
             * 字体家族名称
             */
            font?: string;

            /**
             * 鼠标悬浮提示
             */
            title?: any;

            /**
             * 图标尺寸，长宽相同
             */
            size?: number | string;

            /**
             * 图标颜色
             */
            color?: string;

            /**
             * title的className
             */
            titleClassName?: string;
            /**
           * 鼠标悬停样式
           */
            cursor?: string;
            /**
             * 点击时触发
             */
            onClick?: (event: React.MouseEvent<HTMLSpanElement>) => any;

            /**
             * 鼠标悬浮触发
             */
            onMouseOver?: (event: React.MouseEvent<HTMLSpanElement>) => any;

            /**
             * 鼠标移走触发
             */
            onMouseLeave?: (event: React.MouseEvent<HTMLSpanElement>) => any;

            /**
             * 鼠标移入事件
             */
            onMouseEnter?: (event: React.MouseEvent<HTMLSpanElement>) => any;
        }
    }
}